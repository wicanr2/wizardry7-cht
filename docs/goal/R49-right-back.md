# R49（輪 49）：右轉＋後退排查

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R48——左箭頭僅 UI 高亮無轉向；r48-maze.state（3065M）可用。

## 目標

-load-state r48 續跑到 3100M：-clicks 3070M 右箭頭 (207,177)，
3080M 後退 (160,187)＋3075/3085/3095M shots＋3095M 存點。
oracle：任一視角變化。

## 完成定義

1. 轉向／移動：順勢走出出生點，#16 驗收達成關 issue。
2. 皆無：試鍵盤方向鍵（bios-keys）或查 VMAZE 移動條件，下一輪。
3. 追記＋commit＋push＋issue 留言／關閉。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測）

- 三張 shots md5 全異：右轉（3070M）＋後退（3080M）生效。
- 3085M 起視角變成開闊草地小徑（密林→空地）：迷宮移動確證。
- #16 驗收達成：開機→建隊→迷宮→遭遇戰勝利→迷宮移動，全程打通。
- r49-maze.state（3095M，空地）已存。
- worklist m2-playthrough 移到 done，關 #16。
