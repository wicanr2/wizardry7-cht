# 034：世界人物頁性別覆繪

狀態READY，兩字型兩倍率已量M及TEMPEST女性正常路線限定CONFORMED；2026-10-09；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
入口[R986](../goal/R986-world-identity-labels.md)、[033](033-world-identity-labels.md)；WHO輸出另循033，不由本規格放行。

## 目的與來源

世界人物頁仍顯示M。沿既定男／女字模覆繪性別，使正常人物資料更易讀。
不改原版資料、原文、玩法、輸入或保存；不放行其他caller或一般M鍵。

Archives DS.EXE SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
R984正常錄製SHA-256 `d7204e01e17698e0a03a764e4521f7c391b03c2e2fe326683c690edcb723e453`。
R986三個窄探針各1614發布的原版RAM／Regs／Step／Ticks／indexed／palette全bytes等於raw，
trace SHA-256 `84388c8f8cf1461fc7d6e5bde6fd2d05597ba02855c3ce505ac76ea107c404ae`。
工具Go1.24.13、IDA9.4／py312既有映像；原版、字型與輸入唯讀，原版寫入0。

confirmed於正常M路線：單字元入口runtime0110:3877，返回0110:6159，字碼77、font0、x90／y6、SP65310。
返回SP65312、x96／y6，完整離屏差異13像素，皆在90／6的6×6格、全部改成色號1。
88／4的8×8背景區至晚發布仍保持；擴充區非文字色1，無其他已觀察印字入口寫入。
code來源為該正常執行當刻0110:60c9起320 bytes，SHA-256
`3d8b652dbbcbd17d553a84878a132e1a14a44f6cf0eeff4ec30dc2c7b1148259`。
原operand保留`[bx+5e2ch]`；讀取值加140後由0110:0582取字串，取目的緩衝首byte，
於0110:6156的near call送0110:3877，返回6159。140／141原文為MALE／FEMALE，首字碼M／F共享此輸出位置。
此code的圖塊資料庫以raw snippet EA記錄，與runtime段位址分開；不把舊DS.EXE靜態資料段當載入overlay code。
F的共同資料流為既有靜態證據；正常F印字及GUI已由下列R990實際收據補量，不外推其他女性人物路線。

## 覆繪契約

1. 保持024角色選單6c42路線。世界只接受6159、font0、90／6及M／F，2或3倍。
   先撤除同類性別舊事件，拒巢狀、一般／複合pending交錯、無資源或壞離屏；不改其他事件或原版。
   世界callback在新增StackWord讀取前核SP≤65530，helper也核相同入口上限；sp＋2不wrap。
   倍率獨立硬限2／3，不能靠舊sexEnabled放行。兩鍵、單漢字、dy0、同高15／16非空字模及位元界完整合格。
2. 沿既有性別pending保存入口SP、caller、完整before與字模；NUL插入仍清pending。
   6159返回只消費自己的pending；6c42也只消費自己的pending，不能用另一個返回提前完成。
3. 返回須SP＋2、96／6、完整64000-byte after。非空差異只許原6×6格且改為色1。
   顯示事件before／after保存88／4／8／8完整守門域，還原仍只改M／F的原筆畫。
4. 世界背景八點須在入口符合：88,4＝115；92,4＝116；95,4＝115；88,7＝117；
   89,9＝115；88,11＝115；91,5＝13；95,5＝14。不符記缺譯，保留原文。
   完整8×8逐發布比對，游標容忍沿既有001；其他區域改寫仍撤除，不保留過期文字。
5. 字模只取已驗`text/profile-sex.zh-TW.tsv`及其兩鍵原生資源的男／女，維持15／16完整字模。
   該資源source欄仍描述024來源，不當世界caller證據；世界來源只據本規格的獨立6159證據。
   不改資源schema、檔名、英文鍵、錄放格式或一般catalog，無須新增字型二進位。
6. 世界中文字在HD／Draw共用位移。令d＝max(字高−6×倍率,0)，畫在90×倍率−d、6×倍率−d。
   2倍倚天為177／9，Unifont176／8；3倍皆270／18。完整字格落在放大後88／4／8／8內，
   不碰原連字號96／6或下一行職業90／12。不縮字、不裁字、不改其他文字位置。
7. clone按值帶sex／原caller及獨立字模；024的2倍整頁重排與3倍性別位置保持。
   不啟用資源、倍率錯誤、原版重畫或守門失配時回原文並記錄。

## 驗收

- 原版實際before／after、八點背景及兩鍵兩字型兩倍率完整字格，由獨立oracle核。
- 錯caller、SP、font、字碼、背景、差異色／範圍、巢狀、返回交錯、缺資源與倍率負測。
- 兩後端兩倍率正常載入、人物頁、法術開／取消／再開、離開／重入；限定像素差只在性別批准域。
- 原版六類狀態及保存bytes／mode／mtime保持；正常HD／Draw核相同文字位置，不用派生圖代替GUI。
- 024原回歸、032五名稱及一般字串覆繪保持。F正常路線未量時據實維持READY。

R986最終獨立契約審查0／0／0於FINAL完整回讀採認READY；正式驗證見下一節，F正常路線仍未量。
JSON SHA-256 `d4a2b7ccc296c63456948ee297060f939c68e7712f41feeaa826bfa51127ce04`，
Markdown `073794822a9a67fade60871da5bba3887bd2689cb13724d377dc60abe92ed4dc`。

