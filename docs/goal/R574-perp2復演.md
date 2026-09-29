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

## R574 結果
- m0＝g0（124B）；m1→m2 17045B（＝R572 同數）；m2≈j2（0B）⟹ 確定性位移，
  X5-perp2 開 → X7！（R572 是真位移，開側壓制。）
- X5 全閉：z1 前向牆、z3 入口→X4、perp1 牆、perp2→X7。
- r574＝（X7，perp2 到達方向）。R575 原地 C4。
