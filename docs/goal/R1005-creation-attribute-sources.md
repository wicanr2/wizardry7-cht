# R1005：早期創角屬性來源補證

2026-10-10，起點008ce17，沿[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)與[R1004原型](R1004-creation-attribute-layout.md)。R1004新增原型及實際重疊量測，屬實際進展。排版及移位箭頭的輸入對應仍待使用者定案，本輪只補不依賴該選擇的來源證據，不接正式重排或輸入映射。

## 既存證據與停止線

- R1003 frame1585的Shown與原版indexed已證八屬性在DOS115,35至77，每6一列；native字高15／16，fg85。R1004量出兩倍相鄰墨點交集。原版STR已選箭頭bbox只量到DOS133,35至153,39，不是全列命中區。
- 024／研究025的R954表記完整人物頁標籤返回6a43、數值返回032a。當時r935／r937／r954／r966原始工作區現已不存在，不能把文件表直接升為早期創角的同路證據。
- 現行main.go在精確0110:3895入口已取得r0、ret、partRaw最多8bytes、compoundInput的字型／字串來源指標／游標／Step。可把這些既有值交給不持有Oracle的私人helper，無需再讀stack、driver、文字、離屏、VGA或原版程式。
- 本輪只定位標籤／數字輸出tuple與時間順序。不證數值before／after、底圖還原、箭頭來源、實際命中區、撤除或正式版面。取得來源tuple後即停止此入口切片；需要後續像素證據時再寫窄契約。

## 私人參數觀察契約 DRAFT

C1：只在私人main副本的既有compoundInput完整建立之後、任何分支早退之前，追加helper呼叫，傳入該in及原已取得的r0.SS。helper不接收Oracle、memory reader或callback，不註冊任何OnCall／hook，不呼叫讀取／執行／輸入API。所有正式語句順序及值保持。只有現成參數的副本可輸出，不能從地址再解引用。

C2：先核Step在332800000至417800000閉區間、in.Ret.Seg為0110、font0、原字串span不越16位offset或物理A0000，再核有限輸出tuple。標籤只收x115、y35＋6×row、ASCII STR／INT／PIE／VIT／DEX／SPD／PER／KAR與對應row。數字只收ret032a、x139或145、y35＋6×row，以及x139或145、y89的bonus候選；必有現成raw內NUL，NUL前僅1或2位十進位ASCII。BONUS標籤只收x103,y89與完整ASCII BONUS。未知返回的八標籤或BONUS可記ret作來源探索，但不因此放行正式路由；未知座標、字型、文字一律不輸出。不得輸出一般文字、原raw尾段、MSG ID、詞庫或防拷內容。

C3：只有WIZ7_R1005_PARAMETERS指向新檔時啟用，O_EXCL0600建立；off時nil helper完全不寫檔。記setup、有限entry及summary，最多512筆含末註；entry保存Step／SS／SP／DS／Off／ret／font／x／y／白名單key及分類，不保存原版位元組。到上限、開檔／寫入／close失敗時固定R1005_PARAMETERS_INVALID類別、停採並保留失敗輸出，原版照常。最末summary記各類計數、總數與是否cap，不把不完整檔當成功收據。helper只由原owner執行，無GUI共享或背景goroutine。

C4：私人來源由現行97來源正式樹加既有已審item_probe_api.go及本helper建立，純建置後獨立來源審查才執行。只沿R1002完整合法wire，SHA e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732，正常API派送，不改header／輸入／seed或追加Run／Click。相同私人binary on／off各一次，Unifont2、原版與R990合法初始保存唯讀、全新scratch分開；核每組1786發布的六類摘要與R1003既存trace bc95b65d4f308b136f4236993ee97db7f57ebb25d43df622d0697a2a421fc1af逐bytes相同，以及實際SAVEGAME／PCFILE bytes／mode／mtime保持。完整wire只留本機，不冒稱跨binary正式-replay。只改私人副本，不新增畫面矩陣或重跑R1003的M像素驗收。

C5：獨立核觀察器輸出確實僅白名單、entry與summary計數相合、無INVALID／上限、tuple時序可回查Step、兩組原版摘要及保存保持。零命中也如實接受，不改window／加寬字串／重跑追求命中。結果若只證入口，明列返回與像素未知；不把新ret命名為confirmed頁面owner或規則。只提交工具與安全摘要，完整參數收據及原版資料只留workplace/r1005。

本契約先交唯讀子代理逐條審查；READY前不寫觀察器，不執行新原版。受影響正式版面仍依R1004提問，待答期間僅此讀取證據分支可前進。
