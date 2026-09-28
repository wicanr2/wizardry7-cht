# 010 — 模擬試玩：主選單 → 隊伍管理畫面（#16 進行中）

日期：2026-09-21。issue #16（進行中）。

## 已走到的位置

主選單（MAINSCRN，menu.state @940M）點第一列熱點 (159,47) →
**隊伍管理畫面**（PCVWSCRN 底圖，55,600 非零）：

- 底部按鈕列：CREATE｜REVIEW｜DELETE｜RENAME｜PICTURE｜EXIT
- 左側鏡框＝隊伍成員格（空）
- 場景：王座殘骸、金幣、骷髏戰士（PCVWSCRN.VGA 底圖）

## 已嘗試、未通的互動

| 嘗試 | 結果 |
|---|---|
| 點 CREATE 位置 (26,187)（含 premove/hold/settle） | 畫面不變（僅游標移動） |
| 掃描按鈕列與全螢幕（-sweep 104 點） | 每點都變（游標精靈），無法區分真熱點 |
| 鍵盤 C/R/D/E/Enter 掃描碼 | 畫面不變 |
| 不動滑鼠等待自動推進 | 停在隊伍管理畫面（1.2B 步以上） |

## 判讀（帶等級）

- 隊伍管理畫面的**按鈕命中測試不在 PLAY.DRV/VGA.DRV 的已知熱點表**
  （驅動段 32KB dump 掃不到矩形表模式）——命中邏輯應在**載入的 overlay
  模組**（PCVW 對應 VPCVW.OVR 或 VBASE）內部，或以指標傳入的另一張表。
- 下一步（輪 4 續）：
  1. 找出隊伍畫面對應的 overlay 模組與其載入段（開檔順序＋段位址）。
  2. 對該 overlay 反組譯，找主迴圈裡的滑鼠命中比較（`shr cx,1` 後的
     矩形比對形狀已在 VGA.DRV 0x1317 見過一次，overlay 內應有同形）。
  3. 或以 `-regs-at` 盯 0x740C:1417（熱點服務）在隊伍畫面的呼叫頻率，
     確認命中測試是否仍經驅動。

## 產物

- `workplace/run/pcvw1.png`（隊伍管理畫面）、`create.png`（CREATE 點擊後）
- `docs/images/wiz7-menu.png`、`wiz7-pcvw.png`（probe 確定性合成，README 用）

## 追記：鍵盤路徑判定（2026-09-21）

- BDA 環形緩衝驗屍：`-bios-keys` 的第一個鍵（K，word 0x1661）**停在
  0040:001E 沒被讀走**（head≠tail、後續鍵因 pending 閘沒推）→
  **名字欄不讀 BIOS 鍵盤佇列（int 16h）**。
- 遊戲也從未掛 int 09h（VecSets 0 次、IVT[9] 維持 0080:0024）。
- 判定（強證據）：**遊戲直接輪詢鍵盤控制器 port 60h（掃描碼）**——
  PLAY.DRV 掛的 INT 08h handler 或主迴圈裡讀。
- golem 的 `keyTick` 在「int 09h 向量＝stub」時刻意 stall IRQ1
  （`internal/machine/keyboard.go:105`，理由：stub 不送 EOI 會卡死
  後續中斷）——對 port 60h 輪詢的程式需要新投遞模式：
  **寫 kbdData（60h 可讀）但不發 IRQ1**。屬 dosgolem 擴充，需 spec＋
  EOI 分析，下一輪做。
- probe 的 `-press` 排程在 `-load-state` 前會被狀態還原覆蓋——已修
  （dosgolem master `d16f03a`：press 設定移到 load 之後）。

## 追記二：對話框熱點表與點擊實測（2026-09-21）

- 對話框（CREATE NEW CHARACTER）的熱點表在 `ds:59FE`（4 矩形＋FFFF 結尾）：
  r1=(148,174)-(252,180)、r2=(63,150)-(135,158)、r3=(63,166)-(135,174)、
  r4=(63,182)-(135,190)。
