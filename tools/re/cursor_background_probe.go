package main

// P1–P6私人唯讀游標來源/背景觀測，不參與正式覆繪。
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/wicanr2/dosgolem/oracle"
)

const r1000DefA = "b8dc0350ff36200eff361e0e8b46ee050c005250e8d8a383c40ab8010050ff76"
const r1000DefB = "8946fa50ff76fcb80d0050b84c0050ff36200eff361e0ee8eee383c40ce89fdf"
const r1000Reset = "c4048b46fa034610508b46fc03460e50b800005050b80000ba00005250e827e5"

type r1000Candidate struct {
	serviceID, nulID int
	sourceDS         uint16
	ptr              oracle.Addr
	reset            bool
}
type r1000Registry struct {
	id, nulID, serviceID    int
	driver                  uint16
	ptr                     oracle.Addr
	picture                 []byte
	pictureFile, pictureSHA string
	start, deadline         uint64
	samples                 int
	reset                   bool
	resetNUL, resetService  int
}
type r1000CursorRow struct {
	Record, Registry, NULID, ServiceID, TriggerNULID int
	Publication, Step                                uint64
	Status                                           string
	Driver                                           uint16
	PicturePointer, BackgroundPointer                oracle.Addr
	X, Y, ClipW, ClipH, FullW, FullH                 int
	Busy                                             byte
	Samples                                          int
	LastRecord                                       bool
	Files, SHA256                                    map[string]string
}

var r1000Cursor struct {
	candidate            *r1000Candidate
	active               *r1000Registry
	definitions, records int
}

