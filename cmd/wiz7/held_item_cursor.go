package main

// spec037：名稱與原游標圖層分開；所有原版資料只讀。
import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"image"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

const heldWidth, heldHeight = 76, 13
const heldPixels = heldWidth * heldHeight
const heldSourceWindow uint64 = 1000000

const heldDefinitionCode = "8946fa50ff76fcb80d0050b84c0050ff36200eff361e0ee8eee383c40ce89fdf"
const heldCopyCode = "b8dc0350ff36200eff361e0e8b46ee050c005250e8d8a383c40ab8010050ff76"

type heldMachineReader interface {
	heldReader
	Regs() oracle.Regs
	Steps() uint64
}

func heldAt(r oracle.Regs, ip uint16) bool { return r.CS == 0x0110 && r.IP == ip }

type heldReader interface {
	VideoMode() uint8
	Bytes(oracle.Addr, int) []byte
	Palette() [256][3]uint8
}

func heldRead(o heldReader, a oracle.Addr, n int) []byte {
	if n <= 0 || n > 65536 || int(a.Off)+n > 65536 || int(a.Seg)*16+int(a.Off)+n > 0xa0000 {
		return nil
	}
	b := o.Bytes(a, n)
	if len(b) != n {
		return nil
	}
	return append([]byte(nil), b...)
}
func heldWord(b []byte, at int) uint16     { return binary.LittleEndian.Uint16(b[at : at+2]) }
func heldPtr(b []byte, at int) oracle.Addr { return oracle.Far(heldWord(b, at+2), heldWord(b, at)) }
func heldZero(a oracle.Addr) bool          { return a.Seg == 0 && a.Off == 0 }
func heldSig(o heldReader, off uint16, value string) bool {
	b, err := hex.DecodeString(value)
	return err == nil && bytes.Equal(heldRead(o, oracle.Far(0x0110, off), len(b)), b)
}
func heldModule(o heldReader) bool {
	return heldSig(o, 0x5c24, heldCopyCode) && heldSig(o, 0x5c9f, heldDefinitionCode)
}
func heldDriver(o heldReader) (uint16, bool) {
	b := heldRead(o, oracle.Far(0x0110, 0x2de7), 4)
	if len(b) != 4 || heldWord(b, 0) != 0x0163 {
		return 0, false
	}
	seg := heldWord(b, 2)
	return seg, bytes.Equal(heldRead(o, oracle.Far(seg, 0x0163), 4), []byte{0xe8, 0x63, 0x13, 0xcb})
}
func heldField(o heldReader, driver uint16, r image.Rectangle) []byte {
	if r.Dx() != heldWidth || r.Dy() != heldHeight || !r.In(image.Rect(0, 0, 320, 200)) {
		return nil
	}
	b := heldRead(o, oracle.Far(driver, 0x01e5), 2)
	if len(b) != 2 || heldWord(b, 0) == 0 {
		return nil
	}
	seg := heldWord(b, 0)
	result := make([]byte, heldPixels)
	for y := 0; y < heldHeight; y++ {
		row := heldRead(o, oracle.Far(seg, uint16((r.Min.Y+y)*320+r.Min.X)), heldWidth)
		if len(row) != heldWidth {
			return nil
		}
		copy(result[y*heldWidth:], row)
	}
	return result
}
func heldText(o heldReader, a oracle.Addr) (string, bool) {
	b := make([]byte, 0, 31)
	for i := 0; i < 32; i++ {
		if int(a.Off)+i > 65535 {
			return "", false
		}
		one := heldRead(o, oracle.Far(a.Seg, a.Off+uint16(i)), 1)
		if len(one) != 1 {
			return "", false
		}
		if one[0] == 0 {
			return string(b), true
		}
		if one[0] < 32 || one[0] > 126 {
			return "", false
		}
		b = append(b, one[0])
	}
	return "", false
}
func heldBandEmpty(b []byte, line int) bool {
	if len(b) != heldPixels || (line != 1 && line != 7) {
		return false
	}
	for y := line; y < line+5; y++ {
		for x := 16; x < heldWidth; x++ {
			if b[y*heldWidth+x] != 255 {
				return false
			}
		}
	}
	return true
}
func heldNameMask(before, after []byte, line int) ([]byte, byte, bool) {
	if !heldBandEmpty(before, line) || len(after) != heldPixels {
		return nil, 0, false
	}
	mask := make([]byte, heldPixels)
	count := 0
	var fg byte
	for i, v := range after {
		inBand := i%heldWidth >= 16 && i/heldWidth >= line && i/heldWidth < line+5
		if !inBand {
			if v != before[i] {
				return nil, 0, false
			}
			continue
		}
		if v == 255 {
			continue
		}
		if count > 0 && v != fg {
			return nil, 0, false
		}
		fg = v
		mask[i] = 1
		count++
	}
	return mask, fg, count > 0
}

