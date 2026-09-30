#!/usr/bin/env python3
# gen_help_names.py — spec 003 R3.4：F1 說明的譯名對照頁。
#
#   tools/py.sh tools/gen_help_names.py
#
# 從 text/story-trans.json 依訊息 ID 範圍取「原文　譯文」，寫 text/help-names.zh-TW.tsv
# （欄位 key、zh、source）。每頁最多 PER_PAGE 行（前端分兩欄排）。
# 英文欄是原文（使用者定案 2026-09-19：原文可進版控）。手冊防拷關鍵字（2500–2574）不在任何範圍內。
import json
from pathlib import Path

CATS = [  # (頁代號, 標題, 訊息 ID 起, 迄)
    ("race", "譯名對照：種族", 100, 110),
    ("class", "譯名對照：職業", 120, 133),
    ("resist", "譯名對照：抗性與屬性", 401, 413),
    ("maze", "譯名對照：迷宮選單", 2005, 2013),
    ("spell", "譯名對照：咒語", 4000, 4101),
]
PER_PAGE = 60
BANNED = range(2500, 2575)

items = json.loads(Path("text/story-trans.json").read_text())["items"]
rows = []
for code, title, a, b in CATS:
    ents = [(int(k), v) for k, v in items.items() if a <= int(k) <= b and int(k) not in BANNED]
    ents.sort()
    for page in range(0, max(len(ents), 1), PER_PAGE):
        chunk = ents[page:page + PER_PAGE]
        pid = f"names.{code}{page // PER_PAGE + 1}"
        suffix = f"（{page // PER_PAGE + 1}）" if len(ents) > PER_PAGE else ""
        rows.append((f"{pid}.00", title + suffix, f"gen_help_names.py {a}-{b}"))
        for n, (mid, v) in enumerate(chunk, 1):
            en = "".join(c for c in v["en"] if c >= " ").strip()
            zh = "".join(c for c in v["zh"] if c >= " ").strip()
            rows.append((f"{pid}.{n:02d}", f"{en}　{zh}", f"msg{mid}"))

out = ["key\tzh\tsource"] + ["\t".join(r) for r in rows]
Path("text/help-names.zh-TW.tsv").write_text("\n".join(out) + "\n")
print(f"help-names：{len(rows)} 行")
