# R48（輪 48）：左轉再前進

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R47——前進箭頭點擊送達但無移動，正面疑似密林無路；
  r47-maze.state（3030M）可用。

## 目標

-load-state r47 續跑到 3070M：-clicks 3035M 左箭頭 (133,178)，
3045M 前進 (160,162)＋3040/3050/3060M shots＋3065M 存點。
oracle：視角轉向／前進（任一變化即輸入有效）。

## 完成定義

1. 轉向或移動：#16 驗收達成，關 issue。
2. 皆無：試右箭頭／後退／鍵盤方向鍵，下一輪。
3. 追記＋commit＋push＋issue 留言／關閉。

## 不在範圍

- 改 golem；#17/#18。
