# 013 — 印字管線（print pipeline）完整解碼

狀態：confirmed（靜態反組譯 + runtime trace/dump 交叉驗證）
日期：2026-09-22
資料庫：`workplace/ida/DS.EXE.i64`（sha256 見 scratch log）、`workplace/ida/VGA.DRV.asm`

## 結論

巫術7 的所有可見文字（選單、訊息）走同一條管線，全部在 **runtime 由 resident 靜態
程式碼逐 bit 展開字形**，不是烤在 `.VGA` 底圖裡（011 號文件確認底圖乾淨）。
M3 文字攔截點因此唯一且明確：**印字串入口 `sub_13895`（runtime 0110:3895）**。

## 管線（上→下）

```
sub_102DA  縮排＋印字串                    DS.EXE IDA 0x102DA = runtime 0110:02DA
  │  [bp+4/6]=字串 far ptr、[bp+8]=x、[bp+0Ah]=字型ID
  │  loc_102F1：x 未達右界就先印空白（' '，0x20）做縮排
  ▼
sub_13895  印字串（NUL 終止，逐字元）      IDA 0x13895 = runtime 0110:3895  ★攔截點
  │  [bp+4]=字串 off（ds 段）、[bp+6]=字型ID
  │  lodsb → or al,al → jz 結束；push 字型ID、push 字元碼 → sub_13877
  ▼
sub_13877  印一字元                        IDA 0x13877 = runtime 0110:3877
  │  ax = 字元碼；bx = cs:364Ah + ax×2 → 字形指標表 cs:0x364A
  │  call cs:dword_12DB3 → driver 服務 11
  ▼
VGA.DRV sub_293  字形展開（服務 11）       file 0x293 = module 0x393
  │  ds:si = 字形 bitmap（lodsb）；es = cs:word_1E5（離屏段）
  │  dl = 前景色、dh = 背景色
  │  loc_347：shl al,1 → bit=0: es:[di]=dh / bit=1: es:[di]=dl
  │  每列 8 bits（[bx+3] 相關）、[bx+6] bit1 = 變體、[bx+0E] 參與高度乘法
  │  add [bp+var_6],140h → 離屏也是 320 寬線性
  ▼
VGA.DRV sub_4B5  blit（服務 16）           file 0x4B5 = module 0x5B5 = runtime 740C:05B5
  │  stack：x起[y起[x迄[y迄]]]（arg_8/arg_A/arg_C/arg_E = bp+0Ch/0Eh/10h/12h）
  │  ds = cs:word_1E5（離屏）、es = 0xA000
  │  目標 = y*140h+x；si 初值 = di（同座標）；每列 rep movsb
  ▼
A0000（mode 13h）
```

## 證據鏈

1. **A0000 寫入者**：`-watch-mem`（run-print3）顯示文字區寫入者 IP 全為
   `740C:0504`＝module offset → file 0x404＝`sub_293` 內部迴圈 `mov es:[di]`。
   blit `sub_4B5` 的 movsb 在 file 0x4F5-0x50x＝module 0x5F5+（另一批寫入者）。
2. **traceA（37.2M）**：`0110:59CB-5A13` 熱點＝字型描述表 `cs:0x364A` 資料區被當碼
   執行的瞬間（loc_300 分支的變體路徑），佐證 0x364A 表在用。
3. **word_1E5 runtime dump 為 0**（246.42M／247M 兩次）：blit 與字形展開是
   事件式，閒置時歸零或位址尚未驗證；不影響結論（driver 內 38 處讀取、0 處
   直接寫入，由 resident 透過服務設定）。
4. **靜態反組譯**：sub_13877 的 `mov bx,364Ah; shl ax,1; add bx,ax` 是字形指標表
   定址的直接證據；sub_13895 的 lodsb/or al,al/jz 是 NUL 終止字串的直接證據。

## 攔截設計（M4 依此寫 spec）

