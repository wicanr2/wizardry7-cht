#!/usr/bin/env python3
"""在隔離媒體工具鏈中，合成介紹卡與完整正常遊玩對話。"""
import argparse
import csv
import gzip
import hashlib
import json
import math
from pathlib import Path
import re
import subprocess

from PIL import Image, ImageDraw

FONT_SHA = "7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574"
BG, CREAM, GOLD = (15, 17, 20), (248, 232, 160), (192, 146, 57)


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def run(command, timeout=180):
    result = subprocess.run(command, capture_output=True, text=True, timeout=timeout)
    if result.returncode:
        raise RuntimeError(result.stderr[-6000:])
    return result


def probe(path):
    return json.loads(run(["ffprobe", "-v", "error", "-show_streams", "-show_format",
                           "-of", "json", str(path)], 20).stdout)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--plan", type=Path, required=True)
    parser.add_argument("--font", type=Path, required=True)
    parser.add_argument("--work", type=Path, required=True)
    parser.add_argument("--copy", type=Path, default=Path(__file__).with_name("promo_story_features.zh-TW.tsv"))
    args = parser.parse_args()
    assert not args.work.exists(), "保留各輪輸出，請使用新的工作目錄"
    assert sha(args.font) == FONT_SHA
    plan = json.loads(args.plan.read_text())
    version = plan["version"]
    assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", version)
    assert plan["local_only"] is True
    with args.copy.open() as file:
        copy = {row["id"]: row for row in csv.DictReader(file, delimiter="\t")}
    copy = {key: {field: value.replace("{version}", version) for field, value in row.items()}
            for key, row in copy.items()}
    assert len(copy) == 10
    required = {ord(c) for row in copy.values() for key, value in row.items()
                if key != "id" for c in value}
    glyphs = {}
    with gzip.open(args.font, "rt", encoding="ascii") as file:
        for line in file:
            cp, bits = line.strip().split(":")
            code = int(cp, 16)
            if code in required:
                width = len(bits) // 4
                assert width in (8, 16)
                glyphs[code] = width, [int(bits[y*width//4:(y+1)*width//4], 16) for y in range(16)]
    assert required <= glyphs.keys(), "影片字型缺字"
    args.work.mkdir()
    safe = []

    def text(image, value, y, scale, color=CREAM):
        width = sum(glyphs[ord(c)][0] for c in value) * scale
        x = (1280-width)//2
        assert 40 <= x and 20 <= y and y+16*scale <= 700, value
        safe.append({"text": value, "rectangle": [x, y, width, 16*scale]})
        draw = ImageDraw.Draw(image)
        for char in value:
            w, rows = glyphs[ord(char)]
            for row, bits in enumerate(rows):
                for col in range(w):
                    if bits & (1 << (w-1-col)):
                        draw.rectangle((x+col*scale, y+row*scale,
                                        x+(col+1)*scale-1, y+(row+1)*scale-1), fill=color)
            x += w*scale

    sources = {}

    def source(item):
        path = args.root / item["path"]
        assert path.is_file() and sha(path) == item["sha256"], str(path)
        sources[item["path"]] = {"path": item["path"], "sha256": item["sha256"], "size": path.stat().st_size}
        return path

    npc_receipt = json.loads(source(plan["npc_receipt"]).read_text())
    music = source(plan["music"])
    music_duration = float(probe(music)["format"]["duration"])
    assert music_duration >= 5
    elapsed, timeline, commands, encoded = 0.0, [], [], []
    segments = plan["segments"]
    assert sum(s["id"] == "npc" for s in segments) == 1
    npc = next(s for s in segments if s["id"] == "npc")
    assert npc["audio"] == "source"
    assert npc_receipt["complete"] is True and npc_receipt["normal_player_path"] is True
    assert npc["path"] == npc_receipt["source_path"]
    assert npc["sha256"] == npc_receipt["source_sha256"]
    start, end = float(npc.get("start", 0)), float(npc.get("start", 0))+float(npc["duration"])
    event_start, event_end = float(npc_receipt["event_start"]), float(npc_receipt["event_end"])
    assert all(math.isfinite(t) for t in (start, end, event_start, event_end))
    assert 0 <= start <= event_start < event_end <= end, "不得裁掉完整事件起訖"
    assert npc_receipt["page_intervals"], "須有實錄逐頁驗收"
    for page_start, page_end in npc_receipt["page_intervals"]:
        assert event_start <= page_start < page_end <= event_end, "文字頁須在完整事件範圍內"
    for index, segment in enumerate(segments):
        label, kind, duration = segment["id"], segment["kind"], float(segment["duration"])
        assert math.isfinite(duration) and duration >= 1
        assert abs(duration*30-round(duration*30)) < 1e-5, "段長須對齊30fps"
        row = copy[segment.get("copy", label)]
        png = args.work / f"{index:02d}-{label}.png"
        command = ["ffmpeg", "-n", "-hide_banner", "-loglevel", "warning", "-filter_threads", "2"]
        if kind == "card":
            image = Image.new("RGB", (1280, 720), BG)
            draw = ImageDraw.Draw(image)
            draw.rectangle((64, 64, 1215, 655), outline=GOLD, width=2)
            text(image, row["heading"], 160, 4 if label != "title" else 5)
            draw.line((200, 280, 1080, 280), fill=GOLD, width=2)
            for y, key in zip((330, 410, 490), ("line1", "line2", "line3")):
                text(image, row[key], y, 2)
            image.save(png)
            command += ["-loop", "1", "-framerate", "30", "-threads", "1", "-i", str(png)]
            video_filter = "[0:v]setsar=1,format=yuv420p[v]"
        else:
            path = source(segment)
            if kind == "still":
                command += ["-loop", "1", "-framerate", "30"]
            else:
                assert kind == "video"
                info = probe(path)
                assert float(info["format"]["duration"]) >= segment.get("start", 0)+duration-.05
                command += ["-ss", str(segment.get("start", 0))]
            command += ["-threads", "1", "-i", str(path)]
            info = probe(path)
            video = next(s for s in info["streams"] if s["codec_type"] == "video")
            assert (video["width"], video["height"]) == (960, 600)
            image = Image.new("RGBA", (1280, 720), (0, 0, 0, 0))
            text(image, row["line1"], 664, 2, CREAM+(255,))
            image.save(png)
            command += ["-threads", "1", "-i", str(png)]
            video_filter = ("[0:v]setpts=PTS-STARTPTS,fps=30,pad=1280:720:160:40:color=0x0f1114,"
                            "setsar=1[base];[base][1:v]overlay=repeatlast=1,format=yuv420p[v]")
        if segment["audio"] == "music":
            audio_index = 1 if kind == "card" else 2
            command += ["-stream_loop", "-1", "-ss", str(elapsed % music_duration), "-threads", "1", "-i", str(music)]
            gain = float(plan["music"].get("gain", 1))
        else:
            assert kind == "video" and segment["audio"] == "source"
            audio_index, gain = 0, float(segment.get("gain", 1))
            assert any(s["codec_type"] == "audio" for s in info["streams"])
        assert 0 < gain <= 2, "固定有界增益，不追逐安靜段響度"
        fade_in = .4 if index == 0 else .05
        fade_out = 1.5 if label == "end" else .05
        audio_filter = (f"[{audio_index}:a]aresample=44100:async=1:first_pts=0,"
                        f"aformat=channel_layouts=stereo,highpass=f=20:p=2,volume={gain},"
                        f"apad=whole_dur={duration},atrim=duration={duration},"
                        f"afade=t=in:d={fade_in},afade=t=out:st={duration-fade_out}:d={fade_out}[a]")
        output = args.work / f"{index:02d}-{label}.mkv"
        command += ["-filter_complex_threads", "2", "-filter_complex", video_filter+";"+audio_filter,
                    "-map", "[v]", "-map", "[a]", "-t", str(duration), "-r", "30", "-threads", "2",
                    "-c:v", "libx264", "-preset", "veryfast", "-crf", "20", "-pix_fmt", "yuv420p",
                    "-c:a", "pcm_s16le", str(output)]
        result = run(command, max(180, math.ceil(duration*3)))
        (args.work / f"{index:02d}-{label}.log").write_text(result.stderr)
        commands.append(command)
        encoded.append(output)
        timeline.append({**segment, "movie_start": elapsed, "movie_end": elapsed+duration})
        elapsed += duration
    listing = args.work / "concat.txt"
    assert all("'" not in str(path) for path in encoded)
    listing.write_text("".join(f"file '{path}'\n" for path in encoded))
    movie = args.work / f"wizardry7-cht-{version}-promo-story-features.mp4"
    command = ["ffmpeg", "-n", "-hide_banner", "-loglevel", "warning", "-threads", "1",
               "-f", "concat", "-safe", "0", "-i", str(listing), "-map", "0:v", "-map", "0:a",
               "-c:v", "copy", "-c:a", "aac", "-b:a", "192k", "-threads", "2", "-t", str(elapsed),
               "-movflags", "+faststart", str(movie)]
    result = run(command, max(180, math.ceil(elapsed)))
    (args.work / "mux.log").write_text(result.stderr)
    final = probe(movie)
    video = next(s for s in final["streams"] if s["codec_type"] == "video")
    audio = next(s for s in final["streams"] if s["codec_type"] == "audio")
    assert (video["width"], video["height"], video["r_frame_rate"]) == (1280, 720, "30/1")
    assert audio["sample_rate"] == "44100" and audio["channels"] == 2
    assert abs(float(video["duration"])-elapsed) < .05
    assert abs(float(video["duration"])-float(audio["duration"])) < .1
    receipt = {"version": version, "local_only": True, "status": "RENDERED_ACCEPTANCE_PENDING",
               "movie": movie.name, "sha256": sha(movie), "size": movie.stat().st_size,
               "duration": elapsed, "timeline": timeline, "sources": list(sources.values()),
               "font_sha256": FONT_SHA, "copy_sha256": sha(args.copy), "script_sha256": sha(Path(__file__)),
               "plan_sha256": sha(args.plan), "text_safe_rectangles": safe, "commands": commands+[command],
               "audio_edit": "原版實錄；20Hz二階去DC、固定有界增益、各段淡入淡出；依PTS補音訊，不作響度正規化",
               "npc_receipt": plan["npc_receipt"], "clearance_included": False}
    (args.work / "ffprobe.json").write_text(json.dumps(final, indent=2)+"\n")
    (args.work / "rights-and-source.json").write_text(json.dumps(receipt, ensure_ascii=False, indent=2)+"\n")
    print(json.dumps({k: receipt[k] for k in ("movie", "sha256", "size", "duration", "status")}))


if __name__ == "__main__":
    main()
