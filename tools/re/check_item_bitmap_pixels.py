from pathlib import Path
import argparse
import base64
import json
import struct

parser = argparse.ArgumentParser(description="獨立核對R999私人服務像素；不作正式路由或參數判斷。")
parser.add_argument("--witness-dir", type=Path, required=True)
parser.add_argument("--output", type=Path, required=True)
args = parser.parse_args()
root = args.witness_dir
rows = []
for path in sorted(root.glob('service24-*-return.json')):
    sample = json.loads(path.read_text())
    if not sample['Accepted']:
        rows.append({'id': sample['ID'], 'accepted': False, 'reason': sample['Reason']})
        continue
    stem = path.name[:-12]
    roi = sample['Indexed']
    record = {
        'id': sample['ID'], 'nul_id': sample['NULID'],
        'entry_step': sample['EntryStep'], 'return_step': sample['ReturnStep'],
        'entry_regs': sample['Entry'], 'return_regs': sample['Returned'],
        'entry_stack_words': struct.unpack('<8H', base64.b64decode(sample['EntryStack'])),
        'return_stack_words': struct.unpack('<8H', base64.b64decode(sample['ReturnStack'])),
        'entry_code_hex': base64.b64decode(sample['EntryCode']).hex(),
        'return_code_hex': base64.b64decode(sample['ReturnCode']).hex(),
        'offscreen_diff': sample['Offscreen'], 'indexed_diff': roi,
    }
    if roi['Changed']:
        before = (root / (stem + '-indexed-before.bin')).read_bytes()
        after = (root / (stem + '-indexed-after.bin')).read_bytes()
        mask = (root / (stem + '-indexed-mask.bin')).read_bytes()
        assert len(before) == len(after) == len(mask) == roi['W'] * roi['H']
        assert sum(mask) == roi['Changed']
        assert all(bool(k) == (b != a) for b, a, k in zip(before, after, mask))
        nul = root / ('longsword-%03d-return.json' % sample['NULID'])
        source = json.loads(nul.read_text())['Offscreen']
        prefix = 'longsword-%03d-offscreen-' % sample['NULID']
        text = (root / (prefix + 'after.bin')).read_bytes()
        text_mask = (root / (prefix + 'mask.bin')).read_bytes()
        ink = [(i % source['W'], i // source['W'], v) for i, (v, k) in enumerate(zip(text, text_mask)) if k]
        candidates = {}
        for phase, data in [('before', before), ('after', after)]:
            matches = []
            for y in range(max(0, roi['H'] - source['H'] + 1)):
                for x in range(max(0, roi['W'] - source['W'] + 1)):
                    if all(data[(y + dy) * roi['W'] + x + dx] == value for dx, dy, value in ink):
                        matches.append({'x': roi['X'] + x, 'y': roi['Y'] + y, 'ink_points': len(ink), 'changed_ink_points':sum(mask[(y + dy) * roi['W'] + x + dx] for dx, dy, value in ink)})
            candidates[phase] = matches
        record['translated_source_ink_matches'] = candidates
    rows.append(record)
output = {'observations': rows, 'limits': '墨點平移吻合僅為資料對齊證據，不用形狀搜尋建立正式中文路由，不推stack參數或bitmap來源語意。'}
with args.output.open("x") as report:
    report.write(json.dumps(output, indent=2) + '\n')
print(json.dumps(output))
