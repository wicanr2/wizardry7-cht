package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

func sexFixture(size int) profileSexData {
	g := make([]uint16, size)
	g[0] = 1
	return profileSexData{Schema: 1, Size: size, Context: "character.review.sex",
		Source: "runtime:0110:6c42", SourceSHA256: strings.Repeat("a", 64),
		Entries: []rosterEntry{{En: "M", Zh: "男", G: [][]uint16{g}}, {En: "F", Zh: "女", G: [][]uint16{g}}}}
}

func TestProfileSexResourceRejectsPartialOrAmbiguousData(t *testing.T) {
	for _, size := range []int{15, 16} {
		data := sexFixture(size)
		b, _ := json.Marshal(data)
		entries, err := decodeProfileSex(bytes.NewReader(b), size)
		if err != nil || len(entries) != 2 || len(entries["M"].G) != 1 {
			t.Fatalf("size %d: %v", size, err)
		}
	}
	cases := map[string]func(*profileSexData){
		"missing":   func(d *profileSexData) { d.Entries = d.Entries[:1] },
		"duplicate": func(d *profileSexData) { d.Entries[1].En = "M" },
		"unknown":   func(d *profileSexData) { d.Entries[1].En = "X" },
		"context":   func(d *profileSexData) { d.Context = "other" },
		"source":    func(d *profileSexData) { d.Source = "runtime:0110:3895" },
		"hash":      func(d *profileSexData) { d.SourceSHA256 = strings.Repeat("z", 64) },
		"backend":   func(d *profileSexData) { d.Size = 16 },
		"text":      func(d *profileSexData) { d.Entries[0].Zh = "男戰士" },
		"ascii":     func(d *profileSexData) { d.Entries[0].Zh = "M" },
		"dy":        func(d *profileSexData) { d.Entries[0].Dy = 1 },
		"empty":     func(d *profileSexData) { d.Entries[0].G = [][]uint16{make([]uint16, 15)} },
		"height":    func(d *profileSexData) { d.Entries[0].G = [][]uint16{{1}} },
		"bits":      func(d *profileSexData) { d.Entries[0].G[0][0] = 1 << 15 },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			d := sexFixture(15)
			change(&d)
			b, _ := json.Marshal(d)
			if entries, err := decodeProfileSex(bytes.NewReader(b), 15); err == nil || entries != nil {
				t.Fatal("不合格資源不能留下部分字模")
			}
		})
	}
	b, _ := json.Marshal(sexFixture(15))
	for _, raw := range [][]byte{append(append([]byte(nil), b...), b...),
		bytes.Replace(b, []byte(`"schema":1`), []byte(`"extra":1,"schema":1`), 1),
		append(append([]byte(nil), b...), bytes.Repeat([]byte(" "), profileSexLimit)...)} {
		if _, err := decodeProfileSex(bytes.NewReader(raw), 15); err == nil {
			t.Fatal("多文件、未知欄及超量輸入須拒絕")
		}
	}
}

func sexOverlay() *overlay {
	d := sexFixture(15)
	return &overlay{sexEnabled: true, sexEntries: map[string]rosterEntry{"M": d.Entries[0], "F": d.Entries[1]}}
}

func texturedSexFrame() []byte {
	b := bytes.Repeat([]byte{115}, 64000)
	for y := 6; y < 12; y++ {
		for x := 90; x < 96; x++ {
			b[y*320+x] = byte(14 + (x+y)%5)
		}
	}
	return b
}

func TestProfileSexOwnsSnapshotAndPreservesTexture(t *testing.T) {
	for _, code := range []uint16{'M', 'F'} {
		ov := sexOverlay()
		before := texturedSexFrame()
		ov.enterProfileSex(profileSexReturn, 100, 0, code, 90, 6, before)
		before[0] = 0
		if ov.sexPending == nil || ov.sexPending.before[0] != 115 {
			t.Fatal("before必須是獨立深複製")
		}
		after := append([]byte(nil), ov.sexPending.before...)
		for _, i := range []int{6*320 + 90, 7*320 + 93, 11*320 + 95} {
			after[i] = 1
		}
		var pal [256][3]uint8
		pal[1] = [3]uint8{9, 20, 31}
		ov.returnProfileSex(102, 96, 6, after, pal)
		if ov.sexPending != nil || len(ov.evs) != 1 || ov.created != 1 {
			t.Fatal("合法返回須單次建立事件並清slot")
		}
		e := ov.evs[0]
		if !e.sex || e.fg != 1 || e.fgRGB != pal[1] || e.bw != 6 || e.bh != 6 || len(e.before) != 36 {
			t.Fatal("性別事件必須保存精確字格及前景")
		}
		for i, v := range e.before {
			if v != byte(14+(90+i%6+6+i/6)%5) {
				t.Fatal("不能把底圖紋理改成固定底色")
			}
		}
		ov.returnProfileSex(102, 96, 6, after, pal)
		if len(ov.evs) != 1 {
			t.Fatal("不可重複消費返回")
		}
	}
}

func TestProfileSexRejectsReturnAndClearsPending(t *testing.T) {
	for _, name := range []string{"sp", "cursor", "missing", "empty", "outside", "color"} {
		t.Run(name, func(t *testing.T) {
			ov := sexOverlay()
			before := texturedSexFrame()
			ov.enterProfileSex(profileSexReturn, 100, 0, 'M', 90, 6, before)
			after := append([]byte(nil), before...)
			after[6*320+90] = 1
			sp, x := uint16(102), 96
			switch name {
			case "sp":
				sp = 104
			case "cursor":
				x = 102
			case "missing":
				after = nil
			case "empty":
				after = before
			case "outside":
				after[0] = 1
			case "color":
				after[6*320+90] = 2
			}
			ov.returnProfileSex(sp, x, 6, after, [256][3]uint8{})
			if ov.sexPending != nil || len(ov.evs) != 0 {
				t.Fatal("拒絕返回不得留slot或中文事件")
			}
		})
	}
}

func TestProfileSexArbitrationAndSignalGuards(t *testing.T) {
	for _, name := range []string{"ret", "font", "code", "x", "y", "scale", "frame", "generic", "compound", "nested"} {
		t.Run(name, func(t *testing.T) {
			ov := sexOverlay()
			ret, font, code, x, y := profileSexReturn, uint16(0), uint16('M'), 90, 6
			before := texturedSexFrame()
			general := &pendingPrint{}
			switch name {
			case "ret":
				ret = oracle.Addr{Seg: 0x0110, Off: 0x389e}
			case "font":
				font = 1
			case "code":
				code = 'X'
			case "x":
				x = 96
			case "y":
				y = 12
			case "scale":
				ov.sexEnabled = false
			case "frame":
				before = nil
			case "generic":
				ov.pend = general
			case "compound":
				ov.compound = &compoundPair{}
			case "nested":
				ov.sexPending = &pendingProfileSex{sp: 88}
			}
			ov.enterProfileSex(ret, 100, font, code, x, y, before)
			if ov.sexPending != nil || len(ov.evs) != 0 || name == "generic" && ov.pend != general {
				t.Fatal("不合格訊號不能佔用slot或改一般pending")
			}
		})
	}
	ov := sexOverlay()
	ov.sexPending = &pendingProfileSex{}
	general := &pendingPrint{}
	ov.pend = general
	ov.clearProfileSex("一般字串插入")
	if ov.sexPending != nil || ov.pend != general {
		t.Fatal("一般入口只可撤性別slot")
	}
}