- 攔 `oracle.Addr{Seg:0x0110, Off:0x3895}`（sub_13895 入口）＝「一則訊息」單位。
- 參數：stack [bp+4]=字串 off、[bp+6]=字型ID；ds 段由呼叫者決定（用 `o.Arg()` 讀）。
- 替換策略：
  1. 讀原字串 → 以（ds,off）為 key 查 `text/` 譯文表；
  2. 有譯文：跳過原函式、以 CJK 點陣字（倚天 stdfont.15，15×16）直接畫離屏再
     觸發 blit；無譯文：放行原文＋記錄缺譯（rulebook：不得靜默略過）。
- 跳過原函式：OnCall callback 內改 IP／ret 模擬（dosgolem oracle API 待查
  `SetIP`/`Skip` 支援；若無則在 sub_13877 層攔、逐字元替換）。
- 字型：倚天 15×16、BIG5 索引（`~/cht/etan_font/stdfont.15`），譯文以 BIG5
  編碼存放，繪製時逐字索引字形。

## 與 psychic-war 的對應

| psychic-war | 巫術7 |
|---|---|
| blitter 入口 OnCall | 同機制，改掛 0110:3895（印字串） |
| 攔 blit 換字元資料 | 攔字串層整串替換（單位更大、更穩） |
| 點陣字烘製 | 同：font/ 子集烘製（BIG5→字形索引） |

## 追記（同日）：字型資源表與參數協定實證

### cs:0x364A＝字型資源段表（confirmed，runtime dump）

`-dump-seg 0110:364a:100`（menu 時點）：

```
項 0: 0x7B0E   項 1: 0x7B80   項 2: 0x7BB6   項 3: 0x7D88
項 4: 0x7DDA   項 5: 0x826C   項 6: 0x83A3   項 7: 0x83F2
項 8-9: 0x0000   項 10: 0xFFFF（表尾）
```

- 8 個已載入字型資源段；開機 trace（VBFONT0-3 載入 728A/7B0E/7B80/7BB6/7D88）
  與表內容吻合（0x728A 不在表內＝項 0 之前另有用途或已重整）。
- 表尾 0xFFFF 由 sub_1357B（int 21h AH=49h 逐項釋放）證實＝段號陣列。

### sub_13877 參數定案（confirmed，-regs-at）

`-regs-at 0110:3877`（#113.7M，開場畫 'W' 時）：DS=1106、[bp+06]=0002。

- `arg_2`＝**字型資源 ID**（0-7），不是字元碼。
- `bx = 0x364A + ID×2` → `ax = cs:[bx]`＝該字型資源的**段號**。
- far call 服務 11 時 `ax`＝字型段號；當前字元/位置由 driver 端狀態傳遞
  （resident 在印字前以其他服務設定；word_1E5 等變數）。

### 字形兩種格式（sub_293）

- `loc_347` 主路徑：每列**寫死 8 bits**（`cmp cx,9; jb; mov cx,8`）→ 8px 字形。
- `loc_371` 變體（`[bx+6]` bit1）：`mul cx`（[bx+0Eh]）→ **寬字形路徑**（16px
  中文可走這裡）——M4 烘製字形時的目標格式。

### 字形描述頭（0x7BB6:0000 dump，強證據）

driver 服務 11 的參數 `ax`＝**字型段號**；driver 內部讓 ds 指向字型段、從
`ds:0000` 讀描述頭。兩個資源的頭對照：

```
0x7BB6（ID2，16px）：0c 0e 00 02 02 80 06 ff 00 00 1c 00 10 0f 1c 00
0x7B0E（ID0， 8px）：06 06 00 01 01 80 02 00 00 00 06 00 10 04 06 00
                     之後 +12 起字元寬度表（00→0x0C/0x06、01→07/06…）
```

- `[+03]`＝每列 bytes（0x7BB6＝0x02→16px；0x7B0E＝0x01→8px）✓ 對應兩路徑
- `[+06]`＝flags（0x7BB6＝0xFF bit1→loc_371；0x7B0E＝0x00→loc_300）✓
- `[+0C]/[+0E]`＝0x0F10/0x001C（16px）、0x0010/0x0004（8px）——進 loc_371 的
  mask 平面算式（`bx = var_A×[+0E] + [+0C]`），精確語意待 M4 實作實驗閉環
