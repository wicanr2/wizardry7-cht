# R319（輪 319）：dump 碼段找分派表

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R318——數字鍵無效；26 輪選單僅 TRUCE；
  ip-truce.bin 定位點擊分派 0110:4096（real-mode 線性 0x5196？）。

## 目標

- load-state r293 → 19240M：空跑＋-dump-mem 4000-7000:/w/code1.bin
  ＋shot 確認。回 host 用 objdump -m i8086 反組譯，看 0x5196
  （檔內偏移 0x1196）周邊是否為滑鼠分派（CMP x/y＋跳表）。
  命中則找 six 詞熱區座標；映射不對則讀 dosgolem 009/003 spec
  找段基址。

## 完成定義

1. 分派點周邊反組譯判讀（THESUS 存活，無需推進）。
2. 追記＋commit＋push（R290–R319 一批）。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- dump 0110:4000 成功；0x4096 周邊是 overlay thunk 表
  （push bp／lcall *%cs:0x2Dxx），RE 需再深挖，暫擱。
- 務實繞行：RHALLICK 固定格記案，r291 右彎一次即來路。
  R320 退回（THESUS 零風險）。
