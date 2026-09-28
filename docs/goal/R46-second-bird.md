# R46（輪 46）：收第二隻鳥

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R45——1 BLACKBIRD (1)，選單在下一輪；
  r45-maze.state（2855M）可用。

## 目標

-load-state r45 續跑到 2990M：三輪 ×（FIGHT(26,171)＋5M 後 START(160,166)），
2857／2900／2945M 起＋shots＋2980M 存點。
oracle：第二隻鳥死／勝利／戰利品／回迷宮任一。

## 完成定義

1. 戰鬥結束：判讀結算畫面，#16 續推（迷宮移動）。
2. 仍在打：再排輪次。
3. 滅團：讀檔改 RUN 脫離。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
