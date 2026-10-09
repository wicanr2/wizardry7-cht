package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

func rosterFixture(owner rosterOwner, index, count, first, selected int) *rosterContextEvidence {
	e := &rosterContextEvidence{
		registers:    rosterRegisters{0x110, 0x1106, 0x1106, 63000, 63014},
		frame:        [9]uint16{uint16(index), 63100, 0x101b, 148, 149, 15, uint16(count), uint16(first), uint16(selected)},
		ancestors:    [3]rosterAncestor{{63100, 63200, 0x5775}, {63200, 63300, 0xb092}, {63300, 63400, 0xb0cb}},
		residentHash: rosterResidentDigest, arenaHash: rosterArenaDigests[owner],
	}
	if owner == rosterOwnerSelect {
		e.ancestors[0].ReturnIP, e.ancestors[1].ReturnIP, e.ancestors[2].ReturnIP = 0x5c2a, 0xba64, 0xbe6c
	}
	return e
}

func rosterFixtureResources(count int) map[string]rosterEntry {
	out := map[string]rosterEntry{}
	for key := range rosterAllowedCodes(count) {
		out[key] = rosterEntry{En: key}
	}
	return out
}

func rosterFixtureEnter(ov *overlay, e *rosterContextEvidence, raw []byte, before []byte) bool {
	return ov.enterRoster(oracle.Far(0x110, 0xfcc), e.registers.SP, 0,
		uint16(0xe38+20*(int(e.frame[7])+int(e.frame[0]))), 150,
		150+6*int(e.frame[0]), e, raw, before)
}

func TestRosterWholeRowIdentityAndSelectionBorder(t *testing.T) {
	before := bytes.Repeat([]byte{195}, 64000)
	after := append([]byte(nil), before...)
	after[174*320+151], after[174*320+212] = 180, 180
	ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(6)}
	rosterFixtureEnter(ov, rosterFixture(rosterOwnerAdd, 4, 6, 0, 65535), []byte("TENLETTERS"+"M-MAG\x00"), before)
	e := finishRoster(ov.pend, after, [256][3]uint8{})
	if e == nil || e.bx != 210 || e.bw != 30 || e.bh != 6 || !e.rosterMatches(after, nil) {
		t.Fatal("normal row missing")
	}
	changed := append([]byte(nil), after...)
	changed[174*320+151] = 195
	if e.rosterMatches(changed, nil) {
		t.Fatal("changed name retained old translation")
	}
	for x := 150; x < 240; x++ {
		after[179*320+x] = 1
	}
	if !e.rosterMatches(after, nil) {
		t.Fatal("measured adjacent-selection border rejected")
	}
	after[179*320+200] = 0
	if e.rosterMatches(after, nil) {
		t.Fatal("unmeasured border ink accepted")
	}
	after[10] = 180
	if finishRoster(ov.pend, after, [256][3]uint8{}) != nil {
		t.Fatal("outside print change accepted")
	}
}

func TestRosterBlackSelectedInkAndUnknownEntryRetirement(t *testing.T) {
	before := bytes.Repeat([]byte{1}, 64000)
	ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(6)}
	rosterFixtureEnter(ov, rosterFixture(rosterOwnerSelect, 0, 6, 0, 0), []byte("ANYNAME   M-FIG\x00"), before)
	after := append([]byte(nil), before...)
	after[150*320+150], after[150*320+210] = 0, 0
	e := finishRoster(ov.pend, after, [256][3]uint8{0: {7, 8, 9}})
	if e == nil || e.fg != 0 || e.fgRGB != [3]uint8{7, 8, 9} || !e.rosterMatches(after, nil) {
		t.Fatal("selected ink missing")
	}
	ov.evs = []*textEv{e}
	if !ov.enterRoster(oracle.Far(0x110, 0xfcc), 8, 0, 0, 150, 150, nil, nil, before) || !e.dead || ov.pend != nil {
		t.Fatal("bad context retained stale row")
	}
	if ov.enterRoster(oracle.Far(0x110, 0x1234), 8, 0, 0, 150, 150, nil, nil, before) {
		t.Fatal("unknown caller intercepted")
	}
}

