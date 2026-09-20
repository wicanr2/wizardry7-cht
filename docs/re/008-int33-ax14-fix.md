# 008 — INT 33h AX=14h 修復：故事畫面可點擊（#13 結案）

日期：2026-09-20。issue #13 結案。前置：[007](007-story-wait-vgadrv.md)。

## 從 007 之後补上的拼圖（依序）

1. **poke 閂鎖驗證因果鏈**（confirmed）：直接寫 `cs:1ECD=1`（按鈕）、
   `1ECF=159`、`1ED1=47`（第一列選單中心）→ 畫面前進
   （d14fcc8c → 55,632 像素的新畫面）。閂鎖之後的鏈全通。
2. **滑鼠偵測其實過關**（勘正 007 的「init 沒跑」推論之一）：
   `sub_135A5`（DS.EXE file 0x37A5）四關全過——模式參數 [bp+4]=1（VGA）、
   `INT 21h AH=35h` 取回向量 0080:00CC、stub 首 byte `CD`≠`CF`、
   `int 33h ax=0` 回 AX=FFFF → `byte_12D66=FF`、旗標 0x1000 有設。
   （007 的懷疑點「golem 向量像 IRET 樁」不成立；期間修正兩次自己的
   線性位址換算錯誤：0080:00CC=線性 0x8CC、cs:12D66=線性 0x3E66。）
3. **driver init 在故事畫面開場時被呼叫**（confirmed）：stub 0x740C:0103 →
   init 入口 file 0x11F8（0x1200 是其續體）於 **#883,096,038** 執行
   （不是開機時；之前只盯開機 5M 步所以以為沒被呼叫）。
4. **init 完整跑完、註冊呼叫真的發生**（confirmed）：同輪 int33 功能號
   `0000/0004/0007/0008` 各 ×1、**`AX=0014×1`**——但 golem 把它記成
   未實作服務靜默吞掉 → handler 沒登記 → 事件不投遞 → 閂鎖永不設定。
   **根因定案：dosgolem 缺 INT 33h AX=0014h。**

## 修復（dosgolem 分支 `wizardry7/int33-ax14`）

- `docs/spec/194-int33-ax14-exchange.md`（READY）：交換語意照 DOSBox-X
  `src/ints/mouse.cpp:1970`——登記 ES:DX＋CX mask，回傳舊
  CX=mask／DX=off／ES=seg（初次全零）。
- `internal/dos/bios.go`：`case 0x0014`（約 15 行，複用既有 Handler 結構
  與 `fireMouseEventMickeys` 投遞）。
- `internal/dos/mouse_contract_test.go`：`TestInt33AX14ExchangesHandler`
  （首次回全零、登記生效、第二次回舊值、換掉後舊遮罩不再投遞）——通過。

## 端到端驗證（confirmed）

從 title-stable.state 跑到 1B 步、`-clicks "950000000:159:47"`（註冊點
#883M 之後才點）：

- `事件回呼：遮罩 000B handler 740C:1ED4，送出 2 次`（premove＋click）。
- 畫面：故事畫面 e141d21a（43,530）→ **142c2139（55,600）**——真實滑鼠
  點擊推進選單。與 poke 路徑的 2d250ae9（55,632）相位不同屬正常
  （事件序不同：hover 高亮先發生）。

## 對 #9 的解鎖

主選單檢查點的阻礙移除：按鍵腳本＋`-clicks`（配 `-click-premove`）可以
一路點過故事畫面往主選單走。注意點擊目標要用**熱點表**（`ds:59FE`，
(105,43)-(213,51)×5 列，驅動讀 `cx>>1` 後比較）。

## 遺留

- DOSBox-X 的輸入送達正對照（xdotool 焦點未驗證）——做 #9 的對拍時順便。
- dosgolem 分支未推上游（規則：推之前回報使用者）。
