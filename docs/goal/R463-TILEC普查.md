# R463：TILE-C 四側「轉＋撞」普查

## 前提
- 起點：`workplace/run/r462-maze.state`（TILE-C；TILE-B→TILE-C 連通）。
- 目的：R449 式四次（140 轉＋前進），每動必驗，畫出 TILE-C 四側。

## 步驟（絕對步數！起點 22040000000）
 clicks: 140@22050M, fwd@22070M, 140@22090M, fwd@22110M,
          140@22130M, fwd@22150M, 140@22170M, fwd@22190M
 shots: e0@22045M, e1@22060M, e2@22080M, e3@22100M, e4@22120M,
        e5@22140M, e6@22160M, e7@22180M, e8@22200M
 save 22200M:/w/r463-maze.state, steps 22230M

## 安全
- 存點 r462 已落；新格若遇蟲戰→RUN 既定，敗／散隊 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。
