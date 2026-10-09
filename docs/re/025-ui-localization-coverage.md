# 按鈕與訊息的中文覆蓋驗收

正常暖戰鬥的按鈕操作與撤除：[R962](#r962-normal-battle-gui)。

<a id="r972-sampling-takeover"></a>
## R972 譯文上屏抽樣接手複核

2026-10-07，沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)及[R972](../goal/R972-sampling.md)，只分析既有收據，不重跑原版。正式程式、譯文、字型保持。
當時輸入是R972本機中英文PNG、JSON、A000及色盤，catalog只作正式文字對照。
完整原版衍生資料、匿名材料、讀值及暫存工具已於驗收後依使用者要求刪除；
以下數字與雜湊保留，不能直接重跑已清理的收據。工具為既有`wizardry7-go-ebiten:latest`，
image ID `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`，Go1.24.13、Python3.11.2。
本輪不新增原執行器binary或EXE來源閉集合驗證。
證據等級為confirmed，限保存收據的摘要、像素與可見字形核對；遮字仍為未知。

原`run_pass2.py`使用小寫欄名及`.get()`，使缺欄變成`None == None`；新複核必須取得大寫欄名，
缺欄立即報錯。舊`settle`使用截短雜湊，這次另逐位元比對完整資料。16對原抽樣及12對等待／診斷取幀，
共28對`Step`、1MiB RAM的FNV32摘要、完整A000與色盤相同。這不等於完整RAM逐位元或新GUI驗證。
要求步數與實際`Step`逐筆同列，檔名不當精確印字停止點。

母體74次翻譯事件、25相異鍵，缺鍵及碰撞排除皆0。獨立重解析seed9720417的16鍵、首次步數及座標與原名單相符。
正取idx3的影片字幕在一次`+5M`仍不可見，依已審方法4用第一筆備用idx12遞補；不採`+8M`診斷當驗收。
五筆原零差異的包絡落在文字之外，重新目視定位後才重算，未由差分自動產生包絡。
12個鍵對應9張不同PNG，局部差異比率17.85%至28.82%，皆大於原門檻5%，
日誌點200px鄰域皆相交，交集寬高已保存。

主代理逐筆重讀後，另外將12張PNG去模式、步數及目標標籤，以固定seed打亂，交第二位唯讀子代理轉錄。
一筆字幕初讀前段未知，原PNG的最近鄰4倍裁片獨立讀出「技術程式設計」後解除。
12筆可見中文字與正式譯文一致；一筆題目有游標遮住一字，匿名讀者的未知保留，不宣稱該字已完整獨立辨讀。
已知目標的包絡修正不稱重新盲測，機器差異不替字形未知補成功。

| 新收據 | SHA-256 |
|---|---|
| takeover-check.py | `e285ec4d62f6caea16bd8850cfe4b6a4b26fb43c5feed64e57b705e7e83b14eb` |
| state-check.json | `3f42c36c7d6f661ad99d8dc9aeef43c90ce95cd00139bc2a088afc9c914f3a8b` |
| localization-check.json | `00a203f2e33acd268bf846c7ce93c4d9a6e8cfa178834363b4dedfc38674dc39` |
| source-sha256.json | `1cf6b50c2e3281bf8b6f1a2a60b1b403149a0d9270d86d7e5e8880602ce7376a` |

當時以離線非root容器內的暫存工具只讀舊收據，獨占建立新輸出。
repo唯讀、只有新收據工作區可寫，舊輸出未覆寫。204份來源SHA由唯讀審查全核同，
四檔catalog與記錄基準相符，xlate凍結SHA `08c1692e28b926dd754ca5f4df32ae3a114626dfdcbc944610edccc3c9fb28f4`保持。
當時輸出為`takeover-20261007/`，人工一致性與額外匿名讀值分別存`visual-consistency.json`及`visual-review.md`。
`FINAL-manifest.json` SHA `e2353da3df1bc41db02c0ce40bb7763e7a4000afebd300af7cd7e20d83d630ea`，
封存時再核204份來源及新輸出1000:1000，舊輸出未覆寫。

結論限既有路線、Unifont3倍及12筆抽樣的可見文字；未抽中、倚天、2倍、正式GUI、完整RAM、
十一組、TSV、存讀檔、音訊及封包均不外推。最終唯讀審查阻擋0、應改0，主代理於FINAL回讀後
採認12筆中文定位與可見字形覆核完成，審查記錄為新目錄的`review-final.md`。不以此段升級spec狀態。

2026-10-07，R973依使用者指示精確刪除已驗完的`workplace/r972/`，共379檔、26,773,719 bytes。
刪除前25個FINAL項目及204份來源SHA再次全同；刪除後目錄不存在。
正式摘要保留，詳細收據已移除，清理及接續入口見[本輪目標](../goal/R973-workspace-cleanup-paper-evidence.md)。

<a id="r973-paper-source-recovery"></a>
## R973 紙頁來源復核與工作區清理

2026-10-07，沿[R973](../goal/R973-workspace-cleanup-paper-evidence.md)，接續R944已定案紙頁內中文方案。
本機R944原型、比較頁與排版模型已不存在；先沿已提交
[片頭唯讀探針](../../tools/re/intro_picture_probe.go)重生正常無輸入原版來源。
來源檢查入口為[check_intro_source.py](../../tools/re/check_intro_source.py)，不呼叫前端或烘製器，
只核有／無觀察的原版摘要、完整indexed／色盤及獨立PNG解幀。正式025增補保持DRAFT。

工具為固定`wizardry7-go-ebiten:latest`映像ID
`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`，Go1.24.13、Python3.11.2。
原版`DS.EXE`為Archives版，SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
dosgolem本機基線`d16f03a2d8d782024d470606e1f4cb5272b89cd6`含既有未提交補丁，
本輪未改SDK；前後完整列出1134份非`.git`檔保持，不把HEAD當乾淨來源。
原版、repo及SDK唯讀，只有本輪輸出可寫，兩程序退出0，無玩家輸入、state或記憶體注入。

confirmed：兩趟各100個固定2M至200M停止點的`Step`、1MiB RAM FNV32與SHA-256摘要、
indexed／色盤FNV32與SHA-256、滑鼠及模式全同。indexed與色盤另逐bytes核同，
214張PNG獨立解析CRC、尺寸及按原色盤解幀相符。沒有保存RAM bytes或Regs，不稱完整CPU／RAM逐位元。
有觀察趟取得14張完整轉送畫格，14個3921入口的原始矩形0／0／319／199與
5B9B／5AFD返回後取幀順序另核同；SP＋2是原探針擷取條件，未新增返回暫存器收據。
既有六張呈獻及第一張紙頁的7個已記錄構圖SHA相同。
六張紙頁來源與當前色盤核對如下，步數只是收據定位，不作runtime畫字條件。

| 完整轉送步數 | indexed SHA-256 |
|---|---|
| 73235759 | `0f098eb7bf96d22777424ac0f8bc263534be1441ecc631413d61a2323714411a` |
| 75322912 | `304fef948a4eacdaab7cf0ba85e3a4e54e5f37d14fddc857e4f2a570b30a9c11` |
| 77352859 | `accbd133c3c5fd194d02ebc9b273dbde7688b7e5a3049df40e5ca66b2cf90ca1` |
| 100999086 | `5c53a427a0eaeea9026a2b6e2846380fe433c344e5b7b8ee4fb907145747e9d0` |
| 103102872 | `56bd0693ef1e4039426474aab61006046150e319201645385179da58f84606ee` |
| 105202085 | `f63f95bd99e5b95da6e732de86ec8a1317fca6c319a2e1f1ad7b2e5e1f5f2d07` |

六張palette SHA皆為`d689ae136905afa64a9329ddd123426f7e0f3cde6108c034e0249dbab0a9727c`。
無觀察控制沒有事件資料，觀察趟只計印字入口次數，不讀字串正文或手冊。
檢查器初稿把固定色盤誤套到最初廠商標誌；實際971527畫格色盤不同。
依025只限制已量構圖的契約修正核對範圍，以相同原收據乾淨重跑，沒有調原版參數。

新檢查器的7個拒絕案例通過：缺RAM摘要欄、全小寫欄名、截短SHA、改動indexed一個byte、
截短色盤、PNG CRC錯與既有輸出拒絕且bytes保持。前後2531份列出輸入全同，
含原版248、R971991、SDK1134、cmd63、既有tools80、text13與font2；新檢查工具另列來源。
root-owned與誤建`.md`目錄0。規格讀取曾誤寫001檔名，原版執行前已補讀實際scope契約。

| 本輪來源／當時報告 | SHA-256 |
|---|---|
| `intro_picture_probe.go` | `144b904d07896aaeef1c4fc63ccf1bbe0e2b3152c6daafac96e767956f241c58` |
| 原版探針binary | `8809954a18e9ab5cfaaa449a1cb1ef1c3d6ef530916ab9eb68dfd86efa8fa01b` |
| `check_intro_source.py` | `62d193d14cf3d0e21ea7aba3b22cb7552650dde4eb4a2b61abf169de933f75a3` |
| source-check-final.json | `7738dd22d5b415310be40dae68b983fedfedf248409a92f896474df8225f6a74` |
| observed/observation.json | `e3bd344632e02703eb2bb83d3a24af2c737525e9009d37f6f1c40bbba7df4522` |
| control/observation.json | `78c47e57639ff67ffa1bf2ed7f36dd490624e48228b103e71e117c832b0ead5c` |
| negative-check-final.json | `6f72f208c8ba884f7b2cf2e99c17a2ad76c228399b98dcb7b13c584ee854e26c` |
| integrity.json | `5d84942eacd1a16dd51470c7608f6364a22217cd8ee1d8ba0108bc52e789c7fa` |
| cleanup-r972.json | `a8c6af69c58f3568a2049d0f428b711da36fdec61cfc1b20ba5bb86f83950648` |

可重生來源的命令如下，只在上述離線非root有界容器內執行；`/src`為repo唯讀、
`/original`為自備`wizar7/`唯讀、`/out`為已存在且UID/GID正確的全新本輪輸出。
Go快取只寫`/out/gocache`，GOPROXY及GOSUMDB關閉，不借用主機runtime。

```text
GOCACHE=/out/gocache GOPROXY=off GOSUMDB=off GOFLAGS=-mod=readonly go build -o /out/intro-source-probe ./tools/re/intro_picture_probe.go
/out/intro-source-probe /original /out/observed on
/out/intro-source-probe /original /out/control off
python3 /src/tools/re/check_intro_source.py /out/observed /out/control --output /out/source-check.json
```

重生只恢復原版六種構圖，無法從舊集合SHA推回中文位置、補底模型或24張選定原型。
025紙頁增補仍DRAFT，須把可審查的字帶、原生字模位置、部分遮蔽及背景／撤除契約寫齊，
經獨立READY審查才接正式程式。本輪沒有新中文上屏、GUI、字型、TSV或全遊戲完成聲明。
獨立唯讀審查阻擋0、應改0、建議1，主代理完整回讀採認原版來源恢復及清理證據。
審查者獨立重算正式結果及7個拒絕案例相符，另驗`Mem=null`拒絕；
建議的建置環境變數已直接補入命令。025紙頁增補保持DRAFT。
依使用者驗收後清理要求，審查採認後已刪除`workplace/r973/`共1427檔、77,286,736 bytes，
包含完整畫格、私人報告、binary與Go快取；刪除前後2531份列出輸入全同。
R972與R973目錄均不存在。正式探針、獨立檢查器與上述摘要保留，可從自備原版重生來源。

<a id="r974-paper-contract-rebuild"></a>
## R974 紙頁中文候選重建

2026-10-07，沿[本輪目標](../goal/R974-paper-contract-rebuild.md)，接續R973原版來源復核。
重啟前已留下六張完整轉送畫格，`source/observation.json` SHA仍為
`e3bd344632e02703eb2bb83d3a24af2c737525e9009d37f6f1c40bbba7df4522`，
六張indexed及固定色盤完整SHA與本檔R973表相同。接手後只重建候選，未再次跑原版。

使用既有Docker映像`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`、
Python3.11.2；repo、原版畫格與兩字型唯讀，只有本輪輸出可寫，UID/GID1000:1000。
原版版本與DS.EXE SHA沿R973；定位是完整畫格的indexed SHA，步數只作來源出處。
沒有讀取MSG、手冊、答案或其他原版文字，也沒有原版記憶體寫入。

重建入口：[prototype_intro_paper.py](../../tools/re/prototype_intro_paper.py)、
[文字與座標候選](../../tools/re/intro_paper_candidate.json)，
獨立檢查入口：[check_intro_paper_candidate.py](../../tools/re/check_intro_paper_candidate.py)。
這些來源進版控，原圖、原字型、候選PNG與詳細收據只留本機。
文案是「呈獻」「奇幻／角色扮演／遊戲」「作者」，作者姓名保留英文。
字帶、可見紙面、字格中心與作者保護區均是新候選，不當原版幾何已證實。

候選背景限制在字帶與可見紙面的交集。同一紙面排除所有候選字帶及作者保護區後，
只取R>G>B、5G≥3R及G−B≥16的紙色，按x／y奇偶分四組，各至少八點。
每個目標點取同組最近八點，權重為`1 / ((sx-x)^2 + (sy-y)^2 + 1)`，
各RGB通道加權平均後採四捨六入五成雙。距離相同時按y／x排序。
這是周邊紙紋近似，未解原版無字底圖；補底仍有明顯接縫。
部分遮蔽頁三行移到可見紙面，不縮字、不旋轉、不裁字。

最終候選為`workplace/r974/candidate-v4-unifont/`與`candidate-v4-eten/`，
每套六張、兩倍率，共24張。原生Unifont16及倚天15字模保持，字距2像素。
第六張沒有中文或補底，完整原圖保持。主代理逐張檢視五張有字Unifont3倍候選，
外觀仍待使用者確認；[本機對照入口](../../workplace/r974/compare.html)保留待查看。

獨立工具直接解析固定HEX及倚天Big5字模，不匯入候選的印字、補底或多邊形函式。
24張共9,984,000像素，12個原生字元、11,458個前景點相符，
字格出候選紙面、字格重疊、候選作者保護區變更、批准補底及字模筆畫外差異皆0。
48張原版／候選PNG核CRC、尺寸及完整原版色盤展開。
不把候選座標本身當獨立原版oracle；補底區內字模零位元的背景、英文殘筆與補底品質沒有獨立期望。
沒有正式GUI、新舊原版同狀態中文覆繪或全遊戲完成聲明。

`source-before.json`列出的1543份既有輸入與兩套固定字型保持。
本機Unifont原路徑已不存在，本輪沿現存相同SHA來源，
`/home/anr2/cht/golden_box/拯救地球/workplace/unifont-src/unifont_all-17.0.05.hex.gz`唯讀掛載；
倚天沿`/home/anr2/cht/etan_font/stdfont.15`。工具依固定SHA拒絕其他版本。
Unifont SHA為`7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574`；
倚天SHA為`39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701`。
輸出擁有權及誤建`.md`目錄異常0；自有容器清理與獨立審查另補。
審查指出先前把本輪根目錄掛到`/out`，讓同目錄`source/`產生可寫別名。
工具未寫該來源，六張indexed／palette與觀察JSON保持既有SHA；仍屬隔離配置缺口。
以下重跑命令改成獨立輸出子目錄，repo及原版畫格只由`/src:ro`讀取。
修正後已在獨立`rebuild-output`掛載乾淨重建雙字型24張；52份PNG／收據／HTML
與已審v4逐位元相同，輸出1000:1000。原版來源僅從repo唯讀掛載讀取，沒有可寫來源別名。

| 來源／本機收據 | SHA-256 |
|---|---|
| `prototype_intro_paper.py` | `216c0d4a497593382abdda32cd6d9831ee4e2d55f7a37d3e0c0ddd3b7c822660` |
| `intro_paper_candidate.json` | `1a9632b661b106367adffc2a8336808fe38202e60080620e897072e2fdcf5a22` |
| `check_intro_paper_candidate.py` | `fbb98f530606d540f5da13642461909411e4a0a8f9f00934411106386a91f9e3` |
| `candidate-check-v4.json` | `501c10dc6c35395d4df6eb3fa898d5ad375909972612376b0f21efc3d3dd01ab` |
| `integrity-v4.json` | `161900242f90c24e7e656f194b9646d48370b8882607e791e65e415ef8acca53` |
| Unifont `receipt.json` | `e93d35efc9a3e4e3c3fd000d91090ea3b0697d6087c5ba9e803cbbeabda22ea9` |
| 倚天 `receipt.json` | `aab6a4ed8f008c57a6f40951fa4c96b21df11ce36fd77cb8b91e598871a1e37b` |
| `rebuild-output/rebuild-check.json` | `a752aaf4043d294f42e7754a6395304c372cbbc17e505c087c1d402aebebec95` |
| `review-prototype.json` | `56448a977e0214949a51cd8182a98acba630da593b92660590e0a124f8ba30e8` |
| `review-prototype.md` | `ad1b86c5c554ec61daae0f999fd32a4a77a291847f8892ab8cd4180b3b6b9a36` |

重跑命令只在非root、離線、有界Docker內執行。`/src`為repo唯讀、`/font`及`/eten`為固定原字型唯讀，
`/out`為已存在的新輸出目錄；各候選目錄或報告已存在時拒絕覆寫。

```text
PYTHONDONTWRITEBYTECODE=1 WIZ7_FONT=unifont WIZ7_UNIFONT_HEX=/font/unifont_all-17.0.05.hex.gz python3 /src/tools/re/prototype_intro_paper.py /src/workplace/r974/source /out/candidate-unifont --eten-font /eten/stdfont.15
PYTHONDONTWRITEBYTECODE=1 WIZ7_FONT=eten python3 /src/tools/re/prototype_intro_paper.py /src/workplace/r974/source /out/candidate-eten --eten-font /eten/stdfont.15
python3 /src/tools/re/check_intro_paper_candidate.py /src/workplace/r974/source /out/candidate-unifont /out/candidate-eten --unifont-hex /font/unifont_all-17.0.05.hex.gz --eten-font /eten/stdfont.15 --output /out/candidate-check.json
```

本機有界容器命令如下。先在非root容器內建立獨立`rebuild-output`並核UID/GID；
輸出目錄不存在時，下面的`test -d`會停止，不讓dockerd代建。所有掛載來源均先驗型態。
以全新候選目錄代替既有`candidate-new-unifont`，不覆寫先前結果。

```sh
set -eu
test -d /home/anr2/cht/wizardry7
test -d /home/anr2/cht/wizardry7/workplace/r974/rebuild-output
test -d /home/anr2/cht/etan_font
test -d '/home/anr2/cht/golden_box/拯救地球/workplace/unifont-src'
timeout --kill-after=5s 120s docker run --rm --name wiz7-r974-rebuild-unifont \
  --network none --memory 1g --cpus 2 --pids-limit 64 \
  --log-opt max-size=10m --log-opt max-file=3 --user "$(id -u):$(id -g)" \
  -e PYTHONDONTWRITEBYTECODE=1 -e WIZ7_FONT=unifont \
  -e WIZ7_UNIFONT_HEX=/font/unifont_all-17.0.05.hex.gz \
  -v /home/anr2/cht/wizardry7:/src:ro \
  -v /home/anr2/cht/wizardry7/workplace/r974/rebuild-output:/out \
  -v /home/anr2/cht/etan_font:/eten:ro \
  -v '/home/anr2/cht/golden_box/拯救地球/workplace/unifont-src:/font:ro' \
  -w /src --entrypoint python3 wizardry7-go-ebiten:latest \
  tools/re/prototype_intro_paper.py /src/workplace/r974/source \
  /out/candidate-new-unifont --eten-font /eten/stdfont.15
```

025紙頁增補仍DRAFT，正式程式與R943六張呈獻保持。
外觀確認、共同HD／Draw、完整來源見證、背景回退及撤除契約、獨立READY審查與正式驗收仍待。
十一組、TSV、錄放、冷存讀、音訊及封包門檻保持。待查看候選不提前刪除。

唯讀子代理最終報告[JSON](../../workplace/r974/review-prototype.json)與
[Markdown](../../workplace/r974/review-prototype.md)已完整回讀、核SHA並採認。
候選展示未解除阻擋／應改／建議為0／0／1；建議後續新候選把工具與字型識別內嵌receipt，
本輪外部研究及報告已固定，不重寫已生成候選。隔離配置應改已修正並由主代理乾淨重建核同。
審查自寫PNG、整數面積幾何及兩字型讀法，獨立重算48PNG／24候選／9,984,000像素／11,458前景點相符，
75份列出審查輸入前後SHA保持；沒跑主檢查器、重建工具或原版。
正式READY仍有外觀確認、完整技術契約及可獨立驗收顯示契約三項閘門；
正常GUI及同狀態中文覆繪須於正式實作後完成，不以候選通過跳關。
本輪容器及擁有權異常0，未追蹤R433保持；候選及來源保留待查看。

<a id="r975-paper-letter-weight"></a>
## R975 紙頁字樣放大、加粗與提亮

2026-10-07，沿[本輪目標](../goal/R975-paper-letter-weight.md)，處理使用者指出R974字偏小、
要求原版般厚實及稍亮的字樣。R974的15／16像素字形在960×600畫布仍為1:1，
原版圖片文字則隨畫布放大，是相對大小不同的直接原因。
本輪不改正式字型或前端；新紙頁效果由[可重跑工具](../../tools/re/prototype_intro_paper_weight.py)產生。

原版來源沿R973／R974六張紙頁完整SHA，實際只用前五張有字畫格。
原版DS.EXE版本與SHA沿本檔R973；未再次執行原版，不讀MSG、手冊或答案。
背景固定R974的`1a9632b6…`模型與`216c0d4a…`工具，完整SHA由工具守門。
原版畫格、原字模及R974比較來源只讀。背景仍有條紋，此輪不改它，供字樣並列比較。

兩套原字模均完整整數放大2倍為倚天30及Unifont32像素。
甲將放大後筆畫向右增1像素，乙不額外加粗；字距4，中文字保持水平。
紅字由R974的156／20／12改為196／36／20，深色字由原版palette[0]改為40／28／16。
亮度相對R974提高，並不宣稱已量原版所有圖片字樣的同一亮度。
五張構圖、雙字型、兩厚度共20候選，另10原版及10R974對照，均為960×600。
傾斜頁將「角色扮演」分兩行；呈獻及翻起頁作者中心微移，完整字格避免出紙面與作者區。
這是新候選幾何，沒有原版排版精確聲明。

本機入口：[字樣對照](../../workplace/r975/compare.html)，最終輸出在
`workplace/r975/final-output/unifont/`與`eten/`。
完整收據保存字模、工具SHA、Python版本、固定R974比較收據SHA及逐張PNG SHA。
原圖、字型與候選PNG只留本機，譯文與正常遊戲資料保持。

原型先由既有wizardry7-go-ebiten映像`1d7a342d…`與Python3.11.2產生。
其後Docker回報原名稱及完整ID均不存在，原因未知，沒有自行刪除或另建映像。
本輪標準函式庫工具改沿現存`python:3.13-bookworm`，實測Python3.13.15，固定完整ID
`933b46a028fd786c9c3d426ebabc237e29a15912231ea8de576e95f0e4f41a4c`。
同來源乾淨重建44份檔案與舊版逐位元相同；補來源守門與metadata後，再核40PNG完全相同，
收據變動只及工具SHA、Python版本及R974比較來源SHA。
這個映像只支援本輪Python原型／檢查，沒有Go／GUI工具鏈恢復或前端驗收聲明。
正式工具建置來源仍為[go-ebiten.Dockerfile](../../tools/docker/go-ebiten.Dockerfile)。

生成及重建均為離線、UID/GID1000:1000、1GiB／2CPU／64pids、有界一次性容器。
repo、來源及字型唯讀，僅獨立輸出子目錄可寫，不產生來源可寫別名。
本輪所有輸出擁有權及誤建`.md`目錄異常0。
候選尚未定版，025仍DRAFT，沒有正式GUI、原版中英文同狀態或2倍畫布驗收。

獨立唯讀審查[JSON](../../workplace/r975/review-weight.json)與
[Markdown](../../workplace/r975/review-weight.md)已完整回讀、核SHA並採認。
展示無未解除阻擋／應改；原先R974比較來源缺固定SHA已補工具及收據守門。
建議由使用者看甲乙後選字重，沒有替使用者決定新外觀或放行正式READY。

審查自寫PNG、Big5／HEX、整數winding幾何及固定背景公式，不匯入被測工具。
40PNG共23,040,000像素，20候選11,520,000像素逐點相符。
12字×2來源共24原字模、168次出現、11,458個來源筆畫點完整放大，前景共51,674點。
72字條完整，字格越界／交疊、候選作者保護區及宣告範圍外差異0；
76份列出審查輸入的SHA及UID/GID前後保持。兩字重與R974的固定背景一致。
這些結果只證明新字樣符合候選契約，不裁決補底品質、完整作者幾何或正式玩家顯示路徑。

| 來源／本機收據 | SHA-256 |
|---|---|
| `prototype_intro_paper_weight.py` | `01206ea5ed2e99d570fca5f42fd59a33b5e969ae944a1a9856c2911627825f44` |
| 最終Unifont `receipt.json` | `e9beb480bc01eb6a18962451e597d74373590161e0b1808f17baa5d8d78b43a9` |
| 最終倚天 `receipt.json` | `1981685b5f95aea56819f6bb8fbd345ee848717d667601e5eb1bdc235da5aece` |
| `rebuild-313/cross-runtime.json` | `4754859755bc461215e50cf06d1811866cfa68daab5daec2dbc6185bda8d2513` |
| `final-output/final-equality.json` | `78ae62c21260c93485ac3f6379824e592dee66b9c10ac351ebe360b63aff4453` |
| `review-weight.json` | `a99c9dc8281e2f7159a97ba477ef1e866b26fb0e1d2b34935de27fc1b06e6f5f` |
| `review-weight.md` | `5b0cc8e821092d4dbe57ffa55736bb22f3fd33ce11e6dd0ded9497e0a0eff39e` |

重跑沿上述固定Python映像。repo與來源唯讀掛`/src`，兩字型唯讀掛`/eten`及`/font`；
已存在且UID/GID正確的全新輸出子目錄掛`/out`，不能掛含R974來源的父目錄作可寫輸出。
沿本輪1GiB／2CPU／64pids、離線、非root、log上限及150秒外層逾時；掛載來源先核型態。

```text
PYTHONDONTWRITEBYTECODE=1 WIZ7_FONT=unifont WIZ7_UNIFONT_HEX=/font/unifont_all-17.0.05.hex.gz python3 /src/tools/re/prototype_intro_paper_weight.py /src/workplace/r974/source /src/workplace/r974/candidate-v4-unifont /out/unifont --eten-font /eten/stdfont.15
PYTHONDONTWRITEBYTECODE=1 WIZ7_FONT=eten python3 /src/tools/re/prototype_intro_paper_weight.py /src/workplace/r974/source /src/workplace/r974/candidate-v4-eten /out/eten --eten-font /eten/stdfont.15
```

唯一獨立核對程式保存在審查JSON的`verification_program`，只在同固定Python容器執行，
repo及兩原字型只讀；實際掛載與本輪範圍見審查Markdown，不作主機執行入口。
字樣選擇、補底改善及正式025契約仍待；候選保留供使用者查看，自有容器清理與推送另補。

<a id="r976-bright-paper-brush-menu"></a>
## R976 紙頁亮字與六按鈕毛筆字樣

2026-10-08，沿[本輪目標](../goal/R976-bright-paper-brush-menu.md)。使用者採R975甲的大小、加粗及排版，
指出文字太暗，要求更亮；毛筆體範圍另確認為主選單六個按鈕。
紙頁只改已採甲的前景RGB；主選單先做符合020按鈕尺寸的字樣樣張，未改正式前端。
本機入口為[亮字與毛筆對照](../../workplace/r976/compare.html)。

### 紙頁亮字

[亮度工具](../../tools/re/prototype_intro_paper_brightness.py)只讀固定R975兩套最終收據及五張甲圖。
Unifont收據SHA `e9beb480bc01eb6a18962451e597d74373590161e0b1808f17baa5d8d78b43a9`，
倚天收據SHA `1981685b5f95aea56819f6bb8fbd345ee848717d667601e5eb1bdc235da5aece`。
深色字40／28／16改為米金248／232／160，紅字196／36／20改為236／72／40。
只在原甲字模前景換色，30／32字高、向右增1像素、字距、位置與背景全部保持。

