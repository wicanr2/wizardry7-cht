package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

func TestRosterWholeRowIdentityAndSelectionBorder(t *testing.T) {
	before := bytes.Repeat([]byte{195}, 64000)
	after := append([]byte(nil), before...)
	// Two independent source marks: name and code. The sixth row is blank.
	after[174*320+151], after[174*320+212] = 180, 180
	p := &pendingPrint{s: "TREON     M-MAG", y0: 174, before: before, roster: &rosterRows[4]}
	e := finishRoster(p, after, [256][3]uint8{})
	if e == nil || e.bx != 210 || e.bw != 30 || e.bh != 6 || !e.rosterMatches(after, nil) {
		t.Fatal("normal roster event missing")
	}
	changed := append([]byte(nil), after...)
	changed[174*320+151] = 195
	if e.rosterMatches(changed, nil) {
		t.Fatal("name changed while identical code retained old translation")
	}
	for x := 150; x < 240; x++ {
		after[179*320+x] = 1
	}
	if !e.rosterMatches(after, nil) {
		t.Fatal("measured adjacent-selection white border rejected")
	}
	after[179*320+200] = 0
	if e.rosterMatches(after, nil) {
		t.Fatal("unmeasured border ink accepted")
	}
	after[10] = 180
	if finishRoster(p, after, [256][3]uint8{}) != nil {
		t.Fatal("outside print change accepted")
	}
}

func TestRosterBlackSelectedInkAndUnknownEntryRetirement(t *testing.T) {
	before := bytes.Repeat([]byte{1}, 64000)
	after := append([]byte(nil), before...)
	after[150*320+150], after[150*320+210] = 0, 0
	p := &pendingPrint{s: "THESUS    M-FIG", y0: 150, before: before, roster: &rosterRows[0]}
	e := finishRoster(p, after, [256][3]uint8{0: {7, 8, 9}})
	if e == nil || e.fg != 0 || e.fgRGB != [3]uint8{7, 8, 9} || !e.rosterMatches(after, nil) {
		t.Fatal("black selected ink missing")
	}
	ov := &overlay{rosterEnabled: true, evs: []*textEv{e}}
	ret := oracle.Addr{Seg: 0x0110, Off: 0x0fcc}
	if !ov.enterRoster(ret, 8, 0, 0, 150, 150, []byte("UNKNOWN   M-FIG\x00"), before) || !e.dead || ov.pend != nil {
		t.Fatal("unknown full source retained previous translation")
	}
	if ov.enterRoster(oracle.Addr{Seg: 0x0110, Off: 0x1234}, 8, 0, 0, 150, 150, nil, before) {
		t.Fatal("unknown caller intercepted")
	}
}

func TestRosterPayloadRejectsPartialWrongFontAndTrailingData(t *testing.T) {
	g := make([]uint16, 16)
	g[1] = 1
	data := rosterData{Schema: 1, Size: 16, Context: "character.add.roster", Source: "runtime:0110:0fcc",
		SourceSHA256: "0000000000000000000000000000000000000000000000000000000000000000"}
	for _, row := range rosterRows {
		data.Entries = append(data.Entries, rosterEntry{En: row.code, Zh: "男", G: [][]uint16{g}})
	}
	raw, _ := json.Marshal(data)
	if got, err := decodeRoster(bytes.NewReader(raw), 16); err != nil || len(got) != 6 {
		t.Fatal("complete resource rejected", err)
	}
	if _, err := decodeRoster(bytes.NewReader(raw), 15); err == nil {
		t.Fatal("wrong font accepted")
	}
	if _, err := decodeRoster(bytes.NewReader(append(raw, []byte("{}")...)), 16); err == nil {
		t.Fatal("trailing JSON accepted")
	}
	data.Entries[5].En = "M-FIG"
	raw, _ = json.Marshal(data)
	if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
		t.Fatal("duplicate replacing missing key accepted")
	}
	data.Entries[5].En = "F-ALC"
	g[15] = 1
	raw, _ = json.Marshal(data)
	if _, err := decodeRoster(bytes.NewReader(raw), 16); err == nil {
		t.Fatal("nonblank native bottom row accepted")
	}
}
