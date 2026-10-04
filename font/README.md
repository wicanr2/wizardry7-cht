# 中文字型重建

本機可選倚天15×15；發行包採 GNU Unifont 17.0.05 原生16×16。
字元取正式譯文與F1說明，包含[text/roster.zh-TW.tsv](../text/roster.zh-TW.tsv)，不能為缺字改譯或裁掉字形筆畫。
本目錄只保存碼點清單與重建說明，不提交第三方字型或烘製資料。
圖片的「呈獻」取自[片頭譯文表](../text/intro-picture.zh-TW.tsv)，規格見[025](../docs/spec/025-intro-picture-labels.md)。

## 固定來源

[GNU Unifont 17.0.05](https://www.unifoundry.com/pub/unifont/unifont-17.0.05/) 的
`unifont_all-17.0.05.hex.gz`，SHA-256：
`7b182454966046d35482469b979edce7d262fab5c53c2180e9b1fbb5d0b5e574`。
來源字型、`COPYING`、`OFL-1.1.txt`與著作權說明保存於本機
`workplace/fonts/unifont-17.0.05/`；正式封包須附上授權與著作權。
發行字模採SIL Open Font License 1.1，第三方字型不受本專案RRSAL涵蓋。
尺寸與驗收契約見[規格007](../docs/spec/007-native-font.md)，
輸入與證據見[研究018](../docs/re/018-distributable-font.md)。

## 重建入口

所有命令須在既有`wizardry7-go-ebiten:latest` Docker容器內、專案根目錄執行。
容器資源、UID/GID、唯讀輸入與清理依[AGENTS.md第10節](../AGENTS.md)；
將上列本機來源目錄唯讀掛到`/font`，僅輸出目錄可寫。
以下輸出目錄須為本輪新目錄，不覆寫已驗收的收據：

```sh
export WIZ7_FONT=unifont
export WIZ7_UNIFONT_HEX=/font/unifont_all-17.0.05.hex.gz
python3 tools/re/bake_font.py workplace/bake-unifont
python3 tools/re/bake_flow.py workplace/bake-unifont/flow.json
python3 tools/re/bake_help.py workplace/bake-unifont/help.json
```

主程式以`-xlate`、`-flow`、`-help`指定這三份JSON；字型選項只影響烘製與顯示。
`bake_font.py`同時從職業碼TSV產生`roster.json`，主程式依`-xlate`的目錄載入它。
六個已量角色列的中文目前限互動與HD皆3倍，規格見[022](../docs/spec/022-character-roster-text.md)。
只重建此家族可用下列入口，不需原版訊息表；環境變數沿上列設定，輸出必須是新檔：

```sh
python3 tools/re/bake_roster.py /out/roster.json
```

`bake_font.py`也產生`intro-picture.json`，主程式依`-xlate`的目錄載入。
只重建片頭家族可用下列入口；環境變數沿上列設定，輸出須為新檔：

```sh
python3 tools/re/bake_intro_picture.py /out/intro-picture.json
```

片頭六張圖片的驗證入口為[check_intro_picture.py](../tools/re/check_intro_picture.py)。
需要英文及雙字型兩倍率的正常冷啟動收據、已選甲案原型和合法字型；輸出須為新檔。
參數及驗收範圍見[規格025](../docs/spec/025-intro-picture-labels.md)。

倚天後端可另傳`--eten-font /eten/stdfont.15`指定唯讀原字型。
六列獨立驗證入口為[check_character_roster.py](../tools/re/check_character_roster.py)，
需要本輪新舊程式的固定路線收據及合法原字型；輸出報告須是新檔。

```sh
python3 tools/re/check_character_roster.py workplace/r928/implementation \
  --unifont-hex /font/unifont_all-17.0.05.hex.gz \
  --eten-font /eten/stdfont.15 --output /out/new-roster-validation.json
```
倚天後端改用`WIZ7_FONT=eten`，並將使用者自備字型目錄唯讀掛到
容器`$HOME/cht/etan_font`；不將倚天加入可散布封包。

碼點清單[charset.txt](charset.txt)不含字形。重建時從新輸出取得清單，核對後
更新本檔索引；同時保留JSON探針收據以記錄來源字型雜湊與缺字結果：

```sh
python3 tools/re/probe_font.py /font/unifont_all-17.0.05.hex.gz \
  workplace/rNNN-font-probe.json --charset workplace/rNNN-charset.txt
```

獨立位圖驗證入口為[check_native_font.py](../tools/re/check_native_font.py)；
`--baked-only`只驗烘製資料，完整同狀態模式還需要研究018所列原版收據。
