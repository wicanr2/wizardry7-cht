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

## 未完（M3 剩餘）

- [ ] oracle 是否支援 callback 內改寫控制流（SetIP/Skip）——決定攔 13895 或 13877
- [ ] `cs:0x364A` 字形指標表的 runtime 內容 dump（字形描述結構 16 bytes 版式）
- [ ] MSG.DBS 訊息（#19）與本管線的接點：誰把訊息字串搬到 ds:si
