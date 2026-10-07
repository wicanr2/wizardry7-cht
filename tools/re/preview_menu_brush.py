#!/usr/bin/env python3
"""R976：按已量主選單按鈕尺寸製作毛筆楷體字樣，不冒稱遊戲畫面。"""
import argparse
import hashlib
import json
import sys
from pathlib import Path

from font_preview_freetype import Font
from prototype_intro_paper import png_write

LABELS = ('CHARACTER MENU','ADD CHARACTER','LOAD SAVEGAME','CONFIGURATION','IMPORT (BANE)','QUIT GAME')
SOURCE = Path(__file__).resolve().parents[2]/'text/story-trans.json'
FG = (248,232,160)
BG = (46,32,18)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('font',type=Path)
    p.add_argument('license',type=Path)
    p.add_argument('output',type=Path)
    args = p.parse_args()
    if args.output.exists():
        p.error('輸出已存在，請用新目錄')
    raw = SOURCE.read_bytes()
    items = json.loads(raw)['items']
    labels = []
    for en in LABELS:
        found = [(key,r) for key,r in items.items() if r['en'].strip()==en]
        values = {r['zh'] for _,r in found}
        assert len(values)==1, (en,'選單譯文缺項或不同')
        zh = next(iter(values))
        assert zh and not any(ord(c)<32 for c in zh)
        labels.append(dict(en=en,zh=zh,source_ids=[key for key,_ in found]))
    font = Font(args.font,2)
    assert font.family=='AR PL UKai TW'
    records,prepared,glyphs = [],[],{}
    try:
        for k in (2,3):
            pixel_size = 8*k
            bw,bh,gap = 109*k,9*k,4*k
            w,h = bw+48,6*bh+5*gap+48
            image = bytearray(bytes((14,18,23))*(w*h))
            for i,label in enumerate(labels):
                x,y = 24,24+i*(bh+gap)
                for yy in range(y,y+bh):
                    for xx in range(x,x+bw):
                        edge = yy in (y,y+bh-1) or xx in (x,x+bw-1)
                        color = bytes((128,104,52)) if edge else bytes(BG)
                        at = (yy*w+xx)*3;image[at:at+3] = color
                pixels = []
                pen = 0
                for c in label['zh']:
                    key = f'{pixel_size}:{c}'
                    g = glyphs.setdefault(key,font.glyph(c,pixel_size))
                    for gy in range(g['h']):
                        for gx in range(g['w']):
                            alpha = g['alpha'][gy*g['w']+gx]
                            if alpha:
                                pixels.append((pen+g['left']+gx,-g['top']+gy,alpha))
                    pen += g['advance']+1
                xmin,xmax = min(v[0] for v in pixels),max(v[0] for v in pixels)
                ymin,ymax = min(v[1] for v in pixels),max(v[1] for v in pixels)
                width,height = xmax-xmin+1,ymax-ymin+1
                assert width<=bw-4 and height<=bh-2,(label['zh'],k,'筆畫超按鈕')
                dx = x+(bw-width)//2-xmin
                dy = y+(bh-height)//2-ymin
                for gx,gy,alpha in pixels:
                    xx,yy = dx+gx,dy+gy
                    assert x<xx<x+bw-1 and y<yy<y+bh-1
                    at = (yy*w+xx)*3
                    old = image[at:at+3]
                    image[at:at+3] = bytes((a*alpha+b*(255-alpha)+127)//255 for a,b in zip(FG,old))
                records.append(dict(**label,scale=k,pixel_size=pixel_size,button=[x,y,bw,bh],
                                    ink_bounds=[dx+xmin,dy+ymin,width,height],offset=[dx,dy],pixels=len(pixels)))
            prepared.append((f'menu-brush-{k}x.png',w,h,image))
        args.output.mkdir()
        for name,w,h,data in prepared:
            png_write(args.output/name,w,h,data)
        data = dict(schema=1,status='DRAFT font sample, not game capture',python=sys.version.split()[0],
                    family=font.family,style=font.style,face=2,freetype=font.version,
                    font_sha256=sha(args.font.read_bytes()),license_sha256=sha(args.license.read_bytes()),
                    catalog_sha256=sha(raw),tool_sha256=sha(Path(__file__).read_bytes()),
                    renderer_sha256=sha(Path(__file__).with_name('font_preview_freetype.py').read_bytes()),
                    colors=dict(foreground=list(FG),background=list(BG)),labels=records,glyphs=glyphs,
                    font_notice='2026-10-08以FreeType產生16／24像素字模樣張；來源AR PL UKai TW，衍生字模依ARPHIC PUBLIC LICENSE，不受程式RRSAL涵蓋。',
                    limits='僅六按鈕字樣。尺寸依020，背景為自製樣張，非遊戲截圖；沒有原版同狀態、GUI、互動或正式字型接線驗收。')
        with (args.output/'receipt.json').open('x') as f:
            json.dump(data,f,ensure_ascii=False,indent=2);f.write('\n')
        with (args.output/'font-copyright.txt').open('xb') as f:
            f.write(args.license.read_bytes())
        page = ('<!doctype html><html lang="zh-Hant"><meta charset="utf-8"><title>主選單毛筆字樣</title>'
                '<style>body{background:#181818;color:#eee;font:18px sans-serif;margin:24px}img{image-rendering:pixelated;margin:20px}</style>'
                '<h1>主選單六按鈕毛筆楷體字樣</h1><p>AR PL UKai TW。按既有109×9原版按鈕尺寸排2倍、3倍字樣；目前為字型樣張，不是遊戲畫面。</p>'
                '<img src="menu-brush-2x.png"><img src="menu-brush-3x.png"><p>背景為樣張。文字只取既有六則譯文，未改操作或正式字型。</p></html>')
        (args.output/'compare.html').write_text(page)
        print(f'{font.family}／FreeType{font.version}：六按鈕、兩倍率，字形缺字0，筆畫不出按鈕')
    finally:
        font.close()


if __name__ == '__main__':
    main()
