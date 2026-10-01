# 018：發行字型的字元與尺寸證據

日期：2026-10-01；關聯 [#25](https://github.com/wicanr2/wizardry7-cht/issues/25)、[規格007](../spec/007-native-font.md)、[R672](../goal/R672-delivery-prerequisites.md)。

## 輸入與可重現範圍

- 程式基線：`632a05c`；Go 1.24.13、Python 3.11.2；Docker image `sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
- GNU Unifont 17.0.05，官方輸入 `unifont_all-17.0.05.hex.gz`；SHA-256 `7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574`。沿用本機已有來源，沒有修改其他專案；已將固定HEX、原始來源壓縮包與完整授權／作者README存到本專案`workplace/fonts/unifont-17.0.05/`，逐檔雜湊在`SOURCE-SHA256.json`。
- 官方入口：[GNU Unifont](https://www.unifoundry.com/unifont/index.html)；固定版來源：[17.0.05](https://www.unifoundry.com/pub/unifont/unifont-17.0.05/)。字型 `COPYING` 說明雙重授權；發行字模採 SIL Open Font License 1.1，保留字型著作權與授權全文，不將第三方字型納入 RRSAL。
- 譯文基線：`text/story-trans.json`（SHA-256 `6609094d412f7484c54dced1680c3da5a6abbd978f3d59719ecbb820dcbbb996`）及物品、怪物與兩份 F1 TSV。完整稽核收據在 `workplace/r681/font-probe.json`。
- 其他輸入SHA-256：`items-trans.json`＝`ed23ec8f82543d1f12bf9561511939d7fa66773385a725a936010059b8e2e390`；`monsters-trans.json`＝`3f0816be28fcd1c90e1e2248fe48c54289142fda541ee99ebb8923089d9a8323`；`help.zh-TW.tsv`＝`0ddc05b783715d7a589e83e94f2aeb9c8755415e9caee736538e313c87629dac`；`help-names.zh-TW.tsv`＝`18f33d4e238fd875b8f589c8641036fe983ef84af02a55350ac9d0e5f9dfface`。
- 本輪未做反組譯；字形位址基準是 HEX 行首 Unicode 碼點，後續位圖每列 8 或 16 位、MSB 在左，16 列。原版位址與覆繪帶沿用研究017，不將字形索引冒充 DOS 位址。

## 已證實結果

重跑入口：在既有Docker容器中從專案根執行 `python3 tools/re/probe_font.py /font/unifont_all-17.0.05.hex.gz workplace/rNNN-font-probe.json`；將既有字型來源目錄唯讀掛到 `/font`，輸出檔須為新名稱，使用AGENTS.md第10節的容器資源、UID與生命週期規則。此探針不讀取或輸出防拷原文。

1. 目前譯文所需非 ASCII 字元共 2,877 個，固定版 HEX 全部涵蓋，缺字 0。
2. 104 個字的實際筆畫高達 16 像素；裁成 15 列會丟筆畫。不能縮字或刪筆畫。發行字型應保留原生 16×16，本機倚天保持原生 15×15。
3. 以 `r674` 既有故事／旁白收據重算，37 個不同頁文案在 2×／3×、字距 2 像素下全部可放入已批准文字帶。故事最多 2 列；旁白最多 3 列；均不超過原有 4 列上限。此項是資料／幾何探針，尚非正式遊戲驗收。
4. 原版列距 12 像素；2×中文字高16小於列距24，3×放大兩倍後字高32小於列距36。最後一列也在 y146–194 的已批准帶內。
5. 獨立審查重算F1：`ops.04`在2×倚天沒有裁切，Unifont若沿用原繪製會裁掉97個實際筆畫像素；原有兩欄容量在2×／3×另省略46／12行。spec003原本允許截斷，並非字型新增的垂直回歸；spec007 R5已明示取代該限制，以原生字寬折行與延續頁保留全部內容。

## 可丟棄畫面探針

`workplace/r681/hdat-363113713-unifont16.png` 與 `hdat-1516000000-unifont16.png` 以已驗收合成畫面為底，重畫故事第3頁與旁白一頁，確認完整字形與列距可讀。探針沒有改原版或正式程式，不能當作啟動、Draw、同狀態或完成中文化的收據。

正式驗收須由 dosgolem 重生：原版雜湊相同、批准帶之外畫面相同、每字筆畫依原始 HEX 獨立核對，並保留本機倚天回歸與 F1、一般覆繪的真實視窗檢查。

## 正式實作與驗收

依READY的規格007提供兩個字型後端；Unifont完整保留16列與原始Unicode碼點，
倚天15×15位元保持原樣。流文以資料尺寸決定字距、列寬、置中與畫筆；
F1按原生字寬折行與延續頁，不再截掉說明。原版資料及翻譯內容均未改。

| 證據 | 結果與實際範圍 |
|---|---|
| `r687/native-independent-final.txt` | 7,523項、失敗0。完整732份固定收據，故事31頁／旁白6頁，兩倍率及缺譯回退；原版Mem／A000／Pal／Step、raw色盤與整張PNG均核對，字形直接解析固定HEX |
| `r688/eten-regression.json` | 549份收據、5,490項、失敗0；三組模式的原版狀態、文字事件、完整PNG／A000／合成色號與r674逐位元相同 |
| `r690/main-final/report.json` | 395項、失敗0。r694兩後端2×各12頁、3×各10頁，全頁、循環與真實F1重開；內容直接取TSV，獨立折行、欄界與完整筆畫；另驗r691主選單／防拷說明與兩張故事Draw |
| `r689/negative-summary.json` | 26個載入負向案例、失敗0、沒有skip；非法尺寸、0／null、列數／高位、交叉全形尺寸等拒載，兩後端真實合法資料可載 |
| `r689/strict-font-missing-report.json` | 隔離U+10FFFF缺字退出1，目標輸出不存在；正式TSV與固定HEX雜湊不變 |
| `r682/`後端／雜湊負向記錄 | 兩類錯誤各在三個烘製器中非零退出，沒有新產物；不修改正式字型 |
| `font/charset.txt` | 正式非ASCII、F1 ASCII、數字與序數所需碼點2,924個；不含第三方字形二進位 |

F1實際呈現的「啟、說、阱」、故事的「偷」、旁白的「內、脫」包含第16列非零筆畫；
獨立位圖核對包含這些最底列。一般覆繪每份13事件／75字形，第16欄非零51例，
但本次一般字形第16列皆0；不能把空列當成已量到最底筆畫。

### 失敗分類與時間勘誤

- 第一版數字字模沒有在16字格正確置中；獨立檢查指出後修正，舊r683收據及其三份烘製輸入保留。
- 測試收據相對路徑與首次子代理Go模組快取設定失敗，屬驗證環境問題；在同image以正確完整路徑／既有快取重跑通過。
- r686部分3× F1圖片擷取早於實際頁面更新，最終oracle分類395項、30失敗；全部失敗收據保留。
  擷取器改為等待實際說明字色、新頁與連續穩定畫格；r694新目錄乾淨重跑，完整395項通過。
- 原生字型oracle首次7,235項有3失敗，皆錯用`hdat-1531000000`請求名判定旁白第5頁。
  r674與r687的實際Step都是1,536,000,000，已進第6頁；原版第5頁最後返回1,524,727,476，
  下一頁首入口1,532,204,150，1528M／1529M／1530M三份實際樣本仍在第5頁。
  修正oracle以原版20筆真正返回配對、下一頁入口與實際Step定窗，窗口內全部樣本核對，
  每頁另要求完整穩定覆蓋；故事6M／8M要求不變。時間與原版像素證據在`r687/native-time-window-evidence.json`，
  原`native-independent.txt`保留，現行為`native-independent-final.txt`。

### 重跑入口與限制

在既有Docker／Xvfb內、專案根目錄執行：

```sh
WIZ7_FONT=unifont WIZ7_UNIFONT_HEX=/font/unifont_all-17.0.05.hex.gz \
  python3 tools/re/run_story_checks.py workplace/rNNN --bake-dir workplace/bake-unifont
python3 tools/re/check_native_font.py --receipts workplace/rNNN \
  --font /font/unifont_all-17.0.05.hex.gz --baked workplace/bake-unifont \
  --eten workplace/r682/eten
python3 tools/re/run_native_gui.py workplace/rNEW
python3 tools/re/run_native_gui.py workplace/rOVERLAY --overlay-only
python3 tools/re/check_native_gui.py --gui workplace/rNEW \
  --overlay workplace/rOVERLAY --out workplace/rCHECK
```

GUI驗收器與單元測試目前使用r682兩後端烘製作固定回歸輸入；新資料基準須另列，不覆寫舊收據。
GUI程式須依AGENTS.md在有界Xvfb程序內執行；所有新輸出名稱必須不存在。
倚天回歸可用`run_story_checks.py --modes en,zh,zh2`，另掛使用者字型來源。
來源／程式／結果雜湊及工具版本存於`r682/verification-inputs-final.json`。

Draw沒有actual Step收據，不另冒稱精確同狀態。r694兩故事Draw與同請求HD圖完全相同；
r686首張Draw上方原版畫面曾差2,169像素，文字帶差0，屬非精確快照的時序差，保留原結果。
Mem仍是runner的全記憶體雜湊；未新增獨立完整記憶體dump。
防禦分支以單元／隔離烘製覆蓋，原版正常路徑未觸發，規格維持READY。
本驗收不代表全文路徑、完整通關、音訊或正式封包已完成。
本輪容器均已清理，輸出與抽查源檔1000:1000，未見root-owned檔或誤建`.md`目錄。
