package main

// 024：只覆繪0110:3877直接印出的角色性別，不改一般字串鍵。
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/wicanr2/dosgolem/oracle"
)

const profileSexLimit = 1 << 20

var profileSexReturn = oracle.Addr{Seg: 0x0110, Off: 0x6c42}

type profileSexData struct {
	Schema       int           `json:"schema"`
	Size         int           `json:"size"`
	Context      string        `json:"context"`
	Source       string        `json:"source"`
	SourceSHA256 string        `json:"source_sha256"`
	Entries      []rosterEntry `json:"entries"`
}

func decodeProfileSex(r io.Reader, size int) (map[string]rosterEntry, error) {
	raw, err := io.ReadAll(io.LimitReader(r, profileSexLimit+1))
	if err != nil || len(raw) > profileSexLimit {
		return nil, fmt.Errorf("性別資源讀取失敗或過大")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	var data profileSexData
	if err := d.Decode(&data); err != nil {
		return nil, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("性別資源須為單一JSON文件")
	}
	hash, err := hex.DecodeString(data.SourceSHA256)
	if data.Schema != 1 || size != 15 && size != 16 || data.Size != size ||
		data.Context != "character.review.sex" || data.Source != "runtime:0110:6c42" ||
		err != nil || len(hash) != 32 || len(data.Entries) != 2 {
		return nil, fmt.Errorf("性別資源標頭不合格")
	}
	entries := map[string]rosterEntry{}
	for _, e := range data.Entries {
		if e.En != "M" && e.En != "F" || entries[e.En].En != "" || e.Dy != 0 ||
			!utf8.ValidString(e.Zh) || utf8.RuneCountInString(e.Zh) != 1 || len(e.G) != 1 || len(e.G[0]) != size {
			return nil, fmt.Errorf("性別資源鍵或文案不合格")
		}
		for _, c := range e.Zh {
			if !unicode.Is(unicode.Han, c) {
				return nil, fmt.Errorf("性別譯文須為單一漢字")
			}
		}
		nonempty := false
		for _, row := range e.G[0] {
			if uint32(row) >= uint32(1)<<size {
				return nil, fmt.Errorf("性別字模位元越界")
			}
			nonempty = nonempty || row != 0
		}
		if !nonempty {
			return nil, fmt.Errorf("性別字模空白")
		}
		entries[e.En] = e
	}
	return entries, nil
}

func loadProfileSex(xlate string, trans map[string]entry) map[string]rosterEntry {
	if xlate == "" || len(trans) == 0 {
		return nil
	}
	size := 0
	for _, e := range trans {
		for _, glyph := range e.G {
			h := len(glyph)
			if h != 15 && h != 16 || size != 0 && h != size {
				log.Print("性別資源停用：一般字模高度不一致")
				return nil
			}
			size = h
		}
	}
	f, err := os.Open(filepath.Join(filepath.Dir(xlate), "profile-sex.json"))
	if err == nil {
		defer f.Close()
		var entries map[string]rosterEntry
		entries, err = decodeProfileSex(f, size)
		if err == nil {
			log.Print("性別資源：載入兩鍵")
			return entries
		}
	}
	log.Printf("性別資源停用：%v", err)
	return nil
}

type pendingProfileSex struct {
	s              string
	sp             uint16
	ret            oracle.Addr
	x, y           int
	before         []byte
	g              [][]uint16
	ss             uint16
	entryStep      uint64
	driver, buffer uint16
}

func (ov *overlay) clearProfileSex(reason string) {
	if ov.sexPending != nil {
		if ov.sexPending.ret == creationSexReturn {
			ov.clearCreationSex()
		} else {
			ov.sexPending = nil
		}
		log.Printf("性別缺譯：%s", reason)
	}
}

func (ov *overlay) enterProfileSex(ret oracle.Addr, sp, font, code uint16, x, y int, before []byte) {
	if (ret != profileSexReturn && ret != worldSexReturn) || x != 90 || y != 6 {
		return
	}
	if ret == worldSexReturn && (sp > 65530 || (flowK != 2 && flowK != 3) || !worldSexResourcesValid(ov.sexEntries)) {
		ov.clearWorldSex()
		log.Print("世界性別缺譯：stack、倍率或資源不符")
		return
	}
	for _, e := range ov.evs {
		if e.sex {
			e.dead = true
		}
	}
	if ov.sexPending != nil {
		ov.clearProfileSex("巢狀印字")
		return
	}
	if (!ov.sexEnabled && !(ret == worldSexReturn && (flowK == 2 || flowK == 3))) || font != 0 || code != 'M' && code != 'F' ||
		ov.pend != nil || ov.compound != nil || len(before) != 64000 {
		log.Print("性別缺譯：倍率、資源、印字訊號或離屏不符")
		return
	}
	if ret == worldSexReturn && !worldSexBackgroundMatches(before) {
		log.Print("世界性別缺譯：背景不符")
		return
	}
	s := string(rune(code))
	e, ok := ov.sexEntries[s]
	if !ok {
		log.Print("性別缺譯：缺少字模")
		return
	}
	ov.sexPending = &pendingProfileSex{s: s, sp: sp, ret: ret, x: x, y: y,
		before: append([]byte(nil), before...), g: e.G}
}

func (ov *overlay) returnProfileSex(sp uint16, x, y int, after []byte, pal [256][3]uint8) {
	p := ov.sexPending
	ov.sexPending = nil // 本次返回只消費一次，失敗也不留slot。
	if p == nil {
		return
	}
	if (p.ret != profileSexReturn && p.ret != worldSexReturn && p.ret != creationSexReturn) || sp != p.sp+2 || x != 96 || y != 6 ||
		len(p.before) != 64000 || len(after) != 64000 {
		log.Print("性別缺譯：返回或離屏不符")
		return
	}
	changes := 0
	for i, a := range after {
		if a == p.before[i] {
			continue
		}
		x, y := i%320, i/320
		if x < p.x || x >= p.x+6 || y < p.y || y >= p.y+6 || a != 1 {
			log.Print("性別缺譯：印字差異不符")
			return
		}
		changes++
	}
	if changes == 0 {
		log.Print("性別缺譯：無印字差異")
		return
	}
	bx, by, bw, bh := p.x, p.y, 6, 6
	if p.ret == worldSexReturn || p.ret == creationSexReturn {
		bx, by, bw, bh = 88, 4, 8, 8
	}
	e := &textEv{originRet: p.ret, s: p.s, x0: p.x, y0: p.y, bx: bx, by: by, bw: bw, bh: bh,
		fg: 1, fgRGB: pal[1], g: p.g, sex: true}
	for y := by; y < by+bh; y++ {
		start := y*320 + bx
		e.before = append(e.before, p.before[start:start+bw]...)
		e.after = append(e.after, after[start:start+bw]...)
	}
	ov.evs = append(ov.evs, e)
	ov.created++
}
