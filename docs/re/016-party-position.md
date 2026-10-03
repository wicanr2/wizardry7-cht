# 016 — 隊伍面向與位置（即時地圖前置，#36）

日期：2026-09-30。探針 `cmd/posprobe`（只讀），沿 `cmd/wiz7 -autoplay` 路線進迷宮，
點擊翻完開場旁白（第 6 次點擊無新印字＝羅盤出現）後做動作序列，每步拍原版記憶體。
原始輸出 `workplace/bake/r666-posprobe5.txt`、`r666-posprobe6.txt`、`r666-walk1.txt`–`r666-walk3.txt`。

線性位址換段：DS＝1106（`docs/re/014` 探針量到的印字時 DS）。位址對 golem 載入配置與本重印版成立。

## 結論

| # | 事實 | 等級 |
|---|---|---|
| P1 | 動作有效的正對照：翻完旁白後，轉向、前進、後退各使 A000 變化約 17,500–20,600 px；撞牆時約 118–154 px 且位置不變 | confirmed |
| P2 | 面向＝byte `1106:82EE`（線性 1934E），值 0–3；右轉 +1（3→0 繞回），左轉 −1；前進後退不變 | confirmed（兩組序列、共 6 次右轉、2 次左轉） |
| P3 | word `1106:82F0`（線性 19350）＝面向 × 2 | confirmed（同上） |
| P4 | 目前位置記錄＝`1106:8F3C`（線性 19F9C）起 6 個 word：`[+0, +2, +4, +6, +8, +10]`，起點值 `[0, 11, 0, 4, 56, 92]` | confirmed（位置） |
| P5 | 面向 3 前進：`+6` 與 `+10` 同時 −1；往面向 0 的方向移動：`+4` 與 `+8` 同時 +1；後退回原值 | confirmed（兩組序列） |
| P6 | 初測同區段內，`+8`／`+10` 與 `+4`／`+6` 一起變化；推測前者為全域座標、後者為區塊內座標 | 假說；後續單一接縫見P12，完整格式未證 |
| P7 | 初測`+2`＝11不變；推測為地圖或區域編號 | 假說；後續跨界值見P12，區域名稱未證 |
| P8 | `1106:8F30`（線性 19F90）起 6 個 word 在每次移動後等於移動前的位置記錄（上一步） | 強證據 |
| P9 | 面向 1 前進：`+6`／`+10` +1（(56,91)→(56,92)）；面向 2 前進兩次都撞牆，但面向 2 後退使 `+8` +1，故面向 2 前進為 `+8` −1。四向：0＝`+8`+1、1＝`+10`+1、2＝`+8`−1、3＝`+10`−1 | confirmed（面向 2 由後退反推） |
| P10 | 走到 (57,91) 觸發劇情旁白，需點擊翻頁才能繼續操作 | confirmed（`r666-pos-5.png`） |
| P11 | 上一步記錄的 `+0`＝那次移動時的面向（移動時面向 2 → 記 2；面向 1 → 記 1） | 強證據 |
| P12 | R759面向1正常前進，`[0,10,1,7,57,87]`→`[0,11,1,0,57,88]`；局部`+6`由7繞至0、`+2`由10至11，全域候選`+10`連續加1 | confirmed（此數值接縫）；區塊大小、區域語意及全部連接仍未知 |

## 對即時地圖的意義

- 自動繪圖用到 P2、P4、P5、P9：以 `(+8, +10)` 當格座標、`1934E` 當面向，記錄走過的格（spec 005）。
- 整張地圖需要迷宮資料格式，另案。
- P12量到一個正常接縫；尚不能推廣全部區塊大小或連接，整張地圖格式仍未完成。

## R759／R764：後續正常接縫

來源`workplace/r759/report.json` SHA-256
`0c63d566677d93b267c76121413989a7d24c46f4294043a5d1724f00b82f0a94`；
DS.EXE雜湊與本文基線相同，執行器dosgolem `d16f03a2`，位置仍按runtime DS=1106及線性19F9C觀察。
4050→4065M的P12接縫有212對A/B及204張前綴無差異，沒有座標注入或新狀態匯出。
R764再前進至4110M `[0,11,1,3,57,91]`，31對A/B及30張前綴無差。
這次未觸發P10在同座標曾見的旁白；來源檢查點不同、原因未知，保留P10原收據。
完整路線及角色查證回連[研究019](019-catalog-migration.md)，不由數字11推論已進新城。

