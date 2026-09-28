// wiz7 是《巫術7》的可遊玩前端（#15）：原版 DS.EXE 跑在 dosgolem 上，
// 開視窗、接鍵盤與滑鼠。
//
// 架構：機器在背景 goroutine 用大塊 Run 推進（probe 同款快路徑——
// 60fps 小塊 Run 在 credits→主選單轉場會卡死，docs/re/009），
// 前端執行緒只抓畫面快照與轉發輸入（鎖保護）。
// 速度以「每秒指令數」節流（預設 24M ≈ DOSBox-X cycles=20000 的體感：
// 開機到主選單約 30–40 秒）。
//
//	wiz7 -root workplace/run [-scale 3] [-ips 24000000]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wicanr2/dosgolem/oracle"
)

// scanCode 把前端會用到的鍵映射成 Set 1 掃描碼（make code）。
func scanCode(k ebiten.Key) (uint8, bool) {
	m := map[ebiten.Key]uint8{
		ebiten.KeyEscape:    0x01,
		ebiten.KeyEnter:     0x1C,
		ebiten.KeySpace:     0x39,
		ebiten.KeyUp:        0x48,
		ebiten.KeyLeft:      0x4B,
		ebiten.KeyRight:     0x4D,
		ebiten.KeyDown:      0x50,
		ebiten.KeyF1:        0x3B,
		ebiten.KeyF2:        0x3C,
		ebiten.KeyF3:        0x3D,
		ebiten.KeyF4:        0x3E,
		ebiten.KeyF5:        0x3F,
		ebiten.KeyA:         0x1E,
		ebiten.KeyB:         0x30,
		ebiten.KeyC:         0x2E,
		ebiten.KeyD:         0x20,
		ebiten.KeyE:         0x12,
		ebiten.KeyF:         0x21,
		ebiten.KeyG:         0x22,
		ebiten.KeyH:         0x23,
		ebiten.KeyI:         0x17,
		ebiten.KeyJ:         0x24,
		ebiten.KeyK:         0x25,
		ebiten.KeyL:         0x26,
		ebiten.KeyM:         0x32,
		ebiten.KeyN:         0x31,
		ebiten.KeyO:         0x18,
		ebiten.KeyP:         0x19,
		ebiten.KeyQ:         0x10,
		ebiten.KeyR:         0x13,
		ebiten.KeyS:         0x1F,
		ebiten.KeyT:         0x14,
		ebiten.KeyU:         0x16,
		ebiten.KeyV:         0x2F,
		ebiten.KeyW:         0x11,
		ebiten.KeyX:         0x2D,
		ebiten.KeyY:         0x15,
		ebiten.KeyZ:         0x2C,
		ebiten.KeyBackspace: 0x0E,
	}
	sc, ok := m[k]
	return sc, ok
}

// normalizeKey＝查表 key 正規化（與 bake_font.py 的 normalize 同一套）：
// 剝 <0x20 控制碼、連續數字→'$'、序數詞 FIRST..TENTH→'$'、'^'/'%' 佔位→'$'。
// token（數字/序數詞）另存，供 '$' 槽填回。
func normalizeKey(s string) (string, []string) {
	var b []byte
	var nums []string
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 {
			continue
		}
		if c >= '0' && c <= '9' {
			j := i
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			b = append(b, '$')
			nums = append(nums, s[i:j])
			i = j - 1
			continue
		}
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			j := i
			for j < len(s) && (s[j] >= 'A' && s[j] <= 'Z' || s[j] >= 'a' && s[j] <= 'z') {
				j++
			}
			word := strings.ToUpper(s[i:j])
			if _, ok := ordinalsPub[word]; ok {
				b = append(b, '$')
				nums = append(nums, word)
				i = j - 1
				continue
			}
		}
		if c == '^' || c == '%' {
			b = append(b, '$')
			i++
			continue
		}
		b = append(b, c)
	}
	return string(b), nums
}

