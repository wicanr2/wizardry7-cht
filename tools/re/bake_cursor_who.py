"""將正式WHO字模增補既有xlate字典，不改其他鍵或資源schema。"""
import argparse
import csv
import hashlib
import json
from pathlib import Path

def bake(xlate_path, flow_path, catalog_path, output_path):
    xlate_path,flow_path,catalog_path,output_path=map(Path,[xlate_path,flow_path,catalog_path,output_path])
    if output_path.resolve() in [xlate_path.resolve(),flow_path.resolve(),catalog_path.resolve()]:
        raise ValueError('output must differ from inputs')
    with catalog_path.open(encoding='utf-8',newline='') as f:
        reader=csv.DictReader(f,delimiter='\t')
        if reader.fieldnames!=['en','zh','context','dy','source']:
            raise ValueError('WHO catalog columns')
        rows=list(reader)
    if len(rows)!=1 or rows[0]['en']!='WHO?' or rows[0]['zh']!='誰？' or rows[0]['context']!='cursor.party.member.selection' or rows[0]['dy']!='0':
        raise ValueError('WHO catalog identity')
    data=json.loads(xlate_path.read_text(encoding='utf-8'))
    if not isinstance(data,dict) or 'WHO?' in data:
        raise ValueError('WHO key already exists or bad xlate')
    font=json.loads(flow_path.read_text(encoding='utf-8'))
    size=font.get('size')
    if type(size)!=int or size not in [15,16]:
        raise ValueError('WHO font height')
    glyphs=[]
    for c in rows[0]['zh']:
        g=font.get('glyphs',{}).get(c)
        if not isinstance(g,list) or len(g)!=size or not any(g):
            raise ValueError('WHO glyph missing')
        if any(type(v)!=int or v<0 or v>=(1<<size) for v in g):
            raise ValueError('WHO glyph bit bounds')
        glyphs.append(g)
    data['WHO?']={'dy':0,'g':glyphs}
    output_path.write_text(json.dumps(data,ensure_ascii=False,separators=(',',':')),encoding='utf-8')
    return {'english_key':'WHO?','height':size,'added_keys':1,
            'catalog_sha256':hashlib.sha256(catalog_path.read_bytes()).hexdigest(),
            'xlate_input_sha256':hashlib.sha256(xlate_path.read_bytes()).hexdigest(),
            'flow_input_sha256':hashlib.sha256(flow_path.read_bytes()).hexdigest(),
            'output_sha256':hashlib.sha256(output_path.read_bytes()).hexdigest()}

if __name__=='__main__':
    p=argparse.ArgumentParser()
    for name in ['xlate','flow','catalog','out']:p.add_argument('--'+name,required=True)
    a=p.parse_args()
    print(json.dumps(bake(a.xlate,a.flow,a.catalog,a.out)))
