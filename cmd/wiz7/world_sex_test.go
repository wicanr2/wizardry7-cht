package main

import (
	"bytes"
	"github.com/wicanr2/dosgolem/oracle"
	"os"
	"testing"
)

func worldSexFixture(t *testing.T, size int) ([]byte, []byte, *overlay) {
	t.Helper()
	root := os.Getenv("WIZ7_WORLD_SEX_FIXTURE_DIR")
	if root == "" {
		t.Skip("requires private dosgolem world sex receipt")
	}
	before, err := os.ReadFile(root + "/M-before.bin")
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(root + "/M-after.bin")
	if err != nil {
		t.Fatal(err)
	}
	data := sexFixture(size)
	return before, after, &overlay{sexEntries: map[string]rosterEntry{"M": data.Entries[0], "F": data.Entries[1]}}
}
func TestWorldSexObservedReturnAndNativeGeometry(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	for _, k := range []int{2, 3} {
		flowK = k
		for _, size := range []int{15, 16} {
			before, after, ov := worldSexFixture(t, size)
			ov.enterProfileSex(worldSexReturn, 65310, 0, 'M', 90, 6, before)
			if ov.sexPending == nil {
				t.Fatal("observed world entry rejected", k, size)
			}
			before[0] ^= 1
			if ov.sexPending.before[0] == before[0] {
				t.Fatal("pending did not own before")
			}
			ov.returnProfileSex(65312, 96, 6, after, [256][3]uint8{})
			if ov.sexPending != nil || len(ov.evs) != 1 {
				t.Fatal("world return rejected")
			}
			e := ov.evs[0]
			if e.originRet != worldSexReturn || !e.sex || e.bx != 88 || e.by != 4 || e.bw != 8 || e.bh != 8 || len(e.before) != 64 {
				t.Fatal("wrong guard domain")
			}
			if !e.rectMatches(after, nil) {
				t.Fatal("observed after mismatch")
			}
			altered := append([]byte(nil), after...)
			altered[4*320+88] ^= 1
			if e.rectMatches(altered, nil) {
				t.Fatal("extended background rewrite accepted")
			}
			clone := cloneVisualEvents(ov.evs)[0]
			shift := worldSexShift(clone, k)
			x, y := clone.x0*k-shift, clone.y0*k-shift
			if x < 88*k || y < 4*k || x+size > 96*k || y+size > 12*k {
				t.Fatal("native glyph outside header field")
			}
			if worldSexShift(&textEv{originRet: profileSexReturn, sex: true, g: e.g}, k) != 0 {
				t.Fatal("024 moved")
			}
			ov.returnProfileSex(65312, 96, 6, after, [256][3]uint8{})
			if len(ov.evs) != 1 {
				t.Fatal("return consumed twice")
			}
		}
	}
}
func TestWorldSexRejectsBadEntryAndEachBackgroundPoint(t *testing.T) {
	prior := flowK
	flowK = 2
	defer func() { flowK = prior }()
	for _, change := range []func(*overlay, *oracle.Addr, *uint16, *uint16, *int, *int, *[]byte){
		func(_ *overlay, r *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) { r.Off++ },
		func(_ *overlay, _ *oracle.Addr, f *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) { *f = 1 },
		func(_ *overlay, _ *oracle.Addr, _ *uint16, c *uint16, _ *int, _ *int, _ *[]byte) { *c = 'X' },
		func(_ *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, x *int, _ *int, _ *[]byte) { *x = 91 },
		func(_ *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, y *int, _ *[]byte) { *y = 7 },
		func(ov *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) {
			ov.sexEntries = nil
		},
		func(ov *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) {
			ov.pend = &pendingPrint{}
		},
		func(ov *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) {
			ov.compound = &compoundPair{}
		},
		func(ov *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, _ *[]byte) {
			ov.sexPending = &pendingProfileSex{}
		},
		func(_ *overlay, _ *oracle.Addr, _ *uint16, _ *uint16, _ *int, _ *int, b *[]byte) { *b = (*b)[:100] },
	} {
		before, _, ov := worldSexFixture(t, 16)
		ret, font, code, x, y := worldSexReturn, uint16(0), uint16('M'), 90, 6
		change(ov, &ret, &font, &code, &x, &y, &before)
		ov.enterProfileSex(ret, 65310, font, code, x, y, before)
		if ov.sexPending != nil {
			t.Fatal("malformed entry accepted")
		}
	}
	for _, p := range worldSexBackground {
		before, _, ov := worldSexFixture(t, 16)
		before[p[1]*320+p[0]] ^= 1
		ov.enterProfileSex(worldSexReturn, 65310, 0, 'M', 90, 6, before)
		if ov.sexPending != nil {
			t.Fatal("wrong background accepted", p)
		}
	}
}
func TestWorldSexRejectsBadReturnAndInterrupts(t *testing.T) {
	prior := flowK
	flowK = 3
	defer func() { flowK = prior }()
	for _, kind := range []string{"SP", "X", "Y", "empty", "short", "outside", "color"} {
		before, after, ov := worldSexFixture(t, 15)
		ov.enterProfileSex(worldSexReturn, 65310, 0, 'M', 90, 6, before)
		sp, x, y := uint16(65312), 96, 6
		switch kind {
		case "SP":
			sp++
		case "X":
			x++
		case "Y":
			y++
		case "empty":
			after = before
		case "short":
			after = after[:100]
		case "outside":
			after[4*320+88] ^= 1
		case "color":
			for i, a := range after {
				if a != before[i] {
					after[i] = 2
					break
				}
			}
		}
		ov.returnProfileSex(sp, x, y, after, [256][3]uint8{})
		if ov.sexPending != nil || len(ov.evs) != 0 {
			t.Fatal("bad return retained state", kind)
		}
	}
	before, _, ov := worldSexFixture(t, 15)
	copyOfBefore := append([]byte(nil), before...)
	ov.enterProfileSex(worldSexReturn, 65310, 0, 'F', 90, 6, before)
	if ov.sexPending == nil {
		t.Fatal("statically proven first-letter F rejected")
	}
	ov.clearProfileSex("NUL interrupted")
	if ov.sexPending != nil || !bytes.Equal(before, copyOfBefore) {
		t.Fatal("interruption changed original input")
	}
}

