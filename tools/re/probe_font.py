#!/usr/bin/env python3
"""只讀字型／現有譯文的尺寸探針；不產生正式烘製資料。須在 Docker 執行。"""
import argparse
import gzip
import hashlib
import json
from collections import Counter
from pathlib import Path

ap = argparse.ArgumentParser()
ap.add_argument("hex", type=Path)
ap.add_argument("out", type=Path)
ap.add_argument("--charset", type=Path, help="另存正式譯文與F1所需碼點清單，不含字模")
args = ap.parse_args()
if args.out.exists():
    raise SystemExit("輸出已存在，禁止覆寫研究收據")
if args.charset and args.charset.exists():
    raise SystemExit("字元清單已存在，禁止覆寫")
raw = args.hex.read_bytes()
glyphs = {}
for line in gzip.decompress(raw).decode("ascii").splitlines():
    cp, bits = line.split(":")
    if len(bits) not in (32, 64):
        raise SystemExit("HEX 字形尺寸不符")
    ch = chr(int(cp, 16))
    if ch in glyphs:
        raise SystemExit("HEX 碼點重複")
    glyphs[ch] = bits

chars = set()
ascii_help = set("0123456789")
for name in ("story-trans.json", "items-trans.json", "monsters-trans.json"):
    for key, value in json.loads(Path("text", name).read_text())["items"].items():
        if name == "story-trans.json" and 2500 <= int(key) <= 2574:
            continue
        zh = value["zh"] if isinstance(value, dict) else value
        chars.update(c for c in zh if ord(c) > 127)
for name in ("help.zh-TW.tsv", "help-names.zh-TW.tsv"):
    for line in Path("text", name).read_text().splitlines()[1:]:
        text = line.split("\t")[1]
        chars.update(c for c in text if ord(c) > 127)
        ascii_help.update(c for c in text if 32 <= ord(c) <= 126)

boxes, over = Counter(), []
for ch in sorted(chars & glyphs.keys()):
    bits = glyphs[ch]
    width = len(bits) // 4
    rows = [int(bits[i:i + width // 4], 16)
            for i in range(0, len(bits), width // 4)]
    points = [(x, y) for y, row in enumerate(rows) for x in range(width)
              if row & (1 << (width - x - 1))]
    if not points:
        continue
    x0, x1 = min(x for x, _ in points), max(x for x, _ in points)
    y0, y1 = min(y for _, y in points), max(y for _, y in points)
    boxes[x1 - x0 + 1, y1 - y0 + 1] += 1
    if x1 - x0 + 1 > 15 or y1 - y0 + 1 > 15:
        over.append([ch, x0, x1, y0, y1])
report = {"font": args.hex.name, "sha256": hashlib.sha256(raw).hexdigest(),
          "glyphs": len(glyphs), "required": len(chars),
          "missing": sorted(chars - glyphs.keys()), "over15": over,
          "sizes": {str(k): v for k, v in boxes.items()}}
args.out.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
if args.charset:
    required = chars | ascii_help | set("一二三四五六七八九十")
    if required - glyphs.keys():
        raise SystemExit("正式字元清單含缺字，不產生清單")
    args.charset.write_text("".join(f"U+{ord(ch):04X}\n" for ch in sorted(required)), encoding="utf-8")
print(f"所需字元 {len(chars)}，缺字 {len(report['missing'])}，超過15像素 {len(over)}")
