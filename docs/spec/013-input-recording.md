# 013：以原版指令數錄製與重播前端輸入

狀態DRAFT；2026-10-02；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
證據[022](../re/022-input-recording-boundary.md)，完成定義[M2](../goal/M2-playable-frontend.md)。
尚未通過READY閘門，不進正式實作；下面格式仍待技術契約定稿與獨立審查。

## 為什麼要做

使用者要求完成中文化並驗證可玩；既有#17要求錄一次正常試玩、重播三次，畫面雜湊一致。
牆鐘時間容易受負載與音訊節流影響，應保存機器goroutine實際交付API的Step與順序。
錄放不新增按鍵語意，不用重播檔替玩家解答防拷，也不改原版規則或存檔。

## 擬定契約

- I1：CLI提供明確的錄製／重播入口，兩模式互斥；初版從冷啟動Step0開始，
  使用同一版dosgolem、原版EXE、初始工作副本與機器設定。現有前端沒有Restore入口，
  本規格不假裝已支持從任意存檔開始。原版玩家存讀檔操作仍是普通輸入。
  錄放時禁止混用內建autoplay、測試clicks-at／help-at或其他未記錄的輸入生產者；
  原版合法輸入只由現有前端或重播佇列交付。
  亦禁止HD牆鐘截圖／hd-at／hd-after-frag等額外有狀態compose來源；
  驗收擷取只能讀已發布快照，不再呼叫compose或更新地圖。
- I2：草稿檔案包含schema、工具／程式版本、原版EXE SHA-256、初始檔案雜湊清單、
  LoadWith機器設定、start_step、結束Step及有序命令列表。
  初始檔案清單保存封閉集合的相對路徑、大小、SHA-256及原版可見的DOS日期／時間與屬性；
  複製前、複製後、LoadWith前皆核對，不含原版位元組、不接受路徑逃逸或symlink。
  前端倍率、譯文／flow／help產物SHA及影響畫面的設定另列；有不同GUI配置不能稱同合成畫面。
  重播啟動前全部核對；不相符即拒絕，不以只驗DS.EXE宣稱遊戲資料與存檔相同。
  既有原型的欄位、編碼與容量已由R815明列於下節；正式schema、完整啟動配置及操作入口尚待定稿。
  未知欄位與型別必須拒載，不沿用私人原型名稱冒稱正式相容。
- I3：命令由機器goroutine在實際API呼叫處記錄，而非GUI收到按鍵的牆鐘時刻。
  每筆保存嚴格遞增seq、before_step、after_step、kind及該kind限定的參數；
  同Step多筆保持seq順序。對Run、Click記錄實際步數效果與有限錯誤類別，
  payload與任意原版Console不寫log；錄放模式的既有逐掃描碼及任意error字串log也須停用，
  錯誤分類不可呼叫含Console的Error()。不能將有錯誤的錄製默認當成功。
  下列命令完整涵蓋目前前端：
  `press_key`（原樣uint8 API引數）、`type`（原樣位元組）、`move_mouse`（原版座標）、
  `click`（現有左鍵預設巨集）、`press_button`／`release_button`（現有右鍵路徑）、
  `run`（原有步數）、`snapshot`（原有發布邊界）。
  不將PressKey的引數拆解為新的make／break事件，也不把Type改送BIOS。
- I4：重播依命令順序在exact before_step執行同一API，再核對after_step與錯誤類別。
  before_step不可倒退、略過未記錄執行或落在前一個Click內部；命令間的推進只來自run／click。
  若不符立即停止重播並明確報第幾筆，不補跑、重擲、偷偷丟棄或篡改事件。
  Run與Click期間的音訊進度回呼及退出沿[009 READY](009-player-audio.md)A3／A7共用契約；
  正式音訊已接線，錄放仍待移植與冷整合驗收。不能因錄放改原版Step或事件順序。
  重播只重建原有snapshot邊界，不增加／合併以讓像素結果通過。
- I5：F1以機器goroutine在安全邊界觀察並承認helpOpen的Step生效，承認後原版Step與輸入不前進；
  GUI要求不會中斷已執行Click或改變已抽出佇列的原有順序。F1／F2是前端狀態，不作原版輸入命令。
  其牆鐘停留時間不參與原版決定性。GUI只讀畫面，使用者在重播期間的遊戲輸入不得混入佇列。
  若要比較F1／F2自身畫面，另保存其面板狀態與對應snapshot，不用原版雜湊冒充前端驗收。
  每份Draw收據須帶發布序號、Step與前端配置SHA，原子綁定同一份快照；
  009正式發布已深拷，R780／R784亦有私人暖flow競爭檢查；
  錄放的Publication／配置SHA／PNG完成確認仍未移植，不把暖旁證當正式冷含音驗收。
