# R461：TILE-B 完整普查

## 前提
- 起點：`workplace/run/r460-maze.state`（TILE-A，面 F0）。
- 目的：兩個 140 到 G2 後前進踩進 TILE-B（驗 15KB），再四個 180
  普查 TILE-B（同前綴 C4），存點 TILE-B。

## 步驟（絕對步數！起點 21800000000）
 clicks: 140@21810M, 140@21830M, fwd@21850M,
          180@21870M, 180@21890M, 180@21910M, 180@21930M
 shots: b0@21805M, b1@21820M, b2@21840M, b3@21860M,
        c0@21880M, c1@21900M, c2@21920M, c3@21940M
 save 21940M:/w/r461-maze.state, steps 21970M
- fwd 靜默（b3≈b2 無位移無字）⟹ 停，R462 重探，不普查。

## 安全
- 存點 r460（TILE-A）已落；TILE-B 若遇蟲戰→RUN 既定，敗／散隊 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。
