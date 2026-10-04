# R961 戰鬥控制的GUI色號傳遞

2026-10-05；[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)。上一輪R960已推送正式置中及固定HD證據，本輪接實際GUI驗證。

路由命中正常試玩、規格閘門及文件職責，已載入入口及[026](../spec/026-battle-control-labels.md)。R961修正前Go源碼顯示cloneVisualEvents只複製fgRGB，沒有複製fg；026守門要求fg1，Draw使用這份副本，因此修正前固定HD不能證明GUI置中。

1. 只用R960已核七控制metadata與現行正式Go，在私人Go測試重現clone前後fg及目的Y，沒有原版執行或新正文讀取。
2. 對026增加GUI副本傳遞色號的窄補充，獨立唯讀審查到READY後才修正式clone與接線回歸。只補既有fg欄，不改原版、字模、資料格式或使用者體驗取捨。
3. 正式cmd/wiz7回歸及乾淨建置後，以已確認GUI來源鏈核Draw七控制字格。真正視窗、暖輸入及撤除如需新的執行工具，另審有界契約；不以來源測試冒稱實際GUI或可玩完成。
4. Docker沿既有wizardry7-go-ebiten映像，離線1000:1000、2CPU／2GiB／96pids、明確逾時及10MiB×3 logs；輸出限workplace/r961，repo／SDK／原版／字型輸入唯讀。原版及完整收據留本機，禁止MSG、手冊答案、新鄰文或原版寫入。
5. 獨立結果審查後更新025、唯一現況表與WORKLOG，推送及議題回讀，核擁有權與自有容器清理。十一組、TSV、冷存讀、錄放、音訊、人物版面待選及封包門檻保持。

已採用R953重排不代表同意R954蠟燭重疊或其他待選；本輪只修已定案七控制的GUI接線缺口。
