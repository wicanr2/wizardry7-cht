# R985：法術視窗遮擋守門與裝備中文修補

2026-10-08，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[032](../spec/032-profile-popup-occlusion.md)。
上一輪有進展：R984正常錄製與原版對照已完成，兩次整矩形不符撤除原因已證實。

1. 以R984合法正常錄製重播，唯讀記錄既有印字入口資料的caller／font及事件前景色號。
2. 核原版法術視窗邊界與未遮區域，不新增原版字串、MSG或答案表讀取。
3. 檢查2／3倍與一般／複合名稱、024重排的互動，先做可丟棄原型。
4. 完成032契約及獨立審查轉READY後才改正式程式；原版記憶體、indexed、palette與保存保持。
5. 驗正常開／取消／再開與兩字型兩倍率，修補只影響批准名稱顯示，其他像素與守門保持。

原版、SDK、字型及R984收據唯讀；輸出workplace/r985。Docker固定映像、network none、UID1000:1000、限額與逾時。
WHO／M與31異譯仍待，CLAYMORE選名未回覆，不將本修補當作全部中文化完成。R433保持未追蹤。
