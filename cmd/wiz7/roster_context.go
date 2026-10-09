package main

// R994 disposable prototype. This file only reads the bounded witnessed context.
import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/wicanr2/dosgolem/oracle"
)

type rosterRegisters struct{ CS, DS, SS, SP, BP uint16 }
type rosterAncestor struct{ BP, SavedBP, ReturnIP uint16 }
type rosterContextEvidence struct {
	registers               rosterRegisters
	frame                   [9]uint16 // Index, SavedBP, ReturnIP, Xbase, Ybase, Width, Count, First, Selected.
	ancestors               [3]rosterAncestor
	residentHash, arenaHash [32]byte
}
type rosterOwner uint8

const (
	rosterOwnerAdd rosterOwner = iota + 1
	rosterOwnerSelect
)

type rosterScope struct {
	owner                      rosterOwner
	first, count, xbase, ybase int
}
type rosterContext struct {
	scope       rosterScope
	index, slot int
	selected    bool
}

func rosterDigest(s string) [32]byte {
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		panic("invalid pinned roster code hash")
	}
	var out [32]byte
	copy(out[:], b)
	return out
}

var rosterResidentDigest = rosterDigest("34ad5b07cf59b4c779d5ba8e980b33978281040d590e9a7aedaf262cbe127e3f")
var rosterArenaDigests = map[rosterOwner][32]byte{
	rosterOwnerAdd:    rosterDigest("12eb7a9cf85d8423d5c7a021af2df23e04dd97cae3d4b4386dd02ccaeb2a511d"),
	rosterOwnerSelect: rosterDigest("32d7f3d24617df671f43648e69bb06b30fc9058f29bece4b78d3d697a7672d73"),
}

func rosterFrameBounds(r rosterRegisters) bool {
	return r.CS == 0x0110 && r.DS == 0x1106 && r.SS == r.DS &&
		r.SP <= 65530 && r.BP >= 2 && r.BP <= 65520 && int(r.BP) == int(r.SP)+14
}

func rosterRawValid(raw []byte) bool {
	if len(raw) != 16 || raw[15] != 0 {
		return false
	}
	for _, b := range raw[:10] {
		if b < 32 {
			return false
		}
	}
	return rosterAllowedCodes(28)[string(raw[10:15])]
}

func readRosterInvocationContext(r rosterRegisters, ret oracle.Addr, off, font uint16, x, y int, raw []byte, read func(oracle.Addr, int) []byte) (*rosterContextEvidence, error) {
	if ret != oracle.Far(0x0110, 0x0fcc) || !rosterFrameBounds(r) || off > 65520 || font != 0 || x != 150 || !rosterKnownY(y) || !rosterRawValid(raw) {
		return nil, fmt.Errorf("basic row signal")
	}
	return readRosterContext(r, read)
}

// Exactly 18+12 context bytes and 166+325 code bytes. No unbounded stack search.
func readRosterContext(r rosterRegisters, read func(oracle.Addr, int) []byte) (*rosterContextEvidence, error) {
	if !rosterFrameBounds(r) {
		return nil, fmt.Errorf("frame bounds")
	}
	e := &rosterContextEvidence{registers: r}
	b := read(oracle.Far(r.SS, r.BP-2), 18)
	if len(b) != 18 {
		return nil, fmt.Errorf("short frame")
	}
	for i := range e.frame {
		e.frame[i] = binary.LittleEndian.Uint16(b[i*2:])
	}
	parent, previous := e.frame[1], r.BP
	for i := range e.ancestors {
		if parent <= previous || parent > 65532 {
			return nil, fmt.Errorf("ancestor bounds")
		}
		b = read(oracle.Far(r.SS, parent), 4)
		if len(b) != 4 {
			return nil, fmt.Errorf("short ancestor")
		}
		e.ancestors[i] = rosterAncestor{parent, binary.LittleEndian.Uint16(b), binary.LittleEndian.Uint16(b[2:])}
		previous, parent = parent, e.ancestors[i].SavedBP
	}
	b = read(oracle.Far(0x0110, 0x0f38), 166)
	if len(b) != 166 {
		return nil, fmt.Errorf("short resident code")
	}
	e.residentHash = sha256.Sum256(b)
	b = read(oracle.Far(0x0110, 0xaf6f), 325)
	if len(b) != 325 {
		return nil, fmt.Errorf("short owner code")
	}
	e.arenaHash = sha256.Sum256(b)
	return e, nil
}

// Pure validation; numeric sums remain int until all 16-bit spans are checked.
func validateRosterContext(e *rosterContextEvidence, off, font uint16, x, y int) (rosterContext, error) {
	var out rosterContext
	if e == nil || !rosterFrameBounds(e.registers) || font != 0 || x != 150 || !rosterKnownY(y) {
		return out, fmt.Errorf("entry bounds")
	}
	f := e.frame
	index, count, first := int(f[0]), int(f[6]), int(f[7])
	if f[2] != 0x101b || f[5] != 15 || count < 1 || count > 6 || index >= count {
		return out, fmt.Errorf("frame values")
	}
	slot := first + index
	address := 0x0e38 + 20*slot
	if address > 65520 || address != int(off) || int(f[3])+2 != x || int(f[4])+6*index+1 != y {
		return out, fmt.Errorf("slot or coordinates")
	}
	parent, previous := f[1], e.registers.BP
	for _, a := range e.ancestors {
		if parent <= previous || parent > 65532 || a.BP != parent {
			return out, fmt.Errorf("ancestor chain")
		}
		previous, parent = parent, a.SavedBP
	}
	a, b, c := e.ancestors[0].ReturnIP, e.ancestors[1].ReturnIP, e.ancestors[2].ReturnIP
	var owner rosterOwner
	switch {
	case (a == 0x5775 || a == 0x5d24) && b == 0xb092 && c == 0xb0cb:
		owner = rosterOwnerAdd
	case (a == 0x5c2a || a == 0x61d9) && b == 0xba64 && c == 0xbe6c:
		owner = rosterOwnerSelect
	default:
		return out, fmt.Errorf("unknown owner")
	}
	if e.residentHash != rosterResidentDigest || e.arenaHash != rosterArenaDigests[owner] {
		return out, fmt.Errorf("code identity")
	}
	out = rosterContext{scope: rosterScope{owner, first, count, int(f[3]), int(f[4])}, index: index, slot: slot, selected: slot == int(f[8])}
	return out, nil
}
