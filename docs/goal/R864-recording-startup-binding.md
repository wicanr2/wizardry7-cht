# R864：錄放起始設定與實際啟動的綁定

日期：2026-10-03。基線cc0db32；工作來源為[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)。

## 目標

承接R862候選，讓同一LoadWith實例、實際EXE、編譯來源、資源載入結果及兩層檔案清單共同產生設定。
先在私人副本驗證技術契約，移除只由呼叫端聲明pristine與任意DTO的綁定缺口。

## 範圍

- 013維持DRAFT；只寫workplace/r864候選，正式前端及SDK不改。
- 沿用實際資源載入器及014存檔守門。候選擁有者在LoadWith／SetScratch後擷取純量，完成設定前不交出Oracle。
- 以實際Go編譯套件列出本專案與SDK來源，連同模組鎖定檔及實際候選binary核對。不用舊RuntimeSHA。
- 原版目錄唯讀，只雜湊與LoadWith合法EXE；不Run、不解碼MSG、手冊或防拷答案，不載入或解碼原版state。
- 原版正文不輸出；資源內容及完整設定只留本機。自製守門測試不當原版替身。
- Docker採既有離線、非root、有界工具鏈；來源唯讀，只寫本輪研究目錄，完成後核對owner及容器清理。
- 故障紀錄及正常錄製結束仍待產品決策，不代選，不新增正式CLI或成功錄製聲明。

## 驗收

實際選定EXE、資源讀入位元組及載入結果可回查；owner擷取前沒有Run或輸入入口，釋出後不得重新擷取。
初始Root／Scratch與選定EXE載入前後一致；正常及明示加速／靜音設定如實區分，Step／Ticks仍0。
既有格式及純量測試不重開；只驗本輪新增的綁定、生命週期及拒絕條件。
獨立唯讀審查後回填013 DRAFT與研究022，驗證通過才推送及更新#17。

## 入口

- [目前狀態](../../CONTEXT.md)
- [013錄放草稿](../spec/013-input-recording.md)
- [022錄放證據](../re/022-input-recording-boundary.md)
- [R862候選](R862-recording-initial-config.md)

## 結果

12個綁定案例及1個不掛載Go來源的隔離案例通過。R865獨立審查阻擋0、應改0、建議0；
主代理完整讀取JSON及Markdown，只採認私人啟動設定與本輪文件回填，013維持DRAFT。
76來源中74項保持，013與本goal為批准文件變動；170封存產物保持。
審查JSON SHA-256為`7ad888ba516e01cfa827d38dcf7e2bf22801d0bbd2a44832efcd99b4835e7a76`，
Markdown為`095b481fd5c6c83e4b9c5e907baa77db24f07c54296b23810f813a856eb296da`。
根收尾入口為`workplace/r864/final-check.json`；不重跑產品或覆寫凍結索引。
完整錄放、正式冷存讀還原、十一組、TSV及封包仍未完成。
