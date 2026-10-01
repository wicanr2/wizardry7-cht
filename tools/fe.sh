#!/usr/bin/env bash
# 前端開發入口；主機只做Docker控制面，工作負載一律留在容器。
# tools/fe.sh build ./cmd/wiz7
# tools/fe.sh xvfb 'go test ./cmd/wiz7'
# WIZ7_FE_SH='指令' tools/fe.sh
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
IMAGE="${WIZ7_FE_IMAGE:-wizardry7-go-ebiten:latest}"
BOUND="${WIZ7_FE_TIMEOUT:-10m}"
test -d "$ROOT"
test -f "$ROOT/tools/docker/fe-entry.sh"
[[ "$BOUND" =~ ^[1-9][0-9]*[smhd]?$ ]] || { echo '[fe.sh] 逾時須為正整數，可加s/m/h/d。' >&2; exit 2; }
MODE=go
PAYLOAD="${WIZ7_FE_SH:-}"
if [[ "${1:-}" == xvfb ]]; then
  shift
  if [[ -n "$PAYLOAD" ]]; then
    [[ $# == 0 ]] || { echo '[fe.sh] xvfb指令與WIZ7_FE_SH不可並用。' >&2; exit 2; }
  else
    [[ $# == 1 ]] || { echo "[fe.sh] 用法：tools/fe.sh xvfb '指令'" >&2; exit 2; }
    PAYLOAD="$1"
    shift
  fi
fi
if [[ -n "$PAYLOAD" ]]; then
  [[ $# == 0 ]] || { echo '[fe.sh] WIZ7_FE_SH不可與Go引數並用。' >&2; exit 2; }
  MODE=xvfb
else
  [[ $# -gt 0 ]] || { echo '[fe.sh] 請指定Go命令或xvfb指令。' >&2; exit 2; }
fi
if ! timeout --kill-after=1s 5s docker image inspect "$IMAGE" >/dev/null 2>&1; then
  echo '[fe.sh] 缺少既有工具鏈映像檔；請依README.md的開發入口準備，不自動開網路建置。' >&2
  exit 1
fi

INPUTS=()
if [[ -e "$ROOT/Wizardry_VII_-_Crusaders_of_the_Dark_Servant_1992.zip" ]]; then
  test -f "$ROOT/Wizardry_VII_-_Crusaders_of_the_Dark_Servant_1992.zip"
  INPUTS+=(-v "$ROOT/Wizardry_VII_-_Crusaders_of_the_Dark_Servant_1992.zip:/src/Wizardry_VII_-_Crusaders_of_the_Dark_Servant_1992.zip:ro")
fi
if [[ -e "$ROOT/workplace/original" ]]; then
  test -d "$ROOT/workplace/original"
  INPUTS+=(-v "$ROOT/workplace/original:/src/workplace/original:ro")
fi
NAME="wiz7-fe-$$-$RANDOM"
INIT_NAME="$NAME-init"
for candidate in "$NAME" "$INIT_NAME"; do
  if ! existing=$(timeout --kill-after=1s 5s docker container ls -a \
    --filter "name=^/${candidate}$" --format '{{.Names}}'); then
    echo '[fe.sh] 無法完成容器名稱檢查；未操作既有容器。' >&2
    exit 1
  fi
  if [[ -n "$existing" ]]; then
    echo '[fe.sh] 本趟容器名稱已存在；未操作既有容器。' >&2
    exit 1
  fi
done
CONTROL_PID=''
cleanup() {
  timeout --kill-after=1s 5s docker rm -f "$NAME" "$INIT_NAME" >/dev/null 2>&1 || true
  if [[ -n "$CONTROL_PID" ]]; then
    kill -TERM "$CONTROL_PID" 2>/dev/null || true
    wait "$CONTROL_PID" 2>/dev/null || true
  fi
}
wait_control() {
  local result=0
  wait "$CONTROL_PID" || result=$?
  CONTROL_PID=''
  return "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
COMMON=(--rm --init --network none --memory 2g --cpus 2 --pids-limit 256
  --log-opt max-size=10m --log-opt max-file=3 -u "$(id -u):$(id -g)"
  -e HOME=/tmp -v "$ROOT:/src" "${INPUTS[@]}" -w /src --entrypoint bash)
# 掛入尚未存在的快取前，先在已驗證的工作樹掛載內建立並檢查擁有權。
timeout --kill-after=5s 30s docker run --name "$INIT_NAME" "${COMMON[@]}" \
  "$IMAGE" /src/tools/docker/fe-entry.sh init &
CONTROL_PID=$!
wait_control
for directory in "$ROOT/workplace" "$ROOT/workplace/gocache" "$ROOT/workplace/gomodcache"; do
  test -d "$directory"
  test ! -L "$directory"
done
FE_ENV=(-e GOFLAGS=-mod=readonly -e GOPROXY=off -e GOSUMDB=off
  -e GOCACHE=/src/workplace/gocache -e GOMODCACHE=/src/workplace/gomodcache
  -e DISPLAY=:99 -v "$ROOT/workplace/gomodcache:/src/workplace/gomodcache:ro")
if [[ "$MODE" == xvfb ]]; then
  ARGS=(xvfb "$PAYLOAD")
else
  ARGS=(go "$@")
fi
timeout --kill-after=5s "$BOUND" docker run --name "$NAME" "${COMMON[@]}" "${FE_ENV[@]}" \
  "$IMAGE" /src/tools/docker/fe-entry.sh "${ARGS[@]}" &
CONTROL_PID=$!
wait_control
