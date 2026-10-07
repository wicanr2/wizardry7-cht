package main

// 028：已確認紙頁顯示幾何。原版畫格、記憶體及作者英文字保持。
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type paperPolygon [][2]float64
type paperLabel struct {
	En     string       `json:"en"`
	Zh     string       `json:"zh"`
	Center [2]float64   `json:"center"`
	Color  string       `json:"color"`
	Safe   paperPolygon `json:"safe"`
}
type paperFrame struct {
	Step     uint64         `json:"step"`
	SHA      string         `json:"indexed_sha256"`
	Protect  []paperPolygon `json:"protect"`
	Bands    []paperPolygon `json:"bands"`
	Surfaces []paperPolygon `json:"band_surfaces"`
	Labels   []paperLabel   `json:"labels"`
}

var paperFrames = func() []paperFrame {
	h := sha256.Sum256([]byte(introPaperModelJSON))
	if hex.EncodeToString(h[:]) != paperModelSHA {
		panic("紙頁幾何來源不符")
	}
	var m struct {
		Frames []paperFrame `json:"frames"`
	}
	if json.Unmarshal([]byte(introPaperModelJSON), &m) != nil || len(m.Frames) != 6 {
		panic("紙頁幾何不合格")
	}
	return m.Frames
}()

func (p paperPolygon) contains(x, y float64) bool {
	hit := false
	for i, a := range p {
		b := p[(i+1)%len(p)]
		if (a[1] > y) != (b[1] > y) && x < (b[0]-a[0])*(y-a[1])/(b[1]-a[1])+a[0] {
			hit = !hit
		}
	}
	return hit
}
func paperProtected(f *paperFrame, x, y int) bool {
	for _, p := range f.Protect {
		if p.contains(float64(x)+.5, float64(y)+.5) {
			return true
		}
	}
	return false
}
func isPaperFrame(h [32]byte) bool {
	for _, f := range paperFrames {
		if hex.EncodeToString(h[:]) == f.SHA {
			return true
		}
	}
	return false
}

type introPaperView struct {
	sha    [32]byte
	pixels []introPixel
	glyphs []approvedPatch
}

