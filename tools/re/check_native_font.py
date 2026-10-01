#!/usr/bin/env python3
"""spec007 獨立字型驗收；僅讀取資料，不執行遊戲或烘製程式。

Docker 內執行：
  python3 tools/re/check_native_font.py --baked-only --font /font/unifont_all-17.0.05.hex.gz
  python3 tools/re/check_native_font.py --receipts workplace/r687 \
    --font /font/unifont_all-17.0.05.hex.gz

原版與字型目錄唯讀掛載；正式收據固定要求183/183/183/125/58份。
HEX直接解析，字形、來源、流文換行均不匯入受測程式。PNG解碼器沿用
check_story.py的資料解碼契約，不匯入其字型供應器或命令列。
低解析英文筆畫還原以已驗收r674的comp作回歸基線；完整PNG期望由該色號
背景、原版raw色盤、直接HEX與本工具獨立排版重建，包含一般覆繪批准字格。
本工具不驗互動Draw、F1實際翻頁、負向烘製/載入測試、封包授權或清理。
"""
import argparse
import base64
import gzip
import hashlib
import json
import re
import struct
import sys
import zlib
from collections import Counter
from functools import lru_cache
from pathlib import Path

FONT_SHA = '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
STORY = 0xBE0C
NARRATOR = 0xBCFC
SEMANTIC = ('Site', 'Fragments', 'Para', 'A', 'N', 'Lines', 'Zh')
ORDINALS = dict(zip(('FIRST SECOND THIRD FOURTH FIFTH SIXTH SEVENTH EIGHTH NINTH TENTH').split(), '一二三四五六七八九十'))
# 自製5×7數字的設計資料，非字型供應器輸出；水平/垂直各擴2倍。
DIGITS = {
    '0': (14, 17, 19, 21, 25, 17, 14), '1': (4, 12, 4, 4, 4, 4, 14),
    '2': (14, 17, 1, 2, 4, 8, 31), '3': (30, 1, 1, 14, 1, 1, 30),
    '4': (2, 6, 10, 18, 31, 2, 2), '5': (31, 16, 30, 1, 1, 17, 14),
    '6': (6, 8, 16, 30, 17, 17, 14), '7': (31, 1, 2, 4, 8, 8, 8),
    '8': (14, 17, 17, 14, 17, 17, 14), '9': (14, 17, 17, 15, 1, 2, 12),
}


class Audit:
    def __init__(self, verbose=False):
        self.total = self.failed = 0
        self.verbose = verbose
        self.general = set()

    def check(self, ok, label):
        self.total += 1
        self.failed += not ok
        if not ok or self.verbose:
            print(('通過 ' if ok else '失敗 ') + label)
        return ok


def load(path):
    return json.loads(Path(path).read_text())


def protected(i):
    return 2500 <= int(i) <= 2574


def clean(text):
    text = re.sub(r'[\x00-\x1f]', '', text)
    text = text[1:] if text.startswith(('@', '#')) else text
    return re.sub(r'\s*(?:&[%#|\[]*)+\]$', '', text)


def key(text):
    text = re.sub(r'[\x00-\x1f]', '', text)
    # 保留既有normalization的尾詞行為（如headfirst→head$），此輪只改字型。
    def ordinal_suffix(m):
        for word in ORDINALS:
            if m[0].upper().endswith(word):
                return m[0][:-len(word)] + '$'
        return m[0]
    text = re.sub(r'[A-Za-z]+', ordinal_suffix, text)
    return re.sub(r'\d+', '$', text).replace('^', '$').replace('%', '$')