- I6：原版工作副本由玩家自備，錄製與重播各使用乾淨副本；允許原版自然寫檔，
  不改原檔、記憶體或原版存檔格式。錄放檔只寫玩家明確指定的本機位置，
  不覆寫已有錄製；權限及有界容量、寫入失敗與退出收尾須先由探針量到。
  完整輸入可能含防拷答案，檔案及payload不能進Git、issue或正式交付包；
  公開測試僅可用自製不含原版資料的命令型別fixture。
- I7：結束點在機器goroutine安全邊界核對結束Step及命令數，再輸出有限hash摘要。
  每次Run／Click的partial error、終止、容量上限或落盤失敗均留下未完成狀態，
  不把被強制kill且沒有結束紀錄的檔案當可驗收重播。
  錄製觀察不改現有Click的NoResponseError後仍snapshot／續行；記錄錯誤與成功驗收資格分離。
  部分推進的錯誤亦保留原有後續順序，不自行補release／Run；重播遇類別或Step不符才中止。
  CPU／DOS完整狀態序列化目前未提供；至少獨立核對全1MiB記憶體、A000、palette、
  暫存器、原版可寫檔案及實際live合成PNG，聲稱範圍限定這些欄位。

## 草稿缺口

1. R763／R776已有實際API與三次重播；R775同輸入的觀察開／關2384對原版雜湊相同。
   已解除該500M讀檔介面路線的觀察未知；有效flow、正常存讀檔及009整合仍待驗。
2. R755／R761冷啟動248檔與設定已固定，R777獨立核對候選binary／profile／header綁定；
   正式CLI與schema仍待READY審查，不將錄製結束後的角色檔當重播初始輸入。
3. 正式009已有200k安全停止與owner／播放器Close共用5秒監督；
   尚缺錄放Writer、Draw確認與成功末註在同一期限內的收尾順序，不另延長退出預算。
4. R815已明列原型命令、參數、錯誤類別及容量，既有R761／R777獨立收據仍有效；
   正式格式與新配置差異需定稿及唯讀審查，不重跑已足夠的格式大表。

## 擬定驗收

真實視窗錄製一條正常操作，包含按鍵、DOS文字、滑鼠移動／左右鍵及F1暫停；
各使用乾淨初始工作副本重播三次。同停止Step、原版完整記憶體／A000／palette／暫存器、
原版可寫檔案及代表live合成畫面一致；至少一筆操作確實改變畫面，不能以空輸入錄製冒充。
不推導防拷答案；涉及防拷只重播使用者既有的本機合法輸入收據。

負向驗收含檔案hash不符、未知命令／欄位、序列／Step不連續、Click跨步、非法引數、
容量／寫入失敗及不完整結尾；拒絕應發生在未投遞該筆不合法輸入之前。
另改一筆可觀察輸入並確認結果不同，證明比較能辨識路線差異。
沒有原版資料就skip，不製造假原版；未量到的路徑維持READY，不由單測升CONFORMED。

## 獨立審查

R732為8阻擋、4應改、3建議，報告`workplace/r732/review.json` SHA-256
`51439dbb79971b113fbdbb55d35e5c2564f5924169c1b8a4e1c920002594c071`。
已依靜態證據澄清F1承認邊界、有狀態合成、API交付時點、檔案日期與封閉集合、
前端配置／Draw關聯、保密log及Click錯誤分流；仍須schema、獨立動態oracle與有界退出探針。
維持DRAFT，不以文字補正代替新一輪審查或原版三次重播。

R739隔離格式原型已通過自製命令與損毀／寫入失敗fixture；候選格式、上限、
原子完成與來源SHA見[022](../re/022-input-recording-boundary.md)。
尚未與原版API、GUI或真實初始資料連接，不轉READY；容量、完整配置、
錯誤分類及跨平台收尾仍需獨立審查。不能用這次格式測試代替I1–I7的動態鏈。

