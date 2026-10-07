package main

import (
	"crypto/sha256"
	"encoding/json"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
)

func syntheticApproved(t *testing.T) (brushMenuData, introPaperData, map[string]any) {
	t.Helper()
	b := brushMenuData{Schema: 1, Context: "main.menu.brush", FontSHA256: brushFontSHA, LicenseSHA256: brushLicenseSHA, SourceSHA256: strings.Repeat("0", 64)}
	for i, label := range mainMenuLabels {
		e := brushEntry{En: label.s, Zh: approvedMenuText[i]}
		for _, scale := range []int{2, 3} {
			e.Stripes = append(e.Stripes, brushStripe{scale, 2, 2, []byte{0, 128, 255, 1}})
		}
		b.Entries = append(b.Entries, e)
	}
	p := introPaperData{Schema: 1, Context: "intro.mon61.paper", Size: 15, SourceSHA256: strings.Repeat("1", 64), ModelSHA256: paperModelSHA}
	for _, label := range approvedPaperText {
		e := rosterEntry{En: label.en, Zh: label.zh}
		for range []rune(label.zh) {
			g := make([]uint16, 15)
			g[0] = 1
			e.G = append(e.G, g)
		}
		p.Entries = append(p.Entries, e)
	}
	r := map[string]any{"synthetic": entry{G: [][]uint16{make([]uint16, 15)}}, "__main_menu_brush__": b, "__intro_paper__": p}
	return b, p, r
}
func TestApprovedXlateFamilies(t *testing.T) {
	_, _, r := syntheticApproved(t)
	raw, _ := json.Marshal(r)
	d, err := decodeXlate(raw)
	if err != nil || d.brush == nil || d.paper == nil || len(d.trans) != 1 {
		t.Fatalf("完整家族：%v", err)
	}
	delete(r, "__main_menu_brush__")
	delete(r, "__intro_paper__")
	raw, _ = json.Marshal(r)
	old, err := decodeXlate(raw)
	if err != nil || old.brush != nil || old.paper != nil || len(old.trans) != 1 {
		t.Fatal("舊xlate回退")
	}
	b, p, r := syntheticApproved(t)
	b.Entries[0].Stripes[0].A = nil
	r["__main_menu_brush__"] = b
	raw, _ = json.Marshal(r)
	d, err = decodeXlate(raw)
	if err != nil || d.brush != nil || d.paper == nil {
		t.Fatal("非法毛筆家族影響紙頁或主資源")
	}
	r = map[string]any{"__intro_paper__": p, "__main_menu_brush__": b}
	raw, _ = json.Marshal(r)
	d, err = decodeXlate(raw)
	if err != nil || d.brush != nil || d.paper != nil {
		t.Fatal("英文模式啟用新家族")
	}
	if _, err := decodeXlate([]byte(`{"synthetic":[]}`)); err == nil {
		t.Fatal("普通字模錯誤未保持失敗")
	}
}
func TestApprovedParserGuards(t *testing.T) {
	b, p, _ := syntheticApproved(t)
	for _, mutate := range []func(*brushMenuData){
		func(v *brushMenuData) { v.Schema++ },
		func(v *brushMenuData) { v.Context = "elsewhere" },
		func(v *brushMenuData) { v.FontSHA256 = strings.Repeat("0", 64) },
		func(v *brushMenuData) { v.Entries[0].En = "elsewhere" },
		func(v *brushMenuData) { v.Entries[0].Stripes[0].W = 1000 },
		func(v *brushMenuData) { v.Entries[0].Stripes[0].Scale = 3 },
		func(v *brushMenuData) { v.Entries[0].Stripes[0].A = make([]byte, 4) },
	} {
		copyRaw, _ := json.Marshal(b)
		var invalid brushMenuData
		json.Unmarshal(copyRaw, &invalid)
		mutate(&invalid)
		raw, _ := json.Marshal(invalid)
		if _, err := decodeBrush(raw); err == nil {
			t.Fatal("非法毛筆資源通過")
		}
	}
	raw, _ := json.Marshal(b)
	for _, invalid := range [][]byte{append(append([]byte(nil), raw...), []byte(` {}`)...), []byte(`{"unknown":1}`), []byte(strings.Repeat(" ", (1<<20)+1))} {
		if _, err := decodeBrush(invalid); err == nil {
			t.Fatal("JSON或容量守門失效")
		}
	}
	p.Entries[0].G[0][0] = 0x8000
	raw, _ = json.Marshal(p)
	if _, err := decodePaper(raw, 15); err == nil {
		t.Fatal("紙頁位元越界")
	}
	if _, err := decodePaper(raw, 16); err == nil {
		t.Fatal("纸頁後端不符")
	}
}
func TestBrushPublicationAndIntegerBlend(t *testing.T) {
	b, _, _ := syntheticApproved(t)
	var shown []*textEv
	for _, l := range mainMenuLabels {
		shown = append(shown, &textEv{s: l.s, x0: l.x, y0: l.y, shown: true, match: true})
	}
	comp := make([]byte, 64000)
	var pal [256][3]uint8
	pal[0] = [3]uint8{30, 20, 10}
	v := b.view(shown, comp, pal)
	if v == nil || !v.suppress(shown[0], 3) || v.suppress(shown[0], 4) {
		t.Fatal("完整訊號／倍率守門")
	}
	patch := v.patches[1][0]
	if patch.img.RGBAAt(1, 0) != (color.RGBA{139, 126, 85, 255}) || patch.img.RGBAAt(0, 1) != (color.RGBA{248, 232, 160, 255}) {
		t.Fatal("整數alpha合成")
	}
	pal[0] = [3]uint8{1, 2, 3}
	b.Entries[0].Stripes[1].A[1] = 0
	if patch.img.RGBAAt(1, 0) != (color.RGBA{139, 126, 85, 255}) {
		t.Fatal("發布view仍依賴可變底色或字模")
	}
	shown[0].dead = true
	if b.view(shown, comp, pal) != nil {
		t.Fatal("死事件仍啟用毛筆")
	}
	shown[0].dead = false
	shown[0].x0++
	if b.view(shown, comp, pal) != nil || b.view(shown[1:], comp, pal) != nil {
		t.Fatal("其他座標／缺事件仍啟用毛筆")
	}
	wrong := image.NewRGBA(image.Rect(0, 0, 960, 640))
	v.paintHD(wrong, 3)
	for _, c := range wrong.Pix {
		if c != 0 {
			t.Fatal("頁尾整窗未排除")
		}
	}
}

