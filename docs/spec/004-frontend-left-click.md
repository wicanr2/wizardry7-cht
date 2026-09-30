# 004 — 互動前端的滑鼠左鍵

狀態：READY（2026-09-30，審查一輪：應改 1／建議 2，已處理）
issue：#37
證據：`cmd/wiz7/main.go`（Update 記錄 `sh.buttons[0]`，機器 goroutine 只讀 `sh.buttons[1]`；
`sh.clicks` 無互動來源）；`o.Click` 的座標系＝320×200 遊戲座標（autoplay 點 (160,62) 進建角、
(82,169) 點掉匯入詢問後 `Mouse()`＝(82,169)，`docs/re/014` E8）。

## 為什麼要做

巫術7 的選單、對話框、迷宮羅盤、戰鬥指令都靠左鍵（`docs/re/015` H7–H9）。互動前端目前點左鍵不會送進原版，
玩家無法遊玩；F1 說明（spec 003）也無從教操作。

## 規格

- R1：Update 偵測左鍵「剛按下」時，以當下游標換算的遊戲座標 `(mx/scale, my/scale)`（裁切到 0–319、0–199）
  排進 `sh.clicks`。
- R2：機器 goroutine 既有的 `sh.clicks` 消費路徑以 `o.Click(x, y)` 送出（與 autoplay 同一 API，
  含 hover、按住、放開與 settle）。`o.Click` 回傳錯誤只記 log，不中止。
- R3：`sh.buttons[0]` 不再使用，移除（核對過無其他讀取點）；`g.lbtn` 可保留作判斷中介；右鍵維持既有 Press／Release 路徑。
- R4：F1 說明開啟時（spec 003）不排入點擊。本規格實作時說明尚未存在，此條留待 003 實作。

## 驗收

1. Xvfb 下以 `-scale 3` 啟動前端（不帶 `-autoplay`），輪詢前端 log 直到出現 `翻譯 … "CHARACTER MENU"`
   （主選單已印字；不用固定 sleep，開機到主選單的牆鐘時間會漂移，`docs/re/010` R8），再用 xdotool 在視窗座標
   (480,186)（遊戲座標 (160,62)，「新增角色」）按一下左鍵。
2. `-draw-at` 截取點擊後的畫面；與不點擊的同一跑法比較：點擊組出現建角對話框（log 有 `ADD CHARACTER` 對話框的
   `OK`、`CANCEL` 印字），未點擊組停在主選單。
3. log 記錄每次左鍵排入的遊戲座標與 `o.Click` 結果。

## 未決

- 牆鐘與步數的對應受 `-ips` 節流影響，xdotool 點擊時點不決定性；驗收只比「有無對話框」，不比雜湊。
- 非字元鍵經 IRQ1 的問題（015 H1、H5）另案。
- 本驗收只涵蓋主選單單擊。羅盤按鈕很小，歷史上需要較長間隔或重試（R100、R101），其可靠度另案驗證；
  004 通過不代表迷宮點擊已保證可用。

## 驗收紀錄

2026-09-30：`tools/fe-leftclick-check.sh` 點擊組 `左鍵 (160,62) err=<nil>`，建角對話框 OK／CANCEL 印字 2；
對照組 0（`docs/re/015` T1）。只涵蓋主選單單擊，維持 READY。
