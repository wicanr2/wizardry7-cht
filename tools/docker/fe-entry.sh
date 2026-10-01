#!/usr/bin/env bash
# 僅在專案工具鏈容器內執行；由tools/fe.sh控制掛載與生命週期。
set -euo pipefail
assert_directory() {
  local directory="$1"
  test ! -L "$directory"
  test -d "$directory"
  test "$(stat -c '%u:%g' "$directory")" = "$(id -u):$(id -g)"
}
case "${1:-}" in
  init)
    assert_directory /src
    for directory in /src/workplace /src/workplace/gocache /src/workplace/gomodcache; do
      test ! -L "$directory"
      if [[ ! -e "$directory" ]]; then mkdir "$directory"; fi
      assert_directory "$directory"
    done
    ;;
  go)
    shift
    assert_directory /src/workplace/gocache
    assert_directory /src/workplace/gomodcache
    exec go "$@"
    ;;
  xvfb)
    shift
    [[ $# == 1 ]]
    assert_directory /src/workplace/gocache
    assert_directory /src/workplace/gomodcache
    command -v Xvfb >/dev/null
    command -v xdotool >/dev/null
    Xvfb :99 -screen 0 1024x768x24 -nolisten tcp >/tmp/wiz7-fe-xvfb.log 2>&1 &
    xvfb_pid=$!
    cleanup_x() { kill "$xvfb_pid" 2>/dev/null || true; wait "$xvfb_pid" 2>/dev/null || true; }
    trap cleanup_x EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    ready=false
    for ((attempt=0; attempt<30; attempt++)); do
      if DISPLAY=:99 xdotool getdisplaygeometry >/dev/null 2>&1; then ready=true; break; fi
      if ! kill -0 "$xvfb_pid" 2>/dev/null; then break; fi
      sleep 0.1
    done
    if [[ "$ready" != true ]]; then echo '[fe.sh] 容器Xvfb未就緒。' >&2; exit 1; fi
    export DISPLAY=:99
    bash -c "$1"
    ;;
  *) echo '[fe.sh] 容器入口模式錯誤。' >&2; exit 2 ;;
esac
