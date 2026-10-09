package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

func partyTestEntry(size int) entry {
	ent := entry{}
	for range 4 {
		glyph := make([]uint16, size)
		glyph[12] = 1 << uint(size-1)
		ent.G = append(ent.G, glyph)
	}
	return ent
}
func partyTestRaw(s string) []byte { b := make([]byte, 64); copy(b, s); return b }
func partyTestOverlay(size, scale int) *overlay {
	return &overlay{partyEnabled: true, partySize: size, partyScale: scale}
}

func TestPartyClassConfigAndGlyphBounds(t *testing.T) {
	for _, size := range []int{15, 16} {
		for _, scale := range []int{2, 3} {
			trans := map[string]entry{partyClassKey: partyTestEntry(size)}
			if ok, got := partyClassConfig(trans, scale, scale); !ok || got != size {
				t.Fatal("valid configuration rejected")
			}
			if ok, _ := partyClassConfig(trans, scale, 5-scale); ok {
				t.Fatal("mixed HD/window scales accepted")
			}
			trans["synthetic"] = partyTestEntry(31 - size)
			if ok, _ := partyClassConfig(trans, scale, scale); ok {
				t.Fatal("mixed native heights accepted")
			}
		}
	}
	if ok, _ := partyClassConfig(nil, 2, 2); ok {
		t.Fatal("missing translation enabled")
	}
	for _, mutate := range []func(*entry){
		func(e *entry) { e.Dy = 1 }, func(e *entry) { e.D = []int{0} },
		func(e *entry) { e.G = e.G[:3] }, func(e *entry) { e.G[0] = make([]uint16, 16) },
		func(e *entry) { e.G[0] = make([]uint16, 15) },
		func(e *entry) { e.G[0][0] = 0x8000 },
	} {
		e := partyTestEntry(15)
		mutate(&e)
		if partyClassEntryValid(e, 15) {
			t.Fatal("invalid class resource accepted")
		}
	}
}

func TestPartyClassResourceCannotUseGenericNormalizedLookup(t *testing.T) {
	trans := map[string]entry{partyClassKey: partyTestEntry(15), "ALCHEMIST": partyTestEntry(15)}
	enabled, size := partyClassConfig(trans, 3, 3)
	reserved := takePartyClassEntry(trans)
	if !enabled || size != 15 || !partyClassEntryValid(reserved, 15) || len(trans) != 1 || len(trans["ALCHEMIST"].G) != 4 {
		t.Fatal("dedicated resource extraction changed unrelated translations")
	}
	for _, text := range []string{"ALCHEMI", "\x01ALCHEMI", "\tALCHEMI", "\nALCHEMI"} {
		key, _ := normalizeKey(text)
		if key != partyClassKey {
			t.Fatal("fixture no longer reaches normalized reserved key")
		}
		if _, ok := trans[key]; ok {
			t.Fatal("reserved short code available to generic lookup")
		}
	}
	ov := partyTestOverlay(15, 3)
	if !ov.enterPartyClass(partyClassReturn, 100, 0, 200, 9, 56, partyTestRaw(partyClassKey), 10, make([]byte, 64000), reserved) || ov.pend == nil {
		t.Fatal("resource extraction disabled exact dedicated source")
	}
	for _, text := range []string{"\x01ALCHEMI", "\tALCHEMI", "\nALCHEMI"} {
		ov.pend = nil
		ov.enterPartyClass(partyClassReturn, 100, 0, 200, 9, 56, partyTestRaw(text), 10, make([]byte, 64000), reserved)
		if ov.pend != nil {
			t.Fatal("near-match entered dedicated source")
		}
	}
}

func TestPartyVisualSnapshotPreservesIndependentShiftMetadata(t *testing.T) {
	e := &textEv{x0: 9, y0: 56, g: partyTestEntry(16).G, party: &partyClassWitness{slot: 0, size: 16}}
	e.party.nameRow[13] = 1
	visual := cloneVisualEvents([]*textEv{e})[0]
	dx, dy := partyClassShift(visual, 2)
	if dx != 2 || dy != -1 || visual.party == e.party {
		t.Fatal("visual publication lost or aliased party position metadata")
	}
	e.party.size, e.party.nameRow[13] = 15, 180
	if visual.party.size != 16 || visual.party.nameRow[13] != 1 {
		t.Fatal("published party metadata still depends on owner mutation")
	}
	if dx, dy := partyClassShift(visual, 3); dx != 0 || dy != 0 {
		t.Fatal("three-times position changed")
	}
}