type heldPart struct {
	s                    string
	ds, off, ss, sp, ret uint16
	x, y, line           int
	entry, returned      uint64
	before, after, mask  []byte
	fg                   byte
}
type heldName struct {
	gen         uint64
	parts       []heldPart
	rect        image.Rectangle
	deadline    uint64
	after, mask []byte
	fg          byte
}

func (n *heldName) key() string {
	if len(n.parts) == 2 {
		return n.parts[0].s + "/" + n.parts[1].s
	}
	return n.parts[0].s
}

type heldCapture struct {
	name    *heldName
	picture []byte
	step    uint64
}
type heldLabel struct {
	key               string
	ptr               oracle.Addr
	picture, mask     []byte
	glyphs            [][]uint16
	fg                byte
	source            *heldName
	captured, defined uint64
}
type heldDefinition struct {
	label          *heldLabel
	ss, sp         uint16
	step, deadline uint64
}
type heldItemCursor struct {
	catalog                                                         map[string]entry
	gen                                                             uint64
	pending                                                         *heldPart
	name                                                            *heldName
	capture                                                         *heldCapture
	definition                                                      *heldDefinition
	label                                                           *heldLabel
	lastRect                                                        image.Rectangle
	rejectedBand, rejectedPair, rejectedCapture, rejectedDefinition uint64
}

