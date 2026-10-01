#!/usr/bin/env bash
# spec 006：在單一既有 Docker 工具鏈重生故事與旁白收據，再做獨立驗收。
# tools/flow-check.sh workplace/rNNN（新目錄；不覆寫舊基準）
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${1:-workplace/r674}"
[[ "$OUT" =~ ^workplace/r[0-9]+$ ]] || { echo '輸出須為新的 workplace/rNNN' >&2; exit 2; }
test -d "$ROOT"
test -d "$ROOT/workplace"
test -d "$ROOT/workplace/original"
test -d "$HOME/cht/etan_font"
test ! -e "$ROOT/$OUT" || { echo '收據目錄已存在，請換新編號' >&2; exit 2; }
docker image inspect wizardry7-go-ebiten:latest >/dev/null
for target in workplace workplace/bake workplace/bake/bake.json workplace/bake/fontseg.bin workplace/bake/xlate.json workplace/bake/flow.json; do
  if test -e "$ROOT/$target"; then
    owner="$(stat -c '%u:%g' "$ROOT/$target")"
    test "$owner" = "$(id -u):$(id -g)" || { echo "擁有權不符：$target $owner" >&2; exit 2; }
  fi
done
exec timeout 35m docker run --rm --name "wiz7-flow-${OUT##*/}" --network none \
  --memory 2g --cpus 2 --pids-limit 256 --log-opt max-size=10m --log-opt max-file=3 \
  -u "$(id -u):$(id -g)" -e HOME=/tmp -e GOFLAGS=-mod=readonly \
  -e GOCACHE=/src/workplace/gocache -e GOMODCACHE=/src/workplace/gomodcache \
  -e DISPLAY=:99 -e PYTHONDONTWRITEBYTECODE=1 \
  -v "$ROOT:/src:ro" -v "$ROOT/workplace:/src/workplace" \
  -v "$ROOT/workplace/original:/src/workplace/original:ro" \
  -v "$HOME/cht/etan_font:/tmp/cht/etan_font:ro" -w /src \
  wizardry7-go-ebiten:latest bash -c '
set -e
Xvfb :99 -screen 0 1024x768x24 >/tmp/x.log 2>&1 & xp=$!
trap "kill $xp 2>/dev/null || true" EXIT
sleep .5
python3 tools/re/run_story_checks.py "$1"
python3 tools/re/check_story.py --en "$1/en" --zh "$1/zh" --zh2 "$1/zh2" --drop "$1/drop" --log "$1/en.log" >"$1/story-check.txt"
python3 tools/re/check_flow.py --en "$1/en" --zh "$1/zh" --zh2 "$1/zh2" --drop "$1/narr-drop" --log "$1/en.log" --frag-dumps 81,82 >"$1/narr-check.txt"
tail -2 "$1/story-check.txt"
tail -2 "$1/narr-check.txt"
find "$1" -user root -o -type d -name "*.md"
stat -c "%u:%g %n" "$1" workplace/bake/flow.json workplace/bake/xlate.json
' bash "$OUT"
