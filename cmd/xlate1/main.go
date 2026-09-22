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
type entry struct {
	Dy int         `json:"dy"`
	G  [][]uint16  `json:"g"`
}

var trans = map[string]entry{}

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
	collect := flag.String("collect", "", "無譯文串寫到這個檔（一行一串）")
	dy := flag.Int("dy", 0, "中文 y 微調像素")
	clicks := flag.String("clicks", "", "step:x:y 逗號分隔（絕對步數時點擊）")
	keysAt := flag.String("keys-at", "", "step:key 逗號分隔（Esc=0x1b）")
	typeAt := flag.String("type-at", "", "step:字串（於指定步數輸入字串+CR）")
	shotsAt := flag.String("shots-at", "", "step 逗號分隔（純截圖）")
	bkeysAt := flag.String("bkeys-at", "", "step:鍵名 逗號分隔（BIOS 緩衝路徑，如 step:Up）")
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
	collected := map[string]bool{}

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
		ent, ok := trans[s]
		if !ok {
			if *collect != "" && !collected[s] {
				collected[s] = true
				f, _ := os.OpenFile(*collect, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
				f.WriteString(s + "\n")
				f.Close()
			}
			return // 無譯文：放行
		}
		// 位置：driver 座標變數；畫到離屏（blit 的來源段）
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})) + *dy + ent.Dy
		scr := lastScr
		fg := lastFG
		fmt.Printf("翻譯 %d %q → %d 字形 @(%d,%d) fg=%02x 離屏=%04x\n", o.Steps(), s, len(ent.G), x, y, fg, scr)

		cx := x
		for _, g := range ent.G {
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

	type clk struct{ at uint64; x, y int }
	var clicksQ []clk
	var keysQ []struct {
		at uint64
		b  byte
	}
	var shotsQ []uint64
	for _, k := range strings.Split(*shotsAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		if _, err := fmt.Sscanf(k, "%d", &at); err == nil {
			shotsQ = append(shotsQ, at)
		}
	}
	var bkeysQ []struct {
		at uint64
		n  string
	}
	for _, k := range strings.Split(*bkeysAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var n string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &n); err == nil {
			bkeysQ = append(bkeysQ, struct {
				at uint64
				n  string
			}{at, n})
		}
	}
	var typeQ []struct {
		at uint64
		s  string
	}
	for _, k := range strings.Split(*typeAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var s string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &s); err == nil {
			typeQ = append(typeQ, struct {
				at uint64
				s  string
			}{at, s})
		}
	}
	for _, k := range strings.Split(*keysAt, ",") {
		if k == "" {
			continue
		}
		var at uint64
		var hb string
		if _, err := fmt.Sscanf(k, "%d:%s", &at, &hb); err == nil {
			var v int
			fmt.Sscanf(hb, "%x", &v)
			keysQ = append(keysQ, struct {
				at uint64
				b  byte
			}{at, byte(v)})
		}
	}
	for _, c := range strings.Split(*clicks, ",") {
		if c == "" {
			continue
		}
		var at uint64
		var x, y int
		if _, err := fmt.Sscanf(c, "%d:%d:%d", &at, &x, &y); err == nil {
			clicksQ = append(clicksQ, clk{at, x, y})
		}
	}
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
		for _, tq := range typeQ {
			if tq.at <= total && tq.at > total-n {
				fmt.Printf("輸入 %q @%d\n", tq.s, tq.at)
				o.Type(tq.s + "\r")
			}
		}
		for _, sa := range shotsQ {
			if sa <= total && sa > total-n {
				_ = o.Run(1_000_000)
				_ = o.WritePNG(fmt.Sprintf("workplace/bake/shot-%d.png", sa))
			}
		}
		for _, bq := range bkeysQ {
			if bq.at <= total && bq.at > total-n {
				fmt.Printf("BIOS鍵 %q @%d\n", bq.n, bq.at)
				if bq.n == "w" || bq.n == "a" || bq.n == "s" || bq.n == "d" || bq.n == "x" {
					if err := o.TypeKeys(bq.n); err != nil {
						fmt.Println("TypeKeys：", err)
					}
				} else if o.SendKeys(bq.n) != nil {
					fmt.Println("鍵名失敗：", bq.n)
				}
			}
		}
		for _, kk := range keysQ {
			if kk.at <= total && kk.at > total-n {
				fmt.Printf("按鍵 0x%02x @%d\n", kk.b, kk.at)
				o.Type(string([]byte{kk.b}))
				_ = o.Run(3_000_000)
				_ = o.WritePNG(fmt.Sprintf("workplace/bake/key-%d.png", kk.at))
			}
		}
		for _, c := range clicksQ {
			if c.at <= total && c.at > total-n {
				fmt.Printf("點擊 (%d,%d) @%d\n", c.x, c.y, c.at)
				if err := o.Click(c.x, c.y); err != nil {
					fmt.Println("點擊失敗：", err)
				}
				// 推進一點讓畫面穩定，存 PNG
				_ = o.Run(3_000_000)
				png := fmt.Sprintf("workplace/bake/click-%d.png", c.at)
				if err := o.WritePNG(png); err != nil {
					fmt.Println("PNG：", err)
				}
			}
		}
	}
	fmt.Println("畫了", drawn, "串；見過的離屏段：", scrs)
	if lastScr != 0 { // 驗屍：離屏文字區（y146-170）快照
		_ = os.WriteFile("workplace/bake/offscr-end.bin",
			o.Bytes(oracle.Addr{Seg: lastScr, Off: 46720}, 2000), 0o644)
	}
	if err := o.WritePNG("workplace/bake/xlate1.png"); err != nil {
		panic(err)
	}
	_ = os.WriteFile("workplace/bake/xlate1.vram", o.Indexed(), 0o644)
	fmt.Println("輸出 workplace/bake/xlate1.png")
}
