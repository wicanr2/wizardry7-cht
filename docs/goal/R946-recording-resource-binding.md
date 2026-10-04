# R946 錄放六類資源的起始綁定

2026-10-04；沿[#17](https://github.com/wicanr2/wizardry7-cht/issues/17)、[013草稿](../spec/013-input-recording.md)及[研究022](../re/022-input-recording-boundary.md)。使用者已定案單一錄放檔內嵌完整起始設定，本輪補齊目前正式前端新增的三類中文資源，不更改玩家選擇或正式格式。

核對xlate、flow、help、roster、profile-sex及intro-picture的實際載入條件，建立可丟棄的私人候選。每份資源由同一次讀取供解析與SHA核對；記錄請求、讀到及解析成功的差異，驗證失敗回退及重播拒絕不相符設定。候選不得直接接入正式程式，013保持DRAFT至契約與獨立審查完成。

只寫workplace/r946及本輪文件；原版、正式cmd、SDK、字型與既有封存保持。不讀MSG、手冊或防拷答案。分析、建置與測試在既有離線、有界、非root Docker內執行；原始輸入唯讀，快取及收據有明確輸出目錄。完成後核擁有權、自有容器與來源保持，更新目前真相、研究證據及GitHub議題。
