# 目前狀態與交接

核對日期：2026-10-01。本輪起點：`6fcb8d1`（加速音樂定案及錄放格式修正證據）；正式程式基線仍為`2d22364`，本機 dosgolem：`d16f03a2d8d782024d470606e1f4cb5272b89cd6`。本輪後續文件提交以git log核對。
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
| 本機／發行字型雙軌 | [spec 007](docs/spec/007-native-font.md)、[018](docs/re/018-distributable-font.md)：倚天原生15×15與GNU Unifont原生16×16；原字型驗收7,523／5,490／395項失敗0。訊息統一後所需2,872個非ASCII字元、缺字0；字元清單2,919項。F1完整折行及延續頁；未觸發的防禦分支維持READY。 |
| 物品名稱校正 | [spec 010](docs/spec/010-item-name-corrections.md)、[019](docs/re/019-catalog-migration.md)：28個名稱及兩個共用訊息同步正式表與產生器，重建逐位元相同；雙字型靜態64項、正常角色頁兩倍率1,366項失敗0。實跑校正鍵0／28，僅稱名稱校正及靜態驗證完成。 |
| 同義訊息文案統一 | [spec 011](docs/spec/011-message-wording-unification.md)、[019](docs/re/019-catalog-migration.md)：171組／373列同步，獨立驗收16,038項、失敗0；兩字型兩倍率31故事頁、6旁白頁、代表選單與F1通過。實際觸發修改訊息10筆，其餘363筆只完成靜態驗證；剩餘11組仍是交付門檻。 |
| 三組複合物品名稱 | [spec 012](docs/spec/012-compound-item-output.md)、[021](docs/re/021-compound-item-output.md)：皮胸甲、毛皮護腿、小圓盾在正常角色／物品頁通過兩字型兩倍率、470項獨立核對與代表Draw；原版A000與色盤相同。隔離負向58項通過，未量玩家故障分支保持READY。 |

## 收尾尚缺的證據與交付

