# R112（輪 112）：打甲蟲

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R111——2 BUGS 遭遇；r111-maze.state（5170M）可用。

## 目標

-load-state r111 續跑：-clicks FIGHT(26,171)@5175M，
START FIGHTING(160,166)@5182M＋5178/5185/5195/5205M shots
＋5210M 存點＋跑到 5215M。
oracle：FINAL ORDERS／回合結算（R43–R44 循環）。

## 完成定義

1. 回合執行：續打到贏。
2. 無反應：查按著計數，微調。
3. 追記＋commit＋push。

## 不在範圍

- 改 golem 本體；#17/#18。
