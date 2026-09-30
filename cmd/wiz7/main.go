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
	"sort"
	"strconv"
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

// ── 不寫原版的印字覆繪（docs/spec/001，證據 docs/re/014 E1–E9）──────────
//
// 原版每一則字串照常畫出英文（R1）。轉譯層在印字入口與返回各取一張離屏
// 快照，差異矩形就是英文筆畫所在（E4）；合成時以 A000 判定該處是否仍是
// 這串英文（E7），是的話把筆畫還原成印字前的背景，再在放大畫布上畫中文。

// evStrings 取 shown 事件的英文原文（畫面辨識用，spec 003 R4）。
func evStrings(evs []*textEv) []string {
	out := make([]string, 0, len(evs))
	for _, e := range evs {
		out = append(out, e.s)
	}
	return out
}

// offscreen 回傳 VGA.DRV 離屏緩衝（blit 來源段＝740C:01E5 的值，E5）。
// 段值為 0 時回 nil（尚未初始化）。
func offscreen(o *oracle.Oracle) []uint8 {
	seg := o.Word(oracle.Addr{Seg: 0x740C, Off: 0x01E5})
	if seg == 0 {
		return nil
	}
	return o.Bytes(oracle.Addr{Seg: seg, Off: 0}, 320*200)
}

// textEv 是一則有譯文的印字事件（R2.4）。
type textEv struct {
	s              string
	x0, y0, dy     int
	bx, by, bw, bh int     // 英文筆畫框（原版座標）
	before, after  []uint8 // 框內離屏色號：印字前／後
	fg             uint8   // 筆畫色號（差異像素眾數）
	fgRGB          [3]uint8
	g              [][]uint16
	shown          bool
	img            *ebiten.Image // Draw 時 lazily 預渲染的字形條
}

// pendingPrint 是已進入、尚未返回的印字（R2.1）。
type pendingPrint struct {
	s      string
	x0, y0 int
	dy     int
	sp     uint16
	ret    oracle.Addr
	before []uint8
	g      [][]uint16
}

// overlay 保存事件與統計；呼叫端持有 evMu。
type overlay struct {
	pend    *pendingPrint
	evs     []*textEv
	created int
	shown   int
	dropped int
	empty   int
	nested  int
	noseg   int
}

// cursorBox 是 A000 上游標可能佔的框（E8 量到 10×10，未量形狀先放寬到 16×16）。
const cursorBox = 16

// finish 在返回 hook 取 after 快照、算差異框並建立事件（R2.3–R2.4）。
func (ov *overlay) finish(o *oracle.Oracle) {
	p := ov.pend
	ov.pend = nil
	after := offscreen(o)
	if after == nil || p.before == nil {
		ov.noseg++
		return
	}
	x0, y0, x1, y1 := 320, 200, -1, -1
	for i := range after {
		if after[i] != p.before[i] {
			x, y := i%320, i/320
			x0, y0 = min(x0, x), min(y0, y)
			x1, y1 = max(x1, x), max(y1, y)
		}
	}
	if x1 < 0 {
		ov.empty++
		return
	}
	w, h := x1-x0+1, y1-y0+1
	e := &textEv{s: p.s, x0: p.x0, y0: p.y0, dy: p.dy, bx: x0, by: y0, bw: w, bh: h, g: p.g}
	e.before = make([]uint8, w*h)
	e.after = make([]uint8, w*h)
	var cnt [256]int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y0+y)*320 + x0 + x
			k := y*w + x
			e.before[k], e.after[k] = p.before[i], after[i]
			if after[i] != p.before[i] {
				cnt[after[i]]++
			}
		}
	}
	for c := 1; c < 256; c++ {
		if cnt[c] > cnt[e.fg] {
			e.fg = uint8(c)
		}
	}
	e.fgRGB = o.Palette()[e.fg] // 返回 hook 當下的調色盤（R2.4）
	ov.evs = append(ov.evs, e)
	ov.created++
}

