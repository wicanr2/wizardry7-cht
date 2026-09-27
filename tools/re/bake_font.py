#!/usr/bin/env python3
# bake_font.py — 巫術7 M4：中文字形烘製（VBFONT2 重寫版）
#
# 輸入：譯文字串清單 → 每字經 stdfont.15 抽 15×15 位圖
# 輸出：
#   fontseg.bin   7,440 bytes 完整 VBFONT2 段映像（頭+寬度表+顏色平面+mask 平面）
#   bake.json     字形索引表（byte 0x80+i ↔ 中文字）＋譯文索引串
#
# VBFONT2 格式（docs/re/013 追記二，runtime+檔案大小雙重驗證）：
#   +0000 頭 16B（+03=每列 bytes、+06 flags bit1=寬字形、+0C=mask 偏移、+0E）
#   +0010 寬度表 256B（字元碼 → px 寬）
#   +0110 顏色平面（每字形 寬px→bytes/列 × 15列；本檔 = 2×15 = 30B）
#   +0F10 mask 平面（同幾何；bit1=畫、bit0=透明；畫時顏色 bit1=前景 bit0=背景）
#
# 字形索引：字串 byte 0x80+i → 第 i 個中文字形（0x00-0x7F 填 0 寬＝無字形）。
# 上限：3584 bytes 平面 ÷ 30B = 119 字形（byte 0x80-0xFF 共 128 槽足夠）。
import json
import sys
from pathlib import Path

FONT15 = Path.home() / "cht/etan_font/stdfont.15"
SEG_SIZE = 7440
HDR = 16
WTAB = 256
COLOR_OFF = HDR + WTAB          # 0x110
COL_MASK = 30                   # 每字形：15 列 × 2 bytes
HEIGHT = 15
BPR = 2                         # bytes per row
SCALE = 15                      # 字高（15＝原版倚天字形，不縮）


def big5_bytes(s: str) -> bytes:
    return s.encode("big5")


def glyph_scaled(b5: bytes) -> list[int]:
    """stdfont.15 → SCALE x SCALE 位圖（區間 OR 重採樣：目標像素覆蓋的來源
    像素任一亮即亮。點陣字筆畫僅 1px，最近鄰取樣會跳過筆畫而破相；OR 重採樣
    保留筆畫、代價是略粗）。"""
    rows = glyph_of(b5)
    out = []
    for y in range(SCALE):
        y0, y1 = y * 15 // SCALE, max((y + 1) * 15 // SCALE, y * 15 // SCALE + 1)
        row = 0
        for x in range(SCALE):
            x0, x1 = x * 15 // SCALE, max((x + 1) * 15 // SCALE, x * 15 // SCALE + 1)
            bit = 0
            for sy in range(y0, y1):
                for sx in range(x0, x1):
                    if (rows[sy] >> (15 - sx)) & 1:
                        bit = 1
            if bit:
                row |= 1 << (SCALE - 1 - x)
        out.append(row)
    return out


def glyph_of(b5: bytes) -> list[int]:
    """stdfont.15 → 15 個 row word（MSB first，15px 用高 15 bits）"""
    assert len(b5) == 2
    hi, lo = b5
    idx = (hi - 0xA4) * 157 + ((lo - 0x40) if lo <= 0x7E else (lo - 0xA1 + 63))
    off = idx * 30
    d = FONT15.read_bytes()
    return [int.from_bytes(d[off + r * 2:off + r * 2 + 2], "big") for r in range(HEIGHT)]


