# R18（輪 18）：反過來打——先按住＋連續搬移

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R12–R17——park 先讀後寫（應為先 park 後讀）贏了所有 Move-then-press；
  R15 race 連 transient 偏離都沒有；VMAZE 無直呼

## 目標

反轉順序：PressButton(0) 先按住（1434M）→ MoveMouse 前進 driver (284,165)
每 2M 一次 ×10（1436–1454M）→ 放開（1456M）。若 poll 讀 pos＋button 的時機
不在 park 正後方，必有一次讀到注入點。另首動先點頭像 driver (60,40)
（排查手形 modal 吃掉羅盤輸入）。

## 完成定義

1. 全跑無 panic；move 觀測行顯示 pos 是否釘住。
2. pos 釘住＋mv 互異：#16 驗收達成，關 issue。
3. pos 釘不住：park-first 順序 confirmed，R19 逮呼叫者（probe＋stack）。
4. pos 釘住但不動：modal（手形）成立，R19 排查解除條件（頭像／ESC／右鍵
   在正確座標）。
5. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 golem；#17/#18。

## 記錄

2026-09-28 實測：pos 全程釘住 (284,165)（10 次觀測），按住 22M，
park 18M 內僅 +1——送達雙全仍不動，park 出局。新證據：底文區域描述在屏
＋1403M 後零事件（疑非攔截管線）。完成定義 1 達成（觀測）；2–3 未達
（R19 查 modal）。詳見 `docs/re/010` R18 追記。
