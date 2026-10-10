#!/usr/bin/env python3
"""R1002 P1–P5：只建立私人98來源及純建置，不執行原版或讀取wire。"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess


BASE_MAIN = '7d1b069fc1edc3706dcd8a4d2ae5b6bc250236b2693a1ead75ee836186dd8c18'
API_SHA = '9833d39ad7a12943d0106b9f1ac8d0540cabcf15c85a632667a1c5f2c9f211ff'
FACTS_SHA = 'c75123dd6885ea6832a73c11647ff92167f55259365a677582c949ef967add10'
REVIEW_SHA = '50a735cd4b1938005aa0e871f49fa267a0edd577180613076ce8a394f99a6a3a'
EXTENSION_REVIEW_SHA = '0951843e98a6555899ae7111bd3cdad5efeb609fec21f81a89dbfedba1c5bc9e'


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--src', type=Path, default=Path('/src'))
    p.add_argument('--out', type=Path, default=Path('/out'))
    a = p.parse_args()
    src, out = a.src, a.out
    assert src.is_dir() and out.is_dir() and out.stat().st_uid == os.getuid() == 1000
    main_path = src/'cmd/wiz7/main.go'
    api_path = src/'tools/re/item_probe_api.go'
    facts_path = src/'workplace/r1002/observer-contract-facts.json'
    review_path = src/'workplace/r1002/probe-contract-review.json'
    assert sha(main_path) == BASE_MAIN and sha(api_path) == API_SHA
    assert sha(facts_path) == FACTS_SHA and sha(review_path) == REVIEW_SHA
    facts = json.loads(facts_path.read_bytes())
    review = json.loads(review_path.read_bytes())
    extension_review_path = src/'workplace/r1002/probe-extension-review.json'
    assert sha(extension_review_path) == EXTENSION_REVIEW_SHA
    extension_review = json.loads(extension_review_path.read_bytes())
    assert extension_review['counts'] == {'blocking': 0, 'should_change': 0, 'suggestion': 0}
    assert review['counts'] == {'blocking': 0, 'should_change': 0, 'suggestion': 0}
    expected = {i: key for key, ids in facts['authority']['groups_from_existing_evidence'].items() for i in ids}
    assert len(expected) == 29 and set(expected) == set(facts['authority']['whitelist_ids'])
    assert not any(2500 <= i <= 2574 for i in expected)
    observer = out/'context_probe.go'
    assert observer.is_file() and observer.stat().st_uid == os.getuid()
    subprocess.run(['gofmt', '-w', str(observer)], check=True, timeout=20)
    text = observer.read_text()
    actual = {int(i): json.loads(value) for i, value in re.findall(r'(\d+):\s*("(?:[^"\\]|\\.)*")', text)}
    assert actual == expected
    tree, build = out/'probe-source-r1', out/'probe-build-r1'
    report_path = out/'implementation-report.json'
    assert not tree.exists() and not build.exists() and not report_path.exists()
    table = src/'cmd/wiz7/recording_source_paths.go'
    paths = re.findall(r'^\s*"([^"\n]+)",\s*$', table.read_text(), re.M)
    assert len(paths) == len(set(paths)) == 96
    baseline = []
    tree.mkdir()
    build.mkdir()
    for relative in paths:
        path = Path(relative)
        assert not path.is_absolute() and '..' not in path.parts
        original = src/'worktrees/dosgolem'/Path(*path.parts[1:]) if path.parts[0] == 'sdk' else src/path
        assert original.is_file() and not original.is_symlink()
        baseline.append({'path': relative, 'sha256': sha(original), 'size': original.stat().st_size})
        if path.parts[0] != 'sdk':
            target = tree/path
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(original, target)
    sdk = src/'worktrees/dosgolem'
    assert sdk.is_dir()
    (tree/'worktrees').mkdir()
    (tree/'worktrees/dosgolem').symlink_to(sdk, target_is_directory=True)
    tools = tree/'tools/re'
    tools.mkdir(parents=True)
    shutil.copy2(src/'tools/re/build_recording.py', tools/'build_recording.py')
    shutil.copy2(api_path, tree/'cmd/wiz7/r982_probe.go')
    shutil.copy2(observer, tree/'cmd/wiz7/r1002_context_probe.go')
    main_target = tree/'cmd/wiz7/main.go'
    main_text = main_target.read_text()
    anchor = '\tlog.Printf("原版存檔目錄：%s", preparedSaveDir)\n'
    setup = (anchor + '\tif os.Getenv("WIZ7_R982_API_DOCUMENT") != "" {\n'
             '\t\tlive = r982Control(o, *root, preparedSaveDir)\n'
             '\t\tlog.SetOutput(recordingPrivacyLog{})\n\t}\n'
             '\tcontextObserver := r1002Setup(o)\n\tdefer contextObserver.Close()\n')
    snapshot = '\tsnapshot := func() {\n'
    assert main_text.count(anchor) == main_text.count(snapshot) == 1
    main_text = main_text.replace(anchor, setup).replace(snapshot, snapshot+'\t\tr982Observe(o, live)\n')
    main_target.write_text(main_text)
    subprocess.run(['gofmt', '-w', str(main_target), str(tree/'cmd/wiz7/r1002_context_probe.go')], check=True, timeout=20)
    env = os.environ.copy()
    env.update({'GOCACHE': '/gocache', 'GOMODCACHE': '/gomodcache', 'GOFLAGS': '-mod=readonly',
                'GOPROXY': 'off', 'GOSUMDB': 'off', 'GOTOOLCHAIN': 'local', 'GOMAXPROCS': '2'})
    completed = subprocess.run(['python3', str(tools/'build_recording.py'),
                                '--output', str(build/'wiz7-context-probe'), '--receipt-dir', str(build)],
                               cwd=tree, env=env, check=True, capture_output=True, text=True, timeout=160)
    identity = json.loads((build/'build-result.json').read_bytes())
    assert identity['source_count'] == 98
    # 再核所有正式/SDK輸入未變，來源表只在私人樹重生。
    for item in baseline:
        path = Path(item['path'])
        original = src/'worktrees/dosgolem'/Path(*path.parts[1:]) if path.parts[0] == 'sdk' else src/path
        assert sha(original) == item['sha256'] and original.stat().st_size == item['size']
    for output in (observer, main_target, tree/'cmd/wiz7/r1002_context_probe.go', build/'wiz7-context-probe'):
        assert output.stat().st_uid == os.getuid() and output.stat().st_gid == os.getgid()
    report = {'status': 'PRIVATE_IMPLEMENTATION_BUILT_SOURCE_REVIEW_REQUIRED', 'original_executed': False,
              'wire_or_original_data_read': False, 'production_or_sdk_modified': False,
              'formal_baseline_source_count': 96, 'formal_main_sha256': BASE_MAIN,
              'extension_contract_review_sha256': EXTENSION_REVIEW_SHA,
              'private_api_sha256': API_SHA, 'facts_sha256': FACTS_SHA, 'contract_review_sha256': REVIEW_SHA,
              'observer_sha256': sha(observer), 'constructor_sha256': sha(Path(__file__)),
              'private_main_sha256': sha(main_target), 'build_tool_sha256': sha(tools/'build_recording.py'),
              'build': identity, 'baseline_sources': baseline,
              'limits': {'fixed_hooks': 4, 'dynamic_return_hooks': 8, 'character_entries': 16, 'events': 256,
                         'source_window_steps': 1000000, 'text_bytes': 64,
                         'character_step_window': [332800000, 417800000], 'exported_character_pixels': [6, 6]},
              'enable_env': 'WIZ7_R1002_CONTEXT', 'off_installs_hooks': False,
              'main_delta': ['setup and deferred Close', 'approved r982Control/privacy log', 'snapshot r982Observe'],
              'runtime_validation': 'NOT_RUN; independent source review and root approval required',
              'build_stdout': json.loads(completed.stdout)}
    with report_path.open('x') as handle:
        handle.write(json.dumps(report, ensure_ascii=False, indent=2)+'\n')
    print(json.dumps({k: report[k] for k in ('status', 'original_executed', 'observer_sha256', 'constructor_sha256', 'build')}, ensure_ascii=False))


if __name__ == '__main__':
    main()
