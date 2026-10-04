package main

import "testing"

func reviewSignals(height int) []*textEv {
	glyphs := func(n int) [][]uint16 {
		g := make([][]uint16, n)
		for i := range g {
			g[i] = make([]uint16, height)
			g[i][0] = 1
		}
		return g
	}
	var shown []*textEv
	// R937實跑的選人訊號；不從正式標籤表生成期望或加入已消失的父選單。
	for _, label := range []struct {
		s    string
		x, y int
	}{
		{"REVIEW WHO?", 67, 152}, {"OK", 94, 168}, {"CANCEL", 82, 184},
	} {
		shown = append(shown, &textEv{s: label.s, x0: label.x, y0: label.y, g: glyphs(2)})
	}
	return shown
}

func TestReviewButtonsStayInsideNativeRectangles(t *testing.T) {
	for _, c := range []struct{ height, scale, okY, cancelY int }{
		{15, 2, 334, 366}, {16, 2, 333, 365}, {15, 3, 504, 552}, {16, 3, 504, 552},
	} {
		s := reviewSignals(c.height)
		if got := buttonLabelY(s[1], s, c.scale); got != c.okY {
			t.Fatalf("height%d scale%d OK: %d", c.height, c.scale, got)
		}
		if got := buttonLabelY(s[2], s, c.scale); got != c.cancelY {
			t.Fatalf("height%d scale%d CANCEL: %d", c.height, c.scale, got)
		}
		if got := buttonLabelY(s[0], s, c.scale); got != 152*c.scale {
			t.Fatal("標題位置須保持")
		}
	}
}

func TestReviewRequiresCompleteConsistentSignals(t *testing.T) {
	for _, name := range []string{"missing", "source", "coordinate", "dy", "height", "bits", "blank", "wide", "scale"} {
		t.Run(name, func(t *testing.T) {
			s := reviewSignals(15)
			e, scale := s[1], 2
			switch name {
			case "missing":
				s = s[1:]
			case "source":
				s[0].s = "ADD CHARACTER"
			case "coordinate":
				s[0].x0++
			case "dy":
				s[0].dy = 1
			case "height":
				s[0].g[0] = make([]uint16, 16)
			case "bits":
				s[0].g[0][0] = 1 << 15
			case "blank":
				s[0].g[0] = make([]uint16, 15)
			case "wide":
				e.g = append(e.g, e.g[0], e.g[0], e.g[0], e.g[0])
			case "scale":
				scale = 4
			}
			if got := buttonLabelY(e, s, scale); got != 168*scale {
				t.Fatalf("%s須保留舊畫筆，得到%d", name, got)
			}
		})
	}
}

func TestReviewDoesNotRequireUnrelatedShownEvents(t *testing.T) {
	s := reviewSignals(16)
	for i := 0; i < 6; i++ {
		s = append(s, &textEv{s: "UNRELATED", x0: 210, y0: 150 + 6*i})
	}
	if got := buttonLabelY(s[1], s, 2); got != 333 {
		t.Fatalf("其他事件不應阻止完整選人訊號置中，得到%d", got)
	}
}
