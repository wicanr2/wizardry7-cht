"""022：原字型與原版色號獨立核對六列中文及同狀態A/B。"""
import argparse
import gzip
from functools import lru_cache
import hashlib
import importlib.util
import json
import struct
from pathlib import Path

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('root', type=Path)
p.add_argument('--unifont-hex', type=Path, required=True)
p.add_argument('--eten-font', type=Path, required=True)
p.add_argument('--output', type=Path, required=True)
args = p.parse_args()
if args.output.exists():
    p.error('輸出已存在，請指定新檔案')
project = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('pngcheck', project/'tools/re/check_native_gui.py')
decoder = importlib.util.module_from_spec(spec)
spec.loader.exec_module(decoder)
png = lru_cache(maxsize=128)(decoder.png)
rows = []
for line in (project/'text/roster.zh-TW.tsv').read_text().splitlines()[1:]:
    en, zh, context, dy, source = line.split('\t')
    rows.append((en, zh))
names = ['THESUS    M-FIG', 'TEMPEST   F-FIG', 'LYSAND    M-THI',
         'NOBAL     M-PRI', 'TREON     M-MAG', 'PENTAS    F-ALC']
assert [en for en, _ in rows] == [s[-5:] for s in names]
needed = set(''.join(zh for _, zh in rows))
fonts = {'eten': {}, 'unifont': {}}
raw = args.unifont_hex.read_bytes()
assert hashlib.sha256(raw).hexdigest() == '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
for line in gzip.decompress(raw).decode('ascii').splitlines():
    cp, bits = line.split(':')
    ch = chr(int(cp,16))
    if ch in needed:
        assert len(bits) == 64
        fonts['unifont'][ch] = [int(bits[i:i+4],16) for i in range(0,64,4)]
raw = args.eten_font.read_bytes()
assert hashlib.sha256(raw).hexdigest() == '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
for ch in needed:
    lead, trail = ch.encode('big5')
    column = trail-64 if trail < 127 else trail-98
    record = (lead-164)*157+column-(408 if lead >= 201 else 0)
    fonts['eten'][ch] = [v>>1 for v in struct.unpack_from('>15H', raw, record*30)]
assert all(set(f) == needed for f in fonts.values())
report = dict(same_state_pairs=0, glyph_checks=0, pixel_checks=0, full_frame_diff=[],
              primary_pixel_checks=0, untranslated_roster_events={}, limits='six default rows only; no 2x Chinese or generic name/scroll claim')

def check_rows(dest, backend, step):
    height = 15 if backend == 'eten' else 16
    meta = json.loads((dest/f'hdat-{step}.json').read_text())
    a = (dest/f'hdat-{step}.a000').read_bytes()
    pal = (dest/f'hdat-{step}.pal').read_bytes()
    total = 0
    for row, ((en, zh), full) in enumerate(zip(rows, names)):
        y = 150+row*6
        e = next(e for e in meta['Shown'] if e['S'] == full)
        glyphs = [fonts[backend][ch] for ch in zh]
        assert (e['X0'],e['Y0'],e['BX'],e['BY'],e['BW'],e['BH']) == (210,y,210,y,30,6)
        assert e['G'] == glyphs
        fg, bg = (0,1) if (step == 528000000 and row == 0) or (step == 656000000 and row == 5) else (180,195)
        assert e['Fg'] == fg
        ink = {(630+i*(height+2)+x, y*3+yy) for i,g in enumerate(glyphs)
               for yy,word in enumerate(g) for x in range(height) if word & (1<<(height-1-x))}
        assert all(630 <= x < 720 and y*3 <= yy < (y+6)*3 for x,yy in ink)
        report['glyph_checks'] += len(glyphs)
        for prefix in ('hdat','drawat'):
            w,h,rgb = png(dest/f'{prefix}-{step}.png')
            assert (w,h) == (960,600)
            for yy in range(y*3,(y+6)*3):
                for x in range(630,720):
                    index = a[(yy//3)*320+x//3]
                    if index == fg:
                        index = bg
                    if (x,yy) in ink:
                        index = fg
                    at = (yy*w+x)*3
                    assert rgb[at:at+3] == pal[index*3:index*3+3], (dest.name,step,prefix,row,x,yy)
                    total += 1
    return total

for backend in ('eten','unifont'):
    baked = json.loads((args.root/'bake'/backend/'roster.json').read_text())
    assert [e['g'] for e in baked['entries']] == [[fonts[backend][c] for c in zh] for _,zh in rows]
    for scale in (2,3):
        new = args.root/f'aux-new-{backend}-{scale}'
        old = args.root/f'aux-old-{backend}-{scale}'
        newcap, oldcap = [json.loads((d/'captures.json').read_text()) for d in (new,old)]
        assert newcap['returncode'] == oldcap['returncode'] == 0
        assert newcap['save_snapshots']['after_normal_exit'] == oldcap['save_snapshots']['after_normal_exit']
        for step in (272000000,400000000,528000000,656000000,784000000):
            assert newcap['original'][str(step)] == oldcap['original'][str(step)]
            for ext in ('a000','pal'):
                assert (new/f'hdat-{step}.{ext}').read_bytes() == (old/f'hdat-{step}.{ext}').read_bytes()
            report['same_state_pairs'] += 1
            if scale == 3 and step in (400000000,528000000,656000000):
                report['pixel_checks'] += check_rows(new,backend,step)
            for prefix in ('hdat','drawat'):
                w,h,rgb = png(new/f'{prefix}-{step}.png')
                ow,oh,original = png(old/f'{prefix}-{step}.png')
                assert (w,h) == (ow,oh)
                changed = outside = 0
                for i in range(w*h):
                    if rgb[i*3:i*3+3] != original[i*3:i*3+3]:
                        changed += 1
                        x,y = i%w,i//w
                        allowed = scale == 3 and step in (400000000,528000000,656000000) and 630 <= x < 720 and 450 <= y < 558
                        outside += not allowed
                assert outside == 0
                report['full_frame_diff'].append(dict(backend=backend,scale=scale,step=step,kind=prefix,changed=changed,outside=outside))
        log = (new/'run.log').read_text()
        roster_missing = [line for line in log.splitlines() if '無譯文' in line and 'ret=0110:0fcc' in line]
        assert len(roster_missing) == (18 if scale == 2 else 0)
        assert '職業碼缺譯' not in log
        report['untranslated_roster_events'][f'{backend}-{scale}'] = len(roster_missing)
primary = args.root/'primary'
for step in (400000000,528000000,656000000):
    report['primary_pixel_checks'] += check_rows(primary,'unifont',step)
for d in [primary]+list(args.root.glob('aux-*')):
    meta=json.loads((d/'hdat-784000000.json').read_text())
    assert not set(names)&{e['S'] for e in meta.get('Shown') or []}
with args.output.open('x',encoding='utf-8') as f:
    json.dump(report,f,ensure_ascii=False,indent=2)
    f.write('\n')
print(json.dumps({k:v for k,v in report.items() if k!='full_frame_diff'},ensure_ascii=False))
