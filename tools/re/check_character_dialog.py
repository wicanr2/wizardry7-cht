"""Independent native-font and original-indexed checks for R921 observed labels."""
import argparse
import gzip
import hashlib
import importlib.util
import json
import struct
from pathlib import Path

parser=argparse.ArgumentParser(description='021新增角色小按鈕原字型、原版及差異遮罩驗證')
parser.add_argument('root',type=Path)
parser.add_argument('--baseline-root',type=Path,required=True)
parser.add_argument('--unifont-hex',type=Path,required=True)
parser.add_argument('--eten-font',type=Path,required=True)
parser.add_argument('--output',type=Path,required=True)
args=parser.parse_args()
if args.output.exists():
    parser.error('輸出已存在，請指定新報告路徑')
root,out=Path(__file__).resolve().parents[2],args.root
spec = importlib.util.spec_from_file_location('pngcheck',root/'tools/re/check_native_gui.py')
decoder = importlib.util.module_from_spec(spec)
spec.loader.exec_module(decoder)
needed=set('新增角色確定取消')
raw=args.unifont_hex.read_bytes()
assert hashlib.sha256(raw).hexdigest()=='7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'
unifont={}
for line in gzip.decompress(raw).decode('ascii').splitlines():
    cp,bits=line.split(':')
    ch=chr(int(cp,16))
    if ch in needed:
        assert len(bits)==64
        unifont[ch]=[int(bits[i:i+4],16) for i in range(0,64,4)]
assert set(unifont)==needed
raw=args.eten_font.read_bytes()
assert hashlib.sha256(raw).hexdigest()=='39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
eten={}
for ch in needed:
    lead,trail=ch.encode('big5')
    column=trail-64 if trail<127 else trail-98
    record=(lead-164)*157+column-(408 if lead>=201 else 0)
    eten[ch]=[value>>1 for value in struct.unpack_from('>15H',raw,record*30)]
catalog=json.loads((root/'text/story-trans.json').read_text())['items']
for source,zh in [('ADD CHARACTER','新增角色'),('OK','確定'),('CANCEL','取消')]:
    candidates=[v['zh'] for k,v in catalog.items() if v['en']==source and not(2500<=int(k)<=2574 or 7000<=int(k)<=7197)]
    assert candidates and set(candidates)=={zh}

