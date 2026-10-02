#!/usr/bin/env python3
# check_flow.py — docs/spec/006 驗收（流動文字）。
#
#   tools/py.sh tools/re/check_flow.py --zh <zh 輸出目錄> --en <en 輸出目錄> --log <en 模式 log> \
#        --scale 3 [--zh2 <2× 輸出目錄> --scale2 2] [--drop <缺段落的 zh 輸出目錄> --drop-para 15005] \
#        [--transition-from <頁第一個片段步數>]
#
# 獨立於被測程式：只讀 text/messages.json、text/story-trans.json、en 模式 log 與 hdat 輸出，
# 依 spec 006 R1–R4 用 Python 另寫一遍（與 bake_flow.py／cmd/wiz7/flow.go 除字形資料
# bake_font.glyph_scaled 外不共用判定程式）。
import argparse
import base64
import json
import re
import struct
import sys
import zlib
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from bake_font import glyph_scaled  # 只借字形資料

ap = argparse.ArgumentParser()
ap.add_argument("--zh", required=True)
ap.add_argument("--en", required=True)
ap.add_argument("--log", required=True)
ap.add_argument("--scale", type=int, default=3)
ap.add_argument("--zh2")
ap.add_argument("--scale2", type=int, default=2)
ap.add_argument("--drop")
ap.add_argument("--drop-para", type=int, default=15005)
ap.add_argument("--frag-dumps", default="", help="hdfrag 序號逗號分隔（新頁前兩個片段返回的當下，驗收 10）；第一個序號是新頁的首片段")
ap.add_argument("--flow-json", default="workplace/bake/flow.json")
args = ap.parse_args()

fails = []


def check(cond, msg):
    print(("通過 " if cond else "失敗 ") + msg)
    if not cond:
        fails.append(msg)


# ── R1 段落（獨立實作）────────────────────────────────────────────
recs = sorted(json.loads(Path("text/messages.json").read_text())["records"], key=lambda r: r["id"])
trans = json.loads(Path("text/story-trans.json").read_text())["items"]
BADSET = set('&%]^$#!"')


def ctl_free(s):
    return re.sub(r"[\x00-\x1f]", "", s)


def big5_ok(z):
    if not z:
        return False
    for ch in z:
        try:
            if len(ch.encode("big5")) != 2:
                return False
        except UnicodeEncodeError:
            return False
    return True


paras = []  # [{"words": [...], "lines": [(id, end, zh)]}]
cur = None
prev_id = None
for r in recs:
    i = r["id"]
    t = ctl_free(r["en"])
    lead = t.startswith("@")
    body = t[1:] if lead else t
    excluded = 2500 <= i <= 2574 or any(c in BADSET for c in body) or "@" in body or any(c.isdigit() for c in body)
    if excluded:
        cur, prev_id = None, i
        continue
    if cur is None or prev_id is None or i != prev_id + 1 or lead:
        cur = {"words": [], "lines": [], "raw": []}
        paras.append(cur)
    cur["raw"].append(i)
    ws = body.split()
    if ws:
        cur["words"] += ws
        z = ctl_free(trans.get(str(i), {}).get("zh", ""))
        cur["lines"].append((i, len(cur["words"]), z if big5_ok(z) else ""))
    prev_id = i
paras = [p for p in paras if p["lines"]]

# 驗收 1：flow.json 的統計與獨立算的相同
own = (len(paras), sum(len(p["lines"]) for p in paras), sum(1 for p in paras for l in p["lines"] if l[2]),
       sum(1 for p in paras if all(l[2] for l in p["lines"])))
fj = json.loads(Path(args.flow_json).read_text())
theirs = (len(fj["paras"]), sum(len(p["lines"]) for p in fj["paras"]),
          sum(1 for p in fj["paras"] for l in p["lines"] if l["zh"]),
          sum(1 for p in fj["paras"] if all(l["zh"] for l in p["lines"])))
check(own == theirs == (2906, 6307, 4923, 2168), f"驗收1 段落／行／可用行／整段可譯 獨立={own} flow.json={theirs} 規格=(2906, 6307, 4923, 2168)")


# ── R3 比對（獨立實作）────────────────────────────────────────────
def expect_zh(W):
    n = len(W)
    outs = set()
    for p in paras:
        w = p["words"]
        for a in range(len(w) - n + 1):
            if w[a:a + n] != W:
                continue
            done = [l for l in p["lines"] if a <= l[1] - 1 < a + n]
            zh = None if (not done or any(not l[2] for l in done)) else "".join(l[2] for l in done)
            outs.add((zh, p["lines"][0][0], a, tuple(l[0] for l in done)))
    zhs = {o[0] for o in outs}
    if not outs or len(zhs) != 1 or None in zhs:
        return None, outs
    return zhs.pop(), outs


