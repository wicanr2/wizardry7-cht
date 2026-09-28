# R34（輪 34）：噪聲屏等的是鍵盤

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R33——噪聲屏零滑鼠輪詢、首動後鍵盤阻塞退出；硬體鍵盤已完整
  （spec 012/014，010 缺口已關）；-keys-at 經 feedKeys 雙路投放
  （Stdin＋硬體掃描碼）；r32-maze.state（2610M）可用

## 目標

-load-state 續跑：-keys-at 送 X（2615M、2620M）＋vram（2625–2655M）。
看噪聲屏是否為等鍵盤（press any key）。

## 完成定義

1. 按鍵後畫面變化：觸發條件確認，繼續跟進。
2. 回迷宮＋移動：#16 驗收達成，關 issue。
3. 無反應：查該屏阻塞源（-watch ports／KeyPolls／AX 分工），開 READY spec。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
