# R996：剩餘譯名與兩倍顯示缺口

2026-10-09，起點60f70fc，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。

1. 以現行正式catalog與既有安全出處重核剩餘異譯，不讀MSG防拷表或額外鄰文；保持已定BASELARD。
2. 未定詞只產生本機對照與可撤回提案，產品取捨依共同決策；不以last-write消除怪物23鍵衝突。
3. 窄核已量六個職業錨點在Unifont2倍的實際晚寫行，判定是否存在其他可見遮蓋；不由ALCHEMI四字外推其他字形。
4. 若新缺陷成立，先正常GUI驗證與DRAFT增補，獨立READY後才改正式路徑；維持原生字模及原版資料。
5. 每項結果寫目前狀態與唯一歷程；驗證後推送及回讀Issue，未量範圍照實保留。

入口：[008](../spec/008-english-key-tsv.md)、[031](../spec/031-monster-name-slots.md)、[036](../spec/036-party-class-short.md)、[025](../re/025-ui-localization-coverage.md#r995-full-party-ui)。
輸出只在ignored workplace/r996，Docker沿既有工具鏈、UID1000／network none及有界資源；原始輸入唯讀。
