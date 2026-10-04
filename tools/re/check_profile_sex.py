"""024限定驗收：正常GUI、獨立性別探針、原生字模與修前／修後像素。"""
import argparse
import base64
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
p.add_argument('--control-root', type=Path)
p.add_argument('--unifont-hex', type=Path, required=True)
p.add_argument('--eten-font', type=Path, required=True)
p.add_argument('--eten-symbol-font', type=Path, required=True)
p.add_argument('--output', type=Path, required=True)
args = p.parse_args()
if args.output.exists():
    p.error('輸出已存在，請指定新報告路徑')
repo = Path(__file__).resolve().parents[2]
runs = args.root / 'retry2'
control_root = args.control_root or args.root
control_runs = control_root / 'retry2'
sha = lambda path: hashlib.sha256(path.read_bytes()).hexdigest()
spec = importlib.util.spec_from_file_location('pngcheck', repo/'tools/re/check_native_gui.py')
decoder = importlib.util.module_from_spec(spec)
spec.loader.exec_module(decoder)
png = lru_cache(maxsize=12)(decoder.png)
needed = set('男女確定取消')
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
symbols = args.eten_symbol_font.read_bytes()
assert hashlib.sha256(symbols).hexdigest() == 'f32049ba2a7a21db908878a488a2c1d93c389d17398cf46db1390ba89e247605'
for ch in needed:
    lead, trail = ch.encode('big5')
    col = trail-64 if trail < 127 else trail-98
    record = (lead-161)*157+col if lead < 164 else (lead-164)*157+col-(408 if lead >= 201 else 0)
    source = symbols if lead < 164 else raw
    assert record >= 0 and record*30+30 <= len(source)
    fonts['eten'][ch] = [n >> 1 for n in struct.unpack_from('>15H', source, record*30)]
assert all(set(value) == needed for value in fonts.values())
report = dict(same_state_pairs=0, probe_snapshot_pairs=0, native_glyph_checks=0,
              independent_pixel_checks=0, unchanged_outside_regions=0,
              source_field_checks=0, primaries={}, resources_unchanged=0,
              other_missing_events_unchanged=[],
              limits='024限定：3倍兩性別與兩倍率選人按鈕；2倍性別保留英文並記錄，其他人物頁及全遊戲未完成。')
report['current_binary_sha256'] = sha(args.root/'wiz7')
report['control_binary_sha256'] = sha(control_root/'wiz7')
assert report['current_binary_sha256'] == json.loads((args.root/'build-artifacts.json').read_text())['wiz7']
assert report['control_binary_sha256'] == json.loads((control_root/'build-artifacts.json').read_text())['wiz7']
steps = tuple(272000000+i*128000000 for i in range(7))


def load(path):
    return json.loads(path.read_text())


def compare_state(current, control, current_dir, control_dir):
    assert current['returncode'] == control['returncode'] == 0
    assert current['save_snapshots']['after_normal_exit'] == control['save_snapshots']['after_normal_exit']
    for step in steps:
        assert current['original'][str(step)] == control['original'][str(step)]
        for ext in ('a000', 'pal'):
            assert (current_dir/f'hdat-{step}.{ext}').read_bytes() == (control_dir/f'hdat-{step}.{ext}').read_bytes()
        report['same_state_pairs'] += 1


def outside_equal(new, old, rectangles, scale):
    w, h, a = png(new)
    ow, oh, b = png(old)
    assert (w, h) == (ow, oh) == (320*scale, 200*scale)
    for y in range(h):
        spans = sorted((x0, x1) for x0, y0, x1, y1 in rectangles if y0 <= y < y1)
        start = 0
        for left, right in spans + [(w, w)]:
            assert a[(y*w+start)*3:(y*w+left)*3] == b[(y*w+start)*3:(y*w+left)*3], (new, y, start, left)
            report['unchanged_outside_regions'] += left-start
            start = right