- `[+12]` 起＝字元寬度表（每字元一 byte，px）

**勘誤**：先前假設 driver 讀 `ds:bx`（1106:364A+ID×2）為描述結構——runtime
dump 1106:364A 內容（遞增長度＋far ptr 串）不合結構語意，**推翻**。描述頭在
字型段開頭，`bx=0x364A+ID×2` 只用於 sub_13877 查 cs 段號表。

### 追記四（同日）：方案 B 打通——第一個中文選單（實作完成）

#### 架構勘誤與定案

- **沒有 word_1E5 離屏**：word_1E5（740C:02E5）全程 = 0（閒置標記，先前解讀作廢）。
- **真正的繪圖架構**：字形展開（sub_293）直接寫**離屏段 0x6289**；
  blit（sub_4B5）＝離屏 → A0000 同座標搬移。**0x6289 由 blit movsb 的 DS 暫存器
  實證**（blit hook 讀 Regs().DS）。
- 印字座標：x=740C:0295、y=740C:0297（driver self-modify 變數、印字前由遊戲設定；
  x 每字元 += 該字 advance——實測 'W'=11px、'r'=6px ✓ 比例字寬）。
- 前景色：dl（字元事件 DX 低 byte，選單 = 0x22 黑）；黃條底 = 0xC0。

#### 實作（cmd/xlate1，證據 workplace/bake/xlate1.png）

```
OnCall(0110:3895)：
  1. 先 unpoke（寫回 baseline byte 0x55）
  2. 讀字串（ds:off、NUL 終止）查譯文表（xlate.json：原文→每字 15 rows 位圖）
  3. 有譯文：從 (x,y) 以 16px 步進把字形寫進離屏 6289（筆畫=0x22、
     非筆畫=原像素讀回）→ poke 0xC3（ret）到 3895 讓原印失效
  4. 無譯文：放行＋不 poke
OnCall(740C:012F)：記前景色 dl
OnCall(740C:05F5)（blit movsb）：記離屏段 DS
```

- 執行視窗：240M 掛 hook → 940M 收（主選單印字在 **246.8M**）。
- **結果：主選單六項全數中文上屏**（角色選單／新增角色／載入紀錄／
  系統設定／匯入隊伍／結束遊戲），寫後即讀驗證位元一致。
- 字形來源：stdfont.15 經 `tools/re/bake_font.py` 烘成 xlate.json
  （Go 側無 BIG5 編碼器，字形位圖在 Python 側烘好）。

#### 已知瑕疵（下輪打磨）

- y 對位微偏（行距 13px vs 字高 15px 輕微重疊）、x 可加置中補正
- 顏色取「最後字元事件的 dl」，若譯文串先於任何原印會 fallback 0x22
- 譯文 key 用字串內容（stack 字串無固定位址）——MSG.DBS 文本仍照
  rulebook 用來源位置 key

### 攔截點修正（Stub 相容性）

- `Stub()` 只支援 **far** 常式（fireStub 讀 [sp+2] 當 CS）；sub_13895/sub_13877
  都是 near——**不可直接 Stub**。
- 可 Stub 的是 driver thunk（far call 進入，如服務 11 的 740C:012F）。
- M4 候選設計：
  - **方案 A（資料替換，原版碼零修改）**：烘中文字形段 → 改 cs:0x364A 表項指向
    新段 → 譯文以自訂字形索引 poke 進原字串位址。倚天式做法，最貼「原版邏輯不動」。
  - **方案 B（Stub 服務 11）**：攔 740C:012F，fn 內自行畫字形。彈性大但改執行流。

### 追記三（同日）：字型資源用途、格式 v2、方案 A 死路、方案 B 轉正

#### stdfont.15 公式勘誤（confirmed，渲染驗證：木/國/巫全對）

