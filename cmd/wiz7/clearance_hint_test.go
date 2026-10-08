package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
)

// 合成字形與合成英文只驗契約，不保存手冊內容。
func clearanceFixture(height int) (map[string]entry, *helpData, []*textEv) {
	gs := make([][]uint16, 13)
	for i := range gs {
		gs[i] = make([]uint16, height)
		gs[i][0] = 1
	}
	h := &helpData{Glyphs: map[string]helpGlyph{}}
	for _, r := range clearanceLabel + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		w := height
		if r >= 'A' && r <= 'Z' {
			w = 8
		}
		rows := make([]int, height)
		rows[0] = 1
		h.Glyphs[string(r)] = helpGlyph{w, rows}
	}
	evs := []*textEv{
		{s: "SECURITY CLEARANCE", x0: 106, y0: 148, shown: true, match: true},
		{s: "PLEASE TURN TO PAGE 7 OF YOUR", x0: 70, y0: 158, shown: true, match: true},
		{s: "MANUAL AND ENTER THE FIRST WORD", x0: 64, y0: 164, shown: true, match: true},
		{s: "LOCATED IN PARAGRAPH 2  LINE 4", x0: 70, y0: 170, shown: true, match: true, g: gs, fgRGB: [3]uint8{255, 255, 255}},
	}
	return map[string]entry{clearanceLastKey: {G: gs}}, h, evs
}

