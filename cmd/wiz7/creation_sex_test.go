package main

import (
	"bytes"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

// 合成記憶體只驗守門及資料所有權，不宣稱原版畫面parity。
type creationFake struct {
	*heldFake
	mode uint8
}

func (f *creationFake) VideoMode() uint8 { return f.mode }

func creationFixture(size int) (*creationFake, *overlay, []byte, []byte) {
	f := &creationFake{heldFake: &heldFake{mem: map[int][]byte{}, regs: oracle.Regs{CS: 0x0110, IP: 0x3877, SS: 0x6000, SP: 100}, step: 100}, mode: 0x13}
	f.put(oracle.Far(0x6000, 100), heldWords(0x9595, 'M', 0))
	f.put(oracle.Far(0x0110, 0x9592), []byte{0xe8, 0xe2, 0xa2})
	f.put(oracle.Far(0x0110, 0x2de7), heldWords(0x0163, 0x5000))
	f.put(oracle.Far(0x5000, 0x0163), []byte{0xe8, 0x63, 0x13, 0xcb})
	f.put(oracle.Far(0x5000, 0x0295), heldWords(90, 6))
	f.put(oracle.Far(0x5000, 0x01e5), heldWords(0x7000))
	before := texturedSexFrame()
	for _, point := range worldSexBackground {
		before[point[1]*320+point[0]] = byte(point[2])
	}
	after := append([]byte(nil), before...)
	for i := 0; i < 13; i++ {
		after[(6+i/6)*320+90+i%6] = 1
	}
	f.put(oracle.Far(0x7000, 0), before)
	d := sexFixture(size)
	ov := &overlay{sexEntries: map[string]rosterEntry{"M": d.Entries[0], "F": d.Entries[1]}}
	return f, ov, before, after
}
func creationReturnState(f *creationFake, after []byte) {
	f.regs.IP, f.regs.SP, f.step = 0x9595, 102, 492
	f.put(oracle.Far(0x5000, 0x0295), heldWords(96, 6))
	f.put(oracle.Far(0x7000, 0), after)
}

func TestCreationSexRejectsAliasAndInvalidStackBeforeRead(t *testing.T) {
	prior := flowK
	flowK = 3
	defer func() { flowK = prior }()
	for _, kind := range []string{"alias", "SP", "physical"} {
		t.Run(kind, func(t *testing.T) {
			f, ov, _, _ := creationFixture(16)
			old := &pendingProfileSex{ret: creationSexReturn}
			ov.sexPending = old
			switch kind {
			case "alias":
				f.regs.CS, f.regs.IP = 0x0111, 0x3867 // 相同linear，錯誤CS:IP。
			case "SP":
				f.regs.SP = 65531
			case "physical":
				f.regs.SS = 0xa000
			}
			if !ov.enterCreationSex(f) || f.reads != 0 {
				t.Fatal("invalid ABI caused a memory read")
			}
			if kind == "alias" && ov.sexPending != old {
				t.Fatal("alias consumed pending")
			}
			if kind != "alias" && ov.sexPending != nil {
				t.Fatal("true reentry left old creation pending")
			}
		})
	}
}

func TestCreationSexRejectsEntryTupleResourcesAndSources(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	cases := []string{"ret", "F", "font", "x", "y", "scale", "missingM", "missingF", "dy", "height", "bits", "empty", "signature", "driver_signature", "driver_pointer", "zero_buffer", "physical_buffer", "short_frame", "mode", "generic", "compound", "other_pending"}
	for _, kind := range cases {
		t.Run(kind, func(t *testing.T) {
			flowK = 2
			f, ov, _, _ := creationFixture(15)
			legacy := &pendingProfileSex{ret: profileSexReturn}
			switch kind {
			case "ret":
				f.put(oracle.Far(0x6000, 100), heldWords(0x9596, 'M', 0))
			case "F":
				f.put(oracle.Far(0x6000, 100), heldWords(0x9595, 'F', 0))
			case "font":
				f.put(oracle.Far(0x6000, 100), heldWords(0x9595, 'M', 1))
			case "x":
				f.put(oracle.Far(0x5000, 0x0295), heldWords(91, 6))
			case "y":
				f.put(oracle.Far(0x5000, 0x0295), heldWords(90, 7))
			case "scale":
				flowK = 4
			case "missingM":
				delete(ov.sexEntries, "M")
			case "missingF":
				delete(ov.sexEntries, "F")
			case "dy":
				v := ov.sexEntries["M"]
				v.Dy = 1
				ov.sexEntries["M"] = v
			case "height":
				v := ov.sexEntries["F"]
				v.G = [][]uint16{{1}}
				ov.sexEntries["F"] = v
			case "bits":
				ov.sexEntries["M"].G[0][0] = 1 << 15
			case "empty":
				ov.sexEntries["M"].G[0][0] = 0
			case "signature":
				f.put(oracle.Far(0x0110, 0x9592), []byte{0xe8, 0xe2, 0xa3})
			case "driver_signature":
				f.put(oracle.Far(0x5000, 0x0163), []byte{0xe8, 0x63, 0x12, 0xcb})
			case "driver_pointer":
				f.put(oracle.Far(0x0110, 0x2de7), heldWords(0x0164, 0x5000))
			case "zero_buffer":
				f.put(oracle.Far(0x5000, 0x01e5), heldWords(0))
			case "physical_buffer":
				f.put(oracle.Far(0x5000, 0x01e5), heldWords(0x9fff))
			case "short_frame":
				f.put(oracle.Far(0x7000, 0), []byte{1})
			case "mode":
				f.mode = 3
			case "generic":
				ov.pend = &pendingPrint{}
			case "compound":
				ov.compound = &compoundPair{}
			case "other_pending":
				ov.sexPending = legacy
			}
			ov.enterCreationSex(f)
			if kind == "other_pending" {
				if ov.sexPending != legacy {
					t.Fatal("legacy pending consumed")
				}
			} else if ov.sexPending != nil {
				t.Fatal("bad entry accepted")
			}
			if f.end > 0xa0000 {
				t.Fatal("physical boundary read")
			}
			if (kind == "F" || kind == "font" || kind == "ret") && f.reads != 1 {
				t.Fatal("tuple rejected after extra reads")
			}
		})
	}
	for _, point := range worldSexBackground {
		flowK = 3
		f, ov, before, _ := creationFixture(16)
		before[point[1]*320+point[0]] ^= 1
		f.put(oracle.Far(0x7000, 0), before)
		ov.enterCreationSex(f)
		if ov.sexPending != nil {
			t.Fatal("bad background accepted", point)
		}
	}
}

func TestCreationSexReentryNULAndLegacyHandoff(t *testing.T) {
	prior := flowK
	flowK = 3
	defer func() { flowK = prior }()
	for _, kind := range []string{"nested", "character", "NUL", "review", "world"} {
		t.Run(kind, func(t *testing.T) {
			f, ov, _, _ := creationFixture(15)
			ov.enterCreationSex(f)
			if ov.sexPending == nil {
				t.Fatal("initial entry")
			}
			old := &textEv{sex: true, originRet: creationSexReturn}
			legacy := &textEv{sex: true, originRet: worldSexReturn}
			ov.evs = []*textEv{old, legacy}
			switch kind {
			case "NUL":
				ov.clearProfileSex("一般字串插入")
			case "character":
				f.put(oracle.Far(0x6000, 100), heldWords(0x1234, 'X', 0))
				ov.enterCreationSex(f)
			case "review", "world":
				ret := profileSexReturn
				if kind == "world" {
					ret = worldSexReturn
				}
				f.put(oracle.Far(0x6000, 100), heldWords(ret.Off, 'M', 0))
				if ov.enterCreationSex(f) {
					t.Fatal("legacy branch swallowed")
				}
			default:
				ov.enterCreationSex(f)
			}
			if ov.sexPending != nil || !old.dead || legacy.dead {
				t.Fatal("reentry cleanup crossed source boundary")
			}
		})
	}
}

func TestCreationSexReturnRejectsIdentityPixelsAndDoesNotRevive(t *testing.T) {
	prior := flowK
	flowK = 3
	defer func() { flowK = prior }()
	for _, kind := range []string{"SS", "SP", "backward", "timeout", "driver", "buffer", "x", "y", "signature", "short", "outside", "color", "empty"} {
		t.Run(kind, func(t *testing.T) {
			f, ov, before, after := creationFixture(16)
			ov.enterCreationSex(f)
			if ov.sexPending == nil {
				t.Fatal("entry")
			}
			old := &textEv{sex: true, originRet: creationSexReturn}
			ov.evs = []*textEv{old}
			creationReturnState(f, after)
			switch kind {
			case "SS":
				f.regs.SS++
			case "SP":
				f.regs.SP++
			case "backward":
				f.step = 99
			case "timeout":
				f.step = 1000101
			case "driver":
				f.put(oracle.Far(0x0110, 0x2de7), heldWords(0x0163, 0x5100))
				f.put(oracle.Far(0x5100, 0x0163), []byte{0xe8, 0x63, 0x13, 0xcb})
				f.put(oracle.Far(0x5100, 0x0295), heldWords(96, 6))
				f.put(oracle.Far(0x5100, 0x01e5), heldWords(0x7000))
			case "buffer":
				f.put(oracle.Far(0x5000, 0x01e5), heldWords(0x8000))
				f.put(oracle.Far(0x8000, 0), after)
			case "x":
				f.put(oracle.Far(0x5000, 0x0295), heldWords(97, 6))
			case "y":
				f.put(oracle.Far(0x5000, 0x0295), heldWords(96, 7))
			case "signature":
				f.put(oracle.Far(0x5000, 0x0163), []byte{0xe8, 0x63, 0x12, 0xcb})
			case "short":
				f.put(oracle.Far(0x7000, 0), after[:10])
			case "outside":
				after[4*320+88] ^= 1
				f.put(oracle.Far(0x7000, 0), after)
			case "color":
				after[6*320+90] = 2
				f.put(oracle.Far(0x7000, 0), after)
			case "empty":
				f.put(oracle.Far(0x7000, 0), before)
			}
			ov.returnCreationSex(f)
			if ov.sexPending != nil || !old.dead || ov.created != 0 {
				t.Fatal("bad return left state")
			}
			creationReturnState(f, after)
			ov.returnCreationSex(f)
			if !old.dead || ov.created != 0 {
				t.Fatal("old event revived after rejected return")
			}
		})
	}
	f, ov, _, after := creationFixture(15)
	ov.enterCreationSex(f)
	pending := ov.sexPending
	creationReturnState(f, after)
	f.regs.CS, f.regs.IP = 0x0111, 0x9585
	f.reads = 0
	ov.returnCreationSex(f)
	if f.reads != 0 || ov.sexPending != pending {
		t.Fatal("return alias read or consumed pending")
	}
	f.regs.CS, f.regs.IP = 0x0110, profileSexReturn.Off
	ov.returnCreationSex(f)
	if ov.sexPending != pending {
		t.Fatal("other return consumed creation source")
	}
}

func TestCreationSexOwnedEventNativeGeometryAndClone(t *testing.T) {
	prior := flowK
	defer func() { flowK = prior }()
	for _, k := range []int{2, 3} {
		for _, size := range []int{15, 16} {
			flowK = k
			f, ov, before, after := creationFixture(size)
			original := append([]byte(nil), before...)
			ov.enterCreationSex(f)
			if ov.sexPending == nil || ov.sexPending.ss != 0x6000 || ov.sexPending.entryStep != 100 || ov.sexPending.driver != 0x5000 || ov.sexPending.buffer != 0x7000 {
				t.Fatal("source identity missing")
			}
			f.mem[0x70000][0] ^= 1
			ov.sexEntries["M"].G[0][0] = 2
			if ov.sexPending.before[0] != original[0] || ov.sexPending.g[0][0] != 1 {
				t.Fatal("pending aliases source")
			}
			creationReturnState(f, after)
			ov.returnCreationSex(f)
			if ov.sexPending != nil || len(ov.evs) != 1 || ov.created != 1 {
				t.Fatal("valid return")
			}
			e := ov.evs[0]
			if e.originRet != creationSexReturn || !e.sex || e.bx != 88 || e.by != 4 || e.bw != 8 || e.bh != 8 || len(e.before) != 64 || len(e.after) != 64 || !e.rectMatches(after, nil) {
				t.Fatal("wrong creation guard field")
			}
			shift := worldSexShift(e, k)
			x, y := e.x0*k-shift, e.y0*k-shift
			wantX, wantY := 270, 18
			if k == 2 {
				wantX, wantY = 192-size, 24-size
			}
			if x != wantX || y != wantY || x < 88*k || y < 4*k || x+size > 96*k || y+size > 12*k {
				t.Fatal("native glyph clipped or hyphen overlap")
			}
			clone := cloneVisualEvents([]*textEv{e})[0]
			if clone.originRet != creationSexReturn || clone.bw != 8 || !bytes.Equal(clone.before, e.before) || !bytes.Equal(clone.after, e.after) {
				t.Fatal("clone lost source field")
			}
			e.before[0] ^= 1
			e.after[0] ^= 1
			e.g[0][0] = 4
			if clone.before[0] == e.before[0] || clone.after[0] == e.after[0] || clone.g[0][0] != 1 {
				t.Fatal("clone aliases owner")
			}
			changed := append([]byte(nil), after...)
			changed[4*320+88] ^= 1
			if clone.rectMatches(changed, nil) {
				t.Fatal("extended background rewrite accepted")
			}
			ov.clearCreationSex()
			creationReturnState(f, after)
			ov.returnCreationSex(f)
			if !e.dead || ov.created != 1 {
				t.Fatal("withdrawn event revived")
			}
		}
	}
}
