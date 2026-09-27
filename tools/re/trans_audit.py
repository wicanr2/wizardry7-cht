#!/usr/bin/env python3
# trans_audit.py — 翻譯覆蓋率對帳（rulebook/83：完整性要有數字）
#
# 未翻 = messages.json 全量 - story-trans 已含 - 明確不翻類。
# 分類輸出統計＋各類抽樣，供人工核對邊界後訂出真實工作量。
import json
import sys

MSG = json.load(open('text/messages.json'))['records']
ST = json.load(open('text/story-trans.json'))['items']
done = {int(k) for k in ST}

# 不翻段（013 追記三十一：原因明確）
PROTECTED = [range(2500, 2575), range(7000, 7198)]  # 手冊防拷關鍵字、NLP 詞庫


def in_protected(i):
    return any(i in r for r in PROTECTED)


def classify(en):
    import re
    if en.strip() == '':
        return 'empty'
    if en.lstrip().startswith('*'):
        return 'star'          # credits 人名 / 職級終極道具（留原文）
    if '\x7f' in en:
        return 'deco'          # 裝飾碼
    if en.strip().startswith('!@#'):
        return 'symbol'        # 防拷符號列
    if re.match(r'^[\x00-\x1f](NO|YE|SURE|OK|MAN|WOMAN)/', en):
        return 'yesno'         # \x0N(NO|YES)/ 玩家回應詞庫（疑似輸入比對）→ 待查，不進批次
    if en.strip().isdigit():
        return 'number'        # 座標/價格 → 不翻
    body = re.sub(r'^[\x00-\x1f]+', '', en).strip('!#@%~\t\n\r\x0b\x0c-')
    if not re.search(r'[A-Za-z0-9]', body):
        return 'ctrl'          # 剝控制碼/符號後無字母數字 → 引擎內部碼，不翻
    lead = re.sub(r'^[\x00-\x1f]+', '', en)
    # 全大寫縮寫/密文（PS/AC/NONE/HEDGHFLDFRM）不再單獨分類：NPC 名、按鈕、
    # 縮寫進批次（指令檔有判斷規則），子音密文由 subagent 保留原文記 notes。
    return 'todo'  # 其餘有可見文字者一律要翻（含控制碼開頭句、敘事、對話）


todo, cls_cnt, samples = [], {}, {}
for m in MSG:
    i = m['id']
    if i in done:
        continue
    if in_protected(i):
        cls_cnt['protected'] = cls_cnt.get('protected', 0) + 1
        continue
    c = classify(m['en'])
    cls_cnt[c] = cls_cnt.get(c, 0) + 1
    samples.setdefault(c, []).append((i, m['en']))
    if c == 'todo':
        todo.append(i)

total = len(MSG)
print(f'總數 {total}  已翻 {len(done)}  不翻/待核 {total - len(done) - len(todo)}  要翻 {len(todo)}')
print('分類:', json.dumps(cls_cnt, ensure_ascii=False, sort_keys=True))
for c, lst in sorted(samples.items()):
    print(f'\n[{c}] {len(lst)} 則，抽樣:')
    for i, en in lst[:6]:
        print(f'  {i}\t{en[:60]!r}')

# 要翻的 id 桶分布
todo.sort()
buckets = {}
for i in todo:
    b = i // 1000 * 1000
    buckets[b] = buckets.get(b, 0) + 1
print('\n要翻 id 桶分布:')
for b, n in sorted(buckets.items()):
    print(f'  {b}-{b+999}: {n}')

json.dump(todo, open('workplace/todo_ids.json', 'w'))
print(f'\n要翻清單 -> workplace/todo_ids.json')
