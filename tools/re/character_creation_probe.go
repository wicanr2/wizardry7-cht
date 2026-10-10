package main

// R1002 P1–P5 私人唯讀觀察器；不註冊正式覆繪路由。
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/wicanr2/dosgolem/oracle"
)

const contextWindow uint64 = 1000000
const contextCharStart uint64 = 332800000
const contextCharEnd uint64 = 417800000

var contextExpected = map[uint16]string{
	127: "PSIONIC", 410: "PSIONIC", 707: "PSIONIC",
	406: "POISON", 477: "POISON", 4026: "POISON",
	478: "STONE", 3829: "STONE", 1492: "FORCE", 2101: "FORCE",
	1538: "LEAVE", 5102: "LEAVE", 5204: "LEAVE",
	3330: "TRUNK", 3337: "TRUNK", 3451: "HER", 3456: "HER",
	27222: "NAME:", 27307: "NAME:",
	8910: "Ort!", 15550: "Ort!", 15810: "Ort!", 17035: "Ort!",
	8915: "Hoyo...", 17950: "Hoyo...", 18720: "Hoyo...",
	8916: "Bodendoi...", 17955: "Bodendoi...", 18725: "Bodendoi...",
}

type contextReader struct {
	Generation, Step            uint64
	ID, DS, Dst, SS, SP, Return uint16
}
type contextDecoded struct {
	Reader contextReader
	Step   uint64
	Key    string
}
type contextCharacter struct {
	ExtensionCode                        []byte
	Number                               int
	Step                                 uint64
	SS, SP, Return, Code, Driver, Buffer uint16
	Before                               []byte // 自有64000-byte副本；不傳入JSON。
}
type contextRow struct {
	CodeAddressSpace string            `json:"code_address_space,omitempty"`
	CodeBase         []uint16          `json:"code_base,omitempty"`
	CodeBytes        []byte            `json:"code_bytes,omitempty"`
	CodeSHA256       string            `json:"code_sha256,omitempty"`
	EightRectangle   []int             `json:"eight_rectangle,omitempty"`
	Before8          []byte            `json:"before8,omitempty"`
	After8           []byte            `json:"after8,omitempty"`
	BackgroundPoints [][2]int          `json:"background_points,omitempty"`
	BackgroundBefore []byte            `json:"background_before,omitempty"`
	BackgroundAfter  []byte            `json:"background_after,omitempty"`
	Status           string            `json:"status"`
	Reason           string            `json:"reason,omitempty"`
	Step             uint64            `json:"step"`
	Generation       uint64            `json:"generation,omitempty"`
	ID               uint16            `json:"id,omitempty"`
	Key              string            `json:"key,omitempty"`
	DS               uint16            `json:"ds,omitempty"`
	Dst              uint16            `json:"dst,omitempty"`
	SS               uint16            `json:"ss,omitempty"`
	SP               uint16            `json:"sp,omitempty"`
	Return           uint16            `json:"return,omitempty"`
	Font             *uint16           `json:"font,omitempty"`
	X                *int              `json:"x,omitempty"`
	Y                *int              `json:"y,omitempty"`
	Number           int               `json:"number,omitempty"`
	Code             uint16            `json:"code,omitempty"`
	Driver           uint16            `json:"driver,omitempty"`
	Buffer           uint16            `json:"buffer,omitempty"`
	AfterDriver      uint16            `json:"after_driver,omitempty"`
	AfterBuffer      uint16            `json:"after_buffer,omitempty"`
	EntryStep        uint64            `json:"entry_step,omitempty"`
	Before           []byte            `json:"before,omitempty"`
	After            []byte            `json:"after,omitempty"`
	BBox             []int             `json:"difference_bbox,omitempty"`
	Differences      *int              `json:"differences,omitempty"`
	Outside          *int              `json:"outside_differences,omitempty"`
	Events           int               `json:"events,omitempty"`
	ReaderEntries    *int              `json:"reader_entries,omitempty"`
	DecodedEntries   *int              `json:"decoded_entries,omitempty"`
	PrintedEntries   *int              `json:"printed_entries,omitempty"`
	CharEntries      *int              `json:"character_entries,omitempty"`
	CharReturns      *int              `json:"character_returns,omitempty"`
	ReturnHooks      *int              `json:"return_hooks,omitempty"`
	ReaderPending    *bool             `json:"reader_pending,omitempty"`
	DecodedPending   *bool             `json:"decoded_pending,omitempty"`
	CharPending      *bool             `json:"character_pending,omitempty"`
	Rejects          map[string]uint64 `json:"rejects,omitempty"`
}
type contextProbe struct {
	extensionAttempted                                                      bool
	mu                                                                      sync.Mutex
	file                                                                    *os.File
	encoder                                                                 *json.Encoder
	stopped                                                                 bool
	events                                                                  int
	generation                                                              uint64
	pending                                                                 *contextReader
	decoded                                                                 *contextDecoded
	character                                                               *contextCharacter
	hooks                                                                   map[uint16]bool
	rejects                                                                 map[string]uint64
	readerEntries, decodedEntries, printedEntries, charEntries, charReturns int
	charCap, hookCap                                                        bool
}

