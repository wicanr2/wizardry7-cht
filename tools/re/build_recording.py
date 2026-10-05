#!/usr/bin/env python3
"""於專案 Docker 內建置包含錄放來源身份的 wiz7。"""
import argparse
import base64
import hashlib
import json
import os
from pathlib import Path
import subprocess


def run(args, cwd):
    return subprocess.run(args, cwd=cwd, check=True, capture_output=True,
                          text=True, timeout=120).stdout


def digest(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', required=True)
    parser.add_argument('--receipt-dir', required=True)
    args = parser.parse_args()
    root = Path(__file__).resolve().parents[2]
    output, receipts = Path(args.output).absolute(), Path(args.receipt_dir).absolute()
    for directory in [root, output.parent, receipts]:
        if not directory.is_dir() or directory.stat().st_uid != os.getuid():
            raise RuntimeError('build_output_directory')
    if output.exists() or output.is_symlink():
        raise RuntimeError('build_output_exists')
    source_receipt = receipts / 'build-source-inputs.json'
    result_receipt = receipts / 'build-result.json'
    for receipt_path in [source_receipt, result_receipt]:
        if receipt_path.exists() or receipt_path.is_symlink():
            raise RuntimeError('build_receipt_exists')
    if run(['go', 'env', 'GOVERSION'], root).strip() != 'go1.24.13':
        raise RuntimeError('build_go_version')
    stream = run(['go', 'list', '-deps', '-json', './cmd/wiz7'], root)
    decoder = json.JSONDecoder()
    packages = []
    position = 0
    while position < len(stream):
        if stream[position].isspace():
            position += 1
            continue
        package, position = decoder.raw_decode(stream, position)
        packages.append(package)
    paths = {}
    for package in packages:
        module = package.get('Module', {})
        module_name = module.get('Path')
        if module_name not in ['github.com/wicanr2/wizardry7-cht', 'github.com/wicanr2/dosgolem']:
            continue
        module_dir = Path(module['Dir'])
        package_dir = Path(package['Dir'])
        prefix = '' if module_name.endswith('wizardry7-cht') else 'sdk/'
        for name in package.get('GoFiles', []) + package.get('CgoFiles', []):
            actual = package_dir / name
            relative = prefix + actual.relative_to(module_dir).as_posix()
            paths[relative] = actual
        for name in ['go.mod', 'go.sum']:
            actual = module_dir / name
            if actual.exists():
                paths[prefix + name] = actual
    table = root / 'cmd/wiz7/recording_source_paths.go'
    if paths.get('cmd/wiz7/recording_source_paths.go') != table:
        raise RuntimeError('build_source_table_missing')
    table_text = ('package main\n\n'
                  '// 自動由 tools/re/build_recording.py 依實際選定來源重建；此表納入雜湊。\n'
                  'var recordingSourcePaths = []string{\n' +
                  ''.join('    ' + json.dumps(path, ensure_ascii=False) + ',\n'
                          for path in sorted(paths)) + '}\n')
    if table.stat().st_uid != os.getuid():
        raise RuntimeError('build_source_table_owner')
    formatted = subprocess.run(['gofmt'], input=table_text, cwd=root, check=True,
                               capture_output=True, text=True, timeout=20).stdout
    table.write_text(formatted)
    sources = []
    for relative, actual in sorted(paths.items()):
        if actual.is_symlink() or not actual.is_file():
            raise RuntimeError('build_source_type')
        data = actual.read_bytes()
        sources.append({'path': relative, 'size': len(data), 'sha256': digest(data)})
    if not 1 <= len(sources) <= 4096:
        raise RuntimeError('build_source_capacity')
    encoded = json.dumps(sources, ensure_ascii=False, separators=(',', ':')).encode()
    injection = base64.b64encode(encoded).decode('ascii')
    flags = '-X main.compiledRecordingSources=' + injection + ' -X main.compiledRecordingSourceSHA=' + digest(encoded)
    run(['go', 'build', '-p', '2', '-trimpath', '-buildvcs=false', '-ldflags', flags,
         '-o', str(output), './cmd/wiz7'], root)
    # 不把編譯期間改動過的輸入冒充同一次建置來源。
    for entry in sources:
        data = paths[entry['path']].read_bytes()
        if len(data) != entry['size'] or digest(data) != entry['sha256']:
            raise RuntimeError('build_source_changed')
    with source_receipt.open('xb') as receipt_file:
        receipt_file.write(encoded + b'\n')
    receipt = {'schema': 'wiz7-recording-build/1', 'go_version': 'go1.24.13',
               'goos': run(['go', 'env', 'GOOS'], root).strip(),
               'goarch': run(['go', 'env', 'GOARCH'], root).strip(),
               'source_count': len(sources), 'source_sha256': digest(encoded),
               'binary_sha256': digest(output.read_bytes()),
               'binary_size': output.stat().st_size,
               'sources': 'build-source-inputs.json'}
    with result_receipt.open('x') as receipt_file:
        receipt_file.write(json.dumps(receipt, indent=2) + '\n')
    print(json.dumps(receipt))


if __name__ == '__main__':
    main()
