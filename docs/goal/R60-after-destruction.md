# R60（輪 60）：destruction 後翻完

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R59——「destruction...」頁；r59-maze.state（3355M）可用。

## 目標

-load-state r59 續跑到 3390M：-clicks 3358/3365/3372/3379M (160,150)
＋3361/3368/3375/3382/3389M shots＋3385M 存點。
oracle：事件收尾（回迷宮／開打／新頁任一）。

## 完成定義

1. 收尾：順勢續推。
2. 仍同頁：試長按／雙擊／其他座標，下一輪。
3. 追記＋commit＋push。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測）

- 四頁：「And how can you be certain of the difference...」→
  「These thoughts begin to fade...」→
  「Perhaps somewhere here on this new world lie the answers...」→
  「And then again, perhaps not...」（s3/s4 同頁，事件尾）。
- r60-maze.state（3385M）已存。
- 下一步 R61：兩擊看回迷宮。