- 點擊 r2 中心 (99,154)（含 premove、短按 300K 兩種）都**不推進**；
  每次點擊畫面只有微幅變化（55,906–55,939，游標＋小狀態）。
- 名字輸入已通：`-type "KAOKI"`（DOS stdin → int 21h AH=07h）欄位顯示成功。
  **名字輸入走 AH=07h，不是 int 16h／BDA**（BDA 環形緩衝裡的鍵沒被讀）。
- 待解：對話框點擊的消費邏輯（wait fn 回傳 -index 之後、party 畫面程式
  對它的處理）；以及 r1-r4 矩形與畫面上可見元素的對應。
- probe 工具備忘：`-shots`/`-save-state` 的步數等於 `-steps` 上限時不會
  觸發（off-by-one，勿把缺檔誤讀為程式退出）；`-press` 排程需在
  `-load-state` 之後（dosgolem `d16f03a` 已修）。

## 追記三：繪圖服務架構定案（2026-09-21 深夜）

- **driver 服務分發表**：VGA.DRV file 0x003 起 = `E8 rel16 CB` 序列
  （near call 服務函式＋retf）——resident 遠呼驅動入口後，dispatcher
  near-call 對應服務。blit 函式 = file 0x4B5（module 0x5B5）。
- **sub_10954（輸入+繪圖泵，IDA 0x10954）** 在 +0x94 呼叫 sub_1408D
  （= `call cs:dword_12D8F` 間接遠呼）——指標表 cs:12D8B/12D8F/12DE3…
  每槽對應一個 driver 服務 thunk（module 0x103+4k）。
- **文字字形也經這條路**：blit 的 DS:SI 參數指向「展開後的字形緩衝」
  （246.885M 的 DS=6289 是底圖案例；文字 blit 的來源＝字形組裝緩衝）。
- 印字常式（字形組裝＋字元索引來源）＝ sub_10954 的呼叫端一帶
  （resident 0110，IDA 0x12E00-13000 附近的 sub_12EA3/sub_1308）——
  下一輪 IDA 細讀目標。
- 工具備忘：OnCall 是 IP-reach hook（名稱有誤導性）；Caller() 在
  hook 內讀 [sp]/[sp+2]——near/far 呼叫的語意不同，兩種都要印。

## 追記四：UI overlay 載入架構（2026-09-21 深夜）

- segment 0x392A 內容在 450M 與 940M 一致（同一 overlay 駐留）。
- blit 呼叫端 392A:0146：Caller() 讀到 [sp]=0146, [sp+2]=392A（遠呼返回），
  但 -regs-at 392a:146 為 0 命中——**far return 從未執行**，代表被呼叫的
  函式（blit）不會 retf 回來，或 0x392A 段的內容是編碼 UI 資源而非程式碼。
- seg392a-450.bin（線性 0x392A0-0x42800，38,240 bytes）大小 ≈ VBASE.OVR
  （38,171 bytes），**確認 0x392A 段＝VBASE.OVR**（UI/選單 overlay）。
- 段內容有 `5d ee cd ce` 重複模式＋`ff ff ff ff` 分隔線＋`00 00 00` 空白
  ——**編碼後的 UI 資源**，非直接可反組譯的程式碼。
- sub_13334（loader）的搬移邏輯（interleave + 段表寫入）需要完整逆向
  才能把編碼 UI 資源還原為可分析的程式碼。

## #20 剩餘工作（下一輪）

1. 逆向 sub_13334 完整邏輯（loader interleave/relocation 演算法）
2. 還原 VBASE.OVR 的執行期程式碼（印字常式在其中）
3. 找到字元索引→VBFONT 字形→A0000 的對應式
4. 確認 MSG.DBS 的字串如何傳給印字常式

## 追記（2026-09-27 R7 後複驗）：autoplay 序列完備、迷宮段停格待調

R7 翻譯完成（8,825 筆）後複跑 autopl2：

- **主選單→建隊→START→防拷→旁白**路徑仍通（log 18 個 AUTO 至 40.7s；
  110s 視窗只夠跑到旁白第 12 次點擊，後續時序未及執行）。
