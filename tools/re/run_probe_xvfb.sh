#!/bin/sh
set -eu
Xvfb :99 -screen 0 1280x900x24 -nolisten tcp > /tmp/xvfb.log 2>&1 &
viz_pid=$!
trap 'kill "$viz_pid" 2>/dev/null || :; wait "$viz_pid" 2>/dev/null || :' EXIT INT TERM
export DISPLAY=:99
tries=0
until xdotool getdisplaygeometry >/dev/null 2>&1; do
  tries=$((tries+1))
  test "$tries" -le 50 || exit 1
  sleep .1
done
"$@"
