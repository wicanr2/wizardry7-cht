# R630：2 BUGS 遭遇 RUN 脫離

## 前提
- 起點：`workplace/run/r629-maze.state`（遭遇選單開；resume 28260M）。
- BUGS 可逃（舊 R500/R501＋本系）；RUN(262,166)@28280M＋START(160,166)@28320M。
- THESUS 舊傷 7：新增 HIT 若致死 ⟹ 讀回 r626；成功 ⟹ 定位 t3。

## 步驟（絕對步數！起點 28260000000）
 clicks: 262,166@28280M, 160,166@28320M
 shots: t0@28270M, t1@28300M, t2@28340M, t3@28380M
 save 28380M:/w/r630-maze.state, steps 28410M

## 安全
- 存點 r629（遭遇選單）＋r626（X8）已落；敗 → 讀回 r626。
- THESUS 存活是完成定義。結束追記＋commit＋push。