R744獨立格式審查為6阻擋、4應改、2建議，SHA及九個漏拒案例見研究022。
R752新副本已重跑自製攻擊案例，精確欄名、明示null及非法Unicode拒載，
將完整結構解析與成功驗收資格分開，保留I4／I7的錯誤續行契約。
原型完成前同步正文，失敗不發布正式名稱並盡力移除footer；
不能由注入故障測試冒稱真正磁碟／斷電或跨平台耐久。
Type候選canonical寫出base64，讀取仍允等價uint8陣列，長度按解碼位元組計；
正式schema仍待完整初始化配置與實際API序列探針，尚未轉READY。

R758新增Type陣列null漏拒由R760隔離修正；R761獨立21文件及1清理雙故障案例通過。
R755已量到248檔／DOS中繼資料、59來源、3資源與有效純量設定；R761獨立重算相符。
完整CLI／profile／header啟動綁定與負向拒絕尚未接入，不把中繼資料候選當I1／I2完成。
R763私人API／GUI串接完成三次限定路線重播，7152對快照無差；
單筆正常滑鼠變更使295對快照不同，詳見研究022。這條路線只到主選單讀檔介面，
不能當作名字輸入、保存或載入成功。R765獨立重現Draw舊畫面消費新要求，
另確認-exe與標頭雜湊來源不一致；未解除Draw發布關聯、快照所有權及有界交付閘門。
觀察關閉的同輸入A/B、有效flow競爭檢查、正常存讀檔與009整合仍待驗；維持DRAFT。
R772私人發布綁定修正完成三次同配置重播，7152對快照與同發布PNG確認無差；
獨立交錯／錯誤審查另行，不以作者測試將I5或完整規格升為READY。
R775另跑同冷啟動與原輸入的觀察關閉runner，2384對原版雜湊及終態248檔相同；
只解除這條500M路線的觀察副作用未知，正常存讀檔及完整前端退出仍待驗。
R773獨立重現record摘要資格及late-ok確認兩失敗；R776新副本修正後的三次重播與
7152對快照／同發布PNG確認相同。錯誤續行、完成時點及有限定位仍待R777獨立複核，
不能由正常路線皆none而忽略錯誤資格；5秒不代表阻塞I/O或整個退出有界。

R777獨立十三群組與來源／配置綁定核對通過，R773兩個失敗解除，見研究022。
只採認該路線的發布、資格及有限診斷；有效flow、009與完整退出仍待驗，維持DRAFT。

R780以既有合法checkpoint與明示NoCursorWait量到有效中文flow真GUI-race，
F1穩定／恢復同Step、原版248檔未改，詳見研究022；獨立審查R784另行。
這是暖GUI所有權旁證，含一次no_response，錄放資格false；
不升格冷啟動schema、正式預設Click或原版保存載入完成，013仍DRAFT。

R784已獨立核對R780有效flow／Draw、F1與終態同發布確認，暖啟動及資格false界線維持。
首次敘事圖及F1視窗圖沒有同發布確認，不能外推終態ACK；013仍待正式冷接線及正常存讀檔。
R786普通點擊後未開保存／載入選單，不解除此閘門；詳見研究022。

## R814：目前正式接線契約審查

2026-10-02獨立唯讀審查為4阻擋、3應改、2建議；報告SHA-256
`a10682f83ca908bfd73dbecc683c78054cc62ca6e11c55258a70e937e5af7c8f`，
來源索引`06ef1b2c3df85f919db12a71dc3de62f63f6aea668a822ccf4ce64912cfb655e`，
輸出索引`0865ad072959b3e772df02501f9c27594a10e5570efeeea40d1606e3a7acfb5f`。
R776／R777的資格與晚到確認、R779／R784的回呼故障，以及暖flow所有權局部缺口已解除，
不重開既有格式表或原版RE。現有技術阻擋是：

| 項目 | 必須定稿的正式契約 |
|---|---|
| B1 | CLI、schema版本、精確欄位／編碼／容量、有限錯誤碼及profile雜湊公式；候選格式已有證據，不能一直留「待探針」。 |
| B2 | 新正式音訊的H／初始K／mute／Adlib／binary配置；Root與Scratch各自的初始可見閉集合及DOS中繼資料。014的系統使用者資料目錄與精確路徑已經R860審查轉READY，正式接線已完成；錄放對初始可見資料的綁定尚未實作。 |
| B3 | 原API錯誤、文件結構完成與整條鏈的成功資格分開；音訊observer／consumer fault即使原操作回nil，也須否定末註資格，不能偽造原API error。 |
| B4 | 正常結束、關窗、原操作部分錯誤、Writer失敗、同發布Draw確認及owner／Close的末註順序；沿009單一總期限，逾時不發布成功末註。 |

