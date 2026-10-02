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
