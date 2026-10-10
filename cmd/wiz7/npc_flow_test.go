package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

// These frames exercise adapter state only. They are not original-run parity.
func npcTestData() *flowData {
	fd := &flowData{Size: 16, Glyphs: map[string][]uint16{"甲": make([]uint16, 16)}, Msgs: map[int]flowMessage{}}
	id := 1
	for _, def := range npcDefinitions {
		keys := def.sources
		if len(keys) == 0 {
			keys = []string{npcDefinitionKey(def)}
		}
		for _, key := range keys {
			fd.Msgs[id] = flowMessage{W: key, Z: "甲"}
			id++
		}
	}
	return fd
}

func npcTestOverlay() *overlay {
	fd := npcTestData()
	return &overlay{flow: fd, npc: newNPCFlow(fd)}
}

func npcTestFrame() []byte {
	a := make([]byte, 320*200)
	for i := range a {
		a[i] = byte(1 + i%17)
	}
	return a
}

func npcTestEvent(signature npcSignature) *textEv {
	e := &textEv{originRet: npcTextSite, originFont: flowFont, s: signature.text, x0: signature.x, y0: signature.y,
		bx: 150, by: signature.y + 2, bw: 20, bh: 3, fg: flowFg, before: make([]byte, 60), after: make([]byte, 60)}
	for i := range e.before {
		e.before[i] = byte(1 + i%17)
		e.after[i] = e.before[i]
	}
	e.after[0], e.after[39], e.after[59] = flowFg, flowFg, flowFg
	return e
}

func npcTestPaint(a []byte, e *textEv, after bool) {
	data := e.before
	if after {
		data = e.after
	}
	for y := 0; y < e.bh; y++ {
		copy(a[(e.by+y)*320+e.bx:], data[y*e.bw:(y+1)*e.bw])
	}
}

func npcTestEntry(ov *overlay, definition, row int, ret oracle.Addr) *pendingPrint {
	sig := npcDefinitions[definition].rows[row]
	return ov.npcEntry(append([]byte(sig.text), 0), 0x200, 0x100, flowFont, ret, sig.x, sig.y)
}

func npcTestParagraph(t *testing.T, ov *overlay, definition int) (*flowPage, []byte) {
	t.Helper()
	a := npcTestFrame()
	for i, sig := range npcDefinitions[definition].rows {
		p := npcTestEntry(ov, definition, i, npcTextSite)
		if p == nil {
			t.Fatal("known entry rejected", definition, i)
		}
		e := npcTestEvent(sig)
		npcTestPaint(a, e, true)
		if !ov.finishNPC(p, e, a, [3]uint8{255, 247, 0}, 100+uint64(i)) {
			t.Fatal("known return rejected")
		}
	}
	pg := ov.npc.current.page
	for _, e := range pg.evs {
		e.match = true
	}
	return pg, a
}

func TestNPCCatalogUsesEnglishNotIDs(t *testing.T) {
	fd := npcTestData()
	fd.Msgs[95001] = fd.Msgs[1]
	nf := newNPCFlow(fd)
	if nf.translations[0].res.Zh != "甲" || !reflect.DeepEqual(nf.translations[0].res.Lines, []int{1, 95001}) {
		t.Fatal("same-key provenance")
	}
	if nf.translations[7].res.Zh != "甲甲" {
		t.Fatal("final two parts not joined")
	}
	shifted := npcTestData()
	shifted.Msgs = map[int]flowMessage{}
	for id, m := range fd.Msgs {
		m.W = "  " + strings.ReplaceAll(m.W, " ", "\n\t") + " "
		shifted.Msgs[id+100] = m
	}
	if newNPCFlow(shifted).translations[7].res.Zh != "甲甲" {
		t.Fatal("ID or whitespace selected translation")
	}
	for _, test := range []struct {
		name   string
		change func(*flowData)
	}{
		{"ambiguous", func(fd *flowData) { fd.Msgs[999] = flowMessage{W: fd.Msgs[1].W, Z: "甲甲"} }},
		{"missing", func(fd *flowData) { delete(fd.Msgs, 1) }},
		{"empty", func(fd *flowData) { m := fd.Msgs[1]; m.Z = ""; fd.Msgs[1] = m }},
		{"punctuation", func(fd *flowData) { m := fd.Msgs[1]; m.W = strings.TrimSuffix(m.W, "."); fd.Msgs[1] = m }},
		{"missing-glyph", func(fd *flowData) { delete(fd.Glyphs, "甲") }},
		{"short-glyph", func(fd *flowData) { fd.Glyphs["甲"] = []uint16{1} }},
		{"non-native", func(fd *flowData) { fd.Size = 14 }},
	} {
		t.Run(test.name, func(t *testing.T) {
			fd := npcTestData()
			test.change(fd)
			if newNPCFlow(fd).translations[0].res.Zh != "" {
				t.Fatal("unsafe source enabled")
			}
		})
	}
	fd = npcTestData()
	delete(fd.Msgs, 9)
	if newNPCFlow(fd).translations[7].res.Zh != "" || newNPCFlow(fd).translations[0].res.Zh == "" {
		t.Fatal("missing final part scope")
	}
	fd = npcTestData()
	fd.Size = 15
	fd.Glyphs["甲"] = make([]uint16, 15)
	fd.Glyphs["甲"][0] = 0x8000
	if newNPCFlow(fd).translations[0].res.Zh != "" {
		t.Fatal("out-of-width bit accepted")
	}
}

