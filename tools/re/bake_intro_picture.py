#!/usr/bin/env python3
"""025：烘製片頭呈獻獨立家族，不讀原版訊息。"""
import argparse
import hashlib
import json
from pathlib import Path

CATALOG = Path(__file__).resolve().parents[2] / 'text/intro-picture.zh-TW.tsv'


def build(catalog, glyph_char, size):
    raw = Path(catalog).read_bytes()
    if size not in (15, 16) or raw.startswith(b'\xef\xbb\xbf'):
        raise ValueError('片頭字高或TSV編碼不合格')
    text = raw.decode('utf-8').replace('\r\n', '\n')
    if '\r' in text or any((ord(c) < 32 and c not in '\t\n') or 127 <= ord(c) <= 159 for c in text):
        raise ValueError('片頭TSV含禁用控制字元')
    lines = text.split('\n')
    if lines[-1] == '':
        lines.pop()
    if len(lines) != 2 or lines[0] != 'en\tzh\tcontext\tdy\tsource':
        raise ValueError('片頭TSV須有表頭及唯一鍵')
    fields = lines[1].split('\t')
    if fields != ['presents', '呈獻', 'intro.mon60.presents', '0', 'runtime:0110:5a9d']:
        raise ValueError('片頭TSV鍵、譯文或情境不合格')
    glyphs = [glyph_char(ch) for ch in fields[1]]
    if any(len(g) != size or not any(g) or any(type(v) is not int or not 0 <= v < 1 << size for v in g) for g in glyphs):
        raise ValueError('片頭原生字模不合格')
    return dict(schema=1, size=size, context=fields[2], source=fields[4],
                source_sha256=hashlib.sha256(raw).hexdigest(),
                entries=[dict(en=fields[0], zh=fields[1], dy=0, g=glyphs)])


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
    print(f'片頭：一鍵兩字，字高{payload["size"]}，缺字0')


if __name__ == '__main__':
    main()
