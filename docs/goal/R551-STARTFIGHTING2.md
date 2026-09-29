# R551：START FIGHTING 脫離 BIRDS

## 前提
- 起點：`workplace/run/r550-maze.state`（FINAL ORDERS 開；resume 25610M）。
- START FIGHTING＝(160,166)；成功字串 `THESUS RUNS AWAY!`→`RUNNING!`。
- 步數放寬吃回合（R497/R500 經驗：回合結算需 ~150M）。

## 步驟（絕對步數！起點 25610000000）
 clicks: 160,166@25630M
 shots: s0@25620M, s1@25660M, s2@25700M
 save 25700M:/w/r551-maze.state, steps 25760M
- 未脫離（還在戰鬥）⟹ R552 首輪再 RUN；傷亡 ⟹ 讀回 r548。

## 安全
- 存點 r550（FINAL ORDERS）已落；THESUS 存活是完成定義。
- 結束追記＋commit＋push。

## R551 結果
- `BLACK BIRD CLAWS THESUS／HIT! HEAD 7／HIT! HAND 6／THESUS DIES!`→角色選單。
- THESUS 亡（此分支）；r551 態廢棄。存活線退回 r548（X3 迷宮）。
- 遭遇觀：R502（90M 空跑無怪）證非純計時；須迷宮點擊推進；
  同路徑同歩重放（R495/R496）；首輪 RUN 1/3（PHOOTS✓／BUGS✗／BIRDS✗）。
- R552：讀 r548，單轉短輪（25510M 收，避開 25512M 確定性點）。