func (h *heldItemCursor) clearCandidate() { h.pending = nil; h.name = nil; h.capture = nil }
func heldOrigin(x, y, ret int) (image.Rectangle, int, bool) {
	line := 1
	firstY := y
	if ret == 0xba63 {
		line = 7
		firstY -= 6
	} else if ret != 0xba40 {
		return image.Rectangle{}, 0, false
	}
	if (x != 18 && x != 97) || firstY < 114 || firstY > 170 || (firstY-114)%14 != 0 {
		return image.Rectangle{}, 0, false
	}
	return image.Rect(x-16, firstY-1, x-16+heldWidth, firstY-1+heldHeight), line, true
}
func (h *heldItemCursor) install(o *oracle.Oracle, mu *sync.Mutex) {
	if h == nil || len(h.catalog) == 0 {
		return
	}
	guard := func(f func(*oracle.Oracle)) func(*oracle.Oracle) {
		return func(o *oracle.Oracle) { mu.Lock(); defer mu.Unlock(); f(o) }
	}
	o.OnCall(oracle.Far(0x0110, 0x3895), guard(func(o *oracle.Oracle) { h.enterName(o) }))
	for _, ret := range []uint16{0xba40, 0xba63} {
		target := ret
		o.OnCall(oracle.Far(0x0110, ret), guard(func(o *oracle.Oracle) { h.returnName(o, target) }))
	}
	o.OnCall(oracle.Far(0x0110, 0x396f), guard(func(o *oracle.Oracle) { h.captureName(o) }))
	o.OnCall(oracle.Far(0x0110, 0x40a7), guard(func(o *oracle.Oracle) { h.enterDefinition(o) }))
	o.OnCall(oracle.Far(0x0110, 0x40b0), guard(func(o *oracle.Oracle) { h.returnDefinition(o) }))
}
func (h *heldItemCursor) enterName(o heldMachineReader) {
	r := o.Regs()
	if !heldAt(r, 0x3895) || o.VideoMode() != 0x13 || !heldModule(o) {
		return
	}
	b := heldRead(o, oracle.Far(r.SS, r.SP), 6)
	if len(b) != 6 {
		h.clearCandidate()
		return
	}
	ret, off, font := heldWord(b, 0), heldWord(b, 2), heldWord(b, 4)
	if ret != 0xba40 && ret != 0xba63 {
		return
	}
	if ret == 0xba40 {
		h.gen++
		h.clearCandidate()
	}
	driver, ok := heldDriver(o)
	if !ok || font != 0 || r.SP > 65533 {
		h.clearCandidate()
		return
	}
	xy := heldRead(o, oracle.Far(driver, 0x0295), 4)
	if len(xy) != 4 {
		h.clearCandidate()
		return
	}
	x, y := int(heldWord(xy, 0)), int(heldWord(xy, 2))
	rect, line, ok := heldOrigin(x, y, int(ret))
	if !ok {
		h.clearCandidate()
		return
	}
	before := heldField(o, driver, rect)
	if !heldBandEmpty(before, line) {
		h.rejectedBand++
		h.clearCandidate()
		return
	}
	s, ok := heldText(o, oracle.Far(r.DS, off))
	if !ok || ret == 0xba40 && s == "" {
		h.clearCandidate()
		return
	}
	if ret == 0xba63 {
		n := h.name
		if n == nil || len(n.parts) != 1 || h.pending != nil {
			h.rejectedPair++
			h.clearCandidate()
			return
		}
		first := n.parts[0]
		if r.DS != first.ds || r.SS != first.ss || r.SP != first.sp || x != first.x || y != first.y+6 ||
			int(first.off)+len(first.s)+1 > 65535 || int(off) != int(first.off)+len(first.s)+1 ||
			o.Steps() < first.returned || o.Steps()-first.returned != 36 || !bytes.Equal(before, n.after) {
			h.rejectedPair++
			h.clearCandidate()
			return
		}
	}
	h.pending = &heldPart{s: s, ds: r.DS, off: off, ss: r.SS, sp: r.SP, ret: ret, x: x, y: y, line: line, entry: o.Steps(), before: before}
}
func (h *heldItemCursor) returnName(o heldMachineReader, ret uint16) {
	if !heldAt(o.Regs(), ret) {
		return
	}
	p := h.pending
	if p == nil || p.ret != ret {
		return
	}
	h.pending = nil
	r := o.Regs()
	if r.SS != p.ss || r.SP != p.sp+2 || o.Steps() < p.entry || o.Steps()-p.entry > heldSourceWindow {
		h.clearCandidate()
		return
	}
	driver, ok := heldDriver(o)
	if !ok {
		h.clearCandidate()
		return
	}
	firstY := p.y
	if p.line == 7 {
		firstY -= 6
	}
	rect := image.Rect(p.x-16, firstY-1, p.x-16+heldWidth, firstY-1+heldHeight)
	after := heldField(o, driver, rect)
	// spec037 E3: an empty second print must leave the complete source unchanged.
	if p.line == 7 && p.s == "" {
		n := h.name
		if n == nil || n.gen != h.gen || len(n.parts) != 1 || o.Steps() > n.deadline || !bytes.Equal(p.before, after) {
			h.clearCandidate()
		}
		return
	}
	mask, fg, ok := heldNameMask(p.before, after, p.line)
	if !ok {
		h.clearCandidate()
		return
	}
	p.after, p.mask, p.fg, p.returned = after, mask, fg, o.Steps()
	if p.line == 1 {
		if p.entry > ^uint64(0)-heldSourceWindow {
			return
		}
		h.name = &heldName{gen: h.gen, parts: []heldPart{*p}, rect: rect, deadline: p.entry + heldSourceWindow, after: after, mask: append([]byte(nil), mask...), fg: fg}
	} else {
		n := h.name
		if n == nil || n.gen != h.gen || fg != n.fg || o.Steps() > n.deadline {
			h.clearCandidate()
			return
		}
		n.parts = append(n.parts, *p)
		n.after = after
		for i, v := range mask {
			if v != 0 {
				n.mask[i] = 1
			}
		}
	}
}
func (h *heldItemCursor) captureName(o heldMachineReader) {
	if !heldAt(o.Regs(), 0x396f) {
		return
	}
	r := o.Regs()
	b := heldRead(o, oracle.Far(r.SS, r.SP), 10)
	if len(b) != 10 || heldWord(b, 0) != 0x5c1b {
		return
	}
	h.capture = nil
	n := h.name
	if r.CS != 0x0110 || n == nil || h.pending != nil || n.gen != h.gen || o.Steps() > n.deadline || !heldModule(o) {
		h.rejectedCapture++
		h.clearCandidate()
		return
	}
	rect := image.Rect(int(heldWord(b, 2)), int(heldWord(b, 4)), int(heldWord(b, 6))+1, int(heldWord(b, 8))+1)
	if rect != n.rect || r.DS != n.parts[0].ds || r.SS != n.parts[0].ss {
		h.rejectedCapture++
		h.clearCandidate()
		return
	}
	driver, ok := heldDriver(o)
	if !ok {
		h.clearCandidate()
		return
	}
	pic := heldField(o, driver, rect)
	if len(pic) != heldPixels {
		h.clearCandidate()
		return
	}
	for i, v := range n.mask {
		if v != 0 && pic[i] != n.after[i] {
			h.rejectedCapture++
			h.clearCandidate()
			return
		}
	}
	h.capture = &heldCapture{name: n, picture: pic, step: o.Steps()}
}
func heldGlyphs(e entry) ([][]uint16, bool) {
	if e.Dy != 0 || len(e.D) != 0 || len(e.G) == 0 || len(e.G) > 32 {
		return nil, false
	}
	size := len(e.G[0])
	if size != 15 && size != 16 {
		return nil, false
	}
	g := make([][]uint16, len(e.G))
	for i, rows := range e.G {
		if len(rows) != size {
			return nil, false
		}
		var ink uint16
		for _, row := range rows {
			if uint32(row) >= uint32(1)<<size {
				return nil, false
			}
			ink |= row
		}
		if ink == 0 {
			return nil, false
		}
		g[i] = append([]uint16(nil), rows...)
	}
	return g, true
}
func (h *heldItemCursor) enterDefinition(o heldMachineReader) {
	if !heldAt(o.Regs(), 0x40a7) {
		return
	}
	h.label = nil
	h.definition = nil
	r := o.Regs()
	b := heldRead(o, oracle.Far(r.SS, r.SP), 14)
	c := h.capture
	h.capture = nil
	h.pending = nil
	h.name = nil
	if len(b) != 14 || r.CS != 0x0110 || heldWord(b, 0) != 0x5cb9 || heldWord(b, 6) != 76 || heldWord(b, 8) != 13 || r.SP < 2 ||
		c == nil || c.name.gen != h.gen || o.Steps() > c.name.deadline || !heldModule(o) {
		h.rejectedDefinition++
		return
	}
	globals := heldRead(o, oracle.Far(r.DS, 0x0e1e), 4)
	ptr := heldPtr(b, 2)
	if len(globals) != 4 || ptr != heldPtr(globals, 0) || r.DS != c.name.parts[0].ds || r.SS != c.name.parts[0].ss || heldZero(ptr) {
		h.rejectedDefinition++
		return
	}
	pic := heldRead(o, ptr, heldPixels)
	if !bytes.Equal(pic, c.picture) {
		h.rejectedDefinition++
		return
	}
	key := c.name.key()
	normalized, _ := normalizeKey(key)
	g, ok := heldGlyphs(h.catalog[normalized])
	if !ok {
		log.Printf("物品游標缺譯：%s 字模不可用", key)
		return
	}
	label := &heldLabel{key: key, ptr: ptr, picture: pic, mask: append([]byte(nil), c.name.mask...), glyphs: g, fg: c.name.fg, source: c.name, captured: c.step}
	h.definition = &heldDefinition{label: label, ss: r.SS, sp: r.SP, step: o.Steps(), deadline: c.name.deadline}
}
func (h *heldItemCursor) returnDefinition(o heldMachineReader) {
	if !heldAt(o.Regs(), 0x40b0) {
		return
	}
	p := h.definition
	h.definition = nil
	if p == nil {
		return
	}
	r := o.Regs()
	if r.SS != p.ss || r.SP != p.sp-2 || o.Steps() < p.step || o.Steps() > p.deadline {
		h.rejectedDefinition++
		return
	}
	p.label.defined = o.Steps()
	h.label = p.label
}

