#!/usr/bin/env python3
"""R1001：兩趟正常有聲路徑的獨立音訊計數及裝置PCM驗收。

只讀指定的audio-final-lines.json、operations.jsonl、device.s16le及列明程式契約。
選用viewport收據只核正常保存/未保存轉向/冷讀檔三張PNG，不讀保存內容。
不讀run.log、DB、答案、完整輸入或存檔；不執行原版，不匯入被測音訊實作。
PCM是44100Hz signed16LE雙聲道。操作區間只沿收集塊標記，不稱精確音符/聲效邊界。
"""
import argparse
from array import array
import hashlib
import json
import math
from pathlib import Path
import re
import sys


RATE = 44100
H_NUMERATOR = 196875030000
H_DENOMINATOR = 65536
PREFILL = 4410
MAX_RING = 88200
MAX_FADE = 221
FORMAL_BINARY = 'c7da1a38a0a690de1529d0ebadf13c3378e1b51bf3063a522a707962c88fe224'
STREAM_FIELDS = ('EndStep SourceStep Frames CanonicalFrames Received Applied Pending Rejected '
                 'Bank1 CSM Unsupported Stopped Incomplete FaultCode').split()
READER_FIELDS = ('Queued MaximumQueued Underflow Overflow Recovery DeviceFrames DeliveredBytes '
                 'SourceFrames LastSource StopQueued StopFade StopDiscard FadeDeliveredBytes '
                 'Stopped ACK EOF').split()
SOURCE_FIELDS = ('Stream Pending Recorded Dropped FaultCode').split()
KINDS = {'menu', 'click', 'pause', 'resume', 'operation_start', 'operation_end', 'close_request', 'process_exit'}


def digest(data):
    return hashlib.sha256(data).hexdigest()


class Audit:
    def __init__(self):
        self.checks = 0
        self.failures = []
        self.identities = {}

    def check(self, ok, label):
        self.checks += 1
        if not ok:
            self.failures.append(label)

    def read(self, path):
        data = path.read_bytes()
        self.identities[str(path)] = digest(data)
        return data


def parse_struct(text, fields):
    result = {}
    position = 0
    while position < len(text):
        m = re.match(r'([A-Za-z][A-Za-z0-9]*):(map\[[^\]]*\]|true|false|[0-9]+|[a-z_]*)($| )', text[position:])
        if not m:
            raise ValueError('safe_counter_format')
        key, value = m.group(1, 2)
        if key in result or key not in fields:
            raise ValueError('safe_counter_field')
        if key == 'Unsupported':
            # 這兩條正常路線要求空；不輸出任何未知map文字。
            result[key] = {} if value == 'map[]' else {'nonempty': True}
        elif value in ('true', 'false'):
            result[key] = value == 'true'
        elif value.isdecimal():
            result[key] = int(value)
        else:
            result[key] = value
        position += m.end()
    if set(result) != set(fields):
        raise ValueError('safe_counter_missing')
    return result


def safe_summary(directory, audit):
    lines = json.loads(audit.read(directory / 'audio-final-lines.json'))
    if not isinstance(lines, list) or len(lines) != 2 or not all(isinstance(s, str) for s in lines):
        raise ValueError('safe_summary_shape')
    m = re.search(r'音訊結束 fault=([a-z_]*) close=([a-z_]*) stream=\{([^}]+)\} reader=\{([^}]+)\} 畫格率=([0-9.]+)$', lines[0])
    n = re.search(r'原版完成Step=([0-9]+) 來源=\{([^}]+)\}；Received為合成器合法接受，送入合成器=Received\+Rejected；停止後未送入=([0-9]+)$', lines[1])
    if not m or not n:
        raise ValueError('safe_summary_format')
    return {'fault': m[1], 'close_fault': m[2], 'stream': parse_struct(m[3], STREAM_FIELDS),
            'reader': parse_struct(m[4], READER_FIELDS), 'fps': float(m[5]),
            'original_end_step': int(n[1]), 'source': parse_struct(n[2], SOURCE_FIELDS), 'stop_unsent': int(n[3])}


