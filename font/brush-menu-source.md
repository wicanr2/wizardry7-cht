# 主選單毛筆字樣來源

目前是[R976](../docs/goal/R976-bright-paper-brush-menu.md)的外觀候選，只限主選單六個按鈕。
一般訊息仍沿既有點陣字型。正式字型接線、020增補與玩家路徑尚未完成。

| 項目 | 固定來源 |
|---|---|
| 字型 | AR PL UKai TW，Book字面；ukai.ttc索引2 |
| 本機輸入 | `/usr/share/fonts/truetype/arphic/ukai.ttc`，唯讀 |
| 字型SHA-256 | `cfc758e69386a21bade8cb3b373fda20d7f2fac6e8904ced2ded98588b63e274` |
| 授權 | ARPHIC PUBLIC LICENSE；字模素材保留此授權，不改為程式RRSAL |
| 本機授權原文 | `/usr/share/doc/fonts-arphic-ukai/copyright`，唯讀 |
| 授權SHA-256 | `1f11c4a22865e1289f4694dad56b3561790a4b7341a511ddcd0edddcc3b95a46` |
| 原型渲染 | 固定Python容器內FreeType2.12.1，直接渲染16／24像素，沒有字模縮放 |

原作者為文鼎科技，CJKUnifonts版本由Arne Goetje修改。
原授權允許複製、改作及散布，保留授權全文及修改聲明；公開散布衍生字模時須依原條款提供改作。
來源為[Debian套件保存的原授權](https://sources.debian.org/src/fonts-arphic-ukai/0.2.20080216.2-5/debian/copyright/)。
本輪只作本機原型，沒有發布字型或衍生資源。

[預覽工具](../tools/re/preview_menu_brush.py)只讀正式catalog六則英文鍵及現存字型，
輸出文字、source IDs、字型／授權SHA、原alpha字模與2倍／3倍的字格收據。
[FreeType工具](../tools/re/font_preview_freetype.py)只支援已核Linux64 ABI，標頭與原生庫均來自固定容器，
不掛主機runtime或library。新輸出目錄不存在於輸入來源，原版與字型保持唯讀。

本機結果與授權全文在`workplace/r976/menu/`，
驗證與重跑入口見[研究025](../docs/re/025-ui-localization-coverage.md#r976-bright-paper-brush-menu)。
第三方字型、衍生字模二進位及本機圖檔不加入Git。
