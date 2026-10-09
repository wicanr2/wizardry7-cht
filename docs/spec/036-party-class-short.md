# 036：側欄職業短碼

狀態：已量左上及右下ALCHEMI路線限定CONFORMED；其餘分支READY。入口[R995](../goal/R995-full-party-ui.md)、[025](../re/025-ui-localization-coverage.md#r995-full-party-ui)。

使用者要求完整中文介面；ALCHEMI沿既定職業譯名「鍊金術士」，維持原生15／16像素，不擴大一般訊息或改原版資料。
文案唯一來源為text/party-class.zh-TW.tsv，英文鍵ALCHEMI、context party.class.short、dy0、source runtime:0110:1c71。
本規格只增此短碼；其他既有職業全名的普通覆繪保持。

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