func TestWorldSexRejectsUnsupportedScaleAndStack(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	for _, k := range []int{0, 1, 4} {
		flowK = k
		before, _, ov := worldSexFixture(t, 15)
		ov.sexEnabled = true
		ov.enterProfileSex(worldSexReturn, 100, 0, 'M', 90, 6, before)
		if ov.sexPending != nil {
			t.Fatal("world branch inherited old enabled flag", k)
		}
	}
	flowK = 2
	before, _, ov := worldSexFixture(t, 15)
	ov.enterProfileSex(worldSexReturn, 65531, 0, 'M', 90, 6, before)
	if ov.sexPending != nil {
		t.Fatal("stack span overflow accepted")
	}
}
func TestWorldSexRejectionDropsWorldStateOnly(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	for _, kind := range []string{"scale", "stack", "resource"} {
		for _, pendingRet := range []oracle.Addr{worldSexReturn, profileSexReturn} {
			flowK = 2
			before, _, ov := worldSexFixture(t, 15)
			ov.sexEnabled = true
			world := &textEv{sex: true, originRet: worldSexReturn}
			review := &textEv{sex: true, originRet: profileSexReturn}
			ov.evs = []*textEv{world, review}
			pending := &pendingProfileSex{ret: pendingRet}
			ov.sexPending = pending
			sp := uint16(100)
			switch kind {
			case "scale":
				flowK = 4
			case "stack":
				sp = 65531
			case "resource":
				delete(ov.sexEntries, "F")
			}
			ov.enterProfileSex(worldSexReturn, sp, 0, 'M', 90, 6, before)
			if !world.dead || review.dead {
				t.Fatal("rejection must drop only world event")
			}
			if pendingRet == worldSexReturn && ov.sexPending != nil {
				t.Fatal("old world pending survived")
			}
			if pendingRet == profileSexReturn && ov.sexPending != pending {
				t.Fatal("024 pending consumed")
			}
		}
	}
}
func TestWorldSexRejectsPartialAndCorruptResources(t *testing.T) {
	prior := flowK
	flowK = 2
	defer func() { flowK = prior }()
	for _, change := range []func(map[string]rosterEntry){
		func(m map[string]rosterEntry) { delete(m, "F") },
		func(m map[string]rosterEntry) { v := m["M"]; v.En = "F"; m["M"] = v },
		func(m map[string]rosterEntry) { v := m["F"]; v.G = nil; m["F"] = v },
		func(m map[string]rosterEntry) { v := m["F"]; v.G = [][]uint16{make([]uint16, 15)}; m["F"] = v },
		func(m map[string]rosterEntry) { v := m["M"]; v.G[0][0] = 1 << 15; m["M"] = v },
		func(m map[string]rosterEntry) {
			v := m["F"]
			v.G = [][]uint16{make([]uint16, 16)}
			v.G[0][0] = 1
			m["F"] = v
		},
	} {
		before, _, ov := worldSexFixture(t, 15)
		change(ov.sexEntries)
		ov.enterProfileSex(worldSexReturn, 100, 0, 'M', 90, 6, before)
		if ov.sexPending != nil {
			t.Fatal("invalid resource accepted")
		}
	}
}
func TestWorldSexBothKeysCloneAndWaitingForBlit(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	for _, k := range []int{2, 3} {
		flowK = k
		for _, size := range []int{15, 16} {
			for _, code := range []uint16{'M', 'F'} {
				before, after, ov := worldSexFixture(t, size)
				// F uses a synthetic output boundary fixture; this is not normal F parity.
				ov.enterProfileSex(worldSexReturn, 100, 0, code, 90, 6, before)
				if ov.sexPending == nil {
					t.Fatal("valid entry missing")
				}
				if len(ov.evs) != 0 {
					t.Fatal("painted before return")
				}
				ov.returnProfileSex(102, 96, 6, after, [256][3]uint8{})
				e := ov.evs[0]
				if e.shown || e.rectMatches(before, nil) {
					t.Fatal("before blit incorrectly shown")
				}
				if !e.rectMatches(after, nil) {
					t.Fatal("after blit rejected")
				}
				changed := append([]byte(nil), after...)
				changed[4*320+88] ^= 1
				if e.rectMatches(changed, nil) || !e.rectMatches(changed, func(x, y int) bool { return x == 88 && y == 4 }) {
					t.Fatal("cursor/domain guards differ")
				}
				c := cloneVisualEvents([]*textEv{e})[0]
				e.originRet = oracle.Far(0, 0)
				e.g[0][0] = 9
				if c.originRet != worldSexReturn || !c.sex || c.g[0][0] != 1 {
					t.Fatal("visual clone shared owner")
				}
				s := worldSexShift(c, k)
				if c.x0*k-s < 88*k || c.x0*k-s+size > 96*k || c.y0*k-s < 4*k || c.y0*k-s+size > 12*k {
					t.Fatal("dash or role overlap")
				}
			}
		}
	}
}
