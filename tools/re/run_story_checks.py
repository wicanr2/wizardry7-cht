#!/usr/bin/env python3
"""在既有前端 Docker／Xvfb 內重跑故事與旁白收據；不覆寫既有目錄。

入口見 docs/spec/006-flow-text.md。只重播既有 autoplay，不推導防拷輸入。
"""
import argparse
import json
import subprocess
import time
from pathlib import Path

ap = argparse.ArgumentParser()
ap.add_argument('out', nargs='?', default='workplace/r674')
ap.add_argument('--bake-dir', type=Path, default=Path('workplace/bake'),
                help='字型後端烘製目錄（spec007；原版輸入不變）')
ap.add_argument('--modes', default='en,zh,zh2,drop,narr-drop',
                help='重跑模式，逗號分隔；預設完整五組')
args = ap.parse_args()
out = Path(args.out)
bake = args.bake_dir
modes = args.modes.split(',')
if not modes or len(set(modes)) != len(modes) or set(modes) - {'en', 'zh', 'zh2', 'drop', 'narr-drop'}:
    raise SystemExit('重跑模式不合法')
if out.exists():
    raise SystemExit('收據目錄已存在；請使用新的 workplace/rNNN，禁止覆寫基準')
out.mkdir()
story = Path('workplace/r673/samples.txt').read_text().strip().split(',')
narr = [str(s) for s in range(1483000000, 1539000000, 1000000)]
clicks = ','.join(f'{s}:160:170' for s in range(1440000000, 1541000000, 10000000))
common = ['/tmp/wiz7', '-hdonly', '-autoplay', '-autoplay-stop', '1435000000',
          '-root', 'workplace/run', '-clicks-at', clicks,
          '-xlate', str(bake/'xlate.json'), '-flow', str(bake/'flow.json'),
          '-help', str(bake/'help.json')]
with (out/'build.log').open('w') as log:
    for cmd in [['go','build','-o','/tmp/wiz7','./cmd/wiz7'],
                ['python3','tools/re/bake_font.py',str(bake)],
                ['python3','tools/re/bake_flow.py',str(bake/'flow.json')],
                ['python3','tools/re/bake_help.py',str(bake/'help.json')],
                ['python3','tools/re/bake_flow.py',str(out/'flow-drop.json'),'--drop','31654'],
                ['python3','tools/re/bake_flow.py',str(out/'flow-narr-drop.json'),'--drop','15005']]:
        subprocess.run(cmd, stdout=log, stderr=subprocess.STDOUT, check=True)

def run(name, samples, fragments, extra):
    dest = out/name
    dest.mkdir()
    expected = [f'hdat-{s}' for s in samples] + [f'hdfrag-{s}' for s in fragments]
    cmd = common + ['-hd-dir', str(dest), '-hd-at', ','.join(samples),
                    '-hd-after-frag', ','.join(map(str, fragments))] + extra
    started = time.monotonic()
    with (out/f'{name}.log').open('w') as log:
        proc = subprocess.Popen(cmd, stdout=log, stderr=subprocess.STDOUT)
        try:
            while True:
                done = all(all((dest/(s+suffix)).is_file() for suffix in ['.json','.a000','.comp','.png'])
                           for s in expected)
                if done:
                    # JSON／PNG 全部寫完才停止不自行結束的 hdonly runner。
                    for s in expected:
                        json.loads((dest/(s+'.json')).read_text())
                    break
                if proc.poll() is not None:
                    raise RuntimeError(f'{name} 提早結束：{proc.returncode}')
                if time.monotonic()-started > 360:
                    raise RuntimeError(f'{name} 超時，已完成 {len(list(dest.glob("*.json")))}/{len(expected)}')
                time.sleep(.25)
        finally:
            if proc.poll() is None:
                proc.terminate()
            try:
                proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                proc.kill()
                proc.wait()
    print(f'{name}：{len(expected)} 份完整收據，耗時 {time.monotonic()-started:.1f} 秒', flush=True)

for name, extra in [('en',['-flow','']), ('zh',[]), ('zh2',['-hd-scale','2'])]:
    if name in modes:
        run(name, story+narr, [4,81,82], extra)
if 'drop' in modes:
    run('drop', story, [4], ['-flow',str(out/'flow-drop.json')])
if 'narr-drop' in modes:
    run('narr-drop', narr, [81,82], ['-flow',str(out/'flow-narr-drop.json')])
