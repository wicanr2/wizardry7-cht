package main

// 024：完整檢視選人三訊號才調兩個按鈕畫筆；父選單已被子視窗覆蓋。
var reviewDialogLabels = [...]struct {
	s    string
	x, y int
}{
	{"REVIEW WHO?", 67, 152},
	{"OK", 94, 168},
	{"CANCEL", 82, 184},
}

func reviewDialogLabelY(e *textEv, shown []*textEv, scale int) (int, bool) {
	if e == nil || scale != 2 && scale != 3 {
		return 0, false
	}
	button := false
	for _, label := range reviewDialogLabels[1:] {
		button = button || mainMenuLabelMatch(e, label.s, label.x, label.y)
	}
	if !button {
		return 0, false
	}
	height := 0
	for _, label := range reviewDialogLabels {
		var matched *textEv
		for _, other := range shown {
			if mainMenuLabelMatch(other, label.s, label.x, label.y) {
				matched = other
				break
			}
		}
		if matched == nil || len(matched.g) == 0 {
			return 0, false
		}
		for _, glyph := range matched.g {
			h := len(glyph)
			if h != 15 && h != 16 || height != 0 && h != height {
				return 0, false
			}
			height = h
			nonempty := false
			for _, row := range glyph {
				if uint32(row) >= uint32(1)<<height {
					return 0, false
				}
				nonempty = nonempty || row != 0
			}
			if !nonempty {
				return 0, false
			}
		}
	}
	if len(e.g) == 0 {
		return 0, false
	}
	width := len(e.g)*(height+2) - 2
	for _, glyph := range e.g {
		if len(glyph) != height {
			return 0, false
		}
		for _, row := range glyph {
			if uint32(row) >= uint32(1)<<height {
				return 0, false
			}
		}
	}
	if e.x0*scale < 63*scale || e.x0*scale+width > 136*scale {
		return 0, false
	}
	return (e.y0-2)*scale + (9*scale-height+1)/2, true
}
