# R578：START FIGHTING 脫離 BUGS2

## 前提
- 起點：`workplace/run/r577-maze.state`（FINAL ORDERS 開；resume 26430M）。
- START FIGHTING＝(160,166)；成功 `THESUS RUNS AWAY!`→`RUNNING!`。
- BUGS 首逃曾 BLOCKED（R500）；失敗 ⟹ R579 再 RUN；傷亡 ⟹ 讀回 r575。

## 步驟（絕對步數！起點 26430000000）
 clicks: 160,166@26450M
 shots: s0@26440M, s1@26470M, s2@26510M
 save 26510M:/w/r578-maze.state, steps 26570M

## 安全
- 存點 r577（FINAL ORDERS）已落；THESUS 存活是完成定義。
- 結束追記＋commit＋push。

## R578 結果
- SF 生效（投递 ✓）；零 hook 文本（戰鬥訊息走別路）；戰鬥進行中（雙 BUG 在，
  s2 ASCII 驗；CS:IP 0110:5A45 存活）。
- r578＝戰鬥中（RUN 令已下）。R579 再 SF 一輪（逃脱再 roll）。
