#!/usr/bin/env python3
"""spec007 真實 GUI 的獨立字形／說明排版核對；只在 Docker 執行。

F1 內容直接取正式 TSV；Unifont 直接解析固定 HEX，不匯入前端或烘製程式。
倚天字形採 r682 已確認等同舊版的回歸資料，並非重新解析原始倚天字型。
PNG 使用標準庫解碼 palette/gray/RGB/RGBA。所有圖片與詳細收據只留本機。
"""
import argparse
import csv
import gzip
import hashlib
import json
import math
import re
import struct
import zlib
from collections import Counter
from pathlib import Path

FONT_SHA = '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
COLORS = ((255, 221, 68), (240, 240, 240), (160, 160, 160))
ORDINALS = dict(zip('FIRST SECOND THIRD FOURTH FIFTH SIXTH SEVENTH EIGHTH NINTH TENTH'.split(), '一二三四五六七八九十'))
DIGITS = {'0': (14,17,19,21,25,17,14), '1': (4,12,4,4,4,4,14),
          '2': (14,17,1,2,4,8,31), '3': (30,1,1,14,1,1,30),
          '4': (2,6,10,18,31,2,2), '5': (31,16,30,1,1,17,14),
          '6': (6,8,16,30,17,17,14), '7': (31,1,2,4,8,8,8),
          '8': (14,17,17,14,17,17,14), '9': (14,17,17,15,1,2,12)}


def load(p):
    return json.loads(Path(p).read_text())


