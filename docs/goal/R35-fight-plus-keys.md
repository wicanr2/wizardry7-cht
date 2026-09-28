# R35（輪 35）：戰鬥存檔重打＋噪聲屏送鍵

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R33/R34——r32 存檔（轉場中途）續跑必死（2613.5M，EMS 未實作後跳飛，
  存讀保真問題）；r29 戰鬥存檔續跑正常；噪聲屏零滑鼠輪詢＋鍵盤阻塞退出；
  -keys-at 雙路投放已驗

## 目標

-load-state r29 續跑：FIGHT＋點鳥（2485–2510M）→噪聲屏出現後 -keys-at
送鍵（2565M 起四次）＋vram＋存點。答：噪聲屏是否等鍵盤。

## 完成定義

1. 鍵後畫面變化／落點可判讀。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. 存讀保真（轉場中途存檔必死）記 upstream 候選（golem issue＋分支另案，
   推上游前回報使用者）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem 本輪不動（只記錄）；#17/#18。
