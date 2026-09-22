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
