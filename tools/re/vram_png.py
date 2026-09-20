#!/usr/bin/env python3
"""vram_png.py — 模式 13h 的 320×200 色號陣列＋768 bytes 6-bit 色盤 → PNG。

    tools/py.sh tools/re/vram_png.py <frame.vram> <palette.pal> <out.png>
"""
import struct
import sys
import zlib
from pathlib import Path

W, H = 320, 200


def write_png(path, w, h, rgb):
    raw = b"".join(b"\0" + rgb[y * w * 3:(y + 1) * w * 3] for y in range(h))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)

    Path(path).write_bytes(
        b"\x89PNG\r\n\x1a\n"
        + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
        + chunk(b"IDAT", zlib.compress(raw))
        + chunk(b"IEND", b"")
    )


def main(vram_path, pal_path, out_path):
    idx = Path(vram_path).read_bytes()
    pal = Path(pal_path).read_bytes()
    if len(idx) != W * H:
        raise SystemExit(f"{vram_path}：{len(idx)} bytes，不是 {W}×{H}")
    if len(pal) != 768:
        raise SystemExit(f"{pal_path}：{len(pal)} bytes，不是 768")
    rgb = bytearray(W * H * 3)
    for i, c in enumerate(idx):
        p = c * 3
        rgb[i * 3] = (pal[p] << 2) | (pal[p] >> 4)
        rgb[i * 3 + 1] = (pal[p + 1] << 2) | (pal[p + 1] >> 4)
        rgb[i * 3 + 2] = (pal[p + 2] << 2) | (pal[p + 2] >> 4)
    write_png(out_path, W, H, bytes(rgb))
    nz = sum(1 for x in idx if x)
    print(f"{out_path}: 非零像素 {nz}/{W * H}，使用色彩 {len(set(idx))}")


if __name__ == "__main__":
    main(*sys.argv[1:4])
