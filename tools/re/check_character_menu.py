"""023角色選單：由原字型及原版色號獨立驗證中文、界限與畫面差異。"""
import argparse
import gzip
import hashlib
import importlib.util
import json
import struct
from functools import lru_cache
from pathlib import Path

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('root', type=Path)
p.add_argument('--baseline-root', type=Path, required=True)
p.add_argument('--unifont-hex', type=Path, required=True)
p.add_argument('--eten-font', type=Path, required=True)
p.add_argument('--output', type=Path, required=True)
args = p.parse_args()
if args.output.exists():
    p.error('輸出已存在，請指定新報告路徑')
repo = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('pngcheck', repo / 'tools/re/check_native_gui.py')
decoder = importlib.util.module_from_spec(spec)
spec.loader.exec_module(decoder)
png = lru_cache(maxsize=64)(decoder.png)
labels = [('CREATE', 11, '創建', 4), ('REVIEW', 63, '檢視', 56),
          ('DELETE', 115, '刪除', 108), ('RENAME', 167, '改名', 160),
          ('PICTURE', 216, '頭像', 212), ('EXIT', 277, '離開', 264)]
needed = set(''.join(s[2] for s in labels))
fonts = {'eten': {}, 'unifont': {}}
raw = args.unifont_hex.read_bytes()
assert hashlib.sha256(raw).hexdigest() == '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
for line in gzip.decompress(raw).decode('ascii').splitlines():
    cp, bits = line.split(':')
    ch = chr(int(cp, 16))
    if ch in needed:
        assert len(bits) == 64
        fonts['unifont'][ch] = [int(bits[i:i+4], 16) for i in range(0, 64, 4)]
raw = args.eten_font.read_bytes()
assert hashlib.sha256(raw).hexdigest() == '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
for ch in needed:
    lead, trail = ch.encode('big5')
    col = trail-64 if trail < 127 else trail-98
    record = (lead-164)*157+col-(408 if lead >= 201 else 0)
    fonts['eten'][ch] = [word >> 1 for word in struct.unpack_from('>15H', raw, record*30)]
assert all(set(f) == needed for f in fonts.values())
report = {'same_state_pairs': 0, 'glyph_checks': 0, 'pixel_checks': 0,
          'buttons': [], 'before_after': []}
for scale in (2, 3):
    english = args.baseline_root / f'aux-english-{scale}'
    base = json.loads((english / 'captures.json').read_text())
    for backend in ('eten', 'unifont'):
        dest = args.root / f'aux-{backend}-{scale}'
        prior = args.baseline_root / f'aux-{backend}-{scale}'
        current = json.loads((dest / 'captures.json').read_text())
        old = json.loads((prior / 'captures.json').read_text())
        assert current['returncode'] == base['returncode'] == old['returncode'] == 0
        assert current['save_snapshots']['after_normal_exit'] == base['save_snapshots']['after_normal_exit'] == old['save_snapshots']['after_normal_exit']
        height = 15 if backend == 'eten' else 16
        gy = (378 if backend == 'eten' else 377) if scale == 2 else 570
        for step in (272000000, 400000000, 528000000):
            assert current['original'][str(step)] == base['original'][str(step)] == old['original'][str(step)]
            for ext in ('a000', 'pal'):
                assert (dest / f'hdat-{step}.{ext}').read_bytes() == (english / f'hdat-{step}.{ext}').read_bytes() == (prior / f'hdat-{step}.{ext}').read_bytes()
            assert (dest / f'hdat-{step}.comp').read_bytes() == (prior / f'hdat-{step}.comp').read_bytes()
            assert json.loads((dest / f'hdat-{step}.json').read_text())['Shown'] == json.loads((prior / f'hdat-{step}.json').read_text())['Shown']
            report['same_state_pairs'] += 1
            mask = set()
            if step == 400000000 and scale == 2:
                for _, x, zh, _ in labels:
                    width = len(zh)*(height+2)-2
                    mask.update((xx, yy) for start in (190*scale, gy)
                                for yy in range(start, start+height)
                                for xx in range(x*scale, x*scale+width))
            for prefix in ('hdat', 'drawat'):
                w, h, rgb = png(dest / f'{prefix}-{step}.png')
                ow, oh, before = png(prior / f'{prefix}-{step}.png')
                assert (w, h) == (ow, oh) == (320*scale, 200*scale)
                changed = outside = 0
                for yy in range(h):
                    for xx in range(w):
                        at = (yy*w+xx)*3
                        if rgb[at:at+3] != before[at:at+3]:
                            changed += 1
                            outside += (xx, yy) not in mask
                assert outside == 0
                assert (changed > 0) == (step == 400000000 and scale == 2)
                report['before_after'].append(dict(backend=backend, scale=scale, step=step,
                    prefix=prefix, changed_pixels=changed, outside_approved_regions=outside))
        a = (english / 'hdat-400000000.a000').read_bytes()
        pal = (english / 'hdat-400000000.pal').read_bytes()
        meta = json.loads((dest / 'hdat-400000000.json').read_text())
        for source, x, zh, bx in labels:
            event = next(e for e in meta['Shown'] if e['S'] == source and (e['X0'], e['Y0']) == (x, 190))
            assert event['Dy'] == 0 and event['G'] == [fonts[backend][ch] for ch in zh]
            assert list(a[188*320+bx:188*320+bx+49]) == [81]+[176]*47+[182]
            assert list(a[196*320+bx:196*320+bx+49]) == [182]+[197]*47+[198]
            assert all(a[yy*320+bx] == 176 and a[yy*320+bx+48] == 197 for yy in range(189, 196))
            assert event['Fg'] == 1
            ink = {(x*scale+i*(height+2)+xx, gy+yy) for i, ch in enumerate(zh)
                   for yy, row in enumerate(fonts[backend][ch]) for xx in range(height)
                   if row & 1 << (height-1-xx)}
            assert bx*scale <= x*scale and x*scale+len(zh)*(height+2)-2 <= (bx+49)*scale
            assert 188*scale <= gy and gy+height <= 197*scale
            overflow = sum(not (bx*scale <= xx < (bx+49)*scale and 188*scale <= yy < 197*scale) for xx, yy in ink)
            assert overflow == 0
            report['glyph_checks'] += len(zh)
            for prefix in ('hdat', 'drawat'):
                w, h, rgb = png(dest / f'{prefix}-400000000.png')
                for yy in range(188*scale, max(197*scale, 190*scale+height)):
                    for xx in range(bx*scale, (bx+49)*scale):
                        ox, oy = xx//scale, yy//scale
                        index = a[oy*320+ox]
                        if event['BX'] <= ox < event['BX']+event['BW'] and event['BY'] <= oy < event['BY']+event['BH'] and index == 1:
                            index = 192
                        if (xx, yy) in ink:
                            index = 1
                        at = (yy*w+xx)*3
                        assert rgb[at:at+3] == pal[index*3:index*3+3], (backend, scale, source, prefix, xx, yy)
                        report['pixel_checks'] += 1
            report['buttons'].append(dict(backend=backend, scale=scale, source=source,
                source_xy=[x, 190], paint_y=gy, bounds=[bx, 188, 49, 9],
                glyph_height=height, geometry_overflow=0, ink_outside=0))
        returned = json.loads((dest / 'hdat-528000000.json').read_text())
        assert not {s[0] for s in labels} & {e['S'] for e in returned.get('Shown') or []}
with args.output.open('x', encoding='utf-8') as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
    f.write('\n')
print({k: v for k, v in report.items() if k not in ('buttons', 'before_after')})
print('full_images', len(report['before_after']), 'overflow 0; outside approved regions 0')
