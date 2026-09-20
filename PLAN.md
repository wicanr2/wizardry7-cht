# PLAN — wizardry7-cht 全景計畫

《巫術7》（*Wizardry VII: Crusaders of the Dark Savant*，Sir-Tech 1992）DOS 版繁體中文化。
本文件是路線與分期的全景；逐期完成定義在 [`docs/goal/`](docs/goal/README.md)，
未完成項的權威是 [`docs/worklist.json`](docs/worklist.json)。

## 終點（MVP）

**玩家拿自己的 DOS 版《巫術7》，能從開頭玩到結局，所有遊戲內文字都是繁體中文。**
同時提供原版沒有的輔助功能（說明、切換語言、即時存檔），並能在 Linux、Windows、macOS 上執行。

中文化的目的：縮短語言隔閡與鴻溝——讓不諳英文的玩家也能玩到這款經典。

## 路線

比照 [`~/cht/psychic-war`](../psychic-war/) 的做法：原版執行檔跑在
[dosgolem](https://github.com/wicanr2/dosgolem) 上，轉譯層在原版畫字的那一刻攔下來，
把中文畫在放大後的畫布上。**原版邏輯一行不動。**

psychic-war 是這條路的第一個案例（已完整走完 M0–M6）。本專案是第二個案例，
要回答的問題是：**這條路對「overlay 架構＋大文本量」的 90 年代大作也成立嗎？**
所以每一期驗收都要附數字，不能只說「可以用」。

## 與 psychic-war 的差異（＝主要風險）

| 面向 | psychic-war | 巫術7 | 對策 |
|---|---|---|---|
| 執行檔結構 | 單一 `PW.EXE`（EXEPACK，解壓即用） | `DS.EXE` ＋ 11 個以上 `.OVR` overlay | RE 先讀 `borland-tpov-overlay-re` 規則；位址表以「執行檔＋overlay 名＋偏移」定址 |
| 編譯器 | Microsoft C 5.x／Covox（已知） | 未知 | M0 #2 普查定案，不預設 |
| 文本量 | 1,389 則 | 預期大一個量級以上（對話、咒語、道具、NPC） | M3 抽文先量數字；翻譯用批次 sub-agent 流程 |
| 顯示 | EGA 320×200 16 色 | 未知（`.VGA` 檔名暗示 VGA） | M0 #2 確認模式 |
| 音效 | PC 喇叭／AdLib | 未知 | M0 #2 普查 |
| 防拷 | 手冊式（已被原版關掉） | 未知 | M0 查，M5 處理 |

## 分層與程式碼歸屬

判準：**「換一款遊戲之後，這段程式碼還成立嗎？」** 成立就往 dosgolem 上游送，不成立就留在本 repo。

| 層 | 放哪 | 內容 |
|---|---|---|
| 機器 | dosgolem `internal/` | CPU、DOS／BIOS、顯示卡、PIT、喇叭、滑鼠。缺的就補，依據是 DOSBox-X 原始碼 |
| 觀測 | dosgolem `oracle/` | `Load`／`RunUntil`／`Save`／`Restore`／`OnCall`／`Indexed` |
| 轉譯層機制 | dosgolem（新增套件，spec 先寫） | 印字攔截掛鉤、放大畫布與疊字、熱鍵分派、即時存檔、輸入錄放 |
| 遊戲專屬 | **本 repo** | 印字常式位址、文本檔、譯名表、CJK 字型子集、作弊位址 |
| 前端 | **本 repo**（通用部分穩定後抽回 golem） | 視窗、音訊輸出、按鍵對應 |

⚠ 位址永遠是遊戲專屬的；overlay 讓定址多一個維度（載入段會動），設計位址表時就要把這點放進介面。

## 分期

| 期 | 目標 | goal 文件 | issue |
|---|---|---|---|
| M0 | 探勘與基線 | [M0](docs/goal/M0-baseline.md) | #1–#6 |
| M1 | 原版在 golem 上完整可跑 | 待 M0 後寫 | 待開 |
| M2 | 可遊玩前端 | 待 M0 後寫 | 待開 |
| M3 | 文字攔截與文本抽取 | 待 M0 後寫 | 待開 |
| M4 | 中文繪製與全文翻譯 | 待 M0 後寫 | 待開 |
| M5 | 輔助功能與作弊 | 待 M0 後寫 | 待開 |
| M6 | README、主題、打包、授權與發行 | 待 M0 後寫 | #7（readme 已先開） |

依賴：M0 → M1 → M2；M3 可在 M1 後與 M2 並行；M4 需要 M2 前端與 M3 文本；M5 大多需要 M4 疊字層；M6 最後。
M1–M6 的細項 issue 等 M0 事實出來再開，避免開一堆寫著「未知」的空 issue（使用者定案 2026-09-20）。

## 授權與發行

- repo 私有（GitHub `wicanr2/wizardry7-cht`）；轉公開須使用者明示。
- 授權條款 RRSAL-1.0（`rulebook/85`，已定案不重問）。
- 原版素材不散布：repo 與發行包只含自製程式、譯文、字型；原版由玩家自備，工具驗 SHA-256。

## 工作追蹤

- `docs/worklist.json`：未完成項的權威，一條對應一個 issue；做完改 JSON 並關 issue。
- `docs/goal/`：只寫目標與驗收，不記進度。
- 討論與指派在 GitHub issue／milestone。
