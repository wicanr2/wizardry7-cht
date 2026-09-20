# 004 — IDA 建庫與函式普查（M0 #3）

日期：2026-09-20。issue #3。工具：`tools/ida.sh`（image `ida-pro-9.4-idapython:locked-v1`）、
`tools/ida/census.py`（沿用 psychic-war）、`tools/ida/seed_ovr.py`（本輪新寫）、
`tools/re/ovr_split.py`。

## 輸入與資料庫

- `DS.EXE` SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`
  （其餘檔案見 [001](001-unzip-baseline.md)）。
- 資料庫在 `workplace/ida/`（gitignore）：`DS.EXE.i64`（MZ 正常建庫，entry 0x14B3A
  ＝MZ IP 0x4b3a＋載入基礎，符合檔頭）＋ 13 顆 `*.code.i64`
  （`ovr_split.py` 切出的 code 塊，`idat -A -B -p8086 -b0` raw 建庫，
  段 bitness 修 16-bit，種子見下）。

## 種子方法（.code 塊；自訂容器沒有 entry stub 表）

1. 線性解碼：整段 `create_insn`（嵌入式資料可能解錯——`linear_failed` 如實記錄）。
2. 只在「被 call 的目標」與 code offset 0 上 `add_func`；函式邊界只來自 call 圖。
3. `auto_wait` 後普查。

⚠ `ida.sh script` 模式跑的是**副本**（筆記用雜湊的穩定性優先）——種子要用
`tools/ida.sh raw idat -A -S/tools/seed_ovr.py /work/<名>.i64` 就地跑才會留下來
（本輪踩過：`-S` 綁在建庫時會跑在自動分析之前，數字不穩；一律兩段式）。

## 數字（confirmed＝跑得出來；函式數是**下界**）

| 模組 | code bytes | 函式 | INT |
|---|---:|---:|---:|
| DS.EXE（resident） | 19,106（其餘是資料） | 139 | 67 |
| VBASE | 37,855 | 70 | 0 |
| VDOPT | 18,209 | 42 | 0 |
| VINIT | 5,641 | 17 | 0 |
| VMAZE | 44,128 | 109 | 1（int 3h，疑誤解碼資料） |
| VMELE | 14,181 | 42 | 0 |
| VMEXE | 42,295 | 55 | 0 |
| VMEXT | 21,067 | 52 | 0 |
| VMNPC | 44,751 | 82 | 0 |
| VPCLV | 21,666 | 45 | 0 |
| VPCMK | 28,707 | 20 | 0 |
| VPCVW | 38,966 | 63 | 0 |
| VPOPS | 20,173 | 31 | 0 |
| VTREA | 39,437 | 76 | 0 |
| **合計** | 390,321 | **843** | 68 |

**架構結論（強證據）：所有 DOS/BIOS/EMS 存取都在 resident。**
DS.EXE 的 67 個 INT：21h×51、10h×7、**67h×4（EMS）**、16h×3、33h×2。
overlay 幾乎零 INT ⇒ overlay 經 **far call 進 resident 服務層** ⇒
`.tail` 塊＝far call fixup／重定位表的假說再強化（30–1,164 bytes）。
（INT F6h＝golem EMS trampoline，誤判已修，見 [003](003-golem-boot.md)。）

## Resident 高 fan-in helper（分群起點）

| 線性 | callers | bytes | 備註 |
|---|---:|---:|---|
| 0x014452 | 74 | 41 | 最高 fan-in；服務分派熱點？ |
| 0x011171 | 18 | 27 | |
| 0x013850 | 16 | 13 | 13-byte 級的小 helper 一整群 |
| 0x014dec | 14 | 24 | |
| 0x013948 / 0x014179 / 0x0143ab | 13 | 13–44 | |

編譯器家族：**未知**。證據：printf 式格式字串（C 家族）、
runtime 字串只有 `STACK OVERFLOW, INCREASE STACK SIZE` 與 `8087/80287 is required!`
（後者＝遊戲自己的 x87 檢查）、無任何 vendor 字串、MZ 僅 1 個 relocation。
（`compiler-runtime-helper-fingerprints` 的完整分群留給需要時。）

## Overlay 名表（confirmed）

DS.EXE 檔案 offset `0x10d5a` 起連續 13 個 6-byte 字串（5 字元＋NUL）：
`VINIT VBASE VPCVW VPCLV VPCMK VMAZE VMEXT VMELE VPOPS VMEXE VDOPT VTREA VMNPC`，
緊接 `.ovr` 與 `Error %d loading overlay: %s$`（file 0x10dad，IDA 線性 0x20BAD）。
**沒有靜態 xref**（引用走暫存器基底計算），名表也不是立即數指向
（`0x0B5A` 立即數掃描 0 命中）⇒ loader 反推要用動態追蹤（golem oracle OnCall，
下一輪）。附近字串還有 `SCORE00 .CSF .PBA .PBL`（音樂檔名碎片）與
`INSERT DISK (^)`／`INTO DRIVE ^`（換片 UI）。

## 殘量與邊界（誠實記帳）

- 函式數是**下界**：種子只來自 call 圖＋offset 0；add_func 失敗
  74–144 處/檔（多半是誤解碼資料區裡的假 call 目標）。
- `linear_failed`（無法解碼的 byte 數）：VMAZE 20、VDOPT/VMEXT/VMNPC/VPCLV/VPOPS/VTREA 各 1。
- x87：遊戲要求 8087/80287；golem 開機未踩到該錯誤路徑（需求判定方式未查，
  M1 若撞到再處理）。

## 下一步

- 動態追 loader：oracle 攔 INT 21h AH=3D/3F 抓 `.OVR` 開檔與載入目標段，
  反推 `.tail` fixup 語意＋overlay 載入基準段。
- 拿 runtime 執行覆蓋（golem）回種 IDA（census.py 支援 seed JSON）。
