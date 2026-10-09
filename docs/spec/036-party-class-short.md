# 036：側欄職業短碼

狀態：R995已量左上及右下ALCHEMI路線、R996已量六欄Unifont2正常路線限定CONFORMED；其餘分支READY。入口[R995](../goal/R995-full-party-ui.md)、[R996](../goal/R996-remaining-text-ui-gates.md)、[025](../re/025-ui-localization-coverage.md#r996-career-late-line)。

使用者要求完整中文介面；ALCHEMI沿既定職業譯名「鍊金術士」，維持原生15／16像素，不擴大一般訊息或改原版資料。
文案唯一來源為text/party-class.zh-TW.tsv，英文鍵ALCHEMI、context party.class.short、dy0、source runtime:0110:1c71。
R995只增此短碼；R996另將四個已量職業全名接入相同幾何保護。其他呼叫端的普通覆繪保持。

## 原版與字形證據

- R995 r3原API正常路線：六職業欄均font0、ret0110:1c71、x9／261與y56／101／146，印字差異只在英文五行，前景132。
- 字段42×8內，前四欄於字串返回後會重寫y＋7；不能用中文字遮住這些後寫筆畫。鍊金術士原生倚天末墨列13，Unifont為14／14／14／13。
- 2倍Unifont採物理x＋2、y－1，原生字形不變，避開y＋7。頂列六墨點投影至姓名底行column3／6／14／21／23／32。
- 姓名入口ret0110:1c21，x9／261與y50／95／140；R4正常逗號姓名實測左上x9／y50、font0、返回X51／原Y／SP＋2，76個改變點都在42×6內，末行改變恰為column13／31。
- R4執行中字型段7B0E:0000的1040B SHA與原檔相同，末列遮罩獨立等原版差分。raw／probe1825發布trace完全同，SHA-256 `f561d094b4c53a9f366fe4a6eafd0407f7f956aee87b9fae3872cc2aadafa026`；正常改名後保存bytes／SHA／mode相同，自然mtime分列。
- VBFONT0.VGA為1040B，SHA-256 `937eda16ebf61f28506122c465cd3a3e112aa001cf0c37644b001bb0cbc35ae3`。16-byte header及128組字元映射之後，128字模各6列，位元7至2由左至右。這是本機原版身份，字型bytes不提交。

## 顯示契約

1. 只有普通字型一致15或16、視窗倍率及HD倍率同為2或同為3時啟用新短碼。資源不合格、錯caller／字型／座標／offset／SP保留英文並撤除舊事件。
2. 六個姓名錨點保持原版墨點，不經譯文查表。Name NUL限1..10個ASCII32..126，重用一般raw64的前32，off≤65504、SP≤65530。姓名開始時撤除同slot舊職業事件及舊遮罩。
3. 為2倍Unifont保護原姓名，從CS0110:364a唯讀2B字型段指標，再讀fontseg:0的1040B。每次讀前核段內不wrap及物理終址≤A0000；核固定SHA、header、mapped glyph<128及advance6，推導姓名末列前42欄。只留遮罩，不留姓名文字或字型bytes。
4. 姓名遮罩須在同ret／SP＋2／最多1M步返回後才有效。沒有有效遮罩或中文字頂列投影碰到姓名原墨點時，該2倍短碼完整退回英文，記錄原因；3倍與倚天不需要向上移字。
5. ALCHEMI以exact NUL英文鍵辨識，僅ret1c71、font0、六職業錨點、SP／offset界限及合格資源可建立pending。返回同ret／SP＋2且不超1M步；完整64K差異只能位於42×5英文域且所有改變後色132，至少一點改變。
6. 只還原原英文差異墨點，原版與VGA不寫。事件保存同slot姓名末行42B身份，compose沿既有游標ignore核對；不符即撤除中文。重新印姓名先撤舊事件，避免跨人物沿用遮罩。
7. 四字原生尺寸及字距2保持。倚天或3倍位移0；2倍Unifont物理位移＋2／－1。逐墨點驗字段右界及y＋7前界，向上一列只許經姓名遮罩證實不重疊的點。HD與Draw共用位移函式。
8. text/*.zh-TW.tsv正式文案烘製普通xlate鍵，既存其他entries、紙頁／毛筆及三個附屬家族保持。ALCHEMI在未證其他caller不啟用普通全局查表。

事件接線細節：Name pending按slot獨立保存ret／SP／entrySteps，既有arm先核Name返回再走原pendingPrint；失配不提交遮罩。
Class pending保存專用標記及entrySteps，在finish前獨立核完整delta，不沿一般事件猜前景色。
Name末行42B身份從class入點既有before取，compose以當次同列42B及既有游標ignore核對；位移只套class標記、size16及k2。
發布GUI副本時獨立拷貝partyClassWitness的純量及42B內嵌陣列，讓Draw保留相同位移與資料所有權。

## 驗收與停止線

獨立READY前不接production。內部驗界限、字型SHA錯誤／缺遮罩／頂列碰撞、事件撤除及雙字型兩倍率字形不縮。
沿兩份既存正常wire做原API同狀態，原版記憶體／indexed／palette保持；抽字段字形、姓名及後寫行域外像素，另以正常GUI確認滿隊與逗號姓名側欄。
未量其他職業短碼、其他人物位置、OEM／長名字、非2／3倍率及平台各自保留限制，不改存檔或跳關。

獨立READY FINAL 0／0／0完整採認，15來源身份相同，29見證指紋
`a654fc56943778f2408ad9c0a6a34ecec14df07c063d72b3cda0e5ccaf9f18b4`。
兩報告SHA-256 `55e8edfe80a9e19e1fb92d58c5a389eef1083b1e5c1d541b82f20208e31d4876`／
`60b9d21ac909914bec0afc4373a7a9ddcbbbc0e5df5e5298bd8c9f72ed816a7e`。
以上為實作前READY閘門，限定正式驗收結果如下。

## R995 限定驗收

正式資源只新增ALCHEMI，其四字g與既定ALCHEMIST相同，其他xlate值及三附屬家族保持。
普通正規化不得取得此鍵；初始化先驗配置，再抽至partyEntry並從一般trans刪除，控制碼近形負測已通過。

R2四個Name新配置、兩個Name舊3倍及兩個滿隊新代表配置共八趟原API0，14234發布各等相同輸入的raw。
兩3倍新舊16完整HD共9216000像素相符；獨立重建名單底行及側欄兩幅完整圖1152000像素，批准域外差異0。
10個字段ROI21840像素相符，與全圖分母分列。自然改名PCFILE bytes／SHA／mode相同，各次自然mtime另列；SAVEGAME保持。

兩後端3倍滿隊正常GUI及同879d重播0。R2 Uni2目視後的獨立畫素驗到404點失配，原因為GUI事件副本漏party標記。
R3只補四行獨立值拷貝，單測驗不alias及2／3倍位移，cmd2.284秒及兩包快取通過。R2失敗畫格保留，不當成功。
R3新Uni2正常錄製2844事件／1419快照／末307M，同b4a5正式重播0，保存metadata保持。
兩R2 3倍與新R3 Uni2真正GUI九原行ROI共8316像素符合獨立原版期望，包含姓名末行及y＋7後寫行。
R3另補一趟Name Uni2原API0，1825發布全同raw，8幀48檔逐bytes同R2，確認owner／HD保持。

最終獨立動態FINAL 0／0／0完整採認，31來源／結果身份及八個精確mtime值相符。
JSON／Markdown SHA-256 `bb57a03d906b7eb5259e86c28c34ba24a684dc0a7f6f344208f791da94468a9b`／
`cd438e46112285b9a965c6ed9b8486e38e32f78789b48e125dd7f938fed94c2f`。
只標已量左上逗號姓名及右下PENTAS的ALCHEMI路線限定CONFORMED；GUI3沿R2、Uni2採R3，來源界線見[025](../re/025-ui-localization-coverage.md#r995-full-party-ui)。
其他slot、OEM／長名字、錯資源等正常GUI與其他caller未量，仍READY，不稱全95來源或全遊戲parity。

## R996 既有職業的兩倍後寫行保護

狀態：本次已量六欄Unifont2正常路線限定CONFORMED；其他分支READY。入口[R996](../goal/R996-remaining-text-ui-gates.md)。不改譯名、字形、色彩或原版資料。
正常b4a5／Unifont2倍滿隊真GUI已證8個原版後寫像素遭覆蓋：FIGHTER左上2點、THIEF左中2點、PRIEST右中4點。
位置及原RGB／現GUI RGB在workplace/r996/gui-late-pixel-evidence.json；相同輸入的原indexed／palette來自R995原API。
MAGE原生字形同樣有第14墨列，但其已量底部位置沒有晚寫差異；不宣稱其他位置已重現。

顯示契約：

1. 只新增已量FIGHTER／THIEF／PRIEST／MAGE四個exact英文鍵，caller1c71、font0、六個職業錨點及SP／offset／返回1M步界限沿本規格。
2. 四鍵復用當次普通trans原生g，不另存譯文或重新烘製；它們在其他caller仍沿一般查表。ALCHEMI繼續用已隔離partyEntry。
   現行party配置不可用時，四個既有鍵保留舊普通路徑，不因可選短碼資源失效退回英文。
3. 在class caller，正規化能到這五個鍵但raw不是exact NUL時，保留英文、撤除舊同slot事件；此拒絕優先於配置不可用時的四exact鍵回退，不得由一般查表繞過新專路。
   其他原文／caller不新增專路，維持既有行為，不猜未量職業截短字。
4. 四鍵資源須為兩個一致15／16列非空glyph、dy0及無數字槽；ALCHEMI仍為四字。完整64K印字差異仍只許42×5及後色132、至少一點。
5. 原生字形不縮。Uni2的FIGHTER採物理位移＋5／−1，其他四鍵含ALCHEMI採＋2／−1；倚天與3倍位移0。
   固定字模頂列獨立核算：FIGHTER原HD列11／13／26，＋5後投影不碰任何逗號／分號的原字格column1；其他四鍵＋2同樣安全。
   仍須依本次Name真實末列遮罩逐墨點驗碰撞，不以全大寫姓名假定安全。字寬及y＋7前界沿本規格。
6. Name返回才提交遮罩、42B身份、同slot撤除及GUI獨立副本保持；HD與Draw共用依party標記及原英文鍵選位移的函式。
7. 不更動xlate schema、六個錄放roles、一般career字模、其他caller、ALCHEMI既有字條或原版存檔。

READY前只作可丟棄樣張及來源審查。正式驗收需新正常Uni2六人路線、8個晚寫像素恢復、Name墨點保持、同狀態批准域外0；3倍／倚天原位置保持。

獨立READY FINAL 0／0／0完整採認，13來源身份核同，18份原版見證指紋
`e20c4f30096d0c42634a6433a1424c77dd5401ebec5daad4440376111c7de44b`。
兩報告SHA-256 `f9e6c7df6a8caf1980bf6ee833b461e4a709a92c11664a581a724c88659c03dc`／
`35a1e673bfa17a84bfa748aaaa56b2c8093133be0023bad2c561c012795ea2b5`。
以上為四鍵幾何修補的實作前READY閘門，限定正式驗收如下。

### R996 限定驗收

source獨立FINAL 0／0／0完整採認，三包回歸通過，cmd2.072秒；實際原版字型與Unifont字模正向測試PASS。
正式95實源逐項核同，binary SHA-256 `5d6fe6f396fc714bf54e59aefa3c8b81dcdcd2021627b8931ba121d38e6df013`，資源沿R995。
兩配置私人old／new四趟0、3806對原版發布全同，indexed／palette及保存metadata保持。
30完整派生HD共12480000 RGB像素符合原版印字差分及既定字模期望，批准域外0；倚天3全圖保持。
新正常Uni2六人GUI錄製3691事件／1837快照／454.4M，同正式binary重播0，正常Stop續玩及關窗0，保存metadata保持。
真GUI六字段九原行共9072 RGB像素相符，8個被蓋晚寫像素全恢復原RGB；原姓名末行與ALCHEMI既有路線保持。

最後獨立動態FINAL 0／0／0完整採認，29來源／結果身份核同。獨立重建一張Uni2完整HD256000像素及六GUI字段9072像素全同；其他29張HD只另核實際RGB身份，不稱全部獨立重建。
JSON／Markdown SHA-256 `f40071af75b727e1bc0f891900312f77ca07ef1eebaa540042788af57c70f221`／
`8ae288d31a58ec532e624cca4b7b928fe282b996bc86e72c0d20ae788ac76d32`。
限定本次六欄正常Uni2路線CONFORMED；其他slot、姓名、caller、錯資源、NoDiff／重印與平台仍READY。MAGE只驗已量底部正常位置，不稱其他位置已重現缺陷。
