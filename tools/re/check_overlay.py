#!/usr/bin/env python3
# check_overlay.py — docs/spec/001 驗收（不寫原版的印字覆繪）。
#
#   tools/py.sh tools/re/check_overlay.py <中文模式 hd-dir> <純英文模式 hd-dir>
#
# 兩個目錄都要有同一組步數的 hdat-<步數>.{json,a000,comp}（cmd/wiz7 -hd-at 產生）。
# 期望值來自純英文模式的原版輸出與 docs/re/014 E9，不取自被測程式。
import base64
import json
import sys
from collections import Counter
from pathlib import Path

ZH, EN = Path(sys.argv[1]), Path(sys.argv[2])
STEPS = [272000000, 286000000, 330000000]
CURSOR = 16

# E9：防拷四行進入時的 (x0, y0)，以英文原文前綴辨識
COPYPROT = {
    "SECURITY CLEARANCE": (106, 148),
    "PLEASE TURN TO PAGE": (70, 158),
    "MANUAL AND ENTER THE": (64, 164),
    "LOCATED IN PARAGRAPH": (70, 170),
}
MENU = ["CHARACTER MENU", "ADD CHARACTER", "REVIEW MEMBER", "DISMISS MEMBER",
        "START NEW GAME", "CONFIGURATION", "QUIT GAME", "FIGHTER"]
# spec 002 R5：隊員名不譯，單獨的名字不得有覆繪事件
UNTRANSLATED = ["THESUS"]
XLATE = Path(sys.argv[3]) if len(sys.argv) > 3 else Path("workplace/bake/xlate.json")

fails = []


def check(cond, msg):
    print(("通過 " if cond else "失敗 ") + msg)
    if not cond:
        fails.append(msg)


def load(d, st):
    j = json.loads((d / f"hdat-{st}.json").read_text())
    for e in j.get("Shown") or []:
        e["Before"] = base64.b64decode(e["Before"])
        e["After"] = base64.b64decode(e["After"])
    return j, (d / f"hdat-{st}.a000").read_bytes(), (d / f"hdat-{st}.comp").read_bytes()


def find(shown, prefix):
    return [e for e in shown if e["S"].startswith(prefix)]


for st in STEPS:
    zj, _, zcomp = load(ZH, st)
    ej, ea000, _ = load(EN, st)
    shown = zj.get("Shown") or []
    print(f"== {st}（實際步數 {zj['Step']}）{zj['Stats']}")

    # 1. A/B：原版記憶體、A000、調色盤雜湊相同
    for k in ("Step", "Mem", "A000", "Pal"):
        check(zj[k] == ej[k], f"{st} 驗收1 {k} 中={zj[k]} 英={ej[k]}")

    # 3. 英文不殘留：純英文 A000 是筆畫色且≠before 的像素，合成後必須是 before
    mx, my = zj["Mouse"]
    bad = total = 0
    for e in shown:
        bx, by, bw = e["BX"], e["BY"], e["BW"]
        for k, (b, a) in enumerate(zip(e["Before"], e["After"])):
            x, y = bx + k % bw, by + k // bw
            if mx <= x < mx + CURSOR and my <= y < my + CURSOR:
                continue
            i = y * 320 + x
            if ea000[i] == a and a != b:
                total += 1
                if zcomp[i] != b:
                    bad += 1
    check(bad == 0, f"{st} 驗收3 英文殘留 {bad}／{total} 筆畫像素")

    # 字色＝筆畫像素眾數（由 before/after 獨立重算）
    for e in shown:
        cnt = Counter(a for b, a in zip(e["Before"], e["After"]) if a != b)
        top = max(cnt.values())
        check(cnt.get(e["Fg"], 0) == top, f"{st} 字色 {e['S'][:20]!r} fg={e['Fg']} 眾數={cnt.most_common(1)}")

    if st in (272000000, 286000000):
        # 4. 選單與隊員名
        for m in MENU:
            check(any(e["S"] == m for e in shown), f"{st} 驗收4 shown 含 {m!r}")
        for m in UNTRANSLATED:
            check(not any(e["S"] == m for e in shown), f"{st} spec002 R5 shown 不含 {m!r}")
    if st == 286000000:
        # 2. 防拷四行與 E9 座標
        for p, xy in COPYPROT.items():
            hit = find(shown, p)
            check(len(hit) >= 1 and all((e["X0"], e["Y0"]) == xy for e in hit),
                  f"{st} 驗收2 {p!r} (x0,y0)={[(e['X0'], e['Y0']) for e in hit]} 期望 {xy}")
    if st == 286000000:
        # spec 002 R5：1052 的序數槽展開為 __ordinals__["SECOND"]（成員比對）
        sec = json.loads(XLATE.read_text())["__ordinals__"]["SECOND"]
        hit = find(shown, "MANUAL AND ENTER THE")
        check(len(hit) == 1 and sec in hit[0].get("G", []), f"{st} spec002 R5 1052 含序數字形 SECOND")
    if st == 330000000:
        # 5. 旁白時防拷四行不殘留
        for p in COPYPROT:
            check(not find(shown, p), f"{st} 驗收5 無 {p!r}")

print(f"\n失敗 {len(fails)} 項")
sys.exit(1 if fails else 0)
