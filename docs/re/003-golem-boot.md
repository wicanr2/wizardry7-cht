# 003 — 原版在 dosgolem 上開機到標題（M0 #4）

日期：2026-09-20。issue #4。dosgolem `d0ad8b9`（worktrees/dosgolem）。

## 結論

**DS.EXE 在 dosgolem 上零修改開機，一路走到全螢幕標題美術頁。**
沒有補任何 golem 機制；只用了既有 probe 旗標（`-adlib`）。

## 怎麼跑

```sh
cd worktrees/dosgolem
DOSGOLEM_ORIG=../../../wizardry7/workplace/original/wizar7 \
DOSGOLEM_EXTRA_MOUNT="/home/anr2/cht/wizardry7/workplace:/out" \
tools/go.sh run ./cmd/probe -exe /orig/DS.EXE -root /orig \
  -steps 400000000 -adlib \
  -press enter,space,esc -press-at 20000000 -press-every 300000 \
  -dump-vram /out/pfinal.vram -dump-palette /out/pfinal.pal
```

證據檔：`workplace/pfinal.vram`（320×200 色號，31,997 非零像素）、
`pfinal.vram.png`、`pfinal.pal`；過程 log `workplace/golem-run2.log`、`golem-run3.log`。

## 開機時間軸（以指令數計，confirmed，三次實跑一致）

| 步數 | 畫面 | 非零像素 |
|---|---|---:|
| ~1.1M 前 | 載入：開檔、讀檔（AH=3D/3F/42）、配置記憶體（AH=48） | 0 |
| ~8M | 開始畫橫幅（部分） | 7,374 |
| 8M–100M+ | **停住**（等計時器；無鍵則不動） | 7,374 |
| ~130M | 橫幅長大 | 14,722 |
| ~200M | 橫幅完成＋副標文字（讀出 "…rusaders of the Dark Savant"） | 16,673 |
| 按鍵後 | 內容變化 | 16,187 |
| ~380M | **全螢幕標題美術頁**（大字＋插圖） | 31,997 |
| ~388M | 靜默退場（最後呼叫 AH=3F 讀 18 bytes；無 AH=4C，推測走 PSP:0 的 INT 20h） | — |

按鍵（enter/space/esc 從 20M 起每 300K 循環送）能把動畫推快；
第三次實跑最後遊戲自行退出，原因未查（候選：ESC 走到離開路徑、
或寫檔失敗——素材是唯讀掛載）。

## 觀測到的組態（confirmed，log 證據）

| 面向 | 觀察 |
|---|---|
| 視訊 | **模式 13h**（320×200 256 色線性 A0000）；INT 10h AH=00 ×1、AH=0F ×1、AH=10 ×2；DAC 埠 3C8/3C9 |
| 音效 | 埠 **388/389（OPL2）**；要加 `-adlib` 偵測才走音效路徑（不加時流程不同，第一次實跑提早退出）；INT 21h AH=44 ioctl ×35 |
| 計時 | PIT 通道 0 分頻 2376；INT 1Ch ×481（18.2 Hz 族）；動畫推進被計時器閘住 |
| 輸入 | INT 16h AH=01 鍵盤輪詢 ×2658；INT 33h AH=00 滑鼠初始化 ×1、輪詢 0 次 |
| 檔案 | AH=3D 開檔 ×34、AH=3F 讀 ×125、AH=42 seek ×30、AH=3E 關 ×32 |
| 記憶體 | AH=48 ×23、AH=49 ×2、AH=4A ×1 |
| **自訂中斷** | **INT F6h AH=41/42/43**（各 1 次）——候選：overlay 載入服務（#3 驗證） |

## 遇到的坑（給之後的輪次）

1. probe 的 `-dump-screen-png` 只支援平面模式，模式 13h 會報
   「現在不是平面模式」exit 1。改用 `-dump-vram`＋`-dump-palette`，
   再以 `tools/re/vram_png.py`（本 repo）合成 PNG。
2. `-adlib` 有無會讓遊戲走不同路徑（無音效時提早退出）——之後跑檢查點一律加。
3. 動畫以指令數計非常慢（380M 步才到標題）；做 #5 檢查點時要靠
   按鍵腳本與 `-load-state` 縮短重跑時間。
4. 模型看不圖時，用 ASCII 降取樣＋反相 1:1 放大讀字（本輪實作，無工具化；
   進 #5 再整理成腳本）。

## 下一步

- #3 IDA：追 INT F6h 服務＝overlay loader 反推 `.OVR` 區塊語意；
  確認「本安裝用哪個顯示模式」的組態來源（字串只見驅動名）。
- #5 檢查點：標題頁（~380M 步）存狀態檔；往主選單推進要按鍵腳本。
