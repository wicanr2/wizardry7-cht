// copyprot — 防拷檢查（SECURITY CLEARANCE）RE：輸入猜測、觀察反應。
//
//	copyprot -exe <DS.EXE> -root <素材目錄> [-word WIZARDRY]
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	word := flag.String("word", "WIZARDRY", "要輸入的猜測")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	for total := uint64(0); total < 240_000_000; {
		n := uint64(2_000_000)
		if rest := 240_000_000 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}

	// 記錄防拷之後的印字事件
	seen := map[string]bool{}
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(o *oracle.Oracle) {
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		if !seen[s] {
			seen[s] = true
			fmt.Printf("印字： %q\n", s)
		}
	})

	// 主選單 → LOAD SAVEGAME → 防拷框
	for total := 240_000_000; total < 940_000_000; {
		n := uint64(2_000_000)
		if rest := 940_000_000 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}
	if err := o.Click(160, 75); err != nil {
		fmt.Println("點擊：", err)
	}
	for total := uint64(0); total < 40_000_000; {
		n := uint64(2_000_000)
		if rest := 40_000_000 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}
	if err := o.WritePNG("workplace/bake/cp-before.png"); err != nil {
		panic(err)
	}

	// 輸入猜測 + Enter（AH=07h 的 stdin；Enter = CR 0x0D）
	o.Type(*word)
	o.Type("\r")
	fmt.Println("輸入：", *word)
	for total := uint64(0); total < 60_000_000; {
		n := uint64(2_000_000)
		if rest := 60_000_000 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += n
	}
	if err := o.WritePNG("workplace/bake/cp-after.png"); err != nil {
		panic(err)
	}
	fmt.Println("輸出 cp-before.png / cp-after.png")
}
