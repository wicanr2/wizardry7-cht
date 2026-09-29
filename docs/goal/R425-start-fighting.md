# R425（輪 425）：START FIGHTING 確認開輪

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R424——r424（20095M）RUN 已選（選單靜止）；R114 序列：
  RUN(262,173)＋START FIGHTING(160,166)。
- 目標：load-state r424 → 20170M：START (160,166)@20105M＋
  shots@20125/20145/20165M＋hook＋r425 存點@20165M。
  oracle：逃回迷宮（RUN 成功，R426 續探）／戰鬥輪文字（RUN 失敗或接戰，
  R426 續戰；THESUS HP 危則讀檔）／選單不變（R426 換式）。
## 結果（實測）

- EXIT=0；START 生效（選單重印＋s1→s2 差 2690B＝高亮級變化，
  閃爍 500／換視角 5000 之間＝選項高亮）；但未開輪（無 ROUND／DAMAGE）。
- 疑因：R424 RUN 與 R425 START 相隔 65M（跨輪），選擇逾時或根本未選上。
  R426 單輪重做 R114 時序：RUN@T＋START@T+8M。

- 完成定義：開輪定性（THESUS 存活第一）；追記＋commit＋push（R424–R425 一批）。
- 不在範圍：改 golem 本體；#17/#18。
