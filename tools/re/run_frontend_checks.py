#!/usr/bin/env python3
"""spec 006 回歸：在前端 Docker／Xvfb 內重播既有玩家測試，保留舊收據。

入口見 docs/spec/006-flow-text.md；不得在主機執行。
"""
import json
import os
import re
import signal
import subprocess
import sys
import time
from pathlib import Path

out = Path(sys.argv[1] if len(sys.argv)>1 else 'workplace/r675')
if out.exists():
    raise SystemExit('收據目錄已存在；請換新編號')
out.mkdir()
build_source='./cmd/wiz7'
if '--baseline' in sys.argv:
    source=out/'source'
    source.mkdir()
    names=subprocess.check_output(['git','ls-tree','-r','--name-only','a63e4f8','cmd/wiz7'],text=True).splitlines()
    for name in names:
        if name.endswith('.go'):
            (source/Path(name).name).write_bytes(subprocess.check_output(['git','show',f'a63e4f8:{name}']))
    build_source='./'+str(source)
binary = str((out/'wiz7').resolve())
with (out/'build.log').open('w') as build_log:
    subprocess.run(['go','build','-o',binary,build_source],stdout=build_log,stderr=subprocess.STDOUT,check=True)

def bounded(cmd, logpath, done=None, limit=360):
    started = time.monotonic()
    with logpath.open('w') as log:
        p = subprocess.Popen(cmd, stdout=log, stderr=subprocess.STDOUT, start_new_session=True)
        try:
            while True:
                if done and done():
                    break
                rc = p.poll()
                if rc is not None:
                    if done or rc:
                        raise RuntimeError(f'{logpath}: 提早／錯誤結束 {rc}')
                    break
                if time.monotonic()-started>limit:
                    raise RuntimeError(f'{logpath}: 超時')
                time.sleep(.25)
        finally:
            if p.poll() is None:
                os.killpg(p.pid, signal.SIGTERM)
                try: p.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    os.killpg(p.pid, signal.SIGKILL)
                    p.wait()
    print(f'{logpath.name} 完成（{time.monotonic()-started:.1f} 秒）',flush=True)

def snapshots(name, flags, steps):
    dest=out/name
    dest.mkdir()
    names=[dest/f'hdat-{s}.json' for s in steps]
    cmd=[binary,'-hdonly','-root','workplace/run','-autoplay','-hd-dir',str(dest),'-hd-at',','.join(map(str,steps))]+flags
    bounded(cmd,out/f'{name}.log',lambda:all(p.exists() for p in names))
    return [json.loads(p.read_text()) for p in names]

if '--t6-wall' in sys.argv:
    source=Path('tools/fe-input-check.sh').read_text()
    inner=source.split("<<'EOS' || true\n",1)[1].split('\nEOS',1)[0]
    inner=inner.split('for t in $TESTS; do',1)[0]
    inner=inner.replace('go build -o /tmp/wiz7 ./cmd/wiz7 || exit 3',f'cp "{binary}" /tmp/wiz7')
    bounded(['bash','-c',f'OUT={out};\n'+inner+'\nt6 input yes\nt6 ctl no\n'],out/'input-check.txt',limit=800)
    text=(out/'input-check.txt').read_text()
    samples=[float(x) for x in re.findall(r'T6-input：第 \d 下.*差異像素=(\d+) 新印字=0',text)]
    assert len(samples)==2 and all(x>10000 for x in samples),text
    print('T6既有視窗翻頁／前進／轉向回歸通過',flush=True)
    sys.exit(0)

if '--t6-only' in sys.argv:
    late='--late' in sys.argv
    times=([1550000000,1610000000,1660000000,1715000000,1773000000,1830000000] if late
           else list(range(1560000000,1606000000,9000000)))
    base=1880000000 if late else 1618000000
    clicks=','.join(f'{s}:160:170' for s in times)
    stage=[0]
    def t6_done():
        if stage[0]==0 and (out/f'hdat-{base}.json').exists():
            time.sleep(1)
            subprocess.run(['import','-window','root',str(out/'before.png')],check=True)
            stage[0]=1
            window=subprocess.check_output(['xdotool','search','--name','巫術7'],text=True).splitlines()[0]
            for n,(x,y) in enumerate([(480,486),(540,498)],1):
                previous=(out/'t6.log').read_text().count('左鍵 (')
                subprocess.run(['xdotool','mousemove','--window',window,str(x),str(y),'click','1'],check=True)
                until=time.monotonic()+30
                while (out/'t6.log').read_text().count('左鍵 (')<=previous:
                    if time.monotonic()>until:raise RuntimeError('左鍵未處理')
                    time.sleep(.1)
                time.sleep(2)
                subprocess.run(['import','-window','root',str(out/f'after-{n}.png')],check=True)
            return True
        return False
    bounded([binary,'-root','workplace/run','-autoplay','-autoplay-stop','1435000000','-scale','3',
             '-clicks-at',clicks,'-hd-dir',str(out),'-hd-at',str(base)],out/'t6.log',t6_done)
    diffs=[]
    before=out/'before.png'
    for n in [1,2]:
        after=out/f'after-{n}.png'
        r=subprocess.run(['compare','-metric','AE',str(before),str(after),'null:'],capture_output=True,text=True)
        diffs.append(float(r.stderr.strip()))
        before=after
    (out/'t6-result.json').write_text(json.dumps({'late':late,'baseline':'--baseline' in sys.argv,'diffs':diffs}))
    print('T6同時序像素差異：'+str(diffs),flush=True)
    sys.exit(0)

