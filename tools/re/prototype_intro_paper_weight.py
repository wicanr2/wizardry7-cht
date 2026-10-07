#!/usr/bin/env python3
"""R975：紙頁文字放大、加粗及提亮的可撤回3倍畫布原型。"""
import argparse
import hashlib
import html
import json
import sys
from pathlib import Path

import prototype_intro_paper as baseline

MODEL_SHA = '1a9632b661b106367adffc2a8336808fe38202e60080620e897072e2fdcf5a22'
TOOL_SHA = '216c0d4a497593382abdda32cd6d9831ee4e2d55f7a37d3e0c0ddd3b7c822660'
ETEN_SHA = '39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701'
COMPARISON_SHA = {
    'unifont': 'e93d35efc9a3e4e3c3fd000d91090ea3b0697d6087c5ba9e803cbbeabda22ea9',
    'eten': 'aab6a4ed8f008c57a6f40951fa4c96b21df11ce36fd77cb8b91e598871a1e37b',
}
RED = (196, 36, 20)
DARK = (40, 28, 16)
STYLES = (('bold', '甲：2倍字模＋水平加粗1像素', 1),
          ('plain', '乙：2倍字模', 0))


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def layout(row):
    labels = [dict(label) for label in row['labels']]
    if row['step'] == 73235759:
        labels[0]['center'] = [199,68]
    if row['step'] == 75322912:
        safe = labels[0]['safe']
        labels = [dict(en=en, zh=zh, center=center, color='red', safe=safe)
                  for en, zh, center in (
                      ('fantasy', '奇幻', [93,32]),
                      ('role-playing', '角色', [106,52]),
                      ('role-playing', '扮演', [126,73]),
                      ('simulation', '遊戲', [145,94]))] + labels[-1:]
    if row['step'] == 100999086:
        labels[-1]['center'] = [173.5,88]
    return labels


