package main

// spec026：R959已量七個控制同時上屏，只調HD／Draw中文畫筆Y。
// 邊界採原版39×9外框，原事件、字模與按鈕命中區保持。
var battleControlLabels = [...]struct {
	s               string
	x, y, left, top int
}{
	{"FIGHT", 13, 167, 8, 165},
	{"PARRY", 13, 177, 8, 175},
	{"SPELL", 55, 167, 50, 165},
	{"USE", 61, 177, 50, 175},
	{"RUN", 241, 167, 230, 165},
	{"MOVE", 238, 177, 230, 175},
	{"EQUIP", 235, 187, 230, 185},
}

func battleControlLabelY(e *textEv, shown []*textEv, scale int) (int, bool) {
	if e == nil || scale != 2 && scale != 3 {
		return 0, false
	}
	height, top, target := 0, 0, false
	for _, label := range battleControlLabels {
		var found *textEv
		for _, other := range shown {
			if mainMenuLabelMatch(other, label.s, label.x, label.y) {
				found = other
				break
			}
		}
		if found == nil || found.fg != 1 || len(found.g) != 2 {
			return 0, false
		}
		for _, glyph := range found.g {
			h := len(glyph)
			if h != 15 && h != 16 || height != 0 && height != h {
				return 0, false
			}
			height = h
			nonempty := false
			for _, row := range glyph {
				if uint32(row) >= uint32(1)<<h {
					return 0, false
				}
				nonempty = nonempty || row != 0
			}
			if !nonempty {
				return 0, false
			}
		}
		width := len(found.g)*(height+2) - 2
		if found.x0*scale < label.left*scale || found.x0*scale+width > (label.left+39)*scale {
			return 0, false
		}
		if e == found {
			top, target = label.top, true
		}
	}
	if !target {
		return 0, false
	}
	return top*scale + (9*scale-height+1)/2, true
}
