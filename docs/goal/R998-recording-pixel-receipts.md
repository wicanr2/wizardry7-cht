# R998：三次重播的實際畫面收據

2026-10-10，起點d8fe063，沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。

1. 依[013 READY契約](../spec/013-input-recording.md#r951-production-contract)，在複製的現行來源中，於game.Draw已完成遊戲區／面板、既有Draw通知之前，只讀同份publication的遊戲區RGBA並計SHA。不得增加compose、輸入或原版記憶體寫入。
2. 觀測副本及來源以既有build_recording.py綁定，與正式5d6來源差異逐項核對。錄製與重播均用同一觀測binary，重新正常錄製，不改舊wire或繞過配置核對。
3. 每個publication只留一次Draw收據，包含ID、Step、panel、配置及實際binary身份、畫布尺寸與像素SHA；不把配置SHA當畫面SHA。遊戲區不包含錄製頁尾。
4. 正常六人組隊路線錄製一次、停止續玩，再用三份乾淨保存重播；逐publication比較三次實際畫面及原版發布摘要。全檔與收據核SHA、輸入及保存保持。
5. 私人觀測源碼先經獨立唯讀審查，再執行。結果不符時保留失敗並找真因，不刪動畫差異或改比較範圍來過關。充分故障測試不重跑。
6. 獨立採認後更新#17完成狀態及現況，整體中文化與平台交付仍依原目標續行。

入口：[R997](R997-recording-completion-audit.md)、[研究022](../re/022-input-recording-boundary.md#r997-recording-completion)、[目前狀態](../../CONTEXT.md)。
私人來源／binary、完整wire、圖片及保存只留ignored workplace/r998。沿既有Go1.24.13離線Docker、UID1000及有界資源；原版與原始資源唯讀。保留R433，正式Go與SDK不改。

可重建入口：[觀測來源產生器](../../tools/re/prepare_recording_pixel_observer.py)、[面板時刻修補](../../tools/re/update_recording_pixel_observer_panel.py)、[實際PNG核對器](../../tools/re/check_recording_pixel_receipts.py)。
容器以唯讀儲存庫掛/src、空且UID1000的研究輸出掛/out，依序執行`python3 /src/tools/re/prepare_recording_pixel_observer.py`、`python3 /src/tools/re/update_recording_pixel_observer_panel.py`，再執行`python3 /out/source-r4/tools/re/build_recording.py --output /out/wiz7-observer --receipt-dir /out/build`建置。
核對器限本輪640×400、正常加三重播布局，命令為`python3 /src/tools/re/check_recording_pixel_receipts.py --input <player目錄> --binary <同觀測binary> --output <新結果JSON>`；需既有Go映像的Pillow9.4.0，不讀或散布原版素材。