func TestRosterNameDescendersUseWholeBottomRowIdentity(t *testing.T) {
	for _, selected := range []bool{false, true} {
		bg, fg, selection := byte(195), byte(180), 65535
		if selected {
			bg, fg, selection = 1, 0, 0
		}
		before := bytes.Repeat([]byte{bg}, 64000)
		after := append([]byte(nil), before...)
		// 原版comma末列bit64是字格column1，兩個逗號落姓名column13／31。
		after[155*320+163], after[155*320+181] = fg, fg
		after[150*320+210] = fg
		ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(6)}
		rosterFixtureEnter(ov, rosterFixture(rosterOwnerAdd, 0, 6, 0, selection), []byte("AB,CD,E   F-ALC\x00"), before)
		e := finishRoster(ov.pend, after, [256][3]uint8{})
		if e == nil || !e.rosterMatches(after, nil) {
			t.Fatal("original name descenders rejected")
		}
		for _, column := range []int{13, 20, 31} {
			changed := append([]byte(nil), after...)
			at := 155*320 + 150 + column
			changed[at] ^= 1
			if e.rosterMatches(changed, nil) {
				t.Fatal("changed bottom-row name ink accepted")
			}
			if !e.rosterMatches(changed, func(x, y int) bool { return x == 150+column && y == 155 }) {
				t.Fatal("cursor ignore changed")
			}
		}
		partial := append([]byte(nil), after...)
		partial[155*320+163] = bg
		if e.rosterMatches(partial, nil) {
			t.Fatal("partial descender deletion accepted")
		}
		for _, color := range []byte{195, 1} {
			solid := append([]byte(nil), after...)
			for x := 150; x < 240; x++ {
				solid[155*320+x] = color
			}
			if e.rosterMatches(solid, nil) == selected {
				t.Fatal("ordinary/selected solid-border contract changed")
			}
			if !selected {
				solid[155*320+200] = 196 - color
				if e.rosterMatches(solid, nil) {
					t.Fatal("mixed background border accepted")
				}
			}
		}
	}
}

func TestRosterPayloadRejectsPartialWrongFontAndTrailingData(t *testing.T) {
	for _, count := range []int{6, 28} {
		g := make([]uint16, 16)
		g[1] = 1
		data := rosterData{Schema: 1, Size: 16, Context: "character.add.roster", Source: "runtime:0110:0fcc", SourceSHA256: "0000000000000000000000000000000000000000000000000000000000000000"}
		for key := range rosterAllowedCodes(count) {
			data.Entries = append(data.Entries, rosterEntry{En: key, Zh: "男", G: [][]uint16{g}})
		}
		raw, _ := json.Marshal(data)
		if got, err := decodeRoster(bytes.NewReader(raw), 16); err != nil || len(got) != count {
			t.Fatal("complete resource rejected", err)
		}
		if _, err := decodeRoster(bytes.NewReader(raw), 15); err == nil {
			t.Fatal("wrong font accepted")
		}
		if _, err := decodeRoster(bytes.NewReader(append(raw, []byte("{}")...)), 16); err == nil {
			t.Fatal("trailing JSON accepted")
		}
		original := data.Entries[count-1].En
		data.Entries[count-1].En = data.Entries[0].En
		raw, _ = json.Marshal(data)
		if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
			t.Fatal("duplicate accepted")
		}
		data.Entries[count-1].En = "X-XYZ"
		raw, _ = json.Marshal(data)
		if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
			t.Fatal("same-count unknown replacement accepted")
		}
		data.Entries[count-1].En = original
		g[15] = 1
		raw, _ = json.Marshal(data)
		if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
			t.Fatal("nonblank native bottom row accepted")
		}
		g[15] = 0
		data.Entries = data.Entries[:count-1]
		raw, _ = json.Marshal(data)
		if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
			t.Fatal("partial set accepted")
		}
	}
}

