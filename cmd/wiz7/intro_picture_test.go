package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
)

// 自製字模只測資源邊界，原生字模與遊戲像素另由合法本機輸入驗證。
func introTestData(size int) introPictureData {
	g := make([]uint16, size)
	g[0] = 1
	return introPictureData{1, size, "intro.mon60.presents", "runtime:0110:5a9d", strings.Repeat("a", 64), []rosterEntry{{En: "presents", Zh: "呈獻", G: [][]uint16{g, g}}}}
}

func TestIntroPictureLocalEvidence(t *testing.T) {
	root, resource := os.Getenv("WIZ7_INTRO_EVIDENCE"), os.Getenv("WIZ7_INTRO_RESOURCE")
	if root == "" || resource == "" {
		t.Skip("需要玩家本機的原版圖片收據與已烘製字模")
	}
	read := func(name string) []byte {
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	b, err := os.ReadFile(resource)
	if err != nil {
		t.Fatal(err)
	}
	data, err := decodeIntroPicture(bytes.NewReader(b), 16)
	if err != nil {
		t.Fatal(err)
	}
	s := &introPictureState{data: data}
	ret := oracle.Addr{Seg: 0x110, Off: 0x5afd}
	names := []string{"38938757", "40696847", "42466770", "44378262", "46331493", "48267654"}
	for _, name := range names {
		indexed, rawPal := read("picture-"+name+".a000"), read("picture-"+name+".pal")
		if len(rawPal) != 768 {
			t.Fatal("palette length")
		}
		var pal [256][3]uint8
		for i := range pal {
			copy(pal[i][:], rawPal[i*3:i*3+3])
		}
		s.enter(ret, 100, [4]uint16{0, 0, 319, 199})
		s.returned(ret, 102, indexed)
		if s.view(indexed, pal, 0x13) == nil {
			t.Fatal("known complete frame rejected", name)
		}
		// 下一檔載入不會產生完成通知，末圖必須仍有效。
		if s.view(indexed, pal, 0x13) == nil {
			t.Fatal("last view lost", name)
		}
		badPal := pal
		badPal[0][0] ^= 1
		if s.view(indexed, badPal, 0x13) != nil || s.view(indexed, pal, 3) != nil || s.view(indexed[:63999], pal, 0x13) != nil {
			t.Fatal("invalid source accepted", name)
		}
		partial := append([]byte(nil), indexed...)
		partial[0] ^= 1
		if s.view(partial, pal, 0x13) != nil {
			t.Fatal("stale view accepted", name)
		}
		if name == names[len(names)-1] {
			next := read("picture-73235759.a000")
			s.enter(ret, 100, [4]uint16{0, 0, 319, 199})
			s.returned(ret, 102, next)
			if s.view(next, pal, 0x13) != nil || s.view(indexed, pal, 0x13) != nil {
				t.Fatal("complete next-page did not remove label")
			}
		}
	}
}
func TestIntroPictureResourceBoundaries(t *testing.T) {
	for _, size := range []int{15, 16} {
		v := introTestData(size)
		b, _ := json.Marshal(v)
		if _, err := decodeIntroPicture(bytes.NewReader(b), size); err != nil {
			t.Fatal(err)
		}
		for name, mutate := range map[string]func(*introPictureData){
			"backend": func(v *introPictureData) { v.Size = 17 }, "source": func(v *introPictureData) { v.Source = "other" },
			"key": func(v *introPictureData) { v.Entries[0].En = "other" }, "partial": func(v *introPictureData) { v.Entries[0].G = v.Entries[0].G[:1] },
			"blank":    func(v *introPictureData) { v.Entries[0].G[0] = make([]uint16, size) },
			"overflow": func(v *introPictureData) { v.Entries[0].G[0][0] = 1 << size },
			"height":   func(v *introPictureData) { v.Entries[0].G[0] = v.Entries[0].G[0][:size-1] },
		} {
			t.Run(name, func(t *testing.T) {
				v := introTestData(size)
				mutate(&v)
				b, _ := json.Marshal(v)
				if _, err := decodeIntroPicture(bytes.NewReader(b), size); err == nil {
					t.Fatal("accepted invalid resource")
				}
			})
		}
		for _, bad := range [][]byte{append(append([]byte{}, b...), b...), []byte(`{"extra":1}`), bytes.Repeat([]byte{' '}, (1<<20)+1)} {
			if _, err := decodeIntroPicture(bytes.NewReader(bad), size); err == nil {
				t.Fatal("accepted malformed JSON")
			}
		}
	}
}
func TestIntroPictureWitness(t *testing.T) {
	ret := oracle.Addr{Seg: 0x110, Off: 0x5afd}
	rect := [4]uint16{0, 0, 319, 199}
	frame := make([]byte, 64000)
	s := &introPictureState{}
	if !s.enter(ret, 100, rect) {
		t.Fatal("shared return rejected")
	}
	s.returned(ret, 102, frame)
	if !s.hasComplete || s.complete != sha256.Sum256(frame) {
		t.Fatal("complete frame not witnessed")
	}
	// 單純開檔沒有入口或返回通知，不應清掉末圖見證。
	if s.complete != sha256.Sum256(frame) {
		t.Fatal("witness lost")
	}
	if !s.enter(ret, 100, rect) {
		t.Fatal("enter")
	}
	s.returned(ret, 104, frame)
	if s.hasComplete {
		t.Fatal("wrong SP accepted")
	}
	if s.enter(ret, 100, [4]uint16{0, 0, 318, 199}) {
		t.Fatal("wrong rectangle accepted")
	}
	s.enter(ret, 100, rect)
	if s.enter(ret, 90, rect) || s.pending != nil {
		t.Fatal("reentry accepted")
	}
	s.enter(ret, 100, rect)
	frame[0] = 1
	s.returned(ret, 102, frame)
	if s.view(frame, [256][3]uint8{}, 0x13) != nil {
		t.Fatal("unknown complete frame translated")
	}
	frame[1] = 1
	if s.view(frame, [256][3]uint8{}, 0x13) != nil {
		t.Fatal("partial frame translated")
	}
}
func TestIntroPictureGeometry(t *testing.T) {
	q := introFrames[5].quad
	if !q.contains(255, 111) || q.contains(193, 111) {
		t.Fatal("closed edge geometry")
	}
	d := introTestData(16)
	v := &introPictureView{quad: introFrames[3].quad, glyph: introGlyph(&d, [3]uint8{})}
	if v.anchor(3) != image.Pt(604, 320) {
		t.Fatalf("round-to-even anchor %v", v.anchor(3))
	}
	rgba := make([]byte, 64000*4)
	v.pixels = []introPixel{{10, 10, [3]uint8{1, 2, 3}}}
	v.patchRGBA(rgba, 1)
	if !bytes.Equal(rgba, make([]byte, len(rgba))) {
		t.Fatal("unsupported scale patched")
	}
	if _, err := introSolve([3][3]float64{}, [3]float64{}); err == nil {
		t.Fatal("singular model accepted")
	}
}
