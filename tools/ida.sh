#!/usr/bin/env bash
# IDA Pro 9.4 headless 包裝。image 來源 ~/ida_94_official。
#
#   tools/ida.sh build  <執行檔>              從 workplace/original 複製一份並建庫（.i64 + .asm）
#   tools/ida.sh script <腳本.py> [資料庫]     對 .i64 的**副本**跑腳本
#   tools/ida.sh raw    <idat 參數…>          直接下 idat（會改寫 .i64）
#
# 產物在 workplace/ida/（gitignore）。原版檔案唯讀，複製一份進去建庫——
# IDA 會在輸入檔旁邊寫 .i64，而原版目錄要保持唯讀。
#
# script 模式跑的是副本：idat 一開啟 .i64 就會改寫它的雜湊，而筆記要靠
# 雜湊標明「這個結論是在哪一份資料庫上驗的」。對副本跑，原始 .i64 的
# 身分才穩定。
#
# ⚠ headless 的 print 不進 stdout，exit code 也不可信（同一種失敗在不同
# image 上分別回 0 與 1）。**唯一可信的訊號是輸出檔**——腳本一律寫檔，
# 收工前驗檔案存在、非空、schema 對。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WIZ7_IDA_IMAGE:-ida-pro-9.4-idapython:locked-v1}"
WORK="$ROOT/workplace/ida"
mkdir -p "$WORK"

docker image inspect "$IMAGE" >/dev/null 2>&1 || {
  echo "[ida.sh] 找不到 $IMAGE。IDAPython 在沒修過的 image 上是零輸出的" >&2
  echo "         靜默失敗，不要退回別顆 image 重試。" >&2
  exit 3; }

# 容器內以呼叫者的 UID 跑，產物就不必事後 chown。
# HOME 要指到 image 裡 idapyswitch 寫設定的那個家目錄，否則 IDAPython
# 載不起來——而症狀同樣是「什麼都沒有」。
run() {
  local dir="$1"; shift
  timeout "${WIZ7_IDA_TIMEOUT:-30m}" docker run --rm \
    --network none --memory "${WIZ7_IDA_MEM:-4g}" --cpus "${WIZ7_IDA_CPUS:-2}" \
    --pids-limit 256 --log-opt max-size=10m --log-opt max-file=3 \
    -u "$(id -u):$(id -g)" -e HOME=/home/ubuntu \
    -v "$dir:/work" -v "$ROOT/tools/ida:/tools:ro" -w /work \
    "$IMAGE" "$@"
}

MODE="${1:-}"; shift || true

case "$MODE" in
  build)
    BIN="${1:-}"
    [[ -n "$BIN" ]] || { echo "[ida.sh] 要建哪一支執行檔" >&2; exit 2; }
    SRC="${WIZ7_ORIG:-$ROOT/workplace/original/wizar7}/$BIN"
    [[ -f "$SRC" ]] || { echo "[ida.sh] 找不到 $SRC" >&2; exit 2; }
    cp -f "$SRC" "$WORK/$BIN"; chmod u+w "$WORK/$BIN"
    echo "輸入："; sha256sum "$WORK/$BIN"
    rm -f "$WORK/$BIN".i64 "$WORK/$BIN".id0 "$WORK/$BIN".id1 \
          "$WORK/$BIN".nam "$WORK/$BIN".til
    run "$WORK" idat -A -B "$BIN"
    ls -la "$WORK/$BIN".i64 "$WORK/$BIN".asm
    ;;
  script)
    PY="${1:-}"; DB="${2:-PW.EXE.i64}"; shift 2 || true
    [[ -n "$PY" ]] || { echo "[ida.sh] 要跑哪一支腳本" >&2; exit 2; }
    SCRATCH="$WORK/scratch"; mkdir -p "$SCRATCH"
    # 中斷過的執行會留下沒打包的資料庫（.id0/.id1/.nam/.til）。IDA 看到它們
    # 就不肯開同名的 .i64，訊息是「Failed to initialize IDA as library
    # (error code 4)」——看起來像授權或 image 壞掉，其實只是殘檔。
    rm -f "$SCRATCH/${DB%.i64}".id0 "$SCRATCH/${DB%.i64}".id1 \
          "$SCRATCH/${DB%.i64}".nam "$SCRATCH/${DB%.i64}".til
    echo "來源資料庫："; sha256sum "$WORK/$DB"
    cp -f "$WORK/$DB" "$SCRATCH/$DB"; chmod u+w "$SCRATCH/$DB"
    run "$SCRATCH" idat -A "-S/tools/$(basename "$PY") $*" "$DB"
    echo "原始資料庫雜湊（應與上方相同）："; sha256sum "$WORK/$DB"
    ;;
  raw) run "$WORK" "$@" ;;
  *) echo "用法: ida.sh {build|script|raw} …" >&2; exit 2 ;;
esac
