import argparse
import gzip
import hashlib
import importlib.util
import json
import struct
from collections import Counter
from pathlib import Path

parser = argparse.ArgumentParser(description='020主選單字格、筆畫與原版同狀態驗證')
parser.add_argument('root', type=Path)
parser.add_argument('--bounds', type=Path, required=True)
parser.add_argument('--baseline-bake-root', type=Path, required=True)
parser.add_argument('--unifont-hex', type=Path, required=True)
parser.add_argument('--eten-font', type=Path, required=True)
parser.add_argument('--output', type=Path, required=True)
args = parser.parse_args()
if args.output.exists():
    parser.error('輸出已存在；請另指定新檔，保留舊收據')
out = args.root
spec = importlib.util.spec_from_file_location('pngcheck', str(Path(__file__).with_name('check_native_gui.py')))
pngcheck = importlib.util.module_from_spec(spec)
spec.loader.exec_module(pngcheck)
bounds = json.loads(args.bounds.read_text())
words = ['角色選單', '新增角色', '載入存檔', '系統設定', '匯入巫術六', '結束遊戲']
needed = set(''.join(words))
raw_hex = args.unifont_hex.read_bytes()
assert hashlib.sha256(raw_hex).hexdigest() == '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
unifont = {}
for line in gzip.decompress(raw_hex).decode('ascii').splitlines():
    cp, bits = line.split(':')
    ch = chr(int(cp, 16))
    if ch in needed:
        assert len(bits) == 64
        unifont[ch] = [int(bits[i:i + 4], 16) for i in range(0, 64, 4)]
assert unifont.keys() == needed
font = args.eten_font.read_bytes()
eten = {}
for ch in needed:
    lead, trail = ch.encode('big5')
    column = trail - 64 if trail < 127 else trail - 98
    record = (lead - 164) * 157 + column - (408 if lead >= 201 else 0)
    eten[ch] = [value >> 1 for value in struct.unpack_from('>15H', font, record * 30)]

def image(p):
    return pngcheck.png(p)

report = {'cases': [], 'resource_changes': {}, 'checks': 0}
for backend in ('eten', 'unifont'):
    old_flow = json.loads((args.baseline_bake_root / backend / 'flow.json').read_text())
    new_flow = json.loads((out / backend / 'flow.json').read_text())
    assert old_flow.keys() == new_flow.keys()
    assert old_flow['size'] == new_flow['size']
    assert old_flow['msgs'].keys() == new_flow['msgs'].keys()
    changed_lines, changed_msgs = [], []
    for a, b in zip(old_flow['paras'], new_flow['paras'], strict=True):
        assert {k: v for k, v in a.items() if k != 'lines'} == {k: v for k, v in b.items() if k != 'lines'}
        for x, y in zip(a['lines'], b['lines'], strict=True):
            if x != y:
                assert x['id'] in (1007, 1415) and x['zh'] == '' and y['zh'] == words[4]
                assert {k: v for k, v in x.items() if k != 'zh'} == {k: v for k, v in y.items() if k != 'zh'}
                changed_lines.append(x['id'])
    for key in old_flow['msgs']:
        a, b = old_flow['msgs'][key], new_flow['msgs'][key]
        if a != b:
            assert key in ('1007', '1415') and a['z'] == '' and b['z'] == words[4]
            assert {k: v for k, v in a.items() if k != 'z'} == {k: v for k, v in b.items() if k != 'z'}
            changed_msgs.append(key)
    removed = old_flow['glyphs'].keys() - new_flow['glyphs'].keys()
    added = new_flow['glyphs'].keys() - old_flow['glyphs'].keys()
    assert not removed and added <= set(words[4])
    assert all(old_flow['glyphs'][k] == new_flow['glyphs'][k] for k in old_flow['glyphs'])
    assert set(changed_lines) == {1007, 1415} and set(changed_msgs) == {'1007', '1415'}
    report['resource_changes'][backend] = {'changed_flow_lines': changed_lines,
        'changed_flow_msgs': changed_msgs, 'added_glyph_codepoints': [ord(c) for c in sorted(added)],
        'other_flow_fields_equal': True}