檔案**無符號區**、從 0xA440（「一」）起。正確索引：

```
idx = (hi-0xA4)*157 + (lo-0x40 若 lo<=0x7E，否則 lo-0xA1+63)
字形 = idx*30 起 30 bytes = 15 列 × 2 bytes（big-endian bits、15px 寬）
```

先前筆記的 `(hi-0xA1)*157...` 會多算 471（=3 區），**作廢**。

#### 字型 ID 用途（confirmed，OnCall 0110:3877 直方圖 520M-945M）

| ID | 段 | 呼叫數 | 字元集 | 用途 |
|---|---|---|---|---|
| 0 | 0x7B0E | 495 | ` (),-.129ABCDEFGHILMNOPQRSTUVWY` | **主選單/大寫 UI** |
| 2 | 0x7BB6 | 1,905 | 大小寫+標點 44 種 | 對話/故事文字 |

ID0 的字元集與主選單六項字串聯集**完全吻合**。

#### 資源格式 v2（強證據，560/7440 bytes 容量對帳）

- 頭 16 bytes：**+0E = 每字形 bytes**（ID0=6＝1byte/列×6列；VBFONT2=28＝2bytes/列×14列）
- 字形流＝**等寬**：流偏移 = 碼 × 每字形 bytes（先前「寬度累積」說法作廢）
- 寬度表 +10 起＝(碼,寬) 對錶（只管顯示寬、不管資料位置）
- VBFONT2：128 字形×28B×2 平面 = 3,584×2 ✓；ID0：128 條表+48 字形×6B = 560 ✓

#### 方案 A 死路（定案）

- 字形槽 = ASCII 全集（128 或 48），**碼空間零**——塞中文必殺 ASCII
- 格式 6 列/14 列，**中文 15px 塞不進**
- 結論：**資料替換路線不可行**，#21 改走方案 B

#### 方案 B（轉正）：xlate 攔印字＋自繪

```
OnCall(0110:3895)（印字串入口，near）
  → 讀參數：字串 off/[bp+6] 字型ID（x,y 待追：上層 sub_102DA [bp+8]=x）
  → 查譯文表（key = 字串來源位址）
  → 有譯文：以 Eten 15×15 自繪（WriteBytes 畫離屏或 VRAM）、poke `ret`(0xC3)
    於 0110:3895 跳過原印（1-byte 攔截，恢復時寫回原 byte）
  → 無譯文：放行原文＋記錄（rulebook：不得靜默）
```

剩餘 RE（M4 下一輪）：印字時 x,y 的來源（102DA 層或 driver 游標狀態）、
重繪觸發點（方向鍵/滑鼠皆未觸發選單重印——20M 步靜態）。


## 未完（M3 剩餘）

- [ ] oracle 是否支援 callback 內改寫控制流（SetIP/Skip）——決定攔 13895 或 13877
- [ ] `cs:0x364A` 字形指標表的 runtime 內容 dump（字形描述結構 16 bytes 版式）
- [ ] MSG.DBS 訊息（#19）與本管線的接點：誰把訊息字串搬到 ds:si

### 追記五（同日）：PCVW 推進、防拷結案、排版殘項

- **防拷結案（confirmed）**：目標單字由遊戲印在框頂（實測 WIZARDRY）；
  **任意輸入可通過**（XXXX 對照組同樣進 LOAD 畫面）——重印版防拷失效，
  玩家不會卡死，無需繞過 RE。LOAD 畫面：LOAD/CANCEL/DIRECTORY/
  C:\DSAVANT\FILENAME/SAVEGAME.DBS（術語已譯：載入/目錄/檔名）。
- **PCVW 推進**（clicks：ADD CHARACTER→點 THESUS→OK）：THESUS+FIGHTER
  頭像入隊、選單出現 REVIEW MEMBER/DISMISS MEMBER/START NEW GAME。
  「新增角色」「系統設定」「結束遊戲」等中文字形在 PCVW 黃條上屏。
