# R978 譯文碰撞收斂與完成門檻核對

2026-10-08。整體目標維持「完成巫術7中文化」，不以字樣或局部測試替代。
沿[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)，
規格入口[008](../spec/008-english-key-tsv.md)、[011](../spec/011-message-wording-unification.md)。

## 本輪工作

1. 以118477b工作樹及遠端Issue重核剩餘門檻。存檔持久層已接上；冷LOAD及後續玩家路線仍缺合法既存收據。
   工作區找不到r733/two-forward.state，使用者確認沒有備份；正常遊玩由R979接續。
2. 只讀正式story-trans、monsters-trans、items-trans與既有規格，重算原樣鍵碰撞及文案差異。
   不讀MSG.DBS、messages.json、鄰文、詞庫或手冊內容，不復活已完成的171組。
3. 優先處理可由原文與已審同義契約判明的新增NO組，核所有有效產生來源與控制碼。
   候選及樣張只寫workplace/r978；DRAFT補充經唯讀獨立審查READY後才改正式文案。
4. 真正異義的情境保留缺口，不以最後寫入、多數或英文回退代替完成。
   譯名或玩家呈現取捨先做可撤回比較，再依已讀grilling提問。
5. 完成有意義的來源、烘製、字型／版面驗證後，逐檔提交、push及Issue進度回讀。
   目前狀態只寫CONTEXT，歷程追加WORKLOG；完整中文化仍須逐項完成現行交付門檻。

## 邊界與入口

資料、原版、SDK及字型唯讀；Docker離線、UID/GID1000:1000、限額、有界逾時及清理。
只用既有固定Python933與Go1.24.13 GUI映像，配置沿[R977](R977-approved-ui-wiring.md)。
不修改原版記憶體／資料、輸入驗證、遊戲規則、存檔、音訊或錄放schema。
保留未追蹤R433，不新增交付目錄，不修改~/.claude或其他專案。

路由命中規格閘門與跨規格回填，載入local/retro-remake-spec-gated-workflow.md及
local/re-resolution-backlinks.md；文件職責沿local/project-document-responsibilities.md。
R969證據及逐組裁決沿[研究019](../re/019-catalog-migration.md#r969-adjudication)，只補本輪新證據。