- 全文來源仍是 JSON 與訊息 ID；TSV＋英文鍵遷移見 [#24](https://github.com/wicanr2/wizardry7-cht/issues/24)。譯文覆蓋與正常路徑觸發對帳見 [#22](https://github.com/wicanr2/wizardry7-cht/issues/22)，資料已翻譯不能當作已上屏。
- 發行字型已完成實作與上述驗收（[#25](https://github.com/wicanr2/wizardry7-cht/issues/25)）；倚天字型仍只可本機使用，Unifont子集採OFL-1.1並保留上游全文。正式封包與三平台可啟動性、存讀檔及通關收據仍未完成。
- 輸入錄放、音訊完成條件仍以 [#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[#18](https://github.com/wicanr2/wizardry7-cht/issues/18) 為準；輸入錄放已查證現有API邊界，見[研究022](docs/re/022-input-recording-boundary.md)及[013草稿](docs/spec/013-input-recording.md)，尚無正式實作。早期 OPL2 事件比較通過不等於前端聲音與播放對齊驗收完成。
- 作弊功能見 [#35](https://github.com/wicanr2/wizardry7-cht/issues/35)，尚未有正式驗收；不可用作弊建立中文化回歸收據。
- 規格覆蓋與 adapter 位置仍待 [#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#29](https://github.com/wicanr2/wizardry7-cht/issues/29) 驗收。既有六份 spec 並不代表所有玩家文字路徑都已規格化。

## 有效決定與限制

以 [AGENTS.md](AGENTS.md) 第 13 節為完整決定來源：放大畫布、不縮中文字；預設 3×；TSV 英文鍵；字型雙軌；防拷維持原版驗證；F1 說明、自動／全圖兩種地圖與作弊範圍。
防拷關鍵字資料移出版控的 [#26](https://github.com/wicanr2/wizardry7-cht/issues/26) 依 2026-09-29 指示擱置，不能自行處理。
對話關鍵字與防拷答案仍以英文輸入；不得推導、輸出或提交防拷答案。
儲存庫保持 private，原版與不可散布字型不進公開交付。
2026-10-01使用者選擇先逐組統一同原文不同譯文，再遷移TSV；不採保留最後寫入文案的方案。
2026-10-01使用者進一步定案：剩餘八組異義情境與三組短句全部補齊後才交付，
不採十一組保留英文先發版；這十一組是正式封包交付門檻。
2026-10-01使用者定案音訊／前端預設依原版虛擬時鐘同步，保留手動加速設定，
不採目前24M加速作預設再設計變速音樂；原版CPU／PIT時基保持原樣，音訊實作仍受009規格閘門。
同日補定手動加速的聲音：音樂節奏同步加快、維持音高，排除加速期間靜音；
此分支尚待原型與獨立驗收，不能把原研究固定H丟幀當成已實作的加速行為。
同日補定退出採快速淡出並結束，捨棄未播尾段並計數，不等待整段尾音。
Reader原型約5毫秒不等於裝置退出時間；009的有界停止、預讀與實際淡出交付仍待驗證。

## 本輪交接

[R672](docs/goal/R672-delivery-prerequisites.md)的字型實作與驗收通過；固定字型來源存於本機`workplace/fonts/unifont-17.0.05/`，字元清單及重建入口在[font](font/README.md)。完整收據r682、r687–r691、r694與來源雜湊見研究018；舊失敗收據保留，詳細勘誤追加WORKLOG。
TSV規格[008](docs/spec/008-english-key-tsv.md)仍為DRAFT；[011](docs/spec/011-message-wording-unification.md)的同義文案及限定回歸已CONFORMED，171組／373列同步正式表與有效產生器，重建逐位元相同，正常玩家畫面獨立驗收16,038項失敗0。排除三組短句超寬與八組情境歧異，剩餘原文11組；排除空鍵／純槽後烘製25／執行期24組碰撞，控制碼問題仍待核。盤點與界線見[019](docs/re/019-catalog-migration.md)。
物品原文417筆與原版SCENARIO逐筆相符；28筆校正已按READY010套用，兩個安全共用名稱一併同步。r714靜態與r715正常角色頁驗收通過，但不宣稱全部校正名稱已上屏。三個複合名稱已經r730獨立正常路線驗收：470項與五項代表回歸零失敗、四張Draw與HD全圖相同；未量防禦分支保持READY。LONGSWORD搬移及新背景仍待證據，未納入這三組實作。
前端音訊[009](docs/spec/009-player-audio.md)仍DRAFT：R718確認200k進度水位可解決整包Click的交付間隙，絕對節流與快照所有權仍只在隔離副本。R743首次播放門檻消除初始三次Read欠載，但整段仍欠載663,447／溢出649,475；同二進位的R750只將執行槽2→4，仍欠載72,265、停止未播78,990、實讀淡出0。停止環境調參，續作正式API、有界收尾及加速維持音高原型；原版雜湊、事件與PCM相同不代表裝置播放通過。R749未設定參數探針的14槽不能外推已有2槽設定的R743。舊GUI真文字靜音競爭檢查通過；新版本含音訊及敘事頁仍待驗。
輸入錄放[013](docs/spec/013-input-recording.md)仍DRAFT：R760格式修正經R761獨立審查通過；R755已量冷啟動中繼資料。R763私人前端錄到八種API，三次重播7152對快照無差，單筆滑鼠變更負向295對不同。範圍只到主選單讀檔介面，未驗證名字輸入或存讀檔成功；R765找到Draw發布交錯及實際-exe綁定缺口。新副本修正與獨立審查另行，研究及剩餘閘門見[022](docs/re/022-input-recording-boundary.md)。
R733正常檢查點已成功恢復；R747自然醒來及兩頁旁白完成。R759正常前進量到`[0,10,1,7,57,87]`→`[0,11,1,0,57,88]`接縫；R764續至`[0,11,1,3,57,91]`，沒有新NPC或目標文案。R768依使用者明示範圍讀10出處及各前後兩則、共38唯一ID，先排除防拷區段；只留下角色候選，未量中文上屏。原文與完整收據只留本機，八組異義詞的近鄰不在此次授權範圍；見[019](docs/re/019-catalog-migration.md)。R738的401–413仍是物品效果欄強推論；十一組門檻未解除，續查具體正常路線。
R766私人音訊修正採先驗批次再提交，舊拒絕留下Pending的缺口已在局部測試解除；R767自製精確動態相位與240566個取樣核對通過。完整OPL進度API、播放器、有界退出與有效flow仍待接線；原型與獨立審查界線見[020](docs/re/020-player-audio.md)，009保持DRAFT。
R771獨立複核上述原子性及動態相位通過；R772私人Draw發布／EXE綁定修正三次7152對無差，同發布PNG確認相同。R775觀察關閉對照2384對原版雜湊及終態248檔相同。R773另重現錯誤錄製摘要與late-ACK兩個缺口，R776新副本修正後完成三次7152對無差，R777獨立複核另行。R770真正進度API已有冷啟動監聽開／關對照，R778離線消費端整合另行；尚未移入正式程式。
已封存批次輸出1000:1000、自有容器／root-owned／誤建`.md`目錄0；在跑批次另留清理收據。R433仍不混入提交。

[R671](docs/goal/R671-story-overlay.md) 完成故事覆繪與獨立驗收，收據在`workplace/r674`；6M／8M全部31頁有效，旁白六頁回歸通過且56份原版狀態與`r670`全同。
實作`65e62ad`已推送；[#40](https://github.com/wicanr2/wizardry7-cht/issues/40)已關閉並回讀確認。後續仍以遠端Issue為工作清單，本輪不代表全遊戲收尾完成。
最後故事頁3M有前頁筆畫淡出殘留，靜態閘門保持英文，6M／8M消退；證據F26與規格訂正經唯讀審查無阻擋。正式判準未放寬。
前端覆繪、F1暫停、左鍵、T2／T4、地圖兩格與視窗F2均通過（`r675`、`r676`）。T6初次混跑差異為0，未採為通過；`r680`同工具乾淨重跑前進／轉向為157,563／159,975像素，控制組各279，通過。初次原因未定，完整失敗收據保留於r675。
舊工作清單 [docs/worklist.json](docs/worklist.json) 保留為歷史快照，不再更新；未追蹤的 R433 goal 不混入提交。
歷史基線`workplace/r670/`、原版探針`r673`已保留；[tools/flow-check.sh](tools/flow-check.sh) 要求新的`workplace/rNNN`，拒絕覆寫舊基準，所有工作負載在Docker執行。
