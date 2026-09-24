// msgall — 全文抽取：msglist.json 的每則訊息，用對應樹＋壓縮塊解壓明文。
//
//	msgall -dbs <MSG.DBS> -list <msglist.json> -treedir <trees/>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func word(t []byte, off int) (uint16, bool) {
	if off < 0 || off+1 >= len(t) {
		return 0, false
	}
	return uint16(t[off]) | uint16(t[off+1])<<8, true
}

// decompress 走訪樹（位元流 MSB-first、葉回根、內部 neg 轉正 ×4）。
// 起點＝塊 pos+8（結構 header 之後）。
func decompress(dbs []byte, tree []byte, pos int64, maxOut int) []byte {
	start := int(pos) + 8
	if start+1 >= len(dbs) {
		return nil
	}
	out := []byte{}
	bx := 0
	al := 0
	cx := 0
	i := start
	// 葉後的「下一個 bit」跳過規則：bit 流連續、不額外跳
	_ = i
	for len(out) < maxOut {
		if cx == 0 {
			if i >= len(dbs) {
				break
			}
			al = int(dbs[i])
			cx = 8
			i++
		}
		bit := (al >> 7) & 1
		al = (al << 1) & 0xFF
		cx--
		var v uint16
		var ok bool
		if bit != 0 {
			v, ok = word(tree, bx+2)
		} else {
			v, ok = word(tree, bx)
		}
		if !ok {
			break
		}
		if v&0x8000 != 0 {
			bx = ((-int(v)) & 0xFFFF) * 4
		} else {
			out = append(out, byte(v))
			bx = 0
		}
	}
	return out
}

type msgRec struct {
	Step    uint64 `json:"step"`
	Text    string `json:"text"`
	TreeSeg uint16 `json:"tree_seg"`
	ChunkPo int64  `json:"chunk_pos"`
}

func main() {
	dbsF := flag.String("dbs", "workplace/original/wizar7/MSG.DBS", "")
	listF := flag.String("list", "workplace/bake/msglist.json", "")
	treeDir := flag.String("treedir", "workplace/bake/trees", "")
	maxOut := flag.Int("max", 400, "每則輸出上限")
	flag.Parse()

	dbs, err := os.ReadFile(*dbsF)
	if err != nil {
		panic(err)
	}
	b, err := os.ReadFile(*listF)
	if err != nil {
		panic(err)
	}
	var recs []msgRec
	if err := json.Unmarshal(b, &recs); err != nil {
		panic(err)
	}
	trees := map[uint16][]byte{}
	for _, r := range recs {
		if _, ok := trees[r.TreeSeg]; !ok {
			fn := fmt.Sprintf("%s/tree-%04x.bin", *treeDir, r.TreeSeg)
			if d, err := os.ReadFile(fn); err == nil {
				trees[r.TreeSeg] = d
			}
		}
	}
	seen := map[string]bool{}
	okN, failN := 0, 0
	for _, r := range recs {
		if seen[r.Text] {
			continue
		}
		seen[r.Text] = true
		tr, ok := trees[r.TreeSeg]
		if !ok {
			fmt.Printf("%q → 無樹 %04x\n", r.Text, r.TreeSeg)
			failN++
			continue
		}
		if r.ChunkPo < 0 {
			fmt.Printf("%q → 無塊位置\n", r.Text)
			failN++
			continue
		}
		// 兩層架構：塊 → 樹走訪#1 → 中間層（可能仍是壓縮流）→ 樹走訪#2 → 明文
		plain := decompress(dbs, tr, r.ChunkPo, *maxOut)
		layer2 := decompress(plain, tr, 0, *maxOut)
		if len(layer2) > 0 {
			plain = layer2
		}
		// 去除控制碼（<0x20 且非 NUL）便讀
		vis := strings.Map(func(rr rune) rune {
			if rr == 0 {
				return '␀'
			}
			if rr < 32 {
				return -1
			}
			return rr
		}, string(plain))
		fmt.Printf("ID %q\n  原文 %q\n  解壓 %q\n", r.Text, r.Text, vis)
		okN++
	}
	fmt.Printf("\n成功 %d、失敗 %d\n", okN, failN)
}
