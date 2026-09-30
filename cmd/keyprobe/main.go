// keyprobe — #34 證據探針：原版是否使用某個鍵（F1）。
//
// 在指定步數存快照，分三個變體各跑 n 步：不按鍵、按受測鍵、按對照鍵（正對照）。
// 比較 A000 與不按鍵變體的差異像素數，以及期間印字（0110:3895）次數。
// 送鍵路徑二選一：-path irq（oracle.PressKey，IRQ1 掃描碼）或 -path stdin
// （oracle.Type，DOS stdin 位元組；擴充鍵是 00＋掃描碼，F1＝003b）。
//
//	keyprobe -exe <DS.EXE> -root <素材目錄> -clicks step:x:y,... -at 272000000 -path stdin -key 003b -ctl 41
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	clicksF := flag.String("clicks", "", "step:x:y 逗號分隔")
	atF := flag.String("at", "", "量測起點步數，逗號分隔（遞增）")
	path := flag.String("path", "stdin", "送鍵路徑：irq 或 stdin")
	keyF := flag.String("key", "003b", "受測鍵（16 進位；irq 為一個掃描碼，stdin 為位元組串）")
	ctlF := flag.String("ctl", "41", "對照鍵（格式同 -key，預期會有反應）")
	n := flag.Uint64("n", 30_000_000, "每個變體跑幾步")
	flag.Parse()

	type click struct {
		at   uint64
		x, y int
	}
	var clicks []click
	for _, c := range strings.Split(*clicksF, ",") {
		var k click
		if _, err := fmt.Sscanf(c, "%d:%d:%d", &k.at, &k.x, &k.y); err == nil {
			clicks = append(clicks, k)
		}
	}
	var ats []uint64
	for _, f := range strings.Split(*atF, ",") {
		if v, err := strconv.ParseUint(strings.TrimSpace(f), 10, 64); err == nil {
			ats = append(ats, v)
		}
	}
	hexBytes := func(h string) []byte {
		b := make([]byte, len(h)/2)
		for i := range b {
			v, _ := strconv.ParseUint(h[i*2:i*2+2], 16, 8)
			b[i] = byte(v)
		}
		return b
	}
	key, ctl := hexBytes(*keyF), hexBytes(*ctlF)

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	prints := 0
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(*oracle.Oracle) { prints++ })

	ci := 0
	runTo := func(stop uint64) {
		for o.Steps() < stop {
			for ci < len(clicks) && o.Steps() >= clicks[ci].at {
				_ = o.Click(clicks[ci].x, clicks[ci].y)
				ci++
			}
			step := uint64(1_000_000)
			if rest := stop - o.Steps(); rest < step {
				step = rest
			}
			if err := o.Run(step); err != nil {
				panic(err)
			}
		}
	}

	for _, at := range ats {
		runTo(at)
		snap := o.Save()
		start := o.Steps()
		type res struct {
			a000     []uint8
			prints   int
			consumed int
		}
		variant := func(b []byte) res {
			o.Restore(snap)
			o.ClearInput()
			p0, k0 := prints, o.Pending()
			switch {
			case len(b) == 0:
			case *path == "irq":
				o.PressKey(oracle.Key(b[0]))
			default:
				o.Type(string(b))
				k0 = o.Pending()
			}
			for o.Steps() < start+*n {
				step := uint64(1_000_000)
				if rest := start + *n - o.Steps(); rest < step {
					step = rest
				}
				if err := o.Run(step); err != nil {
					panic(err)
				}
			}
			return res{append([]uint8(nil), o.Indexed()...), prints - p0, k0 - o.Pending()}
		}
		base := variant(nil)
		diff := func(r res) int {
			d := 0
			for i := range r.a000 {
				if r.a000[i] != base.a000[i] {
					d++
				}
			}
			return d
		}
		k := variant(key)
		c := variant(ctl)
		fmt.Printf("起點 %d（+%d 步，%s）：無鍵 印字=%d｜鍵 %s 讀走=%d 差異=%d 印字=%d｜對照 %s 讀走=%d 差異=%d 印字=%d\n",
			at, *n, *path, base.prints, *keyF, k.consumed, diff(k), k.prints, *ctlF, c.consumed, diff(c), c.prints)
		o.Restore(snap)
	}
}