# ── R4 排版（獨立實作）────────────────────────────────────────────
NOSTART = set("，。、！？；：）」』》…—")
NOEND = set("（「『《")


def layout(zh, N):
    units, k = [], 0
    while k < len(zh):
        if zh[k] in "…—" and k + 1 < len(zh) and zh[k + 1] == zh[k]:
            units.append(zh[k:k + 2]); k += 2
        else:
            units.append(zh[k]); k += 1
    rows, i = [], 0
    while i < len(units):
        j, w = i, 0
        while j < len(units) and w + len(units[j]) <= N:
            w += len(units[j]); j += 1
        if j == i:
            j = i + 1
        for _ in range(3):
            if j - i > 1 and j < len(units) and (units[j][0] in NOSTART or units[j - 1][0] in NOEND):
                j -= 1
            else:
                break
        rows.append("".join(units[i:j]))
        i = j
    return rows


def gs_of(k):
    """字形放大（spec 006 R4，#41）：hdFrame 端 hdf＝1，倍率 3 以上放大 2 倍。"""
    return 2 if k >= 3 else 1


def cols(k):
    return max(1, 184 * k // (17 * gs_of(k)))


# ── log 分頁（R2）─────────────────────────────────────────────────
pages = []
last_y = None
for line in open(args.log, encoding="utf-8"):
    m = re.search(r'流文 片段 (\d+) ret=(\S+) font=(\d+) y0=(\d+) x0=(\d+) "(.*)"$', line)
    if not m:
        continue
    if m.group(2) != '0110:bcfc':
        continue
    step, y, s = int(m.group(1)), int(m.group(4)), m.group(6)
    if last_y is None or y != last_y + 12 or len(pages[-1]["frags"]) >= 4:
        pages.append({"frags": [], "steps": []})
    pages[-1]["frags"].append(s)
    pages[-1]["steps"].append(step)
    last_y = y
check(len(pages) == 6 and [len(p["frags"]) for p in pages] == [4, 4, 1, 4, 4, 3],
      f"驗收4 en log 分頁 {[len(p['frags']) for p in pages]}（期望 6 頁 4,4,1,4,4,3）")
for p in pages:
    p["words"] = " ".join(p["frags"]).split()
    p["zh"], _ = expect_zh(p["words"])


# ── 載入 hdat ─────────────────────────────────────────────────────
def load(d):
    out = {}
    for f in sorted(Path(d).glob("hd*.json")):
        if f.stem.startswith('hdat-') and int(f.stem[5:]) < 1483000000:
            continue  # 故事收據另由 check_story.py 驗收。
        if f.stem == 'hdfrag-4':
            continue
        j = json.loads(f.read_text())
        for e in j.get("Shown") or []:
            e["Before"] = base64.b64decode(e["Before"])
            e["After"] = base64.b64decode(e["After"])
        j["_a000"] = f.with_suffix(".a000").read_bytes()
        j["_comp"] = f.with_suffix(".comp").read_bytes()
        j["_png"] = f.with_suffix(".png")
        out[f.stem] = j
    return out


def read_png(path):
    d = path.read_bytes()
    i, idat, w = 8, b"", 0
    while i < len(d):
        n = struct.unpack(">I", d[i:i + 4])[0]
        t, c = d[i + 4:i + 8], d[i + 8:i + 8 + n]
        if t == b"IHDR":
            w, h, bd, ct = struct.unpack(">IIBB", c[:10])
        elif t == b"IDAT":
            idat += c
        i += 12 + n
    assert bd == 8 and ct in (2, 6)
    bpp = 3 if ct == 2 else 4
    raw, stride, rows, prev, p = zlib.decompress(idat), w * bpp, [], bytearray(w * bpp), 0
    for _ in range(h):
        ft = raw[p]; p += 1
        line = bytearray(raw[p:p + stride]); p += stride
        for x in range(stride):
            a = line[x - bpp] if x >= bpp else 0
            b = prev[x]
            c = prev[x - bpp] if x >= bpp else 0
            if ft == 1: line[x] = (line[x] + a) & 255
            elif ft == 2: line[x] = (line[x] + b) & 255
            elif ft == 3: line[x] = (line[x] + (a + b) // 2) & 255
            elif ft == 4:
                pp = a + b - c
                pa, pb, pc = abs(pp - a), abs(pp - b), abs(pp - c)
                line[x] = (line[x] + (a if pa <= pb and pa <= pc else b if pb <= pc else c)) & 255
        rows.append(bytes(line)); prev = line
    return w, h, bpp, rows


glyph_cache = {}


def glyph(ch):
    if ch not in glyph_cache:
        glyph_cache[ch] = glyph_scaled(ch.encode("big5"))
    return glyph_cache[ch]


BAND = (68, 251, 146, 194)


def check_dumps(zh, en, k, tag):
    seen = {}  # 頁索引 → 出現的 dump 數
    band_bad = 0
    for name, j in zh.items():
        for fl in j.get("Flow") or []:
            idx = next((i for i, p in enumerate(pages) if p["frags"] == fl["Fragments"]), None)
            check(idx is not None, f"{tag} {name} Flow 頁的片段對得上 en log 的某一頁 {fl['Fragments'][:1]}")
            if idx is None:
                continue
            seen[idx] = seen.get(idx, 0) + 1
            exp = pages[idx]["zh"]
            check(fl["Zh"] == exp, f"驗收5 {tag} {name} 頁{idx + 1} Zh 獨立={exp!r} 程式={fl['Zh']!r}")
            check("".join(r["Text"] for r in fl["Rows"]) == fl["Zh"], f"驗收5 {tag} {name} 頁{idx + 1} Rows 相接等於 Zh")
            # 驗收 6：排版
            rows = layout(exp, cols(k))
            check([r["Text"] for r in fl["Rows"]] == rows and len(rows) <= 4, f"驗收6 {tag} {name} 頁{idx + 1} 列文字 獨立={rows}")
            gs = gs_of(k)
            check(fl["Scale"] == k and fl.get("GS") == gs, f"驗收6 {tag} {name} Scale={fl['Scale']} GS={fl.get('GS')}（期望 {gs}）")
            for ri, (r, text) in enumerate(zip(fl["Rows"], rows)):
                x0 = 160 * k - len(text) * 17 * gs // 2
                check(r["Y0"] == (146 + 12 * ri) * k and abs(r["X0"] - x0) <= gs + 1, f"驗收6 {tag} {name} 頁{idx + 1} 列{ri} 位置 ({r['X0']},{r['Y0']}) 期望 ({x0},{(146 + 12 * ri) * k})")
                if len(text) > 1:
                    check(text[0] not in NOSTART or ri == 0, f"驗收6 {tag} {name} 頁{idx + 1} 列{ri} 避頭")
                check(text[-1] not in NOEND or ri == len(rows) - 1, f"驗收6 {tag} {name} 頁{idx + 1} 列{ri} 避尾")
            # 驗收 7：像素（PNG）
            w, h, bpp, prow = read_png(j["_png"])
            miss = 0
            for r in fl["Rows"]:
                for ci, ch in enumerate(r["Text"]):
                    for row, bits in enumerate(glyph(ch)):
                        for bit in range(15):
                            if (bits >> (14 - bit)) & 1:
                                for dy in range(gs):
                                    for dx in range(gs):
                                        xx = r["X0"] + (ci * 17 + bit) * gs + dx
                                        px = prow[r["Y0"] + row * gs + dy][xx * bpp:xx * bpp + 3]
                                        if list(px) != fl["FgRGB"]:
                                            miss += 1
            check(miss == 0, f"驗收7 {tag} {name} 頁{idx + 1} 字形像素不符數 {miss}")
            # 不變式：有效頁的英文仍在 A000（底圖像素等於 After、至少一半筆畫像素與 Before 不同）
            mx0, my0 = j["Mouse"]
            base_bad = stroke = present = 0
            for e in j["Shown"]:
                if e["S"] not in fl["Fragments"]:
                    continue
                for kk, (b, a) in enumerate(zip(e["Before"], e["After"])):
                    x, y = e["BX"] + kk % e["BW"], e["BY"] + kk // e["BW"]
                    if mx0 <= x < mx0 + 16 and my0 <= y < my0 + 16:
                        continue
                    v = en[name]["_a000"][y * 320 + x]
                    if b == a:
                        base_bad += v != a
                    else:
                        stroke += 1
                        present += v != b
            check(base_bad == 0 and stroke > 0 and present * 2 >= stroke,
                  f"不變式 {tag} {name} 頁{idx + 1} 英文仍在 A000（底圖不符 {base_bad}、筆畫 {stroke}、仍在 {present}）")
            # 驗收 8：英文不殘留（筆畫像素合成後等於 Before）
            enj = en[name]
            comp, a000 = j["_comp"], enj["_a000"]
            mx, my = j["Mouse"]
            resid = total = 0
            for e in j["Shown"]:
                if e["S"] not in fl["Fragments"]:
                    continue
                for kk, (b, a) in enumerate(zip(e["Before"], e["After"])):
                    if b == a:
                        continue
                    x, y = e["BX"] + kk % e["BW"], e["BY"] + kk // e["BW"]
                    if mx <= x < mx + 16 and my <= y < my + 16:
                        continue
                    total += 1
                    if comp[y * 320 + x] != b:
                        resid += 1
            check(resid == 0 and total > 0, f"驗收8 {tag} {name} 頁{idx + 1} 英文殘留 {resid}／{total}")
    # 中英合成差異只在文字帶內
    for name, j in zh.items():
        enj = en[name]
        bad = sum(1 for i in range(64000) if j["_comp"][i] != enj["_comp"][i] and not (BAND[0] <= i % 320 <= BAND[1] and BAND[2] <= i // 320 <= BAND[3]))
        band_bad += bad
    check(band_bad == 0, f"驗收8 {tag} 中英合成的差異全在文字帶內（帶外差異像素 {band_bad}）")
    return seen


zh, en = load(args.zh), load(args.en)
check(set(zh) == set(en) and len(zh) > 0, f"驗收2 zh 與 en 有相同的 {len(zh)} 個 dump")
# 驗收 3：A/B
ab_bad = [n for n in zh if any(zh[n][k_] != en[n][k_] for k_ in ("Step", "Mem", "A000", "Pal"))]
check(not ab_bad, f"驗收3 原版記憶體／A000／調色盤雜湊 zh＝en（{len(zh)} 個 dump，不符 {len(ab_bad)}）")
seen = check_dumps(zh, en, args.scale, f"{args.scale}×")
check(sorted(seen) == list(range(len(pages))), f"驗收4 六頁都至少在一個 dump 有效（有效頁 {sorted(i + 1 for i in seen)}）")

# 驗收 10：換頁過渡（hdfrag：新頁前兩個片段返回的當下）
if args.frag_dumps:
    seqs = [int(x) for x in args.frag_dumps.split(",")]
    for n in seqs:
        name = f"hdfrag-{n}"
        check(name in zh, f"驗收10 有 {name} 的輸出")
        if name not in zh:
            continue
        flows = zh[name].get("Flow") or []
        new_page = [fl for fl in flows if fl["Fragments"] == pages[3]["frags"]]
        check(not new_page, f"驗收10 {name} 新頁（第 4 頁）在片段未印完時不畫（Flow 頁 {[fl['Fragments'][0][:12] for fl in flows]}）")

# 驗收 9：負向
if args.drop:
    dz = load(args.drop)
    p_idx = next(i for i, p in enumerate(pages) if expect_zh(p["words"])[1] and next(iter(expect_zh(p["words"])[1]))[1] == args.drop_para)
    seen_d = {}
    for name, j in dz.items():
        for fl in j.get("Flow") or []:
            idx = next((i for i, p in enumerate(pages) if p["frags"] == fl["Fragments"]), None)
            seen_d[idx] = seen_d.get(idx, 0) + 1
    check(p_idx not in seen_d, f"驗收9 缺段落 {args.drop_para} 的頁{p_idx + 1}沒有 Flow 項")
    check(sorted(seen_d) == [i for i in range(len(pages)) if i != p_idx], f"驗收9 其餘頁不受影響（有效頁 {sorted(i + 1 for i in seen_d if i is not None)}）")
    # 缺頁的合成帶內等於 A000
    st0, st1 = pages[p_idx]["steps"][-1] + 5_000_000, pages[p_idx]["steps"][-1] + 8_000_000
    for name, j in dz.items():
        if st0 <= j["Step"] <= st1:
            diff = sum(1 for i in range(64000) if j["_comp"][i] != j["_a000"][i]
                       and BAND[0] <= i % 320 <= BAND[1] and BAND[2] <= i // 320 <= BAND[3])
            check(diff == 0, f"驗收9 {name} Step={j['Step']} 缺段落頁在文字帶內 comp 等於 A000（差異 {diff}）")

if args.zh2:
    zh2 = load(args.zh2)
    seen2 = check_dumps(zh2, en, args.scale2, f"{args.scale2}×")
    check(sorted(seen2) == list(range(len(pages))), f"驗收11 {args.scale2}× 六頁都至少在一個 dump 有效（{sorted(i + 1 for i in seen2)}）")

print(f"\n失敗 {len(fails)} 項")
sys.exit(1 if fails else 0)