func TestPartyNameReadBoundsAndSourceRejection(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, mutate := range []func(*uint16, *uint16, *uint16, []byte){
			func(font, sp, off *uint16, b []byte) { *font = 1 },
			func(font, sp, off *uint16, b []byte) { *sp = 65531 },
			func(font, sp, off *uint16, b []byte) { *off = 65505 },
			func(font, sp, off *uint16, b []byte) { b[0] = 0 },
			func(font, sp, off *uint16, b []byte) { b[0] = 31 },
			func(font, sp, off *uint16, b []byte) { b[0] = 127 },
			func(font, sp, off *uint16, b []byte) {
				for i := 0; i < 11; i++ {
					b[i] = 'X'
				}
			},
		} {
			ov := partyTestOverlay(16, 2)
			ov.partyEnabled = enabled
			old := &textEv{party: &partyClassWitness{slot: 0}}
			ov.evs = []*textEv{old}
			ov.partyNames[0].valid = true
			font, sp, off := uint16(0), uint16(100), uint16(200)
			raw := partyTestRaw("SYNTH")
			mutate(&font, &sp, &off, raw)
			if !ov.enterPartyName(partyNameReturn, sp, font, off, 9, 50, raw, 10, func(oracle.Addr, int) []byte { t.Fatal("invalid name read RAM"); return nil }) {
				t.Fatal("name fell through to translation")
			}
			if !old.dead || ov.partyNames[0].valid || ov.partyPendingNames[0] != nil {
				t.Fatal("stale slot state survived")
			}
		}
	}
	for _, cfg := range [][2]int{{15, 2}, {15, 3}, {16, 3}} {
		ov := partyTestOverlay(cfg[0], cfg[1])
		ov.enterPartyName(partyNameReturn, 100, 0, 200, 9, 50, partyTestRaw("SYNTH"), 10, func(oracle.Addr, int) []byte { t.Fatal("unneeded font read"); return nil })
	}
	if partyRAMSpan(oracle.Far(0x9fc0, 0), 1040) || partyRAMSpan(oracle.Far(0x0110, 65535), 2) || !partyRAMSpan(oracle.Far(0x9fbf, 0), 1040) {
		t.Fatal("RAM/end-span boundary")
	}
	for _, test := range []struct {
		pointer []byte
		font    []byte
		want    int
	}{
		{[]byte{0}, nil, 1}, {[]byte{0xc0, 0x9f}, nil, 1},
		{[]byte{0, 1}, make([]byte, 1039), 2}, {[]byte{0, 1}, make([]byte, 1040), 2},
	} {
		calls := 0
		bytesRead := 0
		state := partyNameMask([]byte("SYNTH"), func(addr oracle.Addr, n int) []byte {
			calls++
			bytesRead += n
			if calls == 1 {
				if addr != oracle.Far(0x0110, 0x364a) || n != 2 {
					t.Fatal("wrong pointer read")
				}
				return test.pointer
			}
			if n != 1040 || addr.Off != 0 {
				t.Fatal("wrong bounded font read")
			}
			return test.font
		})
		if state.valid || calls != test.want || bytesRead > 1042 {
			t.Fatal("wrong font SHA/short span accepted or unbounded read")
		}
	}
}

func TestPartyNameOriginalFontOptional(t *testing.T) {
	path := os.Getenv("WIZ7_TEST_PARTY_FONT")
	if path == "" {
		t.Skip("optional original font path not supplied")
	}
	if filepath.Base(path) != "VBFONT0.VGA" {
		t.Fatal("explicit path must name VBFONT0.VGA")
	}
	font, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skip("optional original font absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	state := partyNameMask([]byte(","), func(addr oracle.Addr, n int) []byte {
		calls++
		if calls == 1 {
			if n != 2 {
				t.Fatal("pointer size")
			}
			return []byte{0x0e, 0x7b}
		}
		if addr != oracle.Far(0x7b0e, 0) || n != 1040 {
			t.Fatal("font span")
		}
		return font
	})
	if !state.valid || calls != 2 {
		t.Fatal("confirmed font did not verify")
	}
	for col, on := range state.mask {
		if on != (col == 1) {
			t.Fatal("MSB-first comma bottom row changed")
		}
	}
}

func TestPartyNameReturnCommitsOnlyMatchingFrame(t *testing.T) {
	for _, test := range []struct {
		sp    uint16
		step  uint64
		valid bool
	}{{102, 20, true}, {101, 20, false}, {102, 9, false}, {102, 1000011, false}} {
		ov := partyTestOverlay(16, 2)
		state := partyNameState{valid: true}
		state.mask[3] = true
		ov.partyPendingNames[2] = &pendingPartyName{ret: partyNameReturn, sp: 100, entrySteps: 10, state: state}
		ov.returnPartyNames(partyClassReturn, 102, 20)
		if ov.partyPendingNames[2] == nil {
			t.Fatal("unrelated return consumed pending name")
		}
		ov.returnPartyNames(partyNameReturn, test.sp, test.step)
		if ov.partyNames[2].valid != test.valid || ov.partyPendingNames[2] != nil {
			t.Fatal("return guard did not control mask commit")
		}
	}
}

