# 002 — DS.EXE 與 overlay 初探（M0 #2）

日期：2026-09-20。issue #2。工具：`tools/re/ovr_probe.py`（docker python:3.13-alpine）。

OVR檔頭的後續定位見[R881](016-party-position.md#r881-map-loader-route)：首四位元組包含二位元組magic與code目的偏移，`+12`為入口控制值，舊checksum是假說。下文保留本輪原始觀察。

## DS.EXE 檔頭（confirmed）

| 欄位 | 值 |
|---|---|
| 大小 | 69,206 bytes；**load module ＝ 檔案大小**（MZ 表 0x88 頁 × 512 − 426），**EXE 後面沒掛 tail** |
| relocation | **只有 1 個**：file 0x20 處 `3b4b 0000`（image offset 0x4b3b、seg 0） |
| entry | CS:IP = **0:0x4b3a**（reloc 緊鄰其後＝進入點附近有個要修的 seg word） |
| SS:SP | 0x0ff6:0xf198 |
| minalloc | 58,544 bytes（＋模組 ≈127K，與 BAT 的「580K 傳統記憶體」相容） |
| 字串風格 | `Error %d loading overlay: %s$`＝printf 式格式＋`$` 結尾（INT 21h AH=9h 慣例？）→ C 家族 runtime 的味道，**編譯器未知** |

## .OVR 容器格式（13 檔全驗）

**不是 Borland TPOV**：沒有 `TPOV` magic、DS.EXE 全檔 **0 處 `CD 3F`**（Borland
overlay manager 的觸發指令）。borland-tpov-overlay-re 的建庫管線不能照抄，
但「raw code 沒有 entry point、IDA 要自己種」的教訓直接適用。

自訂格式，13/13 檔滿足同一算式（強證據）：

```text
檔案大小 = u16@4 + u16@8 + 14
```

| 檔頭欄位 | 值 | 觀察 |
|---|---|---|
| +0 | `f2 00 55 50` | magic（13/13 相同） |
| +4 | u16 | 最大區塊的大小（5,641–44,751） |
| +6 | `98 90`（0x9098） | 常數（13/13 相同，語意未知） |
| +8 | u16 | 次區塊大小（2–1,164） |
| +10 | u16 | 第三區塊大小（2–7,962） |
| +12 | u16 | 逐檔不同，無規律 → 可能是 checksum（假說） |
| +14 | code 起點 | **未壓縮 x86**（肉眼看得到 `e8` call、`74` jz、`c3` ret） |

 區塊語意（哪塊是 reloc、哪塊是資料）未知 → M0 #3 從 DS.EXE 的 overlay
 loader（有 "Error %d loading overlay: %s$" 錯誤路徑可定位）反推。

## 驅動與硬體組態（confirmed，字串）

- **顯示**：DS.EXE 參照 `vga.drv/ega.drv/cga.drv/tandy.drv/herc.drv` 五種；
  VINSTALL 選單：VGA/MCGA(256 Color)、EGA(16)、CGA(4)、HERCULES B/W、TANDY(16)。
  zip 內含 **EGA.DRV 與 VGA.DRV 兩個驅動**＋`PALETTE.DEF`；**本安裝實際用哪個模式未知**
  （VINSTALL 的組態寫在哪、DS.EXE 怎麼選，#3 查）。
- `.DRV` 檔格式：`e9 00 00` 起手＋一排 `e8 xx xx cb`（near call＋retf）＝
  far-call 分派表（confirmed，PLAY.DRV／VGA.DRV 同形）。
- **音效**：VINSTALL 選單 `SOUND BLASTER`、`ROLAND LAPC-1/MT-32`；
  BAT 錯誤碼有 Score／Patch Bank／Resident Sounds；`.SND`×62＋`.PBL/.PBA/.CSF`。
  `PLAY.DRV` 內含 "PBAPBLU" 字串。具體硬體暫存器存取＝#3 的 INT/census 工作。
- **滑鼠**：BAT err6「Unable to Load Mouse Pointers」、VINSTALL 有 MOUSE: 選項、
  EGA.DRV 有 "Mouse pointer size out of range." → 遊戲支援滑鼠（psychic-war
  沒有的新東西，M2 前端與 golem mouse 要接）。

## 文字與資料（初步）

| 檔案 | 大小 | 觀察 |
|---|---:|---|
| `MSG.DBS` | 262,144 | 開頭不是明文（`0705 effd df7f...`），**像壓縮／編碼資料**——文字量與格式是 M3 的重頭戲 |
| `SCENARIO.DBS` | 368,320 | 同上待查 |
| `NEWGAME.DBS` | 10,720 | 初始隊伍／狀態（猜） |
| `PCFILE.DBS` | 9,368 | 角色檔（猜） |
| `VBFONT0-3.VGA`、`VCFONT1-3.VGA` | 7 檔 | 遊戲字型（**印字常式的字型來源線索**，#3 追） |

TP 經驗提醒：UI 字串可能在 code 段而非資料檔（長度前綴 or `$` 結尾），
掃描時兩種都要試、且要「有指令指著它」才算數。

## 防拷

- BAT err21 "Password Verification Error" ＝**密碼式防拷存在**。
- 形態（手冊題？磁片簽章？）未知，#3 找密碼驗證常式。

## 下一步（→ #3 IDA 普查）

1. DS.EXE 建庫：entry 0:4b3a、1 個 reloc，重點先找 overlay loader
   （"Error %d loading overlay" 字串的 xref）反推 .OVR 區塊語意。
2. INT 呼叫點全表（顯示模式切換、音效硬體、檔案存取、DOS exec）。
3. runtime helper 分群 → 編譯器家族定案。
4. 字型檔（VBFONT/VCFONT）與印字常式的配對線索。

## 推論等級摘要

| 結論 | 等級 |
|---|---|
| `.OVR` 是自訂容器（非 Borland TPOV），檔頭算式 13/13 成立 | confirmed（算式）／強證據（欄位語意待 loader 反推） |
| DS.EXE 無 tail、1 reloc、entry 0:4b3a | confirmed（檔頭解析） |
| 顯示五模式、音效 SB＋MT-32、滑鼠支援 | confirmed（VINSTALL/DS 字串） |
| 存檔走外部 PKZIP | 假說（PKZIP.EXE 在目錄裡＋INSERT SAVEGAME DISK） |
| +12 是 checksum | 假說 |
| 編譯器家族 | 未知（C 家族味道，待 #3） |
| 本安裝用哪種顯示模式 | 未知 |