- **迷宮移動八點截圖（mv-*.png，19:49 輪）md5 全同**＝畫面停格。
  追記十四的移動實證是 probe 環境（較快節拍）；前端 Ebiten 節流下
  旁白 40 次推進可能沒走完、或中文化後防拷/旁白的回應時序改變。
- **下一輪除錯順序**：(1) autopl2 的 sleep 110 加長到旁白走完為止；
  (2) 每次點擊後 dump 幀雜湊確認轉場；(3) 迷宮段起點先驗
  -hdshot 的單幀（不靠 Xvfb import）；(4) 座標表沿用 main.go
  dirT 八點，先驗 (142,165) 單點配狀態檔 diff。

## 追記補（同日）：卡點定位——savegame 載入畫面等輸入

dirT 時序修正（30-58s→82-110s，避開旁白重疊）後重跑：遊戲在
33.5s 後無進一步 AUTO（主循環停）——log 尾端為 DIRECTORY/FILENAME/
SAVEGAME.BCF @25.9 億步＝**savegame 載入畫面等輸入**。

判讀：(160,101) 註解寫 START NEW GAME，實際觸發的是**匯入隊伍**
路徑（013 追記二十九的路徑），與新開遊戲不同線。下一輪二選一：
(1) 在 savegame 畫面補輸入（Enter/y 確認）讓匯入走完；
(2) 找真正的 START NEW GAME 熱點改走新路線（迷宮從頭）。
建議 (2)——驗收要求「建隊伍→進遊戲」不含匯入。

## 追記補二（同日）：五輪時序試錯——START 已通、卡在巫6匯入詢問的 CANCEL

逐步修正 autoplay 序列後的實測（每輪 autopl2 全跑）：

| 輪 | 修法 | 結果 |
|---|---|---|
| 1 | dirT 平移避開旁白重疊 | 卡 savegame 載入畫面（(160,101) 在主選單＝匯入） |
| 2-3 | 插入 (159,47) 角色選單→PCVW | 到 PCVW（RENAME/PICTURE/EXIT @y190）但 (160,101) 無效 |
| 4 | 插入 PCVW EXIT (277,190)→七項選單 | **START 生效**：防拷 AAAA 打中、進「巫6 匯入詢問」 |
| 5 | 匯入詢問 CANCEL (82,169) | CANCEL 後仍到 DIRECTORY/FILENAME（語意待驗） |

**已確認路徑**（confirmed）：主選單→(160,62) 建隊→(159,47) 角色選單→
PCVW→EXIT (277,190)→七項選單 (160,101) START NEW GAME→巫6 匯入詢問
（LOAD SAVEGAME @(121,71)／IMPORT (BANE) @(121,97)／CANCEL @(82,184)）。

**下一輪**：換檢查點法——在匯入詢問畫面存狀態檔，直接載入測三選項
熱點（y-15/y-13 兩規則各試），不用整條路徑重跑。CANCEL 的實際語意
（取消匯入 vs 取消新遊戲）以畫面轉場判。

## 追記（2026-09-28 R8）：IMPORT 快速路線實測——時序漂移，全程停在主選單

本輪把 autoplay 切到 IMPORT 路線（未提交：主選單 (121,97) IMPORT→防拷 AAAA→
LOAD 舊進度→旁白→迷宮八點），autopl2 全跑（58 個 AUTO 全執行，panic 0）。
**結論：一次有效點擊都沒打中——全程停在主選單。**

錨點（log 時間＝容器 UTC＝本地−8h，啟動 17:52:13）：

- AUTO 16s 的 (121,97) 落在片頭 credits（同期文字＝"IBM, Amiga, & Macintosh"，
  steps ~118M）——主選單文字（CHARACTER MENU…IMPORT (BANE)）到 17:52:58
  （啟動+45s，steps ~246M）才出現。**開機到主選單約 45s wall**，比脚本假設的
  16s 慢近 3 倍。
- 之後 40 次旁白點擊 (25,15) 全打在主選單左上角空白；迷宮八點點擊同樣打在
  主選單上。mv-*.png×8、hd-*.png 全數 md5 相同（`1eb248d2…`）＝停格。
  對照 `docs/images/wiz7-menu.png` 確認停格畫面＝主選單（金條無字是事件緩衝
   aged-out，另案）。