class RawFont:
    def __init__(self, path):
        raw = Path(path).read_bytes()
        if hashlib.sha256(raw).hexdigest() != FONT_SHA:
            raise ValueError('固定HEX雜湊不符')
        self.raw = {}
        for line in gzip.decompress(raw).decode('ascii').splitlines():
            cp, bits = line.split(':')
            if not re.fullmatch(r'[0-9A-Fa-f]{32}|[0-9A-Fa-f]{64}', bits):
                raise ValueError('HEX字形格式不符')
            ch, width = chr(int(cp, 16)), len(bits) // 4
            if ch in self.raw:
                raise ValueError('HEX重複碼點')
            self.raw[ch] = (width, tuple(int(bits[n:n + width // 4], 16)
                                        for n in range(0, len(bits), width // 4)))
        if len(self.raw) != 128071:
            raise ValueError('固定HEX字形數不符')

    def square(self, ch):
        w, rows = self.raw[ch]
        return [v << 4 if w == 8 else v for v in rows]

    def help(self, ch):
        w, rows = self.raw[ch]
        return {'w': w, 'rows': list(rows)}


def digit(ch):
    rows = []
    for n in DIGITS[ch]:
        wide = sum(3 << (11 - 2 * x) for x in range(5) if n & (16 >> x))
        rows.extend((wide, wide))
    return [0] + rows + [0]


def source_pairs(root):
    # 在取得英文/譯文欄之前排除受保護id；工具不輸出任何來源文字。
    trans = {int(i): v for i, v in load(root / 'text/story-trans.json')['items'].items() if not protected(i)}
    pairs = [(v['en'], v['zh'], v.get('dy', 0)) for v in trans.values()]
    for name in ('monsters-trans.json', 'items-trans.json'):
        pairs.extend((en, zh, 0) for en, zh in load(root / 'text' / name)['items'].items())
    return trans, pairs


def own_xlate(font, pairs):
    result = {}
    for en, zh, dy in pairs:
        if zh.isascii() or not key(en).replace('$', '').strip():
            continue
        glyphs, slots = [], []
        for ch in zh:
            if ch == '$':
                slots.append(len(glyphs))
                glyphs.append([0] * 16)
            elif len(ch.encode('big5')) == 2:
                glyphs.append(font.square(ch))
        v = {'dy': dy, 'g': glyphs}
        if slots:
            v['d'] = slots
        result[key(en)] = v
    result['__digits__'] = {ch: digit(ch) for ch in DIGITS}
    result['__ordinals__'] = {en: font.square(zh) for en, zh in ORDINALS.items()}
    return result


def usable(text):
    try:
        return bool(text) and all(len(ch.encode('big5')) == 2 for ch in text)
    except UnicodeEncodeError:
        return False


def own_flow(font, messages, trans):
    paragraphs, previous, current = [], None, None
    for i, en in sorted(messages.items()):
        raw = re.sub(r'[\x00-\x1f]', '', en)
        lead, body = raw.startswith('@'), raw.removeprefix('@')
        if any(ch in '&%]^$#!"@' or ch.isdigit() for ch in body):
            current, previous = None, i
            continue
        if current is None or previous != i - 1 or lead:
            current = {'words': [], 'lines': []}
            paragraphs.append(current)
        previous = i
        if body.split():
            current['words'].extend(body.split())
            z = trans.get(i, {}).get('zh', '')
            current['lines'].append({'id': i, 'end': len(current['words']), 'zh': z if usable(z) else ''})
    paragraphs = [{'id': p['lines'][0]['id'], **p} for p in paragraphs if p['lines']]
    msgs = {}
    for i, en in messages.items():
        body, z = clean(en), clean(trans.get(i, {}).get('zh', ''))
        if body.split() and all(32 <= ord(c) <= 126 and c not in '@^$%&]#' for c in body):
            msgs[str(i)] = {'w': ' '.join(body.split()), 't': int(body != re.sub(r'[\x00-\x1f]', '', en).lstrip('@#')),
                            'z': z if usable(z) else ''}
    used = {ch for p in paragraphs for ln in p['lines'] for ch in ln['zh']}
    used |= {ch for m in msgs.values() for ch in m['z']}
    return {'size': 16, 'paras': paragraphs, 'msgs': msgs,
            'glyphs': {ch: font.square(ch) for ch in used}}


def own_help(root, font):
    groups, chars = {}, set()
    for name in ('help.zh-TW.tsv', 'help-names.zh-TW.tsv'):
        lines = (root / 'text' / name).read_text().splitlines()
        if lines[0] != 'key\tzh\tsource':
            raise ValueError('說明TSV表頭不符')
        for line in lines[1:]:
            k, zh, _ = line.split('\t')
            page, seq = k.rsplit('.', 1)
            if seq in groups.setdefault(page, {}):
                raise ValueError('說明TSV重複鍵')
            groups[page][seq] = zh
            chars.update(zh)
    pages = [{'id': p, 'title': s['00'], 'lines': [s[k] for k in sorted(s) if k != '00']}
             for p, s in groups.items()]
    return {'pages': pages, 'glyphs': {ch: font.help(ch) for ch in chars}}


def verify_baked(args, audit, font, trans, pairs, messages):
    root, directory, baseline = args.root, args.baked, args.eten
    expected = {'xlate': own_xlate(font, pairs), 'flow': own_flow(font, messages, trans),
                'help': own_help(root, font)}
    for name, want in expected.items():
        got = load(directory / (name + '.json'))
        audit.check(got == want, f'{name}.json全部來源／字形／槽位與獨立HEX期望完全相同')
        if name == 'flow':
            old = load(baseline / 'flow.json')
            audit.check(old['size'] == 15 and got['size'] == 16, '字型尺寸15／16明示')
            audit.check(all(got[k] == old[k] for k in ('paras', 'msgs')), '流文資料譯文／英文鍵／來源欄與倚天相同')
            historic = load(args.historic / 'flow.json')
            historic.pop('size', None)
            old.pop('size', None)
            audit.check(old == historic, '倚天流文重建與r674烘製基線逐位元相同')
        elif name == 'help':
            audit.check(got['pages'] == load(baseline / 'help.json')['pages'], 'F1邏輯來源頁文字未變')
            audit.check(load(baseline / 'help.json') == load(args.historic / 'help.json'), '倚天F1字形與既有烘製相同')
        else:
            audit.check(load(baseline / 'xlate.json') == load(args.historic / 'xlate.json'), '倚天一般覆繪與既有烘製逐位元相同')
            audit.check(got.get('__digits__') == want['__digits__'], '10個自製數字均在16字格x3..12/y1..14居中')
            audit.check(got.get('__ordinals__') == want['__ordinals__'], '10個序數中文字形逐位元採固定HEX')
    chars = {c for _, zh, _ in pairs for c in zh if ord(c) > 127} | set(expected['help']['glyphs'])
    audit.check(all(c in font.raw for c in chars), f'正式譯文及F1所需 {len(chars)} 字全部存在固定HEX')
    audit.check(all(expected['help']['glyphs'][str(n)] == font.help(str(n)) for n in range(10)), 'F1頁尾ASCII數字採原生8×16')
    print('烘製範圍：一般覆繪 %d 串、流文 %d 字形、F1 %d 字形（含ASCII）、自製數字10槽' %
          (len(expected['xlate']) - 2, len(expected['flow']['glyphs']), len(expected['help']['glyphs'])))
    return expected


def read_png(path):
    return decode_png(Path(path).read_bytes())


@lru_cache(maxsize=24)
def decode_png(data):
    """獨立解碼8位RGB/RGBA；快取相同原始PNG，不共享任何排版/字型程式。"""
    if data[:8] != b'\x89PNG\r\n\x1a\n':
        raise ValueError('PNG簽名不符')
    pos, chunks, header = 8, [], None
    while pos < len(data):
        n = struct.unpack('>I', data[pos:pos + 4])[0]
        kind, block = data[pos + 4:pos + 8], data[pos + 8:pos + 8 + n]
        crc = struct.unpack('>I', data[pos + 8 + n:pos + 12 + n])[0]
        if zlib.crc32(kind + block) & 0xffffffff != crc:
            raise ValueError('PNG區塊CRC不符')
        if kind == b'IHDR':
            header = struct.unpack('>IIBBBBB', block)
        elif kind == b'IDAT':
            chunks.append(block)
        elif kind == b'IEND':
            break
        pos += n + 12
    w, h, depth, color, compression, filtering, interlace = header
    if depth != 8 or color not in (2, 6) or (compression, filtering, interlace) != (0, 0, 0):
        raise ValueError('PNG格式不支援')
    bpp = 3 if color == 2 else 4
    raw, stride = zlib.decompress(b''.join(chunks)), w * bpp
    if len(raw) != h * (stride + 1):
        raise ValueError('PNG解壓長度不符')
    previous, rows = bytearray(stride), []
    for y in range(h):
        start = y * (stride + 1)
        ft, line = raw[start], bytearray(raw[start + 1:start + 1 + stride])
        if ft > 4:
            raise ValueError('PNG濾波格式不符')
        if ft == 2:
            line = bytearray((a + b) & 255 for a, b in zip(line, previous))
        elif ft:
            for i in range(stride):
                a, b, c = line[i - bpp] if i >= bpp else 0, previous[i], previous[i - bpp] if i >= bpp else 0
                if ft == 1:
                    v = a
                elif ft == 2:
                    v = b
                elif ft == 3:
                    v = (a + b) // 2
                else:
                    p = a + b - c
                    ds = abs(p - a), abs(p - b), abs(p - c)
                    v = (a, b, c)[ds.index(min(ds))]
                line[i] = (line[i] + v) & 255
        previous = line
        if bpp == 4:
            if any(v != 255 for v in line[3::4]):
                raise ValueError('原版PNG必須不透明')
            line = bytearray(b''.join(line[x:x + 3] for x in range(0, stride, 4)))
        rows.append(bytes(line))
    return w, h, b''.join(rows)


def dump(directory, name):
    base = directory / name
    m, a, c = load(base.with_suffix('.json')), base.with_suffix('.a000').read_bytes(), base.with_suffix('.comp').read_bytes()
    if len(a) != 64000 or len(c) != 64000:
        raise ValueError('色號檔長度不符')
    return {'meta': m, 'a': a, 'c': c, 'png': base.with_suffix('.png')}


def own_rows(text, columns):
    # 依spec006標點契約切成不可分的雙省略號/破折號，貪婪裝入後最多回退3單位。
    units = re.findall(r'……|——|.', text)
    rows = []
    while units:
        end, used = 0, 0
        for u in units:
            if used + len(u) > columns:
                break
            used += len(u)
            end += 1
        end = max(1, end)
        for _ in range(3):
            if 1 < end < len(units) and (units[end][0] in '，。、！？；：）」』》…—' or units[end - 1][0] in '（「『《'):
                end -= 1
            else:
                break
        rows.append(''.join(units[:end]))
        units = units[end:]
    return rows


def positions(fl, scale):
    gs = 2 if scale == 3 else 1
    width = 296 if fl['Site'] == STORY else 184
    text = own_rows(fl['Zh'], max(1, width * scale // (18 * gs)))
    if len(text) > 4:
        raise ValueError('流文超過4列')
    return [{'Text': t, 'X0': 160 * scale - len(t) * 18 * gs // 2,
             'Y0': (146 + 12 * i) * scale} for i, t in enumerate(text)]


def rect_points(x, y, w, h):
    return {yy * 320 + xx for yy in range(max(0, y), min(200, y + h))
            for xx in range(max(0, x), min(320, x + w))}


def fnv(raw):
    h = 2166136261
    for v in raw:
        h = ((h ^ v) * 16777619) & 0xffffffff
    return h


def palette_from_original(en, raw_palette_path):
    w, h, rgb = read_png(en['png'])
    if (w, h) != (960, 600):
        raise ValueError('r674/en必須是3×畫布')
    excluded = set()
    for e in en['meta'].get('Shown') or []:
        if e.get('G'):
            excluded |= rect_points(e['X0'], e['Y0'] + e['Dy'], (len(e['G']) * 17 + 2) // 3 + 1, 6)
    raw_palette = Path(raw_palette_path).read_bytes()
    if len(raw_palette) != 768 or fnv(raw_palette) != en['meta']['Pal']:
        raise ValueError('原版RGB色盤768bytes及FNV收據不符')
    palette = {i: raw_palette[3 * i:3 * i + 3] for i in range(256)}
    for i, index in enumerate(en['c']):
        if i in excluded:
            continue
        p = ((i // 320) * 3 * w + (i % 320) * 3) * 3
        c = rgb[p:p + 3]
        if palette[index] != c:
            raise ValueError('原版可見色號RGB不唯一')
    return palette


def background(comp, palette, scale):
    rows = []
    for y in range(200):
        row = b''.join(palette[c] * scale for c in comp[y * 320:(y + 1) * 320])
        rows.extend([row] * scale)
    return bytearray(b''.join(rows))


@lru_cache(maxsize=None)
def bits(rows):
    width = len(rows)
    return tuple((x, y) for y, row in enumerate(rows) for x in range(width) if row & (1 << (width - 1 - x)))


def paint(image, w, h, x0, y0, rows, fg, gs=1):
    for x, y in bits(tuple(rows)):
        for dy in range(gs):
            yy = y0 + y * gs + dy
            if not 0 <= yy < h:
                raise ValueError('字形高度越出畫布')
            for dx in range(gs):
                xx = x0 + x * gs + dx
                if not 0 <= xx < w:
                    raise ValueError('字形寬度越出畫布')
                off = (yy * w + xx) * 3
                image[off:off + 3] = fg


def source_flow(fl, paragraphs, messages, trans, story_pages):
    fragments = fl['Fragments']
    words = ' '.join(fragments).split()
    if fl['Site'] == STORY:
        found = [p for p in story_pages if p['fragments'] == fragments]
        if len(found) != 1:
            raise ValueError('故事不是原版完整頁')
        ids = found[0]['ids']
        want = {'Site': STORY, 'Fragments': fragments, 'Para': ids[0], 'A': 0,
                'N': len(words), 'Lines': ids, 'Zh': ''.join(clean(trans[i]['zh']) for i in ids)}
        return want
    if fl['Site'] != NARRATOR:
        raise ValueError('非故事/旁白流文')
    matches = []
    for p in paragraphs:
        if not 15000 <= p['id'] < 15100:
            continue
        for start in range(len(p['words']) - len(words) + 1):
            if p['words'][start:start + len(words)] == words:
                done = [ln for ln in p['lines'] if start < ln['end'] <= start + len(words)]
                if not done or any(not ln['zh'] for ln in done):
                    continue
                matches.append({'Site': NARRATOR, 'Fragments': fragments, 'Para': p['id'], 'A': start,
                                'N': len(words), 'Lines': [ln['id'] for ln in done], 'Zh': ''.join(ln['zh'] for ln in done)})
    if len(matches) != 1:
        raise ValueError('旁白英文來源不唯一或不完整')
    return matches[0]


def story_sources(args, messages, audit):
    raw = (args.root / 'workplace/bake/r667-story-msgs.txt').read_text()
    entries = [(int(s), json.loads(t), int(y)) for s, t, y in re.findall(
        r'印字 (\d+) ("(?:[^"\\]|\\.)*") font=2 ret=0110:be0c .*?xy0=\(\d+,(\d+)\)', raw)]
    reads = [(int(s), int(i)) for s, i in re.findall(r'訊息 (\d+) id=(\d+).*ret=bfeb', raw)]
    traces = load(args.root / 'workplace/r673/timing/traces.json')
    if len(entries) != 71 or len(reads) != 68 or len(traces) != 71:
        raise ValueError('故事原版觀測不完整')
    pages = []
    for (step, text, y), t in zip(entries, traces):
        if y == 146:
            pages.append({'fragments': [], 'entry': step})
        p = pages[-1]
        if t['Page'] != len(pages) or t['Y'] != y or step >= t['Step']:
            raise ValueError('故事真正返回觀測不符')
        p['fragments'].append(text)
        p.update(entry=step, returned=t['Step'])
    if len(pages) != 31:
        raise ValueError('故事必須31頁')
    lower = 0
    for n, p in enumerate(pages, 1):
        p['ids'] = [i for step, i in reads if lower < step < p['entry'] and 31500 <= i <= 31899
                    and i in messages and clean(messages[i]).split()]
        audit.check([w for i in p['ids'] for w in clean(messages[i]).split()] == ' '.join(p['fragments']).split(),
                    f'故事第{n}頁以前頁真正返回為下界，原版英文與來源完整相等')
        lower = p['returned']
    return pages


def inspect(audit, name, new, old, en, scale, font, xlate, paragraphs, messages, trans, story_pages):
    m, b = new['meta'], old['meta']
    audit.check(all(k in m and k in b and m[k] == b[k] for k in ('Step', 'Mem', 'A000', 'Pal', 'Mouse'))
                and new['a'] == old['a'] == en['a'], f'{name}全部原版狀態及raw A000與r674同等樣本一致')
    audit.check(new['c'] == old['c'], f'{name}低解析英文筆畫還原與倚天基線相同')
    audit.check(fnv(new['a']) == m['A000'], f'{name}raw A000的FNV與收據欄相符')
    fs, bs = m.get('Flow') or [], b.get('Flow') or []
    audit.check([{k: f[k] for k in SEMANTIC} for f in fs] == [{k: f[k] for k in SEMANTIC} for f in bs],
                f'{name}完整流文語意與原有淡入/過渡/回退閘門相同')
    es, be = m.get('Shown') or [], b.get('Shown') or []
    audit.check([{k: v for k, v in e.items() if k != 'G'} for e in es] ==
                [{k: v for k, v in e.items() if k != 'G'} for e in be], f'{name}一般印字事件、原始離屏快照與基線相同')
    palette = en['palette']
    image = background(old['c'], palette, scale)
    w, h = 320 * scale, 200 * scale
    for e in es:
        if not e.get('G'):
            continue
        entry = xlate.get(key(e['S']))
        # 新增動態槽案例須擴充獨立oracle；不拿runtime G充當期望。
        if entry is None or entry.get('d'):
            raise ValueError('一般覆繪尚無獨立完整槽位來源')
        glyphs = entry['g']
        audit.check(e['G'] == glyphs and e['Dy'] == entry['dy'], f'{name}一般覆繪逐字採來源HEX完整16列')
        fg = palette[e['Fg']]
        x0, y0 = e['X0'] * scale, (e['Y0'] + entry['dy']) * scale
        audit.general.add((scale, x0, y0, len(glyphs) * 18, 16))
        for i, g in enumerate(glyphs):
            paint(image, w, h, x0 + i * 18, y0, g, fg)
    seen = set()
    for f in fs:
        want = source_flow(f, paragraphs, messages, trans, story_pages)
        audit.check(all(f[k] == want[k] for k in SEMANTIC), f'{name}流文英文比對、完成來源id及正式譯文獨立相等')
        rs = positions(want, scale)
        gs = 2 if scale == 3 else 1
        audit.check(f['Rows'] == rs and f['GS'] == gs and f['Scale'] == scale, f'{name}16px字格、18pitch換行／座標與倍率相同')
        fg = palette[0x55]
        audit.check(bytes(f['FgRGB']) == fg, f'{name}流文字色由原版0x55獨立導出')
        for row in rs:
            for i, ch in enumerate(row['Text']):
                paint(image, w, h, row['X0'] + i * 18 * gs, row['Y0'], font.square(ch), fg, gs)
        seen.add((f['Site'], tuple(f['Fragments'])))
    pw, ph, pixels = read_png(new['png'])
    audit.check((pw, ph) == (w, h) and pixels == image,
                f'{name}整張PNG等於倚天還原色號背景＋原版色盤＋獨立字形（含空白及批准區域之外）')
    return seen


def verify_receipts(args, audit, font, expected, messages, trans):
    pages = story_sources(args, messages, audit)
    story = {f"hdat-{p['entry'] + d * 1000000}" for p in pages for d in (1, 3, 6, 8)} | {'hdfrag-4'}
    narr = {f'hdat-{n}' for n in range(1483000000, 1538000001, 1000000)} | {'hdfrag-81', 'hdfrag-82'}
    modes = {'en': (story | narr, 3), 'zh': (story | narr, 3), 'zh2': (story | narr, 2),
             'drop': (story, 3), 'narr-drop': (narr, 3)}
    # 旁白20片段由原版印字觀測分成6頁，不以中文Flow輸出推算頁數。
    narrator_entries = [(int(s), int(y), json.loads(t)) for s, y, t in re.findall(
        r'流文 片段 (\d+) ret=0110:bcfc font=2 y0=(\d+) x0=\d+ ("(?:[^"\\]|\\.)*")',
        (args.baseline / 'en.log').read_text())]
    narrator_returns = [int(s) for s in re.findall(r'流文 頁 (\d+) 片段=',
                         (args.baseline / 'en.log').read_text())
                        if narrator_entries and int(s) >= narrator_entries[0][0]]
    if len(narrator_entries) != 20 or len(narrator_returns) != 20:
        raise ValueError('旁白原版入口／真正返回觀測不完整')
    narrator_pages = []
    for index, ((step, y, text), returned) in enumerate(zip(narrator_entries, narrator_returns)):
        if not step < returned or (index + 1 < len(narrator_entries) and
                                   returned >= narrator_entries[index + 1][0]):
            raise ValueError('旁白原版入口／真正返回配對順序不符')
        if y == 146:
            narrator_pages.append({'fragments': [], 'first_entry': step})
        p = narrator_pages[-1]
        p['fragments'].append(text)
        p['entry'] = step
        p['returned'] = returned
    if len(narrator_entries) != 20 or len(narrator_pages) != 6:
        raise ValueError('旁白原版片段／完整頁來源不足')
    for index, p in enumerate(narrator_pages):
        # 第6頁終窗開放；僅量既有截至1538M的請求，不外推下一屏。
        p['next_entry'] = narrator_pages[index + 1]['first_entry'] if index < 5 else None
    narr_keys = {(NARRATOR, tuple(p['fragments'])) for p in narrator_pages}
    narr_drop_keys = set()
    for p in narrator_pages:
        want = source_flow({'Site': NARRATOR, 'Fragments': p['fragments']}, expected['flow']['paras'], messages, trans, pages)
        if want['Para'] != 15005:
            narr_drop_keys.add((NARRATOR, tuple(p['fragments'])))
    audit.check([len(v[0]) for v in modes.values()] == [183, 183, 183, 125, 58], '固定必要樣本183/183/183/125/58')
    for mode, (names, _) in modes.items():
        for directory in (args.receipts / mode, args.baseline / mode):
            present = {p.stem for p in directory.glob('hd*.json')}
            if not names <= present:
                raise ValueError(f'{mode}必要同狀態收據不足')
            for name in names:
                suffixes = ('.json', '.a000', '.comp', '.png', '.pal') if directory.parent == args.receipts else ('.json', '.a000', '.comp', '.png')
                if any(not (directory / (name + suffix)).is_file() for suffix in suffixes):
                    raise ValueError(f'{mode}必要收據配套檔缺漏')
    seen = {m: set() for m in modes}
    window_seen = {m: set() for m in ('zh', 'zh2', 'narr-drop')}
    for counter, name in enumerate(sorted(story | narr), 1):
        en = dump(args.baseline / 'en', name)
        en['palette'] = palette_from_original(en, args.receipts / 'en' / (name + '.pal'))
        # 頁窗口只來自原版log與實際Step；不以請求檔名或中文Flow推算。
        window_page = next((i for i, p in enumerate(narrator_pages)
                            if name in narr and en['meta']['Step'] >= p['returned'] and
                            (p['next_entry'] is None or en['meta']['Step'] < p['next_entry'])), None)
        original_stable = False
        if window_page is not None:
            p = narrator_pages[window_page]
            original = dump(args.baseline / 'zh', name)
            events = [e for e in original['meta'].get('Shown') or [] if e['S'] in p['fragments']]
            base_bad = stroke = present = 0
            mx, my = original['meta']['Mouse']
            for e in events:
                before, after = base64.b64decode(e['Before']), base64.b64decode(e['After'])
                if len(before) != e['BW'] * e['BH'] or len(after) != len(before):
                    raise ValueError('原版旁白離屏快照長度不符')
                for pos, (b, a) in enumerate(zip(before, after)):
                    x, y = e['BX'] + pos % e['BW'], e['BY'] + pos // e['BW']
                    if mx <= x < mx + 16 and my <= y < my + 16:
                        continue
                    value = en['a'][y * 320 + x]
                    if b == a:
                        base_bad += value != a
                    else:
                        stroke += 1
                        present += value != b
            # 返回後仍可能尚未blit；只有完整原始快照和已證實的背景／筆畫
            # 條件成立才強制覆繪。窗口內其餘樣本仍逐份驗全部狀態和整張PNG。
            original_stable = (sorted(e['S'] for e in events) == sorted(p['fragments']) and
                               en['meta']['Step'] >= p['returned'] + 200000 and
                               base_bad == 0 and stroke > 0 and present * 2 >= stroke)
        for mode, (names, scale) in modes.items():
            if name not in names:
                continue
            old, new = dump(args.baseline / mode, name), dump(args.receipts / mode, name)
            pal_raw = (args.receipts / mode / (name + '.pal')).read_bytes()
            audit.check(len(pal_raw) == 768 and fnv(pal_raw) == new['meta']['Pal'] and
                        pal_raw == (args.receipts / 'en' / (name + '.pal')).read_bytes(),
                        f'{mode}/{name}全部原版RGB色盤bytes與同狀態en相同且FNV符合')
            actual = inspect(audit, mode + '/' + name, new, old, en, scale, font, expected['xlate'],
                             expected['flow']['paras'], messages, trans, pages)
            seen[mode] |= actual
            for idx, p in enumerate(pages):
                if name in {f"hdat-{p['entry'] + d * 1000000}" for d in (6, 8)} and mode in ('zh', 'zh2', 'drop'):
                    need = 0 if mode == 'drop' and idx == 2 else 1
                    audit.check(sum(f['Site'] == STORY and f['Fragments'] == p['fragments'] for f in new['meta'].get('Flow') or []) == need,
                                f'{mode}/{name}第{idx + 1}頁6/8M完整穩定覆繪強制存在（回退頁除外）')
            if mode == 'drop' and name in {f"hdat-{pages[2]['entry'] + d * 1000000}" for d in (1, 3, 6, 8)}:
                audit.check(not new['meta'].get('Flow'), f'{mode}/{name}缺31654全頁英文回退')
            if mode in ('zh', 'zh2', 'narr-drop'):
                if window_page is not None:
                    p = narrator_pages[window_page]
                    k = NARRATOR, tuple(p['fragments'])
                    fs = [f for f in new['meta'].get('Flow') or [] if f['Site'] == NARRATOR]
                    audit.check(all(f['Fragments'] == p['fragments'] for f in fs),
                                f'{mode}/{name}實際Step位於原版第{window_page + 1}頁最後返回／下一入口窗口，所見頁完整相符')
                    if original_stable:
                        need = 0 if mode == 'narr-drop' and k not in narr_drop_keys else 1
                        audit.check(len(fs) == need,
                                    f'{mode}/{name}原版完整快照及背景／筆畫穩定時，旁白第{window_page + 1}頁強制覆蓋（回退頁除外）')
                        if len(fs) == 1 and k in actual:
                            window_seen[mode].add(k)
        if counter % 12 == 0:
            print(f'收據進度：{counter}/183組', flush=True)
    story_keys = {(STORY, tuple(p['fragments'])) for p in pages}
    audit.check({k for k in seen['zh'] if k[0] == STORY} == story_keys and
                {k for k in seen['zh2'] if k[0] == STORY} == story_keys, '故事31完整頁兩倍率皆覆蓋')
    for mode in ('zh', 'zh2'):
        audit.check({k for k in seen[mode] if k[0] == NARRATOR} == narr_keys, f'{mode}旁白6完整頁皆穩定覆蓋')
        audit.check(window_seen[mode] == narr_keys, f'{mode}旁白6頁各至少一份完整覆蓋落在原版返回／下一入口窗口且原版背景筆畫穩定')
    audit.check({k for k in seen['drop'] if k[0] == STORY} == story_keys - {(STORY, tuple(pages[2]['fragments']))}, '故事缺訊息僅影響第3頁，其他30頁全覆蓋')
    audit.check(len(narr_drop_keys) == 5 and {k for k in seen['narr-drop'] if k[0] == NARRATOR} == narr_drop_keys and
                not any(f.get('Para') == 15005 for name in narr for f in load(args.receipts / 'narr-drop' / (name + '.json')).get('Flow') or []),
                 '旁白缺15005僅該頁回退，其他5頁全覆蓋')
    audit.check(window_seen['narr-drop'] == narr_drop_keys, '旁白缺15005其他5頁均有原版時窗內完整穩定覆蓋')
    print('一般批准字格（倍率,x,y,寬,高）：' + repr(sorted(audit.general)))


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument('--root', type=Path, default=Path('.'))
    ap.add_argument('--font', type=Path, default=Path('/font/unifont_all-17.0.05.hex.gz'))
    ap.add_argument('--baked', type=Path, default=Path('workplace/r682/unifont'))
    ap.add_argument('--eten', type=Path, default=Path('workplace/r682/eten'))
    ap.add_argument('--historic', type=Path, default=Path('workplace/bake'))
    ap.add_argument('--receipts', type=Path, default=Path('workplace/r687'))
    ap.add_argument('--baseline', type=Path, default=Path('workplace/r674'))
    ap.add_argument('--baked-only', action='store_true', help='僅驗烘製資料，明示不算同狀態收據通過')
    ap.add_argument('--verbose', action='store_true')
    args = ap.parse_args()
    audit = Audit(args.verbose)
    try:
        font = RawFont(args.font)
        trans, pairs = source_pairs(args.root)
        messages = {r['id']: r['en'] for r in load(args.root / 'text/messages.json')['records'] if not protected(r['id'])}
        expected = verify_baked(args, audit, font, trans, pairs, messages)
        if not args.baked_only:
            verify_receipts(args, audit, font, expected, messages, trans)
        for source in (args.font, args.root / 'text/story-trans.json', args.root / 'text/help.zh-TW.tsv',
                       args.root / 'text/help-names.zh-TW.tsv', args.root / 'text/monsters-trans.json', args.root / 'text/items-trans.json',
                       args.root / 'text/messages.json', *(args.baked / (n + '.json') for n in ('xlate', 'flow', 'help'))):
            print(f'來源 {source} SHA-256={hashlib.sha256(source.read_bytes()).hexdigest()}')
    except (OSError, ValueError, KeyError, IndexError, TypeError, struct.error, zlib.error) as exc:
        audit.check(False, f'驗收輸入/格式受阻：{type(exc).__name__}（不輸出來源文字）')
    print(f'原生字型驗收：{audit.total}項，失敗{audit.failed}項；範圍：' + ('僅烘製資料' if args.baked_only else '烘製＋全部固定同狀態收據'))
    print('未量限制：Mem僅比較runner雜湊；互動Draw、F1實際翻頁/PNG、負向烘製與載入分支、封包、全遊戲；不據此宣稱CONFORMED。')
    return 1 if audit.failed else 0


if __name__ == '__main__':
    sys.exit(main())
