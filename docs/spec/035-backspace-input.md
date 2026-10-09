# 035 姓名輸入的退格字元

狀態：已量Linux3倍姓名輸入限定CONFORMED；其他畫面與平台未量，維持READY。R993最終獨立FINAL 0／0／0完整採認。入口[R993](../goal/R993-backspace-stdin.md)、[#46](https://github.com/wicanr2/wizardry7-cht/issues/46)。

## 問題與證據

正常改名時Backspace不能刪除預填姓名。R992兩次GUI及來源複核已確認：

- 前端`scanCode`有0Eh，但`Update`的DOS字元佇列只收ASCII32..126及Enter的0Dh。
- SDK `PressKey`只排IRQ1掃描碼；`Type`才追加DOS Stdin，兩路徑彼此獨立。
- 名字輸入走int21h AH07而不讀BIOS字元，沿[010](../re/010-playthrough.md)既有證據。
- 原版DS.EXE SHA-256為`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`，Archives版本。

以上缺口為confirmed。R993私人原型已正常量到ASCII08逐字刪除，七字清空及空欄額外退格保持空欄；
F1開啟時按退格及Q，關閉後THESUSR保持，才用正常新按退格刪為THESUS。
之後逐鍵輸入ROWTEST，Return及原版確定後人物名牌、ADD選取／取消均顯示新名。
上述範圍為confirmed；不外推其他輸入畫面或按住重複行為。證據見[025 R993](../re/025-ui-localization-coverage.md#r993-backspace-input)。

## 擬定契約

1. `Update`在既有F1說明開啟的早退之後，Backspace的JustPressed事件另追加一個byte08到`sh.stdin`，使用既有鎖。
2. 原掃描碼流程保持。只新增此控制字元，Enter及可列印字元的順序與路徑保持，不推定其他控制鍵。
3. owner沿原`o.Type`與錄放`type`事件派送；不直接修改原版姓名、字元緩衝、PCFILE、保存或VGA。
4. F1開啟時退格不得送到原版；關閉後新按退格才生效，不殘留暫停期間的輸入。
5. 以JustPressed維持既有按鍵事件契約，不新增按住重複、快捷鍵或前端文字編輯器。
6. 正常刪字與空欄多按一次退格，輸入合法新姓名、Return結束編輯、原版確定，再重開名單及取消。
7. 同binary完整重播與無覆繪原API沿相同完整輸入控制，全部發布原版六類摘要相同；保留正常改名的新PCFILE，SAVEGAME保持。
8. 兩字型3倍名單及姓名原墨點維持。2倍排版及022通用姓名覆繪仍由原規格處理，不因本修補宣告完成。

## READY與驗收

私人原型先量第6點及F1隔離，附binary／來源／完整收據SHA與正常畫面；獨立證據及契約審查後才READY。
正式接線後跑既有cmd／錄放回歸、正常GUI與正式重播、獨立原版控制；未實測分支保留限制。
本修補只恢復標準鍵的可用輸入，不變更玩家規則、資料格式或翻譯文案。完整輸入、原版、字型與保存只留本機。

私人完整錄製4,441事件／2,190快照，末Step490,200,000，頁尾Stop後續玩及關窗0。
wire SHA-256為`b5ee4abd64ff6d5a91ab8edd15c4843f8829a85a2bd079b87543af0e88e1fb2d`。
正常路線type byte08共8、Q共0，IRQ1 0Eh／8Eh各8，Q掃描碼0；正常影像另驗F1關閉後原名保持。
原型binary SHA-256為`d6fe340c924cef40a33798998d2d5086ef23ca45a5ee7a6740b17c2a989c4cd0`，Go1.24.13／Linux amd64、93來源。
正式改動須與已審原型五行完全相同；public main及重建93來源均核同才可把同binary原型收據納入正式範圍。
獨立READY審查JSON／Markdown SHA-256為`51a1cca2ff54895f5a8fff8f89d456dc8fde2edd275089834845e1bcf762c638`／
`149be16593073645c05e15e6a615174d38257964537ea2186cf0c57fe50e699e`。
七來源、13原圖、兩報告、wire及binary共24項身份核同後採認，只放行已證五行與後續正式驗收。

## 正式接線與身份

正式main僅新增原型的五行，SHA-256為`a889ea08d9ced6559f2c8bb08486a7a415393b9bfe840cf404d1b86e127fbaad`。
固定建置器由public來源重建，93項path／size／SHA清單與原型完全相同，兩binary逐bytes相同，仍為上列D6。
起始設定的Program身份不含frontend檔名或mtime，故完整正常原型錄製及同binary重播可依本READY契約採入正式倚天範圍，未改wire或檔頭。
前端cmd回歸2.356秒通過，initconfig與recording為有效快取結果。

## 正式限定驗收

倚天前述正常路線4,441事件／2,190快照與同D6重播完成0。正式Unifont正常一次退格ROWTEST→ROWTES、T復原，
原版確定後直接回角色選單，再ADD選取及取消；3,167事件／1,574快照，末367M、Stop後續玩及關窗0，同D6重播0。
兩份不同輸入各自原API有／無覆繪的2,190／1,574發布六類原版摘要與完整trace相同，共3,764對。
倚天保留合法ROWTEST改名的四份人物檔bytes／mode相同，mtime分開記錄；Unifont同名確認分支的四份人物檔全部metadata保持。
兩組SAVEGAME bytes／mode／mtime保持，保存只SHA/stat，未解資料格式或改寫收據。

既有唯讀frame helper另取姓名欄九個唯一狀態，原DOS `[150,210)×[175,181)`整數放大3倍，
29,160個RGB像素與原indexed／palette的獨立期望相同。此為衍生HD姓名域，正常GUI另列，不宣稱live Draw全圖驗收。
最終獨立代理以Pillow另重建原indexed／palette期望，九樣本全同；並重算兩組原版trace、保存、source及正式身份。
FINAL JSON／Markdown SHA-256：`074217f7fcce5c423e00a5cc46b3ebe331899ad17ba42ab64c6f158c8b913a20`／
`444072aaa6f59b4de54d6e1298cca118a2334c6685336df6ba5f2c5a97900d2b`。
完整回讀後75項來源／結果／影格／wire／binary身份核同才採認。未驗其他畫面、2倍與原生Windows／macOS，不外推按住重複。
