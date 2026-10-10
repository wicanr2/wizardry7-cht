# R1003：創角早期M正式覆繪

2026-10-10，起點75cbcd6，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)與[024 READY增補](../spec/024-character-review-output.md#r1002-創角早期m覆繪增補-ready)。上輪R1002是實質進展：新9595來源、128byte程式、8×8背景及兩組原版保持已足，正式八條審查0／0／0，直接實作，不重開來源考古。

1. 在私人候選接入已量M來源，沿已定字模與2／3倍率位置。精確CS:IP、讀前SP／完整span、code／driver簽章、背景、SS／步數／driver／buffer返回配對及pending重入清理照READY八條；原版與既有6c42／6159路由保持。
2. 只改必要Go來源及有意義負測。候選完整cmd/wiz7回歸及純建置完成後獨立來源審查；來源未採認前不執行候選原版。正式接入前後來源與binary身份明列，英文鍵、字型資源、保存及錄放schema不改。
3. 用R1002合法完整wire作私人原API修前／修後／無覆繪控制，兩字型兩倍率核1786原版摘要、保存及早期M批准域。固定輸入不改header／footer，不稱跨binary正式重播相容。
4. 獨立由R1002原字元差分、8×8背景及原始字型重建期望，核HD／真正Draw及域外保持；沿原版重畫或既有性別來源接手核撤除。另跑新正常GUI確認早期M中文、完整字格、正常收尾及保存保持，新正式錄製用新binary自身身份。
5. 根代理完整採認獨立結果後才限定CONFORMED、正式接入／乾淨重建、逐檔提交、推送及Issue全文回讀。未量F／其他caller／背景、取消／新人保存及全案交付仍保留限制，不用本組代替全文完成。

現行正式基準為R1000 c7da binary、96來源；資源沿R995雙字型、R983 flow及R979 help。原版、SDK與舊wire唯讀，研究只寫`workplace/r1003`，Docker離線UID1000、相稱限資源與有界Xvfb；完整遊戲輸入／原畫面只留本機。未定怪物／物品名字及人耳問題維持待答，不重新問已定BASELARD。

## 私人畫面驗收契約 READY

V1：沿已審`item_probe_api.go`派送R1002完整wire，SHA `e5137dd4453f2394a1d476dcdccb1734484d1fc8ed93ebe6da758b9fb4af8732`，初始R990合法保存及原版唯讀。私人來源只加原API／六欄trace與owner快照擷取，正式版本不加入觀察器。來源審查後才執行；不改輸入、種子、header／footer，不追加Run／Click／Stub或原版寫入。

V2：修前／修後各倚天及Unifont、2／3倍共八組，無覆繪控制按2／3倍共兩組。控制組設xlate／flow／help空值，須實際確認Shown空、comp等於A000及HD等於原indexed／palette放大才算無覆繪。每組核1786六欄trace及保存metadata，並等於R1002基準bc95b65d…c1af。純控制沒有字型差別，兩同倍率字型共用一組控制，不重複相同控制。

V3：只抽完整wire中help關閉的五個既存publication：1519／Step332800000前態、1520／338800000早期M、1585／363400000屬性頁、1686／406800000既有來源接手、1712／417800000後態。owner持鎖的snapshot輸出既有HD驗收資料A000／comp／palette／Shown／六欄摘要，另記同發布的sex事件originRet／原座標／守門域／before／after／字模及shown狀態；不讀新原版code、文字表或任意記憶體。當頁正常可見輸出只留本機，沒有防拷畫面。

V4：真正Draw請求沿R1000已審方式，於Draw第一個鎖同時取得rgba、事件／publication及drawReq並清請求，避免畫格錯配；只改私人副本，不重開正式錄放或除錯鍵語意。每publication最多一次、全組最多五樣本。每樣本先建唯一子目錄，已存在即有限INVALID並停採，不覆寫；新增metadata以O_EXCL0600保存。沿既有HD輸出helper只寫該新目錄。既有helper會吞I/O錯誤，擷取後必驗全部必要檔存在、完整可解析、尺寸／長度正確及Step／publication配對相符；所有I/O或配對失敗均固定INVALID停採、保留該失敗組，原版照常，不用panic中斷原版。

V5：獨立核全部原版摘要與保存、五圖HD／Draw發布配對，原版M的13點只取R1002已驗before／after差分，字模由既有兩原字型獨立解析，不拿候選自身G作期望。早期兩樣本批准域為DOS88,4起8×8放大；其他三樣本完整修前／修後畫面應相同。新M不得有英文墨點／裁字／域外差異；接手後無新9595事件，原6c42頁依既有契約。畫面或收據不齊照實FAIL，保留全部結果，不篩掉首幀或把私人API當新正常GUI。

獨立審查阻擋0／文字應改1／建議0完整採認。V4已補全部必要檔後驗證及所有I/O／配對失敗的固定INVALID／停採／保留失敗組要求，沒有技術或產品分支變更，依AGENTS3採READY。JSON SHA `d83e676503ed23aa0464ff8951091334cdc87448f96043f5f380c0fefbe2c59c`，Markdown SHA `b143cd09f54d7b7611aa28078416bfc7fc83de1b59520aedc755c490a4d9bd00`。私人實作仍須來源複核後才執行原版。

## 候選與驗收進度

候選新增creation來源及負測，完整cmd/wiz7為349 pass／22個明列舊opt-in skip，無FAIL；97來源SHA `18c9d3bbcbb5f1e4164f18597f63e65edd6c02fd6253c74b77024bf01299579d`，binary SHA `631926caa447010014ff70dada2dc56179a010eaa6fd5573394084b511484c99`。正式來源保持。候選獨立0／0／0採認，JSON SHA `3e12f07b7876cf5d286195a3ca57ed416b0d20139842092d2013662f532c4d64`後才跑候選正常GUI。

私人擷取R1先擋下publish之前的舊Panel及逾時FAIL保存，R2修正後再補trace IO只讀一次，R3獨立0／0／0採認，JSON SHA `8340cddf6f37da53b9f7bc772cf9a7b086f58116cb092bfe041956e4c1d5972b`。已審修前98來源binary6ff7da29…f6f9、修後99來源binaryd29168b1…3c04才執行；前兩版收據／binary保留未run，不用R1冒稱目前觀測。

正常GUI控制器獨立0／0／0採認，JSON SHA `1e665e7f83a36137e943933b56db83e4201f0b89e1c5a5b69132593c60eb0685`。第一Uni3從正常新建／Human／Male、實際可用戰士及9點原版獎勵，正常分配剩餘6點、到頭像頁接手、Stop續玩與WMclose0；兩保存metadata保持，完整wire SHA `b7289dea437046b3222202a148d9cd2bd683bc408e997ed853ea5837c6bbbbde`。主代理實看早期男及接手圖，尚待獨立像素／錄放驗收。

第一Uni2與倚天2正常M均上屏，根代理寫比對器期間達120秒閒置而正常Stop、續玩與close0，兩保存保持；只作M上屏收據，未算頭像接手。兩完整wire e531148f…d21f／7645c5f6…8d6a保留。另開新輸出補尚缺接手，不延伸或覆寫已終止會話。其餘正常與十組控制仍在驗收中，正式Go未接入。

Uni2第一次補接手時，把剩餘7點全點力量，第六次已為18／bonus1，第七次無數值變化且原API回no_response，錄製失效後partial停止更新。F1等待逾時及最後WMclose發生其後，不能判成F1產品故障；exit2及partial保留，保存保持。窄audit JSON SHA `1c10bab3a5f4355976c724505fbe7d2e2db794fa988280401d1dae51bd987f8f`。200ms F1原型未審／未run／未採用；沒有改錄放契約。

後續沿原R1控制器正常逐步操作：力量到18即停，點智力欄120,43可切換箭頭；Down／Tab未觀察到切欄，不據此改規則。最終Uni2新case及倚天2皆到頭像頁、Stop續玩／close0／完整錄製且保存保持。四組正常wire為Uni3 b7289dea…bbde、倚天3 e11d114d…cf94、Uni2 e2d35a4a…a4e5、倚天2 83e25033…7abf。主代理已實看四組早期男與接手圖；Reader、正式重播及獨立GUI像素採認另記。

三倍五組控制已全完成0，每1786原版摘要／保存同R1002；純像素checker275項、28,802,304 RGB像素符合獨立期望，原M13點擦除、原生男字及完整域外／HD／Draw皆核。兩倍另五組續驗，checker兩項frame綁trace／oldHD批准域基準的來源複核待採認，不提前CONFORMED。

實圖另見兩倍創角一般屬性標籤仍重疊，連頭像選擇頁亦未套成R966已驗完整人物頁布局。這是另一路中文排版缺口，不重新開啟已完成R966世界人物頁，另沿#22／#28處理；本次只驗新M來源與既有來源接手的保持。

## 保存的本輪工具

工具只供本機固定收據驗收，輸入與原版素材維持唯讀：

- [owner擷取器](../../tools/re/creation_capture.go)與[私人副本產生器](../../tools/re/prepare_creation_sex_capture.py)：以`--source`、已核`--main-sha`、不存在的`--target`及`--repo`建立；兩Go副本同名observer保留，只有成功publish後抽五點及同鎖Draw請求。純建置沿副本內`tools/re/build_recording.py`，輸出／receipt目錄先建立並核UID1000。
- [固定輸入執行器](../../tools/re/run_creation_sex_capture.py)：`WIZ7_R1003_CASE`指定上述十組，`/src`原來源與wire、`/original`、`/initial`皆唯讀，空白`/out`只寫該組；每組880秒，外層沿[run_probe_xvfb.sh](../../tools/re/run_probe_xvfb.sh)及910秒強制截止。
- [完整像素比對器](../../tools/re/check_creation_sex_pixels.py)：`--root`指十組本機收據，`--out`為新JSON；`--eten`／`--unifont`指固定原字型，`--font`／`--scale`只能作子集檢查，完整驗收不傳過濾條件。
- [正常GUI窄量測器](../../tools/re/check_creation_sex_normal_gui.py)：讀本輪四成功組、原字型及`--source-events`固定R1002 P6事件。只驗穩定8×8文字域，實際核背景／palette，不把普通PNG當同publication或頭像整頁oracle。量測待獨立採認。
- [完整錄製Reader](../../tools/re/creation_normal_wire_reader.go)：`-wire`、固定`-sha256`、不存在的`-output`，沿正式recording.Reader只輸出安全摘要，不執行原版或輸出header／輸入內容。

完整新舊wire、PNG、保存、字模及fail片段只留`workplace/r1003`；缺任一本機前置檔就停止，不造代用品。正常GUI及正式重播包裝的已審實際腳本與SHA留本機本輪審查報告，不把未採用的F1原型當入口。

十組全部退出0，每1786原版摘要／保存同基準。完整像素checker已審R2（JSON SHA `3e7d9c4c638bf43101b11415884dcbd8e70e86b67e3b6ed911b3173e584d6812`），最終550項／41,603,328 RGB像素符合獨立期望，結果SHA `943ae90bed0a637f4d47673d6858f2d60ed4c6c3533aee14e1d8074cfef4ca45`。正常四成功組的12張早期圖4992 RGB像素量測0差，測量JSON SHA `c7ef6027006b58aca776b53b777f1afaf3a625abdb7774464a7cda726dc5ce99`，不宣稱普通PNG同publication或頭像整頁oracle。

四份正常完整wire均通過正式Reader。新Uni3同6319正式`-replay`在包裝獨立0／0／0後執行，3627事件／1800發布／末435.4M、成功字串及退出0、保存與binary／wire保持；結果SHA `50ef7b68c5ab93635716b102d3dc87f13c8fe5e709ed1f4f20dd774112a11eb7`，無私人API。保存工具純重建99來源manifest及binary同已驗d291，未追加原版。最終獨立矩陣／普通GUI／重播採認尚待，正式Go仍未接入。

十組矩陣獨立0／0／0已完整採認，JSON SHA `a4c96a0fbb9f9f08232fb247a134ffa34bb7d8b7452777a9737841d1fde0769e`、Markdown SHA `e6412b83ae777842a44a3fc36010727caf2f5268ab564165576bc065ddab95fd`。獨立核17,860原版rows、20份實際保存、50 HD／Draw、443個來源manifest及八張原字型全圖重建，所有差異0；50個Regs亦獨立綁同pub trace。審查額外build目錄glob錯誤如實列工具問題，不重跑原版。普通GUI及正式重播的獨立gate續核，尚未宣告整輪或全案完成。

最後普通GUI／正式重播獨立0／0／0完整採認，JSON SHA `c631e52c22fb121cc17d9f1410a3af3ddb510606450a7c54992640447d0e33fb`、Markdown SHA `92d1892b703ac20558b1102f36160036899cf86d69a3875474cd756d8ffcfb7b`。另獨立核151斷言、12PNG／4992 ROI、四Reader／8個正常保存、實際strict成功log／保存及三binary身份。才逐bytes接入六正式Go檔，重跑cmd/wiz7無FAIL、349測試pass／22舊opt-in skip，套件pass；乾淨97來源與binary6319逐位元同已驗候選，promotion-result.json PASS。024本M路線限定CONFORMED，未量條件與一般二倍排版另保留，未宣稱全案或發行完成。
