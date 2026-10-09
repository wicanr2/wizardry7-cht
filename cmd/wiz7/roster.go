package main

// 022：六個完整來源及已量位置。姓名及遊戲狀態只讀，還原只限職業碼。
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/wicanr2/dosgolem/oracle"
)

type rosterRow struct {
	source, code string
	y, nameLen   int
}

var rosterRows = []rosterRow{
	{"THESUS    M-FIG", "M-FIG", 150, 6},
	{"TEMPEST   F-FIG", "F-FIG", 156, 7},
	{"LYSAND    M-THI", "M-THI", 162, 6},
	{"NOBAL     M-PRI", "M-PRI", 168, 5},
	{"TREON     M-MAG", "M-MAG", 174, 5},
	{"PENTAS    F-ALC", "F-ALC", 180, 6},
}

type rosterEntry struct {
	En string     `json:"en"`
	Zh string     `json:"zh"`
	Dy int        `json:"dy"`
	G  [][]uint16 `json:"g"`
}

type rosterData struct {
	Schema       int           `json:"schema"`
	Size         int           `json:"size"`
	Context      string        `json:"context"`
	Source       string        `json:"source"`
	SourceSHA256 string        `json:"source_sha256"`
	Entries      []rosterEntry `json:"entries"`
}

func decodeRoster(r io.Reader, size int) (map[string]rosterEntry, error) {
	d := json.NewDecoder(io.LimitReader(r, 1<<20+1))
	d.DisallowUnknownFields()
	var data rosterData
	if err := d.Decode(&data); err != nil {
		return nil, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("職業碼須為單一JSON文件")
	}
	hash, err := hex.DecodeString(data.SourceSHA256)
	if data.Schema != 1 || (size != 15 && size != 16) || data.Size != size ||
		data.Context != "character.add.roster" || data.Source != "runtime:0110:0fcc" ||
		err != nil || len(hash) != 32 || len(data.Entries) != 6 {
		return nil, fmt.Errorf("職業碼資源標頭不合格")
	}
	allowed := map[string]bool{}
	for _, row := range rosterRows {
		allowed[row.code] = true
	}
	entries := map[string]rosterEntry{}
	for _, e := range data.Entries {
		count := utf8.RuneCountInString(e.Zh)
		if !allowed[e.En] || entries[e.En].En != "" || e.Dy != 0 ||
			!utf8.ValidString(e.Zh) || count == 0 || count != len(e.G) || count*(size+2)-2 > 90 {
			return nil, fmt.Errorf("職業碼鍵或文案不合格")
		}
		for _, c := range e.Zh {
			if !unicode.Is(unicode.Han, c) {
				return nil, fmt.Errorf("職業碼譯文須為中文字元")
			}
		}
		for _, g := range e.G {
			if len(g) != size || size == 16 && g[15] != 0 {
				return nil, fmt.Errorf("職業碼字模高度或底邊不合格")
			}
			nonempty := false
			for _, word := range g {
				if uint32(word) >= uint32(1)<<size {
					return nil, fmt.Errorf("職業碼字模位元越界")
				}
				nonempty = nonempty || word != 0
			}
			if !nonempty {
				return nil, fmt.Errorf("職業碼字模空白")
			}
		}
		entries[e.En] = e
	}
	return entries, nil
}

func loadRoster(xlate string, trans map[string]entry) map[string]rosterEntry {
	if xlate == "" || len(trans) == 0 {
		return nil
	}
	size := 0
	for _, e := range trans {
		if len(e.G) > 0 {
			size = len(e.G[0])
			break
		}
	}
	f, err := os.Open(filepath.Join(filepath.Dir(xlate), "roster.json"))
	if err == nil {
		defer f.Close()
		var st os.FileInfo
		st, err = f.Stat()
		if err == nil && st.Size() > 1<<20 {
			err = fmt.Errorf("職業碼資源過大")
		}
	}
	if err == nil {
		var entries map[string]rosterEntry
		entries, err = decodeRoster(f, size)
		if err == nil {
			log.Print("職業碼：載入六鍵")
			return entries
		}
	}
	log.Printf("職業碼停用：%v", err)
	return nil
}

type rosterWitness struct {
	after  []byte
	normal bool
}

func rosterKnownY(y int) bool {
	for _, row := range rosterRows {
		if row.y == y {
			return true
		}
	}
	return false
}

