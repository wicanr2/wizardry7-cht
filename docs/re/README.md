# 研究證據索引

本目錄保存可回查的原版定位、量測與推論等級。[目前狀態](../../CONTEXT.md)、[規格](../spec/)、[工作歷程](../../WORKLOG.md) 各有獨立職責。
較早研究中的舊方案與假說不直接成為現況；以後續勘誤、目前程式及同狀態收據判定。

| 文件 | 用途 |
|---|---|
| [001 原版基線](001-unzip-baseline.md) | 使用者本機輸入的大小與雜湊 |
| [002 DS.EXE 初探](002-ds-exe-first-look.md) | executable 與 overlay 初步定位 |
| [003 開機](003-golem-boot.md) | dosgolem 開機與執行期組態 |
| [004 IDA 普查](004-ida-census.md) | 原始位址、函式及 DOS／BIOS／EMS 分布 |
| [005 DOSBox-X 輔助對拍](005-dosboxx-parity.md) | 標題比較與檢查點，非全遊戲驗收 |
| [006 M1 量測](006-m1-progress.md) | 顯示與執行環境進展 |
| [007 故事等待](007-story-wait-vgadrv.md) | 驅動段與事件等待研究 |
| [008 滑鼠回呼](008-int33-ax14-fix.md) | INT 33h AX=14h 缺口與修復 |
| [009 決定性與早期回歸](009-round3.md) | 顯示、狀態、音訊事件比較 |
| [010 正常試玩](010-playthrough.md) | 建隊、輸入與早期迷宮路線 |
| [011 VGA 格式](011-vga-format.md) | 原版影像資產格式研究 |
| [012 MSG.DBS 初探](012-msgdbs-format.md) | 訊息格式早期證據，後續以 013 為準 |
| [013 印字管線與全文抽取](013-print-pipeline.md) | resident 入口與訊息資料；歷史寫入方案已由 014 更正 |
| [014 唯讀覆繪](014-overlay-evidence.md) | 印字前後差異、原版 A/B 雜湊及顯示驗收 |
| [015 F1 說明](015-f1-help-evidence.md) | 前端送鍵與說明面板證據 |
| [016 隊伍位置](016-party-position.md) | 移動與面向量測；未證實欄位語意維持原等級 |
| [017 流動文字](017-flow-text-evidence.md) | 旁白分頁、故事31頁來源與返回邊界、淡入／淡出原版像素證據；驗收見spec006 |
| [M3 早期筆記](m3-notes.md)、[後續筆記](m3-notes2.md) | 歷史線索，不能取代後續資料庫或量測證據 |
