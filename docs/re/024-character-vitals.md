# 024：隊員 HP／MP 顯示資料

日期：2026-10-02。對應[#35](https://github.com/wicanr2/wizardry7-cht/issues/35)與[R849](../goal/R849-character-vitals-evidence.md)。
研究涵蓋原版顯示讀取與第一隊員戰鬥回寫，尚未有作弊實作或寫入規格。

## 輸入與工具

原版與狀態唯讀；IDA資料庫只分析副本，保留原名、位址、operand與bytes。
本機重查入口為`workplace/r849/inputs-final.json`、`frozen-index.json`及`report.json`。
完整匯出與原版畫格只留本機，不提交原版素材。

| 輸入 | SHA-256 |
|---|---|
| `DS.EXE` | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |
| `workplace/ida/DS.EXE.i64` | `477f73a17b7abcaf3eef70d9baae16c77d9730f3bf7e0b4c0c9c9ed91c612412` |
| `workplace/ida/VPCVW.code` | `f8e2ed6c99860106a09e5ac48e6f1c3dbb964066a9c1ba949971872e0563425e` |
| `workplace/ida/VPCVW.code.i64` | `9502e774ce4a42772cfdadf312b85bcf81be6b77bbd7ef8c4f4c01627a37beb3` |
| `workplace/r733/two-forward.state` | `beecc0c54302d6304628d60f385fb4f3f6748e77c6bc859dfc153f669e337e6c` |

- IDA Pro 9.4，映像`ida-pro-9.4-idapython:locked-v1`，UID/GID1000:1000；最小探針核對8086與63個overlay函式。
- Go1.24.13，映像`wizardry7-go-ebiten:latest`。編譯採既有`workplace/r839/sdk`與私有狀態載入器；53份實際編譯SDK來源逐檔索引。正式SDK未改。
- IDA的`VPCVW.code`位址從0起算；`DS.EXE`線性EA從`0x10000`起算。執行期resident為`0110:IP`，資料位址以下列實測DS為準。

## 已證實的第一位隊員 HP 顯示讀取

在既有合法暖狀態3475M步，先執行3M步，點一次可見頭像`(20,29)`，再執行15M步。
點擊內部另外推進6M步，終點3499M步。沒有修改初態、存檔、原版記憶體或亂數。
角色頁顯示HP5／10；精確像素放大亦為5／10。小圖首次被誤讀為40／40，已更正。

| IDA `VPCVW.code` EA | 原始bytes／operand | runtime CS:IP | DS:offset／值 | 等級 |
|---|---|---|---|---|
| `0x14AD` | `FF B7 0C 5C`，`push word ptr [bx+5C0Ch]` | `0110:6502` | `1106:5C0C`／5 | confirmed：本畫面HP分子讀取 |
| `0x14DF` | `FF B7 0E 5C`，`push word ptr [bx+5C0Eh]` | `0110:6534` | `1106:5C0E`／10 | confirmed：本畫面HP分母讀取 |

兩讀取端前均為`B8 48 02 F7 6E 04 8B D8`，即`mov ax,248h; imul word ptr [bp+4]; mov bx,ax`。
記錄步距584 bytes為confirmed；本次BX=0，只驗第一位隊員。
同一原始函式`sub_1038`的`0x11BF`及`0x15CE`分別對到`0110:6214`、`0110:6623`，讀到1與114。
四個原始指令命中都支持本次overlay載入偏移`0x5055`，不外推其他overlay或載入狀態。

`ida-hp-caller-v2-window.json`保留`0x14A1–0x14F7`的有限呼叫端：

- 呼叫前先push區域緩衝`[bp-58h]`，再push上述數值，呼叫後清除4 bytes。
- 中間`0x14C9`的`B8 2F 00`印出斜線。
- 原始near call目標`0xB26F`在此次載入對到runtime`0110:02C4`；resident線性EA為`0x102C4`。
- resident入口讀取`[bp+4]`與`[bp+6]`，呼叫端之後傳同一緩衝給印字入口；IDA資料庫未定義此入口的函式邊界。

AX返回C字串的初始假設不採用。四次`NumericReturn=false`不代表原版轉換或印字失敗。
證據入口為`ida-resident-location.json`、`ida-resident-window.json`及原始四筆`on-v2/result.json`。

## MP 與其他未知

| 定位 | 證據與限制 |
|---|---|
| IDA `sub_36D`，EA`0x36D–0x513` | 靜態讀取`5C0C/5C0E`、`5C10/5C12`，並以步距4累加六組`5C1C/5C1E`值供第三個bar呼叫。原始operand與累加流程confirmed，MP語意仍是假說。 |
| 本次bar觀察 | 命中0；沒有MP畫面數值與資料欄位的動態對照。 |
| HP寫入與規則 | R851已驗第一隊員正常戰鬥回寫10→6，見後節；死亡、上限、其他五位隊員及存讀檔資料流未量。顯示與本輪回寫不能單獨當作作弊安全寫入契約的完整證據。 |

靜態入口為`ida-vitals-index.json`與`ida-vitals-slice-v2.json`。
不以函式名、數字搜尋命中或同一數值推定寫入位置。

## 觀察不改原版的核對

第一版只在外層Run觀察，沒有涵蓋Click內部執行，所以讀取命中0。
第二版用既有`oracle.Watch`觀察Click，四筆讀取全在點擊內部6M步出現；沒有新增SDK能力。
兩版的觀察開／關共四趟，三個唯一停止點的RAM、indexed、palette與暫存器SHA全部相同，原始PNG亦相同。
終點RAM SHA為`02dc6ad0f3622514b81b211fd1b0fd14cc9fc72a6a84a41bb36902352e107da9`。
保護段訊息取得命中0，四個scratch無檔；沒有MSG正文或手冊答案讀取。

這是私有暖狀態、`NoCursorWait`、非GUI的原版讀取證據。
它不構成冷啟動、正常存讀檔、作弊、全文中文化或可玩封包驗收。
十一組、TSV、完整錄放、人耳與三平台交付門檻保持。

## 工具勘誤

- IDA第一版slice使用不存在的`UA_MAXOP`；保留失敗後改為列舉`ins.ops`。
- 數值入口探針錯把函式邊界當成前提；有限metadata確認位址有效但沒有函式定義。
- HP呼叫端第一版起點落在指令中間，該片段不採用；第二版從已確認的`0x14A1`開始。
- 像素放大首建置超過程序限額；沿同一映像限制Go並行後成功。這些失敗均為工具問題，不記成產品缺陷。

## 第一位隊員的戰鬥 HP 回寫（R851）

日期2026-10-02，對應[#35](https://github.com/wicanr2/wizardry7-cht/issues/35)及[R851](../goal/R851-hp-writer-evidence.md)。
本輪沒有正式作弊實作。新增證據將HP的原版顯示讀取接到正常戰鬥回寫，不證實死亡、復活、上限規則、其他隊員或存讀檔。

### 原始資料與位址基準

| 輸入 | SHA-256 |
|---|---|
| `workplace/ida/VMELE.code` | `433d3e925b8413b7c5da5a22803e648fdd0700730a5de3424e937aa56295cde2` |
| `workplace/ida/VMELE.code.i64` | `9e9092f83b4a6a9079380bd93716365d548c22542c5a10b37a9dfa42c5c3031c` |
| `workplace/run/r43-maze.state` | `8438009616992ece7187407d68c742d0765dfc765bf2cc41018e3848e43dced0` |

沿用前表DS.EXE與resident資料庫；IDA Pro9.4、Go1.24.13、既有非root離線映像與私有r839 SDK。
VMELE位址是抽出code的IDA線性EA，自0起算，不能直接當resident或執行期位址。
`DS.EXE.hp-references.json`只查139個已定義函式、5100個指令項目，直接欄位命中0；不代表沒有間接讀寫或未定義程式碼。
`VMELE.code.hp-references.json`查42個已定義函式、4902個指令項目，取得兩個讀取與一個寫入；不外推其他模組。

### 最小寫入鏈

IDA原名`sub_1EA6`，主範圍`0x1EA6–0x1F2B`，133 bytes。
原始函式、operand及bytes保留在`workplace/r851/VMELE.hp-writer-slice.json`。

| IDA VMELE EA | 原始指令 | 證據等級 |
|---|---|---|
| `0x1EAC–0x1EB8` | `mov ax,2Ch; imul word ptr [bp+4]; mov si,ax; mov bx,ds:8482h; mov ax,[bx+si+2]` | confirmed：由DS指標指定的44-byte記錄讀取word |
| `0x1EBC–0x1EC4` | `mov ax,248h; imul word ptr [bp+4]; mov bx,ax; pop ax` | confirmed：同一引數選取584-byte目的記錄 |
| `0x1EC5` | `89 87 0C 5C`，`mov [bx+5C0Ch],ax` | confirmed：本次正常戰鬥第一隊員HP分子回寫 |

有限讀取窗口在`0xD08`讀`[bx+5C0Eh]`後寫至44-byte記錄`+4`，在`0xD2C`讀`[bx+5C0Ch]`後寫至`+2`。
原始IDA歸屬為`sub_F7A`，主範圍`0xF7A–0xF87`，另有尾段`0xCDD–0xE36`；不能因讀取位址在主範圍外而刪除該證據。
`[bx-7B7Eh]`保留原operand；其位元組位移為`0x8482`，不改寫原始定位。

五個直接call位於`0x1F75`、`0x1FB2`、`0x32B0`、`0x32DC`、`0x36C0`，分屬`sub_1F2B`、`sub_3193`、`sub_36AA`。
各call以迴圈索引作引數，迴圈與`ds:69C8h`比較；此處保留原欄位，不單憑比較推定完整隊伍數量契約。
控制流原始匯出為`VMELE.hp-writer-callers.json`。沒有深入與回寫無關的helper或硬體driver。

### 正常戰鬥的動態核對

來源為[R43](../goal/R43-fight-target.md)已記錄的正常`FINAL ORDERS`狀態；正常開打位置沿[R44](../goal/R44-start-fighting.md)。
先零步數、零輸入還原並檢視畫格，確認按鈕。兩側固定同一完整state，未改初態或重擲。
2690M起點先執行2M步，畫格仍是`FINAL ORDERS`；點一次`(160,166)`，Click內部6M，再執行2M及31M步。
四停止點為2690M、2692M、2700M、2731M。畫格顯示正常開打及敵方咬擊；不比較舊R44骰序或傷害公式，也未識別或改寫seed欄位。

| 執行步數 | runtime CS:IP | DS／目的offset | DS暫存指標／來源word | AX／前值→後值 |
|---|---|---|---|---|
| 2707136130 | `0110:5371` | `1106:5C0C` | `1106:D8B2`／`+2`=10 | 10／10→10 |
| 2720084008 | `0110:5371` | `1106:5C0C` | `1106:D8B2`／`+2`=6 | 6／10→6 |
| 2724396717 | `0110:5371` | `1106:5C0C` | `1106:D8B2`／`+2`=6 | 6／6→6 |

三次均為BX=0、SI=0，下一步為原步數+1、IP=`0110:5375`，目的word等於AX及來源word。
每次從`1106:8482`讀到近指標`0xD8B2`，有效來源word為`1106:D8B4`。
原始寫入前綴`B8 48 02 F7 6E 04 8B D8 58`及指令`89 87 0C 5C`與IDA逐位元相符。
此模組載入偏移為`0x5371-0x1EC5=0x34AC`，與先前VPCVW的`0x5055`不同；兩者只適用各自實測載入。
第一隊員HP分子與原版回寫語意為confirmed；44-byte記錄的本次HP暫存角色亦為confirmed。
該指標的配置、釋放、有效期、其他隊員與跨戰鬥使用仍未知。

觀察開／關兩趟的四個停止點RAM、indexed、palette、暫存器及原始PNG相同；初點亦與零步數inspect相同。
保護段訊息取得0、三個scratch無檔，沒有新MSG正文、手冊／答案／autoplay讀取或探針原版寫入。
完整runtime輸出只留`workplace/r851/on/result.json`與`off/result.json`。

### 對後續實作的限制

戰鬥存在獨立HP暫存記錄，原版會回寫隊員欄位。只更新隊員欄位的作弊方案，不能未驗證就假設會保持效果。
本輪未施作作弊或修改兩份HP，也未驗證指標生命週期與正式前端暫停時的安全寫入契約。
因此仍須撰寫並審查作弊DRAFT規格，不因找到原版writer就宣稱作弊已完成。
十一組、TSV、正常冷存讀、錄放、人耳與三平台交付門檻保持。

本機重查入口為`workplace/r851/report.json`、`inputs-final.json`、`frozen-index.json`。
圖片第一次轉送的額度不足，兩張未成功處理；後續增加額度並核完整base64長度後檢視同一原圖，不重跑戰鬥或改圖。
