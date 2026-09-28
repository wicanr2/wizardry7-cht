# R159（輪 159）：連走穿花田

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R158——睏意延遲假說（連走穿過），r158（7455M）可用。

## 目標

- load-state r158 → 7500M：前進 (160,162)@7460/7470/7480M
  ＋shots＋hook＋7495M 存點。
  oracle：連續推進（穿過）／RESTING（發作→等 180M）／OUCH!。

## 完成定義

1. 穿過：R160 續探。
2. 睡著：空跑等自然醒（~180M）。
3. 追記＋commit＋push（R135–R159 一批）。

## 結果（實測，NEW ENCOUNTER confirmed）

- @7460M 推進；@7470M `OUCH!`；@7480M 前進觸發
  `AN ENCOUNTER!`→`2 EERIE MOTHS(2)`＋戰鬥選單（`THESUS-VS-`＋七鈕）。
- 花粉蛾（呼應花田）；r159（7495M 開打前）可用；R160 首輪 RUN。

## 不在範圍

- 改 golem 本體；#17/#18。
