#!/usr/bin/env python3
"""R1004: render a disposable layout component from existing R1003 receipts.

No game execution, memory reads, image retouching, or production wiring.  The
blank background deliberately leaves the original-art restoration gate open.
"""
import argparse
import hashlib
import json
from pathlib import Path

from PIL import Image, ImageDraw


KEYS = ("STR", "INT", "PIE", "VIT", "DEX", "SPD", "PER", "KAR")


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def read_font_case(root, backend):
    frame = root / f"api-new-{backend}-2/capture/frame-1585"
    meta_path = frame / "hd.json"
    raw_path = frame / "hd.a000"
    pal_path = frame / "hd.pal"
    meta = json.loads(meta_path.read_text())
    raw, pal = raw_path.read_bytes(), pal_path.read_bytes()
    assert len(raw) == 64000 and len(pal) == 768
    assert meta["Step"] == 363400000
    labels = []
    for row, key in enumerate(KEYS):
        found = [e for e in meta["Shown"] if e["S"] == key]
        assert found
        e = found[0]
        signature = lambda v: (v["X0"], v["Y0"], v["Dy"], v["Fg"], v["G"])
        assert all(signature(v) == signature(e) for v in found)
        assert e["X0"] == 115 and e["Y0"] == 35 + row * 6
        assert e["Dy"] == 0 and e["Fg"] == 85 and len(e["G"]) == 2
        height = len(e["G"][0])
        assert height == (15 if backend == "eten" else 16)
        assert all(len(g) == height and any(g) and
                   all(0 <= word < (1 << height) for word in g) for g in e["G"])
        # Numeric ink is copied from the actual indexed frame, never parsed or
        # recalculated.  White pixels are the two original selector arrows.
        y0 = 35 + row * 6
        dots = [[x * 2, y * 2, list(pal[raw[(y0+y)*320+133+x]*3:][:3])]
                for y in range(6) for x in range(22)
                if raw[(y0+y)*320+133+x] in (1, 85)]
        assert any(color == list(pal[85*3:85*3+3]) for _, _, color in dots)
        labels.append({"key": key, "height": height, "g": e["G"],
                       "rgb": list(pal[85*3:85*3+3]), "dots": dots})
    return {"labels": labels, "inputs": {str(p): sha(p) for p in
            (meta_path, raw_path, pal_path)}}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--receipts", type=Path, required=True)
    ap.add_argument("--out", type=Path, required=True)
    args = ap.parse_args()
    assert args.receipts.is_dir() and not args.out.exists()
    data = {backend: read_font_case(args.receipts, backend)
            for backend in ("eten", "unifont")}
    # Native labels, numeric ink and selector arrows have separate safe cells.
    # New positions use the accepted R953 18px rhythm, without asserting that
    # the early creation path shares R966's source or background contract.
    cells = []
    for row in range(8):
        cells.extend([[230, 70 + row * 18, 34, 18],
                      [266, 70 + row * 18, 44, 18]])
    for i, a in enumerate(cells):
        assert a[0] >= 0 and a[1] >= 0 and a[0]+a[2] <= 640 and a[1]+a[3] <= 400
        for b in cells[i+1:]:
            assert a[0]+a[2] <= b[0] or b[0]+b[2] <= a[0] or a[1]+a[3] <= b[1] or b[1]+b[3] <= a[1]
    payload = json.dumps(data, ensure_ascii=False, separators=(",", ":"))
    page = """<!doctype html><html lang="zh-Hant"><meta charset="utf-8">
<title>R1004 兩倍創角屬性排版原型</title>
<style>body{background:#181818;color:#eee;font:18px sans-serif;margin:24px;max-width:1100px}button{font:inherit;padding:8px;margin-right:8px}.pair{display:flex;gap:24px;flex-wrap:wrap}canvas{background:#242424;image-rendering:pixelated;border:1px solid #555}p{line-height:1.6}.muted{color:#bbb}</style>
<h1>兩倍創角屬性排版原型</h1>
<p>原生字型、實際數值墨點及箭頭。這是欄位示意，灰底不代表正式遊戲背景。</p>
<p><button onclick="show('unifont')">Unifont 16 像素</button><button onclick="show('eten')">倚天 15 像素</button></p>
<div class="pair"><div><h2>目前：12 像素列距</h2><canvas id="old" width="150" height="196"></canvas></div><div><h2>提案：18 像素列距</h2><canvas id="new" width="150" height="196"></canvas></div></div>
<p>提案把文字、原版數值墨點與箭頭按同一列移動。獎勵點數預計移至八列下方。文字保持原生大小，三倍畫面維持現行排版。</p>
<p class="muted">正式接線仍須量到來源、底圖還原、更新與撤除。若採用，移位後的滑鼠操作須由前端對應回原版位置；此頁沒有操作遊戲，也沒有證明映射已完成。</p>
<script>const data=PAYLOAD;
function show(backend){for(const [id,pitch] of [['old',12],['new',18]]){const c=document.getElementById(id),ctx=c.getContext('2d');ctx.clearRect(0,0,c.width,c.height);for(let n=0;n<8;n++){const e=data[backend].labels[n],h=e.height,top=10+n*pitch;ctx.fillStyle=`rgb(${e.rgb})`;for(let i=0;i<e.g.length;i++)for(let y=0;y<h;y++)for(let x=0;x<h;x++)if((e.g[i][y]>>(h-1-x))&1)ctx.fillRect(10+i*(h+2)+x,top+y,1,1);for(const [x,y,rgb]of e.dots){ctx.fillStyle=`rgb(${rgb})`;ctx.fillRect(46+x,top+y+2,2,2)}}}}
show('unifont');</script></html>""".replace("PAYLOAD", payload)
    args.out.mkdir(mode=0o700)
    with (args.out / "compare.html").open("x", encoding="utf-8") as f:
        f.write(page)
    # Render only this newly constructed UI component.  No source PNG is
    # opened or retouched; the original game art is deliberately absent.
    for backend, case in data.items():
        img = Image.new("RGB", (300, 180), (36, 36, 36))
        draw = ImageDraw.Draw(img)
        draw.text((10, 1), "Current: 12px", fill=(238, 238, 238))
        draw.text((160, 1), "Proposal: 18px", fill=(238, 238, 238))
        for left, pitch in ((0, 12), (150, 18)):
            for row, e in enumerate(case["labels"]):
                h, top = e["height"], 20 + row * pitch
                for i, glyph in enumerate(e["g"]):
                    for y, word in enumerate(glyph):
                        for x in range(h):
                            if word & (1 << (h - 1 - x)):
                                img.putpixel((left+10+i*(h+2)+x, top+y), tuple(e["rgb"]))
                for x, y, rgb in e["dots"]:
                    for dx in (0, 1):
                        for dy in (0, 1):
                            img.putpixel((left+46+x+dx, top+y+2+dy), tuple(rgb))
        img.resize((900, 540), Image.Resampling.NEAREST).save(args.out / f"compare-{backend}.png")
    report = {"status": "PROTOTYPE_ONLY", "production_changed": False,
              "game_executed": False, "native_font_heights": [15, 16],
              "current_pitch": 12, "proposed_pitch": 18,
              "proposed_cells": cells, "cell_overlap": 0,
              "inputs": {b: d["inputs"] for b, d in data.items()},
              "generator_sha256": sha(Path(__file__)),
              "html_sha256": sha(args.out / "compare.html"),
              "png_preview_zoom": 3,
              "png_sha256": {b: sha(args.out / f"compare-{b}.png") for b in data},
              "unknown": ["early creation source tuple and return",
                          "background restoration at moved source/destination",
                          "input mapping and lifecycle"]}
    with (args.out / "report.json").open("x", encoding="utf-8") as f:
        json.dump(report, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(json.dumps({k: v for k, v in report.items() if k != "inputs"}, ensure_ascii=False))


if __name__ == "__main__":
    main()