class PCMStats:
    def __init__(self):
        self.frames = 0
        self.nonzero_frames = 0
        self.left_nonzero = self.right_nonzero = 0
        self.left_squares = self.right_squares = 0
        self.left_peak = self.right_peak = 0
        self.different_frames = self.max_channel_delta = 0

    def add(self, data):
        samples = array('h')
        samples.frombytes(data)
        if sys.byteorder != 'little':
            samples.byteswap()
        if len(data) % 4:
            raise ValueError('pcm_partial_frame')
        left, right = samples[::2], samples[1::2]
        self.frames += len(left)
        self.left_nonzero += sum(v != 0 for v in left)
        self.right_nonzero += sum(v != 0 for v in right)
        self.nonzero_frames += sum(a != 0 or b != 0 for a, b in zip(left, right))
        self.left_squares += sum(v*v for v in left)
        self.right_squares += sum(v*v for v in right)
        self.left_peak = max(self.left_peak, max(map(abs, left), default=0))
        self.right_peak = max(self.right_peak, max(map(abs, right), default=0))
        self.different_frames += sum(a != b for a, b in zip(left, right))
        self.max_channel_delta = max(self.max_channel_delta, max((abs(a-b) for a, b in zip(left, right)), default=0))

    def merge(self, other):
        for key in ('frames', 'nonzero_frames', 'left_nonzero', 'right_nonzero', 'left_squares', 'right_squares', 'different_frames'):
            setattr(self, key, getattr(self, key)+getattr(other, key))
        for key in ('left_peak', 'right_peak', 'max_channel_delta'):
            setattr(self, key, max(getattr(self, key), getattr(other, key)))

    def report(self):
        def channel(squares, nonzero, peak):
            rms = math.sqrt(squares/self.frames) if self.frames else 0.0
            return {'nonzero_samples': nonzero, 'sum_squares': squares, 'rms_int16': rms,
                    'rms_dbfs': 20*math.log10(rms/32768) if rms else None, 'peak_abs_int16': peak,
                    'peak_dbfs': 20*math.log10(peak/32768) if peak else None}
        return {'frames': self.frames, 'duration_seconds': self.frames/RATE,
                'nonzero_frames': self.nonzero_frames, 'nonzero_fraction': self.nonzero_frames/self.frames if self.frames else 0,
                'left': channel(self.left_squares, self.left_nonzero, self.left_peak),
                'right': channel(self.right_squares, self.right_nonzero, self.right_peak),
                'different_channel_frames': self.different_frames, 'max_channel_delta': self.max_channel_delta}


