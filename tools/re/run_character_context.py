from pathlib import Path
import hashlib,json,os,shutil,stat,subprocess

src=Path('/src');out=Path('/out');variant=os.environ['WIZ7_R1002_VARIANT'];assert variant in ('control','observe')
assert list(out.iterdir())==[]
save=out/'save';save.mkdir();witness=out/'witness';witness.mkdir()
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
def meta(p):
 s=p.stat();return {'size':s.st_size,'mode':stat.S_IMODE(s.st_mode),'mtime_ns':str(s.st_mtime_ns),'sha256':sha(p)}
initial={p.name:meta(p) for p in Path('/initial').iterdir()}
for p in Path('/initial').iterdir():shutil.copy2(p,save/p.name)
wire=src/'workplace/r1002/profile-context.wiz7-input';expected_wire='e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732';assert sha(wire)==expected_wire
binary=Path(os.environ['WIZ7_R1002_BINARY']);expected_binary=os.environ['WIZ7_R1002_BINARY_SHA'];assert sha(binary)==expected_binary
env=os.environ.copy();env['WIZ7_R982_API_DOCUMENT']=str(wire);env['WIZ7_R982_TRACE']=str(out/'states.jsonl')
if variant=='observe':env['WIZ7_R1002_CONTEXT']=str(witness/'events.jsonl')
else:env.pop('WIZ7_R1002_CONTEXT',None)
args=[str(binary),'-root','/original','-save-dir',str(save),'-mute','-ips','24000000','-scale','3','-hd-scale','3','-xlate',str(src/'workplace/r995/party-assets/unifont/xlate.json'),'-flow',str(src/'workplace/r983/final/unifont/flow.json'),'-help',str(src/'workplace/r979/unifont/help.json')]
with (out/'run.log').open('x') as log:p=subprocess.run(args,env=env,stdout=log,stderr=subprocess.STDOUT,timeout=880)
final={p.name:meta(p) for p in save.iterdir()};invalid='R1002_CONTEXT_INVALID_' in (out/'run.log').read_text()
j={'variant':variant,'exit_code':p.returncode,'wire_sha256':expected_wire,'binary_sha256':expected_binary,'trace_sha256':sha(out/'states.jsonl'),'publications':len((out/'states.jsonl').read_bytes().splitlines()),'initial':initial,'final':final,'invalid':invalid,'private_api_only':True,'formal_cross_binary_compatibility':False}
(out/'result.json').write_text(json.dumps(j,indent=2)+'\n');print({k:v for k,v in j.items() if k not in ('initial','final')},flush=True)
assert p.returncode==0 and not invalid and initial==final and j['publications']==1786 and sha(wire)==expected_wire
