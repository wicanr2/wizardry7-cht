# 巫術7 繁體中文化

《巫術7》（*Wizardry VII: Crusaders of the Dark Savant*，Sir-Tech 1992）DOS 版的繁體中文化。
原版 `DS.EXE` 在 [dosgolem](https://github.com/wicanr2/dosgolem) 上執行，轉譯層辨識玩家可見輸出，
在放大後的畫布覆繪中文。目的是讓不諳英文的玩家也能理解迷宮、對話與角色資訊。

玩家須自備合法的原版遊戲。儲存庫與可散布發行包不提供原版程式、資料、美術、音樂、手冊或防拷答案。

## 遊戲與中文化方式

《巫術7》是第一人稱迷宮角色扮演遊戲：建立隊伍、探索、對話、戰鬥與培養角色。
本專案驗證的版本是 Interplay《Wizardry Archives》合輯重印版；其他版本的位址與資料尚未驗證。

中文只影響顯示。英文原文、條件比較、檔名、存檔與遊戲規則保持原樣；
覆繪層只讀原版記憶體與畫面，缺譯文時顯示英文並記錄。
對話關鍵字與手冊防拷答案仍照原版要求輸入英文；題目可覆繪中文，驗證機制維持原樣。
作弊功能另有玩家明確觸發與已證實欄位的例外規則，見 [AGENTS.md](AGENTS.md)。

## 目前能做什麼

前端已能從建隊走進迷宮，中文選單、31頁開場故事、說明與迷宮開場旁白已有驗收；
F1 開啟說明，F2 開關自動地圖，預設畫布為 960×600（3×）。
唯讀覆繪與旁白在已測路線有中英兩模式的原版狀態雜湊比較，見 [覆繪證據](docs/re/014-overlay-evidence.md) 與 [流動文字規格](docs/spec/006-flow-text.md)。

尚未完成全遊戲文字路徑與結局驗收、正式 TSV 遷移、發行字型及三平台發行包。
開場故事與迷宮旁白使用訊息譯文表並依原版分頁覆繪；資料表中有譯文不代表所有玩家畫面都已中文化。
完整狀態、證據範圍與限制見 [CONTEXT.md](CONTEXT.md)，待辦以 [GitHub Issue](https://github.com/wicanr2/wizardry7-cht/issues) 為準。

## 開發與本機執行

目前提供開發入口，尚無完成驗收的正式下載包。
需要 Docker、本機 `worktrees/dosgolem` 相容版本、使用者自備原版與本機字型；
詳細依賴基線見 [CONTEXT.md](CONTEXT.md)。建置與執行經 [tools/fe.sh](tools/fe.sh) 的 Docker 工具鏈。

```bash
tools/fe.sh build ./cmd/wiz7
```

主程式入口為 `cmd/wiz7`，原版工作副本以 `-root` 指定，預設 `workplace/run`；
譯文與字形重建方式見 [tools/re/bake_font.py](tools/re/bake_font.py)、[tools/re/bake_flow.py](tools/re/bake_flow.py) 及 [tools/re/bake_help.py](tools/re/bake_help.py)。
目前字型使用倚天 `stdfont.15` 與符號字型，僅供本機；可散布字型尚待 [#25](https://github.com/wicanr2/wizardry7-cht/issues/25)。
這些開發工具不構成已驗收的跨平台安裝流程。

原版檔案基線見 [001 原版基線](docs/re/001-unzip-baseline.md)，其中：

| 檔案 | 大小 | SHA-256 |
|---|---:|---|
| `DS.EXE` | 69,206 bytes | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |

## 文件入口

- [CONTEXT.md](CONTEXT.md)：目前真相、限制與最近交接。
- [PLAN.md](PLAN.md)：穩定路線與分期完成條件。
- [WORKLOG.md](WORKLOG.md)：歷程與舊結論勘誤。
- [目標與驗收](docs/goal/README.md)、[輸出路徑規格](docs/spec/)、[研究證據索引](docs/re/README.md)。
- [譯名表](docs/names.md) 與 [翻譯指令](text/TRANS-INSTRUCTIONS.md)。

## 授權

本專案採 [RRSAL-1.0](LICENSE)：非商業使用免費，允許修改與再散布，
實況、錄影、截圖與教學明示允許；修改後散布須提供原始碼與修改說明，貢獻回授本專案。
商業使用另洽 wicanr2@gmail.com，完整條款以 LICENSE 為準。

授權不涵蓋原版素材與第三方字型。本專案與 Sir-Tech、Interplay 沒有隸屬、合作、贊助或授權關係。
