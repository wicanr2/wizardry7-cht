#!/usr/bin/env python3
"""正常GUI的窄像素量測，交另一reviewer採認；不執行原版或讀wire。

期望只取R1002 source_extension的原M8×8/13點與直接原字型。
正常PNG沒有publication，私有frame1520只提供palette來源；實際核非中文字背景
RGB，不能把同pub當正常PNG的像素oracle。接手圖只核身份，另由root人工判讀。
"""
import argparse
import base64
import gzip
import hashlib
import json
from pathlib import Path
import stat
import sys

from PIL import Image


FONT_SHA = {'eten': '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701',
            'unifont': '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574'}
CASES = [('gui-unifont-3', 'unifont', 3, 'creation-bonus-str.png'),
         ('gui-eten-3', 'eten', 3, 'creation-bonus-str.png'),
         ('gui-unifont-2-handoff-r2', 'unifont', 2, 'creation-str-two.png'),
         ('gui-eten-2-handoff', 'eten', 2, 'creation-str-two.png')]
BINARY_SHA = '631926caa447010014ff70dada2dc56179a010eaa6fd5573394084b511484c99'
POINTS = [[88,4],[92,4],[95,4],[88,7],[89,9],[88,11],[91,5],[95,5]]


class InputError(Exception):
    pass


class Measure:
    def __init__(self):
        self.checks = 0
        self.failures = []
        self.sources = {}

    def check(self, ok, label):
        self.checks += 1
        if not ok:
            self.failures.append(label)

    def read(self, path):
        data = path.read_bytes()
        self.sources[str(path)] = hashlib.sha256(data).hexdigest()
        return data

    def load(self, path):
        return json.loads(self.read(path))


def fonts(eten, unifont, measure):
    result = {}
    for backend, path in [('eten', eten), ('unifont', unifont)]:
        data = measure.read(path)
        if hashlib.sha256(data).hexdigest() != FONT_SHA[backend]:
            raise InputError('fixed_font_sha')
        if backend == 'eten':
            if len(data) != 392820:
                raise InputError('eten_length')
            hi, lo = '男'.encode('big5')
            col = lo-(0x40 if lo<=0x7e else 0x62)
            index = (hi-0xa4)*157+col-(408 if hi>=0xc9 else 0)
            if not 0<=index<13094:
                raise InputError('eten_index')
            raw = data[index*30:(index+1)*30]
            result[backend] = tuple(tuple(bool(int.from_bytes(raw[y*2:y*2+2], 'big') & (0x8000>>x))
                                          for x in range(15)) for y in range(15))
            one = data[:30]
            measure.check(one == bytes.fromhex('0000000000000000000000000004fffe0000000000000000000000000000'),
                          '倚天A440一字索引正對照')
        else:
            selected = []
            for line in gzip.decompress(data).decode('ascii').splitlines():
                cp, bits = line.split(':')
                if int(cp,16)==ord('男'):
                    selected.append(bits)
            if len(selected)!=1 or len(selected[0])!=64:
                raise InputError('unifont_native_glyph')
            rows=[int(selected[0][n:n+4],16) for n in range(0,64,4)]
            result[backend]=tuple(tuple(bool(row&(0x8000>>x)) for x in range(16)) for row in rows)
    return result


