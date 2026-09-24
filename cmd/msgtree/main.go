// msgtree — 動態樹捕捉：印字時記錄（明文、樹段、壓縮塊位置）並 dump 樹。
//
// 輸出：workplace/bake/trees/tree-XXXX.bin（各訊息群的樹）＋msglist.json
//
//	msgtree -exe <DS.EXE> -root <素材目錄>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

type msgRec struct {
	Step    uint64 `json:"step"`
	Text    string `json:"text"`
	TreeSeg uint16 `json:"tree_seg"`
	ChunkPo int64  `json:"chunk_pos"`
}

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	outDir := flag.String("out", "workplace/bake/trees", "")
	stop2 := flag.Uint64("stop2", 950_000_000, "跑到這個步數（主選單後）")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	os.MkdirAll(*outDir, 0o755)

	recs := []msgRec{}
	lastTree := uint16(0)
	dumped := map[uint16]bool{}

	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(o *oracle.Oracle) {
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		ts := o.Word(oracle.Addr{Seg: 0x0110, Off: 0x2D4E})
		// 最新 MSG.DBS 讀取位置
		pos := int64(-1)
		ops := o.FileOps()
		for i := len(ops) - 1; i >= 0; i-- {
			if strings.Contains(strings.ToUpper(ops[i].Name), "MSG.DBS") && ops[i].Op == "read" {
				pos = ops[i].Pos
				break
			}
		}
		recs = append(recs, msgRec{
			Step: o.Steps(), Text: s, TreeSeg: ts, ChunkPo: pos,
		})
		if ts != lastTree && !dumped[ts] {
			dumped[ts] = true
			lastTree = ts
			data := o.Bytes(oracle.Addr{Seg: ts, Off: 0}, 8192)
			fn := fmt.Sprintf("%s/tree-%04x.bin", *outDir, ts)
			if err := os.WriteFile(fn, data, 0o644); err != nil {
				panic(err)
			}
			fmt.Printf("樹段 %04x → %s\n", ts, fn)
		}
	})

	for total := uint64(0); total < *stop2; {
		n := uint64(2_000_000)
		if rest := *stop2 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += n
	}

	b, _ := json.MarshalIndent(recs, "", " ")
	os.WriteFile("workplace/bake/msglist.json", b, 0o644)
	fmt.Printf("訊息 %d 則、樹 %d 棵 → msglist.json + trees/\n", len(recs), len(dumped))
}
