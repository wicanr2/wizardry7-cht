# 028 定版紙頁亮字與六按鈕毛筆資源

狀態：READY。2026-10-08，沿#22／#28及[R977](../goal/R977-approved-ui-wiring.md)。
使用者確認R976兩種字樣；本規格只整理正式資料與接線契約，不重新選外觀。
READY前不得接production。原版不寫、控制語意隔離沿001，來源見證沿025，主選單訊號沿020。
資源邊界預審已採認：`workplace/r977/review-boundary.md`，SHA
`016d405caf42fdc1cf10e407309d01e1c15972e98f9e42b8173eae86a17405a2`。
其阻擋由R2同次解碼及失敗語意、R3遊戲區、R5實際像素驗收承接；此預審不代替完整READY審查。

## 來源與證據等級

- confirmed：020六個英文鍵／精確座標及109×9按鈕尺寸，原文字路徑resident0110:3895。
- confirmed：025的3921完整轉送見證，返回0110:5b9b／5afd、完整矩形0／0／319／199、SP＋2。
- confirmed：R973／R974六張紙頁indexed／固定palette完整SHA；步數只為出處，不作辨識條件。
- 定案顯示契約：R975甲的30／32字高、水平加粗1、位置及R976色值；紙面、字帶與作者區屬已採顯示幾何，不稱原版無字底圖或精確紙頁模型。
- 定案顯示契約：R976 AR PL UKai TW／Book字面2，2倍16／3倍24像素、字條置中及米金字色。
- 已量輸入工具為Python3.13.15／FreeType2.12.1，固定映像933b46a028fd786c9c3d426ebabc237e29a15912231ea8de576e95f0e4f41a4c；Go接線工具復原1.24.13。
- DS.EXE Archives版SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
- 未量：新正式HD／Draw、正常GUI、故障回退及資源錄放。紙頁2倍尚無大字排版契約，保持原圖並記錄。

## R1 正式文案及衍生資源

新家族分別用`text/main-menu.zh-TW.tsv`與`text/intro-paper.zh-TW.tsv`，
表頭沿獨立家族`en／zh／context／dy／source`。英文原文為鍵，MSG ID只記出處。
主選單六文案與現有catalog保持，紙頁五鍵為presents／fantasy／role-playing／simulation／by，
依序為呈獻／奇幻／角色扮演／遊戲／作者，作者姓名維持原版英文。
Context分別`main.menu.brush`及`intro.mon61.paper`，dy=0，source標示已量輸出入口。
UTF8、欄數、鍵唯一、控制碼、原文及中文完整性必須lint；缺字中止烘製。

兩族構建時烘製，不在正式遊戲呼叫FreeType或讀主機字型。
毛筆字模源SHA cfc758e69386a21bade8cb3b373fda20d7f2fac6e8904ced2ded98588b63e274，
授權SHA1f11c4a22865e1289f4694dad56b3561790a4b7341a511ddcd0edddcc3b95a46。
字型及衍生字模保持Arphic原授權，與程序RRSAL分開；原始及產生字模不進Git。
紙頁源字模依007兩後端，與一般字模高度一致15／16，不改一般字型資料。

## R2 資源容器與解析

在已有全遊戲字形字典xlate.json追加`__main_menu_brush__`與`__intro_paper__`兩個具名顯示元資料值。
這是兩族字形／顯示參數，普通英文鍵、g／數字槽／序數槽保持；兩族各有schema=1及context，未知欄位拒絕。
錄放原六個resource roles與wiz7-input/1 wire schema不擴項，仍從同次xlate bytes抽取新家族。
現parser已有digits／ordinals的同層分支，兩族需新增明確分支並從普通trans移除。
普通與錄放共用同一xlate bytes解碼函式；錄放仍只捕六檔的原bytes與SHA，不重開檔讀新族。
舊xlate沒有這兩鍵時沿舊顯示；原schema與resource roles／順序不變，不更新wiz7-input/1。
頂層JSON或普通字形／digits／ordinals解碼失敗沿原錯誤處置；兩可選家族驗證失敗只關其家族，
正常／錄放使用相同結果並記限定缺字樣。xlate Loaded表示主資源已解碼，不代表所有可選字樣已上屏。