type heldCursorView struct {
	x, y, w, h          int
	picture, background []byte
	palette             [256][3]uint8
	label               *heldLabel
}

func heldClipRect(x, y int) image.Rectangle {
	if x < 0 || y < 0 || x >= 320 || y >= 200 {
		return image.Rectangle{}
	}
	return image.Rect(x, y, min(x+heldWidth, 320), min(y+heldHeight, 200))
}
func (h *heldItemCursor) view(o heldReader, a []byte) (*heldCursorView, []image.Rectangle) {
	if h == nil {
		return nil, nil
	}
	if o.VideoMode() != 0x13 || len(a) != 64000 {
		h.label = nil
		h.lastRect = image.Rectangle{}
		return nil, nil
	}
	if !heldModule(o) {
		h.label = nil
		h.clearCandidate()
		h.definition = nil
	}
	driver, ok := heldDriver(o)
	if !ok {
		h.label = nil
		h.lastRect = image.Rectangle{}
		return nil, nil
	}
	desc := heldRead(o, oracle.Far(driver, 0x15ac), 24)
	full := heldRead(o, oracle.Far(driver, 0x1e5c), 12)
	busy := heldRead(o, oracle.Far(driver, 0x1ed3), 1)
	if len(desc) != 24 || len(full) != 12 || len(busy) != 1 {
		h.label = nil
		return nil, nil
	}
	ptr := heldPtr(full, 8)
	if heldWord(full, 0) != 76 || heldWord(full, 2) != 13 || heldZero(ptr) {
		h.label = nil
		h.lastRect = image.Rectangle{}
		return nil, nil
	}
	if h.label != nil && h.label.ptr != ptr {
		h.label = nil
	}
	x, y := int(int16(heldWord(desc, 0))), int(int16(heldWord(desc, 2)))
	current := heldClipRect(x, y)
	uncertain := []image.Rectangle{h.lastRect, current, heldClipRect(int(int16(heldWord(desc, 4))), int(int16(heldWord(desc, 6))))}
	if busy[0] != 0 || current.Empty() || int(heldWord(desc, 12)) != current.Dx() || int(heldWord(desc, 14)) != current.Dy() {
		return nil, uncertain
	}
	pic := heldRead(o, ptr, heldPixels)
	if len(pic) == heldPixels && h.label != nil && !bytes.Equal(h.label.picture, pic) {
		h.label = nil
	}
	bg := heldRead(o, oracle.Far(driver, heldWord(desc, 8)), current.Dx()*current.Dy())
	if len(pic) != heldPixels || len(bg) != current.Dx()*current.Dy() {
		return nil, uncertain
	}
	for cy := 0; cy < current.Dy(); cy++ {
		for cx := 0; cx < current.Dx(); cx++ {
			v := pic[cy*heldWidth+cx]
			if v == 255 {
				v = bg[cy*current.Dx()+cx]
			}
			if a[(y+cy)*320+x+cx] != v {
				return nil, uncertain
			}
		}
	}
	if h.label != nil && (h.label.ptr != ptr || !bytes.Equal(h.label.picture, pic)) {
		h.label = nil
	}
	h.lastRect = current
	return &heldCursorView{x: x, y: y, w: current.Dx(), h: current.Dy(), picture: pic, background: bg, palette: o.Palette(), label: h.label}, nil
}
func (v *heldCursorView) base(a []byte) []byte {
	if v == nil {
		return a
	}
	b := append([]byte(nil), a...)
	for y := 0; y < v.h; y++ {
		copy(b[(v.y+y)*320+v.x:], v.background[y*v.w:(y+1)*v.w])
	}
	return b
}
func heldIntersects(e *textEv, rects []image.Rectangle) bool {
	r := image.Rect(e.bx, e.by, e.bx+e.bw, e.by+e.bh)
	for _, box := range rects {
		if !box.Empty() && r.Overlaps(box) {
			return true
		}
	}
	return false
}
func (v *heldCursorView) translated(k int) bool {
	if v == nil || v.label == nil || (k != 2 && k != 3) {
		return false
	}
	g := v.label.glyphs
	size := len(g[0])
	width := len(g)*(size+2) - 2
	if width > 60*k || size+1*k > 13*k {
		return false
	}
	for gi, glyph := range g {
		for y, bits := range glyph {
			for x := 0; x < size; x++ {
				if bits&(1<<(size-1-x)) == 0 {
					continue
				}
				px, py := 16*k+gi*(size+2)+x, k+y
				i := (py/k)*heldWidth + px/k
				if i < 0 || i >= heldPixels || v.picture[i] != 255 && v.label.mask[i] == 0 {
					return false
				}
			}
		}
	}
	return true
}
func (v *heldCursorView) raster(k int) *image.RGBA {
	if v == nil || k < 1 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, heldWidth*k, heldHeight*k))
	zh := v.translated(k)
	for y := 0; y < heldHeight*k; y++ {
		for x := 0; x < heldWidth*k; x++ {
			i := (y/k)*heldWidth + x/k
			value := v.picture[i]
			if value == 255 || zh && v.label.mask[i] != 0 {
				continue
			}
			c := v.palette[value]
			img.SetRGBA(x, y, color.RGBA{c[0], c[1], c[2], 255})
		}
	}
	if zh {
		c := v.palette[v.label.fg]
		fg := color.RGBA{c[0], c[1], c[2], 255}
		for gi, glyph := range v.label.glyphs {
			size := len(glyph)
			for y, bits := range glyph {
				for x := 0; x < size; x++ {
					if bits&(1<<(size-1-x)) != 0 {
						img.SetRGBA(16*k+gi*(size+2)+x, k+y, fg)
					}
				}
			}
		}
	}
	return img
}
func (v *heldCursorView) paint(dst *image.RGBA, k int) {
	if v == nil {
		return
	}
	src := v.raster(k)
	if src == nil {
		return
	}
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			c := src.RGBAAt(x, y)
			if c.A != 0 {
				dst.SetRGBA(v.x*k+x, v.y*k+y, c)
			}
		}
	}
}
func (g *game) drawHeldCursor(dst *ebiten.Image, v *heldCursorView) {
	if v == nil || dst.Bounds() != image.Rect(0, 0, 320*g.scale, 200*g.scale) {
		return
	}
	img := v.raster(g.scale)
	if img == nil {
		return
	}
	h := sha256.Sum256(img.Pix)
	key := hex.EncodeToString(h[:])
	if g.heldImage == nil || g.heldImageKey != key {
		if g.heldImage != nil {
			g.heldImage.Dispose()
		}
		g.heldImage = ebiten.NewImageFromImage(img)
		g.heldImageKey = key
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(v.x*g.scale), float64(v.y*g.scale))
	dst.DrawImage(g.heldImage, &op)
}

func heldEventStill(e *textEv, off []byte) bool {
	still := e.rectMatches(off, nil)
	if e.profile != nil {
		still = e.profile.matches(off, nil)
	}
	if e.roster != nil {
		still = e.rosterMatches(off, nil)
	}
	if e.party != nil {
		still = still && e.partyNameMatches(off, nil)
	}
	return still
}