func TestNPCEntryBoundsAndGenericFallback(t *testing.T) {
	sig := npcDefinitions[0].rows[0]
	raw := append([]byte(sig.text), 0)
	for _, test := range []struct {
		name          string
		raw           []byte
		off, sp, font uint16
		ret           oracle.Addr
		x, y          int
	}{
		{"alias", raw, 0x200, 0x100, 2, oracle.Far(0x111, 0xe396), sig.x, sig.y},
		{"caller", raw, 0x200, 0x100, 2, storySite, sig.x, sig.y},
		{"font", raw, 0x200, 0x100, 0, npcTextSite, sig.x, sig.y},
		{"x", raw, 0x200, 0x100, 2, npcTextSite, sig.x + 1, sig.y},
		{"y", raw, 0x200, 0x100, 2, npcTextSite, sig.x, sig.y + 1},
		{"stack", raw, 0x200, 65531, 2, npcTextSite, sig.x, sig.y},
		{"no-nul", []byte(sig.text), 0x200, 0x100, 2, npcTextSite, sig.x, sig.y},
		{"cross-offset", raw, 65530, 0x100, 2, npcTextSite, sig.x, sig.y},
		{"over-64", append(make([]byte, 65), 0), 0x200, 0x100, 2, npcTextSite, sig.x, sig.y},
		{"high-byte", []byte{0xff, 0}, 0x200, 0x100, 2, npcTextSite, sig.x, sig.y},
	} {
		t.Run(test.name, func(t *testing.T) {
			ov := npcTestOverlay()
			if ov.npcEntry(test.raw, test.off, test.sp, test.font, test.ret, test.x, test.y) != nil {
				t.Fatal("non-signature claimed")
			}
		})
	}
	ov := npcTestOverlay()
	if ov.npcEntry(raw, uint16(65536-len(raw)), 65530, 2, npcTextSite, sig.x, sig.y) == nil {
		t.Fatal("valid bounded tail rejected")
	}
	if _, ok := npcString([]byte{0}, 65535); !ok {
		t.Fatal("last-byte NUL")
	}
	if _, ok := npcString([]byte{'x'}, 65535); ok {
		t.Fatal("unterminated tail")
	}
	if _, ok := flowSites[npcTextSite]; ok {
		t.Fatal("NPC entered generic flow router")
	}
	if flowGeometry(storySite) != flowSites[storySite] {
		t.Fatal("story geometry changed")
	}
}