- `import -window root` 的 play2-shot.png 是 Xvfb 根視窗雜訊——此路徑不可靠，
  今後以 hdFrame（mv-/hd-）為準（舊建議重申）。

兩個獨立教訓（confirmed）：

1. **wall-clock 時序表漂移**（confirmed：主選單 +45s 才出現，舊脚本 16s 點擊
   打在 credits）。機制是假說：R7 把 xlate 從 ~1.8k 串加到 8.5k 串，每 blit
   查表成本上升、前台節拍變慢；舊時序（probe 環境／R7 前量得）全失效。
   下一輪 autoplay 時序必須按本輪實測錨點重排，或改為事件門控
   （見到主選單文字事件才點）。
2. **IMPORT 路線缺存檔**：`workplace/run/` 無 `SAVEGAME.BCF`，即使點中 IMPORT→
   LOAD 也無檔可載（正是舊輪 DIRECTORY/FILENAME 卡點的同一根因）。
   #16 驗收本來就要求「建隊伍→進遊戲」，IMPORT 繞過建隊——**退回 START NEW GAME
   路線（追記補二輪 4 的 confirmed 路徑），下一輪走檢查點法解巫6匯入詢問**。

worktree 未提交 main.go 的 IMPORT 改法處置：驗證失敗，下一輪改回 START 序列
（重排時序）或直接 revert；本輪不動，留給下一輪動手。

## 追記（2026-09-28 R9）：START 路線重排後大幅推進，終點＝DIRECTORY 載入畫面

autoplay 切回 START 序列（主選單 48s 起，窗口 205s），65 個 AUTO 全執行、
panic 0。**路徑大幅推進：建隊→七項選單→PCVW→主選單→防拷通過→DIRECTORY 卡點**
（mv-130＝DIRECTORY/FILENAME 對話框，C:\DSAVANT\／SAVEGAME.BCF，中文字正常疊畫）。

步數錨點（決定性，以 steps 為準；wall 僅供參考，啟動 18:00:36 UTC）：

| steps | 畫面 |
|---|---|
| 246M | 主選單（48s 建隊點擊命中，建隊對話框 OK/CANCEL） |
| 254M | THESUS／FIGHTER——建隊流程自動建出預設角色（免打名） |
| 260M | 七項選單（ADD／REVIEW／DISMISS／**START NEW GAME**／CONFIG／QUIT） |
| 266M | PCVW（CREATE…EXIT @y190；58s (159,47) 到達） |
| 273M | 主選單（64s EXIT 回到的是**主選單，不是七項選單**——拓樸修正） |
| 278M | 防拷 SECURITY CLEARANCE（70s (160,101) 在主選單＝匯入，舊教訓重演） |
| 287M | AAAA @(98,184) 被接受（76s Type 命中） |
| 289M | IMPORT 對話框＋DIRECTORY／FILENAME 等輸入——**終點卡點** |

三個結論（confirmed，以翻譯事件＋截圖為據）：

1. **重排時序有效**：R8 零命中 → R9 建隊、防拷全命中，THESUS 自動建出。
2. **拓樸修正**：PCVW EXIT→主選單；七項選單在建隊對話框後直接出現（260M）。
   R10 不必繞 PCVW：建隊→七項選單→START。
3. **wall-clock 不可靠**：同 binary 主選單到達 R8 45s、R9 約 55–65s（host 負載漂移）。
   R10 改步數門控（steps 決定性已由 M1 保證），按上表 steps＋餘裕觸發。

## 追記（2026-09-28 R10）：步數門控成功，START→防拷→引言直通，剩引言長度

autoplay 改步數門控（63 個 AUTO 全執行、panic 0）。**路徑確認**：
主選單 246M→建隊→THESUS 264M→七項選單 270M→START（274M 發火）→
**防拷直達 276M（中間無巫6匯入詢問）**→AAAA 288M 通過→引言旁白
（Guardia／Astral Dominae／Dark Savant／Vi Domina…）一路到 1225M＋。
mv-568＝引言星空頁（無字時的正常畫面）。

