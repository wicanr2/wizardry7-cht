# R13（輪 13）：迷宮輸入層診斷——Click error＋變體

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R12——著陸決定性重現，八點點擊零反應；input.go 確認 Click 回
  `NoResponseError`（附 Polls），按住太短會被 int33h 輪詢跳過；
  輸入模型：滑鼠＝移動（追記十四），無鍵盤快捷鍵

## 目標

同一輪內測迷宮輸入變體（前綴不變，著陸點不變），誰先讓 mv md5 變化誰是答案：

1. `cl()` 印出 Click error（含 Polls）——區分「座標錯／遊戲沒準備好」。
2. 變體序列（1430M 起，間距 4M，截圖＋9M）：中央視野 (160,100) plain →
   中央 Hold(5M) → 羅盤前進 (142,165) Hold(5M) → 前進 plain → 左轉 → 右轉 →
   後退 → 前進。

## 完成定義

1. 全跑無 panic；AUTO點擊行（error 有無＋Polls）可判讀。
2. 任一變體致 mv 互異：#16 驗收達成，記錄有效輸入＋轉場步數，關 issue。
3. 全滅：以 Polls 數＋畫面定位下一層（游標等待？熱點表？overlay 內命中邏輯）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- #17 錄放；#18 音訊；鍵盤輸入（強證據不支持）。
