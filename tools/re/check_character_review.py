"""R935修正前人物檢視：原版摘要、原字型及3倍局部像素核對。"""
import argparse
import base64
import gzip
import hashlib
import importlib.util
import json
import struct
from functools import lru_cache
from pathlib import Path

p = argparse.ArgumentParser()
p.add_argument('root', type=Path)
p.add_argument('--unifont-hex', type=Path, required=True)
p.add_argument('--eten-font', type=Path, required=True)
p.add_argument('--eten-symbol-font', type=Path, required=True)
p.add_argument('--output', type=Path, required=True)
args = p.parse_args()
if args.output.exists():
    p.error('輸出已存在，請指定新報告路徑')
root, out = Path(__file__).resolve().parents[2], args.root
runs = out / 'retry'
sha = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
spec = importlib.util.spec_from_file_location('pngcheck', root / 'tools/re/check_native_gui.py')
decoder = importlib.util.module_from_spec(spec)
spec.loader.exec_module(decoder)
png = lru_cache(maxsize=64)(decoder.png)
profile = [('HUMAN',102,6,'人類'),('FIGHTER',90,12,'戰士'),('LVL',90,18,'等級'),
           ('RNK',228,6,'階級'),('NONE',288,6,'無'),('EXP',228,12,'經驗'),('MKS',228,18,'擊殺'),
           ('HP',6,93,'生命'),('ST',6,99,'體力'),('STR',115,35,'力量'),('INT',115,41,'智力'),
           ('PIE',115,47,'虔誠'),('VIT',115,53,'體質'),('DEX',115,59,'敏捷'),
           ('SPD',115,65,'速度'),('PER',115,71,'人格'),('KAR',115,77,'因果')]
selector = [('REVIEW WHO?',67,152,'檢視誰？'),('OK',94,168,'確定'),('CANCEL',82,184,'取消')]
needed = set(''.join(v[3] for v in profile+selector))
fonts = {'eten':{}, 'unifont':{}}
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
symbols = args.eten_symbol_font.read_bytes()
assert len(symbols) == 12240 and hashlib.sha256(symbols).hexdigest() == 'f32049ba2a7a21db908878a488a2c1d93c389d17398cf46db1390ba89e247605'
for ch in needed:
    lead,trail = ch.encode('big5')
    col = trail-64 if trail<127 else trail-98
    record = (lead-161)*157+col if lead<164 else (lead-164)*157+col-(408 if lead>=201 else 0)
    source = symbols if lead<164 else raw
    assert record>=0 and record*30+30<=len(source)
    fonts['eten'][ch] = [n>>1 for n in struct.unpack_from('>15H',source,record*30)]
assert all(set(v)==needed for v in fonts.values())
report = dict(same_state_pairs=0, native_glyph_checks=0, selector_buttons=[], profile_labels=17,
              profile_pixel_checks_3x=0, before_after_source_pixels_checked=0,
              raw_male_probe_snapshots_equal=0, limits='Current renderer audit; M remains English; female GUI and other tabs not completed.')
