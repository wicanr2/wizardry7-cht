#!/usr/bin/env python3
"""028 READY：烘製已確認字樣；只在離線Docker內執行。"""
import argparse
import base64
import hashlib
import json
import os
from pathlib import Path

from font_preview_freetype import Font

ROOT = Path(__file__).resolve().parents[2]
FONT_SHA = 'cfc758e69386a21bade8cb3b373fda20d7f2fac6e8904ced2ded98588b63e274'
LICENSE_SHA = '1f11c4a22865e1289f4694dad56b3561790a4b7341a511ddcd0edddcc3b95a46'
MODEL_SHA = '1a9632b661b106367adffc2a8336808fe38202e60080620e897072e2fdcf5a22'
MENU = [('CHARACTER MENU','角色選單'),('ADD CHARACTER','新增角色'),
        ('LOAD SAVEGAME','載入存檔'),('CONFIGURATION','系統設定'),
        ('IMPORT (BANE)','匯入巫術六'),('QUIT GAME','結束遊戲')]
PAPER = [('presents','呈獻'),('fantasy','奇幻'),('role-playing','角色扮演'),
         ('simulation','遊戲'),('by','作者')]


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def catalog(name, context, expected, source):
    raw = (ROOT/'text'/name).read_bytes()
    lines = raw.decode('utf-8').splitlines()
    assert lines[0].split('\t') == ['en','zh','context','dy','source']
    result = []
    for line in lines[1:]:
        en,zh,ctx,dy,src = line.split('\t')
        assert ctx == context and dy == '0' and src == source
        assert en and zh and all(ord(c)>=32 for c in en+zh)
        result.append((en,zh))
    assert result == expected and len(set(en for en,_ in result)) == len(result)
    return result,sha(raw)


def build_menu(font_path, license_path):
    assert sha(font_path.read_bytes()) == FONT_SHA
    assert sha(license_path.read_bytes()) == LICENSE_SHA
    labels,source_sha = catalog('main-menu.zh-TW.tsv','main.menu.brush',MENU,'runtime:0110:3895')
    font = Font(font_path,2)
    assert font.family == 'AR PL UKai TW' and font.version == '2.12.1'
    entries = []
    try:
        for en,zh in labels:
            stripes = []
            for scale in (2,3):
                pixels,pen = {},0
                for ch in zh:
                    g = font.glyph(ch,8*scale)
                    for y in range(g['h']):
                        for x in range(g['w']):
                            a = g['alpha'][y*g['w']+x]
                            if a:
                                key = (pen+g['left']+x,-g['top']+y)
                                assert key not in pixels, '字模重疊需另核合成契約'
                                pixels[key] = a
                    pen += g['advance']+1
                left,top = min(x for x,y in pixels),min(y for x,y in pixels)
                w,h = max(x for x,y in pixels)-left+1,max(y for x,y in pixels)-top+1
                assert 0<w<=109*scale-4 and 0<h<=9*scale-2
                alpha = bytearray(w*h)
                for (x,y),a in pixels.items():
                    alpha[(y-top)*w+x-left] = a
                stripes.append(dict(scale=scale,w=w,h=h,a=base64.b64encode(alpha).decode('ascii')))
            entries.append(dict(en=en,zh=zh,stripes=stripes))
    finally:
        font.close()
    return dict(schema=1,context='main.menu.brush',font_sha256=FONT_SHA,
                license_sha256=LICENSE_SHA,source_sha256=source_sha,entries=entries)


def build_paper(eten_path):
    import bake_font
    from font_source import UNIFONT_SHA256
    labels,source_sha = catalog('intro-paper.zh-TW.tsv','intro.mon61.paper',PAPER,'runtime:0110:3921')
    bake_font.FONT15 = eten_path
    size = bake_font.HEIGHT
    expected_font = UNIFONT_SHA256 if size == 16 else '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
    if size == 15:
        assert sha(eten_path.read_bytes()) == expected_font
    entries = []
    for en,zh in labels:
        glyphs = [bake_font.glyph_char(c) for c in zh]
        assert all(len(g)==size and any(g) and all(0<=v<1<<size for v in g) for g in glyphs)
        entries.append(dict(en=en,zh=zh,dy=0,g=glyphs))
    return dict(schema=1,context='intro.mon61.paper',size=size,source_sha256=source_sha,
                model_sha256=MODEL_SHA,entries=entries)


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('xlate',type=Path)
    ap.add_argument('--font',required=True,type=Path)
    ap.add_argument('--license',required=True,type=Path)
    ap.add_argument('--eten-font',required=True,type=Path)
    ap.add_argument('--go-model',type=Path)
    args = ap.parse_args()
    assert args.xlate.is_file() and not args.xlate.is_symlink() and args.xlate.stat().st_uid == os.getuid()
    for name in ('xlate-legacy.json','brush-font-copyright.txt','approved-ui-receipt.json'):
        path = args.xlate.with_name(name)
        assert not path.exists() and not path.is_symlink(), '輸出已存在，請使用新目錄'
    if args.go_model:
        assert not args.go_model.is_symlink()
        if args.go_model.exists():
            assert args.go_model.is_file() and args.go_model.stat().st_uid == os.getuid()
    model = (ROOT/'tools/re/intro_paper_candidate.json').read_bytes()
    assert sha(model) == MODEL_SHA
    old = args.xlate.read_bytes()
    data = json.loads(old)
    assert '__main_menu_brush__' not in data and '__intro_paper__' not in data
    menu,paper = build_menu(args.font,args.license),build_paper(args.eten_font)
    for item in (menu,paper):
        assert len(json.dumps(item,ensure_ascii=False).encode()) <= 1<<20
    data.update(__main_menu_brush__=menu,__intro_paper__=paper)
    backup = args.xlate.with_name('xlate-legacy.json')
    with backup.open('xb') as f:
        f.write(old)
    args.xlate.write_text(json.dumps(data,ensure_ascii=False)+'\n')
    notice = args.xlate.with_name('brush-font-copyright.txt')
    with notice.open('xb') as f:
        f.write(args.license.read_bytes())
    receipt = dict(schema=1,legacy_sha256=sha(old),xlate_sha256=sha(args.xlate.read_bytes()),
                   model_sha256=MODEL_SHA,menu_source_sha256=menu['source_sha256'],
                   paper_source_sha256=paper['source_sha256'],size=paper['size'],
                   font_sha256=FONT_SHA,license_sha256=LICENSE_SHA,
                   notice='2026-10-08以FreeType產生16／24像素字模，來源AR PL UKai TW，衍生字模依ARPHIC PUBLIC LICENSE。')
    with args.xlate.with_name('approved-ui-receipt.json').open('x') as f:
        json.dump(receipt,f,ensure_ascii=False,indent=2);f.write('\n')
    if args.go_model:
        assert args.go_model.parent.is_dir() and args.go_model.parent.stat().st_uid == os.getuid()
        code = '// 由 tools/re/bake_approved_ui.py 固定模型重建；不含原版像素。\npackage main\n\nconst introPaperModelJSON = '+json.dumps(model.decode(),ensure_ascii=False)+'\n'
        args.go_model.write_text(code)
    print('六按鈕×兩倍率、紙頁五鍵已烘製；一般字形保持，授權另存。')


if __name__ == '__main__':
    main()
