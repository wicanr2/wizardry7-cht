# R845：明示存檔目錄的前端接線候選

日期2026-10-02；基線c30bd5e；[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)仍OPEN。
入口[014草稿](../spec/014-original-save-storage.md)、[022研究](../re/022-input-recording-boundary.md#r802存檔落盤反例與正常保存證據)。

## 目標與邊界

- 沿現行正式前端與R807可丟棄守門，做明示save-dir的啟動接線副本；只寫workplace/r845。
- 原型強制明示目錄，沿兩個預設方案共有的SetScratch契約，不代選預設位置。
- 啟動先Prepare再LoadWith，成功載入後且Run／輸入／安裝音訊播放器之前SetScratch；失敗不回退只記帳。
- 原有R807十八路徑fixture不重跑；只驗此次新增CLI啟動邊界與副本建置。
- 自製最小檔案只用於守門與拒載案例，明列非原版；不當遊戲或存檔代用品。
- 暫不執行原版、讀state／RAM／MSG／手冊／答案，沒有新的保存或載入還原聲明。
- R807 guard使用Linux UID／Nlink，原型只限Linux；不猜Windows守門或冒稱跨平台可玩。
- 014仍DRAFT，正式Go／SDK不改；經獨立審查後才回填足夠契約，不以本輪升READY。

## 驗收

固定正式來源與SDK索引；候選只有明示flag、Prepare、SetScratch及R807守門副本。
CLI不提供save-dir、與自製Root重疊、非目錄及不可寫案例，在LoadWith之前以有限碼拒絕。
合法獨立目錄則守門成功，再由LoadWith拒絕明示的自製非MZ輸入；不啟動遊戲或GUI。
缺參數／危險路徑不先建寫原版子路徑；自製Root雜湊保持。實際編譯輸入與候選一致。
全部工作負載在既有離線、有界、非root Docker；Xvfb有trap及外層逾時，收尾自有容器0。
預設存檔位置與正式錄放格式仍待使用者回答，十一組及正常冷存讀／可玩封包門檻保持。