for scale in (2,3):
    english = runs / f'aux-english-{scale}'
    base = json.loads((english/'captures.json').read_text())
    for backend in ('eten','unifont'):
        dest = runs / f'aux-{backend}-{scale}'
        current = json.loads((dest/'captures.json').read_text())
        assert current['returncode']==base['returncode']==0
        assert current['save_snapshots']['after_normal_exit']==base['save_snapshots']['after_normal_exit']
        for step, expected in base['original'].items():
            assert current['original'][step]==expected
            assert all((dest/f'hdat-{step}.{ext}').read_bytes()==(english/f'hdat-{step}.{ext}').read_bytes() for ext in ('a000','pal'))
            report['same_state_pairs'] += 1
        height = 15 if backend=='eten' else 16
        for step in (528000000,656000000,784000000):
            meta = json.loads((dest/f'hdat-{step}.json').read_text())
            expected_labels = profile if step==784000000 else selector
            for source,x,y,zh in expected_labels:
                e = next(e for e in meta['Shown'] if (e['S'],e['X0'],e['Y0'])==(source,x,y))
                assert e['Dy']==0 and e['G']==[fonts[backend][ch] for ch in zh], (backend, scale, source, [i for i,(actual,wanted) in enumerate(zip(e['G'],[fonts[backend][ch] for ch in zh])) if actual!=wanted])
                report['native_glyph_checks'] += len(zh)
                if source in ('OK','CANCEL'):
                    top = 166 if source=='OK' else 182
                    ink = {(x*scale+i*(height+2)+xx,y*scale+yy) for i,ch in enumerate(zh)
                           for yy,row in enumerate(fonts[backend][ch]) for xx in range(height) if row & (1<<(height-1-xx))}
                    report['selector_buttons'].append(dict(backend=backend,scale=scale,step=step,source=source,
                        geometry_overflow=max(0,y*scale+height-(top+9)*scale),
                        ink_outside=sum(not(63*scale<=xx<136*scale and top*scale<=yy<(top+9)*scale) for xx,yy in ink)))
        if scale==3:
            step=784000000
            meta=json.loads((dest/f'hdat-{step}.json').read_text())
            a=(english/f'hdat-{step}.a000').read_bytes()
            pal=(english/f'hdat-{step}.pal').read_bytes()
            # The expectation uses raw original pixels and sampled original before/after bytes,
            # then draws independently decoded font bits. No renderer helper is imported.
            restored=bytearray(a)
            for e in meta['Shown']:
                before,after=base64.b64decode(e['Before']),base64.b64decode(e['After'])
                assert len(before)==len(after)==e['BW']*e['BH']
                for i,(b,c) in enumerate(zip(before,after)):
                    x,y=e['BX']+i%e['BW'],e['BY']+i//e['BW']
                    assert a[y*320+x]==c
                    report['before_after_source_pixels_checked'] += 1
                    if b!=c:
                        restored[y*320+x]=b
            ink={}
            regions=set()
            for source,x,y,zh in profile:
                e=next(e for e in meta['Shown'] if e['S']==source)
                for i,ch in enumerate(zh):
                    for yy,row in enumerate(fonts[backend][ch]):
                        for xx in range(height):
                            if row & (1<<(height-1-xx)):
                                ink[(x*3+i*(height+2)+xx,y*3+yy)]=e['Fg']
                regions.update((xx,yy) for yy in range(y*3,y*3+max(e['BH']*3,height))
                               for xx in range(x*3,x*3+max(e['BW']*3,len(zh)*(height+2)-2)))
            for prefix in ('hdat','drawat'):
                w,h,rgb=png(dest/f'{prefix}-{step}.png')
                assert (w,h)==(960,600)
                for x,y in regions:
                    index=ink.get((x,y),restored[(y//3)*320+x//3])
                    at=(y*w+x)*3
                    assert rgb[at:at+3]==pal[index*3:index*3+3],(backend,prefix,x,y)
                    report['profile_pixel_checks_3x'] += 1
for name in ('sex-observation','sex-full-male'):
    probe=json.loads((out/name/'observation.json').read_text())
    base=json.loads((runs/'aux-english-3/captures.json').read_text())
    for s in probe['Snapshots']:
        assert s==base['original'][str(s['Step'])]
        assert (out/name/f'original-{s["Step"]}.a000').read_bytes()==(runs/'aux-english-3'/f'hdat-{s["Step"]}.a000').read_bytes()
        report['raw_male_probe_snapshots_equal'] += 1
primary=json.loads((runs/'primary/captures.json').read_text())
assert primary['returncode']==0 and primary['returned_to_menu'] and primary['save_contents_unchanged']
report['primary_close_seconds']=primary['close_seconds']
with args.output.open('x') as f:
    json.dump(report,f,ensure_ascii=False,indent=2);f.write('\n')
print({k:v for k,v in report.items() if k!='selector_buttons'})
print('selector_overflow',[(v['backend'],v['scale'],v['source'],v['geometry_overflow'],v['ink_outside']) for v in report['selector_buttons']])
