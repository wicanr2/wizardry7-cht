# R30（輪 30）：打黑鳥——戰鬥選項

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R29——金條放大讀出 FIGHT/SPELL/2 BLACKBIRDS/RUN/PARRY/USE/MOVE/EQUIP：
  老鷹戲是**戰鬥畫面**（2 隻黑鳥），不是過場；手形＝輪到行動的隊員；
  r29-maze.state（2480M，戰鬥中）可用

## 目標

-load-state 續跑：FIGHT（driver (70,165)，THESUS 行動）→數輪補點→
RUN 備選（driver (200,165)？右側 RUN 鈕，game x~250 → driver ~500？
先打左側再說，見機調整）＋vram＋存點。看戰鬥是否推進／結束回迷宮。

## 完成定義

1. 戰鬥有反應（敵人名數變化／傷害數字／行動輪轉／戰鬥結束）。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. 戰鬥不動：讀戰鬥 UI 機制（行動選擇流程、按鍵路徑），開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
