#!/usr/bin/env python3
"""spec 006 故事驗收：以原版返回觀測、取訊息與字形獨立計算期望值。

Docker 內從專案根目錄執行，字型掛載方式同 check_flow.py：
  python3 tools/re/check_story.py --en workplace/r674/en --zh workplace/r674/zh \
    --zh2 workplace/r674/zh2 --drop workplace/r674/drop --log workplace/r674/en.log

只驗故事 31 頁、124 個請求及 hdfrag-4；容許目錄另含旁白收據。
不讀 flow.json、不呼叫受測的比對／排版函式。共享來源只有字形抽取。
記憶體與色盤依 runner 的雜湊收據比較，A000 另比較原始 bytes。
此工具不代替互動 Draw 目視、旁白回歸、失效分支單測或容器清理。
"""
import argparse
import hashlib
import json
import re
import struct
import sys
import zlib
from collections import Counter
from functools import lru_cache
from pathlib import Path

from bake_font import glyph_scaled

SITE = 0xBE0C
STORY_ROUTE_END = 1435000000
BAND = set(y * 320 + x for y in range(146, 195) for x in range(12, 308))
NOSTART = set("，。、！？；：）」』》…—")
NOEND = set("（「『《")
TERMINATOR = re.compile(r"\s*(?:&[%#|\[]*)+\]$")


class Audit:
    def __init__(self, verbose=False):
        self.total = 0
        self.failed = 0
        self.verbose = verbose
        self.color_missing = set()

    def check(self, ok, label):
        self.total += 1
        if not ok:
            self.failed += 1
        if not ok or self.verbose:
            print(("通過 " if ok else "失敗 ") + label)
        return ok


def clean(text):
    text = re.sub(r"[\x00-\x1f]", "", text)
    if text.startswith(("@", "#")):
        text = text[1:]
    return TERMINATOR.sub("", text)


@lru_cache(maxsize=None)
def glyph(ch):
    encoded = ch.encode("big5")
    if len(encoded) != 2:
        raise ValueError("譯文含非雙位元組字")
    rows = tuple(int(n) for n in glyph_scaled(encoded))
    if len(rows) != 15 or any(not 0 <= n < 32768 for n in rows):
        raise ValueError("字形不是 15×15")
    return rows


def layout(text, width):
    units = re.findall(r"……|——|.", text)
    result, start = [], 0
    while start < len(units):
        end, count = start, 0
        while end < len(units) and count + len(units[end]) <= width:
            count += len(units[end])
            end += 1
        if end == start:
            end += 1
        for _ in range(3):
            if (end - start > 1 and end < len(units)
                    and (units[end][0] in NOSTART or units[end - 1][0] in NOEND)):
                end -= 1
            else:
                break
        result.append("".join(units[start:end]))
        start = end
    return result


def read_png(path):
    data = path.read_bytes()
    if data[:8] != b"\x89PNG\r\n\x1a\n":
        raise ValueError("PNG 簽名不符")
    offset, compressed, header = 8, bytearray(), None
    while offset < len(data):
        n = struct.unpack(">I", data[offset:offset + 4])[0]
        kind = data[offset + 4:offset + 8]
        chunk = data[offset + 8:offset + 8 + n]
        crc = struct.unpack(">I", data[offset + 8 + n:offset + 12 + n])[0]
        if zlib.crc32(kind + chunk) & 0xffffffff != crc:
            raise ValueError("PNG 區塊 CRC 不符")
        if kind == b"IHDR":
            header = struct.unpack(">IIBBBBB", chunk)
        elif kind == b"IDAT":
            compressed.extend(chunk)
        elif kind == b"IEND":
            break
        offset += n + 12
    w, h, depth, color, compression, filtering, interlace = header
    if depth != 8 or color not in (2, 6) or (compression, filtering, interlace) != (0, 0, 0):
        raise ValueError("只接受非交錯 8-bit RGB／RGBA PNG")
    bpp = 3 if color == 2 else 4
    raw = zlib.decompress(compressed)
    stride = w * bpp
    if len(raw) != h * (stride + 1):
        raise ValueError("PNG 解壓長度不符")
    rows, prev, offset = [], bytearray(stride), 0
    for _ in range(h):
        ft = raw[offset]
        if ft > 4:
            raise ValueError("PNG filter 不符")
        line = bytearray(raw[offset + 1:offset + 1 + stride])
        offset += stride + 1
        if ft == 0:
            rows.append(bytes(line))
            prev = line
            continue
        for i in range(stride):
            a = line[i - bpp] if i >= bpp else 0
            b, c = prev[i], prev[i - bpp] if i >= bpp else 0
            if ft == 1:
                v = a
            elif ft == 2:
                v = b
            elif ft == 3:
                v = (a + b) // 2
            elif ft == 4:
                p = a + b - c
                distances = (abs(p - a), abs(p - b), abs(p - c))
                v = (a, b, c)[distances.index(min(distances))]
            else:
                v = 0
            line[i] = (line[i] + v) & 255
        rows.append(bytes(line))
        prev = line
    return w, h, bpp, rows


