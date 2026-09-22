# 輪 6 — M4：中文字形烘製與資料替換管線

- issue：[#21](https://github.com/wicanr2/wizardry7-cht/issues/21)
- 分期：M4（中文繪製與全文翻譯）第一步
- 前置證據：`docs/re/013-print-pipeline.md`（含兩則追記，2026-09-22）

## 目標

不攔截執行流、不 poke VRAM——**只替換資料**，讓原版印字管線畫出中文：

1. 烘製：倚天 15×16 字形 → 巫術7 寬字形格式（16px、兩平面）
2. 造一個新字型資源段（描述頭＋寬度表＋字形資料）
3. 改 `cs:0x364A` 表項指向新段
4. 譯文以字形索引編碼，poke 進原字串位址
5. 原版 `sub_13895 → sub_13877 → driver sub_293` 照跑，畫出中文

## 已知（013 定案）

| 項 | 內容 | 等級 |
|---|---|---|
| 印字管線 | 102DA→13895→13877→cs:0x364A→driver 293→blit 4B5 | confirmed |
| 服務 11 參數 | ax＝字型段號 | confirmed（regs-at） |
| 描述頭 | +03 每列 bytes、+06 flags bit1、+12 寬度表 | confirmed（dump 對照） |
| 寬字形 | loc_371 兩平面（顏色 ds:si＋mask bx）、每組 8 bits | confirmed（反組譯） |
| mask 幾何 | bx = var_A×[+0E] + [+0C]（0x7BB6：0x0F10/0x001C） | **未知→本輪實驗解** |

## 步驟

1. **字形烘製器** `tools/re/bake_font.py`：stdfont.15 → 取 N 個中文字 →
   巫術7 格式段映像（頭＋寬度表＋字形資料）
2. **實驗閉環**：menu.state 載入 → poke 新字型段到自由記憶體 → 改表項 7
   （0x83F2，選單未用的字型 ID）→ 用 `-poke` 把某選單項字串換成索引 →
   跑到畫面 → 看 PNG 是不是中文
3. mask 幾何從實驗結果反推定案
4. 追記 013 → spec 001-m4-glyph-substitution 標 READY

## 驗收

- 主選單至少一項顯示中文字形（原版管線畫出，非 VRAM poke）→ issue #21 关
- 產出：tools/re/bake_font.py、013 追記三、PNG 證據

## 記錄

（開工後填）
