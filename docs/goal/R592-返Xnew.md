# R592：重返 Xnew

## 前提
- 起點：`workplace/run/r591-maze.state`（X2 朝 G1；resume 25470M）。
- T1＋T1（→G3）＋fwd（→Xnew，二次）；g1≈e1-R590（≤150B）⟹ 落 Xnew（確定性）。

## 步驟（絕對步數！起點 25470000000）
 clicks: 140@25490M, 140@25510M, fwd@25530M
 shots: g0@25480M, g1@25550M
 save 25550M:/w/r592-maze.state, steps 25580M

## 安全
- 存點 r591（X2）已落；若遇蟲戰→首輪 RUN，敗 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。