// compose 依 R3 判定事件狀態，回傳合成後的低解析色號與 shown 事件（R4.1）。
// 只在機器 goroutine 的快照節點呼叫（R4.3）。
func (ov *overlay) compose(o *oracle.Oracle) ([]uint8, []*textEv) {
	a := o.Indexed()
	comp := append([]uint8(nil), a...)
	off := offscreen(o)
	mx, my := o.Mouse()
	inCursor := func(x, y int) bool {
		return x >= mx && x < mx+cursorBox && y >= my && y < my+cursorBox
	}
	kept := ov.evs[:0]
	var shown []*textEv
	for _, e := range ov.evs {
		match := true
		for y := 0; y < e.bh && match; y++ {
			for x := 0; x < e.bw; x++ {
				px, py := e.bx+x, e.by+y
				if inCursor(px, py) {
					continue
				}
				if a[py*320+px] != e.after[y*e.bw+x] {
					match = false
					break
				}
			}
		}
		switch {
		case match:
			if !e.shown {
				e.shown = true
				ov.shown++
			}
			for y := 0; y < e.bh; y++ {
				for x := 0; x < e.bw; x++ {
					k := y*e.bw + x
					i := (e.by+y)*320 + e.bx + x
					if a[i] == e.after[k] && e.before[k] != e.after[k] {
						comp[i] = e.before[k]
					}
				}
			}
			kept = append(kept, e)
			shown = append(shown, e)
		case e.shown:
			ov.dropped++
		default:
			// 尚未上屏：離屏框仍是 after 才保留（R3 表第 3、4 列）
			still := off != nil
			for y := 0; y < e.bh && still; y++ {
				for x := 0; x < e.bw; x++ {
					if off[(e.by+y)*320+e.bx+x] != e.after[y*e.bw+x] {
						still = false
						break
					}
				}
			}
			if still {
				kept = append(kept, e)
			} else {
				ov.dropped++
			}
		}
	}
	for i := len(kept); i < len(ov.evs); i++ {
		ov.evs[i] = nil
	}
	ov.evs = kept
	return comp, shown
}

func (ov *overlay) stats() string {
	return fmt.Sprintf("事件 建立=%d 上屏=%d 刪除=%d 無差異=%d 巢狀捨棄=%d 無離屏=%d 存活=%d",
		ov.created, ov.shown, ov.dropped, ov.empty, ov.nested, ov.noseg, len(ov.evs))
}

func fnv(b []byte) uint32 {
	var h uint32 = 2166136261
	for _, v := range b {
		h = (h ^ uint32(v)) * 16777619
	}
	return h
}

// hdScale 是 hdFrame 的畫布倍率（與互動視窗 -scale 同義）。行距 6 的畫面
// （防拷三行）在 2× 只有 12 px，放不下 15 px 的中文字；3× 為 18 px。
var hdScale = 3

// hdFrame 在快照節點合成高解析幀：低解析合成（已還原英文筆畫）放大 hdScale 倍，
// 中文 15×15 原字形畫在 (x0×k, (y0+dy)×k)。呼叫端持有 evMu。
func hdFrame(o *oracle.Oracle, ov *overlay, path string) (comp []uint8, shown []*textEv) {
	comp, shown = ov.compose(o)
	pal := o.Palette()
	k := hdScale
	img := image.NewRGBA(image.Rect(0, 0, 320*k, 200*k))
	for py := 0; py < 200*k; py++ {
		for px := 0; px < 320*k; px++ {
			c := pal[comp[(py/k)*320+(px/k)]]
			img.Set(px, py, color.RGBA{c[0], c[1], c[2], 255})
		}
	}
	for _, e := range shown {
		fr := color.RGBA{e.fgRGB[0], e.fgRGB[1], e.fgRGB[2], 255}
		cx := e.x0 * k
		for _, glyph := range e.g {
			for row, w := range glyph {
				for bit := 0; bit < len(glyph); bit++ {
					if (w>>(len(glyph)-1-bit))&1 != 0 {
						img.Set(cx+bit, (e.y0+e.dy)*k+row, fr)
					}
				}
			}
			cx += len(glyph) + 2
		}
	}
	if f, err := os.Create(path); err == nil {
		png.Encode(f, img)
		f.Close()
	}
	return comp, shown
}

