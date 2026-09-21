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