func TestClearanceCapturedHelpAndRecordingArea(t *testing.T) {
	trans, h, evs := clearanceFixture(16)
	raw, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "help.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	r := &recordingResources{metadata: []ic.Resource{{Role: "help"}}}
	r.afterRead = func(role, path string) {
		if err := os.WriteFile(path, []byte("replaced after capture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	c := newClearanceHint(trans, r.loadHelp(path))
	if c.font == nil || r.metadata[0].SHA256 != ic.Hash(raw) {
		t.Fatal("same captured help lost")
	}
	c.selected = &clearanceSelection{prefix: [4]byte{7, 2, 4, 1}, answer: "TEST"}
	v := c.view(evs, 3)
	g := &game{scale: 3, screen: ebiten.NewImage(320, 200), live: &inputRecording{font: &helpData{}}, sh: &shared{rgba: make([]byte, 64000*4), clearance: v}}
	g.Draw(ebiten.NewImage(960, 640))
	if g.clearanceView != v || g.clearanceImage == nil {
		t.Fatal("recording footer disabled valid hint")
	}
	g.sh.clearance = nil
	g.Draw(ebiten.NewImage(960, 640))
	if g.clearanceView != nil || g.clearanceImage != nil {
		t.Fatal("withdrawal retained GUI cache")
	}
}

func TestClearanceActualCapturedResources(t *testing.T) {
	root := os.Getenv("WIZ7_CLEARANCE_RESOURCES")
	if root == "" {
		t.Skip("本機正式字模未提供")
	}
	for _, backend := range []string{"unifont", "eten"} {
		xlate := filepath.Join(root, "..", "r977", backend, "xlate.json")
		help := filepath.Join(root, backend, "help.json")
		bundle, err := loadRecordingBundle(xlate, "", help)
		if err != nil || len(bundle.resources.metadata) != 6 {
			t.Fatalf("same six roles: %v", err)
		}
		c := newClearanceHint(bundle.trans, bundle.help)
		if c.font == nil {
			t.Fatal("actual captured hint font rejected")
		}
		raw, err := os.ReadFile(help)
		if err != nil {
			t.Fatal(err)
		}
		m := bundle.resources.metadata[bundle.resources.index("help")]
		if !m.Loaded || m.SHA256 != ic.Hash(raw) {
			t.Fatal("captured help identity")
		}
	}
}

func TestClearanceRecordBoundsAndCopy(t *testing.T) {
	raw := append([]byte{7, 2, 4, 1}, []byte("TEST\x00")...)
	s := decodeClearanceRecord(raw)
	if s == nil || s.answer != "TEST" || s.prefix != [4]byte{7, 2, 4, 1} {
		t.Fatal("valid record")
	}
	raw[0] = 99
	raw[4] = 'X'
	if s.answer != "TEST" || s.prefix[0] != 7 {
		t.Fatal("record aliases original buffer")
	}
	for _, bad := range [][]byte{nil, {7, 2, 4, 1, 0}, {7, 2, 4, 1, 'a', 0}, {7, 2, 4, 1, 'A'}, {7, 2, 4, 1, 'A', ' ', 0}} {
		if decodeClearanceRecord(bad) != nil {
			t.Fatal("malformed record accepted")
		}
	}
}

func TestClearanceAcquisitionPair(t *testing.T) {
	for _, v := range []struct {
		sp uint16
		ok bool
	}{{0, false}, {0x111, false}, {0x112, true}, {64594, true}, {65530, true}, {65531, false}, {65535, false}} {
		if clearanceEntrySP(v.sp) != v.ok {
			t.Fatal("entry argument span")
		}
	}
	p := &clearancePending{id: 2500, ds: 0x1100, dst: 60000, ss: 0x1100, sp: 64000, step: 100}
	if !p.paired(p.ds, p.ss, p.sp-0x112, p.sp-2, p.id, p.dst, 200) {
		t.Fatal("valid pair")
	}
	for _, mutate := range []func(*clearancePending){
		func(p *clearancePending) { p.id++ }, func(p *clearancePending) { p.ds++ }, func(p *clearancePending) { p.ss++ },
		func(p *clearancePending) { p.sp-- }, func(p *clearancePending) { p.dst = 65500 }, func(p *clearancePending) { p.step = 201 },
	} {
		bad := *p
		mutate(&bad)
		if bad.paired(p.ds, p.ss, p.sp-0x112, p.sp-2, p.id, p.dst, 200) {
			t.Fatal("wrong acquisition accepted")
		}
	}
	if p.paired(p.ds, p.ss, p.sp-0x112, p.sp-2, p.id, p.dst, 1000101) {
		t.Fatal("timeout accepted")
	}
}

func TestClearanceHintGeometryAndWithdrawal(t *testing.T) {
	for _, height := range []int{15, 16} {
		for _, scale := range []int{2, 3} {
			trans, h, evs := clearanceFixture(height)
			c := newClearanceHint(trans, h)
			c.selected = &clearanceSelection{prefix: [4]byte{7, 2, 4, 1}, answer: "TEST"}
			v := c.view(evs, scale)
			if v == nil {
				t.Fatal("valid view missing")
			}
			wantWidth := 3*height + 4*8 + 6*2
			wantX := 70*scale + 13*height + 12*2 + 8
			if v.x != wantX || v.y != 170*scale || v.pixels.Bounds() != image.Rect(0, 0, wantWidth, height) {
				t.Fatal("geometry")
			}
			// 記錄控制列在遊戲區外，提示只在原題末行帶內。
			img := image.NewRGBA(image.Rect(0, 0, 320*scale, 200*scale+40))
			v.paint(img)
			for _, b := range img.Pix {
				if b != 0 {
					t.Fatal("wrong canvas accepted")
				}
			}
			v.paint(img.SubImage(image.Rect(0, 0, 320*scale, 200*scale)).(*image.RGBA))
			ink := 0
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					if img.RGBAAt(x, y).A != 0 && !image.Pt(x, y).In(image.Rect(wantX, 170*scale, wantX+wantWidth, 170*scale+height)) {
						t.Fatal("paint escaped")
					}
					if img.RGBAAt(x, y).A != 0 {
						ink++
					}
				}
			}
			if ink != 7 {
				t.Fatal("complete synthetic glyphs missing")
			}
			if c.view(evs, scale) != v {
				t.Fatal("immutable view not reused")
			}
			if c.view(evs[:3], scale) != nil || c.selected != nil {
				t.Fatal("stale hint retained")
			}
		}
	}
}

func TestClearanceHintRejectsIncompleteContext(t *testing.T) {
	for _, change := range []func([]*textEv){
		func(e []*textEv) { e[0].dead = true }, func(e []*textEv) { e[1].match = false }, func(e []*textEv) { e[2].dy = 1 },
		func(e []*textEv) { e[3].x0++ }, func(e []*textEv) { e[1].s = "PLEASE TURN TO PAGE 8 OF YOUR" },
		func(e []*textEv) { e[3].s = "LOCATED IN PARAGRAPH -2  LINE 4" },
	} {
		trans, h, evs := clearanceFixture(16)
		c := newClearanceHint(trans, h)
		c.selected = &clearanceSelection{prefix: [4]byte{7, 2, 4, 1}, answer: "TEST"}
		change(evs)
		if c.view(evs, 3) != nil {
			t.Fatal("bad context accepted")
		}
	}
	trans, h, evs := clearanceFixture(16)
	delete(h.Glyphs, "Q")
	if newClearanceHint(trans, h).font != nil {
		t.Fatal("partial alphabet accepted")
	}
	if newClearanceHint(nil, h).font != nil {
		t.Fatal("English mode accepted")
	}
	trans, h, evs = clearanceFixture(16)
	c := newClearanceHint(trans, h)
	c.selected = &clearanceSelection{prefix: [4]byte{7, 2, 4, 1}, answer: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
	if c.view(evs, 2) != nil {
		t.Fatal("overflow accepted")
	}
	c.selected.answer = "TEST"
	if c.view(evs, 4) != nil {
		t.Fatal("unsupported scale accepted")
	}
}
