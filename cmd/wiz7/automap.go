package main

// 即時地圖：自動繪圖模式（docs/spec/005，證據 docs/re/016 P1–P11）。
// 只讀原版：面向 1106:82EE、位置記錄 1106:8F3C 起 6 個 word（+2 區域、+8 X、+10 Y）。

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

var (
	addrFacing = oracle.Addr{Seg: 0x1106, Off: 0x82EE} // 線性 1934E（016 P2）
	addrPosRec = oracle.Addr{Seg: 0x1106, Off: 0x8F3C} // 線性 19F9C（016 P4）
)

type cell struct{ Area, X, Y int }

// automap 保存本次執行走過的格與目前位置。只在機器 goroutine 的快照節點更新。
type automap struct {
	have    bool
	cur     cell
	face    int
	visited map[cell]bool
	order   []cell // 記錄順序（JSON 輸出用，保持決定性）
}

func newAutomap() *automap { return &automap{visited: map[cell]bool{}} }

// update 讀原版記憶體並依 R2 記格：同區域、恰一軸 ±1 才記移動前後兩格；其他情形只更新目前位置。
func (m *automap) update(o *oracle.Oracle) {
	rec := o.Bytes(addrPosRec, 12)
	w := func(i int) int { return int(rec[i]) | int(rec[i+1])<<8 }
	now := cell{Area: w(2), X: w(8), Y: w(10)}
	face := int(o.Bytes(addrFacing, 1)[0])
	if m.have && now != m.cur && now.Area == m.cur.Area {
		dx, dy := now.X-m.cur.X, now.Y-m.cur.Y
		if (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1)) {
			m.mark(m.cur)
			m.mark(now)
		}
	}
	m.cur, m.face, m.have = now, face, true
}

func (m *automap) mark(c cell) {
	if !m.visited[c] {
		m.visited[c] = true
		m.order = append(m.order, c)
		log.Printf("地圖 記格 區域=%d (%d,%d)", c.Area, c.X, c.Y)
	}
}

// mapView 是交給 Draw 的地圖快照（避免 Draw 與機器 goroutine 共用 map）。
type mapView struct {
	cur   cell
	face  int
	cells []cell // 目前區域走過的格
}

func (m *automap) view() mapView {
	v := mapView{cur: m.cur, face: m.face}
	for _, c := range m.order {
		if c.Area == m.cur.Area {
			v.cells = append(v.cells, c)
		}
	}
	return v
}

const (
	mapCell  = 12
	mapCells = 15 // 面板邊長（格），目前位置在正中央
)

// drawMap 在視窗右上角畫自動地圖面板（R3）：X 向右、Y 向下；面向 3 上、0 右、1 下、2 左。
func (g *game) drawMap(dst *ebiten.Image, v mapView) {
	size := mapCell * mapCells
	img := image.NewRGBA(image.Rect(0, 0, size+16, size+40))
	for i := 3; i < len(img.Pix); i += 4 {
		img.Pix[i-3], img.Pix[i-2], img.Pix[i-1], img.Pix[i] = 0, 0, 16, 0xC8
	}
	ox, oy := 8, 32
	half := mapCells / 2
	fill := func(cx, cy int, c color.RGBA, inset int) {
		for y := inset; y < mapCell-inset; y++ {
			for x := inset; x < mapCell-inset; x++ {
				img.Set(ox+cx*mapCell+x, oy+cy*mapCell+y, c)
			}
		}
	}
	title, empty := "自動地圖", "尚未記錄移動"
	if g.help != nil {
		for _, p := range g.help.Pages {
			if p.ID == "ui.map" && len(p.Lines) >= 1 {
				title = p.Title
				empty = p.Lines[0]
			}
		}
		g.help.drawText(img, 8, 8, title, color.RGBA{0xFF, 0xDD, 0x44, 0xFF})
	}
	for _, c := range v.cells {
		cx, cy := c.X-v.cur.X+half, c.Y-v.cur.Y+half
		if cx >= 0 && cx < mapCells && cy >= 0 && cy < mapCells {
			fill(cx, cy, color.RGBA{0x70, 0x90, 0x70, 0xFF}, 1)
		}
	}
	// 目前位置：箭頭（三角形）指向面向
	arrow := color.RGBA{0xFF, 0x50, 0x40, 0xFF}
	bx, by := ox+half*mapCell, oy+half*mapCell
	for i := 0; i < mapCell/2; i++ {
		for j := -i; j <= i; j++ {
			var x, y int
			switch v.face {
			case 3: // 上
				x, y = mapCell/2+j, 1+i
			case 1: // 下
				x, y = mapCell/2+j, mapCell-2-i
			case 0: // 右
				x, y = mapCell-2-i, mapCell/2+j
			default: // 2 左
				x, y = 1+i, mapCell/2+j
			}
			img.Set(bx+x, by+y, arrow)
		}
	}
	if len(v.cells) == 0 && g.help != nil {
		g.help.drawText(img, 8, oy+size+2-18, empty, color.RGBA{0xC0, 0xC0, 0xC0, 0xFF})
	}
	b := dst.Bounds()
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(b.Dx()-img.Bounds().Dx()-8), 8)
	dst.DrawImage(ebiten.NewImageFromImage(img), &op)
}

// mapJSON 是 -hd-at 輸出的地圖狀態（R4）。
type mapJSON struct {
	Area, X, Y, Face int
	Visited          [][3]int
}

func (m *automap) json() mapJSON {
	j := mapJSON{Area: m.cur.Area, X: m.cur.X, Y: m.cur.Y, Face: m.face}
	for _, c := range m.order {
		j.Visited = append(j.Visited, [3]int{c.Area, c.X, c.Y})
	}
	return j
}

func (c cell) String() string { return fmt.Sprintf("%d:(%d,%d)", c.Area, c.X, c.Y) }