- **排版殘項**：PCVW 黃條高度矮於主選單，15px 字形上下溢出與鄰行疊加
  （「系」字頭被蓋形似「至」）。待 dy 微調或以 M2 前端放大畫布時一併處理。
- 對話框關閉用 ESC（AH=07h stdin）可繞過 #16 的對話框點擊消費問題。
- 譯文表現有 14 串（xlate.json，bake_font.py 烘製；per-string dy 欄位）。

### 追記六（同日）：開場故事中文化、進入遊戲本體

- **熱點地圖定案（confirmed，多次實測）**：
  - 主選單（六項）熱點 y＝文字 y−15：(159,47) 角色選單、(160,62) 新增角色、
    (160,112) 結束遊戲——點文字位置不會中。
  - 對話框：名單 (170,153) 選人（直接入隊）、OK (90,168) **按兩次**
    （第一次加入、第二次完成關框）、ESC 回主選單。
  - PCVW 七項選單：START NEW GAME 熱點 (160,101)、系統設定 (160,114)
    （熱點＝文字−13px）。
  - START NEW GAME → 防拷（**AAAA 任意輸入可過**）→ 開場故事 → 遊戲本體。
- **開場故事 75 行逐行中文化**（workplace/bake/story.txt，1:1 行映射）：
  譯名集中 docs/names.md（宇宙鍛爐/宇宙之環/星界統領/黑暗賢者/薇·多米娜/
  阿萊賽德斯/嘉迪亞/特蘭格/恩帕尼）。xlate.json 共 89 串、324 字形。
- **遊戲本體到達**：地牢走廊＋Aletheides NPC＋頂部地圖列（shot-1700000000.png）。
  開場故事翻頁極快（數 M 步/頁），文字幀截圖待更密 shots；翻譯命中由
  hook 事件（94 次）與寫後讀驗證證實。
- 方案 A 遺物限制移除（bake_font.py 119 字形上限，方案 B 不受限）。

### 追記七（同日）：迷宮到達、故事幀時序、餘項

- 迷宮起點（shot-1565000000：無 NPC 的走廊＝玩家起點）、開場動畫序列
  （星空→黑船→Dark Savant 王座→Vi Domina 插圖，1900-2300M 實照）全部到達。
- 開場旁白文字頁：印字 hook 命中（81-94 次/輪）、離屏寫入驗證過，但插圖輪播
  為計時器驅動、文字幀壽命極短（印後數百步即清），probe 停點難捕捉——
  視覺留檔待 M2 前端（Ebiten 即時畫面自然可見）。
- 迷宮移動輸入路徑未解（Stdin w 鍵無效；移動可能走 port 60h 輪詢或
  int 16h——下一步 RE）。
- 離屏文字區驗屍法（-dump 離屏 y 區、計數 fg 像素）已入 xlate1
  （offscr-end.bin）。

### 追記八（同日）：迷宮移動輸入 RE（未完，記錄現況）

- **BIOS 鍵盤路徑已通**：oracle.SendKeys/TypeKeys（d.PushKey → BDA 緩衝）
  送鍵成功（遊戲有讀——旁白被 w 跳過的行為差異實證）。
- **開場旁白/動畫序列**：AAAA 過防拷 → 迷宮起點 → 星空動畫 → 黑船插圖 →
  Dark Savant 王座 → Vi Domina → 行星 Guardia 停格頁。頁控制機制不明：
  w×60 連發全部落在行星停格頁（3800-5200M 四張同畫面）。
- **待解**：迷宮內移動鍵（W 前進？方向鍵？）與動畫翻頁鍵。下一步：
  1. 004 的 INT 普查找 int 16h 呼叫者位址（迷宮狀態下 -regs-at 盯著看）。
  2. 或探討 keyboard.go 的 IRQ1 stall 機制是否擋住掃描碼（probe -keys-at
     走 IRQ1 硬體路徑，跟 Stdin/BDA 是三條不同的路）。
- 產物：collectM~collectS（各輪字串）、story-sheet/move-sheet/sweep2-sheet（拼圖）。

