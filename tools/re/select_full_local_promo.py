#!/usr/bin/env python3
"""Select two real intervals from the local runtime capture; no fake gameplay."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--capture", type=Path, required=True)
    p.add_argument("--out", type=Path, required=True)
    a = p.parse_args()
    assert a.capture.is_dir() and not a.out.exists()
    source = a.capture / "world.mkv"
    probed = subprocess.run(["ffprobe", "-v", "error", "-show_format", "-of", "json", str(source)], check=True, capture_output=True, text=True, timeout=15)
    assert float(json.loads(probed.stdout)["format"]["duration"]) >= 119.1
    a.out.mkdir()
    filters = ("[0:v]split=2[vh][vt];[vh]trim=duration=22,setpts=PTS-STARTPTS[v0];"
               "[vt]trim=start=85:duration=34,setpts=PTS-STARTPTS[v1];"
               "[0:a]asplit=2[ah][at];[ah]atrim=duration=22,asetpts=PTS-STARTPTS[a0];"
               "[at]atrim=start=85:duration=34,asetpts=PTS-STARTPTS[a1];"
               "[v0][a0][v1][a1]concat=n=2:v=1:a=1[v][a]")
    cmd = ["ffmpeg", "-n", "-hide_banner", "-loglevel", "warning", "-i", str(source),
           "-filter_complex_threads", "2", "-filter_complex", filters, "-map", "[v]", "-map", "[a]",
           "-t", "56", "-threads", "2", "-c:v", "libx264", "-preset", "veryfast", "-crf", "18",
           "-pix_fmt", "yuv420p", "-c:a", "pcm_s16le", str(a.out/"world.mkv")]
    encoded = subprocess.run(cmd, check=True, capture_output=True, text=True, timeout=150)
    (a.out/"encode.log").write_text(encoded.stderr)
    shutil.copy2(a.capture/"menu.mkv", a.out/"menu.mkv")
    result = {"source": str(source), "source_sha256": sha(source), "selected_world_sha256": sha(a.out/"world.mkv"),
              "menu_source_sha256": sha(a.capture/"menu.mkv"), "segments_seconds": [[0,22],[85,34]],
              "reason": "normal exploration then actual character/equipment UI; idle gap omitted",
              "new_or_repeated_gameplay_created": False, "audio": "same selected real captured intervals"}
    (a.out/"edit-receipt.json").write_text(json.dumps(result, indent=2)+"\n")
    print(json.dumps(result))


if __name__ == "__main__":
    main()
