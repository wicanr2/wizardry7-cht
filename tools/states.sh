#!/usr/bin/env bash
# states.sh — M0 #5：檢查點狀態檔產生器。
#
#   tools/states.sh            # 跑到標題（380M 步，按鍵腳本推進），存狀態檔＋幀
#   tools/states.sh --check    # 從標題狀態恢復，驗證畫面與繼續執行
#
# 檢查點（依 docs/re/003 的開機時間軸）：
#   banner  200M 步：橫幅完成＋副標文字
#   title   380M 步：全螢幕標題美術頁（等待輸入）
# 產物在 workplace/states/（gitignore）。步數一律絕對值。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLEM="$ROOT/worktrees/dosgolem"
ORIG="$ROOT/workplace/original/wizar7"
OUT="$ROOT/workplace/states"
mkdir -p "$OUT"
SHOTS="200000000:/out/states/banner.vram,380000000:/out/states/title.vram"
STATES="200000000:/out/states/banner.state,380000000:/out/states/title.state"

if [[ "${1:-}" == "--check" ]]; then
  cd "$GOLEM"
  DOSGOLEM_ORIG="$ORIG" DOSGOLEM_EXTRA_MOUNT="$ROOT/workplace:/out" \
    tools/go.sh run ./cmd/probe -load-state /out/states/title.state \
    -steps 381000000 -adlib -dump-vram /out/check-restore.vram \
    -dump-palette /out/check-restore.pal >/dev/null
  "$ROOT/tools/py.sh" - "workplace/states" <<'EOF'
import hashlib, json, sys
from pathlib import Path
out = Path(sys.argv[1])
meta = json.loads((out / "meta.json").read_text())
new = hashlib.sha256((out.parent / "check-restore.vram").read_bytes()).hexdigest()
old = meta["frames"]["title"]["vram_sha256"]
print("title 幀雜湊：", old)
print("恢復後雜湊：", new)
if new != old:
    print("FAIL：恢復後畫面不同（標題應是等待輸入的靜態畫）")
    sys.exit(1)
print("OK：從狀態檔恢復，畫面一致")
EOF
  exit 0
fi

cd "$GOLEM"
DOSGOLEM_ORIG="$ORIG" DOSGOLEM_EXTRA_MOUNT="$ROOT/workplace:/out" \
  DOSGOLEM_TIMEOUT=30m tools/go.sh run ./cmd/probe \
  -exe /orig/DS.EXE -root /orig -steps 400000000 -adlib \
  -press enter,space,esc -press-at 20000000 -press-every 300000 \
  -shots "$SHOTS" -save-state "$STATES" \
  -dump-vram /out/final.vram -dump-palette /out/final.pal \
  > "$OUT/run.log" 2>&1 || true

for f in banner title; do
  [[ -s "$OUT/$f.state" && -s "$OUT/$f.vram" ]] || { echo "缺 $f 產物，看 $OUT/run.log"; exit 1; }
done

"$ROOT/tools/py.sh" - "workplace/states" <<'EOF'
import hashlib, json, sys
from pathlib import Path
out = Path(sys.argv[1])
meta = {
    "command": "probe -exe /orig/DS.EXE -root /orig -steps 400000000 -adlib "
               "-press enter,space,esc -press-at 20000000 -press-every 300000",
    "frames": {
        "banner": {"steps": 200000000,
                   "vram_sha256": hashlib.sha256((out / "banner.vram").read_bytes()).hexdigest()},
        "title": {"steps": 380000000,
                  "vram_sha256": hashlib.sha256((out / "title.vram").read_bytes()).hexdigest()},
    },
    "states": {
        n: {"sha256": hashlib.sha256((out / f"{n}.state").read_bytes()).hexdigest(),
            "bytes": (out / f"{n}.state").stat().st_size}
        for n in ("banner", "title")
    },
}
json.dump(meta, open(out / "meta.json", "w"), indent=1)
print(json.dumps(meta, indent=1))
EOF
