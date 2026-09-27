# R8（輪 8）：試玩 IMPORT 快速路線驗證

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2（可遊玩前端）
- 前置：R7 翻譯入袋（story-trans 8,825 筆）；010 追記補二確認 START NEW GAME 全路徑，
  卡在巫6匯入詢問 CANCEL 語意；worktree 未提交修改把 autoplay 切到 IMPORT 快速路線

## 目標

驗證 worktree 未提交的 IMPORT 路線（`cmd/wiz7/main.go` autoplay 段）：
主選單 (121,97) IMPORT (BANE) → 防拷 AAAA → 載入畫面 → LOAD 舊進度進迷宮 →
旁白 40 次推進 → 迷宮方向盤八點移動，每點截圖 md5 互異（移動生效）。

## 完成定義

1. autopl2 全跑（build＋150s 窗口）無 panic，log 有 AUTO 序列到迷宮段。
2. mv-*.png 八點截圖 md5 不全同（至少出現移動轉場）；若仍停格，定位卡點
   （旁白未走完／載入畫面等輸入／時序重疊）並記錄步數與畫面。
3. 結果追記進 `docs/re/010-playthrough.md`；main.go 的處置（留／修／退回）有結論。
4. 不在範圍：issue #17 錄放、#18 音訊、巫6匯入 CANCEL 語意（START 路線另輪）。

## 記錄

2026-09-28 實測：autopl2 全跑（58 AUTO、panic 0），但 IMPORT 路線零命中——
16s 點擊落在 credits，主選單 +45s（246M steps）才出現；之後全部點擊打在主選單
上，mv-/hd- 截圖全同 md5。另發現 `workplace/run/` 無 SAVEGAME.BCF（IMPORT 無檔
可載）與 `import -window root` 截圖不可靠。詳見 `docs/re/010-playthrough.md`
R8 追記。完成定義 1–2 未達成（未進迷宮）；處置建議：退回 START NEW GAME 路線，
下一輪檢查點法解巫6匯入詢問＋按 +45s 錨點重排時序。
