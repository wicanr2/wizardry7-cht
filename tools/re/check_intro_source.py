#!/usr/bin/env python3
"""片頭正常冷啟動的有／無觀察來源復核；只在Docker執行，不驗中文覆繪。"""
import argparse
import hashlib
import json
import re
from pathlib import Path

from check_native_gui import png


PAPER_STEPS = (73235759, 75322912, 77352859, 100999086, 103102872, 105202085)
KNOWN_FRAMES = {
    38938757: 'ef8a35221a7405572804a277e965a7fdcce7768c7f646866b271b45af35737a1',
    40696847: '3858a8793984ff7a7754d11308d441b06e4e3695a2f651cca28f7e7ba973fc2f',
    42466770: 'dbc1d7aa1471aaad52e87bf8c644c9d6d28eb888aad3932d2da2b190f652d091',
    44378262: '7d1ecb00f365bea1298518803ebadf6ea9ab04f788dcaea151a37885844e8301',
    46331493: '2996bda97e5591818aecdee6fe856fd1629c6ba6882983f7a7b8173771fc3e66',
    48267654: '69220ee76660b723ca0c52d7c54ff86b0735192963ebac5c4962176081aa28de',
    73235759: '0f098eb7bf96d22777424ac0f8bc263534be1441ecc631413d61a2323714411a',
}
PICTURE_PALETTE = 'd689ae136905afa64a9329ddd123426f7e0f3cde6108c034e0249dbab0a9727c'
FIELDS = ('Step', 'Mem', 'A000', 'Pal', 'MemorySHA256', 'IndexedSHA256',
          'PaletteSHA256', 'Mouse', 'Mode')


def sha(data):
    return hashlib.sha256(data).hexdigest()


def require(condition, message):
    if not condition:
        raise ValueError(message)


def fnv(data):
    value = 2166136261
    for byte in data:
        value = ((value ^ byte) * 16777619) & 0xffffffff
    return value


def record_shape(row):
    require(isinstance(row, dict) and all(key in row for key in FIELDS), '快照缺必要欄位')
    for key in ('Step', 'Mem', 'A000', 'Pal', 'Mode'):
        require(type(row[key]) is int, f'{key}須為整數')
    require(0 < row['Step'] <= 200000000, '步數越界')
    require(all(0 <= row[key] <= 0xffffffff for key in ('Mem', 'A000', 'Pal')), '摘要越界')
    for key in ('MemorySHA256', 'IndexedSHA256', 'PaletteSHA256'):
        require(isinstance(row[key], str) and re.fullmatch('[0-9a-f]{64}', row[key]),
                f'{key}須為完整SHA-256')
    require(isinstance(row['Mouse'], list) and len(row['Mouse']) == 2
            and all(type(n) is int for n in row['Mouse']), '滑鼠欄位不合法')


def verify_files(folder, row, prefix):
    record_shape(row)
    base = folder / f'{prefix}-{row["Step"]}'
    indexed = base.with_suffix('.a000').read_bytes()
    palette = base.with_suffix('.pal').read_bytes()
    require(len(indexed) == 64000 and len(palette) == 768, '畫面或色盤長度不符')
    require(sha(indexed) == row['IndexedSHA256'] and fnv(indexed) == row['A000'], 'indexed摘要不符')
    require(sha(palette) == row['PaletteSHA256'] and fnv(palette) == row['Pal'], '色盤摘要不符')
    width, height, rgb = png(base.with_suffix('.png'))
    require((width, height) == (320, 200), 'PNG尺寸不符')
    colors = [palette[i:i+3] for i in range(0, 768, 3)]
    require(rgb == b''.join(colors[n] for n in indexed), 'PNG與原版色盤解幀不符')
    return indexed, palette


