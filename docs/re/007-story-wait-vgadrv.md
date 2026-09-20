# 007 — 故事畫面等待條件：VGA.DRV 的事件閂鎖與 INT 33h AX=14h（#13 進展）

日期：2026-09-20。issue #13（進展，未結）。前提與誤導线索見 [006](006-m1-progress.md)。

## 結論（依證據等級標注）

1. **段 0x740C 不是 overlay，是「驅動段」**（confirmed，16-gram 最長共同段）：
   `VGA.DRV` 載在 module+0x100（9,141 bytes 連續命中）、`PLAY.DRV` 在
   module+0x5410（2,684 bytes 命中）。**006 的「loader 有解壓/重組轉換」推論錯誤**
   ——對不上是因為那是驅動不是 `.OVR`；`.OVR` 載在別的段，無轉換之謎。
2. **故事畫面的等待函式＝VGA.DRV file 0x1317 的 UI 熱點服務**（confirmed，
   反組譯＋記憶體行為對上）：雙欄矩形熱點表（(105,43)-(213,51)×5 列）由
   呼叫端當參數傳入；命中回傳正 index（hover）、**點擊回傳負 index**。
3. **點擊的唯一路徑是事件閂鎖**（confirmed）：輪詢路徑把 `int 33h ax=3` 的
   BX 按鈕狀態**直接丟棄**（`xor bx,bx`）；只有 `cs:1ECD`（按鈕）＋`1ECF/1ED1`
   （x,y）閂鎖非零才走「已點擊」分支。**這解釋了「AX=3 看到按著 11,345 次仍無反應」。**
4. **閂鎖寫入者＝VGA.DRV file 0x1DD4 的事件處理常式**（confirmed）：
   `mov cs:[1ECD],bx; shr cx,1; mov cs:[1ECF],cx; mov cs:[1ED1],dx`，
   防重入旗標 byte_1ED3。它由 driver init 註冊：
5. **driver init（file 0x1200–0x1266）用 `INT 33h AX=0014h` 註冊**（confirmed，
   bytes `B8 14 00 … B9 0B 00 BA D4 1E CD 33`＝AX=14h、mask CX=0x0B、
   handler ES:DX=0x1ED4）。AX=14h 是 **MS Mouse v3.0+「EXCHANGE INTERRUPT
   SUBROUTINES」**——DOSBox-X 有實作（`src/ints/mouse.cpp:1970`：設
   sub_mask/sub_seg/sub_ofs 並回傳舊值）。
6. **這個 init 在 golem 從未執行**（強證據）：init 內的 int 33h 呼叫
   （AX=0/7/8/4/3/14）全部沒出現在 golem 的功能號統計（只有 DS.EXE 自己的
   0000×1，出處 DS.EXE file 0x37C1）。init 開頭 `cmp cs:[module 0x205],0 /
   jz skip`：旗標在檔案與執行期都是 0 → 自我跳過。**誰該設這個旗標＝未解**。
7. **golem 的 int 33h 沒有 AX=14h 分支**（confirmed，`internal/dos/bios.go`
   的 int33 switch 無 case）：即使 init 跑了，註冊也會掉在地上。
8. DOSBox-X 上故事畫面同樣不動（006）——與「init 沒跑」一致；
   但 DOSBox-X 的輸入送達未做正對照（xdotool 焦點未驗證），保留疑問。

## 因果鏈（現況理解）

```
故事畫面 43,530
  └ 等待：熱點服務回傳「負 index」＝點擊
      └ 唯一來源：cs:1ECD 閂鎖 ≠ 0
          └ 唯一寫入者：driver 事件處理常式（file 0x1DD4）
              └ 唯一註冊：driver init 的 INT 33h AX=14h（file 0x1257）
                  └ init 沒跑（旗標 module 0x205==0 自跳過）★斷點
                      + golem 沒實作 AX=14h ★第二斷點
```

## 下一步

1. **找 module 0x205 旗標的設定者**（候選：另一個 driver 服務先被呼叫；
   追遠呼進 driver stub 表的呼叫序）。也可能是「AX=0000 偵測到滑鼠後」才走
   的路徑（DS.EXE 0x37C1 的檢查後續未讀）。
2. **dosgolem 補 INT 33h AX=0014h**（依據 DOSBox-X `mouse.cpp:1970` 的交換語意；
   走 SDD：先寫 spec，分支上做，推上游前回報）＋事件投遞到交換來的 handler。
3. DOSBox-X 輸入正對照（在已知可互動的畫面驗證 xdotool 送鍵有效），讓
   oracle 比較公平。
4. 補驗：`.OVR` 實際載入段與載入位址（現在知道 0x740C 是驅動，overlay 段
   還沒定位過——用同樣的 16-gram 法對記憶體 dump 找）。

## 工具與產物

- `tools/ida/decode_at.py`（強制解碼指定位址/範圍）
- `workplace/ida/VGA.DRV.i64`（種子＋decode_at 後）
- 記憶體驅動段 dump：`workplace/run/ovr740c.bin`（32KB，其實是驅動段）
- 離線指紋法：E8/E9 (位置,目標) 對、rel16 多重集、16-gram 最長共同段
  （本次方法的教訓：**先確認段裝的是誰再談格式**——006 把驅動段誤當 overlay，
  差點推導出「loader 會轉換」的錯規格）
