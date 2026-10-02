# R847：存檔守門的平台中繼資料候選

日期2026-10-02；基線61d6070；[#42](https://github.com/wicanr2/wizardry7-cht/issues/42)仍OPEN。
入口[014](../spec/014-original-save-storage.md)、[研究022](../re/022-input-recording-boundary.md#r845-explicit-save-wiring)。

## 目標與邊界

- R807守門直接引用Linux syscall.Stat_t；在workplace/r847分離共用路徑與平台中繼資料檢查。
- Linux／macOS沿UID、權限與Nlink；Windows以官方檔案handle、owner SID及NumberOfLinks核對候選契約。
- 只讀官方Go／Microsoft API文件與鎖版模組原始碼，不重做原版RE。
- 不代選存檔預設位置，不改正式前端、SDK、catalog、字型或原版資料；014仍DRAFT。
- 使用自製檔案，核對相同Linux拒絕行為與跨平台編譯；Wine存在時只作輔助API驗證，不冒稱Windows／macOS真機。
- Windows無法取得必要中繼資料時保留拒絕，不偷偷移除擁有權或連結檢查。
- 不執行遊戲、不讀原版、MSG、state、RAM、手冊或autoplay輸入；不是存讀檔或可玩替代品。

## 驗收

保存輸入與候選雜湊；共用接線及平台後端差異經唯讀獨立審查。
Linux來源不變性、路徑重疊、一般檔案、符號／硬連結、不可寫及寫探針故障，以既有自製fixture核對。
macOS amd64／arm64與Windows amd64的獨立守門編譯，正式前端GUI與平台原生文件系統仍未驗。
工具限定已有Docker映像、有界、離線、非root；只寫明確輸出目錄，收尾檢查擁有權及自有容器。
新證據只回填014 DRAFT與既有研究022，不解除十一組、正常冷存讀、錄放或封包門檻。
