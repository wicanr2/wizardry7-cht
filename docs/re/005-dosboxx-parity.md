# 005 — 狀態檢查點與 DOSBox-X 對拍（M0 #5、#6）

日期：2026-09-20。issue #5、#6。工具：`tools/states.sh`、`tools/dosboxx-ref.sh`、
`tools/re/frame_cmp.py`。

## 檢查點（#5）

| 檢查點 | 步數 | vram SHA-256（前 12） | 狀態檔 |
|---|---:|---|---|
| banner（橫幅＋副標） | 200,000,000 | `69220ee76660` | banner.state（147,286 B） |
| title（標題動畫中） | 380,000,000 | `f0280f0997a5` | title.state（173,462 B） |

- 跨 run 決定性：banner 幀雜湊與先前獨立跑（run3 p200）**相同**（confirmed）。
- 恢復驗證（`tools/states.sh --check`）：從 title.state 載入、續跑 1M 步，
  畫面雜湊不變 → **通過**。載入到可實驗是秒級（400M 步的全程跑約 13 分鐘）。
- ⚠ 標題的**穩定相位**在更晚（~460M 步後、31,054 非零），380M 的 title 檢查點
  是動畫中幀。穩定相位的狀態檔留給 M1（按鍵腳本要再調：space 按太密會讓
  遊戲在 ~470M 提早退出）。主選單檢查點同樣待按鍵腳本（#5 的「主選單、
  第一個可操作畫面」項目依 #4 能走多遠調整——目前最遠到標題）。

## DOSBox-X 參照（#6）

設定（`tools/dosboxx-ref.sh`）：machine=svga_s3、cputype=386、cycles=fixed 20000、
oplmode=opl2（對齊 golem `-adlib`）、EMS 預設開（遊戲用 EMS）、scaler=none＋
aspect=false → 640×400（320×200 的 2 倍）。24 張快照每 5 秒一張。

### 開機時間軸對照（confirmed，兩邊的非黑像素數逐步吻合）

| 相位 | golem 非零 | DOSBox-X 非黑 | 快照 |
|---|---:|---:|---|
| 橫幅（部分） | 7,374 | 7,374 | title-01/02 |
| 橫幅完成 | 16,673 | 16,673 | title-04 |
| 副標出現 | 16,187 | 16,187 | title-05/06 |
| 標題動畫中 | 26,484→31,997 | 31,997 | title-07/08 |
| **穩定標題** | **31,054** | **31,054** | title-09–15 |

### 對拍結果

- 2×2 區塊一致性：**全數 24 張通過**（純整數倍放大，無插值）。
- **版面（遮罩）逐像素差異：0**（golem 穩定標題 vs DOSBox-X title-09/12）。
- 顏色（RGB）差異：**未驗**——golem 穩定相位沒伴隨色盤傾印（`-dump-palette`
  沒開），且動畫中幀的調色盤相位本來就對不上。屬 M1，比照 psychic-war
  的色盤工作另開。
- DOSBox-X 自身：title-09＝10＝15 逐像素相同（靜態相位、無時鐘隨機成分）。

### 方法備註（踩過的坑）

1. probe 的 `-shots` 不伴出 `.pal`（文件說會）——要色盤就用 `-dump-palette`。
2. probe 的 `-dump-palette` 輸出**已是 8-bit RGB**（首 bytes 0,255,255,227…），
   不要再做 6-bit 展開（`frame_cmp.py` 初版因此全幀誤判差異）。
3. python 比較 `bytes 切片 != tuple` 恆真——遮罩比較要先把切片轉 tuple
   （本輪除錯實錄：一個恆真比較讓遮罩全 True、差異恆 32,946）。
4. 動畫是調色盤／逐段推進，對拍要挑**同相位**：以非黑像素數對齊再比。

## 產物

- `workplace/states/`（banner/title 的 state＋vram＋meta.json＋run.log）
- `workplace/dosboxx/title-01..24.{rgb,x11.png}`
- 比對腳本 `tools/re/frame_cmp.py`（遮罩＋RGB雙軌）
