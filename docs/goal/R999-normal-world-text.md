# R999：正常世界遊玩與未量文字

2026-10-10，起點0b97000；沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#24](https://github.com/wicanr2/wizardry7-cht/issues/24)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。

1. 從R990由原版正常SAVE取得的合法進度冷LOAD，在新Scratch副本正常操作；不注入state、座標、角色數據或作弊。
2. 先核物品與互動選單，尋找STONE、FORCE、LEAVE等實際輸出，再沿可見路線探索；每次新畫面先讀再操作，不盲點對話或戰鬥。
3. 文字觀測限既有已授權29出處×64B目的緩衝及已上屏輸出，不讀鄰文、整份手冊或答案表。防拷如遇到，依現行029提示與原版正常輸入，不繞過判定。
4. 以正式現行5d6、Unifont3倍、既有資源新正常錄製，保存真正GUI與有限操作／進度收據；錄放已完成，不藉本輪重開#17。
5. 發現缺譯或語意分流時，先保留原版來源／呼叫／座標證據，再依RE→DRAFT→獨立READY→實作閘門處理。不替未回答的譯名問題選字。
6. 有界正常路線及新發現經獨立複核後提交、推送、回讀Issue。沒有觸發即只記本路線未量，不降低全部中文化完成條件。

入口：[目前狀態](../../CONTEXT.md)、[十一組證據卡](../re/019-catalog-migration.md)、[既有窄觀測授權](R857-targeted-maze-route.md)、[正常保存來源](R990-normal-female-profile.md)、[029](../spec/029-clearance-answer-hint.md)。
輸出只留ignored workplace/r999；既有原版／保存／字型唯讀，Docker沿既有離線、UID1000及有界資源。完整輸入、保存與原版畫格只留本機，R433不動。

## 長劍服務24窄觀察契約 READY

目的只補[研究021](../re/021-compound-item-output.md#longsword來源訂正)的新位置背景缺口。舊r717至r730工具已不在本機；保留文獻證據，不假稱重用了舊來源或已取得新服務參數。

- 以新完整short-item錄製作私人原API派送；不改header或footer，不執行逾時partial，不聲稱不同binary正式重播相容。沿已審R3 controller及有限隱私日誌。
- 只在現有有界NUL見證精確接受LONGSWORD、font0／x18／y114後，開啟至多1,000,000原版步的窗口。每次NUL最多接受16個服務呼叫，全程至多64個。窗外不取背景、不讀stack或code。
- 在已記錄的resident `0110:40a7`入口與`0110:40b0`返回觀察數值Regs、Step、SS:SP的至多16 bytes及各位址至多12 bytes code。讀取前核16-bit範圍及物理終址小於A0000；只讀數值stack，不追指標、不讀DB／MSG或鄰文。不得將舊driver segment、SP+2或任何矩形參數解釋寫成gate。
- 一次只有一個pending呼叫，綁觸發NUL身份與固定deadline，新NUL不得延長舊pending，須先丟棄舊件。入口複製當下Indexed及現行offscreen()，返回只在同SS、Step未倒退且仍在原窗口才接受；SP及暫存器變化先照實記錄，不推定呼叫語意。重入或未配對返回只記有限原因並丟棄。只輸出實際差異bbox的before／after／mask及palette，完整64K留在記憶體比較。
- 不新增Run、輸入、Stub、原版記憶體／離屏／A000寫入或產品覆繪。輸出O_EXCL、0600，任何I/O失敗輸出固定INVALID碼並中止觀測。沿現有來源副本與Docker隔離，不修改正式Go或SDK。
- 同一私人binary啟用／停用觀察跑完整wire，原版每次發布的六類摘要、起始保存與最終保存須全同。0命中或未知來源如實記錄；只有數值及frame差分，不稱服務24語意或中文修補完成。後續正式路由須另走RE→DRAFT→READY。

獨立審查回報阻擋0、應改0；唯一提醒為重疊窗口不得延長舊pending，已明寫第4條。只放行本節私人觀察，不放行產品路由。

### 工具入口

可重現來源保存在[產生器](../../tools/re/prepare_item_bitmap_probe.py)、[NUL見證](../../tools/re/item_nul_probe.go)、[服務觀察](../../tools/re/item_bitmap_probe.go)、[私人原API派送](../../tools/re/item_probe_api.go)。只在專案Docker內使用，正式來源及SDK唯讀，輸出指定到新的workplace子目錄：

```sh
python3 tools/re/prepare_item_bitmap_probe.py --source /src --output /out/private-source
python3 /out/private-source/tools/re/build_recording.py --output /out/build/wiz7-item-service --receipt-dir /out/build
```

建置前由同UID建立`/out/build`。原版執行用[有界Xvfb包裝](../../tools/re/run_probe_xvfb.sh)，容器內以`timeout 900s sh /src/tools/re/run_probe_xvfb.sh <診斷命令>`啟動；原R979包裝與本檔相同。`WIZ7_R982_API_DOCUMENT`只接受成功end的新完整本機wire，`WIZ7_R982_TRACE`指定新摘要檔。`WIZ7_R999_ITEM_PROBE_DIR`指定既存空目錄才啟用觀察，控制組不設此變數。原版及初始保存唯讀、scratch新複製，正常root／save-dir、字型及時鐘需與wire相同；原API先驗儲存及初始純量。本工具不支援正式跨binary重播，不接受partial，也不提供中文修補。

像素數字由[獨立核對器](../../tools/re/check_item_bitmap_pixels.py)重算，容器內執行`python3 tools/re/check_item_bitmap_pixels.py --witness-dir <私有witness目錄> --output <新JSON檔>`。比較墨點實際色號，再另計改值點數；不要求既已同色的像素發生改值。它只產生資料對齊收據，不作bitmap參數、背景還原或中文路由判斷。
