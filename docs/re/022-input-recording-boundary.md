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

## R819／R820：起始設定純量候選與獨立審查

2026-10-02；僅讀固定Go來源及自製初始化資料，未讀原版、RAM、state或原版輸入。
本地dosgolem基線d16f03a、隔離分支wiz7-cht-opl-stream的既有變更以逐檔SHA綁定，
不是以基線commit冒稱完整來源。56份來源副本加4份對照來源全同；正式前端未修改。
沿用既有Docker image `sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`，
Go1.24.13、Python3.11.2；定位為Go來源檔與行號，不是原版IDA位址。

| 本機收據 | SHA-256 |
|---|---|
| r819/contract-candidate.md（訂正後） | `489b96fdbc9dca20aea8cb88838edfec67f25db113f753f543b7951f296a1516` |
| r819/report.json | `38fcf1382f9facd0f701af6ee47040c7ce3bc11c0ebb0dd07fa035f502a9b5b5` |
| r819/source-manifest.json | `1d2f85b8964dfca731d013e49d2c80689f417a5f4fd3faba4531cacc3c82e491` |
| r819/output-manifest.json（69項） | `291f0239f045527fc9544489120957f2416153523e02a4fb2c7922a850bdbf7f` |
| r820/report.json | `759991744f78a9033a654fdff02c86bc78c67eafd374d637977b64d7d69d4805` |
| r820/review.md | `d9a210538de5ce715ba505978c6b27c1d908280e5b5aa7ec2e0fc7d7c73bed28` |
| r820/source-manifest.json（71項） | `5d9ea043bbf337ead40900cf742d1d860db707d051105a512401ba28234ec6cc` |
| r820/output-manifest.json（6項） | `62bde993c516e6da7f1f10dfda3a9b7c74db378d1f825c3bfe8da95a99bf5d19` |

路徑均在`workplace/`。通用候選getter直接複製DOS設定值，遊戲配置組合只在自製測試；
精確欄位及Go序列化候選回填[013 DRAFT](../spec/013-input-recording.md)。
confirmed僅涵蓋來源及自製初始化設定；日期／RTC與時間模型明列來源政策，
不是實測中斷結果。stored_mouse_x_scale=0為自動，有效倍率與BDA未觀察。
pristine是單一擁有者生命週期聲明，Step0不是無輸入證據，正式綁定仍未寫入。

唯一自製群組`go test -race -run '^TestR819InitialScalarFixture$' -count=1 ./oracle`
通過（1.015秒），獨立手寫預設JSON及第二程式重算654 bytes同SHA
`be68ecaf157861ac3ae4391f0c7951061bf7769dcb3f9dbdac57c6ebb946e735`。
20次重讀、回傳值修改無別名、三邊界拒絕及自製覆寫通過；沒有整體RAM／CPU A/B、
GUI／裝置或原版三次重播。有限純量不取代完整啟動配置，B2及013仍未完成。

獨立審查阻擋0、文字應改1已解除、未解除應改0；僅批准DRAFT回填。
勘誤：日期來源實際為`internal/dos/int21.go:125–129`，舊候選誤列133–137，
日期語意原已正確；只修定位，未改Go或重跑測試。
舊候選`contract-before-review.md` SHA `15f27e1ead6f0dbf8b76c64fabd1e727ae5ef40400ef7415a5a09bd6fa12e6d3`，
舊67項索引`output-manifest-before-review.json` SHA
`e8ced53c963946a64e6e1a5278e72a4e0810b9f3ae41a144ae07a71b4f9991e0`保留；來源索引依快照回查。
主代理逐項核對69份候選輸出、60項來源、71項審查來源與6項輸出相符，
013原SHA `e6b387fb75a3c5138f4809a91674c0082534d39b219c28b9c2d288a1036dd976`
因本輪回填會變更，不能當成原型漂移。輸出UID/GID1000；檢查210609項的lstat衛生異常0，
自有分析容器均清理，不升READY或可玩交付。

<a id="r845-explicit-save-wiring"></a>

## R845：明示存檔目錄的啟動接線候選

2026-10-02，基線c30bd5e，目標[R845](../goal/R845-explicit-save-wiring.md)，
[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)仍OPEN。
沿SDK009 READY及014 DRAFT，候選只在workplace/r845；本輪沒有執行原版、讀state／RAM／MSG／手冊。

現行main副本只加明示save-dir旗標、Prepare有限錯誤分流及LoadWith後SetScratch。
SetScratch在Run、輸入及音訊播放器安裝之前；成功載入分支只有編譯證據，未動態執行。
守門與R807逐位元同，其他13前端及go.mod／sum保持；原型要求明示目錄，不代選預設。

| 自製CLI案例 | 實際結果 | 邊界 |
|---|---|---|
| 缺save-dir | explicit_save_required | LoadWith前拒絕 |
| 存檔目錄與Root重疊 | overlap | LoadWith前拒絕 |
| 存檔路徑是一般檔案 | mkdir | LoadWith前拒絕 |
| 存檔目錄不可寫 | directory_not_writable | LoadWith前拒絕 |
| 獨立合法fresh目錄 | Prepare成功，再拒載自製非MZ | 沒有到SetScratch |

五案例的stdout、stderr及退出碼收據已核對，自製來源初末SHA相同，fresh空目錄。
Xvfb僅供圖形套件初始化，沒有玩家視窗流程；finally收尾、退出0。沒有原版Run或輸入，
不把合法目錄的預檢正分支當原版保存／冷LOAD、GUI遊玩、音訊或錄放驗收。
R807十八案例只沿既有收據，不再重跑。UID／Nlink僅Linux，檢查後外部替換與跨平台未量。

85初來源保持，91最終來源及43產物凍結；55已索引SDK與53實際編譯來源分開核對同。
Go1.24.13、Ebitengine2.9.9、Linux amd64及本地SDK replace由實際build-metadata核對，
映像Id沿既有Go／Python工具鏈。根lstat224506項異常0，產物UID/GID1000:1000、自有容器0。

| 本機收據（workplace/r845/） | SHA-256 |
|---|---|
| report.json | `dc01b28251cc9075fe3942dd11f42dcd135f2f736ddcf05ef5452dc4bdd8b972` |
| final-source-index.json，91項 | `d52dde3a0b86c003bf4ea1ba577094d848b6a8d34a8c648fde04a226b030f27d` |
| final-output-index.json，43項 | `f0f1868354821f8e386b4da2fbe3c90700705dbfc1ce410cdf9551f943dd0b5e` |
| freeze-guard.json | `14f9d094bb12210aa1ad22a468b6df37d9142af4dcee360c73ab8a6023fd38a8` |
| review-r846-report.json | `da8dbaae24f43e8eb64b37aa86baae29ce48db827c6c77bcfe91d07f938156ae` |
| review-r846-report.md | `9ca021a26eabf68009d9023e0dfca23235d2a8284e6d0724d752e8e9f1abbc33` |

R846獨立核對初／末來源、輸出、編譯來源與五CLI，阻擋0／應改0／建議0。
主代理完整讀取兩報告，只採認限定DRAFT回填；014仍DRAFT，預設位置待答，正式程式未改。
原始失敗描述與freezer v1保留，不改凍結收據來驗文件；十一組及完整可玩門檻保持。

R846四文件文字補審阻擋0／應改0／建議0，主代理完整讀取JSON／Markdown。
根doc-guard確認四份批准回填、其餘87來源與43產物保持，六文件擁有權與入口通過。
封存只追加驗收文字，保留受審文件前綴，不改凍結索引或產品資格。
補審JSON SHA為`5ec295e1a9edff35932ba58b1805a34b61225fff6e65f3526db340396c4cc60a`；
Markdown SHA為`fb4ae1efc0e95cd14ba29cf7aeda01daee7fcfa92d2ce04115f0f5b2f8212e71`；
doc-guard SHA為`0a707e197833b70a28a4224c8b0d43df5bda1cc6160d3158af9993a0a10df1a4`。

<a id="r847-save-platform-metadata"></a>

## R847：平台存檔預檢候選

2026-10-02，基線61d6070；[goal](../goal/R847-save-platform-metadata.md)、
[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)。只使用自製檔案及公開平台API，
沒有原版執行、MSG、state、RAM、手冊或autoplay輸入讀取，正式Go、SDK、譯文與字型未變。

