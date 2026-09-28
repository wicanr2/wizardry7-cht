# R287（輪 287）：再試喚醒

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R286——仍睡（RESTING 未斷），r286（17275M）可用。

## 目標

- load-state r286 → 17440M：前進 (160,162)@17320M（喚醒）＋
  shots@17300/17370/17430M＋hook＋r287 存點@17435M。
  oracle：醒來文本／持續 RESTING（則空跑等自然醒）／被襲。

## 完成定義

1. 醒來續行或排程等待。
2. 蟲遭遇：首輪 RUN。
3. 蛾遭遇：RUN。
4. 追記＋commit＋push（R285–R287 一批）。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- EXIT=0 到 17440M；前進@17320M 喚醒：`The last thing you
  remember was walking amongst the pretty flowers...`
  （＝messages.json id 15062，實機觸發覆蓋）。無遭遇。
  r287（17435M）已醒。R288 行走驗證。
