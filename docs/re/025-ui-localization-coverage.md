# 按鈕與訊息的中文覆蓋驗收

狀態：主選單、角色選單及新增角色對話框限定驗收通過，其餘介面重查中；2026-10-04；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
入口：[R915](../goal/R915-ui-localization-audit.md)、[R921](../goal/R921-character-dialog-ui.md)、[R928](../goal/R928-character-roster-text.md)、[R932](../goal/R932-character-menu-audit.md)、[目前狀態](../../CONTEXT.md)。

使用者要求把圖片內的英文按鈕納入完整中文化。本文保存目前正式來源及實際畫面的覆蓋矩陣；譯文存在、原版畫格及舊測試完成數量不直接當作本次介面已驗證。

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
| 冷啟動2M，廠商標誌 | 正式Draw可見圖像；Shown事件0；實際素材檔與繪製呼叫尚未追到 | 無操作按鈕；品牌字樣仍是原版 | 已量畫面，圖片來源待定位 |
| 冷啟動60M，廠商標誌及呈獻字樣 | 正式Draw可見英文PRESENTS；Shown事件0；圖像來源為強證據，確切檔案與呼叫未知 | 呈獻字樣尚無中文覆繪；沒有按鈕操作收據 | 圖片文字缺口，不能以catalog查表通過 |
| 主選單六個標籤 | 正常無輸入272M，六個resident文字覆繪事件；不是圖片內固定字樣 | 六標籤完整；「匯入巫術六」五字；雙字型2×／3×的HD及實際Draw都在按鈕內 | 角色選單與新增角色入口／返回通過，其餘入口未量；防禦分支僅單測，見[020](../spec/020-main-menu-labels.md) |
| 片頭製作名單的未翻譯事件 | 四組中文run.log各80次；75次精確對應15個去低位碼後的正式英文鍵，5次空字串 | 15筆來源的zh與en相同；尚未逐項確認保留內容的分類及實際畫面 | 不當75個新漏譯，也不計全文覆蓋通過 |
| 新增角色對話框標題及確定／取消 | R921、R925正常GUI及固定步數已量；既有resident文字覆繪 | 雙字型兩倍率字格與筆畫在框內；正常取消返回 | 限定通過，[021](../spec/021-character-dialog-labels.md)仍READY |
| 新增角色清單六列 | 完整名字＋性別職業動態文字，返回位置0110:0fcc | 名字依既定規則保留英文；六個性別職業3倍雙字型完整，選取及取消正常 | 六個實例通過；2倍、任意姓名與捲動未完成，見[022](../spec/022-character-roster-text.md) |
| 角色選單六按鈕 | R932正常GUI及resident0110:5686、六個精確原文／座標 | 創建／檢視／刪除／改名／頭像／離開完整；雙字型兩倍率字格及筆畫在框內，離開返回正常 | 限定通過；五個子畫面未量，見[023](../spec/023-character-menu-labels.md) |
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