### 追記九（同日）：遊戲內訊息窗實證、敘事推進機制

- **迷宮移動的鍵盤路徑**：int 16h 全 game 只有 3 處（004 普查；一處是
  sub_13C5E 的清緩衝+服務 10）——移動讀鍵**不走 BIOS**。候选＝AH=07h
  Stdin（名字/防拷同路徑）或 port 60h。
- **Stdin 的 w 實證有效**（`-type-at`，迷宮狀態下）：每個 w 推進一段
  訊息窗敘事（降落森林→Aletheides 說明→偵測到移動人影→…）→
  **傳送法陣動畫**（3750M，森林傳送）。先前測試全落在開場動畫期＝時序錯誤。
- **遊戲內訊息窗**：底部視窗黃色文字（shot-2900000000），走印字管線＝
  可攔截翻譯。遊戲本文的第一批畫面內文字實證。
- 傳送法陣後的狀態與迷宮移動鍵驗證＝下輪續（從 4050M 接續送鍵/觀察）。

### 追記十（同日）：迷宮卡死根因定位——int08h handler 消失

- **IVT runtime 實證**：動畫期（2300M）int08h＝793D:0BB4（PLAY.DRV tick
  handler，活著）；法陣期（4200M）int08h＝**0080:0420（golem stub）**——
  handler 在 2300-4200M 之間被拔。
- **拔掉的方式**：IVT[08h] 的 OnWrite 監看（linear 0x20-0x27）在 2200M 後
  **零命中**＝寫入不走 CPU 指令；AH=25h 的 VecSets 紀錄也無 int08h＝
  不走 DOS 服務。嫌疑＝golem 內部「程式終止還原向量」邏輯
  （internal/dos 的 AH=00h 路徑直接寫 IVT，繞過 WatchWrites）。
- **影響**：沒有 timer tick＝遊戲的時間驅動輪詢（讀鍵、動畫）全停＝
  法陣殭屍畫面。這解釋「所有鍵都沒反應」。
- **修復方向（dosgolem 上遊）**：找 internal/dos 終止路徑的向量還原碼，
  查它為何在迷宮載入流程被觸發（AH=00h 被誰呼叫？Exited flag 核對）、
  或讓 PLAY.DRV handler 的卸載/重裝週期正確。
- 移動鍵本身的謎已解：**AH=07h Stdin 是唯一活路**（int16h 幾乎不呼叫、
  IRQ1 因 int09h＝stub 永遠 stall、方向鍵 Stdin＝00+scan 兩 byte）。

### 追記十一（同日）：法陣殭屍的全面體檢（未解，證據完備待對照）

「法陣定格、所有輸入無反應」的體檢結果（4200M 時點）：

| 檢查 | 結果 | 判讀 |
|---|---|---|
| IVT int08h | 0080:0420 | ＝golem 的 BIOS 預設 timer stub（正常，非 PLAY.DRV） |
| int09h | 0080:0024 | stub → IRQ1 永遠 stall（鍵進不來） |
| port 60h 讀取 | **0 次** | 不輪詢硬體鍵盤（三國演義模式排除） |
| AH=01/07/08 呼叫 | **0 次**（KeyWaits=0） | 不讀 DOS 鍵盤服務 |
| Stdin (AH=3Fh) | Consumed=0 | 法陣期不讀 Stdin |
| IP 取樣 2000 步 | 全落在 0110:4452-446E | **sub_14452（Turbo C 堆疊檢查呼叫樁，74 callers）的呼叫風暴**＝程式活著、在跑某個小函式迴圈 |
| 缺檔／未實作 | 無 | 資源與服務齊 |
| tick 0040:006C | 14427→26278（在走） | BIOS timer 正常 |

**程式活著、在跑、但不讀任何輸入**。可能：等一個非鍵盤非滑鼠的事件、
主迴圈狀態機停在不預期分支、或開場流程其實還有未完成的步驟
（我們用 w 跳過旁白可能跳掉了「正常流程會設定的旗標」）。

