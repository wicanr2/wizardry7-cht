# 009 — 輪 3：AX=14h 上游化、handler 持久化、選單點擊地圖（#14、#9、#13 收尾）

日期：2026-09-20。issue #14（結）、#13（已結）、#9（進度）。

## dosgolem 上游化（#14，結）

- `wizardry7/int33-ax14` 分支（spec 194＋`case 0x0014`＋contract test）
  **沒有直接 merge**：clone 時的 origin/master（d0ad8b9 線）與現行
  origin/master（ff35717 線）**已分叉**（上游歷史被重寫）——原 merge 會
  倒退 224 個 commit。改用 **cherry-pick** 重放 3 檔修復：
  dosgolem master `3c21907`。
- 追加：**狀態快照補存 `Mouse.Handler`**（`internal/dos/state.go`；
  master `ec71176`）。漏存的話從快照展開的機器「滑鼠有裝、事件不來」
  （巫術7 主選單實測）。新增 `TestStateKeepsMouseEventHandler`。
- dosgolem 無實質 CI workflow（僅 Dependency Graph bot）；
  以本地 docker `tools/go.sh test ./internal/dos` 全綠替代。
- 教訓記帳：**worktree 是共享資源**——推上游前先 `git fetch`＋確認
  自己分支的基底還在 origin/master 的祖先內。

## 主選單點擊地圖（#9 核心，confirmed）

`menu.state`（940M 步，AX=14h 註冊完成後）恢復＋點擊五列
（座標來自驅動熱點表 `(105,43)-(213,51)` 起、每列 +13px）：

| 點擊 | 結果畫面 | 底圖 | 非零 |
|---|---|---|---:|
| 列1 (159,47) | PCVWSCRN | `142c2139` | 55,600 |
| 列2 (159,60) | PCVWSCRN（不同面板） | `f79827c9` | 55,600 |
| 列3 (159,73) | PCVWSCRN（不同面板） | `f73f5f53` | 55,600 |
| 列4 (159,86) | PCVWSCRN（不同面板） | `02096452` | 55,600 |
| 列5 (159,99) | PCVWSCRN（不同面板） | `08cc8fba` | 55,600 |

- **MAINSCRN.VGA（43,530 相位）＝主選單**（5 列選單項）；
  五列各自選中不同項，都畫在 PCVWSCRN 底圖上（55,600 家族）。
- 事件投遞實錄：遮罩 000B、handler 740C:1ED4，每個點擊 down＋up 各一筆。

## DOSBox-X 輸入正對照（#9 收尾條件，通過）

主選單畫面 `xdotool` 點遊戲座標 (159,47)（視窗 2 倍 318,94）：
點擊前後差 61,267 像素、點擊後非黑 **55,600**、與 golem PCVW 幀
**遮罩差異 0**——兩個模擬器在同一輸入下收斂到同一張畫面。

## 檢查點（states.sh v2）

| 檢查點 | 步數 | 狀態檔 | 備註 |
|---|---:|---|---|
| banner | 200M | —（舊格式棄用） | 幀 `69220ee7` |
| title | 520M | title.state | 穩定標題，RGB 對拍 0 |
| menu | 940M | menu.state | **新格式（含 Handler）**，可點擊 |
| pcvw | menu＋點擊 | —（由 menu 派生） | `142c2139…`，--check 參考 |

`tools/states.sh --check`：恢復 menu.state → 點 (159,47) → 雜湊與
PCVW 參考一致 → **通過**。

## 決定性成對測試（#11，通過）

同一 binary、同一輸入、並行兩次跑到 940M：

| 產物 | pairA | pairB | 判定 |
|---|---|---|---|
| 狀態檔（machine＋dos） | `e895df54…` | `e895df54…` | **完全相同** |
| 驅動段記憶體 dump（0x740C0-0x7C0C0） | `4c690774…` | `4c690774…` | **完全相同** |
| 跨 run 幀雜湊 | — | — | banner `69220ee7`×3、menu `e141d21a`×多輪、PCVW `142c2139`×2 一致 |

**dosgolem 附帶修復（master `de5e2d2`）**：`Ports/PortsIn` 快照原本是 map，
gob 編碼順序隨機 → 兩次一模一樣的執行存出的狀態檔 bytes 不同。改成排序過的
slice 後狀態檔才可能用雜湊驗決定性。畫面決定性另有跨 run 幀雜湊為證
（工具小毛病：`-dump-mem` 的解析在收尾階段才跑，失敗訊息會蓋在正常結束之後；
`-dump-vram` 模式 13h 於本輪部分 run 未產出，以既有幀雜湊替代）。

## 遺留

- 五個選單項的**文字名稱**未解（MSG.DBS 壓縮格式，M3）。
- 從 PCVW 往下走（建隊伍→進遊戲）＝模擬試玩，M2/輪 4。
- DOSBox-X 輸入送達的正對照已做（本輪）；之後每次換輸入法都要重驗。
