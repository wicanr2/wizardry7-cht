# R19（輪 19）：迷宮 modal 排查——點底文、點頭像

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R18——pos 釘住 (284,165)＋按住 22M 仍不動；著陸幀有區域描述底文
  （"And with memory…"）且 1403M 後零文字事件（疑非攔截管線）；
  鍵盤免測（010：不讀 BIOS 佇列、port60＋IRQ stall）

## 目標

迷宮十動全滑鼠：底文中／上（描述關閉？）→頭像（手形解除？）→羅盤一圈＋
中央＋前進。任一動致文字事件／md5 變化即找到門。

## 完成定義

1. 全跑無 panic。
2. 底文消失／新事件／mv 互異：定位解除條件，R20 乘勝進迷宮，#16 關 issue。
3. 全滅：modal 假說死，R20 轉 probe PressReads（AX=5 呼叫者 CS:IP）查按鍵
   讀取路徑（golem 缺口 vs 遊戲邏輯）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；鍵盤；#17/#18。

## 記錄

2026-09-28 實測：十動零事件、十張全同（`2daf13…`）；但底文＋手形消失——
點底文關閉區域描述 confirmed。乾淨迷宮羅盤仍不動；重測幾何發現舊座標
y 偏低未對準。完成定義 1 達成；2 未達（R20 精確箭頭）。詳見 `docs/re/010`
R19 追記。
