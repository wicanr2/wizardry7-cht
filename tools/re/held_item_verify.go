package main

// 私人驗收：沿已審原API與owner快照，輸出有界實際Draw及原版游標資料。
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/oracle"
)

var formalSeen = map[string]bool{}
var formalSeenProof = map[uint64]bool{}
var formalHadCursor bool
var formalReset int

func formalSHA(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func formalWrite(name string, b []byte) {
	root := os.Getenv("WIZ7_FORMAL_OBSERVE")
	if root == "" {
		return
	}
	f, e := os.OpenFile(filepath.Join(root, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		panic("FORMAL_OBSERVER_IO")
	}
	n, e := f.Write(b)
	e2 := f.Close()
	if e != nil || e2 != nil || n != len(b) {
		panic("FORMAL_OBSERVER_IO")
	}
}
func formalJSON(name string, x any) {
	b, e := json.Marshal(x)
	if e != nil {
		panic("FORMAL_OBSERVER_JSON")
	}
	formalWrite(name, append(b, '\n'))
}
func formalBefore(o *oracle.Oracle, ov *overlay) {
	if os.Getenv("WIZ7_FORMAL_OBSERVE") == "" || ov.held == nil || ov.held.label == nil || len(formalSeenProof) >= 8 {
		return
	}
	l := ov.held.label
	if formalSeenProof[l.defined] {
		return
	}
	formalSeenProof[l.defined] = true
	name := fmt.Sprintf("proof-%d", l.defined)
	files := map[string]string{}
	for i, p := range l.source.parts {
		for role, b := range map[string][]byte{"before": p.before, "after": p.after, "mask": p.mask} {
			file := fmt.Sprintf("%s-part%d-%s.bin", name, i, role)
			formalWrite(file, b)
			files[fmt.Sprintf("part%d-%s", i, role)] = file
		}
	}
	formalWrite(name+"-picture.bin", l.picture)
	formalWrite(name+"-mask.bin", l.mask)
	formalJSON(name+".json", map[string]any{"key": l.key, "ptr": l.ptr, "source_rect": l.source.rect, "source_step": l.source.parts[0].entry, "captured_step": l.captured, "defined_step": l.defined, "picture_sha256": formalSHA(l.picture), "files": files, "module_valid": heldModule(o)})
}
func formalObserve(o *oracle.Oracle, ov *overlay, sh *shared, l *inputRecording, k int) {
	root := os.Getenv("WIZ7_FORMAL_OBSERVE")
	if root == "" || l == nil || sh.helpOpen || sh.mapOpen || len(formalSeen) >= 16 || ov.held.definition != nil {
		return
	}
	driver, ok := heldDriver(o)
	if !ok {
		return
	}
	full := heldRead(o, oracle.Far(driver, 0x1e5c), 12)
	desc := heldRead(o, oracle.Far(driver, 0x15ac), 24)
	busy := heldRead(o, oracle.Far(driver, 0x1ed3), 1)
	if len(full) != 12 || len(desc) != 24 || len(busy) != 1 || busy[0] != 0 {
		return
	}
	pub := l.publication + 1
	step := o.Steps()
	if heldWord(full, 0) != 76 || heldWord(full, 2) != 13 || heldZero(heldPtr(full, 8)) {
		if !formalHadCursor {
			return
		}
		formalHadCursor = false
		formalReset++
		tag := fmt.Sprintf("reset-%d", formalReset)
		formalSeen[tag] = true
		sh.drawReq = filepath.Join(root, tag+"-draw.png")
		formalJSON(tag+".json", map[string]any{"status": "reset", "pub": pub, "step": step, "label_present": ov.held.label != nil, "view_present": sh.held != nil})
		return
	}
	ptr := heldPtr(full, 8)
	picture := heldRead(o, ptr, heldPixels)
	if len(picture) != heldPixels {
		return
	}
	sha := formalSHA(picture)
	if sha != "98e96e235de0701c61b2629c6a78bc6ba1f49fced39efad81e094d8ca498e863" && sha != "5922e142ae731e2e041ea1f0aa16cb869e345bd60b00f954f413b4479ffa532c" {
		return
	}
	x, y := int(int16(heldWord(desc, 0))), int(int16(heldWord(desc, 2)))
	rect := heldClipRect(x, y)
	if rect.Empty() || int(heldWord(desc, 12)) != rect.Dx() || int(heldWord(desc, 14)) != rect.Dy() {
		return
	}
	bg := heldRead(o, oracle.Far(driver, heldWord(desc, 8)), rect.Dx()*rect.Dy())
	if len(bg) != rect.Dx()*rect.Dy() {
		return
	}
	a := o.Indexed()
	if len(a) != 64000 {
		return
	}
	indexed := make([]byte, len(bg))
	for yy := 0; yy < rect.Dy(); yy++ {
		copy(indexed[yy*rect.Dx():], a[(y+yy)*320+x:(y+yy)*320+x+rect.Dx()])
	}
	key := ""
	if sh.held != nil && sh.held.label != nil {
		key = sh.held.label.key
	}
	tag := fmt.Sprintf("%s-%d-%d-%s", sha[:8], x, y, key)
	if formalSeen[tag] {
		return
	}
	formalSeen[tag] = true
	formalHadCursor = true
	name := fmt.Sprintf("frame-%03d", len(formalSeen))
	sh.drawReq = filepath.Join(root, name+"-draw.png")
	for role, b := range map[string][]byte{"picture": picture, "background": bg, "indexed": indexed} {
		formalWrite(name+"-"+role+".bin", b)
	}
	pal := o.Palette()
	palette := make([]byte, 0, 768)
	for _, c := range pal {
		palette = append(palette, c[:]...)
	}
	formalWrite(name+"-palette.bin", palette)
	hdFrame(o, ov, filepath.Join(root, name+"-hd.png"))
	formalJSON(name+".json", map[string]any{"pub": pub, "step": step, "scale": k, "x": x, "y": y, "w": rect.Dx(), "h": rect.Dy(), "picture_sha256": sha, "key": key, "translated": sh.held != nil && sh.held.translated(k), "module_valid": heldModule(o), "rejected_band": ov.held.rejectedBand, "rejected_capture": ov.held.rejectedCapture, "rejected_definition": ov.held.rejectedDefinition})
}