// hdAtDump 輸出 -hd-at 驗收資料（spec 001 R5 與驗收）：HD 合成 PNG、
// A000 原始色號、低解析合成色號、shown 事件與雜湊。呼叫端持有 evMu。
func hdAtDump(o *oracle.Oracle, ov *overlay, am *automap, dir string, st uint64) {
	base := fmt.Sprintf("%s/hdat-%d", dir, st)
	a := append([]uint8(nil), o.Indexed()...)
	comp, shown := hdFrame(o, ov, base+".png")
	pal := o.Palette()
	palB := make([]byte, 0, 768)
	for _, c := range pal {
		palB = append(palB, c[0], c[1], c[2])
	}
	mem := fnv(o.Bytes(oracle.Addr{Seg: 0, Off: 0}, 0x100000))
	type evOut struct {
		S              string
		X0, Y0, Dy     int
		BX, BY, BW, BH int
		Before, After  []uint8
		Fg             uint8
		G              [][]uint16 // 字形（spec 002 R5 的序數檢查用）
	}
	out := struct {
		Step           uint64
		Mem, A000, Pal uint32
		Stats          string
		Mouse          [2]int
		Shown          []evOut
		Map            mapJSON
	}{Step: o.Steps(), Mem: mem, A000: fnv(a), Pal: fnv(palB), Stats: ov.stats()}
	out.Mouse[0], out.Mouse[1] = o.Mouse()
	out.Map = am.json()
	for _, e := range shown {
		out.Shown = append(out.Shown, evOut{e.s, e.x0, e.y0, e.dy, e.bx, e.by, e.bw, e.bh, e.before, e.after, e.fg, e.g})
	}
	_ = os.WriteFile(base+".a000", a, 0o644)
	_ = os.WriteFile(base+".comp", comp, 0o644)
	if b, err := json.Marshal(out); err == nil {
		_ = os.WriteFile(base+".json", b, 0o644)
	}
	log.Printf("hdat %d steps=%d mem=%08x a000=%08x pal=%08x shown=%d %s",
		st, out.Step, out.Mem, out.A000, out.Pal, len(shown), out.Stats)
}

type entry struct {
	Dy int        `json:"dy"`
	G  [][]uint16 `json:"g"`
	D  []int      `json:"d,omitempty"` // '$' 數字槽索引（g 內位置，位圖全 0）
}

type shared struct {
	mu         sync.Mutex
	rgba       []byte  // 320×200×4 的畫面快照（Draw 用）
	keys       []uint8 // 待注入的掃描碼（make＋break 成對排程）
	mouseX     int
	mouseY     int
	rbtn       bool      // 右鍵狀態（左鍵走 clicks，spec 004）
	evs        []*textEv // shown 事件（快照節點判定，Draw 只消費，R4.3）
	drawReq    string    // 非空＝下一次 Draw 完成後把視窗畫面存到這個路徑（-draw-at）
	helpOpen   bool      // F1 說明開啟中＝機器暫停、輸入不送原版（spec 003 R2）
	helpScreen string    // 開啟當下辨識的畫面（R4）
	helpPage   int
	mapOpen    bool                 // F2 自動地圖面板（spec 005 R3）
	mapView    mapView              // 快照節點產生的地圖狀態
	stdin      []byte               // 待注入的 Stdin 字元（AH=3Fh 路徑）
	clicks     []struct{ x, y int } // 待執行的整包點擊（Click 路徑）
}

type game struct {
	help    *helpData     // F1 說明資料（help.json；nil＝缺檔）
	helpImg *ebiten.Image // 目前說明頁的預渲染
	helpKey string
	sh      *shared
	scale   int
	screen  *ebiten.Image
	lastMX  int
	lastMY  int
	lbtn    bool
	rbtn    bool
}

