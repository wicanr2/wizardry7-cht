#!/usr/bin/env python3
"""故事時序蒐證：可丟棄探針準備與原版像素核對（spec 006、研究 017）。

Docker 內執行：story_evidence.py prepare|analyze。正式程式唯讀；輸出 workplace/r673。
準備模式複製 cmd/wiz7 作純觀察副本，沒有接入故事中文實作。
"""
import hashlib
import json
import re
import subprocess
import sys
from collections import Counter
from pathlib import Path

OUT = Path(sys.argv[2]) if len(sys.argv) > 2 else Path('workplace/r673')
BASELINE = 'a63e4f8'
LOG = Path('workplace/bake/r667-story-msgs.txt')


def pages():
    result = []
    for line in LOG.read_text().splitlines():
        m = re.match(r'印字 (\d+) .*font=2 ret=0110:be0c .*?xy0=\(\d+,(\d+)\)', line)
        if m:
            if int(m[2]) == 146:
                result.append([])
            result[-1].append(int(m[1]))
    assert len(result) == 31 and sum(map(len, result)) == 71
    return result


def prepare():
    OUT.mkdir(exist_ok=True)
    if (OUT/'timing/traces.json').exists():
        raise SystemExit('已有原始收據；請指定新的 workplace/rNNN 目錄，禁止覆寫基準')
    probe = OUT / 'probe'
    probe.mkdir(exist_ok=True)
    files = subprocess.check_output(['git','ls-tree','-r','--name-only',BASELINE,'cmd/wiz7'],text=True).splitlines()
    for fn in files:
        if fn.endswith('.go'):
            (probe/Path(fn).name).write_bytes(subprocess.check_output(['git','show',f'{BASELINE}:{fn}']))
    p = probe / 'main.go'
    s = p.read_text()
    assert s.count('after := offscreen(o)\n') == 1
    s = s.replace('after := offscreen(o)\n', 'after := offscreen(o)\n\tstoryProbePrint(o, p, after)\n', 1)
    s = s.replace('a := append([]uint8(nil), o.Indexed()...)',
                  'a := append([]uint8(nil), o.Indexed()...)\n\t_ = os.WriteFile(base+".off", offscreen(o), 0o644)', 1)
    p.write_text(s)
    code = '''package main
import("encoding/json";"fmt";"os";"github.com/wicanr2/dosgolem/oracle")
type storyTrace struct { Step uint64; Page, Y int; Before, After string }
var storyTraces []storyTrace
var storyPage int
func storyProbePrint(o *oracle.Oracle,p *pendingPrint,after []uint8) {
 if p.ret != (oracle.Addr{Seg:0x0110,Off:0xbe0c}) || p.before==nil || after==nil {return}
 if p.y0==146 {storyPage++}
 base:=fmt.Sprintf("workplace/r673/timing/frag-%d",len(storyTraces)+1)
 if err:=os.WriteFile(base+".before",p.before,0644);err!=nil {panic(err)}
 if err:=os.WriteFile(base+".after",after,0644);err!=nil {panic(err)}
 storyTraces=append(storyTraces,storyTrace{o.Steps(),storyPage,p.y0,base+".before",base+".after"})
 b,err:=json.Marshal(storyTraces);if err!=nil {panic(err)}
 if err:=os.WriteFile("workplace/r673/timing/traces.json",b,0644);err!=nil {panic(err)}
}
'''
    (probe / 'story_probe.go').write_text(code.replace('workplace/r673', str(OUT)))
    samples = sorted({pg[-1]+d*1_000_000 for pg in pages() for d in [1, 3, 6, 8]})
    (OUT / 'samples.txt').write_text(','.join(map(str, samples)))
    (OUT / 'timing').mkdir(exist_ok=True)
    inputs = ['workplace/run/DS.EXE', 'workplace/run/MSG.DBS', 'text/messages.json',
              'text/story-trans.json', 'workplace/bake/story.txt', str(LOG),
              'workplace/bake/r667-narr-msgs.txt', 'cmd/ovlprobe/main.go', 'cmd/wiz7/main.go',
              'tools/re/story_evidence.py', str(probe/'main.go'), str(probe/'story_probe.go')]
    hashes = {fn: hashlib.sha256(Path(fn).read_bytes()).hexdigest() for fn in inputs}
    (OUT / 'input-hashes.json').write_text(json.dumps(hashes, ensure_ascii=False, indent=2))
    print('31頁、124個原版取樣點已準備；副本僅供觀察')


