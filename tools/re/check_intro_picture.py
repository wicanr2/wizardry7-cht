#!/usr/bin/env python3
"""025：獨立原生字模、已選原型、GUI／HD及批准區域驗證。只在Docker執行。"""
import argparse
import gzip
import hashlib
import json
from pathlib import Path

from check_native_gui import png


def inside(x, y, polygon):
    # 射線及線段判定，獨立於前端的凸多邊形判定。
    hit = False
    for a, b in zip(polygon, polygon[1:] + polygon[:1]):
        px, py = 2*x+1, 2*y+1
        ax, ay, bx, by = 2*a[0], 2*a[1], 2*b[0], 2*b[1]
        if (bx-ax)*(py-ay) == (by-ay)*(px-ax) and min(ax, bx) <= px <= max(ax, bx) and min(ay, by) <= py <= max(ay, by):
            return True
        if (a[1] > y+.5) != (b[1] > y+.5):
            if x+.5 < (b[0]-a[0])*(y+.5-a[1])/(b[1]-a[1])+a[0]:
                hit = not hit
    return hit


def font_rows(unifont, eten):
    raw = unifont.read_bytes()
    assert hashlib.sha256(raw).hexdigest() == '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
    result = {'unifont': [], 'eten': []}
    chars = {}
    for line in gzip.decompress(raw).decode('ascii').splitlines():
        cp, bits = line.split(':')
        if int(cp, 16) in map(ord, '呈獻'):
            assert len(bits) == 64
            chars[chr(int(cp, 16))] = [int(bits[i:i+4], 16) for i in range(0, 64, 4)]
    result['unifont'] = [chars[c] for c in '呈獻']
    data = eten.read_bytes()
    assert len(data) == 392820
    for c in '呈獻':
        hi, lo = c.encode('big5')
        col = lo-64 if lo <= 126 else lo-161+63
        index = (hi-164)*157+col-(408 if hi >= 201 else 0)
        assert 0 <= index*30 < len(data)-29
        result['eten'].append([int.from_bytes(data[index*30+y*2:index*30+y*2+2], 'big') >> 1 for y in range(15)])
    return result


