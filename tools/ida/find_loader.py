#!/usr/bin/env python3
# IDAPython：找 overlay loader 與向量掛鉤。
#
#   tools/ida.sh script find_loader.py DS.EXE.i64
#
# 輸出 /work/loader-report.json：
#   - "Error %d loading overlay" 字串位址與 xref 函式
#   - 該函式的完整反組譯（就在 loader 附近）
#   - 所有 mov ah,25h（DOS 設向量）前後的指令（看掛了哪些中斷）
# ⚠ headless print 不進 stdout；唯一可信訊號是輸出檔。
import json

import ida_auto
import ida_bytes
import ida_funcs
import ida_idaapi
import ida_ua
import idautils
import idc

OUT = "/work/loader-report.json"


def disasm_range(lo, hi):
    out = []
    ea = lo
    while ea < hi:
        out.append({"ea": ea, "text": idc.generate_disasm_line(ea, 0)})
        ea = idc.next_head(ea, hi)
    return out


def main():
    ida_auto.auto_wait()
    rep = {}

    # 1) overlay error 字串（bin_search 直接搜 bytes；Strings() 對 $ 結尾收不到）
    pat = b"Error %d loading overlay"
    target = ida_bytes.find_bytes(pat, 0)
    if target is None or target == ida_idaapi.BADADDR:
        rep["error"] = "找不到字串"
        json.dump(rep, open(OUT, "w"), indent=1)
        return
    rep["string_ea"] = target
    rep["string"] = pat.decode()

    # 2) xref 的所屬函式
    funcs = set()
    for xref in idautils.DataRefsTo(target):
        f = ida_funcs.get_func(xref)
        if f:
            funcs.add((f.start_ea, f.end_ea, xref))
    rep["xref_functions"] = [
        {"start": s, "end": e, "xref_site": x} for s, e, x in sorted(funcs)]

    # 3) 倒出這些函式與其 callee 一層
    disasm = {}
    for s, e, x in sorted(funcs):
        disasm[hex(s)] = disasm_range(s, e)
    rep["disasm"] = disasm

    # 4) mov ah,25h（B4 25）→ 前一條指令通常是 mov al,向量號
    setters = []
    for ea in idautils.Heads():
        if not idc.is_code(idc.get_full_flags(ea)):
            continue
        b = ida_bytes.get_byte(ea)
        if b == 0xB4 and ida_bytes.get_byte(ea + 1) == 0x25:
            prev = idc.prev_head(ea)
            nxt = idc.next_head(ea)
            setters.append({
                "ea": ea,
                "prev": idc.generate_disasm_line(prev, 0),
                "here": idc.generate_disasm_line(ea, 0),
                "next": idc.generate_disasm_line(nxt, 0),
            })
    rep["vector_setters_ah25"] = setters

    json.dump(rep, open(OUT, "w"), indent=1)


main()
