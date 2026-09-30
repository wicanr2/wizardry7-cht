#!/usr/bin/env bash
# fe-help-check.sh — spec 003 R6 驗收（F1 說明）。
#
#   tools/fe-help-check.sh
#
# 1 暫停不改原版：同路線 -hd-at 300M，有無 -help-at 286M -help-for 3000 的雜湊相同
# 2 暫停期間步數凍結：log 的「說明 暫停=true／false」步數相同
# 3 說明畫面截圖：Xvfb 互動視窗 -help-at 286M -help-page 2 -draw-at 286M（防拷畫面說明頁）
# 4 缺字時 bake_help.py 非零結束且不產檔；5 lint：ops 行引用未通過的實測編號時失敗
# 6 畫面辨識單元測試（go test）
# 產物在 workplace/r667/（gitignore）。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT=workplace/r667
cd "$ROOT"; rm -rf "$OUT"; mkdir -p "$OUT/a" "$OUT/b" "$OUT/draw"

py() { timeout 5m docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
  --log-opt max-size=10m --log-opt max-file=3 -u "$(id -u):$(id -g)" -e HOME=/tmp \
  -e PYTHONDONTWRITEBYTECODE=1 "$@"; }
FONT=(-v "$HOME/cht/etan_font:/tmp/cht/etan_font:ro")

echo "== 4 缺字、5 lint"
printf 'key\tzh\tsource\nkeys.00\t缺字測試\tx\nkeys.01\t𠮷\tx\n' > "$OUT/missing.tsv"
printf 'key\tzh\tsource\nops.00\t基本操作\tx\nops.01\t沒有實測的操作\t015 T99\n' > "$OUT/lint.tsv"
set +e
py -v "$ROOT:/src" "${FONT[@]}" -w /src -e HELP_TSV=$OUT/missing.tsv python:3.13-alpine \
  python3 tools/re/bake_help.py $OUT/missing.json; e4=$?
py -v "$ROOT:/src" "${FONT[@]}" -w /src -e HELP_TSV=$OUT/lint.tsv python:3.13-alpine \
  python3 tools/re/bake_help.py $OUT/lint.json; e5=$?
set -e
echo "4 缺字 exit=$e4 產檔=$([ -f $OUT/missing.json ] && echo 有 || echo 無)"
echo "5 lint exit=$e5"

WIZ7_FE_TIMEOUT=20m WIZ7_FE_SH="set -u
go test ./cmd/wiz7 -run DetectScreen -v 2>&1 | grep -E '^(=== RUN|--- |ok|FAIL|PASS)|skip'
go build -o /tmp/wiz7 ./cmd/wiz7 || exit 3
timeout 90 /tmp/wiz7 -hdonly -autoplay -root workplace/run -hd-dir $OUT/a -hd-at 300000000 > $OUT/a.log 2>&1
timeout 90 /tmp/wiz7 -hdonly -autoplay -root workplace/run -hd-dir $OUT/b -hd-at 300000000 -help-at 286000000 -help-for 3000 > $OUT/b.log 2>&1
timeout 90 /tmp/wiz7 -autoplay -root workplace/run -scale 3 -hd-dir $OUT/draw -help-at 286000000 -help-page 2 -help-for 6000 -draw-at 286000000 > $OUT/draw.log 2>&1
true" tools/fe.sh 2>&1 | grep -v "XSERV\|XGB\|Xauthority" || true

echo "== 1 雜湊"
grep -h "hdat hdat-300000000 " $OUT/a.log $OUT/b.log | sed 's/.*steps=/steps=/' | cut -d' ' -f1-4
echo "== 2 暫停步數"
grep -h "說明 暫停" $OUT/b.log
echo "== 3 截圖"
grep -h "drawat 輸出\|說明 true" $OUT/draw.log
