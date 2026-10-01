# R816：前端Docker開發入口收尾

- 日期2026-10-02；程式基線f4871c2；工作[#43](https://github.com/wicanr2/wizardry7-cht/issues/43)。
- 狀態：工具入口驗收完成，`6656d6b`已推送；#43已關閉並回讀，非完整遊戲交付。
- 目標是修正README現行`tools/fe.sh`入口的主機工作負載、模組修改模式與逾時容器殘留。
- 沿用既有`wizardry7-go-ebiten:latest`及鎖版離線模組快取，不重建重複映像檔。
- 保留Go build／run與`WIZ7_FE_SH`用法；讓已記載的`xvfb '指令'`在容器內檢查與啟動Xvfb。
- 主機只做Docker控制面、參數及掛載型態檢查；快取建立、UID/GID檢查與命令執行留在容器。
- 原版ZIP及已存在的`workplace/original`額外唯讀掛載；所有可寫掛載先驗來源存在及擁有權。
- 每趟有唯一自有容器名、`--rm`、`--init`、資源／日誌上限、離線網路及有界外層逾時；EXIT／INT／TERM清理只動本趟容器。
- 不改原版輸入、顯示、音訊、存檔格式或預設位置；不新增主機X桌面映射。
- 驗收以自製fixture、現行Go build／Xvfb命令、失敗返回及逾時後容器0為準；由唯讀子代理獨立審查。
- 不跑原版或輸出手冊／答案／完整輸入，不把開發入口驗證當跨平台遊戲、存讀檔或中文化完成。
- README只保存穩定用法；結果／首工具失敗與清理追加WORKLOG，現況回填CONTEXT。

載入路由：規格閘門、README與文件職責；入口分別為
`~/.codex/knowledge-base/local/retro-remake-spec-gated-workflow.md`、
`~/.codex/knowledge-base/local/retro-remake-readme-standard.md`與
`~/.codex/knowledge-base/local/project-document-responsibilities.md`。
本案為不改遊戲行為的工具維護，沒有新增原版RE或玩家規則spec；獨立驗證不可略過。

控制面訊號依[GNU Bash的Signals契約](https://www.gnu.org/software/bash/manual/bash.html#Signals)：
等待前景命令可能延後trap，使用背景控制子程序與內建wait使INT／TERM可即時進入清理。
容器刪除亦有獨立5秒控制面逾時；這不是播放器A7的5秒玩家退出聲明。

## 本輪收據與範圍

R816實際Docker測試報告`workplace/r816/report.json` SHA-256
`69b987a8a75e7eb524bd4f87640a72c66ecfa22ea242d7e6587581cc834ca543`；
最新版來源清單`59f3279cba94874113173b46ff363656ae242c04e5caf25f8a3e66859f2ba9c6`；
11項輸出索引`ca570c0ec1bcdf852282e6495b2f8bddd337a42dcd6e97cfb4996c27140fa93c`。
Go1.24.13建置、建置完成後序列CLI及容器Xvfb 1024×768皆退出0；
模組快取／原版目錄／ZIP的唯讀掛載已實測，UID/GID1000。
原指令失敗碼17、非法逾時2、缺映像1均保留；2秒逾時fixture經5秒強制終止為137。
實際單次TERM退出143並移除同名容器；工具調度時間與write_stdin等待時間不當退出延遲。
248個原版SHA初末全同，沒有遊戲執行；所有輸出1000:1000、root-owned與誤建md目錄0。

獨立R817核對最新版訊號、名稱碰撞、查詢錯誤／逾時、原版來源型態與四個既有呼叫者；
自製fixture與實際Docker測試分開，不把假Docker登錄表當實際daemon證據。
前景wait造成的首清理延遲已保留，改為背景控制程序＋內建wait後同fixture解除。
尚未驗全新環境建置映像檔／離線快取；Dockerfile只修註解，沒有重建或聲稱新鎖版image。

R817最終獨立報告SHA-256
`9df40265228417bf75cc06c46317132849bb32fb7a6ba91464589e51a70b94b3`；
79項manifest `557b14a53d41c1ee8766b6001afba7450d7c136f067bb1f8f0eab8a709f07632`；
輸入索引`78e4cbad96c4600b2f9cde4d8732e18d816f8a16454d97e30ff4f9553d5daa87`；
輸出索引`8ae67ae0c7af212d00635d46be4c13c5c660fe4e4cf4708aaee4d89d64bf3237`。
14控制面與3路徑案例、34參數／名稱斷言失敗0，兩個修前對照保留。
功能阻擋0／必要修正0；保留一條映像查詢錯誤診斷文案建議，README已明列查詢失敗也停止。
來源無漂移、輸出1000:1000，專案root-owned／誤建md目錄與自有容器0。
