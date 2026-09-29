# R479：TILE-B（三向 junction）全普查

## 前提
- 起點：`workplace/run/r461-maze.state`（TILE-B；已知邊：r452／TILE-A／TILE-C，
  第四側未知）。
- 目的：R449 式四次（140 轉＋前進），靜默＋不動即判牆，畫完 TILE-B 四側，
  定三向 junction 全貌。

## 步驟（絕對步數！起點 21940000000）
 clicks: 140@21950M, fwd@21970M, 140@21990M, fwd@22010M,
          140@22030M, fwd@22050M, 140@22070M, fwd@22090M
 shots: s0@21945M, s1@21960M, s2@21980M, s3@22000M, s4@22020M,
        s5@22040M, s6@22060M, s7@22080M, s8@22100M
 save 22100M:/w/r479-maze.state, steps 22130M

## 安全
- 存點 r461 已落；若遇蟲戰→RUN 既定，敗／散隊 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R479 結果（2026-09-29 追加）
- T1 有效；fwd1 OUCH #21970541685（側一牆）；T2 有效；
  fwd2 無聲位移 18184B；T3 有效；fwd3 靜默 115B；T4 有效；
  fwd4 無聲位移 17334B；s8≈s0 179B **折返 TILE-B**。
- s4≈TILE-A-a0 79B ⟹ X1＝TILE-A。TILE-B 四側全滿：
  G1 牆／G2→TILE-A／G3→r452？／G4→TILE-C？
  （r452/TILE-C 归属 G3/G4 待定；R480 测 G3→r452）。
- THESUS 存活。
