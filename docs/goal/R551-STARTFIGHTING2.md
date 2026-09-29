# R551：START FIGHTING 脫離 BIRDS

## 前提
- 起點：`workplace/run/r550-maze.state`（FINAL ORDERS 開；resume 25610M）。
- START FIGHTING＝(160,166)；成功字串 `THESUS RUNS AWAY!`→`RUNNING!`。
- 步數放寬吃回合（R497/R500 經驗：回合結算需 ~150M）。

## 步驟（絕對步數！起點 25610000000）
 clicks: 160,166@25630M
 shots: s0@25620M, s1@25660M, s2@25700M
 save 25700M:/w/r551-maze.state, steps 25760M
- 未脫離（還在戰鬥）⟹ R552 首輪再 RUN；傷亡 ⟹ 讀回 r548。

## 安全
- 存點 r550（FINAL ORDERS）已落；THESUS 存活是完成定義。
- 結束追記＋commit＋push。
