# 009：原版OPL事件的前端音樂播放

狀態READY；2026-10-02；[#18](https://github.com/wicanr2/wizardry7-cht/issues/18)。
證據[020](../re/020-player-audio.md)，工作入口[R701](../goal/R701-audio-playability.md)。
本版經R790逐條獨立審查，規格阻擋0；准許依契約正式實作，尚未CONFORMED。

目前正式接線及R1001正常玩家鏈已量：預設速度有聲冷LOAD、世界轉向／前進、正常SAVE及重開LOAD，兩趟164項通過、15對F1同Step恢復及正常關窗0。來源／幀／byte守恆成立，Overflow12092／48509及定案的停止捨尾照實保留，不稱零斷流。保存恢復viewport177408點差0；人耳試聽回覆仍待，未量平台、其他退出／加速及封包保持。完整範圍、工具與固定收據見[020的R1001](../re/020-player-audio.md#r1001-預設速度有聲冷load與世界操作)，不由此升格整份009。

## 為什麼要做

使用者要求完成中文化、打包並驗證可玩。前端未接聲音，原版已產生OPL事件；
#18第一階段為音樂可聽、時間戳與取樣可量，不要求完整保真調校。
2026-10-01使用者定案：預設依原版虛擬時鐘同步，保留手動加速；
加速時音樂節奏同步提高、維持音高；退出快速淡出、捨棄未播尾段並計數。
因此不採24M加速作預設、加速靜音或等整段尾音的方案。
原版CPU／PIT、遊戲規則、輸入與canonical 1×原始PCM不因播放或節流而改變；
速度時間軸生成的tempo PCM另按A2核對。
合成依固定版YM3812近似核心，屬硬體規格近似（hardware-spec approximation），
不聲稱原版逐週期、逐波形或實體喇叭保真，不深入原版硬體driver。

## 輸入與證據入口

原版輸入為Archives重印版DS.EXE，SHA-256
`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`；原版素材由玩家自備。
Go1.24.13、Ebitengine2.9.9、Oto3.4.0；Docker及固定合成核心來源／SHA見研究020。
定位以Go檔案／方法及原版Step為準，不將Go行號當IDA位址。
時基與PCM來源R703、固定倍率R753、原子性R766、精確相位R767／R771、
真正進度介面R770／R779／R784、消費端R778／R781及裝置鏈R783各自保留收據。
R790已確認契約充分；私人來源的移植應改須在正式實作及驗收時完成。

## 正式接線候選契約

### A1：通用能力與唯讀邊界

通用進度介面與固定合成核心置於dosgolem的隔離分支`wiz7-cht-opl-stream`，
不操作主機上游工作區、不升級CPU／PIT／DOS、不覆蓋既有觀測欄位修改。
巫術前端只負責已知時基、參數、GUI狀態與播放接線；通用合成器不內嵌本遊戲H。
原版EXE、資料、存檔、RAM、indexed與palette均不寫入；覆繪仍遵守spec001。

初始化採正值有理時基`H=N/D`、取樣率44100與正值有理速度`K`；複製配置及輸入，
不保存呼叫端可變別名。只允許Step0、空OPL trace的冷安裝；任意Oracle.Restore或暖重建不在本版。
原版玩家存讀檔由原版程式自然執行，Step持續前進，不重設合成器、不把PCM寫進原版存檔。

### A2：取樣與速度

本版LoadWith開Adlib、CycleClock=false、IRQ0Base=165000，
`H=196875030000/65536`步／秒；不能用前端-ips或當下PIT TickRate代替H。
未明示-ips時採此精確H；明示時沿既有Float64參數取得有限正值I，速度`K=I/H`。
無效、零、負數、NaN、Inf於LoadWith前拒絕；不新增快捷鍵或改原版時鐘。

正常速度的全域取樣n對應`floor(n*N/(44100*D))`步，先施加該步以前含端點的事件。
合成到T是`[0,T)`；總幀為`ceil(T*44100*D/N)`。
速度K=A/B時全域相位`P(T)=T*44100*D*B/(N*A)`，幀與事件索引均為`ceil(P(T))`；
初始化速度自Step0生效，不用丟幀或改取樣率實作加速。
如內部需要切換，僅於已完成水位T切換，保存該點精確有理相位，再以新斜率累加；
不能逐塊重設索引、每段先ceil或重播舊事件。暫存器順序與原始音高頻率不變。

採可驗證的整數／有理中間積，最後轉索引才做範圍檢查；
沒有研究用的全域1000萬幀／64段限制。過去已完成段不長期保存，
待施加事件保存其已計算的全域幀；失敗在提交核心／索引／Pending前拒絕整批。
若索引、Step或counter超可表示範圍，有限故障停止音訊，不移動／略過事件。

### A3：進度介面、供應者與唯一節流

公開候選沿已審原型的型別與方法：

```go
type CompletedOPL struct {
    Step, FirstSequence uint64
    Writes []OPLWrite
    Final bool
    OperationErrorCode, ObserverFaultCode string
}
func (o *Oracle) InstallOPLProgress(quantum uint64, fn func(CompletedOPL), opts ...OPLProgressOpt) (OPLProgressToken, error)
func WithOPLStream() OPLProgressOpt
func (o *Oracle) RemoveOPLProgress(token OPLProgressToken) error
func (o *Oracle) OPLProgressFault() string
```

quantum固定200000；註冊／移除僅於owner停機安全邊界、token須屬同Oracle與當前代。
正式前端採WithOPLStream：保留至多262144筆未交付來源事件、全域sequence及counter，
在成功dispatch後釋放已交付前綴，不能只移cursor仍保留整段m.OPL。
OPLRegs、CPU、PIT與原版記憶體不因研究trace退場而改變；OPLWrites在此明示模式只回尚存事件。
來源側append／copy之前做容量守門，超量記有限event_capacity觀察故障、停音並否定資格，
不改原版暫存器寫入或方法返回；故障後來源紀錄仍有界，不持續積存。
未啟用stream的既有完整trace API保持相容，本版不宣稱它有界；
本前端完整trace只允明示、有容量限制的本機診斷，超量不能假成功。
同一機器goroutine是唯一Oracle操作／讀取端；每完成200k步交付新增Writes獨立副本，
含「已收完Step≤T事件」水位與來源sequence；每次Run／RunUntil／整包Click返回另交付Final。
錯誤但有進度時仍交付已完成部分，原方法返回及後續輸入順序不變。
同Step維持來源順序；逆序、晚到、倒退、來源截短或溢位不得補移事件。
有限操作碼為空、budget、exit、no_response、reentrant、execution，不格式化Console或任意Error()。

回呼不可重入Oracle、ClearOPL、Restore或寫原版資料。只捕捉同步觀察回呼panic，
保存首個callback_panic並卸下listener；不吞原版CPU／條件函式panic，不保證Goexit或違約阻塞回呼。
每次原操作返回，owner必須讀fault；非空立即否定音訊／Final成功資格，nil原返回不等於音訊成功。

每批Writes在copy前限262144筆，總Pending亦限262144筆；
超量分別記來源event_capacity觀察故障或pending_capacity消費拒絕，音訊／錄放資格false。
先完整驗事件、ceil、容量與核心守門，再合成／交付至多4410幀的固定PCM塊；
不因極慢正值K而分配整批巨量PCM，也不暗中限制速度選項。
無效批次在任何輸出／核心／Pending提交前拒絕，不留下部分污染；
每個PCM塊間查explicit停止，中止剩餘合成記不完整，不能冒充原子拒絕或成功錄製。
每塊依全域輸出相位作中間絕對節流，最後一塊對準P(Step)；
有多塊時不先全塞2秒ring才等待，不以overflow丟低速音符宣稱同步。
交付後隨即釋放批次PCM與Applied審計；
只保留核心、精確相位、待施加事件、全域索引及counter。播放器Reader不碰Oracle。
唯一牆鐘目標為`epoch+P(Step)/44100`，F1恢復只平移暫停時長；
不疊加舊外層相對sleep，不用改原版CPU或溢出丟幀消除落後。
音訊consumer／觀察fault後，正常owner仍按相同epoch與速度節流；
若回呼已卸下，原操作返回時以同一絕對目標fallback Pace，不把故障變成遊戲加速。
只有明示關窗後跳過剩餘Pace；所有pacing timer可由停止要求喚醒，不等原sleep結束。

### A4：裝置輸出與有界緩衝

格式為44100Hz、signed16LE立體聲、左右相同、無檔頭。
100ms人工靜音預填、50ms播放器緩衝、2秒ring（88200幀）；
首200k來源水位完成且該速度第一個PCM塊入隊、完成其絕對節流後才首次Play。
正常H首塊2937幀；加速首塊依A2獨立ceil計算，不硬套2937；
低速多塊時首塊最多4410幀，預填加首塊不超過2秒ring，後續按A3分塊供應。
欠載補靜音、溢出丟最舊幀並計數；這些是裝置復原，不是加速算法。

復原淡入最多221幀、gain=k/220（k=0..220），只改裝置輸出、不改原始PCM／晶片包絡。
淡化乘積採整數half-away捨入，不能由特定浮點樣本通過外推所有int16。
Read允許任意byte長度，保留至多3bytes的已開始幀，不拆掉原立體聲相位。
原始canonical索引、tempo幀、裝置輸出幀分開計數；已預讀幀不能回寫淡化。
預設不開PCM檔或逐批完整日誌，不累積Spans／Drops陣列；
本機明示診斷模式另有容量／失敗標記，不拿可能阻塞的檔案I/O保證退出上界。

### A5：F1、靜音與裝置錯誤

依spec003由owner在原安全邊界承認F1，保存生效Step並凍結機器／輸入；
不重排已抽出佇列、不擅自截斷Click。Pause／Play只暫停裝置，保留尚未消費音訊。
暫停穩定後Step及合成索引不前進；恢復不重播／跳過取樣。
Pause前裝置預讀及可聽尾音另量，不能把Step凍結稱為喇叭立即無聲。
沿用／補上明示-mute，無頭驗收預設不開裝置；裝置失敗有限記錄，不能當有聲通過。

### A6：計數與權利

保存起迄Step、canonical／tempo／device幀、收到／施加／待施加／拒絕事件、
欠載／溢出、恢復、捨棄尾段、bank1／CSM／節奏模式未支援及退出畫格率。
固定核心Write只收Reg／Val，bank1與CSM守門拒絕並計數，節奏模式沿核心Unsupported計數；
量到0次只能稱這條路線未使用，不外推支援全部OPL功能。
Received指合成器合法接受，Received＋Rejected只代表送入合成器；來源Recorded與容量Dropped分開。
A7停止要求後原操作仍可能產生事件，前置停止分支未送合成器的事件另計，不能稱已合成或容量拒絕。
合成核心按dosgolem授權；原版音樂、PCM、完整輸入／state僅留本機，不進Git或可散布包。

### A7：快速退出與總期限

關窗／RunGame返回提出atomic停止要求；不直接從GUI操作Oracle或讀其狀態。
owner在下一200k回呼／原返回安全水位停止供應，凍結播放器佇列，
保留最多221幀快速淡出、其餘明確捨棄；短／空佇列按實際長度，不造新原始音訊。
先補完已開始的至多3bytes幀，再交付淡出；完整幀的Read ACK、EOF與實體出聲各自分開。
後續原操作不再合成／節流，保持原Run／Click返回；未完成操作不偽造成功錄放footer。

總退出期限為關窗要求後5秒；fade ACK最多等500ms且不能越過總期限。
總期限包含owner返回與Player.Close結果；不得持ring／player鎖等待owner。
Close可用受控單次工作goroutine，主程式select等待結果／總期限；
逾時記有限timeout、錄放資格false、不寫成功末註，main返回結束process。
這是退出失敗界線，不是假造操作完成或取消原版CPU規則；不得留下Xvfb／容器／外部工作程序。
Reader預設沒有檔案I/O，有限pure Read仍不外推所有排程有牆鐘保證。
Ebitengine Player.Close不代表audio context／ALSA硬體全部關閉的ACK。

原預設Click來源最多100M游標等待＋6M hover／hold／settle，僅為步數界線；
不由此或自製Run關窗約92ms推導普遍牆鐘上界。超過總期限的原操作不等到106M才退。
任意Watch／Cond無窮阻塞與玩家自訂不受控回呼不在此正式前端接線範圍。

## 已量與仍待驗

R784解除R777回呼panic截斷原操作的局部阻擋；R778／R781確認冷消費鏈與雙索引，
R781另保留淡出1LSB反例，修正在R783。R780有效敘事GUI所有權經R784確認，
只屬暖啟動旁證、不計正式冷錄放或音訊GUI整合；完整來源／界線見研究020／022。
R783裝置鏈與R790獨立審查均已封存。已量私人fixture不能直接搬其總幀數上限、
累積審計、同步檔案I/O或未受控Close為production；正式實作須遵守A1–A7並完成下面驗收。

## 正式實作驗收

- 自製：獨立整數／大有理oracle核對ceil、同Step順序、跨批原子拒絕、長索引、
  動態相位與1／2／4倍440／660Hz；整段／分塊逐取樣相同。
- 來源：原版冷啟動有聲／靜音同狀態，RAM／A000／palette／列明CPU純量相同，
  正常速度PCM與R703相同；固定原始事件只留本機。
- 玩家：真GUI有聲、F1凍結／恢復、有效文字／flow深拷、首次Play、正常及手動加速、
  Run與Click內關窗；明示欠載／溢出，不要求靠環境調參歸零。
- 收尾：短／空佇列、任意byte、精確淡化、ACK／EOF、owner／Close故障注入、
  總期限確實停止主程式；未完成錄放不能寫成功末註，正常存讀檔另驗不重設時間軸。
- 聲明：裝置非靜音與人耳分開；未聽不能稱人耳確認，Linux不能稱其他平台實機通過。

正式實作後依上述收據才標CONFORMED；未量玩家分支照實保留READY，
不得由綠色單測、某次92ms關窗或檔案存在宣稱全遊戲／封包完成。

## READY獨立審查

R790報告SHA-256 `e0acf7945ef00771127a8ef9d0c272cff1d4140b67910b697555ca230eb5c4af`，
42項本機索引SHA-256 `5b600114b63b11a22467ee498c0f28e3d71ed0c65037daf44380003b1004909d`；
受審契約SHA-256 `c62d633c37ac4498485329f4724898705151ff968901e1c969d872433521eece`。
逐A1–A7確認規格充分，阻擋0、需要猜補的未證原版語意0；四群組獨立競爭檢查通過。
正式移植應改：來源trace退場、固定容量／分塊、故障後節流、完整退出監督，
以及復原淡入也須採整數helper。R783停止淡出已修，但復原淡入仍有1LSB反例；
不修改舊收據、不稱私人原型已符合全部契約，正式實作後再驗CONFORMED。
