# R41（輪 41）：bios-keys 直寫 BDA 測戰鬥

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R40——-keys-at 在戰鬥中 stall（無 int09）；遊戲用 BIOS AH=01 輪詢，
  -bios-keys 才是正確通道。r40-maze.state（2635M）可用。

## 目標

-load-state r40 續跑到 2670M：-bios-keys「 RrF」（-bios-key-from 2640M，
-bios-key-every 5M）＋-log-calls＋shots。oracle：int16 AH=00 出現或畫面推進。

## 完成定義

1. 取鍵／推進：順著有效鍵深挖。
2. BDA 寫了仍無 AH=00：遊戲讀鍵走非常規路（直接掃埠 60h？），
   下一輪 -watch 埠或靜態找鍵盤常式。
3. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。
