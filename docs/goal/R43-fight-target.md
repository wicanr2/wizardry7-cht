# R43（輪 43）：FIGHT→點鳥完整序列

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R42——AH=07 吃鍵但無推進；R39 證明遊戲座標點擊可投遞；
  R31 的 FIGHT＋點鳥因越界白測。r42-maze.state（2658M）可用。

## 目標

-load-state r42 續跑到 2695M：-clicks 2660M:FIGHT(26,171)，
2670M:鳥(160,80)＋-log-calls＋2665/2675/2685M shots＋2690M 存點。
oracle：「看到按著」＞0 ＋ 目標選取反白／回合結算／敵方減員任一。

## 完成定義

1. 回合執行／敵方減員：戰鬥機制打通，#16 續推（脫離或打贏）。
2. 有選取反應但未結算：順勢補下一動。
3. 全無反應：戰鬥輸入另有機制，轉 IDA 靜態（VMELE/VMEXE）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
