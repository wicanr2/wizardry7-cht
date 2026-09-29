# R492：side1 重演（定 (a)/(b)）

## 前提
- 起點：`workplace/run/r491-maze.state`（N8 朝 F；resume 22040M）。
- 判據（新截圖 t2）：t2≈s2-R491（17KB 級異）⟹ (b) 到 X；
  t2≈s1-R491（≤2KB）⟹ (a) 牆。

## 步驟（絕對步數！起點 22040000000）
 clicks: 140@22050M, fwd@22070M
 shots: t0@22045M, t1@22060M, t2@22080M
 save 22080M:/w/r492-maze.state, steps 22110M

## 安全
- 存點 r491 已落；若遇蟲戰→RUN 既定，敗／散隊 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R492 結果
- t2 vs s2-R491＝84B ⟹ (b)：side1 開口 → X，可重现。
- r492＝（X，朝 side1-global）。N8：遠端牆、side1→X、side2 未測、入口 N7。
- R493 普查 X：遠端先行（無 turn＝無 quirk，靜默即位移）。
