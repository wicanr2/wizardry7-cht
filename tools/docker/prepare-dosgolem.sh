#!/usr/bin/env bash
# 僅為Docker控制面；來源核對與套用都在隔離容器中。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
test -d "$ROOT"
test -d "$ROOT/worktrees/dosgolem"
test ! -L "$ROOT/worktrees"
test ! -L "$ROOT/worktrees/dosgolem"
test -f "$ROOT/tools/docker/prepare-dosgolem.py"
test -f "$ROOT/tools/docker/dosgolem-audio-manifest.json"
test -f "$ROOT/tools/docker/dosgolem-audio.patch"
IMAGE="${WIZ7_FE_IMAGE:-wizardry7-go-ebiten:latest}"
docker image inspect "$IMAGE" >/dev/null
NAME="wiz7-prepare-dosgolem-$$"
cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
timeout --kill-after=5s 60s docker run --rm --init --name "$NAME" --network none \
  --memory 768m --cpus 1 --pids-limit 64 -u "$(id -u):$(id -g)" \
  --log-opt max-size=10m --log-opt max-file=3 \
  -v "$ROOT:/src:ro" -v "$ROOT/worktrees/dosgolem:/src/worktrees/dosgolem" \
  --entrypoint python3 "$IMAGE" /src/tools/docker/prepare-dosgolem.py
