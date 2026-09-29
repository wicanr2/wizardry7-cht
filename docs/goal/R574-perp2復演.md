# R574：X5-perp2 乾淨復演錘

## 前提
- 起點：`workplace/run/r567-maze.state`（X5 朝 z1；resume 26140M；乾淨）。
- 三轉（→perp2，跳過入口不停）＋fwd：OUCH ⟹ 牆 ✓✓（R572 是 quirk）；
  靜默＋m2≈j2-R572 ⟹ 開（X7 確定性）；靜默＋m2≈j1（凍結）⟹ 牆。

## 步驟（絕對步數！起點 26140000000）
 clicks: 140@26160M, 140@26180M, 140@26200M, fwd@26220M
 shots: m0@26150M, m1@26210M, m2@26230M
 save 26230M:/w/r574-maze.state, steps 26260M

## 安全
- 存點 r567（X5）已落；若遇蟲戰→首輪 RUN，敗 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。
