// msgdump — MSG.DBS runtime 對照法：監看 stack 字串區寫入、找旁白明文源頭。
//
//	msgdump -exe <DS.EXE> -root <素材目錄>
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}

	// 監看 stack 字串區（旁白印字時 ds:off 落在 0x210DC 附近）——
	// 掛在推進之前，抓「誰把明文搬進來」。
	writes := 0
	o.OnWrite(0x20E00, 0x21200, func(o *oracle.Oracle, h oracle.WriteHit) {
		if o.Steps() < 200_000_000 {
			return // 開機期雜訊
		}
		if h.Val < 0x20 || h.Val > 0x7E {
			return // 只看 ASCII 明文寫入
		}
		if writes >= 200 {
			return
		}
		writes++
		seg := uint16(h.Addr >> 4)
		off := uint16(h.Addr & 0xF)
		fmt.Printf("寫入 #%d step=%d 線性=%#x 舊=%02x 新=%02x by IDA=%#x\n",
			writes, o.Steps(), h.Addr, h.Old, h.Val, h.At)
		raw := o.Bytes(oracle.Addr{Seg: seg, Off: off &^ 0xF}, 64)
		vis := strings.Map(func(r rune) rune {
			if r < 32 || r > 126 {
				return '·'
			}
			return r
		}, string(raw))
		fmt.Printf("   內容 %04x:%04x %q\n", seg, off&^0xF, vis)
	})

	// 印字明文 log（對齊 FileOps 的 MSG.DBS 讀取 step）
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(o *oracle.Oracle) {
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		fmt.Printf("PRINT step=%d %q\n", o.Steps(), s)
	})
	// 解壓函式入口監看：3E17 = mov ax,cs:[2D4E]（ax=樹段）、si=來源
	caps := 0
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3E1B}, func(o *oracle.Oracle) {
		if caps >= 5 {
			return
		}
		caps++
		r := o.Regs()
		fmt.Printf("解壓#%d step=%d si=%04x ax(樹段)=%04x ds=%04x\n",
			caps, o.Steps(), r.SI, r.AX, r.DS)
		// dump 來源區（ds:si 起 128 bytes）
		raw := o.Bytes(oracle.Addr{Seg: r.DS, Off: r.SI}, 128)
		vis := strings.Map(func(rr rune) rune {
			if rr < 32 || rr > 126 {
				return '·'
			}
			return rr
		}, string(raw))
		fmt.Printf("   來源 %04x:%04x % x\n", r.DS, r.SI, raw[:64])
		fmt.Printf("   可讀 %q\n", vis)
	})

	// 建隊 + START NEW GAME + 防拷 + 旁白
	seq := []struct{ x, y int }{{160, 62}, {170, 153}, {90, 168}, {90, 168}, {160, 101}}
	for _, s := range seq {
		_ = o.Click(s.x, s.y)
		_ = o.Run(120_000_000)
	}
	o.Type("AAAA\r")
	for i := 0; i < 24; i++ {
		_ = o.Click(25, 15)
		_ = o.Run(150_000_000)
	}
	_ = o.Run(200_000_000)
	fmt.Println("steps =", o.Steps(), "監看命中：", writes)

	// 樹段：解壓函式讀 cs:[0110:2D4E] 當 es
	ts := o.Word(oracle.Addr{Seg: 0x0110, Off: 0x2D4E})
	fmt.Printf("樹段 cs:[2D4E] = %04x\n", ts)
	if ts != 0 {
		tree := o.Bytes(oracle.Addr{Seg: ts, Off: 0}, 8192)
		_ = os.WriteFile("/w/workplace/bake/tree.bin", tree, 0o644)
		fmt.Println("樹區 8KB → workplace/bake/tree.bin")
		// 前幾個節點
		for i := 0; i < 8; i++ {
			w := uint16(tree[i*2]) | uint16(tree[i*2+1])<<8
			fmt.Printf("  node[%d] = %04x（bit15=%d）\n", i, w, w>>15)
		}
	}
	fmt.Println("== MSG.DBS 檔案操作 ==")
	for _, fop := range o.FileOps() {
		if strings.Contains(strings.ToUpper(fop.Name), "MSG") {
			fmt.Printf("  step=%d %s %s pos=%#x arg=%#x len=%d\n", fop.Step, fop.Op, fop.Name, fop.Pos, fop.Arg, fop.Len)
		}
	}

	// 旁白明文搜尋
	probes := []string{
		"Inside the ruins of an abandoned temple",
		"I am Aletheides",
		"The Dark Savant...",
		"CHARACTER MENU",
	}
	for _, p := range probes {
		hits := o.Search([]byte(p))
		fmt.Printf("%q：%d 命中 %v\n", p, len(hits), hits)
	}
}
