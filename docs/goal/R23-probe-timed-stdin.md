# R23（輪 23）：probe 定時投放 stdin

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R22——5 bytes 在 7.8M 步被開機讀者吃掉（noteKey 鐵證），防拷餓死；
  conIn 空佇列阻塞不 halt（前端晚打字能通）；feedKeys 同時餵 Stdin；
  -keys-at 支援定時＋`\\r` 轉義

## 目標

去 -type，改 -keys-at 四次定時投放（285M／350M／450M／600M 各 AAAA＋CR，
走檔案避轉義）。第一次若被開機阻塞者吃掉，後面的餵防拷（阻塞等鍵，
晚到也收）。餘量被名字欄誤食也無害（多一隊員 AAAA 即證據）。

## 完成定義

1. noteKey 顯示 AAAA 在防拷步數段被讀走（非 7.8M）。
2. 防拷通過→著陸→真檢查點（vram 對前端 mv）。
3. 迷宮點擊＋log-calls＋PressReads 定按鍵路徑。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
