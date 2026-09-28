# R31（輪 31）：FIGHT 後點鳥選目標

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R30——FIGHT 疑似接受（手形回＋老鷹換姿）但敵人仍 2 隻；
  讀出 "THESUS -vs- 2 BLACKBIRDS (2)"；r30-maze.state（2540M）可用

## 目標

-load-state 續跑：FIGHT→點鳥（driver (320,80)，邏輯 (160,80) 鳥身）×2 輪。
看敵數減少／傷害／勝利回迷宮。

## 完成定義

1. 敵數/HP/戰況變化：戰鬥機制確認，繼續打完。
2. 勝利回迷宮＋移動：#16 驗收達成，關 issue。
3. 全無反應：讀戰鬥選擇流程（目標選擇／確認鍵），開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
