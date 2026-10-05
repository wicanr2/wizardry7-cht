package main

// 024 R966：2倍人物整頁覆繪。只搬原版墨點，不改原版資料。
import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"image"
	"image/color"
	"log"
	"sort"
	"strconv"

	"github.com/wicanr2/dosgolem/oracle"
)

type profileRect struct{ x, y, w, h int }
type profileChange struct {
	at            int
	before, after byte
}
type profileDot struct {
	x, y int
	c    byte
}
type profileField struct {
	role string
	key  string
	cell profileRect
	x, y int
	fg   byte
	g    [][]uint16
	dots []profileDot
}
type profileView struct {
	fields []profileField
	pal    [256][3]uint8
}
type profilePage struct {
	view  profileView
	guard []int
	final []byte
	mask  []profileChange
}
type profileSource struct {
	role string
	key  string
	box  profileRect
	mask []profileChange
	fg   byte
	g    [][]uint16
}
type profilePending struct {
	role       string
	key        string
	ret        oracle.Addr
	sp         uint16
	box        profileRect
	endX, endY int
	fg         byte
	g          [][]uint16
	before     []byte
}
type profileReflowState struct {
	begin   []byte
	pending *profilePending
	health  []byte
	sources map[string]*profileSource
	event   *textEv
	failed  bool
}

var profileAttributes = []string{"STR", "INT", "PIE", "VIT", "DEX", "SPD", "PER", "KAR"}

type profileLabelSite struct {
	role, key string // 空key為catalog精確鍵的動態欄。
	ret       uint16
	x, y      int
	fg        byte
}

var profileLabelSites = []profileLabelSite{
	{"race", "", 0x6c7a, 102, 6, 1},
	{"class", "", 0x6cc0, 90, 12, 132},
	{"LVL", "LVL", 0x6cf8, 90, 18, 132},
	{"RNK", "RNK", 0x6e13, 228, 6, 253},
	{"rank-name", "", 0x6e6b, 288, 6, 253},
	{"EXP", "EXP", 0x6ea3, 228, 12, 71},
	{"MKS", "MKS", 0x6ef9, 228, 18, 71},
	{"HP", "HP", 0x6f5d, 6, 93, 29},
	{"ST", "ST", 0x7007, 6, 99, 29},
}

func profileLabels() []profileLabelSite {
	out := append([]profileLabelSite(nil), profileLabelSites...)
	for i, key := range profileAttributes {
		out = append(out, profileLabelSite{key, key, 0x6a43, 115, 35 + 6*i, 85})
	}
	return out
}

func profileCells() map[string]profileRect {
	m := map[string]profileRect{
		"sex": {176, 8, 16, 16}, "race": {204, 7, 124, 18}, "class": {180, 29, 148, 18},
		"LVL": {204, 51, 34, 18}, "level-value": {240, 51, 36, 18},
		"RNK": {456, 7, 60, 18}, "rank-name": {576, 7, 48, 18},
		"EXP": {456, 29, 48, 18}, "experience-value": {504, 29, 120, 18},
		"MKS": {456, 51, 48, 18}, "kills-value": {504, 51, 120, 18},
		"HP": {12, 184, 34, 18}, "hp-value": {48, 184, 120, 18},
		"ST": {12, 204, 34, 18}, "st-value": {48, 204, 120, 18},
	}
	for i, key := range profileAttributes {
		m[key] = profileRect{230, 70 + 18*i, 34, 18}
		m["attribute-value-"+strconv.Itoa(i)] = profileRect{266, 70 + 18*i, 36, 18}
	}
	return m
}

func (p *profileReflowState) fail(reason string) {
	if !p.failed {
		log.Printf("人物重排缺譯：%s", reason)
	}
	p.failed, p.pending, p.health = true, nil, nil
	if p.event != nil {
		p.event.dead = true
	}
}

func profileGlyphCopy(g [][]uint16) [][]uint16 {
	out := make([][]uint16, len(g))
	for i := range g {
		out[i] = append([]uint16(nil), g[i]...)
	}
	return out
}

