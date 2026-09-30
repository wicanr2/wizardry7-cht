#!/usr/bin/env bash
# fe-input-check.sh — spec 003 R3.2 的前端實測（docs/re/015 T2–T4）。
#
#   tools/fe-input-check.sh [T2|T3|T4]...   # 預設全跑
#
# Xvfb 裡開真正的前端視窗，autoplay 跑到指定步數後停手（-autoplay-stop），
# 再用 xdotool 送真實鍵盤／滑鼠事件。每項都有不送輸入的對照組。
#   T2 防拷畫面打字 AAAA＋Enter → 進入旁白（log 出現旁白第一頁的印字）
#   T3 旁白印完（印字靜止 3 秒）後按 Enter → 2 秒內有新印字（對照組 2 秒內無新印字）
#   T4 迷宮開場旁白期間點擊 → 畫面改變（翻下一頁旁白；對照組不變）
#   T5 （無效，#38）開場旁白期間點羅盤座標，量到的是旁白翻頁，不能證明轉向；由 T6 取代
#   T6 先點到旁白不再有新印字，再點羅盤前進、右轉 → 畫面大幅改變且無新印字（對照組不點羅盤）
# 產物在 workplace/r666/（gitignore）。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT=workplace/r666
mkdir -p "$ROOT/$OUT"
TESTS="${*:-T2 T3 T4 T5}"

read -r -d '' INNER <<'EOS' || true
set -u
go build -o /tmp/wiz7 ./cmd/wiz7 || exit 3
start() { # $1＝log $2＝autoplay-stop
  /tmp/wiz7 -root workplace/run -scale 3 -autoplay -autoplay-stop "$2" >"$1" 2>&1 &
  PID=$!
}
waitfor() { # $1＝log $2＝grep 樣式 $3＝秒
  for _ in $(seq 1 $(( $3 * 2 ))); do grep -qE "$2" "$1" && return 0; sleep 0.5; done
  return 1
}
prints() { grep -cE '(翻譯|無譯文) [0-9]+ ' "$1" || true; }  # 有譯文與無譯文的印字都要算（T6 曾因只數翻譯而誤判）
win() { xdotool search --name '巫術7' | head -1; }
stop() { kill $PID 2>/dev/null; wait $PID 2>/dev/null || true; }

