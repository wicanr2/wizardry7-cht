# R618：2 PHOOTS 遭遇 RUN 脫離（二）

## 前提
- 起點：`workplace/run/r617-maze.state`（遭遇選單開；resume 27440M）。
- RUN 協議同 R600（成功前例）；THESUS 無傷。
- 脫離附帶位移，h3 需定位；b5-d3 仍未測（R619 續）。

## 步驟（絕對步數！起點 27440000000）
 clicks: 262,166@27460M, 160,166@27500M
 shots: h0@27450M, h1@27480M, h2@27520M, h3@27560M
 save 27560M:/w/r618-maze.state, steps 27590M

## 安全
- 存點 r617（遭遇選單）＋r616（b5 格）已落；敗 → 讀回 r616（b5 續測 d3）。
- THESUS 存活是完成定義。結束追記＋commit＋push。

## R618 結果
- RUN→FINAL ORDERS ✓；PHOOT 攻擊 MISSED；`THESUS RUNS AWAY!`→`RUNNING!` ⟹
  脫離成功，無傷 ✓。
- h3 無 corpus 匹配（最近 Z-b4 11KB）⟹ 又位移到未知視角。
- r618＝（逃後點2，未定位）。R619 原地 C4 試掛鉤（逃部位移若為鄰格當有面朝舊格）。
