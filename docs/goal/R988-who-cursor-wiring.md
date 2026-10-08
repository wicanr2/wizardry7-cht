# R988：WHO游標正式中文接線

2026-10-09，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
上一輪R987來源與字樣證據47aa473及交接aa4081f已推送。

1. 以033已證24×20游標、文字來源尾圖與原生字樣建立正式唯讀顯示契約。
2. WHO?→誰？以正式UTF-8 TSV供字模生成，沿既有xlate字典及15／16字模，不加新資源schema或CLI。
3. 獨立契約審查READY後才接正式程式；HD／Draw共用布局，原圖手尾及圖外像素保持。
4. 正常冷LOAD、選人、移動、取消、再入，雙字型兩倍率驗證；固定原API對照原版狀態及保存層。
5. 字模缺失、來源不符、越界、未繪完、取消及圖案碰撞保原圖，不寫原版、離屏或A000。

研究輸出只留workplace/r988，原版、SDK與R984／R987輸入唯讀；Docker沿已驗image及UID1000:1000。
原圖、字模、堆疊、完整輸入與答案不提交。入口[033](../spec/033-world-identity-labels.md)、[R987](R987-who-call-witness.md)。
