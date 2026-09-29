# R571：X5 perp 側（一）

## 前提
- 起點：`workplace/run/r570-maze.state`（X5 朝 z1；resume 26170M）。
- T1（→perp）＋fwd：OUCH＝牆；靜默 ⟹ quirk／位移／吞（R572 往返／復演）；
  埋伏 ⟹ 首輪 RUN。

## 步驟（絕對步數！起點 26170000000）
 clicks: 140@26190M, fwd@26210M
 shots: i0@26180M, i1@26200M, i2@26220M
 save 26220M:/w/r571-maze.state, steps 26250M

## 安全
- 存點 r570（X5）已落；若遇蟲戰→首輪 RUN，敗 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R571 結果
- i0＝h1（41B）；T1 exact（18112B）；fwd→OUCH #26210468575＋105B ⟹ perp1＝牆 ✓✓。
- r571＝（X5，perp1；有 pending-bump）。R572：T1＋T1（跳入口）＋fwd 測 perp2。
