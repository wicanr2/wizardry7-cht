# 023：adapter 分層與交付來源

2026-10-02；入口 [R823](../goal/R823-adapter-boundary-audit.md)；
議題 [＃29](https://github.com/wicanr2/wizardry7-cht/issues/29)。
本文件是目前來源的盤點，尚未選定搬移方案，也不是 READY 架構規格。

## 輸入、工具與定位

專案提交 `649d367a8c558e7588dc2811e1bd4a049ff4d677`；SDK 分支
`wiz7-cht-opl-stream`，提交 `d16f03a2d8d782024d470606e1f4cb5272b89cd6` 加本機修改。
採用既有 `wizardry7-go-ebiten:latest`，映像 SHA-256
`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`；
Go 1.24.13、Linux amd64、`CGO_ENABLED=1`，離線、唯讀來源。
定位全為目前 Go 檔案行號及匯入路徑，不是 IDA、DOS runtime 或原版檔案位移。
原始遊戲、手冊及防拷答案沒有作為本輪輸入；未執行原版或建置產物。

## 已證實的目前依賴

`go list -deps -json ./cmd/wiz7 ./cmd/xlate1` 成功：263 個套件，
其中標準庫 208、非標準庫 55；實際使用 9 個模組，9 個套件位於目前兩個本機工作樹。
其中 SDK 7 個套件、53 個編譯來源檔；正式前端 8 檔、歷史原型 1 檔。
這是當前 Linux 建置條件的集合，不能外推 Windows／macOS 或所有測試來源。

根 `go.mod:12` 的 `replace` 指向 `./worktrees/dosgolem`，不是可下載的完整版本。
SDK 有 13 個未提交檔案；其中 9 個被本輪兩個入口的依賴圖使用。
其餘 4 個是三份測試與 `opl2/SOURCE-LICENSE.txt`；不參與編譯不能當作交付時可忽略。
來源清單另含模組及 Docker 入口，共 71 項，各有檔名、大小及 SHA-256。
SDK 沒有 `go.sum`，清單明確記錄不存在；根目錄的 `go.sum` 保持原樣。

## 責任分類與回查入口

下表的檔案、匯入及函式定位為 confirmed；「可重用／應留遊戲層」的分層判斷為強推論，
不代表已對其他遊戲驗收或已批准拆包。

| 責任 | 目前入口 | 本輪判斷與界線 |
|---|---|---|
| CPU、DOS、檔案、記憶體與通用觀察 | SDK 的 `internal/cpu`、`internal/cpu386`、`internal/dosfile`、`internal/dos`、`internal/machine`、`oracle` | 已在通用層；不依巫術7 文字表做比較或存檔。沒有重驗整個 SDK。 |
| 時鐘設定及 OPL 事件完成通知 | SDK `oracle/clock_settings.go:10`、`oracle/opl_progress.go:51` | API 不帶巫術7 呼叫位址；屬通用能力。實際 3M 時脈與前端節奏選擇仍見規格 009。 |
| OPL 音訊串流 | SDK `opl2/stream.go:57`、`:261` 及 `opl2/opl2.go` | 合成器與事件串流已在 SDK；與音訊授權資料、測試來源一同追蹤。此處不新增音訊完成聲明。 |
| 巫術7 字串、印字返回及覆繪 | `cmd/wiz7/main.go:163`、`:228`、`:911`；`slots.go`、`flow.go`、`compound.go` | resident／驅動位址、字串情境及清除條件屬遊戲專屬，不能僅因移入 SDK 就視為通用。 |
| 前端視窗、說明、地圖及音訊操作 | `cmd/wiz7/main.go:571`、`:662`；`help.go`、`automap.go`、`audio_player.go`、`audio_reader.go` | 與遊戲及未匯出的 `game`／`overlay` 狀態耦合，目前不是可直接匯入的 adapter 套件。分離 API 尚未規格化。 |
| 歷史印字替換原型 | `cmd/xlate1/main.go:103`、`:105`、`:136` | callback 會寫回原入口並暫設 `ret`；不能當作現行唯讀覆繪實作，也不能與正式前端合併搬移後宣稱符合規格 001。 |

正式前端的印字觀察見 `main.go:911`，複合物品／flow／一般查表仍由現有規格仲裁。
`xlate1` 的歷史證據仍保留；本輪沒有執行它或改動其行為。

## 裸提交與完整受控來源的編譯對照

同一容器工具鏈、快取、前端來源與 Linux 條件：

- 目前 SDK：`go build -o /out/current-wiz7 ./cmd/wiz7` 成功，退出碼 0。
- 固定提交：以 `git archive d16f03a...` 的乾淨 SDK 及獨立暫存模組替換檔建置，
  退出碼 1；在 `audio_player.go:11` 因缺少 `dosgolem/opl2` 且模組唯讀而停止，沒有產物。
- 正式 `go.mod`／`go.sum` 未修改。乾淨 SDK 的 564 個封存檔另有完整雜湊清單。

confirmed：目前前端不能只依未套補丁的 SDK 提交重建；搬移 adapter 目錄本身不會增加音訊來源。
完整來源入口已由 README 指向 [prepare-dosgolem.sh](../../tools/docker/prepare-dosgolem.sh)、
[固定清單](../../tools/docker/dosgolem-audio-manifest.json) 及同目錄的受控補丁。

R825 在新的乾淨 SDK clone 執行原有準備器，13 個固定來源雜湊均符合；
559 個未被補丁修改的基線檔與 R823 封存相同，正式工作樹的 71 個已記錄來源亦保持不變。
完整受控來源建置成功，Linux 產物與目前 SDK 產物逐位元相同。
因此本機「未提交」與「沒有可重建來源」不能畫等號；完整準備流程的這個 Linux 範圍已通過，
不另列未解來源缺口。未驗全新映像檔／快取建置、原版試玩、聲音流暢、存讀檔或三平台封包。

## 尚未定案與本輪停止線

＃29 的搬移與否仍需使用者決定。若搬移，需先定義遊戲專屬 API、文字／字型資產入口及
前端模組邊界，再依 DRAFT → 獨立審查 → READY 實作；不能把本盤點直接當成搬移授權。
本輪沒有修改或推送 SDK，也沒有建立正式封包、tag 或 Release。

## 收據

本機 `workplace/r823/`：來源 71 項、固定提交封存 564 項、輸出索引 587 項。

- `report.json`：`df4d8f60db207fda0046b9899e35097b37089db80be7f9f4c62b317eeda1504b`
- `source-manifest.json`：`548ad6d9f254f9f87db79ba62f64083ed31edbf85c8a80a97e7eff84f6ff6b03`
- `pristine-source-manifest.json`：`2432ed74cf9361aa48f872a341ce0141109e9a875ebddd8a5dc1120385b27bce`
- `output-manifest.json`：`5cdcf7e9fa268d427ed9c154717010105898834788ba98d57f1aca4debc774f8`
- Linux 中間產物：`12c2e0dcfee0206fba3e4581aa2a0a7ec4f51374ee97af69f533b0f9d2692b1c`

廣義 `go list -m -json all` 因離線快取未齊失敗，保留空輸出；未假造缺少的模組資料。
映像檔查詢缺欄位與盤點器誤以為 SDK 必有 `go.sum` 均屬工具問題，已修正後重跑；
未套補丁的 SDK 建置失敗如實保留，以下述完整來源重建訂正其交付解讀。
獨立核對與容器／擁有權清理收據完成後追加於本節。

R825 訂正及完整流程收據：`workplace/r825/` 的輸出索引 580 項；
`report-final.json` SHA-256 `11e483dba491df4d9bc1aaaae14a1ab4f92241b3edd0138f2d0efcbc34536a78`，
`output-manifest.json` SHA-256 `602eb0bb869ba2f0af579bd9f0bbea122a952aed9c86183601126f39d00979c0`。
初次盤點把未套補丁的失敗外推到完整流程，已由上述準備器重建結果訂正；
R823 原始失敗收據仍完整保留，僅將其結論限定為裸提交對照。

R824 獨立核對最終阻擋 0、應改 0、建議 0；初次文字應改及圖中 compose／Draw 分工發現
均保留為已解除歷史。主代理完整讀取報告，另核對最終 94 項來源及 7 項輸出。
`correction-review.json` SHA-256 `c93c7af60c911757fc54352cb4d91661dc957d5ce9275c06a154ed2a298391a0`，
`correction-review.md` SHA-256 `85ad529c75f563a8e9ea05725282588ff984466c9a5c87b524a869dcd2785df4`，
來源清單 `12fd07522824da4f51f158956d98b3f28fe74a54add0a9733503202c3863cd0f`，
輸出清單 `ebb1ee22153f34f95e5173676c37e15ab765cf0a75a5fa5eb3363588bc74cf09`。
SVG 已核對 XML、可存取標題、箭頭與目前 API 分工；字型度量只為粗估，未驗中文像素渲染。
輸出均為 UID/GID 1000:1000；212,477 項中繼資料檢查的 root-owned／誤建 `.md` 目錄為 0。
自有分析與建置容器均採 `--rm`，收尾回查結果另記工作歷程。
