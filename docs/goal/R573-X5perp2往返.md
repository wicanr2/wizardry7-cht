# R573：X5-perp2 往返驗

## 前提
- 起點：`workplace/run/r572-maze.state`（X? 或 X5，朝 perp2；resume 26290M）。
- T1＋T1（→perp1）＋fwd；判 k1：≈j2-R572 ⟹ 新格開口往返（入口雙驗）；
  落 X5 既知視角 ⟹ 牆（R572 是 quirk，未動）。

## 步驟（絕對步數！起點 26290000000）
 clicks: 140@26310M, 140@26330M, fwd@26350M
 shots: k0@26300M, k1@26370M
 save 26370M:/w/r573-maze.state, steps 26400M

## 安全
- 存點 r572 已落；若遇蟲戰→首輪 RUN，敗 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R573 結果
- k0＝j2（0B）；k1 未回 j2（18104B opp 級）；k1≈g1-R567（141B）⟹ 落 X5 朝 perp1。
- R572 去／R573 回（X7 往返？）或 R572 quirk＋R573 解彈吞 perp1 牆。
  兩讀皆合（141B 殘差／凍結難分）。
- R574 乾淨復演錘：m2≈j2 ⟹ 開（X7 確定性）；m2≈j1（凍結）或 OUCH ⟹ 牆。
