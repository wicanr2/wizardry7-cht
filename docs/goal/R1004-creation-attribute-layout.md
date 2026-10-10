# R1004：兩倍創角屬性排版

2026-10-10，起點60accf5，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28#issuecomment-6092121590)。R1003早期男性覆繪已完成限定驗收，不重跑其矩陣；R966完整人物檢視頁保持。

1. 先核R1003既存正常配點／頭像頁實圖與同發布indexed、palette、Shown，量出原生字高、列距、數值及箭頭位置。只讀既存本機收據，不先新增原版記憶體讀取或重播。
2. 依共同決策技能做可撤回的版面原型，維持640×400及原生15／16字型；展示文字排版與可點擊箭頭的關係。涉及新顯示配置或輸入映射的分支，取得使用者定案後才進正式規格。
3. 新來源、背景與生命週期只補實作所需的最小證據。DRAFT經獨立唯讀審查到READY後，才接正式程式；原版規則、角色資料、存檔、字型、譯名及錄放schema不改。
4. 定版後驗正常GUI配點、欄位切換、頭像頁及退出；同狀態原版摘要／保存、獨立字模與批准域外保持各自驗收。3倍與R966原路線不受影響。未量條件明列，不將原型稱為正式完成。

目前正式執行檔及資源依[CONTEXT](../../CONTEXT.md)。原型生成器入口為[tools/re/prototype_creation_attributes.py](../../tools/re/prototype_creation_attributes.py)，容器內以`--receipts /src/workplace/r1003 --out /out/<不存在的子目錄>`執行，`/src`唯讀、`/out`只掛本輪`workplace/r1004`。只讀既存收據、生成欄位示意；不執行遊戲、不修圖、不驗正式背景或滑鼠映射。完整原圖、字模與原型只留本機。Docker沿既有GUI image、離線、UID1000及相稱限資源。新目錄沿既有`workplace/`，不新增交付結構。

比較入口為[本機HTML原型](../../workplace/r1004/prototype-r2/compare.html)，可切換兩字型；[Unifont示意](../../workplace/r1004/prototype-r2/compare-unifont.png)與[倚天示意](../../workplace/r1004/prototype-r2/compare-eten.png)放大3倍供辨讀。灰底不代表正式遊戲背景。原型只沿R953已選18像素節奏，不把早期創角來源視為R966同一路。是否採用及移位箭頭的輸入對應待使用者確認，未寫成正式規則。
