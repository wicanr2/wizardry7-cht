# R987：WHO選人提示呼叫見證

2026-10-09，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
上一輪R986已正式驗收世界人物頁M，來源57c78bf與交接8219f8b均推送。

1. 沿R984正常選人／取消收據，定位R986已證runtime740c:20fe位圖寫入的上層呼叫。
2. 只在已見WHO色號22寫入當刻，窄讀最多64-byte現行堆疊及有限候選caller code；不解碼鄰文字。
3. 保留原始runtime、operand及工具位址基準。候選返回框與caller在動態／IDA證據核定前維持假說。
4. 補033的WHO來源、背景與字格契約。獨立審到READY後才接正式覆繪；不重開034已量M。
5. 正常選人／取消、兩字型兩倍率、原版狀態與保存層驗證，完成後push及更新Issue。

原版、SDK及既存收據唯讀，研究輸出只在workplace/r987，Docker沿既有固定映像與UID1000:1000。
不改原版記憶體、state、wire、資料格式或玩法。完整堆疊、code、輸入及畫面只留本機。
不讀答案表、不提交手冊答案；合法防拷只重用當前已取得收據或正常窄讀／作答。
入口[033](../spec/033-world-identity-labels.md)、[025](../re/025-ui-localization-coverage.md#r986-world-identity-labels)。

## 已得結果

WHO屬游標圖形文字，整圖24×20、文字末區24×6，45個紅色墨點已量。
實際740c:1568／0166返回與來源圖已取得；已完成的完整1614發布原版trace均與控制相同，保存層保持。
首輪無視窗缺Draw確認只到一發布，另記環境；第五探針已有1614發布及重播完成，外層180秒收尾相遇，正常程序退出尚未採認。

原生「誰？」四派生樣張已在原24×20圖內放下，完整15／16字模與手尾黑輪廓保持，未擴背景。
原型入口workplace/r987/native-compare.html；先前黑底擴張樣張未採用。
033維持DRAFT，正式接線、正常WHO移動／取消及同狀態驗收仍待。
