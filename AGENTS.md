# wizardry7-cht — 給代理的專案規則

## 授權（使用者定案 2026-09-20）

- **允許在 host 上使用 `gh`**（wicanr2 帳號）：
  - 建立 GitHub repo **`wizardry7-cht`**（private）；
  - 對該 repo push、開立與管理 issues、檢視 Actions。
- `gh` 只在 host 跑；docker 一律 `--network none` 的規則不變。
- repo 預設 private。要轉公開必須使用者明示；授權條款依 RRSAL-1.0
  （`rulebook/85`，已定案不重問）。
- git 身分一律 `wicanr2@gmail.com`；進 repo 先檢查 `git config user.email`，
  再以 `git log --format=%ae | sort -u` 檢查歷史。

## 這是什麼

《巫術7》（*Wizardry VII: Crusaders of the Dark Savant*，Sir-Tech 1992）DOS 版的繁體中文化。

路線**不是重寫遊戲**，是比照 `~/cht/psychic-war` 的做法：原版執行檔跑在
[dosgolem](https://github.com/wicanr2/dosgolem) 上，轉譯層在原版畫字的那一刻攔下來，
把中文畫在放大後的畫布上。原版邏輯一行不動；輔助功能（說明、切換語言、即時存檔）由轉譯層提供。

中文化的目的：**縮短語言隔閡與鴻溝**——讓不諳英文的玩家也能玩到這款經典，
其他以此為原則取捨。

全景計畫（路線、分層、與 psychic-war 的差異、風險）見 [`PLAN.md`](PLAN.md)。

### 原版素材（使用者自備）

根目錄的 `Wizardry_VII_-_Crusaders_of_the_Dark_Servant_1992.zip`（已 gitignore）。
清單（confirmed，unzip -l）：

- 250 檔、解壓縮共 5,414,248 bytes，全在 `wizar7/` 下（另有 `wizar7/cd/` 子目錄）。
- 主執行檔 `DS.EXE`（69,206 bytes）、安裝程式 `VINSTALL.EXE`、啟動 `DSAVANT.BAT`／`WIZARDRY.BAT`。
- **overlay 架構**：`VMAZE.OVR`、`VMELE.OVR`、`VMEXE.OVR`、`VMEXT.OVR`、`VMNPC.OVR`、
  `VPCLV.OVR`、`VPCMK.OVR`、`VPCVW.OVR`、`VPOPS.OVR`、`VTREA.OVR`、`VINIT.OVR` 等。
- 資料檔：`*.VGA`（BASEDATA／CAVEDATA／MAZEDATA／MAINSCRN／VIDSCN01–05 等）、
  `*.PIC`（MON00– 至少 MON21、VIDEO00–15）、`*.HDR`（DISK／MASTER／MISC）。
- 檔案日期都是 1996-12-24（散布版的重製時間，非原始版）；實際版本與出品資訊待 M0 從
  遊戲內標題與檔案內容確認（zip 檔名裡的 "Dark Servant" 是副標誤拼，正式副標是 Dark Savant——常識，待驗證）。

## 分期目標（沿用 psychic-war 的骨架，逐期驗收定義在 `docs/goal/`）

| 期 | 目標 |
|---|---|
| M0 | 探勘與基線：解壓與雜湊、IDA 普查（overlay 是重點）、狀態檔檢查點、DOSBox-X 參照 |
| M1 | 原版在 golem 上完整可跑：顯示模式、音效、決定性、全程無缺口 |
| M2 | 可遊玩前端：節拍、視窗、鍵盤、滑鼠、音訊、輸入錄放 |
| M3 | 文字攔截與文本抽取 |
| M4 | 中文繪製與全文翻譯 |
| M5 | 輔助功能與作弊 |
| M6 | README、主題、打包、授權與發行 |

## 預定目錄

```
cmd/wiz7/          可遊玩的主程式（前端）
apps/wiz7/         位址、狀態、攔截點（import dosgolem）
text/              文本檔：原文、譯文、來源位置（JSON，一則一筆）
font/              CJK 點陣字子集與烘製腳本
docs/spec/         規格，標 DRAFT／READY
docs/re/           反組譯與量測紀錄
docs/worklist.json 未完成項的權威
tools/             docker 包裝與腳本
worktrees/         相關 repo 的本地 clone（gitignore）
workplace/         解開的原版、快照、暫存輸出（gitignore）
```

## 動手前

1. **SDD：spec 齊了才實作。只有標 `READY` 的規格可以動手**（RE 證據 → READY spec → 實作 → 同狀態驗證）。
2. 跑原版或改 golem 之前，先讀 `worktrees/dosgolem/CLAUDE.md` 與其 `docs/spec/000-index.md`。
3. 這款是 **overlay 架構**：解讀 `.OVR` 前先讀
   `~/.claude/knowledge-base/retro/borland-tpov-overlay-re.md`（是否 Borland 產物待 M0 確認，別預設）。
4. 任務分類先查 skill `re-retro-cht-rulebook` 的路由表。

## `[HARD]` 硬規則（沿用 psychic-war，違反即停）

- **不得散布原版素材。** repo 不含 zip 內任何檔案（`DS.EXE`、`*.OVR`、`*.VGA`、`*.PIC`、`*.HDR`）。
  根目錄 `*.zip` 由使用者自備，已列入 `.gitignore`。需要原版的測試**缺檔就 skip**，不做代用品。
- **文本檔的原文欄位進版控**（比照 psychic-war 使用者定案 2026-09-19）：轉譯層靠原文算長度與對位；
  沒有原版遊戲本來就跑不起來，原文欄位照留。
- **相關 repo 一律 clone 到 `worktrees/<repo>` 工作**，不動 `~/cht/dosgolem` 等主機上的原始工作區。
  在 worktree 內開分支，推上游前先回報使用者。
- **建置、分析、轉檔、測試一律走 docker**：`docker run --rm`、目前 UID/GID、限制記憶體與 CPU、
  預設 `--network none`、原版素材唯讀掛載。只清理自己建立的 container；禁止任何
  `docker image/system/volume/builder prune` 或 `rmi`。（`gh` 是唯一例外，見「授權」。）
- **推論標籤要誠實**：confirmed／強證據／假說／未知。導覽名稱不是證據，
  保留原始位址、operand、bytes。
- **驗證以原版為準，不以自己的內部訊號為準。** 拿快照解幀跟原版同一時點比，不用截圖量像素。

## 中文化的設計約束（沿用 psychic-war）

- **畫布放大，不縮字**（`rulebook/81`）：原版底圖整數倍放大，中文用 16×16 以上點陣畫在放大後的畫布。
- **替換單位是「一則訊息」**，key 用來源位置（檔名＋偏移），不用字串內容。
- **沒有翻譯的字串要看得出來**：缺譯文顯示原文並記錄，不靜默略過。
- **完整性要有數字**（`rulebook/83`）：抽出幾則、翻了幾則、正常路徑觸發幾則、觸發了卻不在文本檔的有幾則（必須是 0）。
- 譯名統一：人名、地名、咒語／道具名稱集中一份譯名表；台灣既有譯名（當年代理或攻略本）優先，來源要記。

## 工作追蹤

- 未完成項的權威：`docs/worklist.json`（一條對應一個 GitHub issue；`verify` 為真＝仍未完成）。
  條目做完改 JSON 並關 issue，不在 markdown 打勾。
- **每一輪開工前先在 `docs/goal/` 建該輪的 goal markdown**（使用者定案 2026-09-20）；
  對齊分期的輪次沿用該分期文件，跨分期輪次用 `Rn-<主題>.md`。
- README.md 已排入工作（`docs/worklist.json` 的 `readme` 條）：遊戲介紹＋中文化目的（縮短語言隔閡與鴻溝）＋原版自備，
  寫法依 `rulebook/80`。

## 目前已知

| 事實 | 等級 |
|---|---|
| 原版是 DOS 磁片版，`DS.EXE` ＋ overlay（`.OVR`）＋資料檔（`.VGA`／`.PIC`／`.HDR`），250 檔 | confirmed（zip 清單） |
| 顯示模式、字型與印字常式、文字存放位置、音效硬體 | 未知（M0 探勘） |
| overlay 的編譯器家族（Borland？） | 未知（M0 探勘，先讀 tpov 規則再判） |
| 防拷形態 | 未知 |