func (g *game) Update() error {
	// ── F1 說明（spec 003 R1）：F1 永遠由前端處理；說明開啟時其他輸入一律不送原版 ──
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.sh.mu.Lock()
		if g.sh.helpOpen {
			g.sh.helpOpen = false
		} else {
			g.sh.helpOpen, g.sh.helpPage = true, 0
			g.sh.helpScreen = detectScreen(evStrings(g.sh.evs))
		}
		log.Printf("說明 %v 畫面=%s", g.sh.helpOpen, g.sh.helpScreen)
		g.sh.mu.Unlock()
	}
	g.sh.mu.Lock()
	open := g.sh.helpOpen
	if open {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
			g.sh.helpOpen = false
			log.Printf("說明 false（Esc）")
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), inpututil.IsKeyJustPressed(ebiten.KeyPageUp):
			g.sh.helpPage--
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight), inpututil.IsKeyJustPressed(ebiten.KeyPageDown):
			g.sh.helpPage++
		}
	}
	g.sh.mu.Unlock()
	if open {
		ebiten.AppendInputChars(nil) // 丟棄本幀字元，不排入 stdin
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) { // spec 005 R3：前端處理，不送原版
		g.sh.mu.Lock()
		g.sh.mapOpen = !g.sh.mapOpen
		log.Printf("地圖 面板=%v", g.sh.mapOpen)
		g.sh.mu.Unlock()
	}
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
		// spec 004 R1：左鍵剛按下＝排一次整包點擊（遊戲座標，裁切到 320×200）
		cx, cy := min(max(sx, 0), 319), min(max(sy, 0), 199)
		g.sh.mu.Lock()
		g.sh.clicks = append(g.sh.clicks, struct{ x, y int }{cx, cy})
		g.sh.mu.Unlock()
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
	g.sh.rbtn = g.rbtn
	g.sh.mu.Unlock()
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	g.sh.mu.Lock()
	rgba := append([]byte(nil), g.sh.rgba...)
	evs := append([]*textEv(nil), g.sh.evs...)
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
	for _, e := range evs {
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
							img.Set(cx+bit, row, color.RGBA{e.fgRGB[0], e.fgRGB[1], e.fgRGB[2], 255})
						}
					}
				}
				cx += len(glyph) + 2
			}
			e.img = ebiten.NewImageFromImage(img)
		}
		op2 := ebiten.DrawImageOptions{}
		op2.GeoM.Scale(float64(hd), float64(hd))
		// 平移在縮放之後，單位是視窗像素：原版座標 × 倍率
		op2.GeoM.Translate(float64(e.x0*g.scale), float64((e.y0+e.dy)*g.scale))
		op2.Filter = ebiten.FilterNearest
		dst.DrawImage(e.img, &op2)
	}

	g.sh.mu.Lock()
	hOpen, hScreen, hPage := g.sh.helpOpen, g.sh.helpScreen, g.sh.helpPage
	mOpen, mView := g.sh.mapOpen, g.sh.mapView
	g.sh.mu.Unlock()
	if mOpen {
		g.drawMap(dst, mView)
	}
	if hOpen {
		g.drawHelp(dst, hScreen, hPage)
	}

	// -draw-at：把互動視窗實際畫面存檔（驗證 Draw 路徑，spec 001 R4.3）
	g.sh.mu.Lock()
	req := g.sh.drawReq
	g.sh.drawReq = ""
	g.sh.mu.Unlock()
	if req != "" {
		w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		dst.ReadPixels(img.Pix)
		if f, err := os.Create(req); err == nil {
			png.Encode(f, img)
			f.Close()
			log.Printf("drawat 輸出 %s（%dx%d）", req, w, h)
		}
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
	hdAt := flag.String("hd-at", "", "步數逗號分隔：到該步數輸出 HD 合成幀與驗收資料 <hd-dir>/hdat-<步數>.*（spec 001 R5）")
	hdDir := flag.String("hd-dir", "workplace/bake", "-hd-at 的輸出目錄")
	flag.IntVar(&hdScale, "hd-scale", 3, "HD 幀（-hd-at、-hdshot、自動遊玩截圖）的畫布倍率（使用者定案 2026-09-30：預設 3）")
	helpPath := flag.String("help", "workplace/bake/help.json", "F1 說明資料（tools/re/bake_help.py 產生）")
	helpAt := flag.Uint64("help-at", 0, "測試用：到該步數開啟 F1 說明（spec 003 R6）")
	helpFor := flag.Int("help-for", 3000, "測試用：-help-at 開啟後幾毫秒（牆鐘）關閉")
	helpPageF := flag.Int("help-page", 0, "測試用：-help-at 開啟時停在第幾頁（0 起算）")
	clicksAt := flag.String("clicks-at", "", "測試用：step:x:y 逗號分隔，到該步數以 o.Click 點擊（遊戲座標；spec 005 R4）")
	mapAt := flag.Uint64("map-at", 0, "測試用：到該步數開啟自動地圖面板")
	autoplayStop := flag.Uint64("autoplay-stop", 0, "測試用：autoplay 只執行門檻 ≤ 此步數的動作（0＝全部），之後交給真實輸入")
	drawAt := flag.String("draw-at", "", "步數逗號分隔：到該步數後把互動視窗畫面存到 <hd-dir>/drawat-<步數>.png")
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
	// 印字串入口（0110:3895）：只讀，不寫原版（spec 001 R1、R2）。
	ov := &overlay{}
	am := newAutomap() // spec 005：自動地圖（只在快照節點更新）
	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	retHooked := map[uint32]bool{}
	o.OnCall(entry, func(o *oracle.Oracle) {
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
			log.Printf("無譯文 %d %q @(%d,%d)", o.Steps(), s,
				o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}), o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
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
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
		r := o.Regs()
		p := &pendingPrint{
			s: s, x0: x, y0: y, dy: ent.Dy, g: ent.G,
			sp:  r.SP,
			ret: oracle.Addr{Seg: r.CS, Off: o.StackWord(0)}, // near call（E2）
		}
		if b := offscreen(o); b != nil {
			p.before = append([]byte(nil), b...)
		}
		evMu.Lock()
		if ov.pend != nil {
			// 防禦性：E2 實測未見巢狀（R2.5）
			ov.nested++
			log.Printf("巢狀印字：捨棄未返回的 %q", ov.pend.s)
		}
		ov.pend = p
		evMu.Unlock()
		if lin := p.ret.Linear(); !retHooked[lin] {
			retHooked[lin] = true
			ret := p.ret
			o.OnCall(ret, func(o *oracle.Oracle) {
				evMu.Lock()
				defer evMu.Unlock()
				// 只認本次返回：IP＝返回位址且 SP＝進入 SP＋2（R2.2）
				if ov.pend == nil || ov.pend.ret != ret || o.Regs().SP != ov.pend.sp+2 {
					return
				}
				ov.finish(o)
			})
		}
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
	})
	g := &game{sh: sh, scale: *scale, help: loadHelp(*helpPath)}
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
		add := func(s uint64, fn func(*oracle.Oracle)) {
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{s, fn})
		}
		// R14：游標 park 觀測。每次動作後記 pos／sets／range。
		logMouse := func(o *oracle.Oracle, op string, x, y int, err error) {
			mx, my := o.Mouse()
			sets := o.MouseSets()
			last := "none"
			if len(sets) > 0 {
				l := sets[len(sets)-1]
				last = fmt.Sprintf("(%d,%d)@%d", l.X, l.Y, l.Step)
			}
			n0, n1, n2, n3 := o.MouseRange()
			log.Printf("AUTO滑鼠 %s (%d,%d) err=%v pos=(%d,%d) sets=%d last=%s range=[%d-%d,%d-%d]",
				op, x, y, err, mx, my, len(sets), last, n0, n1, n2, n3)
		}
		cl := func(x, y int) func(*oracle.Oracle) {
			// R13：印出 Click error（含 Polls）——畫面不變時區分座標錯／沒準備好。
			// R14：附游標觀測（park 假說）。
			return func(o *oracle.Oracle) {
				err := o.Click(x, y)
				if err != nil {
					log.Printf("AUTO點擊 (%d,%d)：%v", x, y, err)
				}
				logMouse(o, "click", x, y, err)
			}
		}
		domouse := func(op string, x, y int) func(*oracle.Oracle) {
			return func(o *oracle.Oracle) {
				var err error
				switch op {
				case "press":
					o.MoveMouse(x, y)
					o.PressButton(0)
				case "release":
					o.ReleaseButton(0)
				case "right":
					err = o.Click(x, y, oracle.Right())
				case "racepress":
					// R15：Move 後立即零 hover 按住 20M，賭第一次 park 前的 poll。
					o.MoveMouse(x, y)
					err = o.Press(oracle.Hold(20_000_000))
				case "parkpress":
					// R15 對照組：不 Move，原地按（讀到的必是 park 點，預期不動）。
					err = o.Press(oracle.Hold(20_000_000))
				case "log":
					// 純觀測，不動輸入。
				case "move":
					// R18：只搬移（按住另由 press/release 管），觀測 pos 釘住否。
					o.MoveMouse(x, y)
				default:
					err = o.Click(x, y)
				}
				logMouse(o, op, x, y, err)
			}
		}
		// START NEW GAME 路線（R9 拓樸修正：建隊對話框後直達七項選單，不繞 PCVW）：
		// 主選單 246M→建隊→THESUS 254M→七項選單 260M→START→巫6匯入詢問 CANCEL→
		// 防拷→旁白 40 次→迷宮八點。門檻＝錨點＋餘裕（約 4–5M steps）。
		add(250_000_000, cl(160, 62)) // 主選單「新增角色」直達建隊
		add(253_000_000, cl(170, 153))
		add(255_000_000, cl(90, 168))
		add(257_000_000, cl(90, 168))
		add(265_000_000, cl(160, 101))                                // 七項選單: START NEW GAME（y-13 熱點系統）
		add(276_000_000, cl(82, 169))                                 // 巫6 匯入詢問: CANCEL（語意待驗）
		add(288_000_000, func(o *oracle.Oracle) { o.Type("AAAA\r") }) // 防拷
		for i := 0; i < 133; i++ {                                    // 旁白推進（R11：引言延續到 1225M＋，每 8M 一次）
			add(295_000_000+uint64(i)*8_000_000, cl(25, 15))
		}
		// 迷宮精確箭頭（R20）：先點底文清場，再打實測幾何的箭頭。
		// driver 座標（x 已加倍；羅盤 game px 見 010 R19 追記）。
		type mvAct struct {
			s    uint64
			x, y int
			op   string
		}
		dirS := []mvAct{
			{1_430_000_000, 320, 170, "click"}, // 底文（清場，R19 confirmed 有效）
			{1_434_000_000, 320, 150, "click"}, // 底文上部（保險）
			{1_438_000_000, 310, 118, "click"}, // 北（上箭頭）
			{1_442_000_000, 230, 140, "click"}, // 西（左箭頭）
			{1_446_000_000, 390, 140, "click"}, // 東（右箭頭）
			{1_450_000_000, 310, 158, "click"}, // 南（下箭頭）
			{1_454_000_000, 310, 138, "click"}, // 中央
			{1_458_000_000, 320, 100, "click"}, // 中央視野
			{1_462_000_000, 310, 118, "click"}, // 北再一次
			{1_466_000_000, 230, 140, "click"}, // 西再一次
		}
		for _, d := range dirS {
			add(d.s, domouse(d.op, d.x, d.y))
		}
		// 每點後截圖（移動實證）
		for _, d := range dirS {
			shotS := d.s + 9_000_000
			add(shotS, func(o *oracle.Oracle) {
				evMu.Lock()
				hdFrame(o, ov, fmt.Sprintf("workplace/bake/mv-%d-%d-%d.png", d.s, d.x, d.y))
				evMu.Unlock()
			})
		}
	}
	if *autoplay && *autoplayStop > 0 {
		kept := autos[:0]
		for _, a := range autos {
			if a.at <= *autoplayStop {
				kept = append(kept, a)
			}
		}
		autos = kept
		log.Printf("autoplay 停在 %d 步之後（autoplay 剩 %d 個動作）", *autoplayStop, len(autos))
	}
	// -hd-at：依步數輸出 HD 合成幀（決定性，取代 wall-clock 的 -hdshot 做驗收）。
	if *hdAt != "" {
		for _, f := range strings.Split(*hdAt, ",") {
			st, err := strconv.ParseUint(strings.TrimSpace(f), 10, 64)
			if err != nil {
				log.Fatalf("-hd-at %q：%v", f, err)
			}
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{st, func(o *oracle.Oracle) {
				evMu.Lock()
				hdAtDump(o, ov, am, *hdDir, st)
				evMu.Unlock()
			}})
		}
	}
	for _, f := range strings.Split(*clicksAt, ",") {
		var st uint64
		var x, y int
		if f == "" {
			continue
		}
		if _, err := fmt.Sscanf(f, "%d:%d:%d", &st, &x, &y); err != nil {
			log.Fatalf("-clicks-at %q：%v", f, err)
		}
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{st, func(o *oracle.Oracle) {
			err := o.Click(x, y)
			log.Printf("測試點擊 (%d,%d) steps=%d err=%v", x, y, o.Steps(), err)
		}})
	}
	if *mapAt > 0 {
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{*mapAt, func(o *oracle.Oracle) {
			sh.mu.Lock()
			sh.mapOpen = true
			sh.mu.Unlock()
			log.Printf("地圖 面板=true（-map-at）")
		}})
	}
	if *helpAt > 0 {
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{*helpAt, func(o *oracle.Oracle) {
			sh.mu.Lock()
			sh.helpOpen, sh.helpPage = true, *helpPageF
			sh.helpScreen = detectScreen(evStrings(sh.evs))
			log.Printf("說明 true（-help-at）畫面=%s shown=%q", sh.helpScreen, evStrings(sh.evs))
			sh.mu.Unlock()
			time.AfterFunc(time.Duration(*helpFor)*time.Millisecond, func() {
				sh.mu.Lock()
				sh.helpOpen = false
				sh.mu.Unlock()
				log.Printf("說明 false（-help-for）")
			})
		}})
	}
	if *drawAt != "" {
		for _, f := range strings.Split(*drawAt, ",") {
			st, err := strconv.ParseUint(strings.TrimSpace(f), 10, 64)
			if err != nil {
				log.Fatalf("-draw-at %q：%v", f, err)
			}
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{st, func(o *oracle.Oracle) {
				sh.mu.Lock()
				sh.drawReq = fmt.Sprintf("%s/drawat-%d.png", *hdDir, st)
				sh.mu.Unlock()
			}})
		}
	}
	sort.SliceStable(autos, func(i, j int) bool { return autos[i].at < autos[j].at })
	autoIdx := 0

	// snapshot 是快照節點（spec 001 R4.3）：R3 判定一次，低解析合成（已還原英文筆畫）
	// 轉 RGBA 給 Draw，並更新 shown 事件。呼叫端持有 sh.mu。
	// 每輪迴圈結尾與每個 auto 執行後都要呼叫：autoplay 的 o.Click 內部會跑數百萬步，
	// 多個 auto 可能在同一輪連鎖執行，只在輪末快照會讓 sh.evs 過期（R667 實測：
	// 240M→286M 之間零快照，-help-at 讀到 247M 的舊主選單）。
	snapshot := func() {
		am.update(o) // spec 005 R1／R2：只讀
		sh.mapView = am.view()
		evMu.Lock()
		comp, shown := ov.compose(o)
		evMu.Unlock()
		pal := o.Palette()
		for i, v := range comp {
			c := pal[v]
			sh.rgba[i*4], sh.rgba[i*4+1], sh.rgba[i*4+2], sh.rgba[i*4+3] = c[0], c[1], c[2], 255
		}
		sh.evs = shown
	}

	// 機器 goroutine：大塊 Run（probe 同款快路徑），節流＋吃前端輸入。
	// ⚠ MoveMouse 只在座標真的變了才送——遊戲把滑鼠活動當「有人在看」，
	// 每塊都送會重置故事頁的推進計時（docs/re/009 的教訓）。
	go func() {
		chunk := uint64(200_000)
		minWall := time.Duration(float64(chunk) / *ips * float64(time.Second))
		last := time.Now()
		lastR := false
		lastX, lastY := -1, -1
		paused := false
		for i := 0; ; i++ {
			sh.mu.Lock()
			// spec 003 R2：說明開啟＝不推進、不投遞任何輸入（含游標位置）
			if sh.helpOpen != paused {
				paused = sh.helpOpen
				log.Printf("說明 暫停=%v steps=%d", paused, o.Steps())
			}
			if paused {
				sh.mu.Unlock()
				time.Sleep(16 * time.Millisecond)
				continue
			}
			for _, sc := range sh.keys {
				o.PressKey(oracle.Key(sc))
				log.Printf("送鍵 %02x steps=%d", sc, o.Steps()) // spec 005 驗收：F2（3c）不得出現
			}
			sh.keys = nil
			if len(sh.stdin) > 0 {
				o.Type(string(sh.stdin))
				sh.stdin = nil
			}
			bx, by := sh.mouseX, sh.mouseY
			r := sh.rbtn
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
				// spec 004 R2／R3：與 autoplay 同一 API；錯誤只記 log
				err := o.Click(c.x, c.y)
				log.Printf("左鍵 (%d,%d) steps=%d err=%v", c.x, c.y, o.Steps(), err)
				sh.mu.Lock()
				snapshot() // spec 005 R1：每下點擊後快照，避免兩步併成一次
				sh.mu.Unlock()
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
				snapshot()
			}
			// HD 幀輸出（每 30 秒）：A0000 2x 放大＋文字事件疊畫
			if hdShot && time.Since(lastHD) > 3*time.Second {
				lastHD = time.Now()
				evMu.Lock()
				hdFrame(o, ov, fmt.Sprintf("workplace/bake/hd-%d.png", time.Now().Unix()))
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
			snapshot()
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
