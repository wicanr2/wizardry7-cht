// ovlprobe — #31 證據探針：不寫原版，量覆繪需要的四件事。
//
//	P1 離屏緩衝（blit 來源段）與 A000 的關係
//	P2 sub_13895（0110:3895）的呼叫形態與返回點
//	P3 每次印字在離屏改了哪個矩形（字高）
//	P4 前景色來源（服務 11 入口 740C:012F 的 DL）
//
//	ovlprobe -exe <DS.EXE> -root <素材目錄> -stop <步數> -from <步數> -clicks step:x:y,...
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

type click struct {
	at   uint64
	x, y int
}

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	stop := flag.Uint64("stop", 286_000_000, "停止步數")
	from := flag.Uint64("from", 0, "只記錄這個步數之後的印字事件")
	clicksF := flag.String("clicks", "", "step:x:y 逗號分隔")
	msgFn := flag.Bool("msg", false, "另記訊息取得入口 0110:0582（sub_10582，docs/re/013 追記二十四）的參數")
	typeF := flag.String("type-at", "", "step:字串（以 DOS stdin 輸入；\\r 代表 Enter）")
	flag.Parse()

	var clicks []click
	for _, c := range strings.Split(*clicksF, ",") {
		var k click
		if _, err := fmt.Sscanf(c, "%d:%d:%d", &k.at, &k.x, &k.y); err == nil {
			clicks = append(clicks, k)
		}
	}

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}

	if *msgFn {
		o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x0582}, func(o *oracle.Oracle) {
			if o.Steps() >= *from {
				fmt.Printf("訊息 %d id=%d arg2=%d ret=%04x\n", o.Steps(), o.StackWord(1), o.StackWord(2), o.StackWord(0))
			}
		})
	}
	var scr uint16 // 最近一次 blit 的來源段
	blitDS := map[uint16]int{}
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x05F5}, func(o *oracle.Oracle) {
		if r := o.Regs(); r.DS != 0xA000 {
			scr = r.DS
			blitDS[r.DS]++
		}
	})

	offscr := func() []byte {
		if scr == 0 {
			return nil
		}
		return o.Bytes(oracle.Addr{Seg: scr, Off: 0}, 320*200)
	}

	type pend struct {
		step       uint64
		s          string
		x0, y0     int
		ret        oracle.Addr
		sp         uint16
		font       uint16
		before     []byte
		fgs        []uint8
		chars      int
		es         map[uint16]bool
	}
	var cur *pend
	retHooked := map[uint32]bool{}

	// P4：服務 11 入口的 DL（前景）與 ES（寫入段）
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x012F}, func(o *oracle.Oracle) {
		if cur == nil {
			return
		}
		r := o.Regs()
		cur.chars++
		if len(cur.fgs) < 4 {
			cur.fgs = append(cur.fgs, uint8(r.DX&0xFF))
		}
		cur.es[r.ES] = true
	})

	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	o.OnCall(entry, func(o *oracle.Oracle) {
		if o.Steps() < *from {
			return
		}
		r := o.Regs()
		off := o.StackWord(1)
		raw := o.Bytes(oracle.Addr{Seg: r.DS, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		p := &pend{
			step: o.Steps(), s: s,
			x0:   int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295})),
			y0:   int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})),
			ret:  oracle.Addr{Seg: r.CS, Off: o.StackWord(0)},
			sp:   r.SP,
			font: o.StackWord(2),
			es:   map[uint16]bool{},
		}
		if b := offscr(); b != nil {
			p.before = append([]byte(nil), b...)
		}
		cur = p
		lin := p.ret.Linear()
		if retHooked[lin] {
			return
		}
		retHooked[lin] = true
		o.OnCall(p.ret, func(o *oracle.Oracle) {
			// near ret 後 SP＝進入時 SP＋2（只彈 IP；參數由呼叫端清）
			if cur == nil || o.Regs().SP != cur.sp+2 || o.IP().Linear() != cur.ret.Linear() {
				return
			}
			q := cur
			cur = nil
			x1 := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
			y1 := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
			bx0, by0, bx1, by1, n := 320, 200, -1, -1, 0
			if after := offscr(); after != nil && q.before != nil {
				for i := range after {
					if after[i] != q.before[i] {
						n++
						x, y := i%320, i/320
						bx0, by0 = min(bx0, x), min(by0, y)
						bx1, by1 = max(bx1, x), max(by1, y)
					}
				}
			}
			ess := []string{}
			for e := range q.es {
				ess = append(ess, fmt.Sprintf("%04x", e))
			}
			fmt.Printf("印字 %d %q font=%d ret=%04x:%04x xy0=(%d,%d) xy1=(%d,%d) 字數=%d fg=% x es=%v 離屏差=%d bbox=(%d,%d)-(%d,%d)\n",
				q.step, q.s, q.font, q.ret.Seg, q.ret.Off, q.x0, q.y0, x1, y1, q.chars, q.fgs, ess, n, bx0, by0, bx1, by1)
		})
	})

	var typed []struct {
		at uint64
		s  string
	}
	if *typeF != "" {
		var at uint64
		var str string
		if _, err := fmt.Sscanf(*typeF, "%d:%s", &at, &str); err == nil {
			typed = append(typed, struct {
				at uint64
				s  string
			}{at, strings.ReplaceAll(str, `\r`, "\r")})
		}
	}
	ti := 0
	ci := 0
	for o.Steps() < *stop {
		for ti < len(typed) && o.Steps() >= typed[ti].at {
			o.Type(typed[ti].s)
			ti++
		}
		for ci < len(clicks) && o.Steps() >= clicks[ci].at {
			if err := o.Click(clicks[ci].x, clicks[ci].y); err != nil {
				fmt.Printf("點擊 (%d,%d)：%v\n", clicks[ci].x, clicks[ci].y, err)
			}
			ci++
		}
		n := uint64(1_000_000)
		if rest := *stop - o.Steps(); rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
	}

	// P2：入口機器碼
	fmt.Printf("0110:3895 碼：% x\n", o.Bytes(entry, 48))
	fmt.Printf("blit 來源段：%v；word 740C:01E5=%04x、740C:02E5=%04x\n", blitDS,
		o.Word(oracle.Addr{Seg: 0x740C, Off: 0x01E5}), o.Word(oracle.Addr{Seg: 0x740C, Off: 0x02E5}))
	// P1：離屏與 A000 逐像素比較
	a := o.Indexed()
	if b := offscr(); b != nil {
		n := 0
		bx0, by0, bx1, by1 := 320, 200, -1, -1
		for i := range a {
			if a[i] != b[i] {
				n++
				x, y := i%320, i/320
				bx0, by0 = min(bx0, x), min(by0, y)
				bx1, by1 = max(bx1, x), max(by1, y)
			}
		}
		pal := o.Palette()
		save := func(src []byte, path string) {
			img := image.NewRGBA(image.Rect(0, 0, 320, 200))
			for i, v := range src {
				c := pal[v]
				img.Set(i%320, i/320, color.RGBA{c[0], c[1], c[2], 255})
			}
			f, _ := os.Create(path)
			png.Encode(f, img)
			f.Close()
		}
		// 對話框區（y≥140，兩邊都該是同一張對話框）的差異＝A000 上的游標
		cx0, cy0, cx1, cy1 := 320, 200, -1, -1
		for i := 140 * 320; i < len(a); i++ {
			if a[i] != b[i] {
				x, y := i%320, i/320
				cx0, cy0 = min(cx0, x), min(cy0, y)
				cx1, cy1 = max(cx1, x), max(cy1, y)
			}
		}
		fmt.Printf("對話框區差異 bbox=(%d,%d)-(%d,%d)\n", cx0, cy0, cx1, cy1)
		for y := cy0; y <= cy1 && cx1 >= 0; y++ {
			row := make([]byte, 0, cx1-cx0+1)
			for x := cx0; x <= cx1; x++ {
				if a[y*320+x] != b[y*320+x] {
					row = append(row, '#')
				} else {
					row = append(row, '.')
				}
			}
			fmt.Printf("  %3d %s\n", y, row)
		}
		save(b, "workplace/bake/r663-offscreen.png")
		save(a, "workplace/bake/r663-a000.png")
		mx, my := o.Mouse()
		fmt.Printf("離屏 %04x vs A000：差 %d 像素 bbox=(%d,%d)-(%d,%d)；滑鼠=(%d,%d)\n", scr, n, bx0, by0, bx1, by1, mx, my)
	}
}
