package main

import (
	"bytes"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

// 自製來源測試生命週期；原字型及正常玩家路徑另驗。
func profileFixture() (*profileReflowState, []byte) {
	start := bytes.Repeat([]byte{7}, 64000)
	final := append([]byte(nil), start...)
	p := &profileReflowState{begin: start, sources: map[string]*profileSource{}}
	g := make([]uint16, 15)
	g[0] = 1 << 14
	put := func(role string, x, y int, fg byte, glyph bool) {
		at := y*320 + x
		mask := []profileChange{{at, 7, fg}}
		final[at] = fg
		s := &profileSource{role: role, box: profileRect{x, y, 6, 6}, mask: mask, fg: fg}
		if glyph {
			s.g = [][]uint16{append([]uint16(nil), g...)}
		}
		p.sources[role] = s
	}
	put("sex", 90, 6, 1, true)
	for _, site := range profileLabels() {
		put(site.role, site.x, site.y, site.fg, true)
	}
	put("level-value", 126, 18, 132, false)
	put("experience-value", 306, 12, 71, false)
	put("kills-value", 306, 18, 71, false)
	for i := range profileAttributes {
		put("attribute-value-"+string(rune('0'+i)), 145, 35+6*i, 85, false)
	}
	health := &profileSource{role: "health-band", box: profileRect{6, 93, 78, 12}}
	health.mask = append(health.mask, p.sources["HP"].mask...)
	health.mask = append(health.mask, p.sources["ST"].mask...)
	for _, dot := range []struct {
		x, y int
		c    byte
	}{{24, 93, 27}, {24, 99, 29}} {
		at := dot.y*320 + dot.x
		health.mask = append(health.mask, profileChange{at, 7, dot.c})
		final[at] = dot.c
	}
	p.sources["health-band"] = health
	return p, final
}

func TestProfileReflowCompletePageAndIndependentView(t *testing.T) {
	p, final := profileFixture()
	page, ok := p.seal(final)
	if !ok || len(page.view.fields) != 31 || !page.matches(final, nil) {
		t.Fatal("完整來源未封存31欄")
	}
	var level, sex profileField
	for _, f := range page.view.fields {
		if f.role == "level-value" {
			level = f
		}
		if f.role == "sex" {
			sex = f
		}
	}
	if sex.x != 177 || sex.y != 9 || len(level.dots) != 4 || level.dots[0].x != 252 || level.dots[0].y != 59 {
		t.Fatal("置中或原數值水平位置改變", sex, level)
	}
	copy := page.visualCopy()
	for i := range copy.view.fields {
		if len(copy.view.fields[i].g) > 0 {
			copy.view.fields[i].g[0][0] = 0
		}
		if len(copy.view.fields[i].dots) > 0 {
			copy.view.fields[i].dots[0].x = -1
		}
	}
	if !profileGlyphsValid(sex.g, sex.cell) || level.dots[0].x != 252 {
		t.Fatal("GUI副本改到來源")
	}
	changed := append([]byte(nil), final...)
	changed[14*320+90]++ // 原型擴大字格的非來源背景。
	if page.matches(changed, nil) {
		t.Fatal("未檢查擴大字格背景")
	}
	if !page.matches(changed, func(x, y int) bool { return x == 90 && y == 14 }) {
		t.Fatal("游標容忍失效")
	}
	comp := append([]byte(nil), final...)
	page.restore(comp, nil)
	if !bytes.Equal(comp, p.begin) || !bytes.Equal(final, page.final) {
		t.Fatal("mask還原修改原版或漏筆畫")
	}
}

func TestProfileReflowRejectsIncompleteOrChangedSource(t *testing.T) {
	for _, mode := range []string{"missing", "background", "source", "overlap", "empty-glyph"} {
		t.Run(mode, func(t *testing.T) {
			p, final := profileFixture()
			switch mode {
			case "missing":
				delete(p.sources, "attribute-value-7")
			case "background":
				final[14*320+90]++
			case "source":
				final[p.sources["LVL"].mask[0].at]++
			case "overlap":
				p.sources["class"].mask = append(p.sources["class"].mask, p.sources["race"].mask[0])
			case "empty-glyph":
				p.sources["race"].g[0][0] = 0
			}
			if _, ok := p.seal(final); ok {
				t.Fatal("失配仍畫整頁", mode)
			}
		})
	}
}

func TestProfileReflowLongValueDropsOldPageWithoutTouchingOtherPending(t *testing.T) {
	p, final := profileFixture()
	page, ok := p.seal(final)
	if !ok {
		t.Fatal("fixture")
	}
	p.event = &textEv{profile: page}
	other := &pendingPrint{s: "unrelated"}
	ov := &overlay{profileEnabled: true, profile: p, pend: other}
	handled := ov.enterProfileReflowText(oracle.Addr{Seg: 0x0110, Off: 0x032a}, 400, 0, 120, 18,
		[]byte("11\x00"), nil, func() []byte { t.Fatal("拒絕長值不應讀離屏"); return nil })
	if handled || !p.failed || !p.event.dead || ov.pend != other {
		t.Fatal("長值殘留舊1或吞其他pending")
	}
	for _, tuple := range []struct {
		ret  uint16
		x, y int
	}{{0x032a, 30, 190}, {0x6cf8, 30, 190}} {
		if ov.enterProfileReflowText(oracle.Addr{Seg: 0x0110, Off: tuple.ret}, 400, 0, tuple.x, tuple.y,
			[]byte("other\x00"), nil, func() []byte { t.Fatal("未認領來源不應讀離屏"); return nil }) || ov.pend != other {
			t.Fatal("失敗後吞其他畫面的未認領入口")
		}
	}
	if !ov.enterProfileReflowText(oracle.Addr{Seg: 0x0110, Off: 0x032a}, 400, 0, 126, 18,
		[]byte("1\x00"), nil, func() []byte { t.Fatal("已失敗來源不應再讀離屏"); return nil }) {
		t.Fatal("已認領來源在失敗後進入舊覆繪")
	}
}

func TestProfileReflowReturnAndHealthFailureAreLocal(t *testing.T) {
	p, final := profileFixture()
	other := &pendingPrint{s: "unrelated"}
	ov := &overlay{profileEnabled: true, profile: p, pend: other}
	p.pending = &profilePending{role: "level-value", ret: oracle.Addr{Seg: 0x0110, Off: 0x032a},
		sp: 100, box: profileRect{126, 18, 6, 6}, endX: 132, endY: 18, fg: 132, before: p.begin}
	ov.returnedProfileReflow(p.pending.ret, 103, 132, 18, final)
	if !p.failed || p.pending != nil || ov.pend != other {
		t.Fatal("錯SP沒有失敗或清到別的slot")
	}
	p, final = profileFixture()
	beforeHealth := append([]byte(nil), final...)
	for _, c := range p.sources["health-band"].mask {
		beforeHealth[c.at] = c.before
	}
	p.health = append([]byte(nil), beforeHealth...)
	p.finishHealth(final)
	if p.failed || p.health != nil {
		t.Fatal("合法生命體力與標籤並存被拒")
	}
	p.health = append([]byte(nil), beforeHealth...)
	final[94*320+6] = 27 // 標籤側有區間以外的未知圖元。
	p.finishHealth(final)
	if !p.failed || p.health != nil {
		t.Fatal("未知圖元被當生命數值")
	}
}