def bake(pairs: list[tuple[str, str, int]], outdir: Path):
    strings = [dst for _, dst, _ in pairs]
    font = FONT15.read_bytes()
    # 收集字元（按出現序）
    chars: list[str] = []
    for s in strings:
        for ch in s:
            b5 = ch.encode("big5")
            if len(b5) == 2 and ch not in chars:
                chars.append(ch)
    seg = bytearray(SEG_SIZE)  # 方案 A 遺物；方案 B 只用 xlate.json
    # --- 頭（仿 0x7BB6 原頭）---
    seg[0x03] = BPR                     # 每列 bytes
    seg[0x06] = 0xFF                    # flags bit1 = 寬字形路徑（loc_371）
    seg[0x0C:0x0E] = (COLOR_OFF + 3584).to_bytes(2, "little")   # mask 偏移 0x0F10
    seg[0x0E:0x10] = (0x001C).to_bytes(2, "little")             # mul 乘數（var_A=0 不生效）

    # --- 寬度表：0x80+i → 15px ---
    for i in range(len(chars)):
        seg[HDR + 0x80 + i] = 15

    # --- 顏色平面 + mask 平面 ---
    col = bytearray()
    msk = bytearray()
    index = {}
    for i, ch in enumerate(chars):
        rows = glyph_of(ch.encode("big5"))
        for w in rows:                                   # 顏色：bit1=前景（方案 A 遺物）
            col += bytes(((w >> 8) & 0xFF, w & 0xFF))
            msk += bytes(((w >> 8) & 0xFF, w & 0xFF))    # mask = 字形（1=筆畫）
        index[0x80 + i] = ch
    if len(col) <= 3584:
        seg[COLOR_OFF:COLOR_OFF + len(col)] = col
    mask_off = COLOR_OFF + 3584
    seg[mask_off:mask_off + len(msk)] = msk

    outdir.mkdir(parents=True, exist_ok=True)
    (outdir / "fontseg.bin").write_bytes(seg)

    (outdir / "bake.json").write_text(json.dumps({
        "chars": chars, "index": {f"0x{k:02x}": v for k, v in index.items()},
    }, ensure_ascii=False, indent=1))

    # xlate.json：譯文串 → 字形位圖（Go 側免 BIG5 編碼器）；key＝normalize(en)
    # zh 的 '$'＝數字槽：位圖全 0、槽索引記在 "d"；前端依 runtime 數字從 digits 填。
    # 序數詞（FIRST..TENTH，防拷第二行）normalize 成 '$'、runtime 填的序數由
    # __ordinals__（中文字形）呈現。
    xlate = {}
    for src, dst, dy in pairs:
        entry = []
        slots = []  # '$' 槽索引
        for ch in dst:
            if ch == "$":
                slots.append(len(entry))
                entry.append([0] * HEIGHT)
                continue
            b5 = ch.encode("big5")
            if len(b5) != 2:
                continue
            rows = glyph_scaled(b5)
            entry.append([int(w) for w in rows])
        xlate[normalize(src)] = {"dy": dy, "g": entry}
        if slots:
            xlate[normalize(src)]["d"] = slots
    xlate["__digits__"] = {c: digit_glyph(c) for c in DIGIT5x7}
    xlate["__ordinals__"] = {
        w: [int(v) for v in glyph_scaled(zh.encode("big5"))]
        for w, zh in ORDINALS.items()
    }
    (outdir / "xlate.json").write_text(json.dumps(xlate, ensure_ascii=False))
    print(f"xlate.json：{len(xlate)} 串")
    print(f"字形 {len(chars)} 個；顏色平面 {len(col)}/3584 bytes")
    print(f"輸出：{outdir/'fontseg.bin'}、{outdir/'bake.json'}")
    print("索引：", " ".join(f"{0x80+i:02x}={ch}" for i, ch in enumerate(chars)))


def normalize(s: str) -> str:
    """查表 key 正規化（與 wiz7 前端同一套）：
    1. 剝掉 <0x20 的控制碼（印字屬性：顏色/停頓，不影響對位）
    2. 連續數字 → '$'（防拷頁數等 runtime 填入的動態數字）
    3. 序數詞 FIRST..TENTH → '$'（id 1052 的 $ runtime 填序數詞；
       槽內容由前端依 token 型別填數字位圖或序數中文字形）
    """
    out = []
    i = 0
    while i < len(s):
        c = s[i]
        if ord(c) < 0x20:
            i += 1
            continue
        if c.isdigit():
            j = i
            while j < len(s) and s[j].isdigit():
                j += 1
            out.append("$")
            i = j
            continue
        if c.isalpha():
            j = i
            while j < len(s) and s[j].isalpha():
                j += 1
            word = s[i:j].upper()
            if word in ORDINALS:
                out.append("$")
                i = j
                continue
        if c in "^%":
            # 原文佔位符（段落/行/物品名等 runtime 填入）——一律正規化成 '$'
            out.append("$")
            i += 1
            continue
        out.append(c)
        i += 1
    return "".join(out)


