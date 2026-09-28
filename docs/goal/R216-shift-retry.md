# R216（輪 216）：平移時序重賭

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R214–R215——同輸入同死亡（RNG 確定）；r213（10195M）可用。

## 目標

- load-state r213 → 10265M：RUN (262,166)@10205M，
  START FIGHTING (160,166)@10213M＋shots＋hook＋10260M 存點。
  oracle：回合敘述變化（RNG 相位改變）／`THESUS RUNS AWAY!`／`DIES!`。

## 完成定義

1. 脫離：R217 續探。
2. 存活纏鬥：再排 RUN。
3. 死亡：換更大平移（±10M）或改訂單（FIGHT 先攻賭蟲減員）。
4. 追記＋commit＋push（R212–R216 一批）。

## 結果（實測，RNG 相位 confirmed 可調）

- 平移 5M 步改命：`BUG BITES`→`MISSED!`（首蟲落空），
  但次蟲 `BUG CLAWS`→`HIT!`→`BODY 14 DAMAGE`→`THESUS DIES!`。
- 三蟲一輪多段攻擊：solo 全落空＋逃脫機率≈0——此戰無解，繞行。
- r213 作戰態保留（蟲情報）；R217 讀 r212（戰前迷宮）左彎繞行。

## 不在範圍

- 改 golem 本體；#17/#18。
