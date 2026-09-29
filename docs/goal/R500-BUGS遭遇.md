# R500：2 BUGS 遭遇＋脫離

## 前提
- 起點：`workplace/run/r499-maze.state`（2 BUGS 選單開；resume 22630M）。
- RUN＝(262,166)→FINAL ORDERS；START FIGHTING＝(160,166) 執行。
- 成功字串 `THESUS RUNS AWAY!`→`RUNNING!`。

## 步驟（絕對步數！起點 22630000000）
 clicks: 262,166@22650M, 160,166@22690M
 shots: b0@22640M, b1@22670M, b2@22710M, b3@22750M
 save 22750M:/w/r500-maze.state, steps 22800M
- RUN 失效／THESUS 傷亡 ⟹ 停手，R501 評估（再戰／讀回 r498）。

## 安全
- 存點 r499（戰鬥選單）已落；敗／散隊 → 讀回 r498（E 格）。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R500 結果
- RUN→FINAL ORDERS→START FIGHTING 生效，回合結算：
  `BUG BITES THESUS／MISSED!／THESUS RUNS／BLOCKED!／BUG CLAWS THESUS／HIT! BODY 9 DAMAGE`。
- 首逃被 BLOCKED（逃脫看臉）；THESUS 存活（-9 BODY）；戰鬥選單開（b3 ASCII 見選項列）。
- r500＝戰鬥中。R501 再 RUN 一次。
