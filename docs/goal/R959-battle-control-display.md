# R959 戰鬥控制的現行覆繪檢查

2026-10-05，沿[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)與[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)，回應使用者重查文字及圖片按鈕。原版入口見[R958](../re/019-catalog-migration.md#r958-forest-transition-wait)；82點路線已見七個戰鬥控制，不再探索或新增輸入。

本輪命中規格閘門、畫面對拍及顯示／語意隔離路由，載入相應入口與[001](../spec/001-overlay-no-write.md)、[015](../spec/015-generic-slot-expansion.md)。現行R943雙字型xlate資源都有七鍵，但只有資源存在，尚未證明本次上屏。正式前端沒有暖載入參數；R859暖工具替換舊音訊與SDK，不沿用該候選。

## 私人驗收工具契約 READY（七鍵訂正）

- 使用現行全部非test的cmd/wiz7 Go來源精確副本及本機隔離SDK；不改正式程式、SDK、catalog或字型。私人工具只加合法R733 checkpoint載入、已核對前綴重播、停止守門與末點顯示匯出，並在讀取樹段及其內容之前停用整個動態樹捕捉分支。現行main.go的tree-*.bin讀取／匯出不得在本模式執行，不能只移寫檔位置或略過寫檔；正常印字／compose及正式程式保持。載入在LoadWith後、SetScratch前；重播在所有現行adapter回呼就緒後、音訊初始化及owner啟動前，輸出後return，不替換任何音訊實作。
- 初態3475M、原seed／時鐘／隊伍、/w原版唯讀保持。先初始摘要、Run10M、兩次Click(160,162,NoCursorWait)各Run20M，再精確重播R958的78命令。每個Run最多100k區塊；不新增點擊、文字、游標移動、state／RAM匯出或原版寫入。三筆歷史NoResponse保留，其他錯誤立即停止。
- 每一固定點核R958的完整原版摘要，包含記憶體、indexed、palette、暫存器、鍵盤／IRQ、位置／面向、取訊息與保護數量；不符即停止，不向後補跑。numeric取訊息回呼只讀ID並計數，手冊2500–2574或29待補ID首次出現即停止；詞庫7000–7197只計數，不新增目的指標或正文讀取。正式印字回呼另加私人停止守門，先於任何正文／compound讀取。
- 唯一新增展示輸出在末5453.1M。直接呼叫現行hdAtDump／hdFrame，不改compose、字形、字色、撤除或buttonLabelY。分倚天／Unifont、2倍／3倍四個獨立程序，沿R943既有完整資源；另同工具純英文模式一個程序，關閉flow／特殊catalog，不拿其他畫面當基準。五側各82點原版摘要必須相同，Scratch只核正常生成檔案的路徑與bytes。
- 戰鬥控制只核實際上屏的七個可見英文鍵FIGHT、SPELL、PARRY、USE、RUN、MOVE、EQUIP。不從MSG解新文案，不讀targets.json或手冊／詞庫正文，不新增29正文觀察。末點shown若少任何一鍵，記缺項並停止擴測，不憑資源存在補畫。全部原始印字日誌、圖片及字模只留本機。
- 獨立像素oracle從原版indexed／palette、事件before／after、正式字模及公開座標規則重建批准區域。核英文差異墨點全部撤除、七鍵字模與所選原生字型相同、中文字不裁切、各控制中文字矩形／墨點不交疊、區域外像素差異0。按鈕背景範圍由當輪原圖確認；無法證實邊界或背景就記未量，不用被測PNG反推期望。圖片控制若無印字事件仍保留原圖、記待補，不用中文截圖宣稱全畫面已完整。
- 所有現行選定來源、完整資源、SDK1132與原版248做前後盲SHA，R958204封存檔獨立核同。只是列出輸入保持，不稱工具鏈完整閉集合。未知的圖片控制、動態生命週期、按鈕操作、真正GUI與音訊仍未驗；001／015維持READY，不新增CONFORMED或完整中文化聲明。
- 輸出只在workplace/r959；repo、原版、SDK與模組快取唯讀。沿現有wizardry7-go-ebiten映像、Go1.24.13／Python3.11.2，Docker前綴wiz7-r959-、離線1000:1000、2CPU／2GiB／96pids、900秒及10MiB×3 logs。Xvfb只用於Ebiten初始化，有界父程序finally停止並wait，不冒稱GUI。建置／執行／分析全在Docker。
- 先由唯讀子代理審查至READY才產生候選；結果另審，主代理等FINAL完整回讀核SHA後才採認封存、提交／push及議題回讀。任何英文殘留、裁切、跨區、state差異或守門即記問題，不在本輪猜修正式版面；需要改版面時沿共同決策。全部十一組、TSV、存讀檔、錄放、音訊與封包門檻保持。

本輪只進行目前戰鬥選單的固定上屏檢查。停止時機不改玩家產品行為，未回答R954等級位置、R944圖片翻頁或R951錄放位置的待選。

獨立契約審查0／0／0，樹段捕捉B1已解除；主代理等FINAL後完整回讀兩報告及核SHA採認。JSON SHA `37b5a9693c6adf8c679a09d05e8ead2c40643571d699f88c6ae70d95af51e132`、Markdown SHA `116ed8aae28ad1572917789b53485f1508e16dc13cc254797257f452493540ea`。READY只限本輪私人工具，正式規格及所有交付門檻保持。

七鍵勘誤：首份目視清單把右側第二控制MOVE誤認為HIDE。純英文與Unifont3原版82點均相同，Unifont3程序exit0，wrapper因錯鍵宣告缺項後停止擴測。當輪末點唯一右側第二Shown是MOVE、x238／y177、dy0；既有譯文3007為「移動」，日誌記原版Step5433995917印字。這是驗收清單錯誤，沒有HIDE漏繪的產品證據。保留原失敗與第一次READY快照；本次只補審七鍵訂正、修私人wrapper／checker，再核已有畫格，不重跑通過的Go程序。不改adapter、glyph、原版路線或產品行為。

七鍵補審0／0／0，主代理等FINAL後完整回讀採認；JSON SHA `278c7751131be4e9e00f0aa9e1d409ae15b6d52219eb40ad7bf9d31cbf31235d`、Markdown SHA `19dc6c37367e528e65712728dbea492f1c60f19e5b10dc04c9a504bf4fd15859`。只有錯鍵與退出資格訂正，像素／其他配置仍待核對。

## 驗收結果

三個程序各82原版摘要、78命令、67輸入及三筆歷史NoResponse相同，退出0、Scratch空。英文3倍完整RGB差異0；Unifont3倍七控制14原生字模、全圖及區外差異0，按鈕邊界通過。Unifont2倍原生字模及全圖重建相符，但每個控制越出下緣2列，合計448字格像素／99墨點，邊界未通過，按契約停止其餘兩個倚天配置。

獨立最終審查0／0／0只採有限證據及已證實缺口，主代理於FINAL後完整回讀兩報告並核SHA。沒有按鈕操作、生命週期或真正GUI驗收；001／015仍READY，不能宣稱戰鬥介面全數完成。完整結果及入口見[025](../re/025-ui-localization-coverage.md#r959-battle-control-display)。110份私人檔已封存，SHA `9c3ab15cea76251bee5910f6d761a76f2f1d9836c3e2b278e38044350a79a072`。下一步依既定原生字模與按鈕內排版要求修正置中，十一組及交付門檻保持。
