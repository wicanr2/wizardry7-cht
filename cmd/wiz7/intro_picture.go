package main

// 025：只在已見證的六張完整片頭圖上補底及畫字，不改原版資料。
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

type introPictureData struct {
	Schema       int           `json:"schema"`
	Size         int           `json:"size"`
	Context      string        `json:"context"`
	Source       string        `json:"source"`
	SourceSHA256 string        `json:"source_sha256"`
	Entries      []rosterEntry `json:"entries"`
}

func decodeIntroPicture(r io.Reader, size int) (*introPictureData, error) {
	b, err := io.ReadAll(io.LimitReader(r, (1<<20)+1))
	if err != nil || len(b) > 1<<20 {
		return nil, fmt.Errorf("片頭資源讀取失敗或過大")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	var v introPictureData
	if err = d.Decode(&v); err != nil {
		return nil, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return nil, fmt.Errorf("片頭資源須為單一JSON")
	}
	h, err := hex.DecodeString(v.SourceSHA256)
	if v.Schema != 1 || (size != 15 && size != 16) || v.Size != size || v.Context != "intro.mon60.presents" || v.Source != "runtime:0110:5a9d" || err != nil || len(h) != 32 || len(v.Entries) != 1 {
		return nil, fmt.Errorf("片頭資源標頭不合格")
	}
	e := v.Entries[0]
	if e.En != "presents" || e.Zh != "呈獻" || e.Dy != 0 || len(e.G) != 2 {
		return nil, fmt.Errorf("片頭資源鍵或文案不合格")
	}
	for _, g := range e.G {
		if len(g) != size {
			return nil, fmt.Errorf("片頭字高不合格")
		}
		any := false
		for _, row := range g {
			if uint32(row) >= uint32(1)<<size {
				return nil, fmt.Errorf("片頭字模位元越界")
			}
			any = any || row != 0
		}
		if !any {
			return nil, fmt.Errorf("片頭字模空白")
		}
	}
	return &v, nil
}

func loadIntroPicture(path string, trans map[string]entry) *introPictureState {
	s := &introPictureState{}
	if path == "" || len(trans) == 0 {
		s.english = true
		return s
	}
	size := 0
	for _, e := range trans {
		for _, g := range e.G {
			h := len(g)
			if h != 15 && h != 16 || size != 0 && size != h {
				s.report("一般字模高度不一致")
				return s
			}
			size = h
		}
	}
	f, err := os.Open(filepath.Join(filepath.Dir(path), "intro-picture.json"))
	if err == nil {
		defer f.Close()
		s.data, err = decodeIntroPicture(f, size)
	}
	if err != nil {
		s.report("資源不可用")
		return s
	}
	log.Print("片頭圖片資源：載入呈獻")
	return s
}

type introQuad [4]image.Point
type introFrame struct {
	sha  string
	quad introQuad
}

var introFrames = []introFrame{
	{"ef8a35221a7405572804a277e965a7fdcce7768c7f646866b271b45af35737a1", introQuad{{148, 107}, {205, 97}, {208, 109}, {148, 122}}},
	{"3858a8793984ff7a7754d11308d441b06e4e3695a2f651cca28f7e7ba973fc2f", introQuad{{154, 105}, {216, 97}, {219, 111}, {154, 121}}},
	{"dbc1d7aa1471aaad52e87bf8c644c9d6d28eb888aad3932d2da2b190f652d091", introQuad{{161, 104}, {226, 98}, {229, 112}, {161, 121}}},
	{"7d1ecb00f365bea1298518803ebadf6ea9ab04f788dcaea151a37885844e8301", introQuad{{173, 104}, {240, 100}, {242, 113}, {173, 121}}},
	{"2996bda97e5591818aecdee6fe856fd1629c6ba6882983f7a7b8173771fc3e66", introQuad{{189, 105}, {252, 104}, {254, 115}, {189, 121}}},
	{"69220ee76660b723ca0c52d7c54ff86b0735192963ebac5c4962176081aa28de", introQuad{{194, 106}, {255, 106}, {256, 117}, {194, 120}}},
}

const introPaletteSHA = "d689ae136905afa64a9329ddd123426f7e0f3cde6108c034e0249dbab0a9727c"

type introPending struct {
	ret oracle.Addr
	sp  uint16
}
type introPictureState struct {
	english     bool
	data        *introPictureData
	paper       *introPaperData
	paperCached *introPaperView
	pending     *introPending
	complete    [32]byte
	hasComplete bool
	cached      *introPictureView
	lastReason  string
}

func (s *introPictureState) report(reason string) {
	if reason != s.lastReason {
		log.Printf("片頭圖片：%s", reason)
		s.lastReason = reason
	}
}
func (s *introPictureState) enter(ret oracle.Addr, sp uint16, rect [4]uint16) bool {
	if ret.Seg != 0x110 || ret.Off != 0x5b9b && ret.Off != 0x5afd {
		return false
	}
	if s.pending != nil {
		s.pending = nil
		s.hasComplete = false
		s.report("重入，撤除來源見證")
		return false
	}
	if rect != [4]uint16{0, 0, 319, 199} {
		s.hasComplete = false
		s.report("完整矩形不符")
		return false
	}
	s.pending = &introPending{ret, sp}
	return true
}
func (s *introPictureState) returned(ret oracle.Addr, sp uint16, indexed []byte) {
	if s == nil || s.pending == nil {
		return
	}
	p := s.pending
	if p.ret != ret {
		return
	}
	s.pending = nil
	if sp != p.sp+2 || len(indexed) != 64000 {
		s.hasComplete = false
		s.report("返回見證不符")
		return
	}
	s.complete = sha256.Sum256(indexed)
	s.hasComplete = true
}
func introPaletteHash(pal [256][3]uint8) [32]byte {
	b := make([]byte, 0, 768)
	for _, c := range pal {
		b = append(b, c[:]...)
	}
	return sha256.Sum256(b)
}
func (s *introPictureState) view(indexed []byte, pal [256][3]uint8, mode uint8) *introPictureView {
	if s == nil || len(indexed) != 64000 || mode != 0x13 || !s.hasComplete {
		return nil
	}
	h := sha256.Sum256(indexed)
	if h != s.complete {
		s.report("部分或未知畫格，保留原圖")
		return nil
	}
	var f *introFrame
	for i := range introFrames {
		if hex.EncodeToString(h[:]) == introFrames[i].sha {
			f = &introFrames[i]
			break
		}
	}
	if f == nil {
		if isPaperFrame(h) {
			return nil
		}
		s.report("非已量字樣畫格")
		return nil
	}
	if s.data == nil {
		if s.english {
			s.report("英文模式：保留片頭原圖")
		} else {
			s.report("呈獻缺譯：資源不可用")
		}
		return nil
	}
	ph := introPaletteHash(pal)
	if hex.EncodeToString(ph[:]) != introPaletteSHA {
		s.report("呈獻缺譯：未量色盤")
		return nil
	}
	if s.cached != nil && s.cached.sha == h {
		return s.cached
	}
	v, err := buildIntroPicture(indexed, pal, f.quad, s.data)
	if err != nil {
		s.report("呈獻缺譯：補底樣本不合格")
		return nil
	}
	v.sha = h
	s.cached = v
	s.report("呈獻：完整畫格已辨識")
	return v
}

type introPixel struct {
	x, y int
	rgb  [3]uint8
}
type introPictureView struct {
	sha    [32]byte
	quad   introQuad
	pixels []introPixel
	glyph  *image.RGBA
}

func (q introQuad) contains(x, y int) bool {
	for i, a := range q {
		b := q[(i+1)%4]
		if (b.X-a.X)*(2*y+1-2*a.Y)-(b.Y-a.Y)*(2*x+1-2*a.X) < 0 {
			return false
		}
	}
	return true
}
func (q introQuad) bounds() image.Rectangle {
	r := image.Rect(q[0].X, q[0].Y, q[0].X, q[0].Y)
	for _, p := range q {
		r.Min.X = min(r.Min.X, p.X)
		r.Min.Y = min(r.Min.Y, p.Y)
		r.Max.X = max(r.Max.X, p.X)
		r.Max.Y = max(r.Max.Y, p.Y)
	}
	return r
}
func introSolve(m [3][3]float64, rhs [3]float64) ([3]float64, error) {
	var a [3][4]float64
	for i := 0; i < 3; i++ {
		copy(a[i][:3], m[i][:])
		a[i][3] = rhs[i]
	}
	for k := 0; k < 3; k++ {
		q := k
		for i := k + 1; i < 3; i++ {
			if math.Abs(a[i][k]) > math.Abs(a[q][k]) {
				q = i
			}
		}
		a[k], a[q] = a[q], a[k]
		den := a[k][k]
		if math.Abs(den) <= 1e-9 || math.IsNaN(den) || math.IsInf(den, 0) {
			return [3]float64{}, fmt.Errorf("奇異補底矩陣")
		}
		for j := 0; j < 4; j++ {
			a[k][j] /= den
		}
		for i := 0; i < 3; i++ {
			if i != k {
				f := a[i][k]
				for j := 0; j < 4; j++ {
					a[i][j] -= f * a[k][j]
				}
			}
		}
	}
	v := [3]float64{a[0][3], a[1][3], a[2][3]}
	for _, n := range v {
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return v, fmt.Errorf("非有限補底係數")
		}
	}
	return v, nil
}
func buildIntroPicture(indexed []byte, pal [256][3]uint8, q introQuad, data *introPictureData) (*introPictureView, error) {
	r := q.bounds()
	if len(indexed) != 64000 || r.Min.X < 4 || r.Min.Y < 4 || r.Max.X+4 >= 320 || r.Max.Y+4 >= 200 {
		return nil, fmt.Errorf("補底範圍不合格")
	}
	var groups [4][]introPixel
	for y := r.Min.Y - 4; y <= r.Max.Y+4; y++ {
		for x := r.Min.X - 4; x <= r.Max.X+4; x++ {
			if q.contains(x, y) {
				continue
			}
			c := pal[indexed[y*320+x]]
			if c[0] > c[1] && c[1] > c[2] && int(c[0])-int(c[2]) >= 30 {
				g := x%2 + 2*(y%2)
				groups[g] = append(groups[g], introPixel{x, y, c})
			}
		}
	}
	var models [4][3][3]float64
	for group, samples := range groups {
		if len(samples) < 30 {
			return nil, fmt.Errorf("補底樣本不足")
		}
		var m [3][3]float64
		var rhs [3][3]float64
		for _, p := range samples {
			v := [3]float64{1, float64(p.x - 200), float64(p.y - 110)}
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					m[i][j] += v[i] * v[j]
				}
				for c := 0; c < 3; c++ {
					rhs[c][i] += v[i] * float64(p.rgb[c])
				}
			}
		}
		for c := 0; c < 3; c++ {
			var err error
			models[group][c], err = introSolve(m, rhs[c])
			if err != nil {
				return nil, err
			}
		}
	}
	v := &introPictureView{quad: q, glyph: introGlyph(data, pal[0])}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if !q.contains(x, y) {
				continue
			}
			group := x%2 + 2*(y%2)
			nearest := groups[group][0]
			distance := math.MaxInt
			for _, p := range groups[group] {
				d := (p.x-x)*(p.x-x) + (p.y-y)*(p.y-y)
				if d < distance {
					nearest = p
					distance = d
				}
			}
			p := introPixel{x: x, y: y}
			for c := 0; c < 3; c++ {
				coef := models[group][c]
				n := float64(nearest.rgb[c]) + coef[1]*float64(x-nearest.x) + coef[2]*float64(y-nearest.y)
				p.rgb[c] = uint8(max(0, min(255, math.RoundToEven(n))))
			}
			v.pixels = append(v.pixels, p)
		}
	}
	return v, nil
}
func introGlyph(data *introPictureData, fg [3]uint8) *image.RGBA {
	n := data.Size
	img := image.NewRGBA(image.Rect(0, 0, 2*n+2, n))
	for j, g := range data.Entries[0].G {
		for y, row := range g {
			for x := 0; x < n; x++ {
				if row&(1<<(n-1-x)) != 0 {
					img.SetRGBA(j*(n+2)+x, y, color.RGBA{fg[0], fg[1], fg[2], 255})
				}
			}
		}
	}
	return img
}
func (v *introPictureView) anchor(scale int) image.Point {
	x, y := 0, 0
	for _, p := range v.quad {
		x += p.X
		y += p.Y
	}
	return image.Pt(int(math.RoundToEven(float64(x*scale)/4))-v.glyph.Bounds().Dx()/2, int(math.RoundToEven(float64(y*scale)/4))-v.glyph.Bounds().Dy()/2)
}
func (v *introPictureView) patchRGBA(rgba []byte, scale int) {
	if v == nil || scale != 2 && scale != 3 || len(rgba) != 64000*4 {
		return
	}
	for _, p := range v.pixels {
		i := (p.y*320 + p.x) * 4
		copy(rgba[i:i+3], p.rgb[:])
	}
}
func (v *introPictureView) paintHD(img *image.RGBA, scale int) {
	if v == nil || scale != 2 && scale != 3 || img.Bounds() != image.Rect(0, 0, 320*scale, 200*scale) {
		return
	}
	for _, p := range v.pixels {
		c := color.RGBA{p.rgb[0], p.rgb[1], p.rgb[2], 255}
		for y := 0; y < scale; y++ {
			for x := 0; x < scale; x++ {
				img.SetRGBA(p.x*scale+x, p.y*scale+y, c)
			}
		}
	}
	a := v.anchor(scale)
	draw.Draw(img, v.glyph.Bounds().Add(a), v.glyph, image.Point{}, draw.Over)
}
func (g *game) drawIntroPicture(dst *ebiten.Image, v *introPictureView) {
	if v == nil || g.scale != 2 && g.scale != 3 || dst.Bounds() != image.Rect(0, 0, 320*g.scale, 200*g.scale) {
		return
	}
	if g.pictureGlyph != v.glyph {
		g.pictureImage = ebiten.NewImageFromImage(v.glyph)
		g.pictureGlyph = v.glyph
	}
	a := v.anchor(g.scale)
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(a.X), float64(a.Y))
	dst.DrawImage(g.pictureImage, &op)
}
