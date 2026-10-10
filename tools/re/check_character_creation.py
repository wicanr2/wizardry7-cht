from pathlib import Path
import base64
import hashlib
import json
import os

root = Path('/out')
base = Path('/src/workplace/r1002')
sha = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
load = lambda p: json.loads(p.read_text())
core = load(root / 'probe-api-comparison.json')
baseline = load(base / 'probe-api-comparison.json')
assert core['status'] == 'PASS' and core['original_trace_sha256'] == baseline['original_trace_sha256']
expected_trace = 'bc95b65d4f308b136f4236993ee97db7f57ebb25d43df622d0697a2a421fc1af'
assert core['original_trace_sha256'] == expected_trace
assert sha(base / 'probe-api-control/states.jsonl') == expected_trace
assert sha(base / 'probe-api-observe/states.jsonl') == expected_trace
events_path = root / 'probe-api-observe/witness/events.jsonl'
assert sha(events_path) == core['events_sha256']
events = [json.loads(s) for s in events_path.read_text().splitlines()]
previous = [json.loads(s) for s in (base / 'probe-api-observe/witness/events.jsonl').read_text().splitlines()]
extensions = [r for r in events if r['status'] == 'source_extension']
assert len(extensions) == 1
original_rows = [dict(r) for r in events if r['status'] != 'source_extension']
original_rows[-1]['events'] -= 1
assert original_rows == previous
row = extensions[0]
partner = next(r for r in events if r['status'] == 'character_return' and r['number'] == row['number'])
index = events.index(row)
assert index > 0 and events[index-1] == partner
for key in ('step', 'entry_step', 'number', 'return', 'code', 'font', 'ss', 'sp', 'x', 'y', 'driver', 'buffer', 'after_driver', 'after_buffer'):
    assert row[key] == partner[key]
assert row['return'] == 0x9595 and row['code'] == ord('M') and row['font'] == 0
assert 334900000 <= row['entry_step'] <= row['step'] <= 335000000
assert row['driver'] == row['after_driver'] and row['buffer'] == row['after_buffer']
assert row['code_address_space'] == 'dosgolem runtime real-mode CS:IP'
assert row['code_base'] == [0x0110, 0x9515]
fragment = base64.b64decode(row['code_bytes'], validate=True)
assert len(fragment) == 128 and fragment[-3:] == bytes.fromhex('e8e2a2')
assert hashlib.sha256(fragment).hexdigest() == row['code_sha256']
before = base64.b64decode(row['before8'], validate=True)
after = base64.b64decode(row['after8'], validate=True)
assert row['eight_rectangle'] == [88, 4, 96, 12] and len(before) == len(after) == 64
crop = lambda b: bytes(b[y*8+x] for y in range(2, 8) for x in range(2, 8))
assert crop(before) == base64.b64decode(partner['before'], validate=True)
assert crop(after) == base64.b64decode(partner['after'], validate=True)
changed = [(i, a, b) for i, (a, b) in enumerate(zip(before, after)) if a != b]
assert len(changed) == 13 and all(b == 1 and i%8 >= 2 and i//8 >= 2 for i, a, b in changed)
points = [[88, 4], [92, 4], [95, 4], [88, 7], [89, 9], [88, 11], [91, 5], [95, 5]]
assert row['background_points'] == points
point_values = lambda b: bytes(b[(y-4)*8+x-88] for x, y in points)
assert base64.b64decode(row['background_before'], validate=True) == point_values(before)
assert base64.b64decode(row['background_after'], validate=True) == point_values(after)
background = list(point_values(before))
source_dir = root / 'ida-source'
source_dir.mkdir()
code_path = source_dir / 'creation-code.bin'
with os.fdopen(os.open(code_path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600), 'wb') as handle:
    handle.write(fragment)
metadata = {'runtime_segment': 0x0110, 'runtime_offset': 0x9515, 'length': 128, 'sha256': row['code_sha256'],
            'address_space': row['code_address_space'], 'source_event_sha256': core['events_sha256'],
            'source_entry_step': row['entry_step'], 'source_return_step': row['step']}
with os.fdopen(os.open(code_path.with_suffix('.json'), os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600), 'w') as handle:
    handle.write(json.dumps(metadata, indent=2) + '\n')
result = {'status': 'PASS', 'publication_pairs': 1786, 'trace_matches_r1': True,
          'original_trace_sha256': core['original_trace_sha256'], 'binary_sha256': core['binary_sha256'],
          'r1_events_preserved': True, 'extension_samples': 1, 'code_sha256': row['code_sha256'],
          'code_runtime_base': row['code_base'], 'code_length': len(fragment),
          'eight_by_eight_changed_pixels': len(changed), 'background_points': points,
          'background_before': background, 'background_after': list(point_values(after)),
          'matches_existing_world_background': background == [115, 116, 115, 117, 115, 115, 13, 14],
          'formal_route_implemented': False}
with (root / 'extension-comparison.json').open('x') as handle:
    handle.write(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
print(json.dumps(result, ensure_ascii=False, indent=2))
