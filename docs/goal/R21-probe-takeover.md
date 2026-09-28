# R21（輪 21）：probe 接手——click-polls＋log-calls＋存檢查點

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R12–R20 前端動態用盡（精確箭頭亦滅）；probe 有 -click-polls
  （按住到遊戲讀 N 次才放）、-log-calls（int/AH 直方圖）、-save-state、
  結尾自動印 PressReads（含呼叫者 CS:IP）；前綴已定、Adlib 對齊

## 目標

把前綴搬到 probe（同座標、同 Type AAAA、-adlib）：主選單→建隊→START→
CANCEL→防拷→旁白 133 次→迷宮（點底文＋精確箭頭，-click-polls 5）。
收：log-calls（AX=3 vs AX=5 誰在讀）、PressReads（呼叫者 CS:IP）、
vram shots 前後比對、1420M 存檢查點（後輪秒級展開）。

## 完成定義

1. 到著陸（vram 與前端 mv 對得上）。
2. 任一 maze click 致 vram 變化：#16 驗收達成，關 issue。
3. 不動：以 log-calls＋PressReads 判定按鍵路徑（golem 缺口 vs 遊戲邏輯），
   開 READY spec（golem 上游或 overlay 機制）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
