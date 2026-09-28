# R33（輪 33）：噪聲過場是卡住還是長載入

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R32——2555M 起全噪聲帶，2605M 仍在；r32-maze.state（2610M）可用

## 目標

-load-state 續跑：稀疏點擊＋vram 到 2700M。判讀噪聲是長載入（會播完）
還是卡死（ referrals: 按鍵有反應否）。

## 完成定義

1. 噪聲結束畫面可判讀。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. 卡死：回退到 r30/r31 檢查點換打法（FIGHT 目標改點法／RUN 脫離）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
