# R945 錄製停止與續玩擁有者

2026-10-04；沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[013草稿](../spec/013-input-recording.md)及[研究022](../re/022-input-recording-boundary.md)。使用者已定案按鈕停止、當前操作完成後收尾、停止後續玩及故障保留未完成片段，不重新詢問。

本輪補正式規格前的技術契約：錄製截止與整個遊戲退出分開，最終同發布Draw確認與Writer單一擁有者依序完成，錄製停止不呼叫原版或音訊Stop／Close。以可撤回私人候選驗原有API結果保持、後續操作不入檔、軟體音訊續行、故障不發布成功末註，以及期限／晚到完成限制。

正式程式、SDK、原版、字型及已封存候選保持；輸入唯讀、私人輸出只寫workplace/r945，不讀MSG／手冊／答案。分析與測試沿離線、有界、非root Docker，完工核來源、UID/GID及自有容器。完整證據回填研究022、目前真相更新CONTEXT、歷程追加WORKLOG；獨立審查後才決定013能否READY，不以自製測試取代正常GUI或完整可玩驗收。
