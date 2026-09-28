# R10（輪 10）：步數門控 autoplay——START→巫6詢問→防拷→旁白→迷宮

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R9 步數錨點（主選單 246M→七項 260M→防拷 278M→DIRECTORY 289M）；
  wall-clock 漂移確認（R8 45s vs R9 ~60s）

## 目標

autoplay 改步數門控（門檻＝R9 錨點＋餘裕），路線：建隊（250–257M）→七項選單
START（265M，不繞 PCVW）→巫6匯入詢問 CANCEL（276M，語意待驗）→防拷 AAAA
（288M）→旁白 40 次（295M 起，間距 6M）→迷宮八點（540–568M，截圖＋9M）。
驗證巫6詢問 CANCEL 後是否進防拷→旁白→迷宮。

## 完成定義

1. 全跑無 panic；AUTO 行帶 steps（可驗門控精度）。
2. 若進迷宮且 mv-*.png 互異：#16 驗收門檻達成，記錄轉場步數。
3. 若 CANCEL 語意不如預期（回選單／卡住）：以翻譯事件定位，R11 換選項
   （LOAD／IMPORT）或檢查點法。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- #17 錄放；#18 音訊；前端存讀檔檢查點（本輪仍是全程重跑）。

## 記錄

2026-09-28 實測：63 AUTO 全執行、panic 0。步數門控命中全部關鍵屏：
七項選單 270M→START→防拷 276M（無巫6詢問）→AAAA 288M→引言到 1225M＋。
迷宮點擊（540–568M）打在引言中段（mv-568 星空頁存證）。完成定義 1–2 達成
（門控精度＋離開防拷）；進迷宮未達，R11 延長旁白到 ~1360M。
詳見 `docs/re/010-playthrough.md` R10 追記。
