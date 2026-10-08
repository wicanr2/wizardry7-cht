package main

// spec029：只讀本次已選定的防拷記錄，提示不參與原版驗證。
import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

const clearanceLabel = "答案："
const clearanceLastKey = "LOCATED IN PARAGRAPH $  LINE $"

type clearancePending struct {
	id, ds, dst, ss, sp uint16
	step                uint64
}

type clearanceSelection struct {
	prefix  [4]byte
	answer  string
	matched bool
}

type clearanceHint struct {
	pending  *clearancePending
	selected *clearanceSelection
	font     *helpData
	height   int
	reasons  map[string]bool
	views    [4]*clearanceHintView
}

type clearanceHintView struct {
	pixels    *image.RGBA
	x, y      int
	scale     int
	color     [3]uint8
	lineWidth int
}

func newClearanceHint(trans map[string]entry, h *helpData) *clearanceHint {
	c := &clearanceHint{reasons: map[string]bool{}}
	e, ok := trans[clearanceLastKey]
	if !ok || len(e.G) == 0 || h == nil {
		c.fail("font_unavailable")
		return c
	}
	c.height = len(e.G[0])
	if c.height != 15 && c.height != 16 {
		c.fail("font_shape")
		return c
	}
	for _, g := range e.G {
		if len(g) != c.height {
			c.fail("font_shape")
			return c
		}
	}
	font := &helpData{Glyphs: map[string]helpGlyph{}}
	for _, r := range clearanceLabel + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		g, ok := h.Glyphs[string(r)]
		w := c.height
		if r >= 'A' && r <= 'Z' {
			w = 8
		}
		if !ok || g.W != w || len(g.Rows) != c.height {
			c.fail("font_shape")
			return c
		}
		nonempty := false
		for _, row := range g.Rows {
			if row < 0 || row >= 1<<w {
				c.fail("font_shape")
				return c
			}
			nonempty = nonempty || row != 0
		}
		if !nonempty {
			c.fail("font_shape")
			return c
		}
		font.Glyphs[string(r)] = helpGlyph{g.W, append([]int(nil), g.Rows...)}
	}
	c.font = font
	return c
}

func (c *clearanceHint) fail(reason string) {
	if c == nil {
		return
	}
	if !c.reasons[reason] {
		c.reasons[reason] = true
		log.Printf("防拷提示缺譯：%s", reason)
	}
}

func (c *clearanceHint) begin(o *oracle.Oracle) {
	if c == nil || c.font == nil {
		return
	}
	r := o.Regs()
	if !clearanceEntrySP(r.SP) {
		c.pending, c.selected = nil, nil
		c.views = [4]*clearanceHintView{}
		c.fail("acquisition_stack")
		return
	}
	if c.pending != nil {
		c.selected = nil
		c.views = [4]*clearanceHintView{}
		c.fail("acquisition_reentry")
		c.pending = nil
		return
	}
	c.pending = nil
	id := o.StackWord(1)
	if id < 2500 || id > 2574 {
		return
	}
	c.selected = nil
	c.views = [4]*clearanceHintView{}
	c.pending = &clearancePending{id, r.DS, o.StackWord(2), r.SS, r.SP, o.Steps()}
}

func clearanceEntrySP(sp uint16) bool { return sp >= 0x112 && sp <= 65530 }

func (p *clearancePending) paired(ds, ss, sp, bp, id, dst uint16, step uint64) bool {
	return p != nil && clearanceEntrySP(p.sp) &&
		ds == p.ds && ss == p.ss && sp == p.sp-0x112 && bp == p.sp-2 &&
		id == p.id && dst == p.dst && step >= p.step && step-p.step <= 1000000 && int(p.dst)+64 <= 65536
}

func decodeClearanceRecord(raw []byte) *clearanceSelection {
	end := -1
	for i, b := range raw {
		if b == 0 {
			end = i
			break
		}
	}
	if end < 5 || end > 36 {
		return nil
	}
	for _, b := range raw[:4] {
		if b == 0 {
			return nil
		}
	}
	for _, b := range raw[4:end] {
		if b < 'A' || b > 'Z' {
			return nil
		}
	}
	s := &clearanceSelection{answer: string(raw[4:end])}
	copy(s.prefix[:], raw[:4])
	return s
}

func (c *clearanceHint) finish(o *oracle.Oracle) {
	if c == nil || c.pending == nil {
		return
	}
	p := c.pending
	c.pending = nil
	r := o.Regs()
	if r.BP > 65528 || !p.paired(r.DS, r.SS, r.SP, r.BP,
		o.Word(oracle.Far(r.SS, r.BP+4)), o.Word(oracle.Far(r.SS, r.BP+6)), o.Steps()) {
		c.selected = nil
		c.fail("acquisition_pair")
		return
	}
	c.selected = decodeClearanceRecord(o.Bytes(oracle.Far(p.ds, p.dst), 64))
	if c.selected == nil {
		c.fail("record_shape")
	}
}

