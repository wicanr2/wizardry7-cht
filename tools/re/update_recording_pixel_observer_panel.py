#!/usr/bin/env python3
"""R998 R4：僅修owner摘要與真Draw panel的提交綁定；不跑遊戲。"""
from pathlib import Path
import hashlib
import shutil
import subprocess

work=Path('/out');base=work/'source';target=work/'source-r4'
assert work.stat().st_uid==1000 and not target.exists()
assert hashlib.sha256((base/'cmd/wiz7/r998_observer.go').read_bytes()).hexdigest()=='6fc9b36dc4419b3acddf66948379d18ccc62793b5e0808ebf6f45302c8166a14'
assert hashlib.sha256((base/'cmd/wiz7/main.go').read_bytes()).hexdigest()=='13d3139fd5c7c64d0e07d4066fadc4e7be9e9846d69938c1a7dac779aa011233'
shutil.copytree(base,target,symlinks=True)
p=target/'cmd/wiz7/r998_observer.go';s=p.read_text()
# Only pure scalars and already computed hashes cross from owner to GUI.
anchor='func r998ObserveDraw('
assert s.count(anchor)==1
pending='\ntype r998PendingOriginal struct {\n    mu sync.Mutex\n    valid bool\n    publication, step uint64\n    config string\n    original recordingOriginal\n}\nvar r998PendingState r998PendingOriginal\n\n'
s=s.replace(anchor,pending+anchor,1)
start=s.index('func r998ObserveDraw(');end=s.index('func r998ObserveState(',start)
part=s[start:end];old='\ts.write(p, row)';assert part.count(old)==1
part=part.replace(old,'\tif !r998CommitState(p, panel) {\n\t\ts.invalid("state_commit")\n\t\treturn\n\t}\n'+old)
s=s[:start]+part+s[end:]
start=s.index('func r998ObserveState(');brace=s.index('{',start);depth=0
for i in range(brace,len(s)):
    if s[i]=='{':depth+=1
    elif s[i]=='}':
        depth-=1
        if depth==0:
            end=i+1
            break
replacement='func r998ObserveState(o *oracle.Oracle, p recordingPublication, _ rc.Panel) {\n    s := &r998StateSink\n    s.mu.Lock()\n    defer s.mu.Unlock()\n    defer func() { if recover() != nil { s.invalid("state_panic") } }()\n    if !s.ready(p) { return }\n    original := recordingOriginalHash(o)\n    if original.Step != p.Step { s.invalid("state_step"); return }\n    pending := &r998PendingState\n    pending.mu.Lock()\n    defer pending.mu.Unlock()\n    if pending.valid { s.invalid("pending_overwrite"); return }\n    pending.publication, pending.step, pending.config = p.ID, p.Step, p.Config\n    pending.original = original\n    pending.valid = true\n}\n\n// GUI只提交owner已算好的摘要，絕不讀Oracle；panel取當次真Draw。\nfunc r998CommitState(p recordingPublication, panel rc.Panel) bool {\n    s := &r998StateSink\n    s.mu.Lock()\n    defer s.mu.Unlock()\n    if !s.ready(p) { return false }\n    pending := &r998PendingState\n    pending.mu.Lock()\n    defer pending.mu.Unlock()\n    if !pending.valid { s.invalid("pending_missing"); return false }\n    if pending.publication != p.ID || pending.step != p.Step || pending.config != p.Config {\n        s.invalid("pending_identity")\n        return false\n    }\n    original := pending.original\n    row := struct {\n        Publication uint64 `json:"publication"`\n        Step uint64 `json:"step"`\n        Panel rc.Panel `json:"panel"`\n        ConfigSHA string `json:"config_sha256"`\n        BinarySHA string `json:"binary_sha256"`\n        Ticks uint64 `json:"ticks"`\n        RAM string `json:"ram"`\n        Indexed string `json:"indexed"`\n        Palette string `json:"palette"`\n        Registers string `json:"registers"`\n    }{p.ID,p.Step,panel,p.Config,s.binary,original.Ticks,original.RAM,original.Indexed,original.Palette,original.Registers}\n    s.write(p,row)\n    if s.failed { return false }\n    pending.valid = false\n    pending.publication, pending.step, pending.config = 0, 0, ""\n    pending.original = recordingOriginal{}\n    return true\n}'
s=s[:start]+replacement+s[end:]
p.write_text(s)
subprocess.run(['gofmt','-w',str(p)],check=True)
assert (target/'cmd/wiz7/main.go').read_bytes()==(base/'cmd/wiz7/main.go').read_bytes()
changed=[]
for p in base.rglob('*'):
    if p.is_file() and not p.is_symlink():
        other=target/p.relative_to(base)
        if p.read_bytes()!=other.read_bytes():changed.append(str(p.relative_to(base)))
assert changed==['cmd/wiz7/r998_observer.go'],changed
print('R4 source prepared: helper pending/Draw commit only')
