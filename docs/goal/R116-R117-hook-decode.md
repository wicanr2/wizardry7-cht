# R116–R117（輪 116–117）：鉤子解碼＋凍結判定

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R114–R115（RUN＋START 三輪無結論）。

## 操作

- R116：load r115＋RUN@5237M＋START@5245M → 5285M，EXIT=0；
  但 -save-state 恰等於 -steps 未觸發，存點遺失（邊界教訓：
  -steps 必須嚴格大於所有事件步數）。
- R117：load r115 重排同腳本 → 5340M，EXIT=0，r117-maze.state（5335M）。
  掛 `-call-args 0110:3895→改 0110:3877:1`（sub_13877 逐字元，
  AX＝字元碼）＋-arg-regs，全程 58 次呼叫。

## 解碼（confirmed，AX 字元流）

- #5245267559（START@5245M 後 2.7M）：`2 BUGS(2)`；400K 步後又一次
  `2 BUGS(2)`（DX 5F01，第二遍渲染）。
- 接著 `THESUS`、`-VS-`、`FIGHT`、`PARRY`、`SPELL`、`USE`、`RUN`、
  `MOVE`、`EQUIP`——整套戰鬥選單重繪。
- SI 驗證：'T' 時 SI=5BF5＝DS 字串 `THESUS`(DS:5BF4)＋1 ✓；
  'F' 時 SI=FD79，'R' 時 SI=FDC9（DS dump 內 `RUN` 在 0FDC8 ✓）。
- 「由 SEG:OFF」的 SEG＝字元碼本身（near call，stack 頂是返回位址＋
  push 的字元）；OFF 0x38BE＝sub_13895 內 lodsb 迴圈呼叫點，
  0x54A0/0x550A/0x5538＝空白／括號呼叫點（sub_102DA 縮排路徑，假說）。
- 結算文字（R112 見過的 `THESUS PUNCHES`／`BUG MISSED!`／`BUG BITES`）
  **零出現**；5245.7M 後 95M 步零印字。

## 判定

- 回合從未執行：START 點擊只觸發選單重繪（訂單未確認時 START＝no-op
  重繪，無錯誤訊息——UI 發現）。
- r117-s0 == r117-s2（35M 步凍結）；對 r111 差 8953（選單態），
  對 r110（迷宮）差 24643——仍在戰鬥選單，THESUS 活著。
- DS 段 dump（0x11060，64KB）：`THESUS`@5BF4、`BUG`@D7E8／`BUGS`@D7F8、
  `RUN`@FDC8；escape/flee/miss/dead 無靜態字串（結算訊息應在 overlay
  段，不在 DS）。
- RUN(262,173) 沿用 R38/R39（黑鳥戰）；但 RUN 在本 playthrough 從未
  成功脫離——座標或訂單確認流程二者必有一錯。

## 下一步

- R118：FIGHT(26,171) 校準（R112 已驗證可執行回合）＋hook 看結算文字。
- 若存活：R119 起掃 RUN 列 y（150/158/166/174/182）找真 RUN。

## R118 結果（實測，hook 全字元流 confirmed）

- FIGHT 訂單成立（`THESUS-VS-`）→ START 出現 `FINAL ORDERS／START
  FIGHTING／BACKUP TO PREV／TERMINATE GAME` 確認框→ 回合執行。
- `THESUS PUNCHES BUG`→`HIT!`→`BODY 2 DAMAGE`；
  `BUG CLAWS THESUS`→`HIT!`→`BODY 12 DAMAGE`→`THESUS DIES!`。
- 主選單字串：`CHARACTER MENU／ADD CHARACTER／LOAD SAVEGAME／
  CONFIGURATION／IMPORT (BANE)／QUIT GAME`（BANE＝六代 import，lore 確認）。
- 結論：單人 1 級打 2 蟲必死（蟲一輪 12 傷）；r118 存點是主選單態，
  作戰仍用 r117（5335M，選單，活著）。
- R119：RUN 列掃描（262,158/166/174）＋右鍵關子選單＋hook 判讀。
