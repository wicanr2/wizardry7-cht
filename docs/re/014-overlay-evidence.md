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


<a id="r837-subject-boundary"></a>
## R837：十來源主體顯示邊界（2026-10-02）

入口[R837](../goal/R837-dynamic-subject-display.md)、[015](../spec/015-generic-slot-expansion.md)、
[017](../spec/017-dynamic-subject-source.md)；議題#24／#28。基線5d43d0c。
十來源為1541、1542、1543、1544、1547、1548、1549、2040、2061、5010；
ID只作研究選取及出處，不當正式鍵或執行期路由。正式Go、SDK、catalog及產生器未改。

自有Go helper精確副本的40個合成案例直接匹配0、首詞回退9。只有單詞全大寫主體
的九列回退命中；括號內1547、含空白、混合大小寫與字母數字案例均未命中。
案例名字合法性與原版整串替換格式未驗。九列烘製沒有D槽，5010只有D=[4]數字槽；
十列沒有主體字模或名字token契約。原版可能另成印字事件，不能稱十個實際畫面都漏名字。

原版靜態指令支持[013定位勘誤](013-print-pipeline.md#r837-pipeline-correction)。
IDA線性0x11C1E直接call的記錄指標式已證；名字語意仍強推論。
0x10582只證記錄查找與傳入sub_13E02的引數，供後續十來源窄觀察定位。
目的緩衝語意沿既有decoder證據，不由本輪推定完整記錄容量。

resident及13個overlay僅查IDA既有code heads的cmp/sub立即值5Eh，取得三候選：
VBASE.code 0x6D86、VMAZE.code 0x6C76、VTREA.code 0x0AD6。
VBASE的sub_6D11在末字判斷後截尾，取DS:A4A8並在結尾回存位置；位置重用為強推論，
沒有名字指標證據。零命中不排除switch、字串查找或其他間接路徑，不再盲掃原碼。
位址以各code的IDA seg000起點0為準；OVR檔案偏移為code offset＋14，
不是runtime或DS.EXE線性位址。全部13份code與原OVR的已解碼code區逐位元相同。

| 有限輸入與工具 | SHA-256 |
|---|---|
| DS.EXE | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |
| DS.EXE.i64 | `477f73a17b7abcaf3eef70d9baae16c77d9730f3bf7e0b4c0c9c9ed91c612412` |
| VBASE.OVR | `7919a6eddd0f7b65905ad4ab3a301604cccb8bf0f213ebebc190d5d1227104d6` |
| VBASE.code | `f2af41d5bd0e394f7fab514a4ace8bd0ba8de718efcca5e544879a781fe68e87` |
| IDA9.4 locked-v1映像 | `6f6d59af49d0008c4109a5295b5f374bdc007e2d1ab28cb9de08779584de2780` |
| Go／Python既有工具映像 | `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7` |

私人收據皆在workplace/r837。source-index.json含53項，output-index.json含76項，
各有路徑、長度及SHA；凍結索引不覆寫。

| 收據 | SHA-256 |
|---|---|
| subject-boundary.json | `0228f3c9bf4165e2a339ba0df990bf1a65733b40f5ce5c2cf39e9d3f9ed6eda4` |
| source-index.json | `bed5768821d4395b94ae4b59e2242b350a91e30fde75abcc29aa9f6c415767d5` |
| output-index.json | `3ffcd9722d0828d3acac091d73b4ebb745283caa98ffa4c10748e3645dc66d58` |
| review-r838-report.md | `9938be06689d8ed9f89830302728787f354e78d70cf9b7b734bbfc4ae92adf89` |

R838獨立重算40案例，核對53來源及76凍結輸出，並獨立核對217個不同DS指令位置
及72個VBASE位置，原始bytes全同。新增唯讀MZ頭28 bytes確認header為512 bytes；
DS檔案偏移與IDA線性位址逐項分列於review-r838-check.json。
審查阻擋0、文字應改2、建議1；兩項應改在本節及013／015／017回填。
首詞回退的程式註解建議保持為後續文字整理項目，匹配行為不改。

凍結收據選取描述的追加勘誤：ida-message-fetch.json:8稱sub_10582為sub_13895 caller，
實際本體沒有直接call後者；只證記錄查找及傳入sub_13E02的引數。
ida-vbase-subject.json:8的substitution只是選取標籤；原指令支持行尾標記與位置處理，
不支持名字代入。subject-boundary.json的名字與共同入口結論一律限本節所列範圍，
不含正常畫面漏字或所有單字元必經NUL入口。原收據保持原樣。

017限定來源修正、015限定數字／序數回退仍READY，尚無名字顯示READY契約。
下一步先量實際主體、完整字串與caller，再定契約。本輪沒有新正常遊戲、MSG、手冊、
答案、鄰文、存檔或RAM讀取，沒有新上屏／原版A/B。十一組、TSV、正常存讀檔、
錄放、人耳與三平台封包門檻保持，不能宣稱中文化或可玩交付完成。

文字補審`review-r838-supplement.json`解除兩項應改，阻擋0／未解除應改0；SHA
`66cfd8b1eb7d4c541f5a3042026ea568ce99de0ded341ca66310eaa0a3580f00`。
根核對49未改來源及76凍結產物全同，回填連結通過；215924項lstat異常0、輸出1000:1000。
`final-guard.json` SHA為`14835bd180489259b74143b408ca636c0480758e72e07aa144537dd7c137cd4e`。
本段僅追加已讀補審及清理核對，未覆寫索引、程式或原版收據；自有容器均已清理。

<a id="r839-subject-observation"></a>
## R839：十來源窄觀察與正常暖路線（2026-10-02）

入口[R839](../goal/R839-subject-observation.md)、[017](../spec/017-dynamic-subject-source.md)，
議題#24／#28，基線0d6ca7d。沿既有合法暖checkpoint；原版及輸入state唯讀。
沒有正式Go、SDK、catalog或產生器變更；探針只在workplace/r839。
工具Go1.24.13／Python3.11.2，映像沿本節前列的Go／Python SHA。
來源、原始248檔、暖起點及私有loader的完整路徑與SHA見final-source-index.json。

runtime位址0110:0582先讀SS:SP+2來源ID，排除保護段，再以正文觀察開關及十ID白名單
限定後續引數。0110:0659只採最多64 bytes且與正式十來源英文相符的目的緩衝；
0110:3895再要求已解碼、同DS／目的指標及1M步時間窗。ID只作研究選取，沒有正式路由。
正文觀察關閉仍保留來源ID及保護段檢查。本輪保護段未觸發；不宣稱Click內可即時停止。
同指標護欄可能漏掉另段印字或間接組字，不能據此排除原版另印名字。

| 批次 | 正常輸入 | 停止點對數 | 取得入口／印字入口次數 | 十來源／解碼／同緩衝印字 |
|---|---|---|---|---|
| hardware | 一次PressKey(0x48)，未交付即停止 | 3 | 0／0 | 0／0／0 |
| mouse | 兩次前進箭頭Click(160,162) | 4 | 5／7 | 0／0／0 |
| narrative | 同mouse前綴，再最多三次Click(160,170) | 7 | 54／104 | 0／0／0 |

上述次數是各批正文觀察開啟側，不跨重播相加。三批14對、8個唯一停止點；
mouse四對是narrative前綴。各對RAM、indexed、palette、暫存器SHA、步數、輸入及停止理由相同。
這是限定零命中分支的唯讀觀察比較，沒有中文覆繪或白名單正文正分支驗證。

位置讀數沿005：前進後+8欄61→62，第二次點擊停在花朵旁白。
兩次文字區點擊後六word為[0,10,2,7,58,87]、face=0，畫面顯示RESTING。
第三次文字點擊於3595M回NoResponseError，兩側同狀態停止；不算產品故障。
區域／世界欄位語意維持005假說，不推定傳送、休息規則或完整玩法。

mouse／narrative同暖起點觀察IRQ1向量為0080:0024；SDK StubSeg=0080時keyTick暫不交付。
hardware批沒有向量欄，只量IRQ1Delivered=0；後批觀察不冒稱前批逐點向量量測。
正式前端使用PressKey，但本輪只證此暖條件，沒有冷前端或全原版鍵盤缺陷結論。

所有Click使用NoCursorWait，私有暖checkpoint且非GUI；沒有新存讀還原、冷啟動、
正式音訊／錄放、人耳或跨平台封包驗收。下一條件為同固定路線的RESTING，
不重擲、不強制呼叫、不注入座標。十來源零命中，015／017仍READY，008仍DRAFT。

| 本機收據（workplace/r839/） | SHA-256 |
|---|---|
| report.json | `8dcde708f27582b6fccecfcc9e8d32178d6fb32226e091757f65e913cdb56cc3` |
| final-source-index.json，321項 | `f059ed933dbc9d653b9c429e05852af18c036705a0dfc56460aa9a8e7e50d096` |
| final-output-index.json，53項 | `6ee991e7e1fbc9497d3597a155ccd7625997ebc5f1d2ea083403b359666735b2` |
| review-r840-report.json | `6cc247a6537dec3d23c67adcc7b490f97ce2f19ee74880019a644a670d28cf4e` |
| review-r840-report.md | `691d0172c1d9f8392ce13fe1041157c89c382fb5b95440aa86e8971cff2d3719` |

R840獨立核對321來源／53輸出、55 SDK來源、十列期待值、全部停止點及三張正常畫格，
阻擋0／應改0／建議0；主代理完整讀取兩份報告。只核准記錄有限證據，不核准名字實作。
原PNG、原文、完整輸入與收據只留本機；凍結索引不覆寫。十一組及完整交付門檻保持。

文字補審亦無阻擋／應改／建議，主代理完整讀取；review-r840-supplement.json SHA
`9d46a356d323fdaaeb27ae6a0129f1d1b6ce013566ad7daa5ffc1e02638b4217`。
根核對321來源只有四份批准文件回填，53凍結產物保持；217709項lstat衛生異常0、
六份正式文件UID/GID1000:1000，七個scratch目錄檔案0，自有容器0。
final-guard.json SHA `5058f611afc07545fe964c3dec325cf8b7758a64be50ec2327ff0d6ad6e599c5`。
封存後僅追加本段及WORKLOG驗收文字，不改正式程式、原版或凍結收據。

<a id="r841-rest-route-continuation"></a>
## R841：休息後正常續行（2026-10-02）

入口[R841](../goal/R841-rest-route-continuation.md)、[005](../spec/005-automap.md)，
議題#24／#22；基線17ed16e。沿R839相同合法暖起點、十來源及有限正文護欄。
私有探針只增加最多200M步被動等待及已判圖後的正常翻頁、右轉、前進；
原版、state及正式Go／SDK／catalog均不改，不輸出新state或RAM正文。

每批先重播R839七個停止點，3595M另核對固定RAM SHA
`10d2042d1c805e295793033c3566f16908ba7caacf913f280ff50381e2189ab1`。
正文觀察開／關都保留來源ID與保護段檢查；五批未遇保護段，白名單正分支仍0。
同指標、64byte及1M步護欄沿R839；不以零命中排除另段印字或間接組字。

| 批次 | 新增正常條件與操作 | 每批對數／共享前綴 | 最後步數 | 面向 |
|---|---|---|---|---|
| wait | 無輸入等待200M步，每20M快照 | 17／7 | 3795M | 0 |
| wake1 | 已見醒來第一頁，Click(160,170)一次 | 18／17 | 3821M | 0 |
| wake2 | 已見第二頁，再正常翻頁一次 | 19／18 | 3847M | 0 |
| right | 已見羅盤，Click(180,166)一次 | 20／19 | 3873M | 1 |
| forward | 已見新視角及面向1，Click(160,162)一次 | 21／20 | 3899M | 1 |

五批合計95對含重播前綴，只有21個唯一停止點，其中R839七點、新增14點。
各對RAM、indexed、palette、暫存器SHA及完整snapshot欄位、固定輸入與停止理由相同。
最新forward批取得／印字入口288／568次，十來源／解碼／同緩衝印字0；
各批是獨立重播，不相加入口次數。這只驗零命中觀察，沒有新的中文覆繪或名字正分支驗證。

主代理逐階段看五張既有正常畫格，再安排後續輸入：醒來第一頁、第二頁、羅盤、新視角、
最後同視角。私有observations.json保存判圖與下一個操作的關係，不轉寫完整旁白。
200M是等待預算與取樣終點，不是原版精確休息時長。右轉只改面向0→1，
前進後六word仍[0,10,2,7,58,87]，不推定牆、睡眠flag或傳送規則，本方向不再試。
位置與區域／世界標籤維持005證據等級。數值seed未知，只固定合法同初態與等價輸入，不重擲。

來源初索引326項，最後索引329項，原版248檔及SDK索引55檔SHA保持，
私有loader與既有R780相同，實際編譯探針與受檢來源相同。工具Go1.24.13／Python3.11.2，
Docker映像沿前列Go／Python SHA，根Docker inspect已核對相同。
私有暖checkpoint、NoCursorWait、非GUI資格維持；015／017仍READY、008仍DRAFT。
沒有新冷存讀、正式音訊／錄放、人耳、跨平台或完整可玩驗收，十一組門檻保持。

| 本機收據（workplace/r841/） | SHA-256 |
|---|---|
| report.json | `4fc778d49dc40e1c4744d4574e259246f733101620af8f4dbfcbe29c5d63ad66` |
| final-source-index.json，329項 | `b639b5d61b50ec5ed4796c3ac87a3ed9c19e984f1e8792a500b19952538cd3b8` |
| final-output-index.json，224項 | `75c4e8aae76e4660f082990940db0adbf23a82a5abdb2ff1867cb7bfc95fdd11` |
| final-guard.json | `4898adaa435d2dfd51bc305001ba6d68d7f81e3a8cde711ec13592e9a700a428` |

根凍結核對329來源／224產物全同；219642項lstat衛生異常0、十scratch無檔，
輸出UID/GID1000:1000。原畫格、原文與完整收據只留本機，凍結索引不覆寫。

R842獨立核對329來源／224產物、五批完整前綴、全部95對欄位及五張正常畫格，
阻擋0／應改0／建議0；主代理完整讀取JSON與Markdown。只核准記錄有限零命中證據。
review-r842-report.json SHA `dbcb8de37c7663b69473aeb598b0f47f4697c1bf04cf42878f4845700f4995df`；
review-r842-report.md SHA `178f93c3d83e06fcc7818aa47ad491590935ce39fab4fab5ffcfa86d8273b4bf`。
wait-comparison.json仍用舊r839_prefix_pairs欄名，其餘用prefix_pairs；原始on/off及前綴已直接核對，
不依摘要欄名判定。原摘要與凍結索引保持原樣，不為文字回填重跑遊戲。

三文件文字補審阻擋0／應改0／建議0，主代理完整讀取最終兩報告，
review-r842-supplement.json SHA `df433bea99220e15600248c95cfb8f484cb9c5621c0f55943105b8330aea3f94`。
根doc-guard核對三批准回填、其餘326來源及224產物保持，五正式文件擁有權及新連結通過；
doc-guard.json SHA `f0ee7b0567f50489e2be2911032843ea5a214df97e86bc9ae4cf9a68e5668a4c`。
封存後只追加本段及WORKLOG驗收文字，不改原版、程式或凍結收據；自有容器均已清理。

<a id="r843-next-maze-heading"></a>

## R843：暖迷宮換向續行

2026-10-02，基線1eec97e，目標[R843](../goal/R843-next-maze-heading.md)，沿005已實測正常輸入。
R841完整二十一點前綴與3899M原版RAM SHA核對後，按逐批正常畫格決定下一步。
十來源及有限正文護欄與R839相同；沒有新MSG、鄰文、手冊、答案或原版資料改寫。

| 批次 | 已觀察後的正常操作 | 比較對數／前綴 | 末步數 | 面向／末六word |
|---|---|---:|---:|---|
| next-right | R841羅盤後Click(180,166)一次 | 22／21 | 3925M | 2／[0,10,2,7,58,87] |
| next-forward | 新視角及面向2後Click(160,162)一次 | 23／22 | 3951M | 2／[0,10,1,7,57,87] |
| trail1 | 已見移動及新視角後再前進一次 | 24／23 | 3977M | 2／[0,10,1,7,57,87] |
| branch-right | 面向2未移動即停，正常右轉一次 | 25／24 | 4003M | 3／[0,10,1,7,57,87] |
| branch-forward | 新視角及面向3後只前進一次 | 26／25 | 4029M | 3／[0,10,1,7,57,87] |

五批120對含重播前綴，只有26唯一停止點，其中R841二十一點、新增5點。
原始on/off全部snapshot欄位、RAM、indexed、palette、暫存器SHA、固定輸入及停止理由相同。
每批完整包含前一批全部點；最新取得／印字入口291／572，不把五批次數相加。
十來源／解碼／同緩衝印字／保護段／來源失配皆0，只證明此零命中觀察路線。

面向2第一次前進+8由58到57，六word第三項2到1只保留原始值，不升級+4語意。
接著面向2及3各一次前進不動，停止兩方向；牆、區域／世界及傳送未推定。
五張正常畫格及observations保存逐批決策，不轉寫完整旁白；trail=2未使用，
這些收據不核准第二次自動前進。後續只在觀察到不同已實測方向後續行。

初334來源、v2為336、v3為338；舊索引main按保存v1／v2核SHA，不誤報現行差異。
最後341來源及282產物封存；248原版、55已索引SDK及R780私有loader保持，
現行探針與實際SDK編譯輸入相同。工具Go1.24.13／Python3.11.2，Docker映像Id沿前列並實測相同。
數值seed未知，同合法初態與固定輸入、無重擲；不是亂數規則對拍。

| 本機收據（workplace/r843/） | SHA-256 |
|---|---|
| report.json | `6ee31935be97e072394b4661231c4d2d0ac0993c2ae3a10c38f765729450eb88` |
| final-source-index.json，341項 | `5bd79b62cfaa9c39ce51a268bcee212995dab42dc27cda9fa3c1186d6e6682a6` |
| final-output-index.json，282項 | `6f5f0932f25fa26200dfa0040ada356397ee50c4843489887d7248aa7a865751` |
| final-guard.json | `df5f8b2e1f95ab71e77606c1a23d0f188fde7f2b324c2ca1d84b152686e88915` |

根核對341來源／282產物全同，221649項lstat異常0、十scratch無檔，輸出1000:1000。
R844獨立技術審查0／0／0，主代理完整讀取JSON與Markdown；只批准回填限定路線。
review-r844-report.json SHA `21199db05897e026da7a9f4118db05e0f097aec00f6040800b3517c2fad9a60d`；
review-r844-report.md SHA `e7fabda5f55399e369169d11c2ad9e55147528b390ce154b7e6be1d286a92ef1`。
私有暖、NoCursorWait、非GUI及冷驗收不合格保持，沒有名字正分支或新中文覆繪。
015／017仍READY、008仍DRAFT；冷存讀、正式GUI／音訊／錄放、平台與十一組門檻不變。
完整收據及原始畫格只留本機，凍結索引不覆寫，自有容器已清理。

三文件文字補審0／0／0，主代理完整讀取JSON／Markdown；
review-r844-supplement.json SHA `079c1cdcbe5615909ac5c6bbe661ceecb61bed27420c42c2aa91e7a7030cbb9d`。
根doc-guard核對三份批准回填、其餘338來源與282產物保持，五正式文件擁有權及新入口通過；
doc-guard.json SHA `9945a066bca2ca51022337f81725ade45a96d02a035627536ccad6d7ba2f539a`。
封存後只追加本段與WORKLOG驗收文字，不覆寫凍結收據或重跑原版。

三文件文字補審0／0／0，主代理完整讀取JSON／Markdown；
review-r844-supplement.json SHA `079c1cdcbe5615909ac5c6bbe661ceecb61bed27420c42c2aa91e7a7030cbb9d`。
根doc-guard核對三份批准回填、其餘338來源與282產物保持，五正式文件擁有權及新入口通過；
doc-guard.json SHA `9945a066bca2ca51022337f81725ade45a96d02a035627536ccad6d7ba2f539a`。
封存後只追加本段與WORKLOG驗收文字，不覆寫凍結收據或重跑原版。

<a id="r855-untried-maze-route"></a>

## R855：未走方向的正常迷宮路線

2026-10-03，基線839b97b，目標[R855](../goal/R855-untried-maze-route.md)。
私人唯讀探針完整重播R843的26點前綴，再以正常滑鼠操作接續六次。
兩個程序保留各自原版狀態，每次操作後同步停止並核對，不輸出新state或改原版記憶體。
原版及R843 SDK唯讀掛載，沿用R839十個來源的白名單及保護段護欄。

| 新操作 | 步數 | 面向 | 六個位置word | 結果 |
|---|---:|---:|---|---|
| 右轉 | 4055M | 0 | [0,10,1,7,57,87] | 新面向 |
| 前進 | 4081M | 0 | [0,10,2,7,58,87] | +8由57到58 |
| 右轉 | 4107M | 1 | [0,10,2,7,58,87] | R841此格曾試前進，不重試 |
| 右轉 | 4133M | 2 | [0,10,2,7,58,87] | 看原版畫格，未向此方向前進 |
| 右轉 | 4159M | 3 | [0,10,2,7,58,87] | 此格尚未試前進的方向 |
| 前進 | 4185M | 3 | [0,10,2,7,58,87] | 未移動，停止此方向 |

共32對停止點，其中26點重播、6點新增；RAM、indexed、palette、暫存器與完整snapshot相同。
33對PNG逐位元相同，含final複本；不把final當第33個停止點。
十來源取得、解碼、同緩衝印字、來源失配及保護段命中均0；最新取得入口292次。
兩程序正常退出，248原版檔SHA相同，兩個Scratch皆空。

逐次決策留在observations.json。第一個新面向的畫格轉送與前進在同一工具編排內，
主代理在操作後才完成判讀，因此不聲稱這一步已有事前主代理視覺判讀。
其餘新方向判讀與操作分開；此限制不改原版A/B雜湊結果，也不能補成中文上屏或GUI驗收。
面向3在位置[0,10,2,7,58,87]不移動後停止，不推定牆或區域語意。
先前位置[0,10,1,7,57,87]的面向1尚未在本續行量到；需要另行看畫格再走，
不在目前位置反覆試失敗方向，也不把較早路線當本輪相同初態。

實際編譯SDK來源53項，含唯讀私有checkpoint loader；完整來源64項，封存產物81項。
Go1.24.13、Python3.11.2，沿用wizardry7-go-ebiten映像。
數值seed未知，同合法checkpoint及固定輸入，沒有重擲，不稱亂數規則對拍。
準備首批因Docker工作目錄自動建立root-owned空目錄停止，只移除該已確認空目錄；
第二批原版盤點誤設平面目錄，改包含既有cd子目錄，兩次皆在原版spawn前停止。

| 本機收據（workplace/r855/） | SHA-256 |
|---|---|
| report.json | `b382763203ad4f0bad610b0d7ba30ed5330b7f5d24d72f7332159619e14863f7` |
| frozen-index.json，81項 | `1eff362791f446f8075837a731415662503231b976d2f1d0dd5863779062ba89` |

收尾核對64來源／81產物、248原版及1247項擁有權，異常0；自有容器0。
未觸發的主體與十一組譯文仍未知，017維持READY、008維持DRAFT。
本批是私人暖、NoCursorWait、非GUI量測，不代表中文覆繪、冷存讀檔、錄放、音訊或平台封包通過。
原始素材、正文及原始畫格只留本機，正式Go、SDK、譯文與字型未改。

R856獨立審查阻擋0／未解除應改0／建議0，主代理完整讀取兩報告；只接受有限暖路線證據回填。
每步事前判圖的目標未完全達成，限制保留，沒有新的READY、CONFORMED或可玩完成聲明。
審查JSON SHA為`32bb73ec748d1fbdb20991f892023e3434ae2823c7e0e3900057a5e18819c122`，
Markdown SHA為`bffe47137bd9d6f39b9983ce180f90ced81c823cdb69dc9bc41f2adec416b020`。
根收尾入口workplace/r855/final-check.json，核對64來源、81產物及五文件；
封存後只追加本審查紀錄，不重跑原版或覆寫凍結收據。
