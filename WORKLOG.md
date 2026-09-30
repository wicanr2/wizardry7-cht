# 工作歷程與勘誤

目前真相見 [CONTEXT.md](CONTEXT.md)，工作狀態見 [GitHub Issue](https://github.com/wicanr2/wizardry7-cht/issues)。
本檔從文件遷移起追加記錄；先前逐輪記錄仍保留在 [goal 文件](docs/goal/README.md) 與 [研究索引](docs/re/README.md)，不回填未核實的歷史。

## 2026-09-30：R669，文件職責遷移（#27）

- 核對程式基線 `a63e4f8`、本機回歸 `workplace/r670/check.txt`、六份 spec 與遠端 Issue 狀態。
- 建立 CONTEXT.md、WORKLOG.md、研究索引；修正 README、PLAN、goal 索引，舊 JSON 清單保留內容並標為歷史快照。
- 舊清單的未完成項 #17、#18、#22 均在 GitHub 開啟。沒有依清單重新開啟已完成項；#26 保持擱置。
- 勘誤：PLAN 的 2026-09-21 現況已失效。MSG.DBS 已抽取，前端與中文覆繪已實作；「修改按鈕高度、寫離屏或跳過原印字」不是現行方案。唯讀覆繪由 spec 001 與研究 014 的原版雜湊收據支持。
- 勘誤：README 的「防拷形態待查」「沒有可玩前端」「訊息格式待解」及中文寫離屏的說明已撤回，改指目前證據。舊畫面不再作為現行架構或驗收示例。
- #40 開場故事草稿同步送唯讀獨立審查：最近 32 筆紀錄未限制當前頁，有沿用歷史 id 的風險；S＋1–3M 取樣未證實能覆蓋故事淡入後的有效頁。兩項為阻擋，故事維持 DRAFT。獨立重算確認 31 頁、71 片段、68 筆取訊息，逐頁字詞相等；這不解除未測負向分支。另須補證據中繼資料、統一旁白筆數與確認譯文來源。
- 驗證：本輪七份文件的 88 個本機連結全部存在；新文件已掛索引。與 HEAD 的 JSON 比較僅 note 變更，items／done／dropped 內容完整保留；三項未完成工作與主機 gh 讀取的遠端狀態相符。`git diff --check` 通過。
- Docker：本輪唯讀檢查與獨立審查容器皆已刪除，專案沒有殘留容器；工作根未見 root-owned 檔案或誤建的 `*.md` 目錄。本輪只改文件，不執行遊戲回歸、不覆寫既有收據。
- 已提交 `ac51885` 並 push 到 `origin/main`；[#27](https://github.com/wicanr2/wizardry7-cht/issues/27) 已關閉且回讀確認 CLOSED。[#40 審查留言](https://github.com/wicanr2/wizardry7-cht/issues/40#issuecomment-5913264833) 已保存阻擋與下一閘門。全遊戲與發行目標仍未完成。