func (s *introPictureState) paperView(indexed []byte, pal [256][3]uint8, mode uint8, scale int) *introPaperView {
	if s == nil || !s.hasComplete || len(indexed) != 64000 || mode != 0x13 {
		return nil
	}
	h := sha256.Sum256(indexed)
	if h != s.complete {
		return nil
	}
	var frame *paperFrame
	for i := range paperFrames {
		if hex.EncodeToString(h[:]) == paperFrames[i].SHA {
			frame = &paperFrames[i]
			break
		}
	}
	if frame == nil || len(frame.Labels) == 0 {
		return nil
	}
	if s.english || s.paper == nil {
		s.report("紙頁缺譯：字樣資源不可用或英文模式")
		return nil
	}
	if scale != 3 {
		s.report("紙頁缺譯：只定版3倍畫布，保留原圖")
		return nil
	}
	ph := introPaletteHash(pal)
	if hex.EncodeToString(ph[:]) != introPaletteSHA {
		s.report("紙頁缺譯：未量色盤")
		return nil
	}
	if s.paperCached != nil && s.paperCached.sha == h {
		return s.paperCached
	}
	v, err := buildIntroPaper(indexed, pal, frame, s.paper)
	if err != nil {
		s.report("紙頁缺譯：補底或字格不合格")
		return nil
	}
	v.sha = h
	s.paperCached = v
	s.report("紙頁：完整畫格已辨識")
	return v
}
func buildIntroPaper(indexed []byte, pal [256][3]uint8, f *paperFrame, data *introPaperData) (*introPaperView, error) {
	if len(indexed) != 64000 || data == nil || len(f.Bands) != len(f.Surfaces) {
		return nil, fmt.Errorf("紙頁來源不符")
	}
	masks := make([][]bool, len(f.Bands))
	all := make([]bool, 64000)
	for i, band := range f.Bands {
		masks[i] = make([]bool, 64000)
		for y := 0; y < 200; y++ {
			for x := 0; x < 320; x++ {
				if band.contains(float64(x)+.5, float64(y)+.5) && f.Surfaces[i].contains(float64(x)+.5, float64(y)+.5) {
					if paperProtected(f, x, y) {
						return nil, fmt.Errorf("字帶侵入作者區")
					}
					masks[i][y*320+x], all[y*320+x] = true, true
				}
			}
		}
	}
	v := &introPaperView{}
	for i, surface := range f.Surfaces {
		var groups [4][]introPixel
		for y := 0; y < 200; y++ {
			for x := 0; x < 320; x++ {
				c := pal[indexed[y*320+x]]
				if surface.contains(float64(x)+.5, float64(y)+.5) && !all[y*320+x] && !paperProtected(f, x, y) && c[0] > c[1] && c[1] > c[2] && 5*int(c[1]) >= 3*int(c[0]) && int(c[1])-int(c[2]) >= 16 {
					g := x%2 + 2*(y%2)
					groups[g] = append(groups[g], introPixel{x, y, c})
				}
			}
		}
		for _, group := range groups {
			if len(group) < 8 {
				return nil, fmt.Errorf("紙頁取樣不足")
			}
		}
		for y := 0; y < 200; y++ {
			for x := 0; x < 320; x++ {
				if !masks[i][y*320+x] {
					continue
				}
				// 候選先依y/x遍歷；等距保持先入者，與R974排序相同。
				var nearest [8]introPixel
				var dist [8]int
				for n := range dist {
					dist[n] = math.MaxInt
				}
				for _, p := range groups[x%2+2*(y%2)] {
					d := (p.x-x)*(p.x-x) + (p.y-y)*(p.y-y)
					for n := range dist {
						if d < dist[n] {
							copy(dist[n+1:], dist[n:7])
							copy(nearest[n+1:], nearest[n:7])
							dist[n], nearest[n] = d, p
							break
						}
					}
				}
				p := introPixel{x: x, y: y}
				var weights [8]float64
				total := 0.0
				for n, d := range dist {
					weights[n] = 1 / float64(d+1)
					total += weights[n]
				}
				for c := 0; c < 3; c++ {
					sum := 0.0
					for n, sample := range nearest {
						sum += weights[n] * float64(sample.rgb[c])
					}
					p.rgb[c] = uint8(math.RoundToEven(sum / total))
				}
				v.pixels = append(v.pixels, p)
			}
		}
	}
	labels := append([]paperLabel(nil), f.Labels...)
	if f.Step == 73235759 {
		labels[0].Center = [2]float64{199, 68}
	}
	if f.Step == 75322912 {
		safe := labels[0].Safe
		labels = []paperLabel{
			{"fantasy", "奇幻", [2]float64{93, 32}, "red", safe},
			{"role-playing", "角色", [2]float64{106, 52}, "red", safe},
			{"role-playing", "扮演", [2]float64{126, 73}, "red", safe},
			{"simulation", "遊戲", [2]float64{145, 94}, "red", safe}, labels[len(labels)-1],
		}
	}
	if f.Step == 100999086 {
		labels[len(labels)-1].Center = [2]float64{173.5, 88}
	}
	fonts := map[rune][]uint16{}
	for _, e := range data.Entries {
		for i, c := range []rune(e.Zh) {
			fonts[c] = e.G[i]
		}
	}
	occupied := map[image.Point]bool{}
	for _, l := range labels {
		size := data.Size
		width, height := len([]rune(l.Zh))*(2*size+4)-4+1, 2*size
		at := image.Pt(int(math.RoundToEven(l.Center[0]*3))-width/2, int(math.RoundToEven(l.Center[1]*3))-height/2)
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				p := image.Pt(at.X+x, at.Y+y)
				if p.X < 0 || p.Y < 0 || p.X >= 960 || p.Y >= 600 || !l.Safe.contains((float64(p.X)+.5)/3, (float64(p.Y)+.5)/3) || paperProtected(f, p.X/3, p.Y/3) || occupied[p] {
					return nil, fmt.Errorf("紙頁字格越界或重疊")
				}
				occupied[p] = true
			}
		}
		fg := color.RGBA{248, 232, 160, 255}
		if l.Color == "red" {
			fg = color.RGBA{236, 72, 40, 255}
		}
		for j, c := range []rune(l.Zh) {
			glyph := fonts[c]
			if len(glyph) != size {
				return nil, fmt.Errorf("紙頁缺字")
			}
			for y, bits := range glyph {
				for x := 0; x < size; x++ {
					if bits&(1<<(size-1-x)) != 0 {
						for dy := 0; dy < 2; dy++ {
							for dx := 0; dx < 3; dx++ {
								img.SetRGBA(j*(2*size+4)+x*2+dx, y*2+dy, fg)
							}
						}
					}
				}
			}
		}
		v.glyphs = append(v.glyphs, approvedPatch{at, img})
	}
	return v, nil
}
func (v *introPaperView) patchRGBA(rgba []byte, scale int) {
	if v == nil || scale != 3 || len(rgba) != 64000*4 {
		return
	}
	for _, p := range v.pixels {
		i := (p.y*320 + p.x) * 4
		copy(rgba[i:i+3], p.rgb[:])
	}
}
func (v *introPaperView) paintHD(img *image.RGBA, scale int) {
	if v == nil || scale != 3 || img.Bounds() != image.Rect(0, 0, 960, 600) {
		return
	}
	for _, p := range v.pixels {
		for y := 0; y < 3; y++ {
			for x := 0; x < 3; x++ {
				img.SetRGBA(p.x*3+x, p.y*3+y, color.RGBA{p.rgb[0], p.rgb[1], p.rgb[2], 255})
			}
		}
	}
	for _, p := range v.glyphs {
		draw.Draw(img, p.img.Bounds().Add(p.at), p.img, image.Point{}, draw.Over)
	}
}
func (g *game) drawIntroPaper(dst *ebiten.Image, v *introPaperView) {
	if v == nil || g.scale != 3 || dst.Bounds() != image.Rect(0, 0, 960, 600) {
		return
	}
	if g.paperView != v {
		g.paperView, g.paperImages = v, nil
		for _, p := range v.glyphs {
			g.paperImages = append(g.paperImages, ebiten.NewImageFromImage(p.img))
		}
	}
	for i, p := range v.glyphs {
		var op ebiten.DrawImageOptions
		op.GeoM.Translate(float64(p.at.X), float64(p.at.Y))
		dst.DrawImage(g.paperImages[i], &op)
	}
}