本輪已修正009仍稱草稿、正式淺拷與無有界owner的過時敘述。
R802找到DISK並以獨立Scratch在私人暖GUI正常保存11404 bytes；
正式main仍缺SetScratch，冷LOAD停在未具匹配收據的原版防拷，未驗還原。
正式Draw尚無私人Publication／完成確認機制，錄放模式的逐掃描碼診斷亦須在移植時停用。
NoResponse的Click保留snapshot／續行，原Run錯誤保留直接return，不統一成猜測的錯誤流程。
正式冷含音GUI、三次重播、有效flow、提前關窗與正常存讀還原列入實作後驗收；
它們未量到，不應取代上述可定稿的READY契約。013維持DRAFT。

## R815：已核對原型契約與真正待定分支

下列是固定私人來源的已證契約，供正式規格定稿使用，**不是已採用的正式格式**。
來源為R760的`record.go:22–84,113–374,479–683`與R776的`observe.go:89–136`；
候選、來源清單及唯一新增自製測試的SHA見[研究022](../re/022-input-recording-boundary.md)。
正文不再將已證的型別／容量寫成未知；CLI、格式與操作取捨仍pending，維持DRAFT。

### 已證：原型文件與編碼

原型schema是`wiz7-input-prototype/1`；UTF-8 JSON逐列、每列以LF結束，
順序為一列header、有序事件、最後一列end；end後不允許額外列。
必要欄位全部存在、不得null；拒絕未知／重複／大小寫別名欄位、錯型別、
非法UTF-8、孤立代理碼與額外JSON值。SHA一律64個小寫十六進位字元。

| 物件 | 必要欄位及型別 |
|---|---|
| header | `kind="header"`、`schema`、`program_sha256`、`runtime_sha256`、`exe_sha256`、`frontend_sha256`為字串；`start_step=0`為uint64；`clock`為物件、`files`為陣列。 |
| clock | `cpu_hz`、`irq0_base`、`irq0_every`為非零uint64；`cycle_clock`為bool、`ticks=0`為uint64、`model_sha256`為字串。 |
| files每項 | `path`字串、`size`uint64、`sha256`字串、`dos_date`／`dos_time`為uint16、`attributes`為uint8。 |
| 事件 | `kind`字串、`seq`／`before_step`／`after_step`為uint64、`error_class`字串、`params`物件。 |
| end | `kind="end"`、`event_count`／`end_step`為uint64、`content_sha256`字串、`accepted`為bool。 |

| 命令 | params的精確範圍與Step效果 |
|---|---|
| `press_key` | 僅`key` uint8；原樣PressKey，Step不變。 |
| `type` | 僅`bytes`；解碼後1..1024 bytes。Writer寫base64，Loader亦接受uint8陣列，陣列null拒絕；Step不變。 |
| `move_mouse` | 僅整數`x`／`y`，0..319／0..199；Step不變。 |
| `click` | 同座標範圍；原Click，保留實際after≥before及原錯誤，不硬編為6M。 |
| `press_button`／`release_button` | 僅`button` uint8且恰為1，現有右鍵API；Step不變，不補release。 |
| `run` | 僅`steps` uint64，1..200000；實際差值0..steps，成功時必等於steps。 |
| `snapshot` | 僅`publication` uint64與`frontend_sha256`字串；發布從1逐次+1、配置SHA等於header；Step不變。 |

seq從1連續，before恰等於前一筆after，同Step維持seq順序。
錯誤碼僅`none`、`no_response`、`execution`；後兩者僅click／run可用，run不收no_response。
其他命令只收none。這是原API結果，不能塞入音訊／observer故障。
正式009的有限操作碼與上述三類不同層；如新增區分類別，須明列格式差異並審查。

全檔最多64MiB、單列JSON本體最多1MiB（另加LF）、事件最多1,000,000筆、
初始檔案1..4096筆。path最多260 UTF-8 bytes，Clean後不變，以`/`分隔，
禁止U+0000–U+001F（C0）、反斜線、冒號、絕對與逃逸路徑；依path排序唯一。
DOS屬性僅0x27遮罩；日1..31、月1..12、時0..23、分0..59、二秒欄0..29。
此為實際解析界線，未宣稱完整曆法有效性。

