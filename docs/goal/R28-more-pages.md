# R28（輪 28）：繼續翻頁——老鷹戲後段

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R27——2060M 仍老鷹戲；機制解釋：前端點擊燒 ~7M 步空轉（計時器幫忙翻頁），
  probe 瞬發一頁一頁，故同樣點擊數 probe 落後；r27-maze.state（2040M）可用

## 目標

-load-state 續跑：旁白 30 次（2045–2280M）→點底文→精確箭頭
（2285–2305M）＋vram（2100–2300M 每 40M）＋存點。

## 完成定義

1. 翻出老鷹戲（乾淨迷宮或新場景）。
2. 迷宮點擊致視角變化：#16 驗收達成，關 issue。
3. 不動：PressReads＋log-calls 定按鍵路徑，開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
