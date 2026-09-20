#!/usr/bin/env bash
# DOSBox-X 參照畫面（M0 #6）：第二個實作走到標題，存畫面給對拍。
#
#   tools/dosboxx-ref.sh    # 產出 workplace/dosboxx/title-NN.{x11.png,rgb}
#
# 組態對齊 dosgolem 實跑：oplmode=opl2（golem -adlib）、EMS 開（遊戲用 EMS，
# docs/re/003）、直接跑 DS.EXE（不經 BAT）。動畫以計時器推進，
# 標題出現時點依主機速度而定：每 5 秒截一張、最多 24 張，比對時挑「非零密度
# ≈ golem 標題幀」那張。scaler=none、aspect=false → 640×400＝320×200 的 2 倍。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WIZ7_DOSBOXX_IMAGE:-civ1-dosboxx-input:20260830}"
GAME="${WIZ7_ORIG_DIR:-$ROOT/workplace/original}/wizar7"
OUT="$ROOT/workplace/dosboxx"

[[ -f "$GAME/DS.EXE" ]] || { echo "找不到 $GAME/DS.EXE" >&2; exit 2; }
docker image inspect "$IMAGE" >/dev/null 2>&1 || { echo "找不到 image $IMAGE" >&2; exit 2; }
mkdir -p "$OUT"
rm -f "$OUT"/title-*.png "$OUT"/title-*.rgb 2>/dev/null || true

cat > "$OUT/.run.sh" <<'EOF'
set -euo pipefail
export HOME=/tmp DISPLAY=:99
Xvfb :99 -screen 0 1024x768x24 >/tmp/xvfb.log 2>&1 &
sleep 2
mkdir -p /tmp/game
cp -r /orig/. /tmp/game/
chmod -R u+w /tmp/game
cat > /tmp/dosbox.conf <<CONF
[sdl]
output=surface
autolock=false
windowresolution=original
[dosbox]
machine=svga_s3
memsize=8
[render]
aspect=false
scaler=none
[cpu]
core=normal
cputype=386
cycles=fixed 20000
[sblaster]
sbtype=none
oplmode=opl2
[gus]
gus=false
[speaker]
pcspeaker=true
[autoexec]
mount c /tmp/game
c:
DS.EXE
CONF
dosbox-x -conf /tmp/dosbox.conf -nomenu >/tmp/dbx.log 2>&1 &
DBX=$!
sleep 6
WIN=$(xdotool search --name 'DOSBox-X' | tail -1)
eval "$(xdotool getwindowgeometry --shell "$WIN")"
xdotool mousemove $((X + WIDTH / 2)) $((Y + HEIGHT / 2))
xdotool windowfocus --sync "$WIN" || true
snap() {
  import -window "$WIN" "/out/$1.x11.png"
  convert "/out/$1.x11.png" -depth 8 "rgb:/out/$1.rgb"
  echo "[$1] 擷取" >&2
}
for i in $(seq -w 1 24); do
  snap "title-$i"
  sleep 5
done
echo "視窗 ${WIDTH}x${HEIGHT}" >&2
kill $DBX 2>/dev/null || true
EOF

echo "tools/dosboxx-ref.sh：跑 DOSBox-X（約 2.5 分鐘）" >&2
timeout "${WIZ7_DOSBOXX_TIMEOUT:-400}" docker run --rm --network none \
  --memory 2g --cpus 2 --pids-limit 128 \
  --log-opt max-size=10m --log-opt max-file=3 \
  -u "$(id -u):$(id -g)" \
  -v "$GAME:/orig:ro" -v "$OUT:/out" \
  "$IMAGE" bash /out/.run.sh 2>&1 | grep -v -E 'XGetInputFocus|^$' || true
rm -f "$OUT/.run.sh"
ls -l "$OUT" | head -8
