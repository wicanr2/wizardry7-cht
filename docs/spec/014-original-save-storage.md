# 014：原版存檔的持久可寫層

狀態DRAFT；2026-10-02；[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)，錄放整合另見[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
工作入口[R701](../goal/R701-audio-playability.md)，現況[CONTEXT](../../CONTEXT.md)。
尚未通過獨立規格審查，不准據此寫正式接線。

## 玩家問題與已查證事實

使用者要求完成中文化、打包並驗證可玩。正式前端`cmd/wiz7/main.go`目前未呼叫
`Oracle.SetScratch`；原版保存選單可返回遊戲，DOS寫入只被記帳，沒有存檔落盤。
這是會讓玩家丟失進度的實際缺陷，不是文字覆繪問題。

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

## 待審接線契約

1. 提供明確`-save-dir`覆寫入口。預設位置正待使用者選定：系統使用者資料目錄，
   或遊戲旁獨立資料夾；此處不把建議寫成既定需求。需在回答後定義跨平台實際路徑。
2. 啟動先解析與檢查資料、執行檔及存檔層，再LoadWith／SetScratch、安裝觀察與音訊，
   最後才Run。存檔目錄無法建立／不可寫則啟動明確失敗，不能靜默回退只記帳。
3. 存檔層與Root分離，不能覆蓋原版Root／EXE或祖先；拒絕會把寫入導回原版的符號連結。
   目錄寫前確認權限及擁有權。驗收使用唯讀原版掛載與獨立可寫輸出。
   精確守門及現有非一般檔案處置須以隔離fixture定稿，不猜測DOS已提供同等防護。
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
防拷只用既存且條件匹配的合法收據或使用者正常操作，不讀答案表、不推導或代答。
沒有匹配條件時停止該分支，仍可完成不依賴它的目錄與保存驗證。

已有Linux可丟棄路徑fixture；使用者預設位置、跨平台路徑契約、
正式規格READY審查、保存與冷載入還原仍待完成。014維持DRAFT，不宣稱CONFORMED。

R807明示可寫層的可丟棄路徑原型已封存，18項自製fixture通過；
報告`8c024938548051bf0c454062fe7d51d866c685c86bb4c9e85a9aadac1661514b`，
索引`198a18472ec46e7079c23444465e5e2d9650210e484075d1b46f837060450950`。
原版同名自製來源初末SHA不變，危險目錄／符號連結／硬連結／非一般檔／不可寫項目在探針寫前拒絕。
既有DOS OpenFile會追隨連結，不能假設通用API已守門；Linux UID／Nlink實測，
跨平台及檢查後外部替換未量，不擴張通用DOS策略或把惡意並行修改列為新交付條件。
寫入故障為唯讀descriptor自製注入，不能當作實際磁碟滿／斷電驗收；預設位置仍未選定。

## R845：明示目錄的啟動接線候選

入口[研究022](../re/022-input-recording-boundary.md#r845-explicit-save-wiring)、
[goal](../goal/R845-explicit-save-wiring.md)。私有副本沿R807原樣守門，
只加明示save-dir、Prepare有限拒絕及LoadWith後的SetScratch，原型不採任何預設位置。
五CLI自製預檢通過，四例在載入前拒絕；合法目錄Prepare成功後拒載自製非MZ。
SetScratch成功載入後分支只編譯、未動態量到；沒有原版Run、GUI遊玩或新的保存／載入證據。
SDK索引55及實際編譯53來源副本分開核對相同，其他13前端及go.mod／sum同，正式程式未改。
R846限定審查0／0／0，只批准DRAFT回填；Linux UID／Nlink及檢查後替換限制保持。
跨平台、預設位置、正式READY與正常冷存讀門檻未解除，不以CLI通過移植正式接線。

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

正式前端及SDK未改，014保持DRAFT。預設位置仍待答覆，成功候選不移植正式接線。
正常冷保存／LOAD、錄放、人耳、十一組與三平台封包門檻保持。
