package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

func popupFixture(t *testing.T, pub string) ([]byte, []*textEv) {
	t.Helper()
	root := os.Getenv("WIZ7_POPUP_FIXTURE_DIR")
	if root == "" {
		t.Skip("requires private dosgolem popup receipt; original assets are not committed")
	}
	dir := root + "/" + pub
	a, err := os.ReadFile(dir + "/indexed.bin")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(dir + "/shown.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Rows []struct {
			S             string
			X, Y, DY      int
			Box           [4]int
			FG            uint8
			G             [][]uint16
			Before, After string
			Compound      *compoundTrace
		}
	}
	if err = json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	evs := []*textEv{}
	for _, r := range data.Rows {
		before, _ := base64.StdEncoding.DecodeString(r.Before)
		after, _ := base64.StdEncoding.DecodeString(r.After)
		ret := oracle.Far(0x0110, 0xba40)
		if r.S == "CAST" || r.S == "CANCEL" {
			ret = oracle.Far(0x0110, 0x0ea2)
		}
		evs = append(evs, &textEv{s: r.S, x0: r.X, y0: r.Y, dy: r.DY, bx: r.Box[0], by: r.Box[1], bw: r.Box[2], bh: r.Box[3], fg: r.FG, g: r.G, before: before, after: after,
			shown: true, originRet: ret, originFont: 0, compound: r.Compound})
	}
	return a, evs
}
func TestPopupGuardObservedAndRejects(t *testing.T) {
	old := flowK
	flowK = 3
	defer func() { flowK = old }()
	a, controls := popupFixture(t, "1494")
	if !profilePopupActive(controls, a, nil) {
		t.Fatal("observed popup rejected")
	}
	before, _ := popupFixture(t, "1493")
	if profilePopupActive(controls, before, nil) {
		t.Fatal("inactive background accepted")
	}
	for _, p := range profilePopupBorder {
		changed := append([]byte(nil), a...)
		changed[p[1]*320+p[0]] ^= 1
		if profilePopupActive(controls, changed, nil) {
			t.Fatal("altered border accepted", p)
		}
	}
	if profilePopupActive(controls[:1], a, nil) {
		t.Fatal("missing control accepted")
	}
	if profilePopupActive(controls, a, func(int, int) bool { return true }) {
		t.Fatal("popup accepted without visible border witnesses")
	}
	for _, mutate := range []func(*textEv){func(e *textEv) { e.originRet.Off++ }, func(e *textEv) { e.originFont = 1 }, func(e *textEv) { e.fg = 84 }, func(e *textEv) { e.bw-- }, func(e *textEv) { e.y0-- }} {
		copyEv := *controls[0]
		mutate(&copyEv)
		if profilePopupActive([]*textEv{&copyEv, controls[1]}, a, nil) {
			t.Fatal("wrong control witness accepted")
		}
	}
	flowK = 4
	if profilePopupActive(controls, a, nil) {
		t.Fatal("unsupported scale accepted")
	}
}
func TestPopupEquipmentRejectsWrongMeaningAndGlyph(t *testing.T) {
	old := flowK
	flowK = 3
	defer func() { flowK = old }()
	_, items := popupFixture(t, "1493")
	for _, e := range items {
		if !profilePopupEquipment(e) {
			t.Fatal("observed equipment rejected", e.s)
		}
	}
	for _, mutate := range []func(*textEv){func(e *textEv) { e.s = "OTHER" }, func(e *textEv) { e.originRet.Off++ }, func(e *textEv) { e.originFont = 1 }, func(e *textEv) { e.fg = 104 }, func(e *textEv) { e.dy = 1 }, func(e *textEv) { e.shown = false }, func(e *textEv) { e.dead = true }, func(e *textEv) { e.bx = 49 }, func(e *textEv) { e.bw-- }, func(e *textEv) { e.g = nil }, func(e *textEv) { e.g = [][]uint16{make([]uint16, 15), make([]uint16, 15)} }} {
		copyEv := *items[0]
		mutate(&copyEv)
		if profilePopupEquipment(&copyEv) {
			t.Fatal("invalid equipment accepted")
		}
	}
	e := items[0]
	e.popupClip = profilePopupRect
	outside := func(x, y int) bool { return e.popupContains(x, y) }
	a, _ := popupFixture(t, "1494")
	if !e.rectMatches(a, outside) {
		t.Fatal("uncovered observed field rejected")
	}
	changed := append([]byte(nil), a...)
	changed[e.by*320+e.bx] ^= 1
	if e.rectMatches(changed, outside) {
		t.Fatal("uncovered change accepted")
	}
	if !e.popupContainsHD(49*3, 114*3, 3) || e.popupContainsHD(49*3-1, 114*3, 3) {
		t.Fatal("wrong scaled clip edge")
	}
	e.popupClip = image.Rectangle{}
	if e.popupContains(49, 114) {
		t.Fatal("cleared clip remains active")
	}
}

func TestPopupEquipmentRejectsCorruptGlyphVariants(t *testing.T) {
	old := flowK
	flowK = 2
	defer func() { flowK = old }()
	_, items := popupFixture(t, "1493")
	for _, mutate := range []func(*textEv){
		func(e *textEv) { e.g = [][]uint16{make([]uint16, 15)} },
		func(e *textEv) { e.g[0] = append(e.g[0], 1) },
		func(e *textEv) { e.g[0][0] = 0x8000 },
		func(e *textEv) { e.before = e.before[:len(e.before)-1] },
		func(e *textEv) { e.after = e.after[:len(e.after)-1] },
	} {
		e := *items[0]
		e.g = cloneVisualEvents([]*textEv{items[0]})[0].g
		mutate(&e)
		if profilePopupEquipment(&e) {
			t.Fatal("corrupt event accepted")
		}
	}
	for _, e := range items {
		if e.compound == nil {
			continue
		}
		c := *e.compound
		c.Inputs = e.compound.Inputs
		c.Inputs[1].Ret.Off++
		copyEv := *e
		copyEv.compound = &c
		if profilePopupEquipment(&copyEv) {
			t.Fatal("wrong second caller accepted")
		}
	}
}

