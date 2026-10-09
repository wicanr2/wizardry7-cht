# R1002：怪物名稱衝突與角色頁情境

2026-10-10，起點f3e3863；沿[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[008](../spec/008-english-key-tsv.md)、[031](../spec/031-monster-name-slots.md)與[研究019](../re/019-catalog-migration.md)。

上一輪R1001完成正常有聲存讀檔軟體玩家鏈及推送，屬進展。本輪回TSV前置與中文覆蓋，不重做已量音訊、錄放或手持物品矩陣。

1. 從現行來源與R983固定格證據核23組怪物異譯，先區分已定譯名與仍待裁決項。不以數量、載入順序或最後寫入選名。BLACK BIRD／BLACK BIRDS統一「黑鳥」已展示給使用者，回答前不改正式文案；既有CLAYMORE、BUG／BUGS等問題不重複詢問。
2. 不依賴名稱裁決的分支，沿正式c7da與合法保存副本正常檢查角色／創角頁及職業、能力文字，優先尋找PSIONIC／POISON的實際觸發路徑。正常鍵鼠、逐圖確認，不直接修改角色記錄／state／記憶體或繞過原版限制；若需要開創角，只到可見列表後正常取消，不確認保存新人物。
3. 新正常完整錄製可供後續唯讀原API觀察；需要來源證據時，沿既有29出處與64-byte授權另定有界契約、獨立READY才安裝私人探針。只量當前目標，不讀鄰文／詞庫／防拷表，不從畫面中文字推斷MSG來源。
4. 原始SCENARIO固定名字格只沿031已核區域，保留idx／slot及原始定位。031、008未READY前不正式重生怪物catalog或進行TSV遷移。可撤回候選只供證據與名稱展示，不作正式名稱或完成聲明。
5. 具體結果回填目前CONTEXT、研究019／013與WORKLOG，驗收後逐檔提交、推送及回讀Issue。十一組、控制碼、其他同鍵及平台／交付仍在整體目標內，不用局部通過代替全部完成。

Docker離線、UID1000、限資源、硬截止與清理沿專案規則。原版及初始保存唯讀，只寫`workplace/r1002`及已授權cache。完整遊戲輸入、畫格及保存只留本機；原版題目按當前提示正常輸入，不改驗證。

## 正常路線第一趟

正式c7da已正常開創角，僅一次產生8點獎勵，選目前可用法師，分配剩餘點數及查看學術頁；未重擲或強制選職業。錄製器閒置上限先收尾，沒有完成預定取消，不能稱取消流程通過。完整wire `profile-context.wiz7-input` SHA `e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732`，3606事件／1786發布／末438.4M，成功end、續玩及關窗0；兩保存bytes／mode／mtime與初始相同。這份完整輸入只供本輪診斷，不用partial或改header。

創角早期性別欄實際顯示英文M，接受頭像後同位置才顯示男。此為新的可見覆蓋缺口，出處未確認；另沿#22做窄字元印字觀察，不從畫面猜caller。名字及規則不改。

## 私人來源觀察 READY

沿既有29出處×64-byte授權與已證`0110:0582→0659→3895`鏈，原址與配對欄位見本機`observer-contract-facts.json`，SHA `c75123dd6885ea6832a73c11647ff92167f55259365a677582c949ef967add10`。舊R984 helper缺目前的守門，不能直接執行。

- P1：所有callback先核精確CS:IP；每次Bytes讀前核offset不wrap、物理終址不超A0000h及合法長度。0582先只讀ID word，排除保護／非白名單才准讀ret／dst。每次新reader清舊pending與decoded代，白名單才建立一代；不記錄非目標ID或指標。只允單pending／decoded，不留無界registry。
- P2：0659須同DS／SS、SP＝entrySP−112h、BP＝entrySP−2、BP+4的ID與BP+6的dst相同，最多1M步。核完整64-byte目的span合法後逐byte讀至第一NUL，最多64bytes；只接受該ID對應既授權的完整英文鍵。錯誤只記有限原因，不輸出陌生內容。3895須同代、同DS:dst、1M步內及完整英文相同才記print；引數ret／font及當前driver x／y只讀既有固定欄位，driver由0110:2de7與0163四byte簽章驗證，不硬編碼舊段。
- P3：另限本wire的332800000至417800000步，觀察3877的M／F、font0及游標90,6。此窗口由正常性別點擊至接受頭像的完整事件before／after_step取得，不改輸入或原版時序。只讀六byte入口stack；合格後沿已證driver離屏指標取現有畫格自有副本，返回須同SS、SP+2、精確near ret及1M步內，核游標96,6。保留6×6 before／after、差分bbox／數量及域外差異數，不匯出完整離屏或RAM。只量來源，不因M字形就註冊新正式路由。
- P4：所有3877重入先清舊字元pending，不保留可延長的舊窗；動態return hook最多8個、字元接受最多16次、全部輸出最多256事件。只有固定格式O_EXCL0600 JSON／窄像素，I/O失敗有限INVALID並停採，關閉輸出。原API及原版返回保持，不新增Run／Click／Stub／記憶體寫入，GUI不讀Oracle。
- P5：從正式96來源建立私人副本，沿已審原API派送完整wire，on／off同binary、起點、原版及固定輸入。全部1786發布的六欄摘要及保存metadata須同。私人API不宣稱正式跨binary重播相容；0命中與窗口外分支照實分列。獨立契約到READY及來源複核後才執行，正式008／031與任何新性別路由維持原閘門。

獨立契約0／0／0完整讀回採認，JSON SHA `50a735cd4b1938005aa0e871f49fa267a0edd577180613076ce8a394f99a6a3a`。只放行本私人觀察實作，來源複核前不跑原版；不註冊任何新正式路由。

## 私人觀察結果

98來源SHA `232ad6eed1ff5a50b060d50b9743ae1bcb2af34c6f136d03ab5006775f3780ea`，binary `22c6b9b93e87164886c6df1de16b1d26c502655165bc2da08bf3f5794c85ca61`。來源複核0／0／0完整採認，JSON SHA `0d112c0aed83c4ed442ced140f18287f686ad1436c433842f6c4ea1ae0c79a8b`；才執行指定完整wire的開／關兩組。

兩組退出0、無INVALID，1786發布的Step／Ticks／RAM／Regs／indexed／palette逐筆及完整trace bytes相同，SHA `bc95b65d4f308b136f4236993ee97db7f57ebb25d43df622d0697a2a421fc1af`。兩保存size／SHA／mode／mtime保持。29出處的reader／decoded／print皆0，只代表本次未量到；十一組不解除。

字元觀察兩對：334908491→334908883返回`0110:9595`，402932251→402932643返回既有`0110:6c42`。皆M／font0、90,6→96,6、同SS／SP+2、同driver／buffer；各13點變為色1，bbox為90,6至95,11，原6×6域外0差。9595是創角早期新的已量來源，不推定F、完整資料流或正式覆繪已完成。細節見[研究025](../re/025-ui-localization-coverage.md#r1002-creation-sex-source)。獨立動態驗收另記。

獨立動態110項0失敗，完整讀回採認。JSON SHA `96ea1e9ebb8e6ba79a663e08c38e2478ac800b428f70f94fc4c61b539d39016b`、Markdown SHA `7211eb9e6f8bb8d549e554eaad5bfd9e86e9ce92514f0ea1743f0d6ca78e92d5`。四份保存工具逐byte同已審已跑版本，從保存constructor純建置的98來源manifest與binary全同；不新增原版執行。只採本私人對照與來源診斷，不稱新M路由完成。

## 本輪研究工具

只對本輪固定原始輸入及私人診斷，不是玩家功能入口：

- [觀察器](../../tools/re/character_context_probe.go)只由設定`WIZ7_R1002_CONTEXT`啟用，未設定不安裝hook。
- [產生器](../../tools/re/prepare_character_context_probe.py)純建置98來源；須保留本機R1002的`observer-contract-facts.json`與`probe-contract-review.json`及上述SHA，先把觀察器逐bytes複製至全新輸出目錄的`context_probe.go`，再傳`--out`。
- [執行器](../../tools/re/run_character_context.py)限定本機完整`profile-context.wiz7-input`，原版／初始保存／`/src`唯讀、每組空白`/out`可寫；環境變數指定`WIZ7_R1002_VARIANT=control`或`observe`、上述binary路徑與SHA。Xvfb有界入口沿[run_probe_xvfb.sh](../../tools/re/run_probe_xvfb.sh)，880秒內收尾。
- [比較器](../../tools/re/check_character_context.py)讀本輪`probe-api-control`／`probe-api-observe`，核固定身份、1786七欄含publication、保存、事件上限及字元配對，另存comparison；既有結果不覆寫。

上述工具只在既有離線GUI容器內執行，遵守專案UID／資源／逾時與掛載規則。完整輸入、原版資料與窄像素收據只留本機；缺本機前置檔就停止，不用代用品。

## P6 新9595來源補證 DRAFT

R1已量新caller，只補正式路由缺的來源及背景，不新增遊戲輸入。P1至P5與R1收據保持，以下須另經獨立READY及來源審查才執行新私人R2。

1. 沿P3已接受的M、font0、90,6，額外只接受返回9595及步數334900000至335000000，最多一對。先用既有安全讀取核`0110:9592`三bytes正是`E8 E2 A2`，也就是near call 3877、返回9595；不符只記有限原因，不讀後續區域。
2. 簽章合格才讀當刻`0110:9515`起128bytes至9595前，先核offset／物理完整界線。只保存這個runtime原始code片段、SHA與CS:IP基址於本機；不讀函式以外的文字表、原版MSG或鄰文。用途只為確認已量M來源，不從未知前段猜完整函式。
3. 沿原P3已取得的自有64000-byte before／after，不新增畫格讀取。額外輸出88,4起8×8各64bytes及八背景點，只有該一對新9595樣本可輸出；其餘P3仍只6×6。返回仍核同SS／SP+2／near ret／1M與96,6、同driver／buffer才接受P6，失配保持R1原事件並有限記原因，不推正式路由。
4. P4的事件／hook／樣本上限、固定JSON／0600／INVALID及有界Close保持。R2從未改正式96來源重建，同R2 binary的開／關兩組重播同完整wire，1786六欄及保存相同；另核trace等於R1基準。R1及兩組結果保留不覆寫，不藉新抽樣重擲、直接進入或延長原版執行。
5. code以runtime `0110:9515`基址記錄，交給既有IDA9.4隔離工具鏈只解這128bytes，原址／operand不改名，證據等級分清。足以定義窄M路由後就停止，不研究整段創角。9595的F未量不放行；正式024增補仍另走DRAFT→獨立READY→實作→正常GUI及原版對照。
