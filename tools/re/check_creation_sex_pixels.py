"""R1003獨立字型／原版字格oracle；不匯入候選或烘字程式。"""
from pathlib import Path
import argparse,base64,gzip,hashlib,json
from PIL import Image,ImageChops

p=argparse.ArgumentParser();p.add_argument('--root',type=Path,required=True);p.add_argument('--out',type=Path,required=True)
p.add_argument('--eten',type=Path,required=True);p.add_argument('--unifont',type=Path,required=True)
p.add_argument('--font',choices=['eten','unifont']);p.add_argument('--scale',type=int,choices=[2,3]);a=p.parse_args()
sha=lambda b:hashlib.sha256(b).hexdigest();root=a.root;failures=[];checks=0;inputs={};pixels=0
def check(condition,label):
    global checks
    checks+=1
    if not condition:failures.append(label)
def read(path):
    data=path.read_bytes();inputs[str(path)]=sha(data);return data
def load(path):return json.loads(read(path))
def image(path):
    data=read(path)
    import io
    with Image.open(io.BytesIO(data)) as im:im.load();return im.convert('RGB')
def equal(left,right,label):
    global pixels
    check(left.size==right.size,label+'尺寸')
    if left.size==right.size:
        pixels+=left.width*left.height
        check(ImageChops.difference(left,right).getbbox() is None,label+'全部RGB像素')
fonts={}
for backend,path,digest,size in [('eten',a.eten,'39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701',15),('unifont',a.unifont,'7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574',16)]:
    data=read(path);assert sha(data)==digest
    if backend=='eten':
        assert len(data)==392820
        def glyph(ch):
            hi,lo=ch.encode('big5');index=(hi-0xa4)*157+lo-(0x40 if lo<=0x7e else 0x62)-(408 if hi>=0xc9 else 0)
            b=data[index*30:index*30+30];assert len(b)==30
            return [int.from_bytes(b[y*2:y*2+2],'big')>>1 for y in range(15)]
        positive=glyph('一');assert positive==[0]*6+[2]+[32767]+[0]*7
        rows=glyph('男')
    else:
        lines=gzip.decompress(data).decode('ascii').splitlines();matches=[s.split(':')[1] for s in lines if s.startswith('7537:')]
        assert len(matches)==1 and len(matches[0])==64
        rows=[int(matches[0][i:i+4],16) for i in range(0,64,4)]
    fonts[backend]=(size,rows)