func TestPopupVisualCloneOwnsClipAndGlyphs(t *testing.T) {
	source := &textEv{popupClip: image.Rect(49, 114, 271, 198), g: [][]uint16{{1, 2, 3}}}
	copyEv := cloneVisualEvents([]*textEv{source})[0]
	source.popupClip = image.Rectangle{}
	source.g[0][0] = 99
	if !copyEv.popupContains(49, 114) || copyEv.g[0][0] != 1 {
		t.Fatal("owner mutation changed published visual data")
	}
	copyEv.popupClip = image.Rectangle{}
	if copyEv.popupContainsHD(147, 342, 3) {
		t.Fatal("cleared clip retained")
	}
}

func TestPopupSpecialEventsAndCompoundInputs(t *testing.T) {
	old := flowK
	defer func() { flowK = old }()
	for _, scale := range []int{2, 3} {
		flowK = scale
		_, items := popupFixture(t, "1493")
		for _, size := range []int{15, 16} {
			for _, source := range items {
				e := *source
				e.g = make([][]uint16, len(source.g))
				for i := range e.g {
					e.g[i] = make([]uint16, size)
					e.g[i][0] = 1
				}
				if !profilePopupEquipment(&e) {
					t.Fatal("valid native glyph rejected", scale, size, e.s)
				}
				for _, change := range []func(*textEv){
					func(e *textEv) { e.profile = &profilePage{} },
					func(e *textEv) { e.page = &flowPage{} },
					func(e *textEv) { e.roster = &rosterWitness{} },
					func(e *textEv) { e.sex = true },
					func(e *textEv) { e.g = e.g[:len(e.g)-1] },
				} {
					bad := e
					change(&bad)
					if profilePopupEquipment(&bad) {
						t.Fatal("special or incomplete event accepted", e.s)
					}
				}
				if e.compound == nil {
					continue
				}
				for _, change := range []func(*compoundTrace){
					func(c *compoundTrace) { c.Inputs[0].Ret.Off++ },
					func(c *compoundTrace) { c.Inputs[1].Font++ },
					func(c *compoundTrace) { c.Inputs[1].X++ },
					func(c *compoundTrace) { c.Inputs[1].Y++ },
					func(c *compoundTrace) { c.Inputs[1].S = "OTHER" },
				} {
					c := *e.compound
					change(&c)
					bad := e
					bad.compound = &c
					if profilePopupEquipment(&bad) {
						t.Fatal("invalid pair accepted", e.s)
					}
				}
			}
		}
	}
}

type popupDrawCheck struct {
	done    bool
	failure string
}

func (g *popupDrawCheck) Update() error {
	if g.done {
		return ebiten.Termination
	}
	return nil
}
func (g *popupDrawCheck) Layout(int, int) (int, int) { return 64, 64 }
func (g *popupDrawCheck) Draw(*ebiten.Image) {
	if g.done {
		return
	}
	g.done = true
	for _, k := range []int{2, 3} {
		for _, size := range []int{15, 16} {
			for _, mode := range []string{"game", "footer", "nonzero", "width", "height"} {
				w, h := 320*k, 200*k
				if mode == "footer" {
					h += 40
				}
				if mode == "width" {
					w--
				}
				if mode == "height" {
					h--
				}
				backing := ebiten.NewImage(w+1, h+1)
				r := image.Rect(0, 0, w, h)
				if mode == "nonzero" {
					r = r.Add(image.Pt(1, 1))
				}
				dst := backing.SubImage(r).(*ebiten.Image)
				strip := ebiten.NewImage(4*(size+2)-2, size)
				strip.Fill(color.RGBA{255, 0, 0, 255})
				e := &textEv{x0: 18, y0: 128, popupClip: image.Rect(49, 114, 271, 198), g: [][]uint16{make([]uint16, size)}}
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(18*k), float64(128*k))
				e.drawPopupClipped(dst, strip, op, k)
				pixels := make([]byte, (w+1)*(h+1)*4)
				backing.ReadPixels(pixels)
				for y := 0; y < h+1; y++ {
					for x := 0; x < w+1; x++ {
						want := (mode == "game" || mode == "footer") && x >= 18*k && x < 49*k && x < 18*k+4*(size+2)-2 && y >= 128*k && y < 128*k+size
						v := pixels[(y*(w+1)+x)*4:][:4]
						if (want && (v[0] != 255 || v[3] != 255)) || (!want && (v[0] != 0 || v[3] != 0)) {
							g.failure = fmt.Sprintf("scale=%d size=%d mode=%s pixel=%d,%d", k, size, mode, x, y)
							return
						}
					}
				}
				strip.Deallocate()
				backing.Deallocate()
			}
		}
	}
}

func TestPopupDrawProtectsWindowFooterAndRejectsWrongCanvas(t *testing.T) {
	g := &popupDrawCheck{}
	if err := ebiten.RunGame(g); err != nil {
		t.Fatal(err)
	}
	if !g.done || g.failure != "" {
		t.Fatal("Draw constraints failed", g.failure)
	}
}