`content_sha256`是header與所有事件原始列bytes加各自LF的SHA，不含end；
驗讀不重排JSON。R776的程式SHA來自實際binary，EXE SHA來自實際選定程式；
RuntimeSHA是R755冷初始化metadata bytes的SHA；FrontendSHA是固定Go1.24.13的
`SHA256(json.Marshal(cfg))`。cfg依序為`Scale`、`IPS`、`Resources`、`LoadWith`、
`Profile`、`SelectedProgram`；物件依結構順序、map鍵排序，無尾LF。
舊59項來源與3資源不包含新正式音訊和Scratch，不可直接充作新的完整啟動配置。
私人CLI只有`-prototype-mode`、`-prototype-dir`、`-prototype-input`、`-prototype-stop`，
尚無正式錄製／重播／停止錄製入口。

### 工程草案與尚待共同決策

R815提供完整啟動配置的具體欄序／型別候選：實際binary與版本、目前來源清單、
Root與Scratch各自封閉集合及DOS中繼資料、LoadWith與有效初態、精確約分H／K、
明示IPS及Float64位元、mute／裝置、音訊固定參數、倍率與三類資源SHA。
Scratch空集合合法，兩個集合必須保留遮蔽前的Root項目；每次重播在Run前核對乾淨私有副本。
R819／R820已補齊`initial_scalars_sha256`的固定型別與唯讀來源候選（見下節）；
正式getter、擷取生命週期及完整啟動配置仍未移植，不能拿舊R755摘要或任意map補洞。
以下只列決策前沿，不採用其中任何尚未確認的方案：

| 分支 | 實際取捨／尚待定稿 |
|---|---|
| 014預設存檔位置 | 2026-10-03已選定系統使用者資料目錄，保留明示save-dir；R859已完成READY契約及正式接線，冷存讀還原仍待驗，見[014](014-original-save-storage.md)。 |
| 正式格式及啟動配置 | 2026-10-03已選定單一錄放檔，完整起始設定內嵌header。新schema版本、CLI名稱及精確欄位仍需定稿，不能冒稱prototype/1相容。 |
| 工作階段故障 | 故障即Abort留下無末註片段，或新schema保存原操作完整但不具成功資格的文件；後者需要獨立有限故障欄與資格算法。兩者都不偽造原API error。 |
| 正常錄製結束 | 事先明示Step／安全邊界上限，或玩家「停止錄製」GUI要求及退出／續玩行為。僅關窗Abort無法滿足正常錄製與三次重播，不作完成方案。 |
| 有界收尾 | 先取得最終同發布Draw確認，再完成owner／音訊Stop及單次Close；資格與剩餘期限有效才Writer.Finish。Sync／Close／發布須納入009同一總5秒監督，不能沿R776先寫成功末註再等Draw。 |

唯一新增自製fixture直接呼叫Abort，證實不改原事件的none、不寫end、不發布正式名稱，
partial及人工改名的缺末註文件皆拒載；沒有原版、GUI或裝置驗證。
自製observer_fault標籤未接實際故障回呼。正常finish的總期限起點仍需精確定稿，
最終Draw亦須納入，已觸發009退出期限不得重設。
現有Finish不具取消或整體5秒保證，候選順序尚未實作與審查；不由這個fixture升READY。

## R819／R820：起始設定純量的有限候選

獨立審查阻擋0、未解除應改0，僅批准DRAFT回填；來源與封存SHA見
[研究022](../re/022-input-recording-boundary.md)。本節不是正式schema或已實作API。
通用私人候選`Oracle.DOSSettings()`直接複製DOS設定；巫術7配置組合只留自製測試，
不加入通用Oracle介面。既有`Steps()`、`Ticks()`與`ClockSettings()`直接讀機器純量。

| 物件 | 候選JSON欄序與Go型別 |
|---|---|
| 頂層 | `schema` string（`wiz7-initial-scalars-candidate/1`）、`step` uint64、`ticks` uint64、`clock`、`dos`、`platform`物件。 |
| clock | `cpu_hz` uint64、`cycle_clock` bool、`irq0_base` uint64、`irq0_every` uint64。 |
| dos | `stored_mouse_x_scale` uint16、`font_full`／`font_half` string、`font_full_bytes`／`font_half_bytes` int64、`drive` uint8、`dir` string、`now`物件。 |
| dos.now | `hour`、`minute`、`second`、`hundredth`，皆uint8。 |
| platform | `dos_date`物件、`dos_time_model` string、`bios_rtc_read_time_supported`／`bios_rtc_read_date_supported` bool、`bios_bda_ticks_observation`／`effective_mouse_x_scale_observation` string。 |
| platform.dos_date | `year` uint16、`month`／`day`／`weekday` uint8。 |

