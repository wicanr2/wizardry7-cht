#!/usr/bin/env python3
"""spec007 的真實視窗、F1翻頁與故事Draw抽測；Docker／Xvfb內執行。"""
import argparse
import hashlib
import json
import subprocess
import time
from pathlib import Path

ap = argparse.ArgumentParser()
ap.add_argument("out", type=Path)
ap.add_argument("--overlay-only", action="store_true", help="只抽測主選單與防拷說明的正常Draw及固定步數")
args = ap.parse_args()
out = args.out
out.mkdir(exist_ok=False)
binary = out / "wiz7"
subprocess.run(["go", "build", "-o", str(binary), "./cmd/wiz7"], check=True)


def wait_for(predicate, proc, seconds=180):
    until = time.monotonic() + seconds
    while time.monotonic() < until:
        if predicate():
            return
        if proc.poll() is not None:
            raise RuntimeError(f"視窗程序提前結束：{proc.returncode}")
        time.sleep(.1)
    raise RuntimeError("視窗驗收逾時")


def stop(proc):
    proc.terminate()
    try:
        proc.wait(timeout=5)
    except subprocess.TimeoutExpired:
        proc.kill()
        proc.wait()


def help_frame(wid, dest, proc, previous=None):
    """等實際說明畫格更新且穩定；不把尚未呈現的原版畫面當收據。"""
    candidate = dest.with_suffix(".pending.png")
    last = None
    until = time.monotonic() + 20
    while time.monotonic() < until:
        if proc.poll() is not None:
            raise RuntimeError("說明翻頁時視窗程序提前結束")
        subprocess.run(["import", "-window", wid, str(candidate)], check=True)
        raw = subprocess.check_output(["convert", str(candidate), "-depth", "8", "rgba:-"])
        digest = hashlib.sha256(raw).hexdigest()
        pixels = {raw[i:i + 4] for i in range(0, len(raw), 4)}
        visible = all(c in pixels for c in (bytes((255, 221, 68, 255)),
                                            bytes((240, 240, 240, 255)),
                                            bytes((160, 160, 160, 255))))
        if visible and digest != previous and digest == last:
            candidate.replace(dest)
            return digest
        last = digest
        time.sleep(.15)
    raise RuntimeError("說明畫格未更新或未穩定，禁止登錄為通過")


if args.overlay_only:
    for backend in ("eten", "unifont"):
        bake = Path("workplace/r682") / backend
        for scale in (2, 3):
            dest = out / f"{backend}-{scale}"
            dest.mkdir()
            with (dest / "run.log").open("w") as log:
                proc = subprocess.Popen([str(binary), "-root", "workplace/run", "-autoplay",
                                         "-autoplay-stop", "286000000", "-scale", str(scale),
                                         "-hd-scale", str(scale), "-xlate", str(bake / "xlate.json"),
                                         "-flow", str(bake / "flow.json"), "-help", str(bake / "help.json"),
                                         "-hd-at", "272000000,286000000", "-draw-at", "272000000,286000000",
                                         "-hd-dir", str(dest)], stdout=log, stderr=subprocess.STDOUT)
                try:
                    wait_for(lambda: all((dest / f"hdat-{step}{ext}").is_file()
                                         for step in (272000000, 286000000)
                                         for ext in (".json", ".png", ".pal", ".a000", ".comp"))
                             and all((dest / f"drawat-{step}.png").is_file()
                                     for step in (272000000, 286000000)), proc)
                finally:
                    stop(proc)
            print(f"{backend} {scale}×：選單／防拷說明Draw及同狀態抽測完成", flush=True)
    raise SystemExit(0)


for backend in ("eten", "unifont"):
    bake = Path("workplace/r682") / backend
    for scale in (2, 3):
        dest = out / f"{backend}-{scale}"
        dest.mkdir()
        log = dest / "run.log"
        # 使用既有錄放路線；help-at先暫停，以真實鍵盤翻頁並抓圖。
        cmd = [str(binary), "-root", "workplace/run", "-autoplay",
               "-autoplay-stop", "330000000", "-scale", str(scale),
               "-xlate", str(bake / "xlate.json"), "-flow", str(bake / "flow.json"),
               "-help", str(bake / "help.json"), "-help-at", "286000000",
               "-help-for", "180000", "-hd-at", "300000000", "-hd-dir", str(dest),
               "-draw-at", "286000000"]
        with log.open("w") as stream:
            proc = subprocess.Popen(cmd, stdout=stream, stderr=subprocess.STDOUT)
            try:
                first = dest / "drawat-286000000.png"
                wait_for(first.is_file, proc)
                ids = subprocess.check_output(["xdotool", "search", "--name", "巫術7"], text=True).split()
                wid = ids[-1]
                subprocess.run(["xdotool", "windowfocus", "--sync", wid], check=True)
                layouts = json.loads(Path("workplace/r682/help-layout-receipts.json").read_text())
                total = len(next(r["pages"] for r in layouts if r["backend"] == backend
                                 and r["scale"] == scale and r["screen"] == "copyprot"))
                # 記錄全部顯示頁及一次循環；翻頁使用真實視窗輸入。
                previous = None
                for page in range(total + 1):
                    previous = help_frame(wid, dest / f"help-{page:02}.png", proc, previous)
                    subprocess.run(["xdotool", "key", "--delay", "100", "--window", wid, "Right"], check=True)
                subprocess.run(["xdotool", "key", "--delay", "100", "--window", wid, "Escape"], check=True)
                wait_for((dest / "hdat-300000000.json").is_file, proc)
                # F1再次開／關，驗真實保留鍵；不投遞原版按鍵。
                subprocess.run(["xdotool", "key", "--delay", "100", "--window", wid, "F1"], check=True)
                help_frame(wid, dest / "f1-reopen.png", proc)
                subprocess.run(["xdotool", "key", "--delay", "100", "--window", wid, "F1"], check=True)
            finally:
                stop(proc)
        print(f"{backend} {scale}×：真實F1／Esc及全部{total}頁收據完成", flush=True)

dest = out / "unifont-story"
dest.mkdir()
with (dest / "run.log").open("w") as log:
    bake = Path("workplace/r682/unifont")
    proc = subprocess.Popen([str(binary), "-root", "workplace/run", "-autoplay",
                             "-scale", "3", "-xlate", str(bake / "xlate.json"),
                             "-flow", str(bake / "flow.json"), "-help", str(bake / "help.json"),
                             "-draw-at", "363113713,395347218", "-hd-dir", str(dest)],
                            stdout=log, stderr=subprocess.STDOUT)
    try:
        wait_for(lambda: all((dest / f"drawat-{step}.png").is_file()
                             for step in (363113713, 395347218)), proc)
    finally:
        stop(proc)
print("Unifont故事Draw兩頁完成", flush=True)
