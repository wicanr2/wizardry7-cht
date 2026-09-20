# Ebiten 在 Linux 上要 cgo 連 X11／OpenGL，golang 官方 image 沒有那些標頭檔。
# 這是**本專案自己的 image**（`psychicwar-go-ebiten`），不動任何共用資源。
#
# `xdotool` ＋ `imagemagick` 是給 `tools/frontend-playthrough.sh` 用的：
# 把前端跑在 Xvfb 裡、用真正的 X 鍵盤事件操作、用 `import` 截圖。
# **不是除錯捷徑**——按鍵一樣走 X → GLFW → Ebiten → `inpututil`，
# 與玩家在自己機器上按的是同一條路（`docs/spec/006` §4 第 5 項驗這一條）。
#
# `libasound2-plugins` 是給真實裝置上的音訊量測用的（`docs/spec/020`）：
# ⚠ 它會拉進一大串多媒體相依（ffmpeg、cairo、codec2…），第一次 build 要等很久。
# 只需要其中的 `libasound_module_pcm_pulse.so`，但 Debian 沒有更小的套件。
# 不裝 alsa-utils：`aplay` 只是方便，plugin 本身不需要它。
# 它提供 ALSA 的 pulse plugin，讓容器經由主機的 PipeWire／PulseAudio socket 出聲，
# **不獨佔音效卡**。直接掛 /dev/snd 會跟主機的音訊伺服器搶裝置。
FROM golang:1.24-bookworm
RUN apt-get update && apt-get install -y --no-install-recommends \
        libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev \
        libgl1-mesa-dev libasound2-dev pkg-config xvfb xauth libgl1 libglx-mesa0 \
        xdotool imagemagick \
        libasound2-plugins \
    && rm -rf /var/lib/apt/lists/*
