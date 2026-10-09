#!/usr/bin/env python3
"""037限定獨立驗收；只讀witness、原字型、原版摘要，不讀wire/DB/正式字模。

Docker內執行，--root唯讀收據根目錄，--out為唯一可寫新JSON。
預期字形直接取固定原字型；英文mask直接從原76×13圖片的已證文字帶計算。
不匯入受測程式或烘製模組，也不以proof的mask/glyphs當期望。
"""
import argparse
import gzip
import hashlib
import json
import re
from pathlib import Path

from PIL import Image


FONT_SHA = {
    'eten': '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701',
    'unifont': '7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574',
}
PICTURES = {
    '98e96e235de0701c61b2629c6a78bc6ba1f49fced39efad81e094d8ca498e863':
        ('LONGSWORD', '長劍', (114, 0), (0, 110)),
    '5922e142ae731e2e041ea1f0aa16cb869e345bd60b00f954f413b4479ffa532c':
        ('LEATHER/CUIRASS', '皮胸甲', (89, 85), (0, 124)),
}
EXPECTED_TRACE = 'f6b25ba4beb169cdc174d975dc3c7aecb0fbb6efd0925e8ad3ca239eec047dec'
EXPECTED_WIRE = 'b413c31a633943eb97b9d19aaeac4fd74f06ec9927bcb1761733b9f88ab9dd4c'
SOURCE_REVIEW = 'f2ca9fcabbd55430de7f0cc27e7ddeafaa07884f9fa732d081e8e2226878b8d9'
POSITIONS = {(0, 110), (0, 124), (100, 145), (200, 100), (300, 198)}


def sha(data):
    return hashlib.sha256(data).hexdigest()


class Audit:
    def __init__(self):
        self.checks = 0
        self.failures = []
        self.files = {}

    def check(self, ok, label):
        self.checks += 1
        if not ok:
            self.failures.append(label)

    def read(self, path):
        data = path.read_bytes()
        self.files[str(path)] = sha(data)
        return data

    def load(self, path):
        return json.loads(self.read(path))


def fonts(paths, audit):
    result = {}
    wanted = set('長劍皮胸甲')
    for backend, path in paths.items():
        data = audit.read(path)
        audit.check(sha(data) == FONT_SHA[backend], f'{backend}固定原字型SHA')
        glyphs = {}
        if backend == 'eten':
            audit.check(len(data) == 392820, '倚天裸字型392820bytes')
            for ch in wanted | {'一'}:
                hi, lo = ch.encode('big5')
                # 按Big5兩個字元區分別計數，避開C6A1–C8FE的408個保留碼位。
                col = lo - (0x40 if lo <= 0x7e else 0x62)
                index = (hi - 0xa4) * 157 + col - (408 if hi >= 0xc9 else 0)
                assert 0 <= index < 13094
                raw = data[index * 30:(index + 1) * 30]
                glyphs[ch] = tuple(tuple(bool(int.from_bytes(raw[y*2:y*2+2], 'big') & (0x8000 >> x))
                                         for x in range(15)) for y in range(15))
            one = [(x, y) for y, row in enumerate(glyphs.pop('一')) for x, b in enumerate(row) if b]
            audit.check(set(one) == {(x, 7) for x in range(15)} | {(13, 6)},
                        '倚天A440正對照為15點橫線及右端襯線')
        else:
            for line in gzip.decompress(data).decode('ascii').splitlines():
                cp, bits = line.split(':')
                ch = chr(int(cp, 16))
                if ch not in wanted:
                    continue
                assert ch not in glyphs and len(bits) == 64
                rows = [int(bits[n:n+4], 16) for n in range(0, 64, 4)]
                glyphs[ch] = tuple(tuple(bool(row & (0x8000 >> x)) for x in range(16)) for row in rows)
        audit.check(set(glyphs) == wanted, f'{backend}恰取5個所需字')
        result[backend] = glyphs
    return result


