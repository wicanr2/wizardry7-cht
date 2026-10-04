# 按鈕與訊息的中文覆蓋驗收

## R955 頭像頁文字與完成返回

2026-10-04至10-05，入口：[本輪目標](../goal/R955-character-portrait-ui.md)。沿R952已建正式Linux前端，binary SHA `7363fc5be7c48d10eb272d5bd5900045996bb490b0cb08fd19bb1a93e9d948b0`。38份列出cmd來源保持，不稱新建置或完整編譯閉集合；工具鏈沿固定Docker映像、Go1.24.13／Python3.11.2，原版`DS.EXE` SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。repo、SDK、原版與字型唯讀，只有r955可寫；沒有MSG、防拷、內建autoplay、原版記憶體或角色資料注入。

- 真正3倍Unifont GUI逐張檢視實際Draw後送六次點擊：角色選單160,47、頭像236,193、第一列175,152、選人確定100,170、完成160,155、離開287,193。未選新頭像；畫面提供完成控制，不把它稱為取消。七個停止畫格、正常返回主選單與關窗exit0，耗時0.100491952秒。
- 原版248檔SHA保持，擷取後七份Scratch雜湊及退出後PCFILE.DBS 9368 bytes保持；最後另核八個Scratch與原版逐位元相同。Scratch讀取在擷取後，不冒稱恰在固定Step。兩則文字的10個原生Unifont字模、兩個安全格及HD／Draw局部19,968像素通過；14張PNG的CRC、解碼、尺寸與SHA相符。
- 另有雙字型、兩倍率及英文控制六趟固定輸入輔助，全退出0。四組中文各七停止點，共28對Step、1MiB記憶體FNV、A000、Pal、Mouse摘要與完整indexed／palette bytes相同，不稱完整CPU或新人工GUI對拍。40個原字模實例、八個安全格與局部HD／Draw 64,480像素通過；84張PNG的CRC、解碼、尺寸與SHA相符。
- 兩則文字的原版1524個改動像素另與英文indexed核對，純色背景由英文控制的相鄰空列核實。原倚天按Big5分區與大端15列解字，十字皆為漢字，不讀符號區；原Unifont按16列解字。倚天392820 bytes、SHA `39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701`，Unifont原壓縮檔SHA `7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574`。

| 原版英文鍵 | 原版座標 | 本輪限定結果 |
|---|---|---|
| SELECT PORTRAIT | 116,138 | 「選擇頭像」雙字型兩倍率字格與局部像素通過；未選新頭像 |
| CLICK HERE WHEN DONE | 101,153 | 「完成後點這裡」雙字型兩倍率字格與局部像素通過；160,155完成返回 |

首次正常GUI碰到外層600秒期限exit124，已有頭像頁完成返回角色選單的畫格，但最後離開命令送達未確認，沒有成功退出收據。保留片段，以相同binary、時基、七停止點、六次點擊及600秒期限乾淨重跑；只改工具輸出路徑與ready雜湊日誌。第二次成功，前六張Draw與首次逐位元相同。首次Docker inspect誤用不存在的Entrypoint欄位，改核image ID；均屬工具問題，未調遊戲參數或列產品缺陷。

正常收據位於`workplace/r955/normal-unifont-3-attempt2/captures.json`，SHA `1abeace7d4c1426e542e16141717a2e35d05ecb7db0e911a30840fc582133dec`。私人獨立檢查入口為`check_portrait_gui.py`／`check_aux.py`，SHA依序`1355a28df7121282a27554532d9ea24736f107e8a68ebc863ce28082293da8d0`／`98284ad1971ad31f46cec2f256912f6a1dcbbd002d5a70719a56243f9a58ec14`。沿固定離線非root映像，repo與字型唯讀、r955掛`/out:rw`；後者另需倚天目錄唯讀掛`/eten:ro`及六份aux收據。報告獨占建立，不重跑原版或渲染器。擁有權異常0，正式程式及譯文未改；只有頭像頁兩則文字與保持原頭像返回，不消除已知2倍人物／名單或其他圖片控制缺口。十一組、TSV、錄放、冷存讀、音訊與封包門檻保持，R954／R944／R951待選不代答。

正常報告`portrait-gui-check.json` SHA `0bb142d9465f66e39bfa1f7170b003d4b32f2fc3213c097f46e048f8b399d26e`，輔助報告`aux-check.json` SHA `e8d73d715205494a63acc3be20bb304812595d42a3cdca5d6e2a111b2576d770`；`integrity.json` SHA `e497df462dd43cbce27eabc5bace611099ac62f0c07c0715ffc2c85347aa99ea`回查19份選定來源／收據及八個Scratch完整bytes。工作根衛生異常0、自有容器0。SDK只有既有修改狀態，未宣稱本輪SDK1132份新雜湊核對。

R955獨立唯讀審查阻擋0／應改0／建議0，主代理於FINAL後完整核讀兩報告並核SHA。`workplace/r955/review-results.json`／`.md` SHA依序為`2ea5c896c905c0c79d1c33e258ecf5175ed8d41d10dcc381d9ef25bd7b82213c`／`f2cb18a06319eedf1cf1a5938beb70cb0cf549fa0aa8ad41d5549c0dde4db72e`。審查另獨立抽核兩套原字型、四個字模、六張PNG及2948個局部RGB像素；既有檢查器的19,968與64,480像素範圍分列，不稱審查重跑整批。只採認兩則頭像文字及保持原頭像完成返回，不新增READY或完整介面完成聲明。

現行R943已完成片頭六張「呈獻」的正式覆繪及限定正常GUI驗收，入口見本檔「R943 正式呈獻覆繪驗收」。其他圖片敘述、按鈕與完整交付仍待。

狀態：主選單、角色選單、新增角色對話框及3倍人物檢視限定像素驗收通過，R941結果複核通過，其餘介面重查中；2026-10-04；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
入口：[R915](../goal/R915-ui-localization-audit.md)、[R921](../goal/R921-character-dialog-ui.md)、[R928](../goal/R928-character-roster-text.md)、[R932](../goal/R932-character-menu-audit.md)、[R935](../goal/R935-character-review-audit.md)、[R942](../goal/R942-intro-picture-text.md)、[目前狀態](../../CONTEXT.md)。

使用者要求把圖片內的英文按鈕納入完整中文化。本文保存目前正式來源及實際畫面的覆蓋矩陣；譯文存在、原版畫格及舊測試完成數量不直接當作本次介面已驗證。

## R954 動態背景與等級欄修正

2026-10-04，使用者已選定R953重排，沿[本輪目標](../goal/R954-profile-reflow.md)補背景守門。仍是唯讀原版實驗及離線候選，正式程式／SDK未改；024既有READY範圍保持。原型圖片依原版素材權利只留本機，定案表與雜湊進Git。

沿R953原版`DS.EXE`及固定Docker／Go1.24.13／Python3.11.2工具鏈，另建探針，只補取性別入口64000-byte背景及已量畫面的呼叫metadata，不讀MSG、鄰文、詞庫或答案，不寫原版。男女各13事件及7停止點，observation與R953完整相同，indexed／palette亦沿既有14對控制。兩趟退出0採工具實際終止回報。原版248檔SHA與R953相同，兩個Scratch的PCFILE.DBS保持；不稱新正常GUI或完整CPU驗收。

下表地址為原版runtime `0110`段，font均0；訊息原文仍是覆繪比對鍵，表中返回位址只補辨識結構。

| 標籤 | 原版座標 | near返回位址 |
|---|---|---|
| 種族 | 102,6 | 6c7a |
| 職業 | 90,12 | 6cc0 |
| 等級 | 90,18 | 6cf8 |
| 階級／階級名稱 | 228,6／288,6 | 6e13／6e6b |
| 經驗／擊殺 | 228,12／228,18 | 6ea3／6ef9 |
| 生命／體力 | 6,93／6,99 | 6f5d／7007 |
| 八屬性 | 115,35至77、每6一列 | 6a43 |

生命／體力的下一輸出邊界`0110:3895`、返回位址032a、225,138實測font3，解釋R953 V1的整體font0篩選漏掉該邊界。此處仍是時間順序邊界，不是匹配函式返回。八屬性在該邊界之後才輸出，最後來源返回分別為男666672690／女666594273；不能把生命兩列完成當作整頁完成。

以原版色號、R953十三個來源及R937十七個標籤的精確before／after還原筆畫，再比較性別入口與784M完整畫格的候選格。原選版面7970個DOS像素中，扣來源筆畫後，男7029格有8個變更、女6989格有25個變更，皆位於等級字格內的蠟燭圖案。主代理檢視原版相位與四份修正版；新字模與上述變更像素相交的輸出筆畫為男倚天3／Unifont7、女倚天22／Unifont27。這是confirmed的兩停止相位交集，不代表已量全部蠟燭相位。

原選位置因此無法直接通過背景守門。可撤回甲案只將等級文字右移24輸出像素，字格204,51,34,18；等級數值安全格右移12至240,51,36,18，實際數值墨點保持。其餘30個顯示欄位保持。甲案的7900個DOS目標格，扣來源筆畫後男女6959／6919個背景像素變更0；這是兩停止相位證據，正式每次合成仍需守門與撤除。乙保留使用者原選位置，但須明確接受文字覆在蠟燭上，不能默認這項忠實度取捨。

甲雙字型、男女四份重排及四份現行對照八張PNG，獨立2,048,000像素相符；只有四張重排圖各31格交疊／墨點交集0。字模沿R937已核原字型資源，沒有本輪新rawfont oracle。獨立checker的第一份摘要未標新幾何輸入，補入checker、八張PNG及新layout的SHA後另寫v2報告；未覆寫初份，也不以相同報告SHA聲稱新來源已封閉。完整[並列頁](../../workplace/r954/level-revision.html)為私人離線原型，不是正式HD／Draw驗收。

| 本機workplace/r954收據 | SHA-256 |
|---|---|
| `profile_background_probe.go` | `8a1eac669b70e605978bf4f068e58be28facc9a27b289b228674868f07fb22ac` |
| `profile-background-probe` | `50df293b2416e9c23f22c408a85e5bae474baefd8c9b00ac5245c0219326b7cd` |
| `diagnose_background.py` | `d4ee654f2520865ddee6e06e16c52ba4b225f5dadadf50d51842a69bf5358a64` |
| `candle-conflict.json` | `70bf0c70961aa0b854620aed18555238dc7f4ef4683cd1b4bc3cc0ba833df928` |
| `background-check.json` | `dd7266239d117b1fc4ec011d85b709db5c6f897aa8758d41adf9c62414b458a3` |
| `prototype_level_revision.py` | `1b1375a7aafb4940f964fa3609276ca383ce119573037c3341827097e0fcae56` |
| `check_revision.py` | `3856fb5cbae022510aadd9d994e1da96d4407fbdc21f4858e430b089a84ee857` |
| `independent-revision-check-v2.json` | `04b21baa6eb6d9c360dce39ca5eb414b65a52084bfb0468f690155a787a9293c` |
| `prototype-v2/result.json` | `0399fc54ed9422eef466d955b14a51021269e7d36bc5dc41f3c1ebca3f11e5f6` |
| `level-revision.html` | `76377b5f1c3ae716f921789849a2282410a344bbe3bc454b2394cffbf71200d5` |

重跑沿既有固定映像，repo、原版及R953／R937／R943輸入唯讀，新且已存在的輸出目錄掛`/out:rw`。以相同唯讀go.mod/sum、受控SDK及模組快取、可寫Go快取建`profile_background_probe.go`，男女參數與R953相同。最終候選入口依序為`prototype_level_revision.py`、`check_revision.py`、`check_background_v2.py`，另用`diagnose_background.py`量原選位置；須提供新探針男女輸出、原版基線及上述舊收據。輸出既有則拒絕，不重跑所有草稿，也不拿缺輸入的執行當通過。

首次背景checker因真正動畫差異拒絕；診斷器一度把男的8格／3／7套給女，屬先驗計數錯誤，改成各自實測後量到25格／22／27。初份腳本與比較頁說明保留，最終說明已修。正式位置未自動改成甲，待使用者選定最小右移或明示接受覆在動畫上；技術契約未READY、未正式實作，十一組與其餘交付門檻保持。

