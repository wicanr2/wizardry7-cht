#!/usr/bin/env bash
# flow-check.sh — docs/spec/006 驗收 1–11（流動文字）。
#
#   tools/flow-check.sh
#
# 1 en 模式（-flow ""）走迷宮開場旁白，log 的「流文 片段」是獨立檢查的頁來源
# 2 en／zh(3×)／zh(2×)／缺段落 各跑一輪，同一串 -hd-at（旁白期間每 1M 步）＋ -hd-after-frag（第 4 頁前兩個片段返回的當下）
# 3 tools/re/check_flow.py 獨立比對（A/B 雜湊、頁譯文、排版、像素、英文不殘留、過渡、負向）
# 產物在 workplace/r670/（gitignore）。回歸（check_overlay、fe-help-check、automap-check、fe-leftclick-check、
# fe-input-check）另跑。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT=workplace/r670
cd "$ROOT"; rm -rf "$OUT"; mkdir -p "$OUT"/{zh,zh2,en,drop}

py() { timeout 10m docker run --rm --network none --memory 1g --cpus 1 --pids-limit 64 \
  --log-opt max-size=10m --log-opt max-file=3 -u "$(id -u):$(id -g)" -e HOME=/tmp \
  -e PYTHONDONTWRITEBYTECODE=1 -v "$ROOT:/src" -v "$HOME/cht/etan_font:/tmp/cht/etan_font:ro" -w /src \
  python:3.13-alpine python3 "$@"; }

echo "== 烘製 flow.json 與缺段落測試檔"
py tools/re/bake_flow.py
py tools/re/bake_flow.py $OUT/flow-drop.json --drop 15005

CLK=$(python3 -c "print(','.join(f'{s}000000:160:170' for s in range(1440, 1541, 10)))")
COMMON="-hdonly -autoplay -autoplay-stop 1435000000 -root workplace/run -clicks-at $CLK"

# 新頁（第 4 頁）首片段的序號＝前 3 頁片段數（4＋4＋1）＋1；驗收 10 取第 4 頁前兩個片段返回的當下
FRAG="10,11"
HD=$(python3 -c "
l=list(range(1483000000,1539000000,1000000))
print(','.join(str(x) for x in sorted(set(l))))")

echo "== en／zh／zh 2×／缺段落"
WIZ7_FE_TIMEOUT=25m WIZ7_FE_SH="go build -o /tmp/wiz7 ./cmd/wiz7 && timeout 280 /tmp/wiz7 $COMMON -hd-dir $OUT/en -hd-after-frag $FRAG -hd-at $HD -flow '' > $OUT/en.log 2>&1; timeout 280 /tmp/wiz7 $COMMON -hd-dir $OUT/zh -hd-after-frag $FRAG -hd-at $HD > $OUT/zh.log 2>&1; timeout 280 /tmp/wiz7 $COMMON -hd-dir $OUT/zh2 -hd-after-frag $FRAG -hd-at $HD -hd-scale 2 > $OUT/zh2.log 2>&1; timeout 280 /tmp/wiz7 $COMMON -hd-dir $OUT/drop -hd-after-frag $FRAG -hd-at $HD -flow $OUT/flow-drop.json > $OUT/drop.log 2>&1; true" tools/fe.sh >/dev/null 2>&1

echo "== 比對"
grep -c "流文 帶不符" $OUT/zh.log | sed 's/^/流文 帶不符 次數（期望 0）：/' || true
py tools/re/check_flow.py --zh $OUT/zh --en $OUT/en --log $OUT/en.log --scale 3 \
   --zh2 $OUT/zh2 --scale2 2 --drop $OUT/drop --drop-para 15005 --frag-dumps "$FRAG" | tee $OUT/check.txt | grep -v "^通過" 
