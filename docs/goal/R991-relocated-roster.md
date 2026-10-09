# R991：加入角色後移位名單的中文覆蓋

2026-10-09，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
目前來源71af320；R990已完成女性正常保存與人物頁限定驗收，不重開該路線。

1. 從R990已確認的正常加入TEMPEST路線，窄量移位後的五列名單。只讀既有六個完整原文字串、0110:0fcc、font0、x150及六個已量y，最多16 bytes目的字串。
2. 在匹配入口及返回SP時比較完整離屏before／after，核姓名／碼帶差異、色號、背景、選取及取消生命週期；不改人物、存檔、規則或輸入。
3. 窄探針先獨立唯讀審查，正常玩家收據及原版控制另驗；資料夠後補022 DRAFT並審至READY，再改production。
4. 本輪修已存在完整英文鍵的列移位，維持既定姓名、譯文及3倍布局；任意姓名、其他職業與2倍排版仍需各自證據。
5. 雙字型正常GUI及同狀態像素驗收、保存保持、回歸、獨立複核、push及Issue完整回讀均完成後記限定範圍。

入口：[R990缺口](../re/025-ui-localization-coverage.md#r990-normal-female-profile)、[022](../spec/022-character-roster-text.md)。
ignored輸出只在workplace/r991，原版、字型、SDK及既存輸入／保存唯讀；Docker none、UID1000、資源與有界期限。
完整輸入、答案、存檔及字模不提交。

本輪結果：[移位名單證據與驗收](../re/025-ui-localization-coverage.md#r991-relocated-roster)。
