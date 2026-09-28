# R219（輪 219）：＋10M 相位再賭

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R218——原相位三連死；r213（10195M）可用。

## 目標

- load-state r213 → 10270M：RUN (262,166)@10210M，
  START FIGHTING (160,166)@10218M＋shots＋hook＋10265M 存點。
  oracle：回合敘述變化／`THESUS RUNS AWAY!`／`DIES!`。

## 完成定義

1. 脫離：R220 續探。
2. 死亡：R220 換相位（－5M／＋15M）再賭。
3. 追記＋commit＋push（R212–R219 一批）。

## 結果（實測）

- ＋10M：`BUG CLAWS`→`MISSED!`，次段 `HIT!`→`LEG 14 DAMAGE`→死。
- 需全落空（2–3 段攻擊）；估 ~5%/輪；R220 用 ＋20M 相位。

## 不在範圍

- 改 golem 本體；#17/#18。
