# R577：X7-2 BUGS＋RUN

## 前提
- 起點：`workplace/run/r576-maze.state`（X7 戰鬥選單開；resume 26350M）。
- 2 BUGS 曾殺 THESUS（R501 廢棄分支）；首輪 RUN＋SF（1/3 往績）。
- 成功字串 `THESUS RUNS AWAY!`→`RUNNING!`。

## 步驟（絕對步數！起點 26350000000）
 clicks: 262,166@26370M
 shots: q0@26360M, q1@26390M, q2@26430M
 save 26430M:/w/r577-maze.state, steps 26480M
- 進 FINAL ORDERS ⟹ R578 START FIGHTING；未進 ⟹ 看單停手。

## 安全
- 存點 r576（戰鬥選單）已落；敗 → 讀回 r575（X7 迷宮）。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R577 結果
- RUN 生效 → FINAL ORDERS（#26370320269）；q1＝q2（0B 選單靜止）。
- r577＝FINAL ORDERS 開。R578：START FIGHTING(160,166)。
