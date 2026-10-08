package main

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type whoFake struct {
	mode         uint8
	words        map[uint16]uint16
	tail, prefix []byte
	reads        int
}

func (f *whoFake) VideoMode() uint8          { return f.mode }
func (f *whoFake) Word(a oracle.Addr) uint16 { return f.words[a.Off] }
func (f *whoFake) Bytes(a oracle.Addr, n int) []byte {
	f.reads++
	if n == 144 {
		return append([]byte(nil), f.tail...)
	}
	if n == 48 {
		return append([]byte(nil), f.prefix...)
	}
	return nil
}
func (f *whoFake) Palette() [256][3]uint8 {
	p := [256][3]uint8{}
	p[22] = [3]uint8{200, 24, 0}
	return p
}
func whoFixture(t *testing.T, backend string) (*whoFake, entry, []byte) {
	t.Helper()
	root := os.Getenv("WIZ7_WHO_FIXTURE_DIR")
	if root == "" {
		t.Skip("requires private WHO cursor receipt and legal native glyph resource")
	}
	read := func(n string) []byte {
		b, e := os.ReadFile(filepath.Join(root, n))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	f := &whoFake{mode: 0x13, words: map[uint16]uint16{0x1e5c: 24, 0x1e5e: 20, 0x1e60: 10, 0x1e62: 10, 0x15be: 24, 0x15c0: 20, 0x15ac: 30, 0x15ae: 53, 0x1e64: 2, 0x1e66: 0x84d3}, tail: read("WHO-source-bitmap.bin"), prefix: read("WHO-before-tail.bin")}
	var m map[string]entry
	if json.Unmarshal(read(backend+"-xlate.json"), &m) != nil {
		t.Fatal("glyph resource")
	}
	a := make([]byte, 64000)
	for i, c := range f.tail {
		if c != 255 {
			a[(67+i/24)*320+30+i%24] = c
		} else {
			a[(67+i/24)*320+30+i%24] = 117
		}
	}
	return f, m["WHO?"], a
}
func TestWHONativeSourceAndGeometry(t *testing.T) {
	for _, backend := range []string{"eten", "unifont"} {
		for _, k := range []int{2, 3} {
			f, e, a := whoFixture(t, backend)
			before := append([]byte(nil), a...)
			v := whoCursorFromSource(f, e, k, a)
			if v == nil {
				t.Fatal("valid WHO source rejected", backend, k)
			}
			h := len(e.G[0])
			if v.glyphX < 0 || v.glyphX+v.advance+h > 24*k || v.glyphY < 0 || v.glyphY+h > 20*k {
				t.Fatal("whole glyph cells outside")
			}
			_, hi := whoInkColumns(e.G[0])
			lo, _ := whoInkColumns(e.G[1])
			if v.advance+lo-hi != 2 {
				t.Fatal("ink gap")
			}
			comp := append([]byte(nil), a...)
			v.restore(comp)
			count := 0
			for i, b := range comp {
				if b != before[i] {
					count++
					if b != 0 {
						t.Fatal("wrong restoration")
					}
				}
			}
			if count != 45 || !bytes.Equal(a, before) {
				t.Fatal("original modified or wrong ink count")
			}
			for i, c := range f.tail {
				if c == 255 && comp[(67+i/24)*320+30+i%24] != 117 {
					t.Fatal("transparent point changed")
				}
			}
			e.G[0][0] ^= 1
			if bytes.Equal([]byte{byte(v.glyphs[0][0])}, []byte{byte(e.G[0][0])}) {
				t.Fatal("view shared glyph source")
			}
		}
	}
}
func TestWHORejectsBeforePointerRead(t *testing.T) {
	for _, kind := range []string{"mode", "indexed", "width", "height", "hotspot", "crop", "x", "y", "span", "physical", "glyph", "scale"} {
		f, e, a := whoFixture(t, "unifont")
		k := 2
		switch kind {
		case "mode":
			f.mode = 3
		case "indexed":
			a = a[:100]
		case "width":
			f.words[0x1e5c] = 23
		case "height":
			f.words[0x1e5e] = 19
		case "hotspot":
			f.words[0x1e60] = 9
		case "crop":
			f.words[0x15be] = 23
		case "x":
			f.words[0x15ac] = 65535
		case "y":
			f.words[0x15ae] = 190
		case "span":
			f.words[0x1e64] = 65057
		case "physical":
			f.words[0x1e66] = 0xa000
		case "glyph":
			e.G = nil
		case "scale":
			k = 4
		}
		if whoCursorFromSource(f, e, k, a) != nil || f.reads != 0 {
			t.Fatal("unsafe source read", kind, f.reads)
		}
	}
}
func TestWHORejectsIdentityAndIncompletePaint(t *testing.T) {
	for _, kind := range []string{"tail", "prefix", "ink", "opaque", "short"} {
		f, e, a := whoFixture(t, "eten")
		switch kind {
		case "tail":
			f.tail[0] ^= 1
		case "prefix":
			f.prefix[0] ^= 1
		case "ink":
			for i, c := range f.tail {
				if c == 22 {
					a[(67+i/24)*320+30+i%24] = 0
					break
				}
			}
		case "opaque":
			for i, c := range f.tail {
				if c == 0 {
					a[(67+i/24)*320+30+i%24] = 1
					break
				}
			}
		case "short":
			f.tail = f.tail[:100]
		}
		if whoCursorFromSource(f, e, 2, a) != nil {
			t.Fatal("incorrect WHO image accepted", kind)
		}
	}
}
func TestWHORejectsBadGlyphResource(t *testing.T) {
	for _, kind := range []string{"dy", "slot", "count", "height", "mixed", "bits", "empty"} {
		f, e, a := whoFixture(t, "eten")
		switch kind {
		case "dy":
			e.Dy = 1
		case "slot":
			e.D = []int{0}
		case "count":
			e.G = e.G[:1]
		case "height":
			e.G[0] = e.G[0][:14]
		case "mixed":
			e.G[1] = e.G[1][:14]
		case "bits":
			e.G[0][0] = 1 << 15
		case "empty":
			e.G[1] = make([]uint16, 15)
		}
		if whoCursorFromSource(f, e, 2, a) != nil || f.reads != 0 {
			t.Fatal("bad glyph accepted", kind)
		}
	}
}

func TestWHOOutputChannelsAndCursorOnlyKey(t *testing.T) {
	f, e, a := whoFixture(t, "unifont")
	trans := map[string]entry{"WHO?": e, "HUMAN": {Dy: 0, G: [][]uint16{{1}}}}
	extracted := takeWhoCursorEntry(trans)
	if _, ok := trans["WHO?"]; ok || len(trans) != 1 || trans["HUMAN"].G[0][0] != 1 || !whoGlyphsValid(extracted) {
		t.Fatal("WHO leaked into ordinary NUL map")
	}
	gui := whoCursorFromSource(f, extracted, 2, a)
	hd := whoCursorFromSource(f, extracted, 3, a)
	if gui == nil || hd == nil || gui.k != 2 || hd.k != 3 || gui.glyphY == hd.glyphY {
		t.Fatal("output k not independent")
	}
	rejected := whoCursorFromSource(f, extracted, 4, a)
	untouched := append([]byte(nil), a...)
	rejected.restore(untouched)
	if !bytes.Equal(untouched, a) {
		t.Fatal("rejected output lost original WHO")
	}
	valid := append([]byte(nil), a...)
	gui.restore(valid)
	if bytes.Equal(valid, a) || !bytes.Equal(untouched, a) {
		t.Fatal("output buffers were shared")
	}
}

func TestWHORealDrawWrongBoundsPreservesOriginal(t *testing.T) {
	if os.Getenv("WIZ7_WHO_DRAW_CHILD") != "1" {
		whoFixture(t, "unifont")
		exe, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(exe, "-test.run=^TestWHORealDrawWrongBoundsPreservesOriginal$")
		cmd.Env = append(os.Environ(), "WIZ7_WHO_DRAW_CHILD=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("isolated Draw check: %v\n%s", err, output)
		}
		return
	}
	f, e, a := whoFixture(t, "unifont")
	v := whoCursorFromSource(f, e, 2, a)
	check := &whoDrawCheck{t: t, f: f, a: a, v: v}
	if err := ebiten.RunGame(check); err != nil {
		t.Fatal(err)
	}
	if !check.done {
		t.Fatal("Draw check did not execute")
	}
}

type whoDrawCheck struct {
	t    *testing.T
	f    *whoFake
	a    []byte
	v    *whoCursorView
	done bool
}

func (c *whoDrawCheck) Update() error {
	if c.done {
		return ebiten.Termination
	}
	return nil
}

func (c *whoDrawCheck) Layout(int, int) (int, int) { return 64, 64 }

func (check *whoDrawCheck) Draw(*ebiten.Image) {
	if check.done {
		return
	}
	check.done = true
	t, f, a, v := check.t, check.f, check.a, check.v
	rgba := make([]byte, 320*200*4)
	pal := f.Palette()
	for i, c := range a {
		rgba[4*i], rgba[4*i+1], rgba[4*i+2], rgba[4*i+3] = pal[c][0], pal[c][1], pal[c][2], 255
	}
	cases := []struct {
		name   string
		scale  int
		bounds image.Rectangle
		sub    bool
	}{
		{"height", 2, image.Rect(0, 0, 640, 399), false},
		{"width", 2, image.Rect(0, 0, 639, 400), false},
		{"origin", 2, image.Rect(1, 1, 641, 401), true},
		{"scale", 3, image.Rect(0, 0, 960, 600), false},
	}
	draw := func(c struct {
		name   string
		scale  int
		bounds image.Rectangle
		sub    bool
	}, who *whoCursorView) []byte {
		var dst *ebiten.Image
		if c.sub {
			parent := ebiten.NewImage(650, 410)
			defer parent.Dispose()
			dst = parent.SubImage(c.bounds).(*ebiten.Image)
		} else {
			dst = ebiten.NewImage(c.bounds.Dx(), c.bounds.Dy())
			defer dst.Dispose()
		}
		screen := ebiten.NewImage(320, 200)
		defer screen.Dispose()
		g := &game{scale: c.scale, screen: screen, sh: &shared{rgba: append([]byte(nil), rgba...), who: who}}
		g.Draw(dst)
		pixels := make([]byte, dst.Bounds().Dx()*dst.Bounds().Dy()*4)
		dst.ReadPixels(pixels)
		if g.whoImage != nil {
			g.whoImage.Dispose()
		}
		if !bytes.Equal(g.sh.rgba, rgba) {
			t.Fatal("shared original frame changed")
		}
		return pixels
	}
	for _, c := range cases {
		if v.canDraw(c.bounds, c.scale) {
			t.Fatal("wrong output accepted", c.name)
		}
		base := draw(c, nil)
		actual := draw(c, v)
		if !bytes.Equal(base, actual) {
			t.Fatal("rejected Draw lost original WHO", c.name)
		}
	}
	for _, k := range []int{2, 3} {
		view := whoCursorFromSource(f, entry{G: v.glyphs}, k, a)
		render := func(who *whoCursorView) []byte {
			dst := ebiten.NewImage(320*k, 200*k+40)
			defer dst.Dispose()
			screen := ebiten.NewImage(320, 200)
			defer screen.Dispose()
			g := &game{scale: k, screen: screen, sh: &shared{rgba: append([]byte(nil), rgba...), who: who}, live: &inputRecording{font: &helpData{}}}
			g.Draw(dst)
			pixels := make([]byte, 320*k*(200*k+40)*4)
			dst.ReadPixels(pixels)
			if g.whoImage != nil {
				g.whoImage.Dispose()
			}
			if !bytes.Equal(g.sh.rgba, rgba) {
				t.Fatal("footer Draw changed shared original frame")
			}
			return pixels
		}
		base, actual := render(nil), render(view)
		if !bytes.Equal(base[320*k*200*k*4:], actual[320*k*200*k*4:]) {
			t.Fatal("WHO changed recording footer", k)
		}
		changed := 0
		for i := 0; i < 320*k*200*k; i++ {
			if bytes.Equal(base[4*i:4*i+4], actual[4*i:4*i+4]) {
				continue
			}
			changed++
			x, y := i%(320*k), i/(320*k)
			if x < 30*k || x >= 54*k || y < 65*k || y >= 73*k {
				t.Fatal("WHO changed outside approved text area", k, x, y)
			}
		}
		if changed == 0 {
			t.Fatal("WHO missing with recording footer", k)
		}
	}
}
