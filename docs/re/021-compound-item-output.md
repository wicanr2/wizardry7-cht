# 021：物品頁兩行名稱與另一直接繪字路徑

日期2026-10-01；議題#22、#28；入口[R725](../goal/R725-compound-item-output.md)。

## 輸入與位址基準

| 原版輸入 | SHA-256 |
|---|---|
| DS.EXE | 14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db |
| VGA.DRV | f064349ce5dc694b26ed58615262d1bad893e59971bf259c0a74bdea9ee27a70 |
| SCENARIO.DBS | b2cb0722122724d35d379bb10250d0eec51238d9b147c36be12177ef3d2462f6 |

Go1.24.13／Python3.11.2，dosgolem d16f03a2d8d782024d470606e1f4cb5272b89cd6，
image sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7。
原名定位是SCENARIO檔案偏移896+72×idx的前32byte；idx只作出處。
CPU位址皆為dosgolem實際CS:IP／DS:offset，VGA.DRV檔案偏移為runtime偏移減0100h。
本輪沒有IDA資料庫反組譯，不把ToIDA寫入索引當作IDA語意證據。
DS.EXE的512byte MZ標頭加IP映射只對已核resident wrapper成立；
5cb9／5d19／ba40／ba63在原檔此映射全0，不能冒充載入程式的原檔定位。

## 可重播的正常路線

不透明重播既有合法輸入至260M，再270M(160,74)、300M(25,52)、380M(45,116)，
360M／400M保存安全角色頁。四個只讀探針共八次對照，
Step／Mem／A000／Pal／Mouse／Stats／Shown與r715/off-3相同；
A000／色盤／合成／PNG逐位元相同。Mem只比較FNV32，不冒稱完整記憶體逐位元相同。
新增讀值只限安全裝備窗口，未讀防拷詞庫或推導答案。

完整報告workplace/r717/report.json SHA-256
ef37e0563abe194c1da050928cf10b9565dfdb00acc21c1b62de680270d58d7c，
由報告索引輸入／輸出雜湊、探針與背景收據。輸出1000:1000，自有容器、
root-owned與誤建`.md`目錄均0；原版與研究快照只留本機。

## 三個兩行名稱（confirmed）

入口0110:3895，第一返回0110:ba40、第二返回0110:ba63；font0、x18、前景索引105。
第二y為第一y+6、SPentry均fe38；第一指標1106:fe44，第二指標為第一字串長度+NUL之後。
緩衝每件重用，不能當永久名稱或推定直指物品紀錄。

| 完整英文鍵 | 第一y | 合併域 | 第一返回／第二入口Step |
|---|---|---|---|
| LEATHER/CUIRASS | 128 | (18,128,41,11) | 382007777／382007813 |
| FUR/LEGGING | 142 | (18,142,41,11) | 382013470／382013506 |
| BUCKLER/SHIELD | 170 | (18,170,41,11) | 382025208／382025244 |

417筆原名獨立解析全與現行原名表相同。兩行英文筆畫外變化0；
第二返回→服務16前後的offscreen()觀察域→A000→400M的各段矩形變化皆0。
三域均有紋理，不能單色清底。兩字型×兩倍率十二個候選矩形都在已量域內；
排版可行性為強推論，尚未中文上屏，其他slash名稱或畫面未量。

背景來源界線：r717的first-before／after1／before2／after2與服務16的SourceBefore／After
都呼叫現行同一個`offscreen()`，讀`740C:01E5`的段再取畫格，不是直接讀當時regs.DS。
各stage的01E5數值未記錄；不能把三域相等升格成段值必為6289。
服務16實際movsb的DS6289另由逐次寫入trace證實，屬不同觀察端。
`r726/independent-fixed-evidence.json`獨立重算三域完整41×11、兩mask聯集174／124／160，
第一返回after1與第二入口before2相同；支撐沿用helper的限定域契約，不替換通用來源。

## LONGSWORD來源訂正

印字返回382033177的離屏位置(18,114)，最終A000字形在(25,111)。
新trace證實形成者是服務24：0110:5cb6→40a7→740c:0163→20fe，
20fe在382034681–382040004從DS851c寫向ESA000，130筆變化後114個字形筆畫完整吻合。
返回0166@382050891、40b0@382050892、5cb9@382050895。
其後0110:5d16→3921→740c:0143／05b5的服務16同座標複製71列；
新字形在此前已存在，不能稱為服務16搬移。394筆寫入重播舊值不符0，實際時間鏈為confirmed。

