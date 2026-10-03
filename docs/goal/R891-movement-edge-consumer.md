# R891：移動與格子判定的呼叫鏈

2026-10-03；工作[#36](https://github.com/wicanr2/wizardry7-cht/issues/36)。

## 目標

沿[R887正常前進](../re/016-party-position.md#r887-current-cell-reader)及
[R889值轉換](../re/016-party-position.md#r889-cell-value-consumers)，
追查raw `1184`的短getter及實際caller，區分輸入未取用與遊戲移動判定。

- 原始IDA庫的新副本只查已定位程式碼、必要直接xref與有界程式碼跳表；保留位址、bytes、原名及證據等級。
- 先核`1184`前後的短callee及caller。未取得靜態契約前，不假設其用途或runtime引數。
- 如需原版取樣，僅重播已合法R733起點與R887兩次正常前進按鈕輸入，固定停止步數做觀察開／關A/B。
- 觀察只採已定位的數值參數、呼叫／返回、分支與位置；不讀原始格子、訊息、鄰文、答案或RAM掃描，不注入記憶體、存檔、位置或呼叫。
- 原版、SDK、原始IDA庫、舊收據及快取唯讀，只寫`workplace/r891/`與有索引文件。
- 不重試待授權設定頁或尚未量到的前進路線。不改正式Go、譯文、字型、005或完整地圖資格。
- 取得下一個必要的移動分流證據後停止，封存與獨立唯讀審查，核擁有權及自有容器清理，再提交、推送與更新#36。

結果入口：[研究016](../re/016-party-position.md)、[WORKLOG](../../WORKLOG.md)。
