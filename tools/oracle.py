#!/usr/bin/env python3
"""probe 文字 hook (0110:3877) 的 oracle 解碼。

用法: oracle.py <round.log>

把 -call-args 印出的 AX 高低位元組還原成遊戲印出的文字,按 2M 步
分組去重後列印。同一則訊息只顯示一次。

注意:hook 行格式在不同 probe 版本有兩種 --
  舊: `#123  由 004F:38BE  0  | AX=004F ...`
  新: `#123  由-1 004F:38BE  0  | AX=004F ...`
正則必須容忍 `由` 後面有無 `-1`(2026-09 R262-R270 教訓:只認舊格式
會把之後所有 OUCH! 判成靜默,連帶誤判推進/撞牆)。
"""
import re
import sys

rows = []
for line in open(sys.argv[1], errors='replace'):
    m = re.match(r'\s+#(\d+).*?\|\s*AX=([0-9A-F]+)', line)
    if m:
        rows.append((int(m.group(1)), int(m.group(2), 16) & 0xFF))

cur = []
prev = 0
out = []
for step, ch in rows:
    if step - prev > 2000000 and cur:
        out.append((cur[0][0], cur))
        cur = []
    cur.append((step, ch))
    prev = step
if cur:
    out.append((cur[0][0], cur))

print(sys.argv[1], len(rows), 'chars')
seen = set()
for s, g in out:
    t = ''.join(chr(c) if 32 <= c < 127 else '<%02X>' % c
                for _, c in g)[:160]
    if t not in seen:
        print('  #%d: %s' % (s, t))
        seen.add(t)
