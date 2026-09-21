#!/usr/bin/env python3
"""vga_sheet.py — 把 .VGA 點陣 sheet（64,000 bytes raw）用執行期色盤渲染成 PNG。

    tools/py.sh tools/re/vga_sheet.py <sheet.vga> <palette.pal> <out.png> [放大]
"""
import struct
import sys
from pathlib import Path

W, H = 320, 200


def write_png(path, w, h, rgb):
    import zlib
    raw = b"".join(b"\0" + rgb[y * w * 3:(y + 1) * w * 3] for y in range(h))

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)

    Path(path).write_bytes(
        b"\x89PNG\r\n\x1a\n"
        + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
        + chunk(b"IDAT", zlib.compress(raw))
        + chunk(b"IEND", b"")
    )


def main(sheet, pal, out, scale=1):
    idx = Path(sheet).read_bytes()
    p = Path(pal).read_bytes()
    if len(idx) != W * H:
        raise SystemExit(f"{sheet}: {len(idx)} bytes，預期 {W*H}")
    scale = int(scale)
    w, h = W * scale, H * scale
    rgb = bytearray(W * H * 3 * scale * scale)
    for y in range(H):
        for x in range(W):
            c = idx[y * W + x] * 3
            base = (y * W * scale + x) * scale * 3
            row = W * scale * 3
            for dy in range(scale):
                o = base + dy * row
                for dx in range(scale):
                    rgb[o + dx*3:o + dx*3+3] = p[c:c+3]
    write_png(out, w, h, bytes(rgb))
    print(f"{out}: {w}x{h}")


if __name__ == "__main__":
    main(*sys.argv[1:5])
