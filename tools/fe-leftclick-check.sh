#!/usr/bin/env bash
# fe-leftclick-check.sh — spec 004 驗收：互動前端的滑鼠左鍵送進原版。
#
#   tools/fe-leftclick-check.sh
#
# Xvfb 裡開真正的前端視窗（不帶 -autoplay），輪詢 log 等主選單印字，
# 再用 xdotool 對視窗內 (480,186)（3×；遊戲座標 (160,62)「新增角色」）點左鍵。
# 對照組同樣等到主選單但不點。判定：點擊組 log 出現建角對話框的 OK／CANCEL 印字，
# 對照組沒有。產物在 workplace/r666/（gitignore）。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT=workplace/r666
mkdir -p "$ROOT/$OUT"

read -r -d '' INNER <<'EOS' || true
set -u
go build -o /tmp/wiz7 ./cmd/wiz7 || exit 3
run() { # $1＝組名 $2＝是否點擊
  local log=/src/$OUT/$1.log
  /tmp/wiz7 -root workplace/run -scale 3 >"$log" 2>&1 &
  local pid=$!
  for _ in $(seq 1 240); do
    grep -q '"CHARACTER MENU"' "$log" && break
    sleep 0.5
  done
  grep -q '"CHARACTER MENU"' "$log" || { echo "$1：等不到主選單"; kill $pid; return 1; }
  sleep 1
  if [ "$2" = yes ]; then
    local w
    w=$(xdotool search --name '巫術7' | head -1)
    xdotool mousemove --window "$w" 480 186 click 1
  fi
  sleep 8
  import -window root "/src/$OUT/$1.png" 2>/dev/null || true
  kill $pid 2>/dev/null; wait $pid 2>/dev/null
  local dlg
  dlg=$(grep -cE '翻譯 [0-9]+ "(OK|CANCEL)"' "$log" || true)
  echo "$1：左鍵紀錄=$(grep -c '左鍵 (' "$log" || true) 對話框 OK/CANCEL 印字=$dlg"
}
run click yes
run noclick no
EOS

WIZ7_FE_TIMEOUT=10m WIZ7_FE_SH="OUT=$OUT; $INNER" "$ROOT/tools/fe.sh"
