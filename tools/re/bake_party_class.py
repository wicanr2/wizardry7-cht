#!/usr/bin/env python3
"""036：從唯一正式TSV加入側欄短碼，其他xlate值及附屬檔保持。"""
import argparse
import hashlib
import json
import unicodedata
from pathlib import Path

CATALOG = Path(__file__).resolve().parents[2] / 'text/party-class.zh-TW.tsv'
TEXT_SHA = '07041a64c2f07c52656e3f569db455b27a20a62882ff85fc62ece58e20ea8b0c'


def build(catalog, glyph_char, size):
    raw = Path(catalog).read_bytes()
    if size not in (15, 16) or raw.startswith(b'\xef\xbb\xbf'):
        raise ValueError('側欄職業字高或TSV編碼不合格')
    lines = raw.decode('utf-8').splitlines()
    if len(lines) != 2 or lines[0] != 'en\tzh\tcontext\tdy\tsource':
        raise ValueError('側欄職業TSV須為唯一正式鍵')
    fields = lines[1].split('\t')
    if len(fields) != 5:
        raise ValueError('側欄職業TSV欄數不合格')
    en, zh, context, dy, source = fields
    if (en, context, dy, source) != ('ALCHEMI', 'party.class.short', '0', 'runtime:0110:1c71'):
        raise ValueError('側欄職業鍵或情境不合格')
    if len(zh) != 4 or any(not unicodedata.name(ch, '').startswith('CJK UNIFIED IDEOGRAPH-') for ch in zh):
        raise ValueError('側欄職業譯文須為四個漢字')
    canonical = en.encode()+b'\0'+zh.encode()+b'\0'
    if hashlib.sha256(canonical).hexdigest() != TEXT_SHA:
        raise ValueError('側欄職業文案內容指紋不符')
    glyphs = [glyph_char(ch) for ch in zh]
    for glyph in glyphs:
        if len(glyph) != size or not any(glyph) or any(type(row) is not int or not 0 <= row < 1 << size for row in glyph):
            raise ValueError('側欄職業字模不合格')
    return {'dy': 0, 'g': glyphs}


def clone_xlate(source, output, short, size):
    data = json.loads(Path(source).read_bytes())
    if not isinstance(data, dict) or 'ALCHEMI' in data:
        raise ValueError('來源須為尚無側欄短碼的xlate物件')
    for key, value in data.items():
        if not key.startswith('__') and isinstance(value, dict):
            for glyph in value.get('g', []):
                if len(glyph) != size:
                    raise ValueError('原xlate字高與側欄後端不一致')
    result = dict(data)
    result['ALCHEMI'] = short
    if {key: value for key, value in result.items() if key != 'ALCHEMI'} != data:
        raise ValueError('其他xlate資料不得改變')
    with Path(output).open('x', encoding='utf-8') as stream:
        json.dump(result, stream, ensure_ascii=False, separators=(',', ':'))
        stream.write('\n')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('source', type=Path)
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
    short = build(args.catalog, bake_font.glyph_char, bake_font.HEIGHT)
    clone_xlate(args.source, args.output, short, bake_font.HEIGHT)
    print(json.dumps({'key': 'ALCHEMI', 'glyphs': 4, 'size': bake_font.HEIGHT,
                      'source_sha256': hashlib.sha256(args.catalog.read_bytes()).hexdigest(),
                      'xlate_sha256': hashlib.sha256(args.output.read_bytes()).hexdigest()}))


if __name__ == '__main__':
    main()
