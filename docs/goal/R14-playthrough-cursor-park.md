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

## 記錄

2026-09-28 實測：park 假說**成立**（confirmed）。迷宮裡遊戲以 AX=4 把游標
釘在 (302,164)、範圍 [0-639,0-199]：1430M 時 sets=305；一次中央點擊期間
sets 305→447（約每 40K 步設一次）；每次動作後 pos 都被吸回 (302,164)，
注入的 (160,100)/(142,165) 全滅。旁白段 sets=1（無 park，pos 聽話）——
park 是迷宮 overlay 專屬行為。16M 跨按住＋右鍵全滅（按住讀到的都是 park 點）。
(302,164) 在 320 屏上≈右下角落（與截圖角落游標吻合）＝無效點。
R15：race（Move 後立即零 hover 按住 20M，賭 40K 空窗）＋park 原地按對照組。
