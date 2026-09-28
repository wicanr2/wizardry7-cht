# R27（輪 27）：點過老鷹戲進迷宮

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R26——1780M 仍是老鷹戲（著陸後還有一段，1900M 幀有轉場撕裂）；
  r26-maze.state（1880M）可用；前端是靠快點跳過此段的

## 目標

-load-state 續跑：旁白 20 次（1885–2040M）→點底文→精確箭頭
（2045–2065M）＋vram（1940–2060M）＋存點。找乾淨迷宮。

## 完成定義

1. 到乾淨迷宮（vram 同景）。
2. 任一動致視角變化：#16 驗收達成，關 issue。
3. 不動：PressReads＋log-calls 定按鍵路徑，開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
