# R22（輪 22）：probe 重跑——修 -type 轉義

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R21——probe 全程送達（149 clicks press/release 對全有），但 -type
  `$'AAAA\r'` 被轉義層吃成字面 7 bytes（`$AAAA\r`），防拷 AH=07 讀了 7 個
  垃圾字元，1.4B 步仍在防拷；迷宮段數據作廢；r21-landing.state 實為防拷狀態

## 目標

同 R21 命令，唯 -type 改從檔案讀精確 5 bytes（AAAA＋CR），避開轉義。
驗：防拷通過（IMPORT 對話框／DIRECTORY 或直達旁白）→著陸→迷宮點擊。
收：log-calls、PressReads、vram 比對、1420M 真著陸檢查點。

## 完成定義

1. 防拷通過（與前端 R9 同頁同答案）。
2. 到著陸＋檢查點為真（vram 與前端 mv-1380 同景）。
3. 迷宮點擊 vram 變化與否＋PressReads 呼叫者，定按鍵路徑。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 記錄

2026-09-28 實測：type 轉義修好了（AH=07 ×5），但 5 bytes 在 **7.8M 步被開機
讀者吃掉**（noteKey：`#7800858 int21-AH07 41 A` ×4＋`0D`），防拷照樣餓死
（1420M 仍是空框防拷）。完成定義全未達。教訓：-type 只能開場給，不適合
晚用的鍵；R23 改 -keys-at 四次定時投放（feedKeys 同餵 Stdin）。
詳見 `docs/re/010` R22 追記（併入 R23 前段）。
