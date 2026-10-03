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
