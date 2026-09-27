# R7（輪 7）：翻譯規模化——剩餘全文入袋

- issue：[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)（m4-trans-batch）
- 分期：M4（中文繪製與全文翻譯）主體收尾
- 前置：敘事批次一~五（story-trans 1,786 筆）、物品 417 件全譯、怪物 250 隻全譯

## 目標

messages.json 全量 11,019 則中，**要翻而未翻的 6,348 則全數入袋**。
依 `~/.claude/knowledge-base/workflows/batch-subagent-localization.md` 管線：
prep 切批 → 試作定調（給使用者看）→ 統一指令檔＋詞表 → subagent fan-out
→ 合併逐行驗證 → 全域一致性掃描 → bake＋端到端對帳。

## 盤點（tools/re/trans_audit.py，2026-09-27）

| 類 | 則數 | 處置 |
|---|---|---|
| 已翻（story-trans） | 1,786 | — |
| 要翻（todo） | 6,348 | 本輪全做 |
| 空白/裝飾/符號 | 836 | 不翻（無可見文字） |
| protected（手冊防拷 2500-2574、NLP 詞庫 7000-7197） | 177 | 不翻（追記三十一定案） |
| `*` 前綴（credits 人名/終極道具） | 40 | 不翻（追記三十一定案） |
| upper1（全大寫單詞：PS/AC/LB…） | 1,832 | **待核**：部分是 UI 縮寫該翻，部分是密文/代碼 |

todo 的分布：1000-1999 有 200（UI/選單）、8000-9999 有 210（NPC 對話）、
10000-13999 有 619、15000 之後 ~5,300（敘事主體）。

## 完成定義

1. `trans_audit.py` 重跑：todo = 0（upper1 分類定案後的淨剩＝0 或明列不翻原因）。
2. 每批合併時 en 欄位與 messages.json 逐字一致（gen_trans 既有檢查）＋ BIG5 可編。
3. bake 管線通過、字形數成長記錄。
4. 自動播放端到端對帳：觸發了卻不在文本檔＝0（維持追記三十七水準）。
5. 譯名一致：docs/names.md 更新，全域掃描無一人多名。
6. worklist/issue #22 進度更新、commit push。

## 不在範圍內

- upper1 的最終分類可順帶做，但若涉及新機制（如顯示格式）另開條目。
- M2 三項（試玩/錄放/音訊）、M5 輔助功能、M6 打包——「完成中文化」的後續輪次。
- 怪物 HP/XP runtime 對照（獨立小項，另做）。

## 記錄

### 2026-09-27 完成

**61 批全數入袋**（10 波，batch-005 試作＋波一~波十；期間遇 5h 額度用盡一次、
429 殺掉 4 批後補派）。最終數字：

| 指標 | 值 |
|---|---|
| story-trans | 8,825 筆（corpus 1,786＋批次 7,039） |
| trans_audit todo | **0**（不翻 2,194 全數有因） |
| bake | xlate 8,503 串、字形 2,874 |
| 端到端 | panic 0；無譯文＝credits/檔名/路徑類，該翻未翻 0 |
| 一致性 | 三輪掃描收斂 87 則；殘餘 123 個同 en 多譯屬上下文分流 |

**管線實戰驗證**（batch-subagent-localization 工作流的防線全數發揮作用）：
merge 計數檢查擋下 % 誤刪 52 則、^ 誤加 1 則；Write 跳脫轉壞由 agent round-trip
自檢兜住（batch-030/060 兩案）；同名異譯靠 en 條件式掃描收斂（非盲替換）。

**遺留**（M5/實測輪處理）：yesno 77＋四女巫名選項（`\x0N詞/`）上屏 vs 比對待實測；
`~` 前綴家族（~4/~7/~*）、`34}` 前綴渲染待確認；同 en 多譯 123 個為查表架構
context-blind 已知限制。

### 2026-09-27 盤點定案（trans_audit.py 細分後）

| 類 | 則數 | 處置 |
|---|---|---|
| todo | 6,418 | fan-out 批次 |
| abbrev（NPC名@/YOU@/OK/CANCEL/NONE/密文） | 621 | **併批次**：NPC名依 names.md＋保留@；YOU@→你@；OK→確定、CANCEL→取消；NONE→無；子音密文（HEDGHFLDFRM 類）保留原文記 notes |
| yesno（`\x0N(NO|YES)/`） | 77 | 隔離待查：疑似玩家回應詞庫（輸入比對用，同 7000-7197 性質）；翻了可能破壞比對。查法：反組譯對話樹或對照 DOSBox-X 實機 |
| ctrl/empty/protected/star/number | 2,116 | 不翻 |

控制碼開頭的句子（`\x1aHARDWARE CONFIGURATION`、`\x01...EXITING PROGRAM...`）
已併 todo——zh 保留前導控制碼即可（指令檔規則）。

### 管線

- trans_audit.py（對帳）→ trans_batch.py（切批）→ subagent（TRANS-INSTRUCTIONS.md）
  → trans_merge.py（逐行驗證：id 一致/控制碼次數/無!@開頭/BIG5）
  → text/batches-merged.json → gen_trans.py → story-trans.json → bake。
- 試作批 batch-005（NPC 對話 128 則，sonnet）進行中；確認風格後 fan-out。