func contextSpan(a oracle.Addr, n int) bool {
	return n > 0 && n <= 65536 && int(a.Off)+n <= 65536 && int(a.Seg)*16+int(a.Off)+n <= 0xa0000
}
func contextRead(o *oracle.Oracle, a oracle.Addr, n int) []byte {
	if !contextSpan(a, n) {
		return nil
	}
	return heldRead(o, a, n)
}
func contextLive(now, entry uint64) bool { return now >= entry && now-entry <= contextWindow }
func contextInt(n int) *int              { return &n }
func contextU16(n uint16) *uint16        { return &n }
func contextBool(b bool) *bool           { return &b }
func contextText(o *oracle.Oracle, ds, off uint16) (string, bool) {
	a := oracle.Far(ds, off)
	if !contextSpan(a, 64) {
		return "", false
	}
	b := make([]byte, 0, 63)
	for i := 0; i < 64; i++ {
		one := contextRead(o, oracle.Far(ds, uint16(int(off)+i)), 1)
		if len(one) != 1 {
			return "", false
		}
		if one[0] == 0 {
			return string(b), true
		}
		if one[0] < 32 || one[0] > 126 {
			return "", false
		}
		b = append(b, one[0])
	}
	return "", false
}
func contextXY(o *oracle.Oracle) (uint16, int, int, bool) {
	driver, ok := heldDriver(o)
	if !ok {
		return 0, 0, 0, false
	}
	b := contextRead(o, oracle.Far(driver, 0x0295), 4)
	if len(b) != 4 {
		return 0, 0, 0, false
	}
	return driver, int(heldWord(b, 0)), int(heldWord(b, 2)), true
}
func contextFrame(o *oracle.Oracle, driver uint16) (uint16, []byte) {
	if o.VideoMode() != 0x13 {
		return 0, nil
	}
	b := contextRead(o, oracle.Far(driver, 0x01e5), 2)
	if len(b) != 2 || heldWord(b, 0) == 0 {
		return 0, nil
	}
	buffer := heldWord(b, 0)
	// 完整span先核，不能只核起址；contextRead持有自有副本。
	return buffer, contextRead(o, oracle.Far(buffer, 0), 64000)
}
func contextCell(frame []byte) []byte {
	if len(frame) != 64000 {
		return nil
	}
	cell := make([]byte, 36)
	for y := 0; y < 6; y++ {
		copy(cell[y*6:(y+1)*6], frame[(6+y)*320+90:(6+y)*320+96])
	}
	return cell
}

func (p *contextProbe) invalidLocked(code string) {
	if p.stopped {
		return
	}
	p.stopped = true
	p.pending = nil
	p.decoded = nil
	p.character = nil
	if p.file != nil {
		_ = p.file.Close()
		p.file = nil
	}
	// code只從本檔有限常數而來，不格式化任意Error或原版內容。
	fmt.Fprintln(os.Stderr, "R1002_CONTEXT_INVALID_"+code)
}
func (p *contextProbe) writeLocked(row contextRow) {
	if p.stopped {
		return
	}
	if p.events >= 256 {
		p.invalidLocked("EVENT_CAP")
		return
	}
	if p.encoder.Encode(row) != nil {
		p.invalidLocked("IO")
		return
	}
	p.events++
}
func (p *contextProbe) rejectLocked(reason string, step uint64) {
	p.rejects[reason]++
	p.writeLocked(contextRow{Status: "reject", Reason: reason, Step: step, Generation: p.generation})
}
func (p *contextProbe) Close() {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return
	}
	p.writeLocked(contextRow{Status: "summary", Generation: p.generation, Events: p.events, ReaderEntries: contextInt(p.readerEntries),
		DecodedEntries: contextInt(p.decodedEntries), PrintedEntries: contextInt(p.printedEntries), CharEntries: contextInt(p.charEntries), CharReturns: contextInt(p.charReturns),
		ReturnHooks: contextInt(len(p.hooks)), ReaderPending: contextBool(p.pending != nil), DecodedPending: contextBool(p.decoded != nil), CharPending: contextBool(p.character != nil), Rejects: p.rejects})
	if p.stopped {
		return
	}
	if p.file.Close() != nil {
		p.file = nil
		p.invalidLocked("CLOSE")
		return
	}
	p.file = nil
	p.stopped = true
	p.pending = nil
	p.decoded = nil
	p.character = nil
}

