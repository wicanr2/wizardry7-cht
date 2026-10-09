package main

// R999私人source witness；服務24的新位置背景仍未知。
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/oracle"
)

type r999ItemDiff struct {
	X, Y, W, H, Changed int
}
type r999ItemWitness struct {
	ID                    int
	English               string
	EntryRet              oracle.Addr
	EntryRegs, ReturnRegs oracle.Regs
	Font, SourceOff       uint16
	X, Y                  int
	EntryStep, ReturnStep uint64
	ReturnAccepted        bool
	Reason                string
	Offscreen, Indexed    r999ItemDiff
}

var r999ItemCount int
var r999ItemFailed bool

func r999ItemInvalid(code string) {
	if !r999ItemFailed {
		r999ItemFailed = true
		fmt.Fprintf(os.Stderr, "R999_ITEM_PROBE_INVALID code=%s\n", code)
	}
}
func r999ItemWrite(name string, data []byte) bool {
	if r999ItemFailed {
		return false
	}
	root := os.Getenv("WIZ7_R999_ITEM_PROBE_DIR")
	if root == "" {
		return false
	}
	f, err := os.OpenFile(filepath.Join(root, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		r999ItemInvalid("open")
		return false
	}
	n, writeErr := f.Write(data)
	syncErr := f.Sync()
	closeErr := f.Close()
	if writeErr != nil || n != len(data) || syncErr != nil || closeErr != nil {
		r999ItemInvalid("write_sync_close")
		return false
	}
	return true
}
func r999ItemJSON(name string, value any) {
	b, err := json.Marshal(value)
	if err != nil {
		r999ItemInvalid("encode")
		return
	}
	r999ItemWrite(name, append(b, '\n'))
}
func r999ItemSpan(a oracle.Addr, size int) bool {
	return size > 0 && size <= 65536 && int(a.Off) <= 65536-size && int(a.Seg)*16+int(a.Off)+size <= 0xa0000
}

// 完整64K僅在記憶體比較；只輸出實際changed bbox的before/after/mask。
func r999ItemDiffROI(prefix string, before, after []byte) r999ItemDiff {
	if len(before) != 64000 || len(after) != 64000 {
		r999ItemInvalid("frame_size")
		return r999ItemDiff{}
	}
	x0, y0, x1, y1, count := 320, 200, -1, -1, 0
	for i, v := range after {
		if v == before[i] {
			continue
		}
		x, y := i%320, i/320
		if x < x0 {
			x0 = x
		}
		if x > x1 {
			x1 = x
		}
		if y < y0 {
			y0 = y
		}
		if y > y1 {
			y1 = y
		}
		count++
	}
	if count == 0 {
		return r999ItemDiff{}
	}
	w, h := x1-x0+1, y1-y0+1
	b, a, m := make([]byte, w*h), make([]byte, w*h), make([]byte, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i, j := (y0+y)*320+x0+x, y*w+x
			b[j], a[j] = before[i], after[i]
			if b[j] != a[j] {
				m[j] = 1
			}
		}
	}
	r999ItemWrite(prefix+"-before.bin", b)
	r999ItemWrite(prefix+"-after.bin", a)
	r999ItemWrite(prefix+"-mask.bin", m)
	return r999ItemDiff{x0, y0, w, h, count}
}

func r999ObserveItem(o *oracle.Oracle, ret oracle.Addr, r oracle.Regs, font, off uint16, x, y int, prefix []byte) {
	if font != 0 || x != 18 || y != 114 {
		return
	}
	if os.Getenv("WIZ7_R999_ITEM_PROBE_DIR") == "" || r999ItemFailed || r999ItemCount >= 64 {
		return
	}
	if r.CS != 0x0110 || ret.Seg != r.CS || r.SP > 65530 || !r999ItemSpan(oracle.Far(r.SS, r.SP), 6) || !r999ItemSpan(oracle.Far(r.DS, off), 10) || len(prefix) < 8 {
		return
	}
	// 既有partRaw最多8B；只在完整已讀前綴命中後，逐byte核D和NUL，從不越過NUL讀鄰文。
	const head = "LONGSWOR"
	for i := range head {
		if prefix[i] != head[i] {
			return
		}
	}
	suffix := o.Bytes(oracle.Far(r.DS, off+8), 1)
	if len(suffix) != 1 || suffix[0] != 'D' {
		return
	}
	end := o.Bytes(oracle.Far(r.DS, off+9), 1)
	if len(end) != 1 || end[0] != 0 {
		return
	}
	before, screen := append([]byte(nil), offscreen(o)...), append([]byte(nil), o.Indexed()...)
	if len(before) != 64000 || len(screen) != 64000 {
		r999ItemInvalid("entry_frame_size")
		return
	}
	r999ItemCount++
	r999ServiceArm(o, r999ItemCount, o.Steps())
	p := r999ItemWitness{ID: r999ItemCount, English: "LONGSWORD", EntryRet: ret, EntryRegs: r, Font: font, SourceOff: off, X: x, Y: y, EntryStep: o.Steps()}
	name := fmt.Sprintf("longsword-%03d", p.ID)
	r999ItemJSON(name+"-entry.json", p)
	pending := true
	o.OnCall(ret, func(o *oracle.Oracle) {
		if !pending {
			return
		}
		pending = false
		p.ReturnRegs = o.Regs()
		p.ReturnStep = o.Steps()
		if p.ReturnRegs.SP != r.SP+2 || p.ReturnRegs.SS != r.SS || p.ReturnStep < p.EntryStep || p.ReturnStep-p.EntryStep > 1000000 {
			p.Reason = "return_frame_or_step"
			r999ItemJSON(name+"-return.json", p)
			return
		}
		after, visible := offscreen(o), o.Indexed()
		if len(after) != 64000 || len(visible) != 64000 {
			r999ItemInvalid("return_frame_size")
			return
		}
		p.ReturnAccepted = true
		p.Offscreen = r999ItemDiffROI(name+"-offscreen", before, after)
		p.Indexed = r999ItemDiffROI(name+"-indexed", screen, visible)
		pal := o.Palette()
		palette := make([]byte, 0, 768)
		for _, rgb := range pal {
			palette = append(palette, rgb[:]...)
		}
		r999ItemWrite(name+"-palette.bin", palette)
		r999ItemJSON(name+"-return.json", p)
	})
}
