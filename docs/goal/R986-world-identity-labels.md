# R986：世界人物頁性別與選人提示

2026-10-09，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
上一輪R985已正式修補五個裝備名稱的法術視窗遮擋，限定032驗收與private push完成。

1. 核R984正常輸入的WHO提示及世界人物頁M前綴，先判定來源，不按畫面字面猜攔截路徑。
2. 窄範圍觀察既有一般及單字元印字入口、caller／font／座標／前景色及原文字框。
3. 原版資料、字串比較、indexed與palette保持；不注入存檔／state，不啟用作弊。
4. 有足夠證據後補DRAFT規格並獨立審查READY，再接正式覆繪。
5. 正常選人／取消／人物頁與雙後端兩倍率驗證，同輸入原版狀態及保存層保持。

原版、SDK與R984收據唯讀；workplace/r986為輸出。Docker固定映像、network none、UID1000:1000、限額與逾時。
手冊只可重用已窄讀目前題目的合法輸入，不讀完整答案表；完整輸入與原版素材只留本機。
WHO沿[033](../spec/033-world-identity-labels.md)、性別沿[034](../spec/034-world-profile-sex.md)補證據，證據足夠各自完成，不把未量分支稱已中文化。

## 目前結果

034獨立審查READY後已接正式前端。兩字型兩倍率的正常M人物頁、法術開／取消／再開已核，
31組完整HD／Draw文字位置相同；一組關窗前未配對Draw收據另記，不列通過分母。
四組固定原API診斷各1614發布原版六類狀態全同，16張派生HD只改批准性別域，保存bytes／mode／mtime保持。
正式全cmd回歸2.359秒通過。另58組完整HD／Draw核正常Exit回迷宮撤除及人物頁重入，四保存層保持、關窗0。
操作ACK不當完成證據；頭像／Escape未退出嘗試保留，不算成功鏈。
最終獨立實作FINAL 0／0／0完整採認，034已量M正常路線限定CONFORMED，F正常路線仍READY／未量。

WHO已定位runtime740c:20fe的透明位圖寫入，尚缺上層呼叫及來源見證，033維持DRAFT。
F共享資料流及合成字格已驗，正常F仍未量。結果入口[025](../re/025-ui-localization-coverage.md#r986-world-identity-labels)。