func TestPartyClassCollisionAndNativeInk(t *testing.T) {
	e := partyTestEntry(16)
	e.G[0][0] = 1 << 11 // native column4; shifted physical6 => DOS column3.
	if partyClassInkFits(e, 16, 2, partyNameState{}) {
		t.Fatal("missing name mask accepted")
	}
	name := partyNameState{valid: true}
	name.mask[3] = true
	if partyClassInkFits(e, 16, 2, name) {
		t.Fatal("name bottom-row collision accepted")
	}
	name.mask[3] = false
	if !partyClassInkFits(e, 16, 2, name) || !partyClassInkFits(e, 16, 3, partyNameState{}) || !partyClassInkFits(partyTestEntry(15), 15, 2, partyNameState{}) {
		t.Fatal("safe native ink rejected")
	}
	e.G[0][15] = 1
	if partyClassInkFits(e, 16, 2, name) {
		t.Fatal("late y+7 row intrusion accepted")
	}
}

func TestPartyClassExactSourceAndPendingLifecycle(t *testing.T) {
	for _, test := range []struct {
		ret           oracle.Addr
		sp, font, off uint16
		x, y          int
	}{
		{oracle.Far(0x0110, 0x1234), 100, 0, 200, 9, 56},
		{partyClassReturn, 65531, 0, 200, 9, 56}, {partyClassReturn, 100, 1, 200, 9, 56},
		{partyClassReturn, 100, 0, 65505, 9, 56}, {partyClassReturn, 100, 0, 200, 10, 56},
	} {
		ov := partyTestOverlay(15, 2)
		if !ov.enterPartyClass(test.ret, test.sp, test.font, test.off, test.x, test.y, partyTestRaw(partyClassKey), 10, make([]byte, 64000), partyTestEntry(15)) || ov.pend != nil {
			t.Fatal("invalid source fell through or translated")
		}
	}
	ov := partyTestOverlay(15, 2)
	old := &textEv{party: &partyClassWitness{slot: 0}}
	ov.evs = []*textEv{old}
	page := &flowPage{}
	ov.pend = &pendingPrint{page: page}
	ent := partyTestEntry(15)
	if !ov.enterPartyClass(partyClassReturn, 100, 0, 200, 9, 56, partyTestRaw(partyClassKey), 10, make([]byte, 64000), ent) {
		t.Fatal("exact class not handled")
	}
	if !old.dead || ov.pend == nil || ov.pend.party == nil || ov.nested != 1 || !page.failed {
		t.Fatal("nested/page lifecycle changed")
	}
	ent.G[0][12] = 0
	if ov.pend.g[0][12] == 0 {
		t.Fatal("pending glyph aliases mutable resource")
	}
	ov.enterPartyName(partyNameReturn, 100, 1, 200, 9, 50, partyTestRaw("SYNTH"), 11, nil)
	if ov.pend != nil {
		t.Fatal("new name kept previous slot class pending")
	}
}

func TestPartyClassFinishDeltaIdentityAndShift(t *testing.T) {
	ov := partyTestOverlay(15, 2)
	before := make([]byte, 64000)
	before[55*320+9] = 7
	ov.enterPartyClass(partyClassReturn, 100, 0, 200, 9, 56, partyTestRaw(partyClassKey), 10, before, partyTestEntry(15))
	pending := ov.pend
	after := append([]byte(nil), before...)
	after[56*320+9] = 132
	after[60*320+50] = 132
	var palette [256][3]uint8
	palette[132] = [3]uint8{4, 5, 6}
	e := finishPartyClass(pending, after, palette, 20)
	if e == nil || e.originRet != partyClassReturn || e.originFont != 0 || e.fg != 132 || e.fgRGB != palette[132] || e.bx != 9 || e.by != 56 || e.bw != 42 || e.bh != 5 {
		t.Fatal("class return witness invalid")
	}
	if !e.partyNameMatches(after, nil) {
		t.Fatal("name row identity rejected")
	}
	after[55*320+9] = 8
	if e.partyNameMatches(after, nil) || !e.partyNameMatches(after, func(x, y int) bool { return x == 9 && y == 55 }) {
		t.Fatal("cursor-only identity exception incorrect")
	}
	for _, mutate := range []func([]byte){func(a []byte) { a[56*320+51] = 132 }, func(a []byte) { a[61*320+9] = 132 }, func(a []byte) { a[56*320+9] = 131 }} {
		bad := append([]byte(nil), pending.before...)
		bad[56*320+9] = 132
		mutate(bad)
		if finishPartyClass(pending, bad, palette, 20) != nil {
			t.Fatal("outside/domain/foreground delta accepted")
		}
	}
	if finishPartyClass(pending, pending.before, palette, 20) != nil || finishPartyClass(pending, after, palette, 9) != nil || finishPartyClass(pending, after, palette, 1000011) != nil {
		t.Fatal("empty/timed return accepted")
	}
	if dx, dy := partyClassShift(e, 2); dx != 0 || dy != 0 {
		t.Fatal("Eten shifted")
	}
	e.party.size = 16
	if dx, dy := partyClassShift(e, 2); dx != 2 || dy != -1 {
		t.Fatal("Uni2 shift changed")
	}
	if dx, dy := partyClassShift(e, 3); dx != 0 || dy != 0 {
		t.Fatal("Uni3 shifted")
	}
	if dx, dy := partyClassShift(&textEv{}, 2); dx != 0 || dy != 0 {
		t.Fatal("ordinary event shifted")
	}
}
