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

## 結果（實測，r40.log，2600→2640M）

- 四鍵全未投遞：「硬體鍵盤送出 IRQ1 0 次，佇列還剩 4 個事件」。
  戰鬥 overlay 沒裝 int09（向量仍是 BIOS 預設 0080:0024），
  -keys-at 的 IRQ1 通路在戰鬥中 stall——R33–R35 的鍵盤嘗試同理全落空。
- 遊戲只用 BIOS int16 AH=01 輪詢 BDA：正確通道是 -bios-keys（直寫 0040:001E）。
- r40-maze.state（2635M）已存；畫面仍戰鬥（s0–s4 md5 皆異＝老鷹動畫）。
- 下一步 R41：-bios-keys 餵「 RrF」，看 AH=00 取鍵＋畫面。