四個結論（confirmed）：

1. **步數門控有效**：發火精度啟動段 ±10M（chunk 粒度＋啟動競爭），後段歸零；
   餘裕吸收後全部命中。這是今後 autoplay 的固定打法。
2. **巫6詢問未出現**：新隊伍 START 後直達防拷。詢問只在特定條件出現
   （已有存檔？或在引言之後？）——R11 加長跑到迷宮自然揭曉；CANCEL 語意仍待測。
3. **引言超長**：40 次旁白點擊（295–529M）只走完前半；故事頁延續到 1225M＋
   （頁距約 25–30M）。迷宮八點（540–568M）全打在引言中段。
   R11：旁白點擊延長到 ~1360M（每 8M 一次，跳頁也無妨，目標是進迷宮），
   迷宮段後移到 1380M＋，窗口 480s。
4. CANCEL（280M 發火）落在防拷畫面上，無害（AAAA 照過）。

## 追記（2026-09-28 R11）：引言走完，迷宮第一幀到達（移動待證）

R11 長程（旁白 133 次到 1351M，迷宮八點 1380–1408M，窗口 480s）：
156 個 AUTO 全執行、panic 0。**引言走完**：1403M "A gentle wind blows quietly
through the trees…"→登陸 Guardia；**mv-1380＝迷宮第一幀**（森林 3D 視野＋
THESUS「凱薩」中文名＋右下羅盤）——confirmed 到達第一個可操作畫面（畫面層級）。

未竟（R12）：

1. **移動未證**：八點 mv 全同 md5（`52afb0…`），點擊落在引言尾頁／到達幀上，
   沒有產生可見移動。R12 迷宮點擊改 >1420M 門檻（著陸確認後），每點附幀雜湊。
2. 底文英文是截圖時序假象：xlate.json（8503 串）含 "odd tingling…" 等著陸句
   譯文；事件緩衝 aged-out 致合成幀用英文底。頭像中文名正常＝管線無礙。
3. **巫6詢問全程未出現**（除主選單上的選項文字）：fresh party 無匯入問題，
   CANCEL 語意題自然消解——詢問只在有存檔／匯入來源時出現（待有存檔後驗）。

## 追記（2026-09-28 R12）：著陸幀具決定性，移動零反應（輸入層待查）

迷宮八點後移到 1430–1458M（著陸 ~1403M 之後），156 AUTO 全執行、panic 0。
著陸重現（gentle wind 同步驟）：**前綴序列＋步數門控＝著陸點具決定性**。
mv-1430＝純迷宮幀（森林＋凱薩＋發光羅盤，無底文）。

但八點 mv 全同 md5（`f27e69…`），著陸後零文字事件到 3.2B 步、
idxNZ 凍結——**點擊被送達但遊戲不理，或點擊根本沒生效**。關鍵線索：
autoplay 的 `cl()` 丟棄 `o.Click` 回傳的 error，而 Click 在畫面不變時回
`NoResponseError`（附 Polls 計數，正好區分「座標錯」vs「遊戲沒準備好」；
另見 input.go：按住太短會被 int33h 低頻輪詢跳過）。

R13：`cl()` 印出 Click error（含 Polls）；迷宮段改輸入變體——中央視野
(160,100) plain／長按 Hold(5M)、羅盤前進長按、左右後退，同一輪內誰先讓
md5 變化誰就是答案；鍵盤不測（輸入模型：無快捷鍵，強證據）。

## 追記（2026-09-28 R13）：八變體全滅，但拿到關鍵區分證據

159 個 AUTO（156＋3 error 行）、panic 0。迷宮八變體（中央 plain／長按、
羅盤長按、前後左右）截圖兩種 md5——新變體檔全同（`6407f8…`）、舊 R12 檔名
殘留（`f27e69…`，本輪未覆寫，判讀時排除）。

