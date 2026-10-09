package main

import (
	"fmt"
	"os"

	"github.com/wicanr2/dosgolem/oracle"
)

type r999ServiceSample struct {
	ID, NULID                                      int
	Deadline, EntryStep, ReturnStep                uint64
	Entry, Returned                                oracle.Regs
	EntryStack, ReturnStack, EntryCode, ReturnCode []byte
	Accepted                                       bool
	Reason                                         string
	Offscreen, Indexed                             r999ItemDiff
}
type r999ServicePending struct {
	sample          r999ServiceSample
	before, indexed []byte
}

var r999Service struct {
	installed                 bool
	nulID, windowCount, count int
	unpairedCount             int
	windowUnpaired            bool
	deadline                  uint64
	pending                   *r999ServicePending
}

func r999ServiceBytes(o *oracle.Oracle, a oracle.Addr, n int) []byte {
	if !r999ItemSpan(a, n) {
		return nil
	}
	return append([]byte(nil), o.Bytes(a, n)...)
}
func r999ServiceDiscard(reason string) {
	if p := r999Service.pending; p != nil {
		p.sample.Reason = reason
		r999ItemJSON(fmt.Sprintf("service24-%03d-return.json", p.sample.ID), p.sample)
		r999Service.pending = nil
	}
}
func r999ServiceArm(o *oracle.Oracle, id int, step uint64) {
	if os.Getenv("WIZ7_R999_ITEM_PROBE_DIR") == "" || r999ItemFailed {
		return
	}
	r999ServiceDiscard("new_nul_window")
	if step > ^uint64(0)-1000000 {
		r999ItemInvalid("service_deadline")
		return
	}
	r999Service.nulID = id
	r999Service.deadline = step + 1000000
	r999Service.windowCount = 0
	r999Service.windowUnpaired = false
	if r999Service.installed {
		return
	}
	r999Service.installed = true
	o.OnCall(oracle.Far(0x0110, 0x40a7), func(o *oracle.Oracle) {
		step := o.Steps()
		if r999ItemFailed {
			return
		}
		if step > r999Service.deadline {
			r999ServiceDiscard("entry_after_deadline")
			return
		}
		if r999Service.pending != nil {
			r999ServiceDiscard("reentry")
			return
		}
		if r999Service.windowCount >= 16 || r999Service.count >= 64 {
			return
		}
		r := o.Regs()
		stack := r999ServiceBytes(o, oracle.Far(r.SS, r.SP), 16)
		code := r999ServiceBytes(o, oracle.Far(r.CS, r.IP), 12)
		if len(stack) != 16 || len(code) != 12 {
			return
		}
		before, indexed := append([]byte(nil), offscreen(o)...), append([]byte(nil), o.Indexed()...)
		if len(before) != 64000 || len(indexed) != 64000 {
			r999ItemInvalid("service_entry_frames")
			return
		}
		r999Service.count++
		r999Service.windowCount++
		s := r999ServiceSample{ID: r999Service.count, NULID: r999Service.nulID, Deadline: r999Service.deadline, EntryStep: step, Entry: r, EntryStack: stack, EntryCode: code}
		r999Service.pending = &r999ServicePending{sample: s, before: before, indexed: indexed}
		r999ItemJSON(fmt.Sprintf("service24-%03d-entry.json", s.ID), s)
	})
	o.OnCall(oracle.Far(0x0110, 0x40b0), func(o *oracle.Oracle) {
		p := r999Service.pending
		if r999ItemFailed {
			return
		}
		if p == nil {
			step := o.Steps()
			if step <= r999Service.deadline && !r999Service.windowUnpaired && r999Service.unpairedCount < 64 {
				r999Service.windowUnpaired = true
				r999Service.unpairedCount++
				r999ItemJSON(fmt.Sprintf("service24-unpaired-%03d.json", r999Service.unpairedCount), map[string]any{"Reason": "unpaired_return", "NULID": r999Service.nulID, "Step": step})
			}
			return
		}
		r999Service.pending = nil
		s := p.sample
		s.Returned = o.Regs()
		s.ReturnStep = o.Steps()
		name := fmt.Sprintf("service24-%03d", s.ID)
		if s.NULID != r999Service.nulID || s.Returned.SS != s.Entry.SS || s.ReturnStep < s.EntryStep || s.ReturnStep > s.Deadline {
			s.Reason = "return_window_or_ss"
			r999ItemJSON(name+"-return.json", s)
			return
		}
		s.ReturnStack = r999ServiceBytes(o, oracle.Far(s.Returned.SS, s.Returned.SP), 16)
		s.ReturnCode = r999ServiceBytes(o, oracle.Far(s.Returned.CS, s.Returned.IP), 12)
		if len(s.ReturnStack) != 16 || len(s.ReturnCode) != 12 {
			s.Reason = "return_span"
			r999ItemJSON(name+"-return.json", s)
			return
		}
		after, visible := offscreen(o), o.Indexed()
		if len(after) != 64000 || len(visible) != 64000 {
			r999ItemInvalid("service_return_frames")
			return
		}
		s.Accepted = true
		s.Offscreen = r999ItemDiffROI(name+"-offscreen", p.before, after)
		s.Indexed = r999ItemDiffROI(name+"-indexed", p.indexed, visible)
		palette := make([]byte, 0, 768)
		for _, rgb := range o.Palette() {
			palette = append(palette, rgb[:]...)
		}
		r999ItemWrite(name+"-palette.bin", palette)
		r999ItemJSON(name+"-return.json", s)
	})
}
