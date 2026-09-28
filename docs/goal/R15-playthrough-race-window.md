# R15（輪 15）：搶 40K 空窗——Move 後立即按住

- issue：[#16](https://github.com/wicanr2/wizardry7-cht/issues/16)（m2-playthrough，OPEN）
- 分期：M2
- 前置：R14——park confirmed：迷宮 AX=4 釘 (302,164)（≈右下角落無效點），
  約每 40K 步一次；旁白無 park；範圍 [0-639,0-199]

## 目標

既然 park 週期約 40K 步，MoveMouse 後立即（hover 0）按住 20M：
按住期間第一次 park 前的 poll 若讀到注入座標＋按著，即移動。
同輪設 park 原地按（不 Move）為對照組（預期不動，證 park 點無效）。

## 完成定義

1. 全跑無 panic；AUTO滑鼠行可判讀。
2. race 致 mv 互異：#16 驗收達成，關 issue（後續把 race 寫成固定打法）。
3. race 滅＋對照組亦滅：park 空窗不存在／另有鎖定，下一步 RE 迷宮滑鼠處理
  （VMAZE.OVR 內命中比較＋40K 設游標的呼叫者）。
4. 追記＋commit＋push＋issue 留言。

## 不在範圍

- #17 錄放；#18 音訊；golem 改動（upstream 另案）。
