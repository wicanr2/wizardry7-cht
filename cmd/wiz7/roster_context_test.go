package main

import (
	"encoding/binary"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

func TestRosterContextOwnersOffsetAndSelectedDomain(t *testing.T) {
	for _, owner := range []rosterOwner{rosterOwnerAdd, rosterOwnerSelect} {
		for _, first := range []int{0, 7, 3093} {
			e := rosterFixture(owner, 1, 2, first, 65535)
			off := uint16(0xe38 + 20*(first+1))
			c, err := validateRosterContext(e, off, 0, 150, 156)
			if err != nil || c.slot != first+1 || c.selected {
				t.Fatal("bounded First or outside Selected rejected", err)
			}
			e.frame[8] = uint16(first + 1)
			c, err = validateRosterContext(e, off, 0, 150, 156)
			if err != nil || !c.selected {
				t.Fatal("absolute selection predicate failed")
			}
			if owner == rosterOwnerAdd {
				e.ancestors[0].ReturnIP = 0x5d24
			} else {
				e.ancestors[0].ReturnIP = 0x61d9
			}
			if _, err = validateRosterContext(e, off, 0, 150, 156); err != nil {
				t.Fatal("second owner tuple rejected", err)
			}
		}
	}
}

func TestRosterContextRejectsWrapIdentityAndInvalidFrames(t *testing.T) {
	for _, kind := range []string{"bp-low", "bp-high", "sp-wrap", "ss", "ds", "cs", "count-zero", "count-seven", "index", "width", "return", "first-wrap", "x-wrap", "y-wrap", "ancestor-loop", "ancestor-wrap", "ancestor-mismatch", "owner", "resident", "arena"} {
		e := rosterFixture(rosterOwnerAdd, 1, 6, 0, 65535)
		switch kind {
		case "bp-low":
			e.registers.BP = 1
		case "bp-high":
			e.registers.BP = 65521
		case "sp-wrap":
			e.registers.SP = 65530
			e.registers.BP = 8
		case "ss":
			e.registers.SS++
		case "ds":
			e.registers.DS++
		case "cs":
			e.registers.CS++
		case "count-zero":
			e.frame[6] = 0
		case "count-seven":
			e.frame[6] = 7
		case "index":
			e.frame[0] = 6
		case "width":
			e.frame[5] = 14
		case "return":
			e.frame[2]++
		case "first-wrap":
			e.frame[7] = 65535
		case "x-wrap":
			e.frame[3] = 65535
		case "y-wrap":
			e.frame[4] = 65535
		case "ancestor-loop":
			e.frame[1] = e.registers.BP
		case "ancestor-wrap":
			e.frame[1] = 65533
		case "ancestor-mismatch":
			e.ancestors[1].BP++
		case "owner":
			e.ancestors[2].ReturnIP++
		case "resident":
			e.residentHash[0] ^= 1
		case "arena":
			e.arenaHash = rosterArenaDigests[rosterOwnerSelect]
		}
		if _, err := validateRosterContext(e, 0xe4c, 0, 150, 156); err == nil {
			t.Fatal("invalid context accepted", kind)
		}
	}
}

func TestRosterReaderBoundedAndBasicFailureDoesNotRead(t *testing.T) {
	e := rosterFixture(rosterOwnerAdd, 0, 6, 0, 65535)
	for _, kind := range []string{"caller", "off", "font", "x", "y", "raw", "sp"} {
		r := e.registers
		ret := oracle.Far(0x110, 0xfcc)
		off, font, x, y := uint16(0xe38), uint16(0), 150, 150
		raw := []byte("ARBITRARY M-FIG\x00")
		switch kind {
		case "caller":
			ret.Off++
		case "off":
			off = 65521
		case "font":
			font = 1
		case "x":
			x = 149
		case "y":
			y = 151
		case "raw":
			raw[15] = 1
		case "sp":
			r.SP = 65531
		}
		reads := 0
		_, err := readRosterInvocationContext(r, ret, off, font, x, y, raw, func(oracle.Addr, int) []byte { reads++; return nil })
		if err == nil || reads != 0 {
			t.Fatal("basic rejection expanded read", kind, reads)
		}
	}
	reads, total := 0, 0
	read := func(a oracle.Addr, n int) []byte {
		reads++
		total += n
		b := make([]byte, n)
		if a == oracle.Far(e.registers.SS, e.registers.BP-2) {
			for i, v := range e.frame {
				binary.LittleEndian.PutUint16(b[i*2:], v)
			}
			return b
		}
		for _, p := range e.ancestors {
			if a == oracle.Far(e.registers.SS, p.BP) {
				binary.LittleEndian.PutUint16(b, p.SavedBP)
				binary.LittleEndian.PutUint16(b[2:], p.ReturnIP)
				return b
			}
		}
		if a != oracle.Far(0x110, 0xf38) && a != oracle.Far(0x110, 0xaf6f) {
			t.Fatal("unexpected read", a, n)
		}
		return b
	}
	got, err := readRosterContext(e.registers, read)
	if err != nil || got == nil || reads != 6 || total != 521 {
		t.Fatal("reader budget mismatch", reads, total, err)
	}
	reads = 0
	_, err = readRosterContext(e.registers, func(oracle.Addr, int) []byte { reads++; return nil })
	if err == nil || reads != 1 {
		t.Fatal("short read failed early stop")
	}
}
