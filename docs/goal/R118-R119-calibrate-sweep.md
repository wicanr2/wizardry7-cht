# R118–R119（輪 118–119）：校準＋RUN 列掃描

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R116–R117（hook 證明 START 無訂單＝no-op 重繪）。

## R118（FIGHT 校準，confirmed）

- FIGHT(26,171)@5340M → `THESUS-VS-`（訂單成立）；
  START(160,166)@5348M → `FINAL ORDERS／START FIGHTING／BACKUP TO PREV／
  TERMINATE GAME` 確認框 → 5390M frenzy（14490）＝回合執行。
- `THESUS PUNCHES BUG`→`HIT!`→`BODY 2 DAMAGE`；
  `BUG CLAWS THESUS`→`HIT!`→`BODY 12 DAMAGE`→`THESUS DIES!`→主選單。
- r118 存點＝主選單態（作戰仍用 r117）。

## R119（RUN 列掃描，confirmed）

- (262,158)＋START：僅選單重繪，無訂單無回合——158 非 RUN。
- (262,166)＋START：`-VS-`＋`FINAL ORDERS` → 回合執行；
  但蟲先攻 `BUG BITES THESUS`→`HIT!`→`HAND 17 DAMAGE`→`THESUS DIES!`。
  （262,166) 確實設定了某個訂單；是否為 RUN 無法區分——THESUS 沒活到
  自己的行動階段。先攻隨機（R112/R118 皆 THESUS 先攻，R119 蟲先攻）。
- 死亡 → 主選單（5440M 起 frenzy 43k–46k sustained＝主選單輪詢率，
  新 confirmed 基準：戰鬥選單 ~1k、迷宮 1–8k、主選單 ~46k）。
- 第三對點擊落到主選單：誤入 `IMPORT (BANE)` →
  `LOAD／CANCEL／DIRECTORY／C:\DSAVANT\／FILENAME／SAVEGAME.BCF`
  對話框（M2 讀檔功能未來可用；r119 存點＝此對話框態）。

## 判定

- 單人 THESUS（~10 HP）對 2 蟲（12–17 傷）：任何回合都有先攻死亡風險；
  逃脫需同時滿足 (a)(262,166)＝RUN (b)THESUS 先攻 (c)逃脫判定成功。
- 嘗試成本低（死＝讀 r117 重來）。R120 起逐輪重試同一格，
  hook 可區分（出現 THESUS 行動行＝166 非 RUN；出現逃脫字串／迷宮＝成功）。