t2() { # $1＝組名 $2＝yes/no
  local log=/src/$OUT/T2-$1.log
  start "$log" 280000000
  waitfor "$log" '"SECURITY CLEARANCE"' 120 || { echo "T2-$1：等不到防拷"; stop; return; }
  sleep 1
  if [ "$2" = yes ]; then xdotool type --window "$(win)" --delay 150 AAAA; xdotool key --window "$(win)" Return; fi
  sleep 8
  stop
  echo "T2-$1：旁白第一頁印字=$(grep -c 'Cosmic Forge' "$log" || true)"
}
t3() {
  local log=/src/$OUT/T3-$1.log
  start "$log" 400000000
  waitfor "$log" 'Cosmic Forge' 120 || { echo "T3-$1：等不到旁白"; stop; return; }
  # 等本頁印完：印字數連續 3 秒不變（最多 90 秒）
  local p0 last=-1 still=0
  for _ in $(seq 1 180); do
    p0=$(prints "$log")
    if [ "$p0" = "$last" ]; then still=$((still + 1)); else still=0; last=$p0; fi
    [ $still -ge 6 ] && break
    sleep 0.5
  done
  echo "T3-$1：靜止前印字=$p0 靜止判定=$still"
  if [ "$2" = yes ]; then xdotool key --window "$(win)" Return; fi
  sleep 2
  local p1; p1=$(prints "$log")
  stop
  echo "T3-$1：按鍵後 2 秒新印字=$(( p1 - p0 ))"
}
t4() {
  local log=/src/$OUT/T4-$1.log
  start "$log" 1435000000
  waitfor "$log" 'gentle wind' 300 || { echo "T4-$1：等不到迷宮"; stop; return; }
  sleep 4
  import -window root "/src/$OUT/T4-$1-a.png" 2>/dev/null
  if [ "$2" = yes ]; then xdotool mousemove --window "$(win)" 480 486 click 1; fi
  sleep 4
  import -window root "/src/$OUT/T4-$1-b.png" 2>/dev/null
  stop
  # 視口（遊戲座標約 72–248 × 33–143，3× 且視窗在 Xvfb 內的偏移由 import 全螢幕含入，取整張比較）
  local d; d=$(compare -metric AE "/src/$OUT/T4-$1-a.png" "/src/$OUT/T4-$1-b.png" null: 2>&1 || true)
  echo "T4-$1：點擊前後差異像素=$d 左鍵紀錄=$(grep -c '左鍵 (' "$log" || true) OUCH=$(grep -c 'OUCH' "$log" || true)"
}
t5() { # 羅盤右轉 (180,166)、左轉 (145,166)、後退 (160,184)：每下與前一張比
  local log=/src/$OUT/T5.log
  start "$log" 1435000000
  waitfor "$log" 'gentle wind' 300 || { echo "T5：等不到迷宮"; stop; return; }
  sleep 4
  import -window root "/src/$OUT/T5-0.png" 2>/dev/null
  local n=0 prev=0 d
  for xy in "540 498" "435 498" "480 552"; do
    n=$((n + 1))
    xdotool mousemove --window "$(win)" $xy click 1
    sleep 4
    import -window root "/src/$OUT/T5-$n.png" 2>/dev/null
    d=$(compare -metric AE "/src/$OUT/T5-$prev.png" "/src/$OUT/T5-$n.png" null: 2>&1 || true)
    echo "T5：第 $n 下（視窗 $xy）差異像素=$d"
    prev=$n
  done
  stop
  echo "T5：左鍵紀錄=$(grep -c '左鍵 (' "$log" || true)"
}
t6() { # 羅盤實測（#38）：先翻完開場旁白（點到沒有新印字），再點羅盤
  local log=/src/$OUT/T6-$1.log
  start "$log" 1435000000
  waitfor "$log" 'gentle wind' 300 || { echo "T6-$1：等不到迷宮"; stop; return; }
  sleep 4
  local k p0 p1
  for k in $(seq 1 12); do
    p0=$(prints "$log")
    xdotool mousemove --window "$(win)" 480 510 click 1
    sleep 4
    p1=$(prints "$log")
    [ "$p1" = "$p0" ] && break
  done
  echo "T6-$1：翻旁白 $k 次後無新印字"
  import -window root "/src/$OUT/T6-$1-0.png" 2>/dev/null
  local n=0 xy d q0 q1
  for xy in "480 486" "540 498"; do   # 前進 (160,162)、右轉 (180,166)
    n=$((n + 1))
    q0=$(prints "$log")
    [ "$2" = yes ] && xdotool mousemove --window "$(win)" $xy click 1
    sleep 4
    q1=$(prints "$log")
    import -window root "/src/$OUT/T6-$1-$n.png" 2>/dev/null
    d=$(compare -metric AE "/src/$OUT/T6-$1-$((n - 1)).png" "/src/$OUT/T6-$1-$n.png" null: 2>&1 || true)
    echo "T6-$1：第 $n 下（視窗 $xy）差異像素=$d 新印字=$((q1 - q0))"
  done
  stop
}
for t in $TESTS; do
  case $t in
    T2) t2 input yes; t2 ctl no ;;
    T3) t3 input yes; t3 ctl no ;;
    T4) t4 input yes; t4 ctl no ;;
    T5) t5 ;;
    T6) t6 input yes; t6 ctl no ;;
  esac
done
EOS

WIZ7_FE_TIMEOUT=40m WIZ7_FE_SH="OUT=$OUT; TESTS='$TESTS'; $INNER" "$ROOT/tools/fe.sh"
