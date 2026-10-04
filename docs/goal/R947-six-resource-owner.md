# R947 六類資源與單一啟動擁有者

2026-10-04；沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[013草稿](../spec/013-input-recording.md)及[研究022](../re/022-input-recording-boundary.md)。R946已驗同次讀取與六角色配置，本輪將其與R864／R875的單一啟動擁有者整合為私人候選。既有完整設定決定保持，正式格式與CLI不在此定案。

使用目前六份中文資源建立真正LoadWith初態，確認資源差異在LoadWith前拒絕、不可完整擷取時拒絕配置、初態差異在釋出指標前拒絕；比對擷取前後及原始控制初態的RAM／indexed／palette／暫存器摘要。配置只回傳副本，指標只交出一次。不Run、不送輸入、不載入state，不讀MSG正文、手冊或防拷答案。

原版與正式來源唯讀，只寫workplace/r947及本輪文件。建置、載入與核對使用既有離線、有界、非root Docker，Xvfb有界父程序擁有。驗收、獨立窄審後回填證據與目前真相、推送及更新#17；013保持DRAFT，不以load-only證據宣稱正式錄放或可玩交付。