func profileGlyphsValid(g [][]uint16, cell profileRect) bool {
	if len(g) == 0 || len(g[0]) != 15 && len(g[0]) != 16 {
		return false
	}
	h := len(g[0])
	if h > cell.h || len(g)*(h+2)-2 > cell.w {
		return false
	}
	for _, glyph := range g {
		if len(glyph) != h {
			return false
		}
		nonempty := false
		for _, row := range glyph {
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

func (ov *overlay) enterProfileReflowSex(ret oracle.Addr, sp, font, code uint16, x, y int, before []byte) {
	if ret != profileSexReturn || x != 90 || y != 6 {
		return
	}
	if old := ov.profile; old != nil && old.event != nil {
		old.event.dead = true
	}
	p := &profileReflowState{sources: make(map[string]*profileSource)}
	ov.profile = p
	if font != 0 || code != 'M' && code != 'F' || len(before) != 64000 ||
		ov.pend != nil || ov.compound != nil || ov.sexPending != nil {
		p.fail("開始訊號不符")
		return
	}
	g := ov.sexEntries[string(rune(code))].G
	if !profileGlyphsValid(g, profileCells()["sex"]) {
		p.fail("性別字模不完整")
		return
	}
	p.begin = append([]byte(nil), before...)
	p.pending = &profilePending{role: "sex", key: string(rune(code)), ret: ret, sp: sp, box: profileRect{90, 6, 6, 6},
		endX: 96, endY: 6, fg: 1, g: profileGlyphCopy(g), before: append([]byte(nil), before...)}
}

func profileASCII(raw []byte) (string, bool) {
	n := bytes.IndexByte(raw, 0)
	if n < 1 {
		return "", false
	}
	for _, c := range raw[:n] {
		if c < 32 || c > 126 {
			return "", false
		}
	}
	return string(raw[:n]), true
}

// 此回傳值只表示認領入口；不代替原版API結果。
func (ov *overlay) enterProfileReflowText(ret oracle.Addr, sp, font uint16, x, y int, raw []byte, trans map[string]entry, before func() []byte) bool {
	p := ov.profile
	if !ov.profileEnabled || p == nil || ret.Seg != 0x0110 {
		return false
	}
	if ret.Off == 0x032a && x == 225 && y == 138 && p.health != nil {
		if font != 3 || p.pending != nil {
			p.fail("生命體力邊界不符")
		} else {
			p.finishHealth(before())
		}
		return false // 邊界本身不是人物欄，也不吞一般輸出。
	}
	role, fg, labelKey := "", byte(0), ""
	for _, site := range profileLabels() {
		if ret.Off == site.ret && x == site.x && y == site.y {
			role, fg, labelKey = site.role, site.fg, site.key
			break
		}
	}
	if role == "" && ret.Off == 0x032a {
		switch {
		case x == 126 && y == 18:
			role, fg = "level-value", 132
		case x == 306 && y == 12:
			role, fg = "experience-value", 71
		case x == 306 && y == 18:
			role, fg = "kills-value", 71
		case (x == 139 || x == 145) && y >= 35 && y <= 77 && (y-35)%6 == 0:
			role, fg = "attribute-value-"+strconv.Itoa((y-35)/6), 85
		}
	}
	if role == "" {
		known := ret.Off == 0x032a && font == 0
		for _, site := range profileLabels() {
			known = known || ret.Off == site.ret
		}
		if known {
			p.fail("未量來源座標")
		}
		return false
	}
	if p.failed {
		return true
	}
	if p.event != nil || p.pending != nil || p.sources[role] != nil || ov.pend != nil || ov.compound != nil || ov.sexPending != nil ||
		p.health != nil && role != "HP" && role != "ST" {
		p.fail("來源重複或交錯")
		return true
	}
	s, ok := profileASCII(raw)
	if font != 0 || !ok || labelKey != "" && s != labelKey {
		p.fail("文字tuple不符")
		return true
	}
	var glyphs [][]uint16
	if ret.Off == 0x032a {
		n := 1
		if x == 139 {
			n = 2
		}
		if len(s) != n {
			p.fail("未量數值長度")
			return true
		}
		for _, c := range []byte(s) {
			if c < '0' || c > '9' {
				p.fail("數值非十進位")
				return true
			}
		}
	} else {
		e, exists := trans[s]
		if !exists || e.Dy != 0 || !profileGlyphsValid(e.G, profileCells()[role]) {
			p.fail("欄位字模不完整")
			return true
		}
		glyphs = profileGlyphCopy(e.G)
	}
	b := before()
	if len(b) != 64000 || x+6*len(s) > 320 {
		p.fail("來源離屏或字格不符")
		return true
	}
	if role == "HP" {
		p.health = append([]byte(nil), b...)
	}
	p.pending = &profilePending{role: role, key: s, ret: ret, sp: sp, box: profileRect{x, y, 6 * len(s), 6},
		endX: x + 6*len(s), endY: y, fg: fg, g: glyphs, before: append([]byte(nil), b...)}
	return true
}

func profileChanges(before, after []byte, r profileRect, colors map[byte]bool) ([]profileChange, bool) {
	if len(before) != 64000 || len(after) != 64000 || r.x < 0 || r.y < 0 || r.x+r.w > 320 || r.y+r.h > 200 {
		return nil, false
	}
	var changes []profileChange
	for at, a := range after {
		b := before[at]
		if a == b {
			continue
		}
		x, y := at%320, at/320
		if x < r.x || x >= r.x+r.w || y < r.y || y >= r.y+r.h || !colors[a] {
			return nil, false
		}
		changes = append(changes, profileChange{at, b, a})
	}
	return changes, len(changes) > 0
}

func (ov *overlay) returnedProfileReflow(ret oracle.Addr, sp uint16, x, y int, after []byte) {
	p := ov.profile
	if !ov.profileEnabled || p == nil || p.pending == nil || p.pending.ret != ret {
		return
	}
	q := p.pending
	p.pending = nil
	mask, ok := profileChanges(q.before, after, q.box, map[byte]bool{q.fg: true})
	if sp != q.sp+2 || x != q.endX || y != q.endY || !ok {
		p.fail("真正返回或筆畫不符")
		return
	}
	p.sources[q.role] = &profileSource{role: q.role, key: q.key, box: q.box, mask: mask, fg: q.fg, g: q.g}
	if q.role == "attribute-value-7" {
		page, valid := p.seal(after)
		if !valid {
			p.fail("整頁來源或背景不完整")
			return
		}
		e := &textEv{s: "character.review.reflow", profile: page}
		p.event = e
		ov.evs = append(ov.evs, e)
		ov.created++
	}
}

func (p *profileReflowState) finishHealth(after []byte) {
	mask, ok := profileChanges(p.health, after, profileRect{6, 93, 78, 12}, map[byte]bool{27: true, 29: true})
	p.health = nil
	if !ok || p.sources["HP"] == nil || p.sources["ST"] == nil {
		p.fail("生命體力來源不符")
		return
	}
	labels := map[int]profileChange{}
	for _, key := range []string{"HP", "ST"} {
		for _, c := range p.sources[key].mask {
			labels[c.at] = c
		}
	}
	for _, c := range mask {
		if c.at%320 < 24 {
			if old, found := labels[c.at]; !found || old != c {
				p.fail("生命體力標籤重疊不符")
				return
			}
			delete(labels, c.at)
		}
	}
	if len(labels) != 0 {
		p.fail("生命體力標籤缺筆畫")
		return
	}
	p.sources["health-band"] = &profileSource{role: "health-band", box: profileRect{6, 93, 78, 12}, mask: mask}
}

func (p *profileReflowState) seal(final []byte) (*profilePage, bool) {
	if p.failed || p.pending != nil || p.health != nil || len(p.sources) != 30 || len(final) != 64000 || len(p.begin) != 64000 {
		return nil, false
	}
	page := &profilePage{final: append([]byte(nil), final...)}
	unique, owner := map[int]profileChange{}, map[int]string{}
	guard := map[int]bool{}
	keys := make([]string, 0, len(p.sources))
	for key := range p.sources {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		s := p.sources[key]
		for _, c := range s.mask {
			if final[c.at] != c.after {
				return nil, false
			}
			if old, exists := unique[c.at]; exists {
				o := owner[c.at]
				sharedHealth := key == "health-band" && (o == "HP" || o == "ST") || o == "health-band" && (key == "HP" || key == "ST")
				if !sharedHealth || old != c {
					return nil, false
				}
			} else {
				unique[c.at], owner[c.at] = c, key
			}
		}
		for y := s.box.y; y < s.box.y+s.box.h; y++ {
			for x := s.box.x; x < s.box.x+s.box.w; x++ {
				guard[y*320+x] = true
			}
		}
	}
	clean := append([]byte(nil), final...)
	for at, c := range unique {
		clean[at] = c.before
		page.mask = append(page.mask, c)
	}
	cells := profileCells()
	roles := make([]string, 0, len(cells))
	for role := range cells {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	height := 0
	for _, role := range roles {
		cell := cells[role]
		f := profileField{role: role, cell: cell}
		for yy := cell.y; yy < cell.y+cell.h; yy++ {
			for xx := cell.x; xx < cell.x+cell.w; xx++ {
				at := (yy/2)*320 + xx/2
				if clean[at] != p.begin[at] {
					return nil, false
				}
				guard[at] = true
			}
		}
		if s := p.sources[role]; s != nil && len(s.g) > 0 {
			if !profileGlyphsValid(s.g, cell) || height != 0 && height != len(s.g[0]) {
				return nil, false
			}
			height = len(s.g[0])
			f.g, f.fg, f.x, f.key = profileGlyphCopy(s.g), s.fg, cell.x, s.key
			f.y = cell.y + (cell.h-height+1)/2
			if role == "sex" {
				f.x += (cell.w - height + 1) / 2
			}
		} else {
			s := p.sources[role]
			sourceY := -1
			if role == "hp-value" || role == "st-value" {
				s, sourceY = p.sources["health-band"], 93
				if role == "st-value" {
					sourceY = 99
				}
			}
			if s == nil || len(s.g) > 0 {
				return nil, false
			}
			var ink []profileChange
			minY, maxY := 200, -1
			for _, c := range s.mask {
				x, y := c.at%320, c.at/320
				if sourceY >= 0 && (x < 24 || y < sourceY || y >= sourceY+6) {
					continue
				}
				ink = append(ink, c)
				minY, maxY = min(minY, y), max(maxY, y)
			}
			inkHeight := (maxY - minY + 1) * 2
			if len(ink) == 0 || inkHeight > cell.h {
				return nil, false
			}
			paintY := cell.y + (cell.h-inkHeight+1)/2
			for _, c := range ink {
				for dy := 0; dy < 2; dy++ {
					for dx := 0; dx < 2; dx++ {
						f.dots = append(f.dots, profileDot{2*(c.at%320) + dx, paintY + 2*(c.at/320-minY) + dy, c.after})
					}
				}
			}
		}
		page.view.fields = append(page.view.fields, f)
	}
	occupied := map[int]bool{}
	valid := true
	for _, f := range page.view.fields {
		profilePaintField(f, func(x, y int, _ byte) {
			if x < f.cell.x || x >= f.cell.x+f.cell.w || y < f.cell.y || y >= f.cell.y+f.cell.h || occupied[y*640+x] {
				valid = false
			}
			occupied[y*640+x] = true
		})
	}
	if !valid || len(page.view.fields) != 31 {
		return nil, false
	}
	for at := range guard {
		page.guard = append(page.guard, at)
	}
	sort.Ints(page.guard)
	sort.Slice(page.mask, func(i, j int) bool { return page.mask[i].at < page.mask[j].at })
	return page, true
}

func (p *profilePage) matches(a []byte, cursor func(int, int) bool) bool {
	if len(a) != 64000 || len(p.final) != 64000 || len(p.guard) == 0 {
		return false
	}
	for _, at := range p.guard {
		if cursor != nil && cursor(at%320, at/320) {
			continue
		}
		if a[at] != p.final[at] {
			return false
		}
	}
	return true
}

func (p *profilePage) restore(comp []byte, cursor func(int, int) bool) {
	for _, c := range p.mask {
		if cursor == nil || !cursor(c.at%320, c.at/320) {
			comp[c.at] = c.before
		}
	}
}

func profilePaintField(f profileField, emit func(int, int, byte)) {
	for i, glyph := range f.g {
		h := len(glyph)
		for y, row := range glyph {
			for x := 0; x < h; x++ {
				if row&(1<<(h-1-x)) != 0 {
					emit(f.x+i*(h+2)+x, f.y+y, f.fg)
				}
			}
		}
	}
	for _, dot := range f.dots {
		emit(dot.x, dot.y, dot.c)
	}
}

func (v *profileView) paint(img *image.RGBA) {
	for _, f := range v.fields {
		profilePaintField(f, func(x, y int, c byte) {
			rgb := v.pal[c]
			img.SetRGBA(x, y, color.RGBA{rgb[0], rgb[1], rgb[2], 255})
		})
	}
}

func (p *profilePage) visualCopy() *profilePage {
	out := &profilePage{view: profileView{pal: p.view.pal},
		guard: append([]int(nil), p.guard...), final: append([]byte(nil), p.final...),
		mask: append([]profileChange(nil), p.mask...)}
	for _, f := range p.view.fields {
		f.g = profileGlyphCopy(f.g)
		f.dots = append([]profileDot(nil), f.dots...)
		out.view.fields = append(out.view.fields, f)
	}
	return out
}

type profileFieldMetadata struct {
	Role, Key, GlyphSHA256, DotsSHA256 string
	Cell                               [4]int
	Origin                             [2]int
	NativeHeight, DotPixels            int
}
type profileMetadata struct {
	Fields        []profileFieldMetadata
	GuardSHA256   string
	MaskSHA256    string
	NativeFields  int
	NumericFields int
}

func (p *profilePage) metadata() *profileMetadata {
	if p == nil {
		return nil
	}
	m := &profileMetadata{}
	appendInt := func(dst []byte, n int) []byte {
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(n))
		return append(dst, b[:]...)
	}
	hash := func(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
	var guard, mask []byte
	for _, at := range p.guard {
		guard = appendInt(guard, at)
		guard = append(guard, p.final[at])
	}
	for _, c := range p.mask {
		mask = appendInt(mask, c.at)
		mask = append(mask, c.before, c.after)
	}
	m.GuardSHA256, m.MaskSHA256 = hash(guard), hash(mask)
	for _, f := range p.view.fields {
		item := profileFieldMetadata{Role: f.role, Key: f.key, Cell: [4]int{f.cell.x, f.cell.y, f.cell.w, f.cell.h}, Origin: [2]int{f.x, f.y}, DotPixels: len(f.dots)}
		var glyphs, dots []byte
		if len(f.g) > 0 {
			m.NativeFields++
			item.NativeHeight = len(f.g[0])
			glyphs = appendInt(appendInt(glyphs, len(f.g)), item.NativeHeight)
			for _, g := range f.g {
				for _, row := range g {
					glyphs = append(glyphs, byte(row>>8), byte(row))
				}
			}
		} else {
			m.NumericFields++
		}
		for _, dot := range f.dots {
			dots = appendInt(appendInt(dots, dot.x), dot.y)
			dots = append(dots, dot.c)
		}
		item.GlyphSHA256, item.DotsSHA256 = hash(glyphs), hash(dots)
		m.Fields = append(m.Fields, item)
	}
	return m
}
