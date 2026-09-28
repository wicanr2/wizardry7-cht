# R44（輪 44）：START FIGHTING 開打

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R43——FIGHT 生效，FINAL ORDERS 確認畫面；
  r43-maze.state（2690M）可用。

## 目標

-load-state r43 續跑到 2730M：-clicks 2692M:START FIGHTING(160,166)
＋-log-calls＋2695/2700/2710/2720M shots＋2725M 存點。
oracle：回合動畫／傷害數字／敵方減員／下一輪選單任一。

## 完成定義

1. 回合結算：順勢續打或脫離，#16 續推。
2. 無反應：檢查按著計數，微調座標重來。
3. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測）

- START FIGHTING 生效，回合完整執行：
  2700M「THESUS PUNCHES／BLACKBIRD」（反擊），
  2720M「BLACKBIRD BITES／THESUS／MISSED!」（第二隻鳥落空）。
- 戰鬥機制全通（選動作→FINAL ORDERS→回合結算→下一輪）。
- r44-maze.state（2725M）已存。
- 下一步 R45：一次排三輪 FIGHT＋START FIGHTING，看鳥死不死。
