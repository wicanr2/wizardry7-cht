# R24（輪 24）：檢查點快進——probe 續跑進迷宮

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R23——定時 stdin 生效（285M 讀走）；probe 旁白比前端慢約 76M
  （點擊機制步數成本不同：前端每點燒 ~7M hover/hold/settle，probe 不燒，
  錨點不可跨工具沿用）；r23-landing.state（1420M，真敘事中，可用）；
  probe 著陸頁（gentle wind）約 1479M

## 目標

-load-state 續跑：旁白 18 次（1425–1560M，每 8M）→點底文→精確箭頭
（1570–1606M）＋vram 比對＋log-calls＋PressReads。只跑 200M，秒級。

## 完成定義

1. 到乾淨迷宮（vram 對前端）。
2. 任一動致 vram 變化：#16 驗收達成，關 issue。
3. 不動：PressReads 呼叫者 CS:IP＋log-calls 定按鍵路徑，開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