report={'same_state':[],'labels':[],'pixel_checks':0,'glyph_checks':0}
for scale in (2,3):
    english=args.baseline_root/f'aux-r2-english-{scale}'
    base=json.loads((english/'captures.json').read_text())
    for backend,glyphs in [('eten',eten),('unifont',unifont)]:
        dest=out/f'aux-{backend}-{scale}'
        current=json.loads((dest/'captures.json').read_text())
        for step in (272000000,400000000,528000000):
            assert current['original'][str(step)]==base['original'][str(step)]
            for ext in ('a000','pal'):
                assert (dest/f'hdat-{step}.{ext}').read_bytes()==(english/f'hdat-{step}.{ext}').read_bytes()
            report['same_state'].append(dict(backend=backend,scale=scale,step=step,original_equal=True))
        assert current['save_snapshots']['after_normal_exit']==base['save_snapshots']['after_normal_exit']
        meta=json.loads((dest/'hdat-400000000.json').read_text())
        a=(english/'hdat-400000000.a000').read_bytes()
        pal=(english/'hdat-400000000.pal').read_bytes()
        height=15 if backend=='eten' else 16
        for source,xy,zh,bounds in [('ADD CHARACTER',(61,152),'新增角色',None),
                 ('OK',(94,168),'確定',(63,166,73,9)),('CANCEL',(82,184),'取消',(63,182,73,9))]:
            event=next(e for e in meta['Shown'] if e['S']==source and (e['X0'],e['Y0'])==xy)
            assert event['Dy']==0 and event['G']==[glyphs[ch] for ch in zh]
            report['glyph_checks']+=len(zh)
            gx,gy=xy[0]*scale,xy[1]*scale
            if bounds and scale==2:
                gy -= 2 if height==15 else 3
            ink={(gx+i*(height+2)+xx,gy+yy) for i,ch in enumerate(zh)
                 for yy,row in enumerate(glyphs[ch]) for xx in range(height)
                 if row & (1<<(height-1-xx))}
            overflow=0
            if bounds:
                bx,by,bw,bh=bounds
                assert list(a[by*320+bx:by*320+bx+bw])==[81]+[176]*(bw-2)+[182]
                assert list(a[(by+bh-1)*320+bx:(by+bh-1)*320+bx+bw])==[182]+[197]*(bw-2)+[198]
                assert all(a[y*320+bx]==176 and a[y*320+bx+bw-1]==197 for y in range(by+1,by+bh-1))
                overflow=sum(not(bx*scale<=x<(bx+bw)*scale and by*scale<=y<(by+bh)*scale) for x,y in ink)
                assert overflow==0
                assert by*scale<=gy and gy+height<=(by+bh)*scale
                region=(bx*scale,by*scale,bw*scale,max((by+bh)*scale,gy+height)-by*scale)
            else:
                region=(gx,gy,max(event['BW']*scale,len(zh)*(height+2)-2),max(event['BH']*scale,height))
            rx,ry,rw,rh=region
            for prefix in ('hdat','drawat'):
                w,h,rgb=decoder.png(dest/f'{prefix}-400000000.png')
                assert (w,h)==(320*scale,200*scale)
                for y in range(ry,ry+rh):
                    for x in range(rx,rx+rw):
                        ox,oy=x//scale,y//scale
                        index=a[oy*320+ox]
                        if (event['BX']<=ox<event['BX']+event['BW'] and
                            event['BY']<=oy<event['BY']+event['BH'] and index==event['Fg']):
                            index=192 if bounds else 194
                        if (x,y) in ink:
                            index=event['Fg']
                        expected=pal[index*3:index*3+3]
                        at=(y*w+x)*3
                        assert rgb[at:at+3]==expected,(backend,scale,source,prefix,x,y)
                        report['pixel_checks']+=1
            report['labels'].append(dict(backend=backend,scale=scale,source=source,xy=xy,
                independent_glyphs=len(zh),button_xywh=bounds,glyph_height=height,
                geometric_overflow_rows=max(0,gy+height-(bounds[1]+bounds[3])*scale) if bounds else 0,
                ink_pixels_outside_button=overflow))
report['before_after'] = []
for backend in ('eten','unifont'):
    height = 15 if backend=='eten' else 16
    for scale in (2,3):
        dest=out/f'aux-{backend}-{scale}'
        old=args.baseline_root/f'aux-r2-{backend}-{scale}'
        for step in (272000000,400000000,528000000):
            for ext in ('a000','pal','comp'):
                assert (dest/f'hdat-{step}.{ext}').read_bytes()==(old/f'hdat-{step}.{ext}').read_bytes()
            new_meta=json.loads((dest/f'hdat-{step}.json').read_text())
            old_meta=json.loads((old/f'hdat-{step}.json').read_text())
            assert new_meta['Shown']==old_meta['Shown']
            w,h,rgb=decoder.png(dest/f'hdat-{step}.png')
            ow,oh,before=decoder.png(old/f'hdat-{step}.png')
            assert (w,h)==(ow,oh)
            mask=set()
            if step==400000000 and scale==2:
                for x0,y0 in ((94,168),(82,184)):
                    width=2*(height+2)-2
                    old_y=y0*scale
                    new_y=old_y-(2 if height==15 else 3)
                    mask.update((x,y) for start in (old_y,new_y) for y in range(start,start+height) for x in range(x0*scale,x0*scale+width))
            changed=outside=0
            for y in range(h):
                for x in range(w):
                    i=(y*w+x)*3
                    if rgb[i:i+3]!=before[i:i+3]:
                        changed+=1
                        outside+=(x,y) not in mask
            assert outside==0
            assert (changed>0)==(step==400000000 and scale==2)
            report['before_after'].append(dict(backend=backend,scale=scale,step=step,changed_pixels=changed,outside_approved_regions=outside,events_and_comp_equal=True))

with args.output.open('x',encoding='utf-8') as f:
    f.write(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
print(json.dumps(report,ensure_ascii=False))