**下一步（二選一）**：
1. **DOSBox-X 對照**：同輸入序列在 DOSBox-X 跑——若也卡＝我們的序列錯
  （少按了某鍵）；若前進＝golem 缺機制（此時 IVT/呼叫風暴證據直接指向修法）。
2. **旁白不跳過重跑**：不送任何 w、讓旁白自動播完（先前 collectM 的 run
   曾播到 2443M），看法陣之後的正常流程長什麼樣。

### 追記十二（同日）：交互恢復——點擊是王道、方向盤 UI 現身

- **破關鍵**：**點擊頂部左格 (25,15)＝推進訊息窗旁白**（每點一段）。
  之前「法陣殭屍」＝遊戲在**等點擊**、不是等鍵盤——Stdin 的 w 只在
  文字輸入框（名字/防拷）有效，旁白推進與移動是**滑鼠主權**。
- 連點 18 次穿越旁白（記憶恢復→危險氣息）→ **底部方向盤 UI 出現**
  （中央四箭頭＝前進/左轉/右轉/後退、四角＝書/地圖/調查/拿取）。
- 環境動畫在跑（火焰/樹葉——三幀雜湊各異）＝**遊戲主迴圈完全活著**。
- 移動點擊 (142,165) 未命中熱點——座標待 sweep（下輪：底部 UI 逐點
  點擊＋幀雜湊比對＝熱點地圖）。
- 輸入路徑總結：**鍵盤 Stdin＝文字輸入專用；旁白推進與移動＝滑鼠**。

### 追記十三（同日）：迷宮移動實證——通了

底部方向盤 sweep（十字五點＋外擴四點、每點後幀雜湊比對）：
**九點全部觸發畫面變化**——森林視角逐幀改變＝玩家在移動/轉向。
「迷宮可玩」實證達成（dpad-sweep.png 九宮格）。

- 方向盤熱點不是單一格、**底部整條都是互動區**（精確方向對應待細分：
  固定起點、單點測試、拍視角變化建「點→動作」表）。
- 至此全流程可玩：開機→主選單（中）→建隊（半中）→防拷→開場（中）→
  傳送→迷宮（**可移動**）。
- 迷宮內字串收集（開地圖/物品/法術/戰鬥）＝下輪譯文表擴充的料源。

### 追記十四（同日）：方向盤 sweep 修正版——移動實證強化

- 發現「幀雜湊比對」的兩個污染源：**游標位置**（點後游標在熱點上）與
  **環境動畫**（火焰/樹葉逐幀不同）。修正：點後游標歸位 (0,0) 再截圖；
  比對改人眼判讀「視角結構」（樹叢位置/密度）而非逐像素雜湊。
- sweep3（敘事完後 14700-16000M、九點、游標歸位）：九幀全異且
  **數幀呈結構性視角差異**（樹叢位置/密度明顯變化）——移動/轉向實證強化。
- **輸入模型總結**：AH=3Fh Stdin＝旁白/防拷/名字文字；滑鼠＝旁白推進、
  移動、轉向、UI；鍵盤快捷鍵不存在（全滑鼠 UI）。
- **精確方向表**（點→前進/左轉/右轉/後退）待「點前後無間隔」的乾淨實驗。

### 追記十五（同日）：MSG.DBS 解壓演算法破解——二元樹走訪（#19 突破）

**方法**：xlate1 加 `-trace-step/-trace-len/-trace-out`（精準停點＋逐指令
CS:IP＋機器碼軌跡）。停點＝MSG.DBS seek 0x800 read 完成後
（step 246838060）、長度 2550 步（read 到印字之間）＝解壓常式全程。

**解壓函式**（runtime 0110:3E02）：

