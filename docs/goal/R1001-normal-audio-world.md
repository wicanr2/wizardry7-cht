# R1001：有聲的正常載入與世界操作

2026-10-10，起點c42bb28；沿[#18](https://github.com/wicanr2/wizardry7-cht/issues/18)、[009 READY](../spec/009-player-audio.md)、[音訊證據020](../re/020-player-audio.md)。

上一輪R1000完成手持兩物品中文覆繪、正常錄放及推送，屬進展。R956已量角色選單六次正常點擊及軟體monitor，不重做同一範圍。本輪補目前正式版本從冷LOAD進世界的音訊玩家鏈，已有合法R990保存，不使用state注入、作弊或改原版時鐘。

1. 先沿既有專用音訊image與可重現來源核工具鏈。舊`eob-audio-capture:20260922-r2`現已不存在；先查現用image及其他既有音訊image，不退回主機啟動播放器／PulseAudio。
2. 使用正式96來源c7da執行檔、R995字模／R983 flow／R979 help。不明示`-mute`或`-ips`，沿009預設時鐘；冷啟動後正常LOAD、目前題後提示正常作答、原版載入、世界轉向／前進、正常SAVE／重開冷LOAD及F1暫停恢復，逐畫面確認再輸入。保存入口沿R990已驗的原版磁片選單，不開未知設定頁；只寫本輪scratch，不修改初始保存。已見世界磁片選單沒有LOAD，本次重開走原版載入，不新增入口；各程序由Step0冷安裝播放器，程序內自然LOAD／SAVE不重設時間軸。
3. 以容器內軟體monitor取得實際裝置輸出，保留時間標記、正常操作收據、非靜音、左右聲道、Reader／來源計數及正常退出。監聽輸出不替代實體裝置或人耳；可供本機試聽的節錄只留workplace，不公開原音。
4. 原版素材與初始保存唯讀，正常遊戲副本可寫。音訊不重設原版時間軸、不改規則；保存／載入格式不變。環境或擷取失敗先修工具，在同一環境乾淨重跑，不把裝置缺件寫成產品缺陷。
5. 限定結果由唯讀子代理審查，更新020、CONTEXT、WORKLOG及#18。未量人耳、平台、手動加速／退出分支及正式封包仍保持；不以單趟監聽升格全009。

所有工作在有界Docker內，網路none、UID1000、相稱資源與log限制，只寫`workplace/r1001`及已授權cache。新工具及輸出入口掛本檔，完整音訊、輸入、保存與畫格不提交。其他譯名待答不代選，十一組、TSV、地圖／作弊與完整交付仍在整體目標內。

## 工具入口

[音訊擷取包裝](../../tools/re/run_normal_audio_gui.sh)及[正常GUI控制器](../../tools/re/normal_audio_gui.py)是本輪兩趟實跑的原樣保存。控制器固定已驗c7da執行檔與R995／R983／R979資源，只送正常X11輸入，沒有私人hook。從正式65044e1來源依既有`build_recording.py`重建即可取得該96來源身份；不把新binary冒稱相同錄放版本。

於專案Docker建立新輸出目錄，把兩工具分別複製為`/out/run_audio_gui.sh`及`/out/normal_audio.py`。沿[有界啟動配置](../../tools/re/normal_audio_launch.json)的r3 image及完整ID、UID1000、4CPU／2GiB／160pids、network none、--rm與log上限；`/src`、`/original`及合法`/initial`均唯讀，僅`/out`可寫。重跑時換新的容器名及輸出目錄，不覆寫既有收據；初始保存路徑須是合法正常取得的副本。使用外層`timeout -k 20s 910s sh /out/run_audio_gui.sh`，不能省略硬截止。初始保存正常copy2到scratch，輸入命令寫入本機`command-N.json`後才會送入遊戲；每次先看上一張PNG，不提供答案表或固定答案。

PulseAudio16.1只開容器Unix socket；882-frame人工零音在monitor前預啟動，不作原版聲音樣本。擷取為44100Hz、signed16LE雙聲道，原音及手冊輸入皆只留workplace。檔案大小、Reader輸出與原版時間軸分別核對，不把F1期間monitor的靜音幀算成Reader交付幀。兩趟聲音／數值與保存恢復收據在本輪工作區，不放進Git。

[獨立音訊核對器](../../tools/re/check_normal_audio.py)直接解析安全計數與PCM，不匯入受測實作、不讀DB或完整輸入。沿相同有界Docker與唯讀`/src`，新的可寫`/out`執行：

```sh
python3 /src/tools/re/check_normal_audio.py --first /src/workplace/r1001 --second /src/workplace/r1001/coldload --spec /src/docs/spec/009-player-audio.md --reader-source /src/cmd/wiz7/audio_reader.go --stream-source /src/worktrees/dosgolem/opl2/stream.go --viewport-receipt /src/workplace/r1001/viewport-restore.json --out /out/audio-check-new.json
```

外層`timeout 180s`，只接受新的輸出名稱，不覆寫原收據。工具獨立核ceil、事件／幀／byte守恆、全部raw樣本、F1同Step與關窗時段；操作音訊區間沿收集塊標記，不當精確原操作的聲音邊界。

控制器的命令格式例如`{"label":"turn-right","click":[185,168],"wait":3}`，座標為DOS原座標；只能在已確認世界方向控制畫面時使用此例。`{"finish":true}`會恢復後正常關窗。其他操作可用`key`、`move`或目前畫面合法的`text`；不把手冊答案寫入範例或版本庫。

## 結果

兩趟正式c7da、預設時鐘、無mute／ips覆寫的正常有聲鏈已量。第一趟冷LOAD、右轉、前進、SAVE及保存後轉向；第二趟從該新SAVE重新啟動LOAD並續行，皆WM關窗0。各9／4個正常操作收集區間非靜音，15對F1暫停恢復同Step，關窗至程序退出約65.776／164.687毫秒。

獨立164項通過，兩趟canonical為5390519／4503843幀，原版來源事件96443／50343全接收施加、容量Dropped0。raw共19160008／9933566幀，非零5175426／4212968幀，左右差0；Reader Overflow12092／48509及停止捨尾68566／78188保留，不調參歸零。第二趟原版最後Step307M比音訊EndStep306.8M多一個200k停止邊界，不把這段當已合成。

新SAVE為原版正常寫入，初始保存保持，冷LOAD後bytes／mode／mtime保持；保存與恢復的世界viewport177408像素差0，未存轉向反對照162504差。沒有新RAM同狀態A/B、手動加速或全退出排程聲明。15秒本機試聽片段已交使用者，人耳回覆仍待；009及#18維持READY／OPEN，原生平台及其他交付條件不外推。

最終獨立JSON `workplace/r1001/audio-check-final-v2.json` SHA-256 `4b266c34ecf4c452e444d8a9806f2df4f14d4f3bbcd62bfcbb7debc8e53ed790`。首版PASS只訂正可選viewport的限制文字，驗收條件及全部runs／viewport結果保持；沒有重跑原版。正式Go／SDK／譯文未改，完整聲音、輸入與保存只留本機。

精確READY輸入保存在本機`workplace/r1001/spec009-input.txt`，來源為`65044e1:docs/spec/009-player-audio.md`，SHA及用途由同目錄`contract-input.json`索引。重現原契約來源時用該快照，現行009另追加本輪結果；不以兩份文件雜湊不同宣稱程式漂移。