// ordinalsPub＝序數詞→中文（與 xlate __ordinals__ 對應；查詞用）。
var ordinalsPub = map[string]bool{
	"FIRST": true, "SECOND": true, "THIRD": true, "FOURTH": true, "FIFTH": true,
	"SIXTH": true, "SEVENTH": true, "EIGHTH": true, "NINTH": true, "TENTH": true,
}

// leadingWordDollar＝^＝隊員名 pattern 的 variant key 產生器：
// runtime 字串開頭若是「全大寫詞」（隊員名，如 THESUS），換成 '$'。
// 混有大小寫的普通字（如 "Life"）不觸發。
func leadingWordDollar(s string) string {
	i := 0
	for i < len(s) && s[i] < 0x20 {
		i++
	}
	j := i
	for j < len(s) && s[j] >= 'A' && s[j] <= 'Z' {
		j++
	}
	if j-i < 2 || j >= len(s) && j == i {
		return ""
	}
	return "$" + s[j:]
}

// vramHash＝Indexed 色號畫面中事件區域（x..x+w、y..y+14，越界裁切）的
// FNV-1a。殘影清除：合成時背景區域雜湊與記錄時不同＝畫面已切換→丟棄。
func vramHash(idx []uint8, x, y, w int) uint32 {
	var h uint32 = 2166136261
	for row := y; row < y+15 && row < 200; row++ {
		if row < 0 {
			continue
		}
		x0, x1 := x, x+w
		if x0 < 0 {
			x0 = 0
		}
		if x1 > 320 {
			x1 = 320
		}
		base := row * 320
		for px := x0; px < x1; px++ {
			h = (h ^ uint32(idx[base+px])) * 16777619
		}
	}
	return h
}

// evAlive 過濾掉背景已變（畫面已切換）的事件。
func evAlive(evs []glyphEv, idx []uint8) []glyphEv {
	kept := evs[:0]
	for _, e := range evs {
		w := 0
		for _, g := range e.g {
			w += len(g) + 2
		}
		if vramHash(idx, e.x, e.y, w) == e.hash {
			kept = append(kept, e)
		}
	}
	return kept
}

