package main

// 036：側欄短碼只走已量caller，姓名永遠保持原版。
import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/wicanr2/dosgolem/oracle"
)

const partyClassKey = "ALCHEMI"
const partyNameFontSHA = "937eda16ebf61f28506122c465cd3a3e112aa001cf0c37644b001bb0cbc35ae3"
const partyReturnLimit uint64 = 1000000

var partyNameReturn = oracle.Far(0x0110, 0x1c21)
var partyClassReturn = oracle.Far(0x0110, 0x1c71)

type partyNameState struct {
	valid bool
	mask  [42]bool
}

type pendingPartyName struct {
	ret        oracle.Addr
	sp         uint16
	entrySteps uint64
	state      partyNameState
}

type pendingPartyClass struct {
	slot, size int
	entrySteps uint64
	nameRow    [42]byte
}

type partyClassWitness struct {
	slot, size int
	nameRow    [42]byte
}

func partySlot(x, y int, name bool) int {
	base := 56
	if name {
		base = 50
	}
	if (x != 9 && x != 261) || y < base || y > base+90 || (y-base)%45 != 0 {
		return -1
	}
	side := 0
	if x == 261 {
		side = 1
	}
	return 2*((y-base)/45) + side
}

func partyClassEntryValid(ent entry, size int) bool {
	return partyClassEntryValidForKey(partyClassKey, ent, size)
}

