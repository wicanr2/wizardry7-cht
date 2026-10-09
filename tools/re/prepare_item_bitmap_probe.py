#!/usr/bin/env python3
"""在Docker內準備R999 READY私人長劍觀察副本；本工具不執行原版。"""
import argparse
import os
from pathlib import Path
import shutil
import subprocess


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--source', required=True, type=Path)
    parser.add_argument('--output', required=True, type=Path)
    args = parser.parse_args()
    source, target = args.source.resolve(), args.output.absolute()
    if not source.is_dir() or target.exists() or target.is_symlink():
        raise RuntimeError('source_or_output')
    if not target.parent.is_dir() or target.parent.stat().st_uid != os.getuid():
        raise RuntimeError('output_parent_owner')
    target.mkdir()
    for name in ['cmd/wiz7', 'initconfig', 'recording']:
        shutil.copytree(source / name, target / name, symlinks=True)
    for name in ['go.mod', 'go.sum', 'tools/re/build_recording.py']:
        path = target / name
        path.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source / name, path)
    (target / 'worktrees').mkdir()
    (target / 'worktrees/dosgolem').symlink_to(source / 'worktrees/dosgolem', target_is_directory=True)
    path = target / 'cmd/wiz7/main.go'
    code = path.read_text()
    changes = {
        '\t\t\tY: int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})), raw: partRaw}\n':
        '\t\t\tY: int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})), raw: partRaw}\n'
        '\t\tr999ObserveItem(o, ret, r0, in.Font, off, in.X, in.Y, partRaw)\n',
        '\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)':
        '\tlive = r982Control(o, *root, preparedSaveDir)\n'
        '\tif live != nil {\n\t\tlog.SetFlags(0)\n\t\tlog.SetOutput(recordingPrivacyLog{})\n\t}\n'
        '\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)',
        '\tsnapshot := func() {\n':
        '\tsnapshot := func() {\n\t\tr982Observe(o, live)\n',
    }
    for old, new in changes.items():
        if code.count(old) != 1:
            raise RuntimeError('source_hook_mismatch')
        code = code.replace(old, new)
    path.write_text(code)
    here = Path(__file__).resolve().parent
    helpers = {
        'item_nul_probe.go': 'r999_item_probe.go',
        'item_bitmap_probe.go': 'r999_service24.go',
        'item_probe_api.go': 'r982_probe.go',
    }
    formatted = [path]
    for saved, name in helpers.items():
        dest = target / 'cmd/wiz7' / name
        shutil.copy2(here / saved, dest)
        formatted.append(dest)
    subprocess.run(['gofmt', '-w', *map(str, formatted)], check=True, timeout=30)
    print('private observer source prepared; original not executed')


if __name__ == '__main__':
    main()
