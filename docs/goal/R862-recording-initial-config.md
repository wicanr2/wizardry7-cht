# R862：錄放的完整起始設定候選

日期：2026-10-03。工作來源：[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。
基線f78f050，上一輪正式存檔接線已推送；[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)仍OPEN。

## 目標

依使用者已選定的單檔內嵌完整設定，補齊013的B2。
沿R819已審純量候選與既有檔案清單契約，建立可核對的音訊、原版Root、持久Scratch、
實際程式及前端資源設定候選；不以舊RuntimeSHA代替現行初始化。

## 範圍與閘門

- 013仍DRAFT，只在workplace/r862做可丟棄候選及自製fixture，正式前端與SDK不改。
- 單檔形式及系統使用者資料目錄不重問。工作階段故障處理問題已提出，尚未收到答覆，不代選。
- 音訊沿正式009的有理H／K、明示IPS及浮點位元、mute與固定參數，不改CPU／PIT。
- Root與Scratch分別保存封閉集合及DOS中繼資料，空Scratch合法；保留被Scratch遮蔽的Root項目。
- 初態擷取沿單一擁有者及pristine邊界；Step0不等於未投遞輸入。
- 不執行原版、不讀MSG／手冊／防拷答案、原版存檔正文或原版state；自製內容不用作假原版。
- 分析、建置、測試僅在既有離線、非root、有界Docker；歷史與正式來源唯讀，只寫本輪研究目錄。

## 完成條件

候選欄位、型別、來源及雜湊公式可逐項回查；自製獨立期望能辨識Root／Scratch、
音訊與資源差異，拒絕缺項、別名或未確認初態。
獨立唯讀審查後回填013 DRAFT與唯一研究入口，驗證後才推送及更新#17進度。
不冒稱正式錄放、三次冷GUI重播、冷存讀還原或封包完成。

## 入口

- [目前狀態](../../CONTEXT.md)
- [013錄放草稿](../spec/013-input-recording.md)
- [022錄放證據](../re/022-input-recording-boundary.md)
- [014持久存檔READY](../spec/014-original-save-storage.md)

## 結果

五組自製競爭測試及Python獨立核對通過，83項來源與69項產物已封存。
R863獨立唯讀審查阻擋0／應改0／建議0，僅批准限定DRAFT證據回填。
正式Go、SDK、catalog與字型未改；013仍DRAFT，#17仍OPEN。
來源與審查SHA見研究022，根收尾入口為workplace/r862/final-check.json。
工作根衛生檢查異常0，自有容器已清理。
