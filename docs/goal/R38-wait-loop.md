# R38（輪 38）：log-calls 看戰鬥在等什麼

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R37——RUN×4 無效，戰鬥凍結在 2 BLACKBIRDS (2)；
  r30-maze.state（2540M，乾淨戰鬥 UI）可用。

## 目標

-load-state r30 續跑到 2590M，零輸入，開 -log-calls。
看戰鬥等待迴圈的特徵：INT16（鍵盤）主導還是 INT33（滑鼠）主導。

## 完成定義

1. 中斷統計拿到：能判定等待對象（鍵盤／滑鼠／皆非）。
2. 若等鍵盤：下一輪補 -bios-keys 或 -keys-at 變體。
3. 若等滑鼠：下一輪掃點擊（-sweep）找熱點。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測，r38.log，2540→2590M 零輸入）

- 服務呼叫僅 3 種：int33 AX=3 ×9146（滑鼠輪詢，約 1800／千萬道）、
  int16 AH=01 ×8070（鍵盤狀態查詢）、int1C ×304。
- 無 int10／int21／字型／開檔：戰鬥 UI 在「滑鼠＋鍵盤雙輪詢」等待迴圈，
  不是死機，是在等輸入。
- **R37 失敗原因找到**：-clicks 座標 (560,165) 超出 320×200 遊戲空間，
  點到螢幕外。RUN 鈕遊戲座標約 (240–285, 165–182)。
- 下一步 R39：遊戲座標 (262,173) 點 RUN，以「AX=3 看到按著」為 oracle。
