package main

// spec024 R1002：只接已量創角M來源，不改原版資料或既有兩路性別。
import (
	"bytes"
	"log"

	"github.com/wicanr2/dosgolem/oracle"
)

var creationSexReturn = oracle.Far(0x0110, 0x9595)

func (ov *overlay) clearCreationSexEvents() {
	for _, e := range ov.evs {
		if e.sex && e.originRet == creationSexReturn {
			e.dead = true
		}
	}
}
func (ov *overlay) clearCreationSex() {
	ov.clearCreationSexEvents()
	if ov.sexPending != nil && ov.sexPending.ret == creationSexReturn {
		ov.sexPending = nil
	}
}
func (ov *overlay) rejectCreationSex(reason string) {
	ov.clearCreationSex()
	log.Print("創角性別缺譯：" + reason)
}

// 只讀driver metadata；呼叫端先核身分/游標再決定是否讀完整畫格。
func creationSexSource(o heldMachineReader) (driver, buffer uint16, x, y int, ok bool) {
	if o.VideoMode() != 0x13 {
		return
	}
	driver, ok = heldDriver(o)
	if !ok {
		return
	}
	xy := heldRead(o, oracle.Far(driver, 0x0295), 4)
	buf := heldRead(o, oracle.Far(driver, 0x01e5), 2)
	if len(xy) != 4 || len(buf) != 2 || heldWord(buf, 0) == 0 {
		ok = false
		return
	}
	buffer = heldWord(buf, 0)
	x, y = int(heldWord(xy, 0)), int(heldWord(xy, 2))
	return
}

// 回傳true時主hook不再落入舊分支。false只可能出現在已驗六byte stack。
func (ov *overlay) enterCreationSex(o heldMachineReader) bool {
	r := o.Regs()
	if !heldAt(r, 0x3877) {
		return true
	}
	nested := ov.sexPending != nil && ov.sexPending.ret == creationSexReturn
	if nested {
		ov.clearCreationSex()
	}
	if r.SP > 65530 {
		ov.rejectCreationSex("stack span")
		return true
	}
	stack := heldRead(o, oracle.Far(r.SS, r.SP), 6)
	if len(stack) != 6 {
		ov.rejectCreationSex("stack span")
		return true
	}
	ret := oracle.Far(r.CS, heldWord(stack, 0))
	if ret != creationSexReturn {
		if ret == profileSexReturn || ret == worldSexReturn {
			ov.clearCreationSex()
		}
		return false
	}
	ov.clearCreationSex()
	if nested {
		ov.rejectCreationSex("nested")
		return true
	}
	if heldWord(stack, 2) != 'M' || heldWord(stack, 4) != 0 || (flowK != 2 && flowK != 3) ||
		!worldSexResourcesValid(ov.sexEntries) {
		ov.rejectCreationSex("tuple or resource")
		return true
	}
	if ov.sexPending != nil || ov.pend != nil || ov.compound != nil {
		ov.rejectCreationSex("interleaving")
		return true
	}
	if !bytes.Equal(heldRead(o, oracle.Far(0x0110, 0x9592), 3), []byte{0xe8, 0xe2, 0xa2}) {
		ov.rejectCreationSex("call signature")
		return true
	}
	driver, buffer, x, y, ok := creationSexSource(o)
	if !ok || x != 90 || y != 6 {
		ov.rejectCreationSex("driver or cursor")
		return true
	}
	before := heldRead(o, oracle.Far(buffer, 0), 64000)
	if len(before) != 64000 || !worldSexBackgroundMatches(before) {
		ov.rejectCreationSex("frame or background")
		return true
	}
	g := ov.sexEntries["M"].G[0]
	ov.sexPending = &pendingProfileSex{s: "M", sp: r.SP, ret: ret, x: x, y: y,
		before: before, g: [][]uint16{append([]uint16(nil), g...)},
		ss: r.SS, entryStep: o.Steps(), driver: driver, buffer: buffer}
	return true
}

func (ov *overlay) returnCreationSex(o heldMachineReader) {
	r := o.Regs()
	if !heldAt(r, creationSexReturn.Off) || ov.sexPending == nil || ov.sexPending.ret != creationSexReturn {
		return
	}
	p := ov.sexPending
	ov.clearCreationSexEvents()
	step := o.Steps()
	if r.SS != p.ss || r.SP != p.sp+2 || step < p.entryStep || step-p.entryStep > 1000000 {
		ov.rejectCreationSex("return identity or interval")
		return
	}
	driver, buffer, x, y, ok := creationSexSource(o)
	if !ok || driver != p.driver || buffer != p.buffer || x != 96 || y != 6 {
		ov.rejectCreationSex("return driver or cursor")
		return
	}
	after := heldRead(o, oracle.Far(buffer, 0), 64000)
	// 共用返回驗完整全幀差分、6×6/色1，清slot並建立8×8事件。
	ov.returnProfileSex(r.SP, x, y, after, o.Palette())
}
