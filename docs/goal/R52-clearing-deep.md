# R52（輪 52）：從 clearing 往深處

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R51——出生點是死胡同區；r49-maze.state（3095M，clearing）可用。

## 目標

-load-state r49 續跑：R(3100)＋shot，R(3110)＋shot，F(3120)＋shot，
F(3130)＋shot，存點＋steps 到 3145M。
oracle：小徑深處新視角（脫離出生點林子）。

## 完成定義

1. 新地形：沿路續探。
2. 被擋／繞回：試左轉序列，下一輪。
3. 追記＋commit＋push。

## 不在範圍

- 改 golem；#17/#18。
