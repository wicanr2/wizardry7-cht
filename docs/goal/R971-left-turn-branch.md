# R971：61,99 左轉分支探測（M4，#24）

- 日期：2026-10-06；前置：R970（82 點封存 SHA 3bd8387f…；Wait 清事件文後右轉入戰鬥；
  該點剩左轉未試；Y1／J1／J3 列下輪改善）。
- 授權：既有 29 出處×64B runtime 目的緩衝觀察（2026-10-03，不讀鄰文與防拷答案，
  完整收據僅本機）。擴大授權範圍尚未收到書面說明，視同維持，超範圍先停手。

## 目標

從 R970 已驗 Wait 後狀態續行，先試左轉（140,166），找 61,99 繞行分支；
若再入戰鬥即停。收斂 R970 下輪事項：Y1 正名、J1 停止原因機器收據、J3 命令列收據。

## 範圍

1. 精確重播 R970 封存前 81 點：命令取 commands.jsonl 第 1–77 行（該檔條目 SHA
   76962e62…，r970-new-01 止），捨棄第 78 行入戰鬥點與第 79 行 Finish 尾記；
   守門取 comparison.json 前 81 摘要／PNG 對（兩側一致）、report.json（命令／輸入／
   守門計數）、prefix-guards.json（逐點）各司其職，不符即停。
   len(commands)／len(inputs) 均不計 Finish 尾記。續行起點為兩側同狀態已驗點
   （82／82 PNG 對 sha-equal，review-result.json；捨棄末步續行有 R970 捨棄 R957／R958
   兩步之先例）。
2. 有界新操作：最多 6 次，只用前進 160,162、左轉 140,166、右轉 180,166、
   Wait（Run 上限 20M）；第一動為左轉。不送鍵盤文字，不點 NPC／選項。
   每次新輸入前分開檢視上一畫格；位置不動不重試同畫面同方向。
3. 停止：新選單／旁白／其他文字／可見控制／NPC／戰鬥、NoResponse、29 目標 ID 命中、
   手冊 2500–2574（只留有限數值並停止）、詞庫 7000–7197（排除正文與目的指標，
   64 筆數值及捨棄計數保持）、128 事件上限、6 次用盡。
   操作者先於探針日誌之外依 PNG 獨立判讀並記一筆可見特徵（如敵數／指令列），
   再填寫 stop.json（reason＋label／step＋feature），check.py 核對 label／step 與末快照一致
   且 reason ∈ 枚舉，否則拒收。停止→reason 對照唯一：新戰鬥→battle；新選單／可見控制→menu；
   新旁白／其他文字→text；29 目標命中→target；手冊停止→guard；詞庫 Cap→cap；
   128 事件上限／前綴守門不符→guard；新 NoResponse→noresponse；用滿 6 次無觸發→budget；
   未滿 6 次無觸發提前收尾（目標已達或策略止損）→finish。
4. 29 目標命中即取分流訊號；零命中不解除門檻，亦不推論相對關係不存在。
5. 不修改存檔、不模擬記憶體跳關、不讀鄰文、不碰 #26。
6. 留存 docker 命令列收據（run-log.jsonl：每次調用時間戳＋容器名＋命令＋退出碼）。

## 私人工具契約（READY：首輪 0 阻擋／3 應改／7 建議皆文字層級，已修完直接 READY；
報告 SHA 58f70145…，主代理完整回讀採認；S3 前綴點擊數按實測訂正為 64＋2 暖機）

- 前綴：r970 封存前 81 摘要（sealed SHA 3bd8387f…）＋77 命令重播，逐點守門，不符即停；
  合法 R733 暖初態 3475M，原 seed／時鐘／隊伍保持；不改 RAM／原版／存檔／亂數，不匯出新 state。
- 新命令上限 6（prefixIndex 封頂 77＋6＝83）；探針 panic 三重門（封頂＋標籤 r971-new-%02d＋形狀）。
  命令 schema 為 Label／X／Y／Steps／Wait／Finish；座標白名單見範圍 2 四操作；
  Steps 一律 20M（0 與 16.1M 僅前綴既有命令，不為新操作所用）。
  點擊預算 clickCount＞70（前綴命令點擊數預期值 64，實算不符即停；暖機 2 次不計入
  clickCount，inputs 66＝64＋2），寫死於 probe/main.go click budget 行；
  check.py 核對總命令 ≤83、新標籤連續、new_commands ≤6；Wait 每次 20M。
- 停止規則見範圍 3；觀察開／關兩側首次停止 Step 相同；控制側不裝印字回呼，不讀正文。
  Run 分 100k 區塊；守門／Cap（詞庫 64 筆保存上限）後無 PNG 者只核同停止與完整摘要，
  不稱逐像素通過。
- 29 來源／64B 目的緩衝／DS／時間／同指標／精確英文鍵守門沿 R970；SDK 基底 1132 檔
  對 sdk-before 逐檔核對，overlay 雙檔（checkpoint 72327b6d…、memory ab62152f…）SHA 釘選；
  不新增任何讀取。
- 收據欄位正名：`new_waits` 改 `new_commands`（公式 len(commands)−77）並增列拆分
  （new_clicks＝len(inputs)−66，等待數＝相減）；stop.json 由 check.py 驗 label／step。
- 基線：映像 wizardry7-go-ebiten:latest（核實際 ID 1d7a342d…）、Go1.24.13／Python3.11.2、
  現行選定來源、原版 248 盲 SHA 前後保持；R970 舊 206 封存獨立核對；Scratch 空、兩側退出 0。
- Docker 前綴 wiz7-r971-、非 root、離線、UID/GID 1000:1000、2CPU／2GiB／96pids、900 秒、
  log 10MiB×3；掛載前確認來源存在且型態正確；repo／原版（同時唯讀掛 /original 與 /w）／
  SDK 組裝目錄／模組快取／舊輸入唯讀，
  只有 workplace/r971 可寫；互動 -i＋PTY；收尾不留容器，事後檢查 root-owned 與誤建 .md 目錄。
- 原文／畫格／完整收據只留本機，逐檔 git add 限索引、SHA 與脫敏報告；wicanr2@gmail.com；繁體中文 commit。
- 正式程式、SDK、catalog、字型保持。結果獨立審查＋主代理完整回讀才封存提交；
  只限暖原版、NoCursorWait、非 GUI；008 仍 DRAFT，十一組門檻保持。

## 完成定義

- 探測收據僅本機；左轉結果（繞行／再戰鬥／其他）照實記錄，未達目標記位置與所需條件。
- 命中 29 來源者另立規格走審查；未命中維持 R969 裁決。
- 獨立唯讀審查 0 阻擋。

## 步驟

1. #24 留開題留言。2. 契約審查 READY 後重播前綴＋有界探測（收據新編號不覆寫）。
3. 結果審查封存。4. 提交推送更新 issue。
