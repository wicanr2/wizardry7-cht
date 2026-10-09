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

func TestRelocatedRosterRetainsOriginalRowsAndRestoresOnlyCode(t *testing.T) {
	initial := append([]rosterRow(nil), rosterRows...)
	cases := []struct {
		source, code string
		y            int
	}{
		{"LYSAND    M-THI", "M-THI", 156},
		{"NOBAL     M-PRI", "M-PRI", 162},
		{"TREON     M-MAG", "M-MAG", 168},
		{"PENTAS    F-ALC", "F-ALC", 174},
	}
	for _, tc := range cases {
		for _, pair := range [][2]byte{{195, 180}, {1, 0}} {
			before := bytes.Repeat([]byte{pair[0]}, 64000)
			ov := &overlay{rosterEnabled: true, roster: map[string]rosterEntry{tc.code: {En: tc.code}}}
			if !ov.enterRoster(oracle.Far(0x0110, 0x0fcc), 8, 0, 0, 150, tc.y, append([]byte(tc.source), 0), before) || ov.pend == nil {
				t.Fatalf("relocated source missing: %s/%d", tc.code, tc.y)
			}
			after := append([]byte(nil), before...)
			after[tc.y*320+151], after[tc.y*320+212] = pair[1], pair[1]
			e := finishRoster(ov.pend, after, [256][3]uint8{})
			if e == nil || e.by != tc.y || e.bx != 210 || e.bw != 30 || e.bh != 6 || !e.rosterMatches(after, nil) {
				t.Fatal("relocated event geometry/identity rejected")
			}
			if len(e.before) != 180 || len(e.roster.after) != 540 {
				t.Fatal("name included in restoration or row identity missing")
			}
			after[tc.y*320+151] = pair[0]
			if e.rosterMatches(after, nil) {
				t.Fatal("changed relocated name retained stale code")
			}
		}
	}
	for i, row := range rosterRows {
		if row != initial[i] {
			t.Fatal("relocation changed the original mapping table")
		}
		if got := rosterRowAt(row.source, row.y); got == nil || *got != row {
			t.Fatal("original mapping changed")
		}
	}
}

func TestRelocatedRosterRejectsUnmeasuredMappingAndRetiresOldEntry(t *testing.T) {
	before := bytes.Repeat([]byte{195}, 64000)
	for _, tc := range []struct {
		source string
		y      int
		sp     uint16
	}{
		{"LYSAND    M-THI", 150, 8},
		{"PENTAS    F-ALC", 156, 8},
		{"TEMPEST   F-FIG", 162, 8},
		{"UNKNOWN   M-THI", 156, 8},
		{"LYSAND    M-THI", 156, 65531},
	} {
		old := &textEv{y0: tc.y, roster: &rosterWitness{}}
		ov := &overlay{rosterEnabled: true, evs: []*textEv{old}, pend: &pendingPrint{}}
		if !ov.enterRoster(oracle.Far(0x0110, 0x0fcc), tc.sp, 0, 0, 150, tc.y, append([]byte(tc.source), 0), before) {
			t.Fatal("known row rejection was not handled")
		}
		if !old.dead || ov.pend != nil || ov.nested != 1 {
			t.Fatal("rejected source retained stale entry/pending")
		}
	}
}