event_path=root.parent/'r1002/extension-r2/probe-api-observe/witness/events.jsonl'
event_bytes=read(event_path);assert sha(event_bytes)=='22edc4f82f8e34eaa592c4b5dd8842fad2fad044444a41051bbca45da7d43195'
events=[json.loads(s) for s in event_bytes.splitlines()];proof=next(r for r in events if r['status']=='source_extension')
before=base64.b64decode(proof['before8'],validate=True);after=base64.b64decode(proof['after8'],validate=True)
assert len(before)==len(after)==64
mask=[i for i,(x,y) in enumerate(zip(before,after)) if x!=y]
assert len(mask)==13 and all(after[i]==1 for i in mask)
trace_digest='bc95b65d4f308b136f4236993ee97db7f57ebb25d43df622d0697a2a421fc1af'
expected={1519:332800000,1520:338800000,1585:363400000,1686:406800000,1712:417800000}
backends=[a.font] if a.font else ['eten','unifont'];scales=[a.scale] if a.scale else [2,3]
case_data={};summaries=[]
for scale in scales:
    cases=[f'raw-none-{scale}']+[f'{v}-{b}-{scale}' for b in backends for v in ('old','new')]
    for case in cases:
        try:
            directory=root/('api-'+case);result=load(directory/'result.json')
            check(result['status']=='PASS' and result['exit_code']==0 and not result['invalid'],case+'正常退出')
            check(result['initial']==result['final']==result['source_after'],case+'實際保存metadata保持')
            actual={}
            for f in (directory/'save').iterdir():
                import stat
                st=f.stat();actual[f.name]={'size':st.st_size,'mode':stat.S_IMODE(st.st_mode),'mtime_ns':str(st.st_mtime_ns),'sha256':sha(read(f))}
            check(actual==result['initial'],case+'實際保存檔')
            trace=read(directory/'states.jsonl');rows=[json.loads(s) for s in trace.splitlines()]
            check(sha(trace)==trace_digest==result['trace_sha256'] and len(rows)==result['publications']==1786,case+'固定1786原版摘要')
            check([r['publication'] for r in rows]==list(range(1,1787)),case+'發布完整連續')
            check(all(set(r)=={'publication','step','ticks','ram','indexed','palette','registers'} for r in rows),case+'七欄結構')
            trace_rows={r['publication']:r for r in rows}
            frames={}
            for pub,step in expected.items():
                d=directory/'capture'/f'frame-{pub}';meta=load(d/'capture.json');drawmeta=load(d/'draw.json');info=load(d/'hd.json')
                check(meta['publication']==drawmeta['publication']==pub and meta['step']==drawmeta['step']==info['Step']==step,case+str(pub)+'HD／Draw發布配對')
                check(meta['scale']==drawmeta['scale']==scale and not meta['help_open'] and not meta['map_open'],case+str(pub)+'配置與面板')
                indexed=read(d/'hd.a000');comp=read(d/'hd.comp');palette=read(d/'hd.pal');assert len(indexed)==len(comp)==64000 and len(palette)==768
                original=trace_rows[pub]
                check(original['step']==step and sha(indexed)==original['indexed'] and sha(palette)==original['palette'] and info['MemorySHA256']==original['ram'],case+str(pub)+'真實原幀綁定固定trace')
                hd=image(d/'hd.png');draw=image(d/'draw.png');equal(hd,draw,case+str(pub)+'真正Draw／HD')
                check(hd.size==(320*scale,200*scale),case+str(pub)+'畫布尺寸')
                if case.startswith('raw-'):
                    check(not info['Shown'] and comp==indexed,case+str(pub)+'無覆繪實際資格')
                    original=Image.frombytes('P',(320,200),indexed);original.putpalette(palette)
                    equal(original.convert('RGB').resize(hd.size,Image.Resampling.NEAREST),hd,case+str(pub)+'無覆繪HD原圖')
                frames[pub]={'hd':hd,'draw':draw,'info':info,'meta':meta,'indexed':indexed,'palette':palette}
            case_data[case]=frames;summaries.append({'case':case,'frames':len(frames),'publications':1786})
        except (OSError,ValueError,KeyError,TypeError,AssertionError):
            check(False,case+'缺少或損壞收據')
    for backend in backends:
        old=case_data.get(f'old-{backend}-{scale}');new=case_data.get(f'new-{backend}-{scale}')
        if old is None or new is None:continue
        size,glyph=fonts[backend];d=max(size-6*scale,0);px,py=90*scale-d,6*scale-d
        for pub in expected:
            left,right=old[pub],new[pub];check(left['indexed']==right['indexed'] and left['palette']==right['palette'],f'{backend}/{scale}/{pub}原圖／色盤')
            target=left['hd'].copy()
            if pub in (1520,1585):
                roi=bytes(left['indexed'][y*320+x] for y in range(4,12) for x in range(88,96))
                check(roi==after,f'{backend}/{scale}/{pub}原M與完整8×8來源相符')
                palette=left['palette'];color=lambda index:tuple(palette[index*3:index*3+3])
                original_roi=Image.new('RGB',(8,8));original_roi.putdata([color(index) for index in after])
                equal(original_roi.resize((8*scale,8*scale),Image.Resampling.NEAREST),left['hd'].crop((88*scale,4*scale,96*scale,12*scale)),f'{backend}/{scale}/{pub}修前HD批准域確為原M及背景')
                for i in mask:
                    x,y=88+i%8,4+i//8
                    for yy in range(y*scale,(y+1)*scale):
                        for xx in range(x*scale,(x+1)*scale):target.putpixel((xx,yy),color(before[i]))
                for y,row in enumerate(glyph):
                    for x in range(size):
                        if row&(1<<(size-1-x)):target.putpixel((px+x,py+y),color(1))
                translated=[r for r in right['meta']['sex'] or [] if r['Return']=={'Seg':272,'Off':38293} and r['Shown']]
                check(len(translated)==1,f'{backend}/{scale}/{pub}新來源實際shown')
                if translated:
                    e=translated[0];check((e['BX'],e['BY'],e['BW'],e['BH'])==(88,4,8,8) and e['G']==[glyph],f'{backend}/{scale}/{pub}完整守門與外部字模')
                check(88*scale<=px and px+size<=96*scale and 4*scale<=py and py+size<=12*scale,f'{backend}/{scale}/{pub}完整原生字格無裁切')
            else:
                check(not any(r['Return']=={'Seg':272,'Off':38293} for r in right['meta']['sex'] or []),f'{backend}/{scale}/{pub}無過期新來源')
            equal(target,right['hd'],f'{backend}/{scale}/{pub}獨立字模全HD及域外')
            equal(target,right['draw'],f'{backend}/{scale}/{pub}獨立字模真正Draw及域外')
result={'status':'FAIL' if failures else 'PASS','checks':checks,'failures':failures,'configurations':summaries,'canvas_rgb_pixels_compared':pixels,
        'font_sha256':{k:inputs[str(path)] for k,path in [('eten',a.eten),('unifont',a.unifont)]},'input_sha256':inputs,
        'limits':['Only five sampled pixel frames per case;1786 original summaries are separately verified.','Normal GUI and formal replay remain separate gates.']}
assert not a.out.exists()
a.out.write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n');print(json.dumps({k:v for k,v in result.items() if k!='input_sha256'},ensure_ascii=False,indent=2))
raise SystemExit(1 if failures else 0)
