package main

// 028：可選字形家族只從同次xlate bytes派生，不進普通英文鍵查表。
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const brushFontSHA = "cfc758e69386a21bade8cb3b373fda20d7f2fac6e8904ced2ded98588b63e274"
const brushLicenseSHA = "1f11c4a22865e1289f4694dad56b3561790a4b7341a511ddcd0edddcc3b95a46"
const paperModelSHA = "1a9632b661b106367adffc2a8336808fe38202e60080620e897072e2fdcf5a22"

var approvedMenuText = [...]string{"角色選單", "新增角色", "載入存檔", "系統設定", "匯入巫術六", "結束遊戲"}
var approvedPaperText = [...]struct{ en, zh string }{
	{"presents", "呈獻"}, {"fantasy", "奇幻"}, {"role-playing", "角色扮演"}, {"simulation", "遊戲"}, {"by", "作者"},
}

type brushStripe struct {
	Scale int    `json:"scale"`
	W     int    `json:"w"`
	H     int    `json:"h"`
	A     []byte `json:"a"`
}
type brushEntry struct {
	En      string        `json:"en"`
	Zh      string        `json:"zh"`
	Stripes []brushStripe `json:"stripes"`
}
type brushMenuData struct {
	Schema        int          `json:"schema"`
	Context       string       `json:"context"`
	FontSHA256    string       `json:"font_sha256"`
	LicenseSHA256 string       `json:"license_sha256"`
	SourceSHA256  string       `json:"source_sha256"`
	Entries       []brushEntry `json:"entries"`
}
type introPaperData struct {
	Schema       int           `json:"schema"`
	Context      string        `json:"context"`
	Size         int           `json:"size"`
	SourceSHA256 string        `json:"source_sha256"`
	ModelSHA256  string        `json:"model_sha256"`
	Entries      []rosterEntry `json:"entries"`
}
type xlateData struct {
	trans            map[string]entry
	digits, ordinals map[string][]uint16
	brush            *brushMenuData
	paper            *introPaperData
}

func approvedJSON(raw []byte, v any) error {
	if len(raw) > 1<<20 {
		return fmt.Errorf("字樣家族過大")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return fmt.Errorf("字樣家族須為單一JSON")
	}
	return nil
}
func approvedSHA(s string) bool {
	v, e := hex.DecodeString(s)
	return e == nil && len(v) == 32 && len(s) == 64
}
func decodeBrush(raw []byte) (*brushMenuData, error) {
	var d brushMenuData
	if err := approvedJSON(raw, &d); err != nil {
		return nil, err
	}
	if d.Schema != 1 || d.Context != "main.menu.brush" || d.FontSHA256 != brushFontSHA || d.LicenseSHA256 != brushLicenseSHA || !approvedSHA(d.SourceSHA256) || len(d.Entries) != 6 {
		return nil, fmt.Errorf("毛筆字樣標頭不符")
	}
	for i, e := range d.Entries {
		if e.En != mainMenuLabels[i].s || e.Zh != approvedMenuText[i] || len(e.Stripes) != 2 {
			return nil, fmt.Errorf("毛筆字樣鍵不符")
		}
		for j, s := range e.Stripes {
			if s.Scale != j+2 || s.W <= 0 || s.H <= 0 || s.W > 109*s.Scale-4 || s.H > 9*s.Scale-2 || len(s.A) != s.W*s.H {
				return nil, fmt.Errorf("毛筆字條尺寸不符")
			}
			on := false
			for _, a := range s.A {
				on = on || a != 0
			}
			if !on {
				return nil, fmt.Errorf("毛筆字條空白")
			}
		}
	}
	return &d, nil
}
func ordinaryFontSize(trans map[string]entry) int {
	size := 0
	for _, e := range trans {
		for _, g := range e.G {
			h := len(g)
			if h != 15 && h != 16 || size != 0 && size != h {
				return 0
			}
			size = h
		}
	}
	return size
}
func decodePaper(raw []byte, size int) (*introPaperData, error) {
	var d introPaperData
	if err := approvedJSON(raw, &d); err != nil {
		return nil, err
	}
	if d.Schema != 1 || d.Context != "intro.mon61.paper" || size != 15 && size != 16 || d.Size != size || d.ModelSHA256 != paperModelSHA || !approvedSHA(d.SourceSHA256) || len(d.Entries) != 5 {
		return nil, fmt.Errorf("紙頁字樣標頭不符")
	}
	for i, e := range d.Entries {
		if e.En != approvedPaperText[i].en || e.Zh != approvedPaperText[i].zh || e.Dy != 0 || len(e.G) != len([]rune(e.Zh)) {
			return nil, fmt.Errorf("紙頁字樣鍵不符")
		}
		for _, g := range e.G {
			if len(g) != size {
				return nil, fmt.Errorf("紙頁字高不符")
			}
			on := false
			for _, row := range g {
				if uint32(row) >= uint32(1)<<size {
					return nil, fmt.Errorf("紙頁字模位元越界")
				}
				on = on || row != 0
			}
			if !on {
				return nil, fmt.Errorf("紙頁字模空白")
			}
		}
	}
	return &d, nil
}
func decodeXlate(raw []byte) (xlateData, error) {
	d := xlateData{trans: map[string]entry{}, digits: map[string][]uint16{}, ordinals: map[string][]uint16{}}
	var records map[string]json.RawMessage
	if err := json.Unmarshal(raw, &records); err != nil {
		return d, err
	}
	for k, v := range records {
		var err error
		switch k {
		case "__digits__":
			err = json.Unmarshal(v, &d.digits)
		case "__ordinals__":
			err = json.Unmarshal(v, &d.ordinals)
		case "__main_menu_brush__", "__intro_paper__":
			continue
		default:
			var e entry
			err = json.Unmarshal(v, &e)
			d.trans[k] = e
		}
		if err != nil {
			return d, err
		}
	}
	if len(d.trans) == 0 {
		return d, nil
	}
	if v, ok := records["__main_menu_brush__"]; ok {
		var err error
		d.brush, err = decodeBrush(v)
		if err != nil {
			log.Print("主選單毛筆缺譯：資源不可用，沿一般字樣")
		}
	}
	if v, ok := records["__intro_paper__"]; ok {
		var err error
		d.paper, err = decodePaper(v, ordinaryFontSize(d.trans))
		if err != nil {
			log.Print("紙頁缺譯：字樣資源不可用，保留原圖")
		}
	}
	return d, nil
}

