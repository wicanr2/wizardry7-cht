# R1000：選取物品圖片來源與清除

2026-10-10，起點7fff275；沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。

R999已量四NUL及兩次40a7／40b0當次差分，不重跑相同觀察。當次before含既有英文與游標；正式長劍選取覆繪尚未完成。

1. 先用既有IDA Pro 9.4及資料庫定位40a7、caller5cb9／5b80的原始指令與資料流。保留位址、bytes及推論等級，不從stack值猜矩形或來源。
2. 若原檔／資料庫缺少runtime內容，或靜態資料流的動態指標與背景仍需實測，另定有界唯讀觀察契約，獨立READY後只取得必要code、游標像素及背景；不讀DB／MSG、鄰文或答案表，不泛查整個driver。
3. 證據須說明英文來源如何進入圖片、顯示位置、正常移動及放回如何清除，才寫正式DRAFT、審READY及實作。只靠墨點形狀或舊segment不掛正式中文。
4. 正式覆繪只改合成畫布，正常原版輸入、記憶體、離屏、A000及存檔保持；驗收包含選取／移動／放回的實際GUI、同狀態原版及兩字型兩倍率。未觸發分支不外推。

輸出只留`workplace/r1000`，R999原收據唯讀。IDA技能及專案SDK規則沿既有入口；Docker使用現有image、UID1000、離線、有界資源及明確清理。完整遊戲輸入、畫格與保存不提交。入口：[021](../re/021-compound-item-output.md#r999-服務24的當次畫格)、[R999工具](R999-normal-world-text.md#工具入口)、[目前狀態](../../CONTEXT.md)。

## 自訂游標來源與背景私人觀察 READY

靜態已找到VPCVW兩caller、service24入口14c9及游標繪製1f1a。仍須量到當前像素指標與背景，不把靜態語意直接當正式覆繪oracle。

本節所有新增讀取，包括P2／P3的far pointer、code、descriptor、busy，均在讀取前核段內不wrap及物理終址不超A0000h；失敗只記固定原因，不讀取。獨立審查無技術阻擋，唯一文字應改為明寫此共用守門，已補齊後採READY。

- P1：從R999已審98來源副本起步，正式Go／SDK唯讀。用同合法complete short-item wire私人原API派送，不改header／footer、不跑partial。只在精確LONGSWORD NUL已見及原service見證窗內，重用已讀16-byte入口stack。只接受caller5cb9、尺寸76×13、pointer等於DS:0e1e／0e20。新讀只限兩個全域word、該988-byte像素span；先核段內不wrap及物理終址不超A0000，不追其他pointer或文字。
- P2：driver segment從resident `0110:2de7` far pointer讀取，offset須0163，該處四code bytes須與原VGA.DRV的`e8 63 13 cb`同；不硬編碼舊740c段。定義時只另讀resident5c24起32 bytes、5c9f起32 bytes，reset時5b60起32 bytes，與本版VPCVW相應原始碼逐byte核對；不讀其他caller或鄰文字。入口只留候選，原有40b0見證接受後才啟用註冊；新NUL或配對失敗清候選。圖片持有自有副本，與當次NUL身份、Step綁定；最多8次定義。shape只作事後資料對齊，不作來源註冊條件。
- P3：只在已註冊圖片的30M步窗內、最多128次既有發布取樣。原owner既有snapshot觀察，不新增Run、input、compose或GUI讀Oracle。每次讀driver `15ac..15c3`的24 bytes、`1e5c..1e67`的12 bytes及`1ed3` busy byte；busy非0只記有限status，不讀像素背景。pointer或尺寸不符即停止該註冊，原版保持。
- P4：busy0且pointer／76×13相符時，解析當前x／y、clip width／height及背景pointer `15b4`。必須座標在320×200、clip符合右／下裁切、面積1..988、背景span段內且物理終址不超A0000，才讀恰clip面積的driver背景。輸出獨立副本的descriptor、背景及Indexed相應ROI；圖片取樣仍為已註冊988bytes，不讀鄰區。不符只記固定原因，不猜補或放寬。
- P5：從已有service入口stack量到caller5b80的零pointer／零尺寸時，記reset_request並停止該圖片像素取樣；下一次既有發布只讀同P3的descriptor，pointer零且非busy才記reset_complete，之後停止，不把預設或其他圖片當長劍。離窗、128次、定義超8或I/O錯誤皆有界收尾；沿O_EXCL0600及固定stderr INVALID，無原版寫入、Stub或額外字串讀取。
- P6：同一私人binary啟用／停用，1,514發布的原版六類摘要及兩存檔metadata須同。獨立程式用背景＋圖片的非255像素重建當前ROI，與真正Indexed比較；0命中、busy、失配與未量reset如實分列。只採來源／背景證據，不放行正式中文路由或全物品生命週期。

### 正常移動收據增補 READY

第一組29份固定位置ROI已全部重建0差並獨立採認。再以新正常錄製`movement/movement-item.wiz7-input`核移動及右下裁切，SHA-256 `b413c31a633943eb97b9d19aaeac4fd74f06ec9927bcb1761733b9f88ab9dd4c`，3333事件／1651發布／末370.8M，正常Stop、續玩及關窗0；這組P6發布數改為1651。其餘P1至P5、同99來源私人binary、30M／128上限及讀取範圍完全保持。新wire後半的皮胸甲只作GUI來源範例，此探針只觀察LONGSWORD；窗口外或未量分支仍分列，不延長取樣或泛讀其他物品。增補只擴驗收輸入，不放行正式覆繪。

## 工具入口

私人來源由[游標產生器](../../tools/re/prepare_cursor_background_probe.py)及[凍結helper](../../tools/re/cursor_background_probe.go)重建；產生器沿用[R999來源](R999-normal-world-text.md#工具入口)。只在專案Docker中執行：

```sh
python3 tools/re/prepare_cursor_background_probe.py --source /src --output /out/private-source
python3 /out/private-source/tools/re/build_recording.py --output /out/build/wiz7-cursor-probe --receipt-dir /out/build
```

由同UID先建立`/out/build`。GUI及原API派送依R999有界Xvfb入口，三個`WIZ7_R982_API_DOCUMENT`／`WIZ7_R982_TRACE`／`WIZ7_R999_ITEM_PROBE_DIR`環境變數保持；最後一個只在取樣組設定，控制組不設。只接受本輪成功end的新完整本機wire，原版及起始保存唯讀、scratch新複製，不提供正式跨binary重播相容。保存工具的乾淨99來源及binary已與R2逐bytes核同。

## 單行來源失敗的私人內部觀察 READY

正式R2兩倍率Unifont實跑的原版1651發布及保存保持；皮胸甲中文像素通過，長劍穩定英文回退，尚未採認正式修補。新私人副本只在既有名稱入口、返回、capture及讀到空NUL時，輸出已讀的候選狀態、generation、步數、CS:IP／SS:SP／DS、部件文字與矩形。沒有新增Bytes、原版Run、輸入、Stub或記憶體寫入，不讀鄰文／DB。最多256個排他0600 JSON；沿已審觀察器與相同正常movement wire，原版及保存只讀。只定位候選被清除的階段，不能放行新正式規則。來源及契約獨立審查通過後才執行。

獨立`single-line-static-review.json`／`.md`完整讀回0／0／0採認，允許此一次Uni3私人內部診斷。靜態只證正常第一返回後必經第二NUL，第二內容仍待本次觀察；不以靜態控制流直接補正式判定。

## 私人Draw截圖綁定修正

R3四組中23圖通過，倚天2倍首幀的實際Draw與同狀態HD不符。私人Draw先複製畫面，繪完後才另鎖取截圖請求，存在舊畫面取走新請求的競爭窗口。該窗口由程式證實，本次失配由此造成屬強推論，沒有宣稱量到失配Draw的實際發布編號。

R3b只把讀取及清除drawReq移入Draw複製rgba／held／publication的同一鎖。正式96來源、原版、輸入、覆繪規則與已有界觀察範圍不改。經獨立來源複核後，以新私人binary與新輸出目錄重跑四組，首幀仍須存在且逐像素相符，不刪樣本或放寬期望。

## 正式覆繪驗收重建入口

[私人產生器](../../tools/re/prepare_held_item_verify.py)、[唯讀觀察器](../../tools/re/held_item_verify.go)及[獨立原字型比對器](../../tools/re/check_held_item_pixels.py)只供Docker驗收。產生器以正式R3來源為輸入，不執行原版；沿既有`build_recording.py`產出98來源私人執行檔。原API派送沿R999工具，`WIZ7_FORMAL_OBSERVE`指向新的私人可寫輸出目錄，固定8來源／16畫格與reset上限。只用成功完整wire及合法保存副本，不改header或稱跨binary正式相容。

```sh
python3 /src/tools/re/prepare_held_item_verify.py --source /src --output /out/private-source
mkdir /out/private-build
python3 /out/private-source/tools/re/build_recording.py --output /out/private-build/wiz7-held-verify --receipt-dir /out/private-build
python3 /src/tools/re/check_held_item_pixels.py --root /receipts --eten /eten/stdfont.15 --unifont /unifont/unifont_all-17.0.05.hex.gz --prefix formal-verify-r3b --source-review /receipts/formal-source-review-r3.json --source-review-sha 12260b7447e88ff1b18e2fb64761f88807b5cef3e20c0fa00148cbde6f35c79c --out /out/pixels.json
```

掛載來源與字型均唯讀，僅`/out`及已授權cache可寫；GUI沿既有有界Xvfb入口。舊背景探針的99來源身份以7fff275的凍結來源為準，不把新正式來源重建結果稱為舊binary。完整本機輸入、原畫格、字型及保存不提交。

## 結果

正式R3已接入，037本次兩物品路線限定CONFORMED。空第二部件須原配對與完整988像素保持才保留單行名稱；同位址新圖及未知來源仍退回原圖。雙字型雙倍率24圖、8677項、原圖19920點及RGB129480點全部0差，四組每1651原版發布與無覆繪控制相同，保存保持。真正GUI選取、移動、右下裁切及原格放回通過，既有第三列8892點移開後0差。新3408事件／1686發布／末377.8M正式錄製正常Stop、續玩及關窗0，同版嚴格重播0。

正式建置為`workplace/r1000/promoted-build/wiz7-formal`，與已驗96來源c7da版本逐位元相同；保存工具乾淨重建亦與98來源私人28b3版本逐位元相同。限定數字、身份及其他未量條件見037／研究021；不重開已完成來源或錄放，不解除十一組未量情境及全案交付門檻。
