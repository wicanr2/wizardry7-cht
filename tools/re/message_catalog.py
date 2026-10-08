"""讀取文字來源前先略過防拷記錄正文；不改其他JSON內容。"""
import json
import re
from pathlib import Path

_PREFIX = re.compile(
    rb'\{\s*"id"\s*:\s*(0|[1-9][0-9]*)\s*,\s*"blk"\s*:\s*(0|[1-9][0-9]*)'
    rb'\s*,\s*"off"\s*:\s*(0|[1-9][0-9]*)\s*,\s*"en"\s*:')
_ROOT_KEY = re.compile(rb'"(count|maxID|note|schema|tree|records)"\s*:')
_NUMBER = re.compile(rb'(0|[1-9][0-9]*)')


def _space(raw, pos):
    while pos < len(raw) and raw[pos] in b' \t\r\n':
        pos += 1
    return pos


def _string_end(raw, pos):
    if pos >= len(raw) or raw[pos] != ord('"'):
        raise ValueError('catalog_string')
    pos += 1
    escaped = False
    while pos < len(raw):
        c = raw[pos]
        if escaped:
            escaped = False
        elif c == ord('\\'):
            escaped = True
        elif c == ord('"'):
            return pos + 1
        pos += 1
    raise ValueError('catalog_string_end')


def _array_end(raw, pos):
    stack = [ord(']')]
    pos += 1
    while pos < len(raw):
        c = raw[pos]
        if c == ord('"'):
            pos = _string_end(raw, pos)
            continue
        if c in b'[{':
            stack.append(ord(']') if c == ord('[') else ord('}'))
        elif c in b']}':
            if not stack or c != stack.pop():
                raise ValueError('catalog_brackets')
            if not stack:
                return pos + 1
        pos += 1
    raise ValueError('catalog_array_end')


def _root_records(raw):
    pos = _space(raw, 0)
    if pos >= len(raw) or raw[pos] != ord('{'):
        raise ValueError('catalog_root')
    pos, seen, records = pos + 1, set(), None
    while True:
        pos = _space(raw, pos)
        key = _ROOT_KEY.match(raw, pos)
        if key is None or key[1] in seen:
            raise ValueError('catalog_root_key')
        name = key[1]
        seen.add(name)
        pos = _space(raw, key.end())
        if name == b'records':
            if pos >= len(raw) or raw[pos] != ord('['):
                raise ValueError('catalog_records')
            end = _array_end(raw, pos)
            records, pos = (pos + 1, end - 1), end
        elif name in (b'count', b'maxID'):
            number = _NUMBER.match(raw, pos)
            if number is None:
                raise ValueError('catalog_number')
            pos = number.end()
        else:
            pos = _string_end(raw, pos)
        pos = _space(raw, pos)
        if pos < len(raw) and raw[pos] == ord(','):
            pos += 1
        elif pos < len(raw) and raw[pos] == ord('}'):
            if _space(raw, pos + 1) != len(raw) or records is None:
                raise ValueError('catalog_root_end')
            return records
        else:
            raise ValueError('catalog_root_separator')


def decode_player_catalog(raw):
    """只解碼非排除正文；排除記錄保留三個數值與空en。"""
    if not isinstance(raw, bytes) or len(raw) > 64 * 1024 * 1024:
        raise ValueError('catalog_size')
    start_array, end_array = _root_records(raw)
    pos, copied, parts, skipped = start_array, 0, [], 0
    while True:
        while pos < len(raw) and raw[pos] in b' \t\r\n':
            pos += 1
        if pos > end_array:
            raise ValueError('catalog_array_end')
        if raw[pos] == ord(']'):
            break
        prefix = _PREFIX.match(raw, pos)
        if prefix is None:
            raise ValueError('catalog_record_prefix')
        ident, block, offset = map(int, prefix.groups())
        text_start = _space(raw, prefix.end())
        text_end = _string_end(raw, text_start)
        pos = _space(raw, text_end)
        if pos >= end_array or raw[pos] != ord('}'):
            raise ValueError('catalog_record_end')
        pos += 1
        if 2500 <= ident <= 2574:
            parts.append(raw[copied:text_start])
            parts.append(b'""')
            copied, skipped = text_end, skipped + 1
        while pos < len(raw) and raw[pos] in b' \t\r\n':
            pos += 1
        if pos < len(raw) and raw[pos] == ord(','):
            pos += 1
        elif pos < len(raw) and raw[pos] == ord(']'):
            break
        else:
            raise ValueError('catalog_record_separator')
    parts.append(raw[copied:])
    try:
        data = json.loads(b''.join(parts))
    except (ValueError, UnicodeError):
        raise ValueError('catalog_json') from None
    records = data.get('records') if isinstance(data, dict) else None
    if not isinstance(records, list):
        raise ValueError('catalog_records')
    ids = [r.get('id') for r in records if isinstance(r, dict)]
    if len(ids) != len(records) or any(type(i) is not int for i in ids) or len(set(ids)) != len(ids):
        raise ValueError('catalog_ids')
    if any(not isinstance(r.get('en'), str) for r in records):
        raise ValueError('catalog_text')
    return data, skipped


def load_player_catalog(path):
    return decode_player_catalog(Path(path).read_bytes())[0]
