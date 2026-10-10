#!/usr/bin/env python3
"""Verify the completed packaged Linux normal SAVE / cold LOAD receipts."""
import argparse
import hashlib
import json
from pathlib import Path
import stat

from PIL import Image


def sha(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()


def meta(p):
    s = p.stat()
    return {"size": s.st_size, "sha256": sha(p), "mode": stat.S_IMODE(s.st_mode), "mtime_ns": str(s.st_mtime_ns)}


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--repo", type=Path, required=True)
    p.add_argument("--bundle", type=Path, required=True)
    p.add_argument("--out", type=Path, required=True)
    a = p.parse_args()
    root = a.repo / "workplace/r1006"
    first, second = root / "capture", root / "coldload"
    runs = []
    for case in (first, second):
        operations = [json.loads(s) for s in (case/"operations.jsonl").read_text().splitlines()]
        exits = [r for r in operations if r["kind"] == "process_exit"]
        assert len(exits) == 1 and exits[0]["code"] == 0
        assert all(not r["error"] for r in operations if r["kind"] == "click")
        log = (case/"gui/run.log").read_text()
        assert "wizardry7-cht v.1.0.0-20261010" in log and "panic:" not in log
        assert "原版存檔目錄：/out/gui/save" in log
        runs.append({"case": case.name, "exit_code": 0, "operations": len(operations), "log_sha256": sha(case/"gui/run.log")})
    saved = {p.name: meta(p) for p in (first/"gui/save").iterdir()}
    loaded = {p.name: meta(p) for p in (second/"gui/save").iterdir()}
    source_before = json.loads((second/"initial-source.json").read_text())
    assert set(saved) == {"SAVEGAME.DBS", "PCFILE.DBS"}
    assert saved == loaded == source_before
    initial = a.repo/"workplace/r990/r3/record-gui/save/SAVEGAME.DBS"
    assert saved["SAVEGAME.DBS"]["sha256"] != sha(initial)
    images = [first/"gui/packaged-save-confirmed.png", second/"gui/packaged-world-restored.png", first/"gui/normal-world-loaded.png"]
    frames = [Image.open(p).convert("RGB") for p in images]
    assert all(im.size == (960, 600) for im in frames)
    rect = (216, 99, 744, 435)
    pixels = [list(im.crop(rect).getdata()) for im in frames]
    restored_diff = sum(x != y for x, y in zip(pixels[0], pixels[1]))
    initial_diff = sum(x != y for x, y in zip(pixels[0], pixels[2]))
    assert restored_diff == 0 and initial_diff > 0
    inventory = json.loads((root/"package/input-inventory.json").read_text())
    for item in inventory["original_files"]:
        assert sha(a.bundle/"game"/item["path"]) == item["sha256"]
    assert sha(a.bundle/"bin/wiz7") == "631926caa447010014ff70dada2dc56179a010eaa6fd5573394084b511484c99"
    result = {"status": "NORMAL_PACKAGED_LOAD_SAVE_COLDLOAD_PASS", "version": "v.1.0.0-20261010",
              "runs": runs, "save_files": saved, "saved_progress_changed": True,
              "viewport": list(rect), "viewport_pixels": len(pixels[0]), "restored_difference": restored_diff,
              "initial_state_contrast_difference": initial_diff, "original_files_unchanged": 248,
              "images": [{"path": str(p), "sha256": sha(p)} for p in images],
              "scope": "normal GUI and writable save-dir override; no RAM/same-state or default-path native claim"}
    with a.out.open("x") as f:
        json.dump(result, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(json.dumps(result, ensure_ascii=False))


if __name__ == "__main__":
    main()
