#!/usr/bin/env bash
# automap-check.sh — spec 005 驗收 1、2、4（自動地圖）。
#
#   tools/automap-check.sh
#
# 1 逐字重播 docs/re/016 walk3 的時序，-hd-at 1660M 的地圖狀態須為：走過 {(56,91),(56,92)}、
#   目前 (56,92) 面向 1、區域 11（期望值取自 workplace/bake/r666-walk3.txt）；log 無 F2（3c）送鍵
# 2 同跑法有無 -map-at 的原版記憶體／A000／調色盤雜湊相同
# 3 Xvfb 真實視窗按 F2：面板開啟、F2 未送原版，截圖 x.png
# 產物在 workplace/r668/（gitignore）。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT=workplace/r668
cd "$ROOT"; rm -rf "$OUT"; mkdir -p "$OUT/a" "$OUT/b"
CLK=1560000000:160:170,1569000000:160:170,1578000000:160:170,1587000000:160:170,1596000000:160:170,1605000000:160:170
CLK=$CLK,1614000000:160:184,1623000000:160:162,1632000000:180:166,1641000000:180:166,1650000000:160:162
COMMON="-hdonly -autoplay -autoplay-stop 1435000000 -root workplace/run -hd-at 1660000000 -clicks-at $CLK"
WIZ7_FE_TIMEOUT=20m WIZ7_FE_SH="go build -o /tmp/wiz7 ./cmd/wiz7 || exit 3
timeout 200 /tmp/wiz7 $COMMON -hd-dir $OUT/a -map-at 1600000000 > $OUT/a.log 2>&1
timeout 200 /tmp/wiz7 $COMMON -hd-dir $OUT/b > $OUT/b.log 2>&1
# 3 Xvfb：真實視窗，最後一下測試點擊後用 xdotool 按 F2 開面板並截圖
/tmp/wiz7 -autoplay -autoplay-stop 1435000000 -root workplace/run -scale 3 -clicks-at $CLK > $OUT/x.log 2>&1 &
P=\$!
for _ in \$(seq 1 600); do grep -q '測試點擊 (160,162) steps=165' $OUT/x.log && break; sleep 0.5; done
sleep 2
xdotool key --window \$(xdotool search --name '巫術7' | head -1) F2
sleep 2
import -window root $OUT/x.png 2>/dev/null
kill \$P; wait \$P 2>/dev/null
true" tools/fe.sh >/dev/null 2>&1 || true
tools/py.sh - "$OUT" <<'PY'
import json, sys
out = sys.argv[1]
a = json.load(open(f"{out}/a/hdat-1660000000.json"))
b = json.load(open(f"{out}/b/hdat-1660000000.json"))
m = a["Map"]
fails = []
def check(ok, msg):
    print(("通過 " if ok else "失敗 ") + msg)
    if not ok: fails.append(msg)
vis = sorted(tuple(v) for v in (m["Visited"] or []))
check(vis == [(11, 56, 91), (11, 56, 92)], f"驗收1 走過的格 {vis}")
check((m["Area"], m["X"], m["Y"], m["Face"]) == (11, 56, 92, 1), f"驗收1 目前位置 區域={m['Area']} ({m['X']},{m['Y']}) 面向={m['Face']}")
f2 = sum(1 for line in open(f"{out}/a.log") if "送鍵 3c" in line)
check(f2 == 0, f"驗收1 F2 送鍵紀錄 {f2}")
for k in ("Step", "Mem", "A000", "Pal"):
    check(a[k] == b[k], f"驗收2 {k} 有地圖={a[k]} 無地圖={b[k]}")
xl = open(f"{out}/x.log").read()
check("地圖 面板=true" in xl and "送鍵 3c" not in xl, "驗收3 Xvfb 按 F2 開啟面板且未送原版（截圖 x.png 另行目視）")
print(f"失敗 {len(fails)} 項"); sys.exit(1 if fails else 0)
PY
