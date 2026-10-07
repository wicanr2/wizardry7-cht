#!/usr/bin/env python3
"""028：獨立核正常冷啟動原版摘要及已確認字樣的HD／Draw像素。"""
import argparse
import hashlib
import json
from pathlib import Path

from check_native_gui import png

MENU_SHA = '5d94afe72850c66b4b2526bba048b147988ccbd7a454b68fad92cf9420afd623'
PAPER_SHA = {'unifont':'afb796b46964f907e68c3a40eeb8af79d93264d3bc1426d140863ff12010cd49',
             'eten':'ea19ec75516531815b6f7a654b755e3f8ae22abef4e96c0bf4c44581487e82b9'}
FRAME = {73400000:73235759,75600000:75322912,77800000:77352859,
         101200000:100999086,103400000:103102872,105400000:105202085}
BUTTON_Y = (43,56,69,82,95,108)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def checked_json(path,digest):
    raw = path.read_bytes()
    assert sha(raw)==digest, str(path)+'定版收據不符'
    return json.loads(raw)


def scale_rgb(indexed,pal,k):
    return b''.join(b''.join(pal[i*3:i*3+3]*k for i in indexed[y*320:(y+1)*320])*k for y in range(200))


def expected_menu(comp,pal,k,receipt):
    img = bytearray(scale_rgb(comp,pal,k))
    cells = set()
    for i,en in enumerate(('CHARACTER MENU','ADD CHARACTER','LOAD SAVEGAME','CONFIGURATION','IMPORT (BANE)','QUIT GAME')):
        item = next(r for r in receipt['labels'] if r['en']==en and r['scale']==k)
        glyphs = receipt['glyphs']
        pixels,pen = [],0
        for c in item['zh']:
            g = glyphs[f'{8*k}:{c}']
            for y in range(g['h']):
                for x in range(g['w']):
                    a = g['alpha'][y*g['w']+x]
                    if a:
                        pixels.append((pen+g['left']+x,-g['top']+y,a))
            pen += g['advance']+1
        left,top = min(x for x,y,a in pixels),min(y for x,y,a in pixels)
        w,h = max(x for x,y,a in pixels)-left+1,max(y for x,y,a in pixels)-top+1
        assert [w,h]==item['ink_bounds'][2:]
        bx,by,bw,bh = 105*k,BUTTON_Y[i]*k,109*k,9*k
        origin = (bx+(bw-w)//2-left,by+(bh-h)//2-top)
        for x,y,a in pixels:
            x,y = origin[0]+x,origin[1]+y
            assert bx<x<bx+bw-1 and by<y<by+bh-1
            at = (y*320*k+x)*3
            bg = img[at:at+3]
            img[at:at+3] = bytes((fg*a+old*(255-a)+127)//255 for fg,old in zip((248,232,160),bg))
        cells.update((y*320*k+x) for y in range(by,by+bh) for x in range(bx,bx+bw))
    return bytes(img),cells


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('round',type=Path)
    ap.add_argument('reference',type=Path)
    ap.add_argument('output',type=Path)
    ap.add_argument('--cases',default='english-3,unifont-3,eten-3,unifont-2,eten-2')
    ap.add_argument('--actual',default='gui')
    args = ap.parse_args()
    menu = checked_json(args.reference/'menu/receipt.json',MENU_SHA)
    papers = {b:checked_json(args.reference/f'paper-{b}/receipt.json',s) for b,s in PAPER_SHA.items()}
    records = []
    state_checks,pixels,outside = 0,0,0
    for case in args.cases.split(','):
        backend,k = case.rsplit('-',1);k=int(k)
        before,after = args.round/'before-gui'/case,args.round/args.actual/case
        for step in (*FRAME,272000000):
            name = f'hdat-{step}'
            old,new = json.loads((before/(name+'.json')).read_text()),json.loads((after/(name+'.json')).read_text())
            for key in ('Step','Mem','MemorySHA256','Regs','A000','Pal','Mouse'):
                assert old[key]==new[key], (case,step,key,'原版改變')
            for ext in ('.a000','.pal','.comp'):
                assert (before/(name+ext)).read_bytes()==(after/(name+ext)).read_bytes(),(case,step,ext)
            state_checks += 1
            comp,pal = (after/(name+'.comp')).read_bytes(),(after/(name+'.pal')).read_bytes()
            approved = None
            if backend=='english' or step!=272000000 and (k==2 or FRAME[step]==105202085):
                expected = scale_rgb(comp,pal,k)
            elif step!=272000000:
                row = next(r for r in papers[backend]['frames'] if r['step']==FRAME[step])
                path = args.reference/f'paper-{backend}'/row['output']
                assert sha(path.read_bytes())==row['output_sha256']
                w,h,expected = png(path)
                assert (w,h)==(960,600)
            else:
                assert len(new['Shown'])==6, '未量額外選單事件，不能猜期望'
                expected,approved = expected_menu(comp,pal,k,menu)
            for path in (after/(name+'.png'),after/f'drawat-{step}.png'):
                w,h,actual = png(path)
                assert (w,h)==(320*k,200*k)
                count = sum(actual[i:i+3]!=expected[i:i+3] for i in range(0,len(actual),3))
                assert count==0,(case,step,path.name,count,'與定版全圖不同')
                pixels += w*h
            if approved is not None:
                _,_,old_rgb = png(before/(name+'.png'))
                _,_,new_rgb = png(after/(name+'.png'))
                delta = sum(old_rgb[i*3:i*3+3]!=new_rgb[i*3:i*3+3] for i in range(w*h) if i not in approved)
                assert delta==0,(case,step,'按鈕外變更',delta)
                outside += delta
            records.append(dict(case=case,step=step,original_memory_sha256=new['MemorySHA256'],
                                indexed_sha256=sha((after/(name+'.a000')).read_bytes()),
                                hd_draw_full_pixels=w*h*2,differences=0))
    with args.output.open('x') as f:
        json.dump(dict(schema=1,original_states=state_checks,full_pixels=pixels,
                       outside_menu_changes=outside,records=records,
                       oracle='紙頁固定R976已審全圖；毛筆由固定R976且已經C FreeType交叉驗證的原glyph／metrics重建；不匯入被測烘製／Go顯示函式。',
                       limits='正常無輸入冷啟動，靜音加速。紙頁2倍原圖；沒有防拷、後期玩法、音訊、人耳或公開包驗收。'),f,ensure_ascii=False,indent=2)
        f.write('\n')
    print(f'{state_checks}原版同狀態、{pixels:,} HD／Draw像素通過，選單外差異0')


if __name__=='__main__':
    main()
