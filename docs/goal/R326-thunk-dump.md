# R326（輪 326）：thunk 鏈讀取＋選項表定位

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R325——ram.bin（512KB 線性）已 dump；選單三框視覺框。

## 目標

- -dump-seg 讀 0110:2D80 thunk 區＋掃線性記憶體選項名，找按鈕分派。

## 結果（實測）

- thunk 目標 740C:010B＝滑鼠驅動（int33 遠呼叫），ip-log 追蹤是死路。
- 線性 ~0x20E00 找到遭遇選項表（8B 名＋12B，stride 0x14，順序＝印出順序）；
  TRUCE 一條 handler 0110:406F 有效（R327 續解）。
- r325（ram.bin／seg110c 檔）在 workplace/run，未進版控。
