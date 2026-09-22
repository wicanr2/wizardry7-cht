// xlate1 — 巫術7 M4 方案 B 雛形：攔印字串、譯文畫中文。
//
// 機制（docs/re/013 追記三方案 B）：
//   OnCall(0110:3895)（印字串入口，hook 在第一道指令前）：
//     1. 先 unpoke（寫回原 byte，讓本字串照常印，除非有譯文）
//     2. 讀字串（ds:off，NUL 終止）查譯文表
//     3. 有譯文：driver 座標變數（x=740C:0295、y=0x0297）取位置，
//        Eten 15×15 寫離屏段（blit 的來源，非筆畫保原像素），
//        再 poke 0xC3（ret）到 3895 讓原印失效
//
//	xlate1 -exe <DS.EXE> -root <素材目錄>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

// 譯文表：key = 字串內容（stack 字串無固定位址，內容即 key；rulebook 的
// 來源位置 key 適用 MSG.DBS 抽取文本，這裡是 UI 組合字串，標 DRAFT）。
// 值 = 每字 15 個 row（MSB first、15px 寬），由 tools/re/bake_font.py 產出。
var trans = map[string][][]uint16{}

func loadTrans(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &trans); err != nil {
		panic(err)
	}
}

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	stop1 := flag.Uint64("stop1", 240_000_000, "")
	stop2 := flag.Uint64("stop2", 940_000_000, "")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	loadTrans("workplace/bake/xlate.json")

	for total := uint64(0); total < *stop1; {
		n := uint64(2_000_000)
		if rest := *stop1 - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}

	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	baseline := o.Byte(entry) // 原 3895 第一 byte（push bp = 0x55）
	poked := false
	var lastFG uint8 = 0x22
	var lastScr uint16 = 0
	scrs := map[uint16]bool{}
	drawn := 0

	o.OnCall(entry, func(o *oracle.Oracle) {
		if poked { // 先解除，本字串（無論譯不譯）從乾淨狀態開始
			o.WriteU8(entry, baseline)
			poked = false
		}
		off := o.StackWord(1)
		ds := o.Regs().DS
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
		s := string(raw)
		if i := strings.IndexByte(s, 0); i >= 0 {
			s = s[:i]
		}
		glyphs, ok := trans[s]
		if !ok {
			return // 無譯文：放行
		}
		// 位置：driver 座標變數；畫到離屏（blit 的來源段）
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
		scr := lastScr
		fg := lastFG
		fmt.Printf("翻譯 %q → %d 字形 @(%d,%d) fg=%02x 離屏=%04x\n", s, len(glyphs), x, y, fg, scr)

		cx := x
		for _, g := range glyphs {
			for row := 0; row < 15 && row < len(g); row++ {
				w := g[row]
				if y+row > 199 {
					continue
				}
				offA := uint16(((y + row) * 320 + cx) % 65536)
				px := o.Bytes(oracle.Addr{Seg: scr, Off: offA}, 15)
				for bit := 0; bit < 15; bit++ {
					if (w>>(15-bit))&1 != 0 {
						px[bit] = fg
					}
				}
				o.WriteBytes(oracle.Addr{Seg: scr, Off: offA}, px)
				if drawn == 0 && row == 0 { // 首串首列：寫後即讀驗證
					back := o.Bytes(oracle.Addr{Seg: scr, Off: offA}, 15)
					fmt.Printf("  驗證 scr=%04x off=%04x 寫=% x 讀=% x\n", scr, offA, px, back)
				}
			}
			cx += 16
		}
		drawn++
		o.WriteU8(entry, 0xC3) // ret：原印失效
		poked = true
	})

	// 記住前景色（原印時的字元事件）與離屏段（blit 的來源段）
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x012F}, func(o *oracle.Oracle) {
		lastFG = uint8(o.Regs().DX & 0xFF)
	})
	o.OnCall(oracle.Addr{Seg: 0x740C, Off: 0x05F5}, func(o *oracle.Oracle) {
		if r := o.Regs(); r.DS > 0 && r.DS != 0xA000 {
			lastScr = r.DS // blit 來源段＝離屏緩衝
			scrs[r.DS] = true
		}
	})

	for total := *stop1; total < *stop2; {
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
	fmt.Println("畫了", drawn, "串；見過的離屏段：", scrs)
	if err := o.WritePNG("workplace/bake/xlate1.png"); err != nil {
		panic(err)
	}
	_ = os.WriteFile("workplace/bake/xlate1.vram", o.Indexed(), 0o644)
	fmt.Println("輸出 workplace/bake/xlate1.png")
}
