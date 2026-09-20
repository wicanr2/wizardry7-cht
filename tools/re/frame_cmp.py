#!/usr/bin/env python3
"""frame_cmp.py — DOSBox-X（640×400 RGB）vs dosgolem（模式 13h 色號＋色盤）對拍。

    tools/py.sh tools/re/frame_cmp.py <dosboxx 目錄> <golem vram> <golem pal>

每張 DOSBox-X 截圖：驗 2×2 區塊同色 → 縮回 320×200 → 與 golem 幀比
（RGB 容差 tol=6/通道，吃 DAC 6→8 bit 捨入差）。輸出每張的非黑像素數、
差異像素數、差異最大列，最後挑差異最小的一張當候選。
"""
import sys
from pathlib import Path

W, H = 320, 200
TOL = 6


def load_dbx(p: Path):
    d = p.read_bytes()
    if len(d) != 640 * 400 * 3:
        raise ValueError(f"{p}: {len(d)} bytes")
    px = []
    uniform = True
    for y in range(H):
        row = []
        for x in range(W):
            r0 = (y * 2) * 640 * 3 + (x * 2) * 3
            block = [d[r0:r0 + 3], d[r0 + 3:r0 + 6],
                     d[r0 + 640 * 3:r0 + 640 * 3 + 3], d[r0 + 640 * 3 + 3:r0 + 640 * 3 + 6]]
            if len({tuple(b) for b in block}) != 1:
                uniform = False
            row.append(tuple(block[0]))
        px.append(row)
    return px, uniform


def load_golem(vram: Path, pal: Path):
    idx = vram.read_bytes()
    pal = pal.read_bytes()  # probe 已輸出 8-bit RGB
    return [[tuple(pal[c * 3:c * 3 + 3]) for c in idx[y * W:(y + 1) * W]]
            for y in range(H)]


def main(dbx_dir, vram, pal):
    g = load_golem(Path(vram), Path(pal))
    g_mask = [[c != (0, 0, 0) for c in row] for row in g]
    best = best_mask = None
    for p in sorted(Path(dbx_dir).glob("title-*.rgb")):
        px, uniform = load_dbx(p)
        nonblack = sum(1 for row in px for c in row if c != (0, 0, 0))
        diff = sum(1 for y in range(H) for x in range(W)
                   if max(abs(px[y][x][i] - g[y][x][i]) for i in range(3)) > TOL)
        mask = sum(1 for y in range(H) for x in range(W)
                   if (px[y][x] != (0, 0, 0)) != g_mask[y][x])
        print(f"{p.stem}: 2x2一致={uniform} 非黑={nonblack:6} 差異像素={diff:6} 遮罩差異={mask:6}")
        if best is None or diff < best[1]:
            best = (p.stem, diff)
        if best_mask is None or mask < best_mask[1]:
            best_mask = (p.stem, mask)
    print(f"顏色最佳：{best[0]}（差異 {best[1]}）")
    print(f"版面最佳：{best_mask[0]}（遮罩差異 {best_mask[1]}）")


if __name__ == "__main__":
    main(*sys.argv[1:4])