def operations(directory, raw_size, audit, tag):
    rows = [json.loads(line) for line in audit.read(directory/'operations.jsonl').splitlines()]
    allowed = {'kind', 'monotonic_ns', 'device_bytes', 'x', 'y', 'step', 'error', 'number', 'label', 'code'}
    if not rows or any(set(r)-allowed or r.get('kind') not in KINDS for r in rows):
        raise ValueError('operations_schema')
    audit.check(all(isinstance(r.get('monotonic_ns'), int) and isinstance(r.get('device_bytes'), int) for r in rows), tag+'操作數字欄')
    times = [r['monotonic_ns'] for r in rows]
    offsets = [r['device_bytes'] for r in rows]
    audit.check(times == sorted(times) and offsets == sorted(offsets), tag+'操作時刻及裝置收集位移不倒退')
    audit.check(all(0 <= n <= raw_size and n % 4 == 0 for n in offsets), tag+'所有收集邊界為合法完整雙聲道幀')
    steps = [r['step'] for r in rows if 'step' in r]
    audit.check(steps == sorted(steps), tag+'原版Step不倒退')
    pauses, active_pause, op_pairs, started = [], None, [], {}
    for index, row in enumerate(rows):
        kind = row['kind']
        if kind == 'click':
            audit.check(row.get('error') == '', tag+'正常原版操作無錯誤')
        if kind == 'pause':
            audit.check(active_pause is None, tag+'F1暫停不巢狀')
            active_pause = (index, row)
        if kind == 'resume':
            audit.check(active_pause is not None, tag+'F1恢復有對應暫停')
            if active_pause:
                start_index, start = active_pause
                audit.check(start['step'] == row['step'], tag+'F1穩定暫停及恢復為同一Step')
                pauses.append({'start_marker': start_index, 'end_marker': index, 'step': row['step'],
                               'wall_seconds': (row['monotonic_ns']-start['monotonic_ns'])/1e9,
                               'captured_bytes': row['device_bytes']-start['device_bytes']})
            active_pause = None
        if kind in ('operation_start', 'operation_end'):
            number = row['number']
            # 只保留既定操作識別標籤，不把其他文字傳入輸出。
            if not re.fullmatch(r'audio-[a-z0-9-]+', row['label']):
                raise ValueError('operation_label')
            if kind == 'operation_start':
                audit.check(number not in started, tag+'正常操作起點唯一')
                started[number] = (index, row)
            else:
                audit.check(number in started, tag+'正常操作終點有起點')
                if number in started:
                    start_index, start = started.pop(number)
                    audit.check(start['label'] == row['label'], tag+'正常操作首尾識別相同')
                    op_pairs.append({'number': number, 'label': row['label'], 'start_marker': start_index,
                                     'end_marker': index, 'start_byte': start['device_bytes'], 'end_byte': row['device_bytes'],
                                     'wall_seconds': (row['monotonic_ns']-start['monotonic_ns'])/1e9})
    audit.check(active_pause is None and not started, tag+'所有F1及正常操作完整收尾')
    audit.check(len(pauses) > 0, tag+'確有正常F1暫停/恢復')
    close = [r for r in rows if r['kind'] == 'close_request']
    exit_rows = [r for r in rows if r['kind'] == 'process_exit']
    audit.check(len(close) == len(exit_rows) == 1, tag+'單一正常關窗/程序退出')
    if len(close) != 1 or len(exit_rows) != 1:
        raise ValueError('exit_markers')
    latency_ns = exit_rows[0]['monotonic_ns']-close[0]['monotonic_ns']
    audit.check(exit_rows[0]['code'] == 0 and 0 <= latency_ns <= 5000000000, tag+'正常關窗後5秒內程序退出0')
    audit.check(rows[-1] is exit_rows[0], tag+'程序退出為最後操作標記')
    return rows, pauses, op_pairs, latency_ns


def analyse_pcm(path, rows, audit, tag):
    original = path.stat()
    audit.check(original.st_size > 0 and original.st_size % 4 == 0, tag+'原始裝置檔完整s16LE雙聲道幀')
    boundaries = sorted({0, original.st_size, *(r['device_bytes'] for r in rows)})
    statistics = {}
    hasher = hashlib.sha256()
    total = PCMStats()
    with path.open('rb') as stream:
        for start, end in zip(boundaries, boundaries[1:]):
            current = PCMStats()
            remaining = end-start
            while remaining:
                data = stream.read(min(1048576, remaining))
                if not data:
                    raise ValueError('pcm_short_read')
                hasher.update(data)
                current.add(data)
                remaining -= len(data)
            statistics[start, end] = current
            total.merge(current)
        if stream.read(1):
            raise ValueError('pcm_grew_during_read')
    after = path.stat()
    audit.check((original.st_size, original.st_mtime_ns) == (after.st_size, after.st_mtime_ns), tag+'原始裝置檔驗收期間保持')
    audit.identities[str(path)] = hasher.hexdigest()
    audit.check(total.frames*4 == original.st_size and total.nonzero_frames > 0, tag+'完整原始裝置幀數及軟體非靜音')
    audit.check(total.different_frames == 0 and total.max_channel_delta == 0, tag+'全原始裝置資料左右相同')

    def interval(start, end):
        result = PCMStats()
        for (low, high), part in statistics.items():
            if start <= low and high <= end:
                result.merge(part)
        return result.report()

    segments = []
    for index in range(len(rows)+1):
        low = rows[index-1]['device_bytes'] if index else 0
        high = rows[index]['device_bytes'] if index < len(rows) else original.st_size
        segments.append({'start_marker': index-1 if index else None,
                         'end_marker': index if index < len(rows) else None,
                         'start_kind': rows[index-1]['kind'] if index else 'capture_start',
                         'end_kind': rows[index]['kind'] if index < len(rows) else 'capture_end',
                         'start_byte': low, 'end_byte': high, 'statistics': interval(low, high)})
    audit.check(sum(s['statistics']['frames'] for s in segments) == total.frames, tag+'收集塊分段完整不重疊覆蓋全部PCM')
    return {'bytes': original.st_size, 'sha256': hasher.hexdigest(), 'statistics': total.report(), 'segments': segments}, interval


