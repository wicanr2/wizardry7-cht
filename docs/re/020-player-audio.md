# 020：前端音樂接線的最小證據

日期2026-10-01；[#18](https://github.com/wicanr2/wizardry7-cht/issues/18)；[R701](../goal/R701-audio-playability.md)。

## 輸入與位址基準

- 本案41e7e9d、dosgolem d16f03a2d8d782024d470606e1f4cb5272b89cd6；Go1.24.13、Python3.11.2。
- 原版DS.EXE SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`，本機唯讀來源及容器工作副本。
- image `sha256:1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。
- 探針`workplace/r701/opl_probe.go` SHA-256 `23f9fec27cf6e649e4f72f67cd3657981820d4895bb68343a189cc68b5960550`；
  收據`startup-opl.json` SHA-256 `4d3158f270986a44f791390b026390720e006a7b88e436fd3d3803a1b16243a5`。
- OPL位址是硬體暫存器索引／bank，不是IDA線性位址；Step是dosgolem指令計數，不是毫秒或CPU週期。
- 本輪無反組譯。原版音樂事件、PCM或WAV只留本機，不隨發行包散布。

## 已證實

1. 現行前端開啟Adlib，但沒有聲音播放器或OPL合成接線；`cmd/wiz7/main.go` SHA-256
   `f4c51441b1a435d86116bf7699053c21672aae7e229d0828ccbd3e5a201fc0f7`。
2. 正常冷啟動、無輸入跑到260,000,000步，留下7,835筆OPL寫入；第一筆330,702、
   末筆231,196,498，步數遞增。988筆Key-On寫入、1,283筆Key-Off寫入；這不是音符邊緣計數。
3. bank只有0，節奏模式／CSM開啟0筆、波形只有正弦0；顫音致能146筆，震音／KSL致能0筆。
   只代表此啟動樣本，不能外推後續所有音樂。
4. indexed SHA-256 `e141d21a4d57998cfece1ff2e684383acfe9582702a5422a0030573f1030ee70`；
   Palette JSON SHA-256 `a7d075e62d9e7da4e727b0586d48cd05621f7e3824f6e09cd10c1d3963dd1158`。

## 既有合成器與限制

唯讀參考`/home/anr2/cht/psychic-war/worktrees/dosgolem`，版本
`f8c1a6ec5f069058f163ce8ddf345e24dc851de6`；其規格196為READY，合成器是連續時間近似YM3812。
提供9旋律聲道、FM／相加、回授、ADSR、波形0–3、KSL、震音、顫音；節奏模式不支援且有計數。

| 參考檔 | SHA-256 |
|---|---|
| internal/opl2/opl2.go | 18786416fc1ad060f6f936b2c4c422af0b78e5fa3b0a356ff88321176f9a84d8 |
| internal/opl2/render.go | d19232f79c31674c38a454dee47818c5c86dcf29f0faed1f1ed32cf9845508da |
| internal/opl2/opl2_test.go | 1e0c46b5c19b6c4f6fe0c5663dc6d5ccfdf35e810c959e0e1a4e721fada02d7b |
| docs/spec/196-opl2-synth-skeleton.md | 393c56dbe65dfc24bea81a35e52947d94b3b71be347be78e12de89747411bbfc |

採用合成器不等於升級整個dosgolem；不動CPU、DOS或PIT模型。通用能力放本案dosgolem分支，
既有oracle.go觀測API修改保留，不覆蓋或提交上游主工作區。
時間與硬體契約依`/home/anr2/.codex/knowledge-base/local/retro-platform-spec-first.md`。
事件到取樣的換算是播放時間模型，不能冒稱原版逐週期時鐘；若對不齊須回到規格，不能私改原版時鐘。
[Ebitengine 2.9.9音訊介面](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2@v2.9.9/audio)
支援無檔頭雙聲道PCM，串流取樣率須與context相同；裝置緩衝須獨立量測。

## 重跑與下一閘門

在既有Docker及本案根目錄，將原版唯讀來源複製到容器`/tmp/wiz7-audio-run`，
執行`go run workplace/r701/opl_probe.go workplace/rNEW/startup-opl.json`，收據名稱須不存在。
探針是本機研究檔，不是玩家入口；使用AGENTS.md第10節的資源、UID及清理契約。
播放器契約見[規格009草稿](../spec/009-player-audio.md)，正式接線尚未實作。