func r1002Setup(o *oracle.Oracle) *contextProbe {
	name := os.Getenv("WIZ7_R1002_CONTEXT")
	if name == "" {
		return nil
	} // off時完全不安裝任何私人source hook。
	file, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "R1002_CONTEXT_INVALID_OPEN")
		return nil
	}
	p := &contextProbe{file: file, encoder: json.NewEncoder(file), hooks: map[uint16]bool{}, rejects: map[string]uint64{}}
	p.writeLocked(contextRow{Status: "setup"})
	if p.stopped {
		return p
	}
	o.OnCall(oracle.Far(0x0110, 0x0582), p.begin)
	o.OnCall(oracle.Far(0x0110, 0x0659), p.finish)
	o.OnCall(oracle.Far(0x0110, 0x3895), p.print)
	o.OnCall(oracle.Far(0x0110, 0x3877), p.characterBegin)
	return p
}

func (p *contextProbe) begin(o *oracle.Oracle) {
	r := o.Regs()
	if !heldAt(r, 0x0582) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return
	}
	p.pending = nil
	p.decoded = nil
	if p.generation == ^uint64(0) {
		p.invalidLocked("GENERATION")
		return
	}
	p.generation++
	if r.SP < 0x112 || r.SP > 65530 {
		return
	}
	// P1：只讀ID；保護或非白名單來源在ret/dst讀取之前返回。
	idBytes := contextRead(o, oracle.Far(r.SS, r.SP+2), 2)
	if len(idBytes) != 2 {
		return
	}
	id := heldWord(idBytes, 0)
	if id >= 2500 && id <= 2574 {
		return
	}
	if _, ok := contextExpected[id]; !ok {
		return
	}
	ret := contextRead(o, oracle.Far(r.SS, r.SP), 2)
	dst := contextRead(o, oracle.Far(r.SS, r.SP+4), 2)
	if len(ret) != 2 || len(dst) != 2 {
		p.rejectLocked("reader_stack", o.Steps())
		return
	}
	p.pending = &contextReader{p.generation, o.Steps(), id, r.DS, heldWord(dst, 0), r.SS, r.SP, heldWord(ret, 0)}
	p.readerEntries++
	p.writeLocked(contextRow{Status: "reader", Step: o.Steps(), Generation: p.generation, ID: id, DS: r.DS, Dst: p.pending.Dst, SS: r.SS, SP: r.SP, Return: p.pending.Return})
}
func (p *contextProbe) finish(o *oracle.Oracle) {
	r := o.Regs()
	if !heldAt(r, 0x0659) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped || p.pending == nil {
		return
	}
	reader := p.pending
	p.pending = nil
	p.decoded = nil
	if reader.Generation != p.generation || !contextLive(o.Steps(), reader.Step) || r.DS != reader.DS || r.SS != reader.SS ||
		r.SP != reader.SP-0x112 || r.BP != reader.SP-2 || r.BP > 65528 {
		p.rejectLocked("decode_pair", o.Steps())
		return
	}
	fields := contextRead(o, oracle.Far(r.SS, uint16(int(r.BP)+4)), 4)
	if len(fields) != 4 || heldWord(fields, 0) != reader.ID || heldWord(fields, 2) != reader.Dst {
		p.rejectLocked("decode_fields", o.Steps())
		return
	}
	key, ok := contextText(o, reader.DS, reader.Dst)
	if !ok || key != contextExpected[reader.ID] {
		p.rejectLocked("decode_key", o.Steps())
		return
	}
	p.decoded = &contextDecoded{*reader, o.Steps(), key}
	p.decodedEntries++
	p.writeLocked(contextRow{Status: "decoded", Step: o.Steps(), Generation: p.generation, ID: reader.ID, Key: key, DS: reader.DS, Dst: reader.Dst, SS: r.SS, SP: r.SP, Return: reader.Return, EntryStep: reader.Step})
}
func (p *contextProbe) print(o *oracle.Oracle) {
	r := o.Regs()
	if !heldAt(r, 0x3895) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped || p.decoded == nil {
		return
	}
	d := p.decoded
	if d.Reader.Generation != p.generation || !contextLive(o.Steps(), d.Reader.Step) {
		p.decoded = nil
		p.rejectLocked("print_window", o.Steps())
		return
	}
	if r.DS != d.Reader.DS || r.SP > 65530 {
		return
	}
	stack := contextRead(o, oracle.Far(r.SS, r.SP), 6)
	if len(stack) != 6 {
		p.decoded = nil
		p.rejectLocked("print_stack", o.Steps())
		return
	}
	if heldWord(stack, 2) != d.Reader.Dst {
		return
	}
	p.decoded = nil
	key, ok := contextText(o, r.DS, d.Reader.Dst)
	if !ok || key != d.Key || key != contextExpected[d.Reader.ID] {
		p.rejectLocked("print_key", o.Steps())
		return
	}
	driver, x, y, ok := contextXY(o)
	if !ok {
		p.rejectLocked("print_driver", o.Steps())
		return
	}
	p.printedEntries++
	p.writeLocked(contextRow{Status: "print", Step: o.Steps(), Generation: p.generation, ID: d.Reader.ID, Key: key, DS: r.DS, Dst: d.Reader.Dst, SS: r.SS, SP: r.SP, Return: heldWord(stack, 0), Font: contextU16(heldWord(stack, 4)), X: contextInt(x), Y: contextInt(y), Driver: driver, EntryStep: d.Reader.Step})
}

