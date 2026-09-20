#!/usr/bin/env python3
# IDAPython：在指定位址強制解碼成指令並建函式（經指標進入、沒有 call/xref 的
# 事件處理常式用）。
#
#   tools/ida.sh script decode_at.py <db> <hex 位址…>
import ida_auto
import ida_bytes
import ida_funcs
import ida_pro
import ida_ua
import idc

ida_auto.auto_wait()
for a in idc.ARGV[1:]:
    if "-" in a:  # start-end：整段清掉後從頭線性解碼
        lo_s, hi_s = a.split("-")
        lo, hi = int(lo_s, 16), int(hi_s, 16)
        ida_bytes.del_items(lo, ida_bytes.DELIT_EXPAND, hi - lo)
        ea = lo
        while ea < hi:
            if ida_ua.create_insn(ea):
                ea = idc.next_head(ea, hi)
                if ea <= lo or ea == idc.BADADDR:
                    break
            else:
                ea += 1
        ida_funcs.add_func(lo)
    else:
        ea = int(a, 16)
        ida_bytes.del_items(ea, ida_bytes.DELIT_SIMPLE, 1)
        if ida_ua.create_insn(ea):
            ida_funcs.add_func(ea)
ida_auto.auto_wait()
ida_pro.qexit(0)