if '--map-only' in sys.argv:
    coords=[(s,160,170) for s in range(1560000000,1606000000,9000000)]
    coords += [(1614000000,160,184),(1623000000,160,162),(1632000000,180,166),
               (1641000000,180,166),(1650000000,160,162)]
    clicks=','.join(f'{s}:{x}:{y}' for s,x,y in coords)
    common=['-autoplay-stop','1435000000','-clicks-at',clicks]
    a=snapshots('map-a',common+['-map-at','1600000000'],[1660000000])[0]
    b=snapshots('map-b',common,[1660000000])[0]
    assert all(a[k]==b[k] for k in ['Step','Mem','A000','Pal']), '地圖改變原版狀態'
    m=a['Map']
    assert sorted(tuple(v) for v in m['Visited'])==[(11,56,91),(11,56,92)],m
    assert (m['Area'],m['X'],m['Y'],m['Face'])==(11,56,92,1),m
    sent=[False]
    def map_done():
        log=(out/'map-draw.log').read_text()
        if not sent[0] and re.search(r'測試點擊 \(160,162\) steps=165',log):
            window=subprocess.check_output(['xdotool','search','--name','巫術7'],text=True).splitlines()[0]
            subprocess.run(['xdotool','key','--window',window,'F2'],check=True)
            sent[0]=True
        if '地圖 面板=true' in log:
            time.sleep(2)
            subprocess.run(['import','-window','root',str(out/'map.png')],check=True)
            return True
        return False
    bounded([binary,'-root','workplace/run','-autoplay','-scale','3']+common,out/'map-draw.log',map_done)
    assert '送鍵 3c' not in (out/'map-draw.log').read_text(),'F2送進原版'
    print('地圖兩格、A/B、正常視窗F2回歸通過',flush=True)
    sys.exit(0)

steps=[272000000,286000000,330000000]
snapshots('overlay-en',['-xlate','','-flow',''],steps)
snapshots('overlay-zh',[],steps)
subprocess.run(['python3','tools/re/check_overlay.py',str(out/'overlay-zh'),str(out/'overlay-en')],
               stdout=(out/'overlay-check.txt').open('w'),check=True)
ha=snapshots('help-a',[],[300000000])[0]
hb=snapshots('help-b',['-help-at','286000000','-help-for','3000'],[300000000])[0]
assert all(ha[k]==hb[k] for k in ['Step','Mem','A000','Pal']), 'F1 暫停改變原版狀態'
log=(out/'help-b.log').read_text()
paused=re.findall(r'說明 暫停=(true|false) steps=(\d+)',log)
assert len(paused)>=2 and paused[0][0]=='true' and paused[1][0]=='false' and paused[0][1]==paused[1][1],paused

draw=out/'draw'
draw.mkdir()
ds=[363113713,395347218]
bounded([binary,'-autoplay','-autoplay-stop','400000000','-scale','3','-hd-dir',str(draw),'-draw-at',','.join(map(str,ds))],
        out/'draw.log',lambda:all((draw/f'drawat-{s}.png').exists() for s in ds))

# 沿用原先已驗的正常玩家輸入腳本，只改輸出路徑，不經 host 包裝。
for script,name,tests in [('tools/fe-leftclick-check.sh','left',''),('tools/fe-input-check.sh','input','T2 T4 T6')]:
    dest=out/name
    dest.mkdir()
    source=Path(script).read_text()
    inner=source.split("<<'EOS' || true\n",1)[1].split('\nEOS',1)[0]
    bounded(['bash','-c',f'OUT={dest}; TESTS="{tests}";\n'+inner],out/f'{name}-check.txt',limit=1200)
left=(out/'left-check.txt').read_text()
assert re.search(r'click：左鍵紀錄=[1-9]\d* 對話框 OK/CANCEL 印字=[1-9]',left),left
assert 'noclick：左鍵紀錄=0 對話框 OK/CANCEL 印字=0' in left,left
inp=(out/'input-check.txt').read_text()
assert '等不到' not in inp, inp
assert re.search(r'T2-input：旁白第一頁印字=[1-9]',inp),inp
assert 'T2-ctl：旁白第一頁印字=0' in inp,inp
for name,minimum,maximum,clicks in [('input',1000,10000000,1),('ctl',0,1000,0)]:
    m=re.search(rf'T4-{name}：點擊前後差異像素=(\d+) 左鍵紀錄=(\d+)',inp)
    assert m and minimum<=int(m[1])<=maximum and int(m[2])==clicks,inp
for name in ['input','ctl']:
    assert f'T6-{name}：翻旁白' in inp,inp
for m in re.finditer(r'T6-input：第 \d 下.*差異像素=(\d+) 新印字=(\d+)',inp):
    assert int(m[1])>10000 and int(m[2])==0,m.group()
assert len(re.findall('T6-input：第',inp))==2,inp
print('覆繪、F1、故事 Draw、左鍵與 T2／T4／T6 回歸通過',flush=True)