func TestNPCUnknownInvalidatesOnlyIncompleteAssembly(t *testing.T) {
	ov := npcTestOverlay()
	p := npcTestEntry(ov, 0, 0, npcTextSite)
	e := npcTestEvent(npcDefinitions[0].rows[0])
	a := npcTestFrame()
	npcTestPaint(a, e, true)
	if !ov.finishNPC(p, e, a, [3]uint8{}, 100) {
		t.Fatal("first row")
	}
	current := ov.npc.current
	for _, raw := range [][]byte{{0}, []byte(" \t\x00")} {
		if ov.npcEntry(raw, 1, 0x100, 2, npcRedrawSite, 69, 146) != nil || ov.npc.current != current {
			t.Fatal("blank cleared cache")
		}
	}
	ov.npcEntry([]byte("another NPC\x00"), 1, 0x100, 0, npcTextSite, 69, 146)
	ov.npcEntry([]byte("another NPC\x00"), 1, 0x100, 2, npcTextSite, 5, 146)
	ov.npcEntry([]byte("unterminated"), 1, 0x100, 0, npcTextSite, 69, 146)
	ov.npcEntry([]byte("unterminated"), 1, 0x100, 2, npcTextSite, 5, 146)
	if ov.npc.current != current {
		t.Fatal("unrelated call cleared cache")
	}
	ov.npcEntry([]byte("another NPC\x00"), 1, 0x100, 2, npcTextSite, 69, 146)
	if ov.npc.current != nil || !current.page.failed || npcTestEntry(ov, 0, 1, npcTextSite) != nil {
		t.Fatal("mixed assembly survived")
	}
	pg, _ := npcTestParagraph(t, ov, 2)
	current = ov.npc.current
	ov.npcEntry([]byte("another NPC\x00"), 1, 0x100, 2, npcRedrawSite, 69, 146)
	if ov.npc.current != current || pg.failed {
		t.Fatal("published old page erased early")
	}
}

func TestNPCOrderAndStyleGuards(t *testing.T) {
	ov := npcTestOverlay()
	if npcTestEntry(ov, 0, 1, npcTextSite) != nil || npcTestEntry(ov, 0, 0, npcRedrawSite) != nil {
		t.Fatal("orphan line accepted")
	}
	p := npcTestEntry(ov, 0, 0, npcTextSite)
	ov.pend = p
	if npcTestEntry(ov, 0, 1, npcTextSite) != nil || !p.npc.paragraph.page.failed {
		t.Fatal("nested print accepted")
	}
	for _, change := range []struct {
		name string
		fn   func(*textEv)
	}{
		{"font", func(e *textEv) { e.originFont = 0 }},
		{"color", func(e *textEv) { e.fg = 86 }},
		{"minor-color", func(e *textEv) { e.after[0] = 86 }},
		{"center", func(e *textEv) { e.bx++ }},
		{"vertical", func(e *textEv) { e.by = 145 }},
		{"nodiff", func(e *textEv) { copy(e.after, e.before) }},
		{"extent", func(e *textEv) { e.after = e.after[:1] }},
	} {
		t.Run(change.name, func(t *testing.T) {
			ov := npcTestOverlay()
			p := npcTestEntry(ov, 0, 0, npcTextSite)
			e := npcTestEvent(npcDefinitions[0].rows[0])
			change.fn(e)
			if ov.finishNPC(p, e, npcTestFrame(), [3]uint8{}, 100) || !p.npc.paragraph.page.failed {
				t.Fatal("style failure retained assembly")
			}
		})
	}
	ov = npcTestOverlay()
	p = npcTestEntry(ov, 0, 0, npcTextSite)
	e := npcTestEvent(npcDefinitions[0].rows[0])
	a := npcTestFrame()
	npcTestPaint(a, e, true)
	ov.finishNPC(p, e, a, [3]uint8{}, 100)
	if npcTestEntry(ov, 5, 1, npcTextSite) != nil || ov.npc.current != nil {
		t.Fatal("mixed known rows accepted")
	}
}

