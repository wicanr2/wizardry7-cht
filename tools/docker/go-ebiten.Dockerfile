# Ebiten 在 Linux 上要 cgo 連 X11／OpenGL，golang 官方 image 沒有那些標頭檔。
# 這是本專案工具鏈映像檔`wizardry7-go-ebiten`的建置來源，不動共用資源。
#
# `xdotool` ＋ `imagemagick` 用於前端驗證：
# 把前端跑在 Xvfb 裡、用真正的 X 鍵盤事件操作、用 `import` 截圖。
# **不是除錯捷徑**——按鍵一樣走 X → GLFW → Ebiten → `inpututil`，
# 與玩家在自己機器上按的是同一條前端輸入路徑；驗收範圍依本案規格與收據。
#
# `libasound2-plugins` 是給音訊量測用的（`docs/spec/009-player-audio.md`）：
# ⚠ 它會拉進一大串多媒體相依（ffmpeg、cairo、codec2…），第一次 build 要等很久。
# 只需要其中的 `libasound_module_pcm_pulse.so`，但 Debian 沒有更小的套件。
# 不裝 alsa-utils：`aplay` 只是方便，plugin 本身不需要它。
# 它提供 ALSA 的 pulse plugin，讓容器經由主機的 PipeWire／PulseAudio socket 出聲，
# **不獨佔音效卡**。直接掛 /dev/snd 會跟主機的音訊伺服器搶裝置。
ARG GO_IMAGE=golang:1.24.13-bookworm
ARG GUI_IMAGE=eob-remake-go:1.26.7-ebiten2.9.9-20261008-r2
ARG XAUTH_IMAGE=colonization-verification:20261007-runtime-r2

FROM ${GO_IMAGE} AS go-locked

# 原完整工具鏈；只有明示允許建置時聯網才可用full target。
FROM go-locked AS full
RUN apt-get update && apt-get install -y --no-install-recommends \
        libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev \
        libgl1-mesa-dev libasound2-dev pkg-config xvfb xauth libgl1 libglx-mesa0 \
        xdotool imagemagick \
        libasound2-plugins \
    && rm -rf /var/lib/apt/lists/*

# R977：舊1d7映像已不存在。複用已核Debian12 GUI依賴，乾淨替換Go SDK，
# 全程--network none；只支援建置、回歸及靜音GUI，Pulse ALSA plugin仍未恢復。
FROM ${XAUTH_IMAGE} AS xauth-source
FROM ${GUI_IMAGE} AS recovered-gui
USER root
RUN rm -rf /usr/local/go
COPY --from=go-locked /usr/local/go /usr/local/go
COPY --from=xauth-source /usr/bin/xauth /usr/bin/xauth
ENV GOLANG_VERSION=1.24.13 GOTOOLCHAIN=local \
    PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
RUN test "$(go env GOVERSION)" = go1.24.13 \
    && command -v Xvfb && command -v xauth && command -v xdotool \
    && pkg-config --exists x11 xrandr xcursor xinerama xi xxf86vm gl alsa \
    && ! ldd /usr/bin/xauth | grep 'not found'
USER 1000:1000
WORKDIR /src
