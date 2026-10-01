# 巫術7 譯名表（集中管理；台灣既有譯名優先，來源待考）

| 原文 | 譯名 | 備註 |
|---|---|---|
| Cosmic Forge | 宇宙鍛爐 | Pen and Hand of Destiny＝命運之筆與手 |
| Cosmic Circle | 宇宙之環 | Lords of the Cosmic Circle＝宇宙之環諸主 |
| Astral Dominae | 星界統領 | DRAFT：官方攻略本譯名待考 |
| Dark Savant | 黑暗賢者 | 副標 Crusaders of the Dark Savant |
| Vi Domina | 薇·多米娜 | |
| Aletheides | 阿萊賽德斯 | machine-man＝機器人 |
| Guardia | 嘉迪亞 | 隱藏行星 |
| T'Rang | 特蘭格 | 勢力/種族 |
| Umpani | 恩帕尼 | 勢力/軍事組織 |
| Black Ship | 黑船 | ship of darkness |
| M-FIG / F-FIG | 戰士 | Fighter |
| M-THI / F-THI | 盜賊 | Thief |
| M-PRI / F-PRI | 牧師 | Priest |
| M-MAG / F-MAG | 法師 | Mage |
| M-ALC / F-ALC | 鍊金術士 | Alchemist |
| M-PSI / F-PSI | 心靈術士 | Psionic（待確認） |
| M-NIN / F-NIN | 忍者 | Ninja（待確認） |
| M-BAR / F-BAR | 吟遊詩人 | Bard（待確認） |
| THESUS | THESUS | 預設角色名保留 |
| TEMPEST | TEMPEST | 預設角色名保留 |
| LYSAND | LYSAND | 預設角色名保留 |
| NOBAL | NOBAL | 預設角色名保留 |
| TREON | TREON | 預設角色名保留 |
| PENTAS | PENTAS | 預設角色名保留 |

## 訊息表譯名批次（2026-09-26，來源：無官方中文化，我方定案，全標 DRAFT）

- **種族**（MSG id 100-110）：人類／精靈／矮人／侏儒／半身人／妖精／蜥蜴人／卓肯／貓人／狼人／穆克
- **職業**（120-133）：戰士／法師／牧師／盜賊／遊俠／鍊金術士／吟遊詩人／心靈術士／女武神／主教／領主／武士／武僧／忍者
- **屬性**（200-219、5050-5057）：力量／智力／虔誠／體質／敏捷／速度／人格／因果
- **咒語**（4000-4101，102 條）：意譯（火球／閃電／魔法飛彈／死亡之言……），欄寬限制取 2-5 字
- **職級表**（800-939）：14 職業 × 7 級（見習→征服者、魔術師→魔導士……祖師）；`*` 開頭＝終極道具名留原文
- **戰鬥**（3000-3957）：動詞（揮擊/突刺/重擊）、部位（背/翼/尾/頭……觸手）、狀態（麻痺/石化/重創）、施法結果（法術逆火！/被消音了！）
- **陷阱**（5300-5317）：魔法火花／毒鏢／警報／蛇妖之眼／三臂怪手……
- **物品類別**（6000-6127）：`?` 前綴疑為未鑑定標記，原樣保留（匕首／劍／護甲(上)……）
- **怪物武器**（3800-3859）：科幻系專名音譯（爆能彈／相位手槍／虹光手套……）

**不譯**（原因）：
- 2500-2574 手冊防拷關鍵字表——玩家要對實體英文手冊查字
- 7000-7197 對話解析關鍵字詞庫——引擎按玩家輸入的英文比對
- `*R.GARLAND` 等終極道具名——待物品系統（.VGA）統一處理
- 人名（D.W.Bradley 等 credits 人名）

## 怪物譯名批次（2026-09-27，SCENARIO.DBS 怪物資料庫，全標 DRAFT）

