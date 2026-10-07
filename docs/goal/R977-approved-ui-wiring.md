# R977 定版字樣的正式接線與驗收

2026-10-08；沿[#22](https://github.com/wicanr2/wizardry7-cht/issues/22)、
[#28](https://github.com/wicanr2/wizardry7-cht/issues/28)、[020](../spec/020-main-menu-labels.md)
及[025](../spec/025-intro-picture-labels.md)。

## 已定版

使用者確認R976米金亮字、紅字與六按鈕毛筆楷體，並要求繼續正式接線與驗收。
紙頁維持甲的30／32字高、水平加粗1像素及R975排版，亮度採R976；
主選單採AR PL UKai TW字面、2倍16／3倍24像素及樣張置中。
已確認事項不重問，不由這次繼續外推未展示的2倍紙頁排版或其他選單字型。

## 工作

1. 恢復既有Docker工具鏈的可重現建置與離線模組，維持Go1.24.13。
2. 把定版圖與參數固定為本機驗收基準，完善020毛筆及025紙頁增補的資料、
   來源辨識、共同HD／Draw、回退與撤除契約，獨立唯讀審查READY後才改正式前端。
3. 主選單2倍／3倍接線；紙頁先依已展示3倍定版，2倍可排版性另核，不縮小或裁切大字。
   支援範圍與未量限制逐項記錄，不稱全遊戲完成。
4. 合法正常無輸入冷啟動，固定同狀態新舊對照、原版摘要／indexed／palette及批准像素區；
   正常GUI、倍率切換、撤除與退出抽樣，完整cmd/wiz7回歸與資源拒絕案例。
5. 審查採認後push、更新Issue並全文回讀，核Docker與輸出擁有權。

原版只讀，不讀手冊／答案，不改時基、遊戲規則、存檔或正常輸入。
R976圖片及字模含本機輸入，保留本機；程式來源、規格與重跑腳本進版控。
字型素材依原授權，程序仍RRSAL。其他交付門檻保持。

## 入口

[R976定版來源](R976-bright-paper-brush-menu.md)、[研究025](../re/025-ui-localization-coverage.md)、
[Docker建置來源](../../tools/docker/go-ebiten.Dockerfile)、[字型入口](../../font/README.md)。
研究輸出只放`workplace/r977/`；UID/GID1000:1000、離線、資源上限、有界容器與log上限。
來源／原版／字型唯讀，只有明確輸出或工作樹可寫，名稱前綴`wiz7-r977-`。
保留未追蹤R433，不動原始dosgolem工作區或其main。
資源邊界預審入口：`workplace/r977/review-boundary.json`及`.md`，僅作本輪DRAFT預審。
完整增補入口：[028](../spec/028-approved-ui-fonts.md)，已完成獨立READY審查。暫名027的歷史報告沿SHA回查；改名只消除與既有027的編號碰撞。

## 離線工具鏈輸入

既有image缺失後沿原Dockerfile的recovered-gui target恢復，名稱為
`wizardry7-go-ebiten:go1.24.13-20261008-gui-r1`。固定Go1.24.13，Debian12，未恢復Pulse ALSA plugin。
執行前核`docker image inspect`，不可把其他同名映像當這次來源。

| 來源標籤 | 本機image ID |
|---|---|
| golang:1.24.13-bookworm | sha256:1a6d4452c65dea36aac2e2d606b01b4a029ec90cc1ae53890540ce6173ea77ac |
| eob-remake-go:1.26.7-ebiten2.9.9-20261008-r2 | sha256:9b6b1a2a457eb6937aceac5aa0eea79de88227100e33163644e1f8f4e661949c |
| colonization-verification:20261007-runtime-r2 | sha256:c3457c559d90724a99734edb23f7c4788bbe2a37e39420eb6edc1e0f5dc3abf3 |

重建入口為`docker build --network none --pull=false --target recovered-gui -f tools/docker/go-ebiten.Dockerfile -t wizardry7-go-ebiten:go1.24.13-20261008-gui-r1 tools/docker`。
三個來源ID先核相符；缺來源即停止，不聯網拉取或更換版本。建置context只含tools/docker，不含原版與字型。
完成映像config SHA為17586d1c45eec97186081e0fb0903e3733bbfdd6e9d483e27178f39c434863b9；
BuildKit manifest為baa20f69bccc7dc8ac7a4bbcf93fbcb32019505bf1c31c1312f85aa873bbb9d4。
回歸直接以本標籤執行，Go模組採映像既有快取複製到指定workplace輸出，不掛主機runtime。

本輪正式來源：[baketool](../../tools/re/bake_approved_ui.py)、[無輸入冷啟動驗收](../../tools/re/run_approved_ui.py)、
[獨立像素核對](../../tools/re/check_approved_ui.py)、[共同解碼及毛筆](../../cmd/wiz7/approved_ui.go)、
[紙頁顯示](../../cmd/wiz7/intro_paper.go)、[固定幾何](../../cmd/wiz7/intro_paper_model.go)。

## 已採READY審查

完整獨立審查0／0／0，主代理完整回讀兩份報告並採認。
JSON SHA854b5964bd2bba83de6b488b818e8b0571f5f2c6528ba831d955922bee54780b；
Markdown SHAbb13f8a75320ee5aecbcc9cdd54f3b8cac6d60568bdff44cb8607c7572882bdf。
兩份報告保存於workplace/r977/review-spec，暫名027指本028當時的契約，不指既有027。

## 本機字型及原版輸入

以下目錄先核存在，全部唯讀掛載；只有本輪明確輸出可寫。

| 本機來源 | 容器用途 |
|---|---|
| /home/anr2/cht/wizardry7/workplace/original/wizar7 | 原版冷啟動，不讀MSG／手冊或答案 |
| /home/anr2/cht/golden_box/拯救地球/workplace/unifont-src | 固定Unifont17.0.05來源，WIZ7_UNIFONT_HEX指unifont_all-17.0.05.hex.gz |
| /home/anr2/cht/etan_font | 本機stdfont.15、ascfont.15及ET353S/FILES/SPCFONT.15 |
| /usr/share/fonts/truetype/arphic | 新主選單ukai.ttc，字面2 |
| /usr/share/doc/fonts-arphic-ukai | copyright授權原文 |

原版由玩家自備；字型原檔及烘製字模保留本機，不提交。SHA依028及007。

## 本輪結果

已完成定版字樣接線及限定驗收。兩後端、主選單2／3倍與紙頁3倍；35個同GUI原版狀態及31,360,000 HD／Draw像素通過，F1返回、無字頁撤除與正常關窗通過。
獨立實作／資源審查0／0／0已完整回讀採認。JSON SHAe8666813d3a899e843f5116e982540ca432ad3464bb67da880a17c31599af5cc，
Markdown SHAc3bc7967a8a48a4b72080ecf3d25a73ab4e5b6d8fec808629a08f8cac96a73e6。

最終Go1.24.13／Linux amd64建置89份來源，Source SHA40416d488a4105145d8fbed8b18463984187be512a8cb56528690d6a6cb33ba1；
本機workplace/r977/wiz7-final，Binary SHA9898b24eb157f44a42a351c63f818acd9976c6673fc7734b3c524d3ece03c030。
完整前端測試含真實畫格守門1.546s通過；最後Unifont3倍GUI再抽7點、8,064,000像素全同，WM_DELETE退出0。
最終像素收據pixel-final.json SHA30cd33c3353dfdaa4675d7ad1acd90b6d0a476ca74a681fc176771865c0a7432。

紙頁2倍維持原圖；補底為已採近似。沒有音訊／人耳、後期玩法、完整錄放或發行包驗收，#22／#28保持OPEN。
