#!/usr/bin/env python3
# bake_flow.py — spec 006 R1：流動文字（引擎折行成頁的段落）資料。
#
#   tools/re/bake_flow.py [輸出路徑] [--drop <訊息 id>]      （預設 workplace/bake/flow.json）
#
# 從 text/messages.json 與 text/story-trans.json 建段落、字詞序列、每行結束位置、譯文與字形：
#   行     剝控制碼（直接刪除 <0x20）後去掉行首一個 `@`
#   被排除 含 & % ] ^ $ # ! " 任一字元、含數字、含行首以外的 `@`、訊息 id 2500–2574；被排除的行截斷段落
#   段落   連續訊息 id（相鄰以原始 id 計）的行；剝控制碼後行首是 `@` 的行開新段落；
#          零字詞行維持相鄰但不加字詞、不記 end；沒有任何字詞的段落丟棄
#   譯文   story-trans 同 id 的 zh，非空且每個字都是 Big5 雙位元組字、且取得到字形，才可用（否則 zh 為空）
#   字形   bake_font.glyph_scaled（15 位元靠右，與 xlate.json 同一函式）
# --drop <id>：移除含該 id 的整個段落（驗收 9 的負向測試檔）。
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from bake_font import glyph_scaled  # noqa: E402

MSG = Path("text/messages.json")
TRANS = Path("text/story-trans.json")
BAD = set('&%]^$#!"')
PROT = range(2500, 2575)  # 手冊防拷關鍵字表：一律排除，不寫進任何輸出


def strip_ctl(s: str) -> str:
    return "".join(c for c in s if c >= " ")


def parse_line(en: str):
    """回傳 (文字, 是否行首 @)；被排除的行回傳 None。"""
    s = strip_ctl(en)
    lead = s.startswith("@")
    if lead:
        s = s[1:]
    if any(c in BAD for c in s) or "@" in s or any(c.isdigit() for c in s):
        return None
    return s, lead


def usable_zh(zh: str, glyphs: dict) -> str:
    """可用的譯文：非空、每個字 Big5 雙位元組、取得到字形。順便把字形收進 glyphs。"""
    if not zh:
        return ""
    got = {}
    for ch in zh:
        try:
            b5 = ch.encode("big5")
        except UnicodeEncodeError:
            return ""
        if len(b5) != 2:
            return ""
        if ch not in glyphs and ch not in got:
            try:
                got[ch] = [int(v) for v in glyph_scaled(b5)]
            except ValueError:
                return ""
    glyphs.update(got)
    return zh


def build(drop=None):
    records = sorted(json.loads(MSG.read_text())["records"], key=lambda r: r["id"])
    trans = json.loads(TRANS.read_text())["items"]
    raw_paras = []
    cur, prev = None, None
    for r in records:
        i = r["id"]
        parsed = None if i in PROT else parse_line(r["en"])
        if parsed is None:
            cur, prev = None, i
            continue
        text, lead = parsed
        if cur is None or prev is None or i != prev + 1 or lead:
            cur = []
            raw_paras.append(cur)
        cur.append((i, text))
        prev = i

    glyphs, paras = {}, []
    stats = {"paras": 0, "lines": 0, "usable_lines": 0, "full_paras": 0}
    for p in raw_paras:
        words, lines = [], []
        for i, text in p:
            w = text.split()
            if not w:
                continue  # 零字詞行：不加字詞、不記 end
            words += w
            zh = usable_zh(trans.get(str(i), {}).get("zh", ""), glyphs)
            lines.append({"id": i, "end": len(words), "zh": zh})
        if not lines:
            continue
        if drop is not None and any(i == drop for i, _ in p):
            continue
        paras.append({"id": lines[0]["id"], "words": words, "lines": lines})
        stats["paras"] += 1
        stats["lines"] += len(lines)
        stats["usable_lines"] += sum(1 for l in lines if l["zh"])
        stats["full_paras"] += all(l["zh"] for l in lines)
    used = {ch for p in paras for l in p["lines"] for ch in l["zh"]}
    return {"glyphs": {c: g for c, g in glyphs.items() if c in used}, "paras": paras}, stats


def main():
    args = sys.argv[1:]
    drop = None
    if "--drop" in args:
        k = args.index("--drop")
        drop = int(args[k + 1])
        del args[k:k + 2]
    out = Path(args[0]) if args else Path("workplace/bake/flow.json")
    data, st = build(drop)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(data, ensure_ascii=False, separators=(",", ":")))
    print(f"flow.json：段落 {st['paras']}、可索引行 {st['lines']}、譯文可用的行 {st['usable_lines']}、"
          f"整段都可翻譯的段落 {st['full_paras']}、字形 {len(data['glyphs'])} → {out}")


if __name__ == "__main__":
    main()
