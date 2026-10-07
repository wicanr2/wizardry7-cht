#!/usr/bin/env python3
"""028正常無輸入冷啟動抽測；Docker/Xvfb，禁止autoplay及防拷輸入。"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time

STEPS = (73400000,75600000,77800000,101200000,103400000,105400000,272000000)
CLOSE = r'''
#include <X11/Xlib.h>
#include <stdlib.h>
int main(int n,char **v){
 if(n!=2)return 2; Display*d=XOpenDisplay(0);if(!d)return 3;
 Window w=strtoul(v[1],0,10);XEvent e={0};e.xclient.type=ClientMessage;
 e.xclient.window=w;e.xclient.message_type=XInternAtom(d,"WM_PROTOCOLS",False);
 e.xclient.format=32;e.xclient.data.l[0]=XInternAtom(d,"WM_DELETE_WINDOW",False);
 e.xclient.data.l[1]=CurrentTime;
 int ok=XSendEvent(d,w,False,NoEventMask,&e);XFlush(d);XCloseDisplay(d);return !ok;
}
'''


def wait_for(predicate, proc, seconds=180):
    end = time.monotonic()+seconds
    while time.monotonic()<end:
        if predicate():
            return
        if proc.poll() is not None:
            raise RuntimeError(f'程序提前結束：{proc.returncode}')
        time.sleep(.05)
    raise RuntimeError('驗收逾時')


def capture(wid,path):
    subprocess.run(['import','-window',wid,str(path)],check=True,timeout=15)
    data = subprocess.check_output(['convert',str(path),'-depth','8','rgba:-'],timeout=15)
    return hashlib.sha256(data).hexdigest(),data


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('binary',type=Path)
    ap.add_argument('original',type=Path)
    ap.add_argument('resources',type=Path)
    ap.add_argument('output',type=Path)
    ap.add_argument('--gui',action='store_true')
    ap.add_argument('--cases',default='english-3,unifont-3,eten-3,unifont-2,eten-2')
    ap.add_argument('--legacy',action='store_true')
    args = ap.parse_args()
    assert args.binary.is_file() and (args.original/'DS.EXE').is_file()
    assert args.output.is_dir() and args.output.stat().st_uid == os.getuid()
    closer = args.output/'wm-close'
    if args.gui:
        source = args.output/'wm-close.c'
        with source.open('x') as f:
            f.write(CLOSE)
        subprocess.run(['cc',str(source),'-lX11','-o',str(closer)],check=True,timeout=30)
    records = []
    for case in args.cases.split(','):
        backend,scale = case.rsplit('-',1)
        assert backend in ('english','unifont','eten') and scale in ('2','3')
        dest = args.output/case
        dest.mkdir(exist_ok=False)
        resource = args.resources/backend
        xlate = '' if backend=='english' else str(resource/('xlate-legacy.json' if args.legacy else 'xlate.json'))
        help_path = '' if backend=='english' else str(resource/'help.json')
        cmd = [str(args.binary),'-root',str(args.original),'-save-dir',str(dest/'saves'),
               '-mute','-ips','1000000000','-scale',scale,'-hd-scale',scale,
               '-xlate',xlate,'-flow','','-help',help_path,
               '-hd-at',','.join(map(str,STEPS)),'-hd-dir',str(dest)]
        if args.gui:
            cmd += ['-draw-at',','.join(map(str,STEPS))]
            subprocess.run(['xdotool','mousemove','0','0'],check=True,timeout=5)
        else:
            cmd += ['-hdonly']
        with (dest/'run.log').open('x') as stream:
            proc = subprocess.Popen(cmd,stdout=stream,stderr=subprocess.STDOUT,cwd='/tmp')
            try:
                wait_for(lambda: all((dest/f'hdat-{step}.json').is_file() for step in STEPS),proc)
                record = dict(case=case,mode='GUI' if args.gui else 'headless',steps=list(STEPS))
                if args.gui:
                    wait_for(lambda: all((dest/f'drawat-{step}.png').is_file() for step in STEPS),proc)
                    wid = subprocess.check_output(['xdotool','search','--name','巫術7'],text=True,timeout=5).split()[-1]
                    subprocess.run(['xdotool','windowfocus','--sync',wid],check=True,timeout=5)
                    before,_ = capture(wid,dest/'menu.png')
                    if backend!='english':
                        subprocess.run(['xdotool','key','--window',wid,'F1'],check=True,timeout=5)
                        wait_for(lambda: '說明 true 畫面=mainmenu' in (dest/'run.log').read_text(),proc,10)
                        time.sleep(.2)
                        panel,raw = capture(wid,dest/'f1.png')
                        pixels = {raw[i:i+4] for i in range(0,len(raw),4)}
                        assert panel != before and bytes((255,221,68,255)) in pixels, 'F1未量到說明畫面'
                        subprocess.run(['xdotool','key','--window',wid,'F1'],check=True,timeout=5)
                        wait_for(lambda: '說明 false 畫面=mainmenu' in (dest/'run.log').read_text(),proc,10)
                        time.sleep(.2)
                        after,_ = capture(wid,dest/'menu-return.png')
                        assert after==before, 'F1返回後畫面不同'
                        record.update(f1=True,menu_rgba_sha256=before)
                    subprocess.run([str(closer),wid],check=True,timeout=5)
                    proc.wait(timeout=10)
                    assert proc.returncode==0, f'正常關窗退出碼{proc.returncode}'
                    record['normal_close'] = True
                else:
                    proc.terminate()
                    proc.wait(timeout=5)
                    record['termination'] = '無頭探針收據完成後有界停止，不當GUI正常退出'
                records.append(record)
                print(case+'完成',flush=True)
            finally:
                if proc.poll() is None:
                    proc.terminate()
                    try:
                        proc.wait(timeout=5)
                    except subprocess.TimeoutExpired:
                        proc.kill();proc.wait()
    with (args.output/'run-receipt.json').open('x') as f:
        json.dump(dict(schema=1,binary_sha256=hashlib.sha256(args.binary.read_bytes()).hexdigest(),
                       gui=args.gui,legacy=args.legacy,cases=records,
                       limits='全程無autoplay、原版按鍵／點擊、防拷回答或state注入；手動加速且靜音，不驗音訊或正常速度。'),f,ensure_ascii=False,indent=2)
        f.write('\n')


if __name__=='__main__':
    main()
