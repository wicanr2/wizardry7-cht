package main

import "testing"

func menuFixture(height int) []*textEv {
	var events []*textEv
	// 原版272M收據的固定座標，獨立於production的辨識表。
	labels := []struct {
		s    string
		x, y int
	}{
		{"CHARACTER MENU", 118, 45}, {"ADD CHARACTER", 121, 58},
		{"LOAD SAVEGAME", 121, 71}, {"CONFIGURATION", 121, 84},
		{"IMPORT (BANE)", 121, 97}, {"QUIT GAME", 133, 110},
	}
	for _, label := range labels {
		e := &textEv{s: label.s, x0: label.x, y0: label.y}
		for i := 0; i < 5; i++ {
			e.g = append(e.g, make([]uint16, height))
		}
		events = append(events, e)
	}
	return events
}

func TestMenuLabelsStayInsideButtons(t *testing.T) {
	for _, tc := range []struct{ height, scale, importY int }{
		{15, 2, 192}, {16, 2, 191}, {15, 3, 291}, {16, 3, 291},
	} {
		shown := menuFixture(tc.height)
		if got := menuLabelY(shown[4], shown, tc.scale); got != tc.importY {
			t.Fatalf("height=%d scale=%d import y=%d want=%d", tc.height, tc.scale, got, tc.importY)
		}
		for _, e := range shown {
			y := menuLabelY(e, shown, tc.scale)
			if y < (e.y0-2)*tc.scale || y+tc.height > (e.y0+7)*tc.scale {
				t.Fatalf("label %s outside original button", e.s)
			}
		}
	}
}

func TestMenuLayoutRequiresCompleteMeasuredContext(t *testing.T) {
	shown := menuFixture(16)
	if got := menuLabelY(shown[4], shown[1:], 2); got != 194 {
		t.Fatalf("partial menu moved import to %d", got)
	}
	other := *shown[3]
	other.x0 = 80
	if got := menuLabelY(&other, shown, 2); got != 168 {
		t.Fatalf("same label in other context moved to %d", got)
	}
	other = *shown[4]
	other.dy = 1
	if got := menuLabelY(&other, shown, 2); got != 196 {
		t.Fatalf("explicit offset replaced: %d", got)
	}
	if got := menuLabelY(shown[4], shown, 4); got != 388 {
		t.Fatalf("unmeasured scale moved: %d", got)
	}
}

func TestMenuLayoutRejectsGlyphsThatCannotFit(t *testing.T) {
	shown := menuFixture(16)
	shown[4].g[1] = make([]uint16, 15)
	if got := menuLabelY(shown[4], shown, 2); got != 194 {
		t.Fatalf("mixed glyph height accepted: %d", got)
	}
	shown = menuFixture(16)
	shown[4].g = make([][]uint16, 20)
	for i := range shown[4].g {
		shown[4].g[i] = make([]uint16, 16)
	}
	if got := menuLabelY(shown[4], shown, 2); got != 194 {
		t.Fatalf("overwide label accepted: %d", got)
	}
}
