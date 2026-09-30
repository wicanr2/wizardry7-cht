# 中文化路線與分期

本專案以 dosgolem 執行原版《巫術7》，在輸出層覆繪繁體中文；不重寫遊戲規則、不修改原版資料格式。
目前狀態只在 [CONTEXT.md](CONTEXT.md) 維護，執行工作以 [GitHub Issue](https://github.com/wicanr2/wizardry7-cht/issues) 為準；歷程與舊計畫勘誤見 [WORKLOG.md](WORKLOG.md)。

| 分期 | 完成條件與入口 |
|---|---|
| M0 探勘與基線 | 原版版本、雜湊與初步定位，見 [M0](docs/goal/M0-baseline.md) |
| M1 原版執行 | 開機、顯示、輸入、音效事件與決定性，見 [M1](docs/goal/M1-golem-runs-original.md) |
| M2 可遊玩前端 | 正常操作、節拍、錄放與前端音訊，見 [M2](docs/goal/M2-playable-frontend.md) |
| M3 文字攔截與抽取 | 文字來源與輸出定位可回查，見 [M3](docs/goal/M3-text-capture.md) |
| M4 中文覆繪與全文翻譯 | 譯文 catalog、譯名、字型、玩家路徑覆蓋與原版狀態不變，見 [規格](docs/spec/) |
| M5 輔助功能 | F1 說明、即時地圖、作弊等已定案功能，見 [AGENTS.md 第 13 節](AGENTS.md) 與相關 Issue |
| M6 交付 | README、RRSAL-1.0、合法素材與 Linux／Windows／macOS 可啟動包，依交付契約驗收 |

各輸出路徑依「原版證據 → DRAFT → 獨立審查 → READY → 實作 → 同狀態驗證 → CONFORMED」前進。
只有實際收據覆蓋的畫面能稱已中文化；未到達的後期畫面與未觸發分支維持未知。
建立正式交付目錄前須載入 `~/.codex/knowledge-base/local/retro-remake-dist-all-output.md`，版本依同目錄的 `retro-remake-release-versioning.md`；本檔不另創交付結構。
