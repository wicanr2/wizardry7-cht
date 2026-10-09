#!/usr/bin/env python3
"""在Docker內從正式來源建立037私人驗收副本；本工具不執行原版。"""
import argparse
import os
from pathlib import Path
import shutil
import subprocess

p=argparse.ArgumentParser(description=__doc__)
p.add_argument('--source',type=Path,required=True)
p.add_argument('--output',type=Path,required=True)
a=p.parse_args();source=a.source.resolve();target=a.output.absolute()
if not source.is_dir() or target.exists() or target.is_symlink():raise RuntimeError('source_or_output')
if not target.parent.is_dir() or target.parent.stat().st_uid!=os.getuid():raise RuntimeError('output_parent_owner')
target.mkdir()
for name in ('cmd/wiz7','initconfig','recording'):shutil.copytree(source/name,target/name,symlinks=True)
for name in ('go.mod','go.sum','tools/re/build_recording.py'):
 dest=target/name;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(source/name,dest)
(target/'worktrees').mkdir();(target/'worktrees/dosgolem').symlink_to(source/'worktrees/dosgolem',target_is_directory=True)
here=Path(__file__).resolve().parent
for original,dest in (('item_probe_api.go','r982_probe.go'),('held_item_verify.go','formal_observer.go')):shutil.copy2(here/original,target/'cmd/wiz7'/dest)
path=target/'cmd/wiz7/main.go';s=path.read_text()
changes={
 '\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)':
 '\tlive = r982Control(o, *root, preparedSaveDir)\n\tif live != nil { log.SetFlags(0); log.SetOutput(recordingPrivacyLog{}) }\n\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)',
 '\tsnapshot := func() {\n':'\tsnapshot := func() {\n\t\tformalBefore(o, ov)\n\t\tr982Observe(o, live)\n',
 '\t\tsh.held = held\n':'\t\tsh.held = held\n\t\tformalObserve(o, ov, sh, live, *scale)\n',
 '\tpublication := g.sh.recordingPublication\n\tg.sh.mu.Unlock()':'\tpublication := g.sh.recordingPublication\n\treq := g.sh.drawReq\n\tg.sh.drawReq = ""\n\tg.sh.mu.Unlock()',
 '\tg.sh.mu.Lock()\n\treq := g.sh.drawReq\n\tg.sh.drawReq = ""\n\tg.sh.mu.Unlock()\n\tif req != "" {':'\tif req != "" {',
}
for old,new in changes.items():
 if s.count(old)!=1:raise RuntimeError('source_hook_mismatch')
 s=s.replace(old,new)
path.write_text(s)
subprocess.run(['gofmt','-w',str(path),str(target/'cmd/wiz7/formal_observer.go'),str(target/'cmd/wiz7/r982_probe.go')],check=True,timeout=30)
print('private held-item verifier prepared; original not executed')