兩套各五張新圖與五張before-copy，共20PNG。獨立審查逐點核新圖5,760,000像素，
只有28,758個前景點換色；新舊形狀、位置、背景與R975候選作者區差異0。
原甲圖及收據完整SHA保持，不從使用者附圖的plain檔名推翻明示採甲的指示。
不稱原版CPU、記憶體或正式GUI的同狀態新驗收，R974補底條紋仍在。

### 六按鈕毛筆字樣

[毛筆預覽](../../tools/re/preview_menu_brush.py)只取正式catalog的六個英文鍵，
文案為「角色選單」「新增角色」「載入存檔」「系統設定」「匯入巫術六」「結束遊戲」。
來源ID只作出處；既有catalog沒有修改。六鍵與按鈕尺寸來源為[020](../spec/020-main-menu-labels.md)。

使用本機既有AR PL UKai TW／Book字面，ukai.ttc索引2，
字型SHA `cfc758e69386a21bade8cb3b373fda20d7f2fac6e8904ced2ded98588b63e274`。
原授權SHA `1f11c4a22865e1289f4694dad56b3561790a4b7341a511ddcd0edddcc3b95a46`；
原版套件檔0:0只讀，不修其擁有權，輸出與授權原文副本為1000:1000。
字型／衍生字模依ARPHIC PUBLIC LICENSE，與程式RRSAL分開；來源及條款見
[字型來源](../../font/brush-menu-source.md)與[Debian保存的原授權](https://sources.debian.org/src/fonts-arphic-ukai/0.2.20080216.2-5/debian/copyright/)。

[FreeType工具](../../tools/re/font_preview_freetype.py)使用容器內原生庫及標頭，
Linux64 ABI已與原生C API交叉核對；不掛主機runtime或library。
字模以FreeType2.12.1直接渲染16／24像素，8位alpha、每字advance加1像素，
字色248／232／160；按020的109×9按鈕尺寸排2倍與3倍，字模沒有縮放或裁切。
按鈕及背景為自製字樣板，未使用遊戲截圖，不當正常GUI或原版收據。

獨立C API讀取glyph index後render，交叉核44個字模的完整alpha與metrics一致，
12字條缺字及越界0；兩樣張153,386個合成像素逐點差異0。
此結果證明字樣工具與字型來源一致，未放行正式字型接線、原版輸出辨識、按鈕操作或存讀檔。

### 環境與結果採認

沿固定Python3.13.15映像
`933b46a028fd786c9c3d426ebabc237e29a15912231ea8de576e95f0e4f41a4c`，
FreeType2.12.1；生成為512MiB／1CPU／64pids，離線、UID/GID1000:1000，有界一次性容器與log上限。
repo、R975來源、字型與授權唯讀，僅R976獨立輸出可寫；重啟前只寫工具，恢復後確認容器0才開始產圖。
GitHub開工留言已查到[R976入口](https://github.com/wicanr2/wizardry7-cht/issues/22#issuecomment-6043201880)，沒有重複發送。

獨立唯讀報告[JSON](../../workplace/r976/review.json)及[Markdown](../../workplace/r976/review.md)已完整回讀、核SHA並採認。
審查逐點核22成果PNG共11,673,386像素，49份讀入來源SHA及UID/GID保持。
展示阻擋／應改／建議為0／0／1，建議看過原尺寸字樣後確認；正式020／025增補仍DRAFT。
沒有原版執行、messages.json／MSG.DBS、手冊或答案讀取，原版與正常遊戲語意路徑保持。

| 工具／本機收據 | SHA-256 |
|---|---|
| 紙頁亮度工具 | `7ae7f84d37c1dfb78a353c3a4579874aebd771c4307346fc06d3a81e043ccfe8` |
| 毛筆樣張工具 | `89e3969b33f42c7eeb4fb8d91ca49ef69929669d93d5721424d217f36a17732c` |
| FreeType工具 | `36c6887c8b538defb457015de4917fa770ce8c35c3203b51b9372ec3cdc811ee` |
| 紙頁Unifont收據 | `afb796b46964f907e68c3a40eeb8af79d93264d3bc1426d140863ff12010cd49` |
| 紙頁倚天收據 | `ea19ec75516531815b6f7a654b755e3f8ae22abef4e96c0bf4c44581487e82b9` |
| 毛筆選單收據 | `5d94afe72850c66b4b2526bba048b147988ccbd7a454b68fad92cf9420afd623` |
| 獨立審查JSON | `8e7387cb3dd192eb67aa9e0ffbc446b1f9b4f4a857b55bb492d2161ac844f0ec` |
| 獨立審查Markdown | `71c558edcfe52c8d607c8a60743eaba8c620debafa964890d863a2c56d86a90f` |

重跑在上述固定容器執行。repo唯讀掛`/src`、AR PL字型目錄唯讀掛`/brush-fonts`、
原套件授權目錄唯讀掛`/font-license`；僅全新輸出目錄可寫掛`/out`，掛載前核存在、形態及UID/GID。
輸出已存在時拒絕覆寫，原型只留本機。

```text
PYTHONDONTWRITEBYTECODE=1 python3 /src/tools/re/prototype_intro_paper_brightness.py /src/workplace/r975/final-output/unifont /out/paper-unifont
PYTHONDONTWRITEBYTECODE=1 python3 /src/tools/re/prototype_intro_paper_brightness.py /src/workplace/r975/final-output/eten /out/paper-eten
PYTHONDONTWRITEBYTECODE=1 python3 /src/tools/re/preview_menu_brush.py /brush-fonts/ukai.ttc /font-license/copyright /out/menu
```

唯一獨立重跑程式保存在審查JSON，需同映像內原生C編譯器／FreeType標頭，所有來源唯讀。
紙頁新亮度及六按鈕具體字樣待確認；源字型與原版美術不入Git，候選保留待查看。
紙頁補底、正式背景／辨識／HD／Draw／回退／撤除及正常GUI／同狀態收據仍待，Go／GUI工具鏈未恢復。

## R955 頭像頁文字與完成返回

2026-10-04至10-05，入口：[本輪目標](../goal/R955-character-portrait-ui.md)。沿R952已建正式Linux前端，binary SHA `7363fc5be7c48d10eb272d5bd5900045996bb490b0cb08fd19bb1a93e9d948b0`。38份列出cmd來源保持，不稱新建置或完整編譯閉集合；工具鏈沿固定Docker映像、Go1.24.13／Python3.11.2，原版`DS.EXE` SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。repo、SDK、原版與字型唯讀，只有r955可寫；沒有MSG、防拷、內建autoplay、原版記憶體或角色資料注入。

- 真正3倍Unifont GUI逐張檢視實際Draw後送六次點擊：角色選單160,47、頭像236,193、第一列175,152、選人確定100,170、完成160,155、離開287,193。未選新頭像；畫面提供完成控制，不把它稱為取消。七個停止畫格、正常返回主選單與關窗exit0，耗時0.100491952秒。
- 原版248檔SHA保持，擷取後七份Scratch雜湊及退出後PCFILE.DBS 9368 bytes保持；最後另核八個Scratch與原版逐位元相同。Scratch讀取在擷取後，不冒稱恰在固定Step。兩則文字的10個原生Unifont字模、兩個安全格及HD／Draw局部19,968像素通過；14張PNG的CRC、解碼、尺寸與SHA相符。
- 另有雙字型、兩倍率及英文控制六趟固定輸入輔助，全退出0。四組中文各七停止點，共28對Step、1MiB記憶體FNV、A000、Pal、Mouse摘要與完整indexed／palette bytes相同，不稱完整CPU或新人工GUI對拍。40個原字模實例、八個安全格與局部HD／Draw 64,480像素通過；84張PNG的CRC、解碼、尺寸與SHA相符。
- 兩則文字的原版1524個改動像素另與英文indexed核對，純色背景由英文控制的相鄰空列核實。原倚天按Big5分區與大端15列解字，十字皆為漢字，不讀符號區；原Unifont按16列解字。倚天392820 bytes、SHA `39ba9c8519d75fe11d5988a8a27e6daa5794ad2ea215108390b0d7e9e53ff701`，Unifont原壓縮檔SHA `7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574`。

| 原版英文鍵 | 原版座標 | 本輪限定結果 |
|---|---|---|
| SELECT PORTRAIT | 116,138 | 「選擇頭像」雙字型兩倍率字格與局部像素通過；未選新頭像 |
| CLICK HERE WHEN DONE | 101,153 | 「完成後點這裡」雙字型兩倍率字格與局部像素通過；160,155完成返回 |

首次正常GUI碰到外層600秒期限exit124，已有頭像頁完成返回角色選單的畫格，但最後離開命令送達未確認，沒有成功退出收據。保留片段，以相同binary、時基、七停止點、六次點擊及600秒期限乾淨重跑；只改工具輸出路徑與ready雜湊日誌。第二次成功，前六張Draw與首次逐位元相同。首次Docker inspect誤用不存在的Entrypoint欄位，改核image ID；均屬工具問題，未調遊戲參數或列產品缺陷。

正常收據位於`workplace/r955/normal-unifont-3-attempt2/captures.json`，SHA `1abeace7d4c1426e542e16141717a2e35d05ecb7db0e911a30840fc582133dec`。私人獨立檢查入口為`check_portrait_gui.py`／`check_aux.py`，SHA依序`1355a28df7121282a27554532d9ea24736f107e8a68ebc863ce28082293da8d0`／`98284ad1971ad31f46cec2f256912f6a1dcbbd002d5a70719a56243f9a58ec14`。沿固定離線非root映像，repo與字型唯讀、r955掛`/out:rw`；後者另需倚天目錄唯讀掛`/eten:ro`及六份aux收據。報告獨占建立，不重跑原版或渲染器。擁有權異常0，正式程式及譯文未改；只有頭像頁兩則文字與保持原頭像返回，不消除已知2倍人物／名單或其他圖片控制缺口。十一組、TSV、錄放、冷存讀、音訊與封包門檻保持，R954／R944／R951待選不代答。

正常報告`portrait-gui-check.json` SHA `0bb142d9465f66e39bfa1f7170b003d4b32f2fc3213c097f46e048f8b399d26e`，輔助報告`aux-check.json` SHA `e8d73d715205494a63acc3be20bb304812595d42a3cdca5d6e2a111b2576d770`；`integrity.json` SHA `e497df462dd43cbce27eabc5bace611099ac62f0c07c0715ffc2c85347aa99ea`回查19份選定來源／收據及八個Scratch完整bytes。工作根衛生異常0、自有容器0。SDK只有既有修改狀態，未宣稱本輪SDK1132份新雜湊核對。

R955獨立唯讀審查阻擋0／應改0／建議0，主代理於FINAL後完整核讀兩報告並核SHA。`workplace/r955/review-results.json`／`.md` SHA依序為`2ea5c896c905c0c79d1c33e258ecf5175ed8d41d10dcc381d9ef25bd7b82213c`／`f2cb18a06319eedf1cf1a5938beb70cb0cf549fa0aa8ad41d5549c0dde4db72e`。審查另獨立抽核兩套原字型、四個字模、六張PNG及2948個局部RGB像素；既有檢查器的19,968與64,480像素範圍分列，不稱審查重跑整批。只採認兩則頭像文字及保持原頭像完成返回，不新增READY或完整介面完成聲明。

現行R943已完成片頭六張「呈獻」的正式覆繪及限定正常GUI驗收，入口見本檔「R943 正式呈獻覆繪驗收」。其他圖片敘述、按鈕與完整交付仍待。

狀態：主選單、角色選單、新增角色對話框及3倍人物檢視限定像素驗收通過，R941結果複核通過，其餘介面重查中；2026-10-04；[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)。
入口：[R915](../goal/R915-ui-localization-audit.md)、[R921](../goal/R921-character-dialog-ui.md)、[R928](../goal/R928-character-roster-text.md)、[R932](../goal/R932-character-menu-audit.md)、[R935](../goal/R935-character-review-audit.md)、[R942](../goal/R942-intro-picture-text.md)、[目前狀態](../../CONTEXT.md)。

使用者要求把圖片內的英文按鈕納入完整中文化。本文保存目前正式來源及實際畫面的覆蓋矩陣；譯文存在、原版畫格及舊測試完成數量不直接當作本次介面已驗證。

## R954 動態背景與等級欄修正

2026-10-04，使用者已選定R953重排，沿[本輪目標](../goal/R954-profile-reflow.md)補背景守門。仍是唯讀原版實驗及離線候選，正式程式／SDK未改；024既有READY範圍保持。原型圖片依原版素材權利只留本機，定案表與雜湊進Git。

沿R953原版`DS.EXE`及固定Docker／Go1.24.13／Python3.11.2工具鏈，另建探針，只補取性別入口64000-byte背景及已量畫面的呼叫metadata，不讀MSG、鄰文、詞庫或答案，不寫原版。男女各13事件及7停止點，observation與R953完整相同，indexed／palette亦沿既有14對控制。兩趟退出0採工具實際終止回報。原版248檔SHA與R953相同，兩個Scratch的PCFILE.DBS保持；不稱新正常GUI或完整CPU驗收。

下表地址為原版runtime `0110`段，font均0；訊息原文仍是覆繪比對鍵，表中返回位址只補辨識結構。

| 標籤 | 原版座標 | near返回位址 |
|---|---|---|
| 種族 | 102,6 | 6c7a |
| 職業 | 90,12 | 6cc0 |
| 等級 | 90,18 | 6cf8 |
| 階級／階級名稱 | 228,6／288,6 | 6e13／6e6b |
| 經驗／擊殺 | 228,12／228,18 | 6ea3／6ef9 |
| 生命／體力 | 6,93／6,99 | 6f5d／7007 |
| 八屬性 | 115,35至77、每6一列 | 6a43 |

生命／體力的下一輸出邊界`0110:3895`、返回位址032a、225,138實測font3，解釋R953 V1的整體font0篩選漏掉該邊界。此處仍是時間順序邊界，不是匹配函式返回。八屬性在該邊界之後才輸出，最後來源返回分別為男666672690／女666594273；不能把生命兩列完成當作整頁完成。

以原版色號、R953十三個來源及R937十七個標籤的精確before／after還原筆畫，再比較性別入口與784M完整畫格的候選格。原選版面7970個DOS像素中，扣來源筆畫後，男7029格有8個變更、女6989格有25個變更，皆位於等級字格內的蠟燭圖案。主代理檢視原版相位與四份修正版；新字模與上述變更像素相交的輸出筆畫為男倚天3／Unifont7、女倚天22／Unifont27。這是confirmed的兩停止相位交集，不代表已量全部蠟燭相位。

原選位置因此無法直接通過背景守門。可撤回甲案只將等級文字右移24輸出像素，字格204,51,34,18；等級數值安全格右移12至240,51,36,18，實際數值墨點保持。其餘30個顯示欄位保持。甲案的7900個DOS目標格，扣來源筆畫後男女6959／6919個背景像素變更0；這是兩停止相位證據，正式每次合成仍需守門與撤除。乙保留使用者原選位置，但須明確接受文字覆在蠟燭上，不能默認這項忠實度取捨。

甲雙字型、男女四份重排及四份現行對照八張PNG，獨立2,048,000像素相符；只有四張重排圖各31格交疊／墨點交集0。字模沿R937已核原字型資源，沒有本輪新rawfont oracle。獨立checker的第一份摘要未標新幾何輸入，補入checker、八張PNG及新layout的SHA後另寫v2報告；未覆寫初份，也不以相同報告SHA聲稱新來源已封閉。完整[並列頁](../../workplace/r954/level-revision.html)為私人離線原型，不是正式HD／Draw驗收。

| 本機workplace/r954收據 | SHA-256 |
|---|---|
| `profile_background_probe.go` | `8a1eac669b70e605978bf4f068e58be28facc9a27b289b228674868f07fb22ac` |
| `profile-background-probe` | `50df293b2416e9c23f22c408a85e5bae474baefd8c9b00ac5245c0219326b7cd` |
| `diagnose_background.py` | `d4ee654f2520865ddee6e06e16c52ba4b225f5dadadf50d51842a69bf5358a64` |
| `candle-conflict.json` | `70bf0c70961aa0b854620aed18555238dc7f4ef4683cd1b4bc3cc0ba833df928` |
| `background-check.json` | `dd7266239d117b1fc4ec011d85b709db5c6f897aa8758d41adf9c62414b458a3` |
| `prototype_level_revision.py` | `1b1375a7aafb4940f964fa3609276ca383ce119573037c3341827097e0fcae56` |
| `check_revision.py` | `3856fb5cbae022510aadd9d994e1da96d4407fbdc21f4858e430b089a84ee857` |
| `independent-revision-check-v2.json` | `04b21baa6eb6d9c360dce39ca5eb414b65a52084bfb0468f690155a787a9293c` |
| `prototype-v2/result.json` | `0399fc54ed9422eef466d955b14a51021269e7d36bc5dc41f3c1ebca3f11e5f6` |
| `level-revision.html` | `76377b5f1c3ae716f921789849a2282410a344bbe3bc454b2394cffbf71200d5` |

重跑沿既有固定映像，repo、原版及R953／R937／R943輸入唯讀，新且已存在的輸出目錄掛`/out:rw`。以相同唯讀go.mod/sum、受控SDK及模組快取、可寫Go快取建`profile_background_probe.go`，男女參數與R953相同。最終候選入口依序為`prototype_level_revision.py`、`check_revision.py`、`check_background_v2.py`，另用`diagnose_background.py`量原選位置；須提供新探針男女輸出、原版基線及上述舊收據。輸出既有則拒絕，不重跑所有草稿，也不拿缺輸入的執行當通過。

首次背景checker因真正動畫差異拒絕；診斷器一度把男的8格／3／7套給女，屬先驗計數錯誤，改成各自實測後量到25格／22／27。初份腳本與比較頁說明保留，最終說明已修。正式位置未自動改成甲，待使用者選定最小右移或明示接受覆在動畫上；技術契約未READY、未正式實作，十一組與其餘交付門檻保持。

R954獨立唯讀審查0／0／0，主代理於FINAL後完整回讀JSON與Markdown，只採認限定RE與離線原型。報告位於`workplace/r954/review-background.json`及`.md`，SHA依序為`7da897e752b430f0cc04ce37bf746ec3c2da23574371206a5a1a0f4c80466194`／`aee0b7b124c16b63902426504d6d45772eb834bcf3497d07c54c4261a896fa69`。另獨立重算兩相位背景、八圖完整RGB及最小變更範圍；不批准甲乙取捨、新READY或正式實作。審查工具兩次定位／欄位錯誤已修，不列產品失敗；自有容器0。

## R953 兩倍人物欄位原型

2026-10-04，入口：[本輪目標](../goal/R953-profile-layout-evidence.md)、[排版候選](../spec/024-character-review-output.md#r953-兩倍人物排版候選)。這是唯讀原版測量及離線原型，未接正式前端，也不是新GUI驗收。

原版`DS.EXE`為69,206 bytes，SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。工具為Go1.24.13 linux/amd64、Python3.11.2及既有dosgolem隔離工作樹；Docker沿固定映像`1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`。下表地址皆為原版runtime段址，不是檔案偏移或IDA線性址。

| 已量事件 | 入口、辨識及結果 |
|---|---|
| 性別 | `0110:3877`，返回`0110:6c42`；font0，90,6→96,6，SP增加2。擴大取88,4,8,8的背景，M／F差異13／11像素且範圍外0。 |
| 十一個數值 | `0110:3895`，返回`0110:032a`；八屬性y35至77、每6一列，兩位數x139、一位數x145；等級126,18，經驗306,12，擊殺306,18。font均0，返回SP增加2，游標前進字數×6，差異只在各自6像素高字格。 |
| 生命／體力兩列 | 已量6,93入口至`0110:3895`、返回位址032a、游標225,138的下一輸出入口，取6,93,78,12的差異。未讀字串正文；男女177／181個變更像素、範圍外0。這是時間順序邊界，SP差為−306，不是匹配函式返回。 |

探針只在性別返回已通過後、664M至784M之間讀已量數值座標的目的緩衝，最多8 bytes，只接受1至3位十進位數字及NUL。沒有讀MSG、鄰文、詞庫或手冊答案。男女各一趟固定輔助輸入，各13筆事件及7個停止點，退出0。與R935的indexed、R937英文控制的palette完整bytes及Step／1MiB Mem FNV／A000 FNV／Pal FNV／Mouse相比，14對全同。這不包含完整CPU暫存器或新正常GUI。

四組原型為男女×兩字型，640×400；原生15／16字模保持。31個欄位含17標籤、性別、11數值及生命／體力兩個數值列。屬性改18像素行距，性別字格擴至16像素；數值直接搬原版差異墨點，仍以2倍整數放大，不重算或重新編碼。只還原原版前後差異的像素到精確before，不用固定底色補底。四份現行對照圖完整RGB與R937既有圖相同；四份重排圖由另一個檢查器逐像素重建期望。八份PNG共2,048,000像素相符；四張重排圖各31欄位交疊及墨點交集皆0，四張現行對照圖仍保留既有交疊。主代理實際檢視四份重排圖。

字模沿用R937已獨立原字型核過的資源；本輪沒有重讀第三方原字型，不能稱新原字型oracle。原版248檔與38份列出的cmd/wiz7 Go來源SHA保持，四個Scratch的PCFILE.DBS完整bytes與原版同；這不證SDK完整來源閉集合。其他職業、長數值、更新／撤除與正式HD／Draw尚未驗，正式排版仍待定版與新契約READY。

V1男女均在探針自身13筆計數檢查退出2，14份indexed仍與R935相同。V1把font0篩選套到所有NUL入口，V2移除這項整體篩選後完成；十一個實際數值仍全部font0。故障原因為生命／體力結束訊號被篩掉的強推論，結束訊號的確切font未記錄，不猜補。V2另在計數檢查前保存觀察；兩版來源、binary與初輪輸出均保留。獨立檢查器首跑把正式來源的純SHA欄當成物件，修正欄位讀法後同容器／同輸入重跑通過，屬驗證器問題。

| 私人來源／收據 | SHA-256 |
|---|---|
| `profile_probe.go` | `fa657dbb9a9749e6e2170472b96d45972eafe04ec80c3d3a70cd0da58d4da60d` |
| `profile-probe-v2` | `af88a4835a2afc445d5592cb50ddf66a5f6ad50c1d5be086765d9643ad29af88` |
| `make_layout_prototype.py` | `54d612dab5f50c8e1acb9e15d62a70eea1e3a9740fcffe9a3407c13ff7a2a45f` |
| `prototype/result.json` | `834eaf6be5fccb917094a211ccd0a1bbc8cb703e89c4378d408b99093da37efa` |
| `prototype/compare.html` | `04b82d023fe59cfdcab94cee3346e0faf3ef7d4879e486bac5f9b50b071e5356` |
| `check_layout.py` | `b3391052c112e5d4a15011b38d068747900f401f92d3c95e0aaabce232a8c986` |
| `independent-layout-check.json` | `0fc80ccfd1665a2d63ee7486ee6dfd08ba92589e02ee63dc46fac89f354a6449` |

收據根為`workplace/r953`；[比較頁](../../workplace/r953/prototype/compare.html)只留本機，不提交圖片或原版。重跑沿上述映像，repo及原版唯讀、新且已存在的輸出目錄可寫；在容器內為新輸出建立可寫`gocache`，提供本輪來源及`before.json`，既有離線模組快取掛`/src/workplace/gomodcache:ro`。`go.mod`、`go.sum`與既有SDK隨repo掛`/src:ro`；SDK受控準備入口為[prepare-dosgolem.sh](../../tools/docker/prepare-dosgolem.sh)，不能換成裸上游提交或主機runtime。以`build_probe.py`建置探針，再用`profile-probe-v2 /src/workplace/original/wizar7 /out/male-v2 male`及female對應參數。需同層提供R935／R937／R943已核本機輸入，依序執行`make_layout_prototype.py`、`check_layout.py`；完整性檢查另需保留V1的male／female Scratch收據，缺項不可充作重跑通過。既有輸出拒絕覆寫。所有命令在離線、非root、有界Docker內執行，不用主機工具或注入存檔。

獨立窄審阻擋0／應改0／建議0，主代理於FINAL後完整回讀JSON與Markdown採認，只涵蓋上述唯讀測量及離線原型，不批准正式2倍實作。報告`workplace/r953/review-layout.json`／`.md` SHA依次為`b68546549aa096143b6d62ca6be6d868ed4898eeae4c48adfb94263ca5ee74f1`／`f31e51b03c0a05928f0d8c60c6cbe93945762c5da0f502c201b3b8425518851b`。已補可寫Go快取、唯讀模組與SDK入口，並收窄零交疊只適用四張重排圖。024既有READY範圍保持，2倍候選仍待選；輸出1000:1000、自有審查容器0。

## R952 人物頁圖片操作與2倍行距

入口：[R952](../goal/R952-character-review-pages.md)，2026-10-04。從目前正式來源乾淨重建Linux前端，Go1.24.13 linux/amd64，binary SHA `7363fc5be7c48d10eb272d5bd5900045996bb490b0cb08fd19bb1a93e9d948b0`；正式程式與SDK未修改。使用既有六份Unifont資源、3倍畫布、靜音與手動加速，原版唯讀。這不驗預設節拍、裝置音訊或人耳。

真正GUI由主代理逐張看實際Draw後送10次點擊，共11個停止畫格。兩次從主選單／角色選單進THESUS檢視；在已見人物頁點231,143紅色劍形圖示及207,136蠟燭圖示，下一停止畫格都回角色選單，沒有切出功能頁。此結果為confirmed的兩個點擊觀察，不證明所有圖示皆靜態或任何遊戲內分頁完成。最後經中文「離開」返回主選單，正常關窗exit0，耗時0.106417秒。原版248檔盲SHA保持；各畫格擷取後讀到的Scratch內容與原版相同，讀檔時點不冒稱PNG的固定Step。

正常收據為`workplace/r952/normal-unifont-3-attempt3/captures.json`。首輪觀察器誤等成功碼`err=none`，正式前端實際成功碼為空字串；點擊已成功，失敗在觀察確認。第二輪因主代理未在120秒內回覆畫格而超時。兩份失敗收據保留，屬工具／驗證等待問題，不記為產品缺陷；第三輪沿同一binary、容器與120秒截止值乾淨重跑。

另目視R937既有2倍人物Draw，獨立解GNU Unifont原始字模，核八個屬性兩倍率共32字模。STR至KAR的原版y為35、41、47、53、59、65、71、77；字高16，2倍行距12。七對相鄰字格各交疊4列，實際筆畫交集合計58像素。這是confirmed的R937已驗2倍排版缺陷，未冒稱本輪重跑2倍GUI。R952當輪3倍Shown的同八列行距18，字格及筆畫交集皆0。2倍修正需要同時處理文字與數值行的對齊，不能只補性別或縮字。

私人獨立檢查入口為`workplace/r952/check_gui_and_spacing.py`，在既有`wizardry7-go-ebiten:latest`離線非root Docker中執行；儲存庫、原版及字型掛`/src:ro`，僅`workplace/r952`掛`/out:rw`，執行`python3 /out/check_gui_and_spacing.py`，報告以exclusive方式寫`gui-and-spacing-check.json`。它核當輪binary／來源、10份父代理畫格SHA、22份HD／Draw檔案SHA、原版248檔及存檔內容、字型與字格交集；不重跑完整像素oracle或固定輸入A/B。024仍READY，2倍性別與其他排版、遊戲內分頁、十一組、TSV、正式錄放、冷存讀、音訊及封包均未完成。

正常收據、檢查報告及檢查器SHA依次為`fa6bea32eaad3ef9083896defa613b735a674c1cdae382c8300a1a1904810193`、`7ef08992a61362e8bac66cc89b877fd795a712d48d89cf1a6b2e4545ad3cb858`、`4af4748c577b7d106fc498328c315e1502ed19a0a5b3ae700991ca2f67cf1fc4`。本機140份技術檔一次封存，排除Go快取及後續審查／GitHub回讀；`closed-inputs.json` SHA為`bc2c3a186ba4022177f6178d924236a289218464dfa91be1bc4ad7eb52e8935d`。

獨立結果審查阻擋0／應改0／建議0，主代理於FINAL後完整回讀兩報告及SHA採認。`workplace/r952/review-results.json`／`.md` SHA依次為`61f1c439450fd84d114e2bdf68866d7233cab1478d5acedb40cd2d610eb6bd82`／`3e6b74ad507c712626bca2461121f1383be63bb7159ac54f820b102c699ded42`；另獨立核22張PNG的CRC、解壓及尺寸。只核38份列出的cmd來源，不宣稱完整編譯閉集合或SDK1,132檔本輪獨立重核。擁有權1000:1000，自有容器0。

## R937 正式修正驗收入口

[check_profile_sex.py](../../tools/re/check_profile_sex.py)核對024限定範圍的正常GUI、原版狀態、兩字型原生性別字模、選人按鈕邊界及修前／修後HD與Draw。入口參數為`<候選工作根> --baseline-root <R935工作根> --control-root <英文控制工作根> --unifont-hex <字型> --eten-font <字型> --eten-symbol-font <符號字型> --output <新報告>`；只在既有離線非root Docker執行，輸入唯讀，只開放新報告可寫。省略control-root時使用候選工作根；指定控制根時分別記錄兩個binary雜湊，不能聲稱英文控制由新binary重跑。既有輸出拒絕。限定像素驗收及R941結果複核已通過，未將024升為CONFORMED。

R937動態勘誤：16份男女／兩字型／兩倍率的選人停止點都有標題、確定及取消，但父角色選單六個標籤皆已移除。2倍Shown為3筆，3倍加六筆角色碼為9筆。原先規格要求的九訊號不能同時成立，正式helper返回舊畫筆；2倍16張修前／修後HD及Draw逐位元相同，置中未生效。私人selector-guard-evidence.json SHA `85e1cde8acd6e655f2b9002b292608c4f94780cde0ae17e2cf91e14aae0f1238`。024選人條件回DRAFT後，R940窄審0／0／0由主代理於FINAL後完整回讀採認，再轉READY。正式helper及測試只核三個精確選人訊號，不改性別、原版輸入或資料。修正後驗收另存workplace/r937/repair，初輪失敗收據保持。

驗收工具先後修正舊／新字格遮罩與逐張預期尺寸。首跑另外把短角色選單路線的80筆缺譯數套到人物檢視路線而拒絕；現行修前／修後2倍各121筆、3倍各109筆，逐項相同，不算新漏譯或完整覆蓋。修正計數後，局部像素在倚天2倍確定按鈕198,334不合，才揭露上述正式辨識缺陷。未產成功validation.json，不把失敗收據採認為完成。

修正後正式binary SHA `28f8ea051f59a8399a049c780d10f3457568201da9da960b6c4986040a300c01`，完整cmd/wiz7測試通過。男女兩次真正GUI各由主代理逐頁看Draw再送六次正常點擊，男／女中文、選人及圖片箭頭返回正常，均退出0。真正GUI不當固定步數原版同狀態；關閉耗時0.067154／0.066959秒，持久角色內容保持。再取雙字型、兩倍率、兩性別八條固定輔助路線，每條七點、均正常退出0。舊英文控制binary為`ecd047c8ca61afc683dc4cc1343738780d21b3c36ad41ca6872266f9a49f7e33`，沒有重跑新binary純英文控制。新中文對舊英文及修前中文共112對Step／Mem／A000／Pal／Mouse相同，indexed及palette bytes全同；男女原版窄探針另28對相同。

獨立原字型核68個字模、275,904個性別與按鈕HD／Draw像素，144個性別原版before／after欄位像素通過。112張修前／修後HD及Draw核46,570,368個批准區外像素，差異0；2倍選人字格與筆畫在按鈕內，3倍選人及其他未批准畫面保持。主代理另目視兩字型兩倍率的四張選人Draw。一般八資源及兩個新增性別資源保持烘製版本。

本輪檢查器曾把未印字的背景也要求等於印字後快照。女性3倍656M的THESUS底列各30像素Before＝After＝195，後續原版選取框已改為1，原筆畫失配0，與R928已量選取底邊一致。只在Before≠After處核After並還原Before，保留當前背景。R941檢查器窄複審0／0／0於FINAL後完整回讀採認；JSON／Markdown SHA `1330092346f9021f22255b5c95edbc044e5be571060b700017520c587e22b63a`／`23b02ca475e8d154959edfec48ccb1aadda3899b003d3f39863a0dc6b79c55ef`。正式Go與binary不改，舊失敗收據保持，沒有放寬差異遮罩。

成功報告workplace/r937/repair/validation.json SHA `99e99a3d9c756c9aaae84fc24a71a8867b0fb97a80dadf76307050afb74ce056`；公開檢查器SHA `ba12f128f3cb010cc19ea765ccd3c2dd8de2fcc6a2d7e13d5ea002ea59061012`。拒絕覆寫exit2且既有bytes保持。R941結果最終複核0／0／0於FINAL後完整回讀採認，不升024為CONFORMED。此路線其他無譯文事件2倍121、3倍109逐項保持；2倍性別與角色碼、其他人物分頁、圖片呈獻字樣及全遊戲完整覆蓋未完成。十一組、TSV遷移、錄放、冷存讀、音訊與封包門檻保持。

R941最終結果JSON／Markdown SHA `30c83b9caf5e21dadc6e555ce757eae6db8049001d865f374de46c447c905524`／`e9baef2bfd3df353e46df25585b43d05d55cfadcf4e175c898e783538d3f1e6d`，位於workplace/r937/repair。完整性報告integrity.json SHA `d5885bf87c1a4fb4fe2a17053bcc481eb475ce2e4ab1ebc039c06f8f9e8c22bb`：原版248、SDK1,132、R928／R932／R935舊464／335／503及初輪1,006產物保持，十次退出後持久內容保持、擁有權異常0。工作根find自檢異常0。最終複核獨立重核112對摘要及indexed／palette，另核32,444個原筆畫像素皆符合After；未重跑完整像素oracle。


R937修正後535個技術產物一次封存，workplace/r937/repair/closed-inputs.json SHA `86b6a6ef0fb95bd175b2af973127b854b2b42d3ef42d9855475f90d1b10e901b`。正式修正`450803c`已推送，#22／#28留言[5976034164](https://github.com/wicanr2/wizardry7-cht/issues/22#issuecomment-5976034164)／[5976034765](https://github.com/wicanr2/wizardry7-cht/issues/28#issuecomment-5976034765)與本機正文逐字相同、兩議題OPEN。github-readback.json SHA `2983ceec553858a763d1be24dbe968e067cc58390a6e19898ba80c23d5de483e`。原版248、SDK1,132及535／1,006封存收尾再核保持，自有R937／R941容器0。

## R935 人物檢視與返回

正式R932 binary SHA `9796d372ce53cef7528f9ad5d3a8fc57206af15cfa411cb37ab3f2499a20ed22`及兩後端各四資源保持。真正GUI冷啟動依次點160,47、85,193、175,152、100,170、314,190、287,193；每次先看當前畫格再輸入。從主選單進角色選單、選THESUS、檢視人物頁、右下藍色箭頭返回角色選單、離開回主選單，正常退出0，關閉耗時0.068531934秒，PCFILE.DBS持久內容保持。箭頭沒有英文文字，只證實314,190能返回；20×16色號外接框混有其他像素，不能當完整命中矩形。

首次真正GUI在人物頁看圖超過父代理120秒等待，已保留失敗片段；超時後寫入的返回命令未送進遊戲。主代理承認工具等待錯誤，重查GUI路由後，以相同控制器、截止值及計畫輸入乾淨重跑，不調時基或延長等待。成功收據只採workplace/r935/retry/primary/captures.json，SHA `86b657bdcb38726905ba6eb6d86889265c8e08a7dcb8705f40535891fd3f9c1e`。

其後六組英文／倚天／Unifont、2倍／3倍固定GUI各取272M、400M、528M、656M、784M、912M、1040M七點，全部正常退出0。28對Step、Mem、A000、Pal、Mouse及indexed／palette bytes相同，保存內容保持。3倍人物頁17個標籤為人類、戰士、等級、階級、無、經驗、擊殺、生命、體力及八個屬性；雙字型共59,212項局部HD／Draw像素核對通過。196個原生字模核對通過，3,180個原版after像素與英文indexed相符。一般標籤before來自既有攔截收據，不冒稱另一套原版探針獨立讀得。

選人對話框「檢視誰？」及確定／取消字模正確，但2倍倚天兩按鈕字格各超1列，Unifont各超2列、筆畫分別超17／6像素；3倍全0。人物性別M仍英文，不把17個標籤通過寫成全人物頁完成。其他分頁與人物頁2倍完整排版尚未量。

追加女性正常GUI只將選人點改為175,158，檢視既有TEMPEST；主代理逐次看圖，確認女性蜥蜴人頁仍顯示F、箭頭返回角色選單再離開回主選單。正常退出0、0.071909746秒、持久內容保持。返回命令守衛依本次畫格的LIZARDMAN102,6核對，沒有套用男性HUMAN訊號。收據female-primary/captures.json SHA `c1d4970f43abc211dfa2ac884402d3c267f6ca23d443fec545626f2055aba7a2`。未取得女性英文固定GUI控制，不聲稱女性同狀態或中文性別完成；先前驗證report的female未完成限制對應追加前範圍。

### 單字元性別的窄範圍證據

原版runtime0110:3877的near返回0110:6c42，StackWord(1)為M或F、StackWord(2)為0；入口游標90,6、SP65316，返回SP＋2、游標96,6。一般字串入口0110:3895沒有收到目標M；本證據補充既有013字串管線，不能再把所有玩家文字都等同NUL字串入口。

獨立dosgolem冷啟動、相同正常操作，M選THESUS、F選TEMPEST。只有664M至784M、上述地址／字碼／座標組合被觀察。全64000 bytes離屏差異分類：M有13個、F有11個改動，都在90,6,6,6內、改動後為色號1；改動前含0、14、115–118，底圖有紋理。只清原筆畫並還原精確before，不以純色擦底。這些原始36-byte樣本、原版畫格及完整收據只留本機；未讀訊息檔、鄰文或手冊答案。

M的兩套唯讀探針各七點與英文控制摘要及indexed相同，共14對。F沒有獨立英文GUI控制，不列為正式中文或同狀態GUI完成。私人重跑來源workplace/r935/sex_probe_full.go，SHA `7ec6fa88506bb1011a335cd04e3bf8f24f8b0ca4d05cd66ca6a23ec49e82f10e`；binary SHA `217a153fc302a3314ff0835a6b99ff1b026f37c39c969c30d5646deff4543cae`，入口`sex-probe-full <原版唯讀目錄> <新輸出目錄> male|female`。在既有離線非root image內建置／執行，SDK及原版唯讀，輸出限定workplace/r935新子目錄；不在主機執行。M／F observation.json SHA依序`9491b606f489ca13717fc41a7cac3fe23c339bd39bb8d5a347f7ad5b4bec852e`／`7ec495a48718248598291db66d21b17e7bc7c2c813a39a643c677dd576ede771`。

### 獨立驗證入口

公開[check_character_review.py](../../tools/re/check_character_review.py) SHA `11ba230178b4a0ece4f5eb7b34e281bd73f6e7f577c03df85c98fe519f835227`。工具直接解原字型、原始色號及PNG，不匯入正式renderer。二十個精確標籤只沿正式已知文案核對，沒有讀取MSG.DBS或messages.json。使用既有離線非root Docker，原版及字型唯讀，只有新報告位置可寫：

```text
python3 tools/re/check_character_review.py workplace/r935 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --eten-symbol-font /eten/ET353S/FILES/SPCFONT.15 \
  --output /reports/new-character-review.json
```

倚天全形問號來自SPCFONT.15，12240 bytes、SHA `f32049ba2a7a21db908878a488a2c1d93c389d17398cf46db1390ba89e247605`。STDFONT及Unifont雜湊沿R932，工具內固定核對。首次工具錯把問號從漢字字庫負偏移讀取，核符號字庫契約後修正工具，以同一GUI收據重跑通過，正式字型及文案未改。

私人與公開report bytes相同，SHA `251f530d861fcd9c890612db7347ffe766f71cac68f90a6587a587c6c084641e`；既有輸出拒絕exit2且bytes保持。integrity.json SHA `897234c312c04206708ea4c8d794f1cc8b16391372cce1dc80e85df08dff5793`，核原版248、SDK1132、R928舊464封存、R932舊335封存全部保持，R935擁有權異常0。全工作根find自檢異常0，自有R935容器0；初次Python自檢跟隨既有斷連結而失敗，改用規定find完成，沒有產品重跑。新正式修正見[024](../spec/024-character-review-output.md)，仍DRAFT，不把本次證據計為性別中文或全遊戲完成。

R936初審2阻擋／1應改／0建議，主代理於FINAL後完整回讀兩報告。R935限定證據可採，024資源型別入口及pending仲裁需補明後再審；JSON／Markdown SHA依序`2e3aa2c65bf3bb1ae32395c699582bd98632c49b81286c8b6df54b15c010470b`／`c9dc40d3cd3ac944d3551f470d88f101d7c2cdf665500430604bad077fee5fea`，入口workplace/r935/review-r936.json及.md。新正式工作目標見[R937](../goal/R937-character-review-fix.md)，READY前不實作。

R936窄複審0／0／0，B1／B2及字形／女性現況文字皆解除，主代理於FINAL後完整回讀採認，024轉READY。兩報告SHA見024；只准限定實作，性別中文及按鈕修正尚未正式驗收，完整交付門檻保持。前段DRAFT描述為初審前狀態。

R935的503個本機產物一次封存，closed-inputs.json SHA `bf8b59aebd98026d3d0375c29ffe9674e32f25aa0748f89a021912f1a2bd7908`。完整原始收據只留本機，正式修正另在workplace/r937，不覆寫基準。

## 驗收欄位

每筆記錄畫面與控制、文字輸出或圖片字樣的來源、原版定位及座標空間、辨識訊號、正式中文畫格、點擊與回復結果、同狀態收據、字形及邊界、證據等級與下一個缺口。

正式程式已完成、動態畫面未量到與可選美化分開記錄。無文字圖示沒有可翻譯字樣，但其說明與操作仍須核對。玩家自行輸入的名字與既有原版英文關鍵字規則保持，不把它們誤計為未譯按鈕。

## 本輪取樣

2026-10-04先由修前正式來源重建Linux前端，SHA-256為5fe5f410065ecf14303e896e72d00eccdf93a70f11ff30fbdafd050eef560bd5。原版DS.EXE、SDK與字型輸入唯讀；完整來源摘要在workplace/r915/inputs.json，產物摘要在build-artifacts.json。主選單修正後正式binary為7c2553a39f4fef55efb49c7c7df4eae37b4eccac5761371079a14eedae8eb18d，最後來源重建與GUI所測binary逐位元相同。

兩套後端都由現行來源生成：一般覆繪資源8,425個條目、2,867個中文字形，流動文字2,759個字形，說明14頁、426個字形。字元聯集2,919個，非ASCII字元2,872個。資源生成成功只證明字模可生成，不證明每個按鈕上屏。

第一段以正式GUI跑倚天、Unifont及英文版，各2×／3×，無玩家輸入；每次取2M與60M步，12張HD與12張實際Draw均完整解析且尺寸正確。六次WM_DELETE_WINDOW正常退出碼0，關閉請求至退出0.019至0.183秒。這是靜音視覺取樣，不等同聲音或完整可玩驗收。

相同步數八組中文／英文對照的Step、原版記憶體摘要、indexed與palette摘要、滑鼠狀態均相同，.a000及.pal逐位元組相同。comparison.json SHA-256為cf3fc541a08506040c51b6a70feb67870209632ef21270f978811d0bebd475b0。Draw是實際視窗幀，沒有凍結在標稱步數，不拿它逐像素宣稱固定步數同狀態。

## 目前覆蓋矩陣

| 畫面／控制 | 來源與辨識 | 中文與操作結果 | 狀態／下一步 |
|---|---|---|---|
| 冷啟動2M，廠商標誌 | 舊正式Draw已見圖像；R942正常唯讀路線量到MON60.PIC載入及首張圖片轉送 | 無操作按鈕；品牌字樣仍是原版 | 載入與轉送已定位，codec與品牌分類未驗 |
| 冷啟動60M，廠商標誌及呈獻字樣 | R942定位MON60最後圖；R943依使用者選定的周邊紋理補底接入六張呈獻覆繪 | 雙字型兩倍率、32點原版對照及六趟正常冷GUI通過；品牌依既有規則保留 | 限定呈獻通過，翻頁傾斜字樣與介紹仍沿R944待選，不外推其他圖片 |
| 片頭後段翻頁敘述 | R942的MON61索引組與六張完整轉送已量，圖片內另有英文角色扮演敘述 | 敘述尚無中文；作者與品牌另依既有名稱規則核分類 | 僅原版無輸入觀察；不稱正式中文GUI通過 |
| 主選單六個標籤 | 正常無輸入272M，六個resident文字覆繪事件；不是圖片內固定字樣 | 六標籤完整；「匯入巫術六」五字；雙字型2×／3×的HD及實際Draw都在按鈕內 | 角色選單與新增角色入口／返回通過，其餘入口未量；防禦分支僅單測，見[020](../spec/020-main-menu-labels.md) |
| 片頭製作名單的未翻譯事件 | 四組中文run.log各80次；75次精確對應15個去低位碼後的正式英文鍵，5次空字串 | 15筆來源的zh與en相同；尚未逐項確認保留內容的分類及實際畫面 | 不當75個新漏譯，也不計全文覆蓋通過 |
| 新增角色對話框標題及確定／取消 | R921、R925正常GUI及固定步數已量；既有resident文字覆繪 | 雙字型兩倍率字格與筆畫在框內；正常取消返回 | 限定通過，[021](../spec/021-character-dialog-labels.md)仍READY |
| 新增角色清單六列 | 完整名字＋性別職業動態文字，返回位置0110:0fcc | 名字依既定規則保留英文；六個性別職業3倍雙字型完整，選取及取消正常 | 六個實例通過；2倍、任意姓名與捲動未完成，見[022](../spec/022-character-roster-text.md) |
| 角色選單六按鈕 | R932正常GUI及resident0110:5686、六個精確原文／座標 | 創建／檢視／刪除／改名／頭像／離開完整；雙字型兩倍率字格及筆畫在框內，離開返回正常 | 六按鈕限定通過；人物檢視及頭像頁另見024與R955，其餘子畫面仍待完整驗收，見[023](../spec/023-character-menu-labels.md) |
| 頭像頁標題及完成控制 | R955正常GUI到達SELECT PORTRAIT於116,138、CLICK HERE WHEN DONE於101,153 | 「選擇頭像」及「完成後點這裡」雙字型兩倍率字格完整、局部HD／Draw像素相符；保持頭像完成返回，角色檔相同 | 只驗兩則文字與保持原頭像返回；選新頭像、其他圖片控制及已知2倍人物／名單缺口不計通過 |
| 人物檢視頁 | R935／R937正常GUI、窄印字證據與固定冷啟動 | 3倍雙字型既有標籤及男／女性別通過，選人確定／取消兩倍率在框內，箭頭返回正常 | [024](../spec/024-character-review-output.md)READY；限定結果複核通過，2倍性別與其他分頁未完成 |
| 人物頁右下箭頭 | 無英文圖示；R935實際點314,190 | 正常返回角色選單；只證實此點，未量完整命中矩形 | 圖示無待譯字樣，不外推其他圖片按鈕 |
| 清單上下箭頭 | 已見無英文圖示，操作及說明未量 | 尚未完成控制驗收 | 另沿正常GUI核對 |
| 戰鬥七控制 | R958合法暖入口、R959精確英文／座標及原版39×9外框；既有resident文字事件 | R960固定HD與R961固定Draw雙字型兩倍率保持；R962真正Unifont2倍GUI完成裝備、最後指令、返回，七控制撤除及重現 | 只量裝備與返回這條路線；其他控制操作、後端／倍率及新GUI A/B仍待，見[026](../spec/026-battle-control-labels.md)及[R962](#r962-normal-battle-gui) |
| 建立角色、其餘人物／物品／咒語／對話／存讀及其他圖片按鈕 | 此次重查尚未到達；既有規格與歷史收據須逐項核現行來源 | 未量到，不計本次通過 | 依可見畫面走合法操作 |

首輪建置的Go編譯並行數超過容器程序上限；限制GOMAXPROCS=2及-p=2後，同一工具鏈重跑成功。原失敗日誌保留為build.log，成功日誌為build-p2.log。這是工具環境修正，不登錄產品缺陷。

第二段Unifont與英文3×無輸入至200M、272M，兩次正常退出0，兩組原版摘要與色號／色盤相同。200M是標題與著作權畫面，272M自行出現主選單，不需要猜補按鍵。匯入原文IMPORT (BANE)的兩份正式來源1007／1415都譯為「匯入(巫6)」；bake_font.py:165–167略掉非雙位元組字元，因此只生成三字。這是confirmed漏字，不是圖片按鈕或缺字型。

六個按鈕邊界由英文.a000色號176／192／197逐行、逐欄獨立量測，結果在menu-bounds.json：左105、寬109、高9，頂點依序43、56、69、82、95、108。原文字y為45、58、71、84、97、110；修前2×原生15／16列字格底部各超出1／2像素，實際Draw亦已核對。修後倚天2×上移2像素、Unifont2×上移3像素；3×位置保持。字形不縮小，原版按鈕與點擊範圍不改。

量測腳本初次誤從英文模式讀Shown事件，因英文模式沒有中文事件而失敗，未寫結果。修正為已對拍中文事件提供定位、英文原版色號提供獨立邊界，乾淨重跑成功，不作產品缺陷。

## 正式修正及獨立驗收

正式catalog只改1007、1415兩個zh為「匯入巫術六」，en、ID及其他欄位保持。兩字型xlate各只改匯入鍵；flow各只有同兩筆行／訊息由無可用字形變為五字，增一個「匯」字模，其餘字模筆畫與欄位保持。help逐位元相同，字元聯集不變。完整差異收據在menu-fix/resource-diff.json與pixel-validation.json。

[menu_labels.go](../../cmd/wiz7/menu_labels.go)只在六個精確英文鍵與座標同時存在、dy為0、字高一致15／16且字寬足夠時，調整HD及Draw畫筆的目的y。原事件、源字模、原版記憶體、資料比較、存檔及滑鼠操作均不修改。整個cmd/wiz7測試通過；排除情境未有正常GUI收據，因此020維持READY。

六次新正式GUI及四次修前binary搭配同一新資源的GUI，全部於272M取完整HD及實際Draw，正常關窗退出0。這段採手動-ips 24000000及靜音，兩側CPU／PIT及虛擬時鐘相同，只作視覺回歸，不作預設節奏或音訊驗收。四組中英Step、Mem、A000、Pal、Mouse與.a000/.pal相同；四組修前／修後原版摘要及.comp相同。

獨立解析倚天stdfont.15及固定Unifont17.0.05 HEX，100項字模筆畫與306,072項按鈕像素核對全部通過。HD與實際Draw的六個按鈕均核整個矩形，匯入五字完整。2×修前／修後變更倚天1,940像素、Unifont2,258像素，批准舊／新字格以外差異0；3×同資源位置差異0。主代理已目視四組修後實際Draw。這些數字只證明六標籤的顯示，不證明按鈕操作或全部中文化。

可重跑的獨立檢查器為[check_menu_labels.py](../../tools/re/check_menu_labels.py)。在離線非root Docker內使用下列參數；驗收輸入唯讀，只有新報告所在目錄可寫。既有--output會拒絕覆寫。

```text
python3 tools/re/check_menu_labels.py workplace/r915/menu-fix \
  --bounds workplace/r915/menu-bounds.json \
  --baseline-bake-root workplace/r915 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-menu-validation.json
```

Unifont HEX的SHA-256為7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574。原版畫格、倚天資料及完整收據留本機，不提交。公開檢查器已以新的報告目的地重跑，結果同為306,172項通過。

像素檢查器首跑誤將按鈕斜角色號當成文字；修正為由原版按鈕內部區域獨立辨識前景，乾淨重跑通過。這是驗收工具問題，未修改正式顯示來迎合錯誤oracle。公開檢查器初版的固定報告目的地曾覆寫相同內容，現已改為明示新路徑及拒絕既有檔案。

## 完整性界線

四組中文run.log的80次事件完全一致：75次、15個非空英文鍵與現行catalog去低位碼後相符，且各筆zh等於en；5次為空字串。出處為2901、2902、2904、2906、2908、2909、2911、2913、2915、2918、2920、2921、2922、2923、2924。沒有讀取防拷或詞庫正文。完整metadata收據為workplace/r915/untranslated-events.json，SHA-256為5b6d5ee93cb99b9b2ac3a57bd5c48acde57f3ab51de9c0393d41c9f6e2b02f37。來源相同不代表每項保留英文都適當，後續仍須核人名、平台名、職務標題及畫面來源。

靜態ASCII略過候選由1,604降為1,602；其中Latin115、數字134、標點1,500可重疊。標點包含合法控制碼，候選不等於確定漏字。沒有逐筆查證前不批次改文案；其他畫面、圖片按鈕與十一組交付門檻保持。

## R921 新增角色對話框

沿正常主選單，以實際GUI左鍵160,62開啟新增角色對話框。400M畫格已見標題「新增角色」及「確定／取消」。取消按鈕的原版色號邊界為63,182,73,9；主代理看圖及核邊界後才單擊100,188。528M返回選單，標題及兩個對話按鈕移除，WM_DELETE正常退出0，退出耗時0.067秒。這只驗新增對話框及取消，不驗選取、建立角色或完整遊玩。

原版既有PCFILE.DBS寫開時可在獨立存檔層生成複本。第三次正常GUI收據列9,368 bytes、SHA-256為24d94441e5b32e7abbcbffa9dad972b41544da6a7bb06096b967281f639d9c1b，與唯讀原版全同；沒有變更持久角色資料內容，不聲稱零DOS寫入。控制器舊stdout「saves empty」不準確，已在WORKLOG勘誤，原收據保留。前兩次因日誌格式及空目錄斷言失敗，未正常關窗，不計成功。

正常GUI成功後，倚天／Unifont各2×、3×及英文兩倍率固定280M新增、440M取消，於272M、400M、528M取得完整HD及實際Draw，六次均正常退出0。四組中英文、三個停止點共12對的Step／Mem／A000／Pal／Mouse及.a000／.pal全同，退出後PCFILE內容全同。檔案取樣記控制器實際讀取時間，不能把稍後讀目錄冒充畫格Step的精確檔案快照。兩次早期輔助空白斷言失敗仍保留，後續沒有調整步數、輸入或時基。

獨立解原始兩種字型，32個標籤字模及93,950個HD／Draw像素核對通過，主代理也檢視四組中文Draw。原版確定按鈕63,166,73,9，取消63,182,73,9。2×倚天字格底部超出1像素，但本次兩詞的實際筆畫未超出；2×Unifont字格超出2像素，確定17、取消6個筆畫像素在按鈕外。3×兩字型均無超出。這是confirmed排版缺陷，不能把像素正確呈現舊排版算成按鈕內排版通過；限定修正見[021](../spec/021-character-dialog-labels.md)。

角色清單有六個完整動態字串沒有匹配，來源是resident印字返回位置0110:0fcc，遊戲座標x150、y150／156／162／168／174／180。[013](013-print-pipeline.md)已列「名字＋職業碼」家族，因此不是圖片字樣。預設名字依[names](../names.md)保留英文，性別職業仍需中文。四組缺譯事件同為86次，較主選單基線多此六筆；其他80筆分類界線保持，不宣稱漏譯0。清單上下箭頭為無英文圖示，尚未操作或完成說明核對；片頭圖片字樣與其他畫面仍未完成。

R921正式binary為7c2553a39f4fef55efb49c7c7df4eae37b4eccac5761371079a14eedae8eb18d，本輪正式Go／catalog／生成器沒有變更。原版248、SDK1,132與來源基線保持。私人權威收據：workplace/r921/primary-r3/captures.json SHA為3fda0d582631d9c80be2f21462263ef196709e538234da58df81ef26a727100d；dialog-validation.json SHA為961d769f8bc68035ecc5933e60b454254b951de4c29cf0ca8fba31efa9412d98。R922、R923、R924契約與取樣修訂均經主代理完整回讀後才採用；R926結果複核0／0／0，主代理已採認。十一組、TSV、錄放、音訊、冷存讀及封包門檻保持。

## R925 小按鈕正式修正

[021](../spec/021-character-dialog-labels.md)經R926獨立計畫審查到READY後，增加[dialog_labels.go](../../cmd/wiz7/dialog_labels.go)，HD及Draw共同使用buttonLabelY。只在九個已量畫面訊號、指定兩按鈕、dy0、原生15／16字高及寬度守門成立時置中文字格於9列原版按鈕中間。2×倚天上移2像素、Unifont上移3像素，3×保持；020主選單及標題保持。main.go只換兩個helper呼叫，catalog、生成器及字型資源不改。

cmd/wiz7測試全部通過，正式binary SHA為ce62ac638b89d42de0d10239d08ebbaf2c1b849ce13aa80908f4a655a859cb6b。新binary先以實際GUI完成新增、取消100,188、528M返回及正常退出0，退出0.097秒；主代理逐次看圖後才送操作。PCFILE只有與原版相同的複本，持久內容保持。再跑四組中文固定路線，全部正常退出，與R921英文及修前共12個固定停止點的原版摘要、色號／色盤相同；修前後Shown及.comp也全同。

獨立32字模、92,198個HD／Draw局部像素通過，八個中文按鈕的完整字格與實際筆畫超框均0。固定步數HD全圖比較：2×倚天改324像素、Unifont改298像素，批准舊／新字格聯集外差異0；272M及528M、全部3×HD差異0。Draw只核上述局部字模與按鈕，未作新舊全圖差異比較。主代理檢視四組新Draw。防禦分支僅單測，021維持READY，不外推其他同名按鈕。

可重跑檢查器為[check_character_dialog.py](../../tools/re/check_character_dialog.py)，需現有本機R921基準與兩種合法原字型；報告明示新路徑並以exclusive方式建立。執行一律在離線非root Docker，原版、字型與基準唯讀，只新報告目錄可寫：

```text
python3 tools/re/check_character_dialog.py workplace/r925 \
  --baseline-root workplace/r921 \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-character-dialog.json
```

六列動態性別職業仍缺譯，86次缺譯分類保持。這次修正只解決新增角色兩按鈕的排版，沒有完成圖片字樣、箭頭操作、完整中文化或交付。

R927最終獨立複核0／0／0，主代理完整回讀後採認。私人報告workplace/r925/review/review-r927-final.json／.md的SHA為c8bc00f9ed49b8de9808a3c1c59c910dcd97b6c6b0e7ebe67838476a594ebc81／3c26613057a75d0c4fd91b45b323c1c225fff9389f8af89705cb80edc76ee491。R921的244檔與R925的137檔封存保持，closed-inputs.json SHA依序為4e83af8f928ce462e954e292e29b3f18d112f9ff41b6e25d87263b0bb9736eb4／37c4ffec07794e7f6c865cd9ec97457185657ac524de59d98a770950d52b8f5b；來源、SDK及原版雜湊保持，輸出owner1000:1000，自有容器0。

## R928 角色列的目的區域

沿R921固定路線及同一原版載入設定，以私人唯讀探針在0110:3895限制CS0110、返回0fcc、x150、六個已量y後，只讀16 bytes確認六個精確英文鍵及NUL。0110:3877只在目標呼叫期間讀字元目的座標；返回只接受SP＝入口SP+2。沒有讀MSG、鄰文、手冊答案或詞庫，沒有修改原版記憶體。這是同一正常GUI路線的無頭補充量測，不能當新GUI或角色選取驗收。

confirmed：六列字型均0，各15字元，x依序150、156、162、168、174、180、186、192、198、204、210、216、222、228、234，y在該列保持，返回x240。姓名筆畫只在x150起的前十格；性別職業碼筆畫在x210..239、y..y+5，未選取時差異色號均195→180。碼帶固定x210..240、y..y+6。初次只量90×6區域；後續全64,000像素補量證實六次的已知姓名格、碼帶之外差異0，未知色號0。補量來源SHA `5791f78851bf78f2b58243b9298946a91079c2cd656d95c43708206497fbde15`，whole-offscreen-validation.json SHA `d53dba9969b62966237c8e7d53ebaed6debd056e5d9d74cf7522e91d6ea34093`。

272M及400M原版Step／Mem／A000／Pal／Mouse與R921英文3倍收據全同，indexed逐位元相同。PCFILE只有與原版相同的9,368 bytes複本；原版248與SDK1,132雜湊保持。私人探針來源SHA `04466d99e7b86581256bee67f7b47cd2b7598ae1b6731e9850b26c1d6127981d`、binary SHA `12f0497d7e60aa67657989aeb622ce8ba72236ea258ed34ac50491ad3c12a047`；observation.json SHA `02c2ac5c8bcdba51339e661f43eae502c9c79f25db08a1f7d5c00ac616993c01`、observation-validation.json SHA `c08a9114fe74560b83801b5d13be6392cb07c9888dc62b573bfca7dca72c8b36`、probe-integrity.json SHA `5bcc05374e9608bfb4d7a5070361e79d15ffccad17d26773261a5a618076cbef`。檔案均位於workplace/r928，沒有原版素材加入Git。

原生15／16列字型在3倍的18像素列距內；最長候選「女鍊金術士」寬83／88，在90像素碼欄內。2倍列距12、碼欄60不足，不能縮字或省略職業譯名來掩蓋。私人Canvas原型workplace/r928/roster-layout-prototype.html，SHA `c8a09212cfb172e8f6a2203b0cb124467d23709155294fcc128b89064cda1d8c`，只展示既有Draw底圖的2倍限制與3倍候選，不是正式GUI或玩家操作收據。後續契約為[022](../spec/022-character-roster-text.md)，R930複審後為READY；正式六列中文驗收進行中，原版任意姓名格式、捲動生命週期及2倍完整中文仍未完成。

實際GUI後續完成新增、選取THESUS、選取PENTAS、取消返回及正常退出0，退出0.119秒，未按確定、沒有建立角色，PCFILE內容保持。首輪控制器等待父代理操作超過120秒，failure.json保留；同控制器、同預定操作以新目錄重跑成功，沒有調整遊戲時基或程式。成功收據selection-retry/selection-primary/captures.json SHA `dd4165b53fc0495af2663c1ca6b5568015da425cbe2760a388fb9e679add7355`。

無頭補量三批共18次精確印字，兩個被選取的印字為1→0，另16次為195→180。每次姓名與碼均有差異，其外差異0。實際GUI白底範圍x148..240，首列y149..155、末列y179..185；末列選取時較晚畫的白底會改上一列TREON第六列空白。獨立oracle曾漏算此後續重畫，失敗紀錄保留；先量indexed邊界及原字模，再修正oracle，不改有效收據。12個Unifont原生字模末列均0，3倍中文筆畫不碰此原版底邊。獨立18,818項通過，selection-validation.json SHA `ee67ebc0adc5146840be214b50d31c4b0e407fc04b05c9861af5d36889967c86`；probe-selection binary SHA `aa6b8f755ffb7a18a053af659094ef30e8d3e07abbc2a3c41b079ecf7f7783ce`，觀察SHA `1f0813c724971703fcf1c72cbeeb120a029e1eef5455680179e05ec009921aee`。全部位於workplace/r928。原版248、SDK1,132及持久角色內容保持，root-owned輸出0。無頭與GUI輸入時間不同，未宣稱兩者完整Mem相同；未實作正式六列中文。

## R928 正式六列覆繪

R930契約複審通過後，新增[text/roster.zh-TW.tsv](../../text/roster.zh-TW.tsv)作唯一正式文案；[bake_roster.py](../../tools/re/bake_roster.py)烘製兩種原生字型，既有xlate／flow／help bytes保持。[roster.go](../../cmd/wiz7/roster.go)在一般查表早退前辨識六個完整原文及精確呼叫訊號。專用事件只清除30×6碼帶，完整90×6列只作身分見證，包含已量白底邊界例外；選取前景0合法。未知完整列會撤掉舊事件。譯文不寫入原版或存檔。TSV接線另核路由43的顯示與語意隔離契約，原版流程保持。

正式binary SHA `aca682fdcba57cb56225fd43164961805d53899a9e3847f49356433d70d87d54`。cmd/wiz7測試通過；新TSV與typed資源守門另驗九種格式拒絕及禁止覆寫。真正GUI完成新增、首列與末列選取、取消返回，正常退出0、0.072秒；主代理每步先看實際Draw才送下一個操作。兩種字型的最長五字完整在欄位內，姓名保持。這趟靜音／手動加速GUI只驗顯示及操作，不當完整音訊試玩。

修前／後八組固定路線，各五停止點，共20對Step／Mem／A000／Pal／Mouse摘要相同，原始indexed及色盤bytes全同。HD及Draw各20張修前後全圖比較，批准碼帶之外差異0；2倍及主選單／取消後全圖差異0。3倍三個列表畫面每張倚天改3,185像素、Unifont改3,282像素。獨立原字型核180個字模、輔助116,640及真正GUI58,320個局部像素，共174,960項；中文框與筆畫均在90×18欄位內，取消後六列中文全移除。原版248、SDK1,132、九程序持久角色內容全保持。這是六個已量實例的中文化，不推廣到任意姓名、其他職業或捲動。

3倍兩後端及真正GUI此家族缺譯0，其他缺譯仍80次；2倍此家族仍18次、總98次。2倍不能縮字，目前仍有完整中文閘門，不代表批准交付時留英文。圖片字樣及箭頭操作仍未驗。022防禦分支只有單測，整份維持READY，不標CONFORMED。

獨立檢查器[check_character_roster.py](../../tools/re/check_character_roster.py)需要本輪固定路線與合法原字型，入口見[字型說明](../../font/README.md)。初次60秒期限不足，未產成功報告；只快取PNG解碼、按工作量給300秒後，同收據與判準通過，正式Go及binary未改。私人implementation/roster-validation.json SHA `3bf34ff9ac9875bab759da697680cf676aa2cd2d02d92005e1c60029870acb03`，candidate-integrity.json SHA `5b43b08c4551212eb4a5bc0981348193c0e7b95948b0fd05ca06929b15f9c383`，parser-validation.json SHA `b1ff370437324c214395db3c41927463a3d0660dd1cd317d3d6c532fb42018fc`，位於workplace/r928。

R931實作及結果窄補審均0／0／0，主代理完整回讀後採認。結果JSON／Markdown SHA `34fcad80e19bbee53a474a2bc042ba244b4107989958f8201b64e911d877966f`／`0d300ba1ace41c5583e0a60e3d5f97b3bbe469a7be61b58558102a14f9092557`。不改十一組、完整TSV、冷存讀、音訊／錄放及交付門檻。

## R932 角色選單六按鈕

沿正式R928 binary及Unifont3倍，真正GUI正常點角色選單160,47，400M完整檢視六按鈕，再點離開287,193，528M返回主選單。正常退出0、0.129秒，持久內容保持。首輪父代理再次超過120秒操作期限，失敗收據保留；同控制器／程式／預定操作以新目錄重跑成功，未調遊戲參數。重跑primary/captures.json SHA `ab4d1bd4ff9689b56b5ffdb0e81c4f8e384bb1b4cf4eea728cd2f2fd484a5b20`。

六個標籤為CREATE／REVIEW／DELETE／RENAME／PICTURE／EXIT，現有文案為創建／檢視／刪除／改名／頭像／離開。唯讀正常路線補量確認全部走resident0110:3895，返回0110:5686，字型0，y190；x依序11、63、115、167、216、277。這六個標籤是印字文字，沒有證明其他圖片按鈕已中文化。背景美術及未操作的頭像子畫面不外推。

原版indexed量到六個按鈕皆49×9，左上為x4／56／108／160／212／264、y188。四組中文及兩組純英文正常退出0；三個停止點共12對Step／Mem／A000／Pal／Mouse相同，indexed及色盤bytes全同。獨立原字型核48字模及141,120個HD／Draw局部像素通過，真正GUI400M的indexed／色盤與HD／Draw PNG也與已核Unifont3倍逐位元一致；不宣稱真正GUI與固定輔助路線的Mem相同。

confirmed修前排版缺陷：2倍倚天完整字格超底邊1像素，本次六詞筆畫超框0；2倍Unifont字格超2像素，六詞的筆畫超框依序18／12／9／12／7／16像素。3倍兩後端字格及筆畫超框均0。正確畫出舊排版不能算按鈕內排版通過。限定修正契約見[023](../spec/023-character-menu-labels.md)，R933審查後READY，修後結果見下節。四組中文總缺譯事件仍80，這六詞缺譯0；不能稱全遊戲缺譯0。

私人workplace/r932/menu-validation.json SHA `fb236a9f4d2fba933033d783c47e3a2011310801174b020ea3ca270c7c00fdde`，primary-integrity.json SHA `600dcc3d77bfdcd31710a616a5dab771d36e461d7a7af6d32e287bd000df9145`。探針source／binary／observation SHA `1dac0e28ad6246384836bf2da7b4b69f2d40fa3f4af45ed5786ae650cf78b691`／`b9c1148b639fcacb0f8a954e1f61ac2394d9b16506eef1c450be92109861b10a`／`1541a2bb0a689248f68339056078bcc82d20bb87c7e5250d701ad550445438d3`。三個探針停止點與英文輔助摘要及indexed全同；原版248、SDK1,132保持。這段觀察未改正式Go／catalog／字型，也未進入建立、檢視、刪除、改名或頭像子畫面。

## R932 六按鈕正式置中

R933獨立審查0／0／0、主代理完整回讀後，023轉READY。新增[character_menu_labels.go](../../cmd/wiz7/character_menu_labels.go)，六個精確訊號皆在Shown且字模為一致的原生15／16行才調畫筆y。全家族高度、位元及各自寬度均守門；失敗沿020／021位置。HD及Draw共同走buttonLabelY，不改x、原版矩形、事件生命週期、文案或字型。2倍倚天y378、Unifonty377，3倍皆570保持。

正式binary SHA `9796d372ce53cef7528f9ad5d3a8fc57206af15cfa411cb37ab3f2499a20ed22`，完整cmd/wiz7測試通過。真正GUI正常進角色選單、按離開返回及退出0、0.073秒；持久角色內容保持。主代理每次看圖後才送左鍵，沒有預注入輸入。400M的HD／Draw／indexed／色盤與已獨立核過的Unifont3倍輔助逐位元一致，沒有宣稱兩條路線Mem相同。

四組固定中文程序正常退出0；12個停止點與R932英文及修前原版摘要、indexed／色盤相同，Shown及.comp保持。HD及Draw各12張全圖比較：2倍400M倚天各改892像素、Unifont各改1,108像素，批准舊／新字格聯集之外差異0；主選單、返回及全部3倍全圖差異0。原字型獨立核48字模、141,120個局部像素，24個按鈕／字型／倍率組合的完整字格及筆畫超框全0。主代理檢視四組實際Draw。

原版248、SDK1,132、R928封存464檔及兩後端各四資源保持，五程序正常退出及持久內容保持。四組80次無譯文事件與修前逐列相同，六標籤此路線缺譯0。防禦分支只有單測，023保持READY；其他子畫面、圖片字樣、十一組、TSV、錄放、冷存讀及交付閘門保持。

可重跑獨立檢查器為[check_character_menu.py](../../tools/re/check_character_menu.py)，原版、字型與基準唯讀，新報告目錄可寫；拒絕既有輸出exit2且報告bytes保持已核。一律在既有離線非root Docker，原字型入口見[font/README.md](../../font/README.md)：

```text
python3 tools/re/check_character_menu.py workplace/r932/implementation \
  --baseline-root workplace/r932/retry \
  --unifont-hex workplace/fonts/unifont-17.0.05/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 \
  --output /reports/new-character-menu.json
```

私人workplace/r932/implementation/menu-validation.json SHA `72caceae97226b6acb69d55e8a29b995e1893c71ce298651b0a07cec64072a4b`，integrity.json SHA `ac04632f97c7e45f4bcf840de02c7fce08655498caa061b25c194c320141da31`，primary/captures.json SHA `f32a75ed0f0b3234d24644f5061506a0940c4dbe5a5b15a6386cf747e3af770c`。完整來源SHA與報告拒絕覆寫見同目錄formal-summary.json。

R934最終獨立複核0／0／0，主代理於FINAL後完整回讀JSON／Markdown採認；SHA依序`b60e8a651023924217baf1b80945e2d6db5781fb1afb53134804e6c930804b08`／`079f4426de21d9ac3a850bf363a23dfcd8899677fe11f6a88fced28817ea3963`，本機workplace/r932/review-r934.json及.md。335個產物一次封存，closed-inputs.json SHA `73089cfbdcfe6086a73bb544e11777f1c24311d5ef2aaee9faa9dd19ec852fb5`，原版及SDK保持，工作根root-owned及.md目錄異常0，自有容器0。023仍READY，推送回讀另補。

## 生成來源一致性與審查


R919最終正式實作及生成來源同步補充計畫審查各0／0／0，主代理已完整回讀。報告JSON／Markdown的SHA-256為90a7d415007bfbba287460ce3d506c639d79a5e819d760e5bf5747732ccd416a及f288f1aaaefda67ea0a9f1ed9e7426aa69d17eba235e0166ba622a5d3e6c40fe，兩份owner1000:1000。checker末段現以open('x')獨占建立，新目的地重跑同306,172項，拒絕既有報告退出2且bytes保持；三份像素報告逐位元相同，SHA-256為2dab43c45aed77d681eaf26647dcbb673eaca5f508f5b3de2dd3339f80a1e317。

來源查證曾發現gen_trans.py:Z的1007與batches-merged.json的1415仍保留舊值，會在再生成時覆回漏字文案。依020補充READY同步兩值後，抽象語法樹及既有setdefault合併獨立核有效8,825筆zh與正式catalog全部相同，其他8,823筆保持；批次JSON只另補末尾換行。沒有執行會讀messages.json的產生器，也未讀MSG.DBS或新原文。10個既有GUI產物雜湊保持，不重跑未改的正式顯示。workplace/r915/source-sync.json的SHA-256為4eb14d362fe75ba5848ca28066c4fbc7a62fc99e460c655e8bd767b5eb2c1c80。R920只複核此來源同步及回填，已通過並由主代理採認。

R920最終窄複核0／0／0。JSON／Markdown的SHA-256為7488c8163c589b9e37ab5d99bd8721dd1e51e57ed53df7b40e08ecfb85c78504及55bb357e7978451b861e13002ac066d687cf0a58c3d2a96b91a84ae0efddfd05。286個本機驗收產物已一次封存，workplace/r915/sealed.json／verified.json的SHA-256為f8d6da4f5025b1e85dc6a8d3d99bbe7b02015d3dbd73220643ff6a92b87ae3e4及bed3573cab32eed02ddc554e9a5f94c82fed1d4a34f34e9c98911d2ec719b149。SDK1,132、原版248及四份字型輸入保持，產物owner1000:1000；快取未納入封存。

## R942 片頭圖片來源與切換

正常冷啟動無玩家輸入至200M，原版及SDK唯讀。只觀察開檔身分、繪圖數值、記憶體摘要與畫格，不讀訊息或答案正文，不注入state、種子、作弊或原版寫入。這是原版無視窗觀察，正式中文前端保持。

| 關係 | 實際證據與等級 |
|---|---|
| 圖片選取 | 0110:2663原始堆疊值為[5B57,60,15]及[5CB6,61,15]，正常DOS開啟MON60.PIC於684960、MON61.PIC於48268719。指令及實跑數值confirmed。 |
| 覆蓋程式定位 | 依[016的R881](016-party-position.md)，VINIT code載入0110:5055；IDA code-only EA B02／C61對應兩個返回，OVR offset為EA+14。原OVR的14..5655與既有code逐位元同；不是resident的15B57／15CB6。 |
| 圖片繪製 | VINIT raw9C5／runtime5A1A把索引及座標送resident4018，rawA48／runtime5A9D組合一至三張後送3921全屏轉送。實測參數15一致；尚未解driver內部槽與codec，不給表格格式猜名。 |
| 完整轉送 | 指定返回5B9B／5AFD且SP為進入SP+2，才保存14張完成圖。37次screen-blit及driver呼叫不能當14張全部轉送。第一組8張、第二組6張；未在返回另存SP，matching資格由探針條件保證。 |
| 呈獻字樣 | 第一組第3至8張可見；最後圖於48267654完成，60M仍是同indexed圖。MON61載入沒有立即換上屏。第二組第一張於73235759完成，翻起的紙頁仍有傾斜的presents，R944以既有完整畫格確認；六張平面覆繪不涵蓋此圖。圖像與開檔／服務參數的關聯為強證據，沒有據此宣稱完整PIC解碼。 |
| 其他圖片文字 | 第二組另見角色扮演敘述。後續MON62／63及TIT1SCRN／TIT2SCRN只有正常檔名與服務數值；沒有讀新PIC正文或把它們列為中文化通過。 |

第一張完整轉送時色盤仍黑，2M才見廠商畫面，不把黑圖當漏譯。0..60M字串／單字入口各0次；至200M已各125／2325次，不把後段也稱為零印字。

觀察開／關100對1MiB記憶體、indexed及palette摘要、Mouse、mode相同，300份a000／pal／png逐位元同；沒有比較完整CPU寄存器。前30點同初輪無觀察基線，前8張完成圖同上一版窄探針。獨立審查另核469筆IDA匯出定位、48筆相對E8與287筆VINIT原bytes；聯絡圖CRC、尺寸及14張來源像素相同。主代理在FINAL後完整回讀兩份報告，阻擋／應改／建議0／0／0，採認限定來源證據。

重跑入口為[intro_picture_probe.go](../../tools/re/intro_picture_probe.go)，在既有離線、非root、有界`wizardry7-go-ebiten:latest`容器內，從本repo模組根編譯：

```text
go build -p=2 -buildvcs=false -o /out/intro-probe tools/re/intro_picture_probe.go
/out/intro-probe /original /out/new-on on
/out/intro-probe /original /out/new-off off
```

repo、SDK、模組快取與玩家原版目錄唯讀，只有新輸出根可寫；兩個輸出目錄不得已存在。原版PNG、PIC及完整收據留本機。探針只到已量200M標題，事件上限512／10000；遇其他OVR載入停用VINIT特定位址觀察。來源SHA為`144b904d07896aaeef1c4fc63ccf1bbe0e2b3152c6daafac96e767956f241c58`。

工具為Go1.24.13 linux/amd64、IDA Pro9.4，分別沿既有`wizardry7-go-ebiten:latest`與`ida-pro-9.4-idapython:locked-v1`。所有IDA查詢複製新DB，原DB不修改；不採跨段字面callee名稱作語意。地址空間分別是resident IDA EA＝runtime IP+10000h、VINIT IDA EA＝code raw offset、OVR offset＝EA+14、VINIT runtime IP＝EA+5055h。

| 輸入／收據 | SHA-256 |
|---|---|
| 原版DS.EXE | `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db` |
| 原版VINIT.OVR | `c1dfd28e8e616070028ac7bd87b5aea5b85d6b6c80224a8a88bbe1774d415e67` |
| 原版MON60.PIC | `ad299595ac2bb31318bd00182db920d78e2dad630b973bf9199b55611cfc9e06` |
| 原版MON61.PIC | `ab6705f244df4ff910192203cebfbcfc53100ea574239c08a02f432566a7bd40` |
| 原版VGA.DRV | `f064349ce5dc694b26ed58615262d1bad893e59971bf259c0a74bdea9ee27a70` |
| 原VINIT.code.i64 | `98b9f93fd55219dfbe327ebd676ae0bac595459cdcbc119e30dde070ae71600a` |
| 原DS.EXE.i64 | `477f73a17b7abcaf3eef70d9baae16c77d9730f3bf7e0b4c0c9c9ed91c612412` |
| r5實際觀察binary | `5de81cd9ddea440c87c02c1e6f27299d53c40dc5ca7ede91e832c93318c642c2` |
| probe-validation-r5.json | `cccd96b13c65dce7586fe5760ca335313e0779dc9ff96746517bb1a915eb67d1` |
| integrity.json | `b11cf095580f9784089d02378e31eddce21668f4583397336b25c0aa16e1e734` |
| review-source.json | `d10a12b55a24c59f6d7bd30b5381839bf40b86a51cadadff9e125f1340db9412` |
| review-source.md | `fe16394cb3c6a504ae4c86b90ae2c97b50ec4c3b4554ad5236bfe3eda9e515d7` |

收據根為workplace/r942，程式基線3581335、SDK基線d16f03a2的受控隔離工作樹。原版248檔、SDK1132檔保持；七個scratch目錄無存檔，擁有權及root-owned／誤建md目錄異常0、自有容器0。r4／r5退出0採主代理工具終止回報，observation／validation本身不含退出碼。公開來源與r5逐位元同且可建置；公開路徑重建binary為`f973db7d8c5892b81e959012767f77c9db18473f98dea61ef5f48b0eb640c017`，未重跑本觀察，100／300／14只歸r5實際binary。

本輪1216個技術檔已一次封存，closed-inputs.json SHA為`46f3e663d6707fe7708c6d048bb7224076dff16e3187275d3ac0942823671429`，排除Go建置快取及後續GitHub中繼資料。封存時原版、SDK、公開探針來源與FINAL報告保持，root-owned及自有容器0。

來源路線已可重播，正式覆繪仍未實作。下一閘門是逐張呈獻字樣的精確矩形、可恢復背景、撤除時機與原生中文字格；經DRAFT→獨立READY審查後再接線，驗兩字型／兩倍率的正式GUI與同狀態原版。圖片內的敘述也須補齊，不以catalog通過代替圖片文字驗收。

## R943 六張呈獻覆繪定版

2026-10-04，使用者選定[甲的周邊紋理補底](../../workplace/r943/intro-picture-comparison-v2.html)，不採平滑色帶。兩種外觀都把「呈獻」畫在原圖片字帶內；底色為新繪近似，不是已證實的原版無字底圖。正式契約為[025圖片字樣規格](../spec/025-intro-picture-labels.md)，初稿DRAFT，未接正式前端。

只處理R942第一組第3至8張完整轉送的六份indexed／palette。逐張核對原版座標格線，修正初稿遺留英文筆畫的左界。已選甲來源為`workplace/r943/prototype_captions_v6.py`；同奇偶周邊樣本按RGB通道擬合色階，加入最近樣本的紋理差，再以原生16×16 Unifont畫「呈獻」。六張、兩倍率共12張候選，各有原版／候選並列。這是離線可丟棄覆繪，不是GUI畫面；沒有新執行原版或注入狀態。

乙使用相同完整字帶及字模，只以擬合色階產生平滑底色。獨立標準函式庫PNG解析、CRC、來源雜湊、射線字帶判定及固定版HEX筆畫核對兩種候選：各12張、4,992,000像素中，批准字帶及字格外變更0；各核對1,920個字模前景像素與來源相同。原版RGB總和≥90且最大最小通道差≤24的中性色像素，兩種候選皆變更0；此數字不等於完整品牌oracle。甲六張3倍並列已由主代理逐張檢視，未見英文殘筆。完整字模、兩字型與GUI仍須正式驗收。

| 本機來源／收據 | SHA-256 |
|---|---|
| 已選甲`prototype_captions_v6.py` | `52d4d0229616d6b5b25a0aa070e0ee0df643764d5df749a9a7b02572f22f70af` |
| 甲十二張圖與補底模型`caption-prototype-v6.json` | `d8f755406e94af215cb510e0db9fc3f15b7e0933f3478bfb46b6eb5916df83a0` |
| 甲獨立核對`audit-prototype-v6-closed.json` | `7602a4d4d0284aa88d8a365108112a5d20e9ac4c2d265c267e16d8b29fa4f79f` |
| 乙十二張圖與模型`caption-prototype-v7.json` | `4777cf06119e837911bf5973e3bb3d2779223e1aed36ebcafce97bfd729adbce` |
| 乙獨立核對`audit-prototype-v7-closed.json` | `2b52815649f21b225a91ad23293e28d33f9953963531ea5ca7621369afcbb037` |
| 定版比較頁`intro-picture-comparison-v2.html` | `8681d725185a12d56665eb74fe8459f457b7fd1ea227465cbce5654c3e82b91d` |
| 輸入保持`integrity.json` | `2731b552749390974a2e66e906de16165eba9b66db9517a3b3829f7084ac4cdc` |

可重跑入口為本機`prototype_captions_v6.py`及`audit_prototype_v2.py`。沿既有`wizardry7-go-ebiten:latest`，輸入repo、R942及固定版字型唯讀，`/out`掛新的已存在輸出目錄；非root、離線、有界執行。設定`WIZ7_FONT=unifont`及`WIZ7_UNIFONT_HEX`為既有固定版HEX，候選程式先核字型與六份來源雜湊。逐步準備來源的私人腳本只為保留修改歷程；重跑最終來源即可，不重跑所有失敗草稿。

原版248、SDK1132及R942封存1216全部雜湊保持；原型階段184份私人檔的擁有權皆1000:1000。原型圖、原版、第三方字型及完整收據不進Git。此段記錄定版前置驗證，正式生效結果見下節；圖片內其他敘述與所有交付門檻保持。

## R943 正式呈獻覆繪驗收

2026-10-04。025經獨立規格審查轉READY才實作。正式TSV只含`presents→呈獻`，由`bake_intro_picture.py`產生獨立sidecar；一般烘製入口同步接入，不讀原版MSG。兩套原生字模及11個烘製拒絕案例通過。完成轉送通知共用既有arm，沒有以5AFD當圖片家族，沒有寫原版RAM、indexed、palette或素材。

完整`cmd/wiz7`回歸及合法原圖守門測試通過。最終Go1.24.13 Linux程式SHA為`2755c68b2bd707bd8575ad69e28c8214a7f649e59c33cde4da57c0ef4451f646`，只是中間程式，沒有正式版號或封包。實作獨立審查0／0／0於FINAL後完整回讀採認，兩報告SHA見025規格。初輪GUI保留，修正原因紀錄後以最終程式另建`workplace/r943/final/`乾淨重跑。

正常無輸入冷啟動，沒有注入存檔、direct-entry或作弊，預設虛擬時鐘同步、圖片驗證用靜音。39／41／43／45／47／49M涵蓋六張，60M仍保持MON60末圖，74M為MON61第一張完整圖並撤除中文。英文控制及兩套字型各跑2／3倍，共六個真正GUI程序，均標準關窗退出0；關窗至退出0.0326至0.2035秒，不當音訊裝置保證。主代理實際檢視正式3倍Draw圖片。

32個中英文停止點的原版1MiB記憶體FNV、indexed／palette摘要及滑鼠相同，indexed、palette及comp另比完整bytes，原版圖片indexed與palette另核R942的SHA。這不包含完整CPU暫存器，也不宣稱記憶體逐位元驗證。追加舊正式`450803c`程式的兩倍率控制，16點摘要及上述完整bytes同新英文，英文HD／Draw完整RGB也相同；兩個舊控制正常退出0。

公開[獨立驗收工具](../../tools/re/check_intro_picture.py)直接解析合法原生字型，以封閉射線字帶及定案前模型計算期望，不調用被測Go函式。32點×HD／Draw的26,624,000個RGB像素全部相符；兩字型字格及全部筆畫完整，批准字帶外差異0。十二張Unifont甲案原型與正式完整RGB相同；末圖另核，翻頁後兩路徑完整原圖相同。補底一致只證明選定的新繪近似，沒有升為原版無字底圖證據。

| 最終本機收據 | SHA-256 |
|---|---|
| `build-final.json` | `03c75251f05b9c28d70bdd4ee7d9b108fbf4f5d9e8205e57352c9988d3e03169` |
| `final/validation.json` | `f0407562fe90371c7f17e0574025bbc17f4048fa39dd4423756a2eabb56b852f` |
| `final/gui-results.json` | `177bb7abaca1cbbfe453d96a1fcc0598a0d074f57cc7e5a200a6bc52fcb13a69` |
| `old-new-control.json` | `6c22c8d3a4b36d594d25e15c3a418233d952d33c10cab180624cf257b2dff955` |
| `baker-boundaries.json` | `5c7717cf327ddcb35259a8f1fdf9cbbca9774698441452c0c4b546f486bb014f` |

重跑沿既有離線Docker，原版、SDK、字型及模型唯讀，輸出只掛本輪既存工作目錄；資源與UID/GID照AGENTS。建置以`GOPROXY=off`及既有唯讀模組快取，完整回歸為`go test -p=2 -count=1 ./cmd/wiz7`。合法原圖守門可另設定`WIZ7_INTRO_EVIDENCE`指向R942/on-r5、`WIZ7_INTRO_RESOURCE`指向本輪Unifont sidecar；缺本機輸入時此測試skip，不產生代用品。

本機`workplace/r943/run_gui.py`固定八點及正常關窗，`WIZ7_GUI_OUTPUT`須指向新目錄；`WIZ7_GUI_CONTROLS_ONLY=1`只重跑舊正式英文控制。所有圖片、字型、原版與完整收據留本機。完整像素重跑入口如下，參數改為本輪唯讀輸入及新輸出：

```sh
python3 tools/re/check_intro_picture.py /out/final \
  --prototype /out --original /src/workplace/r942/on-r5 \
  --unifont-hex /font/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 --output /out/new-validation.json
```

`close_checks.py`另核舊新控制、原版248、SDK1132及R942封存1216全部SHA保持，兩字型來源保持，輸出擁有權異常0。025維持READY，未量過渡及其他圖片文字、十一組、錄放、冷存讀、音訊和正式封包仍待；本輪不宣稱全中文化或可玩交付完成。

完整性報告`integrity-final.json` SHA為`e30b7193948bd2c5511d88eacd856e6e740ee3b3c575fdf3ea3b5180649d2452`。結果獨立窄審0／0／0於FINAL後完整回讀採認；另核32點摘要、64張中文及16張英文HD的完整RGB與尺寸，60M保持、74M撤除。結果審查未涵蓋稍後完成的舊正式控制，該16點由主代理另驗，不能併稱審查涵蓋48點。兩報告`review-result.json`／`.md` SHA為`8e6a034c727241286f78000df4bf6eb1d30ebf0299a3ecc78e8a20623a35a055`／`ff739b1621902124cb64d8e3b256f6f85b6cfec23ecf3496f87955f68ccf5c8f`。報告保存前的容器路徑別名錯誤已修正，沒有改產品或重跑成功像素以選結果。

本輪964份私人技術檔一次封存，`closed-inputs.json` SHA `c8e2aa1ad1af7bbfad351a8c779cfa14124520cc29e2d10a6add523881706e60`，排除Go快取及後續GitHub中繼資料。工作根root-owned／誤建md目錄、自有R943容器皆0。封存檔、原型、原圖及字模只留本機，不加入版控。

## R944 翻頁圖片原型與勘誤

2026-10-04，沿#22及[R944目標](../goal/R944-intro-picture-description.md)，以R942已封存的六張MON61完整畫格核對圖片英文。沒有新執行原版、注入state、讀取MSG或手冊；正式前端、catalog、字型清單及025規格保持。作者姓名依[names規則](../names.md)保留英文。

R942把73235759判成「無呈獻」過早。主代理重新檢視完整indexed按原色盤整數放大，以及只供辨字的近似攤平圖，確認翻起紙頁仍有傾斜presents。原始indexed SHA為`0f098eb7bf96d22777424ac0f8bc263534be1441ecc631413d61a2323714411a`，palette SHA為`d689ae136905afa64a9329ddd123426f7e0f3cde6108c034e0249dbab0a9727c`。本檔目前狀態表已修正，WORKLOG追加勘誤；R942原始收據保持，R943六張平面覆繪的74M撤除驗收仍成立，但不證明此頁沒有英文。

| 原版完整畫格步數 | 可見圖片字樣及範圍 |
|---|---|
| 73235759 | 傾斜presents、平面by及作者姓名 |
| 75322912 | 翻起的角色扮演介紹、平面by及作者姓名 |
| 77352859 | 平面角色扮演介紹、by及作者姓名 |
| 100999086 | 平面角色扮演介紹、翻起的by及作者姓名 |
| 103102872 | 紙頁遮住部分角色扮演介紹 |
| 105202085 | 介紹及作者文字消失 |

當時私人比較頁包含兩種原型，各六張、兩倍率、兩字型。甲替換紙頁英文，原生15／16字模維持水平並隨紙頁重排；補底按三條完整英文詞帶分別擬合周邊紙紋，是新繪近似，仍有接縫。乙保持原版圖片，書下加「呈獻」「奇幻角色扮演遊戲」「作者」字幕；圖片英文仍在。文字全消失的第六張兩案皆完整保持原圖。紙頁角點及攤平模型為假說，不能當原版幾何或無字底圖oracle。使用者2026-10-05已採用甲，正式技術契約仍DRAFT；2026-10-07核對時比較頁及模型已不存在，恢復界線見[R973](#r973-paper-source-recovery)。

獨立`audit.py`使用PNG CRC／RGB解析、射線閉合字帶及直接解析固定版Unifont HEX／倚天Big5字模，不呼叫候選的補底或印字函式。48張共19,968,000像素，22,916個原生前景像素全符合，批准補底區域及字模筆畫外變更0；54個圖片參照皆存在。字帶與字格位置仍取候選metadata，只核前景筆畫及字格邊界，補底區內的字模零位元背景沒有獨立期望。此檢查不裁決補底品質、英文殘筆、瀏覽器操作或正式GUI，也不是同狀態中文驗收。主代理逐張檢視甲五張有字3倍圖及倚天翻頁2倍，補底接縫如實保留。

| 本機來源／收據 | SHA-256 |
|---|---|
| `prototype_description_v4.py` | `0c686dcebbb404d153082bcde480c98ebf98eec693e6e4de44ec41c5a0b83009` |
| `prototype-v4-unifont.json` | `ade1692ff59a1f11219332f4ca15fc289bea795441a85f5016bae656928fe74f` |
| `prototype-v4-eten.json` | `a9a332dcb905c8025ba5e1d950545e790a45939431fa274ce66424c5188f71e0` |
| `compare.html` | `2945817e1ac3141ac8f3f4b4d7cf0e8179bf4cb69eb1d94035b3831d5f0ac48d` |
| `audit.py` | `5a6005422eb9b5cec084c172cb192c27f267a881a11de9b29d0ac639c3378355` |
| `audit.json` | `c0351f2fca599f4ad57636c7f2ba8ca40a19ba5ab568f096b40278115b7b7756` |

重跑只需最終候選，不重跑失敗草稿。沿`wizardry7-go-ebiten:latest`的Python3.11標準函式庫，repo、原版畫格及兩字型唯讀，已存在的新輸出目錄掛`/out`，離線、非root、有界執行。候選入口為本機`prototype_description_v4.py --backend unifont`或`--backend eten`；前者設定`WIZ7_FONT=unifont`及`WIZ7_UNIFONT_HEX`至固定版HEX，後者設定`WIZ7_FONT=eten`、倚天資料掛`/eten`。獨立核對入口為本機`audit.py`，依主機相同層次提供R942及字型；輸出檔既有則拒絕覆寫。字型、原圖、候選圖片、完整模型及收據只留本機。

來源保持核對原版248、SDK1132、R942封存1216及R943封存964全部SHA相同，正式程式及審查過的production來源保持，原型輸出擁有權異常0。`integrity.json` SHA為`76f9778bbfb3376569fb097fbf39d126163d120d458a54b1216d8026204544fb`，另核工作根root-owned／誤建md目錄0。翻頁中文字位置依共同決策待選，之後才寫規格並走獨立READY審查、正式接線及正常GUI驗收；十一組、錄放、冷存讀、音訊、其他圖片按鈕及正式封包門檻保持。

原型獨立窄審0／0／0於FINAL後完整回讀採認，另核48份PNG的雜湊、CRC、尺寸及54條參照，不重跑原型、audit或原版。審查採認上述有限字模前景與區外證據，不採認補底品質、完整作者區域或正式實作。`review-prototype.json`／`.md` SHA為`7c3be89d482f2b7c20b578ef543e670b5ebbb46f2f605df5983c84b41a4f98d6`／`9b4be61bf8d949b91a51d3e28c60db93f987406fe1418fcec35b57b0af2fb866`。

199份私人技術檔一次封存，`closed-inputs.json` SHA為`0897ba8705f3fc995b25aaefbfeb0c8533dc76e0c40ce92ecf239485cfd54bee`，排除封存檔本身與後續GitHub中繼資料。封存時原型、審查輸入及來源保持，擁有權異常0，自有容器收尾另核；圖片與完整收據不進Git。

<a id="r959-battle-control-display"></a>

## R959 現行戰鬥控制的固定上屏驗收

2026-10-05，沿[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)及[R959契約](../goal/R959-battle-control-display.md)。只用現行cmd/wiz7非test來源副本、R943完整雙字型資源及[R958](019-catalog-migration.md#r958-forest-transition-wait)固定合法暖前綴，不新增玩家輸入。原版DS.EXE為Archives版，69206 bytes，SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。Docker沿wizardry7-go-ebiten映像、Go1.24.13／Python3.11.2，所有原版與SDK輸入唯讀。

| 配置 | 有限結果 |
|---|---|
| 英文3倍 | 完整960×600 RGB與原版放大相同，覆繪事件0。 |
| Unifont3倍 | 戰鬥、格擋、咒語、使用、逃跑、移動、裝備七控制14原生字模相符。全圖重建、區外差異、字格／墨點交疊均0，按鈕邊界通過。 |
| Unifont2倍 | 七控制14字模及完整640×400 RGB重建相符，區外／交疊0，但各越出原版按鈕下緣2列。合計448字格像素／99墨點，button_bounds_pass=false，未通過。 |
| 倚天3倍／2倍 | 依2倍越框停止擴測，未執行。 |

三程序各82完整原版摘要、78命令、67既有輸入及三筆歷史NoResponse相同，退出0、Scratch空。初3475M、末5453.1M、六word[0,12,5,3,61,99]／面向1保持；取訊息502、印字669、手冊0、詞庫102、29目標旗標false。純英文comp即原始A000；中文字型兩倍率comp相同，只有HD排版檢查發現2倍缺口。293個Shown事件中只對七控制作原始字型獨立核對，其他286事件只重建reported glyph，不外推其原生字型或完整翻譯驗收。

原版按鈕外框39×9、含1像素邊框，七個原座標半開矩形依序為FIGHT(8,165,47,174)、PARRY(8,175,47,184)、SPELL(50,165,89,174)、USE(50,175,89,184)、RUN(230,165,269,174)、MOVE(230,175,269,184)、EQUIP(230,185,269,194)。七鍵均dy0，現行字格Y為原印字Y乘倍率；2倍字型高16超出外框，3倍高16在框內。獨立審查直接解析固定Unifont17.0.05原始字模及原版indexed畫格，核越框墨點依上述次序為8／19／18／11／21／11／11。這是confirmed顯示缺口，不改原版狀態或點擊判定。

首次清單將MOVE目視誤記HIDE，Unifont3程序exit0但wrapper錯誤報缺項。原錯鍵失敗檔保持，補審後資格訂正另存，不重跑已完成Go程序。來源守門另在候選執行前修正：完全停用動態樹段讀取、三筆NoResponse即時核同才Run；第一次未執行候選及所有訂正保留。正式程式、catalog、字型與SDK均未更動。

| 私人收據 | SHA-256 |
|---|---|
| 候選binary | `ecb6a26f45b3c7e945c714a386e64fbb9bccb68714f792e55f5096f9aba35ac7` |
| verified.json | `e932cf8fcd8ae8667a264871b270bafe2acdb99c0887e25c151c6c8e3a8dc703` |
| review-result.json | `21dce7684491baa8187101625c5714a8ab04a735a59420c83d713ec52edb85b8` |
| review-result.md | `6135fdac0cebf8866c6af10c5947f69339192c2c77d389b5de73df674bf9297e` |
| sealed.json，110份產物 | `9c3ab15cea76251bee5910f6d761a76f2f1d9836c3e2b278e38044350a79a072` |

獨立規格、候選、錯鍵補審、英文／3倍oracle及最終有限結果各0／0／0，主代理於FINAL後完整回讀採認。最終0／0／0只表示證據及缺口報告可採，不表示2倍通過。48選定來源含1份批准goal訂正，47其他來源、SDK1132、原版248、12資源、3原始字型及R958舊204封存保持，擁有權異常0。一次封存110份技術檔，不含gocache或後續GitHub中繼資料。

重跑入口在本機workplace/r959：prepare.py、build.py、run_case.py，使用全新已存在輸出目錄，固定fixed-prefix.json並唯讀掛原版、SDK、模組快取及兩字型。check_display.py獨立核畫格，2倍預期拒絕；check_inputs.py核輸入及原版狀態保持。工具及完整圖片、字模、印字日誌僅留本機。只限暖HD／NoCursorWait，未驗真正GUI、七按鈕操作、背景守門或撤除；001／015保持READY，十一組、TSV、錄放、冷存讀、音訊與正式封包仍須完成。

<a id="r960-battle-control-centering"></a>

## R960 七個戰鬥控制正式置中

2026-10-05，入口：[本輪目標](../goal/R960-battle-control-centering.md)、[026規格](../spec/026-battle-control-labels.md)，沿#28／#22。R959原始失敗保持，依既定原生字模與中文在按鈕內要求，沿020／021公式修正畫筆Y。原版DS.EXE、工具版本、地址空間及字型基線同R959。規格獨立審查0／0／0到READY後，新增battle_labels.go、守門測試及dialog_labels.go三行共用接線；HD與Draw既有接線保持。沒有譯文、字型、事件、低解析compose、原版資料或SDK改動。

七個精確英文鍵／座標、dy0、fg1及完整兩個原生字模同時上屏才置中；缺訊號、其他座標／色號／倍率或非法字模沿舊位置。倚天2倍上移2HD像素，Unifont2倍上移3，3倍保持。字距、X、原字模尺寸與原版命中區不變。

| 配置 | 七控制字格與筆畫 | 原版狀態與畫格核對 |
|---|---|---|
| Unifont2倍 | 14原生字模，越框0、字格／墨點交疊0 | 與R959修前1,346像素變更全在七字格舊／新聯集，metadata與comp逐位元保持 |
| Unifont3倍 | 14原生字模，越框0、交疊0 | 與R959 PNG逐位元相同 |
| 倚天2倍 | 14原生字模，越框0、交疊0 | 全圖原字型重建及區外差異0；沒有R959同字型舊實跑畫格 |
| 倚天3倍 | 14原生字模，越框0、交疊0 | 全圖及舊Y重建相符；不能稱新舊版本實跑PNG逐位元同 |

四配置完整RGB共1,664,000像素與56原始字模實例獨立核同；每側293個Shown事件中只七控制作原字型驗收，其他286事件使用reported glyph重建，不外推其原生字型或翻譯完整度。四側各82完整原版摘要／78命令／67輸入及三筆既有NoResponse完全相同，退出0、Scratch空。初3475M、末5453.1M、位置面向保持，取訊息502、印字669、手冊0、詞庫102及29目標旗標false；沒有新玩家操作。

47選定來源、43候選檔、40正式Go來源副本、SDK1132、原版248、12資源、3固定字型輸入與R959舊110封存保持。cmd/wiz7全部測試與正式Linux建置退出0，Xvfb均終止。首次候選只多一個gofmt空白，格式正規化前來源及binary另留unformatted-*；遊戲只執行格式化後候選，沒有失敗重播或改原版條件。一次查容器資源時原程序已結束，Docker回容器不存在；隨後核程序exit0，不列產品故障。

| 本機workplace/r960收據 | SHA-256 |
|---|---|
| 正式Linux binary | `4b027cb58101a5cc3985d8c90bb8c29d084eae6d40ad47e68aa3c41b6df9c281` |
| 私人固定前綴候選 | `a2f4a80b05fa5f44c137fbc9700f1eb0e4d081ee7dfbdb10035a17ebf8d1bcd1` |
| verified.json | `4525f5fac03ff03435daddcb49d9039e7af5b070de8fff5118f0c03f67407272` |
| comparison.json | `90b4d46b651e84464d908be766681852ba8b2b2a6412e22666473f591c0dedf5` |
| review-result.json | `bfd70a6bcaaeafb4a010bf701a4df225de603a9da26642f0fe8540f3189fb076` |
| review-result.md | `d6b313172b20e78352029e2abf3d7c2593b2fcebbbed3827e137aa7e26ccfe4a` |
| sealed.json，168份技術檔 | `e379aa4136e978196b6807fe63b835b1bf402893b3127dd15d620e7cad1d8958` |

最終獨立實作／結果審查0／0／0，主代理於FINAL後完整回讀兩報告並核SHA採認。審查不執行候選或既有checker，直接解原始字型與PNG再重建全RGB，並核固定前綴、原版及來源身分。公共矩陣已將戰鬥入口移出「尚未到達」列；R959原始失敗作歷史保持，現行缺口以本輪為準。

重跑入口在本機workplace/r960：prepare.py、build.py、test_build.py、run_case.py及check_display.py；第一個Unifont2倍通過後，另外三配置各用獨立有界容器執行。沿固定映像、原版／SDK／離線模組／字型唯讀，新且已存在的輸出目錄掛/out:rw，先確認來源及擁有權。compare_updated.py只核Unifont2倍舊新聯集，check_inputs.py核原版、來源及Unifont3倍保持；報告與輸出獨占建立。原始資料、字型、圖片及完整收據不進Git。只是暖HD／NoCursorWait，不採真正GUI、按鈕輸入、背景動態撤除或CONFORMED，026保持READY；十一組、TSV、存讀、錄放、音訊、待選版面及正式封包門檻保持。

<a id="r961-gui-color-copy"></a>

## R961 GUI色號傳遞與固定Draw

2026-10-05，入口：[本輪目標](../goal/R961-battle-gui-color.md)、[026兩項補充](../spec/026-battle-control-labels.md#r961-gui色號副本補充-ready)，沿#28／#22。R960只驗固定HD；接續核正式snapshot→cloneVisualEvents→Draw發現副本遺漏fg，七控制的fg1守門因而回退。私人來源探針重用R960七則metadata，確認來源fg1／副本fg0、Unifont2倍目的Y331／351／371回退為334／354／374。沒有原版執行或新的正文讀取。

兩項窄契約各經獨立審查0／0／0、FINAL後完整回讀到READY。正式main只補來源fg，保留fgRGB與字模深複製；新增真正clone→buttonLabelY回歸，雙字高／倍率保持、異色回退及修改副本字模不影響來源。cmd/wiz7完整測試與乾淨Linux建置退出0，40份Go中只有main及回歸兩檔修改，其他38保持。

私人入口將R960四份已驗來源畫格送入現行clone及真正Ebiten.RunGame，直接呼叫未修改的正式Draw／Layout。四程序各Draw一次、退出0，來源純值保持、副本值完整及字模切片獨立，Xvfb終止。獨立來源期望直接序列化原事件，沒有用被測clone建立期望。

| 配置 | 固定Draw畫布 | 七控制原生字模 | 全圖／區外差異、越框及交疊 |
|---|---|---|---|
| Unifont2倍 | 640×400 | 14 | 0 |
| Unifont3倍 | 960×600 | 14 | 0 |
| 倚天2倍 | 640×400 | 14 | 0 |
| 倚天3倍 | 960×600 | 14 | 0 |

四圖共56原生字模、1,664,000個RGB像素，原始字型解析與硬編七矩形重建相符。Draw與R960固定HD完整RGB也相同，此比較只為補充oracle。其餘286事件只核已有字模metadata，不當原生完整性或新增翻譯證據。SDK1132、原版248、三個原始字型檔及R960封存168產物盲SHA保持；沒有原版運行，不能以本輪增加原版狀態對拍數。

| 本機workplace/r961收據 | SHA-256 |
|---|---|
| 修正前來源探針gui-clone-evidence.json | `66ca370e43322e372ee7bd91525848ddd3bcf89945347bfd5a0198c832878cc6` |
| 正式Linux binary | `aca28c84c43ddb1b083273dd8ea1033994b8cd5a9a1a31b1eb9444ce78aacda2` |
| verified.json | `99663f2003cff532b07639dc5162154df42d319dbc53d72ea24a782e9d7def0a` |
| review-result.json | `1990547b16d121314e737e2013e6f18df5a450817e4db2b49b9b359a4020091c` |
| review-result.md | `7cf62f123373d9d7b7bdfabd7362af119753e887d418c53724a64b414cbb029e` |
| sealed.json，185份技術產物 | `870f2a0d010b8890daa486cd9bfa76dd43e44943ae0f0647a8fe929e650477fb` |

最終獨立結果審查0／0／0，主代理於FINAL後完整回讀兩報告及SHA採認。審查另寫PNG／原生字模解析，不執行既有checker。185份技術產物一次封存，不含gocache及後續GitHub收據。

重跑入口在本機workplace/r961：test_build.py、frozen_draw_test.go、run_frozen_draw.py與check_frozen_draw.py。使用已存在且全新輸出目錄，固定R960四配置資料唯讀掛/src，SDK／離線模組／原字型唯讀；每個私人程序30秒、整批240秒。這是固定資料的真正GUI Draw，沒有正常owner、正式Update、玩家輸入、撤除、音訊或冷啟動；026維持READY。下一步是正常暖視窗及按鈕操作，避免重跑已完成固定HD或Draw。十一組、TSV與正式封包門檻保持。

<a id="r962-normal-battle-gui"></a>

## R962 正常暖戰鬥GUI

2026-10-05，入口：[本輪目標](../goal/R962-battle-normal-gui.md)、[026私人驗收契約](../spec/026-battle-control-labels.md)。正式程式保持82bbeff，40份Go與SDK未修改。工具鏈為wizardry7-go-ebiten映像 `1d7a342db379675408a489bbca3313a5219f7ac506ae7587f876027d80454ce7`、Go1.24.13及Python3.11.2。原版DS.EXE SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`；原始資料、SDK、字型及舊收據唯讀，輸出只workplace/r962。

R733合法checkpoint及R960固定前綴核82摘要、78命令、67既有輸入與三筆舊NoResponse相同後，交給真正owner、正式Update與正常CursorWait Click。前綴仍用舊NoCursorWait，不當新正常操作。正式音訊限冷安裝；私人暖GUI只靜音節流，空Stream僅PhaseAt，不Accept或產生PCM，非串流觀察器沿既有trace，CPU／PIT保持。私人binary SHA `accb13bae59395ea3f907fd4de73b384c48a08b97e00c812574bc2dfad639e6b`。

| 取樣 | 正常live Draw Step | 暫停底層Draw Step | 畫面與核對 |
|---|---:|---:|---|
| 起始 | 5455300000 | 5458100000 | 七中文控制，14原Unifont字模、256,000RGB、字格交疊與局部差異0 |
| 裝備後 | 5464500000 | 5465900000 | 最後指令，七控制已撤除；只核畫面及事件，未新增子畫面原字模完整性聲明 |
| 返回後 | 5472300000 | 5474700000 | 七控制重現，14原Unifont字模、256,000RGB、字格交疊與局部差異0 |

主代理每步view_image真正暫停底層PNG後，送X11左鍵248,189「裝備」，再160,179「退回上一頁」。正常Click分別5458500000→5464500000及5466300000→5472300000，均無錯誤；沒有調API、改seed或重試操作。三組paused Step等於F1 ACK及resume Step。paused PNG從原drawFlow後、map／help前ReadPixels；live保留正常最終Draw，兩時點不稱RAM相同。

六份PNG的CRC、像素、尺寸與SHA核同，來源事件及GUI副本純值相同；起始及返回共28原字模、512,000RGB獨立重建相符。其他事件使用收據字模metadata，不當完整原字模來源證據。50次Run、四次滑鼠移動及兩次正常Click共56操作有前後Step，末5475100000；NumericStop=false、Manual0、Lexicon102、目標29來源0命中、WarmPacerFault空。正常WM_DELETE exit0，owner、watcher、遊戲及Xvfb終止。來源40Go、候選44檔、SDK1132、原版248、六資源與R960168及R961185封存盲SHA保持，輸出擁有權異常0。

| 本機workplace/r962收據 | SHA-256 |
|---|---|
| gui-check.json | `0319c88a943255de0879ae5096a07fbadf9523892995ce3fcbe42f6f01f1c061` |
| invariants.json | `4d60dec2bf31b92b434d59977a8c0b4e2927dd70388d82cad0e0f239edbd394f` |
| gui/run-result.json | `1ecc5e0b386292140e3edb6328af613dc5263fef98d2607e5f7ffca8a478cf92` |
| gui/operations.json | `beb0519174a299a37fec6d808980abf52128a4577da3ce448d72fe34bd7b62e0` |
| review-result.json | `d2b5df9791555f9744ee3dbe6ac55becc7b57cc56c3cf325aaa8c8f0aeae4e9d` |
| review-result.md | `0b622e385263a26a9e064c457414abf3302b9f1362d2f51c517683f9e6d39842` |
| sealed.json，166技術產物 | `12b1677a5abc4ece8f5f4d127d57f76dc5f0f0397f297741ac5f2772978937c1` |

獨立結果審查0／0／0於FINAL完整回讀兩報告及SHA後採認；另以獨立PNG及字模解析重建，不執行既有checker。報告「其他五控制」文字已改為「其餘控制」，修前JSON保留；本輪只量七控制中的裝備，退回上一頁屬子畫面。166技術產物一次封存，不含gocache及後續議題收據。捕捉、暖節流與暫停底層契約均經獨立0／0／0、FINAL完整回讀READY才執行，修前／修後報告保留。缺/w掛載、冷音訊拒絕及較早live動畫身分守門的三份停止收據保持；前兩者沒有GUI操作，第三份正常關窗0且0新點擊。不把私人驗證契約修正列為產品缺陷。

重跑入口在本機workplace/r962：prepare.py、build.py、gui.py、check_gui.py及verify_inputs.py。使用全新輸出，固定source／SDK／離線模組唯讀；原版同時唯讀掛/original及/w，UID/GID1000:1000、2CPU／2GiB／96pids、網路none、建置300秒及GUI1800秒外層界限。controller逐張保存ready，主代理判讀paused圖後獨占公布command；最多12新點擊、13組取樣，本次兩次即返回退出。

這是Unifont2倍的一條合法暖GUI操作與生命週期證據。其餘七控制分支、倚天／3倍正常GUI、新雙側A/B、冷START、音訊、完整通關與封包未量，026保持READY。十一組、TSV、錄放、冷存讀及完整交付門檻保持。

<a id="r963-battle-round-gui"></a>

## R963 戰鬥操作至戰敗返回

2026-10-05，入口：[本輪目標](../goal/R963-battle-round-gui.md)、[026限定契約](../spec/026-battle-control-labels.md)。重用R962私人binary SHA `accb13bae59395ea3f907fd4de73b384c48a08b97e00c812574bc2dfad639e6b`，不重建、不改正式40Go或原版規則。沿R733合法checkpoint及82點前綴，正式Update、owner與正常CursorWait Click；私人暖GUI保持靜音節流及no_pcm。工具映像、版本及原版DS.EXE雜湊與R962相同，原版同時唯讀掛/original與/w。

| 暫停取樣 | Step | 實際畫面 |
|---|---:|---|
| 00 | 5458100000 | 七個中文戰鬥控制 |
| 01 | 5468500000 | 施法面板，施放／取消為中文 |
| 02 | 5476300000 | 取消施法，七控制重現 |
| 03 | 5484100000 | 使用道具，使用／取消及長劍為中文，另兩列英文 |
| 04 | 5494300000 | 取消使用，七控制重現 |
| 05 | 5503700000 | 最後指令，開始戰鬥／退回上一頁／結束遊戲為中文 |
| 06 | 5511900000 | 開始戰鬥後，命令按鈕撤除 |
| 07 | 5520700000 | 敵人攻擊動畫、中文攻擊說明及命中提示 |
| 08 | 5528500000 | THESUS死亡圖示、生命條清空及中文死亡提示 |
| 09 | 5537500000 | 隊員及敵人撤除，文字事件為空 |
| 10 | 5546700000 | 返回主選單，六個選項為中文 |

主代理逐張view_image真正paused底層PNG，再送十次正常左鍵。前六次為咒語55,167、取消239,126、使用61,177、取消98,187、戰鬥13,167及開始戰鬥159,169。後四次均點選場景空白160,145，只記實際輸入，不宣稱該點有跳過或等待完成的功能。單人隊伍沿原版結果死亡並返回主選單，沒有改seed、重試或作弊。最末finish_reviewed_observation只表示本次取樣收尾，不是勝利、通關或正式錄放成功。

十一組paused Step均等於F1 ACK及恢復Step，正常WM_DELETE exit0，owner／watcher／遊戲／Xvfb終止。22份PNG及packet核CRC、尺寸與SHA，Source／Visual純值相同。三張七控制畫格共42原Unifont字模及768,000RGB核對相符，字格交疊、越框及局部差異0；這是既有家族在新流程中的撤除／重現，不增加已完成家族數。其他子畫面只記實際畫面及metadata，未聲稱完整原字模覆蓋。170Run、九mouse及十left共189操作有前後Step，十次Click結果均空錯誤，末5547100000。NumericStop=false、Manual0、Lexicon102、目標29仍0、WarmPacerFault空。

使用面板實際顯示LEATHER CUIRASS及FUR LEGGING兩個英文道具名稱，列為未完成。這張packet只含已生成的覆繪事件，去重後有LONGSWORD、USE、CANCEL及FIGHTER，沒有兩列的覆繪事件。不能據此判定截獲的根因。另窄核既有items-trans.json已含LEATHER/CUIRASS「皮胸甲」與FUR/LEGGING「毛皮護腿」；compound.go的配對條件目前只覆蓋已量角色頁，未由本輪擴大條件。下一輪沿使用面板實際印字位置窄查；不以單測通過掩蓋漏譯。

空文字轉場的Go nil slice序列化為null，原檢查器直接迭代而退出1。修成空清單解析後，同映像、同一批PNG及packet重查退出0，沒有重跑遊戲。修前checker及有限失敗收據保留。本輪正式40Go、候選44檔、SDK1132、原版248、六資源與R962166封存盲SHA相同，擁有權異常0。

| 本機workplace/r963收據 | SHA-256 |
|---|---|
| gui-check.json | `10fc68384034dd9a69b872a1ee2b2e957b2af7ffcc20738324d30bb4c0ee1f16` |
| invariants.json | `3c3bc858ca6727a87d9f3a5145af3d3d63bc9b86cc5c05a56f1c86d172de553f` |
| gui/run-result.json | `07fa1469a77bfe05c6f84f7d8a939aece174f8a7e537ce70415ffc555bf230e9` |
| gui/operations.json | `b6237318223834ad2dc28105432a7705db0243147d0f7250183b51cf6600ad27` |
| check_gui-before-null.py | `d8a3275d4579e2db8d9b215f33cce6e86d5b8127e0a0137f8fbd6ebcec488f65` |
| check_gui.py | `3d853699f7a0ad4d34a738c2717f3f1647c9a0d3c4b6f353def15259d4dfc62d` |
| catalog-check.json | `94a35b0fb91f70e012c91cc6b1f315b70246c79fd58075821a083a21af223fc8` |
| sealed.json，125技術產物 | `b3d9366c309d9412c0d3327d27239cca64d3411c910b53780328d676610a0f88` |

獨立結果審查0／0／0於FINAL完整回讀兩報告及SHA後採認，review-result.json SHA `f682688577a49a8ee17f0e742af43673183bd0bd75ef8904ca2714b442f6ce6d`，Markdown SHA `e6632fa545986d2c49dc8ddea6f5b1ebf20af3eb81c50cb96abf1880d5c0980e`。審查獨立解析PNG與原字模，沒有執行既有checker；子畫面目視引用主代理判讀，獨立圖片轉送截斷未採目視。026保持READY；格擋／逃跑／移動、倚天／3倍正常GUI、新雙側A/B、冷START、音訊、冷存讀、十一組、TSV、正式錄放及封包門檻保持。原版衍生畫面與完整收據只留本機。

125技術產物一次封存，後續GitHub收據不重封。重跑入口在本機workplace/r963的prepare.py、gui.py、check_gui.py及verify_inputs.py，沿026限定契約與全新輸出。prepare只核R962封存並複用binary；GUI仍為離線UID/GID1000:1000、2CPU／2GiB／96pids及1800秒外層界限，原版同掛/original、/w唯讀。沒有新的正式建置或完整測試重跑。


<a id="r966-profile-reflow"></a>

## R966 人物兩倍整頁正式重排（2026-10-05）

入口：[024整頁契約](../spec/024-character-review-output.md#r966-profile-reflow-contract)、[本輪目標](../goal/R966-profile-reflow-contract.md)。使用者已採用R953重排及R954等級文字右移。沿原版30來源、28個去重mask及31格背景證據實作，只改顯示副本；正式辨識不用固定Step、人物名或範例數值。互動與HD倍率皆2時啟用，3倍保持。原字模15／16像素、13數值沿原版墨點及色號，等級文字x204，數值墨點水平位置保持。

原版DS.EXE SHA `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`，runtime resident地址空間0110；工具為Go1.24.13、Python3.11.2及既有wizardry7-go-ebiten映像。SDK基線d16f03a、1132檔，原版248檔及兩字型12資源均保持。正式main.go SHA `7239edea13889db05b5abfbc6cd156959cc51f23a230e5f7f1e23bfcb0e4f4a7`、profile_reflow.go SHA `b5a4d354969c8156dfd5cf3cbe41f4320affa7fc5b0692f5617f2da13827de46`。正式Linux binary SHA `d827e020712a728f454e2833574fb789c86e56572f85ce29ac3e4ef205accb79`；不是發行包。

完整cmd/wiz7測試通過。首次96pids下Go預設平行組譯fork失敗屬環境，加入GOMAXPROCS2及-p2後乾淨重跑。獨立來源審查發現未認領tuple回true而吞一般入口，已改為撤除當代profile後回false；已認領30來源仍回true，避免舊局部覆繪復活。隔離案例、長值撤除、HP／ST合法區間去重、完整背景、游標容忍及GUI深複製均通過。最後測試成功3.302秒；其後只有測試斷言用字及gofmt空白調整，正式production來源保持。

固定輸入輔助共12次正常Click runner、84快照。修後八個中文配置對兩個英文控制56組完整RAM SHA／Regs／Mouse／indexed／palette一致；舊3倍HD28張保持。另從已推送修前53896f283ac86c22f67111afb5823aecf652c5ca乾淨取得40Go，僅加入七點唯讀完整RAM／Regs探針，兩個英文控制對十個修後配置70組完整原版狀態一致。所有runner正常退出0，PCFILE.DBS與原版相同。這些比較是已量停止點，不外推整段每一步或全遊戲。

獨立oracle從倚天／Unifont原字型及R953原版數值mask重建，未匯入production排版或paint helper。四張正式2倍HD共1,024,000像素完全相符，Draw核31字格共120,112局部像素，域外差異0由完整HD確認；31格交疊、筆畫交集及越框皆0。兩個checker錯配為font變數遮蔽及HP／ST role名稱，修正後以同一批輸出通過，未重跑產品。

真正Unifont2倍GUI從正常角色選單連續檢視男女兩頁。主代理逐張查看目前F1暫停底層Draw，才送12次實際X11左鍵；沒有固定點擊排程或原版輸入注入。26張同封包PNG及13組F1暫停／恢復Step一致；男女四幀完整Draw共138個原生字模及1,024,000像素與獨立oracle相符，其他22幀無profile。等級文字避開蠟燭，數值無重疊，返回原子撤除；正常WM_DELETE exit0、game／Xvfb終止，原版及PCFILE保持。首兩次控制器檔名及Serial欄錯配各正常退出，第三次用同一binary乾淨通過；失敗收據留存，不列產品缺陷。

結果獨立審查0／0／0已於FINAL完整回讀兩報告及SHA採認。review-validation.json SHA `42f12884e9119c0222689aabf19caae619f14cd4a7f60991eff8a3cc8c9b87cf`，Markdown SHA `a1597041ca67d0bbd9daa93e6f0efc2c3313f2f37547650ffdece5d4125600d7`。該代理獨立核摘要、SHA、控制流及oracle來源；遵守禁讀raw字型／像素界線，像素結果採主代理獨立工具收據，不冒稱再次重算。修前控制增補另採窄審。

| 本機workplace/r966收據 | SHA-256 |
|---|---|
| pixels-check.json | `c0a35e1819b315559a8cd992b88ead56ad7e6786df4e9cf9c8196761d97669f6` |
| same-state-check.json | `36fad27389062752450db025eb827fab70474c79fddcb978a5443fc033d80a0b` |
| gui-check.json | `c8cd6aace683d15c357355567ec8438624b00ad36a7167954a8cfa40fa488f81` |
| invariants.json | `4231f4b1ebe906dc4be3ae0e44b1af394dc820517bba117aa34d3fc1093e43eb` |

重跑入口為本機workplace/r966的prepare_validation.py、build_r4.py、gui_aux_r4.py、prepare_gui.py、gui_r4_attempt3.py、check_pixels.py、check_same_state.py、check_gui.py及check_before.py。必須使用全新輸出，原版／SDK／原字型唯讀；既有離線Docker UID/GID1000:1000、2CPU、1500MiB／2GiB、96pids，有界外層逾時及Xvfb子程序收尾。私人探針只加快照指紋與Draw觀察，正式owner／Update／Click保持。完整畫面、原字型及原版收據只留本機。

限定結果為本機Linux、靜音手動加速、兩位既有正常人物；不聲稱其他職業／種族、長數值實跑、所有背景相位、防禦頁、冷START／存讀、音訊裝置／人耳、十一組、TSV、正式錄放或完整可玩封包已完成。原型與正式不同動畫相位不當同狀態全圖比較。

修前增補窄審0／0／0於FINAL完整回讀採認。review-before JSON SHA `5da5daab8d83950426fa4be32a46475712514e23631fbd46ecb0414256e618b3`、Markdown SHA `c4b3f47032bb3426c8eccc97be0403d70a5baceb6f1ea2f23531a9df89c00a54`；before-after-check.json SHA `f9da288b25a78732f9f4d273ce266d6a88d60c2eaa01fbe21b9bc7cf7d68f084`。CPU只採已輸出的Regs，不外推內部狀態。1,039技術產物一次封存，sealed.json SHA `c837db5f223d351a9028d0ad15da02e6cb4484e59d35c9adea390d87263311e7`，排除快取及後續GitHub收據；擁有權與誤建目錄異常0、自有容器0。024整頁增補僅本輪兩位人物CONFORMED。

## R977-approved-ui

2026-10-08使用者確認R976亮度及主選單六按鈕毛筆字樣。[R977](../goal/R977-approved-ui-wiring.md)依[028](../spec/028-approved-ui-fonts.md) READY接線。主選單AR PL UKai TW字面2，2倍16／3倍24像素；紙頁3倍採甲的30／32字高、水平加粗1及米金／紅字。兩後端正式一般字模保持，只追加兩個xlate元資料家族；錄放保持原六roles及同次bytes。

復原工具為Go1.24.13／Debian12的離線GUI映像，來源image ID及重建入口見R977。Python3.13.15／FreeType2.12.1沿R976固定映像及來源。DS.EXE Archives版SHA14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db；原版輸入唯讀，未使用autoplay、防拷輸入、存檔或原版state注入。

| 驗收 | 已量結果 |
|---|---|
| 正常無輸入冷啟動，新舊同GUI條件 | 英文3倍、倚天／Unifont各2／3倍，共35停止點完整RAM SHA／Regs／Step／Mouse／indexed／palette及comp相同 |
| 紙頁3倍五構圖，完整無字翻頁撤除；主選單兩倍率 | HD及Draw全圖31,360,000像素符合已確認字樣；選單按鈕外差異0 |
| 真實F1／WM_DELETE | 兩後端各2／3倍F1返回畫面一致，五個新GUI程序及五個控制程序正常退出0 |
| 獨立字模／資源審查 | 兩後端各8,439份legacy資料保持，24紙頁原glyph與R975相同，24毛筆字條共39,962個alpha樣本與R976原glyph／metrics獨立重建一致 |
| 回退及資料邊界 | 舊xlate、非法新家族、JSON／尺寸／空字模、不同訊號、發布後源資料變更、同次讀後換檔及六roleLoaded；40px錄製頁尾g.Draw守門抽測通過 |

主代理獨立核PNG CRC／解碼、原版摘要及完整像素，未匯入被測烘製或Go繪製函式。紙頁期望取固定R976全圖；毛筆期望由R976原glyph及metrics另算，C FreeType交叉證據沿R976。完整cmd/wiz7回歸通過；另以真實冷啟動畫格核紙頁完整來源、部分重寫、色盤／mode／尺寸、SP、倍率、缺資源及無字頁守門。這些負向案例屬內部單測，不冒稱正常GUI皆已觸發。

首次使用無頭控制組時，Mouse為0／0，GUI為160／109或160／131，驗收條件不等價。舊收據保持，改用同樣GUI控制組後重跑通過。此失敗為驗收腳本條件錯配，沒有修改原版或產品來遷就比較。

| 本機workplace/r977收據 | SHA-256 |
|---|---|
| pixel-check.json | 72a71d09eeefefad5bc4426189a3a78c1ec413fbc5430640baedcfdd915d1348 |
| gui/run-receipt.json | 654e856a7eb5f504b7b2b36d6b6cfa970e8ad191ac5d6a93a0ec4b0c742fbf4e |
| before-gui/run-receipt.json | e114b7182334dc1db0a6eb0cf67bfe4a63748fd50ef2beee1a9a558572cfcaa9 |

重跑用[bake_approved_ui.py](../../tools/re/bake_approved_ui.py)、[run_approved_ui.py](../../tools/re/run_approved_ui.py)、[check_approved_ui.py](../../tools/re/check_approved_ui.py)。新輸出、固定來源、原版／字型唯讀、UID/GID1000:1000、離線2CPU／2GiB／256pids、有界逾時、Xvfb trap及容器清理。全圖、字模及原版收據只留本機。

本輪只驗已確認字樣的本機Linux正常無輸入路線，靜音手動加速。紙頁2倍保持原圖並記限定缺譯；紙紋補底仍為已採近似。音訊工具的Pulse ALSA plugin未恢復，不聲稱音訊／人耳、正式完整錄放、其他選單或完整交付完成。

最終實作／字模／已量收據審查0／0／0已完整回讀採認。workplace/r977/review-implementation.json SHAe8666813d3a899e843f5116e982540ca432ad3464bb67da880a17c31599af5cc，Markdown SHAc3bc7967a8a48a4b72080ecf3d25a73ab4e5b6d8fec808629a08f8cac96a73e6。審查者直接獨立核35狀態及16對圖15,872,000像素；全量70圖的31,360,000像素由主代理工具驗收，審查者核讀法及摘要。

初輪GUI後只修規格編號／註解、兩條非法family log及負向測試，正常繪製演算法保持。最終89來源重建Source SHA40416d488a4105145d8fbed8b18463984187be512a8cb56528690d6a6cb33ba1，workplace/r977/wiz7-final Binary SHA9898b24eb157f44a42a351c63f818acd9976c6673fc7734b3c524d3ece03c030；完整回歸1.546s通過。此binary另跑Unifont3倍正常GUI7點、14圖8,064,000像素全同、F1返回及正常關窗，pixel-final.json SHA30cd33c3353dfdaa4675d7ad1acd90b6d0a476ca74a681fc176771865c0a7432。028僅此已量正常路線CONFORMED，未量GUI分支維持READY。

## R979-clearance-answer-hint

2026-10-08；[#45](https://github.com/wicanr2/wizardry7-cht/issues/45)、
[goal](../goal/R979-clearance-answer-hint.md)、[029](../spec/029-clearance-answer-hint.md)。
使用者明示答案接在問題後面，並要求代理從冷啟動正常遊玩。
這項定案取代先前答案讀取／顯示及代理作答限制，原版驗證仍保留。
答案、完整輸入與畫面收據只留本機；不讀整表或對話詞庫，不新增可散布答案資料。

原版Archives DS.EXE SHA14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db。
定位均為dosgolem runtime resident CS:IP，非IDA線性位址。
confirmed於本次冷題目：0110:0582的白名單取得，與0110:0659同次完成相隔1916步；
SS／DS、SP−0112h、BP=入口SP−2及兩引數核同。只讀已選定記錄最多64 bytes，
驗NUL、四個索引與大寫ASCII正文後獨立複製，不晚讀已被覆寫的堆疊緩衝。
四個顯示家族的座標及數值與前四byte全相符才發布提示；原版核驗由正常前端字元輸入接受。
這個題目的語意confirmed，其餘74題正常上屏及核驗未量。

提示以原生15／16中文字與8px英文字，共用同一已捕help資源、2px字距，
接在中文末行後8px。原色號／palette核到安全帶[70,170,263,178)，輸入白框由y183開始。
兩後端help各補Q一字，14說明頁及原426字形逐項保持；沒有烘製任何答案字串。
HD及Draw取相同owner快照的不可變RGBA字條，錄製40px頁尾先切遊戲區。

四組正常冷GUI、相同固定步數及五次正常Click輸入：
RAM SHA／Regs／Step／Mouse／indexed／palette／低解析comp與持久資料全同；
6,656,000個HD／Draw像素由固定Unifont HEX與原倚天字型獨立逐點計算。
Unifont各新增243個白色筆畫，倚天各285；提示範圍外差異0。
兩後端兩倍率的正常GUI另以真正前端按鍵通過原版驗證、F1辨識general及提示撤除，
全部正常關窗exit0。靜音、明示24M ips，不外推音訊、人耳或預設速度路線。

末輪只補讀堆疊前的entry SP守門，不改排版／字模。
最終90來源SHA5d6d47f75087fb4d68078a5c5eb8db7f2861f43eea09255e16b2e17161f5a3d8，
workplace/r979/wiz7-final-r2 SHAff4d9b4668784d63c06be7849a4b7f84976125b045c0d268bb2ba2ea5e244171。
全cmd/wiz7含正式兩後端六roles捕獲及真正Draw頁尾測試2.421秒通過。
此末binary另量Unifont3倍同題A/B、2,304,000 HD／Draw像素全同，只有243個批准提示墨點；
正常前端輸入通過、撤除與關窗exit0，原版完整狀態保持。
四組前輪證據與末binary抽驗分列，不冒稱四組都由末binary產生。

本機證據入口workplace/r979/frame-pair、gui-*-final、gui-unifont-3-*-final-r2，
hint-pixel-check-final.json及hint-pixel-check-final-r2.json。
Python固定933b46…、Go1.24.13既有GUI映像，原版及字型唯讀，
user1000:1000、network none、CPU／記憶體／pids及外層逾時，Xvfb由明確擁有程序與trap清理。
初次Xvfb握手及續玩控制器等待過短皆為驗收環境問題，正常遊戲未改；修正後乾淨重跑。
014正常保存／冷LOAD及008情境驗收繼續，不把此題通過當全遊戲完成。

最終獨立實作審查0／0／0已完整採認：JSON SHAf9947c8366719d517a7b255782ad6ffe80b04456546c4a967ebb308a438b563e，
Markdown SHAddc606d0ba80f3e7419237bb30d5ec8e85953afc68a6edcfe042cfd25cf487ba。
九份程式／API來源保持；審查固定READY契約，後續文件只追加驗收與狀態。
正式摘要hint-pixel-check-final.json SHAe2578496e7f8ba500d77690c6b5f4cd9caad91448fd9719485ca8c9f801d1206，
末版摘要hint-pixel-check-final-r2.json SHAeca3d5a1c603c75f97c827ad9340ea89056d98f5eeafbea0e43c48c245dc71a8。
029只此已量題目及正常路線CONFORMED；原版存讀與全部中文化保持未完成。

## R979-normal-save-load

2026-10-08，沿同一末binary、Archives素材唯讀、Unifont3倍、靜音及明示24M ips。
從冷啟動正常選人、建隊、本題核驗、31頁原版故事及六頁迷宮開場到可操作遊戲。
重啟後flow資源未恢復，本輪明示停用流文；旁白保留英文，不當中文覆蓋完成。

正常DISK 212／166、SAVE & RESUME 160／73、SAVE 100／170，
SAVEGAME.DBS落盤11,404 bytes；PCFILE.DBS 9,368 bytes，皆1000:1000。
保存層workplace/r979/gui-unifont-3-normal-save-probe/save，只留本機。
此為合法原版保存，不是warm state、改檔或注入進度。

| 正常路線 | 已確認位置觀察 |
|---|---|
| 原版保存 | Area11、X56、Y92、Face3 |
| 首次冷LOAD、本題核驗及原版選檔確認 | 同上；畫面還原THESUS及戰士 |
| 畫面前進箭頭160／166 | Area11、X56、Y91、Face3 |
| 關窗、以移動後保存層相同bytes第二次冷LOAD | 回Area11、X56、Y92、Face3；隊員畫面相同 |

位置沿已確認的唯讀Map欄位及正常GUI觀察，不解析SAVE格式。
三個保存層的SAVEGAME全bytes相同，SHA cf9e462a3f8067162e3ff75dbdbb2ae10f571e98b55d3f1cf68f1af581e8e901。
沒有重新保存移動後位置；關窗不自動保存。鍵盤Up雖已送入，移動未證實，只有滑鼠箭頭列入位置證據。
本機總結normal-save-load.json SHA b1cd1da6e8b80e0fc561c3c54c33a1bc7ff047ee6f75355be3347cc1336b7300，
完整畫面、原版檔及輸入收據不提交。原版root唯讀掛載；未聲稱新完整248檔雜湊控制組已跑。

冷LOAD先防拷、後選檔。首次控制器把核驗後的檔案對話框預期成general，
F1實際mainmenu，屬腳本條件錯配；回讀畫面後修順序並乾淨重跑。
兩次成功冷LOAD、真正正常輸入與關窗exit0；首次新遊戲保存的外層逾時碰到收尾，
原版前端已normal-close0，保存已落盤，不當產品退出故障。
本輪只完成此正常保存／冷LOAD與已列欄位。新冷保存無覆繪同狀態控制組仍待，014 READY、#42 OPEN。
自有容器0，研究／程式／docs工作根的root-owned及誤建.md目錄0，其他交付門檻保持。

## R982-normal-save-control

2026-10-08，沿[goal](../goal/R982-normal-save-control.md)、014／013 READY及#42。
confirmed限Linux正常路線：從唯讀Archives原版及空Scratch冷啟動，正常建隊、核目前題目畫面後
輸入原版接受的答案、31頁故事自然推進、六頁旁白、原版DISK／SAVE & RESUME／SAVE。
正式頁尾停止錄製後繼續遊戲並正常關窗0，沒有作弊、state或位置注入。

正式R979最終binary SHA ff4d9b4668784d63c06be7849a4b7f84976125b045c0d268bb2ba2ea5e244171，
Go1.24.13、90份綁定來源全核同。原版DS SHA沿029；所有248個Root檔案末大小／SHA與錄製前完整清單相同。
倚天3倍沿R977一般資源、R979 help及R981 final flow；正常錄製使用13正式wiz7-input/1，
15,400事件、7,687發布、終Step1,618,600,000，來源及完整輸入只留workplace/r982。

正式同binary／同六資源／空Scratch重播完成exit0。另用相同SDK及原API命令作無覆繪控制，
並於私人前端副本重用正式派送及Draw，額外只讀原版雜湊；副本只有main兩插入及一個觀察helper，
全部起始純量／Storage及Before／After／原API錯誤類別核同，原wire不改。
私人API控制不宣稱正式binary相容重播，正式重播與私人控制證據分列。

兩個API控制各7,687發布，ID連續；完整RAM SHA、暫存器、indexed、palette、Step與Ticks逐欄差異0，
終Step1,618,600,000／Ticks234,939，private adapter完成且exit0，真正Draw確認沿原方法。
GUI錄製、正式重播、private adapter及raw四組保存全bytes同：PCFILE.DBS 9,368 bytes，
SAVEGAME.DBS 11,404 bytes，後者SHA3d855227d5d02844b504da26ba01875ea14251a3fa5dbeb481c48d0d807aa0b7。

此新保存再由正式GUI冷LOAD，原版驗證及選檔正常完成；350M／400M為區域11、56／92、面向3，
正常滑鼠前進後450M為56／91，保存bytes保持、正常關窗0。不解析或修改SAVE格式。
第一個LOAD控制器等待外部命令逾時，已分類為驗收器問題，改有界擁有程序自動送已驗命令，
同一映像、同一binary及相同起始保存乾淨重跑；不當產品故障。

來源審查0／0／2，結果獨立審查0／0／0完整回讀採認。
來源JSON SHA3045bd38702992ba0fef0226e3a53919bd30728b299f708550f08f13a5619ce3；
結果JSON SHAdff94e0ae6ac3c1ef5539c16edf36227371551e32b128f6adfe450314c57f2c8，
Markdown SHAbfd77ba188cf3038462b44b65d4acac0b6556d76d7c06a9e7acfc7df1feb2ffd。
本機save-control-result.json、cold-load-result.json及root-integrity.json保存數字；
record_save_gui.py、raw_control.go、prepare_adapter_probe.py及load_move_auto.py為本輪重跑入口。
正式程式／SDK／wire不改；Linux保存鏈限定門檻已補齊，014的原生平台、封包與完整音訊仍READY，#42保持OPEN。

## R984-normal-profile-context

2026-10-08，[goal](../goal/R984-normal-profile-context.md)、[032 DRAFT](../spec/032-profile-popup-occlusion.md)。
本輪不改正式程式，使用R979最終ff4 binary、R977倚天一般字型／定版兩族、R979 help及R983 flow。
從R982合法保存正常冷LOAD到人物頁；10次原API Click、8份GUI操作ACK與正常錄製停止續玩／關窗0。
八份ACK的名字不是畫面語意證據；最後profile-return是第二次法術視窗，未回迷宮。
正式同binary重播、私人API及無覆繪控制，3257事件／1614發布的原版六類狀態全同，
四保存層與初始層全同，完整數字及獨立FINAL結果見[019](019-catalog-migration.md#r984-normal-profile-context)。

confirmed於正常可見畫面：點姓名出現WHO?選人提示；頭像進人物頁，M前綴仍英文。
法術視窗開啟時，左側LONGSWORD、LEATHER/CUIRASS、FUR/LEGGING、SANDALS、BUCKLER/SHIELD
回到原版英文；取消後中文重現，再開法術視窗又回英文。未施法，未改隊員或物品，背景書本圖示點擊沒有新頁。
WHO與M的真正輸出路徑未量，不據NUL日誌沒有WHO就判為圖片。

唯讀事件診斷從點頭像前邊界320400000開始，取compose前既有事件及indexed，無新字串／原版記憶體讀取。
158筆觀察中，發布1494／Step339400000及1586／375200000兩次開窗，
五種名稱已有Shown=true、Dead=false，但完整矩形比對為false，次一發布五種名稱都已不在事件集合。
不符點範圍均為x49–70、各名稱原y範圍；CAST／CANCEL當刻匹配。
source的compose在已Shown而不匹配時撤除事件，與正常GUI可見英文一起證實本路線撤除原因。
這是限定路線confirmed，不能把一般rectEq代替特殊profile／roster／page守門。
額外診斷的全部1614發布原版trace也與raw全bytes同，SHA為84388c8f8cf1461fc7d6e5bde6fd2d05597ba02855c3ce505ac76ea107c404ae。

未知：caller／font／fg色號、完整視窗遮擋守門及2倍整頁重排互動；未作正式修補或上屏擴張聲明。
032先要求補這些窄證據，維持DRAFT，不依字面CAST／CANCEL或單一座標猜production。
工具Go1.24.13、Python3.13.15／GUI映像Python3.11.2，runtime位址沿0110:0582／0659／3895，
本輪僅新增相對快照發布定位，沒有新的IDA或檔案偏移語意。
字型與正式90來源固定，所有原版及完整輸入／圖片／字模只留本機。

兩份診斷來源獨立0／0／0於FINAL採認：29來源探針JSON SHA-256為
`68a89557747ac57f6058201f96967284467cd53d09b51a6507bc0eb3a38966f2`，
overlay來源JSON為`2c08fc670d7136b32c5ebb25a5185fc592a4aa0e89d0ffc6708f64e3aac8c594`。
overlay結果獨立0／0／0完整FINAL回讀採認，JSON SHA-256為
`c636d80ad562f487b89016801f86d882073750479af5b7f3a0f087e250ea0053`，
Markdown為`012c06e5b9e93042c137f5d967cd53d27bd41512a3a5fc19bed04942db1e60d3`。
第二次含每名兩個實例，共10列但只有五種名稱，兩次CAST／CANCEL實際中文字保持。

<a id="r985-popup-occlusion"></a>
## R985正式遮擋驗收

2026-10-08定位與READY審查，2026-10-09完成正式修補及限定驗收；入口[032](../spec/032-profile-popup-occlusion.md)、[R985](../goal/R985-profile-popup-guard.md)。
本輪以R984的正常人物頁收據繼續，修補範圍只有五個已證裝備名稱。

輸入仍是Archives版DS.EXE，SHA-256 `14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`。
Go1.24.13與既有GUI映像、Python3.13.15；字型／亮度沿R977，help沿R979、flow沿R983。
新增定位全為runtime返回位址：一般裝備第一行0110:ba40、複合第二行0110:ba63、CAST／CANCEL為0110:0ea2。
原版font均0，裝備fg105、控制fg85。定位取原已取得的入口片段與完整事件，沒有新增原版字串或答案表讀取。
兩倍率此世界人物頁皆用一般／012事件，未進024整頁重排。以上confirmed僅於這條正常路線。

032先完成來源與契約獨立審查，READY後才修改正式main與popup_occlusion.go。
原版視窗x49..270／y114..197，24個indexed邊框點、最低19個未遮守門點，CAST／CANCEL雙控制另核caller／font／fg／座標／原框。
只接受五個已驗一般／複合名稱與完整原文字框，還原英文筆畫與中文繪製都避開視窗；每次compose清除舊遮擋。
可見區改寫、錯語意、錯字模、錯控制或背景仍依原規則撤除。未改原版記憶體、索引幀、色盤、語意或保存。

固定輸入為R984原檔，SHA-256 `d7204e01e17698e0a03a764e4521f7c391b03c2e2fe326683c690edcb723e453`。
這是私人API診斷，沒有改wire檔頭，沒有宣稱新binary相容重播。
四組雙字型兩倍率各1614發布，RAM／Regs／Step／Ticks／indexed／palette trace全bytes等於raw；
trace SHA-256 `84388c8f8cf1461fc7d6e5bde6fd2d05597ba02855c3ce505ac76ea107c404ae`。
從舊版原事件的before／after筆畫與原生字模建立獨立批准mask；256張開窗派生HD僅變五名稱區，視窗差異0。
倚天2／3倍每張差2388／4659像素；Unifont為2390／4735像素。
376張未開／取消派生HD全bytes同，四組原版幀／色盤及保存bytes／mode／mtime同。
派生HD是已發布資料的診斷圖，不當真正GUI或原版同狀態Draw。
結果`workplace/r985/formal-pixel-result.json`，SHA-256 `774d81ed2f9adf7314e18b10d88009a135d22204c511cf47922c0793d98be1e2`。

四組正式正常GUI由同一合法保存冷LOAD，再正常頭像人物頁、開法術／取消／再開，皆正常關窗0。
46組人物頁HD／Draw抽樣的名稱區相同；其中22組法術開窗全畫面相同。
Unifont2的345M未開窗取樣另有3176個右下原版護盾動畫差異；-draw-at為非同步要求，不能只憑同檔名宣稱同發布。
前人物頁F1面板不在本契約比較範圍。四組保存size／SHA／mode／mtime及UID1000:1000保持。
GUI結果`workplace/r985/gui-pixel-result.json`，SHA-256 `6b4b8ee4eb8f1d4168a7b55ceafc17c3ad29df1469cf2e18f91450e5708a000c`。
三組測binary `eb2eda6721e259f9fc0833039d766ecbffec446538404547a9f0b4a97d0f8909`；
倚天3最後測 `8cec6479c2d1a527b4efac131ecf262dc088d74ff56abd3fb4e50a301ca11fd2`。
兩者只有註解、匯入順序及新增測試的來源差異，正式main與render守門語意保持；完整身份在本機final-identities.json。

正式TestPopup六項通過，含原版fixture的雙倍率／15、16字模及壞框、錯語意、特殊事件、複合caller／font／座標／原文負測，
clone獨立所有權，以及真正RunGame的20組逐像素Draw：正常、40px頁尾、非零subimage、錯寬及錯高。
Draw只畫批准左側，視窗、頁尾及不合畫布全保持。原版fixture不提交，缺本機fixture就skip；通用clone／Draw照常測。
每次compose的開／關與撤除另由158發布固定路線觀察核，未量其他人物切換或裝備組合。
WHO／M、其餘名稱、十一組情境及31異譯仍待，不把本修補當整體中文化完成。

最終獨立實作審查0／0／0於FINAL完整回讀採認，032五名稱正常路線標限定CONFORMED。
JSON SHA-256 `f81ae5c42ba0812cb7eaa10cb946b251a3b162b5e6a987a549795445e3c545d0`，
Markdown `f1e3bd9f121de7d0f4438267751d70935a8a9221f3a6aae0f8965509632fe093`。
新增最後未遮邊框拒絕後，全cmd乾淨重測1.494秒通過；重新建置與正常GUI倚天3測過的8cec二進位全bytes同。
來源與驗收輸出UID/GID1000:1000，自有容器、root-owned及誤建.md目錄0；新goal／spec以R984正對照核有入口。

## R986 世界人物頁性別與選人提示（2026-10-09）

<a id="r986-world-identity-labels"></a>

入口[R986](../goal/R986-world-identity-labels.md)、[033](../spec/033-world-identity-labels.md)、[034](../spec/034-world-profile-sex.md)。
沿R984合法正常3257事件／1614發布收據，三個獨立審查過的唯讀探針各全trace等於原版控制，
RAM／Regs／Step／Ticks／indexed／palette均保持。原版DS.EXE與正常輸入SHA由034收錄，不修改原版或wire。

性別confirmed於runtime0110:3877、返回0110:6159，font0、x90／y6、M字碼77。
返回只改原6×6格13個色號1墨點。0110:60c9起320-byte正常執行code在IDA9.4窄解碼，
保留原operand [bx+5e2ch]，讀值加140、取140／141原訊息首字碼，再由0110:6156呼叫3877。
F與M共享資料流已證，但正常F上屏未量。完整8×8背景守門及原生字格邊界見034。

WHO?在正常選取時有45個色號22像素，取消後消失。窄捕捉到runtime740c:20fe的
mov es:[di],al透明位圖寫入，寬24／高5；這是寫入端，不足以識別上層來源。
WHO code256 bytes SHA-256 580cab56d91d80d609ceeb521e4ecd42d01f0378f8c154013b46fe47bf5bd64b；
M code320 bytes SHA-256 3d8b652dbbcbd17d553a84878a132e1a14a44f6cf0eeff4ec30dc2c7b1148259。
IDA raw snippet EA與runtime段位址分開記；線性oracle At544958不混作IDA EA。
033保持DRAFT，不由位圖writer直接實作覆繪，不深入與目前文字無關的驅動程式。

034契約獨立FINAL 0／0／0完整回讀採認READY後正式接線；來源與READY審查SHA見034。
正式全cmd與七個world負測2.359秒通過，binary SHA-256
57f5e3d2328098969c5e64526ec3aff53fb5defea535fd835e2867bf4e379076。
兩字型兩倍率正常冷LOAD、人物頁、法術開／取消／再開完成。
31組完整HD／Draw在性別批准域逐RGB等於獨立原字模期望，五裝備名稱保持；四保存層bytes／mode／mtime相同。
倚天3倍355M只有HD、未及發布Draw即正常關窗，單獨記為未配對，不列核對分母或產品故障。

正式來源只加三個私人只讀診斷hook，四組各1614發布trace全bytes等於raw。
16個派生HD共6,656,000像素，以舊R985幀、13原墨點及外部原字模獨立重建，逐像素全同，批准域外0。
每張變更像素倚天2倍76／3倍105、Unifont2倍76／3倍123；保存bytes／mode／mtime保持。
固定收據不改wire、不稱正式相容重播；派生HD不稱正常GUI，正常畫面另驗。
安全摘要SHA-256：formal-pixel-result.json為3078a2116ec5ada8ed356184750d157aad64b28619bc110f165b822f7a1eaa27；
formal-gui-result.json為21c3b39a030e37954b44b7c170983556334417096df0d941c3bf0f98be030a98。
完整原版、字型、code、答案及輸入只留ignored workplace/r986。

另兩字型兩倍率正常Exit返回迷宮／人物頁重入共58組完整HD／Draw pair，
退出批准域等於原indexed放大無殘字，重入原生男位置相符；四保存層保持，關窗0。
formal-exit-result.json SHA-256 6152aa0d764d06bbf6df35b54168f71aa50ddfc673124712de89427ff2c3e8b7。
[原版手冊](https://www.mocagh.org/sir-tech/wiz7-manual.pdf)第24／25頁說明盔甲右側Exit；正常點315／92已量回迷宮。
頭像點擊及Escape只有操作收據、畫面未退出，未列成功；驗證器一度全域替換suffix連帶變更標籤，
已限縮只改suffix並保留原嘗試，不混作產品故障。
最終實作採認及提交在收尾追加；WHO、F與整體交付仍待。

R986最終獨立實作審查於FINAL後完整回讀採認0／0／0，已量M正常路線限定CONFORMED，F保持READY／未量。
JSON SHA-256 `4b73e468805d2e4c6e7e99b136660f71bcded169e6ae26f40e75a56f68e11dd7`；
Markdown `4982117d1d553adb15a30690fcbc3f2a4fabbe8b34655a6294791fd4b8a70764`。
六份正式來源／TSV／binary與審查SHA一致；034其後只補實際驗收與採認狀態，不改契約。
89個正常GUI完整pair已獨立核，6個退出pair無殘字、9個重入pair男字位置正確。
自有容器、root-owned及誤建.md目錄0，來源／輸出抽查UID/GID1000:1000；
以既有R985／032作正對照，新goal、033／034與025證據入口均可定位，R433保持未追蹤。

## R987 WHO游標來源與原生字樣（2026-10-09）

<a id="r987-who-call-witness"></a>

入口[R987](../goal/R987-who-call-witness.md)、[033](../spec/033-world-identity-labels.md)。
已得實際740c:1568／0166返回；1568清理於1575 retn，不採40b0 far框假說。
IDA9.4的192-byte返回code SHA-256為8692955a85d421ed8dfbafcad0f6bfe61ede16a1c3d320133e430af6986ada7c。
raw snippet EA與runtime段位址分開；1575以後線性解碼只作候選，不當完整函式語意。

WHO整圖24×20、字母墨點5行、文字末區6行。當刻來源84d3:0002，欄位及來源SHA見033。
45個色22墨點可去為本字帶底色0；透明255保留補底前原RGB，11點與全0背景buffer不同已明記。
新字模可占用文字域的原透明位置；不能把11點最終全不變當oracle。
行12非透明手尾及黑輪廓7點全部保護，行13空文字背景0／255不等於indexed全0。
完整15／16字模以一像素字墨間隔利用標點留白，四派生HD手尾及圖外變化0，無背景擴張。
native-compare.html與native-who-sample-result.json留本機，未接正式UI，不稱正常GUI。

原版DS及正常輸入SHA沿034／R986。已完成的1614 trace逐bytes等控制84388，保存bytes／mode／mtime保持。
首輪hdonly私有API缺Draw確認，已從main的select及afterPublish等待定位；保留一發布r1後同來源改Xvfb正常Draw。
第五探針1614發布及重播完成存在，但180秒收尾相遇，程序退出未採認，不列產品缺陷或正常GUI成功。
IDA locked映像仍查無，沿既有py312／UID1000已驗9.4組合；新code匯出schema及輸入SHA相符。
完整原版、圖、code、堆疊、答案及輸入只留ignored workplace/r987，不提交原圖或字模二進位。
033仍DRAFT；獨立原型審查、正式契約／接線及正常WHO移動／取消仍待。

R987最終獨立來源／原型審查於FINAL後完整回讀採認0／0／0，32份已讀來源SHA保持。
四派生HD共1,664,000像素與獨立期望相符，手尾7非FF及文字域外變化0。
補底階段保留13透明點原RGB；新字覆蓋原非零透明點倚天2／3為3／1、Unifont2／3為2／1，不作全不變宣稱。
JSON SHA-256 d97f39354668e8c7ffa58e1d79ceff4711e00e5574004316b35e5cf22f5c9634；
Markdown 17d977356abab9792ba68712b0f2e3ad08a2539948c8801fdb8a3c7e0c198e0d。
仍僅DRAFT派生原型，正式來源、資源、正常GUI及CONFORMED未放行；WHO字型／畫布沿既定15／16、2／3倍。

## R988 WHO游標正式中文（2026-10-09）

<a id="r988-who-cursor-wiring"></a>

入口[R988](../goal/R988-who-cursor-wiring.md)、[033正式契約](../spec/033-world-identity-labels.md#r988-who-cursor-contract)、
[WHO譯文表](../../text/cursor-who.zh-TW.tsv)、[字模生成入口](../../font/README.md)。
獨立契約FINAL完整回讀採認READY後接線；只唯讀辨識R987已證24×20原版游標，
來源tail／prefix SHA及runtime740c欄位沿033。完整字模保留15／16像素，在2／3倍原圖內排「誰？」。
WHO?資源從同次xlate解析後移出一般NUL字典，保持原版字串及原API、六資源schema與身份綁定。

| 驗收 | 實際結果 | 限制 |
|---|---|---|
| 正式cmd回歸 | 10.711秒通過；來源、位址界限、字模、通道倍率及真Draw拒畫／40px頁尾保持 | 私有原版fixture缺少時skip；防禦分支未稱正常量到 |
| 正常GUI | 雙字型兩倍率，28個正常LOAD後選人、移動、取消及再入操作，四趟正常關窗0 | 3倍第一次1CPU在主選單截圖等待逾時，改2CPU同工具鏈乾淨重跑通過 |
| 正常HD／Draw | 倚天2／3倍6／10組、Unifont2／3倍7／11組，共34組文字域及手尾逐RGB等獨立期望 | 原圖下8行；上方透明背景有既有職業覆繪，不稱完整正常畫格同狀態 |
| 固定原API | 四組各1614發布，RAM／Regs／Step／Ticks／indexed／palette trace全bytes等原版控制 | 私人診斷不改wire，不稱正式相容重播；只加三個唯讀hook |
| 派生HD | 32完整圖共13,312,000像素等R986基準加獨立原墨點及原字模期望 | 不代正常GUI；WHO有效8圖，非WHO24圖保持 |
| 保存層 | 四GUI及四API保存bytes／mode／mtime相同 | 沿合法R982保存，不注入記憶體或修改原版資料 |
| 資源 | 兩後端8441→8442鍵，只增WHO?；所有舊entry及三同目錄家族全同 | 字模及完整資源只留本機，正式文案由TSV供應 |

去墨只改45原色22點，透明及既有黑點不整片塗底；手尾7非透明點與原24×20圖外保持。
八個有效派生HD的變更像素，倚天2／3倍每圖226／411，Unifont2／3倍237／410；
全部落在批准文字域x30..54、y65..73，非WHO六個發布逐完整圖保持。
GUI快照保原WHO；只有實際Draw畫布及倍率合格才改本次RGBA副本，拒畫仍顯示原圖。
HD另用自己倍率及comp副本，原生字條GPU快取有界，取消後不沿用舊view。

正式binary為workplace/r988/wiz7-who-final，Go1.24.13，SHA-256
`6fb3bb7f408ec216c3a572d9f8f6f85aef3170b8d72fec82e763068c86868a29`。
原版DS、R984合法輸入及raw控制trace SHA沿033；原版、SDK、字型及輸入唯讀，未改答案或驗證。
安全摘要SHA-256：API `1b99f4e710d5fa7fe771a60aa8739e615c260d1bd00f444b82adc2209880e73a`；
正常GUI `c4aa228a4a97951462fdecf4d6bff29c9d1dc199c670cb176f9e896937b85413`；
資源 `ae624e50afdc895e71256e0b276925d7a48ad1b9c6cbc02f71cdc8564dea61ae`。
完整原版、位圖、字模、答案及輸入只留ignored workplace/r988，最終獨立採認另追加。
部分裁切、其他人物／流程及資源故障正常路線未量，保持READY與原圖fallback。

R988最終獨立實作及像素複審於FINAL後完整回讀採認0／0／0，22份審查輸入SHA相同。
JSON SHA-256 `369c275a6bbda48ebaf0c5d857b27f43d5fe79e3bbdd206080e922a473f2506f`；Markdown `41b69a5c38dcb7a1715199165890670300278f01914942de78009664b3663242`。
033已量WHO正常選人、移動、取消與再入限定CONFORMED；未量分支仍READY。
源碼及binary保持；033其後只追加驗收與採認狀態，不改契約。
自有容器、root-owned與誤建.md目錄0，來源及輸出UID/GID1000:1000，R433未追蹤保持。

## R990 女性人物頁正常玩家路線（2026-10-09）

<a id="r990-normal-female-profile"></a>

入口[R990](../goal/R990-normal-female-profile.md)、[034](../spec/034-world-profile-sex.md)。
沿R989正式binary正常選取既有TEMPEST並加入隊伍，完成原版驗證、影片、六頁迷宮旁白及SAVE。
新程序正常LOAD恢復女性隊伍，實際人物頁顯示女、蜥蜴人、戰士及等級1；EXIT與重入均有正常實圖。
沒有修改state、人物記錄、原版規則、驗證判定或錄放檔頭。
原版Archives DS.EXE、Go1.24.13及執行環境沿034／R989；正式binary SHA-256
`a9f8db0348f572861664736d330ba0aa3f5b5fb5eb92a73374d14fa8edb540e1`。

新SAVEGAME.DBS為11,404 bytes，SHA-256
`0d1083aeb4cb352cfa7d284b93af97d3f674a53e4a6fed6affbd0a5eacf7b84c`；
PCFILE.DBS為9,368 bytes，保持`24d94441e5b32e7abbcbffa9dad972b41544da6a7bb06096b967281f639d9c1b`。
正常保存錄製16,222事件、末Step1,741,000,000；冷LOAD錄製2,856事件、1,415發布、末Step312,000,000。
兩份均按遊戲區外頁尾停止後續玩、正常關窗0；完整輸入及答案只留本機，不提交。
女性保存層在`workplace/r990/r3/record-gui/save`，兩份正式輸入在`r3/normal-female.wiz7-input`
及`r3/coldload/normal-female.wiz7-input`。起始保存只以正常遊玩取得，不是注入的暖state。

冷LOAD正式同binary重播完成0；另外私人原API覆繪端與無覆繪控制均正常0，
1,415發布的RAM／Regs／Step／Ticks／indexed／palette trace全bytes相同，SHA-256
`17217541aa29d50a7288fbf8d66904dec808df09f7d1c6b464e9f9e7f9a9eae9`。
四個冷端保存層的八檔對新女性保存全bytes／mode／mtime一致。
私人原API不修改wire，不稱正式binary相容重播；原版同狀態證據與正式錄放分列。
安全摘要`workplace/r990/cold-control-safe.json`。

獨立玩家驗收FINAL完整回讀0／0／0，範圍限正常保存、冷LOAD、人物頁EXIT及重入。
JSON SHA-256 `60d712176456fc4b8df2a2502a0240976e5ceb56cfcfc644d02cc9acb280c1e1`；
Markdown `26a9fb253a051ccfd609a0be80559a83913efca0b3d3dbc90e6a4f5c5441eb25`。
此玩家報告只採認正常保存鏈；原生字模、F原版返回及新遊戲完整重播依下列獨立收據，不由玩家報告單獨放行。

窄F探針首版把code／font引數傳反，因守門拒收沒有三個F檔。獨立審查已追加勘誤，舊輸出保留。
只修私人main一行，正式來源保持；修正後FINAL完整採認0／0／0、七來源SHA一致才重跑。
本輪正常冷LOAD實際confirmed於runtime0110:3877→0110:6159，F字碼70、font0、入口90／6，
返回Step279,430,711、SP65312、96／6。完整離屏只變11點，全部在原6×6格並改為色1。
八背景點符合034，修正後1,415發布六類trace仍逐bytes等同無覆繪控制，保存bytes／mode／mtime保持。
before SHA-256 `241521b652af84ebc1341ae156e345d1db39d009e3a2479dcb1826cd61730a18`；
after `18e80a9ac76ec33352dd5f4b1d53b5c747910077fb7eb3c4993de16d645143a5`。
原版code與141訊息來源沿034既有資料流證據，不再讀鄰文或手冊答案表；本輪只增實際F印字見證。
安全結果`workplace/r990/f-witness-safe.json`；私人探針binary SHA-256
`0fb0348341aab68411f2143f40589f06dde845f5457643eec223d2dcad8f2565`，不是正式可相容重播binary。

完整新遊戲／SAVE收據由同正式binary重播至Step1,741,000,000，完成0；
兩保存檔bytes及mode等新女性進度。SAVEGAME.DBS實際重寫的mtime不同於首趟壁鐘時間，
不把正常寫檔mtime當固定oracle；PCFILE.DBS及冷LOAD流程的bytes／mode／mtime保持。
安全摘要`workplace/r990/formal-save-replay-safe.json`。本輪沒有新遊戲整段逐發布無覆繪控制，
原版六類同狀態聲明只限上述冷LOAD的1,415發布，不外推整段SAVE重播。

正常加入TEMPEST後的五列候選名單仍有M-THI、M-PRI、M-MAG及F-ALC英文碼。
靜態confirmed：`rosterRows`把完整原文與初始y綁定，022只批准六個初始位置。
加入後列上移是實圖強證據，新的NUL引數與完整印字差異尚未補量；不直接放寬姓名或座標守門。
安全記錄`workplace/r990/roster-gap.json`。後續須窄量移位列、022增補READY及正常名單驗收。

F原版窄收據獨立FINAL完整回讀0／0／0，七來源SHA保持；審查另寫計算，未執行被核驗的摘要產生器。
JSON SHA-256 `a951588730069e5af937ce5459aaec836bd0afe93e9c86c37f6a87a97a800168`；
Markdown `615b8e3f0cd577340e867821d32bd28c5eb55628cc9c02e3d2e171ea0bcb422a`。
採認限F原版差異、八背景點及冷LOAD六類摘要／保存保持；正常GUI像素及全F路線不由此報告放行。

主代理另目視倚天2倍女性人物頁與Unifont3倍重入實圖。
2倍世界人物頁的一般屬性標籤仍有行距不足，與本輪性別8×8批准域分開記錄；
不由女字正確外推整頁排版已完成。022名單與世界人物頁其餘欄位仍須各自完成覆蓋及排版驗收。

四組正常冷GUI完成兩字型兩倍率女性人物頁、法術開／取消／再開／再取消、EXIT迷宮與重入，正常關窗0。
70個有效HD／Draw配對共58,880像素差異0，原生女字域51,072及退出原indexed還原7,808。
期望獨立解析外部stdfont.15的女字30 bytes及Unifont U+5973的原生16列；
本輪F-before／after各只取批准8×8背景64點，先核當刻原indexed等F-after，再計HD／Draw期望。
倚天2／3倍16／17組，Unifont2／3倍18／19組；六張F1 Draw及一張無Draw不列配對。
四GUI保存bytes／mode／mtime保持。最終結果SHA-256
`5ae1a75a9068ecdec12a95669a041e27f206641a586817aaea89b3aea5efbf42`。
獨立GUI FINAL完整回讀0／0／0：JSON `026be067bc3c6c707bee4e2d7a4873589cd023152f85466c7e3af0fd18efb601`；
Markdown `55b1fbec3c96a820317ea2fb185ea6d6a99eade4d98a3f6904e663d5f246104c`。
結合上述F窄收據，只增034已量TEMPEST女性路線限定CONFORMED，其他人物、全頁排版及全遊戲不外推。

## R991 加入角色後的移位名單（2026-10-09）

<a id="r991-relocated-roster"></a>

入口[R991](../goal/R991-relocated-roster.md)、[022 DRAFT增補](../spec/022-character-roster-text.md#r991-加入角色後移位列增補)。
R990正常加入TEMPEST實圖的五列名單，四個既有碼退回英文；`rosterRows`完整鍵與初始y綁定是已證實的實作限制。
沿未改R990正常新遊戲收據，窄探針只重用既有16-byte來源，六完整字串、六y、font0及x150硬守門。
十來源獨立FINAL 0／0／0完整回讀後執行，原版及正式來源未改。
審查JSON SHA-256 `ac6b92e2b6a5dbfdb9fddc1129160e7ac6fd85eb240db9af83002c08c1a7fc1d`；
Markdown `03f929f24609b236b8f59b86a014af74a4283de791e00370fde7c9b39566840c`。

已取得17返回，四移位列confirmed：LYSAND y162→156、NOBAL168→162、TREON174→168、PENTAS180→174。
resident入口0110:3895、返回0110:0fcc、font0、x150→240及y保持，SP63968→63970。
完整64,000-byte前／後差異只在已知姓名格及原碼帶，姓名差異70／64／63／76，碼差異49／53／59／50。
四列全部195→180、域外0；選取1→0另由正常新錄製補量，不由未選取證據猜補。
欄位及完整字串仍沿022原六來源，不解碼PCFILE或推論任意姓名record格式。
原版DS SHA沿034，Go1.24.13，private probe SHA-256
`a2c36fdfeecc65c205a0ad67de49b6bc961cf6422c763a510801f065c2cb5127`。
完整before／after及收據在workplace/r991/private-front，安全計算入口`inspect_roster_witness.py`。
目前022增補DRAFT，不以此稱正式修補或CONFORMED；原版整段控制及正常選取續驗。

未選取及全新遊戲控制獨立FINAL完整回讀採認0／0／0，五來源及53輸入SHA相同。
8,096發布的原版六類摘要逐筆相同，trace SHA-256
`fc6a8b0bf5fa9c56f000cc9135b491369c530c4fd11010a360772c87dc4aec2f`；四保存檔bytes／mode等R990正常新進度，SAVEGAME重寫mtime另記。
獨立報告JSON／Markdown：`492531e1cde96686806e6ebf61b979e92626e7038703fc6a51ee517212368221`／
`41ec0fdf200dc7f0f419f7ee56c7c5000adc63ebf67fa3da591b93c5b74651bc`。

未修前新正常GUI只加入TEMPEST，四移位列各選一次，再取消、重開確認白底撤除及取消回主選單。
10個正常click，3,104事件／1,546發布、末Step367.2M；按頁尾停止後續玩及關窗0，沒有進新遊戲或防拷。
新收據SHA-256 `6be4b82d291750f8a93bdeb2c6dd898353c3cbd259523b55e31e9e584d535dcd`。
正常玩家獨立FINAL JSON／Markdown：`da8ce3c452ba3126be924ba2faeb51bac3f714662710a1b2543a346be80e5c79`／
`1f1debe1ee7e61a6bb74e5a204073b301c3c2748b3b403e9c0945d39d1b625c2`，全文及四來源身份核對相同。

同一探針沿新正常收據取得42次返回，四移位選取ID19／25／31／37皆1→0，姓名／碼差異同上表、域外0。
原版兩側1,546發布逐筆摘要及完整trace bytes相同，SHA-256
`b63495ee7c41845a4c02de1c043ee70dcfff620486be7d82ddb71169c68a90ab`，三保存層bytes／mode／mtime同R982起始。
獨立選取及契約FINAL各0／0／0完整回讀，九來源與126證據檔fingerprint相同後022增補採認READY，參見022審查SHA。

正式只改四個完整來源與位置配對，metadata以副本帶實際y，保留原六表、姓名範圍、finish及底邊契約。
已知y的錯來源／SP拒絕仍先撤除舊事件；讀前SP守門不多讀原版資料。
全cmd2.816秒、initconfig0.024秒、recording1.213秒通過；正式93實源建置binary SHA-256
`55aa21626a977d750249af677ad8a3257ee6f709167c69dbd83c42541479c73d`。
獨立source FINAL 0／0／0完整採認、八來源SHA保持；JSON／Markdown
`4ad8754e02abe1a0247200c080ad545ca2692046f517f60586ac29cb791d067b`／
`fc00ddee3e082ff5a25e45460e875aa08ea4e6595a132dde4b10e2d20035db07`。
正常兩字型GUI、雙倍率原樣基準與限定像素驗收在本輪收尾另記，不由source審查宣告完成。

修前／修後兩字型兩倍率各一趟私人原API控制完成0，每趟1,546發布全trace等原版控制b63495。
每組五個相同發布點的原indexed及palette全bytes相同。20對完整派生HD共8,320,000像素，
以修前圖、原英碼墨點及外部原生字模獨立建立期望，修後逐RGB相符、四碼帶外差異0。
2倍兩組十張完整圖保持；3倍每個有效名單樣本改動倚天2,149／Unifont2,179像素，取消樣本全圖保持。
八保存層bytes／mode／mtime等R982起始。安全結果在`workplace/r991/fixed-pixel-result.json`。
私人派生HD不代正常GUI、不稱正式相容重播；固定source相對正式main只有既有三個唯讀hook，原roster全bytes相同。

固定驗收最終獨立FINAL完整回讀0／0／0，八來源及168資料指紋相同。
其20對完整old／new範圍核對另證姓名區x150..210、y150..186與取消發布1465全部RGB保持。
獨立代理直接讀原字型，另以Pillow建立四張完整原字模期望，2,304,000像素全同；
這四圖抽樣與主代理20圖期望、20圖差異範圍分列，不擴充獨立抽樣分母。
報告JSON／Markdown SHA-256：`c1967578d202836601664ed6ce4f1477b12909b8153a9fd0ce06986a7b7c1c31`／
`7e15233c4b63deb86eebed56306dd2ebe6c8c540d30fb7207e48e0b94daeadf7`。

兩字型3倍正式正常GUI各10次click完成加入TEMPEST、四移位白底選取、取消、重開五列及關窗0。
最終有效oracle是`formal-gui/normal-gui-oracle.json`，SHA-256
`67096c35145a2ac355dcbf60647ef6048c2df60b28e0d77e0c7da6b761ccd70d`。
64個HD收據按請求標號配對，依實際Step／indexed／palette／Mouse去重後55唯一pair，倚天28／Unifont27。
4張實際F1排除，5重複不重算，HD-only0；完整HD／Draw31,680,000像素相同。
獨立原字模及42原字段核ROI3,207,600像素差異0，其中碼帶703,080、姓名與碼帶外2,504,520。
六保存檔bytes／mode／mtime同R982起始。錯分母暫核`normal-gui-oracle-r2.json`保留，不作最終結果。
獨立GUI FINAL完整回讀0／0／0、七來源及有效oracle／新wireSHA相同；JSON／Markdown
`a4c6781b72c6712e4759751046b60568a7afc912b3ee18d1291e05603ce9570f`／
`2c437120b6f076dde58a6d023fbed1e481f0a0cf5a960b69f7d7fa6af59e1e21`。

新正式55aa binary正常record相同10操作，3,110事件／1,549發布、末367.8M，頁尾停止後續玩、關窗0。
同binary正式-replay一趟完成0，原wire SHA-256
`820b2ed71778f8375733f8496128878d4ad05ea95f12c3c1b8e58591ce42f2d3`保持，重播保存bytes／mode／mtime保持。
source／固定控制／正常GUI獨立結果採認後，022只增四已量移位配對限定CONFORMED。
任意姓名、其他位置／職業、2倍中文及整體交付仍待；沒有原版、字型或輸入檔頭修改。

## R992 名單格式的結構來源（2026-10-09）

<a id="r992-roster-format-source"></a>

入口[R992](../goal/R992-roster-format-source.md)、[022](../spec/022-character-roster-text.md)。R991已量四配對保持，不按未知位置繼續追加特例。
舊`workplace/ida/DS.EXE.i64`已缺失，從同SHA原版EXE重建一次性DB。既有py312-v1 image沿用，UID1000／IDA9.4最小JSON、輸入SHA及139函式核對通過。
原版DS.EXE SHA-256仍為`14a94f85cb4ef4df08566ed165dddc88c363aa816360c6509b84f079866a43db`；資料庫及code只在ignored工作區。

resident IDA MZ EA `[10F38,10FDE)`保留原名`sub_10F38`及operand。near call EA10FC9 bytes e8c928到13895，返回10FCC，
與既有runtime0110:0fcc印字錨點吻合。輸出近指標低16位是`0E38h+14h*([bp+0Ch]+[bp-2])`；
傳座標兩word是`[bp+4]+2`、`[bp+6]+6*[bp-2]+1`。分支選取對應1／0、一般195／180的解讀沿既有正常印字證據。
20-byte步距是顯示文字slot，不能稱PCFILE人物record。resident獨立分級證據FINAL 0／0／0完整回讀、五來源SHA相同。
JSON／Markdown SHA-256：`bfa84b535aea21f2b0675f1e9e545c88fa108834cbd6ef4e531e8e1c817ec87c`／
`8cdf6dcc9a0b4abed8323d0eda3da415196fcbff74592645072d8a1bedbe0583`。

VBASE原OVR SHA-256 `7919a6eddd0f7b65905ad4ab3a301604cccb8bf0f213ebebc190d5d1227104d6`，code段SHA
`f2af41d5bd0e394f7fab514a4ace8bd0ba8de718efcca5e544879a781fe68e87`。以下數值為raw code offset，原檔offset另加14。
raw5F1A有直接call入口錨點；24394起的連續寫端將來源字串搬到每20-byte slot，補空白到10格，
接著在+10寫參數140系列訊息首byte、+11寫2Dh、+12..14寫120系列訊息前三byte、+15寫NUL。
copy／length callee已另窄量resident143AB／14AD7，MSG getter0582沿034已證契約；不讀MSG或PCFILE內容。
VBASE raw查詢首版未建立指令項即print_operand，operand無效，保留原輸出；修正先建指令項後另存r2，不當產品或原版缺陷。
填列分級獨立FINAL 0／0／0完整採認，三匯出SHA相同；JSON／Markdown
`92f78b0c3f8ffca623e192da616bfdb863f1e1fea1b31715609ec494310236c0`／
`eb05cf481a1205d8198c315d80d89374457c6a28e77c4fafc49e5b5f5aea5800`。

新runtime探針先獨立審查、完整FINAL採認後執行，只在既有六個完整列入口重用16-byte字串，
新增每事件最多30 bytes數值frame、一次166＋325 bytes code簽章，cap64；沒有原版、離屏、A000、規則或輸入注入寫入。
42份正常名單frame均符合指標及座標算式，width15、count5／6、First0，DS／SS為1106h；
來源offset是0E38h+20*(First+Index)，祖先返回B092連到raw603D的建構端。
resident code166 bytes及VBASE建構端325 bytes逐byte等靜態輸入，實證此路線的raw→runtime加值5055h。
兩code SHA-256：`34ad5b07cf59b4c779d5ba8e980b33978281040d590e9a7aedaf262cbe127e3f`／
`12eb7a9cf85d8423d5c7a021af2df23e04dd97cae3d4b4386dd02ccaeb2a511d`。
1,546發布原版六類摘要及完整trace同R991無覆繪控制b63495，保存bytes／mode／mtime保持。
安全摘要`workplace/r992/structure-runtime-safe.json`。這些只證明正常顯示列生成與來源身份，任意姓名及其他職業仍須正常路線與022增補閘門。

runtime獨立FINAL 0／0／0已完整回讀採認，七來源及42份frame指紋核對相同。
JSON／Markdown SHA-256：`986eaacbe33ad5d9a0fb2f962a54d95da009fbb89917d42de87a75fd473e9b14`／
`c923be6b808286ea58dd1981fce558344f0f42d0331df8eb67a94a35f32f7deb`。
42份frame的合併指紋為`2d27a2134c1f1e7cf7b2b8f868896a93afd3079489a2b81d2af529d4026836f8`。
祖先第二層的ReturnIP B092h是回到建構端；建構端frame為其SavedBP所指第三層BP64808，不能把第二層BP64084當作建構端。
022新增數值守門DRAFT；尚未放行production、READY或任意姓名CONFORMED。

另從既有`text/messages.json`只解碼120..133的英文職業，遇134即停止，讀取檔案前綴1,980 bytes。
14個前三字母碼為FIG／MAG／PRI／THI／RAN／ALC／BAR／PSI／VAL／BIS／LOR／SAM／MON／NIN，全部唯一。
此為既存抽取文本的窄讀，不新增原版MSG.DBS內容觀察；來源及計數在`workplace/r992/class-code-sources.json`。
中文沿`docs/names.md`既有職業表，正式roster仍只有原六碼，沒有把未驗職業或性別資格推成遊戲規則。

### 正常新姓名與來源分流

正式55aa binary、倚天3倍、R982起始copy2，正常以改名將THESUS追加R為THESUSR。
Return結束編輯後按原版確定，人物名牌及重新開ADD的第一列均顯示新姓名，白底選取及取消也已量到。
沒有創角、重擲、刪角、新遊戲、防拷、戰鬥或SAVE；原資料唯讀，改名由原版自然寫入scratch。
正常錄製3,397事件／1,689快照，末Step395,800,000，頁尾Stop、續玩及關窗0，無partial。
新wire SHA-256：`6d4bed0a2382e2b3259029750deb7dec81ba5030edee5001eb29583a359072bd`。
正常玩家FINAL完整回讀採認，五來源SHA核同；JSON／Markdown為
`5587014e97dece9edd9a2a193686855c0f8b70032f423e85ac0bdd07352c5c53`／
`9dd06dd4d8fee9c177af66dc1a87931c0b443a9bb7557ec0c22f498ced849cef`。

首輪未確認改名及未正常停止的partial保留；控制器300秒逾時後關窗的故障分類見WORKLOG。
退格只送IRQ1掃描碼、未送姓名AH07的DOS字元通道，已另登記[#46](https://github.com/wicanr2/wizardry7-cht/issues/46)，尚未修改正式前端。
025原42frame與022數值DRAFT的獨立複核已採認；祖先第一層的16／26份分流及已修正阻擋見WORKLOG。
該DRAFT複核JSON／Markdown為`9da3032f41a402ff39eb578e3307602ad23839a0301e6486cc47ba68316249a6`／
`ddcf8d7ac70dc469c186643400cd8de6a07f8e1841030a245765d34eb3abdf12`，不是READY放行。

新wire共捕捉24份frame：12份改名選人，祖先5C2A／BA64／BE6C及61D9／BA64／BE6C各6份；
12份ADD，祖先5775／B092／B0CB及5D24／B092／B0CB各6份。
只有ADD第一列兩次分類為未知姓名，y150／Index0，一般SelectedFFFFh及選取Selected0各一次；數值公式均保持。
首個probe在改名選人先取到AF6F code，SHA為`32d7f3d24617df671f43648e69bb06b30fc9058f29bece4b78d3d697a7672d73`，不等ADD已證簽章。
保留此觀察，不把它改成ADD期望；另經獨立source FINAL採認的私人probe限定新姓名及ADD祖先成立後首取原491 bytes。
補證兩code逐bytes等原resident166及VBASE325，SHA仍為上列34ad5b／12eb7a，沒有擴大ADD身份集合。

主驗的新probe／無覆繪原API之1,689發布六類摘要及完整602,552-byte trace全同，SHA-256
`dbe4fe3cceab3cbd270d0cfc8ac64771cb3fbab158b244bde30eb33977860589`。
同55aa正式重播完成0。正常錄製、raw、新probe及正式重播的四份PCFILE bytes／mode相同，
SHA-256 `b3c137ae1ea74986fcd012ade6df797b156533215a06729219b76d89a4e936a9`；mtime隨各趟正常改名重寫，不要求相等。
四份SAVEGAME bytes／mode／mtime保持R982起始。全部保存只SHA／stat，不解人物DB、不撤銷合法新名。
安全結果在`workplace/r992/name-add-control-safe.json`。
新姓名目前仍顯示原碼M-FIG；這是辨識限制的正常測例，不能稱新姓名已中文化或022通用路徑CONFORMED。

動態獨立FINAL 0／0／0已完整回讀採認，29項來源／結果／保存身份核對相同。
新舊24份數值收據逐bytes相同，指紋為`35e429b7bbb728f58a04103bd00b63a961b42624dc0f5b6d3e0a4fd8c2905bd6`。
獨立另從原resident匯出及VBASE.code的24346起精確325-byte切片核code，沒有使用被測checker結果當期望。
JSON／Markdown SHA-256：`07272b1b8d6d0505e64ed6e30529c27d351e6c831006982d88cd6a04fab3cccd`／
`65de50022464ef217d8e5dd239d25688ece151e3ce5d5ee5824f40553a87ab74`。
本輪結束RE，不追改名建構端、人物DB或共用runtime；022保持DRAFT，下一個實作為已定位的退格輸入缺口。

## R993 姓名輸入退格鍵（2026-10-09）

<a id="r993-backspace-input"></a>

入口[R993](../goal/R993-backspace-stdin.md)、[035 DRAFT](../spec/035-backspace-input.md)、[#46](https://github.com/wicanr2/wizardry7-cht/issues/46)。
沿R992已證缺口，不新增原版記憶體讀取或寫入。私人原型main恰增五行Backspace JustPressed→既有鎖內append byte08，
F1早退、原IRQ1及owner的Type／錄放type路徑保持。原型source獨立FINAL 0／0／0完整採認，九來源及兩報告／binary共12項身份核同。
JSON／Markdown SHA-256：`e16990babb3a269d4ab9a8f392bc3d60c4431309bbbf5a7f5c4042e6b9b759fa`／
`6d9157db155a3d21cf9ce7d1deae6403386684d658945735b8c30978a3fe4416`。
原型Go1.24.13、93來源及binary d6fe340c…沿固定工具重建，source清單逐項size／SHA核同。

倚天3倍正常改名以R992合法THESUSR保存copy2為起點，先在F1開啟時按住／放開退格及Q。
關F1實圖姓名仍THESUSR，正常首退格顯示THESUS，再逐字至空欄；空欄額外退格後仍空。
逐鍵輸入ROWTEST，Return結束編輯後按原版確定，人物名牌與重新開ADD之新名、白底選取及取消均已目視。
全程無創角、新遊戲、防拷、戰鬥或SAVE；原版／初始保存唯讀，只由原版正常改名寫入本輪scratch。
頁尾Stop後續玩、正常關窗0。完整wire為4,441事件／2,190快照，末490.2M、962,345 bytes，SHA-256
`b5ee4abd64ff6d5a91ab8edd15c4843f8829a85a2bd079b87543af0e88e1fb2d`。
只對本次已知正常路線計數：type byte08共8、Q0；IRQ1 0Eh／8Eh各8、Q掃描碼0。
安全結果在`workplace/r993/prototype-input-safe.json`，F1 ack與所有畫面只留本機。
正常玩家FINAL、獨立READY閘門及正式接線／原版控制見下段。

正常玩家FINAL已完整採認，六來源／報告SHA、13checkpoint及兩F1 ack核同；JSON／Markdown為
`0e9293dbe9bff926762eea039f9fe693e5c4fc41107bdcb5c133194a68a07bea`／
`946bc5b0b743498b07e59f8853e8a96d71998cd709866e11e65e4f73158c98c9`。
035獨立證據與契約FINAL 0／0／0完整回讀，24項身份核同後採認READY，報告SHA及實作界限見035。
正式main只增五行，public重建93來源及binary與原型逐bytes相同；正式binary仍D6，完整前端身份不含檔名或mtime。
因此本次正常倚天錄製與同D6重播按READY契約納入正式範圍，來源身份證據在`formal-build/public-prototype-identity.json`。
正式cmd回歸2.356秒通過，initconfig及recording為有效快取；沒有為此重跑不變來源或額外擴測。

原API有／無覆繪各2,190發布的原版六類摘要及完整781,910-byte trace全同，SHA-256
`df22e708cd127bd3ae367263273c385d1f53966055de04634d697cd3b51d3824`。
正常／raw／translated／同D6重播四PCFILE bytes／mode同`b3e95d390e9a5f44e330729b22c00ad98148d3723e4d1ac5f3675fc7f056a269`，
mtime各趟正常改名重寫；四SAVEGAME bytes／mode／mtime保持。原APIrunner沿既有R991來源，replay Update早退於實體鍵盤處理前，
所以上述對照以完整已錄輸入派送；不把舊runner稱正式binary相容重播。安全結果在`prototype-control-safe.json`。
正式Unifont正常一次退格ROWTEST→ROWTES，T復原後原版確定，直接返回角色選單，依實圖省略多餘Return。
ADD的新名及其餘五列中文、白底選取及取消已量；3,167事件／1,574快照、末367M，700,405 bytes，
wire SHA-256 `97b47764dcb8a3639649bd0e27cced0092123b89efaf0bbfe2032d296dbef785`。
Stop後續玩、關窗0，同正式D6重播0；正常FINAL完整採認，六來源／報告及九圖核同。
JSON／Markdown：`6bb5a860b7e64078c728c9d08f11e86930c82b0e32b85eecbc7f5792d8d61310`／
`3f6a1b8f0896d5ba79930fe92e95729b90c1075c23a24410408735cb3f3d7841`。
Unifont各原API有／無覆繪1,574發布及561,382-byte trace全同，SHA-256
`5cb385371fb06c2f7b1ed7e615691aba50bb68ac3a27dddb49b907c64944e2d1`；四份PCFILE及SAVEGAME全部bytes／mode／mtime保持。
兩份不同正常輸入各自比較，合計3,764對發布，不把兩組不同輸入稱同狀態互拍。安全結果`final-control-safe.json`。

為補姓名墨點保持，開啟既有R991唯讀frame helper，不改runner或wire，兩份抽樣trace仍與既有控制全bytes同。
按控制器command4實際姓名框開啟及原版確定區間選圖，倚天七個、Unifont兩個唯一Step；
原DOS姓名域`[150,210)×[175,181)`整數3倍後29,160個RGB像素與衍生HD相同，差異0。
主驗直接展開原indexed／palette，以ImageMagick解HD；獨立代理另以Pillow NEAREST3重建，沒有用被測checker輸出當期望。
這些是衍生HD姓名域抽樣，不代替兩份已採認正常GUI，也不稱live Draw全圖比對。安全結果`name-pixel-safe.json`。

最終獨立FINAL 0／0／0完整回讀採認，75項來源／結果／36影格／wire／binary身份核同，035只增已量Linux姓名路線限定CONFORMED。
JSON／Markdown SHA及未量限制見035；36影格合併指紋
`c2353c70fd19033282d930c1fc0ae78b239dcbbef2d59f0d3d3641d2711229ce`。
新姓名職業碼仍英文，022通用辨識及其他中文化／平台交付保持未完成；本修補未修改SDK、譯文或字型。

## R994 通用名單證據（2026-10-09）

<a id="r994-generic-roster"></a>

入口[R994](../goal/R994-generic-roster.md)、[022 READY](../spec/022-character-roster-text.md#r994-通用顯示列增補草案)。
正式D6正常依序加入ROWTEST、TEMPEST、LYSAND、NOBAL、TREON、PENTAS，剩6至1各未選／白底共12張實圖。
最後一人加入後列表撤除、六人隊伍可見，原版移除ADD並切換滿隊主選單；沒有空列表印字或empty提示，不假造0-row路線。
新正常錄製3,455事件／1,720快照，末419.4M、760,850 bytes，Stop後續玩及關窗0。
wire SHA-256 `36ca61b9033942269bf10365b0ca6f559db3eb093c1396a2952b5973b82a6fe5`。
正常玩家FINAL完整回讀採認，六來源／報告及13實圖共19項核同；JSON／Markdown
`9bca86547171217b3d74166662926fbb305bb9f929e6bf20506abd62a970d117`／
`b9eb314016c9284ca43f3f49b830f53b039f9f8f489d6f77dc210fa6f7897108`。

R992改名時arena325 bytes窄比對已精確等VPCMK.code raw`[5F1A,605F)`。
VPCMK.OVR SHA-256 `dd1bcb9a54943163237ff48644c021ea202c36c15858fd7fd71520a5f871d28b`；
code SHA-256 `f48ba588d545f0bc71bdc07715b4b3552456fe85287228a8d36a243894992b12`。
此為來源身份，不把改名32d7簽章換成ADD的12eb，也不重開人物record格式或其餘函式研究。

新私人probe source獨立FINAL 0／0／0完整採認，七source與兩報告共9項核同後執行。
每事件重用raw16，frame最多30B、code491B逐次核兩來源簽章，before／after各64K及原返回XY4B，cap64與SP＋2／1M步限制保持。
66個返回為ADD count6..1各12／10／8／6／4／2，共42；另既存改名／ADD路線24，兩owner各12。
全66個返回X240／原y／SP＋2，完整64K差異都只有名字十格及碼五格；名字／碼非空、域外0，單一195→180或1→0。
安全細表在`workplace/r994/row-delta-safe.json`，這是原版有界格式證據，尚未放行正式通用覆繪。

剩餘人數路線probe／raw全部1,720發布與613,650-byte trace相同，SHA-256
`3869627b5653393296bf58e1bce04dd9a32e17c81971f89bbc606300d8e6919b`；同D6正式重播0，四保存層全部metadata保持。
改名probe的2,190發布全trace同R993原API，SHA仍df22e708…；保存由原版正常改名，未撤銷或改輸入。
安全結果`trace-preservation-safe.json`，動態獨立採認與READY契約另記。

私人資源原型依既存120..133十四職業補M/F共28碼，既有六碼的en／zh／dy／glyph全部相同。
兩字型缺字0，原生15／16，最長83／88px≤90；28個所需漢字都在正式U+碼點清單，沒有新增全局字元。
初版檢查器把U+hex清單當正文，誤報28字缺；保留初版，依正式格式另存`resource-prototype-safe-r2.json`，缺0。
原型TSV SHA-256 `4229a1e2cdfbb230e0442a639826c69c7bce9bfcbaf48d7c54a0a58cbc431099`，尚未搬進正式來源。

滿隊實圖另見兩項後續覆蓋缺口：主選單已有中文但回退普通小字；PENTAS側欄仍英文ALCHEMI。
正常raw已量ALCHEMI於step407118987、x261／y146、ret0110:1c71；滿隊START NEW GAME於ret0110:539b。
後續按原定主選單UKai與職業譯名處理，不因這兩項擴張本輪原版資料讀取；完整畫面與log只留本機。

獨立證據／資源／契約審查0／0／0已完整採認，17來源、兩報告與198證據指紋均相同。
022 R994已READY；原生男／女／鍊字模兩後端獨立抽查相符，First非0只放行已證公式的有界實作，正常未量。
報告身份與限制見022；原型source及正式動態驗收另記。本輪UI缺口留言6077357043全文283字／6個LF已回讀相同。

通用Go source獨立FINAL 0／0／0完整採認，12來源、兩報告及原型binary共15身份核同後正式接線。
回歸cmd／initconfig／recording分別3.185／0.029／1.175秒通過。正式與已審原型94份實源收據及binary全bytes同，
source SHA-256 `a6297a068f54caf0888716c64ce81c1ba041cb3ed1aa1cc003ba7fbd2596ea4d`；
binary `6f31cc079d38cad192157bd0ddc1146792f8fe4574369eb302e7f1a32f32899f`、15,989,192 bytes。
來源審查JSON／Markdown SHA為`33f0eef23565a665a5c6873946940d4c3308c2b921c28302b14777ca91646e90`／
`fba3e01f35677e0b39d404d4e09fc1f5400f24a42a7bb171023d04274ede213a`。

正式兩字型3倍正常GUI皆已完成SELECT六列／ROWTEST選取／取消、ADD count6..1及滿隊列表撤除，
正常頁尾Stop、續玩與關窗0。兩原字型另獨立重建SELECT及ADD白底碼帶，每後端7張，共14實圖、22,680 RGBpx差異0。
此為實際GUI的選取碼帶抽樣，不當作全圖HD／Draw同狀態；姓名保持另由正常視讀及原API完整像素檢查承接。
新完整wire為倚天3,248事件／1,614快照／末427.2M，Unifont3,196／1,588／422M；同6f31正式binary重播皆完成0。
兩GUI及兩正式重播保存bytes／SHA／mode／mtime保持，正常continued圖才確認錄製完成；早一張倚天圖仍在收尾，不單獨當成功。
GUI FINAL兩報告及8身份完整採認；80實圖指紋`b88588097a4455080dde00b91198461cacc7e31340fd0f79406536dc1101aaf9`。

六組私人原API old/new共12趟全0，12,200發布配對全trace同基準；remaining各1,720同3869627b…，rename各2,190同df22e708…。
所有抽樣indexed／palette同，remaining保存全metadata保持，rename原版正常寫入的最終bytes／mode相同，自然mtime另列。
80派生HD共37,760,000 RGBpx符合外部原字型期望；3倍54圖31,104,000px，2倍26圖6,656,000px保持old。
原核對器把游標覆在名字的10／33／2點當身分不符；改以完整英文碼墨點辨識，最後完整RGB比較仍護姓名及碼帶外。
完整六組default矩陣通過，3倍診斷結果不取代它；兩次修正皆為驗收工具，不改production或binary。

最終動態獨立FINAL 0／0／0完整採認，25來源／結果／報告及640畫格檔fingerprint核同；身份見022末節。
只標First0、兩owner、count6..1、六碼及新姓名的雙字型3倍限定CONFORMED。其他22碼、First非0、OEM／空白正常未量；
2倍仍非中文，派生HD不稱live Draw，全案及滿隊選單／ALCHEMI兩個後續缺口仍待。
