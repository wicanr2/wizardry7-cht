#!/usr/bin/env python3
"""ports2opl.py — 把 probe `-dump-ports` 的 TSV 轉成 opl_events.py 吃的 golem opl-log。

    tools/py.sh tools/re/ports2opl.py <ports.tsv> <out.opl-log>

TSV 欄位：step、port（hex）、value（hex）（probe writePortLog 的格式）。
OPL2 定址：388h 寫索引、389h 寫資料——成對還原成 (步數, 索引, 值)。
"""
import sys


def main(src, dst):
    out = []
    pending = None
    for line in open(src, encoding="utf-8"):
        parts = line.split()
        if len(parts) < 3 or parts[0] == "step":
            continue
        step = int(parts[0])
        port = int(parts[1], 16)
        val = int(parts[2], 16)
        if port == 0x388:
            pending = (step, val)
        elif port == 0x389 and pending is not None:
            out.append((pending[0], pending[1], val))
            pending = None
    with open(dst, "w", encoding="utf-8") as f:
        for step, reg, val in out:
            f.write(f"{step} {reg:x} {val:x}\n")
    print(f"{dst}: {len(out)} 筆 (索引,值) 對")


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])
