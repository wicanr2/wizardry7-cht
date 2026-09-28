# R294（輪 294）：初遇選單→FIGHT 轉戰鬥輪

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R293——遭遇 1 TRK. RHALLICK，初遇選單開啟，
  r293（19215M）可用。

## 目標

- load-state r293 → 19380M：FIGHT (26,171)@19260M（R43 既定座標）
  ＋shots@19240/19310/19370M＋hook＋r294 存點@19375M。
  oracle：戰鬥輪選單（-VS-／FINAL ORDERS／THESUS 指令）／
  遭遇文字變化。成功則 R295 下 RUN(262,166)＋START(160,166) 脫離
  （R120 既定）；按錯鈕則讀畫面重來（r293 可重載）。

## 完成定義

1. 進入戰鬥輪（THESUS 存活）。
2. 追記＋commit＋push（R290–R294 一批）。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- EXIT=0 到 19380M；FIGHT(26,171) 誤中 TRUCE（332B＋`TRUCE` 字，
  最左鈕）。按鈕列掃描：TRUCE(12-47)／STEAL(52-87)／…／
  FIGHT(240-275)／LEAVE(284-303)——FIGHT 在右側。
- r294（19375M）TRUCE 對話中；R295 重載 r293 直點 LEAVE(293,173)。
  失敗則改 FIGHT(257,173)→RUN(262,166)→START(160,166)。
