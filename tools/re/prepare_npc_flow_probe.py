#!/usr/bin/env python3
"""R1007：在私人來源量測NPC印字，正式來源及SDK保持唯讀。"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess

MAIN_SHA = "22a0134eda3df0ca42f9b361ef204dbea94ebddabec6e2bd56775d7bbf2b2d7a"
FLOW_SHA = "c73d9c50ae45a85b6e7d853f8e093b5c533833a18b6dacab9f2f483ad585963e"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--src", type=Path, default=Path("/src"))
    parser.add_argument("--out", type=Path, required=True)
    args = parser.parse_args()
    src, out = args.src, args.out
    assert src.is_dir() and not out.exists() and out.parent.stat().st_uid == os.getuid() == 1000
    assert sha(src/"cmd/wiz7/main.go") == MAIN_SHA and sha(src/"cmd/wiz7/flow.go") == FLOW_SHA
    table = src/"cmd/wiz7/recording_source_paths.go"
    paths = re.findall(r'^\s*"([^"\n]+)",\s*$', table.read_text(), re.M)
    assert len(paths) == len(set(paths)) == 97
    out.mkdir()
    tree, build = out/"source", out/"build"
    tree.mkdir()
    build.mkdir()
    baseline = []
    for relative in paths:
        path = Path(relative)
        assert not path.is_absolute() and ".." not in path.parts
        original = src/"worktrees/dosgolem"/Path(*path.parts[1:]) if path.parts[0] == "sdk" else src/path
        assert original.is_file() and not original.is_symlink()
        baseline.append({"path": relative, "sha256": sha(original), "size": original.stat().st_size})
        if path.parts[0] != "sdk":
            target = tree/path
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(original, target)
    (tree/"worktrees").mkdir()
    (tree/"worktrees/dosgolem").symlink_to(src/"worktrees/dosgolem", target_is_directory=True)
    (tree/"tools/re").mkdir(parents=True)
    shutil.copy2(src/"tools/re/build_recording.py", tree/"tools/re/build_recording.py")
    flow = tree/"cmd/wiz7/flow.go"
    anchor = "var flowSites = map[oracle.Addr]flowSite{\n"
    assert flow.read_text().count(anchor) == 1
    flow.write_text(flow.read_text().replace(anchor, anchor+
                    "\t// PRIVATE SELECTORS ONLY: no page, geometry, translation, or style contract.\n"
                    "\t{Seg: 0x0110, Off: 0xE3A6}: {},\n"
                    "\t{Seg: 0x0110, Off: 0x5F1E}: {},\n"))
    entry_anchor = "\t// 頁：首片段 y0＝146；之後 y0＝上一片段＋12；至多 4 個片段\n"
    assert flow.read_text().count(entry_anchor) == 1
    flow.write_text(flow.read_text().replace(entry_anchor,
                    "\tif ret == oracle.Far(0x0110, 0xE3A6) || ret == oracle.Far(0x0110, 0x5F1E) {\n"
                    "\t\treturn &pendingPrint{s: s, x0: x, y0: y, sp: sp, ret: ret, font: font, seq: seq}\n"
                    "\t}\n\n"+entry_anchor))
    main_file = tree/"cmd/wiz7/main.go"
    anchor = "\te.fgRGB = o.Palette()[e.fg] // 返回 hook 當下的調色盤（R2.4）\n"
    assert main_file.read_text().count(anchor) == 1
    diagnostic = '''\tif p.ret == oracle.Far(0x0110, 0xE3A6) || p.ret == oracle.Far(0x0110, 0x5F1E) {
        if dir := os.Getenv("WIZ7_NPC_PROBE_DIR"); dir != "" && p.seq > 0 && p.seq <= 128 {
            reads := []flowRead{}
            for _, read := range ov.reads { if read.ID >= 29750 && read.ID < 29800 { reads = append(reads, read) } }
            value := map[string]any{"step": o.Steps(), "ret": p.ret, "font": p.font, "x": p.x0, "y": p.y0,
                "text": p.s, "bbox": [4]int{x0, y0, w, h}, "before": e.before, "after": e.after,
                "colour_counts": cnt, "fg": e.fg, "fg_rgb": e.fgRGB, "seq": p.seq, "observed_source_reads": reads}
            if data, err := json.Marshal(value); err == nil {
                if err := os.WriteFile(fmt.Sprintf("%s/npc-%03d.json", dir, p.seq), data, 0o644); err != nil { log.Fatal(err) }
            } else { log.Fatal(err) }
        }
        return // PRIVATE: do not register an overlay event or apply any page/style gate.
    }
'''
    main_file.write_text(main_file.read_text().replace(anchor, anchor+diagnostic))
    subprocess.run(["gofmt", "-w", str(flow), str(main_file)], check=True, timeout=20)
    env = os.environ.copy()
    env.update({"GOCACHE": "/gocache", "GOMODCACHE": "/gomodcache", "GOFLAGS": "-mod=readonly",
                "GOPROXY": "off", "GOSUMDB": "off", "GOTOOLCHAIN": "local", "GOMAXPROCS": "2"})
    completed = subprocess.run(["python3", str(tree/"tools/re/build_recording.py"),
                                "--output", str(build/"wiz7-npc-probe"), "--receipt-dir", str(build)],
                               cwd=tree, env=env, check=True, capture_output=True, text=True, timeout=180)
    identity = json.loads((build/"build-result.json").read_text())
    assert identity["source_count"] == 97
    for entry in baseline:
        path = Path(entry["path"])
        original = src/"worktrees/dosgolem"/Path(*path.parts[1:]) if path.parts[0] == "sdk" else src/path
        assert sha(original) == entry["sha256"]
    report = {"status": "PRIVATE_PROBE_BUILT_SOURCE_REVIEW_REQUIRED", "original_executed": False,
              "formal_or_sdk_modified": False, "constructor_sha256": sha(Path(__file__)), "build": identity,
              "formal_main_sha256": MAIN_SHA, "formal_flow_sha256": FLOW_SHA, "baseline_sources": baseline,
              "private_main_sha256": sha(main_file), "private_flow_sha256": sha(flow),
              "deltas": ["two private caller selectors with zero geometry", "page=nil pending after NUL validation",
                         "export existing original diff/source fields", "return before page/style gate and overlay registration"],
              "required_runtime": {"flow": "", "normal_load_then_one_step": True,
                                   "max_fragment_exports": 128, "WIZ7_NPC_PROBE_DIR": "existing private output directory"},
              "claims": "No Chinese output, parity, or production implementation is authorized by this probe."}
    (out/"implementation-report.json").write_text(json.dumps(report, ensure_ascii=False, indent=2)+"\n")
    print(completed.stdout)


if __name__ == "__main__":
    main()
