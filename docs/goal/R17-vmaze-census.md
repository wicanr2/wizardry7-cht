# R17（輪 17）：VMAZE 普查——int33 與命中表

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R16——doubling 靜態成立動態滅；缺 park 呼叫者；
  byte_228＝游標可見旗標；ida.sh＋locked-v1 在庫

## 目標

用現成 `census.py`（不寫新碼先）普查 VMAZE.code.i64：全部 int 指令
（位址、operand、bytes、所屬函式）＋函式表。找：int33 AX=3/4 呼叫點、
命中矩形比對形、position 服務呼叫（driver 服務號＋args 來源）。
不够再寫全文匯出。

## 完成定義

1. census JSON 存在、非空、輸入雜湊相符（ida.sh 契約）。
2. int33 呼叫點清單＋命中表候選（函式位址＋比對形狀）。
3. verdict 或明確的下一輪（全文匯出／probe 動態對照）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- 改 .i64 原件（副本跑）；改 golem；#17/#18。

## 記錄

2026-09-28 實測：`tools/ida.sh script census.py VMAZE.code.i64`（雜湊前後一致
`aa134943…`，原件未動）。VMAZE.code：109 函式、28.5KB code，**int 指令只有
1 個**（int3 @0x2480）——overlay 不直呼 int33，一律經 resident 服務。
golem 自帶文件（bios.go int33 注記）確認 doubling 是設計（640 虛擬寬）。
VGA.DRV 靜態：range 設定、halving 群、poll 迴圈、park 門（byte_228 可見旗標）。
決策：**全文匯出不需要**（無直接呼叫可找，省一次重跑）。
R18：park 呼叫者——probe 逮 AX=4 現行（含 stack 回溯）或 VMAZE far-call 掃描。
完成定義 1–2 達成；3 部分（verdict 待呼叫者）。詳見 `docs/re/010` R17 追記。
