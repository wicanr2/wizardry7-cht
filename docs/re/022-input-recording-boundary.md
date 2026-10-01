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
