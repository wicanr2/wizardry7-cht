# 020：前端音樂接線的最小證據

日期2026-10-01；[#18](https://github.com/wicanr2/wizardry7-cht/issues/18)；[R701](../goal/R701-audio-playability.md)。

## 輸入與位址基準

- 本案41e7e9d、dosgolem d16f03a2d8d782024d470606e1f4cb5272b89cd6；Go1.24.13、Python3.11.2。
- 原版DS.EXE SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`，本機唯讀來源及容器工作副本。
- image `sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
- 探針`workplace/r701/opl_probe.go` SHA-256 `23f9fec27cf6e649e4f72f67cd3657981820d4895bb68343a189cc68b5960550`；
  收據`startup-opl.json` SHA-256 `4d3158f270986a44f791390b026390720e006a7b88e436fd3d3803a1b16243a5`。
- OPL位址是硬體暫存器索引／bank，不是IDA線性位址；Step是dosgolem指令計數，不是毫秒或CPU週期。
- 本輪無反組譯。原版音樂事件、PCM或WAV只留本機，不隨發行包散布。

## 已證實

1. 現行前端開啟Adlib，但沒有聲音播放器或OPL合成接線；`cmd/wiz7/main.go` SHA-256
   `f4c51441b1a435d86116bf7699053c21672aae7e229d0828ccbd3e5a201fc0f7`。
2. 正常冷啟動、無輸入跑到260,000,000步，留下7,835筆OPL寫入；第一筆330,702、
   末筆231,196,498，步數遞增。988筆Key-On寫入、1,283筆Key-Off寫入；這不是音符邊緣計數。
3. bank只有0，節奏模式／CSM開啟0筆、波形只有正弦0；顫音致能146筆，震音／KSL致能0筆。
   只代表此啟動樣本，不能外推後續所有音樂。
4. indexed SHA-256 `e141d21a4d57998cfece1ff2e684383acfe9582702a5422a0030573f1030ee70`；
   Palette JSON SHA-256 `a7d075e62d9e7da4e727b0586d48cd05621f7e3824f6e09cd10c1d3963dd1158`。

## 既有合成器與限制

唯讀參考`/home/anr2/cht/psychic-war/worktrees/dosgolem`，版本
`f8c1a6ec5f069058f163ce8ddf345e24dc851de6`；其規格196為READY，合成器是連續時間近似YM3812。
提供9旋律聲道、FM／相加、回授、ADSR、波形0–3、KSL、震音、顫音；節奏模式不支援且有計數。

| 參考檔 | SHA-256 |
|---|---|
| internal/opl2/opl2.go | 18786416fc1ad060f6f936b2c4c422af0b78e5fa3b0a356ff88321176f9a84d8 |
| internal/opl2/render.go | d19232f79c31674c38a454dee47818c5c86dcf29f0faed1f1ed32cf9845508da |
| internal/opl2/opl2_test.go | 1e0c46b5c19b6c4f6fe0c5663dc6d5ccfdf35e810c959e0e1a4e721fada02d7b |
| docs/spec/196-opl2-synth-skeleton.md | 393c56dbe65dfc24bea81a35e52947d94b3b71be347be78e12de89747411bbfc |

採用合成器不等於升級整個dosgolem；不動CPU、DOS或PIT模型。通用能力放本案dosgolem分支，
既有oracle.go觀測API修改保留，不覆蓋或提交上游主工作區。
時間與硬體契約依`/home/anr2/.codex/knowledge-base/local/retro-platform-spec-first.md`。
事件到取樣的換算是播放時間模型，不能冒稱原版逐週期時鐘；若對不齊須回到規格，不能私改原版時鐘。
[Ebitengine 2.9.9音訊介面](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2@v2.9.9/audio)
支援無檔頭雙聲道PCM，串流取樣率須與context相同；裝置緩衝須獨立量測。

## 重跑與下一閘門

在既有Docker及本案根目錄，將原版唯讀來源複製到容器`/tmp/wiz7-audio-run`，
執行`go run workplace/r701/opl_probe.go workplace/rNEW/startup-opl.json`，收據名稱須不存在。
探針是本機研究檔，不是玩家入口；使用AGENTS.md第10節的資源、UID及清理契約。
播放器契約見[規格009草稿](../spec/009-player-audio.md)，正式接線尚未實作。

## R703隔離原型：時基與正常輸入

報告`workplace/r703/audio-prototype-report.json` SHA-256
`7b512b67e3faa7a3ddc5730d9c9ca25fe955dc68da8bd456598983f63b8a718a`，
189件輸出雜湊由`output-hashes.json`索引，原版音訊只留本機。
目前程式／唯讀runtime證實LoadWith未設定Clock：CPUHz33M、CycleClock=false、IRQ0Base165000。
PIT模型的名目每秒步數為`H=196875030000/65536≈3004074.5544`，
前端24M只控制牆鐘節流，不能作取樣時基；PIT分頻改變亦不改IRQ0Base。
公式依現有模型導出，仍是播放近似，不冒稱原版硬體逐週期parity。

無輸入冷啟動7,835事件，實際H重放3,816,817幀、約86.549秒；
全段與200k分段PCM逐位元相同，獨立ceil索引不符0、RMS0.05270。
原先24M假設477,750幀的收據保留，但該時長假設已撤銷。
正常Click／F1路線另有8,046事件，不能與無輸入冷啟動事件數混比；
有聲／靜音側Step、Mem、indexed與palette雜湊相同，沒有改原版時鐘或預設前端速度。

第二版按H合成卻沿24M節流，溢出1,545,640幀／欠載99,198；
第三版研究實際ΔStep/H節流，溢出88,050幀／欠載7,774。
第三版欠載全在一次6M步Click：執行0.779秒，返回前緩衝僅約0.574秒，
必須分段交付事件才能持續供應；未以擴容量或縮短Click隱藏問題。
原始取樣／輸出／丟棄／尚存索引守恆，兩版原始PCM相同；核心11項、封裝5項測試通過。

F1穩定暫停時機器Step、取樣索引、讀取位置與緩衝不前進；Pause保留已預讀資料，
裝置尾音約53ms只是監聽管線的粗估，不證明清空或精確延遲。
先前失焦假說被來源與實測否定：Ebitengine2.9.9預設背景執行為true，
已聚焦仍有Read停頓。相同虛擬裝置的30秒自有440Hz對照只重現1.828秒啟動停頓，
未重現遊戲版多次晚段停頓，後者仍是工作負載／後端交互未知。
各容器虛擬裝置擷取非靜音，不代表實體裝置或人耳通過。
自有容器已清理，UID／root-owned／誤建`.md`檢查皆通過。

R718續測有界進度水位與GUI鎖，只寫隔離副本；正式009仍DRAFT，
原型結果及未量支線須再審查後才可進正式路徑。

## R718隔離原型：交付水位、GUI與裝置阻塞

封存報告`workplace/r718/audio-progress-probe-report.json` SHA-256
`09a71bb323053602fecaf577aaf861ed78a4f8558107fd9a674949c00c31e3b2`；
827件輸出由`outputs.sha256.json`索引，索引SHA-256
`210b18df9274e60fc5dbb60925c72c3dd6b022710a8b1f72cff1f59bf1792c9d`。
建置image及Go版本同上；虛擬裝置用既有`eob-audio-capture:20260922-r2`，
image SHA-256 `90933bf64c453a75d307d0b1c2b591bb67aa1eca57b82e779859164b92f9cde6`。
本節只記隔離研究，不表示正式接線通過。

**confirmed：交付與時間軸。** 單一機器goroutine每完成200k步交付唯讀Step與新增事件副本，
包含最後餘數及錯誤後已推進的步數；保留原版Click的6M步數與輸入。
進度測試5項、核心11項與封裝5項通過。相對分段節流的V4欠載81,623幀，
其中Click11,145幀；改為唯一累計目標`epoch+Step/H`的V4b，欠載2,205幀、
溢出0、Click欠載0，但Run仍持GUI鎖，實際Draw僅8.743fps。

**confirmed：共享所有權的測試界線。** V4c解開Run期間GUI鎖，發布深拷字形、
背景、筆畫及矩形的快照，GUI自行建立圖片；Draw達38.458fps，
但欠載2,205幀、溢出1,563,643幀，不能把畫格改善當成聲音已修好。
真實GUI／F1／Click的競爭檢查退出0、無報告；執行766.819秒且欠載18,349,712幀，
與非競爭檢查版本分列，不作性能驗收。安全觀察epoch後此路線的textEv、flow及字形圖片均為0，
只量到F1中文，**一般文字與flow的GUI競爭分支未量到**。

**confirmed：裝置停頓定位。** V4d只加有界唯讀看門探針；三份停頓堆疊，
包含206.8M步的晚段，都停在Oto的`driver_unix.go:228`／`snd_pcm_writei`。
此版欠載0、溢出1,307,899幀、Draw45.260fps；底層ALSA／Pulse配置原因仍未知，
不是新的原版硬體逆向題目。唯一絕對目標沒有重設，僅F1恢復平移epoch；
V4c最大遲到16.046秒，後段追趕令Click的牆鐘時間短於名目2秒，原版6M步未改。

V4c／競爭檢查／V4d均為8,764事件、4,110,418原始PCM幀；獨立ceil索引錯誤0，
原始事件、PCM及有聲／靜音Step、記憶體、indexed、palette雜湊相同。
V4c停止仍有88,200幀、V4d有72,729幀未播；未驗完整清尾或Pause清空。
Pulse監聽非靜音仍不代表實體喇叭或人耳確認。

全部原型只寫r718；4,558項檔案擁有權均1000:1000，root-owned及誤建`.md`目錄0，
自有容器前綴清理後0。R728以同GUI負載自有純音及有限配置對照續查裝置環境，
不改H、原版輸入或擴緩衝掩飾溢出；正式009仍DRAFT。

## R728封存：驗證環境、啟動順序與真文字GUI

報告`workplace/r728/audio-environment-probe-report.json` SHA-256
`1e4313594dcb08fbce232ca45f685ed893e4a9917ab9bbd1843b1de16cfa9ebd`；
904件非快取輸出索引`outputs.sha256.json` SHA-256
`2c8fd2afe198f9d97616a5b1b3d8cad70f25d797560621537c1bb22dfa071045`。
工具、時基、取樣率、100ms預填及2秒環形緩衝均沿上述固定條件。

**confirmed：環境對照的界線。** 2CPU純音有配額節流及秒級讀取停頓；4CPU減少節流，
但相同280M遊戲仍欠載51,624幀、溢出531,269幀。容器以有界的全零`pacat`
維持虛擬輸出活躍，須確認設定與實際延遲均不超過20ms才啟動；這是測試環境，
不是正式玩家功能。只看設定值的第一次失敗保留。就緒純音60秒完整取得2,646,000幀，
仍有四次0.10–0.13秒間隙及三次ALSA重建，不能稱零停頓。

**confirmed：啟動先後。** 舊原型先呼叫播放，約0.414秒後才完成GUI資產、約0.417秒後
開始機器工作。新隔離版本僅把音訊初始化移至資產、觀察鉤及快照就緒之後；精確差異
保存`startup-v2-source.diff`。同一400M正常角色路線的單次控制，欠載13,792幀
（首交付前2,205、後續11,587）、溢出0、停止未播672；有效運行133.218秒，
虛擬時間133.152秒，Draw呼叫率83.00次／秒。不能將此改善宣稱全部音訊問題已解決。

`startup-v2-comparison.json` SHA-256
`7a58a7f3556c1ca43373694727c3099b52222a6169a08ec8f7263d0021c9d443`：
12,683事件及5,872,025原始PCM幀與舊400M版本完全相同，ceil索引錯誤0；
5,872,025＝已消耗5,871,353＋丟棄0＋尚存672。原版四項狀態雜湊相同，
360M／400M的A000、色盤、合成與PNG逐位元相同。全程虛擬輸出仍活躍，
ALSA建立一次、重建0；155次環形緩衝欠載恢復與ALSA重建是不同計數。

**confirmed：真文字GUI的範圍。** 舊固定GUI副本的400M靜音競爭檢查退出0、無報告，
實際文字事件10,729、字形10,628；原版狀態及兩張安全畫面相同。
新啟動版本未重跑含音訊競爭檢查；兩者不能合併宣稱。敘事頁數仍0，實體喇叭、
人耳及後期正常路線未驗。`callback_elapsed_seconds`起點實為上一水位，
包含CPU、觀察、快照與節流，只能稱水位間隔，不能當成純回呼成本。

全部原型只寫本機研究區；4,194項擁有權均1000:1000、root-owned與誤建`.md`目錄0，
自有容器清理後0，沒有整個DS、存檔內容或原版資料樹匯出。009保持DRAFT；
剩餘欠載分段、首次播放門檻、有界清尾及完整API仍須查證與獨立審查。
