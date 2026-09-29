#!/usr/bin/env python3
"""probe 截圖位元組差分：閃爍（<100B）vs 真換視角（>5000B）。

用法: shotdiff.py <a.png> <b.png> [more pairs...]
probe 的 .png 實為 320x200 raw 8-bit（64000B）；火把動畫造成數十位元組
抖動，md5 比對會誤判。R420 校準：同視角閃爍 22-56B；換朝向 10278-18123B。
判準：diff < 500 視為同視角（保守閾值，介於 56 與 10278 之間取整）。
"""
import sys

SAME_VIEW_MAX = 500


def diffbytes(a, b):
    A = open(a, 'rb').read()
    B = open(b, 'rb').read()
    n = min(len(A), len(B))
    d = sum(1 for x, y in zip(A[:n], B[:n]) if x != y)
    return d, len(A), len(B)


rc = 0
for i in range(1, len(sys.argv) - 1, 2):
    a, b = sys.argv[i], sys.argv[i + 1]
    d, la, lb = diffbytes(a, b)
    tag = 'SAME' if (d <= SAME_VIEW_MAX and la == lb) else 'DIFF'
    if tag == 'DIFF':
        rc = 1
    print(f'{d:6d}B {tag} {a} vs {b}')
sys.exit(rc)
