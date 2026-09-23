// dirtest — 方向盤熱點的「點→動作」映射表。
//
// 用 Save/Restore 讓每個測試都從同一基準出發：Restore → 點 → 短跑 →
// 截圖 → 下一點。截圖連續看即知每個熱點的方向語意。
//
//	dirtest -exe <DS.EXE> -root <素材目錄>
package main

import (
	"flag"
	"fmt"

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
	// 建隊 + 旁白推進
	seq := []struct {
		x, y int
	}{{160, 62}, {170, 153}, {90, 168}, {90, 168}, {160, 101}}
	for _, s := range seq {
		if err := o.Click(s.x, s.y); err != nil {
			fmt.Printf("旁白點 (%d,%d) 無變化（容忍）\n", s.x, s.y)
		}
		_ = o.Run(120_000_000)
	}
	o.Type("AAAA\r")
	// 旁白 18 點
	for i := 0; i < 24; i++ {
		if err := o.Click(25, 15); err != nil {
			fmt.Printf("旁白點 %d 無變化（容忍）\n", i)
		}
		_ = o.Run(150_000_000)
	}
	_ = o.Run(100_000_000)
	fmt.Println("到達迷宮待機，存基準")

	base := o.Save()

	// 測試點（方向盤區＋視野）
	pts := []struct {
		name string
		x, y int
	}{
		{"up-arrow", 142, 165},
		{"down-arrow", 142, 191},
		{"left-arrow", 113, 178},
		{"right-arrow", 188, 178},
		{"dpad-center", 160, 178},
		{"viewport-center", 160, 100},
		{"viewport-left", 80, 100},
		{"viewport-right", 240, 100},
		{"corner-book", 70, 176},
		{"corner-map", 252, 176},
		{"corner-lens", 70, 196},
		{"corner-hand", 252, 196},
	}
	for i, p := range pts {
		o.Restore(base)
		if err := o.Click(p.x, p.y); err != nil {
			fmt.Printf("%-18s 點擊失敗：%v\n", p.name, err)
			continue
		}
		_ = o.Run(20_000_000)
		_ = o.MoveMouse(0, 0)
		_ = o.Run(2_000_000)
		png := fmt.Sprintf("workplace/bake/dir-%02d-%s.png", i, p.name)
		if err := o.WritePNG(png); err != nil {
			panic(err)
		}
		_ = o.Run(20_000_000)
		_ = o.MoveMouse(0, 0)
		_ = o.Run(2_000_000)
		png2 := fmt.Sprintf("workplace/bake/dir-%02d-%s-late.png", i, p.name)
		if err := o.WritePNG(png2); err != nil {
			panic(err)
		}
		fmt.Printf("%-18s (%d,%d) → dir-%02d-*.png\n", p.name, p.x, p.y, i)
	}
	
}
