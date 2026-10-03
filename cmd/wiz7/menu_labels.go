package main

// spec 020：只在六個已量主選單標籤同時上屏時，調整畫筆的目的y。
// 座標與按鈕邊界出自原版英文色號；事件及原版資料保持原樣。
var mainMenuLabels = [...]struct {
	s    string
	x, y int
}{
	{"CHARACTER MENU", 118, 45},
	{"ADD CHARACTER", 121, 58},
	{"LOAD SAVEGAME", 121, 71},
	{"CONFIGURATION", 121, 84},
	{"IMPORT (BANE)", 121, 97},
	{"QUIT GAME", 133, 110},
}

func mainMenuLabelMatch(e *textEv, s string, x, y int) bool {
	return e != nil && e.s == s && e.x0 == x && e.y0 == y && e.dy == 0
}

func menuLabelY(e *textEv, shown []*textEv, scale int) int {
	oldY := (e.y0 + e.dy) * scale
	if scale != 2 && scale != 3 || len(e.g) == 0 {
		return oldY
	}
	height := len(e.g[0])
	if height != 15 && height != 16 {
		return oldY
	}
	width := -2 // 末字之後的字距不屬於字格。
	for _, glyph := range e.g {
		if len(glyph) != height {
			return oldY
		}
		width += height + 2
	}
	isLabel := false
	for _, label := range mainMenuLabels {
		if mainMenuLabelMatch(e, label.s, label.x, label.y) {
			isLabel = true
		}
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
	if !isLabel || e.x0*scale < 105*scale || e.x0*scale+width > 214*scale {
		return oldY
	}
	return (e.y0-2)*scale + (9*scale-height+1)/2
}