候選雜湊為固定Go1.24.13的`SHA256(json.Marshal(上述struct))`，無縮排、無尾LF，
沿用Go預設HTML跳脫，轉64個小寫十六進位字元；不是跨語言JSON正規化契約。
唯一自製初始化資料為654 bytes，獨立手寫期望及第二程式重算同為
`be68ecaf157861ac3ae4391f0c7951061bf7769dcb3f9dbdac57c6ebb946e735`。

單一擁有者須在LoadWith與設定完成後、投遞任何輸入／Run或可推進回呼之前擷取，
Step與Ticks皆0。私人原型另要求`pristine`生命週期聲明，否則拒絕
`not_pristine_initial_boundary`；PressKey／MoveMouse不改Step，所以Step0不能證明無輸入。
正式呼叫端如何綁定此聲明仍待實作；getter不追蹤歷史、不保證外部並行存取。

已證來源及自製設定：DOS儲存倍率0、字型檔名`END_S13.DAT`／`END_S14.DAT`、
字模bytes 30／15、Drive2、Dir `RICH2`與Now全0；覆寫設定會如實回報。
日期1993-01-01／星期5、時間模型`now_plus_ticks_5493_div_1000_mod_8640000`
及RTC讀時／日期旗標false為固定來源導出的政策（source-bound policy），不是本輪中斷實測。
DOS時間沿來源的Go uint64算術；有效滑鼠倍率及BIOS BDA ticks皆`not_observed`。
儲存倍率0代表自動，不能替換為2；候選不呼叫會讀RAM的PixelWidth／VideoMode，
也不呼叫會清跨日旗標的BIOS時鐘中斷。

自製空白初始化群組競爭檢查通過，包含獨立字面值、20次重讀、無可變別名、
三種邊界拒絕與自製覆寫。未LoadWith原版、Run或做整體RAM／CPU副作用A/B，
沒有GUI、裝置或三次原版重播；不升READY／CONFORMED。
這是有限初始化中繼資料，完整來源／binary／LoadWith、Root與Scratch閉集合、
資源、前端與音訊H／K仍須分別綁定，不能用此SHA取代RuntimeSHA或解除B2全部閘門。

## 為什麼要做

2026-10-03使用者選定單一錄放檔，完整起始設定內嵌檔頭，便於搬移與核對。
仍需玩家自備相符原版資料；原版內容與可能含防拷輸入的錄放檔只留本機。
這項定案解除單檔或附檔的產品決策；完整schema、擁有者與故障收尾契約仍須定稿及審查。

## R862：完整起始設定候選

候選只在`workplace/r862`，依單檔定案組合完整設定；013保持DRAFT。
頂層固定欄序為schema、program、sources、load_with、initial_scalars、selected_program、storage、audio、frontend。
schema為`wiz7-startup-config-candidate/1`，不冒稱沿用prototype/1。

| 物件 | 候選內容與綁定限制 |
|---|---|
| program／sources | 實際binary SHA、GoVersion、GOOS／GOARCH、Main.Version與排序建置設定SHA；來源為相對path、size uint64、SHA。候選只核呼叫端提供的來源集合，正式編譯閉集合仍須建置端交付。 |
| load_with／initial_scalars | 保留LoadWith五項引數與R819純量。正式main目前只傳Adlib=true，其餘零值；DTO驗證不證明宣告來自同一次LoadWith，正式單一擁有者仍須綁定。 |
| selected_program | root／external來源角色及完整File；root宣告須與Root一項逐欄相同。實際-exe選擇與LoadWith讀取仍待接線。 |
| storage | Root與Scratch分開的完整集合，保留同path被遮蔽的Root；空Scratch合法。File沿既有相對path、大小、SHA及DOS日期／時間／屬性型別。重播前逐層核對新增、刪除、內容及DOS中繼資料。 |
| audio | H與初始K為既約分數字串；明示IPS、float64位元16位小寫hex、mute與device_requested。固定取樣／緩衝／淡化／進度與退出參數沿009；device_requested不代表裝置成功開啟。 |
| frontend | Scale int64如實保存flag.Int值，不新增2..6限制；三項停用旗標與xlate／flow／help的requested、loaded、大小及SHA，資源載入結果仍須正式載入器綁定。 |

