package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/wicanr2/dosgolem/oracle"
)

// 位址與三組配對只適用研究021的正常角色頁；契約見READY spec012。
type compoundSite struct {
	first, second string
	y, glyphs     int
}

var compoundSites = [...]compoundSite{
	{"LEATHER", "CUIRASS", 128, 3},
	{"FUR", "LEGGING", 142, 4},
	{"BUCKLER", "SHIELD", 170, 3},
}

var compoundFirstReturn = oracle.Addr{Seg: 0x0110, Off: 0xba40}
var compoundSecondReturn = oracle.Addr{Seg: 0x0110, Off: 0xba63}

type compoundInput struct {
	S                 string
	DS, Off, SP, Font uint16
	Ret               oracle.Addr
	Step              uint64
	X, Y              int
	raw               []byte
}

type compoundHash struct {
	Length int
	SHA256 string
}

type compoundTrace struct {
	Inputs    [2]compoundInput
	Returns   [2]uint64
	Snapshots [4]compoundHash // before0、after1、before2、after2
	Masks     [2]compoundHash
	State     string
}

// 完整域比對同時涵蓋英文、行間隙與非筆畫背景；nil代表不容忍游標。
func (e *textEv) rectMatches(screen []byte, ignore func(int, int) bool) bool {
	if len(screen) != 320*200 || e.bx < 0 || e.by < 0 || e.bw <= 0 || e.bh <= 0 ||
		e.bx+e.bw > 320 || e.by+e.bh > 200 || len(e.after) != e.bw*e.bh {
		return false
	}
	for y := 0; y < e.bh; y++ {
		for x := 0; x < e.bw; x++ {
			px, py := e.bx+x, e.by+y
			if ignore != nil && ignore(px, py) {
				continue
			}
			if screen[py*320+px] != e.after[y*e.bw+x] {
				return false
			}
		}
	}
	return true
}

type compoundPair struct {
	site    compoundSite
	phase   int // 0第一行未返回；1等待下一入口；2第二行未返回
	inputs  [2]compoundInput
	returns [2]uint64
	snaps   [4][]byte
	masks   [2][]byte
	g       [][]uint16
}

func compoundDigest(b []byte) compoundHash {
	return compoundHash{len(b), fmt.Sprintf("%x", sha256.Sum256(b))}
}

func compoundRect(screen []byte, y int) []byte {
	if len(screen) != 320*200 || y < 0 || y+11 > 200 {
		return nil
	}
	b := make([]byte, 41*11)
	for row := 0; row < 11; row++ {
		copy(b[row*41:(row+1)*41], screen[(y+row)*320+18:(y+row)*320+59])
	}
	return b
}

func compoundTextOK(in compoundInput, want string) bool {
	n := len(want)
	return int(in.Off)+n <= 65535 && n < 8 && len(in.raw) > n &&
		in.raw[n] == 0 && string(in.raw[:n]) == want && in.S == want
}

func compoundGlyphs(e entry, count, scale int) ([][]uint16, bool) {
	if (scale != 2 && scale != 3) || e.Dy != 0 || len(e.D) != 0 || len(e.G) != count {
		return nil, false
	}
	size := len(e.G[0])
	if size != 15 && size != 16 || count*(size+2)-2 > 41*scale || size > 11*scale {
		return nil, false
	}
	g := make([][]uint16, count)
	for i, glyph := range e.G {
		if len(glyph) != size {
			return nil, false
		}
		var bits uint16
		for _, row := range glyph {
			if size == 15 && row&0x8000 != 0 {
				return nil, false
			}
			bits |= row
		}
		if bits == 0 {
			return nil, false
		}
		g[i] = append([]uint16(nil), glyph...)
	}
	return g, true
}

// 返回差異必須全部落在所屬行、只寫前景105；整域其他背景不得改動。
func compoundMask(before, after []byte, line int) []byte {
	if len(before) != 451 || len(after) != 451 || (line != 0 && line != 6) {
		return nil
	}
	mask, count := make([]byte, 451), 0
	for i, v := range after {
		if v == before[i] {
			continue
		}
		if i/41 < line || i/41 >= line+5 || v != 105 {
			return nil
		}
		mask[i], count = 1, count+1
	}
	if count == 0 {
		return nil
	}
	return mask
}

func (ov *overlay) rejectCompound(reason string) {
	if ov.compound != nil {
		log.Printf("複合物品 捨棄 %s/%s 原因=%s", ov.compound.site.first, ov.compound.site.second, reason)
		ov.compound = nil
	}
}

