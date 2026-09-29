# R496：4 PHOOTS 遭遇＋RUN

## 前提
- 起點：`workplace/run/r495-maze.state`（Y，朝 side1-global；resume 22300M）。
- 遭遇將重放（RNG 存於態內，預期同一步數觸發）。
- RUN＝(262,166)（R121 既定）；成功字串 `THESUS RUNS AWAY!`→`RUNNING!`。

## 步驟（絕對步數！起點 22300000000）
 clicks: RUN(262,166)@22360M
 shots: x0@22320M, x1@22340M, x2@22380M, x3@22420M
 save 22420M:/w/r496-maze.state, steps 22480M
- x0/x1 確認遭遇選單重放；x2/x3 確認脫離回迷宮。
- RUN 失敗（無脫離字串／還在戰鬥）⟹ 停手，R497 打／再跑。

## 安全
- 存點 r495（遭遇前）已落；敗／散隊 → 讀回。
- THESUS 存活是完成定義。結束追記＋commit＋push。
