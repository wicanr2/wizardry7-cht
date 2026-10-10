#!/usr/bin/env python3
"""Independently inspect and freshly unpack a local-only full ZIP."""
import argparse
import hashlib
import json
from pathlib import Path, PurePosixPath
import stat
import struct
import subprocess
import zipfile


def digest(data):
    return hashlib.sha256(data).hexdigest()


def architecture(data, target):
    if target == "linux-amd64":
        assert data[:6] == b"\x7fELF\x02\x01" and struct.unpack_from("<H", data, 18)[0] == 62
        return "ELF64 x86-64"
    if target == "windows-amd64":
        assert data[:2] == b"MZ"
        offset = struct.unpack_from("<I", data, 60)[0]
        assert data[offset:offset+4] == b"PE\0\0"
        assert struct.unpack_from("<H", data, offset+4)[0] == 0x8664
        assert struct.unpack_from("<H", data, offset+24)[0] == 0x20b
        return "PE32+ x86-64"
    assert data[:4] == b"\xcf\xfa\xed\xfe"
    cpu = struct.unpack_from("<I", data, 4)[0]
    assert cpu == (0x01000007 if target == "darwin-amd64" else 0x0100000c)
    return "Mach-O64 " + ("x86-64" if target == "darwin-amd64" else "arm64")


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--archive", type=Path, required=True)
    p.add_argument("--inventory", type=Path, required=True)
    p.add_argument("--out", type=Path, required=True)
    a = p.parse_args()
    assert not a.out.exists() and a.out.parent.is_dir()
    baseline = json.loads(a.inventory.read_text())
    a.out.mkdir()
    with zipfile.ZipFile(a.archive) as z:
        assert z.testzip() is None
        names = [x.filename for x in z.infolist()]
        assert len(names) == len(set(names))
        roots = {PurePosixPath(n).parts[0] for n in names}
        assert len(roots) == 1
        root = next(iter(roots))
        for i in z.infolist():
            q = PurePosixPath(i.filename)
            assert not q.is_absolute() and ".." not in q.parts and "\\" not in i.filename
            assert not stat.S_ISLNK(i.external_attr >> 16)
            if not i.filename.isascii():
                assert i.flag_bits & 0x800
        manifest = json.loads(z.read(root + "/MANIFEST.json"))
        assert manifest["local_only"] is True and manifest["version"] == baseline["version"]
        files = manifest["files"]
        paths = {x["path"] for x in files}
        assert len(paths) == len(files)
        assert set(names) == {root + "/" + n for n in paths | {"MANIFEST.json"}}
        original = {"game/" + x["path"]: x for x in baseline["original_files"]}
        assert {n for n in paths if n.startswith("game/")} == set(original)
        for path, expected in original.items():
            data = z.read(root + "/" + path)
            assert len(data) == expected["size"] and digest(data) == expected["sha256"]
        for expected in baseline["assets"]:
            filename = Path(expected["source"]).name
            if filename not in ("xlate.json", "roster.json", "profile-sex.json", "intro-picture.json", "flow.json", "help.json"):
                continue
            assert digest(z.read(root + "/assets/" + filename)) == expected["sha256"]
        assert {n for n in paths if n.startswith("assets/")} == {"assets/" + n for n in
               ("xlate.json", "roster.json", "profile-sex.json", "intro-picture.json", "flow.json", "help.json")}
        for item in files:
            data = z.read(root + "/" + item["path"])
            assert len(data) == item["size"] and digest(data) == item["sha256"]
        forbidden = (".wiz7-input", ".partial", "stdfont", "ascfont", "messages.json", ".log", ".state", ".pem")
        assert not any(any(t in n.lower() for t in forbidden) or n.startswith(("workplace/", "text/", ".git/")) for n in paths)
        allowed_top = {"bin", "game", "assets", "licenses", "VERSION.txt", "README.txt", "start.sh", "start.command", "start.bat"}
        assert all(PurePosixPath(n).parts[0] in allowed_top for n in paths)
        target = manifest["target"]
        bin_name = "wiz7.exe" if target.startswith("windows") else "wiz7"
        data = z.read(root + "/bin/" + bin_name)
        kind = architecture(data, target)
        assert digest(data) == manifest["binary_identity"]["binary_sha256"]
        assert z.read(root + "/VERSION.txt").decode().strip() == baseline["version"]
        readme = z.read(root + "/README.txt")
        if target.startswith("windows"):
            assert readme.startswith(b"\xef\xbb\xbf") and b"\r\n" in readme
            bat = z.read(root + "/start.bat")
            assert bat.isascii() and not bat.startswith(b"\xef\xbb\xbf") and b"\r\n" in bat
        for name in ("RRSAL-1.0.txt", "dosgolem-LICENSE.txt", "ARPHIC-copyright.txt", "Unifont-OFL-1.1.txt", "Unifont-COPYING.txt", "Unifont-README.txt", "Go-LICENSE.txt", "THIRD-PARTY.txt"):
            assert root + "/licenses/" + name in names
        for i in z.infolist():
            z.extract(i, a.out)
            (a.out / i.filename).chmod((i.external_attr >> 16) & 0o777 or 0o644)
    tree = a.out / root
    version_output = None
    if not target.startswith("windows"):
        launcher = tree / ("start.command" if target.startswith("darwin") else "start.sh")
        version_output = subprocess.check_output([str(launcher), "--version"], text=True, timeout=10).strip()
        assert version_output == "wizardry7-cht " + baseline["version"]
    result = {"status": "PACKAGE_STRUCTURE_PASS", "target": target, "version": baseline["version"],
              "archive_sha256": digest(a.archive.read_bytes()), "archive_size": a.archive.stat().st_size,
              "files": len(names), "original_files": len(original), "architecture": kind,
              "launcher_version": version_output, "extracted": str(tree), "native_game_started": False}
    with (a.out / "structure-result.json").open("x") as f:
        json.dump(result, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(json.dumps(result, ensure_ascii=False))


if __name__ == "__main__":
    main()