func TestApprovedRecordingCapturedBytes(t *testing.T) {
	_, _, payload := syntheticApproved(t)
	raw, _ := json.Marshal(payload)
	path := filepath.Join(t.TempDir(), "xlate.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	r := &recordingResources{metadata: []ic.Resource{{Role: "xlate"}}}
	r.afterRead = func(role, path string) {
		if err := os.WriteFile(path, []byte("changed-after-capture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	captured, err := r.read("xlate", path)
	if err != nil {
		t.Fatal(err)
	}
	d, err := decodeXlate(captured)
	if err != nil || d.brush == nil || d.paper == nil || r.metadata[0].SHA256 != ic.Hash(raw) || r.metadata[0].Size != uint64(len(raw)) {
		t.Fatal("同次捕獲後重開或遺失新家族SHA")
	}
	if _, err := decodeXlate([]byte("changed-after-capture")); err == nil {
		t.Fatal("替換bytes應拒絕")
	}
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	b, err := loadRecordingBundle(path, "", "")
	if err != nil || len(b.resources.metadata) != 6 || !b.resources.metadata[0].Loaded || b.resources.metadata[0].SHA256 != ic.Hash(raw) || b.brush == nil || b.picture.paper == nil {
		t.Fatalf("正式bundle未沿六roles同次解碼：%v", err)
	}
}

func TestApprovedRecordingGameArea(t *testing.T) {
	b, _, _ := syntheticApproved(t)
	var shown []*textEv
	for _, l := range mainMenuLabels {
		shown = append(shown, &textEv{s: l.s, x0: l.x, y0: l.y, shown: true, match: true})
	}
	brush := b.view(shown, make([]byte, 64000), [256][3]uint8{})
	glyph := image.NewRGBA(image.Rect(0, 0, 2, 2))
	paper := &introPaperView{glyphs: []approvedPatch{{image.Pt(10, 10), glyph}}}
	legacy := &introPictureView{glyph: glyph, quad: introQuad{{1, 1}, {2, 1}, {2, 2}, {1, 2}}}
	g := &game{scale: 3, screen: ebiten.NewImage(320, 200), live: &inputRecording{font: &helpData{}},
		sh: &shared{rgba: make([]byte, 64000*4), brush: brush, paper: paper, picture: legacy}}
	g.Draw(ebiten.NewImage(960, 640))
	if g.paperView != paper || g.brushView != brush || g.pictureGlyph != glyph {
		t.Fatal("40px錄製頁尾令新舊片頭／毛筆守門失敗")
	}
}

func TestIntroPaperCapturedFrameGuards(t *testing.T) {
	root := os.Getenv("WIZ7_APPROVED_UI_CASES")
	if root == "" {
		t.Skip("需要本機正常冷啟動畫格與烘製資源")
	}
	read := func(path string) []byte {
		v, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	indexed := read("before/hdat-73400000.a000")
	palBytes := read("before/hdat-73400000.pal")
	if len(palBytes) != 768 || len(indexed) != 64000 {
		t.Fatal("收據長度")
	}
	var pal [256][3]uint8
	for i := range pal {
		copy(pal[i][:], palBytes[i*3:i*3+3])
	}
	d, err := decodeXlate(read("unifont/xlate.json"))
	if err != nil || d.paper == nil {
		t.Fatal("烘製紙頁資源")
	}
	s := &introPictureState{paper: d.paper}
	ret := oracle.Addr{Seg: 0x110, Off: 0x5b9b}
	s.enter(ret, 0x100, [4]uint16{0, 0, 319, 199})
	s.returned(ret, 0x102, indexed)
	if s.paperView(indexed, pal, 0x13, 3) == nil {
		t.Fatal("已量完整來源未辨識")
	}
	if s.paperView(indexed, pal, 0x13, 2) != nil || s.paperView(indexed, pal, 0, 3) != nil || s.paperView(indexed[:100], pal, 0x13, 3) != nil {
		t.Fatal("倍率／mode／尺寸守門")
	}
	wrongPal := pal
	wrongPal[0][0]++
	if s.paperView(indexed, wrongPal, 0x13, 3) != nil {
		t.Fatal("已快取view跳過色盤守門")
	}
	changed := append([]byte(nil), indexed...)
	changed[100]++
	if s.paperView(changed, pal, 0x13, 3) != nil {
		t.Fatal("部分重寫仍沿用紙頁")
	}
	s.enter(ret, 0x100, [4]uint16{0, 0, 319, 199})
	s.returned(ret, 0x104, indexed)
	if s.paperView(indexed, pal, 0x13, 3) != nil {
		t.Fatal("錯SP未撤除")
	}
	s.complete, s.hasComplete = sha256.Sum256(indexed), true
	s.paper = nil
	if s.paperView(indexed, pal, 0x13, 3) != nil {
		t.Fatal("缺資源沿用舊view")
	}
	s.paper = d.paper
	blank := read("before/hdat-105400000.a000")
	s.complete = sha256.Sum256(blank)
	if s.paperView(blank, pal, 0x13, 3) != nil {
		t.Fatal("第六無字頁未撤除")
	}
}
