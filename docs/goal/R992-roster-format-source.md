# R992：名單格式的結構來源

2026-10-09，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
開工來源7c20597；R991已量四個配對保持，不重開其驗收。

1. 從resident 0110:0fcc及既有名單印字收據定位格式建構端。以既有IDA9.4資料庫及必要窄runtime code／指標觀察追來源，不掃MSG或PCFILE內容。
2. 確認姓名、性別／職業碼、列位置與輸出緩衝的結構關係，分清confirmed、強推論及未知；不以目前六個姓名樣本推論通用record。
3. 如果資料已足以辨識通用顯示欄位，先補022 DRAFT並獨立審READY再實作；不繼續以未知位置追加特例。
4. 由正常組隊／選取／取消及必要新姓名路線補驗，原版、存檔與英文語意保持。圖像、原字模及原版六類狀態各用獨立oracle。
5. 若通用來源仍未知，記最小明確缺口並繼續其他已可完成中文化；不以研究代替正式交付，也不縮整體目標。

入口：[022](../spec/022-character-roster-text.md)、[R991證據](../re/025-ui-localization-coverage.md#r991-relocated-roster)、[004 IDA](../re/004-ida-census.md)。
ignored輸出只在workplace/r992；原版、正式DB、SDK、字型與既存輸入唯讀，Docker UID1000、none、資源及有界期限。
程式不注入原版寫入；完整資料、code、輸入及答案不提交。

本輪結果：[名單格式與生成端證據](../re/025-ui-localization-coverage.md#r992-roster-format-source)。