func (p *contextProbe) characterBegin(o *oracle.Oracle) {
	r := o.Regs()
	if !heldAt(r, 0x3877) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return
	}
	p.character = nil // 每個真正3877入口先丟舊代，窗口外也不可保留。
	step := o.Steps()
	if step < contextCharStart || step > contextCharEnd || r.SP > 65530 {
		return
	}
	stack := contextRead(o, oracle.Far(r.SS, r.SP), 6)
	if len(stack) != 6 {
		return
	}
	ret, code, font := heldWord(stack, 0), heldWord(stack, 2), heldWord(stack, 4)
	if (code != 'M' && code != 'F') || font != 0 {
		return
	}
	driver, x, y, ok := contextXY(o)
	if !ok || x != 90 || y != 6 {
		return
	}
	if p.charEntries >= 16 {
		if !p.charCap {
			p.charCap = true
			p.rejectLocked("char_cap", step)
		}
		return
	}
	if !p.hooks[ret] && len(p.hooks) >= 8 {
		if !p.hookCap {
			p.hookCap = true
			p.rejectLocked("hook_cap", step)
		}
		return
	}
	buffer, before := contextFrame(o, driver)
	if len(before) != 64000 {
		p.rejectLocked("char_before", step)
		return
	}
	if !p.hooks[ret] {
		p.hooks[ret] = true
		o.OnCall(oracle.Far(0x0110, ret), func(o *oracle.Oracle) { p.characterReturn(o, ret) })
	}
	p.charEntries++
	p.character = &contextCharacter{nil, p.charEntries, step, r.SS, r.SP, ret, code, driver, buffer, before}
	p.writeLocked(contextRow{Status: "character_entry", Step: step, Number: p.charEntries, Code: code, SS: r.SS, SP: r.SP, Return: ret, Font: contextU16(font), X: contextInt(x), Y: contextInt(y), Driver: driver, Buffer: buffer, Before: contextCell(before)})
	p.extensionBeginLocked(o, p.character)
}
func (p *contextProbe) characterReturn(o *oracle.Oracle, ret uint16) {
	r := o.Regs()
	if !heldAt(r, ret) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped || p.character == nil || p.character.Return != ret {
		return
	}
	c := p.character
	p.character = nil
	step := o.Steps()
	if !contextLive(step, c.Step) || step < contextCharStart || step > contextCharEnd || r.SS != c.SS || r.SP != c.SP+2 {
		p.rejectLocked("char_return_pair", step)
		return
	}
	driver, x, y, ok := contextXY(o)
	if !ok || x != 96 || y != 6 {
		p.rejectLocked("char_return_driver", step)
		return
	}
	buffer, after := contextFrame(o, driver)
	if len(after) != 64000 {
		p.rejectLocked("char_after", step)
		return
	}
	count, outside := 0, 0
	minX, minY, maxX, maxY := 320, 200, -1, -1
	for i, b := range c.Before {
		if b == after[i] {
			continue
		}
		count++
		x, y := i%320, i/320
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		if x < 90 || x >= 96 || y < 6 || y >= 12 {
			outside++
		}
	}
	var bbox []int
	if count > 0 {
		bbox = []int{minX, minY, maxX + 1, maxY + 1}
	}
	p.charReturns++
	p.writeLocked(contextRow{Status: "character_return", Step: step, Number: c.Number, Code: c.Code, SS: r.SS, SP: r.SP, Return: ret, Font: contextU16(0), X: contextInt(x), Y: contextInt(y),
		Driver: c.Driver, Buffer: c.Buffer, AfterDriver: driver, AfterBuffer: buffer, EntryStep: c.Step,
		Before: contextCell(c.Before), After: contextCell(after), BBox: bbox, Differences: contextInt(count), Outside: contextInt(outside)})
	p.extensionReturnLocked(c, step, driver, buffer, after)
}