type approvedPatch struct {
	at  image.Point
	img *image.RGBA
}
type brushMenuView struct{ patches [2][]approvedPatch }

func (d *brushMenuData) view(shown []*textEv, comp []byte, pal [256][3]uint8) *brushMenuView {
	if d == nil || len(comp) != 64000 {
		return nil
	}
	for _, label := range mainMenuLabels {
		found := false
		for _, e := range shown {
			if !e.dead && e.shown && e.match && mainMenuLabelMatch(e, label.s, label.x, label.y) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	v := &brushMenuView{}
	for i, e := range d.Entries {
		for _, s := range e.Stripes {
			at := image.Pt(105*s.Scale+(109*s.Scale-s.W)/2, (mainMenuLabels[i].y-2)*s.Scale+(9*s.Scale-s.H)/2)
			img := image.NewRGBA(image.Rect(0, 0, s.W, s.H))
			for y := 0; y < s.H; y++ {
				for x := 0; x < s.W; x++ {
					bg := pal[comp[((at.Y+y)/s.Scale)*320+(at.X+x)/s.Scale]]
					a := int(s.A[y*s.W+x])
					fg := [3]uint8{248, 232, 160}
					c := color.RGBA{A: 255}
					channels := []*uint8{&c.R, &c.G, &c.B}
					for j, ch := range channels {
						*ch = uint8((int(fg[j])*a + int(bg[j])*(255-a) + 127) / 255)
					}
					img.SetRGBA(x, y, c)
				}
			}
			v.patches[s.Scale-2] = append(v.patches[s.Scale-2], approvedPatch{at, img})
		}
	}
	return v
}
func (v *brushMenuView) suppress(e *textEv, scale int) bool {
	if v == nil || scale != 2 && scale != 3 {
		return false
	}
	for _, l := range mainMenuLabels {
		if mainMenuLabelMatch(e, l.s, l.x, l.y) {
			return true
		}
	}
	return false
}
func (v *brushMenuView) paintHD(img *image.RGBA, scale int) {
	if v == nil || scale != 2 && scale != 3 || img.Bounds() != image.Rect(0, 0, 320*scale, 200*scale) {
		return
	}
	for _, p := range v.patches[scale-2] {
		draw.Draw(img, p.img.Bounds().Add(p.at), p.img, image.Point{}, draw.Src)
	}
}
func (g *game) drawBrushMenu(dst *ebiten.Image, v *brushMenuView) {
	if v == nil || g.scale != 2 && g.scale != 3 || dst.Bounds() != image.Rect(0, 0, 320*g.scale, 200*g.scale) {
		return
	}
	if g.brushView != v {
		g.brushView = v
		g.brushImages = [2][]*ebiten.Image{}
		for k, patches := range v.patches {
			for _, p := range patches {
				g.brushImages[k] = append(g.brushImages[k], ebiten.NewImageFromImage(p.img))
			}
		}
	}
	for i, p := range v.patches[g.scale-2] {
		var op ebiten.DrawImageOptions
		op.GeoM.Translate(float64(p.at.X), float64(p.at.Y))
		dst.DrawImage(g.brushImages[g.scale-2][i], &op)
	}
}