def run(directory, audit, number):
    tag = f'第{number}趟/'
    summary = safe_summary(directory, audit)
    s, r, source = (summary[k] for k in ('stream', 'reader', 'source'))
    # A7允許停止合成後原版完成當前操作；CPU final不是已合成水位。
    t = s['EndStep']
    expected = (t*RATE*H_DENOMINATOR+H_NUMERATOR-1)//H_NUMERATOR
    audit.check(t > 0 and s['SourceStep'] == t and summary['original_end_step'] >= t,
                tag+'來源/合成水位相同且原版最終Step不早於停止合成')
    audit.check(s['Frames'] == s['CanonicalFrames'] == expected, tag+'獨立整數ceil的正常速度幀數')
    audit.check(s['Received'] == s['Applied']+s['Pending'], tag+'合法接受事件等於施加加待施加')
    audit.check(source['Recorded'] == s['Received']+s['Rejected']+summary['stop_unsent']+source['Pending'], tag+'來源事件到消費/停止未送守恆')
    audit.check(s['Pending'] == s['Rejected'] == s['Bank1'] == s['CSM'] == 0 and s['Unsupported'] == {}, tag+'本路線無待施加/拒絕/未支援功能')
    audit.check(source['Stream'] is True and source['Pending'] == source['Dropped'] == summary['stop_unsent'] == 0, tag+'來源stream完整交付且無容量掉事件')
    audit.check(s['Stopped'] is True and s['Incomplete'] is False and source['FaultCode'] == s['FaultCode'] == summary['fault'] == summary['close_fault'] == '', tag+'正常音訊停止且無有限故障')
    audit.check(r['Queued'] == 0 and 0 < r['MaximumQueued'] <= MAX_RING, tag+'有界ring並已收尾清空')
    audit.check(r['SourceFrames']+r['Overflow']+r['StopQueued'] == s['Frames'], tag+'來源PCM幀到消費/溢出/停止佇列守恆')
    audit.check(PREFILL+r['SourceFrames']+r['Underflow']+r['StopFade'] == r['DeviceFrames'], tag+'預填/來源/欠載/淡出到裝置幀守恆')
    audit.check(r['DeviceFrames']*4 == r['DeliveredBytes'], tag+'Reader完整立體聲byte守恆')
    audit.check(r['StopQueued'] == r['StopFade']+r['StopDiscard'] and r['StopFade'] == min(r['StopQueued'], MAX_FADE), tag+'停止佇列淡出及尾段捨棄守恆')
    audit.check(r['FadeDeliveredBytes'] == r['StopFade']*4, tag+'淡出完整立體聲byte守恆')
    audit.check(r['Stopped'] is True and r['ACK'] is True and r['EOF'] is True, tag+'Reader停止/完整幀ACK/EOF')
    audit.check(summary['fps'] > 0, tag+'正常畫格率有數字')
    raw_path = directory/'device.s16le'
    rows, pauses, op_pairs, latency = operations(directory, raw_path.stat().st_size, audit, tag)
    expected_operations = (['audio-clearance-typed', 'audio-load-dialog', 'audio-world-loaded',
                            'audio-turn-right', 'audio-forward', 'audio-disk-menu', 'audio-save-selected',
                            'audio-save-confirmed', 'audio-after-save-turn'] if number == 1 else
                           ['audio-clearance-typed', 'audio-load-dialog', 'audio-world-restored', 'audio-coldload-running'])
    audit.check([op['label'] for op in op_pairs] == expected_operations,
                tag+'正常冷LOAD/行走/保存或新程序冷讀檔操作鏈完整')
    pcm, interval = analyse_pcm(raw_path, rows, audit, tag)
    for op in op_pairs:
        op['captured_statistics'] = interval(op['start_byte'], op['end_byte'])
    return {'directory': str(directory), 'safe_summary': summary, 'expected_canonical_frames': expected,
            'original_steps_after_synthesis_stop': summary['original_end_step']-t,
            'raw_pcm': pcm, 'pause_resume_pairs': pauses, 'normal_operations': op_pairs,
            'close_to_process_exit_ns': latency, 'close_to_process_exit_ms': latency/1e6,
            'last_operation_capture_byte': rows[-1]['device_bytes'],
            'unassigned_capture_tail_bytes': pcm['bytes']-rows[-1]['device_bytes'],
            'limits': ['raw是軟體裝置捕捉，含F1暫停及收集期間靜音，與Reader DeliveredBytes分開。',
                       'operation區間沿收集塊位移，含佇列/裝置延遲，不是原操作精確聲音邊界。',
                       'Overflow及StopDiscard如實保留；不調參歸零，不以非靜音當作人耳品質確認。',
                       '只量這兩條正常路線，不外推未支援OPL功能、所有退出排程或其他平台。']}


