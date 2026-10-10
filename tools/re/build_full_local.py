#!/usr/bin/env python3
"""Create one local-only full package from a pinned build and verified inputs."""
import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import stat
import zipfile


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def dump(path, value):
    with path.open("x", encoding="utf-8") as f:
        json.dump(value, f, ensure_ascii=False, indent=2)
        f.write("\n")


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--repo", type=Path, required=True)
    ap.add_argument("--dist", type=Path, required=True)
    ap.add_argument("--font", type=Path, required=True)
    ap.add_argument("--inventory", type=Path)
    ap.add_argument("--version", required=True)
    ap.add_argument("--target", choices=("linux-amd64", "windows-amd64", "darwin-amd64", "darwin-arm64"), required=True)
    a = ap.parse_args()
    assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", a.version)
    assert a.dist.is_dir() and a.dist.stat().st_uid == os.getuid() == 1000
    repo, target = a.repo, a.target
    platform, arch = target.split("-")
    inventory = json.loads((a.inventory or repo / "tools/re/full_local_inputs.json").read_text())
    assert inventory["version"] == a.version
    build = repo / "workplace/r1006/platform" / target / "build"
    identity = json.loads((build / "build-result.json").read_text())
    binary = build / ("wiz7.exe" if platform == "windows" else "wiz7")
    assert identity["goos"] == platform and identity["goarch"] == arch
    assert sha(binary) == identity["binary_sha256"]
    if target == "linux-amd64":
        assert sha(binary) == "631926caa447010014ff70dada2dc56179a010eaa6fd5573394084b511484c99"
    root = a.dist / a.version
    root.mkdir(exist_ok=True)
    assert root.stat().st_uid == 1000
    for name in ("full-local", "promo", "smoke"):
        (root / name).mkdir(exist_ok=True)
    label = target.replace("darwin", "macos")
    name = f"wizardry7-cht-{a.version}-full-local-{label}"
    tree = root / "full-local" / name
    archive = root / "full-local" / (name + ".zip")
    assert not tree.exists() and not archive.exists()
    tree.mkdir()
    for directory in ("bin", "game", "assets", "licenses"):
        (tree / directory).mkdir()
    shutil.copy2(binary, tree / "bin" / binary.name)
    (tree / "bin" / binary.name).chmod(0o755)
    original = repo / "workplace/original/wizar7"
    for item in inventory["original_files"]:
        relative = PurePosixPath(item["path"])
        assert not relative.is_absolute() and ".." not in relative.parts
        source = original / Path(*relative.parts)
        assert not source.is_symlink() and sha(source) == item["sha256"] and source.stat().st_size == item["size"]
        dest = tree / "game" / Path(*relative.parts)
        dest.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source, dest)
        dest.chmod(0o444)
    sources = {p["source"]: p for p in inventory["assets"]}
    asset_paths = ["workplace/r995/party-assets/unifont/" + n for n in
                   ("xlate.json", "roster.json", "profile-sex.json", "intro-picture.json")]
    asset_paths += ["workplace/r983/final/unifont/flow.json", "workplace/r979/unifont/help.json"]
    for name_in in asset_paths:
        source = repo / name_in
        assert sha(source) == sources[name_in]["sha256"]
        shutil.copy2(source, tree / "assets" / source.name)
    notices = [(repo / "LICENSE", "RRSAL-1.0.txt"),
               (repo / "worktrees/dosgolem/LICENSE", "dosgolem-LICENSE.txt"),
               (repo / "workplace/r976/menu/font-copyright.txt", "ARPHIC-copyright.txt")]
    for filename in ("COPYING", "OFL-1.1.txt", "README"):
        source = a.font / filename
        expected = next(x for x in inventory["unifont_licenses"] if Path(x["source"]).name == filename)
        assert sha(source) == expected["sha256"]
        notices.append((source, "Unifont-" + filename + ("" if filename.endswith(".txt") else ".txt")))
    notices.append((Path("/usr/local/go/LICENSE"), "Go-LICENSE.txt"))
    for module, version in re.findall(r"^\s*(\S+)\s+(v\S+)\s*(?://[^\n]*)?$", (repo / "go.mod").read_text(), re.M):
        if module == "github.com/wicanr2/dosgolem":
            continue
        folder = repo / "workplace/gomodcache" / (module + "@" + version)
        licenses = [p for p in folder.iterdir() if p.is_file() and p.name.upper().startswith(("LICENSE", "COPYING", "NOTICE"))]
        assert licenses, module
        notices.extend((p, module.replace("/", "_") + "-" + p.name) for p in licenses)
    for source, output in notices:
        shutil.copy2(source, tree / "licenses" / output)
    (tree / "VERSION.txt").write_text(a.version + "\n", encoding="ascii")
    if platform == "windows":
        launcher = "start.bat"
        text = ('@echo off\nchcp 65001 >nul\n'
                f'echo wizardry7-cht {a.version}\n'
                'if "%~1"=="--version" exit /b 0\n'
                '"%~dp0bin\\wiz7.exe" -root "%~dp0game" -xlate "%~dp0assets\\xlate.json" -flow "%~dp0assets\\flow.json" -help "%~dp0assets\\help.json" %*\n'
                'exit /b %ERRORLEVEL%\n')
        (tree / launcher).write_bytes(text.replace("\n", "\r\n").encode("ascii"))
    else:
        launcher = "start.command" if platform == "darwin" else "start.sh"
        text = ('#!/bin/sh\nset -eu\n'
                f'printf "%s\\n" "wizardry7-cht {a.version}"\n'
                'if [ "${1-}" = "--version" ]; then exit 0; fi\n'
                'base=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)\n'
                'exec "$base/bin/wiz7" -root "$base/game" -xlate "$base/assets/xlate.json" -flow "$base/assets/flow.json" -help "$base/assets/help.json" "$@"\n')
        (tree / launcher).write_text(text, encoding="ascii")
        (tree / launcher).chmod(0o755)
    location = {"linux": "$XDG_DATA_HOME/wizardry7-cht/save，或 ~/.local/share/wizardry7-cht/save",
                "windows": "%LOCALAPPDATA%\\wizardry7-cht\\save",
                "darwin": "~/Library/Application Support/wizardry7-cht/save"}[platform]
    readme = f"""巫術7 繁體中文化
{a.version}

本機完整版，已包含使用者自備的原版遊戲。請勿公開分享此封包。
解壓後執行 {launcher}。預設為960×600，使用Unifont中文及毛筆主選單。
Linux可在終端機執行 ./start.sh；macOS請在終端機執行 ./start.command。
F1開啟說明，F2開關地圖。存檔與讀檔使用遊戲原有選單。
存檔預設位置：{location}。
可加 -save-dir 指定獨立存檔目錄。更新程式前請備份整個存檔目錄。
對話關鍵字與防拷輸入使用英文；提示只影響顯示，原版驗證保持。

已知限制：尚未完成全遊戲文字與結局驗收、TSV遷移及作弊功能。
兩倍創角屬性文字仍可能重疊；建議使用預設三倍。
Windows以Wine輔助驗收，不等同Windows真機；macOS原生啟動仍待實機確認。
完整版指包含遊戲資料，不代表上述未完成項目已完成。

程式與中文化：RRSAL-1.0。字型及第三方元件見licenses。
原版遊戲、美術與音樂的權利屬各自權利人。
"""
    data = readme.replace("\n", "\r\n").encode("utf-8-sig") if platform == "windows" else readme.encode()
    (tree / "README.txt").write_bytes(data)
    (tree / "licenses/THIRD-PARTY.txt").write_text(
        "Chinese glyph subset: derived from GNU Unifont 17.0.05, SIL OFL 1.1.\n"
        "Main menu glyph subset: derived from AR PL UKai TW, ARPHIC PUBLIC LICENSE.\n"
        "Original font data remains unmodified; generated glyph subsets change representation.\n"
        "The original game and its assets are not licensed by RRSAL-1.0.\n", encoding="utf-8")
    files = [{"path": str(p.relative_to(tree)), "size": p.stat().st_size, "sha256": sha(p)}
             for p in sorted(tree.rglob("*")) if p.is_file()]
    dump(tree / "MANIFEST.json", {"version": a.version, "target": target, "local_only": True,
         "binary_identity": identity, "original_zip_sha256": inventory["original_zip_sha256"], "files": files})
    with zipfile.ZipFile(archive, "x", compression=zipfile.ZIP_DEFLATED, compresslevel=6) as z:
        for p in sorted(tree.rglob("*")):
            if p.is_file():
                z.write(p, str(p.relative_to(tree.parent)))
    with zipfile.ZipFile(archive) as z:
        assert z.testzip() is None
        for i in z.infolist():
            path = PurePosixPath(i.filename)
            assert not path.is_absolute() and ".." not in path.parts
            assert not stat.S_ISLNK(i.external_attr >> 16)
            if not i.filename.isascii():
                assert i.flag_bits & 0x800
        assert len(z.infolist()) == len(files) + 1
    result = {"version": a.version, "target": target, "local_only": True,
              "archive": str(archive.relative_to(root)), "size": archive.stat().st_size,
              "sha256": sha(archive), "tree": str(tree), "files": len(files)+1,
              "launcher": launcher, "binary_sha256": sha(binary)}
    dump(root / "full-local" / (name + ".json"), result)
    print(json.dumps(result, ensure_ascii=False))


if __name__ == "__main__":
    main()
