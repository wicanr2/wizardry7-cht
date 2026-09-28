# R301（輪 301）：BDA 送 L 跑穿阻塞

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R300——初遇選單阻塞等鍵（INT16 AH=07），提早停帶來
  的教訓；r293（19215M，初遇選單）可重載。

## 目標

- load-state r293 → 19390M：-block-after 0（跑穿鍵盤阻塞）、
  -bios-keys "L"（-bios-key-from 19300M，-bios-key-every 5M）＋
  shots@19340/19370/19385M＋hook＋r301 存點@19385M。
  oracle：L 被讀走後的效果（逃脫字串／戰鬥輪／選單變化／無變化）。
  L 無效則 R302 送 F。

## 完成定義

1. L 鍵效果確認（THESUS 存活）。
2. 追記＋commit＋push（R290–R301 一批）。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- -block-after 0 ＋鍵盤阻塞自旋會被 kill（兩次，s0 後即死），
  此路不通。r301 無存點。
- R302：改回預設 block-after（阻塞停＝按鍵無效判據），提早送鍵。
