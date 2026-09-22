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