def english_mask(picture, expected, audit, label):
    assert len(picture) == 988
    mask = bytearray(988)
    counts = []
    for ys in (range(1, 6), range(7, 12)):
        count = 0
        for y in ys:
            for x in range(16, 76):
                i = y * 76 + x
                if picture[i] != 255:
                    mask[i] = 1
                    count += 1
                    audit.check(picture[i] == 105, f'{label}原英文前景索引105@{x},{y}')
        counts.append(count)
    audit.check(tuple(counts) == expected, f'{label}獨立原英文墨點數{expected}')
    return bytes(mask)


def proof(directory, path, audit):
    info = audit.load(path)
    tag = path.stem
    raw = audit.read(directory / (tag + '-picture.bin'))
    identity = sha(raw)
    assert identity in PICTURES
    key, zh, counts, position = PICTURES[identity]
    audit.check(info['key'] == key and info['picture_sha256'] == identity, tag + '完整原圖與名稱身份')
    audit.check(info['module_valid'] is True, tag + '已定位來源版本')
    audit.check(info['source_step'] <= info['captured_step'] < info['defined_step']
                and info['defined_step'] - info['source_step'] <= 1000000, tag + 'source/capture/define步序')
    rect = info['source_rect']
    audit.check(rect == {'Min': {'X': 2, 'Y': position[1]+3}, 'Max': {'X': 78, 'Y': position[1]+16}},
                tag + '來源抓圖矩形')
    own_mask = english_mask(raw, counts, audit, tag)
    audit.check(audit.read(directory / (tag + '-mask.bin')) == own_mask, tag + '註冊mask等於獨立英文帶')
    part_count = 1 if key == 'LONGSWORD' else 2
    audit.check(set(info['files']) == {f'part{i}-{r}' for i in range(part_count)
                                      for r in ('before', 'after', 'mask')}, tag + '完整單/雙行來源')
    previous = None
    joined = bytearray(988)
    for n in range(part_count):
        values = {}
        for role in ('before', 'after', 'mask'):
            name = info['files'][f'part{n}-{role}']
            assert Path(name).name == name
            values[role] = audit.read(directory / name)
            assert len(values[role]) == 988
        before, after, recorded_mask = (values[r] for r in ('before', 'after', 'mask'))
        if previous is not None:
            audit.check(before == previous, tag + '第二行before等於第一行after')
        band = {y*76+x for y in range(1+6*n, 6+6*n) for x in range(16, 76)}
        audit.check(all(before[i] == 255 for i in band), f'{tag}第{n+1}行完整文字帶全透明')
        changed = bytes(int(a != b) for a, b in zip(before, after))
        expected_mask = bytes(int(i in band and after[i] != 255) for i in range(988))
        audit.check(changed == expected_mask == recorded_mask,
                    f'{tag}第{n+1}行全域差分只含本行英文墨點')
        audit.check(sum(changed) == counts[n] and all(after[i] == 105 for i in band if changed[i]),
                    f'{tag}第{n+1}行前景與墨點數')
        joined = bytearray(a | b for a, b in zip(joined, changed))
        previous = after
    audit.check(all(previous[i] == raw[i] for i in range(988) if own_mask[i]),
                tag + '來源名稱墨點在最後抓圖完整保持')
    # C4允許印字完成後圖示補畫；名稱以外不能拿早期after當capture期望。
    non_name_changes = sum(a != b for i, (a, b) in enumerate(zip(previous, raw)) if not own_mask[i])
    audit.check(bytes(joined) == own_mask, tag + '各行差分聯集等於獨立名稱mask')
    return {'file': path.name, 'key': key, 'english_ink': sum(own_mask),
            'source_step': info['source_step'], 'captured_step': info['captured_step'],
            'defined_step': info['defined_step'], 'picture_sha256': identity,
            'after_to_capture_non_name_changes': non_name_changes}


