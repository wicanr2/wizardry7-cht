package main

// spec 023：完整角色選單才調畫筆y，原生字模與原版事件保持。
var characterMenuLabels = [...]struct {
	s       string
	x, left int
}{
	{"CREATE", 11, 4},
	{"REVIEW", 63, 56},
	{"DELETE", 115, 108},
	{"RENAME", 167, 160},
	{"PICTURE", 216, 212},
	{"EXIT", 277, 264},
}

func characterMenuLabelY(e *textEv, shown []*textEv, scale int) (int, bool) {
	if scale != 2 && scale != 3 || e == nil || len(e.g) == 0 {
		return 0, false
	}
	height := len(e.g[0])
	if height != 15 && height != 16 {
		return 0, false
	}
	fits := func(event *textEv, left int) bool {
		if len(event.g) == 0 {
			return false
		}
		width := len(event.g)*(height+2) - 2
		if event.x0*scale < left*scale || event.x0*scale+width > (left+49)*scale {
			return false
		}
		for _, glyph := range event.g {
			if len(glyph) != height {
				return false
			}
			for _, row := range glyph {
				if uint32(row)>>uint(height) != 0 {
					return false
				}
			}
		}
		return true
	}
	isTarget := false
	for _, label := range characterMenuLabels {
		if mainMenuLabelMatch(e, label.s, label.x, 190) {
			if !fits(e, label.left) {
				return 0, false
			}
			isTarget = true
		}
		found := false
		for _, other := range shown {
			if mainMenuLabelMatch(other, label.s, label.x, 190) {
				if !fits(other, label.left) {
					return 0, false
				}
				found = true
			}
		}
		if !found {
			return 0, false
		}
	}
	if !isTarget {
		return 0, false
	}
	return 188*scale + (9*scale-height+1)/2, true
}
