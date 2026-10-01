#!/usr/bin/env python3
"""spec007：固定版發行字型來源；沒有字型時中止，不偷用本機字型。"""
import gzip
import hashlib
import os
from functools import lru_cache
from pathlib import Path

BACKEND = os.environ.get("WIZ7_FONT", "eten")
if BACKEND not in ("eten", "unifont"):
    raise ValueError("WIZ7_FONT 必須為 eten 或 unifont")
SIZE = 16 if BACKEND == "unifont" else 15
UNIFONT_SHA256 = "7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574"


@lru_cache(maxsize=1)
def unifont():
    path = os.environ.get("WIZ7_UNIFONT_HEX")
    if not path:
        raise ValueError("Unifont 後端需要 WIZ7_UNIFONT_HEX")
    raw = Path(path).read_bytes()
    if hashlib.sha256(raw).hexdigest() != UNIFONT_SHA256:
        raise ValueError("Unifont 固定版 SHA-256 不符")
    glyphs = {}
    for line in gzip.decompress(raw).decode("ascii").splitlines():
        cp, bits = line.split(":")
        width = len(bits) // 4
        if width not in (8, 16) or len(bits) not in (32, 64):
            raise ValueError("Unifont 字形不是8×16或16×16")
        ch = chr(int(cp, 16))
        if ch in glyphs:
            raise ValueError("Unifont 碼點重複")
        glyphs[ch] = (width, tuple(int(bits[i:i + width // 4], 16)
                                  for i in range(0, len(bits), width // 4)))
    return glyphs


def unicode_glyph(ch, square=False):
    try:
        width, rows = unifont()[ch]
    except KeyError:
        raise ValueError(f"Unifont 缺字 U+{ord(ch):04X}") from None
    if square and width == 8:
        return 16, [row << 4 for row in rows]
    return width, list(rows)