def cursor(meta):
    mx, my = meta["Mouse"]
    return set(y * 320 + x for y in range(max(0, my), min(200, my + 16))
               for x in range(max(0, mx), min(320, mx + 16)))


def original_pages(args, audit):
    log = Path(args.msg_log).read_text()
    reads = [(int(s), int(i)) for s, i in re.findall(r"訊息 (\d+) id=(\d+).*ret=bfeb", log)]
    entries = [(int(s), json.loads(text), int(y)) for s, text, y in re.findall(
        r'印字 (\d+) ("(?:[^"\\]|\\.)*") font=2 ret=0110:be0c .*?xy0=\(\d+,(\d+)\)', log)]
    traces = json.loads(Path(args.traces).read_text())
    if not audit.check(len(entries) == len(traces) == 71 and len(reads) == 68,
                       "原版來源必須為 71 個片段、71 個返回、68 筆取訊息"):
        raise ValueError("原版來源不完整，不能降低頁數驗收")
    messages = {r["id"]: r["en"] for r in json.loads(Path(args.messages).read_text())["records"]
                if 31500 <= r["id"] <= 31899}
    translations = json.loads(Path(args.translations).read_text())["items"]
    pages = []
    for i, ((step, text, y), t) in enumerate(zip(entries, traces)):
        if y == 146:
            pages.append({"fragments": [], "entries": [], "traces": []})
        if not pages:
            raise ValueError("原版首列不是 y＝146")
        page = pages[-1]
        audit.check(t["Page"] == len(pages) and t["Y"] == y
                    and y == 146 + 12 * len(page["fragments"]), f"原版片段 {i + 1} 頁序與列距")
        audit.check(step < t["Step"] and (i + 1 == len(entries) or t["Step"] < entries[i + 1][0]),
                    f"原版片段 {i + 1} 真正返回在自身入口與下一入口之間")
        before, after = Path(t["Before"]).read_bytes(), Path(t["After"]).read_bytes()
        if len(before) != 64000 or len(after) != 64000:
            raise ValueError("原版離屏快照長度不符")
        changed = set(j for j, (b, a) in enumerate(zip(before, after)) if b != a)
        if not changed:
            raise ValueError("原版片段沒有差異")
        x0, x1 = min(j % 320 for j in changed), max(j % 320 for j in changed)
        y0, y1 = min(j // 320 for j in changed), max(j // 320 for j in changed)
        box = set(yy * 320 + xx for yy in range(y0, y1 + 1) for xx in range(x0, x1 + 1))
        page["fragments"].append(text)
        page["entries"].append(step)
        page["traces"].append({"return": t["Step"], "before": before, "after": after,
                               "changed": changed, "background": box - changed})
    if not audit.check(len(pages) == 31 and Counter(len(p["fragments"]) for p in pages)
                       == {1: 8, 2: 10, 3: 9, 4: 4}, "原版必須為 31 頁，片段分布 8／10／9／4"):
        raise ValueError("原版頁結構不符")
    lower = 0
    for n, p in enumerate(pages, 1):
        ids = [i for s, i in reads if lower < s < p["entries"][-1]
               and i in messages and clean(messages[i]).split()]
        words = [w for i in ids for w in clean(messages[i]).split()]
        p["words"] = " ".join(p["fragments"]).split()
        audit.check(words == p["words"], f"原版第 {n} 頁時間區間來源與完整英文相等")
        z = "".join(clean(translations.get(str(i), {}).get("zh", "")) for i in ids)
        if not ids or any(not clean(translations.get(str(i), {}).get("zh", "")) for i in ids):
            raise ValueError("正式訊息譯文缺漏")
        for ch in z:
            glyph(ch)
        p.update(ids=ids, zh=z, lower=lower, last_return=p["traces"][-1]["return"],
                 strokes=set().union(*(t["changed"] for t in p["traces"])))
        lower = p["last_return"]
    audit.check(sum(len(p["ids"]) for p in pages) == 58, "原版31頁共58則可見來源訊息")
    return pages, entries, reads


def load_dump(directory, name):
    base = Path(directory) / name
    meta = json.loads(base.with_suffix(".json").read_text())
    a000, comp = base.with_suffix(".a000").read_bytes(), base.with_suffix(".comp").read_bytes()
    if len(a000) != 64000 or len(comp) != 64000:
        raise ValueError(f"{base} 原版／合成色號長度不符")
    return {"meta": meta, "a000": a000, "comp": comp, "png": base.with_suffix(".png")}


def flow_page(fl, pages):
    if fl.get("Site") != SITE:
        return None
    matches = [i for i, p in enumerate(pages) if p["fragments"] == fl.get("Fragments")]
    return matches[0] if len(matches) == 1 else None


def inspect_dump(audit, name, dump, en, pages, scale, tag, blocked=None):
    meta = dump["meta"]
    audit.check(all(meta.get(k) == en["meta"].get(k) for k in ("Step", "Mem", "A000", "Pal"))
                and all(k in meta and k in en["meta"] for k in ("Step", "Mem", "A000", "Pal"))
                and dump["a000"] == en["a000"], f"{tag} {name} 原版記憶體／A000／色盤／步數相同")
    audit.check(meta["Mouse"] == en["meta"]["Mouse"], f"{tag} {name} 游標狀態相同")
    cur = cursor(meta)
    flows = meta.get("Flow") or []
    seen, expected = set(), bytearray(en["comp"])
    png = None
    for fl in flows:
        idx = flow_page(fl, pages)
        if not audit.check(idx is not None, f"{tag} {name} Flow 是獨立來源中的完整故事頁"):
            continue
        if not audit.check(idx != blocked, f"{tag} {name} 缺訊息頁不可覆繪"):
            continue
        audit.check(idx not in seen, f"{tag} {name} 不重複登記同頁")
        seen.add(idx)
        p = pages[idx]
        audit.check(fl.get("Zh") == p["zh"] and fl.get("Lines") == p["ids"]
                    and fl.get("Para") == p["ids"][0] and fl.get("A") == 0
                    and fl.get("N") == len(p["words"]), f"{tag} {name} 第 {idx + 1} 頁譯文／來源 id／字詞數")
        audit.check(meta["Step"] >= p["last_return"] + 200000,
                    f"{tag} {name} 第 {idx + 1} 頁不是片段返回中間態")
        original_match = True
        for t in p["traces"]:
            strokes = t["changed"] - cur
            original_match &= bool(strokes) and all(dump["a000"][i] == t["after"][i]
                                                   for i in t["background"] - cur)
            original_match &= 2 * sum(dump["a000"][i] != t["before"][i] for i in strokes) >= len(strokes)
            for i in strokes:
                expected[i] = t["before"][i]
        audit.check(original_match, f"{tag} {name} 第 {idx + 1} 頁英文仍在原版 A000")
        audit.check(all(dump["a000"][i] == p["traces"][-1]["after"][i]
                        for i in BAND - p["strokes"] - cur), f"{tag} {name} 第 {idx + 1} 頁背景帶靜態相等")
        audit.check(all(dump["comp"][i] == t["before"][i]
                        for t in p["traces"] for i in t["changed"] - cur),
                    f"{tag} {name} 第 {idx + 1} 頁英文筆畫全部還原")
        gs = 2 if scale >= 3 else 1
        rows = layout(p["zh"], max(1, 296 * scale // (17 * gs)))
        positions = [{"Text": text, "X0": 160 * scale - len(text) * 17 * gs // 2,
                      "Y0": (146 + 12 * ri) * scale} for ri, text in enumerate(rows)]
        audit.check(fl.get("Rows") == positions and len(rows) <= 4 and fl.get("Scale") == scale
                    and fl.get("GS") == gs, f"{tag} {name} 第 {idx + 1} 頁獨立排版／位置／倍率")
        if png is None:
            png = read_png(dump["png"])
        w, h, bpp, rgb = png
        audit.check((w, h) == (320 * scale, 200 * scale), f"{tag} {name} PNG 畫布尺寸")
        fg = fl.get("FgRGB")
        valid_fg = isinstance(fg, list) and len(fg) == 3 and all(isinstance(c, int) and 0 <= c <= 255 for c in fg)
        audit.check(valid_fg, f"{tag} {name} 第 {idx + 1} 頁前景色格式")
        if "original_fg" not in en:
            if "decoded_png" not in en:
                en["decoded_png"] = read_png(en["png"])
            ew, eh, ebpp, ergb = en["decoded_png"]
            if (ew, eh) != (960, 600):
                raise ValueError("en 前景色來源必須是3×原版畫布")
            colors = {tuple(ergb[(i // 320) * 3][(i % 320) * 3 * ebpp:
                                                            (i % 320) * 3 * ebpp + 3])
                      for i in BAND - cursor(en["meta"])
                      if en["a000"][i] == en["comp"][i] == 0x55}
            audit.check(len(colors) <= 1, f"en {name} 原版0x55像素有唯一RGB")
            en["original_fg"] = list(next(iter(colors))) if len(colors) == 1 else None
        if en["original_fg"] is not None:
            audit.check(fg == en["original_fg"], f"{tag} {name} 前景色等於en原版0x55獨立RGB")
        else:
            audit.color_missing.add(name)
        miss = 0
        if valid_fg and (w, h) == (320 * scale, 200 * scale):
            for row in positions:
                for ci, ch in enumerate(row["Text"]):
                    for gy, bits in enumerate(glyph(ch)):
                        for gx in range(15):
                            if not bits & (1 << (14 - gx)):
                                continue
                            for dy in range(gs):
                                for dx in range(gs):
                                    x = row["X0"] + (ci * 17 + gx) * gs + dx
                                    y = row["Y0"] + gy * gs + dy
                                    miss += not (0 <= x < w and 0 <= y < h
                                                 and list(rgb[y][x * bpp:x * bpp + 3]) == fg)
        audit.check(valid_fg and miss == 0, f"{tag} {name} 第 {idx + 1} 頁獨立字形像素（不符 {miss}）")
    audit.check(dump["comp"] == bytes(expected), f"{tag} {name} 合成色號只還原批准英文筆畫")
    audit.check(all(dump["comp"][i] == en["comp"][i] for i in range(64000) if i not in BAND),
                f"{tag} {name} 帶外色號沒有差異")
    # 同倍率的整張 PNG 差異必須在批准帶內；2× 沒有獨立 en 2× PNG，不冒稱整張對照。
    if scale == 3:
        if png is None:
            png = read_png(dump["png"])
        if "decoded_png" not in en:
            en["decoded_png"] = read_png(en["png"])
        epng = en["decoded_png"]
        w, h, bpp, rgb = png
        ew, eh, ebpp, ergb = epng
        if audit.check((w, h) == (ew, eh) == (960, 600), f"{tag} {name} 中英 PNG 尺寸相同"):
            bad = sum(rgb[y][x * bpp:x * bpp + 3] != ergb[y][x * ebpp:x * ebpp + 3]
                      for y in range(h) for x in range(w)
                      if not (12 * 3 <= x < 308 * 3 and 146 * 3 <= y < 195 * 3))
            audit.check(bad == 0, f"{tag} {name} PNG 差異只在批准文字帶（帶外 {bad}）")
    return seen


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    for mode in ("en", "zh", "zh2", "drop"):
        ap.add_argument("--" + mode, required=True)
    ap.add_argument("--log", required=True, help="純英文 runner log，核對故事印字與取訊息 hook")
    ap.add_argument("--msg-log", default="workplace/bake/r667-story-msgs.txt")
    ap.add_argument("--traces", default="workplace/r673/timing/traces.json")
    ap.add_argument("--messages", default="text/messages.json")
    ap.add_argument("--translations", default="text/story-trans.json")
    ap.add_argument("--verbose", action="store_true")
    args = ap.parse_args()
    audit = Audit(args.verbose)
    try:
        pages, entries, reads = original_pages(args, audit)
        runtime = Path(args.log).read_text()
        actual = [(int(s), json.loads(text), int(y)) for s, y, text in re.findall(
            r'流文 片段 (\d+) ret=0110:be0c font=2 y0=(\d+) x0=\d+ ("(?:[^"\\]|\\.)*")', runtime)
            if int(s) <= STORY_ROUTE_END]
        actual_reads = [(int(s), int(i)) for s, i in re.findall(r"流文 取 (\d+) id=(\d+)", runtime)
                        if int(s) <= STORY_ROUTE_END]
        audit.check(all(s <= STORY_ROUTE_END for s, _ in reads)
                    and all(s <= STORY_ROUTE_END for s, _, _ in entries), "獨立探針在spec故事1435M路線範圍內")
        audit.check(actual == entries, "純英文 runner 71 個故事片段與獨立探針逐一相同")
        audit.check(actual_reads == reads, "純英文 runner 68 筆故事取訊息與獨立探針逐一相同")
        drop_idx = next(i for i, p in enumerate(pages) if 31654 in p["ids"])
        if drop_idx != 2:
            raise ValueError("缺訊息負向來源不是第 3 頁")
        expected_names = {f"hdat-{p['entries'][-1] + d * 1000000}" for p in pages for d in (1, 3, 6, 8)}
        expected_names.add("hdfrag-4")
        selected = {}
        for mode in ("en", "zh", "zh2", "drop"):
            available = {p.stem for p in Path(getattr(args, mode)).glob("hd*.json")}
            audit.check(expected_names <= available, f"{mode} 必須含124張故事取樣與 hdfrag-4")
            missing = expected_names - available
            if missing:
                raise ValueError(f"{mode} 缺少 {len(missing)} 個必要收據")
            selected[mode] = {n: load_dump(getattr(args, mode), n) for n in expected_names}
        seen = {m: set() for m in ("zh", "zh2", "drop")}
        for name in sorted(expected_names):
            en = selected["en"][name]
            audit.check(not en["meta"].get("Flow"), f"en {name} 沒有流動文字覆繪")
            for mode, scale in (("zh", 3), ("zh2", 2), ("drop", 3)):
                seen[mode] |= inspect_dump(audit, name, selected[mode][name], en, pages, scale, mode,
                                           drop_idx if mode == "drop" else None)
            if name.startswith("hdat-"):
                request = int(name[5:])
                p_idx = next(i for i, p in enumerate(pages)
                             if request in {p["entries"][-1] + d * 1000000 for d in (1, 3, 6, 8)})
                meta = selected["drop"][name]["meta"]
                if p_idx == drop_idx:
                    audit.check(not meta.get("Flow") and all(selected["drop"][name]["comp"][i]
                                == selected["drop"][name]["a000"][i] for i in BAND),
                                f"drop {name} 第3頁全帶回退原版英文")
                else:
                    audit.check(meta.get("Flow") == selected["zh"][name]["meta"].get("Flow")
                                and selected["drop"][name]["comp"] == selected["zh"][name]["comp"],
                                f"drop {name} 其餘頁與正式 zh 行為相同")
                delay = (request - pages[p_idx]["entries"][-1]) // 1000000
                if delay in (6, 8):
                    for mode in ("zh", "zh2", "drop"):
                        fs = selected[mode][name]["meta"].get("Flow") or []
                        current = [f for f in fs if flow_page(f, pages) == p_idx]
                        audit.check(len(current) == (0 if mode == "drop" and p_idx == drop_idx else 1),
                                    f"{mode} {name} {delay}M穩定樣本必須有完整當前頁")
            en.pop("decoded_png", None)
        for mode in ("zh", "zh2"):
            audit.check(seen[mode] == set(range(31)), f"{mode} 全部31頁至少一次有效（實際 {len(seen[mode])}）")
        audit.check(seen["drop"] == set(range(31)) - {drop_idx}, "drop 僅第3頁不覆繪，其他30頁全部有效")
        for mode in ("zh", "zh2", "drop"):
            fs = selected[mode]["hdfrag-4"]["meta"].get("Flow") or []
            audit.check(not any(f.get("Fragments", [])[:1] == pages[1]["fragments"][:1] for f in fs),
                        f"{mode} hdfrag-4 第2頁第一片段返回當下不畫新頁")
        for source in (args.msg_log, args.traces, args.messages, args.translations):
            print(f"來源 {source} SHA-256={hashlib.sha256(Path(source).read_bytes()).hexdigest()}")
    except (OSError, ValueError, KeyError, IndexError, TypeError, StopIteration, struct.error, zlib.error) as exc:
        audit.check(False, f"驗收輸入或格式受阻：{type(exc).__name__}（不輸出輸入文字）")
    if audit.color_missing:
        print("顏色證據限制：下列原版樣本沒有可讀的0x55像素，僅檢查字形與記錄色相符，未獨立驗證該樣本RGB："
              + ",".join(sorted(audit.color_missing)))
    print(f"故事驗收：{audit.total} 項，失敗 {audit.failed} 項；要求31頁、124取樣、兩倍率與缺訊息回退")
    return 1 if audit.failed else 0


if __name__ == "__main__":
    sys.exit(main())
