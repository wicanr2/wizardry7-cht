# 013：以原版指令數錄製與重播前端輸入

狀態DRAFT；2026-10-04；[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
證據[022](../re/022-input-recording-boundary.md)，完成定義[M2](../goal/M2-playable-frontend.md)。
尚未通過READY閘門，不進正式實作。目前待審契約以[R951](#r951-production-contract)為準；較早候選只保存證據來源，不決定正式格式。正式介面已選遊戲區外頁尾，格式與相容契約採認尚未完成。

## 為什麼要做

使用者要求完成中文化並驗證可玩；既有#17要求錄一次正常試玩、重播三次，畫面雜湊一致。
牆鐘時間容易受負載與音訊節流影響，應保存機器goroutine實際交付API的Step與順序。
錄放不新增按鍵語意，不用重播檔替玩家解答防拷，也不改原版規則或存檔。

2026-10-03使用者選定單一錄放檔、完整起始設定內嵌檔頭，並選定正常錄製由玩家按鈕停止。
停止要求等當前操作完成後，在安全邊界收尾，不採預先指定步數作正常停止入口。
同日選定收尾後繼續遊戲，後續輸入不納入本次錄製；正常聲音持續執行須另驗。
錄製故障保留未完成片段，不寫成功末註，重播器拒絕把它當成功錄製；不採結構完整但失敗的診斷文件分支。
這些操作與故障保存取捨已定案，技術生命週期及其驗收仍須定稿，不將013直接轉READY或接入正式程式。

2026-10-05使用者選定遊戲區外常駐頁尾，錄製模式增加40px視窗高度，遊戲區與原版滑鼠座標保持。停止錄製後仍在同一視窗繼續遊戲，完成或失敗狀態留在頁尾。

<a id="r951-production-contract"></a>

## R951：正式錄放契約草案

本節收斂B1至B4，待使用者採認及獨立審查；新名稱、CLI與接線均為設計草案。事實基礎為R819、R862、R869、R947及R948至R950的限定證據。檔案型別、解析界線、同次六資源與停止續玩為confirmed候選行為；正式程式尚未具備這條路徑，不能稱正式完成或原版完整parity。

### 入口、相容與初始資料

| 項目 | 待審正式契約 |
|---|---|
| CLI | 新增`-record <本機檔名>`及`-replay <本機檔名>`，互斥；兩者皆空即走現有普通遊玩。沒有固定Step正常停止入口。 |
| 起點 | 冷啟動Step0／Ticks0；不接受既存Oracle、Restore或呼叫端的pristine布林聲明。原版SAVE／LOAD仍是原有玩家操作。 |
| 相容 | 重播要求同一實際binary、SDK／建置來源、GOOS／GOARCH、機器設定、資料及六資源。新正式格式只讀`wiz7-input/1`；拒絕所有prototype／candidate格式，不自動轉換，也不承諾跨版本或跨平台重播。 |
| 設定 | 仍用既有`-root`、`-exe`、`-save-dir`、`-scale`、`-ips`、`-mute`、`-xlate`、`-flow`、`-help`定位本機輸入。header保存完整有效設定，啟動時逐項比對；不默默覆蓋CLI或用header任意絕對路徑讀檔。 |
| 存檔 | 正式錄製沿014的Scratch目錄正常運行及寫檔；header只記錄開始時兩層完整清單。重播需玩家保留的相符初始Scratch副本，以`-save-dir`明示獨立位置。自然存檔會改變該副本，不把錄製後的存檔冒充起始資料。 |
| 錄放位置 | 錄放檔、partial、成功發布及撤回位置全部在Root／Scratch可見目錄之外。建Oracle或Writer前按實際父目錄身份核祖先，不只用路徑字串前綴；拒絕Root／Scratch內、其別名、已存在輸出或與選定EXE／binary／六資源同inode的目標。重播檔也必在兩層之外，避免把錄製新增partial加入DOS可見集合。 |
| 驗收副本 | 三次重播各用乾淨Root／Scratch，保留DOS中繼資料。前後hash與初始閉集合相符；不修改本機原版或捏造存檔。錄放檔不內嵌原版檔案或存檔bytes。 |
| 禁止混用 | 錄放模式拒絕`-autoplay`、`-autoplay-stop`、`-clicks-at`、`-help-at`、`-map-at`、`-hdonly`、`-hdshot`、`-hd-at`、`-hd-after-frag`、`-draw-at`及其有效排程。不能偷偷忽略要求；驗收工具只讀既有發布快照。 |
| 輸入 | 錄製沿現有GUI API次序；重播不接玩家遊戲輸入、前端切換或另一個輸入生產者。正常關窗可取消重播；取消結果明示未完成，不當成功。 |
| 隱私 | 檔案只寫玩家明示位置，可能含防拷答案。命令payload、Console、逐掃描碼、任意error文字與設定的本機路徑不得進公開log、Git、issue或包。有限碼、事件序號、Step及SHA可作本機診斷。 |

### B1：精確wire契約

UTF-8 JSON逐列，以LF結束：一列header、1..1,000,000列事件、一列end。最後事件必為截止snapshot；end後不允許額外列或bytes。全檔上限64MiB包含所有LF，單列JSON本體上限1MiB；完整config也不得超過1MiB，header整列限制優先。必要欄位不可缺或為null；拒絕未知、重複、大小寫別名、錯型別、非法UTF-8、孤立代理碼、額外JSON值及超過16層巢狀。雜湊為64個小寫hex字元；只容許資源未讀到時的明定空SHA。

| 物件 | 欄位及型別，順序也用於Writer編碼 |
|---|---|
| header | `kind="header"` string、`schema="wiz7-input/1"` string、`start_step=0` uint64、`config_sha256` string、`config`物件。 |
| event | `kind` string、`seq` uint64、`before_step` uint64、`after_step` uint64、`error_class` string、`params`物件。 |
| end | `kind="end"` string、`event_count` uint64、`end_step` uint64、`content_sha256` string、`accepted=true` bool。任何false或欠末註文件皆不接受成功重播。 |
| snapshot.panel | `help_open` bool、`help_screen` string、`help_page` int64、`map_open` bool；screen最多64 UTF-8 bytes，不含NUL／CR／LF。 |

命令、引數與Step界線沿R815型別表，以下為明定差異及補全：snapshot的params恰有`publication` uint64、`frontend_sha256` string及上述`panel`物件，其餘七種命令不得帶panel。type的bytes由Writer輸出base64 string，Loader同候選接受uint8陣列；解碼後1..1024 bytes，禁止null及null元素。seq從1連續，before恰等於前筆after；snapshot發布ID從1逐次加一，配置SHA等於header。`none`、`no_response`、`execution`只表達原API結果；non-run／click只允許none，run不允許no_response。任何API錯誤都否定成功末註資格，保留原版後續處置。

| config物件 | 精確欄位與型別 |
|---|---|
| 頂層 | 依序`schema="wiz7-startup-config/1"` string、`program`、`sources`、`load_with`、`initial_scalars`、`selected_program`、`storage`、`audio`、`frontend`。 |
| program | 依序`sha256`、`go_version`、`goos`、`goarch`、`module_version`、`build_settings_sha256`，皆string。實際binary SHA與Go建置資料；GoVersion固定go1.24.13，其餘文字非空。 |
| sources每項 | `path` string、`size` uint64、`sha256` string；1..4096項，依path排序且唯一。建置端產生專案及SDK選定GoFiles／CgoFiles與模組鎖定檔，經`-X`注入同binary。生成路徑表本身納入來源，不遞迴把生成的完整SHA清單當自身輸入。stdlib／第三方依Go工具鏈與模組鎖定，不宣稱逐檔覆蓋。 |
| load_with | `adlib` bool、`dgroup` uint32、`mouse_x_scale` uint16、`font_full` string、`font_half` string。取同一次實際Options；沿目前Adlib=true，其餘零值，不猜補更換規則。 |
| initial_scalars | `schema="wiz7-initial-scalars/1"` string、`step=0` uint64、`ticks=0` uint64、`clock`、`dos`、`platform`物件。 |
| clock | `cpu_hz` uint64、`cycle_clock` bool、`irq0_base` uint64、`irq0_every` uint64。CPUHz非零、cycle_clock=false、IRQ0兩值165000；從實際Oracle唯讀取得。 |
| dos | `stored_mouse_x_scale` uint16、`font_full` string、`font_half` string、`font_full_bytes` int64、`font_half_bytes` int64、`drive` uint8、`dir` string、`now`物件。取同一次Oracle.DOSSettings值副本；不讀RAM或推進執行。 |
| dos.now | `hour`、`minute`、`second`、`hundredth`，皆uint8；保存初始化原值。 |
| platform | `dos_date`物件、`dos_time_model` string、`bios_rtc_read_time_supported` bool、`bios_rtc_read_date_supported` bool、`bios_bda_ticks_observation` string、`effective_mouse_x_scale_observation` string。沿目前SDK模型，兩RTC為false、兩observation為not_observed。 |
| platform.dos_date | `year` uint16、`month`／`day`／`weekday` uint8，依序1993／1／1／5；time_model為`now_plus_ticks_5493_div_1000_mod_8640000`。這是SDK模型，不稱原版實機時鐘parity。 |
| selected_program | `source_role` string及`file`物件。role只能root／external；root時File逐欄等於Root中同項，external只記相對basename及同一實際EXE，不保存任意絕對路徑。 |
| storage | `resolution="scratch_then_root"` string、`metadata_model="dos_archive_mtime_source_policy"` string、`root`及`scratch` File陣列。Root 1..4096項，Scratch 0..4096項；空slice是`[]`，不可為null；遮蔽前的Root項仍保留。 |
| File | `path` string、`size` uint64、`sha256` string、`dos_date`／`dos_time` uint16、`attributes` uint8。path、排序及時間欄界線沿R815；目前來源屬性必0x20，不將Unix唯讀mode當DOS屬性。 |
| audio | 依序`h`、`initial_k` Rational，`explicit_ips` bool、`ips_float64_bits` string、`mute` bool、`device_requested` bool、`sample_rate` uint32、`device_channels` uint8、`sample_bits` uint8、`ring_frames` uint32、`prefill_frames` uint32、`maximum_pcm_chunk_frames` uint32、`fade_frames` uint32、`device_buffer_milliseconds` uint32、`progress_quantum` uint64、`exit_budget_milliseconds` uint32。 |
| Rational | `numerator`／`denominator` string，正十進位整數、既約分數；H=196875030000/65536，未明示IPS時K=1/1，明示則精確Float64/H。IPS位元為16位小寫hex，須對應有限正float64並與H／K完全相符。 |
| audio固定值 | 依序44100／2／16／88200／4410／4410／221／50／200000／5000；device_requested=!mute，只記要求，不證明裝置開啟或人耳正常。 |
| frontend | `scale` int64、`hd_compose_disabled`／`autoplay_disabled`／`test_inputs_disabled` bool、`resources`陣列。scale保存實際flag.Int值，不新加2..6規則；三停用旗標皆true。 |
| resources每項 | `role` string、`requested` bool、`loaded` bool、`size` uint64、`sha256` string；恰六項，順序xlate／flow／help／roster／profile-sex／intro-picture。loaded必requested且SHA有效；未要求為false／false／0／空SHA，讀不到同為未loaded但requested=true。讀到無效內容仍保留bytes SHA。 |

路徑沿R815的260 UTF-8 bytes界線，禁止C0、反斜線、冒號、絕對與逃逸路徑、symlink及非regular檔。DOS日期／時間沿現有ModTime本地時間政策；錄放不更改時區、SDK日期或原版時間。header只做身份綁定，不是完整CPU／DOS序列化。

`config_sha256`及snapshot的frontend_sha256均為固定Go1.24.13 `SHA256(json.Marshal(Config))`，上表struct欄序、預設HTML跳脫、無縮排、無尾LF。BuildSettingsSHA依Key排序建置`[{"key":string,"value":string}]`，同Marshal公式。ContentSHA則對header及全部event的原始列bytes連同各自LF計算，不含end，不重排來源JSON。原始wire的物件鍵順序可不同；config經完整嚴格解析後按固定struct重算，不以原始header列替代配置SHA。

Loader必須完整解析、驗型別／容量／序列／所有SHA／accepted／最後snapshot後才交出文件；以此文件做靜態配置預核，再LoadWith。解析失敗尚未建Oracle；初態不符只允許關閉未交出的Oracle，任何遊戲API或音訊進度回呼皆未執行。沒有公用測試fixture或合法原版時skip，不製作代用品。

### B2：同次啟動及六資源

1. 依014檢查Root與Scratch分離及可寫性，先核錄放所有讀寫位置在兩層外，再記錄兩層初始封閉集合；既有原版素材唯讀。用父目錄及Root／Scratch的實際檔案身份逐層判斷包含關係，拒絕symlink及現存目標別名，不將檔案新增納入錄放後的原版環境。呼叫者不得並行改起始輸入。這不是抵抗外部並行換檔的承諾。
2. 建置清單、目前binary、實際EXE、Options、H／K、有效GUI設定與六資源一次讀取綁定。每份資源以同bytes解析及計SHA，再將所得bundle交Draw；不按路徑再讀一次。側車超過既有1MiB或擷取不完整時拒錄放，不造完整SHA。
3. 單一啟動擁有者建立Oracle，立即SetScratch，在尚未交出指標時擷取DOSSettings、ClockSettings、Step及Ticks，保存初態摘要；禁止擷取前Run、PressKey、Type、MoveMouse及其他輸入。需要在受控SDK補入已受審DOSSettings唯讀getter，更新可重現來源清單，不動主機原始dosgolem或其main。
4. LoadWith前後、設定封存前再次核Root／Scratch／選定EXE；六資源綁定解析實際buffer。配置與slice凍結為值副本，擷取前後原版RAM／indexed／palette／暫存器摘要同。
5. Config只回bytes副本；Take只交出一次。交出後不可重新擷取，即使Step仍0亦拒絕。機器goroutine獨占Oracle，GUI只持深拷快照；不宣稱Step0本身證明沒有輸入。

### B3：原API結果與錄製資格

先記before，呼叫原API，再記after及原有限錯誤；保持原回傳與現有控制流。Click錯誤後仍走原snapshot／續行，Run錯誤仍沿原return，不補release、不補Run或重擲。完成AfterOperation後再核playerAudio、OPL來源與Stream故障；原API回nil也不能掩蓋callback_panic／observer_fault／sink_failure。

錄製資格是單向狀態：原API錯誤、來源／音訊故障、Append／容量故障或Draw失敗一旦發生就失效，之後普通遊戲可依009 fallback Pace繼續，但本檔不能寫成功end或成功名稱。資格故障不偽造event.error_class。停止截止後的普通遊戲故障不回溯否定已完成錄製；明示關窗造成的取樣中止沿009 A7，不只憑最終Incomplete就推定錄製截止前已有故障。

### B4：停止、Draw、Writer與退出

| 狀態 | 轉移與所有權 |
|---|---|
| 錄製中 | 單一機器owner追加實際API與snapshot。GUI只能提出一次停止要求；重複要求不延長期限。首次要求起算5秒單一絕對期限，要求不送進原版。 |
| 承認截止 | 當前原API返回且AfterOperation完成後，在安全邊界記Step、命令數、最後發布ID與配置SHA。之後API不再錄製，不把同一已抽出佇列補錄或重排。 |
| 收尾中 | 最後同Step發布畫面只繪製一次對應快照，等待真正Draw完成通知；ID／Step／配置SHA及完成時間皆須相符。通知保存當次panel，不拿後來面板或PNG檔案存在代替。GUI可顯示收尾狀態。 |
| Worker | 截止後交給單一Writer worker，追加最後snapshot及FinishAt。它獨占Writer；不與阻塞Append／Sync／Close並行Abort，也不由Draw寫檔。 |
| 完成 | payload Sync→end→Sync→Close→同檔案系統不覆寫發布，且每一返回點都未逾期、資格保持及最終Draw有效；才承認成功。正常停止不呼叫Oracle／Stream／Reader Stop或裝置Close。完成後繼續普通遊戲及原音訊，後續操作不改錄製檔。 |
| 失敗 | 保留`.partial`，不新寫成功末註或成功名。若末註已寫但Sync／Close／發布／期限後檢查失敗，限同regular inode嘗試撤回成功名及截去end，保存有限失敗碼。不可清除他人替換檔。 |
| 逾時 | GUI及owner回失敗狀態並恢復未錄製遊戲。正在阻塞的worker仍獨占Writer，系統呼叫返回後按期限拒絕成功並清理；Go無法保證取消任意磁碟系統呼叫。 |
| 關窗 | 已完成錄製保持；未完成錄製取消並保留片段，不把關窗當正常停止。退出沿009首次退出要求的單一5秒總期限停止owner／音訊及單次Close；任何既存錄製收尾只可使用更早有效期限，不能重設退出預算。 |
| 交錯 | 完成、取消及逾時由同一擁有權狀態協調，晚到worker結果不能把已失敗GUI改為成功。關窗或逾時與發布競爭時，worker在每個返回點重新核共同狀態，若已取消則撤回其同inode成功名及end；不假稱檔案系統多步操作為單一原子交易。 |

Writer以0600及O_EXCL建立`<name>.partial`，輸出名及partial都不得已存在，名不得以.partial結束；父目錄必存在、不得有symlink。發布使用不覆寫的同檔案系統link；不支持時回失敗，不退回可覆寫rename。成功後刪partial僅為清理，同inode的完整副本仍不能以.partial名重播。至少事件數1，最後snapshot的Step等於end_step，完成通知必在期限前。

故障清理也可能遭遇第二個IO錯誤：此時回`cleanup_incomplete`，不宣稱已撤回或一定沒有end。任意改名殘留文件、外部換檔、斷電耐久與所有檔案系統不是已證保證；有限實跑已證Draw、ENOSPC、observer及consumer fault的無末註片段被原名及改名重播拒絕。正式驗收須量到撤回成功及清理故障的正確有限結果；不得用這個限制替已知可清理故障偽報成功。

### 中文操作及有限錯誤

2026-10-05使用者採用[頁尾與F1比較](../../tools/recording-control-placement-prototype.html)的甲，遊戲區外常駐40px頁尾。錄製模式視窗為320×scale寬、200×scale+40高，停止及完成／失敗後保留高度；左鍵區x=8..159、y=200×scale+4..35，所有頁尾滑鼠皆不送原版。原型SHA `7211d62a1d38013ce39a96eb74d190366845c396d06b0f45c8c618290868e4fe`。原型字型及背景只作位置示意，正式原生字模、精確佈局及同狀態GUI仍須驗收；不由選定位置直接宣稱013 READY。

正常錄製中顯示「停止錄製」，要求後顯示「錄製收尾中」，成功顯示「錄製完成」，失敗顯示「錄製失敗」。全部用正式字模重建，缺字或控制區放不下必拒絕錄製啟動，不畫不完整中文。重播完成及取消以中文報明；普通遊玩不新增錄製UI。

對外只回有限碼：啟動為`wire_rejected`／`configuration_mismatch`／`initialization_failed`／`recording_output_failed`，執行為`input_step_mismatch`／`api_result_mismatch`／`recording_ineligible`，收尾為`boundary_invalid`／`draw_timeout`／`draw_mismatch`／`writer_append`／`writer_finish`／`stop_timeout`／`cleanup_incomplete`，退出沿009 `owner_timeout`／`close_timeout`。原診斷碼映射到這些有限類別，不呼叫任意Error文字。中文顯示如「錄放檔無效或未完成」、「起始資料或設定不符」、「錄製失敗，已保留未完成片段」；cleanup_incomplete則明示「錄製失敗，檔案收尾未完成」，不能說片段已安全清理。CLI啟動拒絕exit2；已完成後普通退出沿現有exit規則，未完成重播或錄製在退出時為失敗，具體退出碼2。這些名稱是待審正式設計，不改event的三種原API碼。

### READY與實作後驗收分界

本節不以尚未完成實作的動態測試阻止技術契約定稿。READY前須明列全部型別、生命週期、適用證據、已知差異及操作位置，且獨立審查無阻擋並取得尚未定案的產品採認。現在仍DRAFT，沒有正式相容承諾。

實作後以自製fixture驗嚴格wire、完整配置、同次讀取、API錯誤不變、截止、最終Draw、Writer故障／撤回／阻塞、關窗交錯及競爭檢查。再用真GUI錄製非空正常路線、停止續玩且軟體音訊推進、三次乾淨冷重播及無觀察原API控制，核原版RAM／暫存器／indexed／palette／A000、發布Step／Ticks、live畫面與自然可寫檔案。提前關窗必未完成；修改一筆輸入須出現可見差異。裝置與人耳、合法冷LOAD、十一組、全部圖片／按鈕及三平台封包另有門檻，不由錄放局部收據外推。

## 早期契約候選

以下為R951之前的候選摘要及歷史證據，不覆蓋上節的待審正式契約。

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
- I7：錄製結束點在機器goroutine安全邊界核對結束Step及命令數，再輸出有限hash摘要。
  正常停止錄製與退出遊戲分開；錄製停止不呼叫原版或音訊Stop／Close，完成後繼續普通遊戲。
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
| 工作階段故障 | 2026-10-03已選定Abort保留無成功末註的未完成片段，不採結構完整但失敗的診斷檔；原API結果保持，不偽造error。實際故障來源、片段保存與失敗收尾須正式接線驗收。 |
| 正常錄製結束 | 2026-10-03已選定玩家「停止錄製」GUI要求，等當前操作完成後收尾，完成後繼續遊戲，後續輸入不再錄製。預先指定Step及停止錄製即退出均不作正常入口；錄製收尾與持續音訊須另驗。僅關窗Abort無法滿足正常錄製與三次重播。 |
| 有界收尾 | 正常錄製停止先取得最終同發布Draw確認，再完成Writer，遊戲owner與音訊保持；與整個工作階段退出分開。退出才沿009單一總期限停止owner／音訊及單次Close。錄製收尾的Sync／Close／發布也須使用同一絕對期限，不先寫成功末註再等Draw，不重設已觸發的退出期限。 |

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

## R866：成功收尾順序的私人候選

自製最終發布及Draw通知，依序通過同配置／Step／發布確認、允許owner完成、
音訊shutdownPlayer返回且工作階段無故障，才呼叫私人Writer.FinishAt。
整條鏈沿同一stopRequest.Deadline，沒有重新起算五秒；Reader ACK只證明位元組交付。
Writer在Sync、Close及Link返回後核期限，失敗盡力撤回自己建立的成功名稱與末註。
阻塞中的系統呼叫不能由Go直接取消，須由程序監督結束工作程序；不並行Abort正在Finish的Writer。

18個限定案例通過競爭檢查，含完整順序、錯誤／晚到Draw、owner未回報、音訊Close錯誤、
工作階段不具資格、Writer錯誤、已耗期限、晚Sync／Close與兩個實際阻塞子程序。
另一成功案例使用正式Stream.Stop、Reader.Stop及自製裝置，實交付884 bytes淡出並ACK。
兩阻塞子程序的監督預算皆五秒，實際程序退出各5.058與5.146秒；不稱嚴格五秒實機退出保證。

本輪是自製通知及裝置串接，沒有真正GUI Draw、正式owner循環、音訊硬體、原版Run或錄放三次重播。
自製fixture的Abort只用於釋放descriptor，不採用待定的正式故障策略。
寫入前後檢查不保證抵抗惡意並行替換或斷電，未量真實阻塞Link與清理雙故障。
正式停止入口、故障格式、GUI關窗及全鏈移植仍待定稿，013維持DRAFT。
來源、收據及限制見[研究022](../re/022-input-recording-boundary.md#r866-recording-shutdown-chain)。

R867獨立審查阻擋0、未解除應改0、建議0；僅批准限定DRAFT回填。
發布後仍可能在清理或協調器回報時超時，撤回成功名稱僅盡力，不宣稱所有逾時皆無成功名稱。
完整B4及正式GUI／owner移植仍未完成，不能由18個限定案例轉READY。

<a id="r869-recording-config-wire"></a>

## R869：完整設定的單檔解析候選

依單檔定案接合R862配置及R866格式候選，只在`workplace/r869`，013仍DRAFT。
候選schema為`wiz7-input-config-candidate/2`；檔頭依序為kind、schema、start_step、config_sha256、config。
kind固定header、start_step固定0；config完整保留R862型別及固定Go欄序。
config_sha256為Go1.24.13的SHA256(json.Marshal(Config))，預設HTML跳脫、無尾LF。
解析允許等價欄序與空白，配置SHA以型別化Config重建；全檔內容SHA仍核原始列及LF。

檔頭的精確欄名、重複鍵、必要欄位、null、UTF-8及代理碼檢查遞迴涵蓋完整配置。
外層檔頭JSON本體也須不超過1MiB，包含配置與包裝，LF另計；其他容量沿R815。
NewHeader先驗DTO，再確認JSON往返不改動值，拒絕Go編碼器默默替換非法字串。
這只驗自述配置，不以Build的pristine=true證明實際擁有者生命週期。

snapshot仍用既有frontend_sha256欄名，在本候選綁定檔頭config_sha256，涵蓋全配置。
命令順序、原API錯誤、footer的API資格及R866期限候選保持；未補工作階段故障資格或解決發布後撤回缺口。
LoadExpected先完整驗檔，再核呼叫端提供的預期配置SHA，錯誤時回傳空Document；
預期SHA的實際LoadWith／資源綁定仍須正式擁有者供應，不以自述一致代替。

23個限定競爭案例通過，其中一例覆蓋118個必要欄位各自缺省與null，共236次拒絕。
八個獨立Python建檔案例通過，包含中文／HTML跳脫、等價欄序、合法但不同的配置及非法型別。
沒有原版、GUI、音訊、正式冷啟動或三次重播；CLI、schema採認、B3及B4仍未定稿。
來源、失敗紀錄與收據見[研究022](../re/022-input-recording-boundary.md#r869-recording-config-wire)，獨立審查另列。

R870獨立審查阻擋0、應改0、建議0，只批准上述私人解析候選及限定DRAFT證據回填。
實際啟動、正式格式、停止及工作階段資格仍須定稿，不能由23加8案例轉READY。

<a id="r871-recording-startup-wire"></a>

## R871：單檔解析接到實際啟動

私人候選接合R869完整檔案Load及R864實際啟動擁有者，013仍DRAFT。
外部錄放檔先驗全部header／事件／footer，才進存檔守門與資源載入。
檔頭Config只作預期值；實際程式、已編譯來源、選定EXE、兩層檔案、資源及選項由啟動端取得。
靜態不符在LoadWith前拒絕；實際純量不符在載入後關閉未交出的Oracle。
封存配置SHA再與檔頭比對，通過才交出機器，之後逐命令核原版Step。
設定或內容SHA均不能單獨證明實際啟動相符。

本輪只接受自製三筆成功Run，各200000步；這是蒐證子集，不是正式命令或故障政策。
記錄的檔頭直接來自同一實際擁有者，沒有可信設定附檔或人工提供預期SHA。
測試fixture自行完成Writer，不包含最終Draw、正式owner／音訊收尾鏈，不能解除B4。

競爭檢查版27案例通過：七正向、二十負向。
600000步冷執行記錄一次、重播三次及獨立LoadWith控制在0／200k／400k／600k的
RAM、暫存器、indexed、色盤、Step及Ticks摘要相同。
19個負向於LoadWith前拒絕，初態差異於一次載入後拒絕；全部負向交出及Run計數皆0。
沒有玩家輸入、GUI Draw、中文像素、音訊裝置或正常存讀檔，不能當正式試玩三重播完成。
來源、失敗、收據及後續審查見[研究022](../re/022-input-recording-boundary.md#r871-recording-startup-wire)。

R872獨立審查阻擋0、應改0、建議0，只採認限定DRAFT回填。
短段一次記錄三重播不解除正式玩家試玩、B3／B4或交付閘門，013維持DRAFT。

<a id="r875-recording-live-shutdown"></a>

## R875：真實視窗錄放與成功收尾候選

私人副本接合正式機器循環、真正Ebiten Draw、R871啟動綁定與R866收尾。
候選schema為`wiz7-input-live-candidate/3`，初始設定仍為`wiz7-startup-config-candidate/1`。
snapshot新增必要panel物件，欄序為help_open bool、help_screen string、help_page int64、map_open bool。
help_screen最多64 bytes，拒絕NUL／CR／LF；其餘命令不得帶panel，完整巢狀欄位沿嚴格解析。
這是可丟棄的工程候選，不採認正式格式或舊candidate/2相容性。

每個發布ID只採首次完成的真Draw，保存同次Draw實際使用的面板狀態，然後由機器擁有者寫snapshot事件。
重播逐命令核Step及API結果，依該事件還原面板；不是所有GUI畫格或面板開啟時間的錄影。
最終發布沿同一五秒期限取得Draw確認，再放行owner完成、關閉音訊，最後才完成Writer末註。
正常結束使用私人明示Step邊界，不替使用者選定正式停止入口或故障保存策略。

300M正常冷啟動路線，含F1、鍵盤、滑鼠與主選單左鍵，錄製一次及重播三次。
1,471發布的Step／Ticks、合成像素與PNG雜湊一致；38個抽樣發布另核完整RAM、暫存器摘要、indexed、色盤及A000。
獨立控制直接LoadWith／SetScratch與同命令，沒有覆繪或音訊，全部發布的有限原版摘要一致。
原版自然產生PCFILE.DBS 9,368 bytes，三重播與控制檔案SHA相同；不是遊戲進度冷保存還原證據。
改一筆左鍵選項後119發布像素及Ticks不同，完整雜湊差異落在四個抽樣點，證實這條路線可觀察輸入差異。

短段3M競爭版錄製／一次重播／控制通過，15發布，沒有左鍵。
面板格式17例通過，兩接受／十五拒絕；錯誤最終發布及未完成Draw各以真GUI重播拒絕。
後兩例沒有Writer，不能稱實際Writer負向不發布證據。
靜音Reader無裝置交付／ACK，未驗人耳、正式關窗錄製、冷LOAD或完整B3／B4，013仍DRAFT。
來源、完整數值及獨立審查入口見[研究022](../re/022-input-recording-boundary.md#r875-recording-live-shutdown)。

R876封存後唯讀審查阻擋0、應改0、建議0；主代理完整核讀兩份報告。
只採認上述私人候選與限定DRAFT證據，正式停止、故障保存及完整B3／B4閘門保持，013不轉READY。

<a id="r945-recording-stop-owner"></a>

## R945：錄製停止與續玩候選

2026-10-04，沿已定案的玩家按鈕停止、完成當前操作後收尾、停止後續玩及故障保留片段。R875的Stop→owner完成→音訊Close→Writer.Finish只驗整個候選工作階段退出，不能直接接為正常錄製停止。下面是私人技術候選，013保持DRAFT，不宣稱正式CLI／schema已採用。

| 節點 | 候選契約 |
|---|---|
| 要求 | GUI可提出一次停止要求，重複要求不延長期限。要求本身不改原版Step、輸入或音訊；機器擁有者在當前API實際返回後才承認。 |
| 截止 | 先完整記錄當前API的before／after與有限錯誤碼，保持原API回傳；安全邊界鎖定截止Step、命令數及最後發布。之後原路徑照常執行，後續API不再入檔。 |
| 最終畫面 | 以發布ID、Step及配置SHA核對同次完成的Draw，面板保存為snapshot事件。缺Draw、錯配、故障或晚到確認都不完成Writer。候選使用自製通知，不升為正式GUI證據。 |
| Writer所有權 | 機器擁有者交出截止資料後，單一worker獨占Append最後snapshot及FinishAt；不與正在Finish的Writer並行Abort。已停止或失敗的錄製器不接受後續事件。 |
| 續玩與音訊 | 正常停止不調用Oracle、stopRequest、Stream.Stop、Reader.Stop或裝置Close。收尾完成或監督到期後回到未錄製遊戲；音訊仍走原AfterOperation。候選只驗正式軟體音訊，裝置及人耳仍待。 |
| 失敗資格 | run／click的原API錯誤、錄製寫入故障、Draw拒絕及收尾失敗保留未完成片段，不呼叫成功完成。原API錯誤與錄製錯誤分開，不偽造原操作結果。 |
| 期限與限制 | 使用首次要求起算的一個絕對期限。到期可恢復未錄製的擁有者路徑，阻塞worker仍獨占Writer，等待系統呼叫返回後按期限拒絕並清理；不宣稱Go可取消磁碟呼叫、嚴格五秒實機保證、斷電耐久或清理雙故障全涵蓋。 |

候選在44份複製輸入及R875私人SDK上建置；正式SDK仍缺已證實的DOSSettings getter，所以未改正式SDK或假稱已移植。93項選定編譯來源含本專案、私人SDK及cmd測試檔；stdlib／第三方沿工具鏈及模組鎖定，不逐檔宣稱。

15個限定leaf案例競爭檢查通過，包含操作中要求停止、續玩不改已完成檔、自製run／click結果保持、六種Draw拒絕、最後snapshot寫入與Finish故障、阻塞Finish後續玩與20次並行重複要求。真實候選Writer留下三份完整自製檔、十二份無末註片段，獨立Python核配置SHA、原始列內容SHA、事件順序及截止Step。正式軟體音訊Frames從2937增至5873，EndStep從200000到400000，Stop／裝置Close未呼叫；此例沒有OPL寫入，不證明連續有聲音樂。

這不是原版Run、實際停止按鈕、真正Draw、聲音裝置、人耳或三次重播驗收。原型沿candidate/3及原本三類資源設定，尚未把目前正式載入的roster、profile-sex、intro-picture三份sidecar納入同次啟動綁定；完整正式schema、CLI、資源載入、GUI停止及退出競態仍須定稿。來源、失敗及獨立審查入口見[研究022](../re/022-input-recording-boundary.md#r945-recording-stop-owner)，不以15案例直接升READY。

<a id="r946-six-resources"></a>

## R946：六類資源的同次讀取候選

2026-10-04；單檔完整起始設定已定案。R945揭露目前正式前端另載入三份側車，舊三類資源候選不能作最新完整設定。本輪私人配置`wiz7-startup-config-resource-candidate/2`補齊xlate、flow、help、roster、profile-sex、intro-picture六個固定角色；這個名稱只識別可丟棄候選，不定案正式格式或相容策略。

資源的Requested記錄實際讀取要求；沒有譯文或合法空catalog時，三份側車依原載入條件跳過，維持未要求。讀不到保留Loaded=false與空SHA；讀到但解析無效則保留原始bytes的大小與SHA，Loaded=false。成功解析才Loaded=true。解析與SHA使用同一份已讀bytes，不重新按路徑取hash；側車在讀取後換檔仍解析已讀內容。側車超過既有1MiB上限或無法完整擷取時拒絕錄放配置，不聲稱取得完整SHA。

私人Build驗證六個固定角色的數量與順序，配置凍結後複製資源slice；重播候選先比對所有資源紀錄，任一差異都拒絕呼叫後續載入hook。此hook是自製計數，沒有接Oracle.LoadWith。候選依序執行，使用私人全域紀錄與測試換檔回呼，尚未整合R875單一owner或R945停止候選。

22個限定案例競爭檢查通過；15／16字高使用自製字模，五個載入結果與正式loader對照相同。一般譯文map直接沿候選，沒有獨立一般譯文解析oracle，不能外推為兩套發行字型驗收。獨立Python核19份資源紀錄與兩份六資源配置的原始bytes／SHA。完整來源、重跑及限制見[研究022](../re/022-input-recording-boundary.md#r946-six-resources)。完整正式schema、CLI、owner綁定、GUI停止、實際故障、退出競態及三次重播仍待，013保持DRAFT。

<a id="r947-six-resource-owner"></a>

## R947：六資源與單一啟動擁有者候選

2026-10-04；將R946資源紀錄接入R864／R875的私人啟動生命週期。只有私人建構器呼叫LoadWith，使用實際Options產生配置；先核來源、目前binary、儲存清單、選定EXE、音訊及全部六個資源紀錄。可信expected的靜態差異在LoadWith前拒絕；側車不能完整擷取也在建構前拒絕，不以缺SHA繼續配置。

建構後立即SetScratch，指標尚未釋出，才擷取Step0／Ticks0純量與初始摘要。初態不符、擷取改變摘要或清單變化都拒絕並關閉未釋出的指標。完成後Config只回傳bytes副本，Take僅交出一次；即使仍Step0，釋出後也不能再擷取或再次Take。這是同goroutine私人生命週期，未保證並行呼叫。

17個load-only案例通過：六成功、十一拒絕。目前GNU Unifont及本機倚天六份烘製資源都成功載入與綁定；三側車的換空白、缺檔及無效JSON各在LoadWith前拒絕，超限亦拒絕；初態差異則LoadWith一次後、交出前拒絕。中文、英文、相符expected與內嵌來源清單模式的1MiB RAM／indexed／palette／暫存器初始摘要皆等於原始控制。

這裡只綁定已載入資源的紀錄，bundle尚未交給正式Draw，也未整合R945停止續玩。expected只為可信自製fixture，沒有正式不可信wire decoder；實際SDK仍沿私人getter。沒有Run／輸入計數器，無執行與輸入聲明限於已審的load-only源碼路徑，Step0不能單獨證明它。正式CLI／schema、機器goroutine、GUI／故障／退出競態與三次重播仍待，013保持DRAFT。來源與收據見[研究022](../re/022-input-recording-boundary.md#r947-six-resource-owner)。

<a id="r948-recording-gui-continue"></a>

## R948：停止後續玩的GUI蒐證契約

私人候選沿R947同次六資源解析及起始配置，將Trans、數字／序數字模、Flow、Help、Roster、Sex、Picture直接交給目前正式畫面程式，不再讀資源第二次。R945錄製器持有Writer；正常snapshot與原API事件共用相同序列，GUI停止要求不送進原版。

要求在當前API返回後或迴圈開始的安全邊界承認；同Step發布最後畫面，真正Draw確認後由單一Writer工作程序追加最後snapshot並完成。正常錄製停止不呼叫遊戲或音訊Stop／Close；後續原API與發布照常執行，不再寫錄放檔。關窗另沿正式退出期限收尾音訊。重播到錄製最後snapshot才停止重播，不重播錄製完成後的操作。

私人候選的最後Draw使用獨立通知，整段收尾共用首次要求的絕對期限，不另起等待期限。run／click先完成當次AfterOperation，再檢查收尾資格；進入Boundary前再核音訊、來源及串流故障。一次固定的末點判斷同時決定完整原版摘要與Final標記，晚到要求留待下一安全邊界。

全局關窗後可取消一般畫面等待，不補寫snapshot；最後錄製Draw仍走專用Boundary。重播只有消費全部事件且實際Step等於Footer.End才回報完成，提前關窗保留未完成結果。這些來源修正經獨立補審，動態證據依研究022限定，不由來源審查直接宣稱玩家路線通過。

蒐證副本以遊戲區外40像素頁尾放停止控制，字模取同次Flow；遊戲座標仍為原320×200，頁尾滑鼠不送原版。錄放snapshot只保存遊戲區，頁尾另留GUI圖。此位置為可丟棄操作原型，未定案正式介面；故障與逾時資格沿R945限制，未量到不稱READY。一次非空GUI錄製、停止後續玩、三次重播及無覆繪原API控制列為本輪驗收，不用這些收據代替全遊戲可玩或裝置人耳確認。

revision-7已量到真正按鈕停止、35次續玩、三次完整重播及無覆繪控制：1,366發布的Step／Ticks及畫面相符，36點RAM／暫存器／indexed／palette／A000摘要相同。兩個Draw故障保留片段並續玩，提前關窗重播回未完成，短競態GUI通過；獨立另核配置、六資源、事件與PNG。正式getter／CLI／schema、位置、磁碟及來源故障、完整退出交錯與裝置仍待，013維持DRAFT。完整數字、來源及限制見[研究022](../re/022-input-recording-boundary.md#r948-recording-gui-continue)。

<a id="r949-recording-disk-failure"></a>

## R949：實際寫入故障蒐證契約

私人副本只在既有Writer.Write返回後觀察請求／實寫bytes及有限errno，不讀payload、不替換檔案介面或Write結果。先正常磁碟控制，再以容器隔離tmpfs的獨立填充檔耗盡空間，真正GUI及原API保持同一路徑。確認Write實際返回ENOSPC後，沿既有按鈕停止、安全邊界與單一Writer收尾，故障片段不帶成功末註或正式名，後續遊戲及軟體音訊繼續。

輸出片段另存本機，核停止前後SHA及原版前後248檔盲SHA；重播器拒絕.partial及另名的同bytes片段。此驗收不證明其他磁碟故障、立即故障UI、發布後撤回、阻塞呼叫或斷電耐久，不改正式CLI／schema或013的DRAFT狀態。正常控制若不能通過，先修觀察工具，不採認故障結果。

本輪健康磁碟控制正常停止、14次續玩及關窗退出0。故障側真正Write請求285、實寫27 bytes後返回ENOSPC；保留73,728 bytes且尾列未完成，沒有成功末註或正式名。10.2M安全邊界承認停止後再15次操作，片段SHA保持，軟體音訊Frames149737→184969；原檔及改名副本皆拒絕重播，關窗退出2。這是單次實際空間耗盡證據，不替代其他故障與完整退出驗收；來源、獨立核對及重跑入口見[研究022](../re/022-input-recording-boundary.md#r949-recording-disk-failure)。

<a id="r950-recording-audio-fault"></a>

## R950：音訊故障資格蒐證契約

私人副本在9M後各注入一次真正進度回呼panic或取樣yield錯誤，沿既有Oracle捕捉及Stream錯誤路徑，不替換原API回傳或錄製收尾判斷。故障後原Run仍可成功前進，owner讀到觀察／消費故障後須否定成功末註；正常按鈕停止留下片段，後續遊戲照常，原檔與改名副本皆拒絕重播。健康控制必須能正常停止。

只保存有限故障、Step及來源摘要；新片段只在本機核對。真正音訊故障的續玩不要求繼續音訊，仍依009以同一時間軸fallback Pace。關窗造成的中止取樣與關窗前故障分開，以截止時原版OPL觀察碼、playerAudio.fault及Stream摘要判定。本輪不是裝置故障、原版全狀態parity、完整退出交錯或正式格式驗收，013維持DRAFT。

三個真正GUI案例及獨立核對通過。正常控制9.8M停止、15次續玩及關窗exit0；兩種故障9M觸發，故障後各三筆none／完整200k Run。9.6M按鈕停止均回recording_ineligible，72,895 bytes片段無成功末註或正式名，各15次續玩不改檔，原檔與改名副本都拒播，關窗exit2。觀察panic沿Oracle故障碼callback_panic及前端observer_fault，消費錯誤沿Stream／前端sink_failure；故障側沒有繼續合成資格。本輪只解除這兩種已量到的B3路徑，其他故障及正式契約仍待，完整證據見[研究022](../re/022-input-recording-boundary.md#r950-recording-audio-fault)。