// P6 helpers只在已接受P3樣本上運作；原R1資料與守門保持。
var contextExtensionPoints = [...][2]int{
	{88, 4}, {92, 4}, {95, 4}, {88, 7},
	{89, 9}, {88, 11}, {91, 5}, {95, 5},
}

func contextExtensionCell(frame []byte) []byte {
	if len(frame) != 64000 {
		return nil
	}
	cell := make([]byte, 64)
	for y := 0; y < 8; y++ {
		copy(cell[y*8:(y+1)*8], frame[(4+y)*320+88:(4+y)*320+96])
	}
	return cell
}

func contextExtensionBackground(frame []byte) []byte {
	if len(frame) != 64000 {
		return nil
	}
	values := make([]byte, len(contextExtensionPoints))
	for i, point := range contextExtensionPoints {
		values[i] = frame[point[1]*320+point[0]]
	}
	return values
}

func (p *contextProbe) extensionBeginLocked(o *oracle.Oracle, c *contextCharacter) {
	if p.stopped || c == nil || p.extensionAttempted || c.Code != 'M' || c.Return != 0x9595 ||
		c.Step < 334900000 || c.Step > 335000000 {
		return
	}
	// 單次attempt先消耗，簽章失敗不能在後續重試取得128-byte片段。
	p.extensionAttempted = true
	signature := contextRead(o, oracle.Far(0x0110, 0x9592), 3)
	if !bytes.Equal(signature, []byte{0xe8, 0xe2, 0xa2}) {
		p.rejectLocked("extension_signature", c.Step)
		return
	}
	// contextRead在任何Bytes呼叫前驗完整offset及物理128-byte span。
	fragment := contextRead(o, oracle.Far(0x0110, 0x9515), 128)
	if len(fragment) != 128 || !bytes.Equal(fragment[125:], signature) {
		p.rejectLocked("extension_code", c.Step)
		return
	}
	c.ExtensionCode = fragment
}

// 不接收Oracle，不能新增原版畫格或任何記憶體讀取。
func (p *contextProbe) extensionReturnLocked(c *contextCharacter, step uint64, driver, buffer uint16, after []byte) {
	if p.stopped || c == nil || len(c.ExtensionCode) != 128 {
		return
	}
	if step < 334900000 || step > 335000000 || driver != c.Driver || buffer != c.Buffer {
		p.rejectLocked("extension_return_identity", step)
		return
	}
	before8, after8 := contextExtensionCell(c.Before), contextExtensionCell(after)
	beforePoints, afterPoints := contextExtensionBackground(c.Before), contextExtensionBackground(after)
	if len(before8) != 64 || len(after8) != 64 || len(beforePoints) != 8 || len(afterPoints) != 8 {
		p.rejectLocked("extension_owned_frame", step)
		return
	}
	hash := sha256.Sum256(c.ExtensionCode)
	p.writeLocked(contextRow{Status: "source_extension", Step: step, Number: c.Number, Code: c.Code,
		SS: c.SS, SP: c.SP + 2, Return: c.Return, Font: contextU16(0), X: contextInt(96), Y: contextInt(6),
		Driver: c.Driver, Buffer: c.Buffer, AfterDriver: driver, AfterBuffer: buffer, EntryStep: c.Step,
		CodeAddressSpace: "dosgolem runtime real-mode CS:IP", CodeBase: []uint16{0x0110, 0x9515},
		CodeBytes: c.ExtensionCode, CodeSHA256: hex.EncodeToString(hash[:]), EightRectangle: []int{88, 4, 96, 12},
		Before8: before8, After8: after8, BackgroundPoints: append([][2]int(nil), contextExtensionPoints[:]...),
		BackgroundBefore: beforePoints, BackgroundAfter: afterPoints})
}
