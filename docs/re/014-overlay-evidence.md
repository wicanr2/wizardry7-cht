# 014 — 不寫原版的覆繪：證據（#31）

日期：2026-09-29。探針 `cmd/ovlprobe`（純英文、不寫原版記憶體），
路線 START NEW GAME（點擊 250M:160,62／253M:170,153／255M:90,168／257M:90,168／265M:160,101／276M:82,169），
停在 286M（防拷畫面）。原始輸出 `workplace/bake/r663-probe1.txt`。

## 結論

| # | 事實 | 等級 |
|---|---|---|
| E1 | 前端寫 `0xC3` 到 0110:3895 會改變原版行為：防拷第二行（1052）進入時 driver y（740C:0297）英文模式 164、中文模式 177；其他三行相同 | confirmed（R662 兩次實跑對照） |
| E2 | `sub_13895` 一律 near call；返回點依呼叫端不同（539b、51be、116d、0ea2、0fcc、1c21、1c71、b3d1、b454、b4cc、b56a…），返回時 SP＝進入時 SP＋2 | confirmed |
| E3 | 入口機器碼 `55 8b ec fc 83 ec 04 56 06 1e 8c d8 8e c0 8b 76 04 …`：`[bp+4]`＝字串 off、`[bp+6]`＝字型 ID，逐字 `call` 0110:3877 | confirmed |
| E4 | 印字只改離屏段（6289）的一個矩形：左上＝進入時 (x0,y0)，右＝返回時 x1−2，高＝字型 0 為 5 列；改到的像素只有筆畫（字形透明，不填底） | confirmed（43 次印字，差異 bbox 全數吻合） |
| E5 | 離屏段＝`740C:01E5` 的值＝blit（740C:05F5）時的 DS，本路線恆為 0x6289 | confirmed |
| E6 | 服務 11 入口（740C:012F）的 DL：第一次呼叫恆為 00（或 06），第二次起才是字色 | confirmed；成因未知 |
| E7 | 防拷畫面時，離屏的選單區已被清成黑底，只有對話框區 blit 到 A000；A000 仍顯示選單字。離屏與 A000 差 6,926 px | confirmed（`r663-offscreen.png`／`r663-a000.png`） |
| E8 | 對話框區內離屏與 A000 的唯一差異是游標：以 `Mouse()`（遊戲座標）為左上、10×10 | confirmed（此畫面的游標）；其他游標形狀未量 |
| E9 | 純英文模式防拷四行進入時 (x0,y0)：SECURITY CLEARANCE (106,148)、PLEASE… (70,158)、MANUAL… (64,164)、LOCATED… (70,170) | confirmed（`r663-probe1.txt` 步數 276125711–276170620 四筆） |

## 推論

- 由 E7：「覆繪是否仍在畫面上」只能以 A000 判斷，離屏會被重用。
- 由 E4：印字前後的離屏差異就是英文筆畫；印字前的離屏矩形就是該處背景。
- 由 E6：字色改由「差異像素的值」取得（筆畫像素的眾數），不依賴 DL。
- 由 E1：任何跳過原版常式的做法都會改變原版狀態；覆繪只能讀。

## 追記一（同日）：spec 001 實作驗收

`cmd/wiz7 -hdonly -autoplay -hd-at 272000000,286000000,330000000`，中文模式（`-hd-scale 2` 與 `3`）與
`-xlate ""` 各一次，`tools/re/check_overlay.py` 比對，失敗 0 項：

| 步數（實際） | 原版記憶體／A000／調色盤雜湊 | 英文殘留 | 事件統計（建立／上屏／刪除／無差異／巢狀／無離屏／存活） |
|---|---|---|---|
| 272M（280,000,000） | 中英相同 | 0／2,651 | 72／27／57／0／0／0／15 |
| 286M（286,000,000） | 中英相同 | 0／2,659 | 72／27／57／0／0／0／15 |
| 330M（333,000,000） | 中英相同 | 0／881 | 80／34／78／0／0／0／2 |

- 防拷四行 (x0,y0) 與 E9 相同；選單七則、THESUS、FIGHTER 上屏；字色等於筆畫眾數。
- 1052 原有 dy＝13，正是 E1 的位移（177−164）；覆繪不寫原版後已移除。
- 2× 畫布上行距 6 的四行中文互相壓疊（15 px 字高＞12 px 行距）；3× 行距 18 px，四行清楚。
- 未量到：巢狀捨棄、無離屏兩個分支（實跑計數 0）；互動視窗 `Draw` 路徑未截圖驗證（只驗 `hdFrame`）。
- THESUS 顯示為「第六」是查表鍵衝突，與覆繪無關（#32，spec 002 已修）。

## 追記二（同日）：互動視窗 Draw 路徑

