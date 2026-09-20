#!/usr/bin/env bash
# Python 走 docker，不用主機的 Python。
#
#   tools/py.sh tools/worklist.py verify
#   tools/py.sh tools/re/exe_first_look.py
#
# 整個 repo 掛在 /src（可寫）；原版素材在 workplace/（gitignore），唯讀掛載另見各腳本。
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WIZ7_PY_IMAGE:-python:3.13-alpine}"

exec timeout "${WIZ7_TIMEOUT:-5m}" docker run --rm -i --network none \
  --memory 512m --cpus 1 --pids-limit 64 \
  --log-opt max-size=10m --log-opt max-file=3 \
  -u "$(id -u):$(id -g)" -e HOME=/tmp -e PYTHONDONTWRITEBYTECODE=1 \
  -v "$ROOT:/src" -w /src "$IMAGE" python3 "$@"
