#!/usr/bin/env python3
"""R976：保留已採甲的前景形狀，只改紙頁文字RGB。"""
import argparse
import hashlib
import json
import sys
from pathlib import Path

from check_native_gui import png
from prototype_intro_paper import png_write

SOURCE_SHAS = {
    'unifont': 'e9beb480bc01eb6a18962451e597d74373590161e0b1808f17baa5d8d78b43a9',
    'eten': '1981685b5f95aea56819f6bb8fbd345ee848717d667601e5eb1bdc235da5aece',
}
COLORS = {'red': (236,72,40), 'black': (248,232,160)}


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('source', type=Path)
    p.add_argument('output', type=Path)
    args = p.parse_args()
    if args.output.exists():
        p.error('輸出已存在，請用新目錄')
    raw = (args.source/'receipt.json').read_bytes()
    rec = json.loads(raw)
    backend, size = rec['backend'], rec['size']
    assert sha(raw) == SOURCE_SHAS[backend]
    prepared = []
    for row in rec['frames']:
        if row['style'] != 'bold':
            continue
        before_path = args.source/row['output']
        assert sha(before_path.read_bytes()) == row['output_sha256']
        w,h,before = png(before_path)
        assert (w,h) == (960,600)
        after = bytearray(before)
        mask = set()
        for label in row['labels']:
            x0,y0 = label['xy']
            old = bytes(label['rgb'])
            new = bytes(COLORS[label['color']])
            for j,c in enumerate(label['zh']):
                glyph = rec['glyphs'][c]
                for gy,bits in enumerate(glyph):
                    for gx in range(size):
                        if not bits & (1 << (size-1-gx)):
                            continue
                        for dy in (0,1):
                            for dx in (0,1,2):
                                x = x0+j*(size*2+4)+gx*2+dx
                                y = y0+gy*2+dy
                                at = y*w+x
                                assert before[at*3:at*3+3] == old
                                after[at*3:at*3+3] = new
                                mask.add(at)
        assert all(before[i*3:i*3+3]==after[i*3:i*3+3] for i in range(w*h) if i not in mask)
        prepared.append((row,before,after,len(mask)))
    assert len(prepared)==5
    args.output.mkdir()
    frames,sections = [],[]
    for row,before,after,count in prepared:
        old_name = f'{backend}-{row["step"]}-old-bold.png'
        new_name = f'{backend}-{row["step"]}-bright-bold.png'
        png_write(args.output/old_name,960,600,before)
        png_write(args.output/new_name,960,600,after)
        frames.append(dict(step=row['step'],old=old_name,output=new_name,foreground=count,
                           source_png_sha256=row['output_sha256'],output_sha256=sha((args.output/new_name).read_bytes())))
        sections.append(f'<section><h2>{row["step"]}</h2><div><figure><figcaption>原甲</figcaption><img src="{old_name}"></figure><figure><figcaption>亮字甲</figcaption><img src="{new_name}"></figure></div></section>')
    data = dict(schema=1,status='DRAFT brightness prototype',backend=backend,source_receipt_sha256=sha(raw),
                tool_sha256=sha(Path(__file__).read_bytes()),python=sys.version.split()[0],
                colors={k:list(v) for k,v in COLORS.items()},frames=frames,
                limits='甲的字形、大小、加粗、位置及背景保持。只改前景RGB；不代表正式GUI或原版同狀態覆繪。')
    with (args.output/'receipt.json').open('x') as f:
        json.dump(data,f,ensure_ascii=False,indent=2);f.write('\n')
    page = ('<!doctype html><html lang="zh-Hant"><meta charset="utf-8"><title>R976紙頁亮字</title>'
            '<style>body{background:#181818;color:#eee;font:18px sans-serif;margin:24px}section>div{display:flex;flex-wrap:wrap;gap:20px}figure{margin:0}img{image-rendering:pixelated}</style>'
            '<h1>紙頁亮字甲</h1><p>維持已採用的大小與加粗；深色字改米金亮字，紅字提高亮度。背景條紋仍待改善。這是外觀原型。</p>'+''.join(sections)+'</html>')
    (args.output/'compare.html').write_text(page)
    print(f'{backend}：五頁，前景{sum(r[3] for r in prepared)}像素提亮，形狀及背景保持')


if __name__ == '__main__':
    main()
