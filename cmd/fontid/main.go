// fontid — 找主選單文字用的字型 ID：OnCall sub_13877 記錄（字元,ID）直方圖。
//
//	fontid -exe <DS.EXE> -root <素材目錄>
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
	from := flag.Uint64("from", 930_000_000, "")
	to := flag.Uint64("to", 945_000_000, "")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	type rec struct {
		n     int
		chars map[byte]int
	}
	hist := map[byte]*rec{} // ID → 呼叫次數＋字元分布
	entry := oracle.Addr{Seg: 0x0110, Off: 0x3877}
	o.OnCall(entry, func(o *oracle.Oracle) {
		// 進入口時 [bp] 還是呼叫者 BP——參數在 [sp+2](字元) [sp+4](ID)
		sp := o.StackWord(0)
		_ = sp
		// StackWord(i) 讀 SS:SP+i*2；呼叫 sub_13877 後 [sp]=retIP、[sp+2]=字元、[sp+4]=ID
		ch := byte(o.StackWord(1))
		id := byte(o.StackWord(2))
		r := hist[byte(id)]
		if r == nil {
			r = &rec{chars: map[byte]int{}}
			hist[id] = r
		}
		r.n++
		r.chars[ch]++
	})
	for total := *from; total < *to; {
		n := uint64(2_000_000)
		if rest := *to - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += n
	}
	ids := []int{}
	for idb := range hist {
		id := int(idb)
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for _, id := range ids {
		r := hist[byte(id)]
		fmt.Printf("字型 ID %d：%d 次呼叫\n", id, r.n)
		chars := []byte{}
		for c := range r.chars {
			chars = append(chars, c)
		}
		sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })
		s := ""
		for _, c := range chars {
			s += string(rune(c))
		}
		fmt.Printf("  字元（%d 種）：%q\n", len(chars), s)
	}
}
