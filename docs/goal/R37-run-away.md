# R37（輪 37）：讀檔回開打前點 RUN 脫離

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R36——FIGHT＋點鳥後底 UI 糊掉（2585M 起，疑結算寫壞／轉場卡住）；
  r30-maze.state（2540M，乾淨戰鬥 UI）可用；RUN 鈕在右側（game x~250）

## 目標

-load-state r30 續跑：RUN（driver (560,165)）×4＋vram＋存點。
看能否脫離戰鬥回迷宮（不解決戰鬥機制，繞過它）。

## 完成定義

1. 脫離成功（迷宮畫面）：#16 驗收繼續（移動）。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. RUN 無效：回頭解戰鬥結算（勝利條件／目標選擇），開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測）

- RUN×4（2545/2550/2555/2560M）**無效**：2555M 與 2575M 畫面皆為乾淨戰鬥 UI，
  敵方列完全相同「2 BLACKBIRDS (2)」；r37-mv0/mv1 md5 僅差老鷹待機姿勢。
- 戰鬥凍結：FIGHT 輪不結算、RUN 不脫離；選單（FIGHT/SPELL/PARRY/USE，
  RUN/MOVE/EQUIP）可讀，game 座標系。
- r37-maze.state（2600M）已存，可續跑。
- 下一步 R38：-log-calls 看戰鬥等待迴圈（INT16 鍵盤 vs INT33 滑鼠）。
