# R600：2 PHOOTS 遭遇 RUN 脫離

## 前提
- 起點：`workplace/run/r599-maze.state`（遭遇選單開；resume 26010M）。
- 遭遇時面 h1；THESUS 無傷；RUN 協議同 R496/R497。
- 脫離附帶位移（R497 前例），p3 需定位。

## 步驟（絕對步數！起點 26010000000）
 clicks: 262,166@26030M, 160,166@26070M
 shots: p0@26020M, p1@26050M, p2@26090M, p3@26130M
 save 26130M:/w/r600-maze.state, steps 26160M

## 安全
- 存點 r599（遭遇選單）＋r598（X8 遭遇前）＋r595（X2 乾淨）三層已落；
  敗／散隊 → 讀回 r598（同遭遇？RNG 或進）或 r595。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R600 結果
- RUN(262,166)→FINAL ORDERS（#26030M ✓）→START FIGHTING→
  `THESUS RUNS AWAY!`→`RUNNING!` ⟹ 脫離成功，無傷 ✓✓。
- p3 無 corpus 匹配（最近 f1-R591/f1-R541 wall-family 5.6KB）⟹ 脫離位移到未知視角。
- r600＝（逃後點，未定位）。R601 原地 C4 普查。
