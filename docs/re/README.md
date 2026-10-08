# 研究證據索引

譯文上屏抽樣的接手複核：[R972目標](../goal/R972-sampling.md)，[025結果與清理限制](025-ui-localization-coverage.md#r972-sampling-takeover)。
驗收後工作區清理與紙頁來源復核：[R973目標](../goal/R973-workspace-cleanup-paper-evidence.md)、[025來源與限制](025-ui-localization-coverage.md#r973-paper-source-recovery)。

戰鬥控制的中文覆繪檢查：[R959目標](../goal/R959-battle-control-display.md)，驗收接入[025](025-ui-localization-coverage.md)。

森林遭遇過渡的原版續驗：[R958證據](019-catalog-migration.md#r958-forest-transition-wait)、[目標](../goal/R958-forest-transition-wait.md)。

森林事件翻頁與十一組情境：[R957證據](019-catalog-migration.md#r957-forest-event-context)、[目標](../goal/R957-forest-event-context.md)。

預設速度正常點擊音訊：[R956](020-player-audio.md#r956-預設速度的正常點擊音訊)，[本輪目標](../goal/R956-default-audio-clicks.md)。

新證據：[021 物品頁兩行名稱](021-compound-item-output.md)。
待補譯文正常路線：[R898 東向受阻與南向續行](019-catalog-migration.md#r898-normal-east-translation)。
最遠合法路線與取樣時序勘誤：[R900](019-catalog-migration.md#r900-farthest-maze-continuation)。
正常路線及重複休息流程：[R902](019-catalog-migration.md#r902-normal-route-context)。
正常岔路及排除範圍分類：[R904](019-catalog-migration.md#r904-safe-branch-context)。
對話詞庫守門修正及正常續行：[R906](019-catalog-migration.md#r906-lexicon-guard)。
新事件文字翻頁及返回邊界：[R909](019-catalog-migration.md#r909-event-pages)。
事件返回畫面及正常路線：[R912](019-catalog-migration.md#r912-event-return-route)。
正式中文介面的文字與圖片按鈕覆蓋重查：[025](025-ui-localization-coverage.md)。
世界人物頁性別與選人提示：[R986](025-ui-localization-coverage.md#r986-world-identity-labels)，[034性別契約](../spec/034-world-profile-sex.md)。
選人WHO游標的正式中文與正常移動／取消：[R988](025-ui-localization-coverage.md#r988-who-cursor-wiring)，[033契約](../spec/033-world-identity-labels.md#r988-who-cursor-contract)。
頭像頁兩則文字與保持原頭像返回：[R955](025-ui-localization-coverage.md#r955-頭像頁文字與完成返回)。
低倍率人物數值與原生字型排版原型：[R953](025-ui-localization-coverage.md#r953-兩倍人物欄位原型)。
定版後的動態背景守門與等級欄修正：[R954](025-ui-localization-coverage.md#r954-動態背景與等級欄修正)。

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
| [018 發行字型](018-distributable-font.md) | GNU Unifont 字元覆蓋、原生尺寸與排版探針 |
| [019 譯文來源](019-catalog-migration.md) | 重複鍵、跨列對齊與TSV遷移證據 |
| [020 前端聲音](020-player-audio.md) | 正常啟動OPL事件與播放器接線證據 |
| [021 複合物品名稱](021-compound-item-output.md) | 三組兩行名稱及限定正常上屏驗收；另一bitmap背景未量 |
| [022 輸入錄放](022-input-recording-boundary.md) | 前端實際API順序、Step與錄放草稿界線 |
| [023 adapter 分層](023-adapter-boundary.md) | 現行模組依賴、未提交 SDK 來源與乾淨提交編譯對照；尚未定案搬移 |
| [024 隊員 HP／MP](024-character-vitals.md) | 第一位隊員HP顯示與正常戰鬥10→6回寫；MP、暫存有效期及作弊實作未完成 |
| [025 中文介面覆蓋](025-ui-localization-coverage.md) | 按鈕、訊息及圖片字樣的正式畫面與操作驗收矩陣；[R960戰鬥控制](025-ui-localization-coverage.md#r960-battle-control-centering)、[R966人物兩倍重排](025-ui-localization-coverage.md#r966-profile-reflow) |
| [M3 早期筆記](m3-notes.md)、[後續筆記](m3-notes2.md) | 歷史線索，不能取代後續資料庫或量測證據 |
