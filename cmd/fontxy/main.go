// fontxy — 抓印字時的座標/顏色：OnCall 13895（字串）＋ 740C:012F（字元暫存器）。
//
//	fontxy -exe <DS.EXE> -root <素材目錄>
package main

import (
	"strings"
	"flag"
	"fmt"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	from := flag.Uint64("from", 520_000_000, "")
	to := flag.Uint64("to", 945_000_000, "")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}

	type strEv struct {
		step uint64
		s    string
	}
	strs := []strEv{}
	prints := 0
	n := 0
	charEvt := oracle.Addr{Seg: 0x740C, Off: 0x012F}

	// OnCall 13895：印字串入口。[sp]=retIP、[sp+2]=字串off、[sp+4]=ID
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(o *oracle.Oracle) {
		off := o.StackWord(1)
		id := o.StackWord(2)
		ds := uint16(0)
		if r := o.Regs(); true {
			ds = r.DS
		}
		s := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 48)
		strs = append(strs, strEv{step: o.Steps(), s: fmt.Sprintf("ds=%04x off=%04x id=%d %q", ds, off, id, s)})
	})

	o.OnCall(charEvt, func(o *oracle.Oracle) {
		n++
		if prints >= 10 {
			return
		}
		r := o.Regs()
		di := int(r.DI)
		fmt.Printf("字元#%d step=%d AX=%04x BX=%04x DX=%04x SI=%04x DI=%04x(%d=%dx%d)\n",
			n, o.Steps(), r.AX, r.BX, r.DX, r.SI, r.DI, di, di%320, di/320)
		// driver 座標 self-modify 變數：loc_294+1（x）、loc_296+1（y）
		// runtime：file offset +0x100 = module → loc_294+1 = module 0x295、loc_296+1 = 0x297
		x := o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295})
		y := o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})
		fmt.Printf("   driver x=%d y=%d → 線性=%d\n", x, y, int(y)*320+int(x))
		// arg_8 裁決：sub_293 的 bp frame 此刻 = thunk 入口（push bp 前），
		// 看 stack 上 [sp+6..0xC] 的字元/ID 殘留
		for i := 2; i <= 12; i += 2 {
			fmt.Printf("   [sp+%d]=%04x", i, o.StackWord(i/2))
		}
		fmt.Println()
		// driver 資料段找座標游標：掃 0x1E0-0x200 區的 words
		for woff := 0x1E0; woff < 0x200; woff += 2 {
			v := o.Word(oracle.Addr{Seg: 0x740C, Off: uint16(woff)})
			if v < 320 && v > 0 {
				fmt.Printf("   740C:%04x = %d（候選 x/y）\n", woff, v)
			}
		}
		last := ""
		if len(strs) > 0 {
			last = strs[len(strs)-1].s
		}
		fmt.Printf("   最近字串： %s\n", last)
		prints++
	})

	for total := *from; total < *to; {
		k := uint64(2_000_000)
		if rest := *to - total; rest < k {
			k = rest
		}
		if err := o.Run(k); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += k
	}
	fmt.Println("字串事件：", len(strs), "；字元事件：", n)
	uniq := map[string]bool{}
	for _, e := range strs {
		s := e.s
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		if !uniq[s] {
			uniq[s] = true
			fmt.Printf("  %d %s\n", e.step, s)
		}
	}
}
