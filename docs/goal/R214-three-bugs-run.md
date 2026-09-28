# R214（輪 214）：三蟲戰首輪 RUN

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R213——3 BUGS 遭遇，r213（10195M 開打前）可用。

## 目標

- load-state r213 → 10260M：RUN (262,166)@10200M，
  START FIGHTING (160,166)@10208M＋shots＋hook＋10255M 存點。
  oracle：`THESUS RUNS AWAY!`（脫離）／結算文字（蟲先攻評估）／
  `THESUS DIES!`（讀檔重來）。

## 完成定義

1. 脫離：R215 續探。
2. 纏鬥存活：再排 RUN（蟲多段攻擊，活過首輪再逃）。
3. 死亡：讀 r213 重試先攻賭局。
4. 追記＋commit＋push。

## 結果（實測）

- RUN 訂單＋FINAL ORDERS → 回合執行；但蟲先攻
  `BUG CLAWS THESUS`→`HIT!`→`HEAD 9 DAMAGE`→`THESUS DIES!`→主選單。
- 先攻賭輸；r214 存點＝主選單態（作戰仍用 r213）。
- R215 讀 r213 重試同一 RUN。

## 不在範圍

- 改 golem 本體；#17/#18。
