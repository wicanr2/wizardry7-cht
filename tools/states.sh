#!/usr/bin/env bash
# states.sh — 檢查點狀態檔產生器（輪 3 版）。
#
#   tools/states.sh            # 跑到主選單（940M 步），存狀態檔＋幀
#   tools/states.sh --check    # 從 menu 狀態恢復，點第一列，驗證到達 PCVW
#
# 檢查點（docs/re/006、008、009）：
#   banner  200M 步：橫幅完成＋副標文字
#   title   520M 步：穩定標題（31,054 非零；RGB 對拍 0）
#   menu    940M 步：主選單 MAINSCRN（AX=14h 事件註冊完成後，可點擊）
#   pcvw    點主選單第一列 (159,47) → PCVWSCRN（55,600；遮罩與 DOSBox-X 一致）
#
# 前置：workplace/run/ 是 wizar7 的可寫複本（掛 /w）。
# 產物在 workplace/states/（gitignore）。步數一律絕對值。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLEM="$ROOT/worktrees/dosgolem"
RUN="$ROOT/workplace/run"
OUT="$ROOT/workplace/states"
PCVW_SHA="142c2139b5841929db4e0a06af3abe5f6eb042e6ce604d006de0c3a27ee916c3"
[[ -f "$RUN/DS.EXE" ]] || { echo "缺 $RUN/DS.EXE（先解壓原版到 workplace/original 再複製）"; exit 2; }
mkdir -p "$OUT"

if [[ "${1:-}" == "--check" ]]; then
  cd "$GOLEM"
  DOSGOLEM_EXTRA_MOUNT="$RUN:/w" \
    tools/go.sh run ./cmd/probe -load-state /w/states/menu.state \
    -steps 945000000 -adlib -click-premove 3 -clicks "940500000:159:47" \
    -dump-vram /w/check-restore.vram >/dev/null
  got=$(sha256sum "$RUN/check-restore.vram" | cut -d' ' -f1)
  want=$(python3 -c "import json;print(json.load(open('$OUT/meta.json'))['frames']['pcvw']['vram_sha256'])" 2>/dev/null || echo "$PCVW_SHA")
  echo "PCVW 參考：${want:0:12}"
  echo "點擊後：  ${got:0:12}"
  [[ "$got" == "$want" ]] && echo "OK：menu 狀態恢復＋點擊推進到 PCVW" || { echo "FAIL"; exit 1; }
  exit 0
fi

[[ -d "$GOLEM" && -d "$RUN" ]] || { echo "缺 worktrees/dosgolem 或 workplace/run"; exit 2; }

cd "$GOLEM"
DOSGOLEM_EXTRA_MOUNT="$RUN:/w" \
  DOSGOLEM_TIMEOUT=30m tools/go.sh run ./cmd/probe \
  -exe /w/DS.EXE -root /w -steps 1000000000 -adlib \
  -shots "200000000:/w/states/banner.vram,520000000:/w/states/title.vram,940000000:/w/states/menu.vram" \
  -save-state "520000000:/w/states/title.state,940000000:/w/states/menu.state" \
  > "$OUT/run.log" 2>&1 || true

for f in banner title menu; do
  [[ -s "$OUT/$f.vram" ]] || { echo "缺 $f 產物，看 $OUT/run.log"; exit 1; }
done
[[ -s "$OUT/menu.state" && -s "$OUT/title.state" ]] || { echo "缺狀態檔"; exit 1; }

"$ROOT/tools/py.sh" - "workplace/states" <<'EOF'
import hashlib, json, sys
from pathlib import Path
out = Path(sys.argv[1])
def sha(p): return hashlib.sha256(p.read_bytes()).hexdigest()
meta = {
    "command": "probe -exe /w/DS.EXE -root /w -steps 1000000000 -adlib "
               "（-shots banner/title/menu；-save-state title/menu）",
    "frames": {
        "banner": {"steps": 200000000},
        "title": {"steps": 520000000},
        "menu": {"steps": 940000000},
        "pcvw": {"steps": "menu+點擊(159,47)", "vram_sha256": None},
    },
    "states": {n: {"sha256": sha(out / f"{n}.state"), "bytes": (out / f"{n}.state").stat().st_size}
               for n in ("title", "menu")},
}
for n in ("banner", "title", "menu"):
    meta["frames"][n]["vram_sha256"] = sha(out / f"{n}.vram")
json.dump(meta, open(out / "meta.json", "w"), indent=1)
print(json.dumps(meta, indent=1))
print("PCVW 參考幀跑一次 --check 流程後登記 meta.frames.pcvw")
EOF
