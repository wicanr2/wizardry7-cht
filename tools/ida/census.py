# IDAPython：函式普查（issue #1）。
#
#   tools/ida.sh script census.py PW.EXE.i64 /work/census-PW.json [/seeds/cov.json <runtime_base_hex>]
#
# 輸出一份 JSON：區段、函式（範圍、bytes 開頭、被呼叫數、IDA flags）、每一個 IDA
# 解成指令的 `int`（位址、operand、bytes、所屬函式）、字串數。
#
# 有給覆蓋率時，把實跑執行過、IDA 還不是程式碼的位址種成指令，再普查一次。
# runtime_base 是「執行期線性位址」減掉「映像第一個 byte」的那個值：
#   IDA ea ＝ 映像 base（第一個區段的 start_ea）＋（執行期線性 − runtime_base）
#
# ⚠ headless 的 print 不進 stdout；唯一可信的訊號是輸出檔。
import json
import sys

import ida_auto
import ida_bytes
import ida_funcs
import ida_idaapi
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc


def seg_list():
    out = []
    for s in idautils.Segments():
        seg = ida_segment.getseg(s)
        out.append({
            "name": ida_segment.get_segm_name(seg),
            "start": s, "end": seg.end_ea,
            "sel": seg.sel, "base_para": ida_segment.get_segm_base(seg) // 16,
        })
    return out


def func_list():
    out = []
    for ea in idautils.Functions():
        f = ida_funcs.get_func(ea)
        callers = set()
        for x in idautils.CodeRefsTo(ea, 0):
            callers.add(x)
        out.append({
            "start": f.start_ea, "end": f.end_ea, "size": f.end_ea - f.start_ea,
            "name": ida_funcs.get_func_name(ea),
            "lib": bool(f.flags & ida_funcs.FUNC_LIB),
            "thunk": bool(f.flags & ida_funcs.FUNC_THUNK),
            "far": bool(f.flags & ida_funcs.FUNC_FAR),
            "callers": len(callers),
            "head": ida_bytes.get_bytes(f.start_ea, min(12, f.end_ea - f.start_ea)).hex(),
        })
    return out


def int_sites():
    out = []
    for s in idautils.Segments():
        seg = ida_segment.getseg(s)
        for head in idautils.Heads(s, seg.end_ea):
            if not ida_bytes.is_code(ida_bytes.get_flags(head)):
                continue
            insn = ida_ua.insn_t()
            if ida_ua.decode_insn(insn, head) == 0:
                continue
            if insn.get_canon_mnem() != "int":
                continue
            f = ida_funcs.get_func(head)
            out.append({
                "ea": head,
                "segoff": "%s:%04X" % (ida_segment.get_segm_name(seg), head - ida_segment.get_segm_base(seg)),
                "operand": insn.Op1.value,
                "bytes": ida_bytes.get_bytes(head, insn.size).hex(),
                "func": ida_funcs.get_func_name(f.start_ea) if f else None,
                "disasm": idc.GetDisasm(head),
            })
    return out


def counts():
    code = data = unk = 0
    for s in idautils.Segments():
        seg = ida_segment.getseg(s)
        ea = s
        while ea < seg.end_ea:
            fl = ida_bytes.get_flags(ea)
            size = max(1, ida_bytes.get_item_size(ea))
            if ida_bytes.is_code(fl):
                code += size
            elif ida_bytes.is_data(fl):
                data += size
            else:
                unk += 1
                size = 1
            ea += size
    return {"code_bytes": code, "data_bytes": data, "unknown_bytes": unk}


def snapshot():
    return {
        "functions": func_list(),
        "int_sites": int_sites(),
        "counts": counts(),
        "strings": sum(1 for _ in idautils.Strings()),
    }


def seed(cov_path, runtime_base):
    cov = json.load(open(cov_path))
    img = min(idautils.Segments())
    top = max(ida_segment.getseg(s).end_ea for s in idautils.Segments())
    planted = 0
    skipped_outside = 0
    for sp in cov["spans"]:
        lo = int(sp["start"], 16)
        hi = int(sp["end"], 16)
        for lin in range(lo, hi):
            ea = img + (lin - runtime_base)
            if ea < img or ea >= top:
                skipped_outside += 1
                continue
            if ida_bytes.is_code(ida_bytes.get_flags(ea)):
                continue
            if ida_bytes.is_tail(ida_bytes.get_flags(ea)):
                continue
            ida_bytes.del_items(ea, ida_bytes.DELIT_SIMPLE, 1)
            if ida_ua.create_insn(ea):
                planted += 1
    ida_auto.auto_wait()
    return {"planted_insns": planted, "outside_image_bytes": skipped_outside}


def main():
    ida_auto.auto_wait()
    out_path = idc.ARGV[1]
    result = {
        "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
        "input_file": ida_nalt.get_root_filename(),
        "entry": idc.get_inf_attr(idc.INF_START_EA),
        "image_base": min(idautils.Segments()),
        "segments": seg_list(),
        "before_seed": snapshot(),
    }
    if len(idc.ARGV) >= 4:
        result["seed"] = seed(idc.ARGV[2], int(idc.ARGV[3], 16))
        result["after_seed"] = snapshot()
    with open(out_path, "w") as f:
        json.dump(result, f)
    ida_pro.qexit(0)


main()