def original(path, measure):
    events=[json.loads(line) for line in measure.read(path).splitlines()]
    choices=[r for r in events if r['status']=='source_extension']
    if len(choices)!=1:
        raise InputError('single_r1002_extension')
    row=choices[0]
    before=base64.b64decode(row['before8'],validate=True)
    after=base64.b64decode(row['after8'],validate=True)
    if len(before)!=64 or len(after)!=64 or row['eight_rectangle']!=[88,4,96,12]:
        raise InputError('original_eight_shape')
    mask=bytes(int(a!=b) for a,b in zip(before,after))
    measure.check(row['code']==77 and row['return']==0x9595 and row['font']==0, 'R1002原M來源身份')
    measure.check(sum(mask)==13 and all(after[i]==1 and i%8>=2 and i//8>=2 for i in range(64) if mask[i]),
                  '原M獨立13點／色1／6×6域')
    measure.check(row['background_points']==POINTS, 'R1002八背景點座標')
    points=bytes(before[(y-4)*8+x-88] for x,y in POINTS)
    measure.check(points==bytes([115,116,115,117,115,115,13,14]), 'R1002八背景點值')
    return before, after, mask


def fnv(data):
    h=2166136261
    for value in data:
        h=((h^value)*16777619)&0xffffffff
    return h


def palette(root, backend, scale, measure):
    directory=root/f'api-old-{backend}-{scale}'
    result=measure.load(directory/'result.json')
    measure.check(result['case']==f'old-{backend}-{scale}' and result['exit_code']==0 and not result['invalid'],
                  f'{backend}/{scale}指定原版palette私有控制完成')
    data=measure.read(directory/'capture/frame-1520/hd.pal')
    if len(data)!=768:
        raise InputError('palette_length')
    hd=measure.load(directory/'capture/frame-1520/hd.json')
    measure.check(hd['Step']==338800000 and fnv(data)==hd['Pal'], f'{backend}/{scale}原palette bytes及既有摘要')
    return [tuple(data[i*3:i*3+3]) for i in range(256)]


def saved_metadata(case, result, measure):
    directory=case/'record-gui/save'
    measure.check(set(result['initial'])=={'SAVEGAME.DBS','PCFILE.DBS'}, case.name+'兩保存身份')
    measure.check(result['initial']==result['final']==result['source_after'], case.name+'宣告保存/source metadata保持')
    actual={}
    for name in ('SAVEGAME.DBS','PCFILE.DBS'):
        path=directory/name
        data=measure.read(path)
        s=path.stat()
        actual[name]={'size':s.st_size,'mode':stat.S_IMODE(s.st_mode),'mtime_ns':str(s.st_mtime_ns),
                      'sha256':hashlib.sha256(data).hexdigest()}
        measure.check((s.st_uid,s.st_gid)==(1000,1000), case.name+'/'+name+'保存owner')
    measure.check(actual==result['final'], case.name+'實際保存bytes/mode/mtime與收據相同')
    return actual


def png(path, k, measure):
    measure.read(path)
    with Image.open(path) as image:
        image.load()
        if image.size!=(320*k,200*k+40):
            raise InputError('normal_gui_body_footer_shape')
        rgba=image.convert('RGBA')
    measure.check(rgba.getchannel('A').getextrema()==(255,255), path.name+'整圖不透明')
    return rgba.crop((0,0,320*k,200*k)).convert('RGB')


def expected(before, after, mask, native, pal, k):
    # 從原after擦回13個真正M差分，再整數倍放大；不使用候選G或renderer。
    restored=bytearray(after)
    for i,active in enumerate(mask):
        if active:
            restored[i]=before[i]
    if bytes(restored)!=before:
        raise InputError('original_mask_restoration')
    size=len(native);shift=max(size-6*k,0)
    x0,y0=2*k-shift,2*k-shift
    glyph={(x0+x,y0+y) for y,row in enumerate(native) for x,active in enumerate(row) if active}
    if x0<0 or y0<0 or x0+size>8*k or y0+size>8*k:
        raise InputError('native_glyph_clipped')
    pixels=[];indices=[]
    for y in range(8*k):
        for x in range(8*k):
            index=restored[(y//k)*8+x//k]
            indices.append(index)
            pixels.append(pal[1] if (x,y) in glyph else pal[index])
    obsolete={(x,y) for y in range(8*k) for x in range(8*k)
              if mask[(y//k)*8+x//k] and (x,y) not in glyph}
    return pixels,indices,glyph,obsolete,(90*k-shift,6*k-shift),size


def sample(path,k,backend,before,after,mask,native,pal,measure):
    body=png(path,k,measure)
    box=(88*k,4*k,96*k,12*k)
    actual=list(body.crop(box).getdata())
    want,indices,glyph,obsolete,position,size=expected(before,after,mask,native,pal,k)
    diffs=[i for i,(a,b) in enumerate(zip(actual,want)) if a!=b]
    bg_diffs=[];glyph_diffs=[];english_diffs=[];coverage={}
    for i,(a,b) in enumerate(zip(actual,want)):
        point=(i%(8*k),i//(8*k))
        if point in glyph:
            if a!=b:glyph_diffs.append(i)
        else:
            index=indices[i];coverage[str(index)]=coverage.get(str(index),0)+1
            if a!=b:bg_diffs.append(i)
        if point in obsolete and a!=b:english_diffs.append(i)
    label=f'{backend}/{k}/{path.name}'
    measure.check(not diffs, label+'8×8 ROI逐RGB符合原背景／原M遮罩／原生男字')
    measure.check(not bg_diffs, label+'實際非中文字背景RGB核當前palette相關色號等值')
    measure.check(not glyph_diffs, label+'完整原生男字墨點與色號1')
    measure.check(not english_diffs, label+'未被中文字覆蓋的英文M墨點已還原')
    return {'path':str(path),'body_rectangle':[0,0,320*k,200*k],'roi_rectangle':box,'roi_rgb_pixels':len(want),
            'different_pixels':len(diffs),'background_pixels':len(want)-len(glyph),'background_differences':len(bg_diffs),
            'background_palette_index_pixel_counts':coverage,'native_glyph_pixels':len(glyph),'native_glyph_differences':len(glyph_diffs),
            'native_glyph_size':size,'native_glyph_position':position,'glyph_clipped':False,
            'obsolete_english_only_pixels':len(obsolete),'obsolete_english_differences':len(english_diffs),
            'actual_roi_rgb_sha256':hashlib.sha256(bytes(v for rgb in actual for v in rgb)).hexdigest(),
            'expected_roi_rgb_sha256':hashlib.sha256(bytes(v for rgb in want for v in rgb)).hexdigest(),
            'normal_publication_mapping_used':False}


def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--root',type=Path,required=True)
    p.add_argument('--source-events',type=Path,required=True)
    p.add_argument('--eten',type=Path,required=True)
    p.add_argument('--unifont',type=Path,required=True)
    p.add_argument('--out',type=Path,required=True)
    a=p.parse_args()
    assert a.out.name=='normal-gui-pixel-measurement.json' and not a.out.exists() and not a.out.with_suffix('.md').exists()
    m=Measure();m.read(Path(__file__));cases=[];input_error=None
    try:
        native=fonts(a.eten,a.unifont,m)
        before,after,mask=original(a.source_events,m)
        for name,backend,k,last in CASES:
            case=a.root/name;result=m.load(case/'gui-result.json')
            m.check(result['exit_code']==0 and result['complete_record'] is True and result['partial_exists'] is False,
                    name+'限定成功組／完整record metadata')
            m.check(result['font']==backend and result['scale']==k and result['binary_sha256']==BINARY_SHA,
                    name+'字型／倍率／binary身份')
            saved=saved_metadata(case,result,m)
            pal=palette(a.root,backend,k,m)
            images=[sample(case/'record-gui'/filename,k,backend,before,after,mask,native[backend],pal,m)
                    for filename in ('creation-male.png','creation-fighter.png',last)]
            handoff=case/'record-gui/creation-portrait-handoff.png';body=png(handoff,k,m)
            cases.append({'case':name,'font':backend,'scale':k,'normal_metadata_healthy':result['exit_code']==0 and result['complete_record'],
                          'actual_saved_metadata':saved,'early_samples':images,
                          'handoff_identity':{'path':str(handoff),'image_shape':[320*k,200*k+40],
                          'body_shape':list(body.size),'geometry_oracle_applied':False,'root_human_review':'reported by root; not independently asserted by this tool'}})
    except (OSError,ValueError,KeyError,TypeError,IndexError,InputError) as exc:
        input_error=type(exc).__name__;m.check(False,'量測輸入/格式受阻：'+input_error)
    samples=[s for c in cases for s in c['early_samples']]
    report={'measurement_status':'MATCH' if not m.failures and len(cases)==4 else 'MISMATCH_OR_INCOMPLETE',
            'independent_final_acceptance':False,'needs_other_reviewer':True,'runtime_started_by_tool':False,
            'candidate_go_or_baked_glyph_used_as_oracle':False,'checks':m.checks,'failures':m.failures,'input_error_type':input_error,
            'tool_sha256':m.sources.get(str(Path(__file__))), 'font_sha256':FONT_SHA,'cases':cases,
            'early_png_count':len(samples),'early_rgb_pixels':sum(s['roi_rgb_pixels'] for s in samples),
            'early_different_pixels':sum(s['different_pixels'] for s in samples),
            'background_different_pixels':sum(s['background_differences'] for s in samples),
            'obsolete_english_different_pixels':sum(s['obsolete_english_differences'] for s in samples),
            'input_sha256':m.sources,'assumptions_and_limits':[
                '正常PNG沒有publication；frame1520只供同字型/倍率的原palette來源，不稱正常PNG同pub或同RAM。',
                '早期M的8×8底圖採R1002 source_extension；每張實際非中文字背景RGB均與當前palette映射逐點核對，才支持這個固定背景前提。',
                '中文字墨點與原生完整字格直接取固定原字型；不取候選G、runtime mask或bake輸出作期望。',
                '頭像接手圖只核身份/尺寸並沿root人工觀察，不套早期M幾何，尤其2倍整頁重排。',
                '保存健康只核gui-result及實際保存metadata；完整record的Reader/嚴格重播及原版健康由root另驗，這裡不讀wire。',
                '本工具為量測，候選來源作者不在此自作獨立動態FINAL；結果交另一reviewer採認。']}
    a.out.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
    lines=[f"# R1003 正常GUI像素量測：{report['measurement_status']}",'',
           f"{len(cases)}成功組，{len(samples)}張早期PNG，{report['early_rgb_pixels']}個ROI像素，差異{report['early_different_pixels']}點。",'',
           '|組|早期圖|ROI像素|差異|非字背景差異|殘留英文專屬點差異|','|---|---:|---:|---:|---:|---:|']
    for c in cases:
        ss=c['early_samples'];lines.append(f"|{c['case']}|{len(ss)}|{sum(s['roi_rgb_pixels'] for s in ss)}|{sum(s['different_pixels'] for s in ss)}|{sum(s['background_differences'] for s in ss)}|{sum(s['obsolete_english_differences'] for s in ss)}|")
    lines+=['','失敗：','']+(['- '+s for s in m.failures] if m.failures else ['無量測差異。'])
    lines+=['','前提與限制：','']+['- '+s for s in report['assumptions_and_limits']]
    lines+=['',f"工具SHA-256：{report['tool_sha256']}"]
    a.out.with_suffix('.md').write_text('\n'.join(lines)+'\n')
    print(json.dumps({k:report[k] for k in ('measurement_status','checks','failures','early_png_count','early_rgb_pixels','early_different_pixels','background_different_pixels','obsolete_english_different_pixels','independent_final_acceptance')},ensure_ascii=False))
    return int(report['measurement_status']!='MATCH')


if __name__=='__main__':raise SystemExit(main())
