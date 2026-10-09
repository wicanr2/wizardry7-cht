package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"image"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

type heldFake struct {
	mem   map[int][]byte
	reads int
	end   int
	regs  oracle.Regs
	step  uint64
}

func (f *heldFake) Regs() oracle.Regs { return f.regs }
func (f *heldFake) Steps() uint64     { return f.step }

func (f *heldFake) VideoMode() uint8 { return 0x13 }
func (f *heldFake) Palette() (p [256][3]uint8) {
	for i := range p {
		p[i] = [3]uint8{byte(i), byte(i), byte(i)}
	}
	return
}
func (f *heldFake) Bytes(a oracle.Addr, n int) []byte {
	f.reads++
	at := int(a.Seg)*16 + int(a.Off)
	f.end = max(f.end, at+n)
	for base, b := range f.mem {
		if at >= base && at+n <= base+len(b) {
			return b[at-base : at-base+n]
		}
	}
	return nil
}
func (f *heldFake) put(a oracle.Addr, b []byte) {
	f.mem[int(a.Seg)*16+int(a.Off)] = append([]byte(nil), b...)
}
func heldWords(values ...uint16) []byte {
	b := make([]byte, len(values)*2)
	for i, v := range values {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func heldFixture(x, y int) (*heldFake, *heldItemCursor, []byte) {
	f := &heldFake{mem: map[int][]byte{}}
	for off, value := range map[uint16]string{0x5c24: heldCopyCode, 0x5c9f: heldDefinitionCode} {
		b, _ := hex.DecodeString(value)
		f.put(oracle.Far(0x0110, off), b)
	}
	f.put(oracle.Far(0x0110, 0x2de7), heldWords(0x0163, 0x740c))
	f.put(oracle.Far(0x740c, 0x0163), []byte{0xe8, 0x63, 0x13, 0xcb})
	w, h := min(heldWidth, 320-x), min(heldHeight, 200-y)
	f.put(oracle.Far(0x740c, 0x15ac), heldWords(uint16(x), uint16(y), uint16(x), uint16(y), 0x3000, 0x3400, uint16(w), uint16(h), 0, uint16(w), uint16(h), 0))
	f.put(oracle.Far(0x740c, 0x1e5c), heldWords(76, 13, 10, 10, 0, 0x8500))
	f.put(oracle.Far(0x740c, 0x1ed3), []byte{0})
	bg := make([]byte, w*h)
	for i := range bg {
		bg[i] = byte(20 + i%60)
	}
	f.put(oracle.Far(0x740c, 0x3000), bg)
	pic := bytes.Repeat([]byte{255}, heldPixels)
	pic[0] = 45
	pic[heldWidth+16] = 105
	f.put(oracle.Far(0x8500, 0), pic)
	a := bytes.Repeat([]byte{9}, 64000)
	for cy := 0; cy < h; cy++ {
		copy(a[(y+cy)*320+x:], bg[cy*w:(cy+1)*w])
	}
	a[y*320+x] = 45
	if w > 16 && h > 1 {
		a[(y+1)*320+x+16] = 105
	}
	mask := make([]byte, heldPixels)
	mask[heldWidth+16] = 1
	g := make([]uint16, 15)
	g[0] = 1 << 14
	label := &heldLabel{key: "LONGSWORD", ptr: oracle.Far(0x8500, 0), picture: append([]byte(nil), pic...), mask: mask, glyphs: [][]uint16{g}, fg: 105}
	return f, &heldItemCursor{label: label}, a
}
func TestHeldReadRejectsSpanBeforeReader(t *testing.T) {
	f := &heldFake{mem: map[int][]byte{}}
	for _, c := range []struct {
		a oracle.Addr
		n int
	}{{oracle.Far(0x1000, 0xfffe), 3}, {oracle.Far(0xa000, 0), 1}, {oracle.Far(0x1000, 0), 0}} {
		if heldRead(f, c.a, c.n) != nil {
			t.Fatal("invalid span accepted")
		}
	}
	if f.reads != 0 {
		t.Fatal("reader called before span rejected")
	}
}
func TestHeldTextStopsAtFirstNUL(t *testing.T) {
	f := &heldFake{mem: map[int][]byte{}}
	p := oracle.Far(0x1106, 0x200)
	f.put(p, []byte("LEATHER\x00CUIRASS\x00"))
	s, ok := heldText(f, p)
	if !ok || s != "LEATHER" || f.end != int(p.Seg)*16+int(p.Off)+8 {
		t.Fatal(s, ok, f.end)
	}
}
func TestHeldMaskRejectsUnchangedOldInk(t *testing.T) {
	before := bytes.Repeat([]byte{255}, heldPixels)
	after := append([]byte(nil), before...)
	after[heldWidth+16] = 105
	mask, fg, ok := heldNameMask(before, after, 1)
	if !ok || fg != 105 || mask[heldWidth+16] != 1 {
		t.Fatal("valid source rejected")
	}
	before[heldWidth+20] = 105
	after[heldWidth+20] = 105
	if _, _, ok := heldNameMask(before, after, 1); ok {
		t.Fatal("unchanged old English ink inherited")
	}
}
func TestHeldMaskRejectsPixelsOutsideTextBand(t *testing.T) {
	before := bytes.Repeat([]byte{255}, heldPixels)
	for _, point := range []int{heldWidth + 15, 6*heldWidth + 20, 12*heldWidth + 20} {
		after := append([]byte(nil), before...)
		after[heldWidth+16] = 105
		after[point] = 105
		if _, _, ok := heldNameMask(before, after, 1); ok {
			t.Fatal("non-name pixel accepted", point)
		}
	}
}
func TestHeldClearGenerationDoesNotClearCurrentLabel(t *testing.T) {
	h := &heldItemCursor{label: &heldLabel{key: "LONGSWORD"}, pending: &heldPart{}, name: &heldName{}, capture: &heldCapture{}}
	h.clearCandidate()
	if h.pending != nil || h.name != nil || h.capture != nil || h.label == nil {
		t.Fatal("candidate and current label lifetimes mixed")
	}
}
func TestHeldBackgroundClipAndOwnedCopies(t *testing.T) {
	f, h, a := heldFixture(300, 198)
	v, pending := h.view(f, a)
	if v == nil || len(pending) != 0 || v.w != 20 || v.h != 2 {
		t.Fatal("right/bottom clip rejected")
	}
	base := v.base(a)
	if base[198*320+300] != 20 || base[0] != 9 {
		t.Fatal("background restoration escaped ROI")
	}
	f.mem[0x85000][0] = 77
	if v.picture[0] != 45 {
		t.Fatal("view kept original buffer reference")
	}
	img := image.NewRGBA(image.Rect(0, 0, 960, 600))
	v.paint(img, 3)
	if img.RGBAAt(900, 594).R != 45 {
		t.Fatal("non-text icon changed")
	}
}
func TestHeldSamePointerNewPictureDoesNotKeepName(t *testing.T) {
	f, h, a := heldFixture(100, 100)
	f.mem[0x85000][heldWidth+16] = 31
	a[101*320+116] = 31
	v, _ := h.view(f, a)
	if v == nil || v.label != nil || h.label != nil || v.translated(3) {
		t.Fatal("reused pointer inherited old name")
	}
	if v.raster(3).RGBAAt(48, 3).R != 31 {
		t.Fatal("unknown replacement picture lost original English")
	}
}
func TestHeldBusyKeepsOccludedSourceCandidates(t *testing.T) {
	f, h, a := heldFixture(20, 140)
	f.mem[0x740c0+0x1ed3][0] = 255
	v, pending := h.view(f, a)
	if v != nil || len(pending) == 0 {
		t.Fatal("busy background used")
	}
	e := &textEv{bx: 18, by: 142, bw: 41, bh: 11, after: bytes.Repeat([]byte{8}, 451), shown: true}
	off := make([]byte, 64000)
	for y := 0; y < 11; y++ {
		copy(off[(142+y)*320+18:], e.after[y*41:(y+1)*41])
	}
	if !heldIntersects(e, pending) || !heldEventStill(e, off) {
		t.Fatal("occluded source not retained")
	}
	if e.dead {
		t.Fatal("busy cursor killed old event")
	}
	off[142*320+18] = 7
	if heldEventStill(e, off) {
		t.Fatal("changed underlying source retained")
	}
}
func TestHeldScaleFallbackStillPaintsOriginalCursor(t *testing.T) {
	f, h, a := heldFixture(100, 100)
	v, _ := h.view(f, a)
	g := make([]uint16, 16)
	g[0] = 1 << 15
	v.label.glyphs = make([][]uint16, 7)
	for i := range v.label.glyphs {
		v.label.glyphs[i] = append([]uint16(nil), g...)
	}
	if v.translated(2) || !v.translated(3) {
		t.Fatal("native size was shrunk or fit ignored")
	}
	for _, k := range []int{1, 2, 4} {
		if v.raster(k).RGBAAt(16*k, k).R != 105 {
			t.Fatal("fallback lost original image", k)
		}
	}
}

func heldFlowFixture() (*heldFake, *heldItemCursor, []byte) {
	f, _, _ := heldFixture(0, 110)
	f.put(oracle.Far(0x740c, 0x01e5), heldWords(0x6000))
	f.put(oracle.Far(0x740c, 0x0295), heldWords(18, 114))
	frame := bytes.Repeat([]byte{255}, 64000)
	f.put(oracle.Far(0x6000, 0), frame)
	frame = f.mem[0x60000]
	g := make([]uint16, 16)
	g[0] = 1 << 15
	h := &heldItemCursor{catalog: map[string]entry{"LONGSWORD": {G: [][]uint16{g, g}}, "LEATHER/CUIRASS": {G: [][]uint16{g, g, g}}}}
	f.put(oracle.Far(0x1106, 0x0e1e), heldWords(0, 0x8500))
	return f, h, frame
}
func heldBegin(f *heldFake, h *heldItemCursor, s string, off, ret uint16, y int, step uint64) {
	f.regs = oracle.Regs{CS: 0x0110, IP: 0x3895, DS: 0x1106, SS: 0x1106, SP: 0xfe00}
	f.step = step
	f.put(oracle.Far(0x1106, 0xfe00), heldWords(ret, off, 0))
	f.put(oracle.Far(0x1106, off), append([]byte(s), 0))
	f.put(oracle.Far(0x740c, 0x0295), heldWords(18, uint16(y)))
	h.enterName(f)
}
func heldFinish(f *heldFake, h *heldItemCursor, ret uint16, step uint64) {
	f.regs.IP = ret
	f.regs.SP = 0xfe02
	f.step = step
	h.returnName(f, ret)
}
func heldCaptureAndDefine(f *heldFake, h *heldItemCursor, top int, step uint64) {
	f.regs = oracle.Regs{CS: 0x0110, IP: 0x396f, DS: 0x1106, SS: 0x1106, SP: 0xfd00}
	f.step = step
	f.put(oracle.Far(0x1106, 0xfd00), heldWords(0x5c1b, 2, uint16(top), 77, uint16(top+12)))
	h.captureName(f)
	if h.capture != nil {
		f.put(oracle.Far(0x8500, 0), h.capture.picture)
	}
	f.regs.IP = 0x40a7
	f.regs.SP = 0xfc00
	f.step = step + 100
	f.put(oracle.Far(0x1106, 0xfc00), heldWords(0x5cb9, 0, 0x8500, 76, 13, 0, uint16(top-3)))
	h.enterDefinition(f)
	f.regs.IP = 0x40b0
	f.regs.SP = 0xfbfe
	f.step = step + 200
	h.returnDefinition(f)
}
func TestHeldCompleteTwoNamesAndResetChain(t *testing.T) {
	f, h, frame := heldFlowFixture()
	heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 1000)
	if h.pending == nil {
		t.Fatal("first source not accepted")
	}
	frame[114*320+18] = 105
	heldFinish(f, h, 0xba40, 2000)
	heldBegin(f, h, "", 0x020a, 0xba63, 120, 2036)
	if h.pending == nil {
		t.Fatal("empty second did not wait for return")
	}
	heldFinish(f, h, 0xba63, 2200)
	if h.name == nil || len(h.name.parts) != 1 || h.name.key() != "LONGSWORD" {
		t.Fatal("empty second lost single name")
	}
	heldCaptureAndDefine(f, h, 113, 3000)
	if h.label == nil || h.label.key != "LONGSWORD" || h.label.captured != 3000 {
		t.Fatal("single name chain failed")
	}
	// Ordinary printed LONGSWORD over texture must clear stale candidates.
	for y := 113; y < 126; y++ {
		for x := 2; x < 78; x++ {
			frame[y*320+x] = 117
		}
	}
	heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 4000)
	if h.pending != nil || h.name != nil || h.capture != nil {
		t.Fatal("ordinary NUL5 became scratch")
	}
	if h.label == nil {
		t.Fatal("ordinary repaint cleared live cursor")
	}
	// A new image at the same pointer cannot inherit that old name.
	heldCaptureAndDefine(f, h, 127, 5000)
	if h.label != nil {
		t.Fatal("ordinary NUL5 inherited by replacement definition")
	}
	for y := 127; y < 140; y++ {
		for x := 2; x < 78; x++ {
			frame[y*320+x] = 255
		}
	}
	heldBegin(f, h, "LEATHER", 0x0300, 0xba40, 128, 10000)
	frame[128*320+18] = 105
	heldFinish(f, h, 0xba40, 11000)
	if h.name == nil {
		t.Fatal("first without single dictionary key discarded")
	}
	heldBegin(f, h, "CUIRASS", 0x0308, 0xba63, 134, 11036)
	frame[134*320+18] = 105
	heldFinish(f, h, 0xba63, 12000)
	heldCaptureAndDefine(f, h, 127, 13000)
	if h.label == nil || h.label.key != "LEATHER/CUIRASS" || len(h.label.glyphs) != 3 {
		t.Fatal("two-line identity failed")
	}
	f.regs = oracle.Regs{CS: 0x0110, IP: 0x40a7, DS: 0x1106, SS: 0x1106, SP: 0xfc00}
	f.put(oracle.Far(0x1106, 0xfc00), heldWords(0x5b80, 0, 0, 0, 0, 27, 130))
	h.enterDefinition(f)
	if h.label != nil || h.definition != nil || h.capture != nil {
		t.Fatal("reset retained name")
	}
}
func TestHeldReturnAliasesCannotCompleteOrConsumeSource(t *testing.T) {
	f, h, frame := heldFlowFixture()
	heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 1000)
	frame[114*320+18] = 105
	f.regs.CS = 0x0111
	f.regs.IP = 0xba30
	f.regs.SP = 0xfe02
	f.step = 2000
	h.returnName(f, 0xba40)
	if h.name != nil || h.pending == nil {
		t.Fatal("linear alias treated as source return")
	}
	f.regs.CS = 0x0110
	heldFinish(f, h, 0xba40, 2100)
	f.regs.IP = 0x396f
	f.regs.SP = 0xfd00
	f.step = 3000
	f.put(oracle.Far(0x1106, 0xfd00), heldWords(0x5c1b, 2, 113, 77, 125))
	h.captureName(f)
	f.put(oracle.Far(0x8500, 0), h.capture.picture)
	f.regs.IP = 0x40a7
	f.regs.SP = 0xfc00
	f.step = 3100
	f.put(oracle.Far(0x1106, 0xfc00), heldWords(0x5cb9, 0, 0x8500, 76, 13, 0, 110))
	h.enterDefinition(f)
	f.regs.CS = 0x0111
	f.regs.IP = 0x40a0
	f.regs.SP = 0xfbfe
	f.step = 3200
	h.returnDefinition(f)
	if h.label != nil || h.definition == nil {
		t.Fatal("linear alias registered cursor")
	}
	f.regs.CS = 0x0110
	f.regs.IP = 0x40b0
	h.returnDefinition(f)
	if h.label == nil {
		t.Fatal("real return lost after unrelated alias")
	}
}
func TestHeldBadSPAndCapturePixelsRejectName(t *testing.T) {
	f, h, frame := heldFlowFixture()
	heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 1000)
	frame[114*320+18] = 105
	f.regs.IP = 0xba40
	f.regs.SP = 0xfe04
	f.step = 2000
	h.returnName(f, 0xba40)
	if h.name != nil {
		t.Fatal("wrong return SP accepted")
	}
	frame[114*320+18] = 255
	heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 3000)
	frame[114*320+18] = 105
	heldFinish(f, h, 0xba40, 4000)
	// Changed source letter after its return cannot own final capture.
	frame[114*320+18] = 31
	heldCaptureAndDefine(f, h, 113, 5000)
	if h.label != nil {
		t.Fatal("overwritten source pixels retained name")
	}
}
func TestHeldSourceSignatureLossKeepsOnlyOriginalPicture(t *testing.T) {
	f, h, a := heldFixture(100, 100)
	f.mem[0x01100+0x5c24][0] ^= 1
	v, _ := h.view(f, a)
	if v == nil || v.label != nil || h.label != nil {
		t.Fatal("source version loss retained Chinese or dropped raw layer")
	}
	if v.raster(3).RGBAAt(48, 3).R != 105 {
		t.Fatal("source loss erased original name")
	}
}

