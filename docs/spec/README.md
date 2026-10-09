# 中文化輸出與前端規格索引

新路徑：[012 複合物品名稱](012-compound-item-output.md)（三組正常路線驗收通過；未量防禦分支維持READY）。

每份規格內文保存有效狀態與實測範圍；單元測試不代替正常玩家路徑收據。

| 規格 | 範圍 |
|---|---|
| [001](001-overlay-no-write.md) | 唯讀覆繪 |
| [002](002-bare-slot-key.md) | 純槽鍵排除 |
| [003](003-f1-help.md) | F1 說明 |
| [004](004-frontend-left-click.md) | 前端滑鼠左鍵 |
| [005](005-automap.md) | 自動地圖 |
| [006](006-flow-text.md) | 故事與旁白流動文字 |
| [007](007-native-font.md) | 本機與發行字型原生尺寸 |
| [008](008-english-key-tsv.md) | 英文鍵TSV來源（DRAFT） |
| [009](009-player-audio.md) | 前端音樂播放（READY；正式接線與限定驗收完成，完整玩家路線待驗） |
| [010](010-item-name-corrections.md) | 物品名稱錯配校正（READY；名稱與回歸通過，校正鍵上屏未量） |
| [011](011-message-wording-unification.md) | 同義訊息文案統一（限定CONFORMED；實跑10筆，其他僅靜態驗證） |
| [012](012-compound-item-output.md) | 三組複合物品名稱（正常路線通過；整份READY） |
| [013](013-input-recording.md) | 前端輸入錄製／重播（DRAFT） |
| [014](014-original-save-storage.md) | 原版存檔持久可寫層（READY；精確路徑已定稿，正常存讀及冷載入還原待驗） |
| [015](015-generic-slot-expansion.md) | 一般查表與填槽缺字回退（READY；R822獨立審查通過） |
| [016](016-catalog-control-preservation.md) | 控制碼保留與流文可用性（READY；R830獨立審查通過，正常上屏待驗） |
| [017](017-dynamic-subject-source.md) | 動態主體來源與兩處文意（READY；R832審查通過，主體顯示仍缺） |
| [018](018-hint-numeric-slots.md) | 段落與行號提示數字槽（READY；正式雙後端及R834複核通過，新上屏待驗） |
| [019](019-cross-row-source.md) | 兩列敘述與終止碼來源對齊（READY；R836原型審查通過，正常上屏待驗） |
| [020](020-main-menu-labels.md) | 主選單標籤漏字與按鈕內排版（READY；六標籤雙字型兩倍率、R919實作及R920來源同步複核通過，防禦分支僅單測） |
| [021](021-character-dialog-labels.md) | 新增角色確定與取消按鈕內排版（READY；R925雙字型兩倍率及正常取消驗收通過，防禦分支僅單測） |
| [022](022-character-roster-text.md) | 新增角色清單性別與職業碼（R928初始六列已驗；R991四個移位配對3倍雙字型正常選取／取消限定CONFORMED；其他位置、姓名及2倍仍待） |
| [023](023-character-menu-labels.md) | 角色選單六按鈕中文置中（READY；R932雙字型兩倍率及正常離開通過，防禦分支僅單測） |
| [024](024-character-review-output.md) | 人物檢視性別與選人按鈕（READY；3倍性別與兩倍率按鈕限定像素通過，R941結果複核通過；2倍性別未完成） |
| [025](025-intro-picture-labels.md) | 六張片頭圖片呈獻字樣（READY；周邊紋理補底已接正式前端，雙字型兩倍率及正常GUI限定驗收通過，未知過渡仍待） |
| [026](026-battle-control-labels.md) | 七個戰鬥控制的中文按鈕內排版（READY；正式接線與雙字型兩倍率固定HD通過，GUI／操作／撤除待驗） |
| [027](027-reprint-persistence.md) | 同文重印的覆繪持續性，狀態依正文與R968收據 |
| [028](028-approved-ui-fonts.md) | 定版紙頁亮字及主選單毛筆，兩後端限定正常路線CONFORMED |
| [029](029-clearance-answer-hint.md) | 目前防拷題目後附英文答案（限定CONFORMED；兩後端兩倍率正常GUI與同狀態通過，其餘題目及防禦GUI未量） |
| [030](030-flow-slash-glyph.md) | 流文半形斜線字模（DRAFT；保持原文、schema與排版，新增可用集合待獨立核） |
| [031](031-monster-name-slots.md) | 怪物名字四個16-byte固定格與異譯拒絕守門（DRAFT；23個譯名衝突待定） |
| [032](032-profile-popup-occlusion.md) | 人物法術視窗保留未遮住的裝備中文（限定CONFORMED；五名稱雙字型兩倍率正常路線及原版／像素／Draw已驗，其他裝備／人物切換未量） |
| [033](033-world-identity-labels.md) | 世界人物頁性別與選人提示（R988已量WHO正常選人、移動、取消與再入限定CONFORMED；其他分支仍READY，性別由034承接） |
| [034](034-world-profile-sex.md) | 世界人物頁性別覆繪（M及TEMPEST女性正常路線限定CONFORMED；雙字型兩倍率人物頁、法術視窗、退出／重入及原版／保存已驗，其他路線仍READY） |
| [035](035-backspace-input.md) | 姓名退格字元（Linux3倍兩字型已量路線限定CONFORMED；正常刪字／空欄／F1、錄放及原版控制已驗，其他畫面／平台未量） |