// 前端 ↔ 機器 goroutine 的共享狀態（mu 保護）。
// entry 是譯文表的一筆（xlate.json）：dy＝行距微調、g＝每字 15×15 位圖。
// hdFrame 合成高解析幀：A0000 2x 放大＋文字事件（15×15 原字形）。
func hdFrame(o *oracle.Oracle, evs []glyphEv, path string) {
	pal := o.Palette()
	idx := o.Indexed()
	evs = evAlive(evs, idx)
	img := image.NewRGBA(image.Rect(0, 0, 640, 400))
	for py := 0; py < 400; py++ {
		for px := 0; px < 640; px++ {
			c := pal[idx[(py/2)*320+(px/2)]]
			img.Set(px, py, color.RGBA{c[0], c[1], c[2], 255})
		}
	}
	for _, e := range evs {
		fr := color.RGBA{e.fgR, e.fgG, e.fgB, 255}
		cx := e.x * 2
		for _, glyph := range e.g {
			for row, w := range glyph {
				for bit := 0; bit < len(glyph); bit++ {
					if (w>>(len(glyph)-1-bit))&1 != 0 {
						img.Set(cx+bit, e.y*2+row, fr)
					}
				}
			}
			cx += len(glyph) + 2
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	png.Encode(f, img)
}

type entry struct {
	Dy int        `json:"dy"`
	G  [][]uint16 `json:"g"`
	D  []int      `json:"d,omitempty"` // '$' 數字槽索引（g 內位置，位圖全 0）
}

// glyphEv 是一則待合成的中文字串（15×15 原字形、畫在放大後畫布）。
// Hash＝記錄當下背景（Indexed 色號）事件區域的 FNV——合成時背景已變
// （畫面切換/重繪）即丟棄，殘影清除不依賴 flip 訊號（013 追記三十三）。
type glyphEv struct {
	x, y          int
	fgR, fgG, fgB uint8
	g             [][]uint16
	hash          uint32
	img           *ebiten.Image // Draw 時 lazily 預渲染的字形條
}

type shared struct {
	mu      sync.Mutex
	rgba    []byte    // 320×200×4 的畫面快照（Draw 用）
	idx     []uint8   // 320×200 色號快照（事件殘影過濾用）
	keys    []uint8   // 待注入的掃描碼（make＋break 成對排程）
	mouseX  int
	mouseY  int
	buttons [2]bool   // 左、右
	evs     []glyphEv // 翻譯文字事件（機器 goroutine 生產、Draw 消費）
	stdin   []byte    // 待注入的 Stdin 字元（AH=3Fh 路徑）
	clicks  []struct{ x, y int } // 待執行的整包點擊（Click 路徑）
}

type game struct {
	sh     *shared
	scale  int
	screen *ebiten.Image
	lastMX int
	lastMY int
	lbtn   bool
	rbtn   bool
}

func (g *game) Update() error {
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if sc, ok := scanCode(k); ok {
			g.sh.mu.Lock()
			g.sh.keys = append(g.sh.keys, sc, sc|0x80)
			g.sh.mu.Unlock()
		}
	}
	// 可列印字元 → Stdin（AH=3Fh 路徑：旁白推進/防拷/名字）
	for _, r := range ebiten.AppendInputChars(nil) {
		if r >= 32 && r < 127 {
			g.sh.mu.Lock()
			g.sh.stdin = append(g.sh.stdin, byte(r))
			g.sh.mu.Unlock()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.sh.mu.Lock()
		g.sh.stdin = append(g.sh.stdin, '\r')
		g.sh.mu.Unlock()
	}
	mx, my := ebiten.CursorPosition()
	sx, sy := mx/g.scale, my/g.scale
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.lbtn = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.lbtn = false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.rbtn = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.rbtn = false
	}
	g.sh.mu.Lock()
	if sx != g.sh.mouseX || sy != g.sh.mouseY {
		g.sh.mouseX, g.sh.mouseY = sx, sy
	}
	g.sh.buttons[0], g.sh.buttons[1] = g.lbtn, g.rbtn
	g.sh.mu.Unlock()
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	g.sh.mu.Lock()
	rgba := append([]byte(nil), g.sh.rgba...)
	idx := append([]uint8(nil), g.sh.idx...)
	evs := evAlive(g.sh.evs, idx)
	g.sh.mu.Unlock()
	if len(rgba) == 0 {
		return
	}
	g.screen.WritePixels(rgba)
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(float64(g.scale), float64(g.scale))
	op.Filter = ebiten.FilterNearest
	dst.DrawImage(g.screen, &op)

	// ── 高解析中文合成 ──
	// 語意（AGENTS.md 高解析中文合成）：原版 320×200 終幀放大後，
	// 中文字以原字形尺寸畫在對應位置——字相對畫面變小、剛好放進按鈕。
	// hd＝放大比例的一半：scale 2 時字 15px、scale 3 時 22px。
	hd := g.scale / 2
	if hd < 1 {
		hd = 1
	}
	for i := range evs {
		e := &evs[i]
		if len(e.g) == 0 || len(e.g[0]) == 0 {
			continue // 空字形條（純控制碼/純 ASCII 譯文）跳過
		}
		if e.img == nil {
			// 字形條：整串排成一張（每字 15px＋2px 間距）
			n := 0
			for _, gl := range e.g {
				n += len(gl) + 2
			}
			rect := image.Rect(0, 0, n, len(e.g[0]))
			img := image.NewRGBA(rect)
			cx := 0
			for _, glyph := range e.g {
				for row, w := range glyph {
					for bit := 0; bit < len(glyph); bit++ {
						if (w>>(len(glyph)-1-bit))&1 != 0 {
							img.Set(cx+bit, row, color.RGBA{e.fgR, e.fgG, e.fgB, 255})
						}
					}
				}
				cx += len(glyph) + 2
			}
			e.img = ebiten.NewImageFromImage(img)
		}
		op2 := ebiten.DrawImageOptions{}
		op2.GeoM.Scale(float64(hd), float64(hd))
		op2.GeoM.Translate(float64(e.x*hd), float64(float64(e.y)*float64(hd)))
		op2.Filter = ebiten.FilterNearest
		dst.DrawImage(e.img, &op2)
	}
}

func (g *game) Layout(int, int) (int, int) { return 320 * g.scale, 200 * g.scale }

func main() {
	root := flag.String("root", "workplace/run", "原版素材目錄（可寫複本）")
	exe := flag.String("exe", "", "執行檔（預設 <root>/DS.EXE）")
	scale := flag.Int("scale", 3, "畫面放大倍率")
	ips := flag.Float64("ips", 24_000_000, "每秒指令數（節流基準）")
	xlatePath := flag.String("xlate", "workplace/bake/xlate.json", "譯文表（空字串＝純英文模式）")
	hdShotF := flag.Bool("hdshot", false, "每 30 秒輸出 HD 合成幀（PNG）")
	autoplay := flag.Bool("autoplay", false, "內建自動遊玩時序（建隊→防拷→旁白→迷宮移動）")
	hdOnly := flag.Bool("hdonly", false, "頭less：不開視窗，只輸出 HD 幀（無 WM 環境驗證用）")
	flag.Parse()

	program := *exe
	if program == "" {
		program = *root + "/DS.EXE"
	}
	o, err := oracle.LoadWith(program, *root, oracle.Options{Adlib: true})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(320**scale, 200**scale)
	ebiten.SetWindowTitle("巫術7 繁體中文版 — dosgolem")

	sh := &shared{rgba: make([]byte, 320*200*4), idx: make([]uint8, 320*200)}
	// ── xlate 翻譯層 ──────────────────────────────────────────────
	// 譯文表（原文 → 每字 15×15 位圖，tools/re/bake_font.py 烘製）。
	if *xlatePath == "" {
		*xlatePath = "/nonexistent.json"
	}
	trans := map[string]entry{}
	digits := map[string][]uint16{}
	ordinals := map[string][]uint16{}
	if b, err := os.ReadFile(*xlatePath); err == nil {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(b, &raw); err != nil {
			log.Fatal("xlate.json：", err)
		}
		for k, v := range raw {
			if k == "__digits__" {
				if err := json.Unmarshal(v, &digits); err != nil {
					log.Fatal("xlate.json digits：", err)
				}
				continue
			}
			if k == "__ordinals__" {
				if err := json.Unmarshal(v, &ordinals); err != nil {
					log.Fatal("xlate.json ordinals：", err)
				}
				continue
			}
			var e entry
			if err := json.Unmarshal(v, &e); err != nil {
				log.Fatal("xlate.json entry：", err)
			}
			trans[k] = e
		}
		log.Printf("xlate：載入 %d 串、數字 %d、序數 %d", len(trans), len(digits), len(ordinals))
	} else {
		log.Printf("xlate：無譯文表（%v）——純英文模式", xlatePath)
	}
	var evMu sync.Mutex
	treeSeen := map[uint16]bool{}
	lastHD := time.Now()
	hdShot := *hdShotF
	_ = hdShot
	os.MkdirAll("workplace/bake/trees", 0o755)
	palNow := func() [256][3]uint8 { return o.Palette() }
	_ = palNow
	// 印字串入口（runtime 0110:3895）的 hook：unpoke → 查譯文 →
	// 有譯文就記事件＋poke ret 讓原印失效。沿 cmd/xlate1 的實證邏輯。
	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	baseline := o.Byte(entry)
	poked := false
	o.OnCall(entry, func(o *oracle.Oracle) {
		if poked {
			o.WriteU8(entry, baseline)
			poked = false
		}
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		for i := 0; i < len(s); i++ {
			if s[i] == 0 {
				s = s[:i]
				break
			}
		}
		// 查表 key 正規化（與 bake_font.py 的 normalize 同一套）：
		// 剝 <0x20 控制碼、連續數字/序數詞 FIRST..TENTH → '$'（token 另存）。
		key, nums := normalizeKey(s)
		ent, ok := trans[key]
		if !ok {
			// ^＝隊員名 pattern（2040/2061/5010/1541 等）：runtime 開頭是
			// 隊員名（全大寫）——首詞換成 '$' 再查一次（表側 zh 不含名字）。
			if v := leadingWordDollar(s); v != "" {
				k2, n2 := normalizeKey(v)
				if e2, ok2 := trans[k2]; ok2 {
					ent, ok, nums = e2, true, n2
				}
			}
		}
		if !ok {
			log.Printf("無譯文 %d %q", o.Steps(), s)
			return // 無譯文：放行原文
		}
		// '$' 槽展開：數字→數字位圖、序數詞→中文字形（__ordinals__）
		if len(ent.D) > 0 {
			slotSet := map[int]bool{}
			for _, s2 := range ent.D {
				slotSet[s2] = true
			}
			newG := make([][]uint16, 0, len(ent.G)+4)
			ni := 0
			for i, gl := range ent.G {
				if slotSet[i] && ni < len(nums) {
					tok := nums[ni]
					ni++
					if dg, ok := ordinals[tok]; ok {
						newG = append(newG, dg)
						continue
					}
					for _, dc := range tok {
						if dg, ok := digits[string(dc)]; ok {
							newG = append(newG, dg)
						}
					}
					continue
				}
				newG = append(newG, gl)
			}
			ent.G = newG
		}
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})) + ent.Dy
		dl := uint8(o.Regs().DX & 0xFF)
		fg := palNow()[dl]
		// 事件區域背景雜湊（Indexed 色號；原印被 poke 擋掉，此時是純背景）
		evW := 0
		for _, g := range ent.G {
			evW += len(g) + 2
		}
		hash := vramHash(o.Indexed(), x, y, evW)
		evMu.Lock()
		// 區域重疊去重：同位置重印（jitter）只留最新
		kept := sh.evs[:0]
		for _, e := range sh.evs {
			ow := 0
			for _, g := range e.g {
				ow += len(g) + 2
			}
			if e.x < x+evW && e.x+ow > x && e.y+15 > y && e.y < y+15 {
				continue
			}
			kept = append(kept, e)
		}
		sh.evs = append(kept, glyphEv{
			x: x, y: y,
			fgR: fg[0], fgG: fg[1], fgB: fg[2],
			g:    ent.G,
			hash: hash,
		})
		evMu.Unlock()
		// 大 blit（全屏重繪）＝清文字事件（殘影清除）
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x4F5D}, func(o *oracle.Oracle) {
		// rep movsb 的 cx＝行長度；全屏寬（320）＝大重繪 → 清文字事件
		// ⚠ 實測（9/26）：此位址 OnCall 從未觸發（ovr740c.bin 檔映像該處為零
		// ——runtime 映像與檔映像不同段）——殘影清除實際靠事件區域取代，
		// 序數詞命中修正後再驗（013 追記三十二）。
		if int(o.Regs().CX) >= 320 {
			evMu.Lock()
			if len(sh.evs) > 0 {
				sh.evs = nil
			}
			evMu.Unlock()
		}
	})
	// 動態樹捕捉：樹段變化時 dump（msgtree 模式）
		ts := o.Word(oracle.Addr{Seg: 0x0110, Off: 0x2D4E})
		if ts != 0 {
			lastTree, have := treeSeen[ts]
			_ = lastTree
			if !have {
				treeSeen[ts] = true
				data := o.Bytes(oracle.Addr{Seg: ts, Off: 0}, 8192)
				fn := fmt.Sprintf("workplace/bake/trees/tree-%04x.bin", ts)
				if err := os.WriteFile(fn, data, 0o644); err == nil {
					log.Printf("樹段 %04x → %s", ts, fn)
				}
			}
		}
		log.Printf("翻譯 %d %q @(%d,%d)", o.Steps(), s, x, y)
		o.WriteU8(entry, 0xC3)
		poked = true
	})
	g := &game{sh: sh, scale: *scale}
	if !*hdOnly {
		g.screen = ebiten.NewImage(320, 200)
	}

	// 自動遊玩步數門控表（oracle 層直接注入，繞過 Ebiten 輸入在 Xvfb 的失效）。
	// R10 改用 steps 觸發：wall-clock 在 R8（主選單 +45s）vs R9（約 +60s）間漂移，
	// steps 有 M1 決定性保證。錨點來自 R9 實測（docs/re/010 R9 追記）。
	var autos []struct {
		at uint64 // 機器步數門檻
		fn func(*oracle.Oracle)
	}
	if *autoplay {
		add := func(s uint64, fn func(*oracle.Oracle)) { autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{s, fn}) }
		cl := func(x, y int) func(*oracle.Oracle) {
			return func(o *oracle.Oracle) { o.Click(x, y) }
		}
		// START NEW GAME 路線（R9 拓樸修正：建隊對話框後直達七項選單，不繞 PCVW）：
		// 主選單 246M→建隊→THESUS 254M→七項選單 260M→START→巫6匯入詢問 CANCEL→
		// 防拷→旁白 40 次→迷宮八點。門檻＝錨點＋餘裕（約 4–5M steps）。
		add(250_000_000, cl(160, 62)) // 主選單「新增角色」直達建隊
		add(253_000_000, cl(170, 153))
		add(255_000_000, cl(90, 168))
		add(257_000_000, cl(90, 168))
		add(265_000_000, cl(160, 101)) // 七項選單: START NEW GAME（y-13 熱點系統）
		add(276_000_000, cl(82, 169))  // 巫6 匯入詢問: CANCEL（語意待驗）
		add(288_000_000, func(o *oracle.Oracle) { o.Type("AAAA\r") }) // 防拷
		for i := 0; i < 133; i++ { // 旁白推進（R11：引言延續到 1225M＋，每 8M 一次）
			add(295_000_000+uint64(i)*8_000_000, cl(25, 15))
		}
		// 迷宮方向盤十字 sweep（移動實證）。旁白 133 次推進到 ~1351M，
		// 迷宮點擊必須在其後。
		dirS := []struct {
			s uint64
			x, y int
		}{{1_380_000_000, 142, 165}, {1_384_000_000, 160, 155}, {1_388_000_000, 160, 178},
			{1_392_000_000, 160, 196}, {1_396_000_000, 113, 178}, {1_400_000_000, 188, 178},
			{1_404_000_000, 142, 165}, {1_408_000_000, 142, 165}}
		for _, d := range dirS {
			add(d.s, cl(d.x, d.y))
		}
		// 每點後截圖（移動實證）
		for _, d := range dirS {
			shotS := d.s + 9_000_000
			add(shotS, func(o *oracle.Oracle) {
				evMu.Lock()
				hdFrame(o, sh.evs, fmt.Sprintf("workplace/bake/mv-%d-%d-%d.png", d.s, d.x, d.y))
				evMu.Unlock()
			})
		}
	}
	autoIdx := 0

	// 機器 goroutine：大塊 Run（probe 同款快路徑），節流＋吃前端輸入。
	// ⚠ MoveMouse 只在座標真的變了才送——遊戲把滑鼠活動當「有人在看」，
	// 每塊都送會重置故事頁的推進計時（docs/re/009 的教訓）。
	go func() {
		chunk := uint64(200_000)
		minWall := time.Duration(float64(chunk) / *ips * float64(time.Second))
		last := time.Now()
		lastR := false
		lastX, lastY := -1, -1
		for i := 0; ; i++ {
			sh.mu.Lock()
			for _, sc := range sh.keys {
				o.PressKey(oracle.Key(sc))
			}
			sh.keys = nil
			if len(sh.stdin) > 0 {
				o.Type(string(sh.stdin))
				sh.stdin = nil
			}
			bx, by := sh.mouseX, sh.mouseY
			r := sh.buttons[1]
			if bx != lastX || by != lastY {
				o.MoveMouse(bx, by)
				lastX, lastY = bx, by
			}
			// 整包點擊（Click 路徑——分離的 Press/Release 在部分遊戲失效）
			var clicks []struct{ x, y int }
			for len(sh.clicks) > 0 {
				c := sh.clicks[0]
				sh.clicks = sh.clicks[1:]
				clicks = append(clicks, c)
			}
			// MoveMouse（游標位置）
			if bx != lastX || by != lastY {
				o.MoveMouse(bx, by)
				lastX, lastY = bx, by
			}
			sh.mu.Unlock()
			for _, c := range clicks {
				if err := o.Click(c.x, c.y); err != nil {
					log.Printf("點擊 (%d,%d)：%v", c.x, c.y, err)
				}
			}
			sh.mu.Lock()
			// 右鍵維持 Press/Release 路徑
			if r != lastR {
				if r {
					o.PressButton(1)
				} else {
					o.ReleaseButton(1)
				}
				lastR = r
			}
			if err := o.Run(chunk); err != nil {
				sh.mu.Unlock()
				log.Printf("機器停止於 steps=%d：%v", o.Steps(), err)
				return
			}
			// 自動遊玩：步數過門檻執行（R10 步數門控）
			for autoIdx < len(autos) {
				a := autos[autoIdx]
				if o.Steps() < a.at {
					break
				}
				sh.mu.Unlock()
				fmt.Printf("AUTO %d 執行 steps=%d\n", a.at, o.Steps())
				a.fn(o)
				autoIdx++
				sh.mu.Lock()
			}
			// HD 幀輸出（每 30 秒）：A0000 2x 放大＋文字事件疊畫
			if hdShot && time.Since(lastHD) > 3*time.Second {
				lastHD = time.Now()
				evMu.Lock()
				hdFrame(o, sh.evs, fmt.Sprintf("workplace/bake/hd-%d.png", time.Now().Unix()))
				evMu.Unlock()
			}
			if i%60 == 0 {
				nz := 0
				for _, v := range o.Indexed() {
					if v != 0 {
						nz++
					}
				}
				_, _, rgb := o.ScreenRGB()
				nz2 := 0
				for j := 0; j+2 < len(rgb); j += 3 {
					if rgb[j] != 0 || rgb[j+1] != 0 || rgb[j+2] != 0 {
						nz2++
					}
				}
				pal := o.Palette()
				log.Printf("chunk %d steps=%d mode=%d idxNZ=%d rgbNZ=%d pal[15]=%v", i, o.Steps(), o.VideoMode(), nz, nz2, pal[15])
			}
			_, _, rgb := o.ScreenRGB()
			copy(sh.rgba, rgb)
			copy(sh.idx, o.Indexed())
			sh.mu.Unlock()

			elapsed := time.Since(last)
			if elapsed < minWall {
				time.Sleep(minWall - elapsed)
			}
			last = time.Now()
		}
	}()

	if *hdOnly {
		// 頭less 驗證：不開視窗（無 WM 時 GLFW init 會 panic），
		// 只跑機器 goroutine＋hdFrame 輸出。
		select {}
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
