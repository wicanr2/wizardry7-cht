# R220（輪 220）：＋20M 相位再賭

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R219——落空一次仍死；r213（10195M）可用。

## 目標

- load-state r213 → 10285M：RUN (262,166)@10220M，
  START FIGHTING (160,166)@10228M＋shots＋hook＋10280M 存點。
  oracle：`THESUS RUNS AWAY!`／`DIES!`（再換相位）。

## 完成定義

1. 脫離：R221 續探。
2. 死亡：R221 換相位（＋30M）再賭；五連敗則改退守探左岔。
3. 追記＋commit＋push（R212–R220 一批）。

## 結果（實測）

- ＋20M：新動詞 `BUG STINGS THESUS`→`HIT!`→`LEG 10 DAMAGE`→死
  （蟲攻擊動詞：CLAWS／BITES／STINGS；部位：HEAD／BODY／LEG／HAND）。
- 五連死停損：R221 讀 r211 迴頭，找左側未探岔路。

## 不在範圍

- 改 golem 本體；#17/#18。
