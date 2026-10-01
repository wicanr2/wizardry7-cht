# 022：前端輸入錄放的呼叫邊界

日期2026-10-01；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)；
工作入口[M2](../goal/M2-playable-frontend.md)，規格[013草稿](../spec/013-input-recording.md)。
本節是現有Go程式的唯讀查證，尚無本案正式錄放實作或三次重播收據。

## 輸入與定位

本案最近已推送`d466750`，本機複合名稱實作待獨立驗收；dosgolem基線
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