def frame(directory, path, backend, scale, glyphs, audit):
    info = audit.load(path)
    tag = path.stem
    data = {r: audit.read(directory / f'{tag}-{r}.bin') for r in ('picture', 'background', 'indexed', 'palette')}
    picture, background, indexed, palette = (data[r] for r in ('picture', 'background', 'indexed', 'palette'))
    identity = sha(picture)
    assert identity in PICTURES
    key, zh, counts, initial_position = PICTURES[identity]
    x, y, w, h = (info[n] for n in ('x', 'y', 'w', 'h'))
    audit.check((x, y) in POSITIONS and w == min(76, 320-x) and h == min(13, 200-y), tag + '原描述器幾何與裁切')
    assert len(background) == len(indexed) == w*h and len(palette) == 768
    mask = english_mask(picture, counts, audit, tag)
    reconstructed = bytearray(background)
    for yy in range(h):
        for xx in range(w):
            color = picture[yy*76+xx]
            if color != 255:
                reconstructed[yy*w+xx] = color
    original_diff = sum(a != b for a, b in zip(reconstructed, indexed))
    audit.check(original_diff == 0, tag + '背景加原圖逐byte等於真正IndexedROI')
    audit.check(info['scale'] == scale and info['picture_sha256'] == identity and info['module_valid'] is True,
                tag + '倍率/原圖/版本身份')
    translated = info['translated'] is True
    audit.check((translated and info['key'] == key) or (not translated and info['key'] == ''),
                tag + '無錯譯或帶錯名稱')
    size = (w*scale, h*scale)
    rgb = [tuple(palette[c*3:c*3+3]) for c in background]
    expected = Image.new('RGB', (w, h))
    expected.putdata(rgb)
    expected = expected.resize(size, Image.Resampling.NEAREST)
    pixels = expected.load()
    for yy in range(h):
        for xx in range(w):
            pos = yy*76+xx
            color = picture[pos]
            if color == 255 or translated and mask[pos]:
                continue
            for dy in range(scale):
                for dx in range(scale):
                    pixels[xx*scale+dx, yy*scale+dy] = tuple(palette[color*3:color*3+3])
    native_ink = 0
    if translated:
        fg = tuple(palette[315:318])
        advance = 17 if backend == 'eten' else 18
        for n, ch in enumerate(zh):
            for gy, row in enumerate(glyphs[ch]):
                for gx, active in enumerate(row):
                    if not active:
                        continue
                    px, py = 16*scale+n*advance+gx, scale+gy
                    native_ink += 1
                    # 獨立查中文字點所屬DOS像素，不能撞上原圖非名稱圖素。
                    dos_x, dos_y = px//scale, py//scale
                    audit.check(dos_x < 76 and dos_y < 13 and
                                (picture[dos_y*76+dos_x] == 255 or mask[dos_y*76+dos_x]),
                                f'{tag}中文墨點沒有撞上原圖圖示@{px},{py}')
                    if 0 <= px < size[0] and 0 <= py < size[1]:
                        pixels[px, py] = fg
    paths = {r: directory / f'{tag}-{r}.png' for r in ('draw', 'hd')}
    actual = {}
    for role, image_path in paths.items():
        audit.read(image_path)
        with Image.open(image_path) as image:
            audit.check(image.size == (320*scale, 200*scale), f'{tag}{role}正式遊戲畫布尺寸')
            audit.check(image.convert('RGBA').getchannel('A').getextrema() == (255, 255),
                        f'{tag}{role}整張畫布不透明')
            actual[role] = image.convert('RGB')
    entire_diff = sum(a != b for a, b in zip(actual['draw'].getdata(), actual['hd'].getdata()))
    audit.check(entire_diff == 0, tag + '實際Draw與HD全圖逐pixel相同')
    rectangle = (x*scale, y*scale, (x+w)*scale, (y+h)*scale)
    roi = actual['draw'].crop(rectangle)
    roi_diff = sum(a != b for a, b in zip(roi.getdata(), expected.getdata()))
    audit.check(roi_diff == 0, tag + '實際Draw游標ROI等於獨立原圖/背景/原生中文字形')
    return {'file': path.name, 'pub': info['pub'], 'step': info['step'], 'key': key,
            'actual_key': info['key'], 'translated': translated, 'x': x, 'y': y, 'w': w, 'h': h,
            'original_pixels': w*h, 'original_differences': original_diff,
            'rgb_pixels': w*h*scale*scale, 'roi_differences': roi_diff,
            'draw_hd_differences': entire_diff, 'english_ink': sum(mask), 'native_ink': native_ink,
            'rejected_band': info['rejected_band'], 'draw': str(paths['draw'])}


