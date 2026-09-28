# R302（輪 302）：提早送 L 鍵

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R300——阻塞停；R301——block-after 0 被 kill；
  r293（19215M，初遇選單）可重載。

## 目標

- load-state r293 → 19390M：-bios-keys "L"（-bios-key-from 19250M，
  載入即送，每 5M 一個）＋shots@19330/19360/19385M＋hook＋
  r302 存點@19385M（預設 block-after）。
  oracle：跑滿（L 被讀走，看效果）／阻塞停（L 無效，換鍵）。
  若完成，看截圖＋文字判讀效果。

## 完成定義

1. L 鍵效果確認（THESUS 存活）。
2. 追記＋commit＋push（R290–R302 一批）。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- EXIT=0 但阻塞停於 19250100206（載入後 35M）：BDA 根本是空的，
  -bios-keys 沒送進去（或時序錯過）。r302 無截圖無存點。
- R303 改 -keys-at（stdin＋IRQ1；初遇選單是迷宮 overlay，
  R40 證偽的只是戰鬥 overlay 無 int09）。
