from pathlib import Path
import argparse
import hashlib
import shutil
import subprocess

p=argparse.ArgumentParser()
p.add_argument('--source',type=Path,required=True)
p.add_argument('--main-sha',required=True)
p.add_argument('--target',type=Path,required=True)
p.add_argument('--repo',type=Path,default=Path('/src'))
a=p.parse_args()
source=a.source.resolve();target=a.target.absolute();repo=a.repo.resolve()
sha=lambda f:hashlib.sha256(f.read_bytes()).hexdigest()
assert source.is_dir() and target.parent.is_dir() and not target.exists()
assert target.parent.stat().st_uid==1000
assert sha(source/'cmd/wiz7/main.go')==a.main_sha
api=repo/'tools/re/item_probe_api.go'
assert sha(api)=='9833d39ad7a12943d0106b9f1ac8d0540cabcf15c85a632667a1c5f2c9f211ff'
target.mkdir()
for name in ('cmd/wiz7','initconfig','recording'):
    shutil.copytree(source/name,target/name,symlinks=True)
for name in ('go.mod','go.sum','tools/re/build_recording.py'):
    dest=target/name;dest.parent.mkdir(parents=True,exist_ok=True)
    shutil.copy2(repo/name if name.startswith('tools/') else source/name,dest)
(target/'worktrees').mkdir()
(target/'worktrees/dosgolem').symlink_to(repo/'worktrees/dosgolem',target_is_directory=True)
shutil.copy2(api,target/'cmd/wiz7/r982_probe.go')
shutil.copy2(Path(__file__).with_name('creation_capture.go'),target/'cmd/wiz7/creation_capture.go')
main=target/'cmd/wiz7/main.go';text=main.read_text()
changes={
 '\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)':
 '\tlive = r982Control(o, *root, preparedSaveDir)\n\tif live != nil { log.SetFlags(0); log.SetOutput(recordingPrivacyLog{}) }\n\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)',
 '\tsnapshot := func() {\n':'\tsnapshot := func() {\n\t\tr982Observe(o, live)\n',
 '\t\t\tpublication, published = live.publish(o.Steps(), sh)\n':
 '\t\t\tpublication, published = live.publish(o.Steps(), sh)\n\t\t\tif published { creationCapture(o, ov, am, sh, publication, *scale) }\n',
 '\tpublication := g.sh.recordingPublication\n\tg.sh.mu.Unlock()':
 '\tpublication := g.sh.recordingPublication\n\treq := g.sh.drawReq\n\tg.sh.drawReq = ""\n\tg.sh.mu.Unlock()',
 '\tg.sh.mu.Lock()\n\treq := g.sh.drawReq\n\tg.sh.drawReq = ""\n\tg.sh.mu.Unlock()\n\tif req != "" {':'\tif req != "" {',
 '\t}\n\tif g.live != nil {\n\t\tg.live.drawControls(dst, g.scale)':
 '\t}\n\tcreationCaptureDraw(req, publication, g.scale)\n\tif g.live != nil {\n\t\tg.live.drawControls(dst, g.scale)',
}
for old,new in changes.items():
    assert text.count(old)==1,repr(old)
    text=text.replace(old,new)
main.write_text(text)
subprocess.run(['gofmt','-w',str(main),str(target/'cmd/wiz7/creation_capture.go'),str(target/'cmd/wiz7/r982_probe.go')],check=True,timeout=30)
assert sha(source/'cmd/wiz7/main.go')==a.main_sha
print('Private capture prepared; original not executed.')