func TestGenericRosterNamesPositionsAndCodeOnlyRestoration(t *testing.T) {
	names := [][]byte{[]byte("SEVENXX   "), []byte("TENLETTERS"), {0x82, 0xa1, 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'}}
	for _, owner := range []rosterOwner{rosterOwnerAdd, rosterOwnerSelect} {
		for count := 1; count <= 6; count++ {
			for index := 0; index < count; index++ {
				for _, selected := range []bool{false, true} {
					for _, name := range names {
						bg, fg, sel := byte(195), byte(180), 65535
						if selected {
							bg, fg, sel = 1, 0, index
						}
						before := bytes.Repeat([]byte{bg}, 64000)
						ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(28)}
						raw := append(append([]byte(nil), name...), []byte("M-NIN\x00")...)
						if !rosterFixtureEnter(ov, rosterFixture(owner, index, count, 0, sel), raw, before) || ov.pend == nil {
							t.Fatal("generic name rejected")
						}
						y := 150 + 6*index
						after := append([]byte(nil), before...)
						after[y*320+209], after[y*320+239] = fg, fg
						e := finishRoster(ov.pend, after, [256][3]uint8{})
						if e == nil || len(e.before) != 180 || len(e.roster.after) != 540 || !e.rosterMatches(after, nil) {
							t.Fatal("full ten-cell name or code-only restoration failed")
						}
						after[y*320+209] = bg
						if e.rosterMatches(after, nil) {
							t.Fatal("last name cell failed identity")
						}
					}
				}
			}
		}
	}
}

func TestRosterRejectedSourceRetiresOldEntryAndPending(t *testing.T) {
	before := bytes.Repeat([]byte{195}, 64000)
	for _, kind := range []string{"hash", "sp", "control", "suffix", "short", "nul", "font", "x", "off"} {
		e := rosterFixture(rosterOwnerAdd, 1, 6, 0, 65535)
		raw := []byte("ARBITRARY M-FIG\x00")
		sp, font, off, x := e.registers.SP, uint16(0), uint16(0xe4c), 150
		switch kind {
		case "hash":
			e.arenaHash[0] ^= 1
		case "sp":
			sp = 65531
		case "control":
			raw[0] = 0x1f
		case "suffix":
			raw[10] = 'X'
		case "short":
			raw = raw[:15]
		case "nul":
			raw[15] = 1
		case "font":
			font = 1
		case "x":
			x = 149
		case "off":
			off = 65521
		}
		old := &textEv{y0: 156, roster: &rosterWitness{}}
		ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(28), evs: []*textEv{old}, pend: &pendingPrint{}}
		if !ov.enterRoster(oracle.Far(0x110, 0xfcc), sp, font, off, x, 156, e, raw, before) || !old.dead || ov.pend != nil || ov.nested != 1 {
			t.Fatal("rejected source retained stale entry", kind)
		}
	}
}

func TestRosterScopeAndBlankNameDefense(t *testing.T) {
	before := bytes.Repeat([]byte{195}, 64000)
	for _, kind := range []string{"owner", "first", "count", "base"} {
		e := rosterFixture(rosterOwnerAdd, 0, 6, 0, 65535)
		context, _ := validateRosterContext(e, 0xe38, 0, 150, 150)
		oldScope := context.scope
		switch kind {
		case "owner":
			oldScope.owner = rosterOwnerSelect
		case "first":
			oldScope.first = 1
		case "count":
			oldScope.count = 5
		case "base":
			oldScope.ybase = 143
		}
		old := &textEv{y0: 180, roster: &rosterWitness{scope: oldScope, slot: 5}}
		ov := &overlay{rosterEnabled: true, roster: rosterFixtureResources(28), evs: []*textEv{old}}
		rosterFixtureEnter(ov, e, []byte("          M-FIG\x00"), before)
		if !old.dead || ov.pend == nil {
			t.Fatal("scope change failed retirement", kind)
		}
		after := append([]byte(nil), before...)
		after[150*320+210] = 180
		if finishRoster(ov.pend, after, [256][3]uint8{}) == nil {
			t.Fatal("blank-name defensive boundary rejected")
		}
		ov.pend.roster.blankName = false
		if finishRoster(ov.pend, after, [256][3]uint8{}) != nil {
			t.Fatal("nonblank name with no ink accepted")
		}
		ov.pend.roster.blankName = true
		ov.pend.roster.selected = true
		if finishRoster(ov.pend, after, [256][3]uint8{}) != nil {
			t.Fatal("selected frame accepted ordinary color")
		}
	}
}
