#!/usr/bin/env bash
# 前端（Ebiten）建置與執行。image 是本專案自己的 wizardry7-go-ebiten
# （X11／GL／ALSA 標頭＋Xvfb＋xdotool，第一次自動 build，不動共用資源）。
#   tools/fe.sh build ./cmd/wiz7     # 編譯
#   tools/fe.sh run ./cmd/wiz7 …     # 需要 X；無頭請用 tools/fe.sh xvfb '指令'
#   WIZ7_FE_SH='指令' tools/fe.sh    # 在 Xvfb 裡跑任意指令
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WIZ7_FE_IMAGE:-wizardry7-go-ebiten}"
if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
  echo "[fe.sh] 第一次使用，先 build $IMAGE" >&2
  docker build -t "$IMAGE" -f "$ROOT/tools/docker/go-ebiten.Dockerfile" "$ROOT/tools/docker"
fi
mkdir -p "$ROOT/workplace/gocache" "$ROOT/workplace/gomodcache"
if [[ "${1:-}" == "xvfb" ]]; then shift; fi
ARGS=()
if [[ -n "${WIZ7_FE_SH:-}" ]]; then
  ARGS=(bash -c "$WIZ7_FE_SH")
else
  ARGS=("$@")
fi
exec timeout "${WIZ7_FE_TIMEOUT:-10m}" docker run --rm --network none \
  --memory 2g --cpus 2 --pids-limit 256 \
  -u "$(id -u):$(id -g)" -e HOME=/tmp -e GOFLAGS=-mod=mod \
  -e GOCACHE=/src/workplace/gocache -e GOMODCACHE=/src/workplace/gomodcache \
  -e DISPLAY=:99 -v "$ROOT:/src" -w /src \
  "$IMAGE" "${ARGS[@]}"
