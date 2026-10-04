# 026 戰鬥控制的中文按鈕內排版

狀態：READY。2026-10-05；[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)。證據及雜湊見[025 R959](../re/025-ui-localization-coverage.md#r959-battle-control-display)，目標見[R960](../goal/R960-battle-control-centering.md)。獨立規格審查0／0／0，主代理於FINAL後完整回讀兩報告並核SHA採認；JSON `f37f204de9ae51161a8cf8545df4f020fbfef9f8ad6c95299b37570ab552c8ce`、Markdown `9446c3a45e89baa7a20309511abffd34feab78fd9f1a3f8a9a9bdf5e0f5d4fe6`。正常GUI及撤除尚未驗收，整份不升CONFORMED。

## 問題與已證實輸入

使用者要求完整中文按鈕、原生15／16像素字形及中文字在按鈕內。R959量到七個戰鬥控制，Unifont2倍各越出下緣2列。只沿020／021既有垂直置中公式修正，不縮字、不改X、原版按鈕或命中區，不新增產品選擇。

原版Archives版DS.EXE為69206 bytes，SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。dosgolem隔離SDK及現行resident印字入口0110:3895沿001，沒有新增攔截點。座標為原版320×200；R959三程序82摘要相同，末5453.1M。39×9外框含1像素邊框由原始indexed畫格及獨立審查確認。以下既有英文鍵與上屏座標、dy0、筆畫色號1為confirmed，入口為既有resident管線；同名其他座標不在範圍。

| 英文鍵 | 原印字x,y | 原按鈕x,y,w,h |
|---|---|---|
| FIGHT | 13,167 | 8,165,39,9 |
| PARRY | 13,177 | 8,175,39,9 |
| SPELL | 55,167 | 50,165,39,9 |
| USE | 61,177 | 50,175,39,9 |
| RUN | 241,167 | 230,165,39,9 |
| MOVE | 238,177 | 230,175,39,9 |
| EQUIP | 235,187 | 230,185,39,9 |

## 顯示契約

- HD與正式Draw沿既有buttonLabelY共用純顯示helper，新增七控制判斷。既有選單與對話框判斷保持。
- 只有scale2／3、shown同時包含以上七個精確英文／座標、dy0、fg1且各有兩個非空字模時才調整。全部字模高度一致、15或16列，每列不超出原生字寬、各字至少一個墨點，整串字格在各自原版外框X範圍內。缺訊號、字形不合法或超寬均保留既有位置，不補事件。非七控制事件保持。
- 字距保持height+2，X保持x0×scale，Y＝按鈕頂點×scale＋(9×scale-height+1)/2。2倍倚天上移2HD像素、Unifont上移3；3倍兩字型位置保持。首列倚天2倍Y332、Unifont2倍Y331、兩字型3倍Y501；第二／第三列相應加20／40或30／60。
- 只改HD／Draw目的畫筆。textEv、pendingPrint、字模、原版記憶體、離屏、A000、palette、按鈕圖片、滑鼠輸入／命中區、存檔及序列化不變。原版更新時沿001現行事件撤除，沒有新事件生命週期。
- 垂直鏈是原版印字→字模事件→compose判定上屏→helper→HD／Draw。原版規則、按鈕行為及存檔沿原版，沒有新typed資料格式。

## 驗收與停止線

1. 用獨立硬編七矩形核兩字高、兩倍率全字格在框內及原事件保持；核缺訊號、異鍵／座標／dy／色號、空／混合／越寬字模、其他倍率的回退。跑cmd/wiz7既有測試。
2. 現行正式來源及私人工具副本除helper與共用接線外保持。R959固定82點／78命令／67輸入／三筆舊NoResponse原封重播，無新玩家輸入。SDK1132、原版248、R943資源及舊封存盲SHA保持。
3. 用R959獨立原始字型與PNG解析方法，Unifont／倚天各2／3倍七控制原生筆畫及完整字格都在外框內，全圖差異限批准舊／新字格聯集。3倍全圖保持。記憶體、indexed、palette、暫存器及固定點摘要相同，Scratch保持。英文來源／compose不改。
4. 守門錯誤、原版狀態差異、字格越框或區外差異立即停止並保存失敗，不調seed、不用新輸入找通過畫面。
5. 私人固定HD證據不取代真正GUI、按鈕操作、動態撤除或音訊。未量防禦分支如實記錄；整份保持READY直到正常GUI與撤除證據齊全。不外推全戰鬥、圖片控制或全部中文化。

Docker沿wizardry7-go-ebiten映像及Go1.24.13／Python3.11.2，原版、SDK、字型及舊收據唯讀，輸出限workplace/r960。原版圖片、字模及完整收據不進Git。十一組、TSV、其他介面與交付門檻保持。

## 已量結果

R960正式helper與共用接線已完成，cmd/wiz7守門及既有回歸、正式建置通過。雙字型兩倍率四配置的七控制完整原生字格及筆畫越框0、全圖／區外差異與交疊0；各82原版摘要保持。獨立56字模實例及1,664,000個RGB像素核同。Unifont2倍與R959修前1,346像素變更只在七個舊／新字格，metadata及comp保持；Unifont3倍全PNG逐位元保持。倚天3倍缺舊版實跑PNG，只採舊Y與字模獨立重建保持，不升為新舊實跑對拍。

最終獨立結果審查0／0／0，主代理於FINAL後完整核讀採認；結果、SHA及入口見[025 R960](../re/025-ui-localization-coverage.md#r960-battle-control-centering)。真正GUI、按鈕操作與動態撤除未驗，防禦分支僅單測，整份維持READY。

## R961 GUI色號副本補充 READY

2026-10-05獨立唯讀審查0／0／0，主代理於FINAL後完整回讀採認。審查JSON SHA `d59ad36a70ab2fd10025b91827a1e995a8769d3973ca69a56e20565e6cea2339`，Markdown SHA `304ee0e08b068de0a0a7f1ad6e874ea061739776670b4763926e34fb62c427f5`。此補充READY只准修正欄位與回歸，不表示真正GUI驗收。

main.go的snapshot先呼叫cloneVisualEvents，再將副本傳給Draw。R961修正前的clone只保留fgRGB，遺漏fg；026守門讀fg1，因此HD來源事件置中，不表示修正前Draw副本也置中。本補充只修這個現有欄位傳遞，不新增攔截或資料格式。R960固定HD結果保持，正常玩家GUI完成尚未證實。

- 在cloneVisualEvents既有textEv建構中保留來源fg，沿現有s／座標／fgRGB及每個字模深複製。只複製純顯示數值，不傳原版、owner或可變字模slice的參照。其他不供GUI使用的生命週期欄位仍保持現行省略。
- 原事件、glyph、fgRGB與Draw流程不變，七控制守門和目的Y不改；其他畫面副本也保留其來源fg，不猜色號1或新字色。
- 新回歸必須經真正的cloneVisualEvents→buttonLabelY鏈，核雙字高、兩倍率的七控制Y與來源一致、fg逐項相同，並改動副本字模證實原字模保持。異色fg仍拒絕七控制置中，不能因補傳而放寬守門。跑cmd/wiz7既有回歸及正式建置。
- R961私人來源探針只重用R960已核metadata，不執行原版。來源／資料流結果與真正GUI／Draw收據分列；修完後的正常視窗、七按鈕輸入及撤除仍需後續實跑，不因單測升CONFORMED或覆寫R960歷史證據。
- 先經獨立唯讀審查到補充READY再改正式Go；輸出限workplace/r961、正式修正限上述欄位及必要回歸。不修改SDK、catalog、字型、原版資料或存檔，十一組與交付門檻保持。

## R961 固定畫格的真正Draw驗證 READY

2026-10-05獨立唯讀審查0／0／0，來源oracle與過時文字修正已解除。主代理於FINAL後完整回讀採認；JSON SHA `52890214840f99881fd9eb567f13c04ff4be04a8894ccc81ebb0e42e8c50f8a8`，Markdown SHA `73e295f66093cf33234d175b011c6cc11d313436b23ab1571c6bec72dafd4997`。

本次只消費R960四配置已封存的battle.json／comp／pal，不重新執行原版或重播輸入。私人Go測試由現行正式來源完整副本加一份測試入口組成，不修改正式Draw／Layout／cloneVisualEvents。以既有metadata建來源事件，色值直接取既有8位RGB色盤，經cloneVisualEvents發布至shared，再用真正Ebiten.RunGame視窗呼叫正式game.Draw。私人Update只負責在已完成一張Draw後以ebiten.Termination結束，不呼叫正式Update，不產生玩家輸入或遊戲狀態。

- 雙字型各2／3倍分開四程序，以正式Layout建立640×400／960×600畫布。Draw沿正式drawReq儲存實際目的PNG；私人入口只提供固定shared資料、game.screen及退出條件。F1、地圖、圖片與Flow均關閉，原收據Flow須為空。
- 每程序核來源事件及全部字模在clone前後保持，副本逐項fg／fgRGB／s／座標與字模相同，Draw後來源仍保持。字模副本的img快取可沿正式Draw建立，不取回或共享owner參照。
- 以R960獨立PNG解析及原生字型解析方法重建全部RGB。七控制依原生字型與硬編矩形／026公式核對，其餘286事件只核既有字模metadata，不宣稱新的原生完整性。全圖差異、區外差異、七控制越框及交疊均須0；另核Draw與R960固定HD的完整RGB相同，不拿HD當唯一oracle。
- 原版248檔、SDK1132檔、正式Go、R960封存168產物及原始字型只盲SHA前後核對，不讀MSG.DBS／防拷資料、原版或存檔的新正文，也不載入原版執行器。所有PNG與完整metadata留workplace/r961，source／source-immutability收據不含新原文。
- 離線既有wizardry7-go-ebiten容器，UID/GID1000:1000、2CPU／2GiB／96pids；每程序逾時30秒，整批有界240秒。Xvfb由Python finally終止並wait；失敗留收據、不重挑畫格。完整正式測試／建置保持；私人入口先經獨立審查READY才執行。

這是固定畫格的實際GUI Draw驗證，包含正式clone及GPU繪製，不包含正常owner、正式Update、滑鼠操作、動態撤除、音訊或冷啟動。026整份仍READY；不以此宣稱正常玩家路徑、全中文化或可交付封包。

R961已量：正式欄位及clone回歸通過，四個固定Draw程序退出0；56原生字模、1,664,000RGB像素獨立核同，七控制越框、交疊及區外差異0。最終獨立審查0／0／0經FINAL完整回讀採認；JSON SHA `1990547b16d121314e737e2013e6f18df5a450817e4db2b49b9b359a4020091c`，Markdown SHA `7cf62f123373d9d7b7bdfabd7362af119753e887d418c53724a64b414cbb029e`。完整入口及封存見[025 R961](../re/025-ui-localization-coverage.md#r961-gui-color-copy)，正常玩家GUI與撤除仍待驗。
