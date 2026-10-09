#!/usr/bin/env python3
"""在Docker內準備R1000 READY私人游標背景取樣副本，不執行原版。"""
import argparse
from pathlib import Path
import shutil
import subprocess

parser=argparse.ArgumentParser(description=__doc__)
parser.add_argument('--source',type=Path,required=True)
parser.add_argument('--output',type=Path,required=True)
args=parser.parse_args()
source=args.source.resolve(); target=args.output.absolute()
subprocess.run(['python3',str(source/'tools/re/prepare_item_bitmap_probe.py'),
                '--source',str(source),'--output',str(target)],check=True,timeout=60)
helper=target/'cmd/wiz7/r1000_cursor_probe.go'
shutil.copy2(Path(__file__).resolve().parent/'cursor_background_probe.go',helper)
path=target/'cmd/wiz7/r999_service24.go'; code=path.read_text()
changes={
 'func r999ServiceDiscard(reason string) {\n':'func r999ServiceDiscard(reason string) {\n\tr1000CursorDiscard(reason)\n',
 '\t\tr999Service.pending = &r999ServicePending{sample: s, before: before, indexed: indexed}\n':
 '\t\tr1000CursorServiceEntry(s)\n\t\tr999Service.pending = &r999ServicePending{sample: s, before: before, indexed: indexed}\n',
 '\t\ts := p.sample\n':'\t\ts := p.sample\n\t\tdefer func() { r1000CursorServiceReturn(o, s) }()\n',
}
for old,new in changes.items():
 if code.count(old)!=1:raise RuntimeError('service_hook_mismatch')
 code=code.replace(old,new)
path.write_text(code)
trace=target/'cmd/wiz7/r982_probe.go'; code=trace.read_text()
old='\tif err := r982Trace.Encode(row); err != nil {\n\t\tpanic(err)\n\t}\n'
if code.count(old)!=1:raise RuntimeError('owner_hook_mismatch')
trace.write_text(code.replace(old,old+'\tr1000CursorSnapshot(o, l.publication+1)\n'))
subprocess.run(['gofmt','-w',str(helper),str(path),str(trace)],check=True,timeout=30)
print('private cursor observer prepared; original not executed')
