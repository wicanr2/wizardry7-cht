# R215（輪 215）：三蟲戰重賭先攻

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R214——蟲先攻 9 傷死亡；r213（10195M 開打前）可用。

## 目標

- load-state r213 → 10260M：RUN (262,166)@10200M，
  START FIGHTING (160,166)@10208M＋shots＋hook＋10255M 存點。
  oracle：`THESUS RUNS AWAY!`（先攻或落空＋逃脫）／`THESUS DIES!`（再讀檔）。

## 完成定義

1. 脫離：R216 續探。
2. 死亡：R216 再讀 r213 重試（先攻約 1/2，落空另計）。
3. 追記＋commit＋push（R212–R215 一批）。

## 結果（實測，RNG 確定性 confirmed）

- 與 R214 逐步相同：`BUG CLAWS`→`HIT!`→`HEAD 9 DAMAGE`→`THESUS DIES!`。
- 同狀態同輸入＝同結果（RNG 確定）；R119/R120 證時序可改命。
- R216：RUN@10205M＋START@10213M（平移 5M 步換 timer／RNG 相位）。

## 不在範圍

- 改 golem 本體；#17/#18。
