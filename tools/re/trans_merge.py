#!/usr/bin/env python3
# trans_merge.py — 批次譯文合併＋逐行驗證（batch-subagent-localization 第 5 步）
#
# 驗證 workplace/batches/batch-NNN.out.json：
#   1. id 集合與批檔完全一致（不可缺、不可多、不可改）
#   2. 控制碼 & % ] ^ 出現次數與 en 相同
#   3. zh 不以 ! / @ 開頭
#   4. BIG5 可編
#   5. __notes__ 鍵抽出顯示（新譯名記錄，不入表）
# 全過才併入 text/batches-merged.json（id→zh，進版控；gen_trans.py 讀它）。
import json
import sys
from pathlib import Path

BATCHES = Path('workplace/batches')
MERGED = 'text/batches-merged.json'
CTRL = '&%]^'

merged = json.load(open(MERGED)) if Path(MERGED).exists() else {}
notes_all = []
n_new = 0

for p in sorted(BATCHES.glob('batch-*.out.json')):
    batch = json.load(open(p.with_suffix('').with_suffix('.json')))  # batch-NNN.json
    out = json.load(open(p))
    want = {str(m['id']) for m in batch}
    got = {k for k in out if k != '__notes__'}
    if got != want:
        print(f'{p.name}: id 不一致 缺{sorted(want-got)[:5]} 多{sorted(got-want)[:5]}',
              file=sys.stderr)
        sys.exit(1)
    if '__notes__' in out:
        notes_all.append((p.name, out['__notes__']))
    for m in batch:
        i, en, zh = str(m['id']), m['en'], out[str(m['id'])]
        for c in CTRL:
            if en.count(c) != zh.count(c):
                print(f'{p.name} id {i}: "{c}" 次數不符 '
                      f'(en {en.count(c)} zh {zh.count(c)})\n  en: {en!r}\n  zh: {zh!r}',
                      file=sys.stderr)
                sys.exit(1)
        if zh[:1] in ('!', '@'):
            print(f'{p.name} id {i}: zh 不應以 !/@ 開頭: {zh!r}', file=sys.stderr)
            sys.exit(1)
        try:
            zh.encode('big5')
        except UnicodeEncodeError as e:
            print(f'{p.name} id {i}: 非 BIG5 字: {e}\n  zh: {zh!r}', file=sys.stderr)
            sys.exit(1)
        if i in merged and merged[i] != zh:
            # merged 已含後續修正（一致性掃描/人工精煉）——以 merged 為準，
            # 原始批次輸出回寫成 merged 版，保持兩邊同步。列印供稽核。
            print(f'[同步] {p.name} id {i}: out {zh!r} -> merged {merged[i]!r}')
            out[i] = merged[i]
            json.dump(out, open(p, 'w'), ensure_ascii=False, indent=1)
            continue
        if i not in merged:
            merged[i] = zh
            n_new += 1

json.dump(merged, open(MERGED, 'w'), ensure_ascii=False, indent=1)
print(f'{MERGED}: +{n_new} 筆（共 {len(merged)}）')
for name, nt in notes_all:
    print(f'[notes] {name}: {nt}')