H為196875030000/65536，未明示IPS的初始K=1，明示時為精確Float64/H。
固定參數為44100Hz、雙聲道16位元、88200 ring、4410預填及最大塊、221淡化、
50ms裝置緩衝、200000進度量與5000ms總退出期限。
DOS日期位元沿現行ModTime的本地時間，不強制UTC，不把Unix唯讀mode當DOS屬性；屬性沿來源0x20。
跨平台須核對實際DOS位元，不符即拒，不改原版時間政策。

SHA沿固定Go1.24.13的`SHA256(json.Marshal(Config))`、struct欄序、預設HTML跳脫、無縮排／尾LF。
配置最多1MiB、來源及每層各最多4096項，路徑沿R815界線。
Build只驗具型別DTO及複製slice，沒有新wire decoder；未知／重複／null欄位解析仍屬B1。
掃描要求呼叫端獨占、不變輸入，未保證抵抗並行換檔；pristine仍是呼叫端聲明。

五組自製競爭測試通過，八種配置差異改變SHA、十五種非法DTO拒絕；
獨立Python重算H／K、IPS位元、檔案SHA、DOS位元及實際自製binary SHA。
沒有原版LoadWith／Run、GUI、裝置或三次重播，不解除正式B2全部閘門。
來源、測試及獨立審查入口見[研究022](../re/022-input-recording-boundary.md#r862-recording-initial-config)。

R863獨立審查阻擋0／應改0／建議0，僅批准限定DRAFT回填。
正式owner與wire解析、故障及正常結束契約仍待完成，不轉READY。

## R864：私人啟動綁定

R862的宣告配置已在私人啟動器綁定實際LoadWith、選定EXE與資源讀取，013仍DRAFT。
擁有者自行建立Oracle並立即SetScratch，擷取與封存設定前不交出指標。
沒有外部既存Oracle或pristine引數；R819內部聲明由這條建構流程保證，Step0不單獨當歷史證據。
封存後只能交出一次；即使Step仍0，交出後再取得配置或再次交出均拒絕。
配置bytes回傳副本，指標及資源狀態仍由單一goroutine管理，不保證外部並行呼叫。

| 綁定 | 候選契約與已驗界線 |
|---|---|
| EXE／LoadWith | 同一實際路徑及Options生成配置；EXE與兩層清單於載入前後核對，原版來源唯讀且不變。未觀察LoadWith內部單次ReadFile buffer hash，不保證抵抗外部並行換檔。 |
| 資源 | 沿實際xlate解析區塊、loadFlow及loadHelp；同次ReadFile的bytes計SHA，loaded由解析結果判定。沒有正式Draw或GUI接線驗收。 |
| 編譯來源 | go list取得本專案／SDK的GoFiles／CgoFiles，加三份模組鎖定檔，共73項。JSON與SHA在建置時嵌入同一binary，不要求交付包帶Go原始碼；研究模式可再核實際來源。第三方及stdlib沒有逐檔列入，以模組鎖定、Go建置資料與實際binary SHA綁定。 |
| 預核 | 可信自製expected配置先驗DTO，初態外的全部欄位於LoadWith前比對；實際純量只能在載入後核對，不符就關閉未交出的Oracle。expected的json.Unmarshal不是正式wire decoder。 |

12個起始綁定案例及1個無專案來源掛載案例通過，使用競爭檢查版候選。
正常／中文／明示24M靜音的初始RAM、indexed、palette、暫存器摘要相同，Step／Ticks皆0。
靜音、選定EXE角色、Scratch、資源及已編譯來源差異在LoadWith前拒絕；純量差異在載入後、交出前拒絕。
無來源掛載模式與帶來源模式的配置逐位元相同，248原版資料SHA不變。

LoadWith呼叫有動態計數；Run與輸入未裝計數器，來源明示只走LoadWith，沒有CPU指令執行或玩家輸入。
完整原型與失敗、SHA入口見[研究022](../re/022-input-recording-boundary.md#r864-recording-startup-binding)。
正式機器goroutine、錄放Writer、Draw、冷含音重播及成功收尾仍未接線，不將這些初態收據稱為完整B2或READY。

R865獨立審查阻擋0、應改0、建議0，只批准上述私人候選及DRAFT證據回填。
正式錄放、正常結束與故障策略仍須定稿及接線，不轉READY。
