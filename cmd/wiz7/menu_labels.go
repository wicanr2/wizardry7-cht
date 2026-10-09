package main

// spec 020／028：只在已量主選單完整布局同時上屏時，調整畫筆的目的y。
// 座標與按鈕邊界出自原版英文色號；事件及原版資料保持原樣。
type menuLabel = struct {
	s    string
	x, y int
}

var mainMenuLabels = [...]menuLabel{
	{"CHARACTER MENU", 118, 45},
	{"ADD CHARACTER", 121, 58},
	{"LOAD SAVEGAME", 121, 71},
	{"CONFIGURATION", 121, 84},
	{"IMPORT (BANE)", 121, 97},
	{"QUIT GAME", 133, 110},
}

// 028 R995：保留原六項供既有helper使用，完整布局各自辨識。
var middleMenuLabels = [...]menuLabel{
	{"CHARACTER MENU", 118, 45},
	{"ADD CHARACTER", 121, 58},
	{"REVIEW MEMBER", 121, 71},
	{"DISMISS MEMBER", 118, 84},
	{"START NEW GAME", 118, 97},
	{"CONFIGURATION", 121, 110},
	{"QUIT GAME", 133, 123},
}

var fullMenuLabels = [...]menuLabel{
	{"CHARACTER MENU", 118, 45},
	{"REVIEW MEMBER", 121, 58},
	{"DISMISS MEMBER", 118, 71},
	{"START NEW GAME", 118, 84},
	{"CONFIGURATION", 121, 97},
	{"QUIT GAME", 133, 110},
}

func activeMainMenuLayout(shown []*textEv) []menuLabel {
	var selected []menuLabel
	for _, layout := range [][]menuLabel{mainMenuLabels[:], middleMenuLabels[:], fullMenuLabels[:]} {
		complete := true
		for _, label := range layout {
			found := false
			for _, e := range shown {
				if e != nil && !e.dead && e.shown && e.match && mainMenuLabelMatch(e, label.s, label.x, label.y) {
					found = true
					break
				}
			}
			if !found {
				complete = false
				break
			}
		}
		if complete {
			if selected != nil {
				return nil
			}
			selected = layout
		}
	}
	return selected
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
	layout := activeMainMenuLayout(shown)
	if layout == nil {
		return oldY
	}
	isLabel := false
	for _, label := range layout {
		if mainMenuLabelMatch(e, label.s, label.x, label.y) {
			isLabel = true
		}
	}
	if !isLabel || e.x0*scale < 105*scale || e.x0*scale+width > 214*scale {
		return oldY
	}
	return (e.y0-2)*scale + (9*scale-height+1)/2
}
