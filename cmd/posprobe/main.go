// posprobe — #36 證據探針：隊伍座標與面向在記憶體的位置。
//
// 沿 cmd/wiz7 -autoplay 的路線進迷宮（同一組點擊與輸入），再做已知動作並在每步後
// 拍原版記憶體（1 MB）。只讀，不寫原版。篩選：
//
//	面向：右轉每次變化量相同、右轉四次回原值、左轉變化相反、前進／後退不變
//	座標：前進 ±1、後退回原值、轉向不變
//
//	posprobe -exe <DS.EXE> -root <素材目錄> [-ready 1470000000]
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/wicanr2/dosgolem/oracle"
)

type act struct {
	name string
	x, y int
}

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	ready := flag.Uint64("ready", 1_470_000_000, "進迷宮後開始動作的步數")
	settle := flag.Uint64("settle", 3_000_000, "每個動作後多跑幾步")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}

	// cmd/wiz7 autoplay 的前綴（START NEW GAME 路線）＋迷宮前兩下底文點擊
	type ev struct {
		at   uint64
		x, y int
		typ  string
	}
	evs := []ev{
		{250_000_000, 160, 62, ""}, {253_000_000, 170, 153, ""}, {255_000_000, 90, 168, ""},
		{257_000_000, 90, 168, ""}, {265_000_000, 160, 101, ""}, {276_000_000, 82, 169, ""},
		{288_000_000, 0, 0, "AAAA\r"},
	}
	for i := 0; i < 133; i++ {
		evs = append(evs, ev{295_000_000 + uint64(i)*8_000_000, 25, 15, ""})
	}
	evs = append(evs, ev{1_430_000_000, 320, 170, ""}, ev{1_434_000_000, 320, 150, ""})

	ei := 0
	for o.Steps() < *ready {
		for ei < len(evs) && o.Steps() >= evs[ei].at {
			if evs[ei].typ != "" {
				o.Type(evs[ei].typ)
			} else {
				_ = o.Click(evs[ei].x, evs[ei].y)
			}
			ei++
		}
		n := uint64(2_000_000)
		if rest := *ready - o.Steps(); rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
	}

	shot := 0
	mem := func() []byte { return o.Bytes(oracle.Addr{Seg: 0, Off: 0}, 0x100000) }
	do := func(a act) []byte {
		before := append([]uint8(nil), o.Indexed()...)
		err := o.Click(a.x, a.y)
		if err2 := o.Run(*settle); err2 != nil {
			panic(err2)
		}
		after := o.Indexed()
		px := 0
		for i := range after {
			if after[i] != before[i] {
				px++
			}
		}
		fmt.Printf("動作 %-4s (%d,%d) err=%v steps=%d A000 差異=%d\n", a.name, a.x, a.y, err, o.Steps(), px)
		shot++
		savePNG(o, fmt.Sprintf("workplace/bake/r666-pos-%d.png", shot))
		return mem()
	}
	R := act{"右轉", 180, 166}
	L := act{"左轉", 145, 166}
	F := act{"前進", 160, 162}
	K := act{"後退", 160, 184}

	savePNG(o, "workplace/bake/r666-pos-0.png")
	m0 := mem()
	// 右轉×4 回初始面向，再前進、後退（T4 證實初始面向前進會移動），最後左轉
	seq := []act{R, R, R, R, F, K, L}
	ms := [][]byte{m0}
	for _, a := range seq {
		ms = append(ms, do(a))
	}
	// ms[0] 起點；1..4 右轉；5 前進；6 後退；7 左轉
	d := func(a, b byte) int { return int(int8(b - a)) }

	fmt.Println("== 面向候選（右轉四值相異、四次回原值、前進後退不變、左轉＝右轉三次）")
	for i := range m0 {
		// 四個面向值兩兩不同（不假設編碼與繞回方式）
		v := [4]byte{ms[0][i], ms[1][i], ms[2][i], ms[3][i]}
		distinct := v[0] != v[1] && v[0] != v[2] && v[0] != v[3] && v[1] != v[2] && v[1] != v[3] && v[2] != v[3]
		ok := distinct && ms[4][i] == ms[0][i] &&
			ms[5][i] == ms[4][i] && ms[6][i] == ms[5][i] && ms[7][i] == ms[3][i]
		if ok {
			fmt.Printf("  %05x 值 %v\n", i, []byte{ms[0][i], ms[1][i], ms[2][i], ms[3][i], ms[4][i], ms[5][i], ms[6][i], ms[7][i]})
		}
	}
	fmt.Println("== 面向寬鬆候選（右轉四值相異、四次回原值、前進後退不變；不看左轉）")
	for i := range m0 {
		v := [4]byte{ms[0][i], ms[1][i], ms[2][i], ms[3][i]}
		distinct := v[0] != v[1] && v[0] != v[2] && v[0] != v[3] && v[1] != v[2] && v[1] != v[3] && v[2] != v[3]
		if distinct && ms[4][i] == ms[0][i] && ms[5][i] == ms[4][i] && ms[6][i] == ms[5][i] {
			fmt.Printf("  %05x 值 %v\n", i, []byte{ms[0][i], ms[1][i], ms[2][i], ms[3][i], ms[4][i], ms[5][i], ms[6][i], ms[7][i]})
		}
	}
	fmt.Println("== 座標附近 20e80–20ed0 各步值")
	for i := 0x20e80; i < 0x20ed0; i++ {
		row := make([]byte, len(ms))
		chg := false
		for k := range ms {
			row[k] = ms[k][i]
			if row[k] != ms[0][i] {
				chg = true
			}
		}
		if chg {
			fmt.Printf("  %05x %v\n", i, row)
		}
	}
	fmt.Println("== 座標候選（前進 ±1、後退回原值、轉向不變）")
	for i := range m0 {
		f := d(ms[4][i], ms[5][i])
		if f != 1 && f != -1 {
			continue
		}
		same := true
		for k := 1; k <= 4; k++ {
			if ms[k][i] != ms[0][i] {
				same = false
			}
		}
		if same && ms[6][i] == ms[4][i] && ms[7][i] == ms[6][i] {
			fmt.Printf("  %05x 值 %v\n", i, []byte{ms[0][i], ms[4][i], ms[5][i], ms[6][i], ms[7][i]})
		}
	}
}

func savePNG(o *oracle.Oracle, path string) {
	pal, idx := o.Palette(), o.Indexed()
	img := image.NewRGBA(image.Rect(0, 0, 320, 200))
	for i, v := range idx {
		c := pal[v]
		img.Set(i%320, i/320, color.RGBA{c[0], c[1], c[2], 255})
	}
	if f, err := os.Create(path); err == nil {
		png.Encode(f, img)
		f.Close()
	}
}
