package main

import "testing"

func characterDialogFixture(height int) []*textEv {
	shown := menuFixture(height)
	for _, label := range []struct {
		s    string
		x, y int
		n    int
	}{{"ADD CHARACTER", 61, 152, 4}, {"OK", 94, 168, 2}, {"CANCEL", 82, 184, 2}} {
		e := &textEv{s: label.s, x0: label.x, y0: label.y}
		for i := 0; i < label.n; i++ {
			e.g = append(e.g, make([]uint16, height))
		}
		shown = append(shown, e)
	}
	return shown
}

func TestCharacterDialogButtonsFitOriginalBounds(t *testing.T) {
	for _, tc := range []struct{ height, scale, okY, cancelY int }{
		{15, 2, 334, 366}, {16, 2, 333, 365},
		{15, 3, 504, 552}, {16, 3, 504, 552},
	} {
		shown := characterDialogFixture(tc.height)
		for _, b := range []struct{ index, top, want int }{{7, 166, tc.okY}, {8, 182, tc.cancelY}} {
			e := shown[b.index]
			got := buttonLabelY(e, shown, tc.scale)
			if got != b.want || got < b.top*tc.scale || got+tc.height > (b.top+9)*tc.scale {
				t.Fatalf("height=%d scale=%d %s y=%d want=%d", tc.height, tc.scale, e.s, got, b.want)
			}
			if e.x0 != map[int]int{7: 94, 8: 82}[b.index] || e.y0 != b.top+2 || len(e.g[0]) != tc.height {
				t.Fatal("source event changed")
			}
		}
		for _, e := range shown[:7] {
			if buttonLabelY(e, shown, tc.scale) != menuLabelY(e, shown, tc.scale) {
				t.Fatalf("unrelated label moved: %s", e.s)
			}
		}
	}
}

func TestCharacterDialogRequiresEveryMeasuredSignal(t *testing.T) {
	for missing := 0; missing < 9; missing++ {
		shown := characterDialogFixture(16)
		e := shown[7]
		partial := append(append([]*textEv{}, shown[:missing]...), shown[missing+1:]...)
		partial = append(partial, nil)
		if got := buttonLabelY(e, partial, 2); got != 336 {
			t.Fatalf("missing signal %d moved OK to %d", missing, got)
		}
	}
}

func TestCharacterDialogRejectsOtherLayouts(t *testing.T) {
	for _, tc := range []struct {
		name        string
		edit        func(*textEv)
		scale, want int
	}{
		{"other key", func(e *textEv) { e.s = "YES" }, 2, 336},
		{"other x", func(e *textEv) { e.x0++ }, 2, 336},
		{"other y", func(e *textEv) { e.y0++ }, 2, 338},
		{"offset", func(e *textEv) { e.dy = 1 }, 2, 338},
		{"other scale", func(e *textEv) {}, 4, 672},
		{"empty", func(e *textEv) { e.g = nil }, 2, 336},
		{"other height", func(e *textEv) { e.g = [][]uint16{make([]uint16, 14)} }, 2, 336},
		{"mixed height", func(e *textEv) { e.g[1] = make([]uint16, 15) }, 2, 336},
		{"too wide", func(e *textEv) {
			e.g = make([][]uint16, 8)
			for i := range e.g {
				e.g[i] = make([]uint16, 16)
			}
		}, 2, 336},
	} {
		t.Run(tc.name, func(t *testing.T) {
			shown := characterDialogFixture(16)
			e := *shown[7]
			e.g = append([][]uint16{}, e.g...)
			tc.edit(&e)
			if got := buttonLabelY(&e, shown, tc.scale); got != tc.want {
				t.Fatalf("y=%d want=%d", got, tc.want)
			}
		})
	}
}
