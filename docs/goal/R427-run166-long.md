# R427（輪 427）：RUN(262,166)＋START＋長窗

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R426——(262,173)＋8M 時序仍只重印不開輪；R116 列掃描／
  R160 (262,166) 成功；r423（20025M）戰鬥選單。
- 目標：load-state r423 → 20200M：RUN (262,166)@20035M＋
  START (160,166)@20043M＋shots@20080/20130/20190M＋hook＋
  r427 存點@20190M。
  oracle：逃回迷宮（R428 續探）／戰鬥輪文字（R428 續戰；HP 危讀檔）／
  選單重印（R428 讀 R117 hook 解碼找 RUN 確認回音）。
## 結果（實測）［脫逃成功］

- EXIT=0；全文：`FINAL ORDERS／START FIGHTING／BACKUP TO PREV／
  TERMINATE GAME` → `THESUS RUNS AWAY!` → `RUNNING!` → 隊伍重繪。
  關鍵差異：RUN 列 Y＝166（173 未選上）＋長窗結算。
  r427（20190M）逃回迷宮，THESUS 存活。R428 驗收＋續北（同格重進，
  新骰）。

- 完成定義：開輪定性（THESUS 存活第一）；追記＋commit＋push（R424–R427 一批）。
- 不在範圍：改 golem 本體；#17/#18。
