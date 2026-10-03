package main

// spec 021：既有覆繪事件完整辨識新增角色對話框後，只移中文畫筆。
// 原版按鈕為(63,166,73,9)及(63,182,73,9)，事件及命中區不變。
var characterDialogLabels = [...]struct {
	s    string
	x, y int
}{
	{"ADD CHARACTER", 61, 152},
	{"OK", 94, 168},
	{"CANCEL", 82, 184},
}

func buttonLabelY(e *textEv, shown []*textEv, scale int) int {
	oldY := menuLabelY(e, shown, scale)
	if y, ok := characterMenuLabelY(e, shown, scale); ok {
		return y
	}
	if scale != 2 && scale != 3 || len(e.g) == 0 {
		return oldY
	}
	isButton := false
	for _, label := range characterDialogLabels[1:] {
		if mainMenuLabelMatch(e, label.s, label.x, label.y) {
			isButton = true
		}
	}
	if !isButton {
		return oldY
	}
	for _, labels := range [][]struct {
		s    string
		x, y int
	}{mainMenuLabels[:], characterDialogLabels[:]} {
		for _, label := range labels {
			found := false
			for _, other := range shown {
				if mainMenuLabelMatch(other, label.s, label.x, label.y) {
					found = true
					break
				}
			}
			if !found {
				return oldY
			}
		}
	}
	height := len(e.g[0])
	if height != 15 && height != 16 {
		return oldY
	}
	width := -2
	for _, glyph := range e.g {
		if len(glyph) != height {
			return oldY
		}
		width += height + 2
	}
	if e.x0*scale < 63*scale || e.x0*scale+width > 136*scale {
		return oldY
	}
	return (e.y0-2)*scale + (9*scale-height+1)/2
}
