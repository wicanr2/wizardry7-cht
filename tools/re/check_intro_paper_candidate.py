#!/usr/bin/env python3
"""R974：獨立核原生字模、作者保護區與候選區外像素，不裁決補底品質。"""
import argparse
import gzip
import hashlib
import json
from pathlib import Path

from check_native_gui import png


def contains(x, y, polygon):
    signs = [(b[0]-a[0])*(y-a[1])-(b[1]-a[1])*(x-a[0])
             for a, b in zip(polygon, polygon[1:]+polygon[:1])]
    return all(v >= 0 for v in signs) or all(v <= 0 for v in signs)


def font_data(unifont, eten, chars):
    raw = unifont.read_bytes()
    assert hashlib.sha256(raw).hexdigest() == '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
    tables = {'unifont': {}, 'eten': {}}
    for line in gzip.decompress(raw).decode('ascii').splitlines():
        cp, bitmap = line.split(':')
        ch = chr(int(cp, 16))
        if ch in chars:
            assert len(bitmap) == 64
            tables['unifont'][ch] = [int(bitmap[y*4:y*4+4],16) for y in range(16)]
    raw = eten.read_bytes()
    assert hashlib.sha256(raw).hexdigest() == '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
    assert len(raw) == 392820
    for ch in chars:
        a, b = ch.encode('big5')
        ordinal = (a-0xa4)*157 + (b-0x40 if b<0x7f else b-0xa1+63)
        if a >= 0xc9:
            ordinal -= 408
        assert 0 <= ordinal*30 <= len(raw)-30
        tables['eten'][ch] = [int.from_bytes(raw[ordinal*30+y*2:ordinal*30+y*2+2], 'big')//2 for y in range(15)]
    assert all(set(t) == set(chars) for t in tables.values())
    return tables


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('source', type=Path)
    p.add_argument('unifont_candidate', type=Path)
    p.add_argument('eten_candidate', type=Path)
    p.add_argument('--model', type=Path, default=Path(__file__).with_name('intro_paper_candidate.json'))
    p.add_argument('--unifont-hex', type=Path, required=True)
    p.add_argument('--eten-font', type=Path, required=True)
    p.add_argument('--output', type=Path, required=True)
    args = p.parse_args()
    assert not args.output.exists()
    model_raw = args.model.read_bytes()
    model = json.loads(model_raw)
    chars = set(c for f in model['frames'] for r in f['labels'] for c in r['zh'])
    fonts = font_data(args.unifont_hex, args.eten_font, chars)
    report = []
    for backend, root in (('unifont',args.unifont_candidate),('eten',args.eten_candidate)):
        receipt = json.loads((root/'receipt.json').read_text())
        size = 16 if backend == 'unifont' else 15
        assert receipt['model_sha256'] == hashlib.sha256(model_raw).hexdigest()
        assert receipt['backend'] == backend and receipt['size'] == size
        assert receipt['glyphs'] == fonts[backend]
        assert len(receipt['frames']) == 12
        assert {(r['step'],r['scale']) for r in receipt['frames']} == {(f['step'],k) for f in model['frames'] for k in (2,3)}
        for row in receipt['frames']:
            f = next(f for f in model['frames'] if f['step'] == row['step'])
            k = row['scale']
            stem = args.source/f'picture-{row["step"]}'
            indexed, pal = stem.with_suffix('.a000').read_bytes(), stem.with_suffix('.pal').read_bytes()
            assert hashlib.sha256(indexed).hexdigest() == f['indexed_sha256']
            assert hashlib.sha256(pal).hexdigest() == model['palette_sha256']
            w, h, original = png(root/row['original'])
            cw, ch, current = png(root/row['candidate'])
            assert (w,h,cw,ch) == (320*k,200*k,320*k,200*k)
            assert hashlib.sha256((root/row['candidate']).read_bytes()).hexdigest() == row['candidate_sha256']
            assert len(row['labels']) == len(f['labels'])
            glyph_pixels, cells = {}, set()
            for label, config in zip(row['labels'], f['labels']):
                assert all(label[key] == config[key] for key in config)
                width = len(config['zh'])*(size+2)-2
                x0 = round(config['center'][0]*k)-width//2
                y0 = round(config['center'][1]*k)-size//2
                assert label['xy'] == [x0,y0] and label['width'] == width and label['height'] == size
                color = pal[:3] if config['color']=='black' else bytes((156,20,12))
                assert bytes(label['rgb']) == color
                new_cells = {(x0+x,y0+y) for y in range(size) for x in range(width)}
                assert not cells & new_cells, '字格重疊'
                cells.update(new_cells)
                for x,y in new_cells:
                    assert contains((x+.5)/k,(y+.5)/k,config['safe']), '字格出候選紙面'
                    assert not any(contains(x//k+.5,y//k+.5,q) for q in f['protect'])
                for j,c in enumerate(config['zh']):
                    for gy,bits in enumerate(fonts[backend][c]):
                        for gx in range(size):
                            if bits & (1<<(size-gx-1)):
                                glyph_pixels[(x0+j*(size+2)+gx,y0+gy)] = color
            bands = {(x,y) for y in range(200) for x in range(320)
                     if any(contains(x+.5,y+.5,band) and contains(x+.5,y+.5,surface)
                            for band,surface in zip(f['bands'],f['band_surfaces']))}
            protected = {(x,y) for y in range(200) for x in range(320)
                         if any(contains(x+.5,y+.5,q) for q in f['protect'])}
            assert not bands & protected
            changed = 0
            for y in range(h):
                for x in range(w):
                    off = (y*w+x)*3
                    reference = pal[indexed[(y//k)*320+x//k]*3:indexed[(y//k)*320+x//k]*3+3]
                    assert original[off:off+3] == reference, '原版PNG與色盤不符'
                    actual = current[off:off+3]
                    if (x,y) in glyph_pixels:
                        assert actual == glyph_pixels[(x,y)], '原生字模前景不符'
                    elif (x//k,y//k) not in bands:
                        assert actual == reference, '批准補底及字模筆畫外變更'
                    if (x//k,y//k) in protected:
                        assert actual == reference, '作者保護區变更'
                    changed += actual != reference
            report.append(dict(backend=backend,scale=k,step=row['step'],pixels=w*h,
                               native_ink=len(glyph_pixels),changed=changed,outside=0,protected_changed=0))
    value = dict(schema=1,frames=24,rows=report,native_chars=len(chars),
                 limits='獨立核原字模及候選區外；字帶、紙面與作者區仍為候選座標。未驗補底區內零位元背景的期望、英文殘筆、外觀、正式GUI或原版同狀態中文覆繪。')
    with args.output.open('x') as out:
        json.dump(value,out,ensure_ascii=False,indent=2)
        out.write('\n')
    print(f'24張候選、{sum(r["pixels"] for r in report)}像素：字模前景相符、批准區外及作者保護區差異0')


if __name__ == '__main__':
    main()
