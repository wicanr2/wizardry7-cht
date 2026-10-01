# 022：前端輸入錄放的呼叫邊界

日期2026-10-01；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)；
工作入口[M2](../goal/M2-playable-frontend.md)，規格[013草稿](../spec/013-input-recording.md)。
本節是現有Go程式的唯讀查證，尚無本案正式錄放實作或三次重播收據。

## 輸入與定位

初次靜態查證時，本案已推送`d466750`、本機複合名稱實作尚待獨立驗收；目前已驗並推送`2d22364`，現況見[CONTEXT](../../CONTEXT.md)。dosgolem基線
`d16f03a2d8d782024d470606e1f4cb5272b89cd6`，既有`oracle.go`觀測欄位保持。
Go1.24.13、Python3.11.2；Docker image
`sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
此處定位為Go檔案與方法，不是IDA位址；本輪沒有反組譯或原版記憶體寫入。

| 檔案 | 唯讀查證時SHA-256 |
|---|---|
| `cmd/wiz7/main.go` | d1fde8c779596513fb57ef15e26f0dca2391099d96ae6feb4d827f6975e75b9e |
| `worktrees/dosgolem/oracle/oracle.go` | b68331ee851244a35464671a5036d95f585ea6b08720a017ce09b586f052ea8c |
| `worktrees/dosgolem/oracle/input.go` | 4fe6a8222d0f66c8f5b19829aedc1f13e715fe90f3d8ac115fc8440dc56a8991 |
| `worktrees/dosgolem/oracle/run.go` | 80ec43b5a65676cf715101dfb10631874db62fb11402588311a08a4f32a12a4b |
| `worktrees/dosgolem/oracle/state.go` | 3c756aa9703c866793b8b6e1f810918963079046b9580d4de2130b40560ec6fb |

## confirmed：實際輸入順序

1. `game.Update`先處理F1；說明開啟時丟棄字元並返回，其他輸入不送原版。
   F2只切換地圖面板。兩個前端保留鍵不應成為原版重播指令。
2. 一般按鍵先排`sc`、`sc|0x80`至共享佇列；機器goroutine逐項呼叫
   `o.PressKey(oracle.Key(sc))`。底層`PressKey`對每次呼叫再排make及break。
   錄放應保存現有API呼叫值與順序，不能在這項工作順便重解讀或修正鍵盤語意。
3. 可列印ASCII及Enter的CR排至另一佇列；機器一次`o.Type(string(sh.stdin))`交付。
   `Type`追加DOS handle 0輸入，不等同IRQ1或BIOS鍵盤事件。
4. 機器先交付鍵盤、文字、變動的滑鼠位置，再逐項執行左鍵`o.Click(x,y)`；
   每次Click返回後呼叫snapshot。之後才交付右鍵PressButton／ReleaseButton，再Run200k及snapshot。
5. MoveMouse只在座標變化時呼叫；多送一次也可能改原版故事的活動計時。
   不能為方便重播每塊重新送游標，也不能只記最後一次座標而失去順序。

## confirmed：步數及巨集

`Run(n)`由RunUntil逐指令推進，回傳錯誤時可能已推進部分步數。
`Click`可能先等原版設定游標，之後MoveMouse、hover、按下、hold、放開、settle；
預設路線通常6M步，但不能把所有Click硬編為6M，錯誤返回也可能已改狀態。
錄放至少核對每次呼叫前後Step及錯誤類別，不能用下一筆事件覆蓋不一致。

snapshot不寫原版機器，但會更新覆繪裁決、事件生命週期與地圖狀態；只重播原版輸入並在結尾畫一次，
不足以證明同一組live合成畫面。記錄原有Run／Click返回後的snapshot邊界，
不新增或合併snapshot，是本案草稿的驗收方向。

`Oracle.Save`為記憶體內快照，沒有本版公開落地序列化介面；Restore亦不還原開檔位置。
因此不能寫「已提供完整存檔」或以此承諾任意中途恢復錄放。
第一條錄放垂直鏈從相同原版工作副本冷啟動；即時存檔另依證據補規格。

## 別案參考與停止線

唯讀參考`/home/anr2/cht/psychic-war`版本`527456b10c60104926aaba97d0e3b522d41c4425`：
`docs/spec/019-input-recording.md`為READY，SHA-256
`3d46bad145ee93da836cdad122dafd16d1fc209a1ed3f9c71479ef35b6fd6701`；
`apps/psychicwar/record.go` SHA-256
`f1d31a601246c4fdba236ef7d3599c47d644e15cc447e32cb74519b131a95702`。
可沿用以指令數記錄、執行檔驗hash、前端保留鍵不進遊戲的原則；其鍵名／down格式
未涵蓋本案DOS文字與整包Click，不能直接套用而改變巫術7輸入。

錄製輸入可能含玩家依自己手冊輸入的防拷資料；只留玩家指定的本機檔案，
不將payload寫入log、文件、issue、測試名稱或可散布範例。
規格仍DRAFT，待本案呼叫序列隔離探針與獨立審查；別案READY不代替本案驗收。

R732唯讀審查另確認：目前F1要求由GUI提出，機器只在下一輪安全邊界承認，
整包Click及已抽出的輸入不會在要求瞬間中斷。PressKey／Type的Step是API交付時點，
不是原版真正讀取輸入的時點。舊描述的語意訂正保留於WORKLOG，來源SHA不改。

## R739：隔離格式原型

本機`workplace/r739/report.json` SHA-256
`7b65da53f00e53c1f9a7dec29f23c9c1e2daf6be66c5de891b84fc2a05153a2d`；
只有自製資料與命令fixture，沒有原版或GUI接線。Go1.24.13、同既有建置image，
`go test -race -count=1 -v`為42項含群組、失敗0；這不是42個獨立玩家驗收。

候選採JSON逐列：header、嚴格seq及前後Step的API命令、完整內容SHA與結束footer。
保留原樣掃描碼及DOS位元組；原測試覆蓋未知／重複／缺少欄位、部分null與型別、範圍、序號、跨步、
損毀及未完成結尾。R744後續獨立測試揭露下述漏拒，不能由此首輪綠燈宣稱全部拒載；
原型整檔驗完才回傳命令，沒有輸入投遞API。
總容量64MiB、單列1MiB、事件百萬、每Type最多1,024 bytes、初始檔案最多4,096，
均是待審候選上限，未作正式玩家體驗承諾。完整248檔的自製清單已驗可容納。

原型以0600及排他建立`.partial`；Sync／Close成功後才在同檔案系統原子建立正式名稱，
不覆寫已有檔，失敗不發布正式名稱且`.partial`名稱拒載；內容另改名的限制見下文。同步與關閉故障由自製介面注入，
未量真實磁碟損毀、Windows／macOS或斷電耐久性。
首版過小列上限、漏驗零值必要欄位及收尾錯誤後仍有footer的缺口已訂正；
首版來源與通過收據仍保留，不能由首版綠燈推論這些條件都已驗。
實際初始檔案日期／屬性、F1承認、有狀態快照、API錯誤分類、正常錄製及三次重播仍待查證。
013保持DRAFT，原型與完整錄放資料不提交。

## R744獨立審查與R752隔離修正

R744報告`workplace/r744/review.json` SHA-256
`169bd1a15ebdf83bd5d1b40d0d1d313e6e828b708c5bdc75013d2e48d2a94d7e`：
六阻擋、四應改、兩建議，其中四個原型實質缺口、兩個既列初始化／動態證據缺口。
獨立重算正文SHA的19案例含16契約斷言，確認九個漏拒：七種大小寫欄位別名、
錯kind額外null及非法UTF-8；舊測試仍綠燈。另證實NoResponse完整收據被Load全拒，
與I4／I7同類別、同Step續行契約矛盾。Sync／Close故障雖沒有正式名稱，
當時partial內容改名仍可解析，不能宣稱內容本身已失效。

R752新隔離副本報告`workplace/r752/report.json` SHA-256
`f4f13f33fc4f53cddf7b0abd040e6c2de0b08744a879ff64aeaa53385ee3e3f9`，
來源`record.go` SHA-256 `c985afad3753aa5d4e20180dff8bfe1cfb4fe3ed120e70062c1a66f1328ad20f`。
每層物件依JSON tag逐字核對，拒別名及明示null；解碼前驗UTF-8與surrogate配對，
不默默替換。合法錯誤收據可整檔解析，footer accepted必須符合全部事件錯誤分類；
false表示不具成功驗收資格，不阻止結構解析。真正重播仍須逐筆驗原版結果。

完成前先Sync正文；Finish／發布失敗盡力在同inode截去footer，正式名稱不發布。
獨立四種注入I/O失敗改名內容皆拒載，NoResponse後snapshot／run可解析且資格false；
19案例契約失敗0，新增合法／非法Unicode案例及競爭檢查通過，最新1.297秒。
先前舊NoResponse測試假設失敗、第二輪收尾漏拒與第三／四輪結果各自保留。
這不提供無簽章防篡改、真正磁碟損毀／斷電耐久或路徑競態保證。

Type寫出base64；原型仍接受等價uint8陣列，1,024 bytes按解碼後計，正式格式未定。
完整LoadWith／前端配置、實際資料三階段核對、API交付、F1／Draw、音訊及有界停止、
正常錄製再三次重播與跨平台發布仍待驗。013保持DRAFT；本輪只修自製格式原型。
兩輪原型與原始失敗不改，輸出1000:1000，自有容器及root-owned／誤建`.md`目錄0。

## R755／R758／R760／R761：冷啟動中繼資料與陣列邊界

R755唯讀冷啟動候選報告`workplace/r755/report.json` SHA-256
`9504c3e4631b064890af81ed7aa5fb2d683b17fcbd57b35e4747cbed4bbe75f4`。
248檔封閉集合的相對路徑、大小、SHA及DOS日期／時間與屬性，在來源、保留修改時間的副本及
LoadWith後相同；只觀察時計／DOS純量，沒有Run、投遞輸入或匯出RAM。
另外保存59個正式Go／模組來源與3份實際xlate／flow／help資源摘要；
不存在的dosgolem go.sum明示exists=false。候選只驗UTC修改時間；DOS archive屬性0x20
來自本版AH=43h契約，不能由Unix mode猜測。Mouse儲存值0也不代表客端倍率0。

R758獨立審查`workplace/r758/review.json` SHA-256
`7b4650e00c4b81733739521d81ee7ac05f879f91767306c244dfda636221f488`：
原19案例修正通過，另15個邊界案例找到Type陣列內null會被解碼成零位元組。
R760新副本只讓Base64字串走原解碼路徑，陣列逐項拒null；報告SHA-256
`832e77d720f2b0772de1f46d24ad9913be6bdc652623a349f8f2a9471e0fff1a`。
原有19案例、新15邊界、兩組各6項I/O及8項位元組案例與競爭檢查通過；
舊失敗與前輪來源保持原樣，不以這些格式測試宣稱GUI競態通過。

R761獨立報告`workplace/r761/review.json` SHA-256
`8783069690d25b06c13926b8905d075e021f4e2160a14f33e0e48f7874d5fe53`。
以Python自製原始列並獨立重算正文SHA，21份文件及1個清理雙故障案例通過Go競爭檢查；
合法位元組陣列／Base64原樣，null拒載，完整錯誤續流程可解析但資格false。
清理雙故障時另名內容仍可能解析，符合盡力而為的界線，沒有防篡改承諾。
另獨立核對248檔、59來源、3資源及候選canonical SHA
`43f37791ff35cc5ea31bb8f42183cf462a540c6f8902c2dc904fa34f80f736fa`。
本輪沒有重跑LoadWith，不把原報告斷言充當新的獨立原版執行。

窄修正沒有新阻擋；完整CLI／初始化摘要與header綁定、啟動前拒絕流程、
實際API／GUI呼叫鏈與有界退出仍是兩項接線阻擋，013保持DRAFT。
R763私人前端串接已錄到八種API與F1承認後凍結。兩次短期限未量完的失敗保留；
延長驗證預算後的三次完整結果與獨立審查如下，不能把短期限當成產品死鎖證據。

## R763／R765：實際API錄製與三次重播

R763封存報告`workplace/r763/report.json` SHA-256
`5b6052502364efe950f1ed9c9a5aa1400a366c166a00f7a3c7973ac281639022`；
70項輸出索引SHA-256 `10d3d4bc274ea9eab3763ff6dc89c13383c31914158dc050a8321cc217a9b446`。
私人副本從預設DS.EXE冷啟動，4797筆事件涵蓋八種API，2384個快照至500M；
所有事件錯誤類別為none，完整收據資格true。F1在330.8M由owner承認，關閉前步數不變。
實際終畫面是主選單的讀檔介面；輸入API被錄到不等於隊員改名、保存或載入已成功。

首70秒只到285M，另一65秒觀察至290.6M時CPU仍runnable；沒有該時點鎖等待證據。
同二進位、配置與三CPU配額延長程序期限至140秒，三次皆退出0，耗時約89.1／90.5／98.1秒。
7152對快照的原版RAM、A000、色盤、暫存器、低解析合成、滑鼠及事件數均相同，
三次終態248檔及PNG也相同；未比較全部CPU／DOS內部欄位，不稱完整狀態序列化相同。
只改一筆441M的正常滑鼠X座標加1，第一差異在441.2M，295對快照及終態PNG不同，
終態檔案仍相同。負向測試證明比較能偵測輸入變化，未改規則、亂數或原版記憶體。

R765獨立審查報告SHA-256
`909300e10711146ea6ffd40a1b9c9a86bf1573a5c3ced91ab23728c8554fad78`。
四項自製API案例通過；另以原Draw函式的替身交錯重現舊畫面消費新擷取要求。
Draw的畫面、面板與要求分次取鎖，PNG未綁發布序號；預設路線的像素相同不能解除此缺口。
另一靜態缺口是LoadWith可選-exe，標頭卻固定雜湊root/DS.EXE。
尚缺同發布圖片確認與期限、深拷所有權、非預定EXE拒絕、觀察關閉的同輸入A/B、
有效flow競爭檢查與正常存讀檔；013保持DRAFT。原始輸入正文與圖像只留本機。
本輪輸出1000:1000、自有容器、root-owned與誤建`.md`目錄0。

## R772：發布綁定修正與三次實際重播

私人修正報告SHA-256 `e7dee1ab7fd99ca1a736f80f46db137bda5a55dad933d1a7ec5f407c41b8d90c`；
57項輸出索引SHA-256 `7715e73b719f87143506d9579a03c3009892fe49f116a9a7493fe03c04f6557c`。
單次取鎖擷取深拷發布畫面與前端面板；擷取要求綁發布序號，舊Draw不能消費新要求。
PNG成功Encode、Close並重讀SHA才回報同發布確認；私人確認期限為5秒，非實體裝置上界。
非root/DS.EXE或連結的-exe於LoadWith前拒絕，標頭雜湊實際檔案且profile納入選定檔名。
結構完成、成功資格與錯誤筆數分開，沒有把accepted=false的結構重播當成功驗收。

四類自製交錯、深拷、EXE拒絕與圖片I/O案例通過競爭檢查1.134秒；
這不是有聲真GUI或有效flow競爭驗收。新錄製同樣4797事件、2384快照至500M，F1凍結通過。
三次同二進位及配置重播耗時約126.0／116.1／127.4秒，均退出0；
7152對快照、終態248檔、PNG及發布確認全同。確認序號2384、Step500M，
FrontendSHA `0b2fc3d88806353766f7c90fa7cec57715313ae8ebad6f3dfc15e35598e56793`，
PNGSHA `52625d1747185954500bd8d736f83d8f6b293bd69f515559616ddcbaa74157b8`。
來源清單固定後無漂移；獨立審查另行，013保持DRAFT。
仍缺觀察關閉同輸入A/B、有效flow、正常存讀檔、提早關窗的owner有界收尾及009整合。
本輪沒有修改正式程式、增加CPU／裝置調參或輸出原版內容；輸出1000:1000、清理檢查皆0。

## R775：觀察關閉的同輸入對照

報告SHA-256 `f7bcaf68937e1d055466803a0c2d1bc845761a526cc4942d6996f890a8c6ba94`；
11項輸出索引SHA-256 `53b2bad116ad36bcb3430e8c005893d8aab4c7f6bfa9521b7ac7426660c31656`。
獨立小runner從乾淨248檔及相同EXE冷啟動，LoadWith仍Adlib=true，
實際CPUHz、CycleClock、IRQ0Base／Every、Ticks及開始Step與錄製header逐項核對。
不安裝覆繪、印字或VideoWrite觀察，只依R7724797筆原API及順序執行至500M。
2384對原版RAM、A000、色盤、暫存器及滑鼠雜湊全同，終態248檔含DOS日期／時間／屬性全同。
這補上該限定路線的觀察開／關對照；不比較中文合成，不稱全部CPU／DOS內部欄位相同。
沒有亂數判定驗收、重擲、規則修改、RAM輸出或完整狀態匯出。
首封存腳本錯用attr欄名，實際契約為attributes；修正後核對通過，原失敗保留，非原版檔案差異。
程序正常退出、來源與可重現封存腳本留本機，輸出1000:1000、清理檢查皆0。

## R773／R776：錯誤資格與逾期圖片確認

R773獨立報告SHA-256 `264d17fad36132b45f0483d6518fedb58215c7be724c119df71a5c902f48a912`；
43項輸出索引SHA-256 `717906a6be81568f2abc92061317e5d2602d34ac5a8e4762824fe1a183aa6205`。
抽取真Draw的發布交錯、字形／前後背景／地圖深拷、GUI快取隔離、圖片五階段故障及
LoadWith前EXE拒絕／設定綁定皆限定通過；獨立重算7152對重播及R775對照也相同。
R772的F1承認為330.6M，與歷史R763330.8M不同；各自暫停期間步數不變。
兩個新增自製失敗保留：有錯誤record的footer=false，completion卻為true／0；
圖片確認超過期限，非空ok仍被接受。NoResponse可依原Click後續到finish，
execution測試僅證摘要契約，原main的Run錯誤會直接返回，不能稱該實跑錯誤路徑已觸發。
有限seq／kind定位仍待修，有效flow、真正退出及正常存讀檔仍待驗。

R776新私人副本報告SHA-256 `fbdc0d3c15f7fa887de93739206f880546a7b4c7a022dd0c8b070ebe158d05f2`；
59項輸出索引SHA-256 `c00a208d16d82ca82980d2e0790545f75117ece394cf194471daf0ea342a9310`。
實際Append成功後計數錯誤，摘要資格與footer一致；保存圖片完成時點並拒絕晚到確認。
重播不符另保存seq、kind、期待／實際Step、階段及有限錯誤類，不含params／Console／Error()正文。
自製MZ機器與真正Writer Append／Finish／Load、晚到ok及既有故障案例通過競爭檢查1.065秒。
新record4797事件、八種API及2384發布；三次重播約66.7／82.0／68.0秒退出0，
7152對快照、終態檔案及同發布PNG確認相同。這不承諾阻塞I/O有5秒退出上界；
R777獨立複核另行，013仍DRAFT。來源無漂移，輸出1000:1000、清理檢查皆0。

## R777：錄放修正的獨立複核

報告SHA-256 `299cef81d364800db9dbec7088edbb649cd928b3cefa3dd05658ab81e1d53f1f`；
148項清單SHA-256 `9d9d1c6787f115b96ef9d3078a80f05c4ab012d2967b5feaca2899a85f904370`。
R776十三群組獨立案例通過競爭檢查1.091秒，包含真正Draw替身交錯、
深拷、錯誤後snapshot及續行、有限診斷、完成確認等於／晚於期限與缺失時點。
R773兩個原失敗解除；4797個八類事件、2384發布與7152對快照獨立核對無差，
二進位、標頭、配置、59項輸出及同發布PNG真正綁定一致。

這解除該500M路線的資格與確認缺口，不代表有效flow或完整退出完成；
競爭檢查本輪仍是替身，Flows=0。Create／Encode／Close／Read／WriteFile本身
可能阻塞，五秒僅為確認資格期限，不是整個退出的牆鐘上界。
有效敘事頁GUI另在R780，原版正常存讀檔及009整合仍待驗。
封存來源零漂移、輸出1000:1000、root檔案／誤建md目錄／自有容器0，013仍DRAFT。

## R780：有效敘事頁的真視窗競態驗證

報告SHA-256 `51ff0249efc8c3f07d7ae903e1c02042aa0b8c494f204058a22248e7ce61162d`；
644項本機清單SHA-256 `1e03b53a07f969a63a4dee7bf70a0d5fbc294f6f2dac26aaaca5e9ca39e4b3c1`。
以既有合法checkpoint接續3475M，Ebitengine真視窗及-race至3795M退出0、432.9秒；
1368發布，實際Draw含flow2071次，F1面板開啟期間步數固定3730600000並於同Step恢復。
中文野花旁白兩列清楚且完整，F1面板的頁碼與操作提示均在視窗內；
248個原版檔案SHA未變，沒有新state.Save、原版內容匯出或裝置播放。
最終PNG、發布1368／Step3795M及FrontendSHA確認一致，來源／二進位綁定及核心漂移為0；
race警告與panic為0，R784獨立複核另行。

這是私人暖GUI所有權證據，不計013冷啟動錄放完成：checkpoint未保存歷史Mouse.Sets，
所以明示使用既有NoCursorWait選項，未變更hover／按住／放開／settle。
早期未排除這項歷史等待，首Click以100M execution返回，無完成結尾，已停止並保存；
暖時鐘拒絕與離線modfile校驗檔缺失的工具失敗亦保留。
成功GUI這趟另有一次no_response依原契約續行，因此錄放資格false；
不把送入點擊當成與舊定步數路線完全相同。暖header沿用R755來源索引的限制已明寫，
由實際binary／candidate-source／checkpoint另綁定，不能當正式I2完整配置。
輸出1000:1000、擁有權錯誤0、自有主代理容器已清理。

## R784／R786：有效敘事GUI複核與存檔入口探查

R784獨立報告SHA-256 `70da3db3a03be89ab9768b4117cebb4b3599f931ae19732d8419a562219df775`；
122項索引SHA-256 `316feb4217cc3f7bff3a7ac2b47ef155ac480edc47be15499ac3be7ffe1594aa`。
R780的2743事件、1368發布、327有效flow快照、2071次flow Draw、
F1同Step凍結／恢復、一次no_response使資格false及248檔完整中繼資料，獨立核對相符。
二進位`-race=true`已驗，終態Publication1368／Step3795M的PNG確認與雜湊一致。
敘事首次Draw圖只記FirstFlowStep3485200000，沒有該張Publication或同發布確認；
F1視窗圖亦無此確認。不能把終態確認外推到全部三張圖，不另重跑已足夠的長路線。
暖啟動／NoCursorWait與100共同核心檔加1私人入口的限制維持，正式013仍DRAFT。

R786報告SHA-256 `48cc3f9e48bc44cf63fef3b0a66e10a17e51c4c1b87bd8975185468cb689100a`；
260項本機索引SHA-256 `18dfb4e6c65dc242d4cd22658da9ab7ab9820f0c2785aaefabf5d12346a72757`。
沿用R780二進位與合法checkpoint，在真GUI以普通滑鼠點擊(160,15)，至3495M退出0、30.1秒。
終畫面仍迷宮，未開保存／載入選單；248檔SHA與輸入相同，PNG綁Publication71／Step3495M。
結構完成true、暖錄放資格false、錯誤事件0。此操作沒有完成原版保存或載入，
未重新印出的職業欄也不當中文覆繪驗收；沒有新state匯出或正式程式改動。
後續需查原版實際DISK入口，不能把啟動LOAD SAVEGAME介面當成存讀檔成功。
輸出1000:1000，自有主代理容器、root-owned與誤建md目錄皆0。

## R802：存檔落盤反例與正常保存證據

報告SHA-256 `67d57de5efcb2b16a64a524e0e00a0a23dff5558061ee9c9a231ea24def19179`；
970項本機索引`a775d1fdb01f4ad4c60bea180c747d837e46c68b07a267dcc464cf25ed1f2427`。
以既有R733合法暖checkpoint及私人靜音GUI確認實際DISK入口(212,166)，
再SAVE & RESUME(160,73)與SAVE(100,170)。Scratch空的對照只回成功，沒有實檔；
這證實正式main未設定寫檔層是玩家可見缺陷，不能由選單返回宣稱保存成功。

私人v3在LoadWith／既存checkpoint載入之後、原版Run之前SetScratch獨立輸出，
原版74次DOS寫入生成SAVEGAME.DBS 11404 bytes，SHA-256
`cd3ac3a0bb8f53e5e0c9fe6eebf5f2010be5ab8f63412bc7ffaa1ec11a508738`。
保存位置`[0,11,5,3,61,91]`，正常前進後`[0,11,6,3,62,91]`；
原版Root248檔初末全同，沒有既存同名檔，所以覆寫YES分支未量。
暖fixture的NoCursorWait、無OPL listener／裝置與資格false維持，未做競爭檢查，
不當正式音訊、冷錄放或完整存讀檔完成。

冷啟動v5以正式Click且不用NoCursorWait觸發LOAD，先到原版防拷，立即停止該路線。
沒有讀答案、推導或投遞Type，尚未到檔名介面或還原隊員／位置。
R763／R776的既存錄放僅有Type的kind／seq／Step安全中繼資料，
沒有題目身份／初態匹配或防拷成功收據；accepted=true不等於防拷通過，不盲套舊負載。
首缺/w handle、暖時鐘未初始化、暖cursor等待、保存未落盤與cold沿暖選項等失敗均保留。
完整存檔與圖片留本機；有界driver收尾不是正式原版Quit／A7驗收。

`SetScratch`的既有通用契約適用首次建檔；`AllowFileWrites`僅Root既存basename，
不解fresh首次SAVEGAME建立。本案正式修正見[014草稿](../spec/014-original-save-storage.md)，
Root維持唯讀、預設存檔位置待使用者選定，014與013均未升READY。
來源漂移0，UID/GID1000，自有容器、root-owned及誤建md目錄0。

## R814：錄放與正式音訊的契約核對

唯讀獨立審查4阻擋、3應改、2建議；逐條行號與來源在本機review.json，
SHA-256 `a10682f83ca908bfd73dbecc683c78054cc62ca6e11c55258a70e937e5af7c8f`。
來源索引`06ef1b2c3df85f919db12a71dc3de62f63f6aea668a822ccf4ce64912cfb655e`，
輸出索引`0865ad072959b3e772df02501f9c27594a10e5570efeeea40d1606e3a7acfb5f`。
沒有新跑fixture、原版或裝置；現有靜態契約已足以定位，不重開格式／晚到確認／panic測試。

四缺口為正式schema／CLI／profile定稿、新音訊配置與Root＋Scratch閉集合、
不改原API返回的音訊故障資格、以及Writer／Draw／owner／Close共用退出期限的末註順序。
正式009已READY且有深拷與5秒監督，不再用舊前端沒有有界owner的文字阻擋。
R780暖有效flow、R802私人Scratch保存與冷LOAD停在防拷的界線維持；
真冷含音重播與存讀還原仍是實作後驗收，不需新增遊戲語意逆向。
來源漂移0，輸出UID/GID1000，自有容器／root-owned／誤建md目錄0。

## R815：原型精確契約及正式決策前沿

2026-10-02私人候選`workplace/r815/contract-candidate.md` SHA-256
`aa88c1e990bc8329148618968e0842cbb6fa7845e0f2d0d49cf6fa8fb92f88db`；
報告`35bd124790bd5852df776897117579e381736f2217f5963943b9551c19e0ecd8`；
來源索引`abef1590cf1dc6d87e4d11b669e640f6c0859d8734cf69ed55dfae696519a1f7`；
輸出索引`df0e0fcded5cdc35f57f098ad4f22bb3c2e43e9a7a4efc2404b5066809df6a5d`。
主代理核對固定來源及候選後，將原型的精確欄位／編碼、八命令、有限錯誤碼、
容量和既有SHA公式回填[013 DRAFT](../spec/013-input-recording.md)，不再重開已完成格式探針。

已證：R760 Abort不寫末註、不發布正式錄放名稱；唯一新增自製fixture在合法Run的
原錯誤碼none後直接呼叫Abort，保留none，partial與改名後缺末註皆拒載。
observer_fault為自製標籤，未接實際音訊或觀察故障回呼。
沒有新原版、GUI、格式大表、裝置或硬體研究；不當冷錄放或可玩完成。

工程草案：Root／Scratch兩個初始閉集合、新音訊精確H／K、目前來源及有效初態
必須共同綁定啟動配置。舊R755的59項來源與3資源不足；有效初態純量仍需唯一結構來源。
R776先Finish後Draw的舊順序不能轉入正式路徑；候選先Draw、owner／Stop／Close，
再於同一退出期限內核對資格並完成Writer，底層I/O不能冒稱可取消。

真正待定的是014預設位置、正式格式／配置內嵌或必要附檔、工作階段故障文件策略，
以及正常停止錄製的操作入口。僅關窗Abort無法達成已要求的正常錄製與三次重播；
候選沒有選定方案、沒有正式修改或升READY。後續獨立審查不得把候選當使用者決定。

R818唯讀獨立審查完成文字訂正後，新增阻擋0、未解除應改0；僅批准保留DRAFT整理。
報告SHA-256 `021f471413c6d3e2d314acf417f1b11e33b830c9a9c828589f260962e815d7c5`；
來源索引`e59d45dd55f829824f1a3dad1753969368a20c32c265147f53573a88978b4017`；
5項輸出索引`f8d64bf23a125b2b4aedc0b6e00b65368a4c828eb8614cc6c891697e806d86fe`。
路徑拒絕精確限定C0，未冒稱DEL／C1；Abort直接呼叫與正常finish期限起點限制已明列。
R815的11項輸出全同、9項私人／標準庫來源無漂移；第10項013因本輪回填而變更，
原版、首次回填與文字訂正三個文件SHA均保留，不當原型漂移。
沒有重跑fixture或讀取原版／輸入／state／RAM；首次腳本縮排工具FAIL留存。
輸出UID/GID1000、自有容器及root-owned／誤建md目錄0，規格仍DRAFT。
