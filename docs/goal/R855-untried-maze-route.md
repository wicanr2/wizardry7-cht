# R855：未走方向的正常迷宮路線

日期：2026-10-03。工作來源：[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)及[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。

## 目標

從R843已驗證的4029M停止點，正常換到未試過的方向。先看原版畫面，再決定下一次操作，尋找剩餘譯文及動態主體的合法場景。

## 範圍

- 精確重播R843完整前綴，核對26個停止點及末步原版雜湊。
- 私人探針最多接續十二次正常滑鼠操作，每次固定推進20M，另含原有Click的6M。
- 每次操作前由原版畫面判定；不重試已不移動的方向，不注入座標、記憶體或新存檔。
- 訊息觀察只沿用R839十個來源的白名單，不擴大正文讀取或防拷範圍。
- 觀察開／關使用相同初始狀態、相同輸入及停止步數。比較RAM、indexed、palette、暫存器與PNG。
- 原始資料唯讀，輸出留在workplace/r855。暖探針不當冷啟動、中文覆繪或存讀檔完成的證據。

## 完成條件

新的正常移動或可見阻塞有可重播收據；沒有觸發的譯文照實保留。更新現況與唯一歷程，唯讀獨立審查後推送並回讀Issue進度。原版檔不變，自有容器、root-owned檔案及誤建的.md目錄為零。

## 入口

- [目前狀態](../../CONTEXT.md)
- [R843證據](../re/014-overlay-evidence.md#r843-next-maze-heading)
- [譯文統一規格](../spec/011-message-wording-unification.md)
- [動態主體規格](../spec/017-dynamic-subject-source.md)

R855只建立私人唯讀量測，沒有正式實作或新的READY／CONFORMED聲明。
