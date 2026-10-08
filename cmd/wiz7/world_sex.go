// 034：世界人物頁性別使用既定男女字模，維持原版資料與輸出。
package main

import (
	"github.com/wicanr2/dosgolem/oracle"
	"unicode"
	"unicode/utf8"
)

var worldSexReturn = oracle.Far(0x0110, 0x6159)
var worldSexBackground = [...][3]int{
	{88, 4, 115}, {92, 4, 116}, {95, 4, 115}, {88, 7, 117},
	{89, 9, 115}, {88, 11, 115}, {91, 5, 13}, {95, 5, 14},
}

func worldSexBackgroundMatches(before []byte) bool {
	if len(before) != 64000 {
		return false
	}
	for _, p := range worldSexBackground {
		if before[p[1]*320+p[0]] != byte(p[2]) {
			return false
		}
	}
	return true
}
func worldSexResourcesValid(entries map[string]rosterEntry) bool {
	if len(entries) != 2 {
		return false
	}
	size := 0
	for _, key := range []string{"M", "F"} {
		e, ok := entries[key]
		if !ok || e.En != key || e.Dy != 0 || !utf8.ValidString(e.Zh) || utf8.RuneCountInString(e.Zh) != 1 || len(e.G) != 1 {
			return false
		}
		for _, c := range e.Zh {
			if !unicode.Is(unicode.Han, c) {
				return false
			}
		}
		h := len(e.G[0])
		if (h != 15 && h != 16) || (size != 0 && size != h) {
			return false
		}
		size = h
		nonempty := false
		for _, row := range e.G[0] {
			if uint32(row) >= uint32(1)<<h {
				return false
			}
			nonempty = nonempty || row != 0
		}
		if !nonempty {
			return false
		}
	}
	return true
}
func (ov *overlay) clearWorldSex() {
	for _, e := range ov.evs {
		if e.sex && e.originRet == worldSexReturn {
			e.dead = true
		}
	}
	if ov.sexPending != nil && ov.sexPending.ret == worldSexReturn {
		ov.sexPending = nil
	}
}
func worldSexShift(e *textEv, k int) int {
	if e.sex && e.originRet == worldSexReturn && len(e.g) == 1 && (k == 2 || k == 3) {
		return max(0, len(e.g[0])-6*k)
	}
	return 0
}
