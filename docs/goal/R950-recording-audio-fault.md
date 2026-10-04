# R950 錄製音訊故障資格

2026-10-04；沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[013草稿](../spec/013-input-recording.md)、[009 READY](../spec/009-player-audio.md)及[研究022](../re/022-input-recording-boundary.md)。上一輪真實ENOSPC證據已推送，本輪補B3的音訊故障資格，不重跑完整三次重播。

先核R949關窗時的sink_failure與009 A7：明示退出會中止剩餘取樣，不能以最終串流不完整單獨推定錄製停止後播放已故障。新私人副本沿R948 revision-7，只加入可丟棄觀察回呼panic與取樣yield錯誤。兩者在9M後實際進入既有Oracle／Stream控制流；保持原API回傳，記有限故障純量，不讀正文。

真正冷GUI先正常控制，再兩種故障；在故障後要求正常錄製停止，後續遊戲仍可操作，片段無成功末註或正式名，原檔與改名副本皆拒絕重播。獨立核故障後Run仍以none返回且實際Step前進、配置／六資源／來源與原版248檔盲SHA。故障音訊不宣稱持續播放，裝置及人耳不在本輪。

只寫workplace/r950與本輪文件。正式程式、SDK、舊候選、原版與資源唯讀，禁止commit任何原版、片段或binary。離線有界非root Docker；先獨立來源審查再GUI、核對及結果審查。013仍DRAFT，正式getter／CLI／schema、位置、其他IO、撤回／阻塞與全部交付門檻保持。

來源預審0／0／0於FINAL後完整回讀採認。三個真正GUI案例與獨立核對通過；兩種故障後原Run仍成功前進，錄製均未成功收尾，15次續玩後片段保持，兩副本拒播。原版248檔盲SHA及R948舊來源保持，完整數字與限制見研究022本輪入口；沒有正式接線或交付完成聲明。
