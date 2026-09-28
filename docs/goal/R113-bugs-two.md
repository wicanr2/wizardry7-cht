# R113（輪 113）：甲蟲兩輪

- issue：#16（已關；M2 迷宮探索）
- 分期：M2
- 前置：R112——第一輪中；r112-maze.state（5210M）可用。

## 目標

-load-state r112 續跑到 5290M：FIGHT(26,171)@5215M＋START(160,166)@5220M，
FIGHT@5260M＋START@5265M＋5218/5225/5240/5255/5268/5280M shots
＋5285M 存點。
oracle：回合結算／蟲減員／THESUS 血量任一。

## 完成定義

1. 推進：續打到贏。
2. THESUS 危險：評估 RUN 脫離。
3. 追記＋commit＋push。

## 不在範圍

- 改 golem 本體；#17/#18。

## 結果（實測）

- THESUS 戰死：s0 頭像旁骷髏，s1 黑屏過渡，s2 起回主選單
  （CHARACTER MENU／ADD／LOAD／CONFIG／IMPORT／QUIT）。
  單人 1 級打不過 2 蟲；FIGHT@5215 等點擊落在回合動畫無用。
- 教訓：深區遭遇先評估，RUN 脫離是必備技能。
- r111-maze.state（5170M，開打前）仍可用。
- 下一步 R114：讀 r111，RUN＋START FIGHTING 脫離。
