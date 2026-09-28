# R32（輪 32）：跟過戰鬥轉場看落點

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R31——FIGHT＋點鳥後出現轉場特效（2555M 全噪＋2575M 黑條擠壓），
  戰鬥有反應了；r31-maze.state（2575M）可用

## 目標

-load-state 續跑跟過轉場：FIGHT＋點鳥交替＋vram 每 20M（2585–2615M）。
看落點（勝利回迷宮？下一波？新場景？）。

## 完成定義

1. 落點可判讀。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. 新戰鬥／新謎題：繼續解（下一輪）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
