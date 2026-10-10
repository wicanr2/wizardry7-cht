#!/usr/bin/env python3
"""Prepare a clean private source tree; never execute a game or modify inputs."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--out", type=Path, required=True)
    parser.add_argument("--target", choices=("linux-amd64", "windows-amd64", "darwin-amd64", "darwin-arm64"), required=True)
    args = parser.parse_args()
    repo, dest = args.repo.resolve(), args.out / args.target
    assert repo.is_dir() and args.out.is_dir() and not dest.exists()
    assert args.out.stat().st_uid == 1000
    sources = [repo / "go.mod", repo / "go.sum"]
    for directory in ("cmd/wiz7", "initconfig", "recording"):
        sources.extend(p for p in (repo / directory).rglob("*") if p.is_file())
    before = {str(p.relative_to(repo)): sha(p) for p in sources}
    dest.mkdir()
    tree, build = dest / "source", dest / "build"
    tree.mkdir()
    build.mkdir()
    for directory in ("cmd/wiz7", "initconfig", "recording"):
        shutil.copytree(repo / directory, tree / directory)
    for name in ("go.mod", "go.sum", "tools/re/build_recording.py"):
        output = tree / name
        output.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(repo / name, output)
    sdk = repo / "worktrees/dosgolem"
    assert sdk.is_dir()
    (tree / "worktrees").mkdir()
    (tree / "worktrees/dosgolem").symlink_to(sdk, target_is_directory=True)
    assert before == {str(p.relative_to(repo)): sha(p) for p in sources}
    result = {"target": args.target, "source": str(tree), "build": str(build),
              "formal_inputs": before, "original_executed": False}
    (dest / "prepare-result.json").write_text(json.dumps(result, indent=2) + "\n")
    print(json.dumps({k: v for k, v in result.items() if k != "formal_inputs"}))


if __name__ == "__main__":
    main()