851c的bitmap與原字串建構資料流只是強推論；新位置背景恢復仍未知。
不把所有服務24 bitmap當文字，不以形狀搜尋直接掛中文；本輪未修此路徑。
下一步只解玩家覆繪背景阻塞，不擴張研究整個driver。

## R999：現行長劍選取的NUL見證

2026-10-10；入口[R999](../goal/R999-normal-world-text.md)。新正常冷LOAD、人物頁、選取長劍及原格放回後正常停止錄製，續行一秒再關窗0。正式程式為5d6fe6f3…f013；新完整輸入SHA-256 `c9c62d61f6602e482f9a744f7943fab892a2174ba3269e3fe37c862cc64eabb5`，3,054事件、1,514發布、末步331,800,000。三張實際GUI分別為中文、選取後LONGSWORD、放回中文，選取時四個操作按鈕均中文。這是當前可重現缺口，未修補。

已審私人原API觀察程式SHA-256 `e66055aacc56fc5a33e0461c6c27d06a12c7bec47e549d9c336cce9221b76c33`，97來源摘要 `66e8a4120ee15cfced14317b351b2e426d422999fa05ccc62097e5fc6788c360`。Go1.24.13、Python3.11、既有GUI image revision 20261008-gui-r1。啟用與停用探針用同一程式和新完整輸入，兩趟退碼0；1,514對發布的Step、Ticks、RAM、Indexed、Palette、Registers摘要全同，原版兩存檔的大小、雜湊、mode、mtime皆保持。這是私人原API診斷，沒有跨程式正式錄放相容或新增同狀態中文驗收聲明。

位址基準為dosgolem runtime CS:IP／DS:offset。四次精確LONGSWORD皆為`0110:3895`入口、`0110:ba40`返回，font0、x18／y114、DS1106。入口／返回步如下；每次SS相同且SP+2，各有114個離屏變化，完整bbox為(18,114,53,5)，新色皆105；Indexed差異0。

| 見證 | 入口Step | 返回Step | 原墨點底色 |
|---|---:|---:|---|
| 1 | 298318167 | 298321850 | 八種紋理色 |
| 2 | 309629594 | 309633235 | 全114點為255 |
| 3 | 321202787 | 321206428 | 與見證1相同 |
| 4 | 321768197 | 321771838 | 與見證1相同 |

以上NUL來源、返回及差分為confirmed。它沒有量到服務24的現行目的座標、參數或背景恢復，不把離屏返回當最終上屏證據。r717至r730舊工具及收據在本機缺失，歷史紀錄保留；本輪沒有聲稱重跑舊工具。

完整原像素、輸入及存檔只留本機`workplace/r999`。安全入口為`short-report.json`、`item-api-comparison.json`、`probe-source-review.json`及`svc24-plan.json`；本輪長探索逾時124的partial另留，未補footer、未作成功重播。

### R999 服務24的當次畫格