func TestNPCPublishRequiresWholeAfterAndBackground(t *testing.T) {
	ov := npcTestOverlay()
	pg, a := npcTestParagraph(t, ov, 0)
	none := func(int, int) bool { return false }
	if ok, _ := pg.npcReady(a, pg.lastStep+flowQuiet-1, none); ok {
		t.Fatal("quiet ignored")
	}
	if ok, _ := pg.npcReady(a, pg.lastStep-1, none); ok {
		t.Fatal("backward step accepted")
	}
	// No 5F callback occurred: an ordinary later publication is sufficient.
	if ok, why := pg.npcReady(a, pg.lastStep+flowQuiet, none); !ok {
		t.Fatal("ordinary blit not eligible", why)
	}
	e := pg.evs[1]
	i := e.by*320 + e.bx
	beforePublish := append([]byte(nil), a...)
	npcTestPaint(beforePublish, e, false)
	if ok, _ := pg.npcReady(beforePublish, pg.lastStep+flowQuiet, none); ok {
		t.Fatal("prepublication accepted")
	}
	partial := append([]byte(nil), a...)
	partial[i] = e.before[0]
	if !e.pageMatch(partial, none) {
		t.Fatal("fixture should pass old partial rule")
	}
	if ok, _ := pg.npcReady(partial, pg.lastStep+flowQuiet, none); ok {
		t.Fatal("NPC used partial rule")
	}
	partial[i] = 86
	if ok, _ := pg.npcReady(partial, pg.lastStep+flowQuiet, none); ok {
		t.Fatal("wrong stroke color accepted")
	}
	background := append([]byte(nil), a...)
	background[194*320+68]++
	if ok, _ := pg.npcReady(background, pg.lastStep+flowQuiet, none); ok {
		t.Fatal("outside-row band drift accepted")
	}
	background = append([]byte(nil), a...)
	background[i+1]++
	if ok, _ := pg.npcReady(background, pg.lastStep+flowQuiet, none); ok {
		t.Fatal("row background drift accepted")
	}
	if ok, _ := pg.npcReady(partial, pg.lastStep+flowQuiet, func(x, y int) bool { return x == e.bx && y == e.by }); !ok {
		t.Fatal("existing cursor exclusion lost")
	}
	old := pg
	p := npcTestEntry(ov, 1, 0, npcTextSite)
	if p == nil || old.failed {
		t.Fatal("new pending erased old page")
	}
	if ok, _ := old.npcReady(a, old.lastStep+flowQuiet, none); !ok {
		t.Fatal("old native page no longer valid")
	}
	old.evs[0].dead = true
	if ok, _ := old.npcReady(a, old.lastStep+flowQuiet, none); ok {
		t.Fatal("retired page resurrected without redraw")
	}
}

func TestNPCCacheReplacementAndNoDiff(t *testing.T) {
	ov := npcTestOverlay()
	pg, a := npcTestParagraph(t, ov, 0)
	words := append([]string(nil), pg.words...)
	quiet := pg.lastStep
	old := pg.evs[0]
	p := npcTestEntry(ov, 0, 0, npcRedrawSite)
	if !ov.finishNPC(p, npcTestEvent(npcDefinitions[0].rows[0]), a, [3]uint8{1, 2, 3}, 999) {
		t.Fatal("same cache rejected")
	}
	if len(ov.evs) != 2 || len(pg.evs) != 2 || pg.evs[0] == old || !old.dead || pg.lastStep != quiet || !reflect.DeepEqual(pg.words, words) {
		t.Fatal("cache appended or restarted quiet")
	}
	if pg.evs[0].fgRGB != [3]uint8{1, 2, 3} {
		t.Fatal("replacement current palette")
	}
	// Caller memory and returned event bytes cannot mutate the retained source.
	template := append([]byte(nil), ov.npc.current.rows[0].template.after...)
	pg.evs[0].after[0] = 0
	if !bytes.Equal(ov.npc.current.rows[0].template.after, template) {
		t.Fatal("template aliases event")
	}
	p = npcTestEntry(ov, 0, 0, npcRedrawSite)
	if !ov.finishNPC(p, nil, a, [3]uint8{}, 1000) || !bytes.Equal(pg.evs[0].after, template) {
		t.Fatal("exact NoDiff reconstruction")
	}
	bad := append([]byte(nil), a...)
	bad[pg.evs[0].by*320+pg.evs[0].bx+1]++
	p = npcTestEntry(ov, 0, 0, npcRedrawSite)
	if ov.finishNPC(p, nil, bad, [3]uint8{}, 1001) {
		t.Fatal("NoDiff accepted different background")
	}
	if pg.failed {
		t.Fatal("rejected redraw prematurely erased complete page")
	}
	e := npcTestEvent(npcDefinitions[0].rows[0])
	e.before[1]++
	p = npcTestEntry(ov, 0, 0, npcRedrawSite)
	if ov.finishNPC(p, e, a, [3]uint8{}, 1002) {
		t.Fatal("cache before bytes changed")
	}
	ov = npcTestOverlay()
	p = npcTestEntry(ov, 0, 0, npcTextSite)
	if ov.finishNPC(p, nil, npcTestFrame(), [3]uint8{}, 1) {
		t.Fatal("new NoDiff row created")
	}
}