- **來源**：SCENARIO.DBS offset 225272 起、232 bytes／隻 × 250 隻、名字欄 +72
  ＝[真名單數][真名複數][別名單數][別名複數]（別名＝低鑑定度顯示名）。
  工具 `tools/gen_monsters.py` → `text/monsters.json`（原文）＋
  `text/monsters-trans.json`（461 筆譯文）→ bake 併入 xlate.json。
- **種族系譯名**：PHOOT 噗特系／MOTH 蛾系／WEEVIL 象鼻蟲系／ROOK 鴉系／
  CREEPING CRUD 爬行污穢系／THRAXE 巨蟲系／GIANT 巨人系／NYMPH 寧芙系／
  DRAGON 龍系（綠龍/黑龍/翡翠龍）／GOLEM 巨人像
- **勢力系譯名**：T'RANG 特蘭格／UMPANI 恩帕尼／DANE 丹恩／GORN 高恩／
  RATTKIN 鼠人／SAVANT 賢者／XENOZOID 異形
- **劇情角色**：DARK SAVANT 黑暗賢者、VI DOMINA 薇·多米娜、ALETHEIDES
  阿萊賽德斯（依既有譯名表）；DUNGORE 鄧戈爾、GENERAL YAMO 亞莫將軍、
  KING ULGAR 烏爾加王、H'JENN-RA 赫詹拉等音譯。
- **複數欄**：中文無複數變化——譯文同單數欄（key 四欄各自存在）。

## 物品譯名批次一（2026-09-27，SCENARIO.DBS 物品資料庫，DRAFT）

- **物品表**：SCENARIO.DBS offset 896 起、**72 bytes／件 × 417 件**、名字欄 +0
  （BROKEN/ITEM＝idx 0 佔位）。工具 `tools/gen_items.py` → `text/items.json`
  ＋`text/items-trans.json`（批次一＝武器/防具 idx 0-160、161 筆）→
  bake 併入 xlate.json（合計 1,536 串、字形 1,236）。
- **名格式**：`MATERIAL/TYPE`（LEATHER/CUIRASS 皮胸甲、CHAIN/HAUBERK
  鎖鏈甲）＝材質/類型組合；神器專名（EXCALIBER 王者之劍、MURAMASA/BLADE
  村正刀、THE/AVENGER 復仇者）。
- **注意**：譯名限 BIG5 字集（stdfont.15）——「脇差」→「脅差」
  （脇 U+8107 不在 BIG5）。

## 物品名稱校正（2026-10-01）

目前有效的28筆校正集中於[規格010的完整英文鍵表](spec/010-item-name-corrections.md)。
原文名稱與索引已逐筆回查SCENARIO.DBS；正式譯表及產生器同步，兩字型缺字0。
MACE＝鎚矛、STAFF／QUARTER/STAFF＝長杖、MORNING/STAR＝晨星錘、HAMMER＝錘；
完整採用名稱以上述表為準，不以舊批次的索引錯配作譯名依據。
SLING＝投索是欄寬內壓縮的強推論；WALRIBLADE＝瓦爾里刃依既定音譯規則，原版發音未知。
LT./CROSSBOW＝輕弩、CURE/STONE＝解石化藥水的術語連結為強推論，名稱不證明原版效果。
正常角色頁回歸沒有觸發這28個校正完整鍵，未宣稱全部名稱上屏。

## R7 翻譯規模化批次譯名（2026-09-27 起，subagent 批次沉澱，DRAFT）

**NPC 對話試作批（batch-005，9000 段）**：
貝爾坎佐（Belcanzor）、基馬斯·圖蘭（Kymas Turan，主題詞作「基馬斯圖蘭」）、
救世者（the Deliverer）、大丹恩（Magna Dane）、托倫·丹恩（Tollen Dane）、
阿爾馬戈特修士（Brother Almagorte）、魯萊神父（Father Rulae）、
馮德蘭教授（Professor Wunderland）、唐·巴隆內（Don Barlone，Don＝黑手黨式尊稱）、
鼠拉斯普丁（Ratsputin）、老盲麥斯（Old Blindmeis）、莫瑟修士（Brother Moser）、
羅丹·勒瓦克斯（Rodan Lewarx）、魁薩卡（Quasquet)、波利加德（Boerigard，
Gorn Army 陸軍 Captain——既有「船長」譯法有誤待修，見 013 追記）、
亞莫將軍（General Yamo，I.U.F.＝恩帕尼帝國軍統帥）、烏爾加王（King Ulgar）。

