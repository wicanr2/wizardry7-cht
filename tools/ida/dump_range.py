#!/usr/bin/env python3
# IDAPython：把指定位址範圍的反組譯寫到 /work/range.txt。
#
#   tools/ida.sh script dump_range.py <db> <hex 起> <hex 迄> [名稱]
import sys

import ida_auto
import ida_bytes
import ida_funcs
import ida_pro
import ida_ua
import idautils
import idc

ida_auto.auto_wait()
lo = int(idc.ARGV[1], 16)
hi = int(idc.ARGV[2], 16)
out = f"/work/range-{idc.ARGV[3] if len(idc.ARGV) > 3 else hex(lo)}.txt"
with open(out, "w") as f:
    ea = lo
    while ea < hi:
        if idc.is_code(idc.get_full_flags(ea)):
            f.write(f"{ea:#06x}  {idc.generate_disasm_line(ea, 0)}\n")
            ea = idc.next_head(ea, hi)
        else:
            f.write(f"{ea:#06x}  db {ida_bytes.get_byte(ea):#04x}\n")
            ea += 1
ida_pro.qexit(0)
