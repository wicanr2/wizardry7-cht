// xlate1 — 巫術7 M4 方案 B 雛形：攔印字串、譯文畫中文。
//
// 機制（docs/re/013 追記三方案 B）：
//   OnCall(0110:3895)（印字串入口，hook 在第一道指令前）：
//     1. 先 unpoke（寫回原 byte，讓本字串照常印，除非有譯文）
//     2. 讀字串（ds:off，NUL 終止）查譯文表
//     3. 有譯文：driver 座標變數（x=740C:0295、y=0x0297）取位置，
//        Eten 15×15 寫離屏段（blit 的來源，非筆畫保原像素），
//        再 poke 0xC3（ret）到 3895 讓原印失效
//
//	xlate1 -exe <DS.EXE> -root <素材目錄>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

// 譯文表：key = 字串內容（stack 字串無固定位址，內容即 key；rulebook 的
// 來源位置 key 適用 MSG.DBS 抽取文本，這裡是 UI 組合字串，標 DRAFT）。
// 值 = 每字 15 個 row（MSB first、15px 寬），由 tools/re/bake_font.py 產出。
type entry struct {
	Dy int         `json:"dy"`
	G  [][]uint16  `json:"g"`
}

var trans = map[string]entry{}

func loadTrans(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &trans); err != nil {
		panic(err)
	}
}

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	stop1 := flag.Uint64("stop1", 240_000_000, "")
	stop2 := flag.Uint64("stop2", 940_000_000, "")
	collect := flag.String("collect", "", "無譯文串寫到這個檔（一行一串）")
	dy := flag.Int("dy", 0, "中文 y 微調像素")
	clicks := flag.String("clicks", "", "step:x:y 逗號分隔（絕對步數時點擊）")
	keysAt := flag.String("keys-at", "", "step:key 逗號分隔（Esc=0x1b）")
	typeAt := flag.String("type-at", "", "step:字串（於指定步數輸入字串+CR）")
	rawAt := flag.String("raw-at", "", "step:hexbytes（原始 bytes 排 Stdin，方向鍵＝00 掃描碼）")
	irqAt := flag.String("irq-at", "", "step:鍵名（IRQ1 硬體掃描碼路徑：Up/Down/Left/Right/w…）")
	shotsAt := flag.String("shots-at", "", "step 逗號分隔（純截圖）")
	bkeysAt := flag.String("bkeys-at", "", "step:鍵名 逗號分隔（BIOS 緩衝路徑，如 step:Up）")
	eventsFrom := flag.Uint64("events-from", 0, "只合成這個步數之後的文字事件（排除先前畫面的殘影）")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	loadTrans("workplace/bake/xlate.json")

	for total := uint64(0); total < *stop1; {
		n := uint64(2_000_000)
		if rest := *stop1 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}

	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	baseline := o.Byte(entry) // 原 3895 第一 byte（push bp = 0x55）
	poked := false
	var lastFG uint8 = 0x22
	var lastScr uint16 = 0
	scrs := map[uint16]bool{}
	drawn := 0
	collected := map[string]bool{}
	type drawEvent struct {
		at   uint64
		x, y int
		fg   uint8
		s    string
		g    [][]uint16
	}
	var events []drawEvent

	o.OnCall(entry, func(o *oracle.Oracle) {
		if poked { // 先解除，本字串（無論譯不譯）從乾淨狀態開始
			o.WriteU8(entry, baseline)
			poked = false
		}
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		ent, ok := trans[s]
		if !ok {
			if *collect != "" && !collected[s] {
				collected[s] = true
				f, _ := os.OpenFile(*collect, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
				f.WriteString(s + "\n")
				f.Close()
			}
			return // 無譯文：放行
		}
		// 位置：driver 座標變數；畫到離屏（blit 的來源段）
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})) + *dy + ent.Dy
		scr := lastScr
		fg := lastFG
		fmt.Printf("翻譯 %d %q → %d 字形 @(%d,%d) fg=%02x 離屏=%04x\n", o.Steps(), s, len(ent.G), x, y, fg, scr)

		// 高解析模式：不寫進 320x200 遊戲畫面（會被翻頁清掉、字比按鈕大），
		// 記錄事件，stop2 之後在 2x 放大的終幀上以 15x15 原尺寸合成。
		events = append(events, drawEvent{at: o.Steps(), x: x, y: y, fg: fg, s: s, g: ent.G})
		drawn++
		o.WriteU8(entry, 0xC3) // ret：原印失效
		poked = true
	})

	// 記住前景色（原印時的字元事件）與離屏段（blit 的來源段）
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x012F}, func(o *oracle.Oracle) {
		lastFG = uint8(o.Regs().DX & 0xFF)
	})
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x05F5}, func(o *oracle.Oracle) {
		if r := o.Regs(); r.DS > 0 && r.DS != 0xA000 {
			lastScr = r.DS // blit 來源段＝離屏緩衝
			scrs[r.DS] = true
		}
	})

	type clk struct{ at uint64; x, y int }
	var clicksQ []clk
	var keysQ []struct {
		at uint64
		b  byte
	}
	var shotsQ []uint64
	for _, k := range strings.Split(*shotsAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		if _, err := fmt.Sscanf(k, "%d", &at); err == nil {
			shotsQ = append(shotsQ, at)
		}
	}
	var bkeysQ []struct {
		at uint64
		n  string
	}
	for _, k := range strings.Split(*bkeysAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var n string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &n); err == nil {
			bkeysQ = append(bkeysQ, struct {
				at uint64
				n  string
			}{at, n})
		}
	}
	var irqQ []struct {
		at  uint64
		scan byte
	}
	for _, k := range strings.Split(*irqAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var hs string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &hs); err == nil {
			var v int
			fmt.Sscanf(hs, "%x", &v)
			irqQ = append(irqQ, struct {
				at   uint64
				scan byte
			}{at, byte(v)})
		}
	}
	var rawQ []struct {
		at uint64
		b  []byte
	}
	for _, k := range strings.Split(*rawAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var hs string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &hs); err == nil {
			hexstr := strings.ReplaceAll(hs, " ", "")
			b := make([]byte, len(hexstr)/2)
			for i := 0; i < len(b); i++ {
				fmt.Sscanf(hexstr[i*2:i*2+2], "%02x", &b[i])
			}
			rawQ = append(rawQ, struct {
				at uint64
				b  []byte
			}{at, b})
		}
	}
	var typeQ []struct {
		at uint64
		s  string
	}
	for _, k := range strings.Split(*typeAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var s string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &s); err == nil {
			typeQ = append(typeQ, struct {
				at uint64
				s  string
			}{at, s})
		}
	}
	for _, k := range strings.Split(*keysAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var hb string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &hb); err == nil {
			var v int
			fmt.Sscanf(hb, "%x", &v)
			keysQ = append(keysQ, struct {
				at uint64
				b  byte
			}{at, byte(v)})
		}
	}
	for _, c := range strings.Split(*clicks, ",") {
		if c == "" {
			continue
		}
		var at uint64
		var x, y int
		if _, err := fmt.Sscanf(c, "%d:%d:%d", &at, &x, &y); err == nil {
			clicksQ = append(clicksQ, clk{at, x, y})
		}
	}
	for total := *stop1; total < *stop2; {
		n := uint64(2_000_000)
		if rest := *stop2 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += n
		for _, tq := range typeQ {
			if tq.at <= total && tq.at > total-n {
				fmt.Printf("輸入 %q @%d\n", tq.s, tq.at)
				o.Type(tq.s + "\r")
			}
		}
		for _, sa := range shotsQ {
			if sa <= total && sa > total-n {
				_ = o.Run(1_000_000)
				_ = o.WritePNG(fmt.Sprintf("workplace/bake/shot-%d.png", sa))
			}
		}
		for _, bq := range bkeysQ {
			if bq.at <= total && bq.at > total-n {
				fmt.Printf("BIOS鍵 %q @%d\n", bq.n, bq.at)
				if bq.n == "w" || bq.n == "a" || bq.n == "s" || bq.n == "d" || bq.n == "x" {
					if err := o.TypeKeys(bq.n); err != nil {
						fmt.Println("TypeKeys：", err)
					}
				} else if o.SendKeys(bq.n) != nil {
					fmt.Println("鍵名失敗：", bq.n)
				}
			}
		}
		for _, rq := range rawQ {
			if rq.at <= total && rq.at > total-n {
				fmt.Printf("Raw %x @%d\n", rq.b, rq.at)
				o.Type(string(rq.b))
			}
		}
		for _, iq := range irqQ {
			if iq.at <= total && iq.at > total-n {
				fmt.Printf("IRQ1掃描碼 %02x @%d\n", iq.scan, iq.at)
				o.PressKey(oracle.Key(iq.scan))
			}
		}
		for _, kk := range keysQ {
			if kk.at <= total && kk.at > total-n {
				fmt.Printf("按鍵 0x%02x @%d\n", kk.b, kk.at)
				o.Type(string([]byte{kk.b}))
				_ = o.Run(3_000_000)
				_ = o.WritePNG(fmt.Sprintf("workplace/bake/key-%d.png", kk.at))
			}
		}
		for _, c := range clicksQ {
			if c.at <= total && c.at > total-n {
				fmt.Printf("點擊 (%d,%d) @%d\n", c.x, c.y, c.at)
				if err := o.Click(c.x, c.y); err != nil {
					fmt.Println("點擊失敗：", err)
				}
				// 推進一點讓畫面穩定，存 PNG
				_ = o.Run(3_000_000)
				png := fmt.Sprintf("workplace/bake/click-%d.png", c.at)
				if err := o.WritePNG(png); err != nil {
					fmt.Println("PNG：", err)
				}
			}
		}
	}
	fmt.Println("畫了", drawn, "串；見過的離屏段：", scrs)
	// 監看 IVT int08h（linear 0x20-0x23）的寫入——誰拔了 PLAY.DRV 的 tick handler
	o.OnWrite(0x20, 0x27, func(o *oracle.Oracle, h oracle.WriteHit) {
		if o.Steps() > 2200000000 {
			fmt.Printf("IVT08 寫 step=%d 線性=%#x 舊=%02x 新=%02x by IDA %#x\n",
				o.Steps(), h.Addr, h.Old, h.Val, h.At)
		}
	})
	for _, vs := range o.VecSets() {
		if vs.Int == 8 || vs.Int == 9 {
			fmt.Printf("  AH=25h 設 int%02x → %04x:%04x @step %d\n", vs.Int, vs.Seg, vs.Off, vs.Step)
		}
	}
	fmt.Println("== 缺檔 ==", o.Missing())
	fmt.Println("== 未實作 ==", o.Unimplemented())
	fmt.Println("== 開過的檔 ==", o.Opened())
	// IVT 檢查：int 08h/09h/16h/21h 向量現在指哪
	for _, n := range []struct {
		v   uint32
		nam string
	}{{8, "int08h"}, {9, "int09h"}, {0x10, "int10h"}, {0x16, "int16h"}, {0x21, "int21h"}} {
		vec := o.Word(oracle.Addr{Seg: 0, Off: uint16(n.v * 4)})
		seg := o.Word(oracle.Addr{Seg: 0, Off: uint16(n.v*4 + 2)})
		fmt.Printf("  %s → %04x:%04x\n", n.nam, seg, vec)
	}
	// 高解析終幀：A0000（320x200 色號）→ 2x 最近鄰放大（640x400），
	// 再把翻譯事件以 15x15 原字形畫在 (x*2, y*2)——中文相對畫面變小、
	// 剛好放進按鈕裡，且不會被遊戲翻頁清掉。
	pal := o.Palette()
	idx := o.Indexed()
	big := image.NewRGBA(image.Rect(0, 0, 640, 400))
	for py := 0; py < 400; py++ {
		for px := 0; px < 640; px++ {
			c := pal[idx[(py/2)*320+(px/2)]]
			big.Set(px, py, color.RGBA{c[0], c[1], c[2], 255})
		}
	}
	for _, e := range events {
		if e.at < *eventsFrom {
			continue // 先前畫面的事件，底圖已換，合成會張冠李戴
		}
		fr := pal[e.fg]
		fgc := color.RGBA{fr[0], fr[1], fr[2], 255}
		cx := e.x * 2
		for _, glyph := range e.g {
			for row := 0; row < len(glyph); row++ {
				w := glyph[row]
				for bit := 0; bit < len(glyph); bit++ {
					if (w>>(len(glyph)-1-bit))&1 != 0 {
						big.Set(cx+bit, e.y*2+row, fgc)
					}
				}
			}
			cx += len(glyph) + 2
		}
	}
	f, _ := os.Create("workplace/bake/xlate1-hd.png")
	png.Encode(f, big)
	f.Close()
	fmt.Println("輸出 workplace/bake/xlate1-hd.png（640x400 高解析合成）")
	if err := o.WritePNG("workplace/bake/xlate1.png"); err != nil {
		panic(err)
	}
	_ = os.WriteFile("workplace/bake/xlate1.vram", o.Indexed(), 0o644)
	fmt.Println("輸出 workplace/bake/xlate1.png")
}