**地點/組織/概念**：查封區（Condemned Area）、聖道（Holy Path）、聖工（Holy Work）、
神聖之光（Divine Light）、秘密學院（Secret School）、森林家園（Forest Home）、
紅書（Red Book）、夢土（land of dreams）、僧侶聖地的宮殿（Palace of Munkharama）、
烏克皮爾山脈（Ukpyr Mountains）、高恩大軍（Gorn Army）、邪穢之塔（the Unholy Tower）、
召集之爪（the gathering paw）、格蘭梅蘭（Gran Melange）、瀆神（blasphemous）。

## R7 波一沉澱＋衝突裁決（2026-09-27，batch-001~008）

**裁決（後續批次一律依此）**：
- Xen Xheng＝**鄭**（上屏尊稱鄭大師；XHEN 同）｜TShober＝**喬伯**
- Torquesade＝**托克薩德**｜Mick the Pick＝**撬鎖米克**｜Bertie＝**貝提**
- ORKOGRE（城/種族）＝**奧高高**（8045 定案）｜DIONYSCEUS＝**狄奧尼修斯**
- RATTKIN RUINS＝**鼠人遺跡**｜Helazoid＝**赫詹族**（corpus 8042/8211/8218
  誤譯「異形」待全域掃描修正，與 XENOZOID 撞名）
- Boerigard＝**波利加德**，Captain of Gorn Army＝**上尉**（corpus「船長」誤譯待修）
- Blindmeis＝**盲麥斯**（全名老盲麥斯；Blienmeis 布林麥斯＝同一人簡拼，統一盲麥斯）

**波一新譯（DRAFT）**：Shritis 希里提斯、K'Borra 克波拉、D'Rang 德蘭、
Tracker Rhallick 拉利克、Lt. Gruntwrapper 格蘭特拉普中尉、Sgt. Kabomm 卡邦姆中士、
Lt. Gromo 格羅莫中尉、Sgt. Balbrak 巴布拉克中士、Dame Ke-Li 柯麗夫人、
Jan-Ette 珍妮特、Bela 貝拉、Hector 赫克特、HORATHA 霍拉莎（軍事機密話題）、
Ratskell 拉茨凱爾（鼠人商店兼盜賊公會）、Higardi/Higardia 希加爾迪（亞）、
Munkharam（先知）穆克哈拉姆（城名 MUNKHARAMA 仍作僧侶聖地）、Sogheim 索格海姆、
Brombadeg 布羅姆巴德格、Phoonzang 豐贊（赫詹族聖者）、HIGH FATHER 高階神父、
school of 5 Flowers 五花門、Hall of Preservation 保存殿、Helazoid Legion 赫詹軍團、
GREAT MAKER 偉大造物主、DIVINE PLAN 神聖計畫、WORD（丹恩教義）聖言、
HOLY WORD 聖言、Jonga spice 榮加香料、Fellowship of Dane 丹恩團契、
Holy Covenant 聖約、Isle of Crypt 墓窖之島、Witch Mtns 巫女山脈、
Dragon Mtns 龍之山脈、Great Wilderness 大荒野、Cosmic Forces 宇宙諸力、
BLUNDERBUSS 大口徑火槍、MAGIK 奇術（與 MAGIC 魔法分流）。

**控制碼定案（013 既有實證＋波一驗證）**：詞尾 `%` 一律保留（黑市%/聖工%），
機械計數檢查強制；`^N`＋英文後綴＝「刪前詞 N 字元」的比對鍵壓縮，zh 原樣保留
（顯示詞^N 後綴），引擎對中文顯示詞的刪除行為待端到端實測；`#` 句中＝填槽保留、
行首＝成對標記拿掉；`34}` 前綴（10475）語意待查。