def analyze():
    traces = json.loads((OUT / 'timing/traces.json').read_text())
    assert len(traces) == 71 and traces[-1]['Page'] == 31
    groups = [[t for t in traces if t['Page'] == n] for n in range(1, 32)]
    reads = [(int(s), int(i)) for s, i in re.findall(r'訊息 (\d+) id=(\d+).*ret=bfeb', LOG.read_text())]
    records = {r['id']: r['en'] for r in json.loads(Path('text/messages.json').read_text())['records'] if r['id'] >= 31500}
    fragments = [json.loads(m) for m in re.findall(r'印字 \d+ ("(?:[^"\\]|\\.)*") font=2 ret=0110:be0c', LOG.read_text())]
    def words(s):
        s = ''.join(c for c in s if c >= ' ')
        if s.startswith(('@', '#')):
            s = s[1:]
        return re.sub(r'\s*(?:&[%#|\[]*)+\]$', '', s).split()
    lower, pos = 0, 0
    for entries, group in zip(pages(), groups):
        assert len(entries) == len(group)
        assert all(t['Step'] > entry for t, entry in zip(group, entries))
        ids = [i for s,i in reads if lower < s < entries[-1] and i in records]
        expected = sum((words(records[i]) for i in ids), [])
        observed = sum((s.split() for s in fragments[pos:pos+len(group)]), [])
        assert expected == observed, group[0]['Page']
        lower = group[-1]['Step']
        pos += len(group)
    assert pos == 71
    for group in groups:
        for t in group:
            t['before'] = Path(t['Before']).read_bytes()
            t['after'] = Path(t['After']).read_bytes()
            assert len(t['before']) == len(t['after']) == 64000
    band = {y*320+x for y in range(146, 195) for x in range(12, 308)}
    counts = {d: {'samples': 0, 'event_match': 0, 'band_match': 0} for d in [1, 3, 6, 8]}
    rows = []
    for pg, group in zip(pages(), groups):
        strokes = {i for t in group for i, (b, a) in enumerate(zip(t['before'], t['after'])) if b != a}
        for d in counts:
            name = OUT / f'timing/hdat-{pg[-1]+d*1_000_000}'
            meta = json.loads(name.with_suffix('.json').read_text())
            a = name.with_suffix('.a000').read_bytes()
            mx, my = meta['Mouse']
            cursor = {y*320+x for y in range(max(0,my), min(200,my+16)) for x in range(max(0,mx), min(320,mx+16))}
            matches = []
            for t in group:
                original_changed = {i for i, (b,c) in enumerate(zip(t['before'], t['after'])) if b != c}
                changed = sorted(original_changed - cursor)
                same = sum(a[i] != t['before'][i] for i in changed)
                xs = [i % 320 for i in original_changed]; ys = [i // 320 for i in original_changed]
                # 對照原版印字前後差異，僅在實測外接框內核對背景；不呼叫 runtime 判定。
                box = {y*320+x for y in range(min(ys),max(ys)+1) for x in range(min(xs),max(xs)+1)}
                background = all(a[i] == t['after'][i] for i in box - original_changed - cursor)
                matches.append(background and same*2 >= len(changed) and bool(changed))
            mismatch = sum(a[i] != group[-1]['after'][i] for i in band - strokes - cursor)
            c = counts[d]; c['samples'] += 1; c['event_match'] += all(matches); c['band_match'] += mismatch == 0
            colors = sorted({a[i] for i in strokes-cursor})
            rows.append({'page':group[0]['Page'], 'delayM':d, 'requested':pg[-1]+d*1_000_000,
                         'step':meta['Step'], 'lastReturn':group[-1]['Step'],
                         'eventMatch':all(matches), 'bandMismatch':mismatch, 'strokeColors':colors})
    assert len(rows) == 124
    report = {'counts':counts, 'samples':rows}
    (OUT / 'timing-report.json').write_text(json.dumps(report, ensure_ascii=False, indent=2))
    print(json.dumps(counts, ensure_ascii=False))
    print('前頁實際返回下界：31頁來源與可見字詞完整相等')
    print('逐頁末片段返回步數:', ','.join(str(g[-1]['Step']) for g in groups))


def transient():
    """最後故事頁的帶內暫態：與前頁筆畫及實作前原版樣本交叉驗證。"""
    source = Path('workplace/r673/timing/traces.json')
    traces = json.loads(source.read_text())
    last_page = max(t['Page'] for t in traces)
    strokes = {p:set() for p in [last_page-1,last_page]}
    for t in traces:
        if t['Page'] in strokes:
            b,a = Path(t['Before']).read_bytes(),Path(t['After']).read_bytes()
            strokes[t['Page']].update(i for i,(x,y) in enumerate(zip(b,a)) if x!=y)
    after = Path(traces[-1]['After']).read_bytes()
    entry = pages()[-1][-1]
    rows=[]
    for delay in [1,3,6,8]:
        name=f'hdat-{entry+delay*1_000_000}'
        base=OUT/'en'/name
        meta=json.loads(base.with_suffix('.json').read_text())
        live=base.with_suffix('.a000').read_bytes()
        mx,my=meta['Mouse']
        bad={y*320+x for y in range(146,195) for x in range(12,308)
             if y*320+x not in strokes[last_page] and not(mx<=x<mx+16 and my<=y<my+16)
             and live[y*320+x]!=after[y*320+x]}
        old=Path('workplace/r673/timing')/(name+'.a000')
        rows.append({'requested':entry+delay*1_000_000,'step':meta['Step'],'mismatches':len(bad),
                     'previousPageStrokeOverlap':len(bad&strokes[last_page-1]),
                     'sameAsOriginalBaseline':live==old.read_bytes(),
                     'colors':{f'{x}->{y}':n for (x,y),n in Counter((after[i],live[i]) for i in bad).items()},
                     'zhFlow':bool(json.loads((OUT/'zh'/(name+'.json')).read_text()).get('Flow'))})
    report={'addressSpace':'dosgolem A000／離屏，320×200 色號陣列',
            'tracesSHA256':hashlib.sha256(source.read_bytes()).hexdigest(),'samples':rows}
    (OUT/'last-page-transient.json').write_text(json.dumps(report,ensure_ascii=False,indent=2))
    print(json.dumps(report,ensure_ascii=False))


if __name__ == '__main__':
    {'prepare': prepare, 'analyze': analyze, 'transient': transient}[sys.argv[1]]()
