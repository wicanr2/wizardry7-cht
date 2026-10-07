#!/usr/bin/env python3
"""R974：重建紙頁中文候選；只在Docker執行，不接正式前端。"""
import argparse
import hashlib
import json
import struct
import zlib
from pathlib import Path

MODEL = Path(__file__).with_name('intro_paper_candidate.json')


def inside(x, y, poly):
    hit = False
    for a, b in zip(poly, poly[1:] + poly[:1]):
        if (a[1] > y) != (b[1] > y):
            if x < (b[0]-a[0])*(y-a[1])/(b[1]-a[1])+a[0]:
                hit = not hit
    return hit


def points(poly):
    return {(x, y) for y in range(max(0, min(p[1] for p in poly)), min(200, max(p[1] for p in poly)))
            for x in range(max(0, min(p[0] for p in poly)), min(320, max(p[0] for p in poly)))
            if inside(x+.5, y+.5, poly)}


def png_write(path, width, height, data):
    def chunk(kind, body):
        return struct.pack('>I', len(body))+kind+body+struct.pack('>I', zlib.crc32(kind+body)&0xffffffff)
    assert len(data) == width*height*3
    body = b''.join(b'\0'+data[y*width*3:(y+1)*width*3] for y in range(height))
    value = (b'\x89PNG\r\n\x1a\n'+chunk(b'IHDR', struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0))+
             chunk(b'IDAT', zlib.compress(body, 9))+chunk(b'IEND', b''))
    with path.open('xb') as f:
        f.write(value)


def scaled(rgb, scale):
    out = bytearray()
    for y in range(200):
        line = b''.join(rgb[(y*320+x)*3:(y*320+x+1)*3]*scale for x in range(320))
        out.extend(line*scale)
    return out


