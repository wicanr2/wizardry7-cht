package main

// spec038: only the observed Janette event. English selects display data;
// catalog IDs are provenance labels and never original runtime identities.
import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
)

var (
	npcTextSite     = oracle.Far(0x0110, 0xE3A6)
	npcRedrawSite   = oracle.Far(0x0110, 0x5F1E)
	npcTextGeometry = flowSite{width: 184, x0: 68, x1: 251, centerMin: 317, centerMax: 320}
)

type npcSignature struct {
	text string
	x, y int
}

type npcDefinition struct {
	rows []npcSignature
	// nil means the complete joined English key; the final paragraph has two sources.
	sources []string
}

var npcDefinitions = [...]npcDefinition{
	{rows: []npcSignature{{"Suddenly, from out of the blue,", 69, 146}, {"descends a flying red rocket...", 75, 158}}},
	{rows: []npcSignature{{"Upon the back of which rides a", 68, 146}, {"blazing female warrior...", 91, 158}}},
	{rows: []npcSignature{{"O'Haio, strangers!", 106, 146}}},
	{rows: []npcSignature{{"Welcome to the lands of", 91, 146}, {"Guardia!", 133, 158}}},
	{rows: []npcSignature{{"Many are you, who have come", 73, 146}, {"from the stars...", 114, 158}}},
	{rows: []npcSignature{{"We hail you, for it is the time", 73, 146}, {"of the coming of the", 101, 158}, {"Crusaders!", 128, 170}}},
	{rows: []npcSignature{{"Let destiny prevail!", 103, 146}}},
	{rows: []npcSignature{{"...And with that the bizarre", 80, 146}, {"lass zooms into the sky and", 78, 158}, {"shoots out of sight...", 102, 170}},
		sources: []string{"...And with that the bizarre lass zooms into the sky and shoots out of", "sight..."}},
}

type npcTranslation struct {
	key string
	res flowResult
}

type npcFlow struct {
	translations [len(npcDefinitions)]npcTranslation
	current      *npcParagraph
	reasons      map[string]bool
}

type npcParagraph struct {
	definition int
	page       *flowPage
	rows       []*npcRow
}

// template owns immutable original bytes. Events may be retired or replaced;
// neither their lifetime flags nor later original frames mutate this template.
type npcRow struct {
	template textEv
	event    *textEv
}

type pendingNPC struct {
	paragraph *npcParagraph
	row       int
	redraw    bool
}

func npcWords(s string) string { return strings.Join(strings.Fields(s), " ") }

func npcDefinitionKey(def npcDefinition) string {
	parts := make([]string, len(def.rows))
	for i, row := range def.rows {
		parts[i] = row.text
	}
	return npcWords(strings.Join(parts, " "))
}

func (nf *npcFlow) note(reason string) {
	if nf == nil || nf.reasons[reason] {
		return
	}
	nf.reasons[reason] = true
	log.Printf("珍妮特流文：%s", reason)
}

func newNPCFlow(fd *flowData) *npcFlow {
	nf := &npcFlow{reasons: map[string]bool{}}
	if fd == nil {
		return nf
	}
	wanted := map[string]bool{}
	for _, def := range npcDefinitions {
		keys := def.sources
		if len(keys) == 0 {
			keys = []string{npcDefinitionKey(def)}
		}
		for _, key := range keys {
			wanted[npcWords(key)] = true
		}
	}
	index := map[string][]int{}
	for id, msg := range fd.Msgs {
		key := npcWords(msg.W)
		if wanted[key] {
			index[key] = append(index[key], id)
		}
	}
	for key := range index {
		sort.Ints(index[key])
	}
	for i, def := range npcDefinitions {
		key := npcDefinitionKey(def)
		keys := def.sources
		if len(keys) == 0 {
			keys = []string{key}
		}
		if npcWords(strings.Join(keys, " ")) != key {
			nf.note(fmt.Sprintf("段%d來源英文不符", i+1))
			continue
		}
		zh := ""
		var labels []int
		valid := true
		for _, source := range keys {
			ids := index[npcWords(source)]
			if len(ids) == 0 {
				valid = false
				break
			}
			part := fd.Msgs[ids[0]].Z
			if part == "" {
				valid = false
				break
			}
			for _, id := range ids {
				if fd.Msgs[id].Z != part {
					valid = false
				}
			}
			zh += part
			labels = append(labels, ids...)
		}
		if !valid {
			nf.note(fmt.Sprintf("段%d來源缺少或歧義", i+1))
			continue
		}
		if !npcGlyphsValid(fd, zh) {
			nf.note(fmt.Sprintf("段%d缺字形", i+1))
			continue
		}
		nf.translations[i] = npcTranslation{key: key, res: flowResult{Zh: zh, Lines: labels, Cands: 1, Status: "NPC英文鍵"}}
	}
	return nf
}

