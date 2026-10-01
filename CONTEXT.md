# 目前狀態與交接

核對日期：2026-10-01。最近實作與驗收：R671；本機 dosgolem：`d16f03a2d8d782024d470606e1f4cb5272b89cd6`。
本專案是 dosgolem 執行期覆繪中文化，原版遊戲規則與資料保持原樣。
可執行工作唯一來源是 [GitHub Issue](https://github.com/wicanr2/wizardry7-cht/issues)。
本檔記錄目前真相；[WORKLOG.md](WORKLOG.md) 保存歷程與勘誤。

## 已完成的實作與驗證範圍

| 項目 | 目前證據與界線 |
|---|---|
| 原版啟動、標題、主選單與決定性 | [009](docs/re/009-round3.md)：同輸入狀態／記憶體／幀雜湊相同；標題 DOSBox-X 輔助比較差異 0。這不代表整個遊戲已逐像素驗證。 |
| 文字抽取與印字管線 | [013](docs/re/013-print-pipeline.md)：MSG.DBS 抽出 11,019 則；resident `0110:3895` 是 NUL 字串入口。 |
| 前端與正常玩家操作 | [cmd/wiz7](cmd/wiz7/main.go)、[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)：建隊、進迷宮與早期遭遇戰已有操作收據；未有從開頭到結局的完整驗收。 |
| 唯讀中文覆繪 | [spec 001](docs/spec/001-overlay-no-write.md)、[014](docs/re/014-overlay-evidence.md)：中英兩模式原版記憶體、A000、調色盤雜湊相同；只適用實測畫面。部分防禦分支未觸發，規格仍 READY。 |
| 隊員名鍵衝突修正 | [spec 002](docs/spec/002-bare-slot-key.md)、[#32](https://github.com/wicanr2/wizardry7-cht/issues/32)：單獨隊員名不再誤套序數譯文。 |
| F1 說明與前端滑鼠 | [spec 003](docs/spec/003-f1-help.md)、[spec 004](docs/spec/004-frontend-left-click.md)：已有驗收紀錄；原版保留鍵調查與畫面說明只涵蓋實測範圍。 |
| 自動地圖 | [spec 005](docs/spec/005-automap.md)：已驗兩格移動與 F2 面板，整張地圖與記錄持久化未完成；[#36](https://github.com/wicanr2/wizardry7-cht/issues/36) 仍開啟。 |
| 迷宮開場旁白 | [spec 006](docs/spec/006-flow-text.md)、[017](docs/re/017-flow-text-evidence.md)：6頁同狀態A/B、獨立譯文及像素驗收；最新`workplace/r674/narr-check.txt`為失敗0，56份原版狀態與r670全同。倍率3中文字形放大2倍。 |
| 開場故事 | [spec 006](docs/spec/006-flow-text.md)：31頁、71片段、58則訊息；兩倍率、獨立來源／譯文／像素與缺訊息回退驗收7,096項、失敗0。124份原版狀態與修訂前全同，不再依賴story.txt。3個1M暫態的RGB未獨立量到；未量防禦分支維持READY。 |

## 收尾尚缺的證據與交付

- 全文來源仍是 JSON 與訊息 ID；TSV＋英文鍵遷移見 [#24](https://github.com/wicanr2/wizardry7-cht/issues/24)。譯文覆蓋與正常路徑觸發對帳見 [#22](https://github.com/wicanr2/wizardry7-cht/issues/22)，資料已翻譯不能當作已上屏。
- 發行字型尚未完成（[#25](https://github.com/wicanr2/wizardry7-cht/issues/25)），目前使用的倚天字型只可本機使用。沒有可核實的三平台正式交付與通關收據。
- 輸入錄放、音訊完成條件仍以 [#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[#18](https://github.com/wicanr2/wizardry7-cht/issues/18) 為準；早期 OPL2 事件比較通過不等於前端聲音與播放對齊驗收完成。
- 作弊功能見 [#35](https://github.com/wicanr2/wizardry7-cht/issues/35)，尚未有正式驗收；不可用作弊建立中文化回歸收據。
- 規格覆蓋與 adapter 位置仍待 [#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#29](https://github.com/wicanr2/wizardry7-cht/issues/29) 驗收。既有六份 spec 並不代表所有玩家文字路徑都已規格化。

## 有效決定與限制

以 [AGENTS.md](AGENTS.md) 第 13 節為完整決定來源：放大畫布、不縮中文字；預設 3×；TSV 英文鍵；字型雙軌；防拷維持原版驗證；F1 說明、自動／全圖兩種地圖與作弊範圍。
防拷關鍵字資料移出版控的 [#26](https://github.com/wicanr2/wizardry7-cht/issues/26) 依 2026-09-29 指示擱置，不能自行處理。
對話關鍵字與防拷答案仍以英文輸入；不得推導、輸出或提交防拷答案。
儲存庫保持 private，原版與不可散布字型不進公開交付。

## 本輪交接

[R671](docs/goal/R671-story-overlay.md) 完成故事覆繪與獨立驗收，收據在`workplace/r674`；6M／8M全部31頁有效，旁白六頁回歸通過且56份原版狀態與`r670`全同。
最後故事頁3M有前頁筆畫淡出殘留，靜態閘門保持英文，6M／8M消退；證據F26與規格訂正經唯讀審查無阻擋。正式判準未放寬。
前端覆繪、F1暫停、左鍵、T2／T4、地圖兩格與視窗F2均通過（`r675`、`r676`）。T6初次混跑差異為0，未採為通過；`r680`同工具乾淨重跑前進／轉向為157,563／159,975像素，控制組各279，通過。初次原因未定，完整失敗收據保留於r675。
舊工作清單 [docs/worklist.json](docs/worklist.json) 保留為歷史快照，不再更新；未追蹤的 R433 goal 不混入提交。
歷史基線`workplace/r670/`、原版探針`r673`已保留；[tools/flow-check.sh](tools/flow-check.sh) 要求新的`workplace/rNNN`，拒絕覆寫舊基準，所有工作負載在Docker執行。