def expected(row, original, palette, glyphs):
    polygon, k = row['polygon_candidate'], row['scale']
    width, height = 320*k, 200*k
    colors = [palette[i*3:i*3+3] for i in original]
    after = list(colors)
    xs, ys = [p[0] for p in polygon], [p[1] for p in polygon]
    samples = [[] for _ in range(4)]
    for y in range(min(ys)-4, max(ys)+5):
        for x in range(min(xs)-4, max(xs)+5):
            c = colors[y*320+x]
            if not inside(x, y, polygon) and c[0] > c[1] > c[2] and c[0]-c[2] >= 30:
                samples[x%2+2*(y%2)].append((x, y, c))
    assert list(map(len, samples)) == row['samples']
    for y in range(min(ys), max(ys)):
        for x in range(min(xs), max(xs)):
            if not inside(x, y, polygon):
                continue
            group = x%2+2*(y%2)
            sx, sy, c = min(samples[group], key=lambda p: ((p[0]-x)**2+(p[1]-y)**2, p[1], p[0]))
            # 使用定案前原型收據的固定係數，不取被測前端的矩陣結果。
            after[y*320+x] = bytes(max(0, min(255, round(c[ch]+model[1]*(x-sx)+model[2]*(y-sy)))) for ch, model in enumerate(row['models'][group]))
    result = bytearray()
    for y in range(200):
        line = b''.join(c*k for c in after[y*320:(y+1)*320])
        result.extend(line*k)
    size = len(glyphs[0])
    cx = round(sum(xs)/4*k)-(2*size+2)//2
    cy = round(sum(ys)/4*k)-size//2
    on = 0
    for j, g in enumerate(glyphs):
        for y, bits in enumerate(g):
            for x in range(size):
                if bits & (1 << (size-1-x)):
                    assert inside((cx+j*(size+2)+x)//k, (cy+y)//k, polygon), 'glyph outside approved band'
                    offset = ((cy+y)*width+cx+j*(size+2)+x)*3
                    result[offset:offset+3] = palette[:3]
                    on += 1
    return width, height, bytes(result), on


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('receipts', type=Path)
    p.add_argument('--prototype', type=Path, required=True)
    p.add_argument('--original', type=Path, required=True)
    p.add_argument('--unifont-hex', type=Path, required=True)
    p.add_argument('--eten-font', type=Path, required=True)
    p.add_argument('--output', type=Path, required=True)
    args = p.parse_args()
    assert not args.output.exists()
    rows = json.loads((args.prototype/'caption-prototype-v6.json').read_text())['frames']
    fonts = font_rows(args.unifont_hex, args.eten_font)
    captures = (39000000, 41000000, 43000000, 45000000, 47000000, 49000000, 60000000, 74000000)
    report = []
    for backend in ('unifont', 'eten'):
        for k in (2, 3):
            base = args.receipts/f'gui-english-{k}'
            current = args.receipts/f'gui-{backend}-{k}'
            receipt = json.loads((current/'receipt.json').read_text())
            assert receipt['exit_code'] == 0 and receipt['steps'] == list(captures)
            selected = [r for r in rows if r['scale'] == k]
            for index, step in enumerate(captures):
                name = f'hdat-{step}'
                a, b = [json.loads((d/(name+'.json')).read_text()) for d in (base, current)]
                assert all(a[key] == b[key] for key in ('Step', 'Mem', 'A000', 'Pal', 'Mouse')), 'original state mismatch'
                for ext in ('.a000', '.pal', '.comp'):
                    assert (base/(name+ext)).read_bytes() == (current/(name+ext)).read_bytes(), ext
                original = (current/(name+'.a000')).read_bytes()
                palette = (current/(name+'.pal')).read_bytes()
                w, h, english = png(base/(name+'.png'))
                if index < 7:
                    row = selected[min(index, 5)]
                    assert hashlib.sha256(original).hexdigest() == row['indexed_sha256']
                    assert hashlib.sha256(palette).hexdigest() == row['palette_sha256']
                    source = args.original/f'picture-{row["step"]}'
                    assert source.with_suffix('.a000').read_bytes() == original
                    assert source.with_suffix('.pal').read_bytes() == palette
                    ew, eh, want, on = expected(row, original, palette, fonts[backend])
                    assert (ew, eh) == (w, h)
                    if backend == 'unifont':
                        assert png(args.prototype/row['output'])[2] == want, 'approved prototype mismatch'
                else:
                    want, on = english, 0
                for prefix in ('hdat', 'drawat'):
                    cw, ch, actual = png(current/f'{prefix}-{step}.png')
                    assert (cw, ch) == (w, h) and actual == want, (backend, k, step, prefix, 'full pixel mismatch')
                changed = sum(english[i:i+3] != want[i:i+3] for i in range(0, len(want), 3))
                outside = 0
                if index < 7:
                    for i in range(w*h):
                        if english[i*3:i*3+3] != want[i*3:i*3+3] and not inside((i%w)//k, (i//w)//k, row['polygon_candidate']):
                            outside += 1
                assert outside == 0
                report.append(dict(backend=backend, scale=k, step=step, changed=changed, outside=outside, glyph_on=on, hd_draw_exact=True))
    value = dict(schema=1, rows=report, raw_state_equal=True, native_font_sources_equal=True,
                 limits='記憶體為既有1MiB FNV收據；indexed與palette另驗完整bytes及SHA。補底相同於定案原型，不代表原版無字底圖。')
    with args.output.open('x') as f:
        json.dump(value, f, ensure_ascii=False, indent=2)
        f.write('\n')
    print(f'{len(report)}個同狀態點，HD／Draw完整RGB相符，批准字帶外差異0')


if __name__ == '__main__':
    main()
