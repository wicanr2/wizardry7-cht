# R819：錄放啟動純量契約（#17）

- 日期2026-10-02，基線d0e493b；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)仍OPEN，013仍DRAFT。
- 目標是補齊R815尚缺唯一結構來源的initial_scalars_sha256；不重跑格式大表或原版試玩。
- 路由命中規格閘門與文件職責，已讀knowledge-router.md、對應兩份local入口、目前CONTEXT與dosgolem契約。
- 先核對R755已觀察欄位、目前Oracle公開getter及冷初始化來源，區分設定、來源導出的預設與未提供的觀察值。
- 若需要新getter，只在workplace/r819可丟棄來源副本提出候選，先由唯讀子代理獨立審查；不修改正式執行器或前端。
- 僅以自製初始化fixture驗結構與唯讀性，不載入／執行原版、不讀RAM、state、完整輸入或防拷正文。
- 所有工作負載在既有Docker工具鏈，離線、UID/GID1000、資源與逾時有界；來源唯讀，只寫workplace/r819。
- 不決定預設存檔位置、正式schema、配置內嵌／附檔或停止錄製方式；不升READY、未通過審查不移植。
- 驗證完成後把足夠契約回填013 DRAFT／研究022及CONTEXT，歷程追加WORKLOG，提交推送並更新#17；不當中文化或交付完成。

## 目前成果與審查入口

- 自製純量群組競爭檢查通過，56份來源副本無漂移；正式程式未改。日期與RTC明列來源政策，BDA與有效滑鼠倍率不觀察。
- 候選入口：`workplace/r819/contract-candidate.md`；報告及來源／輸出索引分別為同目錄的`report.json`、`source-manifest.json`、`output-manifest.json`，67份輸出已封存。
- 獨立唯讀審查入口：`workplace/r820/`（僅保存本候選的審查收據）；由`review_r819_initial_scalars`執行，主代理已完整閱讀並核對71項來源與6項輸出。阻擋0、未解除應改0，僅批准DRAFT回填，不升READY。
- 日期來源行號訂正後，舊候選及67項索引保留、新69項索引逐項相符；Go未改、不重跑已通過測試。完整SHA與證據等級集中於[研究022](../re/022-input-recording-boundary.md)。
- 純量定義及擷取邊界回填[013](../spec/013-input-recording.md)，但B2完整啟動配置、正式格式、預設存檔位置及停止錄製方式仍待定稿。本輪沒有原版遊玩或正式封包。
