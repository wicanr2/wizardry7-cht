# 024 人物檢視性別與選人按鈕

狀態：READY。2026-10-04；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。R936窄補審0／0／0，主代理於FINAL後完整回讀JSON與Markdown採認。SHA依序`a30e95673a281c3cac5819741d0f28897b2fcfb1095f6703f55daf4a81563fdc`／`422268f187e0f1051112611139563849c18bce4de83fbd86062dea0f73b6f998`，本機workplace/r935/review-r936-supplement.json及.md。只批准限定正式實作，尚未驗收。證據入口：[025的R935](../re/025-ui-localization-coverage.md#r935-人物檢視與返回)。

## 目的與範圍

使用者要求按鈕與訊息完整中文，包含圖片字樣。R935證實人物頁17個標籤已中文，但性別M仍英文，選人對話框的確定／取消在2倍超框。只補這兩處顯示，不改選人、箭頭操作、角色、原版記憶體、輸入或存檔。字型保持原生15／16列，預設3倍沿既有定案；不縮字。

## 性別輸出證據與契約

- 原版runtime地址空間：0110:3877，near返回0110:6c42，StackWord(1)為M／F，StackWord(2)為字型0，入口游標90,6，返回游標96,6及SP＝入口SP＋2。原版6×6字格位於90,6。M／F冷啟動正常選人路線各量一次。一般NUL字串入口0110:3895沒有收到此M。
- 全離屏差異只有這個6×6字格，M為13、F為11個像素；所有改動後色號1，改動前包含多種底圖色號，不能用固定底色擦除。M探針14個停止點與英文控制相同；女性真正GUI檢視、返回及正常退出已完成，尚無獨立英文固定GUI控制，不能稱女性同狀態或性別中文完成。
- 新增`text/profile-sex.zh-TW.tsv`，五欄表頭精確為`en\tzh\tcontext\tdy\tsource`，UTF-8無BOM、固定順序及欄數，CRLF按行尾處理。英文精確閉集合M／F各一筆，情境只收character.review.sex，來源只收runtime:0110:6c42，譯文依既定性別用字為男／女，dy只收0。拒絕空列、缺欄、重複或未知鍵、控制字元、反斜線跳脫及placeholder；譯文須為非空單一漢字。Go不硬寫中文譯文，不把M鍵塞入一般查表；既有一般M鍵為其他文案。
- 新`tools/re/bake_profile_sex.py`解析TSV，接收既有`bake_font.glyph_char`及原生字高，全部驗證成功後才獨占建立`<outdir>/profile-sex.json`。CLI可指定TSV、輸出及合法原字型路徑，只重建本家族，不讀舊訊息表。`bake_font.bake`在任何寫入前接入本家族的存在檢查與完整建檔驗證；既有xlate／flow／help／roster內容保持，家族輸出已存在即拒絕，不能先寫其他資源再失敗。
- typed JSON只有`schema、size、context、source、source_sha256、entries`。schema固定1、size15或16、情境及來源同上、source_sha256為原始TSV bytes的SHA-256；entries兩筆，各只有`en、zh、dy、g`，g為單一原生字模的row-word陣列。生成端驗UTF-8單字、兩鍵齊全、非空字模、每字高度及位元不越字高，字寬／高均不超18。
- 新`cmd/wiz7/profile_sex.go`的`loadProfileSex`載入選用xlate.json同目錄的profile-sex.json；純英文不載入。檔案上限1MiB，decoder讀取上限1MiB＋1，拒絕未知欄及第二份JSON／尾端資料。核schema、15／16字高與現行一般中文字模後端一致、情境、來源、64位十六進位來源雜湊、M／F各一次、dy0、單一漢字與單字模、原生高度、非空及位元／18像素欄寬。缺檔或任一筆不合格整份停用並記錄，不能留下半套性別資源。載入器只核來源雜湊語法，TSV內容雜湊在生成與驗收核對，不冒稱runtime已驗TSV bytes。
- 只有互動與HD都為3倍時啟用，原生字寬15／16能放入18像素性別欄，且不覆蓋右方連字號。其他倍率保留英文並明確記錄；2倍該欄寬12像素，完整覆蓋仍為交付缺口。
- `overlay.sexPending`為獨立slot，保存來源字碼、入口SP、返回0110:6c42、座標及深複製64000-byte before，不使用或覆寫ov.pend、compound、flow或roster。0110:3877遇上述返回與90,6時先撤除舊性別事件；字碼、font、資源或倍率不符保留原文並記錄。若sexPending已存在，先清舊slot、記錄巢狀，拒絕本次；若一般ov.pend或compound存在，拒絕本次且不改它們。無離屏或長度不符不建立slot。一般0110:3895新入口若發現sexPending，先清性別slot並記錄，然後完整走原一般流程；不修改既有頁、slot或事件。
- 0110:6c42返回hook只消費sexPending一次：先取出並清空slot，再核SP＝入口SP＋2、游標96,6、after完整64000 bytes、非空差異及所有差異在6×6且改動後色號1。錯SP、游標、after、色號或範圍都保留英文、記錄且不留slot等待後次返回。只有通過才建立帶性別標記的textEv，before／after只存6×6、fg1；沿001生命週期等待blit、比對A000、只還原已改筆畫到精確before，原版重畫失配後移除。匹配返回不調一般finish，不改其他pending／頁／roster或原版記憶體。

## 選人對話框置中

- 既有Shown需同時包含023六個角色選單訊號及REVIEW WHO?於67,152、OK於94,168、CANCEL於82,184，dy皆0；所有字模非空，原生高度一致15／16且位元不越界。只調兩按鈕y，其他事件走020／021／023既有邏輯。
- 原版按鈕矩形為63,166,73,9及63,182,73,9。字串寬度為字數×(字高＋2)－2，完整位於自己的矩形，否則不調畫筆。y＝按鈕頂×倍率＋(9×倍率－字高＋1)/2。2倍倚天各上移2、Unifont各上移3；3倍保持。HD及Draw共用同一helper，不改命中區。

## 驗收

- 正式測試覆蓋精確辨識、錯呼叫位置／字型／座標／字碼、巢狀、錯SP／游標、空及錯字模、差異越界／前景錯誤、資源缺失、倍率回退。單測不代替GUI，也不將未量防禦分支標CONFORMED。
- 同一原版及六個正常操作，真正GUI分別檢視THESUS及TEMPEST、返回、正常退出，持久角色資料保持。雙字型3倍補M／F，兩倍率補選人按鈕；英文與修前／修後固定停止點比較原版Mem、indexed及palette皆相同。
- 獨立原字型核男／女及兩按鈕，HD與Draw完整原生字格無裁切；原版性別before／after與獨立探針比對，紋理還原及fg1通過。3倍性別區與2倍按鈕舊／新字格聯集之外差異0；非人物頁、返回及3倍選人畫面保持。
- 缺譯明確記錄；本規格的2倍性別、任意其他人物頁及其他分頁仍未完成。十一組、TSV全面遷移、錄放、冷存讀、音訊及封包交付門檻保持。沒有全遊戲或所有圖片按鈕完成聲明。