主選單元資料含font_sha256／license_sha256／source_sha256、固定字色[248,232,160]及六個entries。
每entry有完整en／zh，且含scale2與scale3字條，width／height及base64的8位alpha。
alpha以原生16／24字模、advance＋1、原R976字模metrics合成完整字條，
不縮放或截字。字條矩型取完整非零alpha的包絡，字色由渲染器套固定RGB。
每條最大寬109×scale−4、高9×scale−2；alpha長度必為width×height，非零且有界。
必須六鍵×兩倍率完整、無重複、無未知倍率，source SHA格式及字型／授權SHA符合定版來源。
正式JSON欄位為`schema,context,font_sha256,license_sha256,source_sha256,entries`；context=`main.menu.brush`。
entry欄位為`en,zh,stripes`，stripe欄位為`scale,w,h,a`；a為JSON base64位元組，不接受額外欄位。
字色為定版渲染常數，不由資源任意改色；en／zh固定六鍵與既有六文案。

紙頁元資料含source_size15／16、source_sha256及五鍵完整原字模g；只允許定版en／zh，
每字行數、有效位元與非空須逐字驗證，不把動態人名或原版素材寫入資源。
R974幾何模型完整SHA為1a9632b661b106367adffc2a8336808fe38202e60080620e897072e2fdcf5a22，
由受版控模型提供；正式場景索引使用六個indexed SHA與固定palette SHA。
正式JSON欄位為`schema,context,size,source_sha256,model_sha256,entries`；context=`intro.mon61.paper`。
entry欄位沿`en,zh,dy,g`；dy0，size與普通trans中的每字來源高度一致15／16。
普通trans為空或-xlate空的英文模式不能啟用新字樣；未知一般高度不猜後端。

家族資料各最多1MiB；檔案、JSON、單一值、未知欄位、base64、尺寸或字模不符，整家族不可用並記錄。
新家族缺失時主選單保留既有020字形，紙頁保留原圖。英文模式兩新家族全部不畫。
元資料鍵不得進入普通翻譯map、語意比較、存檔、檔名或原版記憶體。

## R3 主選單來源、發布及繪製

只接受020的六個shown事件同時有效、en／x0／y0精確相符、dy0。缺一、死事件、不同座標或非2／3倍不使用毛筆族。
只抑制這六個事件的一般點陣字畫筆，不改其原英文筆畫還原、事件／原字模或生命週期。
低解析comp仍由原overlay還原英文墨點；新字條在放大後畫布另畫。
按原按鈕矩型中心置中，左上為button×scale＋(button物理尺寸−字條尺寸)÷2，除法向下取整。
按鈕內邊留至少1物理像素，字格及完整筆畫不能裁切。

快照owner在evMu內依shown建立完整獨立view，與rgba／evs同批短鎖發布。
Draw只讀view，不讀Oracle，不以時間／步數猜畫面；未知或不完整六事件發布nil。
HD與Draw共用alpha字條、中心與RGB；每通道用整數`(fg×a＋bg×(255−a)＋127)÷255`。
字條只畫到正式遊戲區，不覆蓋錄製40px頁尾；倍率或畫格尺寸不符的那一格不畫。
Draw在所有新字樣繪製前取得320×scale／200×scale的遊戲區SubImage，
不能把200×scale＋40錄製整窗傳入精確畫格守門。頁尾及F1／地圖的原繪製順序保持。
下次快照移除訊號即撤除，不沿用上一頁的毛筆圖。

## R4 紙頁完整畫格、補底及字形

沿現025的完整轉送入／返回見證，不新增原版入口，不覆寫已有arm dispatcher。
最新完成indexed SHA必須與當前64000色號相同、mode13h及palette固定SHA相符。
六個紙頁indexed SHA沿R973表；前五張中文，第六無字畫格只撤除，不能把next檔載入當翻頁。
未知SHA、部分轉送、錯palette／mode／長度、錯返回SP、資料不可用，當次保留原圖並發布nil。

R974每字帶與紙面polygon交集才補底。取樣排除所有字帶及作者區，
同紙面內只取R>G>B、5G≥3R、G−B≥16的像素，按x／y奇偶分四組，各至少8點。
每個目標點取同組最近8點，距離相同按y／x排序；權重`1/(dx²＋dy²＋1)`，RGB加權平均取RoundToEven。
補底為新繪近似，存在條紋接縫，不當原版無字底圖；原indexed／palette／A000／離屏全部不寫。

