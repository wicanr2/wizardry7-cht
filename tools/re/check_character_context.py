from pathlib import Path
import base64
import hashlib
import json
import stat

root = Path('/out')
sha = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
expected_binary = '22c6b9b93e87164886c6df1de16b1d26c502655165bc2da08bf3f5794c85ca61'
expected_wire = 'e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732'
results = {}
traces = {}
for variant in ('control', 'observe'):
    path = root / ('probe-api-' + variant)
    result = json.loads((path / 'result.json').read_text())
    assert result['variant'] == variant
    assert result['binary_sha256'] == expected_binary and result['wire_sha256'] == expected_wire
    assert result['exit_code'] == 0 and not result['invalid']
    assert result['publications'] == 1786 and result['initial'] == result['final']
    trace = path / 'states.jsonl'
    assert sha(trace) == result['trace_sha256']
    rows = [json.loads(line) for line in trace.read_text().splitlines()]
    assert len(rows) == 1786
    assert [r['publication'] for r in rows] == list(range(1, 1787))
    assert all(set(r) == {'publication', 'step', 'ticks', 'ram', 'indexed', 'palette', 'registers'} for r in rows)
    assert rows[-1]['step'] == 438400000
    assert all(rows[i]['step'] <= rows[i+1]['step'] for i in range(len(rows)-1))
    results[variant] = result
    traces[variant] = rows
assert traces['control'] == traces['observe']
for field in ('binary_sha256', 'wire_sha256', 'trace_sha256', 'initial', 'final'):
    assert results['control'][field] == results['observe'][field]
assert not list((root / 'probe-api-control/witness').iterdir())
events_path = root / 'probe-api-observe/witness/events.jsonl'
assert stat.S_IMODE(events_path.stat().st_mode) == 0o600
events = [json.loads(line) for line in events_path.read_text().splitlines()]
assert 1 <= len(events) <= 256 and events[-1]['status'] == 'summary'
assert events[0]['status'] == 'setup'
assert all(r['status'] in {'setup', 'reader', 'decoded', 'print', 'character_entry', 'character_return', 'reject', 'summary'} for r in events)
assert sum(r['status'] == 'setup' for r in events) == 1
assert sum(r['status'] == 'summary' for r in events) == 1
summary = events[-1]
assert summary['events'] == len(events)-1
assert summary['character_entries'] == sum(r['status'] == 'character_entry' for r in events)
assert summary['character_returns'] == sum(r['status'] == 'character_return' for r in events)
assert summary['reader_entries'] == sum(r['status'] == 'reader' for r in events)
assert summary['decoded_entries'] == sum(r['status'] == 'decoded' for r in events)
assert summary['printed_entries'] == sum(r['status'] == 'print' for r in events)
assert summary['return_hooks'] <= 8 and summary['character_entries'] <= 16
assert not any(summary[k] for k in ('reader_pending', 'decoded_pending', 'character_pending'))
entries = {}
characters = []
for row in events:
    if row['status'] == 'character_entry':
        assert row['number'] not in entries
        assert 332800000 <= row['step'] <= 417800000
        assert row['font'] == 0 and row['code'] in (ord('M'), ord('F')) and (row['x'], row['y']) == (90, 6)
        assert row['sp'] <= 65530 and len(base64.b64decode(row['before'], validate=True)) == 36
        entries[row['number']] = row
    if row['status'] != 'character_return':
        continue
    entry = entries.pop(row['number'])
    assert row['entry_step'] == entry['step'] and row['ss'] == entry['ss'] and row['sp'] == entry['sp'] + 2
    assert all(row[k] == entry[k] for k in ('return', 'code', 'font', 'before', 'driver', 'buffer'))
    assert 332800000 <= row['entry_step'] <= row['step'] <= 417800000
    assert row['step']-row['entry_step'] <= 1000000
    assert row['font'] == 0 and row['code'] in (ord('M'), ord('F'))
    assert (row['x'], row['y']) == (96, 6)
    before, after = base64.b64decode(row['before'], validate=True), base64.b64decode(row['after'], validate=True)
    assert len(before) == len(after) == 36
    inside = sum(a != b for a, b in zip(before, after))
    assert inside + row['outside_differences'] == row['differences']
    changed = [(90+i % 6, 6+i // 6) for i, (a, b) in enumerate(zip(before, after)) if a != b]
    if row['differences']:
        left, top, right, bottom = row['difference_bbox']
        assert 0 <= left < right <= 320 and 0 <= top < bottom <= 200
        assert all(left <= x < right and top <= y < bottom for x, y in changed)
        if not row['outside_differences']:
            assert [left, top, right, bottom] == [min(x for x, y in changed), min(y for x, y in changed), max(x for x, y in changed)+1, max(y for x, y in changed)+1]
    else:
        assert not row.get('difference_bbox')
    characters.append({**{k: row[k] for k in ('number', 'entry_step', 'step', 'return', 'code', 'driver', 'buffer', 'after_driver', 'after_buffer', 'differences', 'outside_differences')},
                       'difference_bbox': row.get('difference_bbox'), 'inside_differences': inside,
                       'before_colors': sorted(set(before)), 'after_colors': sorted(set(after))})
report = {'status': 'PASS', 'publication_pairs': 1786, 'original_fields': ['step', 'ticks', 'ram', 'indexed', 'palette', 'registers'],
          'original_trace_sha256': results['control']['trace_sha256'], 'wire_sha256': results['control']['wire_sha256'],
          'binary_sha256': results['control']['binary_sha256'], 'save_metadata_unchanged': True,
          'events_sha256': sha(events_path), 'summary': summary, 'characters': characters,
          'limits': ['Private API on/off only; no formal cross-binary replay claim.', 'Diagnostic source observations only; no new formal display route.']}
with (root / 'probe-api-comparison.json').open('x') as handle:
    handle.write(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
print(json.dumps(report, ensure_ascii=False, indent=2))
