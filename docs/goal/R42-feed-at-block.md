# R42（輪 42）：阻塞點前餵鍵看吃哪個

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R41——遊戲在 #2640128703 阻塞等鍵；r40-maze.state（2635M）可用。

## 目標

-load-state r40 續跑到 2650M：-bios-keys「空白EnterEnterRrF」
（-bios-key-from 2636M，-bios-key-every 1M）＋-log-calls＋2642/2645/2648M shots。
oracle：int16 AH=00 出現（吃鍵）＋shots 下一屏。

## 完成定義

1. 吃鍵＋推進：截圖判讀＋順勢續跑，#16 大步推進。
2. 仍阻塞：確認佇列狀態，查 bios-keys 投遞語義（讀 dosgolem 原始碼）。
3. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
