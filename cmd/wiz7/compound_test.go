package main

import (
	"bytes"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

// 合成素材只用於防禦契約；正常遊戲與原始字形另外由固定收據驗收。
func compoundFixture(site compoundSite) (compoundInput, []byte, map[string]entry) {
	in := compoundInput{S: site.first, DS: 0x1106, Off: 0xfe44, SP: 0xfe38,
		Ret: compoundFirstReturn, Step: 100, X: 18, Y: site.y, raw: append([]byte(site.first), 0)}
	screen := make([]byte, 320*200)
	for i := range screen {
		screen[i] = byte(115 + i%5)
	}
	g := make([][]uint16, site.glyphs)
	for i := range g {
		g[i] = make([]uint16, 15)
		g[i][1] = 1
	}
	return in, screen, map[string]entry{site.first + "/" + site.second: {G: g}}
}

func compoundSecond(in compoundInput, site compoundSite) compoundInput {
	in.S, in.Ret, in.Step, in.Y = site.second, compoundSecondReturn, 236, site.y+6
	in.Off += uint16(len(site.first) + 1)
	in.raw = append([]byte(site.second), 0)
	return in
}

func TestCompoundOwnsBackgroundAndWaitsForBothReturns(t *testing.T) {
	for _, site := range compoundSites {
		t.Run(site.first, func(t *testing.T) {
			in, screen, trans := compoundFixture(site)
			original := append([]byte(nil), screen...)
			ov := &overlay{}
			if !ov.enterCompound(in, func() []byte { return screen }, trans, 3) || ov.compound == nil {
				t.Fatal("第一部件沒有建立待配對")
			}
			screen[site.y*320+18] = 105
			if !ov.returnCompound(in.Ret, in.SP+2, 200, screen, [256][3]uint8{}) || len(ov.evs) != 0 {
				t.Fatal("第一返回提前建立中文事件")
			}
			two := compoundSecond(in, site)
			if !ov.enterCompound(two, func() []byte { return screen }, trans, 3) || len(ov.evs) != 0 {
				t.Fatal("第二入口提前建立中文事件")
			}
			screen[(site.y+6)*320+20] = 105
			var pal [256][3]uint8
			pal[105] = [3]uint8{12, 34, 56}
			if !ov.returnCompound(two.Ret, two.SP+2, 300, screen, pal) || len(ov.evs) != 1 || ov.compound != nil {
				t.Fatal("兩返回沒有原子建立一件")
			}
			e := ov.evs[0]
			if e.s != site.first+"/"+site.second || e.bw != 41 || e.bh != 11 || e.fgRGB != pal[105] ||
				!bytes.Equal(e.before, compoundRect(original, site.y)) || e.compound == nil {
				t.Fatal("完整名稱、背景域或返回色盤不符")
			}
			before := append([]byte(nil), e.after...)
			screen[site.y*320+18] = 42
			trans[e.s].G[0][1] = 42
			if !bytes.Equal(e.after, before) || e.g[0][1] != 1 {
				t.Fatal("事件保留了來源緩衝參照")
			}
		})
	}
}

func TestCompoundRejectsInterruptedAndMalformedPairs(t *testing.T) {
	site := compoundSites[2]
	changes := map[string]func(*compoundInput, []byte){
		"錯誤SP":   func(p *compoundInput, _ []byte) { p.SP++ },
		"錯誤DS":   func(p *compoundInput, _ []byte) { p.DS++ },
		"指標不相鄰":  func(p *compoundInput, _ []byte) { p.Off++ },
		"錯誤座標":   func(p *compoundInput, _ []byte) { p.Y++ },
		"錯誤字型":   func(p *compoundInput, _ []byte) { p.Font++ },
		"錯誤時距":   func(p *compoundInput, _ []byte) { p.Step++ },
		"缺NUL":   func(p *compoundInput, _ []byte) { p.raw = []byte(site.second) },
		"NUL跨邊界": func(p *compoundInput, _ []byte) { p.Off = 65535 },
		"背景已改":   func(_ *compoundInput, screen []byte) { screen[site.y*320+58] = 1 },
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			in, screen, trans := compoundFixture(site)
			ov := &overlay{}
			ov.enterCompound(in, func() []byte { return screen }, trans, 3)
			screen[site.y*320+18] = 105
			ov.returnCompound(in.Ret, in.SP+2, 200, screen, [256][3]uint8{})
			two := compoundSecond(in, site)
			change(&two, screen)
			if !ov.enterCompound(two, func() []byte { return screen }, trans, 3) || ov.compound != nil || len(ov.evs) != 0 {
				t.Fatal("錯誤第二部件沒有整件回原文／抑制generic")
			}
		})
	}
	in, screen, trans := compoundFixture(site)
	ov := &overlay{}
	ov.enterCompound(in, func() []byte { return screen }, trans, 3)
	if ov.enterCompound(compoundInput{S: "UNTRANSLATED", Ret: oracle.Addr{Seg: 0x0110, Off: 0x1234}},
		func() []byte { t.Fatal("無關入口不應抓背景"); return nil }, trans, 3) || ov.compound != nil {
		t.Fatal("早退前未仲裁無譯文交錯入口")
	}
}