error 日誌是分水嶺：全輪只有 3 個 Click error，全是旁白 (25,15)
（Polls=128063，頁面已走完故不變，預期內）；**迷宮八點擊零 error**＝
幀在點擊期間變化了（游標移動本身即變化），即**座標送達、按鍵機制走完、
遊戲邏輯不理**。結合 golem 記載的失敗形狀（oracle.go MouseSets 注記：
「程式設過之後我們注入的位置就被蓋掉，症狀是點了沒反應」），首要假說：
**遊戲用 AX=4 自設游標（park），我們的座標被蓋掉**——選單誤打誤中（park
點正在按鈕上），迷宮 park 點不在羅盤上故全滅。

R14：迷宮每動作後記 `Mouse()`＋`MouseSets`（len＋末筆座標步數）＋
`MouseRange()`，驗 park 假說；加跨 AUTO 按住（PressButton(0)／ReleaseButton(0)
隔 20M 步）與右鍵變體。`159=156+3` 對帳了結。

## 追記（2026-09-28 R14）：park 假說成立——迷宮每 40K 步釘游標 (302,164)

AUTO滑鼠行鐵證（迷宮段）：1430M 純觀測 pos=(302,164)、sets=305；
中央 (160,100) 點擊後 pos 照樣 (302,164)、sets 305→447（一次點擊內
142 次＝約每 40K 步一次）；press／release／right／左右上之後 pos 全被
吸回 (302,164)。旁白段對照：sets=1（239M 一次）、pos=(25,15) 聽話——
park 是迷宮專屬。範圍恆 [0-639,0-199]（遊戲設的，非預設）。

判讀（confirmed）：注入座標在按住讀到的那一刻已被蓋掉；16M 跨按住讀到的
全是 park 點；(302,164) 落在 320 屏右下角落（截圖角落游標吻合）＝無效點。
選單誤打誤中之謎解：選單不 park，座標有效。

R15：賭 40K 空窗——MoveMouse 後立即零 hover 按住 20M（`o.Press` hover 恆 0，
Hold 另給），期間第一次 park 前的 poll 若讀到注入座標即中；park 原地按
（不 Move）為對照組。若滅：RE VMAZE.OVR 滑鼠命中＋40K 設游標呼叫者。

## 追記（2026-09-28 R15）：race 滅——pos 從未離開 park，空窗結構性不存在

三次 racepress（前進／上／左）＋parkpress 對照組：全程 pos=(302,164)，
sets 305→942 一路漲（一次 20M 按住內 +161）。新截圖八張全同（`1c1c18…`；
`f27e69`／`6407f8`／`b13b54` 是 R12/R13 殘留檔，本輪座標不同未覆寫）。

判讀：pos 連 transiently 偏離都沒有——MoveMouse 與第一次 poll 之間必有一次
park，即遊戲迴圈順序是**先 park 後讀**（或 park 由更高頻的中斷做）。
40K 空窗不是小，是沒有。動態黑盒輸入至此用盡（R12–R15：plain／長按／跨按住／
右鍵／中央／race／對照組全滅）。

R16 轉靜態 RE（範圍精確）：VMAZE.OVR 內 (1) 滑鼠命中比較（矩形表＋座標映射：
range 640 寬如何對 320 屏）；(2) 40K 設游標的呼叫者（常駐鎖定還是劇本旗標——
若是劇本鎖定沒解開，可能是 golem 未送達的計時器／IRQ，屬上游缺口）；
(3) 按鍵狀態讀取點。輸入：workplace/ida 現成 .i64＋M0 普查＋本輪步數錨點。

## 追記前段（2026-09-28 R16）：靜態確證座標 doubling——R12–R15 打錯了一半

IDA 現成 .asm 直接給答案（confirmed，位址＋operand＋bytes）：

- SET 前 doubling：`shl cx,1` → AX=4（DS.EXE.asm:8675 附近 sub_140A7；
  VGA.DRV.asm:2303、2642、4249、4353）。
- READ 後 halving：AX=3 → `shr cx,1` 存 cs:word_15B0/15B2
  （VGA.DRV.asm:2309–2313；同形 2429、2491、2556、2648、4202）。
