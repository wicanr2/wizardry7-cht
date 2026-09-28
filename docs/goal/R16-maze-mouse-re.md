# R16（輪 16）：座標 doubling 假說——x×2 再點一次

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R12–R15；R16 前段靜態確證 halving（見 010 R16 追記前段）

## 目標

靜態已確證：遊戲 SET 游標前 `shl cx,1`（DS.EXE:8675、VGA.DRV:2303/2642/
4249/4353），READ 後 `shr cx,1`（VGA.DRV:2309/2429/…），range [0-639,0-199]。
即 driver 座標＝邏輯座標×2。R12–R15 全滅的根因很可能是：注入的是邏輯座標，
遊戲讀到的是半值（(142,165)→(71,165)，羅盤左側空地）。

本輪：迷宮八點 x×2（y 不變）再跑一次。前綴不變，著陸點不變。

## 完成定義

1. 全跑無 panic。
2. mv 互異＋視角變化：#16 驗收達成，關 issue（doubling 寫成固定打法；
   選單路徑另行解釋——選單命中未 halving，待註記）。
3. 仍全同：回到靜態——park 呼叫者（byte_228/226 旗標、15AC/15AE 來源）、
   迷宮命中矩形表、按鍵讀取點。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17 錄放；#18 音訊。

## 記錄

2026-09-28 兩段：靜態確證 doubling（shl/shr＋range 行號見 010 R16 前段；
勘誤 park 非角落）；動態八點加倍全同 md5（`2cdc49…`），pos 照吸回，
park 速率本輪約 1/400K。doubling 必要非充分——缺 park 呼叫者。
byte_228＝游標可見旗標。R17 用 ida.sh 匯出 VMAZE 全文找呼叫點＋矩形表。
完成定義 1–2 達成（靜態）；3 未達（動態滅，轉 R17）。詳見 `docs/re/010`。
