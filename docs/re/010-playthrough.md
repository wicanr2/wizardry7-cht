# 010 — 模擬試玩：主選單 → 隊伍管理畫面（#16 進行中）

日期：2026-09-21。issue #16（進行中）。

## 已走到的位置

主選單（MAINSCRN，menu.state @940M）點第一列熱點 (159,47) →
**隊伍管理畫面**（PCVWSCRN 底圖，55,600 非零）：

- 底部按鈕列：CREATE｜REVIEW｜DELETE｜RENAME｜PICTURE｜EXIT
- 左側鏡框＝隊伍成員格（空）
- 場景：王座殘骸、金幣、骷髏戰士（PCVWSCRN.VGA 底圖）

## 已嘗試、未通的互動

| 嘗試 | 結果 |
|---|---|
| 點 CREATE 位置 (26,187)（含 premove/hold/settle） | 畫面不變（僅游標移動） |
| 掃描按鈕列與全螢幕（-sweep 104 點） | 每點都變（游標精靈），無法區分真熱點 |
| 鍵盤 C/R/D/E/Enter 掃描碼 | 畫面不變 |
| 不動滑鼠等待自動推進 | 停在隊伍管理畫面（1.2B 步以上） |

## 判讀（帶等級）

- 隊伍管理畫面的**按鈕命中測試不在 PLAY.DRV/VGA.DRV 的已知熱點表**
  （驅動段 32KB dump 掃不到矩形表模式）——命中邏輯應在**載入的 overlay
  模組**（PCVW 對應 VPCVW.OVR 或 VBASE）內部，或以指標傳入的另一張表。
- 下一步（輪 4 續）：
  1. 找出隊伍畫面對應的 overlay 模組與其載入段（開檔順序＋段位址）。
  2. 對該 overlay 反組譯，找主迴圈裡的滑鼠命中比較（`shr cx,1` 後的
     矩形比對形狀已在 VGA.DRV 0x1317 見過一次，overlay 內應有同形）。
  3. 或以 `-regs-at` 盯 0x740C:1417（熱點服務）在隊伍畫面的呼叫頻率，
     確認命中測試是否仍經驅動。

## 產物

- `workplace/run/pcvw1.png`（隊伍管理畫面）、`create.png`（CREATE 點擊後）
- `docs/images/wiz7-menu.png`、`wiz7-pcvw.png`（probe 確定性合成，README 用）

## 追記：鍵盤路徑判定（2026-09-21）

- BDA 環形緩衝驗屍：`-bios-keys` 的第一個鍵（K，word 0x1661）**停在
  0040:001E 沒被讀走**（head≠tail、後續鍵因 pending 閘沒推）→
  **名字欄不讀 BIOS 鍵盤佇列（int 16h）**。
- 遊戲也從未掛 int 09h（VecSets 0 次、IVT[9] 維持 0080:0024）。
- 判定（強證據）：**遊戲直接輪詢鍵盤控制器 port 60h（掃描碼）**——
  PLAY.DRV 掛的 INT 08h handler 或主迴圈裡讀。
- golem 的 `keyTick` 在「int 09h 向量＝stub」時刻意 stall IRQ1
  （`internal/machine/keyboard.go:105`，理由：stub 不送 EOI 會卡死
  後續中斷）——對 port 60h 輪詢的程式需要新投遞模式：
  **寫 kbdData（60h 可讀）但不發 IRQ1**。屬 dosgolem 擴充，需 spec＋
  EOI 分析，下一輪做。
- probe 的 `-press` 排程在 `-load-state` 前會被狀態還原覆蓋——已修
  （dosgolem master `d16f03a`：press 設定移到 load 之後）。
