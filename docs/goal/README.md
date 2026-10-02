# 專案目標

近期跨分期工作：[R701 前端聲音與可玩性](R701-audio-playability.md)；
開發入口維護：[R816 前端Docker入口](R816-frontend-docker-entry.md)。
錄放純量契約：[R819](R819-recording-initial-scalars.md)。
填槽缺字回退：[R821](R821-slot-fallback.md)。
adapter 分層盤點：[R823](R823-adapter-boundary-audit.md)。
法師建角技能情境：[R827](R827-mage-skill-context.md)。
譯文控制碼核對：[R829](R829-catalog-control-codes.md)。
剩餘控制碼與名字槽：[R831](R831-remaining-controls.md)。
段落與行號數字槽：[R833](R833-hint-numeric-slots.md)。
兩列敘述及終止碼：[R835](R835-cross-row-terminator.md)。

## 終點

**一般玩家拿自己的 DOS 版《巫術7》，能從開頭玩到結局，所有遊戲內文字都是繁體中文。**
同時提供原版沒有的輔助功能（說明、切換語言、即時存檔），並能在 Linux、Windows、macOS 上執行。

達成方式不是重寫遊戲，是讓原版執行檔跑在 dosgolem 上，由轉譯層攔截畫字、換成中文。
psychic-war 已用這條路完成一款；本專案要驗證它對 overlay 架構＋大文本量的遊戲也成立。
所以每一期的驗收都要附數字，不能只說「可以用」。

## 三個不變的判準

1. **原版是唯一的真相。** 畫面、顏色、音樂、文字覆蓋率都要拿原版或交叉 oracle（DOSBox-X）比過；
   只看自己的內部訊號（攔截器有觸發、測試綠燈）不算完成。
2. **不散布原版素材。** 發行物與版控裡只有自製的程式、譯文、字型、替換圖；原版由玩家自備，啟動時驗 SHA-256。
3. **通用的歸 golem，專屬的留這裡。** 判準是「換一款遊戲之後還成立嗎」。

## 分期

| 期 | 目標 | 文件 | issue |
|---|---|---|---|
| M0 | 探勘與基線 | [M0](M0-baseline.md) | #1–#6（完成） |
| M1 | 原版在 golem 上完整可跑 | [M1](M1-golem-runs-original.md) | #8–#12 |
| M2 | 可遊玩前端 | [M2](M2-playable-frontend.md) | #15–#18 |
| M3 | 文字攔截與文本抽取 | [M3](M3-text-capture.md) | #19、#20 |
| M4 | 中文繪製與全文翻譯 | 逐路徑 [規格](../spec/) 與輪次文件 | #22、#24、#25、#40 |
| M5 | 輔助功能與作弊 | 逐功能 [規格](../spec/) 與輪次文件 | #34–#36 |
| M6 | README、主題、打包、授權與發行 | [PLAN](../../PLAN.md) 的交付條件 | #7 僅涵蓋首頁，不能代表發行完成 |

本目錄保存完成定義；分期現況見 [CONTEXT.md](../../CONTEXT.md)，不由舊輪次狀態反推。

## 輪次

分期（M0–M6）是完成定義；輪次是每一輪工作要做到哪裡，可以跨分期。

**約定（使用者定案 2026-09-20）：每一輪開工前，先在本目錄建該輪的 goal markdown。**
輪次直接對齊某個分期時沿用該分期的文件（輪 1 ＝ M0）；跨分期的輪次比照 psychic-war
用 `Rn-<主題>.md` 命名並登記在下表。沒有 goal 文件的輪次不開工。

| 輪 | 文件 | 範圍 |
|---|---|---|
| 1 | [M0](M0-baseline.md) | 探勘與基線（完成：#1–#6 全關） |
| 2 | [M1](M1-golem-runs-original.md) | 原版在 golem 上完整可跑（完成） |
| 3 | [R3](R3-dosgolem-expansion.md) | dosgolem 擴充上游化＋M1 收尾＋M2 前端骨架（完成） |
| 4 | [M2](M2-playable-frontend.md) | 可遊玩前端：試玩、輸入錄放、音訊 |
| 5 | [M3](M3-text-capture.md) | 文字攔截與抽取（與輪 4 並行推進） |
| 667 | [R667](R667-flow-text.md) | 迷宮開場旁白流動文字 |
| 668 | [R668](R668-flow-zoom.md) | 旁白中文字形放大 |
| 669 | [R669](R669-closeout-context.md) | 收尾現況與文件遷移（#27） |
| 670 | [R670](R670-story-spec.md) | 開場故事證據與規格修正（#40） |
| 671 | [R671](R671-story-overlay.md) | 開場故事覆繪實作與獨立驗收（#40） |
| 672 | [R672](R672-delivery-prerequisites.md) | TSV 與發行字型前置（#24、#25） |
| 701 | [R701](R701-audio-playability.md) | 前端聲音與可玩性（#18） |
| 816 | [R816](R816-frontend-docker-entry.md) | 前端Docker開發入口、離線與清理（#43） |
| 819 | [R819](R819-recording-initial-scalars.md) | 錄放啟動純量精確契約（#17） |
| 821 | [R821](R821-slot-fallback.md) | 一般填槽缺字回退（#44、#28） |
| 823 | [R823](R823-adapter-boundary-audit.md) | adapter 分層與乾淨建置依賴盤點（#29） |
| 827 | [R827](R827-mage-skill-context.md) | 正常法師建角技能情境（#24、#22） |
| 829 | [R829](R829-catalog-control-codes.md) | 正式控制碼保留與流文可用性（#24） |
| 831 | [R831](R831-remaining-controls.md) | 剩餘控制碼與名字填槽（#24） |

第 2 輪起依 M0 的結論再定。

## 進度看哪裡

- **可執行工作唯一來源**：[GitHub Issue](https://github.com/wicanr2/wizardry7-cht/issues)。
- **目前真相與歷程**：[CONTEXT.md](../../CONTEXT.md)、[WORKLOG.md](../../WORKLOG.md)。
- [`docs/worklist.json`](../worklist.json) 已退役，內容只作歷史追溯。
- 本目錄只寫目標與驗收，**不記進度**——進度寫在這裡會過期。

## 量測起點

原版輸入與執行證據依 [研究索引](../re/README.md) 回查；版本與限制見 [CONTEXT.md](../../CONTEXT.md)。
本目錄的較早「起點」描述保留該輪歷史語境，不能當作現在尚未實跑的證明。