def restored_frame(meta, original):
    restored = bytearray(original)
    for event in meta.get('Shown') or []:
        before, after = base64.b64decode(event['Before']), base64.b64decode(event['After'])
        assert len(before) == len(after) == event['BW']*event['BH']
        for i, (b, a) in enumerate(zip(before, after)):
            at = (event['BY']+i//event['BW'])*320+event['BX']+i%event['BW']
            if b != a:
                # 原版選取框可重畫未印字的底列；覆繪只還原原筆畫。
                assert original[at] == a, (event['S'], at % 320, at // 320)
                restored[at] = b
    return restored


def check_pixels(dest, step, scale, restored, palette, labels, regions, height):
    ink = {}
    for event, x, y, zh in labels:
        for i, ch in enumerate(zh):
            for yy, row in enumerate(fonts[dest.name.split('-')[2]][ch]):
                for xx in range(height):
                    if row & (1 << (height-1-xx)):
                        ink[(x+i*(height+2)+xx, y+yy)] = event['Fg']
    for prefix in ('hdat', 'drawat'):
        w, h, rgb = png(dest/f'{prefix}-{step}.png')
        assert (w, h) == (320*scale, 200*scale)
        for x0, y0, x1, y1 in regions:
            for y in range(y0, y1):
                for x in range(x0, x1):
                    index = ink.get((x, y), restored[(y//scale)*320+x//scale])
                    at = (y*w+x)*3
                    assert rgb[at:at+3] == palette[index*3:index*3+3], (dest, prefix, step, x, y)
                    report['independent_pixel_checks'] += 1


for backend in fonts:
    for name in ('xlate.json', 'flow.json', 'help.json', 'roster.json'):
        assert (args.root/'bake'/backend/name).read_bytes() == (repo/'workplace/r932/implementation/bake'/backend/name).read_bytes()
        report['resources_unchanged'] += 1
    family = load(args.root/'bake'/backend/'profile-sex.json')
    assert family['source_sha256'] == sha(repo/'text/profile-sex.zh-TW.tsv')
    assert [(e['en'], e['zh']) for e in family['entries']] == [('M', '男'), ('F', '女')]
    assert all(e['g'] == [fonts[backend][e['zh']]] for e in family['entries'])

for sex, code, zh in (('male', 'M', '男'), ('female', 'F', '女')):
    primary_dir = runs / ('primary' if sex == 'male' else 'female-primary')
    primary = load(primary_dir/'captures.json')
    assert primary['returncode'] == 0 and primary['returned_to_menu'] and primary['save_contents_unchanged']
    assert primary['close_seconds'] <= 10
    assert [a['point'] for a in primary['actions'] if a['kind'] == 'gui_left_click'] == [
        [160,47], [85,193], [175,152 if sex == 'male' else 158], [100,170], [314,190], [287,193]]
    report['primaries'][sex] = dict(sha256=sha(primary_dir/'captures.json'), close_seconds=primary['close_seconds'])
    probe = load(args.baseline_root/f'sex-full-{sex}/observation.json')
    assert len(probe['Events']) == 1 and probe['Events'][0]['S'] == code
    raw_event = probe['Events'][0]
    for scale in (2, 3):
        english = control_runs/f'{sex}-aux-english-{scale}-new'
        control = load(english/'captures.json')
        for snapshot in probe['Snapshots']:
            assert snapshot == control['original'][str(snapshot['Step'])]
            assert (args.baseline_root/f'sex-full-{sex}'/f'original-{snapshot["Step"]}.a000').read_bytes() == (english/f'hdat-{snapshot["Step"]}.a000').read_bytes()
            report['probe_snapshot_pairs'] += 1
        for backend, height in (('eten',15), ('unifont',16)):
            dest = runs/f'{sex}-aux-{backend}-{scale}-new'
            old = (args.baseline_root/'retry'/f'aux-{backend}-{scale}') if sex == 'male' else control_runs/f'{sex}-aux-{backend}-{scale}-old'
            current = load(dest/'captures.json')
            new_missing = [line.partition('無譯文')[2] for line in (dest/'run.log').read_text().splitlines() if '無譯文' in line]
            old_missing = [line.partition('無譯文')[2] for line in (old/'run.log').read_text().splitlines() if '無譯文' in line]
            assert new_missing == old_missing and len(new_missing) == (121 if scale == 2 else 109)
            report['other_missing_events_unchanged'].append(dict(sex=sex, backend=backend, scale=scale, count=len(new_missing)))
            compare_state(current, control, dest, english)
            compare_state(current, load(old/'captures.json'), dest, old)
            for step in steps:
                meta = load(dest/f'hdat-{step}.json')
                shown = meta.get('Shown') or []
                sex_events = [e for e in shown if e['S'] == code and (e['X0'],e['Y0']) == (90,6)]
                assert bool(sex_events) == (step == steps[4] and scale == 3)
                regions, labels = [], []
                if step in (steps[2], steps[3]):
                    for source, x, y, top, word in (('OK',94,168,166,'確定'), ('CANCEL',82,184,182,'取消')):
                        event = next(e for e in shown if (e['S'],e['X0'],e['Y0']) == (source,x,y))
                        assert event['Dy'] == 0 and event['G'] == [fonts[backend][ch] for ch in word]
                        report['native_glyph_checks'] += len(word)
                        placed_y = top*scale+(9*scale-height+1)//2
                        assert top*scale <= placed_y and placed_y+height <= (top+9)*scale
                        assert 63*scale <= x*scale and x*scale+2*(height+2)-2 <= 136*scale
                        labels.append((event,x*scale,placed_y,word))
                        regions.append((63*scale,top*scale,136*scale,(top+9)*scale))
                if sex_events:
                    event, = sex_events
                    assert event['Dy'] == 0 and event['Fg'] == 1 and event['G'] == [fonts[backend][zh]]
                    assert (event['BX'],event['BY'],event['BW'],event['BH']) == (90,6,6,6)
                    assert event['Before'] == raw_event['Before'] and event['After'] == raw_event['After']
                    report['source_field_checks'] += 36
                    report['native_glyph_checks'] += 1
                    labels.append((event,270,18,zh))
                    regions.append((270,18,288,36))
                a = (english/f'hdat-{step}.a000').read_bytes()
                pal = (english/f'hdat-{step}.pal').read_bytes()
                if regions:
                    check_pixels(dest,step,scale,restored_frame(meta,a),pal,labels,regions,height)
                allowed = []
                if scale == 2 and step in (steps[2],steps[3]):
                    for event, x, placed_y, word in labels:
                        old_y = event['Y0']*scale
                        width = len(word)*(height+2)-2
                        # 兩字格上下重疊，聯集恰為此窄矩形，含修前外溢下緣。
                        assert abs(old_y-placed_y) < height
                        allowed.append((x,min(old_y,placed_y),x+width,max(old_y,placed_y)+height))
                elif sex_events:
                    allowed = regions
                for prefix in ('hdat','drawat'):
                    outside_equal(dest/f'{prefix}-{step}.png', old/f'{prefix}-{step}.png', allowed, scale)
            if scale == 2:
                assert '性別缺譯：倍率、資源、印字訊號或離屏不符' in (dest/'run.log').read_text()
with args.output.open('x', encoding='utf-8') as file:
    json.dump(report,file,ensure_ascii=False,indent=2)
    file.write('\n')
print(json.dumps(report,ensure_ascii=False))
