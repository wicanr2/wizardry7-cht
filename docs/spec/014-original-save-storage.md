# 014：原版存檔的持久可寫層

狀態READY；2026-10-03；[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)，錄放整合另見[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
工作入口[R859](../goal/R859-original-save-wiring.md)，現況[CONTEXT](../../CONTEXT.md)。
2026-10-08補[R979正常冷保存與LOAD](../re/025-ui-localization-coverage.md#r979-normal-save-load)：
正常原版保存11,404 bytes，兩次冷LOAD及隊員／位置已量；新冷保存同狀態無覆繪控制組仍待，014維持READY。
R860技術阻擋0，唯一歷史狀態文字應改已修正；依AGENTS第3節轉READY，准許正式接線。

## 玩家問題與已查證事實

使用者要求完成中文化、打包並驗證可玩。R859接線前，正式前端`cmd/wiz7/main.go`未呼叫
`Oracle.SetScratch`；原版保存選單可返回遊戲，DOS寫入只被記帳，沒有存檔落盤。
R859已將獨立持久可寫層接入正式前端，正常保存及冷載入還原須依下方矩陣驗收。

原版輸入為Archives重印版`DS.EXE`，SHA-256
`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
工具Go1.24.13／Ebitengine2.9.9，dosgolem基線d16f03a，音訊增量來源另見
[固定清單](../../tools/docker/dosgolem-audio-manifest.json)。定位採公開Go API與原版Step，
不是IDA線性位址；本案不需新硬體或存檔格式反組譯。

R802私人原型（完整收據僅本機）已證實；報告SHA-256
`67d57de5efcb2b16a64a524e0e00a0a23dff5558061ee9c9a231ea24def19179`，
970項索引`a775d1fdb01f4ad4c60bea180c747d837e46c68b07a267dcc464cf25ed1f2427`，
研究入口[022](../re/022-input-recording-boundary.md)：

- 既有合法暖checkpoint，正常點DISK `(212,166)`、SAVE & RESUME `(160,73)`、SAVE `(100,170)`。
- 在獨立Scratch寫出`SAVEGAME.DBS`，11404 bytes／74次DOS寫入，SHA-256
  `cd3ac3a0bb8f53e5e0c9fe6eebf5f2010be5ab8f63412bc7ffaa1ec11a508738`；原版248檔不變。
- 保存時位置純量`[0,11,5,3,61,91]`，正常前進後`[0,11,6,3,62,91]`。
  此純量僅作已確認位置的觀察，不據此改RAM或存檔。
- 冷啟動以正式Click觸發LOAD後到原版防拷，尚未通過或讀取還原；
  沒有可證明題目及初始條件匹配的合法輸入收據，不盲套Type。
- 私人暖GUI使用NoCursorWait、靜音且資格false，不當正式冷音訊或錄放完成證據。

## 既有通用契約

以本地dosgolem的`docs/spec/009-scratch-writes.md`（READY）、
`oracle/oracle.go`的`SetScratch`及`internal/dos/files.go`為依據：

- `Root`只讀原版資料；`Scratch`獨立可寫，讀取先找Scratch、再找Root。
- 原版DOS建檔在Scratch新建，寫開Root既有檔先複製至Scratch；原版Root不改。
- `SetScratch`在原版Run／輸入／寫檔之前設定；空路徑會退回只記帳，不可作正式可玩預設。
- 原版存檔保留原名、位元組與規則；中文只進顯示，不能用翻譯當檔名／序列化欄位。
- `AllowFileWrites`的004例外僅限Root既存basename，不能解決fresh首次新建SAVEGAME；
  不採把原版Root變可寫或把Scratch設為Root的快捷修法。
- Scratch不是Oracle.State的一部分；正式前端本輪仍冷啟動，不新增Restore入口。

## 接線契約

1. 提供明確`-save-dir`覆寫入口。使用者於2026-10-03選定預設採系統使用者資料目錄，
   原版資料保持唯讀；跨平台精確路徑依下方R859已審READY契約。
2. 啟動先解析與檢查資料、執行檔及存檔層，再LoadWith／SetScratch、安裝觀察與音訊，
   最後才Run。存檔目錄無法建立／不可寫則啟動明確失敗，不能靜默回退只記帳。
3. 存檔層與Root分離，不能覆蓋原版Root／EXE或祖先；拒絕會把寫入導回原版的符號連結。
   目錄寫前確認權限及擁有權。驗收使用唯讀原版掛載與獨立可寫輸出。
   精確守門及現有非一般檔案處置依本規格R859契約，不猜測DOS已提供同等防護。
4. 不解析或修補SAVEGAME內容，不新增自動保存／關窗保存；保存與載入仍由玩家操作原版選單。
   原版建角／其他自然寫檔也走既有Scratch契約，不限於一個指定存檔名。
5. 原版保存或載入期間Step持續前進，OPL串流及GUI所有權契約不重設，
   不把PCM、中文畫格、前端設定或錄放資料寫入原版存檔。
6. 目錄及路徑為本機設定；完整原版存檔不入Git、Issue或可散布封包。
   README與平台啟動入口需清楚寫出可備份位置。封包驗收從實際解開的程式與唯讀素材運行。

## 驗收與未解閘門

內部驗收：自製檔案fixture檢查路徑分離、不可寫、非目錄、符號連結及建立失敗；
已知通用Scratch建檔、讀覆蓋層及Root不變契約不由本函式自己的輸出作oracle。

原版驗收：正式冷啟動、合法正常建隊／探索、原版保存、確認非空原版存檔與Root逐檔不變，
移動到不同位置、關窗、重新啟動、正常LOAD、隊員／位置還原。
至少與同起點關閉中文覆繪的原版正常保存比較檔案及原版狀態，限定實際量到欄位；
不把非空檔或暖fixture成功寫成冷正常路徑通過。
2026-10-08使用者明示代理正常遊玩、題目後附英文答案，依[029](029-clearance-answer-hint.md)
取代先前僅重用收據或使用者操作的限制。只讀目前已選定記錄，保留原版驗證，以正常輸入通過。
不讀整個答案表、改判定或跳關；新收據只留本機，保存與冷載入仍須按上述正常路線驗收。

已有Linux可丟棄路徑fixture，預設位置已選定；R859已明列跨平台路徑契約，
R860規格審查已完成，014為READY；正式接線及冷啟動已驗。正常保存與冷載入還原仍須驗收，不宣稱CONFORMED。

R807明示可寫層的可丟棄路徑原型已封存，18項自製fixture通過；
報告`8c024938548051bf0c454062fe7d51d866c685c86bb4c9e85a9aadac1661514b`，
索引`198a18472ec46e7079c23444465e5e2d9650210e484075d1b46f837060450950`。
原版同名自製來源初末SHA不變，危險目錄／符號連結／硬連結／非一般檔／不可寫項目在探針寫前拒絕。
既有DOS OpenFile會追隨連結，不能假設通用API已守門；Linux UID／Nlink實測，
跨平台及檢查後外部替換未量，不擴張通用DOS策略或把惡意並行修改列為新交付條件。
寫入故障為唯讀descriptor自製注入，不能當作實際磁碟滿／斷電驗收；R807當時尚未選定預設位置，2026-10-03已定案。

## R845：明示目錄的啟動接線候選

入口[研究022](../re/022-input-recording-boundary.md#r845-explicit-save-wiring)、
[goal](../goal/R845-explicit-save-wiring.md)。私有副本沿R807原樣守門，
只加明示save-dir、Prepare有限拒絕及LoadWith後的SetScratch，原型不採任何預設位置。
五CLI自製預檢通過，四例在載入前拒絕；合法目錄Prepare成功後拒載自製非MZ。
SetScratch成功載入後分支只編譯、未動態量到；沒有原版Run、GUI遊玩或新的保存／載入證據。
SDK索引55及實際編譯53來源副本分開核對相同，其他13前端及go.mod／sum同，正式程式未改。
R846限定審查0／0／0，只批准DRAFT回填；Linux UID／Nlink及檢查後替換限制保持。
R845當時的跨平台、預設位置與正式READY門檻未解除；2026-10-03已定案預設位置，R859另定精確契約。正常冷存讀仍待驗，不以CLI通過移植正式接線。

## R847：平台存檔預檢候選

入口[研究022](../re/022-input-recording-boundary.md#r847-save-platform-metadata)、
[goal](../goal/R847-save-platform-metadata.md)。共用路徑檢查與平台中繼資料只在本機候選分離。
Linux及macOS核對UID、權限與連結數；Windows核對擁有者SID、連結數及重解析屬性。
Windows擁有者候選接受TokenUser或TokenOwner，不據此宣布正式平台契約已定案。

R848首審指出DOS既有檔案以O_RDWR開啟，只驗寫入權限不足。
候選已修正：Unix以O_RDWR開啟既有檔後關閉；Windows以GENERIC_READ及GENERIC_WRITE、
OPEN_EXISTING查詢既有檔權限。兩者不讀寫內容、不建立或截斷檔案。
共用路徑另以實體檔案身分拒絕大小寫及硬連結別名重疊，仍在建立目錄之前檢查。

Linux二十項自製案例通過，前十九項結果保持，新增只能寫不能讀的檔案在寫探針前拒絕。
Windows的Wine輔助驗證十六項通過；三項符號連結及一項讀寫權限負案例未具有效測試條件。
權限測試改用預先取得的控制代碼還原，已正常還原、自製內容不變；不把未量案例計為通過。
Linux amd64、macOS amd64／arm64及Windows amd64的獨立候選可編譯，未驗原生平台GUI。
正式前端及SDK未改，仍缺SetScratch；014保持DRAFT，預設位置、正常冷存讀及全部交付門檻保持。


## R853：明示可寫層的成功冷啟動候選

2026-10-03，入口[goal](../goal/R853-save-startup-success.md)、
[研究022](../re/022-input-recording-boundary.md#r853-save-startup-success)。
候選沿R847修正後的守門，Prepare成功、LoadWith成功後，在Step0呼叫SetScratch，再啟動既有音訊及owner。
對照側使用相同唯讀觀察與停止控制，僅不接可寫層；兩側實際編譯的53項SDK來源相同。

英文、靜音、無玩家輸入的Linux Xvfb冷GUI到主選單，260M停止並正常退出。
200k／260M原版RAM、indexed、palette、暫存器及PNG相同，248素材檔逐檔不變。
兩目錄實際沒有自然寫檔，只證明啟動成功及有限不變性，不證明保存、載入還原或寫入行為。
私有觀察器初版誤阻斷正常內部取訊息；修為只計數，不解碼或輸出答案、不改原版流程。

正式前端及SDK未改，014保持DRAFT。預設位置已選定，成功候選須待完整契約READY才移植正式接線。
正常冷保存／LOAD、錄放、人耳、十一組與三平台封包門檻保持。

## 為什麼要做

2026-10-03使用者選定系統使用者資料目錄，讓存檔可以跨安裝與更新保留。
保留`-save-dir`覆寫；原版資料唯讀，原版保存及載入規則不變。
這項定案解除預設位置的產品決策，尚未解除精確路徑、守門與正式接線審查。

## R859：正式接線的精確契約

下列是2026-10-03經R860審查的READY契約，適用Linux amd64、Windows amd64、macOS amd64／arm64。
平台路徑屬公開API契約，候選守門及SetScratch啟動屬已證實的有限工程證據；
正常冷保存／載入還原與原生Windows／macOS仍未知，待實作後驗收，不阻止定稿接線契約。

### 路徑與CLI

| 平台 | 未指定-save-dir時的預設 |
|---|---|
| Linux | 非空且絕對的XDG_DATA_HOME加`wizardry7-cht/save`；未設定、空或相對值時用UserHomeDir加`.local/share/wizardry7-cht/save`。 |
| macOS | UserHomeDir加`Library/Application Support/wizardry7-cht/save`，不使用XDG變數。 |
| Windows | 以x/sys v0.36.0的KnownFolderPath(FOLDERID_LocalAppData,0)取得當前使用者實際目錄，再加`wizardry7-cht\save`，尊重系統重新導向。一般位置在LOCALAPPDATA內。 |

`-save-dir`型別string，空值表示使用預設；非空為玩家明示路徑，可以是相對工作目錄的路徑。
不作波浪號、環境變數或引號的第二次展開；以filepath.Abs／Clean解析為平台路徑。
明示路徑成功時不查詢預設資料目錄；查詢預設失敗、空值或非絕對基底時報`default_save_path`，
不回退遊戲旁、暫存目錄或只記帳。所有平台共用同一套後續守門。
固定應用程式目錄名稱`wizardry7-cht/save`，不含程式版本、架構或翻譯字串，更新不改備份位置。
成功預檢後啟動輸出實際存檔路徑，README列預設及`-save-dir`範例。
原版資料參數`-root`說明改為唯讀素材目錄。

平台依據：

- [XDG Base Directory Specification 0.8](https://specifications.freedesktop.org/basedir/0.8/)：資料目錄、相對變數忽略及新目錄0700。
- [Microsoft SHGetKnownFolderPath](https://learn.microsoft.com/en-us/windows/win32/api/shlobj_core/nf-shlobj_core-shgetknownfolderpath)：當前使用者的實際已知資料夾。
- [Apple applicationSupportDirectory](https://developer.apple.com/documentation/foundation/url/applicationsupportdirectory)：非沙盒macOS的使用者Application Support目錄。本包沒有macOS App Sandbox。
- 固定Go1.24.13的os.UserHomeDir、filepath及本機x/sys v0.36.0/windows/syscall_windows.go的KnownFolderPath。Go API座標，不是CPU／IDA位址。

### 預檢及持久可寫層

typed input為Root、實際選定EXE、已解析存檔路徑；成功output為實體canonical存檔目錄。
使用R847已審守門，不改通用DOS API或原版格式。順序如下：

1. 先將Root、EXE及存檔路徑絕對化；已有祖先EvalSymlinks，再附未建立的尾端。
2. 在Mkdir前，以字面ancestor和既有實體檔案身分檢查雙向祖先重疊。
   存檔不能位於Root內、等於Root或是Root祖先；對EXE做同樣檢查。
   指向其他獨立資料目錄的祖先連結可canonical化，導回Root／EXE的連結拒絕。
3. MkdirAll(0700)只建立需要的目錄，既有目錄不改權限或擁有權。
   存檔目錄必須是當前使用者擁有的實際目錄；Windows拒絕重解析屬性。
4. 盤點Root的一般檔案與EXE的實體FileInfo，不讀取其內容。
   存檔目錄每項都必須是一般檔案；拒絕子目錄、符號連結、裝置等非一般項目。
   每項必須由當前使用者擁有、無重解析屬性、連結數恰為1，且不得與Root／EXE同一實體檔。
5. Linux／macOS擁有者為Stat_t.Uid等於Getuid；Windows為擁有者SID等於當前TokenUser或TokenOwner。
   既有檔須以O_RDWR或GENERIC_READ|GENERIC_WRITE、OPEN_EXISTING開啟並正常關閉，
   不讀、寫、建立或截斷。Unix另驗owner-write位；Windows拒絕READONLY。
   Unix存檔目錄需owner-write及owner-execute，實際寫入仍由下一項探針核對。
6. 在存檔目錄CreateTemp，寫一個自製byte、Sync、Close、Remove；任一步失敗則啟動失敗。
   清理不完成也不啟動。探針不使用任何原版檔名或內容，不覆寫已有檔。
7. 在任何原版輸入／Run、可推進回呼、音訊及機器owner啟動之前，
   LoadWith成功後立即SetScratch(成功目錄)。LoadWith失敗時不執行原版，可能留下新建空資料夾。
   正式前端冷啟動，不新增Restore入口；SetScratch不保存在Oracle.State內。

有限拒絕碼沿R847：`root_path`、`exe_path`、`save_path`、`overlap`、`mkdir`、
`not_directory`、`platform_metadata_unknown`、`directory_owner`、`root_inventory`、
`exe_stat`、`entries`、`entry_stat`、`nonregular_entry`、`entry_owner`、`root_file_alias`、
`multiple_hardlinks`、`entry_not_writable`、`directory_not_writable`、`initial_write_failed`。
新增`default_save_path`；失敗顯示繁體中文存檔設定錯誤並附有限碼，程序非零退出。
不輸出原版存檔內容或把存檔路徑寫入公開Issue收據。

檢查後的惡意並行路徑替換不在契約範圍內；不把一次預檢當原子安全沙盒。
原版Run期間的磁碟滿、外部修改或裝置故障依既有DOS錯誤契約，不在本次改寫通用DOS服務。
本修正只恢復持久存檔能力，不新增自動保存、即時保存或關窗保存。

### 驗收矩陣與停止線

| 驗收 | 證據標準 |
|---|---|
| 路徑解析 | 自製獨立期望檢查Linux絕對XDG／空／相對值／缺HOME、明示覆寫、空預設與不可用基底；macOS／Windows在原生平台或Wine輔助檢查精確目錄，不以Linux字串模擬當平台API通過。 |
| 守門 | 沿R847二十項Linux自製測試，含Root／EXE祖先、符號連結、硬連結、非一般項目、只能寫不能讀、建立／探針失敗；增加清理失敗獨立故障注入。Wine原有十六項通過、四負案例未量界線保持。 |
| 正式接線 | 正式前端建置及既有競爭回歸；唯讀248素材，冷啟動到主選單，獨立核對SetScratch設定在Step0及第一個Run／input前；兩停止點RAM、indexed、palette、CPU及PNG與相同來源對照相同。 |
| 原版保存 | 不啟用作弊，中文模式沿合法正常原版SAVE選單，取得非空原版檔與實際DOS寫入；與同起點英文模式比較存檔bytes及原版狀態、素材逐檔不變。暖checkpoint僅作保存支路驗證，不能當冷完整可玩。 |
| 原版冷載入 | 正常移動後關窗，以同一存檔目錄冷啟動LOAD，實際還原隊員及位置；防拷只用已存在且條件匹配的合法收據。沒有匹配收據就停止，不推導或代答，#42保持OPEN。 |
| 封包 | 實際解包、唯讀素材、自訂及預設目錄備份位置、各平台可啟動性另驗；交叉編譯或Linux GUI不當原生Windows／macOS或封包通過。 |

程式路徑、原版選單、自然DOS寫檔、Scratch持久層及重新啟動讀取形成同一垂直鏈。
中文只進畫面；存檔原名、bytes、比較、規則及防拷不變。
READY只批准上述工程接線，不宣稱冷存讀還原或全遊戲CONFORMED。

## R860規格閘門

技術阻擋0，唯一文字應改T1已修正，按AGENTS第3節轉READY。
只批准上述持久可寫層接線，冷LOAD、原生平台與封包仍待驗，不稱CONFORMED。
本機審查快照`workplace/r859/review-r860-at-ready.json` SHA-256
`90a785e8c2eeb93a0a593e679bcf01903f04c2adfb8063342427809dd0ac9916`；Markdown SHA-256
`8f653c40582561dd01fa342671bb8890b5d4a73204f7b00e025a668ab191d3d4`，文字差異核對入口`spec-resolution.json`。
