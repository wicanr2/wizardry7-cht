# R1006：1.0.0本機完整版與推廣影片

2026-10-10，起點a82a936，沿[#47](https://github.com/wicanr2/wizardry7-cht/issues/47)。使用者明確要求先commit／push，再打包含遊戲的1.0.0完整版與錄製推廣影片。先行提交已完成；原中文化Goal仍暫停，R1004排版提問不由本次打包代答。

## 打包契約 CONFORMED（限定本機交付）

1. 版號固定`v.1.0.0-20261010`，使用專案環境明定的Asia/Taipei日期。已核本機及遠端沒有tag，不建立／移動／推送tag或公開Release。`dist-all/<版本>/full-local/`保存含原版的完整版；`promo/`保存影片與檢測；`smoke/`保存實際啟動收據；版本根保存SHA256SUMS.json。來源、建置與擷取中間產物只在workplace/r1006，輸出不存在才建立，失敗保留不覆寫。
2. 每個目標從當前受控正式來源乾淨複製建置，錄放來源身份沿build_recording.py真實選定來源。Linux應與已驗R1003正式97來源／6319 binary逐位元同。Windows amd64及macOS amd64／arm64用既有相容工具鏈交叉建置，不把產出存在視為可啟動。平台差異只在編譯目標及啟動包裝，不改遊戲／覆繪／字型／規則／資料格式。必要包裝器提供相同完整版號、固定包內路徑及轉送參數，存檔仍由正式系統使用者目錄或-save-dir管理。
3. 包內含原版Archives 248檔、原始PCFILE、正式執行檔、R995 Unifont xlate及三附屬家族、R983 flow、R979 help、RRSAL-1.0與Unifont／ARPHIC授權、中文啟動說明及已知限制。原版輸入逐檔SHA／size與自備ZIP基線相符；不帶R990進度、研究wire／片段、日誌、手冊、憑證、倚天字型或完整messages.json答案表。原版MSG.DBS仍是本機遊戲資料，不拆出答案語料。完整遊戲與影片不進Git、雲端或Release。
4. 逐包從乾淨目錄解壓，拒絕路徑穿越／symlink，核檔案集合、位元組／SHA、架構、版本與啟動路徑。Linux用包內launcher實際到主選單與正常關窗；再沿合法R990保存的獨立副本正常LOAD／操作／SAVE／退出／冷LOAD，原版包內容前後保持。Windows以Wine作實際啟動輔助驗收，不能冒稱Windows真機；macOS在無真機時核Mach-O、動態相依、架構、launcher及包完整性，明列原生啟動未驗。不能因此關閉#42或宣稱三平台原生驗收。
5. 影片以本輪包內正式Linux前端正常GUI有聲操作實際錄製，沿既有Xvfb／PulseAudio工具鏈，原版及R990合法保存唯讀；不注入記憶體、跳關、作弊或改亂數。防拷沿既有提示正常輸入，該段不進成片；答案與完整操作收據只留本機。可用既存已驗實機畫面作短轉場，但不得使用R1004原型。影片預計60至75秒，實錄片段為主，字幕只描述可見功能，不聲稱全遊戲／未完成作弊或新排版已完成。
6. 配樂為正式前端執行原版所產生的實際OPL輸出，錄到Pulse monitor；不以自寫近似合成、AI曲或假音效替代。成片核ffprobe音／視訊時長與格式、音量及非靜音、黑幀、凍結、抽樣幀字幕裁切與可辨讀性，留下來源SHA、rights metadata、contact sheet及輸出SHA。影片只本機保存，不執行公開上傳。
7. 工具、乾淨輸入manifest、驗收與已知限制先經獨立審查，主代理完整採認後才列交付。正常GUI、原版素材、字型、完整wire及影片留本機；Git只提交可重現腳本、規格／goal、版本與安全雜湊摘要。現行中文化剩餘問題保留在CONTEXT及包內說明，不把打包視為全目標完成。

本輪命中版號、dist-all、影片合成與原版音源路由；文件職責沿既有入口。Docker離線、UID1000、限CPU／記憶體／pids、有界Xvfb／Pulse／FFmpeg及結束清理照AGENTS。

獨立七條審查0／0／0完整讀回採認，JSON SHA`8b04fa4b7c8eae9502d698a60ced20493921f8bcb8a24d1f598d5f68b0656fb1`。只批准本機包裝與錄片契約，不擴成全遊戲完成或原生三平台驗收。

工具入口：[乾淨來源產生器](../../tools/re/prepare_release_build.py)、[本機完整版封包器](../../tools/re/build_full_local.py)。前者的`--out`只掛workplace/r1006/platform，後者的`--dist`只掛已存在的dist-all；原版、來源、SDK與固定字型授權皆唯讀。每個target使用新目錄；正式來源不因交叉建置重生source table。

安全來源紀錄為[輸入雜湊](../../tools/re/full_local_inputs.json)與[交付雜湊](../../tools/re/full_local_outputs.json)，只含版本、路徑、大小、SHA與驗收範圍，不含原版／字型／音訊bytes。封包器預設讀前者，`--inventory`可指向同內容的受控收據；本輪實際輸入鏡像與該檔逐位元相同。

macOS工具鏈補充來源為[osxcross-python.Dockerfile](../../tools/docker/osxcross-python.Dockerfile)。原有Go1.24.13／SDK15.5映像缺Python，新增明確revision `wizardry7-osxcross:go1.24.13-python3.13.15-15.5-r1`，沿本機已存在的兩個base離線建置，只COPY完整Docker Python runtime到/opt/python，保留Go與CC，不掛主機runtime。以`docker build --network none -f tools/docker/osxcross-python.Dockerfile -t wizardry7-osxcross:go1.24.13-python3.13.15-15.5-r1 .`重建，資源與逾時照AGENTS；精確base IDs記在Dockerfile，新image ID`1aa0f82b2c880d33751bf26259a6453d8c0009fd7f13e593db92897710aa9370`。已驗兩架構CGO、Python imports／相依與純建置，不代表macOS真機啟動。

[封包獨立核對器](../../tools/re/check_full_local.py)以`--archive`、固定`--inventory`及不存在的`--out`驗ZIP集合、逐檔SHA、原版／資源、架構、授權、UTF-8與版本，解壓到smoke新子目錄。結構PASS不等於原生啟動；正常GUI及各平台驗收另列。

[Linux正常玩家鏈核對器](../../tools/re/check_full_local_playtest.py)核本輪兩個包內實跑的正常關窗、實際SAVE新進度、冷LOAD保持保存metadata、世界viewport恢復與初態反對照，以及包內248原版檔保持。只證-save-dir覆寫路線與GUI viewport，不宣稱RAM同狀態或未量原生預設路徑。

正常錄影沿[控制器](../../tools/re/capture_full_local.py)、[Pulse／Xvfb包裝](../../tools/re/run_full_local_capture.sh)及[命令送入器](../../tools/re/send_full_local_capture.py)。包內根唯讀掛`/bundle`、合法起始保存唯讀掛`/initial`、全新擷取目錄掛`/out`；以普通X11鍵鼠及WM關窗操作，沒有私人Oracle／API。`--live`連續錄影時不自動開F1，`--record start/stop`只控制FFmpeg，防拷前後分開，成片只取menu／world兩實錄。每個程序有有界owner及外層910秒加20秒硬截止。

[成片合成器](../../tools/re/compose_full_local_promo.py)接收`--capture`、固定Unifont原檔`--font`、不存在的中間`--work`及版本promo目錄`--out`。只用實錄menu／world片段與其真實音訊，另生標題／字幕卡；留下ffprobe、音量／靜音、黑幀／凍結及抽樣contact sheet，初始狀態明列人工視覺採認待完成。

[實錄選段工具](../../tools/re/select_full_local_promo.py)只剪本輪120秒world實錄的0至22秒探索及85至119秒人物／裝備頁，刪去停留間隔；不重播或捏造遊戲，保留原音訊與來源SHA／剪輯區間。選段收據同時進合成器rights metadata，原始實錄保持本機。

本機交付最後獨立採認0／0／0完整讀回，JSON SHA`77df091626760dcb9c294447fb94b790a193f835cc7ec794ae823b3c0c1af7c6`。四包、Linux正常存讀檔鏈、Wine輔助與成片的限定範圍符合本契約；[交付採認](../../dist-all/v.1.0.0-20261010/smoke/release-artifact-acceptance.json)與SHA清單只留本機。原生Win／Mac、人耳及全中文化限制維持，不由本CONFORMED解除。
