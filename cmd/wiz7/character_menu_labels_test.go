package main

import "testing"

func characterMenuFixture(height int) []*textEv {
	var shown []*textEv
	for _, label := range []struct {
		s string
		x int
	}{{"CREATE", 11}, {"REVIEW", 63}, {"DELETE", 115}, {"RENAME", 167}, {"PICTURE", 216}, {"EXIT", 277}} {
		shown = append(shown, &textEv{s: label.s, x0: label.x, y0: 190,
			g: [][]uint16{make([]uint16, height), make([]uint16, height)}})
	}
	return shown
}

func TestCharacterMenuNativeLabelsFitMeasuredBounds(t *testing.T) {
	for _, tc := range []struct{ height, scale, want int }{{15, 2, 378}, {16, 2, 377}, {15, 3, 570}, {16, 3, 570}} {
		shown := characterMenuFixture(tc.height)
		for _, e := range shown {
			got := buttonLabelY(e, shown, tc.scale)
			if got != tc.want || got < 188*tc.scale || got+tc.height > 197*tc.scale {
				t.Fatalf("%s height=%d scale=%d y=%d", e.s, tc.height, tc.scale, got)
			}
			if e.y0 != 190 || e.dy != 0 || len(e.g[0]) != tc.height {
				t.Fatal("source event changed")
			}
		}
	}
}

func TestCharacterMenuRequiresCompleteValidFamily(t *testing.T) {
	for index := 0; index < 6; index++ {
		for _, tc := range []struct {
			name string
			edit func(*textEv)
		}{
			{"key", func(e *textEv) { e.s = "OTHER" }},
			{"x", func(e *textEv) { e.x0++ }},
			{"y", func(e *textEv) { e.y0++ }},
			{"offset", func(e *textEv) { e.dy++ }},
			{"empty", func(e *textEv) { e.g = nil }},
			{"unknown height", func(e *textEv) { e.g[0] = make([]uint16, 14) }},
			{"mixed height", func(e *textEv) { e.g[1] = make([]uint16, 16) }},
			{"outside bits", func(e *textEv) { e.g[0][0] = 0x8000 }},
			{"too wide", func(e *textEv) {
				for i := 0; i < 8; i++ {
					e.g = append(e.g, make([]uint16, 15))
				}
			}},
		} {
			shown := characterMenuFixture(15)
			target := *shown[5]
			tc.edit(shown[index])
			if y := buttonLabelY(&target, shown, 2); y != 380 {
				t.Fatalf("member %d %s moved target to %d", index, tc.name, y)
			}
		}
		shown := characterMenuFixture(16)
		target := shown[5]
		shown = append(shown[:index:index], shown[index+1:]...)
		shown = append(shown, nil)
		if y := buttonLabelY(target, shown, 2); y != 380 {
			t.Fatalf("missing %d moved target to %d", index, y)
		}
	}
	shown := characterMenuFixture(16)
	shown[0].g = [][]uint16{make([]uint16, 15), make([]uint16, 15)}
	if y := buttonLabelY(shown[5], shown, 2); y != 380 {
		t.Fatalf("mixed family heights moved target to %d", y)
	}
}

func TestCharacterMenuDoesNotMoveUnrelatedOrInvalidTarget(t *testing.T) {
	shown := characterMenuFixture(16)
	for _, e := range []*textEv{{s: "EXIT", x0: 277, y0: 189}, {s: "OTHER", x0: 277, y0: 190}, {s: "EXIT", x0: 278, y0: 190}} {
		e.g = shown[5].g
		if buttonLabelY(e, shown, 2) != menuLabelY(e, shown, 2) {
			t.Fatal("unrelated event moved")
		}
	}
	if y := buttonLabelY(shown[5], shown, 4); y != 760 {
		t.Fatal("unsupported scale moved")
	}
	target := *shown[5]
	target.g = [][]uint16{make([]uint16, 14)}
	if y := buttonLabelY(&target, shown, 2); y != 380 {
		t.Fatal("invalid target moved")
	}
	for _, e := range characterDialogFixture(16) {
		if y := buttonLabelY(e, characterDialogFixture(16), 2); y != buttonLabelY(e, append(shown, characterDialogFixture(16)...), 2) {
			t.Fatal("existing dialog layout changed")
		}
	}
}
