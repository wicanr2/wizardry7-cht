// wiz7 是《巫術7》的可遊玩前端（#15）：原版 DS.EXE 跑在 dosgolem 上，
// 開視窗、接鍵盤與滑鼠，速度以「每秒指令數」對齊（預設 10M，
// 對齊 DOSBox-X cycles=20000 的體感：開機到主選單約 90 秒）。
//
//	wiz7 -root workplace/run [-scale 3] [-ips 10000000]
//
// 畫面是 mode 13h（320×200、256 色、A0000 線性），oracle.ScreenRGB
// 連 DAC 一起合成；放大用 nearest，不縮字（rulebook/81）。
// 音訊殼：本版只留介面（Render 取樣流接 Ebiten audio），合成器屬 M2 後續。
package main

import (
	"flag"
	"fmt"
	"log"
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

// pacer 是 lag 式節拍：目標＝牆上毫秒×每毫秒步數，落後太多就放棄追趕。
// 巫術7 的動畫由計時器推進，指令數↔計時器的換算固定，所以同一份輸入
// 在任何速度下逐位元組相同。
type pacer struct {
	perMs     float64 // 每毫秒指令數
	droppedMs float64
}

func (p *pacer) cycles(wallMs, machineSteps float64) uint64 {
	lag := wallMs - p.droppedMs - machineSteps/p.perMs
	if lag <= 0 {
		return 0
	}
	const maxCatchUpMs = 250
	if lag > maxCatchUpMs {
		p.droppedMs += lag - maxCatchUpMs
		lag = maxCatchUpMs
	}
	return uint64(lag * p.perMs)
}

type game struct {
	o      *oracle.Oracle
	p      pacer
	scale  int
	screen *ebiten.Image
	rgba   []byte
	start  time.Time
	startS uint64
	down   map[ebiten.Key]bool
}

func (g *game) Update() error {
	if g.start.IsZero() {
		g.start, g.startS = time.Now(), g.o.Steps()
	}
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if sc, ok := scanCode(k); ok {
			g.o.PressKey(sc)
			g.down[k] = true
		}
	}
	for _, k := range inpututil.AppendJustReleasedKeys(nil) {
		if sc, ok := scanCode(k); ok {
			g.o.PressKey(sc | 0x80) // break code
			delete(g.down, k)
		}
	}
	// 滑鼠：位置逐幀回報；按住／放開用 PressButton／ReleaseButton。
	mx, my := ebiten.CursorPosition()
	_ = g.o.MoveMouse(mx/g.scale, my/g.scale)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.o.PressButton(0)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.o.ReleaseButton(0)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.o.PressButton(1)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.o.ReleaseButton(1)
	}
	// 節拍：目標步數＝牆上毫秒×每毫秒步數。
	wallMs := float64(time.Since(g.start).Microseconds()) / 1000
	machineSteps := float64(g.o.Steps() - g.startS)
	if n := g.p.cycles(wallMs, machineSteps); n > 0 {
		if err := g.o.Run(n); err != nil {
			log.Printf("機器停止：%v", err)
			return ebiten.Termination
		}
	}
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	w, h, rgb := g.o.ScreenRGB()
	if w != 320 || h != 200 {
		return
	}
	for i := 0; i < w*h; i++ {
		g.rgba[4*i], g.rgba[4*i+1], g.rgba[4*i+2], g.rgba[4*i+3] = rgb[3*i], rgb[3*i+1], rgb[3*i+2], 0xFF
	}
	g.screen.WritePixels(g.rgba)
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(float64(g.scale), float64(g.scale))
	op.Filter = ebiten.FilterNearest
	dst.DrawImage(g.screen, &op)
}

func (g *game) Layout(int, int) (int, int) { return 320 * g.scale, 200 * g.scale }

func main() {
	root := flag.String("root", "workplace/run", "原版素材目錄（可寫複本）")
	exe := flag.String("exe", "", "執行檔（預設 <root>/DS.EXE）")
	scale := flag.Int("scale", 3, "畫面放大倍率")
	ips := flag.Float64("ips", 10_000_000, "每秒指令數（節拍基準）")
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
	g := &game{
		o:     o,
		p:     pacer{perMs: *ips / 1000},
		scale: *scale,
		rgba:  make([]byte, 320*200*4),
		down:  map[ebiten.Key]bool{},
	}
	g.screen = ebiten.NewImage(320, 200)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	_ = fmt.Sprint()
}