3倍畫布：每個原字模on點展開2×2，再向右聯集1像素；單字實際外框為2×size＋1、高2×size。
字起點advance固定2×size＋4，n字完整寬為n×(2×size＋4)−4＋1，高2×size；不因加粗再增加advance。
位置沿R975 layout：732呈獻中心199,68；753四行奇幻93,32／角色106,52／扮演126,73／遊戲145,94；
100999作者173.5,88，其餘中心沿受版控R974模型。中心×3採RoundToEven，再減完整字條寬高整除2。
紅字[236,72,40]，原深色族[248,232,160]；完整字格與筆畫必在定版紙面且不侵作者區。
2倍紙頁大字沒有定版布局，此增補先保留原圖並記限定缺譯；不把此狀態列為已中文化。

與原025六張MON60呈獻分開，不改其圖形、補底或字模。紙頁view與完成畫格一起生成，
HD及Draw用同一組補底RGB／字條／顏色／位置。Draw不重新取樣原版或重算畫面識別。
新snapshot不合格時nil，保留MON61載入前仍有效的上一完整圖直到它真正重寫或完整翻頁。

## R5 最小充分驗收

1. 獨立源字模／FreeType C API與烘製alpha逐點相符，TSV與已采R976文案及10紙頁／2選單樣張一致。
2. 兩族解析守門用自製不含原版語料的負向案例；英文模式、缺資料、錯來源、錯倍率／尺寸、非完整六訊號與撤除。
3. 同原始來源正常無輸入冷啟動，保存修前英／中及修後英／中：原版RAM SHA、Regs、indexed、palette、Step、Mouse同。
   選單2／3倍與紙頁3倍HD／Draw各與獨立期望相符；變化只在批准墨點還原／字帶／字條範圍。
4. 正常GUI3倍片頭五頁／完整翻頁撤除與主選單毛筆，分別以-scale 2／-scale 3啟動，核2倍主選單、尺寸不符守門、F1返回及正常關窗。
   不用原版記憶體／state注入，不回覆防拷或修改存檔；未觸發排除路徑維持未量。
5. 普通模式與正式錄放同次xlate bytes派生，resource Hash與Loaded結果保持規範；完整cmd/wiz7回歸。
6. 字型／原版／SDK來源保持，只寫本輪輸出，自有容器與擁有權清理；不把靜音GUI當音訊／人耳驗收。

新家族的正式完成只限實際驗收的畫面，020／025整份不因局部通過改全路徑CONFORMED。

文案入口：[主選單TSV](../../text/main-menu.zh-TW.tsv)、[紙頁TSV](../../text/intro-paper.zh-TW.tsv)。

完整唯讀審查見`workplace/r977/review-spec.md`與`.json`。主代理完整核閱；無實質阻擋。
字距公式已採定版值，倍率驗收文字改為兩種CLI啟動及尺寸負向測試，不新增切換功能。
依AGENTS文字層應改修完可READY，不重問已確認外觀。實作及正式驗收仍待，未標CONFORMED。
命名勘誤：審查時暫名027-approved-ui-fonts.md，與既有027重印持續性規格碰撞。
本輪改為028，契約保持。歷史唯讀審查報告的027均指本規格當時的SHA，不指既有027。

## 已量範圍

R977正常無輸入冷啟動：兩後端紙頁3倍五構圖及無字撤除、主選單2／3倍已通過同狀態、獨立全圖及真GUI驗收。
此正常路線CONFORMED；其餘未在GUI觸發的排除分支維持READY，不把內部守門單測當全路徑CONFORMED。
完整cmd/wiz7回歸及最後來源身份抽測通過，實作／資源／已量收據獨立審查0／0／0，兩報告完整回讀採認。
報告入口workplace/r977/review-implementation.json及.md：JSON SHAe8666813d3a899e843f5116e982540ca432ad3464bb67da880a17c31599af5cc，
Markdown SHAc3bc7967a8a48a4b72080ecf3d25a73ab4e5b6d8fec808629a08f8cac96a73e6。
研究、源身份、限定範圍及重跑入口見[025的R977](../re/025-ui-localization-coverage.md#r977-approved-ui)。
