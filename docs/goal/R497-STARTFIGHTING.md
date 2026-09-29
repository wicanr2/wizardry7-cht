# R497：START FIGHTING 脫離

## 前提
- 起點：`workplace/run/r496-maze.state`（FINAL ORDERS 開；resume 22420M）。
- START FIGHTING＝(160,166)（R114 既定）。
- 成功字串 `THESUS RUNS AWAY!`→`RUNNING!`＋回迷宮。

## 步驟（絕對步數！起點 22420000000）
 clicks: 160,166@22440M
 shots: y0@22430M, y1@22460M, y2@22500M
 save 22500M:/w/r497-maze.state, steps 22560M
- 未脫離（還在戰鬥／THESUS 傷亡）⟹ 停手，R498 評估再戰或讀回 r495。

## 安全
- 存點 r496（FINAL ORDERS）已落；敗／散隊 → 讀回 r495。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R497 結果
- `THESUS RUNS AWAY!`→`RUNNING!` 脫離成功，THESUS 存活。
- y2 無匹配（vs w1/u1/v2 皆 15–18KB）⟹ 脫離附帶位移。
- r497＝逃後迷宮。R498 原地 C4 定位。
