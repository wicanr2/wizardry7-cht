#!/usr/bin/env python3
# IDAPython：種 raw overlay code 塊的 entry point。
#
#   tools/ida.sh script seed_ovr.py <名>.code.i64
#
# 兩段式（自訂容器沒有 entry stub 表，只能從指令結構長）：
#   1. 線性解碼：對整段 create_insn（嵌入式資料可能解錯——如實記殘量）。
#   2. 只在「被 call 的目標」與 code offset 0 上 add_func（unit 初始化）。
#      函式邊界只來自 call 圖，不硬湊。
# 段 bitness 修成 16-bit。報告寫 /work/seed-report.json。
import json

import ida_auto
import ida_bytes
import ida_funcs
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc
OUT = "/work/seed-report.json"


def main():
    ida_auto.auto_wait()
    rep = {"idb": idc.get_idb_path()}
    for s in idautils.Segments():
        seg = ida_segment.getseg(s)
        seg.bitness = 0  # 16-bit
        seg.update()
    ida_auto.auto_wait()

    end = ida_segment.getseg(0).end_ea

    # 1) 線性解碼
    decoded = failed_decode = 0
    ea = 0
    while ea < end:
        if not idc.is_code(idc.get_full_flags(ea)):
            if ida_ua.create_insn(ea):
                decoded += 1
            else:
                failed_decode += 1
                ida_bytes.create_data(ea, ida_bytes.FF_BYTE, 1, idc.BADADDR)
        ea = idc.next_head(ea, end) if idc.is_code(idc.get_full_flags(ea)) else ea + 1
    ida_auto.auto_wait()

    # 2) call 目標＋offset 0 種函式
    targets = {0}
    for h in idautils.Heads():
        if idc.is_code(idc.get_full_flags(h)):
            for ref in idautils.CodeRefsFrom(h, 0):
                if 0 <= ref < end:
                    targets.add(ref)
    made = dup = fail = 0
    for t in sorted(targets):
        if ida_funcs.get_func(t):
            dup += 1
            continue
        if ida_funcs.add_func(t):
            made += 1
        else:
            fail += 1
    ida_auto.auto_wait()

    funcs = list(idautils.Functions())
    rep.update({
        "linear_decoded": decoded,
        "linear_failed": failed_decode,
        "call_targets": len(targets),
        "func_made": made, "func_dup": dup, "func_fail": fail,
        "functions": len(funcs),
        "segment_end": end,
    })
    json.dump(rep, open(OUT, "w"), indent=1)
    ida_pro.qexit(0)


main()
