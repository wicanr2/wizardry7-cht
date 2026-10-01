package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// spec 003 R6 驗收 6：R4 辨識規則對實測 shown 事件的判定。資料在 workplace/（gitignore），缺檔就 skip。
func TestDetectScreenOnRecordedShown(t *testing.T) {
	cases := []struct {
		path, want string
	}{
		{"../../workplace/r663/zh3/hdat-272000000.json", "copyprot"},
		{"../../workplace/r663/zh3/hdat-286000000.json", "copyprot"},
		{"../../workplace/r665/maze/hdat-1470000000.json", "general"},
	}
	for _, c := range cases {
		b, err := os.ReadFile(c.path)
		if err != nil {
			t.Skipf("缺實測資料 %s", c.path)
		}
		var j struct{ Shown []struct{ S string } }
		if err := json.Unmarshal(b, &j); err != nil {
			t.Fatal(err)
		}
		var strs []string
		for _, e := range j.Shown {
			strs = append(strs, e.S)
		}
		if got := detectScreen(strs); got != c.want {
			t.Errorf("%s：判為 %s，應為 %s", c.path, got, c.want)
		}
	}
}

// spec007：以真實烘製字形驗證全部說明內容；沒有本機資料就skip。
// 期望是來源字元完整保留及欄界，而非複製被測折行演算法。
func TestHelpNativeFontContainsEverySourceCharacter(t *testing.T) {
	var receipts []map[string]any
	for _, backend := range []string{"eten", "unifont"} {
		path := "../../workplace/r682/" + backend + "/help.json"
		b, err := os.ReadFile(path)
		if err != nil {
			t.Skipf("缺字型驗收資料 %s", path)
		}
		var h helpData
		if err := json.Unmarshal(b, &h); err != nil {
			t.Fatal(err)
		}
		for _, scale := range []int{2, 3} {
			w, height := 320*scale, 200*scale
			for _, screen := range []string{"mainmenu", "copyprot", "addchar", "general"} {
				pages := h.layoutHelp(screen, w, height)
				if len(pages) == 0 {
					t.Fatalf("%s/%dx/%s 無可讀說明", backend, scale, screen)
				}
				var rendered []map[string]any
				for _, p := range pages {
					rendered = append(rendered, map[string]any{"id": p.ID, "title": p.Title,
						"lines": p.Lines, "columns": p.Columns})
				}
				receipts = append(receipts, map[string]any{"backend": backend, "scale": scale,
					"screen": screen, "pages": rendered})
				actual := map[string]string{}
				perCol := (height - helpTop - 48) / helpLineH
				for _, p := range pages {
					actual[p.ID] += strings.Join(p.Lines, "")
					if len(p.Lines) > p.Columns*perCol {
						t.Fatalf("%s 頁高超界", p.ID)
					}
					for i, line := range p.Lines {
						x := helpMargin + (i/perCol)*(w/2)
						y := helpTop + (i%perCol)*helpLineH
						limit := w - helpMargin
						if p.Columns == 2 {
							limit = (i/perCol+1)*(w/2) - helpMargin
						}
						for _, r := range line {
							g, ok := h.Glyphs[string(r)]
							if !ok || x+g.W > limit || y+len(g.Rows) > height-48 {
								t.Fatalf("%s/%dx/%s 字格超界或缺字 %U", backend, scale, p.ID, r)
							}
							x += g.W + 2
						}
					}
				}
				for _, source := range h.helpPages(screen) {
					if actual[source.ID] != strings.Join(source.Lines, "") {
						t.Fatalf("%s/%dx/%s 說明內容遺失", backend, scale, source.ID)
					}
				}
			}
		}
	}
	if path := os.Getenv("WIZ7_HELP_LAYOUT_RECEIPT"); path != "" {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("排版收據已存在或無法存取：%s", path)
		}
		b, err := json.MarshalIndent(receipts, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// 規則順序：主選單字串與防拷字串同時存在時必須判為防拷（docs/re/015 H6）。
func TestDetectScreenOrder(t *testing.T) {
	if got := detectScreen([]string{"CHARACTER MENU", "START NEW GAME", "SECURITY CLEARANCE"}); got != "copyprot" {
		t.Errorf("判為 %s", got)
	}
	if got := detectScreen([]string{"CHARACTER MENU", "ADD CHARACTER", "OK", "CANCEL"}); got != "addchar" {
		t.Errorf("判為 %s", got)
	}
	if got := detectScreen([]string{"CHARACTER MENU"}); got != "mainmenu" {
		t.Errorf("判為 %s", got)
	}
}

// spec007：隔離輸入檔只放t.TempDir，不修改正式烘製或字型。
func writeNativeFontLoaderFixture(t *testing.T, data any) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := t.TempDir() + "/font-loader.json"
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func nativeFontLoaderRows(height, first int) []int {
	rows := make([]int, height)
	if height > 0 {
		rows[0] = first
	}
	return rows
}

// 全形字須為原生正方形；ASCII寬8可搭配15或16列。
// 一個壞字形必須拒載整份資料，不留下可被誤用的部分資料。
func TestLoadHelpRejectsInvalidNativeGlyphs(t *testing.T) {
	cases := []struct {
		name          string
		width, height int
		first         int
	}{
		{"寬度零", 0, 16, 0},
		{"寬度七", 7, 16, 0},
		{"寬度十四", 14, 15, 0},
		{"寬度十七", 17, 16, 0},
		{"列數十四", 15, 14, 0},
		{"列數十七", 16, 17, 0},
		{"十六寬只有十五列", 16, 15, 0},
		{"十五寬卻有十六列", 15, 16, 0},
		{"ASCII位元超界", 8, 16, 256},
		{"十五位元超界", 15, 15, 32768},
		{"十六位元超界", 16, 16, 65536},
		{"負位元", 16, 16, -1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := writeNativeFontLoaderFixture(t, map[string]any{
				"glyphs": map[string]any{
					"一": map[string]any{"w": 16, "rows": nativeFontLoaderRows(16, 1)},
					"甲": map[string]any{"w": c.width, "rows": nativeFontLoaderRows(c.height, c.first)},
				},
				"pages": []helpPage{{ID: "general", Title: "一", Lines: []string{"甲"}}},
			})
			if loadHelp(path) != nil {
				t.Fatal("不合法字形仍載入說明資料")
			}
		})
	}
}

func TestLoadHelpAcceptsNativeGlyphBoundaries(t *testing.T) {
	path := writeNativeFontLoaderFixture(t, map[string]any{
		"glyphs": map[string]any{
			"一": map[string]any{"w": 15, "rows": nativeFontLoaderRows(15, 32767)},
			"甲": map[string]any{"w": 16, "rows": nativeFontLoaderRows(16, 65535)},
			"A": map[string]any{"w": 8, "rows": nativeFontLoaderRows(15, 255)},
			"B": map[string]any{"w": 8, "rows": nativeFontLoaderRows(16, 255)},
		},
		"pages": []helpPage{{ID: "general", Title: "一", Lines: []string{"甲AB"}}},
	})
	if got := loadHelp(path); got == nil || len(got.Glyphs) != 4 || len(got.Pages) != 1 {
		t.Fatal("合法15/16全形及8位寬ASCII邊界字形被拒載")
	}
}

func TestLoadHelpAcceptsActualNativeBakes(t *testing.T) {
	for _, backend := range []string{"eten", "unifont"} {
		t.Run(backend, func(t *testing.T) {
			path := "../../workplace/r682/" + backend + "/help.json"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Skip("缺本機實際烘製字型資料")
			}
			h := loadHelp(path)
			if h == nil || len(h.Glyphs) == 0 || len(h.Pages) == 0 {
				t.Fatal("實際烘製F1資料被拒載或為空")
			}
		})
	}
}
