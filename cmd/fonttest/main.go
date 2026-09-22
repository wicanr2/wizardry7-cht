// fonttest — 巫術7 M4 實驗閉環：資料替換版中文字形。
//
// 1. 跑到主選單（940M 步）
// 2. WriteBytes 把 fontseg.bin（烘製的 VBFONT2 重寫版）寫回字型段 0x7BB6
// 3. 把選單字串表（linear 0x20E80，六項各 20 bytes）換成 0x80+ 索引串
// 4. 再跑幾步讓遊戲重繪，dump VRAM+調色盤看結果
//
//	fonttest -exe <DS.EXE> -root <素材目錄> -seg <fontseg.bin>
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	exe := flag.String("exe", "", "")
	root := flag.String("root", "", "")
	seg := flag.String("seg", "workplace/bake/fontseg.bin", "")
	menuAt := flag.Uint64("menu-at", 940_000_000, "")
	extra := flag.Uint64("extra", 20_000_000, "")
	dumpStr := flag.Bool("dump-str", true, "")
	flag.Parse()

	o, err := oracle.LoadWith(*exe, *root, oracle.Options{Adlib: true})
	if err != nil {
		panic(err)
	}
	for total := uint64(0); total < *menuAt; {
		n := uint64(2_000_000)
		if rest := *menuAt - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			panic(err)
		}
		total += n
	}
	fmt.Println("到主選單，steps =", o.Steps())

	// 字串表位置動態變動——搜尋定位
	hits := o.Search([]byte("CHARACTER MENU"))
	fmt.Println("CHARACTER MENU 命中線性位址：", hits)
	if len(hits) == 0 {
		panic("找不到字串表")
	}
	base := hits[len(hits)-1]
	if *dumpStr {
		raw := o.Bytes(oracle.Addr{Seg: uint16(base >> 4), Off: uint16(base & 0xF)}, 128)
		fmt.Printf("字串表 %#x：\n", base)
		for i := 0; i+20 <= 128; i += 20 {
			fmt.Printf("  %#x: %q\n", base+uint32(i), raw[i:i+20])
		}
	}

	font, err := os.ReadFile(*seg)
	if err != nil {
		panic(err)
	}
	if len(font) != 7440 {
		panic(fmt.Sprintf("fontseg %d bytes != 7440", len(font)))
	}
	o.WriteBytes(oracle.Addr{Seg: 0x7BB6, Off: 0}, font)
	fmt.Println("字型段已替換")

	// 六項譯文索引串（bake.json 的 translations）
	menu := [][]byte{
		{0x80, 0x81, 0x82, 0x83},          // 角色選單
		{0x84, 0x85, 0x82, 0x83},          // 新增角色
		{0x86, 0x87, 0x88, 0x89},          // 載入紀錄
		{0x8A, 0x8B, 0x8C, 0x8D},          // 系統設定
		{0x8E, 0x8F, 0x90, 0x91},          // 匯入隊伍
		{0x91, 0x92, 0x93, 0x94},          // 結束遊戲
	}
	for i, s := range menu {
		a := oracle.Addr{Seg: uint16((base + uint32(i*20)) >> 4), Off: uint16((base + uint32(i*20)) & 0xF)}
		buf := make([]byte, 20)
		copy(buf, s)
		o.WriteBytes(a, buf)
	}
	fmt.Println("字串表已替換")

	// 方向鍵移動選擇條 → 反白項重印
	o.PressKey(oracle.KeyUp)
	for total := uint64(0); total < *extra; {
		n := uint64(2_000_000)
		if rest := *extra - total; rest < n {
			n = rest
		}
		if err := o.Run(n); err != nil {
			fmt.Println("停止：", err)
			break
		}
		total += n
		fmt.Printf("重繪等待 %d/%d\n", total, *extra)
	}
	if err := o.WritePNG("workplace/bake/fonttest.png"); err != nil {
		panic(err)
	}
	if err := os.WriteFile("workplace/bake/fonttest.vram", o.Indexed(), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("輸出 workplace/bake/fonttest.png")
}
