# 024：隊員 HP／MP 顯示資料

日期：2026-10-02。對應[#35](https://github.com/wicanr2/wizardry7-cht/issues/35)與[R849](../goal/R849-character-vitals-evidence.md)。
研究限於原版顯示讀取，尚未有作弊實作或寫入規格。

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
| HP寫入與規則 | 未量戰鬥、受傷、死亡、上限、其他五位隊員及存讀檔資料流。顯示分子／分母不能單獨當作作弊可寫欄位的完整證據。 |

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
