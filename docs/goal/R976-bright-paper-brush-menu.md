# R976 紙頁亮字與主選單毛筆字樣

2026-10-08；沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、
[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[025](../spec/025-intro-picture-labels.md)
及[020](../spec/020-main-menu-labels.md)。

## 使用者決定

使用者採用R975甲的大小與水平加粗，指出紙頁字仍太暗，需要再亮。
附圖用於指出深色文字難辨；字重按文字指示採甲。
隨後指定選單用毛筆體，範圍確認為主選單六個按鈕。
紙頁大小／字重不重問，亮度與主選單具體字樣先用原型確認。

## 最小工作

1. 固定R975最終甲案；只在原字模前景提亮，保留背景、位置、30／32字高與加粗。
   深色字候選改米金248／232／160，紅字改236／72／40。兩字型五頁共10張，與舊甲並列。
2. 主選單六則文案只取正式story-trans.json的指定英文鍵，不讀messages.json／MSG.DBS。
   用現存AR PL UKai TW字面做毛筆楷體樣張，依020的109×9按鈕尺寸測2倍及3倍。
   樣張不冒稱原版或正式GUI畫面，不操作按鈕、不改譯文、輸入或原版資料。
3. 所有來源、版本、SHA與原型工具可重跑。字型及原版美術只留本機；不加入Git。
   核來源前景、字格／筆畫包含、色彩差異範圍及來源保持，經唯讀審查後推送。
4. 展示亮字及六按鈕毛筆字樣供確認；正式020／025增補仍需獨立READY及正常玩家驗收。

紙頁補底條紋、2倍紙頁、正式來源見證／HD／Draw／回退／撤除及Go工具鏈恢復另待。
一般訊息字型不在本輪毛筆範圍，既有R975／R974收據保留。

## 入口與隔離

[紙頁亮度工具](../../tools/re/prototype_intro_paper_brightness.py)、
[毛筆字樣工具](../../tools/re/preview_menu_brush.py)、
[FreeType唯讀字模工具](../../tools/re/font_preview_freetype.py)、
[研究025](../re/025-ui-localization-coverage.md#r976-bright-paper-brush-menu)。
[本機對照頁](../../workplace/r976/compare.html)、[毛筆字型來源](../../font/brush-menu-source.md)。
僅寫`workplace/r976/`。沿固定Python3.13.15映像
`933b46a028fd786c9c3d426ebabc237e29a15912231ea8de576e95f0e4f41a4c`，
離線、UID/GID1000:1000、512MiB／1CPU／64pids，有界逾時與log上限。
repo、R975來源、字型及授權資料唯讀，只掛獨立本輪輸出可寫；掛載先驗形態及UID/GID。
收尾核擁有權、誤建md目錄與自有容器。不改R433、R971或既有字型／素材。

## 原型完成

紙頁10新圖只改28,758前景點，形狀、位置、背景與候選作者區保持。
六按鈕44個原生字模與獨立C API一致，2倍／3倍共12字條缺字、越界0。
22PNG共11,673,386像素，49份審查來源SHA及UID/GID保持。
最終唯讀兩報告及SHA完整回讀採認，展示0阻擋／0應改／1建議；正式READY未放行。
建議看過具體亮字與毛筆樣張後確認，不再重問甲的大小與字重。