def ink(glyph, size, right):
    result = set()
    for y, bits in enumerate(glyph):
        for x in range(size):
            if bits & (1 << (size-1-x)):
                for dy in (0,1):
                    for dx in range(2+right):
                        result.add((x*2+dx,y*2+dy))
    return result


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('source', type=Path)
    p.add_argument('r974_candidate', type=Path)
    p.add_argument('output', type=Path)
    p.add_argument('--eten-font', type=Path, required=True)
    args = p.parse_args()
    if args.output.exists():
        p.error('輸出已存在，請指定新目錄')
    assert sha(baseline.MODEL.read_bytes()) == MODEL_SHA
    assert sha(Path(baseline.__file__).read_bytes()) == TOOL_SHA
    assert sha(args.eten_font.read_bytes()) == ETEN_SHA
    model = json.loads(baseline.MODEL.read_text())
    import bake_font
    import font_source
    from check_native_gui import png
    bake_font.FONT15 = args.eten_font
    backend, size, k = bake_font.BACKEND, bake_font.HEIGHT, 3
    comparison_raw = (args.r974_candidate/'receipt.json').read_bytes()
    assert sha(comparison_raw) == COMPARISON_SHA[backend], 'R974比較來源收據不符'
    comparison = json.loads(comparison_raw)
    chars = ''.join(dict.fromkeys(c for f in model['frames'] for l in f['labels'] for c in l['zh']))
    fonts = {c:bake_font.glyph_char(c) for c in chars}
    assert all(len(g)==size and any(g) and all(0<=b<1<<size for b in g) for g in fonts.values())
    generated, receipts, sections = [], [], []
    for row in model['frames']:
        if not row['labels']:
            continue
        path = args.source/f'picture-{row["step"]}'
        indexed, pal = path.with_suffix('.a000').read_bytes(), path.with_suffix('.pal').read_bytes()
        assert len(indexed)==64000 and len(pal)==768
        assert sha(indexed)==row['indexed_sha256'] and sha(pal)==model['palette_sha256']
        rgb = b''.join(pal[v*3:v*3+3] for v in indexed)
        protected = set().union(*(baseline.points(q) for q in row['protect'])) if row['protect'] else set()
        patched, _ = baseline.patch(rgb,row['bands'],protected,row['band_surfaces'])
        base = baseline.scaled(patched,k)
        original = baseline.scaled(rgb,k)
        old_path = args.r974_candidate/f'{backend}-3x-{row["step"]}-candidate.png'
        reference = next(r for r in comparison['frames'] if r['step']==row['step'] and r['scale']==3)
        assert old_path.name == reference['candidate']
        assert sha(old_path.read_bytes()) == reference['candidate_sha256'], 'R974比較PNG不符'
        old_w,old_h,old = png(old_path)
        assert (old_w,old_h)==(960,600)
        names = dict(original=f'{backend}-{row["step"]}-original.png',
                     old=f'{backend}-{row["step"]}-r974.png')
        generated.extend(((names['original'],original),(names['old'],old)))
        figures = [('原版',names['original']),('R974',names['old'])]
        for style,title,right in STYLES:
            image = bytearray(base)
            labels = []
            cells = set()
            for label in layout(row):
                width = len(label['zh'])*(2*size+4)-4+right
                height = size*2
                x0 = round(label['center'][0]*k)-width//2
                y0 = round(label['center'][1]*k)-height//2
                label_cells = {(x0+x,y0+y) for y in range(height) for x in range(width)}
                assert not cells & label_cells, (row['step'],style,'字格重疊')
                cells.update(label_cells)
                assert all(baseline.inside((x+.5)/k,(y+.5)/k,label['safe']) for x,y in label_cells), (row['step'],style,label['zh'],'字格出紙面')
                assert all((x//k,y//k) not in protected for x,y in label_cells), (row['step'],style,'作者區重疊')
                color = RED if label['color']=='red' else DARK
                on = 0
                for j,c in enumerate(label['zh']):
                    pixels = ink(fonts[c],size,right)
                    on += len(pixels)
                    for gx,gy in pixels:
                        off = ((y0+gy)*960+x0+j*(2*size+4)+gx)*3
                        image[off:off+3] = bytes(color)
                labels.append(dict(**label,xy=[x0,y0],width=width,height=height,rgb=list(color),ink=on))
            name = f'{backend}-{row["step"]}-{style}.png'
            generated.append((name,image))
            figures.append((title,name))
            receipts.append(dict(backend=backend,scale=k,step=row['step'],style=style,
                                 magnify=2,embolden_right=right,gap=4,output=name,labels=labels,
                                 indexed_sha256=row['indexed_sha256'],protect=row['protect'],
                                 comparison_png_sha256=reference['candidate_sha256']))
        sections.append('<section><h2>'+str(row['step'])+'</h2><div>'+''.join(
            f'<figure><figcaption>{html.escape(title)}</figcaption><img src="{name}"></figure>' for title,name in figures)+'</div></section>')
    args.output.mkdir()
    for name,data in generated:
        baseline.png_write(args.output/name,960,600,data)
    for row in receipts:
        row['output_sha256'] = sha((args.output/row['output']).read_bytes())
    value = dict(schema=1,status='DRAFT typography prototype, not production',backend=backend,size=size,
                 baseline_model_sha256=MODEL_SHA,baseline_tool_sha256=TOOL_SHA,
                 tool_sha256=sha(Path(__file__).read_bytes()),
                 python_version=sys.version.split()[0],
                 comparison_receipt_sha256=COMPARISON_SHA[backend],
                 font_sha256=ETEN_SHA if backend=='eten' else font_source.UNIFONT_SHA256,
                 colors=dict(red=list(RED),dark=list(DARK)),glyphs=fonts,frames=receipts,
                 limits='僅3倍紙頁外觀。整數放大及加粗是新顯示效果，非原生1:1；背景固定R974且仍有條紋，正式規格與2倍畫布未放行。')
    with (args.output/'receipt.json').open('x') as f:
        json.dump(value,f,ensure_ascii=False,indent=2)
        f.write('\n')
    page = ('<!doctype html><html lang="zh-Hant"><meta charset="utf-8"><title>R975紙頁字樣對照</title>'
            '<style>body{background:#181818;color:#eee;font:18px sans-serif;margin:24px}section>div{display:flex;flex-wrap:wrap;gap:20px}figure{margin:0}img{image-rendering:pixelated}figcaption{margin:8px 0}</style>'
            '<h1>R975紙頁字樣對照</h1><p>字模放大2倍為30／32像素，兩案皆稍提亮。甲另水平加粗1像素；乙保持放大後筆畫。請以原尺寸查看。</p>'
            '<p>紙紋固定R974供比較，補底條紋另待改善。這是3倍畫布的新原型，未接正式前端。</p>'+''.join(sections)+'</html>')
    (args.output/'compare.html').write_text(page)
    print(f'{backend}：五構圖、兩種厚度、共10張3倍字樣候選')


if __name__ == '__main__':
    main()