R986可丟棄原型：兩後端兩倍率共16個固定輸入派生HD樣本，獨立由舊原幀、13個原M墨點及原字模建立期望，
全像素吻合，批准性別域外差異0；每組1614原版trace全同。這些不是正常GUI。
七個world負測與024／032全cmd回歸1.511秒通過，已加讀前SP／獨立倍率／資源硬守門及世界舊狀態拒絕清理。
M／F兩鍵字格與clone及合成F邊界已驗，不以合成F輸出冒稱F正常路線。

## 正式實作驗證

正式`cmd/wiz7`沿已審契約接入6159入口及返回，HD／Draw共用位移；024 caller保持分離。
七個world負測與全部人物頁／法術視窗回歸2.359秒通過。
已測Linux二進位SHA-256 `57f5e3d2328098969c5e64526ec3aff53fb5defea535fd835e2867bf4e379076`。

四組正常冷LOAD與人物頁、法術開／取消／再開完成，31組完整HD／Draw pair在批准域逐RGB符合獨立原字模期望。
倚天3倍355M收據只有HD，關窗前未發布Draw；明記未配對，不列像素核對分母。
四保存層bytes／mode／mtime保持，五個032裝備名稱回歸通過。

四組正式來源的私人固定原API診斷各1614發布trace全bytes等於原版控制。
16張派生HD共6,656,000像素，由舊R985幀、13個原M墨點及外部原生字模獨立重建，全部相符，性別批准域外差異0。
這些固定收據不改錄放檔頭、不稱正式相容重播，派生HD不取代上述正常GUI。
安全結果留`workplace/r986/formal-pixel-result.json`與`formal-gui-result.json`；完整輸入、字型及原版畫面只留本機。
正常離開／重入另由兩字型兩倍率冷LOAD補核，共58組完整HD／Draw pair。
點盔甲右側直排EXIT返回迷宮，原性別域逐像素等於原indexed放大，沒有殘留字；
由迷宮重開人物頁時，完整原生男字重新顯示。四組保存bytes／mode／mtime保持、正常關窗0。
安全摘要`workplace/r986/formal-exit-result.json` SHA-256
`6152aa0d764d06bbf6df35b54168f71aa50ddfc673124712de89427ff2c3e8b7`。
先前頭像點擊及Escape未退出的嘗試不列完成，只依實際迷宮／人物頁畫面判定。
最終獨立實作審查0／0／0於FINAL後完整回讀採認，已量M正常路線限定CONFORMED；F正常路線保持READY／未量。
JSON SHA-256 `4b73e468805d2e4c6e7e99b136660f71bcded169e6ae26f40e75a56f68e11dd7`；
Markdown `4982117d1d553adb15a30690fcbc3f2a4fabbe8b34655a6294791fd4b8a70764`。
正式四Go來源、性別TSV及binary與審查記錄SHA保持；規格審查後只補實際驗收與採認狀態，契約不變。
獨立核89組正常GUI pair，其中6個明確退出pair無殘字、9個重入pair原生男位置正確。

## R990 女性正常路線驗收

沿原READY契約及未變的R989正式binary，正常加入既有TEMPEST、新遊戲、SAVE後，以新程序正常冷LOAD接續。
私人唯讀探針量到F字碼70、font0、runtime0110:3877→0110:6159、入口90／6，返回96／6、SP65312。
before／after只變11個色1點，均在原6×6格；八背景點符合本契約。
兩側1,415發布的RAM／Regs／Step／Ticks／indexed／palette摘要逐筆相同，完整trace bytes相同，保存層保持。
新遊戲／SAVE及冷LOAD兩段各一趟同binary正式重播正常0；此正式重播與私人原API控制分列，沒有改wire。

倚天／Unifont各2及3倍的正常冷GUI，均完成女性人物頁、法術開／取消／再開／再取消、EXIT迷宮及重入。
70個完整HD／Draw配對，58,880性別域像素逐RGB等獨立原生女字及本輪F背景期望，差異0。
其中女字域51,072、退出還原7,808；四保存層bytes／mode／mtime保持，正常關窗0。
六張F1 Draw不列遊戲配對；Unifont2倍410M只有HD、無Draw，另列未配對。
窄F與正常GUI兩份獨立FINAL均0／0／0，主代理完整回讀採認；只增已量TEMPEST女性路線限定CONFORMED。
2倍世界人物頁一般屬性標籤行距不足仍待，不稱全頁排版、其他人物或全遊戲完成。

窄F審查JSON／Markdown SHA-256：
`a951588730069e5af937ce5459aaec836bd0afe93e9c86c37f6a87a97a800168`／
`615b8e3f0cd577340e867821d32bd28c5eb55628cc9c02e3d2e171ea0bcb422a`。
GUI審查JSON／Markdown：
`026be067bc3c6c707bee4e2d7a4873589cd023152f85466c7e3af0fd18efb601`／
`55b1fbec3c96a820317ea2fb185ea6d6a99eade4d98a3f6904e663d5f246104c`。
原版、來源、工具、正常保存及完整限界見[R990證據](../re/025-ui-localization-coverage.md#r990-normal-female-profile)。
