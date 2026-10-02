# 目前狀態與交接

核對日期：2026-10-02。已推送正式音訊與停止事件計數基線：`f4871c2`，獨立審查與前端競爭檢查通過。一般填槽缺字修正`caeac31`已推送，#44關閉且回讀確認；015仍READY、#28仍OPEN。錄放起始設定候選`cb18cd5`已推送，013仍DRAFT；最新提交以git log核對。本機 dosgolem 基線：`d16f03a2d8d782024d470606e1f4cb5272b89cd6`，隔離分支`wiz7-cht-opl-stream`。
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
| 一般填槽缺字回退 | [spec 015](docs/spec/015-generic-slot-expansion.md)、[014](docs/re/014-overlay-evidence.md)：所需數字／序數字模缺項或空時整則保留英文，避免略過字元後畫出部分中文。三群組16配置、完整前端競爭回歸與1980個完整字模修前修後對照通過，獨立實作複核阻擋0。玩家缺字分支與新原版A/B未量，維持READY；[#44](https://github.com/wicanr2/wizardry7-cht/issues/44)限定缺陷已CLOSED，#28整體仍待驗收。 |
| 控制碼保留與流文可用性 | [spec 016](docs/spec/016-catalog-control-preservation.md)：10列與有效產生器同步補碼／去高亮；顯示副本先剝低位碼，再驗Big5及字模。雙後端各8,825列、13邊界案例通過，一般烘製逐位元同；六筆可用值變化尚無正常上屏收據，維持READY。指定控制碼差異32→26，其他槽及十一組門檻保持。 |
| 動態主體來源與文意校正 | [spec 017](docs/spec/017-dynamic-subject-source.md)：11列與有效產生器同步，指定控制碼差異26→15；其他8,814列、英文與欄位保持。雙字型完整一般烘製只有兩個批准鍵改變，原始字模各54碼點獨立核對、缺字0；流文只有17320段落行0→13字、訊息17→13字。十列主體仍未顯示，沒有新正常上屏或原版A/B，017維持READY。 |
| 提示的特殊數字槽 | [spec 018](docs/spec/018-hint-numeric-slots.md)：1053保留原始`^`／`%`，完整英文鍵限定映射到現有數字槽，指定控制碼差異15→14。兩字型正式完整一般烘製與修前逐位元相同、8,825列流文同；槽位置1／4、13格保持。規格／實作審查無未解除項目，新正常上屏與原版A/B未量，018仍READY。 |

## 收尾尚缺的證據與交付

- 全文來源仍是 JSON 與訊息 ID；TSV＋英文鍵遷移見 [#24](https://github.com/wicanr2/wizardry7-cht/issues/24)。譯文覆蓋與正常路徑觸發對帳見 [#22](https://github.com/wicanr2/wizardry7-cht/issues/22)，資料已翻譯不能當作已上屏。
- 發行字型已完成實作與上述驗收（[#25](https://github.com/wicanr2/wizardry7-cht/issues/25)）；倚天字型仍只可本機使用，Unifont子集採OFL-1.1並保留上游全文。正式封包與三平台可啟動性、存讀檔及通關收據仍未完成。
- 輸入錄放、音訊完成條件仍以 [#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[#18](https://github.com/wicanr2/wizardry7-cht/issues/18) 為準；輸入錄放已查證現有API邊界，見[研究022](docs/re/022-input-recording-boundary.md)及[013草稿](docs/spec/013-input-recording.md)，尚無正式實作。早期 OPL2 事件比較通過不等於前端聲音與播放對齊驗收完成。
- 作弊功能見 [#35](https://github.com/wicanr2/wizardry7-cht/issues/35)，尚未有正式驗收；不可用作弊建立中文化回歸收據。
- 規格覆蓋與 adapter 位置仍待 [#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#29](https://github.com/wicanr2/wizardry7-cht/issues/29) 驗收。既有六份 spec 並不代表所有玩家文字路徑都已規格化。
- adapter 現況分層已盤點，見 [研究023](docs/re/023-adapter-boundary.md) 與 [架構圖](docs/architecture.svg)。裸 SDK 提交缺 OPL 套件，但 README 的完整受控準備流程已在乾淨副本重建：13 檔符合固定清單、559 個未改基線檔一致，Linux 產物與現行來源逐位元相同。`f68a00c` 已推送，[＃29進度](https://github.com/wicanr2/wizardry7-cht/issues/29#issuecomment-5942950186)全文回讀相符；搬移方案未選定，議題保持開啟。此結果不代表全新工具鏈、試玩或三平台封包通過。

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
正式串流已按此定案實作並通過獨立數學／音高檢查，實際裝置與玩家路線驗收另列；不能把原研究固定H丟幀當成加速算法。
同日補定退出採快速淡出並結束，捨棄未播尾段並計數，不等待整段尾音。
Reader約5毫秒不等於裝置退出時間；R799正常版已確認221幀淡出／884 bytes交付及退出，實體裝置預讀與人耳仍未驗。

## 本輪交接

| 分支 | 目前真相與下一步 |
|---|---|
| 已完成正式實作 | 故事31頁、旁白6頁、雙字型、171組同義文案與三組複合物品名稱維持前表驗收範圍；不重開已完成項目。來源、失敗與勘誤留在各研究文件及WORKLOG。 |
| 音訊009 | READY，R790規格審查及R798實作審查阻擋0。正式接線含精確H、手動加速維持音高、2秒ring、4410分塊、整數淡化與單一5秒退出監督。R801冷啟動200k／260M原版RAM、CPU、時基、畫面、事件與PCM相同，260M為7835事件／3816817幀。R799四趟競爭版及一趟正常版真GUI正常退出；F1同Step恢復，正常有聲／靜音60M畫面及原版摘要同。正常版欠載3401、溢出0，停止2936幀留221淡出、捨2715，實交付884 bytes／ACK=true／EOF=false。關窗至程序退出0.0624秒，不當5毫秒裝置保證；全玩家音訊鏈／人耳未驗，不稱CONFORMED。R804[固定依賴](tools/docker/dosgolem-audio-manifest.json)與[準備入口](tools/docker/prepare-dosgolem.sh)獨立審查通過，未推dosgolem上游。見[020](docs/re/020-player-audio.md)。 |
| 錄放013 | 仍DRAFT。R776／R777三次7152對快照及R775觀察開／關2384對相同；R780／R784暖有效flow旁證資格false。R815已回填原型精確欄位、八命令、容量與SHA公式，不再重開格式探針。R819／R820補齊起始設定純量的固定型別與唯讀來源候選，自製競爭檢查及獨立審查通過；正式getter與owner生命週期綁定未移植，Step0不證明無輸入。新音訊／Root與Scratch完整啟動配置、正常結束入口與故障文件策略仍待定稿；舊R776先Finish後Draw不得轉入正式路徑，Writer必須沿009同一退出期限。正式冷含音錄放與存讀檔仍待實作後驗收。見[013](docs/spec/013-input-recording.md)及[022](docs/re/022-input-recording-boundary.md)。 |
| 原版存讀檔 | 正式前端缺SetScratch，原版顯示保存卻不落盤，為實際缺陷。R802私人靜音暖GUI已確認DISK(212,166)→SAVE & RESUME(160,73)→SAVE(100,170)，獨立Scratch寫出SAVEGAME.DBS 11404 bytes／74寫入；正常前進改變位置。冷啟動以正式Click進LOAD後停在原版防拷，沒有匹配合法輸入收據，未計載入還原成功。[014草稿](docs/spec/014-original-save-storage.md)要求Scratch與唯讀Root分離；預設存檔位置待使用者共同決策，不擅自採用。私人fixture資格false，不當正式冷音訊／錄放完成。 |
| 剩餘11組／TSV | 使用者要求全部補齊才交付，008維持DRAFT。171組／373列完成，11組仍缺實際情境或短句版面證據；排除空鍵／純槽後烘製25、執行期24組碰撞。R768已授權10出處±2鄰則、38唯一ID的窄讀取完成，只有角色候選；原文留本機，不擴讀其他八組鄰文。 |
| 合法場景路線 | R803正常組成六人隊伍；R805 START進入新防拷題，沒有匹配收據即停止。角色管理／TREON審閱107對A/B與2976檔次SHA同，動態IRQ0Every兩側7069，不誤稱初始時基改變。R813一次固定女性仙靈正常CREATE到職業，獎勵9、六候選無PSIONIC，姓名標籤不是NAME:；25對A/B與2480正常檔次SHA同，未保存、不重擲，該分支停止。11組新增仍0。見[019](docs/re/019-catalog-migration.md)。 |
| 加速與平台前置 | R808非競爭24M真GUI退出0，關窗至退出0.0645秒、60M原版與合成同正常版；欠載54757幀、溢出0，不稱流暢度或人耳通過。來源881筆／合成器865筆的停止差額16已明列，R812獨立審查解除診斷文字缺口。R810最新172項來源無漂移，Windows amd64可編譯、Wine命令列退出0；未驗Windows遊戲視窗或真機、macOS及正式封包。見[020](docs/re/020-player-audio.md)。 |
| Docker開發入口 | [R816](docs/goal/R816-frontend-docker-entry.md)、[#43](https://github.com/wicanr2/wizardry7-cht/issues/43)：工具驗收完成、6656d6b已推送，Issue已關閉回讀。tools/fe.sh只留主機控制面，Go／Xvfb及UID檢查在既有容器內；離線模組唯讀、原版輸入唯讀、名稱查詢與執行均有界。實際建置／序列CLI／Xvfb通過，命令失敗17、逾時137與單次TERM143如實保留，自有容器0。248原版SHA相同，Linux中間binary不當可玩包；沒有新環境工具鏈、原版遊玩或跨平台聲明。 |
| 正式交付 | 尚未建立正式dist-all版本、tag或Release；全部11組、TSV、音訊／錄放、正常存讀檔及平台封包仍須完成。#35、#36、#28、#29依遠端Issue繼續；#26保持擱置。 |
| 控制碼與情境續驗 | R827固定法師技能路線已停止，十一精確事件0，不循環建角資格／亂數。017與018限定來源校正及正式雙後端檢查通過，指定集合剩14筆，008仍DRAFT；十列主體顯示缺口另列，不因補碼而解除。後續依[019](docs/re/019-catalog-migration.md)及既有有限證據核對字面符號、未知槽與跨列情境，不擴讀原版或猜補語意。 |

已封存批次輸出1000:1000，自有容器、root-owned與誤建`.md`目錄均0；在跑批次另留清理收據。
字型重建入口在[font/README.md](font/README.md)，流文回歸入口在[tools/flow-check.sh](tools/flow-check.sh)，
兩者完整來源／驗證分別見[018](docs/re/018-distributable-font.md)、[017](docs/re/017-flow-text-evidence.md)。
所有工作負載在Docker執行；研究收據留在workplace，不當正式下載入口。
歷史快照[docs/worklist.json](docs/worklist.json)不再更新；未追蹤的R433 goal保持原樣，不混入提交。
