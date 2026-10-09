# R993：姓名輸入退格鍵

2026-10-09，沿[#46](https://github.com/wicanr2/wizardry7-cht/issues/46)，開工來源223ef27。

1. 沿R992已證缺口，只補前端Backspace到DOS字元佇列，不改原版、規則或文字資料。
2. 先用私人原型正常改名，確認ASCII08的刪字與空欄行為；原型不進正式路徑。
3. [035](../spec/035-backspace-input.md)經獨立審查READY後接線，保留F1暫停與既有掃描碼。
4. 正常退格、輸入新名、確認、重開名單與取消；正式錄製Stop後續玩、同binary重播及原API控制。
5. 兩字型名單姓名墨點與既有覆繪保持，原版只正常寫入scratch。採認後push並更新#46、#28。

入口：[R992證據](../re/025-ui-localization-coverage.md#r992-roster-format-source)、[035](../spec/035-backspace-input.md)。
輸出只在ignored workplace/r993；Docker沿既有image、UID1000、none及有界owner，原版與初始保存唯讀。
