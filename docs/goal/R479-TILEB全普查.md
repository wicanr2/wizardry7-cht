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
