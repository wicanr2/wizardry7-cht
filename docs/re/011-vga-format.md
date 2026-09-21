# 011 — .VGA 檔案格式：未壓縮點陣 sheet（confirmed）

日期：2026-09-21。相關：#16、[m3-notes2](m3-notes2.md)。

## 格式結論（confirmed）

64,000 bytes 的 `.VGA` 檔（MAINSCRN／PCVWSCRN／TIT1SCRN／TIT2SCRN…）＝
**320×200 的 raw 線性點陣圖**：無檔頭、無壓縮、byte＝DAC 索引、
直接對應 mode 13h 的 A0000 內容。

**證據**：
1. 驅動 blitter（VGA.DRV file 0x4B5 的 `rep movsb`）在 #246,884,928 起
   從載入段 6289:2C9F **連續**搬像素到 A0000——來源偏移線性前進，
   無解壓縮狀態機（trace 實錄 run-trace6.log）。
2. 載入段來源：z-run `#239,897,574 MAINSCRN.VGA → 6289:0000`（整檔 64,000）。
3. 執行期畫面（z400.vram @400M）與 MAINSCRN.VGA 檔案
   **逐位元組 58,055/64,000 相同**——差異處＝之後的執行期覆蓋
   （游標、動態欄位）。
4. 渲染管線：檔案 bytes ＝ DAC 索引，配執行期色盤（probe -dump-palette）
   即可重現畫面（`tools/re/vga_sheet.py`）。

## 涵義

- **選單/party 畫面的框線美術是烤進 sheet 的**；
  但**選單文字不是**——MAINSCRN.VGA sheet 渲染無文字，
  文字由 resident 的字型展開常式執行期繪製（VBFONT 字型）。
- PCVWSCRN.VGA 的烤入文字僅右緣直式 `EXIT` 等少量。
- 中文化：烤入文字＝點陣直接替換（少數）；主體＝執行期文字攔截（M3 主體）。

## 其他 .VGA（非 64,000）

| 檔案 | 大小 | 內容 |
|---|---:|---|
| VBFONT0-3.VGA | 1,024/544/7,424/768 | 執行期字型（載入 728A/7B0E/7B80/7BB6 段，+0x10 起為字形） |
| VCFONT1-3.VGA | 544/768/1,024? | 對話字型 |
| VMPNTR0-2.VGA | 482/128/482 | 滑鼠游標（三種） |
| VCPORT0.VGA | 2,074 | VCFONT 埠/附屬 |

（大小以 workplace/original 實測為準；載入位址見 run-trace5.log 開檔清單。）

## 色盤

- 執行期色盤由遊戲寫 DAC（埠 3C8/3C9）＋ `PALETTE.DEF`（6-bit RGB 三元組，
  前 16 色灰階）——`probe -dump-palette` 輸出**已是 8-bit**。
- .VGA sheet 本身不含色盤；索引值依執行期 DAC 解色。
