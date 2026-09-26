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

    # xlate.json：譯文串 → 字形位圖（Go 側免 BIG5 編碼器）
    xlate = {}
    for src, dst, dy in pairs:
        entry = []
        for ch in dst:
            b5 = ch.encode("big5")
            if len(b5) != 2:
                continue
            rows = glyph_scaled(b5)
            entry.append([int(w) for w in rows])
        xlate[src] = {"dy": dy, "g": entry}
    (outdir / "xlate.json").write_text(json.dumps(xlate, ensure_ascii=False))
    print(f"xlate.json：{len(xlate)} 串")
    print(f"字形 {len(chars)} 個；顏色平面 {len(col)}/3584 bytes")
    print(f"輸出：{outdir/'fontseg.bin'}、{outdir/'bake.json'}")
    print("索引：", " ".join(f"{0x80+i:02x}={ch}" for i, ch in enumerate(chars)))


if __name__ == "__main__":
    pairs = [
        ("CHARACTER MENU", "角色選單", 0),
        ("ADD CHARACTER", "新增角色", 0),
        ("LOAD SAVEGAME", "載入紀錄", 0),
        ("CONFIGURATION", "系統設定", 0),
        ("IMPORT (BANE)", "匯入隊伍", 0),
        ("QUIT GAME", "結束遊戲", 0),
        ("CANCEL", "取消", 0),
        # 防拷三行的數字（頁/段/行）是另一個印字呼叫、動態——譯文不含數字
        ("SECURITY CLEARANCE", "安全檢查", 0),
        ("PLEASE TURN TO PAGE 35 OF YOUR", "請翻開手冊指定頁數", 0),
        ("MANUAL AND ENTER THE SECOND WORD", "找出該頁第二個字輸入", 13),
        ("LOAD", "載入", 0),
        ("DIRECTORY", "目錄", 0),
        ("FILENAME", "檔名", 0),
        ("START NEW GAME", "開始新遊戲", 0),
        ("REVIEW MEMBER", "檢視成員", 0),
        ("DISMISS MEMBER", "開除成員", 0),
        ("FIGHTER", "戰士", 0),
        ("Written & Programmed by", "編寫與程式設計", 0),
        ("INPUT DEVICE", "輸入裝置", 0),
        ("DELAY", "延遲", 0),
        ("PATH", "路徑", 0),
        ("FILE", "檔案", 0),
        ("PLAY LEVEL", "遊戲等級", 0),
        ("FADE", "淡入", 0),
        ("MEMORY", "記憶體", 0),
        ("EFFECTS", "音效", 0),
        ("GRAPHICS", "圖形", 0),
        ("BRIGHT", "亮度", 0),
        ("CONTRAST", "對比", 0),
        ("MUSIC", "音樂", 0),
        ("SAVE", "存檔", 0),
        ("Manual & Documentation", "手冊與文件", 0),
        ("Music Score & Sound Effects", "音樂與音效", 0),
        ("Technical Programming", "技術程式設計", 0),
        ("Sound Programming", "音效程式設計", 0),
        ("Graphic Artists", "美術設計", 0),
        ("Graphics Maestro", "美術總監", 0),
    ]
    # 開場故事（逐行 1:1 映射，tools 產生的 workplace/bake/story.txt）
    story = Path("workplace/bake/story.txt")
    if story.exists():
        for line in story.read_text().splitlines():
            src, _, dst = line.partition("\t")
            if src:
                pairs.append((src, dst, 0))
    bake(pairs, Path("workplace/bake"))
