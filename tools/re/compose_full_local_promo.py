#!/usr/bin/env python3
"""Compose a local-only trailer from newly recorded game video/audio."""
import argparse
import gzip
import hashlib
import json
from pathlib import Path
import shutil
import subprocess

from PIL import Image, ImageDraw

VERSION = "v.1.0.0-20261010"


def sha(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()


def run(command, timeout=120):
    if command[0] == "ffmpeg":
        bounded = [command[0], "-filter_threads", "2"]
        for argument in command[1:]:
            if argument == "-i":
                bounded.extend(["-threads", "1"])
            bounded.append(argument)
        command = bounded
    result = subprocess.run(command, check=False, capture_output=True, text=True, timeout=timeout)
    if result.returncode:
        raise RuntimeError("media_tool_failed: " + result.stderr[-6000:])
    return result


def probe(path):
    return json.loads(run(["ffprobe", "-v", "error", "-show_streams", "-show_format", "-of", "json", str(path)], 15).stdout)


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--capture", type=Path, required=True)
    p.add_argument("--font", type=Path, required=True)
    p.add_argument("--work", type=Path, required=True)
    p.add_argument("--out", type=Path, required=True)
    a = p.parse_args()
    assert not a.work.exists() and a.out.is_dir()
    movie_name = f"wizardry7-cht-{VERSION}-promo.mp4"
    assert not (a.out / movie_name).exists()
    menu, world = a.capture / "menu.mkv", a.capture / "world.mkv"
    pm, pw = probe(menu), probe(world)
    for item in (pm, pw):
        assert any(s["codec_type"] == "audio" for s in item["streams"])
        video = next(s for s in item["streams"] if s["codec_type"] == "video")
        assert (video["width"], video["height"]) == (960, 600)
    assert float(pm["format"]["duration"]) >= 4.8
    world_duration = min(56.0, float(pw["format"]["duration"]) - .2)
    assert world_duration >= 46
    total = 14 + world_duration
    assert 60 <= total <= 75
    assert sha(a.font) == "7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574"
    a.work.mkdir()
    texts = ["巫術7", "繁體中文化", VERSION, "本機完整版", "帶著你的隊伍，走入迷宮",
             "毛筆主選單 · 三倍清晰顯示", "原版規則 · 中文介面與訊息"]
    required = {ord(c) for t in texts for c in t}
    glyphs = {}
    with gzip.open(a.font, "rt", encoding="ascii") as f:
        for line in f:
            cp, bits = line.strip().split(":")
            code = int(cp, 16)
            if code in required:
                width = len(bits) // 4
                assert width in (8, 16)
                glyphs[code] = (width, [int(bits[y*width//4:(y+1)*width//4], 16) for y in range(16)])
    assert required <= glyphs.keys()
    safe = []
    def text(img, value, y, scale, color):
        width = sum(glyphs[ord(c)][0] for c in value) * scale
        x0 = (1280 - width) // 2
        assert x0 >= 40 and y >= 20 and y + 16*scale <= 700
        draw = ImageDraw.Draw(img)
        x = x0
        for c in value:
            w, rows = glyphs[ord(c)]
            for row, bits in enumerate(rows):
                for col in range(w):
                    if bits & (1 << (w-1-col)):
                        draw.rectangle((x+col*scale, y+row*scale, x+(col+1)*scale-1, y+(row+1)*scale-1), fill=color)
            x += w*scale
        safe.append({"text": value, "rectangle": [x0, y, width, 16*scale]})
    cream, gold = (248, 232, 160), (192, 146, 57)
    for label in ("title", "end"):
        img = Image.new("RGB", (1280, 720), (15, 17, 20))
        draw = ImageDraw.Draw(img)
        draw.rectangle((64, 64, 1215, 655), outline=gold, width=2)
        draw.line((340, 365, 940, 365), fill=gold, width=2)
        text(img, "巫術7", 200, 5, cream)
        text(img, "繁體中文化" if label == "title" else "帶著你的隊伍，走入迷宮", 310 if label == "title" else 400, 3 if label == "title" else 2, cream)
        text(img, VERSION, 480, 2, cream)
        text(img, "本機完整版", 540, 2, cream)
        img.save(a.work / (label + ".png"))
    for label, value in (("menu-caption", texts[-2]), ("world-caption", texts[-1])):
        img = Image.new("RGBA", (1280, 720), (0, 0, 0, 0))
        text(img, value, 664, 2, cream + (255,))
        img.save(a.work / (label + ".png"))
    command = ["ffmpeg", "-n", "-hide_banner", "-loglevel", "warning",
               "-loop", "1", "-framerate", "30", "-t", "4", "-i", str(a.work/"title.png"),
               "-i", str(menu), "-i", str(world),
               "-loop", "1", "-framerate", "30", "-t", "5", "-i", str(a.work/"end.png"),
               "-i", str(a.work/"menu-caption.png"), "-i", str(a.work/"world-caption.png")]
    filters = (f"[0:v]setsar=1,format=yuv420p[v0];"
               "[1:v]trim=duration=5,setpts=PTS-STARTPTS,pad=1280:720:160:40:color=0x0f1114,setsar=1[m];"
               "[m][4:v]overlay=shortest=0:repeatlast=1,format=yuv420p[v1];"
               f"[2:v]trim=duration={world_duration},setpts=PTS-STARTPTS,pad=1280:720:160:40:color=0x0f1114,setsar=1[w];"
               "[w][5:v]overlay=shortest=0:repeatlast=1,format=yuv420p[v2];"
               "[3:v]setsar=1,format=yuv420p[v3];[v0][v1][v2][v3]concat=n=4:v=1:a=0[v];"
               "[2:a]asplit=3[aw][at][ae];[at]atrim=duration=4,asetpts=PTS-STARTPTS[a0];"
               "[1:a]atrim=duration=5,asetpts=PTS-STARTPTS[a1];"
               f"[aw]atrim=duration={world_duration},asetpts=PTS-STARTPTS[a2];"
               f"[ae]atrim=start={world_duration-5}:duration=5,asetpts=PTS-STARTPTS[a3];"
               f"[a0][a1][a2][a3]concat=n=4:v=0:a=1,aformat=channel_layouts=stereo,loudnorm=I=-18:TP=-2:LRA=11,aformat=sample_rates=44100:channel_layouts=stereo,apad=whole_dur={total},atrim=duration={total},afade=t=in:d=1,afade=t=out:st={total-2}:d=2[a]")
    render = a.work / "render.mp4"
    command += ["-filter_complex_threads", "2", "-filter_complex", filters, "-map", "[v]", "-map", "[a]",
                "-t", str(total), "-r", "30", "-threads", "2", "-c:v", "libx264", "-preset", "veryfast",
                "-crf", "20", "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", "-movflags", "+faststart", str(render)]
    result = run(command, 180)
    (a.work/"encode.log").write_text(result.stderr)
    final_probe = probe(render)
    assert abs(float(final_probe["format"]["duration"]) - total) < .15
    video = next(s for s in final_probe["streams"] if s["codec_type"] == "video")
    audio = next(s for s in final_probe["streams"] if s["codec_type"] == "audio")
    assert (video["width"], video["height"], video["pix_fmt"]) == (1280, 720, "yuv420p")
    assert abs(float(video["duration"])-float(audio["duration"])) < .15
    scan = run(["ffmpeg", "-hide_banner", "-i", str(render), "-vf", "blackdetect=d=0.1:pix_th=0.05:pic_th=0.98,freezedetect=n=-60dB:d=2", "-an", "-f", "null", "-"], 90)
    volume = run(["ffmpeg", "-hide_banner", "-i", str(render), "-af", "volumedetect,silencedetect=noise=-50dB:d=2", "-vn", "-f", "null", "-"], 90)
    samples = [0, 4.5, 10, 20, 30, 40, total-6, total-2]
    for i, t in enumerate(samples):
        run(["ffmpeg", "-n", "-hide_banner", "-loglevel", "error", "-ss", str(t), "-i", str(render), "-frames:v", "1", str(a.work/f"frame-{i:02d}.png")], 15)
    select = "+".join(f"eq(n,{round(t*30)})" for t in samples)
    run(["ffmpeg", "-n", "-hide_banner", "-loglevel", "error", "-i", str(render), "-vf", f"select='{select}',scale=320:180,tile=4x2", "-frames:v", "1", str(a.work/"contact-sheet.png")], 90)
    metadata = {"version": VERSION, "local_only": True, "status": "RENDERED_VISUAL_ACCEPTANCE_PENDING",
                "duration": total, "movie_sha256": sha(render), "size": render.stat().st_size,
                "sources": [{"path": str(p), "sha256": sha(p)} for p in (menu, world, a.font)],
                "source_type": "new normal packaged frontend GUI and original OPL output captured via Pulse monitor",
                "clearance_segment_included": False, "prototype_used": False,
                "audio_edit": "actual captured audio only; title/end reuse short segments; normalization,44100Hz resampling,bounded end silence padding and fades",
                "static_segments": [[0,4,"title"],[total-5,total,"end"]], "text_safe_rectangles": safe}
    edit = a.capture / "edit-receipt.json"
    if edit.is_file():
        metadata["source_selection"] = json.loads(edit.read_text())
        assert metadata["source_selection"]["selected_world_sha256"] == sha(world)
        assert metadata["source_selection"]["menu_source_sha256"] == sha(menu)
    for name, data in (("ffprobe.json", final_probe), ("rights-and-source.json", metadata)):
        (a.work/name).write_text(json.dumps(data, ensure_ascii=False, indent=2)+"\n")
    (a.work/"black-freeze.log").write_text(scan.stderr)
    (a.work/"volume-silence.log").write_text(volume.stderr)
    shutil.copy2(render, a.out/movie_name)
    for source in a.work.iterdir():
        if source.name in ("ffprobe.json", "rights-and-source.json", "black-freeze.log", "volume-silence.log", "contact-sheet.png") or source.name.startswith("frame-"):
            assert not (a.out/source.name).exists()
            shutil.copy2(source, a.out/source.name)
    print(json.dumps({k:metadata[k] for k in ("version","status","duration","movie_sha256","size")}))


if __name__ == "__main__":
    main()
