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

## R962 正常暖GUI的私人驗收契約 READY

入口：[R962](../goal/R962-battle-normal-gui.md)。正式來源保持82bbeff，私人副本只接以下前綴、數字守門與唯讀收據。原版／SDK／R960及R961封存結果唯讀。Unifont2倍先做一條正常視窗路線，其他字型／倍率不由本輪外推。

- 起點為R733 two-forward.state，SHA `beecc0c54302d6304628d60f385fb4f3f6748e77c6bc859dfc153f669e337e6c`；現行LoadWith後、SetScratch前載入既有合法checkpoint。checkpoint的DOS還原會重開原絕對handle路徑，原版同目錄同時唯讀掛/original及/w。沿R960 privateBattleAudit及fixed-prefix.json SHA `c6b9c95c7f511d7c985a4e01c6eb7eb3727115b2a10626850457f1862dbf2f37`，82摘要／78命令／67輸入／三筆舊NoResponse與末5453.1M必須相同。只有既有前綴使用NoCursorWait。核同後做正式snapshot，再進入私人靜音節流接線、原owner及Ebiten.RunGame；不早退、不啟用hdonly、autoplay或clicks-at。
- 已查隔離SDK opl2.NewStream要求InitialStep及ExistingWrites都0，oracle.InstallOPLProgress的WithOPLStream也限冷安裝；第一次暖音訊初始化因此拒絕，沒有進正常GUI。正式newPlayerAudio保持，不改冷安裝守門。私人privateWarmGUIPacer只接受mute、tempo1、Step5453.1M、CycleClock=false及IRQ0Base165000。建立空冷Stream只用PhaseAt公式，不Accept、不施加OPL或產生PCM；依正式audioClock.Pace同一有理相位取整公式，epoch減去起始相位，使暖起點的等待從0開始。用SDK已允許的非串流OPL進度觀察器在200k量子內只呼叫paceAt；AfterOperation、Pause及退出使用原playerAudio方法。原trace保留，觀察器不改CPU／PIT／RAM或投遞輸入。這只驗預設節拍的靜音暖GUI，不驗正式冷音訊初始化、來源串流退場、合成或裝置；診斷數字明列no_pcm，不冒稱音訊成功。
- 新滑鼠只能xdotool→正式Update→shared→原owner→o.Click，維持正常CursorWait。F1以真正按鍵暫停，等待owner「說明 暫停=true」的Step ACK；恢復時同樣等待false ACK。新的操作故障不重試，不改API或seed找成功結果。
- 沿R960的數字取訊息掛鉤，29目標／2500–2574任一命中就設停止旗標。既有前綴7000–7197只計數；新的GUI階段遇此範圍也數字停止，以排除正文。私人副本的六個既有OnCall callback先查旗標，包含讀取、字元、圖片、返回與性別掛鉤，正文及樹段讀取前停止；原樹段dump整塊停用。snapshot第一句及Run／Click返回後同樣先停，不發布新畫格。只保存數字ID、計數及有限失敗類別，不觀察目標目的緩衝或答案，不讀MSG.DBS、鄰文或改原版。
- owner在每次正常原版操作前後記錄Step、種類與有限結果；只允許Run、移動滑鼠及左鍵，本輪不送玩家文字／原版鍵或右鍵。正常Run仍為200k，不由私人收據改長度或時基。這是私人診斷，沒有成功錄製末註，不當013正式錄放。
- 每次snapshot在owner建立不可變private frame：實際Step、完整RAM的SHA、暫存器SHA、indexed／palette與comp副本、來源事件的座標／fg／字模／before／after及GUI副本純值。Draw在原本同一次shared鎖內取這份frame，PNG與frame綁同一次Draw；不把後來的HD或owner狀態冒充該畫格。來源與副本逐欄核同，字模／before／after都深複製；圖形快取與生命週期欄位不入收據。Flow只記數量，非空時不宣稱完整RGB重建。
- 唯讀捕捉watcher依序接受0到12各live／paused兩個獨占請求，最多26次，持shared鎖設定既有drawReq；不持有oracle或投遞原版輸入。請求寫完、關閉後以同目錄hard link原子公布，已有檔拒絕覆寫；命令另外以120秒有界完整JSON守門。Stage／Phase錯誤仍拒絕。watcher沿stop.done終止，結束以同一退出期限確認，不能留下程序。
- Python控制器先取得完整live PNG／frame，再F1暫停。ACK後另取paused底層實際Draw及packet，paused Step必須等於ACK。新privateBaseReq在Draw首次shared鎖內與RGBA／events／privateFrame同取且清空，於原drawFlow後、前端地圖及說明面板前ReadPixels保存PNG；不修改dst或產品繪製順序。主代理只判讀這張真正paused底層PNG，正常live最終Draw另留。ready明列兩份收據SHA、paused Phase及面板前範圍；兩時點Step分列，不稱RAM同狀態。view_image判讀並明示可見座標／目的後，才解除暫停送一次左鍵。最多12次新點擊、13組取樣；可提早正常WM_DELETE退出。每步等候120秒、正常Click ACK45秒、關窗退出10秒；暖前綴900秒、整批外層1800秒，所有等待工具每次不超60秒。Xvfb及遊戲用程序群finally清理，正常退出與強制清理分列。
- 七控制同時存在的實際Draw，依原始Unifont字模、硬編七矩形及026公式核全字格、筆畫及局部RGB，必須越框／交疊／差異0。其他事件只作metadata及已知位置重建；不猜新增子畫面的排版。每次取樣記錄原版狀態與操作，子畫面撤除先比精確事件家族／字格及實際indexed變更，GUI本身不當新雙側A/B。獨立結果審查後，仍依未量的七控制分支、兩字型兩倍率及A/B缺口維持READY，不升CONFORMED。

暫停底層Draw限定補審0／0／0於FINAL完整回讀後採認READY，報告workplace/r962/review-paused-base.{json,md}，JSON SHA `acdecf3600d4afecda61473cc2b43e2810b923354ddcd22e7e4d688657446b25`；只採契約及來源，未採認新GUI結果。

Docker沿wizardry7-go-ebiten固定映像與既有SDK overlay，離線UID/GID1000:1000、2CPU／2GiB／96pids，輸出只workplace/r962；原始輸入、字型、SDK與舊收據唯讀。正式Go、資源及完整原版檔前後盲SHA核同。完整原文、圖片與收據只留本機。2026-10-05修前審查兩項工具阻擋，修後限定補審0／0／0經FINAL完整回讀採認READY；報告workplace/r962/review-spec-v2.{json,md}，JSON SHA `126b5bef67157ef6410b52eef3e1b8c7df29de5ea1705e46bc52834be06dfa92`。暖音訊契約缺口回DRAFT後，靜音節流限定補審0／0／0重新READY，報告review-warm-pacer.{json,md}，JSON SHA `bd1279f7d73800fe9b4fff97751acfe8dfa84a25cbfdf90b43a5a7845ed41757`。實際GUI的初次live5455.7M及暫停5458.7M來源純值同，但indexed變更5987像素、RGBA不同，身分守門正常停止並關窗exit0，沒有新點擊。改採真正暫停底層Draw，舊收據保留且回DRAFT再補審。正式播放器與冷串流契約保持。
