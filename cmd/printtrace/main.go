// printtrace — 巫術7 M3：找出印字常式。
//
// 在 VGA.DRV 的 blitter 入口（執行期 740C:05B5）掛 OnCall，收集呼叫端
// CS:IP 與區塊參數的直方圖：底圖＝大區塊少呼叫，文字字形＝小區塊連發。
//
//	printtrace -exe <DS.EXE> -root <素材目錄> [-steps 450000000]
package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	steps := flag.Uint64("steps", 450_000_000, "")
	flag.Parse()
	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}

	type rec struct {
		n              int
		w, h           [8]uint16
		x, y           [8]uint16
		nw, nh, nx, ny int
	}
	callers := map[oracle.Addr]*rec{}
	var order []oracle.Addr

	blit := oracle.Addr{Seg: 0x740C, Off: 0x05B5}
	o.OnCall(blit, func(o *oracle.Oracle) {
		c := o.Caller()
		a0, a1, a2, a3 := o.Arg(0), o.Arg(1), o.Arg(2), o.Arg(3)
		r := callers[c]
		if r == nil {
			r = &rec{}
			callers[c] = r
			order = append(order, c)
		}
		r.n++
		if r.nw < len(r.w) {
			r.w[r.nw], r.h[r.nh] = a2, a3
			r.x[r.nx], r.y[r.ny] = a0, a1
			r.nw++
			r.nh++
			r.nx++
			r.ny++
		}
	})

	total := uint64(0)
	for total < *steps {
		n := uint64(2_000_000)
		if rest := *steps - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Printf("停止於 steps=%d：%v\n", o.Steps(), err)
			break
		}
		total += n
		fmt.Printf("進度 %d/%d steps=%d 呼叫端 %d 種\n", total, *steps, o.Steps(), len(callers))
	}

	sort.Slice(order, func(i, j int) bool { return callers[order[i]].n > callers[order[j]].n })
	fmt.Println("\n== blitter 呼叫端直方圖 ==")
	for _, c := range order {
		r := callers[c]
		fmt.Printf("%v  %6d 次  樣本 w/h: %v/%v  x/y: %v/%v\n",
			c, r.n, r.w[:r.nw], r.h[:r.nh], r.x[:r.nx], r.y[:r.ny])
	}
}
