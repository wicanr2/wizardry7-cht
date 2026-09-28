# R47（輪 47）：迷宮移動一格

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R46——戰鬥勝利回迷宮；r46-maze.state（2985M，迷宮＋羅盤）可用。

## 目標

-load-state r46 續跑到 3035M：-clicks 2990M＋3000M 點前進箭頭 (160,162)
＋2995/3005/3015/3025M shots＋3030M 存點。
oracle：迷宮視角變化（前進一格）。

## 完成定義

1. 移動一格：#16 驗收達成，關 issue（m2-playthrough 全通：開機→建隊→
   迷宮→戰鬥勝利→迷宮移動）。
2. 無反應：微調箭頭座標／試左右箭頭，下一輪。
3. 追記＋commit＋push＋issue 留言／關閉。

## 不在範圍

- 改 golem；#17/#18。
