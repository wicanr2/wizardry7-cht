# R70（輪 70）：睡醒翻頁

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R69——睡醒文字；r69-maze.state（3732M）可用。

## 目標

-load-state r69 續跑到 3770M：-clicks 3735/3742/3749M (160,150)
＋3738/3745/3752/3760M shots＋3762M 存點。
oracle：後續頁／回迷宮任一。

## 完成定義

1. 回迷宮：續探。
2. 還有頁：繼續翻。
3. 追記＋commit＋push。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測）

- 睡眠原因：「The last thing you remember was walking amongst the
  pretty flowers and your eyes began to feel drowsy...」——野蘭花睡眠陷阱。
- 3745M 過渡，3752M 起羅盤回來：在哪睡就在哪醒。
- r70-maze.state（3762M）已存。
- 下一步 R71：前進兩格離開花叢。
