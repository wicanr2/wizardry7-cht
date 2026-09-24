// msgx — MSG.DBS 解壓器（二元樹走訪，演算法見 docs/re/013 追記十五）。
//
//	msgx -dbs <MSG.DBS> -tree <tree.bin> -blk 0x2400 -len 400
package main

import (
	"flag"
	"fmt"
	"os"
)

// decompress 走訪 4-byte 節點樹：bit=0 走左 word、bit=1 走右 word。
// bit15=1 的值＝內部節點（負指標，neg 轉正＝下一節點編號）；
// bit15=0 的值＝葉（ASCII 字元，輸出後回根）。
func decompress(block []byte, tree []byte, maxOut int) []byte {
	node := func(idx int) uint16 {
		p := idx * 4
		if p+1 >= len(tree) {
			return 0x8000 // 越界保護：視為內部
		}
		return uint16(tree[p]) | uint16(tree[p+1])<<8
	}
	right := func(idx int) uint16 {
		p := idx*4 + 2
		if p+1 >= len(tree) {
			return 0x8000
		}
		return uint16(tree[p]) | uint16(tree[p+1])<<8
	}
	// 第一 byte＝解壓長度（字元數）；之後每 byte 走一次樹（樹深 ≤8），
	// 葉＝bit15=0 的值（low byte＝字元），內部＝bit15=1 負指標（neg 轉正）。
	if len(block) == 0 {
		return nil
	}
	want := int(block[0])
	out := []byte{}
	for i := 1; i < len(block) && len(out) < want && len(out) < maxOut; i++ {
		b := block[i]
		bx := 0
		for bit := 7; bit >= 0; bit-- {
			var v uint16
			if (b>>bit)&1 != 0 {
				v = right(bx)
			} else {
				v = node(bx)
			}
			if v&0x8000 != 0 {
				bx = -int(int16(v)) // 內部節點：neg 轉正（node() 內部已 ×4）
			} else {
				out = append(out, byte(v)) // 葉：字元
				bx = 0                     // 葉後回根（下一字元從根走）
				break                      // byte 對齊
			}
		}
	}
	return out
}

func main() {
	dbs := flag.String("dbs", "workplace/original/wizar7/MSG.DBS", "")
	tree := flag.String("tree", "workplace/bake/tree.bin", "")
	blk := flag.String("blk", "0x2400", "塊 offset（hex）")
	length := flag.Int("len", 512, "壓縮資料長度")
	maxOut := flag.Int("max", 2048, "輸出上限")
	bitsF := flag.String("bits", "", "實證位元流檔（trace 重建；優先於 blk）")
	flag.Parse()

	d, err := os.ReadFile(*dbs)
	if err != nil {
		panic(err)
	}
	tr, err := os.ReadFile(*tree)
	if err != nil {
		panic(err)
	}
	var off int
	fmt.Sscanf(*blk, "%x", &off)
	var out []byte
	bitsUsed := 0
	if *bitsF != "" {
		// 實證位元流模式：trace 重建的 bits + 樹走訪
		bitsTxt, err := os.ReadFile(*bitsF)
		if err != nil {
			panic(err)
		}
		bx := 0
		for _, c := range string(bitsTxt) {
			if c != '0' && c != '1' {
				continue
			}
			bit := int(c - '0')
			_ = bit
			p := bx*4
			if bit != 0 {
				p += 2
			}
			if p+1 >= len(tr) {
				break
			}
			v := uint16(tr[p]) | uint16(tr[p+1])<<8
			if v&0x8000 != 0 {
				bx = (-int(int16(v))) * 4
			} else {
				out = append(out, byte(v))
				bx = 0
			}
			if len(out) >= *maxOut {
				break
			}
			if bitsUsed < 20 {
				fmt.Printf("  bit=%d p=%d v=%04x out=%d bx=%d\n", bit, p, v, len(out), bx)
			}
			bitsUsed++
		}
		fmt.Printf("bitsUsed=%d\n", bitsUsed)
		vis := make([]byte, 0, len(out))
		for _, b := range out {
			if b == 0 {
				vis = append(vis, []byte("␀")...)
			} else if b >= 32 && b < 127 {
				vis = append(vis, b)
			}
		}
		fmt.Printf("實證位元流解壓 %d 字元：\n%s\n", len(out), vis)
		return
	}
	comp := d[off : off+*length]
	out2 := decompress(comp, tr, *maxOut)
	fmt.Printf("塊 %s 解壓 %d bytes：\n", *blk, len(out2))
	// 可讀化
	vis := make([]byte, 0, len(out2))
	nul := 0
	for _, b := range out {
		if b == 0 {
			vis = append(vis, []byte("␀")...)
			nul++
		} else if b >= 32 && b < 127 {
			vis = append(vis, b)
		} else {
			vis = append(vis, []byte(fmt.Sprintf("<%02x>", b))...)
		}
	}
	os.Stdout.Write(vis)
	fmt.Println()
	fmt.Printf("NUL 數：%d\n", nul)
}
