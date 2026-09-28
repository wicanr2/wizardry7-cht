# R14（輪 14）：游標 park 假說驗證＋跨 AUTO 按住

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R13——迷宮零 error（送達但不理）；golem 記載形狀：AX=4 自設游標蓋掉
  注入座標；`MouseSets()`／`MouseRange()`／`PressButton` 皆 public

## 目標

迷宮每動作後記 `Mouse()`＋`MouseSets`（len＋末筆）＋`MouseRange()`：
若注入座標與 park 點系統性偏離且遊戲讀到的是 park 點，假說成立。
同輪加跨 AUTO 按住（PressButton(0)@1438M／ReleaseButton(0)@1458M，前進箭頭）
與右鍵變體；誰先讓 mv 變化誰是答案。

## 完成定義

1. 全跑無 panic；滑鼠觀測行可判讀 park／range。
2. 任一變體致 mv 互異：#16 驗收達成，關 issue。
3. park 成立但按住亦滅：下一層（overlay 內命中邏輯／-regs-at 盯熱點服務）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- #17 錄放；#18 音訊。
