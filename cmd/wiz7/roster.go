package main

// R994 disposable prototype: bounded context identifies the row, never its name.
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
	code                string
	y, nameLen, slot    int
	blankName, selected bool
	scope               rosterScope
}

var rosterLegacyCodes = []string{"M-FIG", "F-FIG", "M-THI", "M-PRI", "M-MAG", "F-ALC"}
var rosterProfessionCodes = []string{"FIG", "MAG", "PRI", "THI", "RAN", "ALC", "BAR", "PSI", "VAL", "BIS", "LOR", "SAM", "MON", "NIN"}

func rosterAllowedCodes(count int) map[string]bool {
	allowed := map[string]bool{}
	switch count {
	case 6:
		for _, code := range rosterLegacyCodes {
			allowed[code] = true
		}
	case 28:
		for _, profession := range rosterProfessionCodes {
			allowed["M-"+profession] = true
			allowed["F-"+profession] = true
		}
	default:
		return nil
	}
	return allowed
}

func rosterResourcesValid(entries map[string]rosterEntry) bool {
	allowed := rosterAllowedCodes(len(entries))
	if allowed == nil {
		return false
	}
	for key, e := range entries {
		if !allowed[key] || e.En != key {
			return false
		}
	}
	return true
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
		err != nil || len(hash) != 32 || rosterAllowedCodes(len(data.Entries)) == nil {
		return nil, fmt.Errorf("職業碼資源標頭不合格")
	}
	allowed := rosterAllowedCodes(len(data.Entries))
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
			log.Printf("職業碼：載入%d鍵", len(entries))
			return entries
		}
	}
	log.Printf("職業碼停用：%v", err)
	return nil
}

type rosterWitness struct {
	after  []byte
	normal bool
	scope  rosterScope
	slot   int
}

func rosterKnownY(y int) bool {
	return y >= 150 && y <= 180 && (y-150)%6 == 0
}

func (ov *overlay) retireRosterAt(y int) {
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
}

func (ov *overlay) applyRosterScope(c rosterContext) {
	for _, e := range ov.evs {
		if e.roster != nil && (e.roster.scope != c.scope || e.roster.slot < c.scope.first || e.roster.slot >= c.scope.first+c.scope.count) {
			e.dead = true
		}
	}
}

func (ov *overlay) enterRoster(ret oracle.Addr, sp, font, off uint16, x, y int, evidence *rosterContextEvidence, raw []byte, before []byte) bool {
	if !ov.rosterEnabled || ret != oracle.Far(0x0110, 0x0fcc) {
		return false
	}
	ov.retireRosterAt(y)
	c, err := validateRosterContext(evidence, off, font, x, y)
	if sp > 65530 || err != nil || evidence.registers.SP != sp {
		log.Printf("職業碼缺譯：數值或來源身份不符 y=%d", y)
		return true
	}
	ov.applyRosterScope(c)
	if int(off)+16 > 65536 || len(raw) != 16 || raw[15] != 0 || len(before) != 64000 {
		log.Printf("職業碼缺譯：完整列或離屏不符 y=%d", y)
		return true
	}
	blank := true
	for _, b := range raw[:10] {
		if b < 32 {
			log.Printf("職業碼缺譯：名字控制字元 y=%d", y)
			return true
		}
		blank = blank && b == ' '
	}
	code := string(raw[10:15])
	e, ok := ov.roster[code]
	if !rosterResourcesValid(ov.roster) || !ok {
		log.Printf("職業碼缺譯：資源不符 y=%d", y)
		return true
	}
	row := &rosterRow{code: code, y: y, nameLen: 10, slot: c.slot, blankName: blank, selected: c.selected, scope: c.scope}
	ov.pend = &pendingPrint{s: code, x0: 210, y0: y, sp: sp, ret: ret,
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
	if (!row.blankName && nameChanges == 0) || codeChanges == 0 || !havePair ||
		(row.selected && (bg != 1 || fg != 0)) || (!row.selected && (bg != 195 || fg != 180)) {
		return nil
	}
	e := &textEv{s: p.s, x0: 210, y0: row.y, bx: 210, by: row.y, bw: 30, bh: 6,
		g: p.g, fg: fg, fgRGB: pal[fg], roster: &rosterWitness{normal: bg == 195, scope: row.scope, slot: row.slot}}
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
		exact, solidPadding := true, true
		for i := 0; i < 90; i++ {
			x, y := 150+i, e.y0+j
			if ignore != nil && ignore(x, y) {
				continue
			}
			a := screen[y*320+x]
			exact = exact && a == e.roster.after[j*90+i]
			if padding {
				if a != 195 && a != 1 || paddingColor != 255 && a != paddingColor {
					solidPadding = false
				}
				paddingColor = a
			}
		}
		// 022 R995：姓名下伸點屬於同次原版after；只有整行失配才沿舊單色底邊例外。
		if !exact && !(padding && solidPadding) {
			return false
		}
	}
	return true
}