R954獨立唯讀審查0／0／0，主代理於FINAL後完整回讀JSON與Markdown，只採認限定RE與離線原型。報告位於`workplace/r954/review-background.json`及`.md`，SHA依序為`7da897e752b430f0cc04ce37bf746ec3c2da23574371206a5a1a0f4c80466194`／`aee0b7b124c16b63902426504d6d45772eb834bcf3497d07c54c4261a896fa69`。另獨立重算兩相位背景、八圖完整RGB及最小變更範圍；不批准甲乙取捨、新READY或正式實作。審查工具兩次定位／欄位錯誤已修，不列產品失敗；自有容器0。

## R953 兩倍人物欄位原型

2026-10-04，入口：[本輪目標](../goal/R953-profile-layout-evidence.md)、[排版候選](../spec/024-character-review-output.md#r953-兩倍人物排版候選)。這是唯讀原版測量及離線原型，未接正式前端，也不是新GUI驗收。

原版`DS.EXE`為69,206 bytes，SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。工具為Go1.24.13 linux/amd64、Python3.11.2及既有dosgolem隔離工作樹；Docker沿固定映像`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。下表地址皆為原版runtime段址，不是檔案偏移或IDA線性址。

| 已量事件 | 入口、辨識及結果 |
|---|---|
| 性別 | `0110:3877`，返回`0110:6c42`；font0，90,6→96,6，SP增加2。擴大取88,4,8,8的背景，M／F差異13／11像素且範圍外0。 |
| 十一個數值 | `0110:3895`，返回`0110:032a`；八屬性y35至77、每6一列，兩位數x139、一位數x145；等級126,18，經驗306,12，擊殺306,18。font均0，返回SP增加2，游標前進字數×6，差異只在各自6像素高字格。 |
| 生命／體力兩列 | 已量6,93入口至`0110:3895`、返回位址032a、游標225,138的下一輸出入口，取6,93,78,12的差異。未讀字串正文；男女177／181個變更像素、範圍外0。這是時間順序邊界，SP差為−306，不是匹配函式返回。 |

探針只在性別返回已通過後、664M至784M之間讀已量數值座標的目的緩衝，最多8 bytes，只接受1至3位十進位數字及NUL。沒有讀MSG、鄰文、詞庫或手冊答案。男女各一趟固定輔助輸入，各13筆事件及7個停止點，退出0。與R935的indexed、R937英文控制的palette完整bytes及Step／1MiB Mem FNV／A000 FNV／Pal FNV／Mouse相比，14對全同。這不包含完整CPU暫存器或新正常GUI。

四組原型為男女×兩字型，640×400；原生15／16字模保持。31個欄位含17標籤、性別、11數值及生命／體力兩個數值列。屬性改18像素行距，性別字格擴至16像素；數值直接搬原版差異墨點，仍以2倍整數放大，不重算或重新編碼。只還原原版前後差異的像素到精確before，不用固定底色補底。四份現行對照圖完整RGB與R937既有圖相同；四份重排圖由另一個檢查器逐像素重建期望。八份PNG共2,048,000像素相符；四張重排圖各31欄位交疊及墨點交集皆0，四張現行對照圖仍保留既有交疊。主代理實際檢視四份重排圖。

字模沿用R937已獨立原字型核過的資源；本輪沒有重讀第三方原字型，不能稱新原字型oracle。原版248檔與38份列出的cmd/wiz7 Go來源SHA保持，四個Scratch的PCFILE.DBS完整bytes與原版同；這不證SDK完整來源閉集合。其他職業、長數值、更新／撤除與正式HD／Draw尚未驗，正式排版仍待定版與新契約READY。

V1男女均在探針自身13筆計數檢查退出2，14份indexed仍與R935相同。V1把font0篩選套到所有NUL入口，V2移除這項整體篩選後完成；十一個實際數值仍全部font0。故障原因為生命／體力結束訊號被篩掉的強推論，結束訊號的確切font未記錄，不猜補。V2另在計數檢查前保存觀察；兩版來源、binary與初輪輸出均保留。獨立檢查器首跑把正式來源的純SHA欄當成物件，修正欄位讀法後同容器／同輸入重跑通過，屬驗證器問題。

| 私人來源／收據 | SHA-256 |
|---|---|
| `profile_probe.go` | `fa657dbb9a9749e6e2170472b96d45972eafe04ec80c3d3a70cd0da58d4da60d` |
| `profile-probe-v2` | `af88a4835a2afc445d5592cb50ddf66a5f6ad50c1d5be086765d9643ad29af88` |
| `make_layout_prototype.py` | `54d612dab5f50c8e1acb9e15d62a70eea1e3a9740fcffe9a3407c13ff7a2a45f` |
| `prototype/result.json` | `834eaf6be5fccb917094a211ccd0a1bbc8cb703e89c4378d408b99093da37efa` |
| `prototype/compare.html` | `04b82d023fe59cfdcab94cee3346e0faf3ef7d4879e486bac5f9b50b071e5356` |
| `check_layout.py` | `b3391052c112e5d4a15011b38d068747900f401f92d3c95e0aaabce232a8c986` |
| `independent-layout-check.json` | `0fc80ccfd1665a2d63ee7486ee6dfd08ba92589e02ee63dc46fac89f354a6449` |

收據根為`workplace/r953`；[比較頁](../../workplace/r953/prototype/compare.html)只留本機，不提交圖片或原版。重跑沿上述映像，repo及原版唯讀、新且已存在的輸出目錄可寫；在容器內為新輸出建立可寫`gocache`，提供本輪來源及`before.json`，既有離線模組快取掛`/src/workplace/gomodcache:ro`。`go.mod`、`go.sum`與既有SDK隨repo掛`/src:ro`；SDK受控準備入口為[prepare-dosgolem.sh](../../tools/docker/prepare-dosgolem.sh)，不能換成裸上游提交或主機runtime。以`build_probe.py`建置探針，再用`profile-probe-v2 /src/workplace/original/wizar7 /out/male-v2 male`及female對應參數。需同層提供R935／R937／R943已核本機輸入，依序執行`make_layout_prototype.py`、`check_layout.py`；完整性檢查另需保留V1的male／female Scratch收據，缺項不可充作重跑通過。既有輸出拒絕覆寫。所有命令在離線、非root、有界Docker內執行，不用主機工具或注入存檔。

獨立窄審阻擋0／應改0／建議0，主代理於FINAL後完整回讀JSON與Markdown採認，只涵蓋上述唯讀測量及離線原型，不批准正式2倍實作。報告`workplace/r953/review-layout.json`／`.md` SHA依次為`b68546549aa096143b6d62ca6be6d868ed4898eeae4c48adfb94263ca5ee74f1`／`f31e51b03c0a05928f0d8c60c6cbe93945762c5da0f502c201b3b8425518851b`。已補可寫Go快取、唯讀模組與SDK入口，並收窄零交疊只適用四張重排圖。024既有READY範圍保持，2倍候選仍待選；輸出1000:1000、自有審查容器0。

## R952 人物頁圖片操作與2倍行距

入口：[R952](../goal/R952-character-review-pages.md)，2026-10-04。從目前正式來源乾淨重建Linux前端，Go1.24.13 linux/amd64，binary SHA `7363fc5be7c48d10eb272d5bd5900045996bb490b0cb08fd19bb1a93e9d948b0`；正式程式與SDK未修改。使用既有六份Unifont資源、3倍畫布、靜音與手動加速，原版唯讀。這不驗預設節拍、裝置音訊或人耳。

真正GUI由主代理逐張看實際Draw後送10次點擊，共11個停止畫格。兩次從主選單／角色選單進THESUS檢視；在已見人物頁點231,143紅色劍形圖示及207,136蠟燭圖示，下一停止畫格都回角色選單，沒有切出功能頁。此結果為confirmed的兩個點擊觀察，不證明所有圖示皆靜態或任何遊戲內分頁完成。最後經中文「離開」返回主選單，正常關窗exit0，耗時0.106417秒。原版248檔盲SHA保持；各畫格擷取後讀到的Scratch內容與原版相同，讀檔時點不冒稱PNG的固定Step。

正常收據為`workplace/r952/normal-unifont-3-attempt3/captures.json`。首輪觀察器誤等成功碼`err=none`，正式前端實際成功碼為空字串；點擊已成功，失敗在觀察確認。第二輪因主代理未在120秒內回覆畫格而超時。兩份失敗收據保留，屬工具／驗證等待問題，不記為產品缺陷；第三輪沿同一binary、容器與120秒截止值乾淨重跑。

另目視R937既有2倍人物Draw，獨立解GNU Unifont原始字模，核八個屬性兩倍率共32字模。STR至KAR的原版y為35、41、47、53、59、65、71、77；字高16，2倍行距12。七對相鄰字格各交疊4列，實際筆畫交集合計58像素。這是confirmed的R937已驗2倍排版缺陷，未冒稱本輪重跑2倍GUI。R952當輪3倍Shown的同八列行距18，字格及筆畫交集皆0。2倍修正需要同時處理文字與數值行的對齊，不能只補性別或縮字。

私人獨立檢查入口為`workplace/r952/check_gui_and_spacing.py`，在既有`wizardry7-go-ebiten:latest`離線非root Docker中執行；儲存庫、原版及字型掛`/src:ro`，僅`workplace/r952`掛`/out:rw`，執行`python3 /out/check_gui_and_spacing.py`，報告以exclusive方式寫`gui-and-spacing-check.json`。它核當輪binary／來源、10份父代理畫格SHA、22份HD／Draw檔案SHA、原版248檔及存檔內容、字型與字格交集；不重跑完整像素oracle或固定輸入A/B。024仍READY，2倍性別與其他排版、遊戲內分頁、十一組、TSV、正式錄放、冷存讀、音訊及封包均未完成。

正常收據、檢查報告及檢查器SHA依次為`fa6bea32eaad3ef9083896defa613b735a674c1cdae382c8300a1a1904810193`、`7ef08992a61362e8bac66cc89b877fd795a712d48d89cf1a6b2e4545ad3cb858`、`4af4748c577b7d106fc498328c315e1502ed19a0a5b3ae700991ca2f67cf1fc4`。本機140份技術檔一次封存，排除Go快取及後續審查／GitHub回讀；`closed-inputs.json` SHA為`bc2c3a186ba4022177f6178d924236a289218464dfa91be1bc4ad7eb52e8935d`。

獨立結果審查阻擋0／應改0／建議0，主代理於FINAL後完整回讀兩報告及SHA採認。`workplace/r952/review-results.json`／`.md` SHA依次為`61f1c439450fd84d114e2bdf68866d7233cab1478d5acedb40cd2d610eb6bd82`／`3e6b74ad507c712626bca2461121f1383be63bb7159ac54f820b102c699ded42`；另獨立核22張PNG的CRC、解壓及尺寸。只核38份列出的cmd來源，不宣稱完整編譯閉集合或SDK1,132檔本輪獨立重核。擁有權1000:1000，自有容器0。

## R937 正式修正驗收入口

[check_profile_sex.py](../../tools/re/check_profile_sex.py)核對024限定範圍的正常GUI、原版狀態、兩字型原生性別字模、選人按鈕邊界及修前／修後HD與Draw。入口參數為`<候選工作根> --baseline-root <R935工作根> --control-root <英文控制工作根> --unifont-hex <字型> --eten-font <字型> --eten-symbol-font <符號字型> --output <新報告>`；只在既有離線非root Docker執行，輸入唯讀，只開放新報告可寫。省略control-root時使用候選工作根；指定控制根時分別記錄兩個binary雜湊，不能聲稱英文控制由新binary重跑。既有輸出拒絕。限定像素驗收及R941結果複核已通過，未將024升為CONFORMED。

R937動態勘誤：16份男女／兩字型／兩倍率的選人停止點都有標題、確定及取消，但父角色選單六個標籤皆已移除。2倍Shown為3筆，3倍加六筆角色碼為9筆。原先規格要求的九訊號不能同時成立，正式helper返回舊畫筆；2倍16張修前／修後HD及Draw逐位元相同，置中未生效。私人selector-guard-evidence.json SHA `85e1cde8acd6e655f2b9002b292608c4f94780cde0ae17e2cf91e14aae0f1238`。024選人條件回DRAFT後，R940窄審0／0／0由主代理於FINAL後完整回讀採認，再轉READY。正式helper及測試只核三個精確選人訊號，不改性別、原版輸入或資料。修正後驗收另存workplace/r937/repair，初輪失敗收據保持。

驗收工具先後修正舊／新字格遮罩與逐張預期尺寸。首跑另外把短角色選單路線的80筆缺譯數套到人物檢視路線而拒絕；現行修前／修後2倍各121筆、3倍各109筆，逐項相同，不算新漏譯或完整覆蓋。修正計數後，局部像素在倚天2倍確定按鈕198,334不合，才揭露上述正式辨識缺陷。未產成功validation.json，不把失敗收據採認為完成。

修正後正式binary SHA `28f8ea051f59a8399a049c780d10f3457568201da9da960b6c4986040a300c01`，完整cmd/wiz7測試通過。男女兩次真正GUI各由主代理逐頁看Draw再送六次正常點擊，男／女中文、選人及圖片箭頭返回正常，均退出0。真正GUI不當固定步數原版同狀態；關閉耗時0.067154／0.066959秒，持久角色內容保持。再取雙字型、兩倍率、兩性別八條固定輔助路線，每條七點、均正常退出0。舊英文控制binary為`ecd047c8ca61afc683dc4cc1343738780d21b3c36ad41ca6872266f9a49f7e33`，沒有重跑新binary純英文控制。新中文對舊英文及修前中文共112對Step／Mem／A000／Pal／Mouse相同，indexed及palette bytes全同；男女原版窄探針另28對相同。

獨立原字型核68個字模、275,904個性別與按鈕HD／Draw像素，144個性別原版before／after欄位像素通過。112張修前／修後HD及Draw核46,570,368個批准區外像素，差異0；2倍選人字格與筆畫在按鈕內，3倍選人及其他未批准畫面保持。主代理另目視兩字型兩倍率的四張選人Draw。一般八資源及兩個新增性別資源保持烘製版本。

本輪檢查器曾把未印字的背景也要求等於印字後快照。女性3倍656M的THESUS底列各30像素Before＝After＝195，後續原版選取框已改為1，原筆畫失配0，與R928已量選取底邊一致。只在Before≠After處核After並還原Before，保留當前背景。R941檢查器窄複審0／0／0於FINAL後完整回讀採認；JSON／Markdown SHA `1330092346f9021f22255b5c95edbc044e5be571060b700017520c587e22b63a`／`23b02ca475e8d154959edfec48ccb1aadda3899b003d3f39863a0dc6b79c55ef`。正式Go與binary不改，舊失敗收據保持，沒有放寬差異遮罩。

成功報告workplace/r937/repair/validation.json SHA `99e99a3d9c756c9aaae84fc24a71a8867b0fb97a80dadf76307050afb74ce056`；公開檢查器SHA `ba12f128f3cb010cc19ea765ccd3c2dd8de2fcc6a2d7e13d5ea002ea59061012`。拒絕覆寫exit2且既有bytes保持。R941結果最終複核0／0／0於FINAL後完整回讀採認，不升024為CONFORMED。此路線其他無譯文事件2倍121、3倍109逐項保持；2倍性別與角色碼、其他人物分頁、圖片呈獻字樣及全遊戲完整覆蓋未完成。十一組、TSV遷移、錄放、冷存讀、音訊與封包門檻保持。

R941最終結果JSON／Markdown SHA `30c83b9caf5e21dadc6e555ce757eae6db8049001d865f374de46c447c905524`／`e9baef2bfd3df353e46df25585b43d05d55cfadcf4e175c898e783538d3f1e6d`，位於workplace/r937/repair。完整性報告integrity.json SHA `d5885bf87c1a4fb4fe2a17053bcc481eb475ce2e4ab1ebc039c06f8f9e8c22bb`：原版248、SDK1,132、R928／R932／R935舊464／335／503及初輪1,006產物保持，十次退出後持久內容保持、擁有權異常0。工作根find自檢異常0。最終複核獨立重核112對摘要及indexed／palette，另核32,444個原筆畫像素皆符合After；未重跑完整像素oracle。


R937修正後535個技術產物一次封存，workplace/r937/repair/closed-inputs.json SHA `86b6a6ef0fb95bd175b2af973127b854b2b42d3ef42d9855475f90d1b10e901b`。正式修正`450803c`已推送，#22／#28留言[5976034164](https://github.com/wicanr2/wizardry7-cht/issues/22#issuecomment-5976034164)／[5976034765](https://github.com/wicanr2/wizardry7-cht/issues/28#issuecomment-5976034765)與本機正文逐字相同、兩議題OPEN。github-readback.json SHA `2983ceec553858a763d1be24dbe968e067cc58390a6e19898ba80c23d5de483e`。原版248、SDK1,132及535／1,006封存收尾再核保持，自有R937／R941容器0。

## R935 人物檢視與返回

正式R932 binary SHA `9796d372ce53cef7528f9ad5d3a8fc57206af15cfa411cb37ab3f2499a20ed22`及兩後端各四資源保持。真正GUI冷啟動依次點160,47、85,193、175,152、100,170、314,190、287,193；每次先看當前畫格再輸入。從主選單進角色選單、選THESUS、檢視人物頁、右下藍色箭頭返回角色選單、離開回主選單，正常退出0，關閉耗時0.068531934秒，PCFILE.DBS持久內容保持。箭頭沒有英文文字，只證實314,190能返回；20×16色號外接框混有其他像素，不能當完整命中矩形。

首次真正GUI在人物頁看圖超過父代理120秒等待，已保留失敗片段；超時後寫入的返回命令未送進遊戲。主代理承認工具等待錯誤，重查GUI路由後，以相同控制器、截止值及計畫輸入乾淨重跑，不調時基或延長等待。成功收據只採workplace/r935/retry/primary/captures.json，SHA `86b657bdcb38726905ba6eb6d86889265c8e08a7dcb8705f40535891fd3f9c1e`。

其後六組英文／倚天／Unifont、2倍／3倍固定GUI各取272M、400M、528M、656M、784M、912M、1040M七點，全部正常退出0。28對Step、Mem、A000、Pal、Mouse及indexed／palette bytes相同，保存內容保持。3倍人物頁17個標籤為人類、戰士、等級、階級、無、經驗、擊殺、生命、體力及八個屬性；雙字型共59,212項局部HD／Draw像素核對通過。196個原生字模核對通過，3,180個原版after像素與英文indexed相符。一般標籤before來自既有攔截收據，不冒稱另一套原版探針獨立讀得。

選人對話框「檢視誰？」及確定／取消字模正確，但2倍倚天兩按鈕字格各超1列，Unifont各超2列、筆畫分別超17／6像素；3倍全0。人物性別M仍英文，不把17個標籤通過寫成全人物頁完成。其他分頁與人物頁2倍完整排版尚未量。

追加女性正常GUI只將選人點改為175,158，檢視既有TEMPEST；主代理逐次看圖，確認女性蜥蜴人頁仍顯示F、箭頭返回角色選單再離開回主選單。正常退出0、0.071909746秒、持久內容保持。返回命令守衛依本次畫格的LIZARDMAN102,6核對，沒有套用男性HUMAN訊號。收據female-primary/captures.json SHA `c1d4970f43abc211dfa2ac884402d3c267f6ca23d443fec545626f2055aba7a2`。未取得女性英文固定GUI控制，不聲稱女性同狀態或中文性別完成；先前驗證report的female未完成限制對應追加前範圍。

### 單字元性別的窄範圍證據

原版runtime0110:3877的near返回0110:6c42，StackWord(1)為M或F、StackWord(2)為0；入口游標90,6、SP65316，返回SP＋2、游標96,6。一般字串入口0110:3895沒有收到目標M；本證據補充既有013字串管線，不能再把所有玩家文字都等同NUL字串入口。

獨立dosgolem冷啟動、相同正常操作，M選THESUS、F選TEMPEST。只有664M至784M、上述地址／字碼／座標組合被觀察。全64000 bytes離屏差異分類：M有13個、F有11個改動，都在90,6,6,6內、改動後為色號1；改動前含0、14、115–118，底圖有紋理。只清原筆畫並還原精確before，不以純色擦底。這些原始36-byte樣本、原版畫格及完整收據只留本機；未讀訊息檔、鄰文或手冊答案。

M的兩套唯讀探針各七點與英文控制摘要及indexed相同，共14對。F沒有獨立英文GUI控制，不列為正式中文或同狀態GUI完成。私人重跑來源workplace/r935/sex_probe_full.go，SHA `7ec6fa88506bb1011a335cd04e3bf8f24f8b0ca4d05cd66ca6a23ec49e82f10e`；binary SHA `217a153fc302a3314ff0835a6b99ff1b026f37c39c969c30d5646deff4543cae`，入口`sex-probe-full <原版唯讀目錄> <新輸出目錄> male|female`。在既有離線非root image內建置／執行，SDK及原版唯讀，輸出限定workplace/r935新子目錄；不在主機執行。M／F observation.json SHA依序`9491b606f489ca13717fc41a7cac3fe23c339bd39bb8d5a347f7ad5b4bec852e`／`7ec495a48718248598291db66d21b17e7bc7c2c813a39a643c677dd576ede771`。

### 獨立驗證入口

公開[check_character_review.py](../../tools/re/check_character_review.py) SHA `11ba230178b4a0ece4f5eb7b34e281bd73f6e7f577c03df85c98fe519f835227`。工具直接解原字型、原始色號及PNG，不匯入正式renderer。二十個精確標籤只沿正式已知文案核對，沒有讀取MSG.DBS或messages.json。使用既有離線非root Docker，原版及字型唯讀，只有新報告位置可寫：

```text
python3 tools/re/check_character_review.py workplace/r935 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --eten-symbol-font /eten/ET353S/FILES/SPCFONT.15 \
  --output /reports/new-character-review.json
```

倚天全形問號來自SPCFONT.15，12240 bytes、SHA `f32049ba2a7a21db908878a488a2c1d93c389d17398cf46db1390ba89e247605`。STDFONT及Unifont雜湊沿R932，工具內固定核對。首次工具錯把問號從漢字字庫負偏移讀取，核符號字庫契約後修正工具，以同一GUI收據重跑通過，正式字型及文案未改。

私人與公開report bytes相同，SHA `251f530d861fcd9c890612db7347ffe766f71cac68f90a6587a587c6c084641e`；既有輸出拒絕exit2且bytes保持。integrity.json SHA `897234c312c04206708ea4c8d794f1cc8b16391372cce1dc80e85df08dff5793`，核原版248、SDK1132、R928舊464封存、R932舊335封存全部保持，R935擁有權異常0。全工作根find自檢異常0，自有R935容器0；初次Python自檢跟隨既有斷連結而失敗，改用規定find完成，沒有產品重跑。新正式修正見[024](../spec/024-character-review-output.md)，仍DRAFT，不把本次證據計為性別中文或全遊戲完成。

R936初審2阻擋／1應改／0建議，主代理於FINAL後完整回讀兩報告。R935限定證據可採，024資源型別入口及pending仲裁需補明後再審；JSON／Markdown SHA依序`2e3aa2c65bf3bb1ae32395c699582bd98632c49b81286c8b6df54b15c010470b`／`c9dc40d3cd3ac944d3551f470d88f101d7c2cdf665500430604bad077fee5fea`，入口workplace/r935/review-r936.json及.md。新正式工作目標見[R937](../goal/R937-character-review-fix.md)，READY前不實作。

R936窄複審0／0／0，B1／B2及字形／女性現況文字皆解除，主代理於FINAL後完整回讀採認，024轉READY。兩報告SHA見024；只准限定實作，性別中文及按鈕修正尚未正式驗收，完整交付門檻保持。前段DRAFT描述為初審前狀態。

R935的503個本機產物一次封存，closed-inputs.json SHA `bf8b59aebd98026d3d0375c29ffe9674e32f25aa0748f89a021912f1a2bd7908`。完整原始收據只留本機，正式修正另在workplace/r937，不覆寫基準。

## 驗收欄位

每筆記錄畫面與控制、文字輸出或圖片字樣的來源、原版定位及座標空間、辨識訊號、正式中文畫格、點擊與回復結果、同狀態收據、字形及邊界、證據等級與下一個缺口。

正式程式已完成、動態畫面未量到與可選美化分開記錄。無文字圖示沒有可翻譯字樣，但其說明與操作仍須核對。玩家自行輸入的名字與既有原版英文關鍵字規則保持，不把它們誤計為未譯按鈕。

## 本輪取樣

2026-10-04先由修前正式來源重建Linux前端，SHA-256為5fe5f410065ecf14303e896e72d00eccdf93a70f11ff30fbdafd050eef560bd5。原版DS.EXE、SDK與字型輸入唯讀；完整來源摘要在workplace/r915/inputs.json，產物摘要在build-artifacts.json。主選單修正後正式binary為7c2553a39f4fef55efb49c7c7df4eae37b4eccac5761371079a14eedae8eb18d，最後來源重建與GUI所測binary逐位元相同。

兩套後端都由現行來源生成：一般覆繪資源8,425個條目、2,867個中文字形，流動文字2,759個字形，說明14頁、426個字形。字元聯集2,919個，非ASCII字元2,872個。資源生成成功只證明字模可生成，不證明每個按鈕上屏。

第一段以正式GUI跑倚天、Unifont及英文版，各2×／3×，無玩家輸入；每次取2M與60M步，12張HD與12張實際Draw均完整解析且尺寸正確。六次WM_DELETE_WINDOW正常退出碼0，關閉請求至退出0.019至0.183秒。這是靜音視覺取樣，不等同聲音或完整可玩驗收。

相同步數八組中文／英文對照的Step、原版記憶體摘要、indexed與palette摘要、滑鼠狀態均相同，.a000及.pal逐位元組相同。comparison.json SHA-256為cf3fc541a08506040c51b6a70feb67870209632ef21270f978811d0bebd475b0。Draw是實際視窗幀，沒有凍結在標稱步數，不拿它逐像素宣稱固定步數同狀態。

## 目前覆蓋矩陣

| 畫面／控制 | 來源與辨識 | 中文與操作結果 | 狀態／下一步 |
|---|---|---|---|
| 冷啟動2M，廠商標誌 | 舊正式Draw已見圖像；R942正常唯讀路線量到MON60.PIC載入及首張圖片轉送 | 無操作按鈕；品牌字樣仍是原版 | 載入與轉送已定位，codec與品牌分類未驗 |
| 冷啟動60M，廠商標誌及呈獻字樣 | R942定位MON60最後圖；R943依使用者選定的周邊紋理補底接入六張呈獻覆繪 | 雙字型兩倍率、32點原版對照及六趟正常冷GUI通過；品牌依既有規則保留 | 限定呈獻通過，翻頁傾斜字樣與介紹仍沿R944待選，不外推其他圖片 |
| 片頭後段翻頁敘述 | R942的MON61索引組與六張完整轉送已量，圖片內另有英文角色扮演敘述 | 敘述尚無中文；作者與品牌另依既有名稱規則核分類 | 僅原版無輸入觀察；不稱正式中文GUI通過 |
| 主選單六個標籤 | 正常無輸入272M，六個resident文字覆繪事件；不是圖片內固定字樣 | 六標籤完整；「匯入巫術六」五字；雙字型2×／3×的HD及實際Draw都在按鈕內 | 角色選單與新增角色入口／返回通過，其餘入口未量；防禦分支僅單測，見[020](../spec/020-main-menu-labels.md) |
| 片頭製作名單的未翻譯事件 | 四組中文run.log各80次；75次精確對應15個去低位碼後的正式英文鍵，5次空字串 | 15筆來源的zh與en相同；尚未逐項確認保留內容的分類及實際畫面 | 不當75個新漏譯，也不計全文覆蓋通過 |
| 新增角色對話框標題及確定／取消 | R921、R925正常GUI及固定步數已量；既有resident文字覆繪 | 雙字型兩倍率字格與筆畫在框內；正常取消返回 | 限定通過，[021](../spec/021-character-dialog-labels.md)仍READY |
| 新增角色清單六列 | 完整名字＋性別職業動態文字，返回位置0110:0fcc | 名字依既定規則保留英文；六個性別職業3倍雙字型完整，選取及取消正常 | 六個實例通過；2倍、任意姓名與捲動未完成，見[022](../spec/022-character-roster-text.md) |
| 角色選單六按鈕 | R932正常GUI及resident0110:5686、六個精確原文／座標 | 創建／檢視／刪除／改名／頭像／離開完整；雙字型兩倍率字格及筆畫在框內，離開返回正常 | 六按鈕限定通過；人物檢視及頭像頁另見024與R955，其餘子畫面仍待完整驗收，見[023](../spec/023-character-menu-labels.md) |
| 頭像頁標題及完成控制 | R955正常GUI到達SELECT PORTRAIT於116,138、CLICK HERE WHEN DONE於101,153 | 「選擇頭像」及「完成後點這裡」雙字型兩倍率字格完整、局部HD／Draw像素相符；保持頭像完成返回，角色檔相同 | 只驗兩則文字與保持原頭像返回；選新頭像、其他圖片控制及已知2倍人物／名單缺口不計通過 |
| 人物檢視頁 | R935／R937正常GUI、窄印字證據與固定冷啟動 | 3倍雙字型既有標籤及男／女性別通過，選人確定／取消兩倍率在框內，箭頭返回正常 | [024](../spec/024-character-review-output.md)READY；限定結果複核通過，2倍性別與其他分頁未完成 |
| 人物頁右下箭頭 | 無英文圖示；R935實際點314,190 | 正常返回角色選單；只證實此點，未量完整命中矩形 | 圖示無待譯字樣，不外推其他圖片按鈕 |
| 清單上下箭頭 | 已見無英文圖示，操作及說明未量 | 尚未完成控制驗收 | 另沿正常GUI核對 |
| 建立角色、其餘人物／物品／咒語／對話／戰鬥／存讀及圖片按鈕 | 此次重查尚未到達；既有規格與歷史收據須逐項核現行來源 | 未量到，不計本次通過 | 依可見畫面走合法操作 |

首輪建置的Go編譯並行數超過容器程序上限；限制GOMAXPROCS=2及-p=2後，同一工具鏈重跑成功。原失敗日誌保留為build.log，成功日誌為build-p2.log。這是工具環境修正，不登錄產品缺陷。

第二段Unifont與英文3×無輸入至200M、272M，兩次正常退出0，兩組原版摘要與色號／色盤相同。200M是標題與著作權畫面，272M自行出現主選單，不需要猜補按鍵。匯入原文IMPORT (BANE)的兩份正式來源1007／1415都譯為「匯入(巫6)」；bake_font.py:165–167略掉非雙位元組字元，因此只生成三字。這是confirmed漏字，不是圖片按鈕或缺字型。

六個按鈕邊界由英文.a000色號176／192／197逐行、逐欄獨立量測，結果在menu-bounds.json：左105、寬109、高9，頂點依序43、56、69、82、95、108。原文字y為45、58、71、84、97、110；修前2×原生15／16列字格底部各超出1／2像素，實際Draw亦已核對。修後倚天2×上移2像素、Unifont2×上移3像素；3×位置保持。字形不縮小，原版按鈕與點擊範圍不改。

量測腳本初次誤從英文模式讀Shown事件，因英文模式沒有中文事件而失敗，未寫結果。修正為已對拍中文事件提供定位、英文原版色號提供獨立邊界，乾淨重跑成功，不作產品缺陷。

## 正式修正及獨立驗收

正式catalog只改1007、1415兩個zh為「匯入巫術六」，en、ID及其他欄位保持。兩字型xlate各只改匯入鍵；flow各只有同兩筆行／訊息由無可用字形變為五字，增一個「匯」字模，其餘字模筆畫與欄位保持。help逐位元相同，字元聯集不變。完整差異收據在menu-fix/resource-diff.json與pixel-validation.json。

[menu_labels.go](../../cmd/wiz7/menu_labels.go)只在六個精確英文鍵與座標同時存在、dy為0、字高一致15／16且字寬足夠時，調整HD及Draw畫筆的目的y。原事件、源字模、原版記憶體、資料比較、存檔及滑鼠操作均不修改。整個cmd/wiz7測試通過；排除情境未有正常GUI收據，因此020維持READY。

六次新正式GUI及四次修前binary搭配同一新資源的GUI，全部於272M取完整HD及實際Draw，正常關窗退出0。這段採手動-ips 24000000及靜音，兩側CPU／PIT及虛擬時鐘相同，只作視覺回歸，不作預設節奏或音訊驗收。四組中英Step、Mem、A000、Pal、Mouse與.a000/.pal相同；四組修前／修後原版摘要及.comp相同。

獨立解析倚天stdfont.15及固定Unifont17.0.05 HEX，100項字模筆畫與306,072項按鈕像素核對全部通過。HD與實際Draw的六個按鈕均核整個矩形，匯入五字完整。2×修前／修後變更倚天1,940像素、Unifont2,258像素，批准舊／新字格以外差異0；3×同資源位置差異0。主代理已目視四組修後實際Draw。這些數字只證明六標籤的顯示，不證明按鈕操作或全部中文化。

可重跑的獨立檢查器為[check_menu_labels.py](../../tools/re/check_menu_labels.py)。在離線非root Docker內使用下列參數；驗收輸入唯讀，只有新報告所在目錄可寫。既有--output會拒絕覆寫。

```text
python3 tools/re/check_menu_labels.py workplace/r915/menu-fix \
  --bounds workplace/r915/menu-bounds.json \
  --baseline-bake-root workplace/r915 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-menu-validation.json
```

Unifont HEX的SHA-256為7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574。原版畫格、倚天資料及完整收據留本機，不提交。公開檢查器已以新的報告目的地重跑，結果同為306,172項通過。

像素檢查器首跑誤將按鈕斜角色號當成文字；修正為由原版按鈕內部區域獨立辨識前景，乾淨重跑通過。這是驗收工具問題，未修改正式顯示來迎合錯誤oracle。公開檢查器初版的固定報告目的地曾覆寫相同內容，現已改為明示新路徑及拒絕既有檔案。

## 完整性界線

四組中文run.log的80次事件完全一致：75次、15個非空英文鍵與現行catalog去低位碼後相符，且各筆zh等於en；5次為空字串。出處為2901、2902、2904、2906、2908、2909、2911、2913、2915、2918、2920、2921、2922、2923、2924。沒有讀取防拷或詞庫正文。完整metadata收據為workplace/r915/untranslated-events.json，SHA-256為5b6d5ee93cb99b9b2ac3a57bd5c48acde57f3ab51de9c0393d41c9f6e2b02f37。來源相同不代表每項保留英文都適當，後續仍須核人名、平台名、職務標題及畫面來源。

靜態ASCII略過候選由1,604降為1,602；其中Latin115、數字134、標點1,500可重疊。標點包含合法控制碼，候選不等於確定漏字。沒有逐筆查證前不批次改文案；其他畫面、圖片按鈕與十一組交付門檻保持。

## R921 新增角色對話框

沿正常主選單，以實際GUI左鍵160,62開啟新增角色對話框。400M畫格已見標題「新增角色」及「確定／取消」。取消按鈕的原版色號邊界為63,182,73,9；主代理看圖及核邊界後才單擊100,188。528M返回選單，標題及兩個對話按鈕移除，WM_DELETE正常退出0，退出耗時0.067秒。這只驗新增對話框及取消，不驗選取、建立角色或完整遊玩。

原版既有PCFILE.DBS寫開時可在獨立存檔層生成複本。第三次正常GUI收據列9,368 bytes、SHA-256為24d94441e5b32e7abbcbffa9dad972b41544da6a7bb06096b967281f639d9c1b，與唯讀原版全同；沒有變更持久角色資料內容，不聲稱零DOS寫入。控制器舊stdout「saves empty」不準確，已在WORKLOG勘誤，原收據保留。前兩次因日誌格式及空目錄斷言失敗，未正常關窗，不計成功。

正常GUI成功後，倚天／Unifont各2×、3×及英文兩倍率固定280M新增、440M取消，於272M、400M、528M取得完整HD及實際Draw，六次均正常退出0。四組中英文、三個停止點共12對的Step／Mem／A000／Pal／Mouse及.a000／.pal全同，退出後PCFILE內容全同。檔案取樣記控制器實際讀取時間，不能把稍後讀目錄冒充畫格Step的精確檔案快照。兩次早期輔助空白斷言失敗仍保留，後續沒有調整步數、輸入或時基。

獨立解原始兩種字型，32個標籤字模及93,950個HD／Draw像素核對通過，主代理也檢視四組中文Draw。原版確定按鈕63,166,73,9，取消63,182,73,9。2×倚天字格底部超出1像素，但本次兩詞的實際筆畫未超出；2×Unifont字格超出2像素，確定17、取消6個筆畫像素在按鈕外。3×兩字型均無超出。這是confirmed排版缺陷，不能把像素正確呈現舊排版算成按鈕內排版通過；限定修正見[021](../spec/021-character-dialog-labels.md)。

角色清單有六個完整動態字串沒有匹配，來源是resident印字返回位置0110:0fcc，遊戲座標x150、y150／156／162／168／174／180。[013](013-print-pipeline.md)已列「名字＋職業碼」家族，因此不是圖片字樣。預設名字依[names](../names.md)保留英文，性別職業仍需中文。四組缺譯事件同為86次，較主選單基線多此六筆；其他80筆分類界線保持，不宣稱漏譯0。清單上下箭頭為無英文圖示，尚未操作或完成說明核對；片頭圖片字樣與其他畫面仍未完成。

R921正式binary為7c2553a39f4fef55efb49c7c7df4eae37b4eccac5761371079a14eedae8eb18d，本輪正式Go／catalog／生成器沒有變更。原版248、SDK1,132與來源基線保持。私人權威收據：workplace/r921/primary-r3/captures.json SHA為3fda0d582631d9c80be2f21462263ef196709e538234da58df81ef26a727100d；dialog-validation.json SHA為961d769f8bc68035ecc5933e60b454254b951de4c29cf0ca8fba31efa9412d98。R922、R923、R924契約與取樣修訂均經主代理完整回讀後才採用；R926結果複核0／0／0，主代理已採認。十一組、TSV、錄放、音訊、冷存讀及封包門檻保持。

## R925 小按鈕正式修正

[021](../spec/021-character-dialog-labels.md)經R926獨立計畫審查到READY後，增加[dialog_labels.go](../../cmd/wiz7/dialog_labels.go)，HD及Draw共同使用buttonLabelY。只在九個已量畫面訊號、指定兩按鈕、dy0、原生15／16字高及寬度守門成立時置中文字格於9列原版按鈕中間。2×倚天上移2像素、Unifont上移3像素，3×保持；020主選單及標題保持。main.go只換兩個helper呼叫，catalog、生成器及字型資源不改。

cmd/wiz7測試全部通過，正式binary SHA為ce62ac638b89d42de0d10239d08ebbaf2c1b849ce13aa80908f4a655a859cb6b。新binary先以實際GUI完成新增、取消100,188、528M返回及正常退出0，退出0.097秒；主代理逐次看圖後才送操作。PCFILE只有與原版相同的複本，持久內容保持。再跑四組中文固定路線，全部正常退出，與R921英文及修前共12個固定停止點的原版摘要、色號／色盤相同；修前後Shown及.comp也全同。

獨立32字模、92,198個HD／Draw局部像素通過，八個中文按鈕的完整字格與實際筆畫超框均0。固定步數HD全圖比較：2×倚天改324像素、Unifont改298像素，批准舊／新字格聯集外差異0；272M及528M、全部3×HD差異0。Draw只核上述局部字模與按鈕，未作新舊全圖差異比較。主代理檢視四組新Draw。防禦分支僅單測，021維持READY，不外推其他同名按鈕。

可重跑檢查器為[check_character_dialog.py](../../tools/re/check_character_dialog.py)，需現有本機R921基準與兩種合法原字型；報告明示新路徑並以exclusive方式建立。執行一律在離線非root Docker，原版、字型與基準唯讀，只新報告目錄可寫：

```text
python3 tools/re/check_character_dialog.py workplace/r925 \
  --baseline-root workplace/r921 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-character-dialog.json
```

六列動態性別職業仍缺譯，86次缺譯分類保持。這次修正只解決新增角色兩按鈕的排版，沒有完成圖片字樣、箭頭操作、完整中文化或交付。

R927最終獨立複核0／0／0，主代理完整回讀後採認。私人報告workplace/r925/review/review-r927-final.json／.md的SHA為c8bc00f9ed49b8de9808a3c1c59c910dcd97b6c6b0e7ebe67838476a594ebc81／3c26613057a75d0c4fd91b45b323c1c225fff9389f8af89705cb80edc76ee491。R921的244檔與R925的137檔封存保持，closed-inputs.json SHA依序為4e83af8f928ce462e954e292e29b3f18d112f9ff41b6e25d87263b0bb9736eb4／37c4ffec07794e7f6c865cd9ec97457185657ac524de59d98a770950d52b8f5b；來源、SDK及原版雜湊保持，輸出owner1000:1000，自有容器0。

## R928 角色列的目的區域

沿R921固定路線及同一原版載入設定，以私人唯讀探針在0110:3895限制CS0110、返回0fcc、x150、六個已量y後，只讀16 bytes確認六個精確英文鍵及NUL。0110:3877只在目標呼叫期間讀字元目的座標；返回只接受SP＝入口SP+2。沒有讀MSG、鄰文、手冊答案或詞庫，沒有修改原版記憶體。這是同一正常GUI路線的無頭補充量測，不能當新GUI或角色選取驗收。

confirmed：六列字型均0，各15字元，x依序150、156、162、168、174、180、186、192、198、204、210、216、222、228、234，y在該列保持，返回x240。姓名筆畫只在x150起的前十格；性別職業碼筆畫在x210..239、y..y+5，未選取時差異色號均195→180。碼帶固定x210..240、y..y+6。初次只量90×6區域；後續全64,000像素補量證實六次的已知姓名格、碼帶之外差異0，未知色號0。補量來源SHA `5791f78851bf78f2b58243b9298946a91079c2cd656d95c43708206497fbde15`，whole-offscreen-validation.json SHA `d53dba9969b62966237c8e7d53ebaed6debd056e5d9d74cf7522e91d6ea34093`。

272M及400M原版Step／Mem／A000／Pal／Mouse與R921英文3倍收據全同，indexed逐位元相同。PCFILE只有與原版相同的9,368 bytes複本；原版248與SDK1,132雜湊保持。私人探針來源SHA `04466d99e7b86581256bee67f7b47cd2b7598ae1b6731e9850b26c1d6127981d`、binary SHA `12f0497d7e60aa67657989aeb622ce8ba72236ea258ed34ac50491ad3c12a047`；observation.json SHA `02c2ac5c8bcdba51339e661f43eae502c9c79f25db08a1f7d5c00ac616993c01`、observation-validation.json SHA `c08a9114fe74560b83801b5d13be6392cb07c9888dc62b573bfca7dca72c8b36`、probe-integrity.json SHA `5bcc05374e9608bfb4d7a5070361e79d15ffccad17d26773261a5a618076cbef`。檔案均位於workplace/r928，沒有原版素材加入Git。

原生15／16列字型在3倍的18像素列距內；最長候選「女鍊金術士」寬83／88，在90像素碼欄內。2倍列距12、碼欄60不足，不能縮字或省略職業譯名來掩蓋。私人Canvas原型workplace/r928/roster-layout-prototype.html，SHA `c8a09212cfb172e8f6a2203b0cb124467d23709155294fcc128b89064cda1d8c`，只展示既有Draw底圖的2倍限制與3倍候選，不是正式GUI或玩家操作收據。後續契約為[022](../spec/022-character-roster-text.md)，R930複審後為READY；正式六列中文驗收進行中，原版任意姓名格式、捲動生命週期及2倍完整中文仍未完成。

實際GUI後續完成新增、選取THESUS、選取PENTAS、取消返回及正常退出0，退出0.119秒，未按確定、沒有建立角色，PCFILE內容保持。首輪控制器等待父代理操作超過120秒，failure.json保留；同控制器、同預定操作以新目錄重跑成功，沒有調整遊戲時基或程式。成功收據selection-retry/selection-primary/captures.json SHA `dd4165b53fc0495af2663c1ca6b5568015da425cbe2760a388fb9e679add7355`。

無頭補量三批共18次精確印字，兩個被選取的印字為1→0，另16次為195→180。每次姓名與碼均有差異，其外差異0。實際GUI白底範圍x148..240，首列y149..155、末列y179..185；末列選取時較晚畫的白底會改上一列TREON第六列空白。獨立oracle曾漏算此後續重畫，失敗紀錄保留；先量indexed邊界及原字模，再修正oracle，不改有效收據。12個Unifont原生字模末列均0，3倍中文筆畫不碰此原版底邊。獨立18,818項通過，selection-validation.json SHA `ee67ebc0adc5146840be214b50d31c4b0e407fc04b05c9861af5d36889967c86`；probe-selection binary SHA `aa6b8f755ffb7a18a053af659094ef30e8d3e07abbc2a3c41b079ecf7f7783ce`，觀察SHA `1f0813c724971703fcf1c72cbeeb120a029e1eef5455680179e05ec009921aee`。全部位於workplace/r928。原版248、SDK1,132及持久角色內容保持，root-owned輸出0。無頭與GUI輸入時間不同，未宣稱兩者完整Mem相同；未實作正式六列中文。

## R928 正式六列覆繪

R930契約複審通過後，新增[text/roster.zh-TW.tsv](../../text/roster.zh-TW.tsv)作唯一正式文案；[bake_roster.py](../../tools/re/bake_roster.py)烘製兩種原生字型，既有xlate／flow／help bytes保持。[roster.go](../../cmd/wiz7/roster.go)在一般查表早退前辨識六個完整原文及精確呼叫訊號。專用事件只清除30×6碼帶，完整90×6列只作身分見證，包含已量白底邊界例外；選取前景0合法。未知完整列會撤掉舊事件。譯文不寫入原版或存檔。TSV接線另核路由43的顯示與語意隔離契約，原版流程保持。

正式binary SHA `aca682fdcba57cb56225fd43164961805d53899a9e3847f49356433d70d87d54`。cmd/wiz7測試通過；新TSV與typed資源守門另驗九種格式拒絕及禁止覆寫。真正GUI完成新增、首列與末列選取、取消返回，正常退出0、0.072秒；主代理每步先看實際Draw才送下一個操作。兩種字型的最長五字完整在欄位內，姓名保持。這趟靜音／手動加速GUI只驗顯示及操作，不當完整音訊試玩。

修前／後八組固定路線，各五停止點，共20對Step／Mem／A000／Pal／Mouse摘要相同，原始indexed及色盤bytes全同。HD及Draw各20張修前後全圖比較，批准碼帶之外差異0；2倍及主選單／取消後全圖差異0。3倍三個列表畫面每張倚天改3,185像素、Unifont改3,282像素。獨立原字型核180個字模、輔助116,640及真正GUI58,320個局部像素，共174,960項；中文框與筆畫均在90×18欄位內，取消後六列中文全移除。原版248、SDK1,132、九程序持久角色內容全保持。這是六個已量實例的中文化，不推廣到任意姓名、其他職業或捲動。

3倍兩後端及真正GUI此家族缺譯0，其他缺譯仍80次；2倍此家族仍18次、總98次。2倍不能縮字，目前仍有完整中文閘門，不代表批准交付時留英文。圖片字樣及箭頭操作仍未驗。022防禦分支只有單測，整份維持READY，不標CONFORMED。

獨立檢查器[check_character_roster.py](../../tools/re/check_character_roster.py)需要本輪固定路線與合法原字型，入口見[字型說明](../../font/README.md)。初次60秒期限不足，未產成功報告；只快取PNG解碼、按工作量給300秒後，同收據與判準通過，正式Go及binary未改。私人implementation/roster-validation.json SHA `3bf34ff9ac9875bab759da697680cf676aa2cd2d02d92005e1c60029870acb03`，candidate-integrity.json SHA `5b43b08c4551212eb4a5bc0981348193c0e7b95948b0fd05ca06929b15f9c383`，parser-validation.json SHA `b1ff370437324c214395db3c41927463a3d0660dd1cd317d3d6c532fb42018fc`，位於workplace/r928。

R931實作及結果窄補審均0／0／0，主代理完整回讀後採認。結果JSON／Markdown SHA `34fcad80e19bbee53a474a2bc042ba244b4107989958f8201b64e911d877966f`／`0d300ba1ace41c5583e0a60e3d5f97b3bbe469a7be61b58558102a14f9092557`。不改十一組、完整TSV、冷存讀、音訊／錄放及交付門檻。

## R932 角色選單六按鈕

沿正式R928 binary及Unifont3倍，真正GUI正常點角色選單160,47，400M完整檢視六按鈕，再點離開287,193，528M返回主選單。正常退出0、0.129秒，持久內容保持。首輪父代理再次超過120秒操作期限，失敗收據保留；同控制器／程式／預定操作以新目錄重跑成功，未調遊戲參數。重跑primary/captures.json SHA `ab4d1bd4ff9689b56b5ffdb0e81c4f8e384bb1b4cf4eea728cd2f2fd484a5b20`。

六個標籤為CREATE／REVIEW／DELETE／RENAME／PICTURE／EXIT，現有文案為創建／檢視／刪除／改名／頭像／離開。唯讀正常路線補量確認全部走resident0110:3895，返回0110:5686，字型0，y190；x依序11、63、115、167、216、277。這六個標籤是印字文字，沒有證明其他圖片按鈕已中文化。背景美術及未操作的頭像子畫面不外推。

原版indexed量到六個按鈕皆49×9，左上為x4／56／108／160／212／264、y188。四組中文及兩組純英文正常退出0；三個停止點共12對Step／Mem／A000／Pal／Mouse相同，indexed及色盤bytes全同。獨立原字型核48字模及141,120個HD／Draw局部像素通過，真正GUI400M的indexed／色盤與HD／Draw PNG也與已核Unifont3倍逐位元一致；不宣稱真正GUI與固定輔助路線的Mem相同。

confirmed修前排版缺陷：2倍倚天完整字格超底邊1像素，本次六詞筆畫超框0；2倍Unifont字格超2像素，六詞的筆畫超框依序18／12／9／12／7／16像素。3倍兩後端字格及筆畫超框均0。正確畫出舊排版不能算按鈕內排版通過。限定修正契約見[023](../spec/023-character-menu-labels.md)，R933審查後READY，修後結果見下節。四組中文總缺譯事件仍80，這六詞缺譯0；不能稱全遊戲缺譯0。

私人workplace/r932/menu-validation.json SHA `fb236a9f4d2fba933033d783c47e3a2011310801174b020ea3ca270c7c00fdde`，primary-integrity.json SHA `600dcc3d77bfdcd31710a616a5dab771d36e461d7a7af6d32e287bd000df9145`。探針source／binary／observation SHA `1dac0e28ad6246384836bf2da7b4b69f2d40fa3f4af45ed5786ae650cf78b691`／`b9c1148b639fcacb0f8a954e1f61ac2394d9b16506eef1c450be92109861b10a`／`1541a2bb0a689248f68339056078bcc82d20bb87c7e5250d701ad550445438d3`。三個探針停止點與英文輔助摘要及indexed全同；原版248、SDK1,132保持。這段觀察未改正式Go／catalog／字型，也未進入建立、檢視、刪除、改名或頭像子畫面。

## R932 六按鈕正式置中

R933獨立審查0／0／0、主代理完整回讀後，023轉READY。新增[character_menu_labels.go](../../cmd/wiz7/character_menu_labels.go)，六個精確訊號皆在Shown且字模為一致的原生15／16行才調畫筆y。全家族高度、位元及各自寬度均守門；失敗沿020／021位置。HD及Draw共同走buttonLabelY，不改x、原版矩形、事件生命週期、文案或字型。2倍倚天y378、Unifonty377，3倍皆570保持。

正式binary SHA `9796d372ce53cef7528f9ad5d3a8fc57206af15cfa411cb37ab3f2499a20ed22`，完整cmd/wiz7測試通過。真正GUI正常進角色選單、按離開返回及退出0、0.073秒；持久角色內容保持。主代理每次看圖後才送左鍵，沒有預注入輸入。400M的HD／Draw／indexed／色盤與已獨立核過的Unifont3倍輔助逐位元一致，沒有宣稱兩條路線Mem相同。

四組固定中文程序正常退出0；12個停止點與R932英文及修前原版摘要、indexed／色盤相同，Shown及.comp保持。HD及Draw各12張全圖比較：2倍400M倚天各改892像素、Unifont各改1,108像素，批准舊／新字格聯集之外差異0；主選單、返回及全部3倍全圖差異0。原字型獨立核48字模、141,120個局部像素，24個按鈕／字型／倍率組合的完整字格及筆畫超框全0。主代理檢視四組實際Draw。

原版248、SDK1,132、R928封存464檔及兩後端各四資源保持，五程序正常退出及持久內容保持。四組80次無譯文事件與修前逐列相同，六標籤此路線缺譯0。防禦分支只有單測，023保持READY；其他子畫面、圖片字樣、十一組、TSV、錄放、冷存讀及交付閘門保持。

可重跑獨立檢查器為[check_character_menu.py](../../tools/re/check_character_menu.py)，原版、字型與基準唯讀，新報告目錄可寫；拒絕既有輸出exit2且報告bytes保持已核。一律在既有離線非root Docker，原字型入口見[font/README.md](../../font/README.md)：

```text
python3 tools/re/check_character_menu.py workplace/r932/implementation \
  --baseline-root workplace/r932/retry \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-character-menu.json
```

私人workplace/r932/implementation/menu-validation.json SHA `72caceae97226b6acb69d55e8a29b995e1893c71ce298651b0a07cec64072a4b`，integrity.json SHA `ac04632f97c7e45f4bcf840de02c7fce08655498caa061b25c194c320141da31`，primary/captures.json SHA `f32a75ed0f0b3234d24644f5061506a0940c4dbe5a5b15a6386cf747e3af770c`。完整來源SHA與報告拒絕覆寫見同目錄formal-summary.json。

R934最終獨立複核0／0／0，主代理於FINAL後完整回讀JSON／Markdown採認；SHA依序`b60e8a651023924217baf1b80945e2d6db5781fb1afb53134804e6c930804b08`／`079f4426de21d9ac3a850bf363a23dfcd8899677fe11f6a88fced28817ea3963`，本機workplace/r932/review-r934.json及.md。335個產物一次封存，closed-inputs.json SHA `73089cfbdcfe6086a73bb544e11777f1c24311d5ef2aaee9faa9dd19ec852fb5`，原版及SDK保持，工作根root-owned及.md目錄異常0，自有容器0。023仍READY，推送回讀另補。

## 生成來源一致性與審查


R919最終正式實作及生成來源同步補充計畫審查各0／0／0，主代理已完整回讀。報告JSON／Markdown的SHA-256為90a7d415007bfbba287460ce3d506c639d79a5e819d760e5bf5747732ccd416a及f288f1aaaefda67ea0a9f1ed9e7426aa69d17eba235e0166ba622a5d3e6c40fe，兩份owner1000:1000。checker末段現以open('x')獨占建立，新目的地重跑同306,172項，拒絕既有報告退出2且bytes保持；三份像素報告逐位元相同，SHA-256為2dab43c45aed77d681eaf26647dcbb673eaca5f508f5b3de2dd3339f80a1e317。

來源查證曾發現gen_trans.py:Z的1007與batches-merged.json的1415仍保留舊值，會在再生成時覆回漏字文案。依020補充READY同步兩值後，抽象語法樹及既有setdefault合併獨立核有效8,825筆zh與正式catalog全部相同，其他8,823筆保持；批次JSON只另補末尾換行。沒有執行會讀messages.json的產生器，也未讀MSG.DBS或新原文。10個既有GUI產物雜湊保持，不重跑未改的正式顯示。workplace/r915/source-sync.json的SHA-256為4eb14d362fe75ba5848ca28066c4fbc7a62fc99e460c655e8bd767b5eb2c1c80。R920只複核此來源同步及回填，已通過並由主代理採認。

R920最終窄複核0／0／0。JSON／Markdown的SHA-256為7488c8163c589b9e37ab5d99bd8721dd1e51e57ed53df7b40e08ecfb85c78504及55bb357e7978451b861e13002ac066d687cf0a58c3d2a96b91a84ae0efddfd05。286個本機驗收產物已一次封存，workplace/r915/sealed.json／verified.json的SHA-256為f8d6da4f5025b1e85dc6a8d3d99bbe7b02015d3dbd73220643ff6a92b87ae3e4及bed3573cab32eed02ddc554e9a5f94c82fed1d4a34f34e9c98911d2ec719b149。SDK1,132、原版248及四份字型輸入保持，產物owner1000:1000；快取未納入封存。

## R942 片頭圖片來源與切換

正常冷啟動無玩家輸入至200M，原版及SDK唯讀。只觀察開檔身分、繪圖數值、記憶體摘要與畫格，不讀訊息或答案正文，不注入state、種子、作弊或原版寫入。這是原版無視窗觀察，正式中文前端保持。

| 關係 | 實際證據與等級 |
|---|---|
| 圖片選取 | 0110:2663原始堆疊值為[5B57,60,15]及[5CB6,61,15]，正常DOS開啟MON60.PIC於684960、MON61.PIC於48268719。指令及實跑數值confirmed。 |
| 覆蓋程式定位 | 依[016的R881](016-party-position.md)，VINIT code載入0110:5055；IDA code-only EA B02／C61對應兩個返回，OVR offset為EA+14。原OVR的14..5655與既有code逐位元同；不是resident的15B57／15CB6。 |
| 圖片繪製 | VINIT raw9C5／runtime5A1A把索引及座標送resident4018，rawA48／runtime5A9D組合一至三張後送3921全屏轉送。實測參數15一致；尚未解driver內部槽與codec，不給表格格式猜名。 |
| 完整轉送 | 指定返回5B9B／5AFD且SP為進入SP+2，才保存14張完成圖。37次screen-blit及driver呼叫不能當14張全部轉送。第一組8張、第二組6張；未在返回另存SP，matching資格由探針條件保證。 |
| 呈獻字樣 | 第一組第3至8張可見；最後圖於48267654完成，60M仍是同indexed圖。MON61載入沒有立即換上屏。第二組第一張於73235759完成，翻起的紙頁仍有傾斜的presents，R944以既有完整畫格確認；六張平面覆繪不涵蓋此圖。圖像與開檔／服務參數的關聯為強證據，沒有據此宣稱完整PIC解碼。 |
| 其他圖片文字 | 第二組另見角色扮演敘述。後續MON62／63及TIT1SCRN／TIT2SCRN只有正常檔名與服務數值；沒有讀新PIC正文或把它們列為中文化通過。 |

第一張完整轉送時色盤仍黑，2M才見廠商畫面，不把黑圖當漏譯。0..60M字串／單字入口各0次；至200M已各125／2325次，不把後段也稱為零印字。

觀察開／關100對1MiB記憶體、indexed及palette摘要、Mouse、mode相同，300份a000／pal／png逐位元同；沒有比較完整CPU寄存器。前30點同初輪無觀察基線，前8張完成圖同上一版窄探針。獨立審查另核469筆IDA匯出定位、48筆相對E8與287筆VINIT原bytes；聯絡圖CRC、尺寸及14張來源像素相同。主代理在FINAL後完整回讀兩份報告，阻擋／應改／建議0／0／0，採認限定來源證據。

重跑入口為[intro_picture_probe.go](../../tools/re/intro_picture_probe.go)，在既有離線、非root、有界`wizardry7-go-ebiten:latest`容器內，從本repo模組根編譯：

```text
go build -p=2 -buildvcs=false -o /out/intro-probe tools/re/intro_picture_probe.go
/out/intro-probe /original /out/new-on on
/out/intro-probe /original /out/new-off off
```

repo、SDK、模組快取與玩家原版目錄唯讀，只有新輸出根可寫；兩個輸出目錄不得已存在。原版PNG、PIC及完整收據留本機。探針只到已量200M標題，事件上限512／10000；遇其他OVR載入停用VINIT特定位址觀察。來源SHA為`144b904d07896aaeef1c4fc63ccf1bbe0e2b3152c6daafac96e767956f241c58`。

工具為Go1.24.13 linux/amd64、IDA Pro9.4，分別沿既有`wizardry7-go-ebiten:latest`與`ida-pro-9.4-idapython:locked-v1`。所有IDA查詢複製新DB，原DB不修改；不採跨段字面callee名稱作語意。地址空間分別是resident IDA EA＝runtime IP+10000h、VINIT IDA EA＝code raw offset、OVR offset＝EA+14、VINIT runtime IP＝EA+5055h。

| 輸入／收據 | SHA-256 |
|---|---|
| 原版DS.EXE | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |
| 原版VINIT.OVR | `c1dfd28e8e616070028ac7bd87b5aea5b85d6b6c80224a8a88bbe1774d415e67` |
| 原版MON60.PIC | `ad299595ac2bb31318bd00182db920d78e2dad630b973bf9199b55611cfc9e06` |
| 原版MON61.PIC | `ab6705f244df4ff910192203cebfbcfc53100ea574239c08a02f432566a7bd40` |
| 原版VGA.DRV | `f064349ce5dc694b26ed58615262d1bad893e59971bf259c0a74bdea9ee27a70` |
| 原VINIT.code.i64 | `98b9f93fd55219dfbe327ebd676ae0bac595459cdcbc119e30dde070ae71600a` |
| 原DS.EXE.i64 | `477f73a17b7abcaf3eef70d9baae16c77d9730f3bf7e0b4c0c9c9ed91c612412` |
| r5實際觀察binary | `5de81cd9ddea440c87c02c1e6f27299d53c40dc5ca7ede91e832c93318c642c2` |
| probe-validation-r5.json | `cccd96b13c65dce7586fe5760ca335313e0779dc9ff96746517bb1a915eb67d1` |
| integrity.json | `b11cf095580f9784089d02378e31eddce21668f4583397336b25c0aa16e1e734` |
| review-source.json | `d10a12b55a24c59f6d7bd30b5381839bf40b86a51cadadff9e125f1340db9412` |
| review-source.md | `fe16394cb3c6a504ae4c86b90ae2c97b50ec4c3b4554ad5236bfe3eda9e515d7` |

收據根為workplace/r942，程式基線3581335、SDK基線d16f03a2的受控隔離工作樹。原版248檔、SDK1132檔保持；七個scratch目錄無存檔，擁有權及root-owned／誤建md目錄異常0、自有容器0。r4／r5退出0採主代理工具終止回報，observation／validation本身不含退出碼。公開來源與r5逐位元同且可建置；公開路徑重建binary為`f973db7d8c5892b81e959012767f77c9db18473f98dea61ef5f48b0eb640c017`，未重跑本觀察，100／300／14只歸r5實際binary。

本輪1216個技術檔已一次封存，closed-inputs.json SHA為`46f3e663d6707fe7708c6d048bb7224076dff16e3187275d3ac0942823671429`，排除Go建置快取及後續GitHub中繼資料。封存時原版、SDK、公開探針來源與FINAL報告保持，root-owned及自有容器0。

來源路線已可重播，正式覆繪仍未實作。下一閘門是逐張呈獻字樣的精確矩形、可恢復背景、撤除時機與原生中文字格；經DRAFT→獨立READY審查後再接線，驗兩字型／兩倍率的正式GUI與同狀態原版。圖片內的敘述也須補齊，不以catalog通過代替圖片文字驗收。

## R943 六張呈獻覆繪定版

2026-10-04，使用者選定[甲的周邊紋理補底](../../workplace/r943/intro-picture-comparison-v2.html)，不採平滑色帶。兩種外觀都把「呈獻」畫在原圖片字帶內；底色為新繪近似，不是已證實的原版無字底圖。正式契約為[025圖片字樣規格](../spec/025-intro-picture-labels.md)，初稿DRAFT，未接正式前端。

只處理R942第一組第3至8張完整轉送的六份indexed／palette。逐張核對原版座標格線，修正初稿遺留英文筆畫的左界。已選甲來源為`workplace/r943/prototype_captions_v6.py`；同奇偶周邊樣本按RGB通道擬合色階，加入最近樣本的紋理差，再以原生16×16 Unifont畫「呈獻」。六張、兩倍率共12張候選，各有原版／候選並列。這是離線可丟棄覆繪，不是GUI畫面；沒有新執行原版或注入狀態。

乙使用相同完整字帶及字模，只以擬合色階產生平滑底色。獨立標準函式庫PNG解析、CRC、來源雜湊、射線字帶判定及固定版HEX筆畫核對兩種候選：各12張、4,992,000像素中，批准字帶及字格外變更0；各核對1,920個字模前景像素與來源相同。原版RGB總和≥90且最大最小通道差≤24的中性色像素，兩種候選皆變更0；此數字不等於完整品牌oracle。甲六張3倍並列已由主代理逐張檢視，未見英文殘筆。完整字模、兩字型與GUI仍須正式驗收。

| 本機來源／收據 | SHA-256 |
|---|---|
| 已選甲`prototype_captions_v6.py` | `52d4d0229616d6b5b25a0aa070e0ee0df643764d5df749a9a7b02572f22f70af` |
| 甲十二張圖與補底模型`caption-prototype-v6.json` | `d8f755406e94af215cb510e0db9fc3f15b7e0933f3478bfb46b6eb5916df83a0` |
| 甲獨立核對`audit-prototype-v6-closed.json` | `7602a4d4d0284aa88d8a365108112a5d20e9ac4c2d265c267e16d8b29fa4f79f` |
| 乙十二張圖與模型`caption-prototype-v7.json` | `4777cf06119e837911bf5973e3bb3d2779223e1aed36ebcafce97bfd729adbce` |
| 乙獨立核對`audit-prototype-v7-closed.json` | `2b52815649f21b225a91ad23293e28d33f9953963531ea5ca7621369afcbb037` |
| 定版比較頁`intro-picture-comparison-v2.html` | `8681d725185a12d56665eb74fe8459f457b7fd1ea227465cbce5654c3e82b91d` |
| 輸入保持`integrity.json` | `2731b552749390974a2e66e906de16165eba9b66db9517a3b3829f7084ac4cdc` |

可重跑入口為本機`prototype_captions_v6.py`及`audit_prototype_v2.py`。沿既有`wizardry7-go-ebiten:latest`，輸入repo、R942及固定版字型唯讀，`/out`掛新的已存在輸出目錄；非root、離線、有界執行。設定`WIZ7_FONT=unifont`及`WIZ7_UNIFONT_HEX`為既有固定版HEX，候選程式先核字型與六份來源雜湊。逐步準備來源的私人腳本只為保留修改歷程；重跑最終來源即可，不重跑所有失敗草稿。

原版248、SDK1132及R942封存1216全部雜湊保持；原型階段184份私人檔的擁有權皆1000:1000。原型圖、原版、第三方字型及完整收據不進Git。此段記錄定版前置驗證，正式生效結果見下節；圖片內其他敘述與所有交付門檻保持。

## R943 正式呈獻覆繪驗收

2026-10-04。025經獨立規格審查轉READY才實作。正式TSV只含`presents→呈獻`，由`bake_intro_picture.py`產生獨立sidecar；一般烘製入口同步接入，不讀原版MSG。兩套原生字模及11個烘製拒絕案例通過。完成轉送通知共用既有arm，沒有以5AFD當圖片家族，沒有寫原版RAM、indexed、palette或素材。

完整`cmd/wiz7`回歸及合法原圖守門測試通過。最終Go1.24.13 Linux程式SHA為`2755c68b2bd707bd8575ad69e28c8214a7f649e59c33cde4da57c0ef4451f646`，只是中間程式，沒有正式版號或封包。實作獨立審查0／0／0於FINAL後完整回讀採認，兩報告SHA見025規格。初輪GUI保留，修正原因紀錄後以最終程式另建`workplace/r943/final/`乾淨重跑。

正常無輸入冷啟動，沒有注入存檔、direct-entry或作弊，預設虛擬時鐘同步、圖片驗證用靜音。39／41／43／45／47／49M涵蓋六張，60M仍保持MON60末圖，74M為MON61第一張完整圖並撤除中文。英文控制及兩套字型各跑2／3倍，共六個真正GUI程序，均標準關窗退出0；關窗至退出0.0326至0.2035秒，不當音訊裝置保證。主代理實際檢視正式3倍Draw圖片。

32個中英文停止點的原版1MiB記憶體FNV、indexed／palette摘要及滑鼠相同，indexed、palette及comp另比完整bytes，原版圖片indexed與palette另核R942的SHA。這不包含完整CPU暫存器，也不宣稱記憶體逐位元驗證。追加舊正式`450803c`程式的兩倍率控制，16點摘要及上述完整bytes同新英文，英文HD／Draw完整RGB也相同；兩個舊控制正常退出0。

公開[獨立驗收工具](../../tools/re/check_intro_picture.py)直接解析合法原生字型，以封閉射線字帶及定案前模型計算期望，不調用被測Go函式。32點×HD／Draw的26,624,000個RGB像素全部相符；兩字型字格及全部筆畫完整，批准字帶外差異0。十二張Unifont甲案原型與正式完整RGB相同；末圖另核，翻頁後兩路徑完整原圖相同。補底一致只證明選定的新繪近似，沒有升為原版無字底圖證據。

| 最終本機收據 | SHA-256 |
|---|---|
| `build-final.json` | `03c75251f05b9c28d70bdd4ee7d9b108fbf4f5d9e8205e57352c9988d3e03169` |
| `final/validation.json` | `f0407562fe90371c7f17e0574025bbc17f4048fa39dd4423756a2eabb56b852f` |
| `final/gui-results.json` | `177bb7abaca1cbbfe453d96a1fcc0598a0d074f57cc7e5a200a6bc52fcb13a69` |
| `old-new-control.json` | `6c22c8d3a4b36d594d25e15c3a418233d952d33c10cab180624cf257b2dff955` |
| `baker-boundaries.json` | `5c7717cf327ddcb35259a8f1fdf9cbbca9774698441452c0c4b546f486bb014f` |

重跑沿既有離線Docker，原版、SDK、字型及模型唯讀，輸出只掛本輪既存工作目錄；資源與UID/GID照AGENTS。建置以`GOPROXY=off`及既有唯讀模組快取，完整回歸為`go test -p=2 -count=1 ./cmd/wiz7`。合法原圖守門可另設定`WIZ7_INTRO_EVIDENCE`指向R942/on-r5、`WIZ7_INTRO_RESOURCE`指向本輪Unifont sidecar；缺本機輸入時此測試skip，不產生代用品。

本機`workplace/r943/run_gui.py`固定八點及正常關窗，`WIZ7_GUI_OUTPUT`須指向新目錄；`WIZ7_GUI_CONTROLS_ONLY=1`只重跑舊正式英文控制。所有圖片、字型、原版與完整收據留本機。完整像素重跑入口如下，參數改為本輪唯讀輸入及新輸出：

```sh
python3 tools/re/check_intro_picture.py /out/final \
  --prototype /out --original /src/workplace/r942/on-r5 \
  --unifont-hex /font/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 --output /out/new-validation.json
```

`close_checks.py`另核舊新控制、原版248、SDK1132及R942封存1216全部SHA保持，兩字型來源保持，輸出擁有權異常0。025維持READY，未量過渡及其他圖片文字、十一組、錄放、冷存讀、音訊和正式封包仍待；本輪不宣稱全中文化或可玩交付完成。

完整性報告`integrity-final.json` SHA為`e30b7193948bd2c5511d88eacd856e6e740ee3b3c575fdf3ea3b5180649d2452`。結果獨立窄審0／0／0於FINAL後完整回讀採認；另核32點摘要、64張中文及16張英文HD的完整RGB與尺寸，60M保持、74M撤除。結果審查未涵蓋稍後完成的舊正式控制，該16點由主代理另驗，不能併稱審查涵蓋48點。兩報告`review-result.json`／`.md` SHA為`8e6a034c727241286f78000df4bf6eb1d30ebf0299a3ecc78e8a20623a35a055`／`ff739b1621902124cb64d8e3b256f6f85b6cfec23ecf3496f87955f68ccf5c8f`。報告保存前的容器路徑別名錯誤已修正，沒有改產品或重跑成功像素以選結果。

本輪964份私人技術檔一次封存，`closed-inputs.json` SHA `c8e2aa1ad1af7bbfad351a8c779cfa14124520cc29e2d10a6add523881706e60`，排除Go快取及後續GitHub中繼資料。工作根root-owned／誤建md目錄、自有R943容器皆0。封存檔、原型、原圖及字模只留本機，不加入版控。

## R944 翻頁圖片原型與勘誤

2026-10-04，沿#22及[R944目標](../goal/R944-intro-picture-description.md)，以R942已封存的六張MON61完整畫格核對圖片英文。沒有新執行原版、注入state、讀取MSG或手冊；正式前端、catalog、字型清單及025規格保持。作者姓名依[names規則](../names.md)保留英文。

R942把73235759判成「無呈獻」過早。主代理重新檢視完整indexed按原色盤整數放大，以及只供辨字的近似攤平圖，確認翻起紙頁仍有傾斜presents。原始indexed SHA為`0f098eb7bf96d22777424ac0f8bc263534be1441ecc631413d61a2323714411a`，palette SHA為`d689ae136905afa64a9329ddd123426f7e0f3cde6108c034e0249dbab0a9727c`。本檔目前狀態表已修正，WORKLOG追加勘誤；R942原始收據保持，R943六張平面覆繪的74M撤除驗收仍成立，但不證明此頁沒有英文。

| 原版完整畫格步數 | 可見圖片字樣及範圍 |
|---|---|
| 73235759 | 傾斜presents、平面by及作者姓名 |
| 75322912 | 翻起的角色扮演介紹、平面by及作者姓名 |
| 77352859 | 平面角色扮演介紹、by及作者姓名 |
| 100999086 | 平面角色扮演介紹、翻起的by及作者姓名 |
| 103102872 | 紙頁遮住部分角色扮演介紹 |
| 105202085 | 介紹及作者文字消失 |

[比較頁](../../workplace/r944/compare.html)包含兩種私人原型，各六張、兩倍率、兩字型。甲替換紙頁英文，原生15／16字模維持水平並隨紙頁重排；補底按三條完整英文詞帶分別擬合周邊紙紋，是新繪近似，仍有接縫。乙保持原版圖片，書下加「呈獻」「奇幻角色扮演遊戲」「作者」字幕；圖片英文仍在。文字全消失的第六張兩案皆完整保持原圖。紙頁角點及攤平模型為假說，不能當原版幾何或無字底圖oracle；兩案尚未定版，尚未接正式前端。

獨立`audit.py`使用PNG CRC／RGB解析、射線閉合字帶及直接解析固定版Unifont HEX／倚天Big5字模，不呼叫候選的補底或印字函式。48張共19,968,000像素，22,916個原生前景像素全符合，批准補底區域及字模筆畫外變更0；54個圖片參照皆存在。字帶與字格位置仍取候選metadata，只核前景筆畫及字格邊界，補底區內的字模零位元背景沒有獨立期望。此檢查不裁決補底品質、英文殘筆、瀏覽器操作或正式GUI，也不是同狀態中文驗收。主代理逐張檢視甲五張有字3倍圖及倚天翻頁2倍，補底接縫如實保留。

| 本機來源／收據 | SHA-256 |
|---|---|
| `prototype_description_v4.py` | `0c686dcebbb404d153082bcde480c98ebf98eec693e6e4de44ec41c5a0b83009` |
| `prototype-v4-unifont.json` | `ade1692ff59a1f11219332f4ca15fc289bea795441a85f5016bae656928fe74f` |
| `prototype-v4-eten.json` | `a9a332dcb905c8025ba5e1d950545e790a45939431fa274ce66424c5188f71e0` |
| `compare.html` | `2945817e1ac3141ac8f3f4b4d7cf0e8179bf4cb69eb1d94035b3831d5f0ac48d` |
| `audit.py` | `5a6005422eb9b5cec084c172cb192c27f267a881a11de9b29d0ac639c3378355` |
| `audit.json` | `c0351f2fca599f4ad57636c7f2ba8ca40a19ba5ab568f096b40278115b7b7756` |

重跑只需最終候選，不重跑失敗草稿。沿`wizardry7-go-ebiten:latest`的Python3.11標準函式庫，repo、原版畫格及兩字型唯讀，已存在的新輸出目錄掛`/out`，離線、非root、有界執行。候選入口為本機`prototype_description_v4.py --backend unifont`或`--backend eten`；前者設定`WIZ7_FONT=unifont`及`WIZ7_UNIFONT_HEX`至固定版HEX，後者設定`WIZ7_FONT=eten`、倚天資料掛`/eten`。獨立核對入口為本機`audit.py`，依主機相同層次提供R942及字型；輸出檔既有則拒絕覆寫。字型、原圖、候選圖片、完整模型及收據只留本機。

來源保持核對原版248、SDK1132、R942封存1216及R943封存964全部SHA相同，正式程式及審查過的production來源保持，原型輸出擁有權異常0。`integrity.json` SHA為`76f9778bbfb3376569fb097fbf39d126163d120d458a54b1216d8026204544fb`，另核工作根root-owned／誤建md目錄0。翻頁中文字位置依共同決策待選，之後才寫規格並走獨立READY審查、正式接線及正常GUI驗收；十一組、錄放、冷存讀、音訊、其他圖片按鈕及正式封包門檻保持。

原型獨立窄審0／0／0於FINAL後完整回讀採認，另核48份PNG的雜湊、CRC、尺寸及54條參照，不重跑原型、audit或原版。審查採認上述有限字模前景與區外證據，不採認補底品質、完整作者區域或正式實作。`review-prototype.json`／`.md` SHA為`7c3be89d482f2b7c20b578ef543e670b5ebbb46f2f605df5983c84b41a4f98d6`／`9b4be61bf8d949b91a51d3e28c60db93f987406fe1418fcec35b57b0af2fb866`。

199份私人技術檔一次封存，`closed-inputs.json` SHA為`0897ba8705f3fc995b25aaefbfeb0c8533dc76e0c40ce92ecf239485cfd54bee`，排除封存檔本身與後續GitHub中繼資料。封存時原型、審查輸入及來源保持，擁有權異常0，自有容器收尾另核；圖片與完整收據不進Git。

<a id="r959-battle-control-display"></a>

## R959 現行戰鬥控制的固定上屏驗收

2026-10-05，沿[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)及[R959契約](../goal/R959-battle-control-display.md)。只用現行cmd/wiz7非test來源副本、R943完整雙字型資源及[R958](019-catalog-migration.md#r958-forest-transition-wait)固定合法暖前綴，不新增玩家輸入。原版DS.EXE為Archives版，69206 bytes，SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。Docker沿wizardry7-go-ebiten映像、Go1.24.13／Python3.11.2，所有原版與SDK輸入唯讀。

| 配置 | 有限結果 |
|---|---|
| 英文3倍 | 完整960×600 RGB與原版放大相同，覆繪事件0。 |
| Unifont3倍 | 戰鬥、格擋、咒語、使用、逃跑、移動、裝備七控制14原生字模相符。全圖重建、區外差異、字格／墨點交疊均0，按鈕邊界通過。 |
| Unifont2倍 | 七控制14字模及完整640×400 RGB重建相符，區外／交疊0，但各越出原版按鈕下緣2列。合計448字格像素／99墨點，button_bounds_pass=false，未通過。 |
| 倚天3倍／2倍 | 依2倍越框停止擴測，未執行。 |

三程序各82完整原版摘要、78命令、67既有輸入及三筆歷史NoResponse相同，退出0、Scratch空。初3475M、末5453.1M、六word[0,12,5,3,61,99]／面向1保持；取訊息502、印字669、手冊0、詞庫102、29目標旗標false。純英文comp即原始A000；中文字型兩倍率comp相同，只有HD排版檢查發現2倍缺口。293個Shown事件中只對七控制作原始字型獨立核對，其他286事件只重建reported glyph，不外推其原生字型或完整翻譯驗收。

原版按鈕外框39×9、含1像素邊框，七個原座標半開矩形依序為FIGHT(8,165,47,174)、PARRY(8,175,47,184)、SPELL(50,165,89,174)、USE(50,175,89,184)、RUN(230,165,269,174)、MOVE(230,175,269,184)、EQUIP(230,185,269,194)。七鍵均dy0，現行字格Y為原印字Y乘倍率；2倍字型高16超出外框，3倍高16在框內。獨立審查直接解析固定Unifont17.0.05原始字模及原版indexed畫格，核越框墨點依上述次序為8／19／18／11／21／11／11。這是confirmed顯示缺口，不改原版狀態或點擊判定。

首次清單將MOVE目視誤記HIDE，Unifont3程序exit0但wrapper錯誤報缺項。原錯鍵失敗檔保持，補審後資格訂正另存，不重跑已完成Go程序。來源守門另在候選執行前修正：完全停用動態樹段讀取、三筆NoResponse即時核同才Run；第一次未執行候選及所有訂正保留。正式程式、catalog、字型與SDK均未更動。

| 私人收據 | SHA-256 |
|---|---|
| 候選binary | `ecb6a26f45b3c7e945c714a386e64fbb9bccb68714f792e55f5096f9aba35ac7` |
| verified.json | `e932cf8fcd8ae8667a264871b270bafe2acdb99c0887e25c151c6c8e3a8dc703` |
| review-result.json | `21dce7684491baa8187101625c5714a8ab04a735a59420c83d713ec52edb85b8` |
| review-result.md | `6135fdac0cebf8866c6af10c5947f69339192c2c77d389b5de73df674bf9297e` |
| sealed.json，110份產物 | `9c3ab15cea76251bee5910f6d761a76f2f1d9836c3e2b278e38044350a79a072` |

獨立規格、候選、錯鍵補審、英文／3倍oracle及最終有限結果各0／0／0，主代理於FINAL後完整回讀採認。最終0／0／0只表示證據及缺口報告可採，不表示2倍通過。48選定來源含1份批准goal訂正，47其他來源、SDK1132、原版248、12資源、3原始字型及R958舊204封存保持，擁有權異常0。一次封存110份技術檔，不含gocache或後續GitHub中繼資料。

重跑入口在本機workplace/r959：prepare.py、build.py、run_case.py，使用全新已存在輸出目錄，固定fixed-prefix.json並唯讀掛原版、SDK、模組快取及兩字型。check_display.py獨立核畫格，2倍預期拒絕；check_inputs.py核輸入及原版狀態保持。工具及完整圖片、字模、印字日誌僅留本機。只限暖HD／NoCursorWait，未驗真正GUI、七按鈕操作、背景守門或撤除；001／015保持READY，十一組、TSV、錄放、冷存讀、音訊與正式封包仍須完成。
