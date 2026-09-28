# R45（輪 45）：連打三輪

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R44——回合結算全通；r44-maze.state（2725M，第一輪尾聲）可用。

## 目標

-load-state r44 續跑到 2860M：三輪 ×（FIGHT(26,171)＋5M 後 START(160,166)），
分別在 2740／2780／2820M 起＋-log-calls＋每 ~10M shots＋2855M 存點。
oracle：敵方減員（BLACKBIRD 數量／(2) 變化）或勝利／滅團畫面。

## 完成定義

1. 鳥死／戰鬥結束：回迷宮或戰利品畫面，#16 接近驗收。
2. 三輪後仍在打：看傷害節奏，估計還需幾輪再排。
3. 滅團：讀檔換策略（RUN 脫離重測——現在點擊有效了）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