for backend, glyphs in (('eten', eten), ('unifont', unifont)):
    height = 15 if backend == 'eten' else 16
    for scale in (2, 3):
        current = out / 'gui' / f'{backend}-{scale}'
        previous = out / 'gui-before' / f'{backend}-{scale}'
        english = out / 'gui' / f'english-{scale}'
        a = (english / 'hdat-272000000.a000').read_bytes()
        pal = (english / 'hdat-272000000.pal').read_bytes()
        meta = json.loads((current / 'hdat-272000000.json').read_text())
        prevmeta = json.loads((previous / 'hdat-272000000.json').read_text())
        for key in ('Step', 'Mem', 'A000', 'Pal', 'Mouse'):
            assert meta[key] == prevmeta[key]
        for extension in ('a000', 'pal', 'comp'):
            assert (current / f'hdat-272000000.{extension}').read_bytes() == (previous / f'hdat-272000000.{extension}').read_bytes()
        assert len(meta['Shown']) == 6
        w, h, hd = image(current / 'hdat-272000000.png')
        dw, dh, draw = image(current / 'drawat-272000000.png')
        pw, ph, before = image(previous / 'hdat-272000000.png')
        assert (w, h) == (dw, dh) == (pw, ph) == (320 * scale, 200 * scale)
        mask = bytearray(w * h)
        glyph_checks, pixel_checks = 0, 0
        for row, zh in zip(bounds, words, strict=True):
            source, xy, rect = row['english'], row['text_xy'], row['button_xywh']
            event = next(e for e in meta['Shown'] if e['S'] == source)
            assert [event['X0'], event['Y0']] == xy and event['Dy'] == 0
            expected_glyphs = [glyphs[ch] for ch in zh]
            assert event['G'] == expected_glyphs
            glyph_checks += len(zh)
            bx, by, bw, bh = rect
            colors = Counter(a[y * 320 + x] for y in range(by + 1, by + bh - 1) for x in range(bx + 1, bx + bw - 1)
                             if a[y * 320 + x] not in (176, 192, 197))
            assert len(colors) == 1
            foreground = next(iter(colors))
            assert foreground == event['Fg']
            fg = pal[foreground * 3:foreground * 3 + 3]
            x0 = xy[0] * scale
            y0 = xy[1] * scale - (2 if backend == 'eten' else 3) if scale == 2 else xy[1] * 3
            width = len(zh) * (height + 2) - 2
            assert bx * scale <= x0 and x0 + width <= (bx + bw) * scale
            assert by * scale <= y0 and y0 + height <= (by + bh) * scale
            ink = set()
            for letter, bitmap in enumerate(expected_glyphs):
                for yy, value in enumerate(bitmap):
                    for xx in range(height):
                        if value & (1 << (height - 1 - xx)):
                            ink.add((x0 + letter * (height + 2) + xx, y0 + yy))
            for y in range(by * scale, (by + bh) * scale):
                for x in range(bx * scale, (bx + bw) * scale):
                    index = a[(y // scale) * 320 + x // scale]
                    if index == foreground:
                        index = a[(y // scale) * 320 + 106]
                    expected = fg if (x, y) in ink else pal[index * 3:index * 3 + 3]
                    offset = (y * w + x) * 3
                    assert hd[offset:offset + 3] == expected, (backend, scale, source, 'HD', x, y)
                    assert draw[offset:offset + 3] == expected, (backend, scale, source, 'Draw', x, y)
                    pixel_checks += 2
            for start_y in (y0, xy[1] * scale):
                for y in range(start_y, start_y + height):
                    for x in range(x0, x0 + width):
                        mask[y * w + x] = 1
        changed, outside = 0, 0
        for pixel in range(w * h):
            i = pixel * 3
            if hd[i:i + 3] != before[i:i + 3]:
                changed += 1
                outside += not mask[pixel]
        assert outside == 0
        if scale == 3:
            assert changed == 0
        else:
            assert changed > 0
        report['cases'].append({'backend': backend, 'scale': scale,
            'independent_glyph_checks': glyph_checks, 'button_pixel_checks': pixel_checks,
            'changed_pixels': changed, 'outside_approved_regions': outside,
            'original_and_comp_equal': True})
        report['checks'] += glyph_checks + pixel_checks
with args.output.open('x', encoding='utf-8') as receipt:
    receipt.write(json.dumps(report, indent=1))
print(json.dumps(report))
