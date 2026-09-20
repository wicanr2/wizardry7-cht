#!/usr/bin/env python3
"""ovr_split.py — 切開自訂 .OVR 容器（docs/re/002 的格式）。

    tools/py.sh tools/re/ovr_split.py <OVR 檔>...      # 在同目錄出 <名>.code/<名>.tail

格式（13/13 檔成立）：檔頭 14 bytes（magic f2 00 55 50、u16 code_size@4、
u16 0x9098@6、u16 tail_size@8、u16 ?@10、u16 checksum?@12），之後
code_size bytes 的 raw x86 code、再 tail_size bytes 的附加塊（語意未知）。
"""
import struct
import sys
from pathlib import Path

MAGIC = b"\xf2\x00\x55\x50"


def split(p: Path):
    d = p.read_bytes()
    if d[:4] != MAGIC:
        raise SystemExit(f"{p}: magic {d[:4].hex()} 不符")
    a, const, b, c, e = struct.unpack_from("<HHHHH", d, 4)
    assert const == 0x9098, f"{p}: 常數 {const:#x} != 0x9098"
    expect = 14 + a + b
    if expect != len(d):
        raise SystemExit(f"{p}: 14+{a}+{b}={expect} != 檔案大小 {len(d)}")
    code, tail = d[14:14 + a], d[14 + a:]
    out = p.with_suffix(".code")
    out.write_bytes(code)
    p.with_suffix(".tail").write_bytes(tail)
    print(f"{p.name}: code={a} tail={b} c={c} e={e:#x} → {out.name}")


if __name__ == "__main__":
    for a in sys.argv[1:]:
        split(Path(a))