def check(observed, control):
    reports = [json.loads((folder / 'observation.json').read_text()) for folder in (observed, control)]
    require([r['Observer'] for r in reports] == ['on', 'off'], '有／無觀察標記不符')
    rows = [r['Snapshots'] for r in reports]
    expected_steps = list(range(2000000, 200000001, 2000000))
    require(all([r['Step'] for r in group] == expected_steps for group in rows), '固定100點步數不符')
    for a, b in zip(*rows):
        left, right = [verify_files(folder, row, 'original')
                       for folder, row in zip((observed, control), (a, b))]
        require(all(a[key] == b[key] for key in FIELDS), f'{a["Step"]}原版摘要不同')
        require(left == right, f'{a["Step"]}原版indexed或色盤bytes不同')
    require(all(not reports[1][key] for key in ('Renders', 'PictureCaptures', 'PicOpens', 'Opened', 'Blits')),
            '無觀察控制含觀察資料')
    pictures = reports[0]['PictureCaptures']
    require(len(pictures) == 14 and len({r['Step'] for r in pictures}) == 14, '完整轉送畫格數不符')
    require([r['Step'] for r in pictures] == sorted(r['Step'] for r in pictures), '完整轉送畫格未排序')
    transfers = [row for row in reports[0]['Renders'] if row['IP'] == 0x3921
                 and row['Words'][0] in (0x5b9b, 0x5afd)]
    require(len(transfers) == len(pictures), '完整轉送入口與返回畫格數不符')
    for index, (entry, picture) in enumerate(zip(transfers, pictures)):
        require(entry['Words'][1:] == [0, 0, 319, 199], '完整轉送矩形不符')
        require(type(entry['Step']) is int and entry['Step'] < picture['Step'], '返回畫格順序不符')
        if index + 1 < len(transfers):
            require(picture['Step'] < transfers[index + 1]['Step'], '完整轉送見證互相交錯')
    by_step = {r['Step']: r for r in pictures}
    paper = []
    for row in pictures:
        verify_files(observed, row, 'picture')
        require(row['Mode'] == 0x13, '圖片模式不符')
    for step, expected in KNOWN_FRAMES.items():
        require(step in by_step and by_step[step]['IndexedSHA256'] == expected, f'{step}既有構圖SHA不同')
        require(by_step[step]['PaletteSHA256'] == PICTURE_PALETTE, f'{step}已量構圖色盤不符')
    for step in PAPER_STEPS:
        require(step in by_step, f'{step}紙頁來源缺失')
        row = by_step[step]
        require(row['PaletteSHA256'] == PICTURE_PALETTE, f'{step}紙頁色盤不符')
        paper.append({key: row[key] for key in ('Step', 'IndexedSHA256', 'PaletteSHA256')})
    return {
        'schema': 1, 'fixed_pairs': 100, 'decoded_pngs': 214,
        'memory_sha256_equal': True, 'indexed_palette_bytes_equal': True,
        'complete_picture_frames': 14, 'full_transfer_entries_verified': len(transfers),
        'known_picture_hashes_equal': len(KNOWN_FRAMES),
        'paper': paper,
        'receipt_sha256': {name: sha((folder / 'observation.json').read_bytes())
                           for name, folder in zip(('observed', 'control'), (observed, control))},
        'limits': '正常無輸入原版片頭；RAM只核1MiB SHA-256摘要，未保存完整RAM或Regs；'
                  '只恢復原始來源，不恢復R944選定排版模型，不驗中文、正式GUI或完整遊戲。',
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('observed', type=Path)
    parser.add_argument('control', type=Path)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    if args.output.exists():
        parser.error('輸出已存在，拒絕覆寫')
    try:
        result = check(args.observed, args.control)
        with args.output.open('x') as handle:
            json.dump(result, handle, ensure_ascii=False, sort_keys=True, indent=2)
            handle.write('\n')
    except (ValueError, KeyError, TypeError, OSError) as error:
        parser.error(str(error))
    print('100對原版摘要及indexed／色盤bytes相同，214張PNG解幀相符；紙頁6來源可重生。')


if __name__ == '__main__':
    main()
