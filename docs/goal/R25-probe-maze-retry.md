# R25（輪 25）：probe 續跑過老鷹過場進迷宮

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R24——檢查點續跑秒級可用；r24-maze.state（1575M，老鷹過場中）；
  probe 旁白比前端長（点击不烧步数＋自动翻页差）；AX=4 在此段 ×1947

## 目標

-load-state 續跑：旁白 18 次（1580–1720M）→點底文→精確箭頭
（1725–1753M）＋vram 比對（對前端著陸幀）＋log-calls＋PressReads。

## 完成定義

1. 到乾淨迷宮（vram 同景）。
2. 任一動致視角變化：#16 驗收達成，關 issue。
3. 不動：PressReads＋AX=5/3 分工定按鍵路徑，開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
