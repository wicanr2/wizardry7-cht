# R40（輪 40）：戰鬥鍵盤隔離測試

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R39——滑鼠 press＋release 送達 (262,173) 但遊戲忽略 RUN；
  r39-maze.state（2590M，戰鬥 UI，游標在 RUN 上）可用。

## 目標

-load-state r39 續跑到 2630M：-keys-at 餵 R／r／空白／F（各隔 5M），
零點擊。oracle：log 出現 int16 AH=00（取鍵）或畫面推進。

## 完成定義

1. 任一鍵推進（畫面變／AH=00）：順著該鍵深挖，#16 繼續。
2. 四鍵皆無：戰鬥選單鍵鼠皆不吃，回頭查 VMELE/VMEXE 戰鬥迴圈（IDA 靜態）。
3. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