```
push bp / mov bp,sp / cld / sub sp,4 / push ds es si di
mov ax,ds / mov es,ax        ← es = ds（呼叫者資料段）
mov si,[bp+4] / mov di,[bp+6]
mov ax,cs:[2D4E] / mov es,ax ← ★es 改指 cs:[2D4E] 段（樹/狀態區）
lodsb / xor ah,ah / mov [bp-4],ax
mov word[bp-2],8             ← 位元計數＝8
lodsb                        ← 資料 byte
bit_loop（3E49）：
  shl al,1                   ← 高位出來
  jc  3E6C                   ← bit=1 → 右子樹
  左：test word es:[bx],0x8000；jnz +5；mov bx,es:[bx]；neg bx
  右：test word es:[bx+2],0x8000；jnz；mov bx,es:[bx+2]；neg bx
  shl bx,1 / shl bx,1        ← bx×4（word 索引）
  loop bit_loop / mov cx,8 / lodsb（下一 byte）
```

**演算法＝8-bit 逐位元二元樹走訪解壓**（Huffman 樣）：
- 節點＝word、bit15＝葉旗標、內部節點＝負指標鏈（neg＋×4 索引）
- 樹根＝es:[0]、左右子樹＝[bx]/[bx+2]
- 樹區段號存於 driver 附近 cs:[2D4E]（runtime 值待 dump）
- 讀取單位＝byte、每 byte 走 8 次（每 bit 一次）

**MSG.DBS 檔案結構**（confirmed，FileOps 實證）：
- MSG.HDR 11,774 bytes＝訊息索引（開機讀入）
- MSG.DBS＝1KB 壓縮塊（0x400 對齊）、按需讀取：印訊息 N → 查 HDR →
  seek 塊 → read 1024 → 樹走訪解壓 → 印字
- 實測讀取：0x2400（選單六項）、0x2800（credits）、0x800/0xc00（主選單
  重印時）、0x1000（迷宮旁白）

**待做**：dump 樹區（es:[0] 起）＋葉值路徑＋重現解壓器（Go）→ 全文抽取。

### 追記十六（同日）：解壓器初版部分成功、葉語意待精確化

- **Go 初版解壓器**（cmd/msgx）：塊 0x2400 解出大量可辨認字串
  （LIZARDMAN／ALCHEMIST／VALKYRIE／SAMURAI／NINJA／CANCEL／NEXT LVL／
  SWORD 殘形）——**樹走訪主體正確**。
- **殘餘錯位**：部分字元丟失/錯形（SWORD→tH AORD）——**葉處理語意**
  與推論不同。精讀 sub_13E02 全文：
  - 外層：`call sub_13E46 → mov [di],bl`——**輸出＝bl**（每 8 bits 一字元）
  - sub_13E46 每呼叫走 8 bits、**bx 歸零每呼叫**、位元計數跨呼叫延續
  - **葉（bit15=0）**：`mov bx,es:[bx]`＝**葉值成為新節點**（非輸出！）——
    葉值兼具「字元與下一節點」雙重語意，精確行為需更長 trace（含葉輸出）
    或逐指令模擬
- 樹區已 dump（workplace/bake/tree.bin 8KB，樹段 0x73CB、4687M 時點）。
- **判定**：主體正確、葉語意精確化＝下輪專注（更長 trace＋樹對照表）。

### 追記十七（同日）：樹區 dump＋森林結構發現

- 樹段實證：cs:[0110:2D4E]＝0x73CB（4687M 時點）、dump 8KB→tree.bin。
- **樹區＝森林**（多棵樹共存）：字母 S/E/T/W/O/R/D 各有多個葉位置
  （例：S 在節點 27 右子、840 左子、877 左子…）——**每個訊息群一棵樹**。
- 葉結構驗證：正數（bit15=0）＝ASCII 葉（0x20 空格、0x65 'e'、0x74 't'
  ——高頻字母在早期節點＝淺路徑，Huffman 特徵）✓；負數＝內部節點。
- **待解**：塊↔樹根映射在 MSG.HDR（word 對結構：例 (0x64,0x00)——
  推測（訊息ID/樹根 或 (長度,塊偏移)），結構未定）。
- 解壓器主體正確（部分字串已解出）、修正＝HDR 映射＋塊的樹根套用。