var clearancePagePattern = regexp.MustCompile(`^PLEASE TURN TO PAGE ([1-9][0-9]*) OF YOUR$`)
var clearanceWordPattern = regexp.MustCompile(`^MANUAL AND ENTER THE (FIRST|SECOND|THIRD|FOURTH|FIFTH|SIXTH|SEVENTH|EIGHTH|NINTH|TENTH) WORD$`)
var clearanceLinePattern = regexp.MustCompile(`^LOCATED IN PARAGRAPH ([1-9][0-9]*)  LINE ([1-9][0-9]*)$`)
var clearanceOrdinals = map[string]int{"FIRST": 1, "SECOND": 2, "THIRD": 3, "FOURTH": 4, "FIFTH": 5, "SIXTH": 6, "SEVENTH": 7, "EIGHTH": 8, "NINTH": 9, "TENTH": 10}

func clearanceQuestion(shown []*textEv) ([4]byte, *textEv, bool) {
	var prefix [4]byte
	var last *textEv
	var found [4]bool
	positive := func(s string) byte {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 255 {
			return 0
		}
		return byte(n)
	}
	for _, e := range shown {
		if e == nil || e.dead || !e.shown || !e.match || e.dy != 0 {
			continue
		}
		switch {
		case e.x0 == 106 && e.y0 == 148 && e.s == "SECURITY CLEARANCE":
			found[0] = true
		case e.x0 == 70 && e.y0 == 158:
			if m := clearancePagePattern.FindStringSubmatch(e.s); m != nil {
				prefix[0] = positive(m[1])
				found[1] = prefix[0] != 0
			}
		case e.x0 == 64 && e.y0 == 164:
			if m := clearanceWordPattern.FindStringSubmatch(e.s); m != nil {
				prefix[3] = byte(clearanceOrdinals[m[1]])
				found[2] = true
			}
		case e.x0 == 70 && e.y0 == 170:
			if m := clearanceLinePattern.FindStringSubmatch(e.s); m != nil {
				prefix[1], prefix[2] = positive(m[1]), positive(m[2])
				found[3] = prefix[1] != 0 && prefix[2] != 0
				last = e
			}
		}
	}
	return prefix, last, found[0] && found[1] && found[2] && found[3]
}

func (c *clearanceHint) view(shown []*textEv, scale int) *clearanceHintView {
	if c == nil || c.font == nil || c.selected == nil {
		return nil
	}
	prefix, last, complete := clearanceQuestion(shown)
	if !complete || prefix != c.selected.prefix {
		if c.selected.matched {
			c.selected = nil
			c.views = [4]*clearanceHintView{}
		}
		return nil
	}
	c.selected.matched = true
	if scale != 2 && scale != 3 {
		c.fail("scale")
		return nil
	}
	lineWidth := -2
	for _, g := range last.g {
		if len(g) != c.height {
			c.fail("question_font")
			return nil
		}
		lineWidth += len(g) + 2
	}
	if lineWidth < 0 {
		c.fail("question_font")
		return nil
	}
	if v := c.views[scale]; v != nil && v.color == last.fgRGB && v.lineWidth == lineWidth {
		return v
	}
	text := clearanceLabel + c.selected.answer
	width, ok := c.font.textWidth(text)
	x, y := last.x0*scale+lineWidth+8, last.y0*scale
	area := image.Rect(x, y, x+width, y+c.height)
	if !ok || !area.In(image.Rect(70*scale, 170*scale, 263*scale, 178*scale)) {
		c.fail("space")
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, c.height))
	c.font.drawText(img, 0, 0, text, color.RGBA{last.fgRGB[0], last.fgRGB[1], last.fgRGB[2], 255})
	v := &clearanceHintView{pixels: img, x: x, y: y, scale: scale, color: last.fgRGB, lineWidth: lineWidth}
	c.views[scale] = v
	return v
}

func (v *clearanceHintView) paint(img *image.RGBA) {
	if v == nil {
		return
	}
	if img.Bounds() != image.Rect(0, 0, 320*v.scale, 200*v.scale) {
		return
	}
	r := v.pixels.Bounds().Add(image.Pt(v.x, v.y))
	if r.In(img.Bounds()) {
		draw.Draw(img, r, v.pixels, image.Point{}, draw.Over)
	}
}

func (g *game) drawClearanceHint(dst *ebiten.Image, v *clearanceHintView) {
	if v == nil {
		g.clearanceView = nil
		g.clearanceImage = nil
		return
	}
	if dst.Bounds() != image.Rect(0, 0, 320*v.scale, 200*v.scale) || !v.pixels.Bounds().Add(image.Pt(v.x, v.y)).In(dst.Bounds()) {
		return
	}
	if g.clearanceView != v {
		g.clearanceView = v
		g.clearanceImage = ebiten.NewImageFromImage(v.pixels)
	}
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(v.x), float64(v.y))
	dst.DrawImage(g.clearanceImage, &op)
}
