package main

// F1 說明畫面（docs/spec/003）。資料由 tools/re/bake_help.py 烘製成 help.json：
// 逐字字形（全形寬 15、半形寬 8）＋頁面。畫面辨識依開啟當下的 shown 事件（R4）。

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type helpGlyph struct {
	W    int   `json:"w"`
	Rows []int `json:"rows"` // 每列一個整數，位元對齊字寬：像素 i 為 bit (W-1-i)
}

type helpPage struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Lines []string `json:"lines"`
}

type helpData struct {
	Glyphs map[string]helpGlyph `json:"glyphs"`
	Pages  []helpPage           `json:"pages"`
}

func loadHelp(path string) *helpData {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("說明：讀不到 %s（%v）", path, err)
		return nil
	}
	var h helpData
	if err := json.Unmarshal(b, &h); err != nil {
		log.Printf("說明：%s 格式錯誤（%v）", path, err)
		return nil
	}
	log.Printf("說明：載入 %d 頁、%d 字形", len(h.Pages), len(h.Glyphs))
	return &h
}

// detectScreen 依 shown 事件的英文原文判定目前畫面（spec 003 R4，順序即優先序）。
// 防拷畫面時主選單字串也在 shown（docs/re/015 H6），所以防拷必須先判。
func detectScreen(strs []string) string {
	has := func(s string) bool {
		for _, v := range strs {
			if v == s {
				return true
			}
		}
		return false
	}
	switch {
	case has("SECURITY CLEARANCE"):
		return "copyprot"
	case has("OK") && has("CANCEL") && has("ADD CHARACTER"):
		return "addchar"
	case has("CHARACTER MENU") || has("START NEW GAME"):
		return "mainmenu"
	}
	return "general"
}

// helpPages 回說明頁的順序：快捷鍵、基本操作、目前畫面、譯名對照（R3）。
func (h *helpData) helpPages(screen string) []helpPage {
	var keys, ops, scr, general []helpPage
	var names []helpPage
	for _, p := range h.Pages {
		switch {
		case p.ID == "keys":
			keys = append(keys, p)
		case p.ID == "ops":
			ops = append(ops, p)
		case p.ID == "screen."+screen:
			scr = append(scr, p)
		case p.ID == "screen.general":
			general = append(general, p)
		case strings.HasPrefix(p.ID, "names."):
			names = append(names, p)
		}
	}
	if len(scr) == 0 {
		scr = general
	}
	out := append(append(append(keys, ops...), scr...), names...)
	return out
}

// helpUI 是頁尾樣板（TSV 的 ui 頁，不列為說明頁）。
func (h *helpData) footer(n, total int) string {
	for _, p := range h.Pages {
		if p.ID == "ui" && len(p.Lines) > 0 {
			return strings.NewReplacer("{n}", fmt.Sprint(n), "{N}", fmt.Sprint(total)).Replace(p.Lines[0])
		}
	}
	return ""
}

// drawText 以說明字形在 img 上畫一行字，回傳畫到的右緣。缺字形的字跳過並記 log 一次。
func (h *helpData) drawText(img *image.RGBA, x, y int, s string, c color.RGBA) int {
	for _, r := range s {
		g, ok := h.Glyphs[string(r)]
		if !ok {
			if !helpMissingLogged[r] {
				helpMissingLogged[r] = true
				log.Printf("說明：缺字形 %q", r)
			}
			continue
		}
		for row, bits := range g.Rows {
			for i := 0; i < g.W; i++ {
				if (bits>>(g.W-1-i))&1 != 0 {
					img.Set(x+i, y+row, c)
				}
			}
		}
		x += g.W + 2
	}
	return x
}

var helpMissingLogged = map[rune]bool{}

const (
	helpLineH  = 20
	helpTop    = 64
	helpMargin = 32
)

// renderHelp 把一頁說明畫成整個視窗大小的圖（面板＋文字）。行數超過兩欄容量時截斷並記 log。
func (h *helpData) renderHelp(w, hgt int, pages []helpPage, idx int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, hgt))
	for i := 3; i < len(img.Pix); i += 4 {
		img.Pix[i-3], img.Pix[i-2], img.Pix[i-1], img.Pix[i] = 0, 0, 16, 0xE0
	}
	p := pages[idx]
	h.drawText(img, helpMargin, 24, p.Title, color.RGBA{0xFF, 0xDD, 0x44, 0xFF})
	perCol := (hgt - helpTop - 48) / helpLineH
	cols := 1
	if len(p.Lines) > perCol {
		cols = 2
	}
	for i, line := range p.Lines {
		col, row := i/perCol, i%perCol
		if col >= cols {
			log.Printf("說明：頁 %s 有 %d 行，超出 %d 行容量，截斷", p.ID, len(p.Lines), perCol*cols)
			break
		}
		x := helpMargin + col*(w/2)
		h.drawText(img, x, helpTop+row*helpLineH, line, color.RGBA{0xF0, 0xF0, 0xF0, 0xFF})
	}
	h.drawText(img, helpMargin, hgt-32, h.footer(idx+1, len(pages)), color.RGBA{0xA0, 0xA0, 0xA0, 0xFF})
	return ebiten.NewImageFromImage(img)
}

// drawHelp 在 Draw 最後疊上說明。help.json 缺檔時用 Ebiten 內建 ASCII 字提示。
func (g *game) drawHelp(dst *ebiten.Image, screen string, page int) {
	if g.help == nil {
		ebitenutil.DebugPrintAt(dst, "help.json missing (tools/re/bake_help.py)", 16, 16)
		return
	}
	pages := g.help.helpPages(screen)
	if len(pages) == 0 {
		return
	}
	page = ((page % len(pages)) + len(pages)) % len(pages)
	key := fmt.Sprintf("%s/%d", screen, page)
	if g.helpImg == nil || g.helpKey != key {
		b := dst.Bounds()
		g.helpImg = g.help.renderHelp(b.Dx(), b.Dy(), pages, page)
		g.helpKey = key
	}
	dst.DrawImage(g.helpImg, nil)
}
