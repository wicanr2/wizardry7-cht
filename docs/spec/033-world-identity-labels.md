# 033：世界人物頁性別與選人提示

狀態READY；R988已量WHO正常路線限定CONFORMED，未量分支仍READY；正式契約條1至15；2026-10-09；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
入口[R986](../goal/R986-world-identity-labels.md)，既有畫面證據[025](../re/025-ui-localization-coverage.md#r984-normal-profile-context)。

## 為什麼要做

正常點隊員姓名出現WHO?提示，點頭像進世界人物頁後性別仍顯示M。
這兩個已見英文殘留阻礙完整中文覆蓋。只改覆繪，不改選人、人物資料、原版文字或遊戲規則。
沿既定性別譯詞男／女、原生15／16字模與2／3倍畫布；不擴大主選單毛筆字範圍。

## 已量與未知

- confirmed於R984正常路線：WHO?及人物頁M確實可見，取消同隊員選取後WHO?撤除。
- R986第一探針沿同一正常錄製窄觀察，原版RAM／Regs／Step／Ticks／indexed／palette的1614發布trace全bytes等於raw控制。
  目前窄框WHO的一般NUL及0110:3877單字元觀察均0；這只排除本次已觀察框內入口，不據此直接稱圖片。
- confirmed於目前世界人物頁：0110:3877入口，返回0110:6159，SP65310，字碼77、font0、x90／y6。
  後續連字號返回0110:6167，HUMAN一般字串返回0110:6191；與024角色選單的0110:6c42分開。
- M返回已量13個改成色號1的墨點，全部在90／6的6×6格。完整背景守門及兩倍率字格由[034](034-world-profile-sex.md)承接，獨立審READY後已正式接線；正常F仍未量。
- WHO?已量到45個色號22的像素。窄寫入觀察在runtime740c:20fe捕捉透明位圖寫入，來源寬24、高5。
  IDA9.4解碼保留原operand `mov es:[di],al`，raw snippet EA與runtime位址分開。
  R987追加勘誤：5是字母墨點高度，完整WHO文字末區24×6、整張游標圖24×20。
  已量實際返回740c:1568及740c:0166；1568清理以1575的retn結束，舊40b0:0166 far返回假說未成立。
  WHO來源圖、原生字格及完整矩形已補，正式契約與正常GUI仍待。入口[R987](../goal/R987-who-call-witness.md)。

runtime位址沿Archives版，DS.EXE SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
來源Go1.24.13、固定dosgolem worktree；私人API只用於固定輸入診斷，不改wire或宣稱正式相容重播。
沒有修改原版記憶體、索引幀、色盤或保存。完整原版／輸入／畫面只留ignored workplace/r986。

## READY前必要契約

1. 兩路各具英文原文／字碼、caller、font、原座標、前景色與像素／背景守門；不得只按M或畫面字形放行。
2. 性別沿已證單字元入口與返回見證，拒錯SP、座標、字碼／font、非空差異或未批准範圍改寫。
3. WHO先確定實際輸出來源；若為圖片，核來源呼叫及原字墨點與補底，不冒用NUL文字事件。
4. 譯文保持英文鍵及來源情境分離，不加一般M鍵或把譯文送進原版比較／保存。
5. 字格與左右上下文字安全區須用兩後端真字模核；不縮字或遮到種族、職業、其他隊員及數值。
6. 同狀態兩倍率、雙字型正常選人／取消／人物頁、HD與Draw、存檔保持及錯守門負測後才驗收。

性別由034獨立放行。WHO?已定位游標來源圖，R988正式契約已採READY，正常GUI仍待，不擴查無關驅動程式。

## R987 WHO游標圖與字格證據

confirmed於R984正常選人／取消路線：整張游標24×20；資料在runtime740c固定欄位，
15ac／15ae為本次30／53位置，1e5c／1e5e為24／20，1e60／1e62為10／10，
1e64／1e66指向當刻84d3:0002。字母墨點在圖的14至18行，末六行14至19為24×6。
來源144 bytes SHA-256為1e34be319abc13e44f2b8b3eeb998bed1dfeb53c37e8ae369cb415546671aa61，
含45個色22、86個0及13個透明255。原圖可見WHO?，比對鍵沿此英文輸出。

45個色22墨點與indexed吻合，正常樣本的去墨底色為0。
透明255處保留補底前原畫面，不能拿背景buffer全0當完整矩形期望；新字模可在文字域占用原透明位置。
字帶前兩行12／13的48-byte來源SHA為575e226442eb3e8e3d52d5374651cd9c66e202197af0696a7cae32ac2c7d8a6a。
行12非透明手尾及黑輪廓在列3至9，全部保護；行13只有0及255，是空文字背景，不能外推indexed全0。

兩字型原生「誰？」已有正式flow字模，15／16全行保留、不縮不裁。
可丟棄原型以一像素字墨間隔利用標點留白，在原24×20游標矩形內排字。
2倍倚天x11／y25、Unifont x18／y24，3倍x22／y42及x21／y42，皆相對游標放大後左上角。
四派生HD核手尾及黑輪廓變化0、原圖外變化0，不擴黑底。這些不是正常GUI或CONFORMED。
先前加黑底3／4像素及2倍維持英文的樣張只留本機研究，不作正式方案。

R987最終獨立來源／原型審查於FINAL後完整回讀採認0／0／0，32份已讀來源SHA保持。
四派生HD共1,664,000像素與獨立期望相符，手尾7非FF及文字域外變化0。
補底階段保留13透明點原RGB；新字覆蓋原非零透明點倚天2／3為3／1、Unifont2／3為2／1，不作全不變宣稱。
JSON SHA-256 d97f39354668e8c7ffa58e1d79ceff4711e00e5574004316b35e5cf22f5c9634；
Markdown 17d977356abab9792ba68712b0f2e3ad08a2539948c8801fdb8a3c7e0c198e0d。
仍僅DRAFT派生原型，正式來源、資源、正常GUI及CONFORMED未放行；WHO字型／畫布沿既定15／16、2／3倍。

## R988 正式WHO顯示契約（READY）

<a id="r988-who-cursor-contract"></a>

本節只放行已證WHO圖形；性別沿034。原生字型及2／3倍畫布沿既定設定，不新增字型、黑底、控件或輸入功能。
原版DS SHA沿本規格；合法R984輸入SHA為d7204e01e17698e0a03a764e4521f7c391b03c2e2fe326683c690edcb723e453，
控制1614發布trace SHA為84388c8f8cf1461fc7d6e5bde6fd2d05597ba02855c3ce505ac76ea107c404ae。
WHO來源tail及prefix SHA沿R987段；Go1.24.13。完整圖、code、字模及輸入只留本機。

### 辨識與生命週期

1. 快照owner唯讀原版，不新增OnCall替換或Stub。只在mode13h及完整64000-byte indexed讀WHO狀態。
2. runtime740c欄位1e5c／1e5e須24／20、1e60／1e62須10／10；15be／15c0須24／20。
   15ac／15ae按有號16-bit解本圖左上角，須x≥0、y≥0、x+24≤320、y+20≤200。
3. 1e64／1e66為當刻來源offset／segment。讀前核offset+480≤65536及實體來源終點≤A0000，
   不讀VGA或跨段。只讀offset+336的144-byte tail及offset+288的48-byte prefix。
   兩SHA必須同R987已證WHO；不只按尺寸、顏色或通用driver caller放行。
4. 原版字母位置為圖下14至19行；tail所有非透明0／22像素須與indexed相符，255只作透明、不要求背景等0。
   不完整、未繪完、隱藏、來源不符、越界或部分裁切時不建立view，保持原圖。
5. 每次compose重新取得不可變view，不持久保留老坐標／字樣；原版取消／移動後由新快照撤除或重建。
   Who view與source資料不得寫回原版、離屏、A000或存檔。

### 譯文與資源

6. 英文key WHO?、中文誰？由text/cursor-who.zh-TW.tsv唯一供應，context為cursor.party.member.selection、dy0。
   生成工具由正式flow內相同後端的兩個原生glyph增補既有xlate字典WHO?項；
   不改schema、CLI、其他鍵或六資源身份綁定方式。字元已在正式字庫，不另提交字型二進位。
   ordinary及recording沿同次已捕xlate bytes解析WHO，不能另外重開檔案；原始資源SHA仍依六資源既有綁定。
   WHO?只供游標資源，取出後從一般NUL查表移除；不得因此放行未批准的一般WHO?印字事件。
   WHO entry族內不合格只保原圖，根JSON解析錯仍沿既有resource_parse，不改原錯誤契約。
7. WHO entry須dy0、無數值槽、恰兩個同高15或16非空方字模，所有row位元在該字高內。
   缺失或不合格保原WHO並記缺譯；空xlate的明示英文模式保持原圖。

### 布局與像素

8. 完整native glyph不縮、不裁。令H為字高、A＝第一字最大墨列+2−第二字最小墨列；
   A須正且兩字實際墨區至少隔1像素。字條全字格寬A+H，WHO2／3倍內須容納。
9. 相對游標，字條y＝14×k−max(H−6×k,0)。候選x從0至24×k−(A+H)，
   優先最接近中心，距離相同取較小x；字條全字格須在原24×20放大矩形內。
   y落第12行的新增墨點只能用prefix第一行的透明255；手尾7個非透明點連黑輪廓都不覆。
   prefix第二行已證只有0／255，屬文字空背景，不推定背後RGB全0。
10. 去墨只在comp副本把原45個tail色22位置改為本WHO字帶底色0；
    tail既有0及透明255位置不整片改底，新字可在文字域占用原透明點。
    不讀／泛化整片背景buffer為0，不擴黑底。新字前景用同快照palette[22]。
11. HD與Draw共用布局、完整字模及RGBA；Draw在一般中文之後、F1／地圖之前畫WHO。
    依各通道實際k=2／3在owner構造view，GUI-scale與hd-scale可不同，不能用另一通道倍率代替。
    compose保持原WHO字樣；只有該通道有合法view才在自己的comp副本去45墨／畫字。
    該通道view nil、錯倍率或尺寸不符時完全保原WHO；不先擦共用comp再讓拒畫通道失字。
    字條RGBA為advance+H乘H、Alpha只在原生glyph墨點為255，其餘透明0，色取同快照palette22。
    WHO只畫完整gameArea Bounds內；40px錄製頁尾及圖外像素保持。GUI只讀快照的不可變view，
    原生glyph及source陣列皆自有副本，取消後不畫舊view；GPU快取有界且不能代替當刻view資格。
    群字元檢查或布局無安全位置時保持原圖並記缺譯。

### 驗收與停止線

12. source辨識負測涵蓋mode、尺寸、hotspot、坐標、span、SHA、非透明像素及未繪完；
    資源負測涵蓋dy、槽位、字數、高度、位元界、空字及布局碰撞。檢原buffer不被修改。
13. 兩字型兩倍率正常冷LOAD、姓名選人、WHO顯示、正常移動、同成員取消及再入；
    真HD／Draw與獨立原glyph期望核原圖矩形、手尾、控件外及非WHO畫面。
14. 固定原API同輸入比較RAM／Regs／Step／Ticks／indexed／palette及保存bytes／mode／mtime；
    不改wire檔頭、不稱正式相容重播，派生HD不代正常GUI。僅已量WHO路線才CONFORMED。
15. 部分裁切、來源／資源錯誤分支未正常量到時保持READY與英文fallback，不為這個輸出深入mouse driver或硬體時序。
    不改WHO選人、姓名原文、NPC輸入關鍵字、原版資料或判定。

R988正式契約獨立FINAL 0／0／0已完整回讀採認READY，T1／T2文字補足均解除。
JSON SHA-256 1597f4d04422bf39bb642ab0cd7d4311e909b2ac4636704dcb97c4a962d91cb0；Markdown 9503d02e16b2e46db6fab88b6764a3f41ea730b7e8071219c792ab2eb68a32f2。
只放行此契約實作，未稱CONFORMED；R987派生字樣不代正常GUI。

### R988正式驗收採認

最終獨立實作及像素複審於FINAL後完整回讀採認0／0／0，22份審查輸入SHA相同。
兩字型兩倍率正常選人、移動、取消及再入限定CONFORMED；部分裁切與其他未量分支保持READY。
28正常操作、34正常WHO文字域／手尾HD與Draw、四組1614原版發布及保存全同，
32派生完整圖13,312,000像素與獨立期望相符。派生HD不代正常GUI或正式相容重播。
JSON SHA-256 369c275a6bbda48ebaf0c5d857b27f43d5fe79e3bbdd206080e922a473f2506f；Markdown 41b69a5c38dcb7a1715199165890670300278f01914942de78009664b3663242。
結果與來源身份見[025](../re/025-ui-localization-coverage.md#r988-who-cursor-wiring)。
本次只追加實際驗收與採認狀態，不改READY契約或未量停止線。
