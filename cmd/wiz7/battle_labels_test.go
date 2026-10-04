package main

import (
	"reflect"
	"testing"
)

func battleLabelFixture(height int) []*textEv {
	var events []*textEv
	// 獨立寫出R959已量座標，不從正式標籤表生成期望。
	for _, label := range []struct {
		s    string
		x, y int
	}{
		{"FIGHT", 13, 167}, {"PARRY", 13, 177}, {"SPELL", 55, 167},
		{"USE", 61, 177}, {"RUN", 241, 167}, {"MOVE", 238, 177}, {"EQUIP", 235, 187},
	} {
		e := &textEv{s: label.s, x0: label.x, y0: label.y, fg: 1}
		for i := 0; i < 2; i++ {
			g := make([]uint16, height)
			g[0], g[height-1] = 1, 1
			e.g = append(e.g, g)
		}
		events = append(events, e)
	}
	return events
}

func TestBattleNativeLabelsStayInsideOriginalButtons(t *testing.T) {
	for _, c := range []struct{ height, scale, firstY int }{
		{15, 2, 332}, {16, 2, 331}, {15, 3, 501}, {16, 3, 501},
	} {
		s := battleLabelFixture(c.height)
		for i, e := range s {
			before := *e
			got := buttonLabelY(e, s, c.scale)
			top := []int{165, 175, 165, 175, 165, 175, 185}[i]
			want := c.firstY + (top-165)*c.scale
			if got != want || got < top*c.scale || got+c.height > (top+9)*c.scale {
				t.Fatalf("%s height%d scale%d: y%d want%d", e.s, c.height, c.scale, got, want)
			}
			if !reflect.DeepEqual(before, *e) {
				t.Fatal("原事件不得修改")
			}
		}
	}
}

func TestBattleCenteringRejectsOtherScreensAndInvalidGlyphs(t *testing.T) {
	for _, name := range []string{"missing", "key", "x", "y", "dy", "color", "nil", "blank", "height", "mixed", "bits", "wide", "scale", "unrelated"} {
		t.Run(name, func(t *testing.T) {
			s := battleLabelFixture(15)
			e, scale := s[0], 2
			switch name {
			case "missing":
				s = s[:6]
			case "key":
				s[6].s = "QUIT GAME"
			case "x":
				s[6].x0++
			case "y":
				s[6].y0++
			case "dy":
				s[6].dy = 1
			case "color":
				s[6].fg = 2
			case "nil":
				s[6].g = nil
			case "blank":
				s[6].g[0] = make([]uint16, 15)
			case "height":
				s[6].g[0] = make([]uint16, 14)
			case "mixed":
				s[6].g[0] = make([]uint16, 16)
			case "bits":
				s[6].g[0][0] = 1 << 15
			case "wide":
				s[6].g = append(s[6].g, s[6].g[0], s[6].g[0])
			case "scale":
				scale = 4
			case "unrelated":
				e = &textEv{s: "FIGHT", x0: 14, y0: 167}
			}
			if got := buttonLabelY(e, s, scale); got != (e.y0+e.dy)*scale {
				t.Fatalf("%s誤調畫筆到%d", name, got)
			}
		})
	}
	if _, ok := battleControlLabelY(nil, nil, 2); ok {
		t.Fatal("nil不得接受")
	}
}

func TestBattleVisualCloneKeepsCenteringAndOwnsGlyphs(t *testing.T) {
	for _, height := range []int{15, 16} {
		for _, scale := range []int{2, 3} {
			source := battleLabelFixture(height)
			visual := cloneVisualEvents(source)
			for i, e := range visual {
				if e.fg != source[i].fg || buttonLabelY(e, visual, scale) != buttonLabelY(source[i], source, scale) {
					t.Fatalf("%s height%d scale%d: GUI副本失去色號或置中", e.s, height, scale)
				}
				e.g[0][0] = 2
				if source[i].g[0][0] != 1 {
					t.Fatal("GUI副本不得共享來源字模")
				}
			}
		}
	}
	for _, fg := range []uint8{0, 2, 255} {
		source := battleLabelFixture(16)
		source[6].fg = fg
		source[6].fgRGB = [3]uint8{12, 34, 56}
		visual := cloneVisualEvents(source)
		if visual[6].fg != fg || visual[6].fgRGB != source[6].fgRGB {
			t.Fatal("GUI副本必須保留來源色號與色值")
		}
		for _, e := range visual {
			if got := buttonLabelY(e, visual, 2); got != e.y0*2 {
				t.Fatal("異色GUI副本不得繞過完整家族守門")
			}
		}
	}
}