窄契約及修正版來源分別獨立0／0／0採認後，才用同新wire執行服務觀察。私人程式SHA-256 `8e4b8f740c8a90cd7fbc50cdfc014ef6ebc42745ba1e2575034783914943e880`，98來源摘要 `2066bee92de545f1156058766a982581da729b643c9b73a5dad4bbcfbba4a0b3`。保存工具從正式來源乾淨重建的來源收據與程式逐bytes相同，入口見[R999工具](../goal/R999-normal-world-text.md#工具入口)。初版配對邊界被審查擋下，未跑原版；修正版先丟棄重入再檢取樣上限，未配對返回有界記因。

啟用／停用兩趟皆退碼0、無INVALID，1,514對原版六類摘要全同，兩存檔metadata保持。兩件實際觀察均由精確長劍NUL觸發，落在各自固定1M步窗口內；入口`0110:40a7`、返回觀察點`0110:40b0`。下表是原版Indexed的實際差分bbox，並非解出的服務矩形參數。兩件offscreen差異皆0。

| 見證 | NUL身份 | 入口Step | 返回觀察Step | Indexed bbox | 改變點 |
|---|---:|---:|---:|---|---:|
| 1 | 2 | 309634357 | 309650945 | (0,110,69,16) | 186 |
| 2 | 3 | 321231086 | 321242924 | (27,116,69,16) | 204 |

入口12 bytes為`55 8b ec fc 2e ff 1e e7 2d 8b e5 5d`，返回觀察點為`8b e5 5d c3 55 8b ec fc 1e 06 56 57`。兩件入口SP均feb2、返回觀察SP為feb0；此觀察點尚在wrapper內，不套near返回SP+2。入口stack首word分別5cb9、5b80，只記數值，不推全部引數或第二caller語意。工具、SDK revision及位址空間沿本節與本篇輸入表。

獨立像素程式在第一件after找到NUL2的全部114個墨點平移至(16,111)，其中96點在當次改變、18點之前已同色；第二件before找到NUL3的114點平移至(43,117)，全部在當次改變，after沒有相同完整字形。初次分析要求每個墨點都必須改值，導致漏掉同色點；已改為比較實際色號，再另計changed，未改原版或觀察收據。這是資料對齊證據，不以形狀搜尋掛正式中文。

當次before／after已取得，但第一件before含既有英文與游標，不能直接當整行清底。bitmap來源關係、參數與完整游標生命週期仍未證實，正式還原及中文路由維持阻擋；沒有用舊segment／座標猜補。安全結果為`item-service-r2-comparison.json`及`service-data-inspection.json`，完整ROI只留本機。

最後獨立結果複核0／0／0，1,514發布、兩件差分及全部114墨點平移核同，完整讀回兩報告後採認。JSON SHA-256 `51d643cdbef85c94e34d10425544ee4bd14d71d8f2f3a5b9f43667d1b69371ea`；只採當次觀察與原版保持，不放行來源語意或正式中文修補。

## R1000：選取名稱是自訂游標圖片

2026-10-10；入口[R1000](../goal/R1000-held-item-source.md)。工具IDA Pro9.4／Python3.12，既有image `ida-pro-9.4-idapython:py312-v1`，ID `4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab`。新最小探針確認139函式與DS.EXE固定雜湊；所有查詢用一次性副本，原庫及原版唯讀。

MZ EA140a7的wrapper bytes與R999實跑相同；EA15b60／15c90／15cb6為零。依[R881載入基準](016-party-position.md#r881-map-loader-route)，13份OVR的code目的偏移皆5055h；兩個已量call位置唯VPCVW.OVR命中。VPCVW.OVR SHA-256 `8dc7006599844bfd33ae60e11b59f7cedf47d649d84669fdf2546ed15779224d`，其14-byte檔頭後38966-byte code SHA-256 `f8e2ed6c99860106a09e5ac48e6f1c3dbb964066a9c1ba949971872e0563425e`，與R992既存code逐bytes同。NUL入口3895仍是resident；兩caller位於同CS的overlay載入區。

本節OVR定位同列區分runtime候選IP、raw code offset及檔案offset。只有已量wrapper對應為confirmed；VPCVW完整runtime對應待短code簽章確認，維持強推論。IDA原有導航名稱與部分operand文字沿用，不用它們代替raw bytes或直接計算的16-bit near target。

| runtime候選IP | raw code offset | OVR檔案offset | 原始bytes與資料流 |
|---|---|---|---|
| 5b7d | 0b28 | 0b36 | `e8 27 e5`呼叫40a7；offset、segment、兩個尺寸word均0，之後清0ch stack |
| 5c18 | 0bc3 | 0bd1 | `e8 54 dd`呼叫396f；參數為來源矩形兩角，右／下各加75／12，回傳DX:AX保存在[bp-10h]／[bp-12h] |
| 5c24至5c38 | 0bcf至0be3 | 0bdd至0bf1 | `b8 dc 03`後push988；目的[0e20]:[0e1e]、來源DX:([bp-12h]+12)，`e8 d8 a3`呼叫0013 |
| 5cb6 | 0c61 | 0c6f | `e8 ee e3`呼叫40a7；六word為pointer、76、13及兩個位置，之後清0ch stack |

0013轉呼叫resident4317。原始`8b 4e 0c`、`c4 7e 08`、`c5 76 04`分別取長度、目的ES:DI、來源DS:SI，後續按word及尾byte複製；988非零分支的複製方向為confirmed靜態控制流。396f是driver wrapper，來源矩形及12-byte頭與實際像素內容仍須runtime核對，沒有只憑988＝76×13稱背景已解。

VGA.DRV固定雜湊沿本篇。raw offset0063的`e8 63 13 cb`對應runtime0163→14c9；20fe是後面的內部寫入點。14c9在driver CS下把入參pointer存1e64／1e66，尺寸存1e5c／1e5e；限制width小於81、面積小於1025。零pointer會還原預設尺寸；位置經INT33服務4／3，再呼叫1f1a。這是設定自訂滑鼠游標的confirmed靜態路徑，當前pointer仍待實測。

1f1a先把新背景從A000複製至CS:[15b6]，處理前後游標重疊、以255為圖片透明色，再交換兩個背景buffer並發布當前欄位。當前欄位為15ac／15ae的位置、15b4的背景offset、15b8／15ba的裁切寬高；1ed3是繪製中旗標。已找到靜態來源與生命週期欄位，尚未取得當前背景oracle。私人P1至P6契約已獨立0／0／0採READY，取樣後才決定正式覆繪接線。

安全證據索引為本機`workplace/r1000/ida-existing.json`、`overlay-callers.json`、`overlay-input.json`、`ida-overlay.json`、`ida-source-capture.json`、`ida-resident.json`、`ida-memory-copy.json`、`ida-driver-entry.json`及`ida-cursor-background.json`。caller獨立0／0／0已採認，copy與driver補證另列，不外推該審查範圍。沒有修改正式Go或宣稱長劍中文修補完成。

### R1000 動態來源與正式契約

本節接續上面的靜態階段。R2私人取樣核當前背景為driver CS內的near offset15b4；29組固定位置988像素共28652點，以背景加圖片非255像素重建真正Indexed，差異0。正常新movement wire再量右下裁切及移動：控制／觀察各1651發布，完整六欄摘要逐列相同，兩保存metadata保持。124組ROI共101656像素重建差異0，含長劍100組77944點及皮胸甲24組23712點。長劍位置為0,110、100,145、200,100及300,198；最後一處裁切為20×2。長劍重設request及complete均量到。

皮胸甲與長劍使用相同圖片pointer，私人弱註冊誤沿普通長劍NUL5。這24組只證背景重建，不能證名稱。長劍圖片SHA-256 `98e96e235de0701c61b2629c6a78bc6ba1f49fced39efad81e094d8ca498e863`；皮胸甲為`5922e142ae731e2e041ea1f0aa16cb869e345bd60b00f954f413b4479ffa532c`。全60×5名稱帶核定：長劍第一帶114個105墨點，第二帶空；皮胸甲兩帶分別89／85個105墨點，沒有圖示色混入。普通NUL5的before含265個非255舊像素，正式來源必須拒絕。

396f抓圖入口0153的`e8 24 07 cb`轉087a，來源為driver CS:01e5指向的離屏。先寫六word共12-byte頭，再從y×320+x逐列複製inclusive寬高，回DX:AX；VPCVW後續複製pointer+12的988bytes。本機`ida-capture-handler.json`保留bytes及位址，正式037以當次完整源圖相等驗證，不靠墨點shape辨識。

獨立背景結果`movement-background-review.json` SHA-256 `1557b4884684433f1bfe338310d3501ea5bda02bad0acf14aca18fc6c5b1baad`採認上述限定數字；新wire正常Stop、續玩及關窗0。正式[037](../spec/037-held-item-cursor.md)第二輪0／0／0採READY，補齊來源至capture／definition、整圖相等及同位址換圖反例。正式R2來源審查JSON `formal-source-review-r2.json` SHA-256 `f2ca9fcabbd55430de7f0cc27e7ddeafaa07884f9fa732d081e8e2226878b8d9`通過；兩字型兩倍率實際畫面及正常GUI仍在驗收。

### R1000 單行來源缺口

正式R2實跑未通過：Unifont2／3倍各1651原版發布及兩保存metadata保持，12圖原圖9960點、RGB64740點及Draw／HD全圖差異0。但長劍四位置均無名稱來源、持續英文；皮胸甲來源與原生205墨點通過，放回後兩reset均清label／view。獨立`formal-pixel-unifont-r2.json` SHA-256 `8e3461617726cb4cbe72daba71d31e667a2bca89bbc62e38fc9909013f9ab57e`列12阻擋，不把英文當短暫過渡。正式新GUI重現相同行為，成功Stop及正常關窗只證收尾完成，不解除文字失敗。

窄靜態補證：VPCVW.code raw69e8的`e8 55 7e`正常返回後，經69eb的`83 c4 04`、y+6及第二pointer計算，無分支到6a0b的`e8 32 7e`再次印字。沿既存5055h載入基準，兩call為0110:ba3d／ba60，返回ba40／ba63。第二pointer在raw6a02為`8b 46 fe 8d 4e d7 03 c1 50`，即[bp-2]+lea[bp-29h]。工具為人工16-bit解碼與GNU objdump2.40交叉核對，bytes／位址及輸入雜湊見`single-line-static-review.json`。已證原始code內必經第二NUL，尚未由此片段證實其字串為空；新私人內部觀察只記既有讀取結果，不增加Bytes或讀DB。

私人內部觀察後續confirmed：同movement wire首次長劍scratch入口Step305629523、返回305633164，保留完整LONGSWORD及矩形(2,113,76,13)。305633200同CS0110、DS／SS1106、SPfe66，讀取器首byte為NUL；位於既有ba63路由，距第一返回恰36步。此空字串被當非法而清掉第一候選，305633680的capture沒有來源。1651原版摘要與原控制trace逐列相同，保存metadata保持。證據入口`single-empty-evidence.json`索引原state檔與SHA；未量第二pointer值或返回全域像素，不從此摘要冒稱。正式037增補要求第二空部件同原配對守門、精確返回及完整988像素保持後才保留第一部件，來源仍須經capture與整圖比對。

### R1000 正式限定驗收

E1至E4增補及R3來源分別獨立0／0／0後，正式單行與雙行來源鏈已實跑。正式96來源SHA `7a38f5709789f8df8647beb1ff26626c0dae23f67ad3fe1a729e0f4a9bb3e230`，binary SHA `c7da1a38a0a690de1529d0ebadf13c3378e1b51bf3063a522a707962c88fe224`；接入正式四檔後完整cmd/wiz7測試2.630秒通過，乾淨建置與已驗R3逐位元相同。

首次四組23／24圖通過；倚天2倍首圖Draw／HD失配，私人擷取工具存在畫面與請求分兩鎖讀取的競爭窗口。修正僅限私人工具、獨立審查通過後，四組R3b全24圖及8677項通過，沒有排除首幀。原圖19920像素、RGB129480像素及Draw／HD全圖差異0，原字模直接從固定倚天／Unifont獨立解析。每組1651原版六欄摘要同關閉覆繪控制，共6604比較；兩保存metadata保持。兩件來源proof、長劍四位置及右下裁切、皮胸甲兩位置及兩次reset皆通過，沒有英文暫態。私人98來源SHA `d02a96c4fd390ec2b00125e13a032bdf306287d4deaad266dd6478bac62c4703`、binary SHA `28b3db43316bd2834c9d7574219245b2c086a3994e2a318e4878a6a707629d7b`；保存工具重建逐位元相同。

真正Unifont3倍冷LOAD後正常選取／移動／放回兩件皆中文；手持長劍移過第三列再移開，8892個RGBA像素差異0，覆蓋時差370點。新正式錄製3408事件／1686發布／末377.8M正常Stop、續玩及關窗0；同96來源binary嚴格重播0，保存保持。此保留比較不另證第三列原字模正確性。

最終獨立像素JSON `formal-pixel-r3b.json` SHA `39dd1fd414550512ffa5ea6ce350657c8392c5c3771e83daa3489474e28f1a33`，GUI保留JSON SHA `d25937489fef96c7a7d9e8b387d0d5d60ea5d911a1b8f618ced6d69834309ac6`。入口為037、R1000 goal及其中保存工具；原版圖片、字型、存檔與完整wire只留workplace。僅本次長劍／皮胸甲路線限定CONFORMED，其他物品／slot／caller／平台及防禦分支未外推。

最後獨立0／0／0完整採認，`final-review.json` SHA `7bac991457d7997dcdc60e2bc54d6eca3d4e205a691aa731e0bdc8a32ccf95b4`。獨立核24對實際Draw／HD共9984000像素、19920原圖點及259份本機manifest全相符；正式四檔、正式重建及私人保存工具重建皆逐位元等於已審版本。未重新讀外部原字型，原生中文ROI沿獨立字型比對器及其固定收據採認，不外推審查範圍。

## R730：三組複合名稱的獨立正常驗收

封存`workplace/r730/report-final.json` SHA-256
`3935e630de3a93e5b0034971d76c66e2992045afe98a7aad847124b0924ff3f5`；
5,031件輸出索引`outputSHA-final.json` SHA-256
`80f93aaa5a2fd2095a133092db4524f0b06080755fddd929eb90114f0d99b931`。
逐條C1–C5獨立稽核為阻擋0、應改0；原版來源與工具沿本篇固定版本。

兩原始字型獨立解析，皮胸甲／毛皮護腿／小圓盾在2×、3×皆完整上屏；
每個字型與倍率的360M／400M各三唯一完整鍵，歷史事件3／6。
470項獨立全圖、原字形、四階段背景與遮罩核對失敗0；原版A000及色盤逐位元相同，
記憶體沿既有FNV32摘要，未匯出完整RAM。四張400M實際Draw與HD全圖相同。
2×既有通用「盾」移除各站造成25個域外差異，全部落在舊15×15／16×16字格，
三新域加舊字格以外差異0；3×三新域外差異0。

早輪原版唯讀探針證實三對第二返回302686519／302692412／302704140；
DS1106、入口SPfefa、返回fefc，兩行指標連續與36步關係相同。
360M原版三域已等於早輪after2。舊r717從320M才開始追蹤，不能證明此前未印；
初次失敗與乾淨重跑都保留，驗收文字訂正未改配對契約。

58個負向子案例、六群組全部通過，但原版玩家路線未觸發這些故障，012維持READY。
五項代表故事／旁白／F1比較通過；未重跑全732圖、全地圖、音訊或發行包。
正式程式封存SHA：`main.go` d1fde8c779596513fb57ef15e26f0dca2391099d96ae6feb4d827f6975e75b9e，
`compound.go` e84ada085248ecf511a38ec114ddadea9011ab001bfc83c97df775ca7bcb3408，
`compound_test.go` 216b33da40612d0c93f8a755cda226acf19160bbb2cfaa633cdfeb319dd5bc67。
10,137項輸出UID/GID1000:1000、root-owned與誤建`.md`目錄0、自有容器0；
未改原版、正式譯文或字元清單，LONGSWORD背景仍未知。

<a id="r964-battle-use-items"></a>

## R964：戰鬥使用面板的兩組原版印字

2026-10-05；入口：[R964](../goal/R964-battle-use-item-evidence.md)、[012私人探針契約](../spec/012-compound-item-output.md#r964-使用面板私人唯讀探針-ready)。原版DS.EXE、SDK及工具映像沿本篇固定版本；私人binary SHA `959c029d8d8691980aad75eff430b07e4a6b1431a3fe8aca74a4355d2f06c431`。來源色盤為[256][3]uint8，首輪審查在建置前發現錯誤切片型別；修成逐RGB取得獨立768byte，修前來源及索引保留。私人來源0／0／0完整採認READY後才建置與運行。

沿R733合法checkpoint、82點固定前綴及原版既有seed。主代理逐張檢視真正paused Draw，再送正常Update／owner／CursorWait Click：61,177使用、98,187取消。三組F1 ACK／paused／resume Step相同，六份PNG及來源／副本核對通過，正常關窗exit0。單側私人暖節流no_pcm；原版RAM、離屏及A000均未寫入。不讀MSG.DBS、鄰文、手冊或答案，只取已存在8／64byte快取的精確四部件。部件含NUL長度8／8／4／8，ASCII及NUL位置符合且無控制碼；每筆來源stage為part8與generic64。

下表CPU位址皆為dosgolem實際resident CS:IP／DS:offset。入口同為0110:3895，font0、x177、DS1106、SPfd8c。真正返回須通過來源SP+2守門，收據未另存ReturnSP欄，不能稱獨立回傳SP取樣；正式守門只採相等及+2關係。

| 部件 | 呼叫者返回位置 | 第一y | offset | 入口Step | 真正返回Step | 原版改變像素 |
|---|---|---:|---|---:|---:|---:|
| LEATHER | 0110:50eb | 165 | fd98 | 5458718948 | 5458721786 | 89 |
| CUIRASS | 0110:510e | 171 | fda0 | 5458721822 | 5458724648 | 85 |
| FUR | 0110:50eb | 179 | fd98 | 5458726053 | 5458727276 | 36 |
| LEGGING | 0110:510e | 185 | fd9c | 5458727312 | 5458730147 | 88 |

兩對第一真正返回至第二入口皆恰36步，第二指標為第一ASCII長度加NUL之後，無16-bit繞回。兩行所有改變只在各自41×5行，新色全為180。完整域分別(177,165,41,11)、(177,179,41,11)，每域451bytes，兩行間隔一列；兩份mask非空且不交疊，兩行以外變化0。after1==before2；各域after2等於stage01 paused Indexed的相應裁切。完整110×40研究ROI因後續其他印字／動畫不同，沒有稱其after2等於paused全ROI。

| 451byte原版欄位 | before0 SHA-256 | after1及before2 SHA-256 | after2 SHA-256 |
|---|---|---|---|
| LEATHER/CUIRASS | `2bf1ec8417243ff307070e7b2ecfef996cc2e7cf31f7323cc15142b9017aae46` | `ea8e6add9d06e629fc0ac1ab08b5275c9545e656e028fa91743e813607c89cf9` | `cacb27652ef1e2e9bd2160d702f3d3bd7026912c40dedb35e7679b9c71bbbd4c` |
| FUR/LEGGING | `2bf1ec8417243ff307070e7b2ecfef996cc2e7cf31f7323cc15142b9017aae46` | `05613afc88e49aa6bff8d983c7f90f6f16ef7cda5765c8bfb45570ebdd6f4fda` | `70ad5914719765a422a02fbd38f5faa7e32361b613af3e87b32c09469556b4bb` |

根因限定為既有角色頁012的ba40／ba63、x18、前景105不接受這兩個使用面板tuple。現有「皮胸甲／毛皮護腿」譯文與當次原生字模皆存在。R964畫面仍是兩列英文，沒有中文化或新A/B完成聲明；只足夠建立R965正式DRAFT。其餘使用道具及其他caller未知。

| 本機workplace/r964收據 | SHA-256 |
|---|---|
| gui/use-observe.json | `1248186148346662690490435007e8f2fdb11d17df23fefae008aa55336d161a` |
| use-check.json | `47fa0254e2d706a98f204f9a846db0241d6d7dbf9afc7eb11fa10a47d47eef6d` |
| gui-check.json | `c80374e97a779b1793aeddee87b589d2e1a3e539c295dc37647845777d505088` |
| invariants.json | `eee3fcec2de7121eb81756fabae65db2ccb7335f2f8cb769c9dd18c923cce8a3` |
| gui/run-result.json | `be0cf571cf749b124517874e9285de8ea44110d44e52a7b69d1b815480696da7` |
| gui/operations.json | `dfff79557e9c9e9a17a1752adcf4e5974f85a99cee43164fbeb3c41a55cced92` |

原版248檔、SDK1132、正式40Go、六資源、R962166及R963125封存盲SHA保持，候選45檔。真正GUI只量Unifont2倍；初末七控制共28原字模、512,000RGB獨立核對0差異，四命中且溢出0、停止旗標false，末Step5473900000。色號180來自實際像素變化，沒有另量前景ABI欄。原版衍生像素及完整收據只留本機。獨立結果0／0／0於FINAL完整回讀兩報告及SHA後採認；JSON SHA `1ea273e1a07104131b9e61f7409b0f18b661f86470c269121196fb5544c849ea`，Markdown SHA `ed89460c445fff2481c2bb9c5e4ebaf991bce352faa343c923ccd27935b61c6d`。審查另解碼原版ROI及PNG／字模，沒有執行既有checker。

115技術產物一次封存，sealed.json SHA `961811b373ca7573837a085f099b9cc1e2a5ca5feb1bf96238ac1cb918873c90`，排除gocache及後續GitHub收據。重跑入口在本機workplace/r964的prepare.py、build.py、gui.py、check_use_observer.py、check_gui.py及verify_inputs.py；依012私人READY契約及全新輸出，不放寬白名單或數字停止。原版同時唯讀掛/original／w，建置300秒／GUI1800秒、UID/GID1000:1000、2CPU／2GiB／96pids及網路none。限定根因已量到，正式修補由[R965](../goal/R965-battle-use-item-fix.md)接續。

<a id="r965-battle-use-items"></a>

## R965：使用面板兩個完整中文名稱

2026-10-05，入口：[本輪目標](../goal/R965-battle-use-item-fix.md)、[012正式READY增補](../spec/012-compound-item-output.md#r965-使用面板正式增補-ready)。原版、SDK、映像與工具版本沿R964。獨立規格0／0／0於FINAL完整回讀採認後，只改正式compound.go及compound_test.go；main.go SHA `e6f50fd640d0a5d90d53b83fc76ad498d31d8579bd92feb35b228d8c974035a9`保持。

封閉route綁定名稱、兩返回位置、x／y與前景。三個角色頁tuple仍用ba40／ba63、x18及105；兩個使用面板tuple才用50eb／510e、x177、y165／179及180。背景41×11、兩行mask、DS／SP／指標、36步、原生字模及整件撤除沿原契約；未知使用物品不放行，不另讀字串，不改原版返回或記憶體。譯名與雙字型資源保持。完整cmd/wiz7測試通過，含既有角色頁及兩名稱各五個正常／錯返回／錯前景／錯座標／背景改寫案例與未知SHIELD；負向測試仍是隔離證據。

私人驗收副本沿R962，正式差異只有compound.go；另加固定重播helper及私人audit完成後的環境變數分支，正常GUI接線保持。私人binary SHA `c3b818a34520f6da7bd507373e361923d5fc288e4c2c1db0d7a114a95a22ea71`。SDK兩份既有私人唯讀overlay沿R962；沒有改SDK或在正式程式插入私人入口。

英文控制與兩字型×兩倍率共五個新程序，沿R733固定起點、82點前綴及R964實際50筆Run／mouse／正常Click。新操作不添NoCursorWait、不改seed、不重擲。每個中文配置的50點記憶體SHA、CPU暫存器、滑鼠、鍵盤及IRQ1摘要與英文控制一致，共200比較；六個指定時間點的完整indexed／palette亦相同，英文控制與R964六個實際GUI幀的原版摘要相同。這是固定實際輸入的輔助重播，不當新的正常GUI A/B。

四個中文配置各六張HD，共24張；在5462300000／5464700000各有「皮胸甲」「毛皮護腿」，取消後5471100000／5473300000兩件均撤除。獨立原字型及R964四stage／mask重建56字模、46,904欄位像素相符，trace背景與mask SHA、返回／入口Step皆同R964。新域及既有批准覆繪域之外像素差異0。其他文字域仍使用既有事件metadata作域界，不把它們稱為本輪全量原字型驗收。

真正Unifont2倍暖GUI另沿owner、正式Update、正常CursorWait Click及Draw。主代理逐張view_image後送61,177使用、98,187取消，實際看見兩列完整中文，取消後兩列及面板撤除、七中文戰鬥控制重現，正常WM_DELETE exit0。三組paused／F1 ACK／resume Step相同；兩張使用面板Draw以原Unifont與R964背景獨立核14字模、7,216像素相符。初末七控制28原字模、512,000RGB保持，六PNG與來源／副本核同。

| 取樣 | live Draw Step | paused底層Draw Step |
|---|---:|---:|
| 戰鬥控制 | 5455700000 | 5459100000 |
| 中文使用面板 | 5465700000 | 5469100000 |
| 取消返回 | 5475500000 | 5477900000 |

GUI共74操作、兩次正常Click，末5478700000；Manual0、Lexicon102、NumericStop=false及WarmPacerFault空，watcher／owner／game／Xvfb終止。私人節流no_pcm，未新增音訊、冷START、正式錄放或完整通關證據。兩字型各六資源、SDK1132、原版248、R964115封存及本輪凍結來源40／候選45前後盲SHA保持，只有上述預期正式兩Go差異，擁有權異常0。

| 本機workplace/r965收據 | SHA-256 |
|---|---|
| build-result.json | `d1b61c1d245d62c4747cc9e1d67050efd38c947c322cbf1d87e75d91083a291c` |
| replay-check.json | `24b15f71d82c1b4ecfcb93d6b31caa41ed180a87514ff174c14b8873c59dd816` |
| gui-check.json | `0de3ebc73fe3b4378c0d95d8166d441abe88f155354283ae893d0aa6a7b75224` |
| gui-items-check.json | `3d2c472e7898b9a21260d36ed04bf7ba21f669511ac80b8d3ad763a78d085843` |
| invariants.json | `4d6147927bf61b51ca2c4d756746dbaada61b42e0f472d01cd94f3f25f4887cf` |
| gui/run-result.json | `bf29714b656ace1ef7035fe3221a890cd68936d8392f9c61ee838f544d9d8589` |
| gui/operations.json | `9d55bd6d8763f86c7ebad8aed3cb6b2cb786a9548f17330f4678dbd127285658` |

原字型輸入：倚天stdfont.15 SHA `39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701`，本機唯讀掛/eten；Unifont固定17.0.05 gzip SHA沿R964。正式compound.go SHA `ccc0a8f3250466a299d5e3b4e585540db42de3e1b986db1af9aac18707007522`、compound_test.go SHA `1c647507e58f5b8f33c85e0c704d611fb0654c34fcebed23a33dec14ea9b93e0`。限定結果獨立0／0／0於FINAL完整回讀兩報告及SHA採認，JSON SHA `2e3a64a5fb2b47b5fd057c8e6500ba115b4ea2bbb82f23f582d5759c414d2905`、Markdown SHA `addb23ddd964b39ab2e2e8139e9a527be2362412008a58f8acee51dea05ac0e0`。審查獨立解析PNG、兩原字型及原版背景，不執行既有checker；倚天索引及metadata誤讀修正屬工具問題，未重跑產品。310技術產物一次封存，sealed.json SHA `4af40246b5a1ff6b1a319b8db6ff149c1b73cc5e9a2857620e426124741d1bc0`，排除gocache及後續GitHub收據。012維持READY，未知物品、其他正常GUI配置與全部交付門檻保持。

重跑入口為本機workplace/r965的prepare.py、build_test.py、run_replay.py、gui.py、check_replay.py、check_gui.py及verify_inputs.py。準備乾淨輸出、固定R964操作收據與SDK私人overlay；原版同時唯讀掛/original及/w，字型唯讀掛/eten。沿本輪有界離線Docker及UID/GID1000:1000，不放寬caller或數字停止守門；完整原版收據、字型與圖片只留本機。
