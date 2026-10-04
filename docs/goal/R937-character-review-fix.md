# R937 人物檢視中文修正

2026-10-04；沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)及[R935證據](../re/025-ui-localization-coverage.md#r935-人物檢視與返回)，只修已證實的性別單字元漏譯與選人按鈕置中。[024](../spec/024-character-review-output.md)須獨立審查並經主代理完整回讀到READY才實作。

正式文案採TSV，生成兩後端原生字模；新事件只讀原版，恢復精確底圖紋理，原版規則與資料保持。正常GUI分別檢視既有THESUS及TEMPEST，右下箭頭返回，再離開回主選單及正常退出。雙字型／兩倍率固定收據核同狀態、HD／Draw、字模、字格及批准區外差異；2倍性別沿規格記為未完成，不縮字、不隱藏缺譯。

Docker沿既有image，非root、離線、有界資源及清理；原版、SDK、字型與R935等舊收據唯讀，新輸出只寫workplace/r937，既有檔拒絕覆寫。不按設定、不答防拷、不創建／刪除／改名角色、不作弊、不注入狀態或種子。完成正式測試、正常GUI與獨立複核後才逐檔提交、推送及更新兩議題，維持OPEN。025保存證據，CONTEXT保存目前真相，WORKLOG追加歷程；不宣稱全遊戲、錄放、冷載入、音訊或封包完成。
