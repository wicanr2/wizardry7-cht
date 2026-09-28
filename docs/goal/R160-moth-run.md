# R160（輪 160）：蛾戰首輪 RUN

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R159——2 EERIE MOTHS 遭遇，r159（7495M 開打前）可用。

## 目標

- load-state r159 → 7560M：RUN (262,166)@7500M，
  START FIGHTING (160,166)@7508M＋shots＋hook＋7555M 存點。
  oracle：`THESUS RUNS AWAY!`（脫離）／結算文字（蛾強度評估）／
  `THESUS DIES!`（讀檔重來）。

## 完成定義

1. 脫離：R161 續探。
2. 纏鬥：蛾傷害評估（決定打／跑／繞）。
3. 死亡：讀 r159，重試（先攻賭局）或繞行。
4. 追記＋commit＋push（R135–R160 一批）。

## 結果（實測，蛾戰報 confirmed）

- RUN＋FINAL ORDERS → 回合執行。
- `EERIE MOTH LASHES TONGUE THESUS`→`MISSED!`；
  `EERIE MOTH THRUSTS TONGUE THESUS`→`HIT!`→`LEG 3 DAMAGE`。
- `THESUS RUNS AWAY!`→`RUNNING!`——脫離成功。
- 蛾舌 3 傷（蟲 12–17）：蛾弱得多，必要時可打（solo 仍險）。
- r160（7555M）可用；R161 前進驗證＋push R135–R161。

## 不在範圍

- 改 golem 本體；#17/#18。
