# 030 流文保留半形斜線字模

狀態READY；2026-10-08，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、
[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[R981](../goal/R981-no-response-unification.md)。
本契約只補目前被拒絕的ASCII斜線字形，不改原版資料、語意、Go runtime或flow schema。

## 玩家問題與證據

006 R1要求譯文每字Big5雙byte，半形`/`為單byte，導致含斜線的整則譯文退回英文。
17385統一成已審候選時，也會失去原flow中文；一般xlate last-entry未變不能掩蓋這個減少。
confirmed於現行建置：兩後端原型僅允許`/`並取原生字型，不變其他條件，
段落／索引行仍2906／6307，可用行5058、完整段落2292、msgs8779／可用7331，
相對R980可用msgs新增133，舊可用值沒有消失；原型另含17385候選文案改動，
合計134個msgs差異。兩項差異分列，不把文案統一當作字模修補。完整集合須獨立重算再READY。
本組及新增行正常上屏未量，不以靜態可用性宣稱已中文化。

來源text/messages.json與正式譯文沿R980已核SHA；006既有runtime CS:IP與文字帶保持。
固定Unifont17.0.05與合法本機倚天ascfont.15為原生字模來源，後者只本機，不提交。
不解碼答案表，不讀原版旁文或新增RE；字型讀法沿既有007與F1半形字契約。

固定輸入：Unifont壓縮hex SHA7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574；
倚天ascfont.15共3,840 bytes，SHA1d0cf09d0a319a9e7039190688c6a905ba4370bd369fbfcfa43b2078480d6918。
Python3.13.15；字形定位以Unicode U+002F及倚天ASCII第47格、每格15 bytes為準。
只用正式舊譯文的新原型已獨立分離：兩套133個msgs／132行新增，舊值與全部舊字模保持，
17385仍為舊譯文；新增碼點只有`/`與因此首次可用的「囚」。

## 契約

- S1：`usable_zh`仍剝低位控制碼、非空、逐字可用及整則通過才更新字形表。
  唯一新增例外為U+002F `/`。其他ASCII仍按006 R1拒絕；不轉換成全形slash，不改譯文或原文鍵。
- S2：slash採固定原始字模，保持完整8像素筆畫，不縮放或重採樣。
  Unifont原生8×16依既有square=True放在16×16字格中央，左右各4px；
  倚天ascfont.15原生8×15放在15×15字格，左3px、右4px，row byte左移4。
  對應字高／有效位元／非空字模須在建置開始驗；缺檔或壞形狀明確失敗，不產部分新資料。
  slash來源錯誤不得被倚天既有ValueError轉空譯文的分支吞掉。
- S3：維持glyphs為同字高的row整數陣列、size15／16。Go排版仍每字固定size+2字距、
  既有整數放大、換行、四列上限及文字帶守門保持。字格留白不裁原生筆畫，不新增寬度欄位或sidecar。
- S4：獨立比較全部paras／msgs／字形，限定變動為原來因slash被拒且其他字元全可用的譯文，
  加slash字形與因此真正使用的新CJK碼點；原words、t、id、end與其他譯文保持。
  無任何舊可用譯文變空，英文輸入、原版記憶體、檔名及保存不變。
- S5：synthetic tests核正例、其他ASCII拒絕、字模錯形狀、全字形原始像素、整則更新原子性。
  雙後端從乾淨輸入重建，獨立原字型比對；同一正常故事／旁白抽測既有輸出保持。
  新增slash行未正常觸發時只寫READY／靜態可用，不宣稱全部133行上屏或全程parity。
- S6：六列同義統一要在此READY及實作驗收後續接，再獨立核正式來源、產生器及完整差異。
  正式字元清單包含slash字元，原始字型／生成資源／完整原版收據只留本機。
  清單同時補回029已實作的通用A–Z與答案標籤覆蓋，沿bake_help既定集合；不讀或烘製答案字串。

這是缺字／拒畫修補，完整目標與既有語意不變；不為較容易的字典相同結果批准中文覆蓋減少。

## READY審查

唯讀子代理r974_candidate_review已獨立核兩套正式舊譯文原型的133個msgs／132行，
全部舊可用值與舊字模、原w／words token bytes、t／id／end保持，新增`/`與「囚」。
審查回覆沒有實質阻擋，可READY；S2的slash來源錯誤不可吞掉已補入契約。
主代理另以原始字型位元核兩個新增碼點共962像素相符，採認READY後才寫正式工具。
這些證據只證明建置範圍，新增行正常上屏仍未量。

報告入口workplace/r981/review-slash-spec.json，SHA93df623b8d9d2b753904ad94283732993d35b6a3fca0a539653ea7fcf5fdceca；
Markdown SHA507f6c8463430e19b4ac0896559aa7754e3a6397289dc4ec05a1b23edd30639c。
完整回讀0／0／0；正式工具在[bake_flow.py](../../tools/re/bake_flow.py)，
邊界測試在[test_flow_slash.py](../../tools/re/test_flow_slash.py)，以Docker內
`python -m unittest discover -s tools/re -p 'test_*.py'`執行。
正式舊譯文兩套乾淨重建與已審原型逐bytes相同；8項新邊界測試及4項安全讀取回歸通過。

最終實作獨立審查0／0／0已完整回讀採認。
JSON SHA1756f741ad786a267419d2b75a9ba93ea70137622b1515d507e23521c6dd1c05，
Markdown SHA8bed37911f2f14db32108acf2fe8c03b24906b4a0d2069334b35cdea7c6c2b44，入口workplace/r981/review-implementation。
正常故事／旁白抽樣見[研究017](../re/017-flow-text-evidence.md#r981-flow-slash)；
新增133則上屏仍未量，維持READY，不因既有畫格通過標全部CONFORMED。