`-draw-at 286000000,340000000`（Xvfb、`-scale 3`）：`Draw` 完成後以 `ReadPixels` 存下 960×600 的實際視窗畫面，
與 `hdFrame` 3× 目視一致（選單、防拷四行、旁白中文；無英文殘留）。截圖在機器到達步數後的下一幀，
輸入框游標狀態不同，故不做逐像素比對。

## R821／R822：一般填槽缺字回退

2026-10-02；[#44](https://github.com/wicanr2/wizardry7-cht/issues/44)，關聯[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
定位為自有Go來源，不是IDA位址；一般印字入口沿研究013與spec001的既有證據。
本輪只讀Go／規格與自製字模，未讀原版、翻譯正文、RAM、state或防拷資料。
工具Go1.24.13／Python3.11.2，既有Docker image
`sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。

已證（來源與自製案例）：修前缺數字時略過缺項、缺序數時整槽移除、nil序數字模仍被接受。
`TestOldSlotMissingReproduction`的PASS代表重現這三個缺口，字模完整性false，不是修正驗收。
固定修前區塊已獨立核對逐字擷取相符，完整數字對照及來源未變也通過。
既定字型007 R6要求缺字退回英文，故本輪不需新增產品決策。

| 本機收據（均在workplace/） | SHA-256 |
|---|---|
| r821/baseline-report.json | `fc51e6a2d02be34a4048a0546de0026e522221010d71f8a73d5982f34df9ed1c` |
| r821/baseline-output-manifest.json（8項） | `5986ba882024e5665713b6f134d922d3909baa50429dad1e22f8c65af7de20af` |
| r821/implementation-report.json | `99796c7bedf5eea28676cb3a8bde5788da3fb7a223283cbb7d6a04aeacfb1d22` |
| r821/implementation-source-manifest.json（75項） | `bc65ce0a6cfc7e98c9cf2ffece85e6cad396583c1cad67197d07cc95c5d73dc3` |
| r821/implementation-output-manifest.json（16項） | `230db4a564b6cf3535ebe9028dec94fc6c69d06475126d82c9a8e5b8287d9951` |
| r822/correction-review.json | `1171cf5b8215c11038bd4619635d4ef79d098c88cb6159a69fe66073c9fd216d` |
| r822/final-output-manifest.json（6項規格審查收據） | `9705679e85a4ef98e1bc42054ebf122506c8c6ce38cc9483793e7ad9be36b8c4` |
| r822/implementation-review.json | `06cbb17d2e0693e002313c601c3db6e8e0f41d6b7701862646ba2239cb75d201` |
| r822/implementation-source-manifest.json（12項） | `a7c16be4146cb9c1e994b78cfb4ed538680768477466de231d7681379d873786` |
| r822/implementation-output-manifest.json（3項） | `31cbcf19c7649158c8ea46ffe8535459ad4705798a74bbc394e743ac00f7ba14` |

規格015經R822唯讀審查後先升READY，再實作`expandSlots`；原R2文字應改已解除，
明列ordinals key存在但nil／空不得落入digits，已知序數缺key亦整則失敗。
caller缺字return在建立pendingPrint、取離屏背景及註冊本次返回之前；不交付部分中文事件。
英文鍵、完整字模的映射順序、D邊界、flow及複合物品路由維持原行為，原版資料不寫入。

正式三群組16自製配置（完整1／缺字10／既有邊界5）與全部前端競爭回歸通過，
`WIZ7_FE_TIMEOUT=60s tools/fe.sh xvfb 'go test -race -count=1 ./cmd/wiz7'`套件5.042秒。
實際helper僅更換package名稱的隔離副本與固定修前區塊比對1980個完整配置，
字模／Dy等結果全同、輸入未改，隔離競爭群組1.187秒。
自製識別值不是原版或第三方字模；沒有新原版A/B、GUI或玩家缺字故障驗收，不升CONFORMED。

實作獨立複核阻擋0、應改0；75項來源、16實作輸出、8修前輸出相符。
主代理完整閱讀報告並核對12項複核來源與3項輸出；字詞邊界及首大寫段仍是既有啟發式，
不冒稱隊員結構身分。main.go修前SHA `90a1b10947b27377773c7b1396baf02d120eaae3e22d2d2c8f0ee2154b5c170b`，
修後 `c9d5016adfca9a6f81679c38ec6204687d05943b320ee0b72df2bd2b13a38355`；
來源索引保留規格訂正、READY與實作前後分類，既有收據未覆寫。
所有輸出1000:1000，210698項lstat衛生檢查異常0，自有容器均已清理。
本缺陷按#44有限驗收範圍完成；#28整體規格覆蓋、完整可玩交付及其餘原版驗收仍未完成。
