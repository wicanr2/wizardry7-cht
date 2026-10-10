from pathlib import Path
import hashlib
import json
import os
import shutil
import stat
import subprocess

src=Path('/src');out=Path('/out')
case=os.environ['WIZ7_R1003_CASE']
parts=case.split('-');assert len(parts)==3
variant,backend,number=parts;scale=int(number)
assert variant in ('old','new','raw') and scale in (2,3)
assert backend in (('none',) if variant=='raw' else ('eten','unifont'))
assert not list(out.iterdir()) and out.stat().st_uid==1000
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
def meta(p):
    s=p.stat();return {'size':s.st_size,'mode':stat.S_IMODE(s.st_mode),'mtime_ns':str(s.st_mtime_ns),'sha256':sha(p)}
initial_dir=Path('/initial');initial={p.name:meta(p) for p in initial_dir.iterdir()}
assert set(initial)=={'SAVEGAME.DBS','PCFILE.DBS'}
save=out/'save';save.mkdir();capture=out/'capture';capture.mkdir()
for p in initial_dir.iterdir():shutil.copy2(p,save/p.name)
wire=src/'workplace/r1002/profile-context.wiz7-input'
assert sha(wire)=='e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732'
identities={'old':('verify-old-r3','6ff7da29c212d634e2b70d867c26beb9ad0db06054d111bb1eafa2177117f6f9'),
            'new':('verify-new-r2','d29168b137ada07255f44b9445169ef76703c29ccc8f14543ae22ff051f93c04')}
directory,binary_sha=identities['old' if variant=='raw' else variant]
binary=src/'workplace/r1003'/directory/'wiz7-capture';assert sha(binary)==binary_sha
env=os.environ.copy();env.update({'WIZ7_R982_API_DOCUMENT':str(wire),'WIZ7_R982_TRACE':str(out/'states.jsonl'),'WIZ7_R1003_CAPTURE':str(capture)})
args=[str(binary),'-root','/original','-save-dir',str(save),'-mute','-ips','24000000','-scale',str(scale),'-hd-scale',str(scale)]
for flag,path in (('xlate',src/f'workplace/r995/party-assets/{backend}/xlate.json'),('flow',src/f'workplace/r983/final/{backend}/flow.json'),('help',src/f'workplace/r979/{backend}/help.json')):
    if variant!='raw':assert path.is_file()
    args += ['-'+flag,'' if variant=='raw' else str(path)]
timed_out=False
with (out/'run.log').open('x') as log:
    try:
        process=subprocess.run(args,env=env,stdout=log,stderr=subprocess.STDOUT,timeout=880)
        exit_code=process.returncode
    except subprocess.TimeoutExpired:
        timed_out=True;exit_code=None
        print('R1003_CAPTURE_INVALID_TIMEOUT',flush=True)
final={p.name:meta(p) for p in save.iterdir()};source_after={p.name:meta(p) for p in initial_dir.iterdir()}
trace=out/'states.jsonl';rows=[];trace_sha=None
invalid=timed_out or 'R1003_CAPTURE_INVALID_' in (out/'run.log').read_text()
try:
    raw_trace=trace.read_bytes()
    trace_sha=hashlib.sha256(raw_trace).hexdigest()
    rows=[json.loads(s) for s in raw_trace.splitlines()]
except (OSError,ValueError):
    invalid=True;print('R1003_CAPTURE_INVALID_TRACE',flush=True)
expected={1519:332800000,1520:338800000,1585:363400000,1686:406800000,1712:417800000}
artifacts=[]
from PIL import Image
try:
    for pub,step in expected.items():
        p=capture/f'frame-{pub}';sample=json.loads((p/'capture.json').read_text());draw=json.loads((p/'draw.json').read_text());hd=json.loads((p/'hd.json').read_text())
        assert sample['publication']==draw['publication']==pub and sample['step']==draw['step']==hd['Step']==step
        assert sample['scale']==draw['scale']==scale and not sample['help_open'] and not sample['map_open']
        for name,length in (('hd.a000',64000),('hd.comp',64000),('hd.pal',768)):assert (p/name).stat().st_size==length
        for name in ('hd.png','draw.png'):
            with Image.open(p/name) as im:im.load();assert im.size==(320*scale,200*scale)
        artifacts.append({'publication':pub,'step':step,'directory':p.name})
except (AssertionError,OSError,ValueError,KeyError,TypeError):
    invalid=True
    print('R1003_CAPTURE_INVALID_ARTIFACTS',flush=True)
result={'status':'PASS' if exit_code==0 and not invalid else 'FAIL','case':case,'exit_code':exit_code,'invalid':invalid,'timed_out':timed_out,
        'binary_sha256':binary_sha,'wire_sha256':sha(wire),'trace_sha256':trace_sha,'publications':len(rows),
        'initial':initial,'final':final,'source_after':source_after,'artifacts':artifacts,'private_api_only':True}
(out/'result.json').write_text(json.dumps(result,indent=2)+'\n')
print({k:v for k,v in result.items() if k not in ('initial','final','source_after','artifacts')},flush=True)
assert exit_code==0 and not invalid and initial==final==source_after and len(rows)==1786
assert trace_sha=='bc95b65d4f308b136f4236993ee97db7f57ebb25d43df622d0697a2a421fc1af'
