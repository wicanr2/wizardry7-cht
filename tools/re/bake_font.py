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


def big5_bytes(s: str) -> bytes:
    return s.encode("big5")


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
    assert len(chars) <= 119, f"字形數 {len(chars)} > 119"

    seg = bytearray(SEG_SIZE)
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
        for w in rows:                                   # 顏色：bit1=前景
            col += bytes(((w >> 8) & 0xFF, w & 0xFF))
            msk += bytes(((w >> 8) & 0xFF, w & 0xFF))    # mask = 字形（1=筆畫）
        index[0x80 + i] = ch
    assert len(col) <= 3584, f"顏色平面 {len(col)} > 3584"
    seg[COLOR_OFF:COLOR_OFF + len(col)] = col
    mask_off = COLOR_OFF + 3584
    seg[mask_off:mask_off + len(msk)] = msk

    outdir.mkdir(parents=True, exist_ok=True)
    (outdir / "fontseg.bin").write_bytes(seg)

    # 譯文 → 索引串
    trans = []
    for s in strings:
        enc = []
        for ch in s:
            b5 = ch.encode("big5")
            enc.append(0x80 + chars.index(ch) if len(b5) == 2 else ord(ch) & 0x7F)
        trans.append(bytes(enc).hex())
    (outdir / "bake.json").write_text(json.dumps({
        "chars": chars, "index": {f"0x{k:02x}": v for k, v in index.items()},
        "translations": trans,
    }, ensure_ascii=False, indent=1))

    # xlate.json：譯文串 → 字形位圖（Go 側免 BIG5 編碼器）
    xlate = {}
    for src, dst, dy in pairs:
        entry = []
        for ch in dst:
            b5 = ch.encode("big5")
            if len(b5) != 2:
                continue
            rows = glyph_of(b5)
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
        ("SECURITY CLEARANCE", "安全檢查", 0),
        ("PLEASE TURN TO PAGE 35 OF YOUR", "請翻開手冊第35頁", 0),
        ("MANUAL AND ENTER THE SECOND WORD", "找出第二個字輸入", 9),
        ("LOCATED IN PARAGRAPH 3  LINE 5", "在第3段第5行", 18),
    ]
    bake(pairs, Path("workplace/bake"))
