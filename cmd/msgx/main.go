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
				bx = -int(int16(v)) // 內部節點：neg 轉正
			} else {
				out = append(out, byte(v)) // 葉：字元
				break                      // byte 對齊（剩餘 bits 丟棄）
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
	comp := d[off : off+*length]
	out := decompress(comp[1:], tr, *maxOut) // 跳過第一 byte（長度）
	fmt.Printf("塊 %s 解壓 %d bytes：\n", *blk, len(out))
	// 可讀化
	vis := make([]byte, 0, len(out))
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