def patch(rgb, bands, protected, surfaces):
    if len(bands) != len(surfaces):
        raise ValueError('字帶與紙面候選數量不符')
    masks = [points(band) & points(surface) for band, surface in zip(bands, surfaces)]
    all_bands = set().union(*masks) if masks else set()
    if all_bands & protected:
        raise ValueError('補底字帶侵入作者保護區')
    out = bytearray(rgb)
    receipts = []
    for band, surface, cells in zip(bands, surfaces, masks):
        groups = [[] for _ in range(4)]
        for y in range(max(0, min(p[1] for p in surface)), min(200, max(p[1] for p in surface))):
            for x in range(max(0, min(p[0] for p in surface)), min(320, max(p[0] for p in surface))):
                c = rgb[(y*320+x)*3:(y*320+x+1)*3]
                if (inside(x+.5, y+.5, surface) and (x, y) not in all_bands | protected and c[0] > c[1] > c[2] and
                        c[1]*5 >= c[0]*3 and c[1]-c[2] >= 16):
                    groups[x%2+2*(y%2)].append((x, y, c))
        if min(map(len, groups)) < 8:
            raise ValueError('紙紋候選周邊樣本不足')
        for x, y in sorted(cells, key=lambda p: (p[1], p[0])):
            group = x%2+2*(y%2)
            nearest = sorted(groups[group], key=lambda p: ((p[0]-x)**2+(p[1]-y)**2, p[1], p[0]))[:8]
            weights = [1/((sx-x)**2+(sy-y)**2+1) for sx, sy, _ in nearest]
            c = bytes(round(sum(w*sample[2][ch] for w, sample in zip(weights, nearest))/sum(weights)) for ch in range(3))
            out[(y*320+x)*3:(y*320+x+1)*3] = c
        receipts.append(dict(polygon=band, surface=surface, samples=list(map(len, groups)),
                             method='same-parity nearest-eight inverse-distance RGB approximation'))
    return bytes(out), receipts


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('source', type=Path)
    p.add_argument('output', type=Path)
    p.add_argument('--model', type=Path, default=MODEL)
    p.add_argument('--eten-font', type=Path, required=True)
    args = p.parse_args()
    if args.output.exists():
        p.error('輸出已存在，請指定新目錄')
    model = json.loads(args.model.read_text())
    import bake_font
    bake_font.FONT15 = args.eten_font
    backend, size = bake_font.BACKEND, bake_font.HEIGHT
    chars = ''.join(dict.fromkeys(ch for r in model['frames'] for v in r['labels'] for ch in v['zh']))
    fonts = {ch: bake_font.glyph_char(ch) for ch in chars}
    if any(len(rows) != size or not any(rows) or any(not 0 <= bits < 1 << size for bits in rows) for rows in fonts.values()):
        raise ValueError('原生字模不完整')
    prepared = []
    for row in model['frames']:
        path = args.source/f'picture-{row["step"]}'
        indexed, pal = path.with_suffix('.a000').read_bytes(), path.with_suffix('.pal').read_bytes()
        if len(indexed) != 64000 or len(pal) != 768:
            raise ValueError('原版畫格長度不合格')
        if hashlib.sha256(indexed).hexdigest() != row['indexed_sha256'] or hashlib.sha256(pal).hexdigest() != model['palette_sha256']:
            raise ValueError('原版來源雜湊不符')
        rgb = b''.join(pal[i*3:i*3+3] for i in indexed)
        protected = set().union(*(points(q) for q in row['protect'])) if row['protect'] else set()
        patched, bands = patch(rgb, row['bands'], protected, row['band_surfaces'])
        for k in (2, 3):
            original, candidate = scaled(rgb, k), scaled(patched, k)
            labels = []
            for label in row['labels']:
                width = len(label['zh'])*(size+2)-2
                x0 = round(label['center'][0]*k)-width//2
                y0 = round(label['center'][1]*k)-size//2
                color = pal[:3] if label['color'] == 'black' else bytes((156, 20, 12))
                for gy in range(size):
                    for gx in range(width):
                        if not inside((x0+gx+.5)/k, (y0+gy+.5)/k, label['safe']):
                            raise ValueError(f'字格出紙頁：{row["step"]} {backend} {k} {label["zh"]}')
                        if ((x0+gx)//k, (y0+gy)//k) in protected:
                            raise ValueError('字格侵入作者保護區')
                for j, ch in enumerate(label['zh']):
                    for y, bits in enumerate(fonts[ch]):
                        for x in range(size):
                            if bits & (1 << (size-1-x)):
                                off = ((y0+y)*320*k+x0+j*(size+2)+x)*3
                                candidate[off:off+3] = color
                labels.append(dict(**label, xy=[x0,y0], width=width, height=size, rgb=list(color)))
            prepared.append((row, k, original, candidate, bands, labels))
    args.output.mkdir()
    records, sections = [], []
    for row, k, original, candidate, bands, labels in prepared:
        prefix = f'{backend}-{k}x-{row["step"]}'
        before, after = prefix+'-original.png', prefix+'-candidate.png'
        png_write(args.output/before, 320*k, 200*k, original)
        png_write(args.output/after, 320*k, 200*k, candidate)
        records.append(dict(step=row['step'], backend=backend, scale=k, indexed_sha256=row['indexed_sha256'],
                            original=before, candidate=after, bands=bands, protect=row['protect'], labels=labels,
                            candidate_sha256=hashlib.sha256((args.output/after).read_bytes()).hexdigest()))
        sections.append(f'<section id="{row["step"]}-{k}"><h2>{row["step"]}／{k}倍</h2><div><figure><figcaption>原版</figcaption><img src="{before}"></figure><figure><figcaption>新候選</figcaption><img src="{after}"></figure></div></section>')
    with (args.output/'receipt.json').open('x') as f:
        json.dump(dict(schema=1, status=model['status'], model_sha256=hashlib.sha256(args.model.read_bytes()).hexdigest(),
                       backend=backend, size=size, frames=records, glyphs=fonts), f, ensure_ascii=False, indent=2)
        f.write('\n')
    page = ('<!doctype html><html lang="zh-Hant"><meta charset="utf-8"><title>R974紙頁中文候選</title>'
            '<style>body{background:#181818;color:#eee;font:16px sans-serif;margin:24px}h2{font-size:18px}section>div{display:flex;flex-wrap:wrap;gap:16px}figure{margin:0}img{image-rendering:pixelated;max-width:100%;height:auto}section{margin:28px 0}</style>'
            '<h1>R974紙頁中文候選</h1><p>新重建版面，尚待外觀確認。紙紋補底是近似；作者英文保持。雙字型皆採原生字格，中文不縮小、不旋轉。</p>'
            '<p>部分遮蔽頁把三行移到可見紙面。第六張沒有覆繪。候選未接正式前端。</p>'+''.join(sections)+'</html>')
    (args.output/'compare.html').write_text(page)
    print(f'{backend}：12張候選、{len(chars)}字原生字模、未接正式前端')


if __name__ == '__main__':
    main()