// 022 R991：只增四個已量移位，回傳副本以保持原表不變。
func rosterRowAt(source string, y int) *rosterRow {
	for _, row := range rosterRows {
		if row.source != source {
			continue
		}
		allowed := row.y == y
		switch row.code {
		case "M-THI":
			allowed = allowed || y == 156
		case "M-PRI":
			allowed = allowed || y == 162
		case "M-MAG":
			allowed = allowed || y == 168
		case "F-ALC":
			allowed = allowed || y == 174
		}
		if !allowed {
			return nil
		}
		row.y = y
		return &row
	}
	return nil
}

func (ov *overlay) enterRoster(ret oracle.Addr, sp, font, off uint16, x, y int, raw []byte, before []byte) bool {
	if !ov.rosterEnabled || ret != (oracle.Addr{Seg: 0x0110, Off: 0x0fcc}) {
		return false
	}
	if !rosterKnownY(y) {
		return false
	}
	for _, e := range ov.evs {
		if e.roster != nil && e.y0 == y {
			e.dead = true
		}
	}
	if ov.pend != nil {
		ov.nested++
		if ov.pend.page != nil {
			ov.pend.page.fail("巢狀", 0)
		}
		ov.pend = nil
	}
	if sp > 65530 || font != 0 || x != 150 || int(off)+16 > 65536 || len(raw) != 16 ||
		raw[15] != 0 || len(before) != 64000 {
		log.Printf("職業碼缺譯：訊號或完整列不符 y=%d", y)
		return true
	}
	row := rosterRowAt(string(raw[:15]), y)
	if row == nil {
		log.Printf("職業碼缺譯：完整列或位置不符 y=%d", y)
		return true
	}
	e, ok := ov.roster[row.code]
	if !ok {
		log.Printf("職業碼缺譯：資源不符 y=%d", y)
		return true
	}
	ov.pend = &pendingPrint{s: row.source, x0: 210, y0: y, sp: sp, ret: ret,
		font: font, before: append([]byte(nil), before...), g: e.G, roster: row}
	return true
}

func finishRoster(p *pendingPrint, after []byte, pal [256][3]uint8) *textEv {
	if len(p.before) != 64000 || len(after) != 64000 || p.roster == nil {
		return nil
	}
	row, nameChanges, codeChanges := p.roster, 0, 0
	bg, fg, havePair := byte(0), byte(0), false
	for i, a := range after {
		b := p.before[i]
		if a == b {
			continue
		}
		x, y := i%320, i/320
		if y < row.y || y >= row.y+6 || !(b == 195 && a == 180 || b == 1 && a == 0) {
			return nil
		}
		if havePair && (b != bg || a != fg) {
			return nil
		}
		bg, fg, havePair = b, a, true
		switch {
		case x >= 150 && x < 150+6*row.nameLen:
			nameChanges++
		case x >= 210 && x < 240:
			codeChanges++
		default:
			return nil
		}
	}
	if nameChanges == 0 || codeChanges == 0 {
		return nil
	}
	e := &textEv{s: p.s, x0: 210, y0: row.y, bx: 210, by: row.y, bw: 30, bh: 6,
		g: p.g, fg: fg, fgRGB: pal[fg], roster: &rosterWitness{normal: bg == 195}}
	for j := 0; j < 6; j++ {
		start := (row.y + j) * 320
		e.before = append(e.before, p.before[start+210:start+240]...)
		e.after = append(e.after, after[start+210:start+240]...)
		e.roster.after = append(e.roster.after, after[start+150:start+240]...)
	}
	return e
}

func (e *textEv) rosterMatches(screen []byte, ignore func(int, int) bool) bool {
	if len(screen) != 64000 || e.roster == nil || len(e.roster.after) != 540 {
		return false
	}
	for j := 0; j < 6; j++ {
		padding := e.roster.normal && j == 5
		paddingColor := byte(255)
		for i := 0; i < 90; i++ {
			x, y := 150+i, e.y0+j
			if ignore != nil && ignore(x, y) {
				continue
			}
			a := screen[y*320+x]
			if padding {
				if a != 195 && a != 1 || paddingColor != 255 && a != paddingColor {
					return false
				}
				paddingColor = a
			} else if a != e.roster.after[j*90+i] {
				return false
			}
		}
	}
	return true
}
