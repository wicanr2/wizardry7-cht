# R327：RHALLICK 選項表只有 TRUCE 有效 → 點 TRUCE 走停戰路脫離

## 發現（ram.bin 線性 ~0x20E00 選項表，12B/record = [id][flags][a][b][off][seg]）

| 選項 | id | flags | handler |
|---|---|---|---|
| TRUCE | 0038 | 0CC0 | **0110:406F（有效，seg 內唯一）** |
| STEAL | 0002 | FFFF | 0080:9E26（無效 seg = 禁用，FFFF 佐證） |
| LEAVE | 000F | 286F | F9C7:F293（無效 = 禁用） |

- 表窗內 `off+1001` 掃描只命中一處（20E24 = TRUCE），SPELL/USE/FIGHT 無 0110 handler。
- 0110 段 thunk 鏈（lcall *cs:2d8f → 740C:010B）是滑鼠驅動呼叫，追 ip-log 是死路，已放棄。
- R321 的 1 TRK = 時間觸發 RNG（與 FIGHT/LEAVE 無關）；R322-R324 menu locked = 點到禁用選項。

## 判準

- 在 r320-maze.state 點 TRUCE 熱區：遭遇結束（文字換頁或回迷宮）= 脫離成功；
  仍停戰畫面 = 再讀該頁文字找下一層選項。
