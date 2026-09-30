#!/usr/bin/env python3
# bake_help.py — spec 003 R5：F1 說明的字形與頁面。
#
#   tools/re/bake_help.py [輸出路徑]      （預設 workplace/bake/help.json）
#
# 讀 text/help.zh-TW.tsv 與 text/help-names.zh-TW.tsv（欄位 key、zh、source），逐字產生字形：
#   Big5 雙位元組字 → bake_font.glyph_of（SPCFONT.15／stdfont.15），寬 15
#   ASCII 0x20–0x7E → 倚天 ascfont.15（256 字 × 15 列 × 1 byte），寬 8
# 任何字元無字形 → 列出後以非零碼結束，不產出檔案（R5、R6 驗收 4）。
# lint（R6 驗收 5）：ops 頁（序號非 00）每行的 source 必須引用 docs/re/015 實測表中「通過」的 T 編號。
import json
import os
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from bake_font import glyph_of  # noqa: E402

ASC15 = Path.home() / "cht/etan_font/ascfont.15"
SOURCES = [Path(p) for p in os.environ.get(
    "HELP_TSV", "text/help.zh-TW.tsv,text/help-names.zh-TW.tsv").split(",")]  # 驗收可覆寫
EVIDENCE = Path("docs/re/015-f1-help-evidence.md")


def read_tsv(p):
    rows = []
    lines = p.read_text().splitlines()
    if not lines or lines[0].split("\t") != ["key", "zh", "source"]:
        sys.exit(f"{p}：表頭必須是 key<TAB>zh<TAB>source")
    for n, line in enumerate(lines[1:], 2):
        f = line.split("\t")
        if len(f) != 3:
            sys.exit(f"{p}:{n}：欄數 {len(f)}，應為 3")
        rows.append(f)
    return rows


def main():
    out = Path(sys.argv[1]) if len(sys.argv) > 1 else Path("workplace/bake/help.json")
    rows = []
    for p in SOURCES:
        rows += read_tsv(p)
    keys = [k for k, _, _ in rows]
    dup = {k for k in keys if keys.count(k) > 1}
    if dup:
        sys.exit(f"key 重複：{sorted(dup)}")

    # lint：ops 頁的來源必須是已通過的前端實測
    passed = set(re.findall(r"^\| (T\d+) \|.*\| 通過", EVIDENCE.read_text(), re.M))
    bad = [(k, s) for k, _, s in rows
           if k.startswith("ops.") and not k.endswith(".00")
           and not (set(re.findall(r"T\d+", s)) and set(re.findall(r"T\d+", s)) <= passed)]
    if bad:
        for k, s in bad:
            print(f"lint：{k} 的 source {s!r} 不是已通過的前端實測（015 表內通過者：{sorted(passed)}）")
        sys.exit(1)

    asc = ASC15.read_bytes()
    glyphs, missing = {}, []
    for ch in sorted({c for _, zh, _ in rows for c in zh}):
        if " " <= ch <= "~":
            glyphs[ch] = {"w": 8, "rows": list(asc[ord(ch) * 15:ord(ch) * 15 + 15])}
            continue
        try:
            b5 = ch.encode("big5")
            if len(b5) != 2:
                raise ValueError
            # glyph_of：15 px 放在 16 bit 的高 15 位 → 右移 1，改成對齊寬 15
            glyphs[ch] = {"w": 15, "rows": [r >> 1 for r in glyph_of(b5)]}
        except (UnicodeEncodeError, ValueError):
            missing.append(ch)
    if missing:
        print(f"缺字形 {len(missing)} 個：{''.join(missing)!r}")
        sys.exit(1)

    pages = {}
    for k, zh, _ in rows:
        pid, seq = k.rsplit(".", 1)
        pages.setdefault(pid, {})[seq] = zh
    out_pages = []
    for pid, seqs in pages.items():
        if "00" not in seqs:
            sys.exit(f"頁 {pid} 缺標題（序號 00）")
        out_pages.append({"id": pid, "title": seqs["00"],
                          "lines": [seqs[s] for s in sorted(seqs) if s != "00"]})
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps({"glyphs": glyphs, "pages": out_pages}, ensure_ascii=False))
    print(f"help.json：{len(out_pages)} 頁、{len(glyphs)} 字形 → {out}")


if __name__ == "__main__":
    main()
