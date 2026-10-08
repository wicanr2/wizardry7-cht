# 029 防拷題目後的英文答案提示

狀態：限定CONFORMED，未量題目與防禦GUI分支維持READY。2026-10-08，[#45](https://github.com/wicanr2/wizardry7-cht/issues/45)，
入口[R979](../goal/R979-clearance-answer-hint.md)。使用者明示題目後附答案、代理正常遊玩，
取代先前答案讀取／顯示及代理作答限制。原版驗證、資料、規則及輸入API保持。
已量正式結果入口：[研究025](../re/025-ui-localization-coverage.md#r979-clearance-answer-hint)。

## 已量證據與限制

Archives版DS.EXE SHA14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db；
Go1.24.13、dosgolem本地已核來源，原版資料唯讀、獨立Scratch。
正常冷啟動272M、ADD CHARACTER、選既有第一位角色、OK、CANCEL返回七項選單、START NEW GAME。
沒有原版記憶體或state注入。首個實際題目由原版正常字元輸入驗證成功，隨後到開場故事畫面。
本機證據workplace/r979/cold-probe/start-game.json及answer-normal／answer-settle的JSON與PNG；
含答案的收據只留本機，不引用答案值或詞級雜湊。

- confirmed：MSG取得入口resident0110:0582，stack第1字為ID、第2字為目的offset；
  沿R971及013原始定位，只在ID2500–2574時取得目前答案資料。
- confirmed，這個冷題目：0110:0659位於文字解壓完成後，DS採取得入口保存值，
  目的緩衝讀到NUL，最多64 bytes、不跨16-bit offset。
- confirmed，這個題目：記錄前四byte的索引0／1／2／3對應顯示的頁、段落、行、序數。
  第4byte起為英文答案，此答案由原版正常輸入接受。最初未讀頁面行的推論由新GUI與四行原文收據補證。
- confirmed：SECURITY CLEARANCE文字在106,148；PLEASE TURN TO PAGE <頁> OF YOUR在70,158；
  MANUAL AND ENTER THE <序數> WORD在64,164；LOCATED IN PARAGRAPH <段落> LINE <行>在70,170。
  frame-pair/question.json保存四個完整原文事件；gui-question-r2/question-stable.png實際呈現中文四行。
- 未量：其餘74題目的正常上屏及驗證，不以同形資料認作全部原版驗證成功。

## C1 只讀目前選定記錄

新的狀態由機器owner持有。每次0582進入先撤去未完成取得；白名單ID才保存DS、目的、SS、SP、Step。
同次配對在0659核SS／DS一致、SP=進入SP−0112h、BP=進入SP−2，
SS:BP+4與+6須仍為本次ID及目的offset；offset與讀取長度不得16-bit wrap。
這項關係由frame-pair/question.json動態證實：同次1916步，SS=DS=1106h，BP／SP及兩引數一致。
0659沒有active白名單pending時直接忽略，保留selected。只有已存在的白名單pending失效、
重入、SS／SP／BP或引數不符時，清pending及selected並記有限理由，不讀目的正文。
0582的新邏輯併入已有故事讀取callback，在其既有早退前仲裁，不新增覆蓋同位址的hook。
0659只消費未完成取得，完成後清空pending；NUL長度5–36、prefix索引0–3大於0，
答案須為1–32個ASCII大寫字母。越界、空值、錯來源或超時1M步保留原圖並記限定缺譯。
1M步只限制pending。完成後將驗證過的prefix與英文正文獨立複製為不可變selected，不保留Oracle buffer別名，
不晚讀DS:dst。selected可由開機預載，保留至題目首次呈現，不因等待超時丟掉。
新白名單取得開始先撤去舊selected；成功替換為新代次。成功提示後題目撤去即清selected。
其他ID只清pending、不讀目的正文或任意清掉已完成selected；不取得鄰文、詞庫或整表，不預讀75個答案。
不寫原版記憶體、A000、離屏、檔案或亂數，不改原版返回／比較。

## C2 題目完整性與撤除

只有shown事件同時含上述三條及PLEASE TURN TO PAGE <數值> OF YOUR（70,158），且dy=0、無dead，才建提示view。
序數僅FIRST到TENTH，頁／段落／行須無符號正整數；解析不可用即不畫。
view所讀目前prefix索引0–3必須與當前題目頁／段落／行／序數全相符；不依Step猜題目、不按最後印字任意沿用。
缺任何事件、來源／數值不相符、英文模式或缺完整字模，發布nil，下個快照撤除。
首次尚未看到題目時保留開機預載selected；成功匹配題目後才以題目撤除清selected。
提示不改一般問題的英文墨點還原、中文字形、事件生命周期、輸入框或hit rect。

## C3 字形與顯示位置

文字為「答案：」加原版英文答案，不作翻譯或正規化。字色跟本題最後一行的fgRGB。
先直接量目前題目中文末行的字條寬；提示接在末行後8物理像素，同行放不下時整個提示nil並記有限缺譯，不裁切或縮字。
支援2／3倍；尺寸或倍率不符整提示nil。既有問題排版保持，不為提示重排問題。
「答案：」用一般後端15／16字形，英文用8×15／8×16。完整字格、筆畫及字距須在
題目最後一行同帶內，且不進輸入區。安全矩形原座標[70,170,263,178)，物理座標皆乘倍率。
底色由frame-pair/question.a000與.pal確認；原右內側x263為底色，x264是邊界，提示採x263排他上界。
輸入區上邊框由原色號量到y182，白色輸入內框y183，與提示同帶分離。
Unifont末行13字形完整字格232px，倚天219px；兩倍率字形不縮放。
實際3倍樣張在workplace/r979/prototype/hint-3x.png，提示位置450,510、92×16，依目前GUI中文末行接續。
提示字距固定2物理像素，完整寬=sum(各字寬)+2×(字數−1)，高度=後端15／16，不乘畫布倍率。
本題Unifont92×16／倚天89×15；字形、位置與安全矩形皆按此算法共用，不從F1繪製副作用推寬度。

字模取同次help資源中完整字形，不在遊戲讀主機字型或烘製答案字串。
初始化需與一般字模高度一致15／16；英文字寬8、答案標籤字寬為對應15／16，逐字行數／有效位元及非空驗證。
未知一般高度、缺任何A–Z或標籤字形即整功能不可用，普通／錄放共用同一驗證結果。
bake_help補完整A–Z字形及答案標籤的字元覆蓋，不新增F1頁面或文案；原help內容與字模保持。
help的捕獲SHA涵蓋新字模；錄放原六roles、schema及同次bytes保持，不新增未綁定sidecar。
GUI只讀同一快照的不可變view；HD與Draw共用RGBA字條、位置、顏色，先取遊戲區再畫。

## C4 原版驗證及驗收

提示只畫字，不自動輸入或替原版判對；代理遊玩時可用正常鍵鼠與字元API輸入，另存私人收據。
至少一個正式正常GUI題目有提示、完整筆畫、正常輸入通過與撤除；同起點提示開／關原版
RAM SHA、Regs、Step、Mouse、indexed、palette及持久資料相同，像素只變批准提示區。
2倍提示須另抽正式GUI、完整筆畫與撤除；3倍樣張不能代替2倍驗收。已量結果見末節限定正式驗收。
資料、題目、字模、模式、座標、來源不完整的負向案例須整則不畫，不洩露答案到錯誤log。
全cmd/wiz7回歸、同次錄放資源解析及40px頁尾守門；未觸發的GUI分支只記內部測試。
原版存讀檔與後續文字覆蓋按014／008另驗，這一個題目成功不當完整中文化或全部防拷題目完成。

## 已補證據及READY前審查

1. 獨立預審兩阻擋已由frame-pair的SS／SP／BP／引數、原色號安全帶與新GUI／樣張回填。
   JSON SHA4813953907c5f73b257b81ea3972a5a0202863db0f4228c6170e250070836ed5，
   Markdown SHA05445f217feef185b5d964649cc9704cf7b8f894eb88f04e2b0e590de7351cc0，入口workplace/r979/review-evidence。
2. pending與selected壽命已分開，四行匹配且新題目取代；空間／字模失敗整提示nil。
3. 資料只透過已捕help的同一型別與結果初始化，不重開檔；同一view和rgba/events一起發布。
   六roles、原版ID只為目前資料取得來源、譯文英文鍵及原版判定保持。
4. 完整獨立複審阻擋／應改／建議0／0／0，主代理已回讀採認READY。
   本機報告workplace/r979/review-spec.json SHA5b1451000fc679ea106a6a5ae929db222f3cf512aeaacfbe714ad95cf7b57e00，
   Markdown SHA553153362920bed017eaf146a43957f81383fbd2eefdc57d94c189d9f322439f。
   此為實作前READY採認；後續正式驗收如下。

## 限定正式驗收（2026-10-08）

四組兩後端兩倍率正常GUI的原版RAM SHA／Regs／Step／Mouse／indexed／palette／comp及持久資料相同，
6,656,000 HD／Draw像素逐點與獨立原字型期望一致，區外差異0。
兩後端兩倍率均正常前端按鍵核驗、撤除及關窗exit0。
末90來源binary補讀前SP守門，另量Unifont3倍同題A/B與2,304,000像素、正常輸入／撤除；
兩階段明列於研究025，前輪四組不冒稱全部由末binary重跑。回歸2.421秒通過，
包含已捕help替換、正式六roles解析與真正Draw錄製頁尾／撤除。

最終獨立實作複審0／0／0已完整回讀採認，四項必要缺口均解除。
workplace/r979/review-implementation.json SHAf9947c8366719d517a7b255782ad6ffe80b04456546c4a967ebb308a438b563e，
Markdown SHAddc606d0ba80f3e7419237bb30d5ec8e85953afc68a6edcfe042cfd25cf487ba。
只此已量題目與正常路線CONFORMED，74題及防禦GUI未量，不阻擋已定契約實作。
014保存／冷LOAD與008全文情境繼續，整體中文化尚未完成。
