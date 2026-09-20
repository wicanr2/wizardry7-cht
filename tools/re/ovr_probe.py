#!/usr/bin/env python3
"""ovr_probe.py — M0 #2 初探：DS.EXE 檔頭、.OVR 容器格式、字串線索。

docker: tools/py.sh tools/re/ovr_probe.py
"""
import re
import struct
import sys
from pathlib import Path

ORIG = Path("workplace/original/wizar7")


def exe_header(path: Path):
    d = path.read_bytes()
    mz, last_page, pages, relocs, hdr, minalloc, maxalloc, ss, sp, csum, ip, cs, reloc_off, ovl = struct.unpack(
        "<2sHHHHHHHHHHHHH", d[:28]
    )
    load_size = pages * 512 - (512 - last_page) if last_page else pages * 512
    print(f"== {path.name} ({len(d)} bytes)")
    print(f"   sig={mz} last_page={last_page:#x} pages={pages:#x} relocs={relocs}"
          f" hdr={hdr * 16}B minalloc={minalloc * 16}B maxalloc={maxalloc:#x}")
    print(f"   SS={ss:#06x} SP={sp:#06x} CS={cs:#06x} IP={ip:#06x}"
          f" reloc_off={reloc_off:#x} overlay={ovl}")
    print(f"   load_module={load_size}  file_tail={len(d) - load_size}"
          f"  (tail>0 表示 EXE 後面掛了東西)")


def main():
    exe_header(ORIG / "DS.EXE")
    exe_header(ORIG / "VINSTALL.EXE")

    print("\n== .OVR 檔頭（前 16 bytes）與大小")
    for p in sorted(ORIG.glob("*.OVR")):
        head = p.read_bytes()[:16]
        print(f"   {p.name:12} {p.stat().st_size:6}  {head.hex(' ')}")

    print("\n== DS.EXE 內 CD 3F（INT 3Fh，Borland overlay manager 觸發）計數")
    d = (ORIG / "DS.EXE").read_bytes()
    code = d[512:]  # load module
    n_all = sum(1 for i in range(len(d) - 1) if d[i:i + 2] == b"\xcd\x3f")
    n_code = sum(1 for i in range(len(code) - 1) if code[i:i + 2] == b"\xcd\x3f")
    print(f"   全檔 {n_all} 處；load module 內 {n_code} 處")

    print("\n== 字串（ASCII >= 4，含關鍵字的行）")
    pat = re.compile(rb"[\x20-\x7e]{4,}")
    keys = re.compile(
        rb"borland|turbo|pascal|wizard|savant|sir-?tech|copyright|\(c\)|version|"
        rb"mt32|mt-32|roland|adlib|sound blaster|sbpro|gravis|ultrasound|general midi|"
        rb"vga|ega|cgA|hercules|tandy|vesa|mouse|msdos|ms-dos|dos/16m|dos4gw|"
        rb"phar|rational|overla", re.I)
    for name in ("DS.EXE", "VINSTALL.EXE", "EGA.DRV"):
        data = (ORIG / name).read_bytes()
        seen = []
        for m in pat.finditer(data):
            s = m.group()
            if keys.search(s):
                seen.append((m.start(), s))
        print(f"   -- {name}: {len(seen)} 命中")
        for off, s in seen[:60]:
            print(f"      {off:#08x}  {s.decode('ascii', 'replace')}")


if __name__ == "__main__":
    sys.exit(main())