- range：AX=7 [0–639]、AX=8 [0–199]（VGA.DRV.asm:2296–2301）——與觀測一致。
- 命中輪詢形：sub_12C9（AX=3→halve→比 cs:word_15AC/15AE→異動才重畫，
  `shl cx,1`＋call sub_1E1A）；park 門：sub_13C9（byte_228/226 旗標＋
  FFFF 跳過＋15AC/15AE 代入，VGA.DRV.asm:2633–2662）。

即 **driver 座標＝邏輯座標×2**。R12–R15 注入 (142,165) 讀到 (71,164)——
羅盤左側空地，全滅有了第二解釋（park 覆蓋是第一個，兩者疊加）。
**勘誤**：R14/R15 的「park＝右下角落無效點」錯；park driver (302,164)＝
邏輯 (151,164)，在羅盤區內（但按住讀到的仍是 park，非注入點）。

選單為何不加倍也中：未解（選單命中路徑另案；可能是不同服務／cell 座標）。
不擋本輪：迷宮八點 x×2 驗 doubling，動了再回頭解釋選單。

## 追記後段（2026-09-28 R16）：doubling 動態無效——park 照贏，VMAZE 匯出是下步

八點加倍點擊（284,165／320,155／…／320,100）全同 md5（`2cdc49…`）。
AUTO滑鼠行：每次動作後 pos=(302,164)，sets 323→422（本輪 park 約 1/400K，
比 R14 的 1/40K 慢——速率不固定，呼叫者另有節奏）。

判讀：doubling 靜態成立，但 park 在觀測前就把座標吸回去，動態上
irrelevant。R12–R16 的完整因果鏈還缺最關鍵一環：**誰每輪呼叫 position
服務、帶什麼 args**（固定 (151,164)？15AC/15AE 陳舊值？byte_228 常駐 1？）。
另收穫：byte_228＝游標可見旗標（1＝shown，VGA.DRV 2323 設／2365 清），
park 是「可見游標的重定位服務」，不是鎖定器——但呼叫者仍在 VMAZE 側。

R17：ida.sh（locked-v1 image 在庫）匯出 VMAZE.code.i64 全文 .asm（唯讀掛載、
副本跑、輸出驗非空——照 ida.sh 契約），找 position 服務呼叫點＋命中矩形表。
選單免倍之謎併入（對照選單服務的 READ 路徑有無 shr）。

## 追記（2026-09-28 R17）：VMAZE 無直呼——全經 resident 服務，全文匯出免了

`tools/ida.sh script census.py VMAZE.code.i64`（契約全守：副本跑、雜湊
`aa134943…` 前後一致、輸出非空）：109 函式，int 指令**僅 1**（int3
@seg000:2480，bytes `cc`）——overlay 內無 int33、無 driver 直呼。
 park 呼叫鏈必為 VMAZE far-call → resident 服務（sub_13C9 門）→ AX=4。

附帶定案（confirmed）：golem bios.go int33 注記——虛擬屏恆 640 寬，
13h 模式回報值是像素兩倍（與遊戲 `shr` 對合，標準 DOS 行為，非 quirks）。
VGA.DRV 行號存檔：range 設定 2296–2301、halving 群 2309/2429/2491/2556/
2648/4202、poll 迴圈 sub_12C9（2418–）、park 門 sub_13C9（2633–2662：
byte_228 門＋FFFF 跳過＋15AC/15AE 代入）、byte_228 設 2323／清 2365。

決策：全文 .asm 匯出**不做**（無直接呼叫可找；far-call 目標要 loader 映射，
另案）。R18 逮現行：probe 跑到著陸（1.4B 步，前例 16B 步可行）＋
AX=4 park 的呼叫者含 stack 回溯；或查 resident 服務分發表對 service 號。
VMAZE.tail 看過了，只是資料檔名表（MAZEDATA.VGA…），非重定位。

## 追記（2026-09-28 R18）：pos 釘住＋按住 22M 仍不動——park 出局，modal 接棒

反打全程 pos=(284,165) 釘住（連續 10 次觀測），按住 22M（1434–1456M），
park 幾乎停擺（18M 內僅 +1 set）——座標＋按鍵雙雙送達，遊戲照樣不理。
**park 作為根因出局**（R14–R17 的主假說死）。

