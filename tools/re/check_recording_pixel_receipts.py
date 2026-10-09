from pathlib import Path
import argparse, hashlib, json, re
from PIL import Image

p=argparse.ArgumentParser(); p.add_argument('--input',required=True); p.add_argument('--binary',required=True); p.add_argument('--output',required=True); args=p.parse_args()
root=Path(args.input); output=Path(args.output); assert not output.exists()
sha=lambda path:hashlib.sha256(path.read_bytes()).hexdigest()
binary_sha=sha(Path(args.binary)); wire=root/'normal-main/unifont/normal.wiz7-input'
snapshots=[]; content=hashlib.sha256(); event_count=0; footer=None; config=None
with wire.open('rb') as f:
 for number,line in enumerate(f):
  assert line.endswith(b'\n') and footer is None
  if line.startswith(b'{"kind":"end",'):
   footer=json.loads(line); continue
  content.update(line)
  if number==0:
   header=json.loads(line); assert header['kind']=='header' and header['schema']=='wiz7-input/1'; config=header['config_sha256']; continue
  event_count+=1
  if line.startswith(b'{"kind":"snapshot",'):
   e=json.loads(line); q=e['params']; assert e['before_step']==e['after_step'] and q['frontend_sha256']==config
   assert q['publication']==len(snapshots)+1
   snapshots.append({'publication':q['publication'],'step':e['after_step'],'panel':q['panel']})
assert footer and footer['accepted'] and footer['event_count']==event_count
assert footer['content_sha256']==content.hexdigest() and footer['end_step']==snapshots[-1]['step']
assert snapshots and re.fullmatch('[0-9a-f]{64}',config)

cases=[('normal',root/'normal-main/unifont/record-gui',root/'frame-normal')]
cases += [(f'replay-{n}',root/f'replay-{n}',root/f'frame-replay-{n}') for n in (1,2,3)]
traces={}; result={'binary_sha256':binary_sha,'wire_sha256':sha(wire),'events':event_count,'snapshots':len(snapshots),'end_step':footer['end_step'],'cases':[],'comparisons':[]}
all_ok=True
for name,folder,frames in cases:
 draw_path,state_path=folder/'draw.jsonl',folder/'state.jsonl'
 draws=[json.loads(line) for line in draw_path.read_text().splitlines()]
 states=[json.loads(line) for line in state_path.read_text().splitlines()]
 valid=len(draws)==len(states)==len(snapshots)
 problems=[]; pixels=0; png_hashes=[]
 for i,(d,s,witness) in enumerate(zip(draws,states,snapshots)):
  for key in ('publication','step','panel'):
   if d[key]!=witness[key] or s[key]!=witness[key]: problems.append({'index':i,'field':key})
  for row in (d,s):
   assert row['config_sha256']==config and row['binary_sha256']==binary_sha
  assert d['width']==640 and d['height']==400
  assert re.fullmatch('[0-9a-f]{64}',d['pixel_sha256'])
  for key in ('ram','indexed','palette','registers'): assert re.fullmatch('[0-9a-f]{64}',s[key])
  png=frames/f'{d["publication"]:08d}.png'
  with Image.open(png) as img:
   assert img.size==(d['width'],d['height'])
   rgba=img.convert('RGBA').tobytes()
  decoded=hashlib.sha256(rgba).hexdigest()
  if decoded!=d['pixel_sha256']: problems.append({'index':i,'field':'png_rgba_sha','decoded':decoded,'recorded':d['pixel_sha256']})
  pixels+=d['width']*d['height']; png_hashes.append((png.name,sha(png)))
 # 正常停止後不再有同publication的新收據或額外PNG。
 assert len(list(frames.glob('*.png')))==len(draws)
 for path in (folder/'run.log',):
  assert b'R998_OBSERVER_INVALID' not in path.read_bytes()
 valid=valid and not problems; all_ok=all_ok and valid
 result['cases'].append({'case':name,'draws':len(draws),'states':len(states),'pngs':len(png_hashes),'decoded_rgba_pixels':pixels,'valid':valid,'problems':problems[:30],'draw_file_sha256':sha(draw_path),'state_file_sha256':sha(state_path),'png_manifest_sha256':hashlib.sha256(json.dumps(png_hashes,separators=(',',':')).encode()).hexdigest()})
 traces[name]=(draws,states)
normal_draw,normal_state=traces['normal']
for n in (1,2,3):
 name=f'replay-{n}'; draws,states=traces[name]; mismatches=[]; state_mismatches=[]
 for i,(a,b) in enumerate(zip(normal_draw,draws)):
  if a!=b: mismatches.append({'publication':i+1,'step':a['step'],'normal_pixel_sha256':a['pixel_sha256'],'replay_pixel_sha256':b['pixel_sha256'],'changed_fields':[k for k in a if a[k]!=b[k]]})
 for i,(a,b) in enumerate(zip(normal_state,states)):
  if a!=b: state_mismatches.append({'publication':i+1,'step':a['step'],'changed_fields':[k for k in a if a[k]!=b[k]]})
 same=len(draws)==len(normal_draw) and len(states)==len(normal_state) and not mismatches and not state_mismatches
 all_ok=all_ok and same
 result['comparisons'].append({'replay':n,'matched':same,'pixel_mismatches':len(mismatches),'state_mismatches':len(state_mismatches),'first_pixel_mismatches':mismatches[:30],'first_state_mismatches':state_mismatches[:30]})
result['all_match']=all_ok
with output.open('x') as f: json.dump(result,f,ensure_ascii=False,indent=2); f.write('\n')
print(json.dumps({'snapshots':len(snapshots),'decoded_rgba_pixels':sum(r['decoded_rgba_pixels'] for r in result['cases']),'all_match':all_ok,'comparisons':[{k:r[k] for k in ('replay','pixel_mismatches','state_mismatches')} for r in result['comparisons']]}))
if not all_ok: raise SystemExit(1)