func TestNPCActualR4FragmentsAndNativeCatalog(t *testing.T) {
	paths, err := filepath.Glob("../../workplace/r1007/npc/probe-r4-run1/style/npc-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Skip("缺本機R4原色收據")
	}
	if len(paths) != 31 {
		t.Fatal("R4 count", len(paths))
	}
	for _, backend := range []string{"unifont", "eten"} {
		for _, k := range []int{2, 3} {
			t.Run(fmt.Sprintf("%s-%d", backend, k), func(t *testing.T) {
				path := fmt.Sprintf("../../workplace/r983/final/%s/flow.json", backend)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Skip("缺本機正式字模資料")
				}
				fd := loadFlow(path)
				if fd == nil {
					t.Fatal("catalog load")
				}
				priorK, priorH := flowK, flowHdf
				flowK, flowHdf = k, 1
				t.Cleanup(func() { flowK, flowHdf = priorK, priorH })
				ov := &overlay{flow: fd, npc: newNPCFlow(fd)}
				for i, tr := range ov.npc.translations {
					if tr.res.Zh == "" {
						t.Fatal("actual paragraph disabled", i)
					}
				}
				var frame []byte
				newRows, cached, complete := 0, 0, map[int]bool{}
				for _, path := range paths {
					var f struct {
						Text          string
						X, Y          int
						Font          uint16
						Ret           oracle.Addr
						BBox          [4]int
						Before, After []byte
						Fg            byte
						FgRGB         [3]uint8
						Step          uint64
					}
					raw, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					if err = json.Unmarshal(raw, &f); err != nil {
						t.Fatal(err)
					}
					definition, row, known := npcSignatureAt(f.Text, f.X, f.Y)
					if !known {
						t.Fatal("R4 signature", path)
					}
					if f.Ret == npcTextSite && row == 0 {
						frame = npcTestFrame()
					}
					p := ov.npcEntry(append([]byte(f.Text), 0), 0x200, 0x100, f.Font, f.Ret, f.X, f.Y)
					if p == nil {
						t.Fatal("R4 entry rejected", path)
					}
					e := &textEv{originRet: f.Ret, originFont: f.Font, s: f.Text, x0: f.X, y0: f.Y, bx: f.BBox[0], by: f.BBox[1], bw: f.BBox[2], bh: f.BBox[3], before: f.Before, after: f.After, fg: f.Fg}
					npcTestPaint(frame, e, true)
					if !ov.finishNPC(p, e, frame, f.FgRGB, f.Step) {
						t.Fatal("R4 return rejected", path)
					}
					if f.Ret == npcTextSite {
						newRows++
					} else {
						cached++
					}
					pg := ov.npc.current.page
					for _, event := range pg.evs {
						event.match = true
					}
					ok, _ := pg.npcReady(frame, pg.lastStep+flowQuiet, func(int, int) bool { return false })
					if len(pg.evs) == len(npcDefinitions[definition].rows) {
						if !ok {
							t.Fatal("complete reconstructed adapter frame", path)
						}
						complete[definition] = true
					} else if ok {
						t.Fatal("partial paragraph published")
					}
					if len(pg.evs) != pg.frags || len(pg.words) != len(strings.Fields(strings.Join(func() []string {
						var s []string
						for _, r := range npcDefinitions[definition].rows[:pg.frags] {
							s = append(s, r.text)
						}
						return s
					}(), " "))) {
						t.Fatal("cache duplicated content")
					}
				}
				if newRows != 16 || cached != 15 || len(complete) != 8 || len(ov.evs) != 16 {
					t.Fatal("R4 accounting", newRows, cached, len(complete), len(ov.evs))
				}
			})
		}
	}
}
