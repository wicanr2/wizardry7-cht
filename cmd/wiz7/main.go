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

// 前端 ↔ 機器 goroutine 的共享狀態（mu 保護）。
// entry 是譯文表的一筆（xlate.json）：dy＝行距微調、g＝每字 15×15 位圖。
// hdFrame 合成高解析幀：A0000 2x 放大＋文字事件（15×15 原字形）。
func hdFrame(o *oracle.Oracle, evs []glyphEv, path string) {
	pal := o.Palette()
	idx := o.Indexed()
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
}

// glyphEv 是一則待合成的中文字串（15×15 原字形、畫在放大後畫布）。
type glyphEv struct {
	x, y          int
	fgR, fgG, fgB uint8
	g             [][]uint16
	img           *ebiten.Image // Draw 時 lazily 預渲染的字形條
}

type shared struct {
	mu      sync.Mutex
	rgba    []byte    // 320×200×4 的畫面快照（Draw 用）
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
	evs := g.sh.evs
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

	sh := &shared{rgba: make([]byte, 320*200*4)}
	// ── xlate 翻譯層 ──────────────────────────────────────────────
	// 譯文表（原文 → 每字 15×15 位圖，tools/re/bake_font.py 烘製）。
	if *xlatePath == "" {
		*xlatePath = "/nonexistent.json"
	}
	trans := map[string]entry{}
	if b, err := os.ReadFile(*xlatePath); err == nil {
		if err := json.Unmarshal(b, &trans); err != nil {
			log.Fatal("xlate.json：", err)
		}
		log.Printf("xlate：載入 %d 串", len(trans))
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
		ent, ok := trans[s]
		if !ok {
			log.Printf("無譯文 %d %q", o.Steps(), s)
			return // 無譯文：放行原文
		}
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})) + ent.Dy
		dl := uint8(o.Regs().DX & 0xFF)
		fg := palNow()[dl]
		evMu.Lock()
		sh.evs = append(sh.evs, glyphEv{
			x: x, y: y,
			fgR: fg[0], fgG: fg[1], fgB: fg[2],
			g: ent.G,
		})
		evMu.Unlock()
		// 大 blit（全屏重繪）＝清文字事件（殘影清除）
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x4F5D}, func(o *oracle.Oracle) {
		if int(o.Regs().CX) > 0x4000 {
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
	g := &game{sh: sh, scale: *scale, screen: ebiten.NewImage(320, 200)}

	// 自動遊玩時序表（oracle 層直接注入，繞過 Ebiten 輸入在 Xvfb 的失效）
	var autos []struct {
		after float64 // 啟動後秒數
		fn    func(*oracle.Oracle)
	}
	if *autoplay {
		add := func(t float64, fn func(*oracle.Oracle)) { autos = append(autos, struct {
			after float64
			fn    func(*oracle.Oracle)
		}{t, fn}) }
		cl := func(x, y int) func(*oracle.Oracle) {
			return func(o *oracle.Oracle) { o.Click(x, y) }
		}
		add(16, cl(160, 62)) // ADD CHARACTER
		add(18, cl(170, 153))
		add(19.5, cl(90, 168))
		add(21, cl(90, 168))
		add(24, cl(160, 101)) // START NEW GAME
		add(26, func(o *oracle.Oracle) { o.Type("AAAA\r") }) // 防拷
		for i := 0; i < 24; i++ {                            // 旁白推進
			t := 27.5 + float64(i)*1.2
			add(t, cl(25, 15))
		}
		// 迷宮移動
		for i := 0; i < 8; i++ {
			add(60+float64(i)*2, cl(142, 165)) // 上箭頭（前進）
		}
	}

	startWall := time.Now()
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
			// 自動遊玩：啟動秒數過點執行
			for autoIdx < len(autos) {
				a := autos[autoIdx]
				if time.Since(startWall).Seconds() < a.after {
					break
				}
				sh.mu.Unlock()
				fmt.Printf("AUTO %.1fs 執行\n", a.after)
				a.fn(o)
				autoIdx++
				sh.mu.Lock()
			}
			// HD 幀輸出（每 30 秒）：A0000 2x 放大＋文字事件疊畫
			if hdShot && time.Since(lastHD) > 30*time.Second {
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
			sh.mu.Unlock()

			elapsed := time.Since(last)
			if elapsed < minWall {
				time.Sleep(minWall - elapsed)
			}
			last = time.Now()
		}
	}()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