func TestCompoundPartSuppressionIsBounded(t *testing.T) {
	in, _, trans := compoundFixture(compoundSites[2])
	in = compoundSecond(in, compoundSites[2])
	ov := &overlay{}
	if !ov.enterCompound(in, func() []byte { t.Fatal("孤立第二部件不應讀背景"); return nil }, trans, 3) {
		t.Fatal("孤立SHIELD沒有抑制generic")
	}
	in.Ret = oracle.Addr{Seg: 0x0110, Off: 0x1234}
	if ov.enterCompound(in, func() []byte { return nil }, trans, 3) {
		t.Fatal("錯誤抑制另一畫面SHIELD")
	}
	in.Ret, in.S = compoundSecondReturn, "OTHER ITEM"
	if ov.enterCompound(in, func() []byte { return nil }, trans, 3) {
		t.Fatal("錯誤抑制未知物品")
	}
}

func TestCompoundMasksProtectFullDomain(t *testing.T) {
	base := bytes.Repeat([]byte{119}, 451)
	for _, point := range []int{0, 40, 41 * 4, 41*4 + 40} {
		a := append([]byte(nil), base...)
		a[point] = 105
		if compoundMask(base, a, 0) == nil {
			t.Fatal("有效第一行差異被拒")
		}
	}
	for _, point := range []int{41 * 5, 41*5 + 40, 41 * 6, 450} {
		a := append([]byte(nil), base...)
		a[point] = 105
		if compoundMask(base, a, 0) != nil {
			t.Fatal("間隙／另一行變動沒有整件拒絕")
		}
	}
	a := append([]byte(nil), base...)
	a[0] = 104
	if compoundMask(base, a, 0) != nil || compoundMask(base, base, 0) != nil {
		t.Fatal("非105或無差異被接受")
	}
}

func TestCompoundVisibilityChecksGapAndUnusedFirstLineWidth(t *testing.T) {
	site := compoundSites[1]
	_, screen, _ := compoundFixture(site)
	e := &textEv{bx: 18, by: site.y, bw: 41, bh: 11, after: compoundRect(screen, site.y)}
	if !e.rectMatches(screen, nil) {
		t.Fatal("完整域本來相同")
	}
	// FUR第一行原筆畫只寬17；右側與行間隙也在中文批准域內。
	for _, delta := range [][2]int{{40, 0}, {20, 5}, {40, 10}, {0, 0}} {
		a := append([]byte(nil), screen...)
		a[(site.y+delta[1])*320+18+delta[0]] = 42
		if e.rectMatches(a, nil) {
			t.Fatal("域內非筆畫或單行清場沒有整件失配")
		}
	}
	a := append([]byte(nil), screen...)
	a[(site.y-1)*320+18] = 42
	if !e.rectMatches(a, nil) {
		t.Fatal("域外無關改寫誤清")
	}
	a[site.y*320+18] = 42
	if !e.rectMatches(a, func(x, y int) bool { return x == 18 && y == site.y }) || e.rectMatches(a, nil) {
		t.Fatal("游標容忍和離屏完整比對混用")
	}
}

func TestCompoundRejectsInvalidGlyphsAndFirstReturn(t *testing.T) {
	site := compoundSites[0]
	in, screen, trans := compoundFixture(site)
	for _, scale := range []int{1, 4} {
		ov := &overlay{}
		if !ov.enterCompound(in, func() []byte { return screen }, trans, scale) || ov.compound != nil {
			t.Fatal("未核倍率沒有回原文")
		}
	}
	entry := trans[site.first+"/"+site.second]
	entry.G[0][1] = 0x8000
	if _, ok := compoundGlyphs(entry, site.glyphs, 3); ok {
		t.Fatal("15位字模溢位被接受")
	}
	entry.G[0][1] = 1
	ov := &overlay{}
	ov.enterCompound(in, func() []byte { return screen }, trans, 3)
	screen[site.y*320+18] = 105
	if !ov.returnCompound(in.Ret, in.SP+4, 200, screen, [256][3]uint8{}) || ov.compound != nil || len(ov.evs) != 0 {
		t.Fatal("錯誤返回SP沒有整件拒絕")
	}
}
