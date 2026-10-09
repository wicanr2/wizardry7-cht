# R990：女性人物頁正常玩家路線

2026-10-09，沿[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)。
R989正式錄放來源已驗並推送d7fb4f5，交接f2fa80a；已完成WHO與M路線不重開。

1. 沿034 READY及024角色選單規格，由正常角色選單尋找既有女性角色；沒有時才依原版正常建角建立本機測試進度。
2. 不修改state、角色bytes、座標或原版規則；合法起始保存層唯讀，僅本輪Scratch由原版正常寫入。
3. 從正常組隊與進入遊戲走到世界人物頁F，保存並冷LOAD接續；以真正GUI核字形、離開撤除及重入。
4. 能到達的路線正式錄製、同binary重播及原版控制另驗，未到達或未量分支據實保留。
5. 新來源或正式實作缺口回RE／DRAFT，不以猜測補production。測試、獨立複核、push及Issue回讀後更新目前狀態表。

輸出只留ignored workplace/r990；原版、既有存檔、字型及SDK唯讀，Docker有界、無網路、UID1000:1000。
完整答案、角色、存檔及輸入不提交。入口[034](../spec/034-world-profile-sex.md)、[024](../spec/024-character-review-output.md)、[013](../spec/013-input-recording.md)。

本輪結果入口：[女性正常玩家與原版控制證據](../re/025-ui-localization-coverage.md#r990-normal-female-profile)。