<a id="r877-map-resource-boundary"></a>

## R877：地圖資源檔名與兩個呼叫入口

2026-10-03；[#36](https://github.com/wicanr2/wizardry7-cht/issues/36)，[目標](../goal/R877-map-resource-boundary.md)。
本輪只做靜態定位，沒有執行原版、移動、地圖解析或正式UI變更。

### 輸入與工具

| 輸入 | SHA-256 |
|---|---|
| DS.EXE | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |
| MAZEDATA.VGA，95,802 bytes | `af2f15381dbcf2d1595a8c3d496a2e70f7f9e778fa490ad2e0433e4e0eb52091` |
| VMAZE.OVR | `887fe601d8d8bc16c44eccd83555355e4947c671b9360a478952873c5bf1abf6` |
| VMEXT.OVR | `252536e7d9ec4d33a315945b6e17faaa98960ee02b5377afac7ad66490a6587f` |
| VMAZE.code.i64原始資料庫 | `aa134943411c1384102e3821b4d2d1169f90b4d7ca92da711b2e329b48b5e036` |
| VMEXT.code.i64原始資料庫 | `83a5a11165855d3030a773f3999347e0a1746b068774e41795842e77c32d273e` |

主工具IDA Pro 9.4，Python3.12.3，既有`ida-pro-9.4-idapython:locked-v1`，映像ID
`6f6d59af49d0008c4109a5295b5f374bdc007e2d1ab28cb9de08779584de2780`。
最小探針核DS.EXE身分、139個函式、UID/GID1000:1000及既有HOME設定，JSON有效才採認。
所有查詢對從唯讀原始庫新複製的獨立資料庫；不更名、不重種正式庫、不改原版。

### 定位結果

| 證據 | 結果與等級 |
|---|---|
| 精確檔名搜尋 | DS.EXE及13個OVR只搜尋三個檔名，不抽文字表。MAZEDATA只命中VMAZE與VMEXT。confirmed只限原始位元組搜尋。 |
| VMAZE檔名 | 原始OVR檔案偏移`0xAC76`，程式碼區長`0xAC60`，附加區由檔案`0xAC6E`開始，檔名在附加區`+8`。confirmed。 |
| VMEXT檔名 | 原始OVR檔案偏移`0x5259`，程式碼區長`0x524B`，附加區由檔案`0x5259`開始，檔名在附加區`+0`。confirmed。 |
| 指標候選 | 檔頭`+6`的word為`0x9098`。若它是附加區資料基址，兩檔名候選分別為`0x90A0`與`0x9098`。載入映射是假說，尚無runtime目的位址。 |
| VMAZE呼叫 | IDA raw overlay EA `0x0489`的`b8 00 00`與`0x048C`的`50`先push 0，`0x048D`的`b8 a0 90`與`0x0490`的`50`再push `0x90A0`，`0x0491`的`e8 5c fe`指向同一程式碼區`0x02F0`。指令與相對目標confirmed；它接收該檔名仍為強推論。 |
| VMEXT呼叫 | IDA raw overlay EA `0x0301`的`b8 00 00`與`0x0304`的`50`先push 0，`0x0305`的`b8 98 90`與`0x0308`的`50`再push `0x9098`，`0x0309`的`e8 02 ff`指向原名`sub_20E`的`0x020E`。指令與相對目標confirmed；檔名映射仍為強推論。 |
| 共同操作序列 | 兩個callee都使用`[bp+4]`及`[bp+6]`，保留返回值於`[bp-2]`並比較`0xFFFF`。原始near-call數值目標含`0xD153`、`0xF8FB`、`0xF004`、`0xFD97`；實際載入後的服務身分未知，不能直接稱開檔／讀檔／關檔或套runtime段。 |
| 資料庫邊界 | 既有`.code.i64`不含附加區；檔名搜尋空不是原版無引用證據。既有線性種子在部分callee分支界線及函式歸屬不可靠，名稱只供導覽。直接解碼保留先前項目／名稱，不修原始庫。 |

IDA raw overlay EA等於程式碼區偏移；本輪103條直接解碼逐條與原始OVR的`14+EA`位元組相同。
相對call獨立按`(EA+3+signed16(displacement)) & 0xFFFF`核目標，沒有把數值目標當runtime位址。
完整收據在`workplace/r877/`：`raw-filename-anchors.json`、`pointer-candidates.json`、
`VMAZE_decode/exact-decode.json`、`VMEXT_decode/exact-decode.json`及`verified.json`。
解碼範圍VMAZE為`0x02F0–0x0353`及`0x0489–0x04B0`，VMEXT為`0x020E–0x0254`及`0x0301–0x0328`，均為右端不含；第一段請求到`0x0352`，最後一條指令解碼至`0x0353`。

### 下一個必要切片

先量原版loader把附加區放到哪裡，以及上述near-call載入後的實際路由，再只追MAZEDATA的讀取目的與消費端。
尚不能斷定MAZEDATA包含地形、圖像或完整世界；區塊形狀、牆／門、樓層及座標對應皆未知。
自動地圖005保持原已驗範圍，整張地圖未達READY，#36仍OPEN。

封存20份來源與29份產物，索引`workplace/r877/sealed.json`的SHA-256為
`af56650547ec1852aab952f64a86e8d319b3482f4f46e79e9b7016d4c9e7084a`；
`verified.json`為`7280083dc102e8d12cec6923bb6b4c2890575a2aae216d357b01de43fc75f7df`。
原始DS／13個OVR及五份來源IDA資料庫保持，MAZEDATA雜湊與開工時相同。
產物owner1000:1000，R877自有容器0；獨立唯讀審查另列。

R878獨立唯讀審查阻擋0、應改0、建議0，主代理完整核讀兩份報告。
20來源、29產物、103指令及20個相對call目標相符，只採認限定靜態證據，不解除整張地圖的READY閘門。
報告`workplace/r877/review/review-r878.json`及`review-r878.md`的SHA-256依序為
`1b1085ff4d99bd1900463aa58f3e2d9daa31d75b764875724045e2c774877869`及
`ae510d725381e38e459f4b756f4e371e4eda4e3d4824b6c5313fbc9e7616e5c4`。
報告owner1000:1000，審查自有容器0；原版執行、地圖實作及玩家驗收未由本輪補齊。

<a id="r879-map-loader-mapping"></a>

## R879：VMAZE檔名指標的快照比對

2026-10-03；[#36](https://github.com/wicanr2/wizardry7-cht/issues/36)，[目標](../goal/R879-map-loader-mapping.md)。
本輪載入既有合法迷宮檢查點，原版執行0步；不注入呼叫，不匯出記憶體或鄰文。

### 輸入與觀察

- DS.EXE、VMAZE.OVR、VMEXT.OVR的雜湊與R877相同。
- `workplace/r733/two-forward.state`為299,208 bytes，SHA-256
  `beecc0c54302d6304628d60f385fb4f3f6748e77c6bc859dfc153f669e337e6c`。
  來源及3,475,000,000步停止點見`workplace/r733/position-progress.json`。
- 使用dosgolem既有`PrivateLoadCheckpoint`唯讀輔助檔，Go建置覆蓋表只增加本地helper，正式SDK沒有修改。
  工具為Go1.24.13、既有`wizardry7-go-ebiten:latest`，映像ID
  `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
- SDK基底commit為`d16f03a2d8d782024d470606e1f4cb5272b89cd6`，含既有未提交修改，不能只用commit指認建置輸入。
  執行後另記89個SDK來源及helper雜湊於`workplace/r879/toolchain.json`，SHA-256為
  `5b4e1773d7985b874aab5a440aec18793d709ae759c5ad28be592dde70b32906`；來源及helper掛載唯讀。
- 遊戲資料段`1106`依本文P2／P4定位。檢查點當下CPU的DS為`312C`、CS為`740C`，不能把driver的當前DS當遊戲資料段。

| runtime位址 | 唯讀觀察 | 等級 |
|---|---|---|
| `1106:90A0`，線性`0x1A100` | 13 bytes與`MAZEDATA.VGA`加NUL完全相同 | confirmed，限此快照 |
| `1106:9098`，線性`0x1A0F8` | 同一13-byte精確檔名比對不相符，不輸出所讀內容 | confirmed，限比對結果；不證VMEXT載入映射錯誤 |
| VMAZE附加區對應 | 原始檔案`0xAC76`及指令的`0x90A0`與快照檔名相合 | 強推論；loader寫入端尚未實測 |

只有兩次13-byte讀取，共26 bytes。觀察前後步數與全部CPU暫存器相同；
未讀整份RAM作掃描或雜湊，沒有記憶體未變的全量A/B聲明。
零步快照不能證明near-call服務、實際開檔／讀取目的、VMAZE程式碼載入段或地圖格式。
原版程式及檢查點雜湊未變，完整地圖仍未READY。

### 本機收據

探針與收據在`workplace/r879/`，`report.json`的SHA-256為
`157d651cfea19b4485ad44e097ca3f25135e1d89ffdc36772225b2978f699133`；
`sealed.json`為`0ec5b007c13a57a95307cc02b30ceb1ce73cbca3a62d9b189c59f9c7cff217c5`。
封存4份來源、6份產物，產物owner1000:1000；獨立唯讀審查另列。

R880獨立唯讀審查阻擋0、應改0、建議0，主代理完整核讀兩報告。
4來源、6產物與另記89個SDK來源及helper雜湊相符；沒有重跑產品或解析state正文。
僅採認限定快照，loader、服務與完整地圖資格保持未知。
`workplace/r879/review/review-r880.json`及`review-r880.md`的SHA-256依序為
`b1e33eed9086cd89758fdf48347a3b71f906ae02ba372b2b8e33562861f03a19`及
`efbe57c15a57ad4de1e96ac26998ba0b4c4786e1d971d32d147eb3e95dab02ee`。

<a id="r881-map-loader-route"></a>

## R881：共同載入器、VMAZE程式碼與VGA資源路由

2026-10-03；[#36](https://github.com/wicanr2/wizardry7-cht/issues/36)，[目標](../goal/R881-map-loader-route.md)。
本輪補上原版載入位置及服務路由。沒有實際觸發MAZEDATA讀取、取得完整地圖格式或修改正式程式。

### 玩家路線與停止線

- 合法R733檢查點與原版DS基線沿R879。觀察開／關分支先執行10M步穩定，再正常點DISK `(212,166)`，
  固定停止點為3,475M、3,485M及3,511M。三對記憶體、indexed、palette與CPU暫存器SHA相同，三對PNG逐位元同。
  隊伍位置保持`[0,11,5,3,61,91]`、面向0；指定VMAZE／VMEXT／VPOPS／MAZEDATA的讀取紀錄為0。
- CONFIGURATION `(160,60)`的點擊被自動核准兩次拒絕，沒有送入遊戲。已詢問明示授權；兩分支正常退出，沒有選SAVE／LOAD。
- 獨立無輸入冷開機先量VINIT載入。防拷來源閘門在100k指令批次末停止於400k，兩對摘要相同，沒有輸出正文或PNG。
  不能稱指令即刻停止，原版本身仍照正常資料路徑執行；探針沒有讀訊息正文或作答。
- 另以已量到的啟動前綴執行固定1,000步，在VINIT開檔時最多讀四個BP框架共16 bytes，僅記返回位址。
  開／關觀察的記憶體、暫存器及讀取紀錄相同；此分支沒有玩家輸入、防拷正文讀取或遊戲進度還原。

### OVR檔頭與共同載入器

工具為dosgolem既有修改版，基底`d16f03a2d8d782024d470606e1f4cb5272b89cd6`；
89個SDK來源及helper身分按`workplace/r881/toolchain-before.json`核對，正式SDK未改。
Go1.24.13，映像沿R879；IDA Pro9.4／Python3.12.3沿R877映像，查詢全部使用新資料庫副本。
原始DS.EXE.i64及ovr740c.bin.i64的SHA與R877相同，原始庫保持。

VINIT開檔於350步，當下`0110:4B1C`；16-byte返回鏈為`4977 → 41B5 → 4184 → 4D62`，
由此定位原名`sub_14186`及其呼叫端`sub_14179`。IDA resident EA為`0x10000+runtime IP`，
與raw OVR的程式碼偏移分開記錄。

| 檔頭偏移 | 載入器的原始定位 | 結論與等級 |
|---|---|---|
| `+0` | EA `0x141E6`比較`[bp-12h]`與`0x00F2` | 二位元組magic，confirmed；舊002將首四位元組一起列為magic不完整 |
| `+2` | EA `0x1420C`檢查`[bp-10h]`，`0x14225`傳給`sub_142B4` | 程式碼目的偏移；13個OVR皆`0x5055`。confirmed於檔頭／控制流，VINIT實際讀取目的見下 |
| `+4` | EA `0x14222`傳入`[bp-0Eh]` | 程式碼長度，confirmed；VINIT為5,641 bytes |
| `+6`、`+8` | EA `0x14236`／`0x14239`傳入`[bp-0Ah]`及`[bp-0Ch]`，`0x1423F`呼叫讀取helper | 附加區目的偏移與長度；13個OVR目的皆`0x9098`。confirmed於檔頭／控制流及VINIT讀取 |
| `+10` | EA `0x141FB`只見其加到容量檢查 | 完整語意未知，不命名為fixup或BSS |
| `+12` | EA `0x14259`回傳`[bp-6]`，`0x14184`執行`jmp ax` | 入口控制值，confirmed於靜態控制流；舊checksum假說不適用 |

VINIT實際三次讀取為：409步讀檔案`0`起14 bytes到`1106:FFCA`；
440步讀檔案`14`起5,641 bytes到`0110:5055`；461步讀檔案`5655`起288 bytes到`1106:9098`。
全部為DOS現成讀取元資料，不讀其正文。13個OVR的檔案長度均等於`14+code_size+data_size`。
這補上R877／R879的資料基址假說，但VINIT收據不能單獨宣稱VMAZE資源呼叫已執行。

### VMAZE與資源服務定位

| 證據 | 定位與等級 |
|---|---|
| 程式碼快照 | 同R733合法檢查點，raw VMAZE `0x0489`及`0x02F0`各16 bytes，精確符合`0110:54DE`及`0110:5345`。共32-byte唯讀比對，只出布林，前後Step／regs同。confirmed，限兩個anchor。 |
| near-call映射 | 加檔頭code base `0x5055`並取16-bit IP，raw `D153/F8FB/F004/FD97/F3FD`分別指向runtime `21A8/4950/4059/4DEC/4452`。confirmed於算式與原始指令；不代表所有callee本輪已執行。 |
| 檔名與開關檔 | 原名`sub_121A8`回傳組合字串位置，`sub_14950`依flags走DOS開檔wrapper；目前候選呼叫傳flags=0。`sub_14DEC`含AH=3Eh關檔。靜態服務序列confirmed；本輪沒有MAZEDATA開檔收據。 |
| driver橋接 | EA `0x1405D`為`call cs:dword_12DA7`。暖快照只讀`0110:2DA7`四位元組，得到`740C:0123`，前後Step／regs同。指標confirmed，限此快照。 |
| VGA thunk | runtime dump IDA EA `0x0123`的`e8 88 0b cb`與原始VGA.DRV檔案`0x23`相同，near-call指向runtime `0x0CAE`，原始檔案`0x0BAE`。confirmed於位元組／算式。 |
| driver消費端 | 從CAE有界解碼103條指令；含bank索引檢查、DOS AH=3Fh讀取及資料目錄長度計算。經VGA圖形資源服務載入是強證據，不能因此宣稱取得地形、牆／門、樓層或世界地圖。 |

driver原始庫的舊項目與disasm部分失準，保存它們只作導覽。
本輪以IDA直接解碼的mnemonic／operands及原始VGA.DRV逐條位元組核對，不沿用那些假函式／變數名作語意。
runtime dump EA到原始VGA.DRV檔案偏移的`EA-0x100`只在已比對的區段成立。
MAZEDATA的實際payload格式及讀取仍未驗證；下一切片應追隊伍位置的地形消費端，不由檔名推論世界地圖來源。

### 驗證與封存

432筆resident解碼紀錄與103條driver指令核對通過，resident紀錄含重複查詢。
resident僅一處差異：MZ宣告的image `0x4B3B`重定位word，原始`0x0FF6`、IDA為`0x1FF6`，
正好增加IDA載入基址`0x1000` paragraphs。其餘指令與原始檔逐位元同；driver全部103條與原檔同。
沒有改資料庫項目、名稱、正式程式、原版素材、譯文或字型。

VINIT.OVR的SHA、原始VGA.DRV及其餘來源在`workplace/r881/sealed.json`，
20來源、58產物及另記89個SDK來源已核對；來源保持、產物owner1000:1000、自有容器0。
封存索引SHA-256為`a5c93249cabbedfab0eb157d38a045ce71d184974dd228b33ec2827df724bf08`，
`verified.json`為`1c0832755f21005c0ae12f2451d81c8a482ea2dac695d0cf7b12547db940330d`。
獨立唯讀審查另列。005仍READY，整張地圖格式未READY，#36保持OPEN。

R882獨立唯讀審查阻擋0、應改0、建議0，主代理完整核讀JSON及Markdown。
20來源、58產物、89個SDK來源、31個相對呼叫與限定文件回填相符；沒有重跑產品或IDA。
只採認載入器、快照及靜態資源路由，不升級MAZEDATA實際讀取或完整地圖的驗收資格。
`workplace/r881/review/review-r882.json`及`review-r882.md`的SHA-256依序為
`398ce2aef15df8b4ca5d41b584f7189c7692cc0de68be20e4a98ece0101710ff`及
`5ebfd5315f1143e57f1c56b575f39f733c531b3536ae13e6169e80750102d8a8`。
報告owner1000:1000；全工作根root-owned／誤建.md目錄0，自有研究及審查容器0。

<a id="r883-position-terrain-consumer"></a>

## R883：區塊座標與SCENARIO資料路由

2026-10-03；[#36](https://github.com/wicanr2/wizardry7-cht/issues/36)，[目標](../goal/R883-position-terrain-consumer.md)。
本輪只查靜態指令與已知檔名，不執行原版、讀地圖正文或修改正式程式。
設定頁操作仍待明示授權，沒有重試被拒的點擊。

### 輸入與定位基準

DS.EXE及VMAZE.OVR雜湊沿R877。原始DS.EXE.i64為
`477f73a17b7abcaf3eef70d9baae16c77d9730f3bf7e0b4c0c9c9ed91c612412`，
VMAZE.code.i64為`aa134943411c1384102e3821b4d2d1169f90b4d7ca92da711b2e329b48b5e036`。
所有IDA查詢在新副本；IDA940／Python3.12.3、UID/GID1000:1000及映像沿R877，驗證工具沿R879。
原始程式與資料庫四份來源保持。

VMAZE的IDA EA是raw code區偏移，對應原始檔案`14+EA`；
runtime IP依R881的`(EA+0x5055)&0xFFFF`換算。resident IDA EA為`0x10000+runtime IP`。
以下`DS:`操作數保留原始偏移，不拿IDA的預設DS選擇子當遊戲runtime DS。
runtime位置欄位仍依本文P2／P4，不由靜態資料庫重新命名。

### 區塊與座標資料流

| 原始定位 | 靜態事實 | 等級及限制 |
|---|---|---|
| VMAZE原名`sub_9BD`，EA `0x9BD–0xA60` | 在EA `0x9D3`／`0x9E9`傳入選擇值2／3、緩衝`DS:0E14`／`DS:0E12`、`DS:830C`及模式0，呼叫raw `0xB62E` | 指令／引數confirmed；按R881基址指向runtime `0110:0683`，不是VGA資源服務 |
| 同段EA `0x9EF–0xA34` | 從`DS:0E1A/0E1C`及`DS:0E16/0E18`兩個遠指標加`DS:830C<<7`，各傳入128與目的`DS:8B7E`及`DS:8AFE`，呼叫raw `0xAFBE` | 數值及引數confirmed；資料搬移為強推論，callee本輪未核；兩段資料的格子意義未知 |
| 同段EA `0xA3A–0xA59` | `BX=[0E14]`、`SI=[8F3E]`；將無號byte `[BX+SI+490h]`加`[8F42]`寫到`[8F46]`，`[BX+SI+4A0h]`加`[8F40]`寫到`[8F44]` | 控制流confirmed；結合P4／P5支持區塊原點加局部座標的強證據，新增路徑未實跑 |
| 原名`sub_677E`，EA `0x677E–0x67B8` | 用上述兩個原點檢查輸入座標是否介於原點與原點+7，兩軸均含邊界；相合回傳區塊索引，否則`FFFF` | 8×8範圍判定confirmed於靜態控制流；不證所有地圖區塊相連或有效 |
| 原名`sub_67B8`，EA `0x67B8–0x6874` | 先查目前區塊與`[8F48]`，再以模16循環查區塊；成功後減去原點，寫回區塊及兩個局部座標指標，失敗回傳1 | 模16及轉換流程confirmed於靜態控制流；循環有界須目前區塊為0–15，兩次先查另計，未驗非法初態或正常失敗分支 |
| 原名`sub_1294`，EA `0x1294–0x12C5` | 複製目前位置的後五個word到前一位置記錄，另將`[830C]`寫`[8F3C]`、面向`[82EE]`寫`[8F30]` | 指令confirmed；補充P8／P11的寫入端線索，不推廣所有狀態的欄位語意 |

`sub_12FA`的EA `0x1384–0x1395`把全域座標候選及三個位置指標傳給`sub_67B8`。
既有原始名稱與函式界線只供導覽；不重命名欄位或改原版位置。
`sub_6CC`只查前512 bytes，可見`[0E14]+90h`、`+290h`的資料更新及區塊座標組合，
後段未補齊。這些緩衝的牆／門、地形及事件位元意義仍未知，不能僅憑索引算式稱完成地圖解析。

### 檔案來源路由

- resident原名`sub_10683`用選擇值2／3索引`DS:5244`，再依選擇值取得`DS:5524`的記錄長度及`DS:521C/521E`的基址。
  計算檔案位置後，EA `0x10727`將`[DS:0E22]`傳入定位helper，EA `0x1074A`等讀取分支傳同一控制代碼。
- 原生函式指令的`0E22`直接操作數查得七筆讀寫候選，唯一寫入為原名`sub_121EF`的EA `0x125D8`。
  前段`0x125C7`傳入檔名偏移`0x0BC0`，經`sub_121A8`組合路徑與`sub_14950`開檔後保存回傳值。
- DS.EXE的MZ header為512 bytes，宣告SS=`0x0FF6`。只在檔案`0x10D20`比對13 bytes，
  與已知`SCENARIO.DBS`加NUL精確相同，沒有輸出鄰文。此偏移按`512+0x0FF6*16+0x0BC0`得到。
- 字面檔名、開檔返回值及讀取控制代碼的靜態路由相合。SCENARIO記錄是目前位置／區塊資料來源的強證據；
  本輪沒有實際讀取目的、完整記錄格式或牆／門判定收據，不能稱整張地圖已驗收。

### 驗證、工具限制與封存

八份IDA查詢輸入身分相符。313筆resident及583筆VMAZE解碼紀錄逐條核對原始位元組，
76筆相對call算式相符；紀錄含重複查詢及線性解碼的內嵌資料，不將總數稱為有效可達指令數。
另有47筆舊項目邊界導覽。resident的七word跳表獨立取值核對，
以EA `0x106E2`的真控制流查詢補足線性解碼跨過跳表的錯誤界線。
舊IDA item／disasm及raw相同數值位移不單獨作語意證據，不改正式庫來讓它們吻合。

最初兩份查詢沒有JSON，已確認terminal；捕捉例外後定位`idc.str2reg`不存在。
讀映像內API原始碼，改用`ida_idp.str2reg`後從原始庫新副本乾淨重跑，兩份JSON通過。
失敗及私人日誌保留本機，不輸出授權環境或加入版控。

封存4來源、42產物，owner1000:1000；原版執行0、正式程式變更0、自有研究容器0。
`workplace/r883/sealed.json`的SHA-256為
`21a7d8f940edaf450e61220c325600c247ea77faa485f0afe5e30cbbbbf78e3d`，
`verified.json`為`602c1ccffb376fce78e5bb485b43d20906df7f7c73888fa809484217f505c065`。
下一必要切片是SCENARIO的選擇值2／3記錄與格子消費端，按實際玩家位置抽樣驗證。
005保持READY，整張地圖未READY，#36仍OPEN；獨立唯讀審查另列。

R884獨立唯讀審查阻擋0、應改0、建議0，主代理完整核讀兩報告。
4來源／42產物、八份查詢、限定資料流及目前文件回填相符；只採靜態證據，不提升完整地圖資格。
`workplace/r883/review/review-r884.json`及`review-r884.md`的SHA-256依序為
`3b22925b65c55a52ec291bed63d01f61b4d5734441ec853d8519fafc91507384`及
`135447da520691d3c5f1903b7dd62e75fe86f5fb44895c7ebfad75b255216685`。
報告owner1000:1000；全工作根275274項衛生0，自有研究及審查容器0。
