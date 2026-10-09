"""正常GUI有聲試玩；只讀既有日誌，原版由正式前端執行。"""
import ctypes,hashlib,json,os,re,shutil,signal,subprocess,time
from pathlib import Path

out=Path('/out');gui=out/'gui';gui.mkdir();save=gui/'save';save.mkdir()
for p in Path('/initial').iterdir():shutil.copy2(p,save/p.name)
binary=Path('/src/workplace/r1000/promoted-build/wiz7-formal')
assert hashlib.sha256(binary.read_bytes()).hexdigest()=='c7da1a38a0a690de1529d0ebadf13c3378e1b51bf3063a522a707962c88fe224'
args=[str(binary),'-root','/original','-save-dir',str(save),'-scale','3','-hd-scale','3','-xlate','/src/workplace/r995/party-assets/unifont/xlate.json','-flow','/src/workplace/r983/final/unifont/flow.json','-help','/src/workplace/r979/unifont/help.json']
log=(gui/'run.log').open('x');proc=subprocess.Popen(args,stdout=log,stderr=subprocess.STDOUT,cwd='/tmp')
wid=None;paused=False;operations=[];deadline=time.monotonic()+840
def logs():return (gui/'run.log').read_text(errors='replace')
def wait(predicate,seconds=60):
 end=min(time.monotonic()+seconds,deadline)
 while time.monotonic()<end:
  if predicate():return
  if proc.poll() is not None:raise RuntimeError('frontend_exit_'+str(proc.returncode))
  time.sleep(.2)
 raise RuntimeError('normal_gui_timeout')
def event(kind,**data):
 row={'kind':kind,'monotonic_ns':time.monotonic_ns(),'device_bytes':(out/'device.s16le').stat().st_size,**data};operations.append(row)
 with (out/'operations.jsonl').open('a') as f:f.write(json.dumps(row)+'\n')
def capture(name):subprocess.run(['import','-window',wid,str(gui/(name+'.png'))],check=True,timeout=15)
def key(value):subprocess.run(['xdotool','key','--window',wid,value],check=True,timeout=5)
def pause(opened):
 global paused
 if paused==opened:return
 before=len(re.findall(r'說明 暫停='+str(opened).lower()+r' steps=(\d+)',logs()))
 key('F1')
 wait(lambda:len(re.findall(r'說明 暫停='+str(opened).lower()+r' steps=(\d+)',logs()))>before,45)
 step=int(re.findall(r'說明 暫停='+str(opened).lower()+r' steps=(\d+)',logs())[-1]);paused=opened;event('pause' if opened else 'resume',step=step)
def click(x,y):
 assert 0<=x<320 and 0<=y<200
 before=len(re.findall(r'左鍵 \('+str(x)+','+str(y)+r'\) steps=(\d+) err=(\w*)',logs()))
 subprocess.run(['xdotool','mousemove','--window',wid,str(x*3),str(y*3),'click','1'],check=True,timeout=5)
 pattern=r'左鍵 \('+str(x)+','+str(y)+r'\) steps=(\d+) err=(\w*)'
 wait(lambda:len(re.findall(pattern,logs()))>before,65)
 step,err=re.findall(pattern,logs())[-1];event('click',x=x,y=y,step=int(step),error=err)
def crop(p):return subprocess.check_output(['convert',str(p),'-crop','450x230+260+140','+repage','RGB:-'],timeout=15)
def close():
 lib=ctypes.CDLL('libX11.so.6');lib.XOpenDisplay.restype=ctypes.c_void_p
 lib.XInternAtom.argtypes=[ctypes.c_void_p,ctypes.c_char_p,ctypes.c_int];lib.XInternAtom.restype=ctypes.c_ulong
 lib.XSendEvent.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.c_int,ctypes.c_long,ctypes.c_void_p]
 lib.XFlush.argtypes=[ctypes.c_void_p];lib.XCloseDisplay.argtypes=[ctypes.c_void_p]
 class Client(ctypes.Structure):_fields_=[('type',ctypes.c_int),('serial',ctypes.c_ulong),('send',ctypes.c_int),('display',ctypes.c_void_p),('window',ctypes.c_ulong),('message',ctypes.c_ulong),('format',ctypes.c_int),('data',ctypes.c_long*5)]
 class Event(ctypes.Union):_fields_=[('client',Client),('padding',ctypes.c_long*24)]
 d=lib.XOpenDisplay(None);assert d
 e=Event();e.client.type=33;e.client.display=d;e.client.window=int(wid);e.client.message=lib.XInternAtom(d,b'WM_PROTOCOLS',0);e.client.format=32;e.client.data[0]=lib.XInternAtom(d,b'WM_DELETE_WINDOW',0)
 event('close_request');lib.XSendEvent(d,int(wid),0,0,ctypes.byref(e));lib.XFlush(d);lib.XCloseDisplay(d)
 proc.wait(timeout=10);event('process_exit',code=proc.returncode)
def interrupted(signum,frame):raise InterruptedError('bounded player interrupted')
signal.signal(signal.SIGTERM,interrupted)
try:
 wait(lambda:bool(subprocess.run(['xdotool','search','--name','巫術7'],capture_output=True,text=True).stdout.strip()),30)
 wid=subprocess.check_output(['xdotool','search','--name','巫術7'],text=True).split()[-1];subprocess.run(['xdotool','windowfocus','--sync',wid],check=True,timeout=5)
 capture('startup');print('ready=startup',flush=True)
 expected=crop(Path('/src/workplace/r981/gui-eten-3-normal-flow-story/menu.png'))
 def menu():
  capture('waiting-menu');return crop(gui/'waiting-menu.png')==expected
 wait(menu,360);capture('menu');event('menu');click(160,73);time.sleep(.3);capture('cold-load-visible');pause(True);capture('cold-load-visible-f1');print('ready=cold-load-paused',flush=True)
 for number in range(1,31):
  p=out/f'command-{number}.json';wait(lambda:p.exists(),120);c=json.loads(p.read_text())
  label=c.get('label','player-step');assert label.replace('-','').isalnum() and len(label)<=32
  if c.get('finish'):
   pause(False);time.sleep(.5);capture('finished');break
  pause(False);event('operation_start',number=number,label=label)
  if 'click' in c:click(*c['click'])
  elif 'key' in c:
   assert c['key'] in ('Escape','Return','space','Up','Down','Left','Right','Tab');key(c['key'])
  elif 'text' in c:
   v=c['text'];assert v.isascii() and v.isalpha() and 1<=len(v)<=16
   subprocess.run(['xdotool','type','--window',wid,'--delay','120',v],check=True,timeout=10)
  elif 'move' in c:
   x,y=c['move'];assert 0<=x<320 and 0<=y<200;subprocess.run(['xdotool','mousemove','--window',wid,str(x*3),str(y*3)],check=True,timeout=5)
  delay=c.get('wait',1);assert 0<=delay<=30;time.sleep(delay);capture(label);pause(True);capture(label+'-f1');event('operation_end',number=number,label=label)
  (out/f'ack-{number}.json').write_text(json.dumps({'label':label,'normal_input':True,'paused':True})+'\n');print('ready='+label,flush=True)
finally:
 if proc.poll() is None:
  if wid:close()
  else:proc.terminate();proc.wait(timeout=10)
 log.close();print('normal-close='+str(proc.returncode),flush=True)
