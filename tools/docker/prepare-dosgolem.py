#!/usr/bin/env python3
"""容器內重建已審 dosgolem 音訊增量；原版遊戲與主機上游均不讀寫。"""
from pathlib import Path
import hashlib
import json
import os
import subprocess

project = Path(__file__).resolve().parents[2]
repo = project / 'worktrees/dosgolem'
entry = project / 'tools/docker'
manifest = json.loads((entry / 'dosgolem-audio-manifest.json').read_text())
patch = entry / 'dosgolem-audio.patch'
def git(*args):
    return subprocess.check_output(['git', '-C', str(repo), *args], text=True).strip()
def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()
def matches():
    return all((repo / row['path']).is_file() and sha(repo / row['path']) == row['sha256'] for row in manifest['files'])

if not repo.is_dir():
    raise SystemExit('找不到 worktrees/dosgolem；請先依 README 準備指定基線。')
if repo.resolve() != repo or not (repo / '.git').is_dir() or (repo / '.git').is_symlink():
    raise SystemExit('worktrees/dosgolem 必須是專案內的實體clone；不跟隨符號連結或外部Git工作樹。')
for row in manifest['files']:
    path = repo / row['path']
    if path.resolve() != path:
        raise SystemExit('來源目標含符號連結；停止，不讀寫外部目標。')
if git('rev-parse', 'HEAD') != manifest['base_commit']:
    raise SystemExit('dosgolem 基線不符；停止，不切換既有分支或丟棄工作。')
if git('rev-parse','--abbrev-ref','HEAD') != manifest['branch']:
    raise SystemExit('請依 README 準備隔離分支；停止，不修改main或master。')
allowed = {row['path'] for row in manifest['files']}
changed = set(git('diff', '--name-only', 'HEAD').splitlines())
if changed - allowed:
    raise SystemExit('dosgolem 有此增量以外的修改；停止，保留原工作。')
if sha(patch) != manifest['patch_sha256']:
    raise SystemExit('增量雜湊不符；停止。')
if matches():
    print('dosgolem 音訊來源已符合固定清單，沒有新增修改。')
else:
    if changed:
        raise SystemExit('來源不是乾淨基線或完整已驗增量；停止，保留既有修改。')
    # 在任何寫入前核對每個現存目標與其最近現存父目錄；不以root修檔。
    for row in manifest['files']:
        path = repo / row['path']
        parent = path.parent
        while not parent.exists():
            parent = parent.parent
        for check in [parent] + ([path] if path.exists() else []):
            stat = check.stat()
            if (stat.st_uid,stat.st_gid) != (os.geteuid(),os.getegid()):
                raise SystemExit('來源目標或父目錄擁有權不符；停止，不變更擁有權。')
    subprocess.run(['git', '-C', str(repo), 'apply', '--check', str(patch)], check=True)
    subprocess.run(['git', '-C', str(repo), 'apply', str(patch)], check=True)
    if not matches():
        raise SystemExit('套用後來源雜湊不符；停止，不能進行建置。')
    print('dosgolem 音訊增量重建完成，13 個固定來源雜湊相符。')
