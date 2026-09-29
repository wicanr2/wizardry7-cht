# R604：3 MOTHS 戰再 RUN

## 前提
- 起點：`workplace/run/r603-maze.state`（戰鬥中；resume 26380M）。
- THESUS 帶 5 傷；再 RUN(262,166)@26400M＋START(160,166)@26440M。
- 成功 ⟹ 定位 t3；再 BLOCKED／THESUS 亡／重傷 ⟹ 讀回 r601。

## 步驟（絕對步數！起點 26380000000）
 clicks: 262,166@26400M, 160,166@26440M
 shots: t0@26390M, t1@26420M, t2@26460M, t3@26500M
 save 26500M:/w/r604-maze.state, steps 26530M

## 安全
- 存點 r601（逃後點滿血）已落；敗／散隊／再擋 ⟹ 讀回 r601。
- THESUS 存活是完成定義。結束追記＋commit＋push。
