# R218（輪 218）：三蟲再賭先攻

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R214–R217——三連蟲先攻死＋繞行失敗；r213（10195M）可用。

## 目標

- load-state r213 → 10260M：RUN (262,166)@10200M，
  START FIGHTING (160,166)@10208M＋shots＋hook＋10255M 存點。
  oracle：`THESUS RUNS AWAY!`（THESUS 先攻）／`DIES!`（蟲先攻再讀檔）。

## 完成定義

1. 脫離：R219 續探。
2. 死亡：R219 讀 r213 換時序再賭（RNG 相位已證可調）。
3. 追記＋commit＋push（R212–R218 一批）。

## 結果（實測）

- 原相位＝原死亡（`HEAD 9 DAMAGE`，R214/R215/R218 三連同）。
- R219：RUN@10210M＋START@10218M（＋10M 相位）。

## 不在範圍

- 改 golem 本體；#17/#18。