改寫案情的新證據：mv-1454 幀底文 "And with memory comes an exciting
recognition…"（區域描述），且 1403M 後零文字事件——著陸描述疑走非攔截
管線（VMAZE overlay 內文本渲染？），而描述在屏＝迷宮未真正交出輸入。
另：鍵盤免測（010 強證據：不讀 BIOS 佇列、port60 輪詢＋IRQ stall，
golem 缺無 IRQ 投遞模式）。

R19：十動全滑鼠——點底文（關閉？）→點頭像（手形解除？）→羅盤一圈。
全滅則 R20 轉 probe PressReads（AX=5 呼叫者 CS:IP，golem 有記）查按鍵路徑。

## 追記（2026-09-28 R19）：底文點掉了、手形沒了，羅盤仍不動——箭頭座標重測

十動零文字事件、十張全同（`2daf13…`），但畫面變了：mv-1430 起底文消失、
頭像旁手形消失（與 R18 mv-1454 同前綴對照）——**點底文關掉了區域描述**
（1430M 一下即消，1439M 截圖已乾淨）。modal 清除半成功。

但乾淨迷宮上羅盤一圈仍無移動。重看羅盤幾何（mv-1438，game px）：
盤面 x 55–250、y 105–167；N(155,118) W(115,140) E(195,140) S(155,158)——
R12–R19 用的 (142,165)/(160,155)/(113,178)/(188,178)/(160,196) 大半落在
盤面下緣之外（y 偏低 10–30px）！probe sweep3 的舊座標在著陸盤上根本沒對準。
（附帶：鍵盤仍免測，010 強證據不變。）

R20：同前綴＋先點底文清場，再打精確箭頭（driver＝邏輯×2）：
N(310,118) W(230,140) E(390,140) S(310,158)＋中央 (310,138)＋視野 (320,100)。
中任一即 #16 達標；全滅則 probe PressReads 查按鍵路徑。

## 追記（2026-09-28 R21）：probe 送達全通、但 -type 轉義害死整輪

probe 首跑（149 clicks、-click-polls 5、log-calls、1420M 存點）：
press/release 對在 handler `740C:1ED4`（mask B）全有反應；
`AX=0003×83182`（poll）、`AX=0004×1`（**無 park**——前端的 park 是前端
路徑特有？或狀態不同，待定）；`AH=07×7`、`AH=3F×191`。

但 1410M 截圖仍是防拷（框內 "A"）：`-type $'AAAA\r'` 經轉義層變成字面
7 bytes（`$AAAA\r`），AH=07 讀走 7 個垃圾字元，答案錯，整輪卡死防拷。
迷宮段 vram 差異＝游標位置差，無效數據。r21-landing.state 實為**防拷態**
（檔名誤導，後用注意）。附帶：int16 AH=01 ×134671（鍵盤狀態空 poll）。

教訓：跨工具搬輸入脚本時，特殊字元走檔案（R22：5 bytes 定稿），不用
shell 轉義賭運氣。送達機制本身（press/release＋polls＋handler）confirmed。

## 追記（2026-09-28 R20）：精確箭頭亦滅；頭像手形會動——輸入進得去

十點全同（`908b46…`）。但 R20 幀有手形、R19 幀無手形（同前綴唯一差別：
R19 在 1438M 點了頭像）——**頭像點擊切換了手形狀態**，即點擊輸入確實
進了遊戲邏輯（非吞噬）。底文：R20 首動點底文後無文（清場 maintained）。

判讀：輸入通路活著（頭像有反應），唯獨羅盤／移動無反應。剩下兩種可能：
(a) 按鍵讀取路徑特殊（AX=5 計數器？golem 缺口？）；(b) 移動要別的手勢
（雙擊？四角功能鈕先開地圖？）。R21 probe：-click-polls 5（按住到讀 N 次）、
-log-calls（AX=3/5 直方圖）、PressReads（呼叫者 CS:IP）、1420M 存檢查點、
vram 前後比對——一輪回答按鍵路徑。
