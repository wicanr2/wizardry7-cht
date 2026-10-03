# R925 新增角色小按鈕排版修正

2026-10-04；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。契約READY，R926獨立計畫審查0／0／0，主代理完整回讀後採用。

依[025](../re/025-ui-localization-coverage.md#r921-新增角色對話框)的正常新增與取消證據，只修正「確定／取消」原生中文字在2×超出框線的問題。規格[021](../spec/021-character-dialog-labels.md)先獨立審查到READY，再實作、正常GUI及雙字型兩倍率同狀態驗收。

輸入沿[R921](R921-character-dialog-ui.md)及已審取樣修訂，仍實際GUI先成功才做固定輸入輔助。原版、SDK、字型及R915／R921證據唯讀；輸出僅workplace/r925，新case拒絕已存在，原版存檔只允許與原版完全相同的PCFILE複本。Docker非root、offline、資源限制、外層逾時及自有程序清理不變。

不改原版記憶體、規則、資料格式、事件或滑鼠操作；不縮字、不改譯文。六列角色清單與圖片字樣仍另補，沒有完整遊玩或交付聲明。成果回填025、CONTEXT、WORKLOG與議題，未驗證不push。
