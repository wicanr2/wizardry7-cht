// spec033：WHO游標來源圖唯讀辨識與原生中文字形。
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"

	"github.com/wicanr2/dosgolem/oracle"
)

const whoTailSHA = "1e34be319abc13e44f2b8b3eeb998bed1dfeb53c37e8ae369cb415546671aa61"
const whoPrefixSHA = "575e226442eb3e8e3d52d5374651cd9c66e202197af0696a7cae32ac2c7d8a6a"

type whoCursorView struct {
	x, y, glyphX, glyphY, advance, k int
	glyphs                           [][]uint16
	tail                             []byte
	fg                               color.RGBA
	bg                               color.RGBA
}

func whoDigest(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func whoGlyphsValid(e entry) bool {
	if e.Dy != 0 || len(e.D) != 0 || len(e.G) != 2 {
		return false
	}
	h := len(e.G[0])
	if (h != 15 && h != 16) || len(e.G[1]) != h {
		return false
	}
	for _, g := range e.G {
		ink := false
		for _, r := range g {
			if uint32(r) >= uint32(1)<<h {
				return false
			}
			ink = ink || r != 0
		}
		if !ink {
			return false
		}
	}
	return true
}
func whoInkColumns(g []uint16) (int, int) {
	lo, hi := len(g), -1
	for _, r := range g {
		for x := 0; x < len(g); x++ {
			if r&(1<<(len(g)-1-x)) != 0 {
				lo = min(lo, x)
				hi = max(hi, x)
			}
		}
	}
	return lo, hi
}
func whoLayout(e entry, prefix []byte, k int) (int, int, int, bool) {
	if (k != 2 && k != 3) || !whoGlyphsValid(e) || len(prefix) != 48 {
		return 0, 0, 0, false
	}
	h := len(e.G[0])
	_, right := whoInkColumns(e.G[0])
	left, _ := whoInkColumns(e.G[1])
	advance := right + 2 - left
	width := advance + h
	if advance <= 0 || width < h || width > 24*k {
		return 0, 0, 0, false
	}
	cy := 14*k - max(h-6*k, 0)
	if cy < 12*k || cy+h > 20*k {
		return 0, 0, 0, false
	}
	preferred := (24*k - width) / 2
	chosen := -1
	for cx := 0; cx+width <= 24*k; cx++ {
		safe := true
		for gi, g := range e.G {
			for row, bits := range g {
				for col := 0; col < h; col++ {
					if bits&(1<<(h-1-col)) == 0 {
						continue
					}
					px, py := cx+gi*advance+col, cy+row
					if py < 13*k && prefix[(py/k-12)*24+px/k] != 255 {
						safe = false
					}
				}
			}
		}
		distance := cx - preferred
		if distance < 0 {
			distance = -distance
		}
		previous := chosen - preferred
		if previous < 0 {
			previous = -previous
		}
		if safe && (chosen < 0 || distance < previous) {
			chosen = cx
		}
	}
	return chosen, cy, advance, chosen >= 0
}

type whoCursorReader interface {
	VideoMode() uint8
	Word(oracle.Addr) uint16
	Bytes(oracle.Addr, int) []byte
	Palette() [256][3]uint8
}

func whoCursorFromSource(o whoCursorReader, e entry, k int, a []byte) *whoCursorView {
	if o.VideoMode() != 0x13 || len(a) != 64000 || !whoGlyphsValid(e) || (k != 2 && k != 3) {
		return nil
	}
	w := func(off uint16) uint16 { return o.Word(oracle.Far(0x740c, off)) }
	if w(0x1e5c) != 24 || w(0x1e5e) != 20 || w(0x1e60) != 10 || w(0x1e62) != 10 || w(0x15be) != 24 || w(0x15c0) != 20 {
		return nil
	}
	x, y := int(int16(w(0x15ac))), int(int16(w(0x15ae)))
	if x < 0 || y < 0 || x+24 > 320 || y+20 > 200 {
		return nil
	}
	off, seg := w(0x1e64), w(0x1e66)
	if int(off)+480 > 65536 || uint32(seg)*16+uint32(off)+480 > 0xa0000 {
		return nil
	}
	tail := o.Bytes(oracle.Far(seg, off+336), 144)
	prefix := o.Bytes(oracle.Far(seg, off+288), 48)
	if len(tail) != 144 || len(prefix) != 48 || whoDigest(tail) != whoTailSHA || whoDigest(prefix) != whoPrefixSHA {
		return nil
	}
	for i, c := range tail {
		if c != 255 && a[(y+14+i/24)*320+x+i%24] != c {
			return nil
		}
	}
	gx, gy, advance, ok := whoLayout(e, prefix, k)
	if !ok {
		log.Print("WHO缺譯：原生字格無安全位置")
		return nil
	}
	g := [][]uint16{append([]uint16(nil), e.G[0]...), append([]uint16(nil), e.G[1]...)}
	c := o.Palette()[22]
	return &whoCursorView{x: x, y: y, glyphX: gx, glyphY: gy, advance: advance, k: k, glyphs: g, tail: append([]byte(nil), tail...), fg: color.RGBA{c[0], c[1], c[2], 255}, bg: color.RGBA{o.Palette()[0][0], o.Palette()[0][1], o.Palette()[0][2], 255}}
}
func (v *whoCursorView) restore(comp []byte) {
	if v == nil || len(comp) != 64000 {
		return
	}
	for i, c := range v.tail {
		if c == 22 {
			comp[(v.y+14+i/24)*320+v.x+i%24] = 0
		}
	}
}
func (v *whoCursorView) paint(dst *image.RGBA, k int) {
	if v == nil || k != v.k {
		return
	}
	for gi, g := range v.glyphs {
		for y, bits := range g {
			for x := 0; x < len(g); x++ {
				if bits&(1<<(len(g)-1-x)) != 0 {
					dst.SetRGBA(v.x*k+v.glyphX+gi*v.advance+x, v.y*k+v.glyphY+y, v.fg)
				}
			}
		}
	}
}

func (g *game) drawWhoCursor(dst *ebiten.Image, v *whoCursorView) {
	if !v.canDraw(dst.Bounds(), g.scale) {
		return
	}
	raw := []byte{v.fg.R, v.fg.G, v.fg.B, byte(v.advance)}
	for _, gl := range v.glyphs {
		for _, row := range gl {
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], row)
			raw = append(raw, b[:]...)
		}
	}
	key := whoDigest(raw)
	if g.whoImage == nil || g.whoImageKey != key {
		h := len(v.glyphs[0])
		img := image.NewRGBA(image.Rect(0, 0, v.advance+h, h))
		for gi, gl := range v.glyphs {
			for y, bits := range gl {
				for x := 0; x < h; x++ {
					if bits&(1<<(h-1-x)) != 0 {
						img.SetRGBA(gi*v.advance+x, y, v.fg)
					}
				}
			}
		}
		if g.whoImage != nil {
			g.whoImage.Dispose()
		}
		g.whoImage = ebiten.NewImageFromImage(img)
		g.whoImageKey = key
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(v.x*v.k+v.glyphX), float64(v.y*v.k+v.glyphY))
	dst.DrawImage(g.whoImage, &op)
}

func takeWhoCursorEntry(trans map[string]entry) entry {
	e := trans["WHO?"]
	if len(trans) > 0 && !whoGlyphsValid(e) {
		log.Print("WHO缺譯：游標字模不可用，保留原圖")
	}
	delete(trans, "WHO?")
	return e
}

func (v *whoCursorView) canDraw(bounds image.Rectangle, k int) bool {
	return v != nil && k == v.k && (k == 2 || k == 3) && bounds == image.Rect(0, 0, 320*k, 200*k)
}
func (v *whoCursorView) patchRGBA(rgba []byte) {
	if v == nil || len(rgba) != 320*200*4 {
		return
	}
	for i, c := range v.tail {
		if c != 22 {
			continue
		}
		p := ((v.y+14+i/24)*320 + v.x + i%24) * 4
		rgba[p], rgba[p+1], rgba[p+2], rgba[p+3] = v.bg.R, v.bg.G, v.bg.B, 255
	}
}
