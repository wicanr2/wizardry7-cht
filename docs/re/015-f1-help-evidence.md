# 015 — F1 說明：證據（#34）

日期：2026-09-30。探針 `cmd/keyprobe`（原始輸出 `workplace/bake/r665-keyprobe*.txt`）、
`cmd/wiz7 -autoplay -hd-at`（`workplace/r665/maze/`）、既有試玩紀錄。

## 結論

| # | 事實 | 等級 |
|---|---|---|
| H1 | 經 IRQ1（`oracle.PressKey`）送鍵，主選單、防拷、旁白三處連正對照（Enter、`A`、Space）都無反應；此路徑的否定結果無效 | confirmed |
| H2 | 經 DOS stdin（`oracle.Type`）送 F1（`00 3B`）：三處都被讀走 2 bytes。防拷畫面 F1 後 30M 步 A000 差異 0、對照 `A` 差異 36 px；旁白 F1 差異 0、對照 Enter 差異 60,804 px 且印字 6 次 | confirmed |
| H3 | 主選單的正對照（Enter）也無反應，主選單為滑鼠操作；F1 在主選單的結論只到「讀走、無可見變化」 | confirmed（對照無效） |
| H4 | 迷宮沒有已知有效的鍵盤正對照（R303／R304：L、F 被讀走但無作用；`docs/re/010`：部分畫面原版直接輪詢 port 60h）。F1 在迷宮是否有作用：未量到 | 未知 |
| H5 | 前端目前把非字元鍵經 IRQ1 送出（`cmd/wiz7` `Update` → `PressKey`），依 H1 這些鍵原版多半收不到 | 強證據 |
| H6 | 可辨識畫面的上屏字串：主選單＝`CHARACTER MENU`／`START NEW GAME`；建角對話框＝`ADD CHARACTER`＋`OK`＋`CANCEL`（y≥150）；防拷＝`SECURITY CLEARANCE`。迷宮（14.9 億步）上屏只有職業名與底部旁白，沒有專屬字串 | confirmed（`r665/maze/hdat-*.json`、`r663/zh3/hdat-*.json`） |
| H7 | 迷宮移動：羅盤按鈕 前進 (160,162)、後退 (160,184)、右轉 (180,166)、左轉 (145,166)，按鈕很小 | confirmed（R102、R104、R120、R124、R126） |
| H8 | 戰鬥：遭遇後出現指令選單，(262,166)＝RUN 設為逃跑，START 執行回合（`-VS-`、`FINAL ORDERS`） | confirmed（R120） |
| H9 | 旁白：30M 步（約 1.25 秒）窗口內，Enter 讓旁白提前推進而不按鍵不會（H2 對照）；但旁白也會自行推進（T3），前端實測分不出按鍵效果 | confirmed（短窗口）／前端未量到 |
| H10 | 名字與防拷答案經 DOS stdin 輸入，Enter 確認 | confirmed（`docs/re/010` 追記二、R663 防拷） |
| H11 | 倚天 `ascfont.15`：3,840 bytes＝256 字 × 15 列，每列 1 byte（8 px 寬） | confirmed（檔長與 F、1 字形目視） |

## 推論

- 由 H2：防拷與旁白讀到 F1 但不使用；前端攔下 F1 不影響這兩處。
- 由 H4：前端攔下 F1 若在迷宮擋掉原版功能，現有證據看不出來；規格記為風險。
- 由 H6：「目前畫面說明」只能對主選單、建角、防拷三種畫面做專屬內容，其他畫面用通用內容。

## 前端實測紀錄（spec 003 R3.2 的 `source` 引用這裡的編號）

| 編號 | 操作 | 方法 | 結果 |
|---|---|---|---|
| T1 | 主選單左鍵點「新增角色」 | `tools/fe-leftclick-check.sh`（Xvfb＋xdotool 真實點擊，spec 004） | 通過：建角對話框出現（OK／CANCEL 印字 2 次；對照組 0） |
| T2 | 防拷畫面鍵盤打字＋Enter | `tools/fe-input-check.sh T2`（`-autoplay-stop 280000000`） | 通過：旁白第一頁出現（印字 1；對照組 0） |
| T3 | 旁白印完後按 Enter | `tools/fe-input-check.sh T3`（等印字靜止 3 秒） | 無法判定：按鍵組與對照組 2 秒內都有 4 次新印字，旁白會自行推進 |
| T4 | 迷宮開場旁白中點擊（座標 (160,162)） | `tools/fe-input-check.sh T4`（`-autoplay-stop 1435000000`） | 通過（重新定性，#38）：點擊翻到下一頁旁白，畫面差異 13,192 px（對照組 0）。原記為「羅盤前進」有誤：當時旁白尚未結束，羅盤未出現 |
| T5 | 迷宮依序點右轉 (180,166)、左轉 (145,166)、後退 (160,184) | `tools/fe-input-check.sh T5` | 無效（#38）：三下都落在開場旁白期間，量到的是旁白翻頁 |
| T6 | 翻完開場旁白（點到無新印字）後點羅盤前進 (160,162)、右轉 (180,166) | `tools/fe-input-check.sh T6`（印字計數含無譯文） | 通過：差異 157,545／160,047 px 且無新印字（對照組 0／0）；與 016 P1 的 17,500 遊戲像素 ×9 相符 |
