#!/usr/bin/env python3
# trans_batch.py — 翻譯切批（batch-subagent-localization 工作流第 1 步）
#
# 讀 workplace/todo_ids.json（trans_audit.py 產出），切 ~130 則/批，
# 優先在 id 千位桶邊界切開（同批主題相近，術語一致）。
# 輸出 workplace/batches/batch-NNN.json：[{id, blk, en}]。
# subagent 譯完寫同目錄 batch-NNN.out.json：{id: zh}。
import json
import sys
from pathlib import Path

BATCH = 130
OUT = Path('workplace/batches')

ids = json.load(open('workplace/todo_ids.json'))
msg = {m['id']: m for m in json.load(open('text/messages.json'))['records']}

# 切批：桶邊界優先
batches, cur = [], []
for i in sorted(ids):
    if cur and (len(cur) >= BATCH or i // 1000 != cur[-1] // 1000):
        batches.append(cur)
        cur = []
    cur.append(i)
if cur:
    batches.append(cur)

# 碎批（<30 則）輪流併入累積中的小批合併箱，上限同 BATCH。
# 跨桶混合可接受：每則翻譯獨立，主題連貫的損失小於多開一個 agent 的成本。
merged, box = [], []
for b in batches:
    if len(b) >= 30:
        merged.append(b)
        continue
    if len(box) + len(b) > BATCH:
        merged.append(box)
        box = []
    box += b
if box:
    merged.append(box)
batches = merged

OUT.mkdir(parents=True, exist_ok=True)
for n, b in enumerate(batches, 1):
    p = OUT / f'batch-{n:03d}.json'
    p.write_text(json.dumps(
        [{'id': i, 'blk': msg[i]['blk'], 'en': msg[i]['en']} for i in b],
        ensure_ascii=False, indent=1), encoding='utf-8')
print(f'{len(batches)} 批 -> {OUT}/batch-001..{len(batches):03d}.json')
for n, b in enumerate(batches, 1):
    print(f'  batch-{n:03d}: {b[0]}-{b[-1]}（{len(b)} 則）')