# 5×7 數字位樣（經典點陣）；stdfont.15 無半形/符號區（0xA440 起 30B/字，
# 實測 A2 段取到垃圾）→ 自製數字，放大 2x＝10×14 塞進 15×15 置中。
DIGIT5x7 = {
    "0": ["01110", "10001", "10011", "10101", "11001", "10001", "01110"],
    "1": ["00100", "01100", "00100", "00100", "00100", "00100", "01110"],
    "2": ["01110", "10001", "00001", "00010", "00100", "01000", "11111"],
    "3": ["11110", "00001", "00001", "01110", "00001", "00001", "11110"],
    "4": ["00010", "00110", "01010", "10010", "11111", "00010", "00010"],
    "5": ["11111", "10000", "11110", "00001", "00001", "10001", "01110"],
    "6": ["00110", "01000", "10000", "11110", "10001", "10001", "01110"],
    "7": ["11111", "00001", "00010", "00100", "01000", "01000", "01000"],
    "8": ["01110", "10001", "10001", "01110", "10001", "10001", "01110"],
    "9": ["01110", "10001", "10001", "01111", "00001", "00010", "01100"],
}

# 防拷序數詞（id 1052 的 $ runtime 填序數詞，id 1060-1069）→ 中文數字
ORDINALS = {
    "FIRST": "一", "SECOND": "二", "THIRD": "三", "FOURTH": "四",
    "FIFTH": "五", "SIXTH": "六", "SEVENTH": "七", "EIGHTH": "八",
    "NINTH": "九", "TENTH": "十",
}


def digit_glyph(ch: str) -> list[int]:
    """數字 → 15×15 位圖（5×7 放大 2x＝10×14，置中）。"""
    rows5 = DIGIT5x7[ch]
    big = []
    for r in rows5:
        row = 0
        for x, c in enumerate(r):
            if c == "1":
                row |= 0b11 << (14 - x * 2)  # 水平 2x（左對齊 10px）
        big.append(row)
        big.append(row)  # 垂直 2x
    # big = 14 列；上移 0 列、兩側置中由使用者位移（此處左 2px）
    out = [0] * 15
    for i, row in enumerate(big):
        out[i] = row >> 2  # 10px 內容右移置中（14-10)/2=2
    return out


if __name__ == "__main__":
    # 譯文表（單一權威）：text/story-trans.json（id→{en,zh,dy?}）。
    # xlate.json 的 key＝normalize(en)（前端查表同款正規化）。
    trans = json.loads(Path("text/story-trans.json").read_text())
    pairs = [
        (it["en"], it["zh"], it.get("dy", 0))
        for it in trans["items"].values()
    ]
    outdir = Path("workplace/bake")
    # 開場故事（SCENARIO 系統，非 MSG.DBS 訊息 ID 系統；逐行 1:1 映射）
    story = Path("workplace/bake/story.txt")
    if story.exists():
        for line in story.read_text().splitlines():
            src, _, dst = line.partition("\t")
            if src:
                pairs.append((src, dst, 0))
    # 怪物譯文（text/monsters-trans.json，key＝名字原文）併入
    mp = "text/monsters-trans.json"
    if Path(mp).exists():
        mons = json.loads(Path(mp).read_text())["items"]
        for en, zh in mons.items():
            pairs.append((en, zh, 0))
        print(f"怪物譯文 {len(mons)} 筆併入")
    # 物品譯文（text/items-trans.json，key＝名字原文）併入
    ip = "text/items-trans.json"
    if Path(ip).exists():
        its = json.loads(Path(ip).read_text())["items"]
        for en, zh in its.items():
            pairs.append((en, zh, 0))
        print(f"物品譯文 {len(its)} 筆併入")
    bake(pairs, outdir)
