// 032：只保留已驗證裝備文字中未被原版法術視窗遮住的部分。
package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

var profilePopupRect = image.Rect(49, 114, 271, 198)
var profilePopupBorder = [...][3]int{
	{49, 114, 0}, {50, 114, 0}, {51, 114, 0}, {147, 114, 0}, {148, 114, 0}, {149, 114, 0}, {269, 114, 0}, {270, 114, 0},
	{49, 115, 0}, {50, 115, 99}, {49, 150, 0}, {50, 150, 102}, {49, 197, 0}, {50, 197, 113}, {147, 197, 113}, {148, 197, 0},
	{269, 197, 197}, {270, 197, 0}, {49, 198, 0}, {50, 198, 0}, {149, 115, 180}, {150, 115, 180}, {269, 196, 197}, {270, 196, 0},
}

func profilePopupActive(evs []*textEv, a []byte, cursor func(int, int) bool) bool {
	if len(a) != 64000 || (flowK != 2 && flowK != 3) {
		return false
	}
	found := [2]bool{}
	for _, e := range evs {
		if e.dead || e.originRet != oracle.Far(0x0110, 0x0ea2) || e.originFont != 0 || e.fg != 85 || e.dy != 0 || e.y0 != 123 {
			continue
		}
		slot := -1
		if e.s == "CAST" && e.x0 == 168 && e.bx == 168 && e.by == 123 && e.bw == 23 && e.bh == 5 {
			slot = 0
		}
		if e.s == "CANCEL" && e.x0 == 222 && e.bx == 222 && e.by == 123 && e.bw == 35 && e.bh == 5 {
			slot = 1
		}
		if slot >= 0 && e.rectMatches(a, cursor) {
			found[slot] = true
		}
	}
	if !found[0] || !found[1] {
		return false
	}
	checked := 0
	for _, p := range profilePopupBorder {
		if cursor != nil && cursor(p[0], p[1]) {
			continue
		}
		checked++
		if a[p[1]*320+p[0]] != byte(p[2]) {
			return false
		}
	}
	return checked >= 19
}
func profilePopupEquipment(e *textEv) bool {
	if e.dead || !e.shown || e.x0 != 18 || e.dy != 0 || e.fg != 105 || e.profile != nil || e.page != nil || e.roster != nil || e.sex {
		return false
	}
	if e.bx < 0 || e.by < 0 || e.bw <= 0 || e.bh <= 0 || e.bx+e.bw > 320 || e.by+e.bh > 200 || len(e.before) != e.bw*e.bh || len(e.after) != e.bw*e.bh {
		return false
	}
	if len(e.g) == 0 {
		return false
	}
	size := len(e.g[0])
	if size != 15 && size != 16 {
		return false
	}
	for _, g := range e.g {
		if len(g) != size {
			return false
		}
		nonempty := false
		for _, row := range g {
			if size == 15 && row >= 1<<15 {
				return false
			}
		}
		for _, row := range g {
			nonempty = nonempty || row != 0
		}
		if !nonempty {
			return false
		}
	}
	k := flowK
	if (k != 2 && k != 3) || e.y0*k < profilePopupRect.Min.Y*k || e.y0*k+size > profilePopupRect.Max.Y*k || e.x0*k+len(e.g)*(size+2)-2 > profilePopupRect.Max.X*k {
		return false
	}
	if e.compound != nil {
		c := e.compound
		if c.Inputs[0].Ret != oracle.Far(0x0110, 0xba40) || c.Inputs[1].Ret != oracle.Far(0x0110, 0xba63) {
			return false
		}
		if c.Inputs[0].Font != 0 || c.Inputs[1].Font != 0 || c.Inputs[0].X != 18 || c.Inputs[1].X != 18 {
			return false
		}
		for _, s := range compoundSites {
			if len(e.g) == s.glyphs && e.y0 == s.y && e.bx == 18 && e.by == s.y && e.bw == 41 && e.bh == 11 && c.Inputs[0].Y == s.y && c.Inputs[1].Y == s.y+6 && c.Inputs[0].S == s.first && c.Inputs[1].S == s.second {
				return true
			}
		}
		return false
	}
	return e.originRet == oracle.Far(0x0110, 0xba40) && e.originFont == 0 && len(e.g) == 2 && e.bx == 18 && e.by == e.y0 && e.bh == 5 && ((e.s == "LONGSWORD" && e.y0 == 114 && e.bw == 53) || (e.s == "SANDALS" && e.y0 == 156 && e.bw == 41))
}
func (e *textEv) popupContains(x, y int) bool {
	return !e.popupClip.Empty() && image.Pt(x, y).In(e.popupClip)
}
func (e *textEv) popupContainsHD(x, y, k int) bool {
	r := e.popupClip
	return !r.Empty() && x >= r.Min.X*k && x < r.Max.X*k && y >= r.Min.Y*k && y < r.Max.Y*k
}
func (e *textEv) drawPopupClipped(dst, strip *ebiten.Image, op *ebiten.DrawImageOptions, k int) {
	if e.popupClip.Empty() {
		dst.DrawImage(strip, op)
		return
	}
	bounds := dst.Bounds()
	if (k != 2 && k != 3) || bounds.Min != (image.Point{}) || bounds.Dx() != 320*k || bounds.Dy() < 200*k || len(e.g) == 0 || len(e.g[0]) == 0 {
		return
	}
	if e.y0*k < e.popupClip.Min.Y*k || e.y0*k+len(e.g[0]) > e.popupClip.Max.Y*k {
		return
	}
	left := dst.SubImage(image.Rect(0, 0, e.popupClip.Min.X*k, 200*k)).(*ebiten.Image)
	left.DrawImage(strip, op)
}