def viewport(path, audit):
    from PIL import Image
    data = json.loads(audit.read(path))
    box = tuple(data['viewport'])
    audit.check(box == (216, 99, 744, 435), '正常保存/冷讀檔固定world viewport矩形')
    pixels = (box[2]-box[0])*(box[3]-box[1])
    images = {}
    hashes = {}
    for role in ('saved', 'after_turn', 'cold_loaded'):
        name = Path(data['files'][role]['path'])
        if name.is_absolute() or '..' in name.parts or name.suffix != '.png':
            raise ValueError('viewport_png_path')
        image_path = path.parent/name
        raw = audit.read(image_path)
        hashes[role] = digest(raw)
        audit.check(hashes[role] == data['files'][role]['sha256'], '正常viewport來源PNG固定雜湊/'+role)
        with Image.open(image_path) as image:
            audit.check(image.width >= box[2] and image.height >= box[3], '正常viewport在實際PNG內/'+role)
            images[role] = tuple(image.convert('RGBA').crop(box).getdata())
    restored = sum(a != b for a, b in zip(images['saved'], images['cold_loaded']))
    unsaved = sum(a != b for a, b in zip(images['saved'], images['after_turn']))
    audit.check(pixels == len(images['saved']) == data['pixels'] == 177408, '正常viewport完整177408像素')
    audit.check(restored == data['save_to_coldload_differences'] == 0, '正常保存與新程序冷讀檔viewport完全相同')
    audit.check(unsaved == data['save_to_unsaved_turn_differences'] == 162504, '保存後未保存轉向為獨立正反對照')
    return {'rectangle': box, 'rgba_pixels': pixels, 'saved_to_coldload_differences': restored,
            'saved_to_unsaved_turn_differences': unsaved, 'source_sha256': hashes,
            'limits': '只證明正常world viewport恢復；不讀保存內容、不稱RAM或全部原版狀態同狀態對拍。'}


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--first', type=Path, required=True)
    p.add_argument('--second', type=Path)
    p.add_argument('--spec', type=Path, required=True)
    p.add_argument('--reader-source', type=Path, required=True)
    p.add_argument('--stream-source', type=Path, required=True)
    p.add_argument('--viewport-receipt', type=Path)
    p.add_argument('--out', type=Path, required=True)
    a = p.parse_args()
    assert a.out.name.startswith('audio-check') and a.out.suffix == '.json'
    assert not a.out.exists() and not a.out.with_suffix('.md').exists()
    audit = Audit()
    audit.read(Path(__file__))
    for path in (a.spec, a.reader_source, a.stream_source):
        audit.read(path)
    runs = []
    view = None
    error = None
    try:
        for number, directory in enumerate([a.first]+([a.second] if a.second else []), 1):
            runs.append(run(directory, audit, number))
        if a.viewport_receipt:
            view = viewport(a.viewport_receipt, audit)
    except (OSError, ValueError, KeyError, TypeError, IndexError, json.JSONDecodeError) as exc:
        # 不把任意來源文字或錯誤訊息抄入報告。
        error = type(exc).__name__
        audit.check(False, '驗收輸入/格式受阻：'+error)
    report = {'status': 'FAIL' if audit.failures else 'PASS', 'scope': '兩趟' if a.second else '僅第一趟',
              'checks': audit.checks, 'failures': audit.failures, 'input_error_type': error,
              'sample_rate': RATE, 'format': 's16LE stereo', 'h_numerator': H_NUMERATOR,
              'h_denominator': H_DENOMINATOR, 'formal_binary_sha256_provided_by_parent': FORMAL_BINARY,
              'binary_verified_by_this_tool': False, 'runs': runs, 'viewport': view, 'source_sha256': audit.identities,
              'limits': ['此工具不讀原版SAVE內容；可選viewport只核列出畫面恢復，不作RAM同狀態聲明。',
                         '此工具不驗啟動命令，正常無mute/ips/privatehook與正式binary身份沿主代理封存來源。',
                         '軟體非靜音不等於人耳或原生平台試玩；人耳回覆仍待使用者。',
                         'YM3812合成沿既定hardware-spec approximation，不宣稱原版逐週期/逐波形/實體喇叭保真。',
                         'Fade ACK有完成旗標，沒有單獨時刻，不稱已量獨立fade延遲。']}
    a.out.write_text(json.dumps(report, ensure_ascii=False, indent=2, allow_nan=False)+'\n')
    lines = [f"# R1001 獨立音訊驗收：{report['status']}", '',
             f"範圍：{report['scope']}；{audit.checks}項檢查，失敗{len(audit.failures)}項。", '',
             '|趟|原版Step|canonical幀|raw幀|raw非零幀|raw左右差幀|Overflow|StopDiscard|退出ms|',
             '|---|---:|---:|---:|---:|---:|---:|---:|---:|']
    for index, one in enumerate(runs, 1):
        stats = one['raw_pcm']['statistics']
        reader = one['safe_summary']['reader']
        lines.append(f"|{index}|{one['safe_summary']['original_end_step']}|{one['expected_canonical_frames']}|{stats['frames']}|{stats['nonzero_frames']}|{stats['different_channel_frames']}|{reader['Overflow']}|{reader['StopDiscard']}|{one['close_to_process_exit_ms']:.6f}|")
    lines += ['', '失敗：', '']
    lines += ['- '+s for s in audit.failures] if audit.failures else ['無。']
    lines += ['', '限制：', '']+['- '+s for s in report['limits']]
    a.out.with_suffix('.md').write_text('\n'.join(lines)+'\n')
    print(json.dumps({'status': report['status'], 'scope': report['scope'], 'checks': audit.checks,
                      'failures': audit.failures, 'runs': [{'raw': r['raw_pcm']['statistics'],
                      'overflow': r['safe_summary']['reader']['Overflow'], 'close_ms': r['close_to_process_exit_ms']} for r in runs]}, ensure_ascii=False, allow_nan=False))
    return int(bool(audit.failures))


if __name__ == '__main__':
    raise SystemExit(main())