func TestHeldEmptyPartsCannotGuessOrChangeSource(t *testing.T) {
	for _, kind := range []string{"first", "unpaired", "wrong-pointer", "wrong-gap", "wrong-SP", "changed-pixel"} {
		t.Run(kind, func(t *testing.T) {
			f, h, frame := heldFlowFixture()
			if kind == "first" {
				heldBegin(f, h, "", 0x0200, 0xba40, 114, 1000)
				if h.pending != nil || h.name != nil {
					t.Fatal("empty first accepted")
				}
				return
			}
			if kind == "unpaired" {
				heldBegin(f, h, "", 0x020a, 0xba63, 120, 2036)
				if h.pending != nil || h.name != nil {
					t.Fatal("unpaired empty accepted")
				}
				return
			}
			heldBegin(f, h, "LONGSWORD", 0x0200, 0xba40, 114, 1000)
			frame[114*320+18] = 105
			heldFinish(f, h, 0xba40, 2000)
			off := uint16(0x020a)
			step := uint64(2036)
			if kind == "wrong-pointer" {
				off++
			}
			if kind == "wrong-gap" {
				step++
			}
			heldBegin(f, h, "", off, 0xba63, 120, step)
			if kind == "wrong-SP" {
				f.regs.IP = 0xba63
				f.regs.SP = 0xfe04
				f.step = 2200
				h.returnName(f, 0xba63)
			} else {
				if kind == "changed-pixel" {
					frame[125*320+77] = 21
				}
				heldFinish(f, h, 0xba63, 2200)
			}
			if h.name != nil || h.pending != nil || h.capture != nil {
				t.Fatal("invalid empty part retained source")
			}
		})
	}
}
