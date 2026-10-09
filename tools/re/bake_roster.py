#!/usr/bin/env python3
"""022：十四職業與兩性別的獨立TSV家族，不讀原版訊息。"""
import argparse
import hashlib
import json
import unicodedata
from pathlib import Path

KEYS = {'F-ALC', 'F-BAR', 'F-BIS', 'F-FIG', 'F-LOR', 'F-MAG', 'F-MON', 'F-NIN', 'F-PRI', 'F-PSI', 'F-RAN', 'F-SAM', 'F-THI', 'F-VAL', 'M-ALC', 'M-BAR', 'M-BIS', 'M-FIG', 'M-LOR', 'M-MAG', 'M-MON', 'M-NIN', 'M-PRI', 'M-PSI', 'M-RAN', 'M-SAM', 'M-THI', 'M-VAL'}
CONTEXT = 'character.add.roster'
SOURCE = 'runtime:0110:0fcc'
CATALOG = Path(__file__).resolve().parents[2] / 'text/roster.zh-TW.tsv'


def build(catalog, glyph_char, size):
    raw = Path(catalog).read_bytes()
    if size not in (15, 16) or raw.startswith(b'\xef\xbb\xbf'):
        raise ValueError('職業碼字高或TSV編碼不合格')
    lines = raw.decode('utf-8').splitlines()
    if not lines or lines[0] != 'en\tzh\tcontext\tdy\tsource':
        raise ValueError('職業碼TSV表頭不合格')
    entries, seen = [], set()
    for line in lines[1:]:
        fields = line.split('\t')
        if len(fields) != 5:
            raise ValueError('職業碼TSV欄數不合格')
        en, zh, context, dy, source = fields
        if en not in KEYS or en in seen or context != CONTEXT or dy != '0' or source != SOURCE:
            raise ValueError('職業碼TSV鍵或情境不合格')
        if not zh or any(not unicodedata.name(c, '').startswith('CJK UNIFIED IDEOGRAPH-') for c in zh):
            raise ValueError('職業碼譯文須為中文字元')
        if len(zh)*(size+2)-2 > 90:
            raise ValueError('職業碼譯文超出90像素欄寬')
        glyphs = [glyph_char(c) for c in zh]
        for rows in glyphs:
            if len(rows) != size or not any(rows) or any(type(v) is not int or not 0 <= v < 1 << size for v in rows):
                raise ValueError('職業碼字模不合格')
            if size == 16 and rows[-1] != 0:
                raise ValueError('職業碼16列字模末列須空白')
        entries.append(dict(en=en, zh=zh, dy=0, g=glyphs))
        seen.add(en)
    if seen != KEYS:
        raise ValueError('職業碼TSV須有完整28鍵')
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
    import bake_font
    if args.eten_font:
        bake_font.FONT15 = args.eten_font
    if args.eten_symbol_font:
        bake_font.SPC15 = args.eten_symbol_font
    payload = build(args.catalog, bake_font.glyph_char, bake_font.HEIGHT)
    write_payload(args.output, payload)
    print(f'職業碼：28鍵，字高{payload["size"]}，缺字0')


if __name__ == '__main__':
    main()