func partyClassEntryValidForKey(key string, ent entry, size int) bool {
	count := 2
	if key == partyClassKey {
		count = 4
	}
	if !partyClassMeasuredKey(key) || (size != 15 && size != 16) || ent.Dy != 0 || len(ent.D) != 0 || len(ent.G) != count {
		return false
	}
	for _, glyph := range ent.G {
		if len(glyph) != size {
			return false
		}
		nonempty := false
		for _, row := range glyph {
			if uint32(row) >= uint32(1)<<size {
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

func partyClassConfig(trans map[string]entry, scale, hd int) (enabled bool, size int) {
	if (scale != 2 && scale != 3) || scale != hd {
		return false, 0
	}
	size = ordinaryFontSize(trans)
	ent, ok := trans[partyClassKey]
	if !ok || !partyClassEntryValid(ent, size) {
		return false, 0
	}
	return true, size
}

// 036：專路資源不進一般normalizeKey或placeholder查表。
func takePartyClassEntry(trans map[string]entry) entry {
	ent := trans[partyClassKey]
	delete(trans, partyClassKey)
	return ent
}

// slot<0表示尚無可採的座標；撤除全部本家族，不能沿用舊人物。
func (ov *overlay) retirePartyClass(slot int) {
	for _, e := range ov.evs {
		if e != nil && e.party != nil && (slot < 0 || e.party.slot == slot) {
			e.dead = true
		}
	}
	if ov.pend != nil && ov.pend.party != nil && (slot < 0 || ov.pend.party.slot == slot) {
		ov.nested++
		ov.pend = nil
	}
}

func partyRAMSpan(address oracle.Addr, size int) bool {
	return size > 0 && size <= 65536 && int(address.Off) <= 65536-size &&
		int(address.Seg)*16+int(address.Off)+size <= 0xa0000
}

// 只讀已量font指標及1040B；不保存字型bytes或姓名文字。
func partyNameMask(name []byte, read func(oracle.Addr, int) []byte) (state partyNameState) {
	if len(name) < 1 || len(name) > 10 || read == nil {
		return state
	}
	for _, b := range name {
		if b < 32 || b > 126 {
			return state
		}
	}
	pointer := oracle.Far(0x0110, 0x364a)
	if !partyRAMSpan(pointer, 2) {
		return state
	}
	raw := read(pointer, 2)
	if len(raw) != 2 {
		return state
	}
	segment := uint16(raw[0]) | uint16(raw[1])<<8
	addr := oracle.Far(segment, 0)
	if !partyRAMSpan(addr, 1040) {
		return state
	}
	font := read(addr, 1040)
	if len(font) != 1040 || fmt.Sprintf("%x", sha256.Sum256(font)) != partyNameFontSHA {
		return state
	}
	header := [16]byte{6, 6, 0, 1, 1, 128, 2, 0, 0, 0, 6, 0, 16, 4, 6, 0}
	for i, v := range header {
		if font[i] != v {
			return state
		}
	}
	for index, value := range name {
		mapping := 16 + int(value)*2
		if mapping < 16 || mapping+2 > 272 {
			return partyNameState{}
		}
		glyph, advance := int(font[mapping]), int(font[mapping+1])
		if glyph >= 128 || advance != 6 {
			return partyNameState{}
		}
		offset := 272 + 6*glyph + 5
		if offset < 272 || offset >= len(font) {
			return partyNameState{}
		}
		for column := 0; column < 6; column++ {
			position := index*6 + column
			if position < 42 && font[offset]&(1<<uint(7-column)) != 0 {
				state.mask[position] = true
			}
		}
	}
	state.valid = true
	return state
}

func (ov *overlay) enterPartyName(ret oracle.Addr, sp, font, off uint16, x, y int, raw []byte, step uint64, read func(oracle.Addr, int) []byte) bool {
	if ret != partyNameReturn {
		return false
	}
	slot := partySlot(x, y, true)
	ov.retirePartyClass(slot)
	if slot < 0 {
		ov.partyNames = [6]partyNameState{}
		ov.partyPendingNames = [6]*pendingPartyName{}
		log.Print("側欄姓名保留英文：座標不符")
		return true
	}
	ov.partyNames[slot] = partyNameState{}
	ov.partyPendingNames[slot] = nil
	if font != 0 || sp > 65530 || off > 65504 || len(raw) < 32 {
		log.Print("側欄姓名保留英文：font、stack或字串界限不符")
		return true
	}
	n := 0
	for n < 32 && raw[n] != 0 {
		n++
	}
	if n < 1 || n > 10 {
		log.Print("側欄姓名保留英文：長度不符")
		return true
	}
	for _, b := range raw[:n] {
		if b < 32 || b > 126 {
			log.Print("側欄姓名保留英文：字元不符")
			return true
		}
	}
	if !ov.partyEnabled || ov.partySize != 16 || ov.partyScale != 2 {
		return true
	}
	state := partyNameMask(raw[:n], read)
	if !state.valid {
		log.Print("側欄姓名遮罩不可用：字型身份或讀取界限不符")
		return true
	}
	ov.partyPendingNames[slot] = &pendingPartyName{ret: ret, sp: sp, entrySteps: step, state: state}
	return true
}

func (ov *overlay) returnPartyNames(ret oracle.Addr, sp uint16, step uint64) {
	for slot, p := range ov.partyPendingNames {
		if p == nil || p.ret != ret {
			continue
		}
		ov.partyPendingNames[slot] = nil
		if p.sp > 65530 || sp != p.sp+2 || step < p.entrySteps || step-p.entrySteps > partyReturnLimit || !p.state.valid {
			ov.partyNames[slot] = partyNameState{}
			log.Print("側欄姓名遮罩未提交：返回身份不符")
			continue
		}
		ov.partyNames[slot] = p.state
	}
}

// 036 R996：只擴四個已有文案的geometry，其他caller保持一般路徑。
func partyClassMeasuredKey(key string) bool {
	switch key {
	case partyClassKey, "FIGHTER", "THIEF", "PRIEST", "MAGE":
		return true
	}
	return false
}

func partyClassRawKey(raw []byte) (string, bool) {
	if len(raw) < 32 {
		return "", false
	}
	n := 0
	for n < 32 && raw[n] != 0 {
		n++
	}
	if n == 0 || n == 32 {
		return "", false
	}
	key := string(raw[:n])
	if !partyClassMeasuredKey(key) {
		return "", false
	}
	return key, true
}

// near拒絕先於legacy fallback，ALCHEMI始終隔離於一般trans之外。
func partyClassShouldRoute(ret oracle.Addr, s, key string, enabled bool) bool {
	if s == partyClassKey {
		return true
	}
	return ret == partyClassReturn && partyClassMeasuredKey(key) && (enabled || s != key)
}

func partyClassRaw(raw []byte) bool {
	key, exact := partyClassRawKey(raw)
	return exact && key == partyClassKey
}

func partyClassOffset(key string, size, k int) (dx, dy int) {
	if size != 16 || k != 2 || !partyClassMeasuredKey(key) {
		return 0, 0
	}
	if key == "FIGHTER" {
		return 5, -1
	}
	return 2, -1
}

func partyClassInkFits(ent entry, size, scale int, name partyNameState) bool {
	return partyClassInkFitsForKey(partyClassKey, ent, size, scale, name)
}

func partyClassInkFitsForKey(key string, ent entry, size, scale int, name partyNameState) bool {
	if !partyClassEntryValidForKey(key, ent, size) || (scale != 2 && scale != 3) {
		return false
	}
	dx, dy := partyClassOffset(key, size, scale)
	if dy < 0 && !name.valid {
		return false
	}
	for index, glyph := range ent.G {
		for row, bits := range glyph {
			for column := 0; column < size; column++ {
				if bits&(1<<uint(size-1-column)) == 0 {
					continue
				}
				xx, yy := index*(size+2)+column+dx, row+dy
				if xx < 0 || xx >= 42*scale || yy >= 7*scale || yy < -1 {
					return false
				}
				if yy < 0 && (size != 16 || scale != 2 || !name.valid || name.mask[xx/scale]) {
					return false
				}
			}
		}
	}
	return true
}

func (ov *overlay) enterPartyClass(ret oracle.Addr, sp, font, off uint16, x, y int, raw []byte, step uint64, before []byte, ent entry) bool {
	key, exact := partyClassRawKey(raw)
	if !exact && ret != partyClassReturn {
		return false
	}
	// 四既有鍵在其他caller不走專路，缺ALCHEMI配置也仍用原普通路徑。
	if exact && key != partyClassKey && (ret != partyClassReturn || !ov.partyEnabled) {
		return false
	}
	slot := partySlot(x, y, false)
	ov.retirePartyClass(slot)
	if !exact || ret != partyClassReturn || slot < 0 || font != 0 || sp > 65530 || off > 65504 || !ov.partyEnabled || len(before) != 64000 ||
		!partyClassInkFitsForKey(key, ent, ov.partySize, ov.partyScale, func() partyNameState {
			if slot >= 0 {
				return ov.partyNames[slot]
			}
			return partyNameState{}
		}()) {
		log.Print("側欄職業缺譯：來源、資源、倍率或姓名保護不符")
		return true
	}
	// 沿一般pending的巢狀統計與page.fail語意，不默默覆蓋一般印字。
	if ov.pend != nil {
		ov.nested++
		if ov.pend.page != nil {
			ov.pend.page.fail("巢狀", step)
		}
		log.Print("側欄職業：捨棄巢狀未返回印字")
		ov.pend = nil
	}
	state := &pendingPartyClass{slot: slot, size: ov.partySize, entrySteps: step}
	copy(state.nameRow[:], before[(y-1)*320+x:(y-1)*320+x+42])
	glyphs := make([][]uint16, len(ent.G))
	for i, glyph := range ent.G {
		glyphs[i] = append([]uint16(nil), glyph...)
	}
	ov.pend = &pendingPrint{s: key, x0: x, y0: y, sp: sp, ret: ret, font: font, g: glyphs,
		before: append([]byte(nil), before...), party: state}
	return true
}

func finishPartyClass(p *pendingPrint, after []byte, pal [256][3]uint8, step uint64) *textEv {
	if p == nil || p.party == nil || p.ret != partyClassReturn || p.font != 0 || p.sp > 65530 || !partyClassMeasuredKey(p.s) ||
		partySlot(p.x0, p.y0, false) != p.party.slot || p.party.slot < 0 || step < p.party.entrySteps || step-p.party.entrySteps > partyReturnLimit ||
		len(p.before) != 64000 || len(after) != 64000 || !partyClassEntryValidForKey(p.s, entry{Dy: p.dy, G: p.g}, p.party.size) {
		return nil
	}
	x0, y0, x1, y1 := 320, 200, -1, -1
	for i, a := range after {
		if a == p.before[i] {
			continue
		}
		x, y := i%320, i/320
		if x < p.x0 || x >= p.x0+42 || y < p.y0 || y >= p.y0+5 || a != 132 {
			return nil
		}
		if x < x0 {
			x0 = x
		}
		if x > x1 {
			x1 = x
		}
		if y < y0 {
			y0 = y
		}
		if y > y1 {
			y1 = y
		}
	}
	if x1 < x0 {
		return nil
	}
	w, h := x1-x0+1, y1-y0+1
	e := &textEv{originRet: p.ret, originFont: p.font, s: p.s, x0: p.x0, y0: p.y0, bx: x0, by: y0, bw: w, bh: h,
		fg: 132, fgRGB: pal[132], g: p.g, party: &partyClassWitness{slot: p.party.slot, size: p.party.size, nameRow: p.party.nameRow}}
	e.before, e.after = make([]byte, w*h), make([]byte, w*h)
	for y := 0; y < h; y++ {
		start := (y0+y)*320 + x0
		copy(e.before[y*w:(y+1)*w], p.before[start:start+w])
		copy(e.after[y*w:(y+1)*w], after[start:start+w])
	}
	return e
}

func (e *textEv) partyNameMatches(screen []byte, ignore func(int, int) bool) bool {
	if e == nil {
		return false
	}
	if e.party == nil {
		return true
	}
	if len(screen) != 64000 || partySlot(e.x0, e.y0, false) != e.party.slot || e.party.slot < 0 {
		return false
	}
	for column, value := range e.party.nameRow {
		x, y := e.x0+column, e.y0-1
		if ignore != nil && ignore(x, y) {
			continue
		}
		if screen[y*320+x] != value {
			return false
		}
	}
	return true
}

func partyClassShift(e *textEv, k int) (dx, dy int) {
	if e == nil || e.party == nil {
		return 0, 0
	}
	return partyClassOffset(e.s, e.party.size, k)
}
