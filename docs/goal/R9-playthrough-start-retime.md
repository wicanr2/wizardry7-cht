# R9（輪 9）：START 路線按實測錨點重排時序

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R8 實測——IMPORT 路線零命中；錨點：主選單＝啟動 +45s（246M steps）；
  追記補二輪 4 的 START NEW GAME confirmed 路徑

## 目標

autoplay 切回 START 序列（主選單→建隊→角色選單→PCVW→EXIT→七項選單→START→
巫6匯入詢問 CANCEL→防拷→旁白→迷宮八點），全體按 +45s 錨點重排（舊時序＋32s，
窗口 150s→205s）。驗證路徑在 R7 後仍通，並抓到巫6匯入詢問畫面的到達步數
（給 R10 檢查點法用）。

## 完成定義

1. 全跑無 panic；log 翻譯事件出現防拷／旁白／匯入詢問文字（離開主選單的證據）。
2. 若到迷宮：mv-*.png 出現互異 md5（移動生效）即達 #16 驗收門檻。
3. 若仍卡：以翻譯事件＋idxNZ 轉折定位卡點畫面與步數，R10 改檢查點法。
4. 結果追記進 `docs/re/010-playthrough.md`；commit＋push＋issue 留言。

## 不在範圍

- 事件門控 autoplay（時序仍飄才做）；#17 錄放；#18 音訊。

## 記錄

2026-09-28 實測：65 AUTO 全執行、panic 0，大幅推進——建隊命中、THESUS
Fighter 自動建出、七項選單 (260M)→PCVW (266M)→主選單 (273M)→防拷 PAGE35 通過
(AAAA 被接受)→終點 DIRECTORY/FILENAME 載入畫面 (289M，mv-130 存證）。
拓樸修正：PCVW EXIT→主選單（非七項選單）。wall-clock R8 45s vs R9 ~60s，
R10 改步數門控。完成定義 1 達成（離開主選單＋防拷文字）；2 未達（未進迷宮）。
詳見 `docs/re/010-playthrough.md` R9 追記。
