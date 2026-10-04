# R960 戰鬥控制中文字置中

2026-10-05；[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)。承接[R959](../re/025-ui-localization-coverage.md#r959-battle-control-display)已證實七控制2倍越框，只修中文畫筆Y。

路由命中規格閘門與畫面覆繪隔離，沿已載入入口及[020](../spec/020-main-menu-labels.md)、[021](../spec/021-character-dialog-labels.md)。這是按既定原生字模、中文字在按鈕內及既有置中公式修正，沒有新增版面取捨。R953人物重排及R954、R944、R951待選保持。

1. [026](../spec/026-battle-control-labels.md)獨立唯讀審查到READY，才修改正式Go。
2. 新增七訊號完整守門及HD／Draw共用Y計算。事件、字型、文字、低解析compose與原版資料保持。
3. 守門測試及cmd/wiz7回歸通過後，以R959相同82點合法暖前綴、SDK與R943資源執行Unifont2倍，核原版摘要相同及完整HD差異只在七個舊／新字格聯集。再核3倍保持及倚天原生字模。首次實跑發現問題即停止。
4. 正式Draw共用helper接線作回歸檢查；真正暖GUI、按鈕操作與撤除仍需獨立收據，未量不冒稱CONFORMED。不新增防拷、十一組正文或玩家操作。
5. 工作負載用既有Docker映像離線1000:1000執行，名稱前綴wiz7-r960-，原版／SDK唯讀、輸出限workplace/r960，2CPU／2GiB／96pids、有界逾時及10MiB×3 logs。原始圖片、字型與收據留本機；獨立結果審查後回填025與唯一現況表、驗證推送、議題回讀及清理。

本輪修正已量排版問題，不擴稱全戰鬥中文化或全遊戲可玩。十一組、TSV、錄放、冷存讀、音訊、平台封包及發行門檻保持。

## 限定驗收結果

正式共用helper已接HD／Draw，守門及cmd/wiz7回歸、正式Linux建置通過。雙字型兩倍率四配置各七控制14原生字模在框內，56字模實例與1,664,000個完整RGB像素獨立核對相符；原版各82摘要保持。Unifont2倍只改七個舊／新字格內1,346像素，Unifont3倍PNG逐位元保持；倚天3倍只採舊Y獨立重建保持，沒有舊版實跑PNG。

最終獨立結果審查0／0／0，主代理於FINAL後完整回讀兩報告核SHA採認，026仍READY。168份私人檔一次封存，SHA `e379aa4136e978196b6807fe63b835b1bf402893b3127dd15d620e7cad1d8958`。結果及重跑入口見[025](../re/025-ui-localization-coverage.md#r960-battle-control-centering)。真正GUI、按鈕操作及動態撤除仍待驗，不當完整中文化或可玩包。