R807的Linux守門分成共用路徑、Unix及Windows後端。Unix沿UID、權限及Nlink；
Windows以檔案控制代碼核對SID、NumberOfLinks及重解析屬性，SID候選接受TokenUser或TokenOwner。
共用路徑在mkdir之前以os.SameFile核對實體祖先，新增EXE硬連結別名拒絕案例。
依據是Go1.24.13本機原始碼、x/sys v0.36.0及
[CreateFileW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew)、
[GetSecurityInfo](https://learn.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-getsecurityinfo)、
[SetSecurityInfo](https://learn.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-setsecurityinfo)。
位址空間為Go來源及OS API，不涉及原版反組譯位址；自製案例結果已證實，原生平台行為仍未量。

R848首審阻擋0／應改1／建議0。主代理核對SDK internal/dos/files.go的O_RDWR契約，
及Go syscall_windows.go的GENERIC_READ／GENERIC_WRITE映射後，修正候選既有檔案權限查詢。
Unix開啟O_RDWR後關閉；Windows採GENERIC_READ及GENERIC_WRITE、OPEN_EXISTING。
不讀寫內容、不建立或截斷。修正後結果如下。

| 驗證 | 結果與限制 |
|---|---|
| Linux自製案例 | 二十項通過，前十九項保持；只能寫不能讀的檔案在寫探針前拒絕 |
| Wine輔助Windows API | 十六項通過，結果與前批相同；三項符號連結及一項讀寫權限負案例未量 |
| ACL測試收尾 | 首批依路徑還原被拒，失敗收據保留；改以事先取得的控制代碼還原，退出0、自製內容不變 |
| 平台編譯 | Linux amd64、macOS amd64／arm64、Windows amd64可編譯；不代表原生GUI或遊戲可玩 |
| 不變性 | 凍結時二十六項輸入相同；文件回填後四份批准變動，其餘二十二項及十八項依賴保持。首批108產物保存在frozen-r848，僅六份候選檔有批准修正 |

Wine的os.Symlink回傳成功，但實際未建立連結；已在Prepare前驗證，不能據此判守門通過或失敗。
ACL負案例雖拒絕O_RDWR，較窄的寫入控制也被拒絕，無法區分修前修後，因此仍未量。
這些是驗證環境限制，不繼續調整案例來湊通過。索引工具讀取刻意不可讀的自製檔失敗，
保留失敗版，改索引候選、編譯產物及原始結果收據；不改測試權限，不重跑產品測試。

| 本機收據（workplace/r847/） | SHA-256 |
|---|---|
| report-rw.json | `f90ccbffca6bac2d265a4dc026946da071d9c3eaaed3483b7c7da6a1ecab9ae6` |
| final-output-index-rw.json | `b8b21a37c6fdc0a3b2ae9487266a614f2799c300064089b1591b0c716a725221` |
| freeze-guard-rw.json | `843fa5bf6ea9242b52d2d92fc00985cdbc10761002f2fa6d1fbf3b6eaa06d2a8` |

014仍DRAFT；沒有正式SetScratch接線、正常保存／冷載入、錄放或可玩封包證據增加。
存檔預設及錄放格式仍待答覆，十一組、TSV、人耳及三平台交付門檻保持。

R848修正及四文件補審阻擋0／應改0／建議0，主代理完整讀取兩報告。根護欄核對二十二項未變來源、四份批准回填與138產物保持；六文件擁有權及入口通過。補審JSON SHA為`7df829acacb6cfb576a88918590bace96b5f3df3bd4e2d7bcfacdb676158e27a`，Markdown SHA為`a0280bc7f638e775e31bd662a7cff09503997df9551c366cb22dd0184f4f6492`，doc-guard SHA為`0fe37638b4468f9c4e8df9cb032d921a1ec3fb57e8c6262e7a67e29df710ad67`。014保持DRAFT，封存只追加本驗收文字，不重跑產品或改凍結收據。


<a id="r853-save-startup-success"></a>

## R853：明示存檔層的成功冷啟動

2026-10-03，基線9a16708；[goal](../goal/R853-save-startup-success.md)、
[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)。
補齊R845尚未執行的成功LoadWith／SetScratch分支，正式Go、SDK、譯文與字型未改。
候選沿現行前端與R847修正後共用守門、Unix後端，使用明示獨立可寫目錄；沒有選定預設位置。
來源為Archives原版DS.EXE，SHA`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
工具沿Go1.24.13、Ebitengine2.9.9與既有離線非root Docker，位址基準為Go API及原版Step。

| 驗證 | 實測結果與資格 |
|---|---|
| 啟動順序 | Prepare成功、LoadWith成功後，SetScratch在Step0完成；之後才建立音訊及開始owner。對照側沒有接Scratch，其他觀察與停止控制相同 |
| 冷GUI | 英文、靜音、零玩家輸入及autoplay=false；兩側到主選單，取得250M的實際960×600 Draw與260M的320×200原始PNG，主代理完整目視 |
| 原版對拍 | 200k及260M兩點RAM、indexed、palette及暫存器相同；原始PNG精確相同，不能外推未量到的保存及載入 |
| 原版檔案 | 248個素材檔初末SHA全同，原版唯讀掛載；兩個Scratch實際沒有自然寫檔 |
| 來源與凍結 | 19正式／守門來源、53實際編譯SDK及42候選／腳本來源不變，62產物凍結；未把測試／忽略檔當成實際編譯來源 |

260M兩側只計數保護範圍的自然內部取訊息，各1次；觀察器不記ID值、不解碼或輸出答案，沒有代答、跳過或改判定。
固定整份冷初態與零輸入，沒有改RAM、存檔、原版CPU／PIT時基或識別／改寫seed。
此次不驗傷害或RNG規則，不把相同seed數字當作對拍證明。

原始收據與候選只留本機workplace/r853，入口report.json、inputs-final.json及frozen-index.json。
report SHA為`87b59087e6f482155290fe75db1a2cdfb78930722fafa36f8b7d3b6b47f8d741`；
frozen-index SHA為`bd4916ffbdfb46e0927bfdb48bf5da909ccc6610278cae87af982fb6da3804be`。
首次Xvfb就緒查詢缺xdpyinfo，原版尚未開始；改用映像既有xdotool後正常。
原型首個原版執行在200k之後被過嚴ID護欄中止，退出2及早期收據保留；
該護欄誤阻斷正常內部讀取，修為單純計數後沿同映像與命令重跑，不改原版流程或正文。
前兩項為環境／觀察器問題，不當作產品缺陷或保存失敗。

014仍DRAFT，成功候選不代表正式SetScratch已移植、正常保存／冷載入或完整可玩。
預設位置及錄放格式仍待答覆；十一組、TSV、人耳與三平台交付門檻保持。

<a id="r859-original-save-wiring"></a>

## R859：正式存檔接線與限定保存驗證

2026-10-03，基線fd85707；[goal](../goal/R859-original-save-wiring.md)、
[014 READY](../spec/014-original-save-storage.md)、[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)。
使用者已選定系統使用者資料目錄與單檔錄放檔頭。R860技術阻擋0，唯一歷史文字應改修正後才接入正式程式。

正式前端先解析及預檢目錄，LoadWith成功後立即SetScratch，再安裝觀察、音訊及啟動owner。
Linux採絕對XDG_DATA_HOME或HOME下.local/share，macOS採Application Support，Windows採KnownFolderPath的LocalAppData。
保留-save-dir；拒絕與原版雙向重疊、別名、非一般項目及不可讀寫來源，探針清理失敗也拒啟動。
所有原版格式、檔名、規則及防拷保持。正式既有Go只有main變動，另增十份存檔程式及測試；SDK與正式音訊未改。

| 驗證 | 實測結果與限制 |
|---|---|
| 自製候選 | 十三路徑、二十守門及一清理故障通過；候選四平台可編譯，未當原生平台GUI通過 |
| 正式測試及建置 | go test -race ./cmd/wiz7退出0，22.357秒；Linux正式產物可建置，版本及依賴另封存。兩項沿主代理工具收據，沒有另保存完整本機測試日誌 |
| 正式冷GUI | 英文、靜音、明示24M及零玩家輸入，到主選單正常退出；關窗至退出0.1099秒。沒有自然寫檔 |
| 冷對拍 | 私人接線開／關200k與260M的RAM SHA、indexed、palette及CPU相同；正式產物hdat及幀檔兩點與私人對照相同。正式run本身未輸出Step0原始觀察，不混稱同一收據 |
| 暖保存 | 同一R733合法3475M進度，固定DISK、SAVE & RESUME、SAVE正常選單輸入，中英文各寫11404 bytes並退出0；檔案逐位元相同，兩點原版FNV摘要及A000／色盤相同 |
| 原版與畫面 | 248素材檔初末SHA不變；已目視中文DISK五選項與保存後返回森林及隊員畫面。沒有截到檔名對話框 |

新存檔SHA為`cd3ac3a0bb8f53e5e0c9fe6eebf5f2010be5ab8f63412bc7ffaa1ec11a508738`，
與R802既有保存相同。新暖測試未量寫入次數，不沿用R802的74次作本輪計數。
暖工具僅載入既存checkpoint，使用NoCursorWait與靜音相對節拍；沒有新state、RAM或存檔注入。
正式編譯53份SDK來源與暖側共同來源相同，暖側另編checkpoint輔助及其引入的未修改state套件，共55份。
定位採Go來源／API與原版Step，非IDA線性位址；工具Go1.24.13、Ebitengine2.9.9及離線非root Docker。
原版Archives DS.EXE SHA為`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。

首批實際冷GUI以xdotool直接摧毀視窗而遇BadWindow，收據保留；改既有WM_DELETE正常關窗後退出0。
暖測試先因缺歷史/w別名、再因正式cold_required音訊守門於Run前停止；私人工具沿R802相對節拍處理。
v3產生器只改GUI點擊，漏掉固定步數點擊，仍等待checkpoint未保存的Mouse.Sets，沒有寫存檔。
v4只修正該測試分支，固定相同座標、步數及起點重跑成功；不調參數或更換亂數結果。
未安裝的保護段讀取計數欄已從私人觀察器移除，沒有虛構零讀取保證。
完整失敗及封存索引只留本機，160來源與2245產物owner均1000:1000。

| 本機收據（workplace/r859/） | SHA-256 |
|---|---|
| report.json | `ee2924472d3c7a2af86ae912c876806592f008aa140e09dcc27274fcd38abe2d` |
| frozen-index.json | `19c6555041c6418f5b9fac0632a14c7c798483d45b335fc5777d5d6bb5bac563` |
| R860規格審查JSON | `90a785e8c2eeb93a0a593e679bcf01903f04c2adfb8063342427809dd0ac9916` |

confirmed僅適用上表限定工程接線及暖保存支路，014保持READY。
冷正常保存與LOAD隊員／位置還原、原生平台、人耳、錄放、十一組、TSV及正式封包未完成，#42保持OPEN。

<a id="r862-recording-initial-config"></a>

## R862：完整起始設定的私人候選

2026-10-03，基線f78f050；[goal](../goal/R862-recording-initial-config.md)、
[013 DRAFT](../spec/013-input-recording.md)、[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
使用者已定案單檔內嵌完整設定，本輪組合B2技術候選，不重問檔案形式，不代選待答的故障策略。

來源採R859實際編譯的53份現行SDK檔與R819私人getter；正式來源唯讀，僅本輪副本增加getter。
候選保存完整LoadWith引數、純量、binary／來源、兩層檔案、音訊H／K及前端資源宣告，
型別與欄序見013及本機contract.md。沒有正式程式、SDK、譯文或字型變更。
定位採Go來源與API，工具Go1.24.13、Python3.11.2、既有離線非root Docker。
沒有原版LoadWith／Run、MSG、手冊、state或原版存檔正文讀取。

| 驗證 | 結果與界線 |
|---|---|
| Go自製競爭測試 | 五組通過，1.041秒；八種配置差異改變SHA、十五種非法DTO拒絕。未重跑既有21格式文件及R819getter純量表 |
| 檔案閉集合 | Root／Scratch同path不同內容均保留，辨識新增、刪除、被遮蔽Root內容及DOS日期變動；空Scratch合法、symlink拒絕 |
| 獨立Python | 按自製內容重算檔案SHA／DOS位元，獨算H／K／IPS位元及真實自製binary SHA；ASCII fixture重編碼相同不外推Unicode canonical契約 |
| 實際來源核對 | 十項音訊／LoadWith／預設倍率來源核對通過；倍率無2..6限制，候選只記錄原值 |
| 封存 | 83項來源及69項產物，owner1000:1000；包含已保留的驗證腳本首次失敗收據 |

自製配置2844 bytes，SHA為`2fb54a8ddf0395d17f3984a55f25208bb1f776469c34f6a90e217289ebf77c9d`。
H既約為12304689375/4096，24M的K為52428800/6562501；正常與手動加速的既定行為不變。
日期沿DOS來源的ModTime本地時間位元與0x20屬性，不改UTC政策或把Unix唯讀屬性當DOS屬性。
候選掃描要求獨占且不變來源，不保證抵抗並行換檔；所有自製內容明示不是原版替身。

| 本機收據（workplace/r862/） | SHA-256 |
|---|---|
| report.json | `161d4b9ab26564ca9a5fab10a174989cf818431a094cf9192ee5bab7a217c858` |
| contract.md | `d7d440fef04a0b52a67dcf6332781e76616c96de62a4627feabac766117f0da1` |
| source-index.json | `d63e08b3780cca7448fc9587ba7c7ad437f5ed029e4148fff5a2d0db703153ba` |
| frozen-index.json | `d92c6b4b9587825440fe0ea32d4c4a837e4bdcbd9d5e18f37b0b6b52e50b5259` |

首次有界Docker呼叫在進容器前受socket權限限制，依已授權隔離執行權限後正常。
Python核對首次把Reader常數查到Player，收據保留；修正定位後同映像重跑成功，候選及Go測試未改。
這些是環境／驗證腳本問題，沒有原版執行或產品缺陷結論。

confirmed僅涵蓋來源定位與自製測試，固定平台政策仍沿R819 source-bound policy。
Build只驗已具型別DTO，沒有新的wire decoder；pristine仍是呼叫端聲明。
正式編譯閉集合、選定EXE、資源載入結果與同一LoadWith實例尚未由owner綁定。
沒有GUI／裝置／原版三次重播或存讀還原，不解除B2全部閘門；013保持DRAFT，#17保持OPEN。

R863獨立審查阻擋0／應改0／建議0，主代理完整讀取兩報告。
審查另逐項重算83來源與69產物的SHA及大小，owner均1000:1000，沒有重跑產品或原版。
JSON SHA為`a7e3fa7cf69be90c4ad426f2bfc106d3542ed2cc6e97d2db0c50f5d2ba50ab66`，
Markdown SHA為`c7fb7ada3258119537d14b55ac11eeb9009c9829839c46a14359c60651a389ed`。
封存後只追加本審查紀錄，不覆寫凍結收據；根收尾入口workplace/r862/final-check.json。
工作根衛生異常0，自有R862／R863容器0。

<a id="r864-recording-startup-binding"></a>

## R864：起始設定與實際啟動的綁定

2026-10-03，基線cc0db32；[goal](../goal/R864-recording-startup-binding.md)、
[013 DRAFT](../spec/013-input-recording.md)、[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
私人擁有者自行LoadWith及SetScratch，完整設定封存前不交出Oracle。
沒有接受外部pristine聲明或既存指標，交出後即使Step0亦拒絕再次交出及重新取得配置。
實際EXE路徑／Options及同一次資源ReadFile的bytes、解析結果共同生成配置；正式前端與SDK未改。

來源沿現行53份SDK檔，私人副本只增加R819純量getter；資源解析規則沿現行main、flow.go與help.go。
go list實際選用54份SDK、15份前端／私人候選、1份配置Go及3份模組鎖定檔，共73項。
編譯前生成來源JSON／SHA，以ldflags嵌入同一候選binary；執行時核嵌入清單，研究模式可另核實際檔案。
第三方與stdlib不逐檔列入，依模組鎖定、Go建置中繼資料及實際binary SHA綁定。
工具Go1.24.13、Python3.11.2及既有離線非root Docker；定位採Go來源／API與Step，不是IDA線性位址。
原版Archives DS.EXE SHA為`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。

| 驗證 | 結果與資格 |
|---|---|
| 12個起始案例 | 控制、英文、中文、24M靜音、同配置與嵌入來源六個正案例；六個負案例均於交出Oracle前拒絕，沒有CPU指令執行或玩家輸入 |
| 初態 A/B | 各正案例Step／Ticks=0，全1MiB RAM、indexed、palette及暫存器摘要相同；擷取前後亦同，只涵蓋LoadWith初態 |
| 負向預核 | 靜音、EXE角色、Scratch、資源及編譯來源差異在LoadWith前拒絕；純量差異在載入後、交出前拒絕。沒有已寫配置殘留 |
| 無專案來源掛載 | 另1案例只掛候選、配置、資源與原版，未掛本專案／SDK Go來源或go.mod；配置與帶來源模式逐位元相同 |
| 原版不變性 | 248檔先與R859全部SHA核對，初末相同；只必要載入及盲雜湊，不解碼MSG、手冊、防拷或原版state，RAM位元組不保存 |
| 封存 | 76來源及170產物，owner1000:1000；原版EXE的本機外部副本排除正文審查與提交 |

新候選binary SHA為`429e82fcff584d460b264a5ff65692167f414e4e31dc7228266cd0cf83851da1`，
73項編譯來源JSON SHA為`73a54993c522d950ccc1ed1770beed92d0baf51a4f4e5c11c1d19cef56502d9a`。
24M靜音中文配置51321 bytes，SHA為`0d913736e97713958aaca060bb7f24a054477aff748a499b15b84746d54f3533`。
LoadWith呼叫有動態計數；Run及輸入沒有計數器，明列source-controlled load-only，不虛構零次動態觀察。
來源要求唯讀且獨占不變，沒有惡意並行換檔或內部單次EXE buffer hash保證。
可信自製expected解析不當正式wire安全證據，GUI／Draw與音訊裝置未驗。

| 本機收據（workplace/r864/） | SHA-256 |
|---|---|
| report.json | `68b549601957a0507f148feaf9c7c6f97e476d6ffa547ccf1140a0a5b541a432` |
| isolated-no-source/result.json | `a0b073b6dc58f453f0af974405fe0d4d5201dccfb0013485bc28fd1e669d6163` |
| contract.md | `6b8bb916504d8c2d5eb7d315c7483a405fa3f15a1aeb84472276fb4338027d3a` |
| source-index.json | `a9e5fa4bc16454c5347984a89d82d5bb1ce94b3675dc139e7fc8fd67a381ac72` |
| frozen-index.json | `e7d4b7a8137a7b0e2f42ced49892f946c021ec3731b1a5909d39911936453487` |

首次編譯誤用Palette的RGB struct欄位，實際為[256][3]uint8；匯入偵測又誤收flowPos字樣，均修正且保留失敗來源。
首次執行誤掛含6575檔的workplace/run研究目錄，248檔護欄在LoadWith／Xvfb前停止。
改用既有workplace/original/wizar7，先核全部基線，不放寬條件；前次只盲雜湊，不解碼研究資料。
首批11案例及原binary保存在before-baked；因修正外部Go來源依賴才重建、重跑12案例，沒有換亂數或挑結果。
隔離副本首次copyfile未保留執行權限，候選未啟動；核owner後改copy2，同容器及命令重跑通過。
前者為私人候選程式錯誤，後兩者為輸入選擇／產物準備問題，不當原版遊戲缺陷。

confirmed只涵蓋以上私人初始化及來源核對，固定日期／RTC政策仍沿R819 source-bound policy。
沒有正式錄放、原版指令、Draw、裝置、三次完整重播或存讀還原證據，013保持DRAFT，#17保持OPEN。

R865最終獨立審查阻擋0、應改0、建議0，主代理完整讀取兩報告。
審查逐項核對74項未變來源、2份批准文件變動、73項編譯輸入及170項凍結產物；
只採認以上私人初始化與六文件回填。JSON SHA-256為
`7ad888ba516e01cfa827d38dcf7e2bf22801d0bbd2a44832efcd99b4835e7a76`，Markdown為
`095b481fd5c6c83e4b9c5e907baa77db24f07c54296b23810f813a856eb296da`。
根收尾入口為`workplace/r864/final-check.json`；沒有重跑原版或覆寫凍結索引，013維持DRAFT。

<a id="r866-recording-shutdown-chain"></a>

## R866：成功收尾的同一期限與完成順序

2026-10-03，基線d2c5ee7；[goal](../goal/R866-recording-shutdown-chain.md)、
[013 DRAFT](../spec/013-input-recording.md)、[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
工具Go1.24.13、Python3.11.2、既有離線非root Docker；定位採Go函式與單調時間，沒有CPU／IDA位址。
沒有原版輸入或載入；自製header與事件不含原版內容，不能當原版替身或玩家收據。

來源為現行audio_player.go、audio_reader.go與R776的record.go。
前兩檔原樣複製，Writer只在私人副本加入同一絕對期限的FinishAt與返回後核對。
協調器以最終發布、Step及配置SHA確認Draw，通過後才允許owner完成；
同一期限內依序關閉音訊、判定工作階段資格、完成Writer。
沒有將工作階段故障偽造成原API錯誤，也沒有代選正式故障紀錄策略。

| 驗證 | 已量範圍 |
|---|---|
| 收尾鏈11案例 | 正確順序成功；錯誤發布／配置、Draw錯誤、晚到確認、已耗期限、Close錯誤、資格故障、Writer錯誤、owner未回報及共享期限皆無成功檔 |
| Writer4案例 | 正常完成可載入；過期起點、晚payload Sync及晚Close失敗，無成功名稱，清理成功時改名片段亦因缺末註拒載 |
| 正式串流與Reader | 使用真正Stream.Stop及Reader.Stop，自製裝置在F1暫停模型解除暫停後消費淡出；ACK=true、884 bytes，不當實體裝置或人耳確認 |
| 實際阻塞2案例 | 自製Close與Sync分別永不返回，五秒監督結束子程序；Close退出5.058秒、Sync退出5.146秒，沒有Writer成功或被阻塞Sync後的末註 |
| 來源及封存 | 136來源核對項、31凍結產物、6份私人Go來源。105項SDK包含go list列出的依賴測試檔，屬來源核對集合，不宣稱105項皆編譯或測試執行 |

18個限定案例通過，Go競爭檢查無失敗；兩個僅供子程序呼叫的測試入口在主suite顯示skip，
由實際子程序各執行一次，不把skip計為另一個成功案例。群組標題也不重複計數。
首批成功案例沒有實際Stream／Reader停止，新版另加入該成功鏈才重跑；舊來源及收據保存在before-reader-phase。
正式Go、SDK、catalog與字型未改；原版檔案未讀，沒有GUI／冷錄放／存讀還原或可玩完成聲明。

| 本機收據（workplace/r866/） | SHA-256 |
|---|---|
| report.json | `b2e095b53c86ef81a2bc32a53726786ebca21f6dc7b01286d2ca82488422c626` |
| source-index.json | `d8fafd627e48ae3d8b0d5b204ec46694b6818bae439d743547e32369f40bf5ce` |
| frozen-index.json | `5d33977f0024580f443597c7ec9dda794e5ccf651e3e5e04e144dcc547457381` |

confirmed只涵蓋上述私人協調器及自製通知／裝置，沒有真正Draw或正式機器goroutine關閉。
五秒是同一監督預算，實際排程與程序退出量到5.058／5.146秒，不作五秒硬時限保證。
Go不能取消正在阻塞的系統呼叫；晚回失敗要由呼叫程序結束工作程序，不能並行關閉同一Writer。
現有清理是盡力而為，沒有斷電耐久、惡意並行替換、真實阻塞Link或清理雙故障驗收。
故障策略及正常停止入口待使用者定稿，013維持DRAFT，#17保持OPEN。

R867最終獨立審查阻擋0、未解除應改0、建議0；主代理完整讀取兩報告。
原goal的全面發布聲明已限於已量正常案例。Writer先Link發布，再清理片段，
協調器仍可能逾時；撤回僅盡力，不保證所有逾時均無成功名稱，完整B4未完成。
審查逐項重算134項未變來源、31產物及6份私人Go的SHA／size，兩份批准追加的快照前綴相符。
JSON SHA-256為`b61e7953b56893c989603853a79f8f3f198890aa9b0eca93a1025aa84297c9d7`，Markdown為
`4a3a7344668ed1a627bd48231b718fb4136f994aff065d056bd1dbe19579a625`。
根收尾入口`workplace/r866/final-check.json`；不重跑產品或測試、不覆寫凍結索引，013仍DRAFT。

<a id="r868-recording-stop-prototype"></a>

## R868：正常停止錄製的操作比較

2026-10-03；[目標](../goal/R868-recording-stop-prototype.md)、[自製原型](../../tools/recording-stop-prototype.html)、
[瀏覽器驗證程式](../../tools/recording-stop-prototype-test.cjs)。只比較停止入口，沒有遊戲資料、音訊或錄放Writer。
候選A預先指定示意停止步數；候選B以按鈕要求下一安全邊界停止。兩者都示意收尾後結束工作階段，
續玩未製作、未排除，仍是待定分支。示意步數200000及300／200毫秒不能當原版或實機證據。

- confirmed僅限自製UI：初始狀態、自動終點與拒收後續輸入、非法終點拒絕、手動邊界停止與拒收輸入、重設使舊計時器失效、390px寬無水平溢出，六案例通過。
- Node v22.12.0、Chrome 131.0.6778.204；既有映像`ghcr.io/puppeteer/puppeteer:23.11.1`，ID `6a74444753c00dd2420b241d13d653fff9762b31033d6f7f3eb94fc059fb55ca`。
- 原型SHA-256 `0eac79888454230d8aee3bd629a469ef9958503e9d202489c7081adff3d9e04a`；驗證程式`b01865b3ea787425e9a03dc175ab36108a33cad1c0170af0abe606760dd704bd`。
- 收據`workplace/r868/report.json`的SHA-256為`bcdceb365d343fbf6a9530a030da23a372cc75494fc7d3fc94d9f66cd9aee765`；頁面錯誤及外部請求均0。
- 比較圖`workplace/r868/comparison.png`：`933aedbea240cf8a2a1287547564ec8648f1e8d3d16ff0999582d7c5452f7664`。
- 等待邊界圖`workplace/r868/manual-waiting.png`：`23b6691748d0b2795f1009a085977701d9dae0f9da37d28a6936307dce7cc41a`。

主代理目視核對兩張圖，中文可讀、按鈕沒有裁切。第一次啟動因HOME換成/tmp後找不到映像預裝Chrome而失敗；
核對既有快取位置後明示`PUPPETEER_CACHE_DIR=/home/pptruser/.cache/puppeteer`重跑通過，屬環境設定。
第二次成功測試只調整等待圖左側回到預設狀態，六案例仍通過；不是原版回歸基準。

重跑在上述映像內使用`node /src/recording-stop-prototype-test.cjs`；tools目錄掛/src唯讀、
明確輸出目錄掛/out可寫，UID/GID 1000:1000、network none、HOME=/tmp、
NODE_PATH=/home/pptruser/node_modules及上述快取位置。容器需init、60秒外層期限、
1GiB記憶體、2 CPUs、256 PIDs、cap-drop ALL、no-new-privileges及日誌上限。
每次比較另指定新的輸出目錄，不覆寫已採認收據。Chrome的no-sandbox只用於這個無原版／無網路的可丟棄原型容器。

使用者尚未選定停止操作；既有故障保存策略問題保持待答。本輪不新問相依問題，
不把原型當定版樣式或正式停止行為。013仍DRAFT，#17保持OPEN。

<a id="r869-recording-config-wire"></a>

## R869：完整設定的單檔解析

2026-10-03；[目標](../goal/R869-recording-config-wire.md)、[013候選契約](../spec/013-input-recording.md#r869-recording-config-wire)。
confirmed僅限自製格式候選與下面驗證，不是正式格式、擁有者或可玩性證據。

- 固定來源R862 `candidate/initial.go`：14774 bytes，SHA-256 `b698a2bbffbbc6240c67889ebc91b6e3562c92b05272f6737c37a59b581d587d`。
- 固定來源R866 `recording/record.go`：19221 bytes，SHA-256 `65d805c303ded144e10bc40c3063da9b607340117eec9a6674c7afcdb9c6595b`。
- 私人生成入口`workplace/r869/prepare.py`、競爭驗證`run.py 3`、獨立驗證`independent.py 3`；只有自製設定與逐列檔案，不讀原版資料。
- Go1.24.13，既有`wizardry7-go-ebiten:latest`映像ID `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。

候選檔頭將完整Config嵌入，config_sha256沿R862固定struct的Go預設JSON編碼；外層schema改candidate/2，
不宣稱舊prototype/1相容。header／配置依精確JSON標記遞迴驗欄，兩層檔案、來源與資源均涵蓋。
欄序與空白可以不同，配置SHA按型別重建，內容SHA按原始列；snapshot的frontend_sha256綁定完整配置SHA。
外部expected SHA由呼叫端提供，本輪沒有實際啟動擁有者，不能宣稱已核原版Root／Scratch或目前binary。

競爭收據`workplace/r869/report-3.json`：23個限定leaf案例通過，Go的24個pass包含一個群組標頭，不能多計。
其中必要欄位群組有118個路徑，各刪除及置null，共236次負向檢查；所有變異檔重算正確內容SHA。
另有配置SHA不符、重算SHA仍不合法的速度／路徑、完整包裝恰1MiB及超限前不建partial、Writer拒失真Unicode。
本輪不重跑舊命令大表；正常往返抽樣run、press_key、type及snapshot四命令。

獨立Python收據`workplace/r869/independent-3/independent-report.json`含八例：兩個合法檔及六個拒載，
拒載CLI退出2且不回傳文件摘要。設定由Python逐欄手寫，含中文、補充平面字元、HTML字元及U+2028。
獨立配置SHA `b7895f1b49df5cdc84c622b366d22f5acc9f64164938e182b13f7d810adbeb7d`；
競爭版讀檔binary SHA `0b90b8a1d9bedd8eebf0986ef5f21f7a7fc4774af1a0043aa65a1e1c8560b0a4`。
這個binary只解析自製檔案，沒有LoadWith、遊戲、視窗、音訊裝置或投遞API。

第一次競爭跑在測試樣本的曖昧font_full編輯停下，屬測試定位錯誤；保存來源後縮小到DOS字型值。
第二次22個限定案例通過；另開單一探針重現Writer把非法UTF-8轉成替代字元，屬候選編碼缺口。
保存修前來源及收據後只加JSON往返一致性守門，再跑23例及相同八個獨立例，均通過。
初次紀錄`report.json`、`tests.jsonl`，第二次`report-2.json`、`tests-2.jsonl`，失真探針`unicode-probe.jsonl`均保留。

重跑將專案掛/src唯讀，本輪輸出掛/out可寫，使用UID/GID1000:1000、network none、
2GiB／2 CPUs／256 PIDs、init及日誌上限；競爭驗證外層220秒、Go子程序180秒，獨立驗證外層120秒。
Go環境沿`GOPROXY=off GOSUMDB=off GOFLAGS=-mod=readonly GOMAXPROCS=2`，快取只寫/out/cache，模組快取唯讀。
每次run.py用新編號，independent.py另建新編號子目錄，不能覆寫已採認收據。

正式Go、SDK、譯文及字型保持。故障保存、正常停止、正式擁有者／Draw及發布後期限缺口仍未解除；
013維持DRAFT，#17保持OPEN。封存來源、產物與獨立審查另列。

本輪封存68來源及53產物：65來源保持、三份規格／目標／證據只允批准追加。
`compiler-inputs.json`為讀檔binary的55項選定來源，另將本候選兩份TestGoFiles納入來源索引；
不把SDK自己的測試來源或stdlib當本輪全部已測試，工具鏈依上述映像及Go版本綁定。
來源索引`workplace/r869/source-index.json` SHA-256 `ec15d6317734108c2a4dff88a8d3609e4fed1bc769e7a67a7e984c56e22921ca`；
產物索引`frozen-output-index.json`為`12011731017129b2cf033baba9273835621e1da15f70324a33dd277d37b8c1a1`。
最後競爭報告SHA為`b006a0b979aa2289bacb5a10e449e0fe9f157df3baad76be8aa89d447069966c`；
獨立報告為`e6532d1ec7292fe890cd2d4bcf73862cea98d13a408a7eaeb158ca3ab06cadec`。
根核對入口`workplace/r869/final-check.py`及`final-check.json`，驗65來源、三快照前綴及53產物相符，
157個文件連結存在、全工作根衛生0，正式前端保持，沒有重跑產品。

另核同旗標的`go list -race -deps -json ./cmd/checkwire`，55項選定來源與封存清單逐項相同。
追加收據`workplace/r869/race-source-check.json`及`go-list-race.json`，後者SHA-256為
`b4843d4403ac8d703710cffb849c415c161ac9e14f1c737306b14b80d48429ac`；這是metadata核對，未再執行候選測試。

R870獨立審查阻擋0、應改0、建議0，主代理完整讀取兩報告。
審查核對65未變來源、三份快照前綴、53產物、55選定來源、23限定案例、236欄位拒絕及八獨立例相符。
報告`workplace/r869/review/review-r870.json` SHA-256為
`a02b52e195e11b5da8e5df6c89f7830d1b2243ef4ec534a0cb6f4f699b70a19b`；Markdown為
`5d3b5768ede024b925acdfb4415fca7f8f1e348004de99104d891ce43e73c6f9`。
首次審查因「尚未追加」整檔相等假設停下，按明示批准的追加政策核完整前綴及新增段後通過，屬核對時點問題。
審查未重跑產品或測試，自有容器0；只採認DRAFT回填，不解除正式錄放或交付門檻。

<a id="r871-recording-startup-wire"></a>

## R871：單檔與實際啟動核對

2026-10-03；[目標](../goal/R871-recording-startup-wire.md)、[013候選契約](../spec/013-input-recording.md#r871-recording-startup-wire)。
confirmed僅限以下私人接線及冷執行，不是正式格式、GUI或完整可玩性證據。

- 沿R864啟動擁有者與R869嚴格整檔解析；後者原始SHA為`f4a9d825048b43158c386d47e13c267affc353d51da396ef3a5cf3a7f6ad2d32`。
- 原版輸入`workplace/original/wizar7`唯讀，248檔逐檔SHA與R859基線相同；前後完整清單保存本機，不解碼手冊內容或送入答案。
- Go1.24.13／Python3.11.2，既有映像ID `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
- 候選入口`workplace/r871/entry.go`、生成`prepare.py`、建置`build.py`、驗證`run.py`，正式Go與SDK保持。

外部檔案先完整Load，通過才核Root／Scratch及同次讀取的資源結果，並建構實際擁有者。
候選不從錄放檔設定自行改用另一個EXE、資源或選項；預期值與實際值逐項比對。
不使用R864的可信設定json.Unmarshal入口；該入口已從私人副本移除。
內部Config的json.Unmarshal只接本擁有者產生的bytes，外部仍走嚴格wire解析。
本輪沿既有輸入獨占、原版唯讀及啟動單一goroutine假設，不增加並行換檔保證。

`attempt-2/report.json`記錄27案例：七正向、二十負向。
正向為獨立LoadWith控制、記錄一次、重播三次、等價檔頭欄序及使用內嵌來源清單。
三筆Run各200k，全部0／200k／400k／600k的RAM、暫存器、indexed、色盤、Step及Ticks摘要相同。
不是完整CPU內部狀態、中文合成像素或固定亂數判定對拍；沒有玩家輸入及手冊答案。
資源沿實際載入器讀取，但沒有真正GUI Draw或聲音裝置，不能稱已完成正常玩家錄放。

十二種設定或來源差異及七種wire／候選子集差異於LoadWith前拒絕；
一個合法但錯誤的cpu_hz自述於實際載入後、交出前拒絕。
二十個負向owner_take_calls／run_calls均0，十九個load_with_calls=0，初態差異一例為1。
變異檔的配置及內容SHA由獨立Python按既有Go欄序與HTML跳脫重算；
巢狀重複、null、缺省及未知欄位案例有正確內容SHA，並非只測到校驗錯誤。
獨立控制組繞過設定擷取及錄放，直接同Options載入、SetScratch及Run，沒有拿候選輸出當控制期望。

首次負向測試把program.sha256誤寫成binary_sha256，實際回wire_rejected而非預期的靜態不符；
屬驗證樣本錯誤。原腳本`run-before-fixture-fix.py`、cases及fixtures保留，候選程式未改。
修正欄名後在相同映像與參數重跑，另存attempt-2及run-2.stdout／stderr，27例符合預期。
不是因產品故障改規格或加入例外。

編譯閉集合75項：72個本專案／私人SDK選定Go來源及三份模組鎖定檔；
同旗標go list -race清單逐項相同，來源清單與SHA嵌在binary，無需設定附檔。
內嵌來源模式仍掛載研究工作區，本輪不稱無Go來源獨立包；R864的獨立包證據另留原範圍。
stdlib及第三方依賴沿映像、Go建置設定與模組鎖定，不宣稱本輪逐檔驗證。

| 收據 | SHA-256 |
|---|---|
| bind | `9742fe7d1460cd9a50d30d5f03e4db1c8770e939842651c74ecf7743290d8171` |
| compiled-sources.json | `9649be9bbff9e4568ff9c4a5bce032eb76ba885825b4c2fdce02d05eda23a7a6` |
| attempt-2/report.json | `10f79327b7860a8155c556810ce2774711832706b6d76ba6a5d9c072c3665dde` |
| attempt-2/limited-run.jsonl | `cfffd79c5ae5287119be6b12cb0669d953b51c774c97c35baf9658f9ee7d1ec3` |
| race-source-check.json | `8257b4b7d63a98b9150bd6bed291810bb455bfd79e01739deacfa880d13c2559` |

配置SHA為`83096d6ff32f77a0cd0e98e461de640d3fb54e58f6d98de4d5837af00e2045e6`。
重建入口先在上述映像跑prepare.py再build.py；驗證run.py需另選新attempt名稱，不能覆寫採認輸出。
專案掛/src唯讀、研究輸出掛/out可寫、原版掛/original唯讀，UID/GID1000:1000、network none、
2GiB／2 CPUs／128 PIDs、init及日誌上限、外層240秒；建置用256 PIDs及260秒。
Xvfb由有界Python父程序持有，finally逐一終止子程序；模組快取唯讀、Go快取只寫/out/cache。

候選只驗三筆成功Run。測試自行FinishAt，沒有正式停止操作、最終Draw、音訊關閉或故障保存。
B3／B4與發布後撤回限制保持；013仍DRAFT、#17仍OPEN，十一組、TSV及正式封包未完成。
封存及獨立審查另列，不能由這批短段測試解除正式試玩閘門。

封存165來源及115產物，含159個保持來源及六份只允批准追加的完整文件前綴。
來源索引`workplace/r871/source-index.json` SHA-256為`c79b68b07d684e654be2efd7a50e032fee53409e4813c6fd141a01b885058cc0`；
產物索引`frozen-output-index.json`為`b013c494f72e67fbce20add43b0c7c086138344a8ec9d32172e60f7920e0ea88`。
來源包括實際編譯副本與其輸入，不把165當獨立編譯檔數；binary閉集合仍為75。

R872獨立審查阻擋0、應改0、建議0，主代理完整讀取兩報告，只採認上述限定DRAFT回填。
逐項核165來源、115產物、75編譯項、27原始stdout與報告；13個完整變異檔的配置及內容SHA正確，
另兩檔刻意附尾資料或截斷。初次測試欄名修正不改候選，範圍與拒絕順序符合文件。
`workplace/r871/review/review-r872.json` SHA-256為`86521333512b19560f4806248fe6703b5c54723b49c5e26d30ba5d503aa76c4b`；
Markdown為`0599ff7221628c9acfbcc629720c160f4f0fd0ccb07ec4a917a9c3f091c2a220`。
審查首次誤把索引包裝當陣列，核對前退出，按sources／outputs修正後完成，屬審查工具假設錯誤。
審查未重跑產品或讀原版正文，報告owner1000:1000、自有容器0；正式錄放與交付門檻保持。

<a id="r875-recording-live-shutdown"></a>

## R875：真實GUI錄放與收尾

2026-10-03；[目標](../goal/R875-recording-live-shutdown.md)、[013候選契約](../spec/013-input-recording.md#r875-recording-live-shutdown)。
confirmed只限下面私人副本與實測路線，正式Go、SDK、譯文與字型保持；013仍DRAFT。

- Go1.24.13、Python3.11.2、既有wizardry7-go-ebiten映像，ID沿R871；原版248檔與三份既有烘製資源唯讀。
- 本輪未重烘全catalog，既有資源不當最新TSV、控制碼或剩餘十一組的完整驗收。
- 76項編譯集合包含73個本專案／私人SDK Go來源及三個模組鎖定檔。copied inputs包含未選定的中間來源，不能當編譯檔數；stdlib／第三方沿映像、鎖版依賴與實際binary，不逐檔宣稱。
- 私人schema為wiz7-input-live-candidate/3；起始設定、資源同次讀取、實際LoadWith與單次交出沿R871。panel只保存Draw實際使用的說明／地圖開關、說明畫面及頁碼。
- 禁用autoplay、hd-at等額外有狀態合成，不用作弊、狀態注入或防拷答案；輸入從XTest經正常前端佇列交付。完整本機錄放檔不進Git或發行包。

### 已驗路線

| 核對 | 結果與界線 |
|---|---|
| 正常版冷啟動與主選單 | 300M步，2,952事件，包含run、snapshot、press_key、type、move_mouse、click、press_button及release_button。F1在8.2M收到真正暫停確認並同Step恢復，發布41保存說明面板；主代理目視中文可讀且無裁切。 |
| 正常左鍵 | 270.4M於原座標160,61點擊，實際after=276.4M、原API結果none，新增角色面板及游標可見；不由API預設常數推導結果。 |
| 三次重播 | 1,471發布的序號、Step、配置SHA、像素SHA及PNG SHA逐項相同，軟體stream統計也相同。每個ID只採首次完成Draw，不是全部GUI畫格。 |
| 原版控制 | 不經覆繪／compose／音訊，直接相同Options的LoadWith、SetScratch及命令。全部發布Step／Ticks與38抽樣點的RAM、暫存器摘要、indexed、色盤、A000相同；不是全CPU內部狀態或全部發布完整RAM比對。 |
| 自然角色資料寫入 | PCFILE.DBS 9,368 bytes，錄製、三重播及控制SHA相同；248原版檔不變。只核檔名、大小與SHA，不讀正文。不是保存遊戲進度、關窗保存或冷LOAD還原。 |
| 輸入差異 | 只把左鍵y改為43，合法文件重算內容SHA。119發布像素及Ticks不同，1,471發布Step仍相同；完整雜湊差異四抽樣點。 |
| 獨立圖片解碼 | 7,355份PNG全部重算檔案SHA；每路15個發布由ImageMagick解碼為RGBA8，75份完整像素SHA與Draw相同，含F1及最終發布。 |
| 真Draw負向 | 錯發布回draw_mismatch，缺Draw回draw_timeout，均退出2且無成功階段。兩例皆重播、沒有Writer，不採Writer負向發布證據。 |
| 短段競爭檢查 | 3M錄製、一次重播及控制通過，15發布，含F1、鍵盤、移動與右鍵；沒有左鍵。無DATA RACE；左鍵長段為正常版。 |
| 面板wire | 17案例，兩接受／十五拒絕；缺省、null、未知欄、錯型別、非snapshot附panel、非法文字拒載，等價欄序接受。未LoadWith。 |
| 音訊軟體 | stream Received／Applied皆7,835，CanonicalFrames／Frames皆4,404,019，SourceStep／EndStep皆300M，Pending／Rejected皆0。Reader靜音，DeviceFrames／DeliveredBytes皆0、ACK=false；不當裝置或人耳驗收。 |

最終發布先取得同ID／Step／配置SHA的真Draw，再放行owner、關閉正式音訊，最後完成Writer；
使用R866同一Stop.Deadline，不另起五秒。私人明示Step只為蒐證，正式停止操作／故障保存仍待定。
發布後盡力撤回及阻塞系統呼叫的限制沿R866，沒有由本輪解除完整B3／B4。

### 收據與重建

成功錄製在attempt-5/record，第一次重播在attempt-7/replay-1，第二、第三次、控制及變異在attempt-8。
三次成功重播分別核相同起始配置，跨批次只沿用已完成收據，不把逾時的部分執行當成功。
完整彙整workplace/r875/verified-report.json，原版檔案與錄放原文只留本機。

| 檔案或內容 | SHA-256 |
|---|---|
| 正常版bind-normal | `36c7832676f8f3f0a0554e5f32285187eb872a032465c585e26f36bc93b40d4b` |
| 競爭版bind | `5c73cde4701d51760e74bc85de77d6c136d7cb497e03767153a63b0c593fd86d` |
| compiled-sources.json | `f5b51f96e05b5bc0c8de92ec5078e47f750e2faef62c4cce4bdb5b73d820f06d` |
| 300M起始配置 | `b44635d15a49c0b93f22f64161b2ee60d2ec78578cddcc3808aa26fa897b6acd` |
| attempt-5/record.jsonl | `f2bfcc9fed077d3733b474bfdca1657cb9d2997e1894ab20da7c354d70cb79f8` |
| verified-report.json | `e72a51887c5795e13a0f4cfdcb1cb6d6efec6c88d41f806abea55ee76343afc8` |
| attempt-8/report.json | `068ca0f305a1af7586d2b3c92dc431091cbad73d174c6818aacb3223bc8a83cc` |
| 自然PCFILE.DBS | `24d94441e5b32e7abbcbffa9dad972b41544da6a7bb06096b967281f639d9c1b` |
| sealed-sources.json，186項 | `e463b5b871689bc9fd6f927050e11338fa221843a0bce939f083a60cb04247ae` |
| sealed-outputs.json，16,581項 | `f6402fdc93f22db4eb07a98b254346a1bb1f17b2be02be86519c606bd2ca5c69` |

來源索引包含正式輸入、私人編譯副本、工具及六份文件快照，不把186當獨立編譯檔數。
工作樹只允CONTEXT現況修正及規格／證據／goal／WORKLOG批准追加；快照與凍結索引不重寫。
建置依序prepare.py、R864 build.py及build_normal.py，rebuild.py在乾淨輸出驗證兩binary及76來源逐位元相同。
須把已存在的workplace/gomodcache唯讀掛到/mods並設GOMODCACHE=/mods，Go快取只放明確可寫目錄；
其他Go環境為GOPROXY=off、GOSUMDB=off、GOFLAGS=-mod=readonly、GOMAXPROCS=1或2。
專案掛/src唯讀、研究目錄掛/out可寫，原版及烘製資源分別掛/original、/resources唯讀；
UID/GID1000:1000、network none、init、記憶體／CPU／PID及日誌上限。
完整執行用run_left_click.py產生一次冷錄製，再依run_left_complete.py／run_left_remaining.py核重播；
重跑需另選輸出編號，不覆寫採認收據。GUI批次上限2000秒、個案擷取600秒，遊戲Stop總期限仍五秒。
Xvfb及遊戲由Python父程序持有，finally逐一終止子程序；原版前後清單與R859基線一致。

失敗歷程只在WORKLOG；首封存工具誤把/out當/src子路徑，修正索引映射後通過，沒有重跑產品。
來源／產物owner1000:1000，全工作根root-owned及誤建.md目錄0，自有容器0；獨立審查另列。
正式CLI、故障策略、正常停止、冷遊戲進度還原、音訊人耳、十一組、TSV與封包保持未完成。

R876封存後唯讀審查逐項重算186來源與16,581產物，核實際76編譯輸入、兩個乾淨重建程式、三次成功重播及限定狀態抽樣。
阻擋0、應改0、建議0；主代理完整核讀JSON與Markdown，只採認限定DRAFT證據，不重跑產品。
報告留在`workplace/r875/review/review-r876.json`與`review-r876.md`，SHA-256依序為
`c20bd9af79e5c505d97826a40f9033c1de17c72e81e1acc3d152676d01bb87fa`及
`075f54ca4178a3cc4b45204d716883f4077c37a2a96d445e4ce5e72f7243755e`。
報告owner1000:1000，審查自有容器0；封存索引保持，後續審查與文件追加另存，不回寫凍結收據。

<a id="r893-recording-stop-continue"></a>

## R893：停止後退出或續玩的操作比較

2026-10-03；[目標](../goal/R893-recording-stop-continue.md)、[原型](../../tools/recording-stop-continue-prototype.html)、
[瀏覽器驗證](../../tools/recording-stop-continue-prototype-test.cjs)。補齊R868未製作的續玩比較，舊原型與測試保持。
兩候選均由玩家停止，等待示意安全邊界及收尾；一個結束工作階段，另一個回到可操作狀態。
沒有原版、錄放Writer、Draw、音訊或再次錄製，不採為正式行為或停止期限。

confirmed僅限自製UI：七個真瀏覽器案例通過，包含退出後拒收操作、等待與收尾拒收操作、
續玩只增加遊戲操作且錄製截止保持、重設取消等待／收尾舊計時器，以及390px寬無水平溢出。
續玩示意的遊戲操作由2到5、步數400000到1000000，已錄製操作仍2、截止仍400000。
頁面錯誤與外部請求均0；主代理目視核三張圖，文字及按鈕無裁切。
Node v22.12.0、Chrome 131.0.6778.204及Puppeteer映像沿R868，映像ID本輪重新核對相同。

| 證據 | SHA-256 |
|---|---|
| 新原型 | `4624e21403c793511b1f8d3e0e1a1cfb5f85eae2aa042a2b8a85cd9f300cf971` |
| 新瀏覽器驗證 | `73e109656b19d5c3573ef6cb4c2d4c3eac338dcfd23ff4037f8e0e70ced246c4` |
| workplace/r893/report.json | `e12eaa3ae5f0715fb03d931d25a0726380860f4e1f2d1e6f69a12c94103558e2` |
| comparison.png | `e17ac07b99ffa0e270cf759302502662cdc6b3e1da680160f34247ed4fcfc88a` |
| continue-waiting.png | `f3492d6bce64ee7bf6d9a4203e2bb22bc492706ce270aa85b85d173d35b29e18` |
| continue-playing.png | `bf1f5b6bb4393c8ffce50846b01edee172fa602537c33554a85bcb70f4833bf0` |

重跑沿R868隔離瀏覽器契約，改用上述新驗證程式並掛獨立新輸出目錄；不覆寫採認圖與收據。
只掛tools唯讀及輸出可寫，不掛原版、SDK或錄放檔。示意步數與300／200毫秒均非實機證據。
1,159個正式前端、SDK、013與舊原型檔案前後SHA相同，產物owner1000:1000。
`workplace/r893/sealed.json`封存兩來源、七產物與未改輸入，SHA-256為
`f131955eff53028ff11036982fdfc8c069b3dfd5855e2e0a0a8cde29b5ad7833`；
`verified.json`為`b4fff8956f5e079ea09042267a7fb54e6bce2995f4d9cf412a7d19942384721c`。
停止操作與故障保存仍待使用者定案，不重問既有待答問題或先寫正式分支。013保持DRAFT，#17保持OPEN。

<a id="r945-recording-stop-owner"></a>

## R945：錄製停止與續玩擁有者

2026-10-04，沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[目標](../goal/R945-recording-stop-owner.md)及[013候選契約](../spec/013-input-recording.md#r945-recording-stop-owner)。使用者已定案玩家按鈕停止、當前操作完成後收尾、停止後續玩、故障保留無成功末註的片段。R893上節的待答敘述保留為當時歷史，現在依此定案推進。

R875的最後Draw→owner完成→音訊Close→Writer.Finish屬整個工作階段退出，不能直接當正常停止錄製。本輪私人`workplace/r945/stop.go`將錄製截止、Writer單一擁有者及遊戲退出分開：自製API先完成並記錄結果，截止後不再收事件；最後發布以ID／Step／配置SHA確認，才追加snapshot及收尾。沒有Oracle、原版Stop或音訊Close入口，原API結果不被錄製錯誤取代。所有自製資料只表明候選邏輯，不是原版語意證據。

### 已驗範圍

Go1.24.13、Python3.11.2及既有`wizardry7-go-ebiten:latest`，映像ID為`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`；來源只讀，UID/GID1000:1000。初輪未執行測試，因正式SDK沒有DOSSettings getter而建置失敗；核R819既有私人getter後，改唯讀沿R875私人SDK，不改正式SDK或產品。attempt-2通過後發現接收worker成功結果時缺期限再核，保存修前來源，補接收時期限守門；同映像、同命令於attempt-3乾淨重跑，15個leaf案例通過、失敗0、沒有DATA RACE。Go的群組pass不另算案例，晚到接收分支未由決定性交錯觸發，不能將此來源修正當新增動態覆蓋。

| 驗證 | 結果及限制 |
|---|---|
| 正常收尾與續玩 | 真實候選Writer完成截止200000的run＋snapshot；再呼叫三次自製操作，步數到200300，完成檔bytes保持，未多記事件。 |
| 操作中要求停止 | 自製run中提出要求，操作仍完整增加17步，文件截止200017；沒有強制中斷API。 |
| 自製API錯誤 | run與click各保留原error物件及37步部分推進，留下無末註片段；不改既有後續回傳路徑。沒有原版錯誤注入或Oracle.Run。 |
| 最後畫面拒絕 | 六例：錯ID、錯Step、錯配置、故障狀態、晚到完成與缺確認，均未完成文件。全部通知自製，沒有真正GUI Draw。 |
| Writer故障 | 最後snapshot Append與Finish兩例自製錯誤回呼，留下片段，成功名稱不存在；未測真磁碟故障或所有清理雙故障。 |
| 阻塞Finish | 自製回呼阻塞超過40ms測試預算後，擁有者可執行後續未錄製操作；阻塞返回前Abort次數0，返回後真Writer按同期限拒絕、保留無末註片段。沒有在Finish中並行Abort，非正式五秒實機保證。 |
| 重複要求 | 二十個競爭goroutine只有一個期限，重複要求不重設。 |
| 軟體音訊 | 停止前後由正式playerAudio接兩筆自製CompletedOPL時間批次，Frames2937→5873、EndStep200000→400000，Stop／裝置Pause／Close皆未呼叫。沒有OPL寫入，不當連續有聲音樂、裝置或人耳驗收。 |

獨立Python直接解析保存的十五份自製文件，確認三份成功、十二份片段：配置JSON SHA、原始列加LF的內容SHA、seq／before／after及截止Step符合；片段無end且無成功名稱。沒有重用Go Loader當獨立期望，不重跑舊wire攻擊表或宣稱完整格式覆蓋。`go list -race -deps -json`核93項選定編譯來源，含cmd所有編譯測試檔，但實際只執行TestR945；stdlib與第三方沿映像及模組鎖定。

| 本機來源／收據 | SHA-256 |
|---|---|
| `stop.go` | `b392dd38fb9ae88c625e55c0a172652baefb6c29d5ac39a94cc314f384bc2db6` |
| `stop_test.go` | `22d8c7d0d4201f4f64afddacc8ee271b38b8d2e7129fb441e6271a6f59349e55` |
| `attempt-3/report.json` | `e5292563b077f8a2e461752bbaecc3cf5eb2e3e7c0e57d45296f2d1f000ee923` |
| `independent-3.json` | `efa7971d2e2cac47e2abe0d96f65197fd2f23e934cf598b71712a98d9411be8a` |
| `compiled-sources-3.json` | `d5d03291b75ca34b736954848a42d1d1f57065f53f7441a54c07eea9a2054e1c` |

重跑先由本機`prepare.py`在新的已存在輸出目錄複製當前cmd、R875格式及初始化候選，沿R875唯讀私人SDK；再跑`run.py`及`verify.py`，兩者須設同一`R945_ATTEMPT`，預設attempt-1。來源repo掛`/src`唯讀、私人輸出掛`/out`可寫、既有gomodcache掛`/mods`唯讀；沿離線Docker、2GiB／2 CPUs／128 PIDs、init與日誌上限，外層220秒。Go設定GOPROXY／GOSUMDB為off、GOFLAGS為-mod=readonly、GOCACHE只寫`/out/cache`。採認收據為attempt-3，獨立報告以`R945_VERIFY_SUFFIX=-3`另存；不覆寫舊編號。Xvfb由Python有界父程序持有，finally停止並回收，沒有持續背景容器。

此候選仍有正式缺口：未移植getter、GUI按鈕與每個原API安全邊界；未接實際音訊／observer故障資格及關窗競態；沒有連續有聲音樂或三次重播。R862完整設定只列xlate／flow／help，現在正式程式另載入roster.json、profile-sex.json及intro-picture.json；錄放啟動必須納入同次載入結果及SHA，不能因三類舊資源核對通過就稱最新前端配置完整。阻塞系統呼叫不可由Go取消，晚到成功名稱或清理失敗的限制仍沿R867；不由這15案例解除完整B3／B4或將013升READY。

44份複製來源及93項選定編譯來源保持；私人SDK與正式SDK的選定編譯差異只有`oracle/initial_scalars_candidate.go`。原版248、正式SDK1132、R943封存964與R944封存199全部SHA保持，production未改，工作根root-owned／誤建md目錄及輸出擁有權異常0。最終`integrity-final.json` SHA為`daa7943d57d6f05e1998156fc0830a1dfa87fb5edd8c6a989152d78269458511`。

最終獨立唯讀審查為阻擋0、應改0、建議0。主代理於FINAL後完整核讀兩份報告，只採認私人候選及DRAFT回填，不批准正式錄放或READY。`workplace/r945/review.json`與`review.md`的SHA-256依序為`41121fd73b4585aad299fc59081d279ff65411bd5edf2c2e0d0fb63d26fd1475`及`34ff86dc5c0d8930205a86159380bf64ab536c975c3c123480d7cf1b9b640389`；報告UID/GID為1000:1000，自有容器0。

私人111份技術檔已一次封存，`workplace/r945/closed-inputs.json` SHA-256為`b3b585fe3b7f3f83a3edf1d15f656956c2aef9c7553642045f9b6a4c570b3070`。排除可重建快取與後續遠端同步回讀；來源、修前歷史、三輪收據及審查報告保留，不覆寫。擁有權異常0，私人候選與自製錄放檔不進Git。

<a id="r946-six-resources"></a>

## R946：六類資源的同次讀取候選

2026-10-04；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[目標](../goal/R946-recording-resource-binding.md)、[013候選契約](../spec/013-input-recording.md#r946-six-resources)。沿既有完整設定決定補技術缺口，未改正式程式、SDK、catalog或字型，未執行原版。正式載入位置是main.go的xlate／flow／help、roster.go:112、profile_sex.go:94及intro_picture.go:90；正式來源SHA保存在私人source-inputs.json。

私人prepare.py從目前正式來源逐字取出五個完整loader，只替換函式名稱及讀檔入口，仍用正式decodeRoster／decodeProfileSex／decodeIntroPicture。一般譯文解析另沿main.go的map／digits／ordinals規則，解析錯誤改有限代碼，沒有呼叫正式main或log.Fatal。側車一次Open讀到記憶體，保留同一handle的Stat，解析與SHA共用相同bytes；超過既有1MiB上限拒絕錄放資格。xlate／flow／help沿既有ReadFile方式，本輪不宣稱所有資源讀取都有新容量限制或抵抗同一檔案的並行寫入。

Go1.24.13、Python3.11.2；既有映像`wizardry7-go-ebiten:latest` ID為`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。attempt-1的22個leaf案例通過，失敗0，無DATA RACE；只執行TestR946，其他編譯測試不算執行。confirmed限以下自製候選：

| 案例 | 已驗範圍 |
|---|---|
| 六資源成功2例 | 15／16字高自製字模的全部六份資源，各讀一次；五個loader與正式結果DeepEqual，一般譯文map直接沿候選，沒有獨立一般譯文解析oracle。請求、大小與SHA另核符合，非GNU／倚天實際字型驗收。 |
| 側車缺檔／無效6例 | 三類側車各測缺檔與無效JSON；正式與候選均回退。缺檔空SHA，讀到無效資料仍保留其SHA，不誤報Loaded。 |
| 停用／空catalog2例 | 三份側車未要求、未讀取；合法空catalog仍標示解析成功。flow／help依其自身路徑載入。 |
| 換空白3例 | 三類側車各只加JSON空白，解析仍成功，但資源比較拒絕自製後續載入hook，計數0。不是Oracle.LoadWith計數。 |
| 讀後換檔3例 | 單次讀完即以rename換成無效文件；三類解析仍使用已讀有效bytes，大小與SHA指向該份內容，每份候選讀取一次。沒有同檔並行寫入防護證據。 |
| 超限3例 | 三份側車各超過1MiB，錄放配置拒絕；不把未完整擷取當成功hash。 |
| 配置集合3例 | 舊三角色、重複角色及錯序均拒絕。私人候選Build凍結後的資源slice另驗不受呼叫端修改。正式相容策略未定案。 |

獨立Python核19份資源紀錄、兩份完整自製配置：六角色順序、實際或保存的已讀bytes大小／SHA、配置內相同資源紀錄均符合。五個Loaded語意由上述Go正式loader對照證明；一般譯文解析僅沿正式來源規則，未獨立對照，Python沒有重作解析器當新oracle。44份複製來源及94項選定編譯來源保持，私人SDK的選定差異仍只有已證getter`oracle/initial_scalars_candidate.go`；R945封存111檔全部保持。輸出UID/GID1000:1000，工作根root-owned／誤建md目錄與擁有權異常0。

| 私人來源／收據 | SHA-256 |
|---|---|
| resources.go | `bd7f291c45400173c3fa37e87eb9fbce7122a88465db1a78b0da6efc1f64689d` |
| resources_test.go | `10e6eaa72fac9bdcf7fc0f07bdc9148d9bb3f3e8630a6f940e288e14a524943c` |
| source-inputs.json | `5fd090d84e2aef3b1f97aa5fa8dd7c97908afe096ccaee12618ae8366e78457d` |
| attempt-1/report.json | `8961de8efbde4a86a6960372247d223111faf47218235ab43a3db69df02f9196` |
| attempt-1/compiled-sources.json | `07bca9d6fea34613febdeeef3c6f02154809e549bf08fc528691d99cc866cf6d` |
| attempt-1/independent.json | `3ba100680fe6d41839c7070680534c9c38922c02f3c53d5471093cbeab3444d0` |

重跑先在新的已存在輸出目錄依序執行本機prepare.py、run.py、verify.py，後兩者設定同一R946_ATTEMPT，預設attempt-1，不覆寫舊收據。來源repo掛/src唯讀，輸出掛/out可寫，既有gomodcache掛/mods唯讀；離線非root Docker、init、2GiB／2 CPUs／128 PIDs、日誌上限及外層220秒。Go設定GOPROXY／GOSUMDB為off、GOFLAGS為-mod=readonly、GOCACHE只寫/out/cache。Xvfb由有界Python父程序在finally終止回收。

這些證據補六資源的私人解析與配置集合，不是正式起始owner、實際錄放／GUI／音訊／三重播。私人schema只識別候選，不對玩家定案新格式；全域啟動紀錄需在正式單一擁有者契約中替換，與R875／R945尚未整合，013保持DRAFT。

最終獨立窄審阻擋0、應改0、建議0；主代理於FINAL後完整核讀JSON／Markdown，只採認限定DRAFT證據。報告`workplace/r946/review.json`與`review.md` SHA-256依序為`d53db246476395b8b2259a31e12f9783939826d791cc43827d8d113e065946b0`及`61628843fc2e5cda61f052e980fb6beb930085c30b170b897080aa4320199104`。核審查輸入保持後，私人196份技術檔一次封存；`closed-inputs.json` SHA為`684540c2cdbbad8b2aa1ae815b8f31f05625a80efa864fa053261f523f0ee849`。排除快取與後續GitHub回讀，輸出及報告UID/GID1000:1000、自有容器0。

<a id="r947-six-resource-owner"></a>

## R947：六資源與單一啟動擁有者

2026-10-04；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[目標](../goal/R947-six-resource-owner.md)、[013候選契約](../spec/013-input-recording.md#r947-six-resource-owner)。原版為玩家本機`workplace/original/wizar7`的248檔，前後盲SHA與R859基線全部相同；沒有MSG正文、state或手冊解碼。資源使用R943的GNU Unifont與本機倚天各六份烘製資料，原檔唯讀，不重新烘製或改文案。

私人prepare.py先核R946來源與94項選定編譯輸入，再複製目前正式非測試cmd、私人SDK及R946資源／配置候選；原main僅改名formalMainUnused，不呼叫。R864啟動綁定改私人R947定位，移除舊三資源讀取器，改由R946同次載入六份；配置使用相同candidate/2，超限capture fault先拒絕。實際Options、儲存清單、binary與編譯來源、LoadWith及單一指標交出規則沿原候選，不新增正式架構決定。

Go1.24.13、Python3.11.2；既有映像`wizardry7-go-ebiten:latest` ID為`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。建置競爭檢查版binary，83份複製輸入；Go來源與go.mod／go.sum／sdk/go.mod共81項清單由go list重生，來源JSON／SHA由ldflags內嵌到同一binary。第三方與stdlib沿鎖定模組／映像與binary中繼資料，不逐檔宣稱閉集合。

attempt-1的17個私人案例通過，六成功、十一拒絕，無DATA RACE。confirmed限下表的load-only路徑：

| 案例 | 結果 |
|---|---|
| 原始控制 | 直接LoadWith／SetScratch，沒有候選純量／配置擷取，建構一次。 |
| 英文、Unifont、倚天 | 英文六角色未要求；兩套目前烘製資源六角色均讀到且解析成功，大小／SHA與實際資源相同。指標均只交出一次，回傳配置bytes不別名。沒有Draw／字形畫面驗收。 |
| 相符expected、內嵌來源清單 | 相符配置核對成功；不核外部來源目錄時仍用binary內嵌清單，配置與Unifont基準逐位元相同。仍掛同一來源repo，不當無來源封包啟動驗收。 |
| 三側車各換空白／缺檔／無效JSON | 九例皆initial_static_mismatch，實際LoadWith呼叫0，不寫配置、不交出指標。已讀但解析無效的SHA與缺檔不同；換空白仍解析成功也拒絕原始bytes差異。 |
| 片頭側車超過1MiB | resource_capture_fault，LoadWith0，不以未完整SHA建立配置。 |
| 初態CPUHz不符 | 可信expected只改CPUHz，LoadWith1；initial_scalars_mismatch在交出前拒絕，未寫配置，指標由失敗defer關閉。 |

五份成功配置與六條成功初態摘要另由Python獨立核對。擷取前／後／交出時的1MiB RAM、indexed、palette、暫存器SHA與Step0／Ticks0均等於原始控制。只保存摘要，不保存或解碼RAM／畫面bytes。Run與輸入沒有裝計數器，未執行聲明來自私人入口及建構器只呼叫LoadWith／SetScratch／唯讀擷取的源碼；Step0本身不是充分證據。

每個程序的stderr只有完全相同、重複兩次的XGB授權檔缺失與無授權連線通知。首次Python核對誤要求空stderr，修正為只接受這個已保存的精確四行通知，其他訊息仍拒絕；產品17案例不重跑。原失敗腳本與分類收據保留，非產品缺陷。獨立核83來源、81選定編譯清單、binary、五配置、原版248、R945封存111及R946封存196全部相符；正式SDK選定差異仍只有私人getter，擁有權及工作根衛生異常0。

| 私人來源／收據 | SHA-256 |
|---|---|
| cmd/wiz7/r947_binding.go | `4212859c22850e3e0be2c289d8fc0766accdfb08de71a2b0de07a0d0d5e03024` |
| source-inputs.json | `a6106c66954f2150cfcb4d47de0aa381b5964268fc29ceca5759183106c07ac7` |
| compiled-sources.json | `e0971ad5e4b3307f8ceef8b13af056c950d1191050525216443814b2ce361b7b` |
| bind | `a0f9f21ea03e6372441ef97f461e3ef4f8a765d53a44670046916449f02862a5` |
| attempt-1/report.json | `e690ecbf07292b91c27379c895856e9fecdb32a1db492001b1f17f706753cc94` |
| attempt-1/independent.json | `dfe6fe90496aab2a374dfbf5b7f81e32a6fe816ce361b856f8e90ab3ca69fde3` |

重跑依序執行本機prepare.py、build.py，再run.py及verify.py；後兩者設定同一R947_ATTEMPT，預設attempt-1。新的已存在輸出目錄掛/out可寫，repo掛/src及原版掛/original唯讀，建置另掛既有gomodcache到/mods唯讀。離線非root Docker、init、2GiB／2 CPUs／128 PIDs與日誌上限；建置外層220秒，17案例外層600秒，個案30秒，核對50秒。不覆寫採認收據；Go快取只寫/out/cache，GOPROXY／GOSUMDB為off，GOFLAGS為-mod=readonly。Xvfb與個案由Python有界父程序於finally逐一終止回收。

本輪只把六資源紀錄接入真正建構／釋出邊界，bundle未交給正式Draw，R945停止未整合。expected僅可信自製JSON，沒有不可信wire載入；正式getter、CLI／schema、機器goroutine、GUI／音訊／故障／退出競態及三次重播保持未完成。原版與資源的並行換檔防護仍未宣稱，013維持DRAFT，不外推為可玩或封包完成。

最終獨立唯讀窄審阻擋0、應改0、建議0；主代理於FINAL後完整核讀JSON與Markdown，採認限定私人DRAFT證據。審查另由ELF64符號及PT_LOAD獨立定位實際binary的兩個Go字串，內嵌來源清單與SHA逐位元符合。`workplace/r947/review.json`與`review.md`的SHA-256依序為`aae3c40c07f38b577b8671e5b74721dae9ec566258d7a0db0644608ebe6e93c4`及`971e5394fc172bfcb66a64037c99732cf60424e2fd2635cd9faaccc6f51fab02`；這是候選ELF位址空間，不是原版位址。兩報告UID/GID1000:1000，自有容器0，013不升READY。

私人219份技術檔一次封存，`workplace/r947/closed-inputs.json` SHA-256為`ae0aa2079238afcf78ad73ee258a2a2079a865200032ff052ff3ee148d2c0d23`。排除Go快取及後續GitHub回讀；二進位、資源與原版摘要均留本機。全數SHA與UID/GID再核相同，工作根衛生異常0；正式程式與SDK保持。

<a id="r948-recording-gui-continue"></a>

## R948：真正介面停止錄製與續玩

2026-10-04；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[目標](../goal/R948-recording-gui-continue.md)、[013候選契約](../spec/013-input-recording.md#r948-recording-gui-continue)。原版仍為248檔唯讀輸入，資源使用R943的GNU Unifont六份烘製資料；只保存原版檔案盲SHA及狀態摘要，不解碼MSG、state或手冊。正式程式、SDK、字型與原版不改，私人getter仍沿R875。

R947同次解析的Trans、數字／序數字模、Flow、Help、Roster、Sex與Picture交給真正GUI，不重新讀取資源。錄製沿R945單一Writer，原API事件與snapshot共用序列。停止按鈕使用Flow字模，置於原遊戲區外40像素；遊戲區仍為960×600，頁尾另留圖，不混入錄放畫面或原版滑鼠座標。這是已測操作原型，正式位置及格式未定案。

最終Draw只走專用stopDraw，沿首次停止要求的絕對期限。run／click先完成一次AfterOperation，再檢查收尾資格；Boundary前另核音訊、來源與串流故障。publish只取一次final，供完整摘要與末點標記共用。全局關窗取消一般Draw等待，不補事件；錄製最終Boundary仍必須確認真正Draw。重播只有消費全部事件且實際Step等於Footer.End才回報完成。來源修正不等於競態分支均有決定性交錯證據。

獨立來源審查在revision-2至7逐步找出與修正共用期限、接收後期限、錯配旗標、AfterOperation順序、末摘要雙讀、正常關窗及提前重播完成資格。舊完整跑因已找到來源缺口主動停止，或因工具觀察期限／XDestroyWindow錯誤失敗，均保留本機，不採認為全測通過。revision-7來源補審0／0／0，主代理在FINAL後完整回讀JSON與Markdown；只有此副本重新進完整路線。完整歷程見WORKLOG的R948段。

Go1.24.13、Python3.11.2；既有映像`wizardry7-go-ebiten:latest` ID為`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。86項選定Go與模組鎖檔清單嵌入同一binary，stdlib與第三方沿映像／模組及二進位身分鎖定，不逐檔宣稱。一般版及短競態版使用同一選定來源。

| 私人來源／收據 | SHA-256 |
|---|---|
| revision-7/compiled-sources.json | `429c9d4cc6d097ed2bac5a6ca567dfc465c61bc5fda3437f94939dbe6d8d2308` |
| revision-7/bind | `cf7c3a39614ff2b2518d6e482283aab8bde18f46dceab4f6645e1767328971a1` |
| revision-7/bind-race | `4c4fc5fe4441b2abd8b52217c4867598c6627116e481a0c4bbc09be860e62d55` |
| review-revision7-source.json | `0c77315e87cdb429c3a83c48f97da588e4f697103df7a25fab4fbf0f0a70922d` |
| review-revision7-source.md | `4feac5687f08c76a78ff72b8e41832c164613c07e43b0eeba3447fd13dc59779` |

confirmed限下列私人實跑及來源核對。停止後不Close原版或音訊，35次後續原操作不改已完成619,795 bytes的錄放檔；最後錄製snapshot為279,000,000步，2,742事件、1,366發布。當時串流Frames4,095,738，OPL實際寫入7,835、丟棄0、無來源故障，Stopped／Incomplete均false。GUI後續仍經AfterOperation推進軟體音訊；靜音無裝置ACK或人耳證據。正常WM_DELETE_WINDOW關窗退出0，F1暫停／返回、鍵盤及滑鼠真正交付，畫面已見中文主選單、確認、取消及錄製狀態；不外推全遊戲文字覆蓋。

兩個實際GUI故障案例通過：錯配確認回draw_mismatch，缺確認本次回stop_timeout；後者與draw_timeout共用同一首次要求期限，不保證競爭中的哪個有限代碼先到。兩者都無成功檔名、無成功末註，片段在其後15次操作保持；串流Frames分別123,313→155,609與126,249→158,545，重播器拒絕片段，正常關窗退出2。缺確認約5.14秒是這次觀察值，不當實機嚴格五秒保證。未注入實際磁碟、來源故障或全部退出交錯。

短競態版真正GUI在8.6M截止，停止後14次操作、完成檔保持、正常退出0、無DATA RACE。只驗此短路線，未重播或外推全路線競態。另以一般版在8M提前正常關窗重播，回incomplete及退出2，沒有誤報到279M錄製終點；這是已量到的提前退出案例，不外推所有關窗時機。

revision-7完整路線五個程序均退出0：一次錄製、三次重播及不含覆繪／音訊的獨立原API控制。三次重播各1,366發布，全部Step／Ticks與合成像素／PNG摘要相同；36個取樣點另核完整1MiB RAM、暫存器、indexed、palette與A000摘要，獨立控制的全部有限摘要也相同。不是每一發布都擷取完整RAM，不宣稱全機器內部狀態逐步相同。錄製完成後的操作不在三次重播中，PCFILE的自然寫入與遊戲進度冷SAVE／LOAD驗收分開。

獨立Python核2,742事件序列、原始列內容SHA、完整配置SHA與實際六資源bytes、86選定來源及二進位；5,496 PNG逐份核SHA、尺寸與區塊CRC，12份另以獨立PNG過濾器解碼RGBA核像素SHA。來源保持由此核對，不採GUI報告中尚未自行實作的source_unchanged欄位。截止前／後軟體串流Frames為4,095,738→4,186,754，成功檔619,795 bytes、SHA為`594c16cf6bc3dbdccacc2e51aba19f69c6c709564f00331a4d6744f80e73e3f2`。故障兩例、短競態與提前重播關窗另核通過，全部原版前後248檔SHA相同；87複製輸入與先前111／196／219封存保持，工作根衛生異常0。

重跑使用新的已存在輸出目錄；生成及建置對到同一/out，避免父目錄舊清單。repo掛/src、原版/original、資源/resources唯讀，輸出/out可寫，建置另掛鎖定gomodcache到/mods唯讀、快取/cache可寫。離線非root Docker、init、2至3GiB、2 CPUs／128 PIDs與日誌上限；建置240秒、完整GUI外層1,500秒、故障200秒、短競態200秒。run.py正常停止由真正XTest頁尾點擊，預設400M僅防失控守門，不是正常停止入口。Xvfb與子程序由有界父程序在finally回收；觀察期限到期沿同一handle核終止，不另開平行重跑。

目前仍為私人DRAFT。正式getter、CLI／schema、按鈕位置、實際音訊裝置、人耳、磁碟故障、發布後撤回與阻塞系統呼叫限制保持。這條冷啟動到人物選單路線不證明完整遊戲、十一組缺譯、存讀檔冷還原或封包完成；013不升READY。

最終限定結果窄審0／0／0，主代理於FINAL後完整回讀兩報告，只採認上述私人實測。`workplace/r948/review-final.json`與`review-final.md`的SHA-256依序為`0c190ebf9ea9d6bdf681aed29dcec2b3a7f9a67cea7a9a28f494b3fc4e7716c9`、`5dc85c7960fa734c1562aed200509801b81261d96cb8a375175476ea292eb5f5`。六份受審文件另保存同SHA快照，供後續只追加中繼資料回查；16份私人成果輸入再核保持。`independent-gui.json` SHA為`d9e2233acaa24177f7ce400e5dee44464b40ee8d5f698606883aeb616857c29c`，完整路線`revision-7/attempt-1/report.json` SHA為`11316a1d09f0cf8ef74b5cc27c0ef298b76d9c61e29f58bf4ebbead9e82be1ad`。

私人16,044份技術檔一次封存，`workplace/r948/closed-inputs.json` SHA為`cb249f20337f68438b5de097b869f7faef5a6e27146cfb8f9bd315cb9b55e866`。排除快取與後續GitHub回讀，原版、錄放、截圖與候選不進Git。封存包含各修正版、失敗收據及最終報告，不能把16,044檔全部當作通過案例。UID/GID1000:1000、自有容器0，正式程式及SDK保持。
