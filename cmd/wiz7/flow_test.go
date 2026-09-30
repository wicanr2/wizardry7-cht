package main

import (
	"reflect"
	"testing"
)

// R4 排版：期望值取自 tools/re/check_flow.py 的獨立排版（含 …… ── 不可拆、避頭尾、連續禁則字）。
func TestFlowLayout(t *testing.T) {
	cases := []struct {
		zh   string
		n    int
		rows []string
	}{
		{"微風輕輕吹過樹林，數月來你第一次感覺到自己活著……", 32, []string{"微風輕輕吹過樹林，數月來你第一次感覺到自己活著……"}},
		{"終於擺脫那台載你來此的奇異飛行器悶窒的束縛，暢快", 21, []string{"終於擺脫那台載你來此的奇異飛行器悶窒的束", "縛，暢快"}},
		{"一二三四五六，七八九十", 6, []string{"一二三四五", "六，七八九十"}},
		{"一二三四五六七八九十……十一", 10, []string{"一二三四五六七八九", "十……十一"}},
		{"一二三四五（六七八九十）", 6, []string{"一二三四五", "（六七八九", "十）"}},
		{"甲乙丙丁戊己庚辛壬癸」』", 10, []string{"甲乙丙丁戊己庚辛壬", "癸」』"}},
		{"一二三四五六七八九——十", 9, []string{"一二三四五六七八", "九——十"}},
	}
	for _, c := range cases {
		if got := flowLayout(c.zh, c.n); !reflect.DeepEqual(got, c.rows) {
			t.Errorf("flowLayout(%q,%d)＝%q，期望 %q", c.zh, c.n, got, c.rows)
		}
	}
}

// R4：字形放大 gs、每列字數 N 與置中位置（Draw 端 hd＝scale÷2，hdFrame 端 hd＝1）。
func TestFlowColsAndPos(t *testing.T) {
	gsCases := []struct{ k, hdf, gs int }{{2, 1, 1}, {3, 1, 2}, {4, 1, 2}, {4, 2, 2}, {5, 2, 2}, {6, 3, 3}}
	for _, c := range gsCases {
		if got := flowGS(c.k, c.hdf); got != c.gs {
			t.Errorf("flowGS(%d,%d)＝%d，期望 %d", c.k, c.hdf, got, c.gs)
		}
	}
	colCases := []struct{ k, gs, n int }{{3, 2, 16}, {2, 1, 21}, {4, 2, 21}, {3, 1, 32}}
	for _, c := range colCases {
		if got := flowCols(c.k, c.gs); got != c.n {
			t.Errorf("flowCols(%d,%d)＝%d，期望 %d", c.k, c.gs, got, c.n)
		}
	}
	x, y := flowRowPos("一二三四", 1, 3, 2) // n＝4，pitch 34：160×3−68＝412；y＝(146＋12)×3
	if x != 412 || y != 474 {
		t.Errorf("flowRowPos＝(%d,%d)，期望 (412,474)", x, y)
	}
}

func testFlow() *flowData {
	fd := &flowData{
		Glyphs: map[string][]uint16{"一": nil, "二": nil, "三": nil, "四": nil},
		Paras: []flowPara{
			{ID: 100, Words: []string{"a", "b", "c", "d", "e", "f"},
				Lines: []flowLine{{ID: 100, End: 3, Zh: "一"}, {ID: 101, End: 6, Zh: "二"}}},
			{ID: 200, Words: []string{"a", "b", "c"}, Lines: []flowLine{{ID: 200, End: 3, Zh: "三"}}},
			{ID: 300, Words: []string{"x", "y", "z", "w"},
				Lines: []flowLine{{ID: 300, End: 2, Zh: "一"}, {ID: 301, End: 4, Zh: ""}}},
			{ID: 400, Words: []string{"p", "q", "r", "s"},
				Lines: []flowLine{{ID: 400, End: 2, Zh: "四"}, {ID: 401, End: 4, Zh: "二"}}},
		},
		byWord: map[string][]flowPos{},
	}
	for pi, p := range fd.Paras {
		for a, w := range p.Words {
			fd.byWord[w] = append(fd.byWord[w], flowPos{int32(pi), int32(a)})
		}
	}
	return fd
}

// R3：完成的行、頁從行中間開始、歧義、無譯文、缺字形。
func TestFlowMatch(t *testing.T) {
	fd := testFlow()
	cases := []struct {
		name   string
		words  []string
		zh     string
		status string
		lines  []int
	}{
		{"整段兩行", []string{"p", "q", "r", "s"}, "四二", "400:0", []int{400, 401}},
		{"只完成第一行", []string{"p", "q", "r"}, "四", "400:0", []int{400}},
		{"頁從行中間開始：尾巴完成該行", []string{"r", "s"}, "二", "400:2", []int{401}},
		{"頁只有行的中段：沒有完成的行", []string{"x"}, "", "無譯文", nil},
		{"歧義：兩個段落都有 a b c，譯文不同", []string{"a", "b", "c"}, "", "歧義", nil},
		{"零候選", []string{"nope"}, "", "無譯文", nil},
		{"完成的行譯文為空", []string{"x", "y", "z", "w"}, "", "無譯文", nil},
	}
	for _, c := range cases {
		res := fd.match(c.words)
		if res.Zh != c.zh || res.Status != c.status || !reflect.DeepEqual(res.Lines, c.lines) {
			t.Errorf("%s：Zh=%q Status=%q Lines=%v，期望 %q %q %v", c.name, res.Zh, res.Status, res.Lines, c.zh, c.status, c.lines)
		}
	}
	// 多個候選但譯文相同：採用
	fd.Paras = append(fd.Paras, flowPara{ID: 500, Words: []string{"p", "q", "r", "s"},
		Lines: []flowLine{{ID: 500, End: 2, Zh: "四"}, {ID: 501, End: 4, Zh: "二"}}})
	fd.byWord = map[string][]flowPos{}
	for pi, p := range fd.Paras {
		for a, w := range p.Words {
			fd.byWord[w] = append(fd.byWord[w], flowPos{int32(pi), int32(a)})
		}
	}
	if res := fd.match([]string{"p", "q", "r", "s"}); res.Zh != "四二" || res.Cands != 2 {
		t.Errorf("譯文相同的多候選：Zh=%q Cands=%d，期望 四二／2", res.Zh, res.Cands)
	}
	// 缺字形
	delete(fd.Glyphs, "二")
	if res := fd.match([]string{"p", "q", "r", "s"}); res.Zh != "" || res.Status != "缺字形" {
		t.Errorf("缺字形：Zh=%q Status=%q", res.Zh, res.Status)
	}
}

func TestStripCtl(t *testing.T) {
	if got := stripCtl("a\x01b\x1fc\x7f \x00"); got != "abc\x7f " {
		t.Errorf("stripCtl＝%q", got)
	}
}