// 每次3895入口都先呼叫；take只在已知部件需要背景時讀離屏。
// 回true代表已知部件，失敗亦不得再套generic文案。
func (ov *overlay) enterCompound(in compoundInput, take func() []byte, trans map[string]entry, scale int) bool {
	old := ov.compound
	if old != nil {
		ov.compound = nil
		first := old.inputs[0]
		if old.phase == 1 && in.Ret == compoundSecondReturn && compoundTextOK(in, old.site.second) &&
			in.DS == first.DS && in.SP == first.SP && in.SP <= 65533 &&
			int(first.Off)+len(old.site.first)+1 <= 65535 &&
			int(in.Off) == int(first.Off)+len(old.site.first)+1 &&
			in.Font == 0 && in.X == 18 && in.Y == old.site.y+6 &&
			in.Step >= old.returns[0] && in.Step-old.returns[0] == 36 {
			before2 := compoundRect(take(), old.site.y)
			if before2 != nil && bytes.Equal(before2, old.snaps[1]) {
				old.phase, old.inputs[1], old.snaps[2] = 2, in, before2
				old.inputs[1].raw = nil
				ov.compound = old
				return true
			}
		}
		log.Printf("複合物品 捨棄 %s/%s 原因=下一入口不符", old.site.first, old.site.second)
	}
	for _, site := range compoundSites {
		first := in.Ret == compoundFirstReturn && in.S == site.first
		second := in.Ret == compoundSecondReturn && in.S == site.second
		if !first && !second {
			continue
		}
		if ov.pend != nil {
			ov.nested++
			if ov.pend.page != nil {
				ov.pend.page.fail("複合物品交錯", in.Step)
			}
			ov.pend = nil
		}
		key := site.first + "/" + site.second
		if !first || !compoundTextOK(in, site.first) || in.SP > 65533 ||
			in.Font != 0 || in.X != 18 || in.Y != site.y {
			log.Printf("複合物品 回原文 %s 原因=部件入口不符", key)
			return true
		}
		e, ok := trans[key]
		g, valid := compoundGlyphs(e, site.glyphs, scale)
		if !ok || !valid {
			log.Printf("複合物品 回原文 %s 原因=缺譯文或非法字模", key)
			return true
		}
		before0 := compoundRect(take(), site.y)
		if before0 == nil {
			log.Printf("複合物品 回原文 %s 原因=無離屏", key)
			return true
		}
		in.raw = nil
		ov.compound = &compoundPair{site: site, inputs: [2]compoundInput{in}, g: g}
		ov.compound.snaps[0] = before0
		return true
	}
	return false
}

func (ov *overlay) returnCompound(ret oracle.Addr, sp uint16, step uint64, screen []byte, palette [256][3]uint8) bool {
	p := ov.compound
	if p == nil || p.phase == 1 {
		return false
	}
	line, index, want := 0, 0, compoundFirstReturn
	if p.phase == 2 {
		line, index, want = 6, 1, compoundSecondReturn
	}
	if ret != want {
		return false
	}
	if sp != p.inputs[index].SP+2 || step < p.inputs[index].Step {
		ov.rejectCompound("返回SP或Step不符")
		return true
	}
	after := compoundRect(screen, p.site.y)
	mask := compoundMask(p.snaps[index*2], after, line)
	if mask == nil {
		ov.rejectCompound("印字差異或前景不符")
		return true
	}
	p.snaps[index*2+1], p.masks[index], p.returns[index] = after, mask, step
	if index == 0 {
		p.phase = 1
		return true
	}
	trace := &compoundTrace{Inputs: p.inputs, Returns: p.returns, State: "pending"}
	for i, b := range p.snaps {
		trace.Snapshots[i] = compoundDigest(b)
	}
	for i, b := range p.masks {
		trace.Masks[i] = compoundDigest(b)
	}
	e := &textEv{s: p.site.first + "/" + p.site.second, x0: 18, y0: p.site.y,
		bx: 18, by: p.site.y, bw: 41, bh: 11, before: p.snaps[0], after: after,
		fg: 105, fgRGB: palette[105], g: p.g, compound: trace}
	ov.evs = append(ov.evs, e)
	ov.created++
	ov.compound = nil
	log.Printf("複合物品 建立 %s step=%d 域=(18,%d,41,11)", e.s, step, e.y0)
	return true
}