def configuration(root, backend, scale, glyphs, audit, baseline, prefix):
    directory = root / f'{prefix}-{backend}-{scale}'
    result = audit.load(directory / 'result.json')
    audit.check(result['font'] == backend and result['scale'] == scale and result['exit_code'] == 0,
                f'{backend}/{scale}正常結束')
    audit.check(result['initial'] == result['final'], f'{backend}/{scale}兩原版保存metadata不變')
    audit.check(result['wire_sha256'] == EXPECTED_WIRE, f'{backend}/{scale}固定正常movement來源')
    trace = audit.read(directory / 'states.jsonl')
    audit.check(trace == baseline and sha(trace) == result['trace_sha256'] == EXPECTED_TRACE,
                f'{backend}/{scale}全部1651原版六類摘要相同')
    rows = [json.loads(line) for line in trace.splitlines()]
    audit.check(len(rows) == result['publications'] == 1651 and
                [r['publication'] for r in rows] == list(range(1, 1652)), f'{backend}/{scale}全部發布完整且連續')
    log = audit.read(directory / 'run.log')
    audit.check(b'INVALID' not in log and b'panic:' not in log, f'{backend}/{scale}無INVALID/panic')
    witness = directory / 'witness'
    proofs = [proof(witness, p, audit) for p in sorted(witness.glob('proof-*.json'))]
    audit.check({p['key'] for p in proofs} == {'LONGSWORD', 'LEATHER/CUIRASS'}, f'{backend}/{scale}兩件來源證據')
    frames = [frame(witness, p, backend, scale, glyphs, audit) for p in sorted(witness.glob('frame-*.json'))]
    stable = {(r['key'], r['x'], r['y']) for r in frames if r['translated']}
    required = {('LONGSWORD', x, y) for x, y in ((0, 110), (100, 145), (200, 100), (300, 198))}
    required |= {('LEATHER/CUIRASS', x, y) for x, y in ((0, 124), (100, 145))}
    audit.check(required <= stable, f'{backend}/{scale}兩件選取/移動及右下裁切皆穩定中文')
    transient = []
    for r in frames:
        if r['translated']:
            continue
        successors = [s for s in frames if s['translated'] and (s['key'], s['x'], s['y']) == (r['key'], r['x'], r['y'])
                      and r['step'] < s['step'] <= r['step'] + 1000000]
        audit.check(bool(successors), f"{backend}/{scale}/{r['file']}暫態英文在1M步內同位置轉為正確中文")
        transient.append(r['file'])
    resets = [audit.load(p) for p in sorted(witness.glob('reset-*.json'))]
    audit.check(len(resets) == 2 and all(not r['label_present'] and not r['view_present'] for r in resets),
                f'{backend}/{scale}兩件放回撤除label及游標view')
    audit.check(all(1 <= r['pub'] <= 1651 and rows[r['pub']-1]['step'] == r['step'] for r in frames+resets),
                f'{backend}/{scale}每張畫面連結真正原版發布步')
    for path in witness.glob('reset-*-draw.png'):
        audit.read(path)
    return {'font': backend, 'scale': scale, 'binary_sha256': result['binary_sha256'],
            'trace_sha256': sha(trace), 'publications': len(rows), 'saves_unchanged': result['initial'] == result['final'],
            'proofs': proofs, 'frames': frames, 'transient_english': transient, 'resets': resets}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--root', type=Path, required=True)
    parser.add_argument('--eten', type=Path, required=True)
    parser.add_argument('--unifont', type=Path, required=True)
    parser.add_argument('--out', type=Path, required=True)
    parser.add_argument('--prefix', default='formal-verify')
    parser.add_argument('--source-review', type=Path, default=Path('formal-source-review-r2.json'))
    parser.add_argument('--source-review-sha', default=SOURCE_REVIEW)
    parser.add_argument('--font', choices=('eten', 'unifont'))
    parser.add_argument('--scale', type=int, choices=(2, 3))
    args = parser.parse_args()
    assert re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9_.-]*', args.prefix)
    assert re.fullmatch(r'[0-9a-f]{64}', args.source_review_sha)
    assert not args.out.exists() and not args.out.with_suffix('.md').exists()
    audit = Audit()
    source = audit.read(Path(__file__))
    glyphs = fonts({'eten': args.eten, 'unifont': args.unifont}, audit)
    baseline = audit.read(args.root / 'movement-api-control/states.jsonl')
    audit.check(sha(baseline) == EXPECTED_TRACE, '先前原版對照trace固定雜湊')
    review_path = args.source_review if args.source_review.is_absolute() else args.root / args.source_review
    review = audit.read(review_path)
    audit.check(sha(review) == args.source_review_sha, '逐字NUL與控制流來源審查固定收據')
    configs = []
    for backend in ((args.font,) if args.font else ('unifont', 'eten')):
        for scale in ((args.scale,) if args.scale else (2, 3)):
            configs.append(configuration(args.root, backend, scale, glyphs[backend], audit, baseline, args.prefix))
    audit.check(len({c['binary_sha256'] for c in configs}) == 1, '所有限定組使用同一受測binary')
    frames = [r for c in configs for r in c['frames']]
    report = {'status': 'PASS' if not audit.failures else 'FAIL', 'checks': audit.checks, 'failures': audit.failures,
              'tool_sha256': sha(source), 'font_sha256': FONT_SHA, 'source_review_sha256': args.source_review_sha,
              'source_review_path': str(review_path), 'configuration_prefix': args.prefix,
              'configurations': configs, 'frame_count': len(frames),
              'original_pixels': sum(r['original_pixels'] for r in frames),
              'original_differences': sum(r['original_differences'] for r in frames),
              'rgb_pixels': sum(r['rgb_pixels'] for r in frames),
              'roi_differences': sum(r['roi_differences'] for r in frames),
              'draw_hd_differences': sum(r['draw_hd_differences'] for r in frames),
              'source_files': audit.files,
              'limits': ['僅正常movement兩物品及列明座標，不外推全物品/全遊戲。',
                         'proof驗像素鏈及步序；逐字NUL身份與capture到definition完整988守門沿指定的固定雜湊獨立正式來源審查，proof未分別匯出capture及definition副本。',
                         '原生字型直接取固定原檔，未使用受測glyphs/mask/renderer。',
                         '私人API只供同狀態診斷，不宣稱跨binary正式wire相容；新正常錄製重播由主代理另驗。',
                         '本工具不執行GUI操作，不代替正常GUI試玩或人眼確認。']}
    args.out.write_text(json.dumps(report, ensure_ascii=False, indent=2)+'\n')
    lines = [f"# R1000 獨立像素驗收：{report['status']}", '',
             f"{len(configs)}組，{len(frames)}張畫面，{audit.checks}項檢查；失敗{len(audit.failures)}項。", '',
             '|字型|倍率|原版發布|畫面|原圖ROI差異|中文ROI差異|Draw/HD差異|', '|---|---:|---:|---:|---:|---:|---:|']
    for c in configs:
        fs = c['frames']
        lines.append(f"|{c['font']}|{c['scale']}|{c['publications']}|{len(fs)}|{sum(r['original_differences'] for r in fs)}|{sum(r['roi_differences'] for r in fs)}|{sum(r['draw_hd_differences'] for r in fs)}|")
    lines += ['', '英文暫態：' + json.dumps({f"{c['font']}/{c['scale']}": c['transient_english'] for c in configs}, ensure_ascii=False), '', '失敗：']
    lines += ['- '+x for x in audit.failures] if audit.failures else ['無。']
    lines += ['', '限制：', ''] + ['- '+x for x in report['limits']]
    args.out.with_suffix('.md').write_text('\n'.join(lines)+'\n')
    print(json.dumps({k: report[k] for k in ('status', 'checks', 'failures', 'frame_count', 'original_pixels', 'original_differences', 'rgb_pixels', 'roi_differences', 'draw_hd_differences')}, ensure_ascii=False))
    return int(bool(audit.failures))


if __name__ == '__main__':
    raise SystemExit(main())
