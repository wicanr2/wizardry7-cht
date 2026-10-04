#!/usr/bin/env python3
"""024：只烘製人物頁性別的獨立TSV家族，不讀原版訊息。"""
import argparse
import hashlib
import json
import unicodedata
from pathlib import Path

CATALOG = Path(__file__).resolve().parents[2] / 'text/profile-sex.zh-TW.tsv'
CONTEXT = 'character.review.sex'
SOURCE = 'runtime:0110:6c42'


def build(catalog, glyph_char, size):
    raw = Path(catalog).read_bytes()
    if size not in (15, 16) or raw.startswith(b'\xef\xbb\xbf'):
        raise ValueError('性別字高或TSV編碼不合格')
    text = raw.decode('utf-8')
    if any((ord(ch) < 0x20 and ch not in '\t\r\n') or
           0x7f <= ord(ch) <= 0x9f or ch in '\u2028\u2029' for ch in text):
        raise ValueError('性別TSV含禁用控制字元或列分隔符')
    text = text.replace('\r\n', '\n')
    if '\r' in text:
        raise ValueError('性別TSV只接受LF或CRLF分列')
    lines = text.split('\n')
    if lines[-1] == '':
        lines.pop()
    if not lines or lines[0] != 'en\tzh\tcontext\tdy\tsource':
        raise ValueError('性別TSV表頭不合格')
    entries, seen = [], set()
    for line in lines[1:]:
        fields = line.split('\t')
        if len(fields) != 5:
            raise ValueError('性別TSV欄數不合格')
        en, zh, context, dy, source = fields
        if en not in {'M', 'F'} or en in seen or context != CONTEXT or dy != '0' or source != SOURCE:
            raise ValueError('性別TSV鍵或情境不合格')
        if len(zh) != 1 or not unicodedata.name(zh, '').startswith('CJK UNIFIED IDEOGRAPH-'):
            raise ValueError('性別譯文須為單一漢字')
        rows = glyph_char(zh)
        if len(rows) != size or not any(rows) or any(type(v) is not int or not 0 <= v < 1 << size for v in rows):
            raise ValueError('性別原生字模不合格')
        entries.append(dict(en=en, zh=zh, dy=0, g=[rows]))
        seen.add(en)
    if seen != {'M', 'F'}:
        raise ValueError('性別TSV須有完整兩鍵')
    return dict(schema=1, size=size, context=CONTEXT, source=SOURCE,
                source_sha256=hashlib.sha256(raw).hexdigest(), entries=entries)


def write_payload(path, payload):
    with Path(path).open('x', encoding='utf-8') as f:
        json.dump(payload, f, ensure_ascii=False, separators=(',', ':'))
        f.write('\n')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('output', type=Path)
    parser.add_argument('--catalog', type=Path, default=CATALOG)
    parser.add_argument('--eten-font', type=Path)
    parser.add_argument('--eten-symbol-font', type=Path)
    args = parser.parse_args()
    if args.output.exists():
        parser.error('輸出已存在，請指定新路徑')
    import bake_font
    if args.eten_font:
        bake_font.FONT15 = args.eten_font
    if args.eten_symbol_font:
        bake_font.SPC15 = args.eten_symbol_font
    payload = build(args.catalog, bake_font.glyph_char, bake_font.HEIGHT)
    write_payload(args.output, payload)
    print(f'性別：兩鍵，字高{payload["size"]}，缺字0')


if __name__ == '__main__':
    main()
