# R42（輪 42）：阻塞點前餵鍵看吃哪個

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R41——遊戲在 #2640128703 阻塞等鍵；r40-maze.state（2635M）可用。

## 目標

-load-state r40 續跑到 2660M：-keys-at 餵「空白／Enter／R」
（2636/2637/2638M，阻塞點 #2640128703 之前）＋-log-calls＋2642/2648/2655M shots。
根因：遊戲阻塞讀取走 int21 AH=07（R41 實測 ×100000），讀的是 d.Stdin，
正是 feedKeys（-keys-at）寫入的通道；R40 的鍵只是送早了一個階段。
oracle：AH=07 讀到鍵＋shots 下一屏。

## 完成定義

1. 吃鍵＋推進：截圖判讀＋順勢續跑，#16 大步推進。
2. 仍阻塞：確認佇列狀態，查 bios-keys 投遞語義（讀 dosgolem 原始碼）。
3. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 結果（實測，r42.log，2635→2660M）

- 三鍵全被 AH=07 取走（#2636051638／#2637065841／#2638084297），
  無阻塞跑滿；但畫面仍是 2 BLACKBIRDS (2)，無推進。
- 结论：鍵送達但非正確輸入（被忽略或 flush 吃掉）；單點 RUN 不夠，
  戰鬥可能要 FIGHT→點目標完整序列（R31 座標越界，白測了）。
- r42-maze.state（2658M）已存。
- 下一步 R43：遊戲座標點 FIGHT (26,171) →shots→點鳥 (160,80)→shots。
