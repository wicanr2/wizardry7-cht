# R26（輪 26）：穿過過場特效找迷宮

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R25——1734M＋ 是 13h 無誤的雜訊過場特效（planar 寫入全零）；
  r25-maze.state（1720M）可用

## 目標

-load-state 續跑穿過特效：旁白式點擊到 1880M＋迷宮六動＋vram 每 20M
（1780–1900M）＋存點。看特效後是不是乾淨迷宮。

## 完成定義

1. 特效結束畫面可判讀（迷宮／新過場／黑屏）。
2. 若乾淨迷宮＋點擊致視角變化：#16 驗收達成，關 issue。
3. 若仍過場／黑屏：繼續延伸或查該段程式（overlay 劇本）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