def png(p):
    """完整檢查 CRC、scanline filters；拒絕交錯圖，保留 RGB 真值。"""
    data = Path(p).read_bytes()
    if data[:8] != b'\x89PNG\r\n\x1a\n':
        raise ValueError('不是 PNG')
    pos, compressed, palette = 8, bytearray(), None
    while pos < len(data):
        n = struct.unpack('>I', data[pos:pos+4])[0]
        kind, body = data[pos+4:pos+8], data[pos+8:pos+8+n]
        if zlib.crc32(kind + body) & 0xffffffff != struct.unpack('>I', data[pos+8+n:pos+12+n])[0]:
            raise ValueError('PNG CRC 不符')
        if kind == b'IHDR':
            w, h, depth, typ, compression, filtering, interlace = struct.unpack('>IIBBBBB', body)
        elif kind == b'PLTE':
            palette = [body[i:i+3] for i in range(0, len(body), 3)]
        elif kind == b'IDAT':
            compressed.extend(body)
        pos += n + 12
        if kind == b'IEND':
            break
    if interlace or compression or filtering or typ not in (0,2,3,4,6):
        raise ValueError('PNG 未支援格式')
    channels = {0:1, 2:3, 3:1, 4:2, 6:4}[typ]
    stride, bpp = math.ceil(w*channels*depth/8), max(1, math.ceil(channels*depth/8))
    raw = zlib.decompress(compressed)
    if len(raw) != (stride+1)*h:
        raise ValueError('PNG 展開長度不符')
    previous = bytearray(stride)
    rgb = bytearray()
    for y in range(h):
        f = raw[y*(stride+1)]
        row = bytearray(raw[y*(stride+1)+1:(y+1)*(stride+1)])
        for i in range(stride):
            a, b, c = (row[i-bpp] if i>=bpp else 0), previous[i], (previous[i-bpp] if i>=bpp else 0)
            if f == 1: q = a
            elif f == 2: q = b
            elif f == 3: q = (a+b)//2
            elif f == 4:
                t = a+b-c
                distances = (abs(t-a), abs(t-b), abs(t-c))
                q = (a,b,c)[distances.index(min(distances))]
            elif f == 0: q = 0
            else: raise ValueError('PNG filter 不合法')
            row[i] = (row[i]+q)&255
        previous = row
        if depth < 8:
            if typ not in (0,3): raise ValueError('低位深 PNG 色型不合法')
            mask = (1<<depth)-1
            samples = [(row[i*depth//8] >> (8-depth-(i*depth%8))) & mask for i in range(w)]
        elif depth == 8:
            samples = list(row)
        elif depth == 16:
            samples = list(row[::2])
        else: raise ValueError('PNG 位深未支援')
        if typ == 3:
            rgb.extend(b''.join(palette[i] for i in samples))
        elif typ in (0,4):
            values = samples[::channels]
            if depth < 8: values = [v*255//((1<<depth)-1) for v in values]
            rgb.extend(b''.join(bytes((v,v,v)) for v in values))
        elif typ == 2:
            rgb.extend(samples)
        else:
            rgb.extend(b''.join(bytes(samples[i:i+3]) for i in range(0,len(samples),4)))
    return w, h, bytes(rgb)


class Font:
    def __init__(self, path, eten):
        raw = Path(path).read_bytes()
        if hashlib.sha256(raw).hexdigest() != FONT_SHA:
            raise ValueError('固定 HEX 雜湊不符')
        self.raw = {}
        for line in gzip.decompress(raw).decode('ascii').splitlines():
            cp, bits = line.split(':')
            if len(bits) not in (32,64): raise ValueError('HEX 字形格式不符')
            width = len(bits)//4
            self.raw[chr(int(cp,16))] = (width, tuple(int(bits[i:i+width//4],16) for i in range(0,len(bits),width//4)))
        self.eten = load(eten)['glyphs']

    def help(self, backend, ch):
        if backend == 'unifont': return self.raw[ch]
        g = self.eten[ch]
        # r682 回歸資料已完成倚天原位圖右移，不能再次位移。
        return g['w'], tuple(g['rows'])

    def square(self, ch):
        w, rows = self.raw[ch]
        return [v << 4 if w == 8 else v for v in rows]


def catalog(root):
    pages = {}
    for name in ('help.zh-TW.tsv','help-names.zh-TW.tsv'):
        with (root/'text'/name).open() as f:
            for row in csv.DictReader(f, delimiter='\t'):
                key, seq = row['key'].rsplit('.',1)
                pages.setdefault(key,{})[int(seq)] = row['zh']
    return pages


def layout(pages, backend, scale, screen, font):
    w,h = 320*scale, 200*scale
    capacity = max(1,(h-64-48)//20)
    result = []
    order = ['keys','ops','screen.'+screen]+[k for k in pages if k.startswith('names.')]
    def wrap(s, maximum):
        output, line, used = [], '', 0
        for ch in s:
            width = font.help(backend,ch)[0]
            if line and used+2+width > maximum:
                output.append(line); line,used = '',0
            if width > maximum: raise ValueError('單字超寬')
            used += (2 if line else 0)+width
            line += ch
        return output+[line]
    for ident in order:
        base = pages[ident]
        lines = [x for seq,s in sorted(base.items()) if seq for x in wrap(s,w-64)]
        columns = 1
        if len(lines)>capacity:
            columns = 2
            lines = [x for seq,s in sorted(base.items()) if seq for x in wrap(s,w//2-64)]
        for i in range(0,max(1,len(lines)),capacity*columns):
            result.append({'id':ident,'title':base[0],'lines':lines[i:i+capacity*columns],'columns':columns})
    return result


def ink(font,backend,text,x,y,w,h):
    points, boxes = set(), []
    for ch in text:
        width, rows = font.help(backend,ch)
        boxes.append((x,y,width,len(rows)))
        for iy,n in enumerate(rows):
            for ix in range(width):
                if n & (1<<(width-1-ix)):
                    if not 0<=x+ix<w or not 0<=y+iy<h: raise ValueError('字形超出畫布')
                    points.add((y+iy)*w+x+ix)
        x += width+2
    return points,boxes


def help_masks(page, n, total, backend, scale, pages, font):
    w,h = scale*320,scale*200
    masks = [set(),set(),set()]
    bounds = []
    masks[0],boxes = ink(font,backend,page['title'],32,24,w,h)
    bounds += [(b,32,w-32) for b in boxes]
    capacity = max(1,(h-64-48)//20)
    for i,s in enumerate(page['lines']):
        col,row = i//capacity,i%capacity
        x,y = 32+col*(w//2),64+row*20
        pts,boxes = ink(font,backend,s,x,y,w,h)
        masks[1] |= pts
        bounds += [(b,x,x+(w-64 if page['columns']==1 else w//2-64)) for b in boxes]
    footer = pages['ui'][1].replace('{n}',str(n)).replace('{N}',str(total))
    masks[2],boxes = ink(font,backend,footer,32,h-32,w,h)
    bounds += [(b,32,w-32) for b in boxes]
    safe = all(lo<=x and x+ww<=hi and y+hh<=h for (x,y,ww,hh),lo,hi in bounds)
    return masks,safe


def normalize(s):
    s = re.sub(r'[\x00-\x1f]','',s)
    def tail(m):
        for word in ORDINALS:
            if m[0].upper().endswith(word): return m[0][:-len(word)]+'$'
        return m[0]
    return re.sub(r'\d+','$',re.sub(r'[A-Za-z]+',tail,s)).replace('^','$').replace('%','$')


def xlate(root,font):
    t = load(root/'text/story-trans.json')['items']
    pairs = [(v['en'],v['zh'],v.get('dy',0)) for i,v in t.items() if not 2500<=int(i)<=2574 and not 7000<=int(i)<=7197]
    for name in ('monsters-trans.json','items-trans.json'):
        pairs += [(en,zh,0) for en,zh in load(root/'text'/name)['items'].items()]
    result = {}
    for en,zh,dy in pairs:
        if zh.isascii() or not normalize(en).replace('$','').strip(): continue
        g = [None if ch=='$' else font.square(ch) for ch in zh if ch=='$' or len(ch.encode('big5'))==2]
        result[normalize(en)] = (g,dy)
    return result


def digits(ch):
    out=[0]
    for n in DIGITS[ch]:
        row=sum(3<<(11-2*x) for x in range(5) if n&(16>>x))
        out.extend((row,row))
    return out+[0]


def fnv(raw):
    h=2166136261
    for b in raw: h=((h^b)*16777619)&0xffffffff
    return h


def story_rows(text):
    """spec006 避頭尾及 spec007 的16px、18px字距；3×故事列寬24字。"""
    units=re.findall('……|——|.',text)
    result=[]
    while units:
        j=0
        while j<len(units) and sum(len(t) for t in units[:j+1])<=24:
            j+=1
        for _ in range(3):
            if j<=1: break
            if (j<len(units) and units[j][0] in '，。、！？；：）」』》…—') or units[j-1][0] in '（「『《':
                j-=1
            else: break
        line=''.join(units[:j]); units=units[j:]
        result.append({'Text':line,'X0':480-len(line)*36//2,'Y0':(146+12*len(result))*3})
    return result


def bottom_statistics(root, gui, pages, font):
    """只統計已呈現內容；不輸出任何英文或保護字表。"""
    def han(items):
        counts=Counter(ch for text in items for ch in text if '\u3400'<=ch<='\u9fff' and font.square(ch)[-1])
        return [{'char':ch,'occurrences':count} for ch,count in sorted(counts.items())]
    texts=[]
    for scale in (2,3):
        ls=layout(pages,'unifont',scale,'copyprot',font)
        for n,p in enumerate(ls):
            texts += [p['title']]+p['lines']+[pages['ui'][1].replace('{n}',str(n+1)).replace('{N}',str(len(ls)))]
    story=[]
    for step in (363113713,395347218):
        meta=load(root/'workplace/r687/zh'/f'hdat-{step}.json')
        for f in meta.get('Flow') or []:
            if all(not 2500<=i<=2574 and not 7000<=i<=7197 for i in f['Lines']):
                story.append({'requested_step':step,'nonzero_bottom_row_han':han([f['Zh']]),'basis':'真實Draw與獨立正式譯文／完整HEX32px像素核對；以同請求HD來源記錄統計。'})
    narr={}
    receipt_count=0
    for path in sorted((root/'workplace/r687/zh').glob('*.json')):
        meta=load(path)
        for f in meta.get('Flow') or []:
            if f['Site']!=0xBCFC or not all(not 2500<=i<=2574 and not 7000<=i<=7197 for i in f['Lines']): continue
            if not path.with_suffix('.png').exists(): raise ValueError('旁白統計缺少實際PNG')
            receipt_count+=1
            narr[(tuple(f['Lines']),f['Zh'])]=f['Zh']
    return {
        'f1':{'backend':'unifont','display_pages_by_scale':{'2':12,'3':10},'nonzero_bottom_row_han':han(texts),'ascii_nonzero_bottom_row_codepoints':sorted(ord(ch) for ch in set(''.join(texts)) if ord(ch)<128 and font.help('unifont',ch)[1][-1]),'basis':'全22個正常鍵盤顯示頁的完整PNG筆畫已比對；循環與重開不重複計字。'},
        'story_draw':story,
        'narrator_hd':{'visible_receipts':receipt_count,'distinct_visible_page_texts':len(narr),'nonzero_bottom_row_han':han(narr.values()),'basis':'r687已驗收HD實際可見Flow；各不同Lines＋Zh僅計一次，不把未觸發字或全字庫104字當成已呈現。'},
        'general_overlay':{'events_per_receipt':13,'glyphs_per_receipt':75,'nonzero_bottom_row_glyphs':0,'nonzero_right_edge_glyphs':51,'basis':'完整HEX16列資料與PNG相同；本次第16欄非零已量，第16列非零未由一般覆繪樣本量到。'}}


def main():
    ap=argparse.ArgumentParser()
    ap.add_argument('--root',type=Path,default=Path('.'))
    ap.add_argument('--gui',type=Path,default=Path('workplace/r694'))
    ap.add_argument('--overlay',type=Path,default=Path('workplace/r691'))
    ap.add_argument('--out',type=Path,default=Path('workplace/r690'))
    args=ap.parse_args(); args.out.mkdir(exist_ok=False)
    root=args.root
    font=Font(root/'workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz',root/'workplace/r682/eten/help.json')
    pages=catalog(root); report={'inputs':{'gui':str(args.gui),'overlay':str(args.overlay)},'checks':0,'failed':0,'failures':[],'help':[],'general':[],'draw':[],'limitations':[]}
    def check(ok,label):
        report['checks']+=1
        if not ok: report['failed']+=1; report['failures'].append(label); print('失敗',label)
    def observed(p):
        w,h,rgb=png(p); sets=[set(),set(),set()]; colors={bytes(c):i for i,c in enumerate(COLORS)}
        for j in range(w*h):
            q=colors.get(rgb[j*3:j*3+3])
            if q is not None: sets[q].add(j)
        return w,h,rgb,sets
    for backend in ('eten','unifont'):
        for scale in (2,3):
            dest=args.gui/f'{backend}-{scale}'
            ls=layout(pages,backend,scale,'copyprot',font)
            all_masks=[help_masks(page,n+1,len(ls),backend,scale,pages,font)[0] for n,page in enumerate(ls)]
            check(len(ls)==(12 if scale==2 else 10),f'{backend}/{scale} 說明完整顯示頁數')
            check(len(list(dest.glob('help-*.png')))==len(ls)+1,f'{backend}/{scale} 全頁與循環收據數')
            for n in range(len(ls)+1):
                index=n%len(ls); expected,safe=help_masks(ls[index],index+1,len(ls),backend,scale,pages,font)
                w,h,rgb,actual=observed(dest/f'help-{n:02}.png')
                check((w,h)==(320*scale,200*scale),f'{backend}/{scale}/{n} 畫布')
                check(safe,f'{backend}/{scale}/{n} 全字格欄界')
                for k,label in enumerate(('標題','正文','頁尾')):
                    check(actual[k]==expected[k],f'{backend}/{scale}/{n} {label} 全部筆畫與空白區（差 {len(actual[k]^expected[k])}）')
                matches=[i for i,masks in enumerate(all_masks) if actual==masks]
                body_matches=[i for i,masks in enumerate(all_masks) if actual[:2]==masks[:2]]
                report['help'].append({'backend':backend,'scale':scale,'page':n,'id':ls[index]['id'],'ink_difference':[len(a^b) for a,b in zip(actual,expected)],'column_safe':safe,'observed_complete_page_indices':matches,'observed_title_body_page_indices':body_matches,'no_help_ink':not any(actual)})
            _,_,_,reopened=observed(dest/'f1-reopen.png')
            matches=[]
            for screen in ('copyprot','general','mainmenu','addchar'):
                choices=layout(pages,backend,scale,screen,font)
                for n,page in enumerate(choices):
                    masks,_=help_masks(page,n+1,len(choices),backend,scale,pages,font)
                    if masks==reopened: matches.append((screen,n,page['id']))
            check(bool(matches),f'{backend}/{scale} 真實F1重開為完整合法說明頁')
            report['help'].append({'backend':backend,'scale':scale,'f1_reopen_matches':matches})
    # 300M 本批沒有字形事件，僅比原版狀態，不冒稱一般覆繪驗收。
    snapshots=[load(args.gui/f'{b}-{s}/hdat-300000000.json') for b in ('eten','unifont') for s in (2,3)]
    check(all(all(a[k]==snapshots[0][k] for k in ('Step','Mem','A000','Pal')) for a in snapshots),'GUI 300M 四組原版雜湊／步數一致')
    report['limitations'].append(f'{args.gui}的300M實際301M，Shown與Flow皆0，只驗F1之後原版狀態；一般覆繪另用{args.overlay}。')
    expected_xlate=xlate(root,font)
    for scale in (2,3):
        for step in (272000000,286000000):
            name=f'hdat-{step}'
            paths=[args.overlay/f'{b}-{scale}'/name for b in ('eten','unifont')]
            if not all(p.with_suffix('.json').exists() for p in paths):
                check(False,f'一般覆繪{scale}/{step} 收據尚未齊全'); continue
            old,new=[load(p.with_suffix('.json')) for p in paths]
            check(all(old[k]==new[k] for k in ('Step','Mem','A000','Pal')),f'一般覆繪{scale}/{step} 原版狀態相同')
            a=paths[1].with_suffix('.a000').read_bytes(); comp=paths[1].with_suffix('.comp').read_bytes()
            pal=paths[1].with_suffix('.pal').read_bytes()
            check(len(pal)==768 and fnv(pal)==new['Pal'],f'一般覆繪{scale}/{step} raw色盤與雜湊')
            check(a==paths[0].with_suffix('.a000').read_bytes() and fnv(a)==new['A000'],f'一般覆繪{scale}/{step} raw A000與雜湊')
            check(comp==paths[0].with_suffix('.comp').read_bytes(),f'一般覆繪{scale}/{step} 原始英文還原與倚天基線相同')
            # 收據 .pal 是原版 renderer 已展開的 8bit RGB，不是 DAC 暫存器。
            rgbpal=[pal[i:i+3] for i in range(0,768,3)]
            w,h,rgb=png(paths[1].with_suffix('.png'))
            expected=bytearray(b''.join(rgbpal[comp[(y//scale)*320+x//scale]] for y in range(h) for x in range(w)))
            shown=new.get('Shown') or []; check(bool(shown),f'一般覆繪{scale}/{step} 實際可見事件非零')
            allowed=set(); glyph_count=0; bottom_count=0; right_count=0
            for n,e in enumerate(shown):
                entry=expected_xlate.get(normalize(e['S']))
                check(entry is not None,f'一般覆繪{scale}/{step}/{n} 英文命中正式來源')
                if entry is None: continue
                gs,dy=entry; slots=re.findall(r'\d+|\b(?:'+ '|'.join(ORDINALS)+r')\b',e['S'],flags=re.I)
                filled=[]; slot=iter(slots)
                for g in gs:
                    if g is not None: filled.append(g); continue
                    value=next(slot,'')
                    if value.upper() in ORDINALS: filled.append(font.square(ORDINALS[value.upper()]))
                    else: filled.extend(digits(ch) for ch in value)
                check(filled==e['G'] and dy==e['Dy'],f'一般覆繪{scale}/{step}/{n} 完整HEX16列／獨立自製數字槽')
                x0,y0=e['X0']*scale,(e['Y0']+dy)*scale
                fg=rgbpal[e['Fg']]; glyph_count+=len(filled)
                bottom_count+=sum(bool(g[-1]) for g in filled)
                right_count+=sum(any(row&1 for row in g) for g in filled)
                for j,g in enumerate(filled):
                    for yy,row in enumerate(g):
                        for xx in range(16):
                            x,y=x0+j*18+xx,y0+yy
                            if 0<=x<w and 0<=y<h:
                                allowed.add(y*w+x)
                                if row & (1<<(15-xx)): expected[(y*w+x)*3:(y*w+x)*3+3]=fg
                            else: check(False,f'一般覆繪{scale}/{step}/{n} 字格超出畫布')
            diff=[i for i in range(w*h) if rgb[i*3:i*3+3]!=expected[i*3:i*3+3]]
            check(not diff,f'一般覆繪{scale}/{step} 整張PNG、完整筆畫與批准字格外像素（差{len(diff)}）')
            report['general'].append({'scale':scale,'requested_step':step,'actual_step':new['Step'],'events':len(shown),'glyphs':glyph_count,'bottom_row_ink_glyphs':bottom_count,'right_edge_ink_glyphs':right_count,'pixel_difference':len(diff),'outside_difference':sum(i not in allowed for i in diff)})
    # Draw 與 HD 相同請求點逐項回報；沒有 Draw 的原版步數收據，不冒稱同狀態。
    for step in (363113713,395347218):
        draw=args.gui/'unifont-story'/f'drawat-{step}.png'
        hd=root/'workplace/r687/zh'/f'hdat-{step}.png'
        if not hd.exists():
            report['draw'].append({'requested_step':step,'classification':'同點HD收據尚缺，未量'}); continue
        w,h,a=png(draw); ww,hh,b=png(hd)
        check((w,h)==(ww,hh)==(960,600),f'Draw故事{step} 畫布')
        changed=[i for i in range(w*h) if a[i*3:i*3+3]!=b[i*3:i*3+3]]
        band={y*w+x for y in range(438,585) for x in range(36,924)}
        inside=sum(i in band for i in changed)
        meta=load(hd.with_suffix('.json')); f=meta['Flow'][0]
        t=load(root/'text/story-trans.json')['items']
        def clean(s):
            s=re.sub(r'[\x00-\x1f]','',s).lstrip('@#')
            return re.sub(r'\s*(?:&[%#|\[]*)+\]$','',s)
        check(all(not 2500<=i<=2574 and not 7000<=i<=7197 for i in f['Lines']),f'Draw故事{step} 來源不含排除範圍')
        text=''.join(clean(t[str(i)]['zh']) for i in f['Lines'])
        rows=story_rows(text)
        check(text==f['Zh'] and rows==f['Rows'],f'Draw故事{step} 正式譯文與獨立16px排版')
        rawpal=hd.with_suffix('.pal').read_bytes(); fg=rawpal[3*0x55:3*0x55+3]
        check(fnv(rawpal)==meta['Pal'] and fg==bytes(f['FgRGB']),f'Draw故事{step} 原版色盤筆畫色')
        want=set()
        for row in rows:
            for n,ch in enumerate(row['Text']):
                for yy,bits in enumerate(font.square(ch)):
                    for xx in range(16):
                        if bits&(1<<(15-xx)):
                            want|={(row['Y0']+yy*2+dy)*w+row['X0']+n*36+xx*2+dx for dy in range(2) for dx in range(2)}
        for label,rgb in [('Draw',a),('HD',b)]:
            actual={i for i in band if rgb[i*3:i*3+3]==fg}
            check(actual==want,f'{label}故事{step} 正式譯文完整HEX32px筆畫及空白（差{len(actual^want)}）')
        check(inside==0,f'Draw故事{step} 故事文字帶與同請求點HD相同')
        box=None if not changed else [min(i%w for i in changed),min(i//w for i in changed),max(i%w for i in changed),max(i//w for i in changed)]
        report['draw'].append({'requested_step':step,'hd_actual_step':meta['Step'],'bottom_row_ink_characters':sum(bool(font.square(c)[-1]) for c in set(text)),'pixel_difference':len(changed),'story_band_difference':inside,'difference_bbox':box,'classification':'圖片完全相同；Draw沒有實際Step收據，不另宣稱同狀態' if not changed else '差異位於原版上方圖像，故事帶與完整HEX筆畫相同；正常Draw非精確快照，未挑選其他點替換'})
    report['limitations'] += ['倚天F1字形取已獨立確認等於舊烘製的r682回歸資料，不是raw Eten新解析。','F1固定標題y24、頁尾yH-32及三種顏色是既有繪製介面；正文折行、字格、容量與分頁獨立依spec007及TSV計算。','這份驗收不證明全遊戲已中文化、完整通關、音訊或封包驗收。']
    visible={ch for p in layout(pages,'unifont',3,'copyprot',font) for s in [p['title']]+p['lines'] for ch in s}
    report['help_bottom_row_ink_codepoints']=sorted(ord(ch) for ch in visible if font.help('unifont',ch)[1][-1])
    report['limitations'].append('一般覆繪抽樣的第16列均為零，底列實際筆畫驗證由F1的5個碼點及故事第3頁含底列筆畫的字補足；不把空列誤當已量到非零底列。')
    report['actual_glyph_statistics']=bottom_statistics(root,args.gui,pages,font)
    report['source_hashes']={str(p):hashlib.sha256(p.read_bytes()).hexdigest() for p in [root/'text/help.zh-TW.tsv',root/'text/help-names.zh-TW.tsv',root/'workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz',root/'workplace/r682/eten/help.json']}
    (args.out/'report.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
    print(f"核對 {report['checks']} 項，失敗 {report['failed']} 項")
    return int(bool(report['failed']))


if __name__=='__main__':
    raise SystemExit(main())