func r1000Enabled() bool {
	return os.Getenv("WIZ7_R999_ITEM_PROBE_DIR") != "" && !r999ItemFailed && r1000Cursor.records < 128
}
func r1000Read(o *oracle.Oracle, a oracle.Addr, n int) []byte {
	if !r999ItemSpan(a, n) {
		return nil
	}
	return append([]byte(nil), o.Bytes(a, n)...)
}
func r1000Word(b []byte, offset int) uint16 { return uint16(b[offset]) | uint16(b[offset+1])<<8 }
func r1000Ptr(b []byte, offset int) oracle.Addr {
	return oracle.Far(r1000Word(b, offset+2), r1000Word(b, offset))
}
func r1000Null(p oracle.Addr) bool { return p.Seg == 0 && p.Off == 0 }
func r1000SHA(b []byte) string     { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func r1000Signature(o *oracle.Oracle, offset uint16, want string) bool {
	expected, err := hex.DecodeString(want)
	return err == nil && len(expected) == 32 && bytes.Equal(r1000Read(o, oracle.Far(0x0110, offset), 32), expected)
}
func r1000Driver(o *oracle.Oracle) (uint16, bool) {
	p := r1000Read(o, oracle.Far(0x0110, 0x2de7), 4)
	if len(p) != 4 || r1000Word(p, 0) != 0x0163 {
		return 0, false
	}
	seg := r1000Word(p, 2)
	code := r1000Read(o, oracle.Far(seg, 0x0163), 4)
	return seg, bytes.Equal(code, []byte{0xe8, 0x63, 0x13, 0xcb})
}
func r1000Row(reg *r1000Registry, status string, pub, step uint64) r1000CursorRow {
	row := r1000CursorRow{Status: status, Publication: pub, Step: step, Files: map[string]string{}, SHA256: map[string]string{}}
	if reg != nil {
		row.Registry, row.NULID, row.ServiceID = reg.id, reg.nulID, reg.serviceID
		row.Driver, row.PicturePointer, row.Samples = reg.driver, reg.ptr, reg.samples
		row.Files["picture"] = reg.pictureFile
		row.SHA256["picture"] = reg.pictureSHA
		if reg.reset {
			row.TriggerNULID = reg.resetNUL
			row.ServiceID = reg.resetService
		}
	}
	return row
}
func r1000Emit(row r1000CursorRow, data map[string][]byte) {
	if !r1000Enabled() {
		return
	}
	r1000Cursor.records++
	row.Record = r1000Cursor.records
	row.LastRecord = row.Record == 128
	for role, b := range data {
		name := fmt.Sprintf("cursor-%03d-%s.bin", row.Record, role)
		if !r999ItemWrite(name, b) {
			r1000Cursor.active = nil
			r1000Cursor.candidate = nil
			return
		}
		row.Files[role] = name
		row.SHA256[role] = r1000SHA(b)
	}
	r999ItemJSON(fmt.Sprintf("cursor-%03d.json", row.Record), row)
	if r999ItemFailed || row.LastRecord {
		r1000Cursor.active = nil
		r1000Cursor.candidate = nil
	}
}
func r1000CursorDiscard(reason string) {
	// 新NUL/重入/失配只清尚未接受的候選，不延長既有registry期限。
	r1000Cursor.candidate = nil
}

func r1000CursorServiceEntry(s r999ServiceSample) {
	if !r1000Enabled() {
		return
	}
	r1000Cursor.candidate = nil
	if len(s.EntryStack) != 16 || s.Entry.CS != 0x0110 || s.NULID <= 0 {
		return
	}
	word := func(n int) uint16 { return r1000Word(s.EntryStack, 2*n) }
	if word(0) == 0x5cb9 && word(3) == 76 && word(4) == 13 {
		ptr := oracle.Far(word(2), word(1))
		if r1000Null(ptr) || !r999ItemSpan(ptr, 988) {
			return
		}
		r1000Cursor.candidate = &r1000Candidate{serviceID: s.ID, nulID: s.NULID, sourceDS: s.Entry.DS, ptr: ptr}
	} else if word(0) == 0x5b80 && word(1) == 0 && word(2) == 0 && word(3) == 0 && word(4) == 0 {
		r1000Cursor.candidate = &r1000Candidate{serviceID: s.ID, nulID: s.NULID, reset: true}
	}
}

func r1000CursorServiceReturn(o *oracle.Oracle, s r999ServiceSample) {
	c := r1000Cursor.candidate
	r1000Cursor.candidate = nil
	if !r1000Enabled() || c == nil {
		return
	}
	if !s.Accepted || c.serviceID != s.ID || c.nulID != s.NULID {
		return
	}
	row := r1000Row(r1000Cursor.active, "", 0, s.ReturnStep)
	row.TriggerNULID, row.ServiceID = c.nulID, c.serviceID
	if c.reset {
		if !r1000Signature(o, 0x5b60, r1000Reset) {
			row.Status = "reset_signature"
			r1000Emit(row, nil)
			return
		}
	} else {
		if r1000Cursor.definitions >= 8 {
			row.Status = "definition_limit"
			r1000Cursor.active = nil
			r1000Emit(row, nil)
			return
		}
		if !r1000Signature(o, 0x5c24, r1000DefA) || !r1000Signature(o, 0x5c9f, r1000DefB) {
			row.Status = "definition_signature"
			r1000Emit(row, nil)
			return
		}
	}
	driver, ok := r1000Driver(o)
	if !ok {
		row.Status = "driver_identity"
		r1000Emit(row, nil)
		return
	}
	if c.reset {
		reg := r1000Cursor.active
		if reg == nil || reg.driver != driver {
			row.Status = "reset_no_registry"
			r1000Emit(row, nil)
			return
		}
		reg.reset, reg.resetNUL, reg.resetService = true, c.nulID, c.serviceID
		row = r1000Row(reg, "reset_request", 0, s.ReturnStep)
		r1000Emit(row, nil)
		return
	}
	low := r1000Read(o, oracle.Far(c.sourceDS, 0x0e1e), 2)
	high := r1000Read(o, oracle.Far(c.sourceDS, 0x0e20), 2)
	if len(low) != 2 || len(high) != 2 || c.ptr != oracle.Far(r1000Word(high, 0), r1000Word(low, 0)) {
		row.Status = "source_globals"
		r1000Emit(row, nil)
		return
	}
	if s.ReturnStep > ^uint64(0)-30000000 {
		row.Status = "step_overflow"
		r1000Emit(row, nil)
		return
	}
	picture := r1000Read(o, c.ptr, 988)
	if len(picture) != 988 {
		row.Status = "picture_span"
		r1000Emit(row, nil)
		return
	}
	r1000Cursor.definitions++
	reg := &r1000Registry{id: r1000Cursor.definitions, nulID: c.nulID, serviceID: c.serviceID, driver: driver, ptr: c.ptr, picture: picture, start: s.ReturnStep, deadline: s.ReturnStep + 30000000}
	reg.pictureFile = fmt.Sprintf("cursor-picture-%02d.bin", reg.id)
	reg.pictureSHA = r1000SHA(picture)
	if !r999ItemWrite(reg.pictureFile, picture) {
		return
	}
	r1000Cursor.active = reg
	r1000Emit(r1000Row(reg, "registered", 0, s.ReturnStep), nil)
}

func r1000CursorSnapshot(o *oracle.Oracle, pub uint64) {
	reg := r1000Cursor.active
	if !r1000Enabled() || reg == nil {
		return
	}
	step := o.Steps()
	if step < reg.start || step > reg.deadline || reg.samples >= 128 {
		r1000Emit(r1000Row(reg, "window_complete", pub, step), nil)
		r1000Cursor.active = nil
		return
	}
	reg.samples++
	row := r1000Row(reg, "", pub, step)
	desc := r1000Read(o, oracle.Far(reg.driver, 0x15ac), 24)
	full := r1000Read(o, oracle.Far(reg.driver, 0x1e5c), 12)
	busy := r1000Read(o, oracle.Far(reg.driver, 0x1ed3), 1)
	if len(desc) != 24 || len(full) != 12 || len(busy) != 1 {
		row.Status = "descriptor_span"
		r1000Emit(row, nil)
		r1000Cursor.active = nil
		return
	}
	data := map[string][]byte{"descriptor": desc, "picture-info": full}
	row.Busy = busy[0]
	row.X, row.Y = int(r1000Word(desc, 0)), int(r1000Word(desc, 2))
	row.ClipW, row.ClipH = int(r1000Word(desc, 12)), int(r1000Word(desc, 14))
	row.BackgroundPointer = oracle.Far(reg.driver, r1000Word(desc, 8))
	row.FullW, row.FullH = int(r1000Word(full, 0)), int(r1000Word(full, 2))
	actual := r1000Ptr(full, 8)
	if reg.reset {
		row.Status = "reset_unverified"
		if row.Busy == 0 && r1000Null(actual) {
			row.Status = "reset_complete"
		}
		r1000Emit(row, data)
		r1000Cursor.active = nil
		return
	}
	if row.Busy != 0 {
		row.Status = "busy"
		r1000Emit(row, data)
		return
	}
	if actual != reg.ptr || row.FullW != 76 || row.FullH != 13 {
		row.Status = "picture_changed"
		r1000Emit(row, data)
		r1000Cursor.active = nil
		return
	}
	if row.X < 0 || row.X >= 320 || row.Y < 0 || row.Y >= 200 || row.ClipW != min(76, 320-row.X) || row.ClipH != min(13, 200-row.Y) {
		row.Status = "clip_or_coordinates"
		r1000Emit(row, data)
		r1000Cursor.active = nil
		return
	}
	area := row.ClipW * row.ClipH
	if area < 1 || area > 988 || r1000Null(row.BackgroundPointer) {
		row.Status = "background_shape"
		r1000Emit(row, data)
		r1000Cursor.active = nil
		return
	}
	background := r1000Read(o, row.BackgroundPointer, area)
	if len(background) != area {
		row.Status = "background_span"
		r1000Emit(row, data)
		r1000Cursor.active = nil
		return
	}
	screen := append([]byte(nil), o.Indexed()...)
	if len(screen) != 64000 {
		r999ItemInvalid("cursor_indexed_size")
		r1000Cursor.active = nil
		return
	}
	visible := make([]byte, area)
	for y := 0; y < row.ClipH; y++ {
		copy(visible[y*row.ClipW:(y+1)*row.ClipW], screen[(row.Y+y)*320+row.X:(row.Y+y)*320+row.X+row.ClipW])
	}
	data["background"], data["indexed"] = background, visible
	row.Status = "sample"
	r1000Emit(row, data)
}