func npcGlyphsValid(fd *flowData, zh string) bool {
	if fd == nil || (fd.size() != 15 && fd.size() != 16) || zh == "" {
		return false
	}
	for _, ch := range zh {
		rows, ok := fd.Glyphs[string(ch)]
		if !ok || len(rows) != fd.size() {
			return false
		}
		for _, row := range rows {
			if uint32(row) >= uint32(1)<<fd.size() {
				return false
			}
		}
	}
	return true
}

func npcCaller(ret oracle.Addr) bool { return ret == npcTextSite || ret == npcRedrawSite }

func npcString(raw []byte, off uint16) (string, bool) {
	if len(raw) == 0 || len(raw) > min(64, 65536-int(off)) {
		return "", false
	}
	n := bytes.IndexByte(raw, 0)
	if n < 0 {
		return "", false
	}
	s := stripCtl(string(raw[:n]))
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e {
			return "", false
		}
	}
	return s, true
}

func npcSignatureAt(s string, x, y int) (int, int, bool) {
	for i, def := range npcDefinitions {
		for j, row := range def.rows {
			if row.text == s && row.x == x && row.y == y {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func (nf *npcFlow) invalidateIncomplete(reason string) {
	if nf == nil {
		return
	}
	nf.note(reason)
	if current := nf.current; current != nil && len(current.rows) < len(npcDefinitions[current.definition].rows) {
		current.page.failed = true
		nf.current = nil
	}
}

// Called inside the existing NUL arbiter. Blank/unknown calls still use generic;
// no reader, hook, original field or write is introduced here.
func (ov *overlay) npcEntry(raw []byte, off, sp, font uint16, ret oracle.Addr, x, y int) *pendingPrint {
	nf := ov.npc
	if nf == nil || !npcCaller(ret) {
		return nil
	}
	if font != flowFont {
		return nil
	}
	s, bounded := npcString(raw, off)
	if !bounded {
		if x >= npcTextGeometry.x0 && x <= npcTextGeometry.x1 && (y == 146 || y == 158 || y == 170) {
			nf.invalidateIncomplete("內容邊界不符")
		}
		return nil
	}
	if len(strings.Fields(s)) == 0 {
		return nil
	}
	definition, row, known := npcSignatureAt(s, x, y)
	if !known {
		if x >= npcTextGeometry.x0 && x <= npcTextGeometry.x1 && (y == 146 || y == 158 || y == 170) {
			nf.invalidateIncomplete("未知正文行")
		}
		return nil
	}
	if sp > 65530 || ov.pend != nil {
		nf.invalidateIncomplete("返回或巢狀不符")
		return nil
	}
	translation := nf.translations[definition]
	if translation.res.Zh == "" {
		nf.invalidateIncomplete("來源不可用")
		return nil
	}
	if len(flowLayout(translation.res.Zh, ov.flow.cols(npcTextGeometry, flowK, flowGS(flowK, flowHdf)))) > flowMaxRows {
		nf.invalidateIncomplete("中文超列")
		return nil
	}
	current := nf.current
	if ret == npcTextSite {
		if row == 0 {
			current = &npcParagraph{definition: definition, page: &flowPage{site: npcTextSite, npc: true}}
			nf.current = current
		} else if current == nil || current.page.failed || current.definition != definition || row != len(current.rows) {
			nf.invalidateIncomplete("正文行序不符")
			return nil
		}
	} else if current == nil || current.page.failed || current.definition != definition || row >= len(current.rows) {
		nf.invalidateIncomplete("重畫沒有當前來源行")
		return nil
	}
	return &pendingPrint{s: s, x0: x, y0: y, sp: sp, ret: ret, font: font,
		npc: &pendingNPC{paragraph: current, row: row, redraw: ret == npcRedrawSite}}
}

func npcEventStyle(e *textEv, signature npcSignature) bool {
	valid := e != nil && e.originFont == flowFont && e.s == signature.text && e.x0 == signature.x && e.y0 == signature.y &&
		e.fg == flowFg && e.bw > 0 && e.bh > 0 && e.bx >= npcTextGeometry.x0 && e.bx+e.bw-1 <= npcTextGeometry.x1 &&
		e.by >= signature.y+1 && e.by+e.bh-1 <= signature.y+flowRowStep &&
		2*e.bx+e.bw-1 >= npcTextGeometry.centerMin && 2*e.bx+e.bw-1 <= npcTextGeometry.centerMax &&
		len(e.before) == e.bw*e.bh && len(e.after) == e.bw*e.bh
	if !valid {
		return false
	}
	changed := false
	for i, before := range e.before {
		if before != e.after[i] {
			if e.after[i] != flowFg {
				return false
			}
			changed = true
		}
	}
	return changed
}

func cloneNPCEvent(source *textEv, page *flowPage, ret oracle.Addr, rgb [3]uint8) *textEv {
	return &textEv{originRet: ret, originFont: source.originFont, s: source.s, x0: source.x0, y0: source.y0,
		bx: source.bx, by: source.by, bw: source.bw, bh: source.bh, before: append([]uint8(nil), source.before...),
		after: append([]uint8(nil), source.after...), fg: source.fg, fgRGB: rgb, page: page}
}

func (ov *overlay) rejectNPC(p *pendingPrint, reason string) {
	if p == nil || p.npc == nil || ov.npc == nil {
		return
	}
	ov.npc.note(reason)
	current := p.npc.paragraph
	if len(current.rows) < len(npcDefinitions[current.definition].rows) {
		current.page.failed = true
		if ov.npc.current == current {
			ov.npc.current = nil
		}
	}
}

// A nil candidate means the original call produced no new differences. Only a
// matching current redraw can reconstruct its immutable row, from exact after.
func (ov *overlay) finishNPC(p *pendingPrint, candidate *textEv, off []uint8, rgb [3]uint8, step uint64) bool {
	if p == nil || p.npc == nil || ov.npc == nil {
		return false
	}
	pending := p.npc
	current := pending.paragraph
	if ov.npc.current != current || current.page.failed || len(off) != 320*200 {
		ov.rejectNPC(p, "返回來源不符")
		return false
	}
	signature := npcDefinitions[current.definition].rows[pending.row]
	var row *npcRow
	if pending.redraw {
		if pending.row >= len(current.rows) {
			ov.rejectNPC(p, "重畫行不存在")
			return false
		}
		row = current.rows[pending.row]
		original := &row.template
		if candidate == nil {
			if !original.rectMatches(off, nil) {
				ov.rejectNPC(p, "無差異重畫背景不符")
				return false
			}
		} else if !npcEventStyle(candidate, signature) || candidate.bx != original.bx || candidate.by != original.by ||
			candidate.bw != original.bw || candidate.bh != original.bh ||
			!bytes.Equal(candidate.before, original.before) || !bytes.Equal(candidate.after, original.after) {
			ov.rejectNPC(p, "重畫原色不符")
			return false
		}
	} else {
		if pending.row != len(current.rows) || !npcEventStyle(candidate, signature) {
			ov.rejectNPC(p, "正文差異或樣式不符")
			return false
		}
		template := cloneNPCEvent(candidate, nil, npcTextSite, rgb)
		row = &npcRow{template: *template}
		current.rows = append(current.rows, row)
		current.page.words = append(current.page.words, strings.Fields(signature.text)...)
		current.page.frags = len(current.rows)
		current.page.lastY = signature.y
		current.page.lastStep = step // only new E3 content resets quiet
		current.page.evs = append(current.page.evs, nil)
	}
	pg := current.page
	if len(current.rows) == len(npcDefinitions[current.definition].rows) &&
		npcWords(strings.Join(pg.words, " ")) != ov.npc.translations[current.definition].key {
		pg.failed = true
		ov.npc.note("整段英文不符")
		return false
	}
	event := cloneNPCEvent(&row.template, pg, p.ret, rgb)
	if previous := row.event; previous != nil {
		kept := ov.evs[:0]
		for _, e := range ov.evs {
			if e != previous {
				kept = append(kept, e)
			}
		}
		for i := len(kept); i < len(ov.evs); i++ {
			ov.evs[i] = nil
		}
		ov.evs = kept
		previous.dead = true
	}
	row.event = event
	pg.evs[pending.row] = event
	ov.evs = append(ov.evs, event)
	ov.created++
	pg.band = snapBand(off, npcTextGeometry)
	if len(current.rows) == len(npcDefinitions[current.definition].rows) {
		translation := ov.npc.translations[current.definition]
		pg.res = translation.res
		pg.rowsOK = len(flowLayout(pg.res.Zh, ov.flow.cols(npcTextGeometry, flowK, flowGS(flowK, flowHdf)))) <= flowMaxRows
	}
	attached := false
	for _, page := range ov.pages {
		if page == pg {
			attached = true
		}
	}
	if !attached {
		ov.pages = append(ov.pages, pg)
	}
	return true
}

// Called by every compose snapshot, not by a particular redraw callback.
func (pg *flowPage) npcReady(a []uint8, step uint64, inCursor func(x, y int) bool) (bool, string) {
	if !pg.npc || pg.failed || pg.res.Zh == "" || !pg.rowsOK || pg.frags == 0 || len(pg.evs) != pg.frags {
		return false, "頁未完整"
	}
	if step < pg.lastStep || step-pg.lastStep < flowQuiet {
		return false, "等待正文發布"
	}
	for _, e := range pg.evs {
		if e == nil || e.dead || !e.match || !e.rectMatches(a, inCursor) {
			return false, "原色未完整發布"
		}
	}
	if !pg.bandStatic(a, inCursor) {
		return false, "文字帶背景不符"
	}
	return true, ""
}
