#!/bin/sh
set -eu
umask 077
mkdir /out/runtime /out/runtime/config /out/runtime/pulse
chmod 700 /out/runtime /out/runtime/config /out/runtime/pulse
export XDG_RUNTIME_DIR=/out/runtime
export XDG_CONFIG_HOME=/out/runtime/config
export PULSE_SERVER=unix:/out/runtime/pulse/native
export ALSA_CONFIG_PATH=/out/asound.conf
ulimit -f 524288
viz_pid=
pulse_pid=
monitor_pid=
cleanup() {
 for task_pid in "$monitor_pid" "$pulse_pid" "$viz_pid"; do
  if test -n "$task_pid"; then kill "$task_pid" 2>/dev/null || :; wait "$task_pid" 2>/dev/null || :; fi
 done
}
trap cleanup EXIT INT TERM
cat > /out/pulse.pa <<'EOF'
load-module module-native-protocol-unix socket=/out/runtime/pulse/native auth-anonymous=1
load-module module-null-sink sink_name=wiz7 rate=44100 channels=2
set-default-sink wiz7
EOF
cat > /out/asound.conf <<'EOF'
pcm.!default {
 type pulse
 server "unix:/out/runtime/pulse/native"
}
ctl.!default {
 type pulse
 server "unix:/out/runtime/pulse/native"
}
EOF
pulseaudio -n --daemonize=no --use-pid-file=no --exit-idle-time=-1 --disable-shm=yes --log-target=file:/out/pulse.log --file=/out/pulse.pa &
pulse_pid=$!
tries=0
until pactl info > /out/pulse-info.txt 2>/dev/null; do
 tries=$((tries+1));test "$tries" -le 100 || exit 1;sleep .1
done
Xvfb :99 -screen 0 1280x900x24 -nolisten tcp > /out/xvfb.log 2>&1 &
viz_pid=$!
export DISPLAY=:99
tries=0
until xdotool getdisplaygeometry >/dev/null 2>&1; do
 tries=$((tries+1));test "$tries" -le 50 || exit 1;sleep .1
done
python3 -c 'from pathlib import Path;Path("/out/prime.s16le").write_bytes(bytes(882*4))'
paplay --raw --format=s16le --rate=44100 --channels=2 /out/prime.s16le
parec --device=wiz7.monitor --rate=44100 --channels=2 --format=s16le --raw --latency-msec=20 > /out/device.s16le 2> /out/monitor.log &
monitor_pid=$!
python3 /out/normal_audio.py
