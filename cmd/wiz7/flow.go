package main

// 流動文字（docs/spec/006，證據 docs/re/017）：故事與旁白由引擎把連續訊息重新折行、
// 依列數分頁印出，印字入口看到的是折行後的片段。這裡以整頁的字詞序列比對訊息段落，
// 取該頁「完成的行」的譯文，依中文寬度重排、置中；英文筆畫沿用 spec 001 的還原機制蓋掉。

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/dosgolem/oracle"
)

// flowSites 是流動文字呼叫點（返回位址）。新增呼叫點要先有 docs/re/017 F13 那一組量測（spec 006 範圍）。
type flowSite struct {
	width, x0, x1, centerMin, centerMax int
	reader                              uint16
}

var flowSites = map[oracle.Addr]flowSite{
	{Seg: 0x0110, Off: 0xBCFC}: {184, 68, 251, 318, 322, 0},      // 旁白（F1）
	{Seg: 0x0110, Off: 0xBE0C}: {296, 12, 307, 316, 322, 0xBFEB}, // 故事（F12、F17）
}

var storySite = oracle.Addr{Seg: 0x0110, Off: 0xBE0C}

// NPC只共用排版及文字帶幾何，不加入一般流文／故事的入口仲裁。
func flowGeometry(site oracle.Addr) flowSite {
	if site == npcTextSite {
		return npcTextGeometry
	}
	return flowSites[site]
}

func (s flowSite) cols(k, gs int) int { return max(1, s.width*k/(flowGlyphPitch*gs)) }

func (fd *flowData) size() int {
	if fd == nil || fd.Size == 0 {
		return 15 // spec007：相容原有資料
	}
	return fd.Size
}

func (fd *flowData) cols(s flowSite, k, gs int) int {
	return max(1, s.width*k/((fd.size()+2)*gs))
}

func (fd *flowData) rowPos(text string, r, k, gs int) (int, int) {
	return flowCenter*k - len([]rune(text))*(fd.size()+2)*gs/2, (flowY0 + flowRowStep*r) * k
}

type flowMessage struct {
	W string `json:"w"`
	T int    `json:"t"`
	Z string `json:"z"`
}
type flowRead struct {
	Step uint64
	ID   int
}
type storyReturn struct {
	ret oracle.Addr
	sp  uint16
}

const (
	flowFont       = 2       // 字型 ID（F13）
	flowFg         = 0x55    // 筆畫色（F13）
	flowMaxFrags   = 4       // 一頁至多幾個片段（F2）
	flowY0         = 146     // 首列 y0（F2）
	flowRowStep    = 12      // 列距（F2）
	flowCenter     = 160     // 列的水平中心（F3）
	flowWidth      = 184     // 列寬上限（遊戲像素，F5、F13）
	flowMaxRows    = 4       // 中文最多列數
	flowQuiet      = 200_000 // 最後片段之後至少過這麼多步才畫（中間態不畫）
	flowBandX0     = 68      // 文字帶（F2、F3）
	flowBandX1     = 251     // 含
	flowBandY0     = 146     // 含
	flowBandY1     = 194     // 含
	flowGlyphPitch = 17      // 15 px 字形＋2 px 字距
)

type flowLine struct {
	ID  int    `json:"id"`
	End int    `json:"end"`
	Zh  string `json:"zh"`
}

type flowPara struct {
	ID    int        `json:"id"`
	Words []string   `json:"words"`
	Lines []flowLine `json:"lines"`
}

type flowPos struct{ P, A int32 }

// flowData 是 flow.json（tools/re/bake_flow.py）載入後的資料與索引。
type flowData struct {
	Size   int                 `json:"size"` // spec007：原生15／16像素
	Glyphs map[string][]uint16 `json:"glyphs"`
	Paras  []flowPara          `json:"paras"`
	Msgs   map[int]flowMessage `json:"msgs"`
	byWord map[string][]flowPos
}

func loadFlow(path string) *flowData {
	if path == "" {
		log.Printf("流文：已停用（-flow \"\"）")
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("流文：讀不到 %s（%v），流動文字關閉", path, err)
		return nil
	}
	var fd flowData
	if err := json.Unmarshal(b, &fd); err != nil {
		log.Printf("流文：%s 格式錯誤（%v），流動文字關閉", path, err)
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return nil
	}
	_, explicitSize := fields["size"]
	if (explicitSize && fd.Size != 15 && fd.Size != 16) || (fd.size() != 15 && fd.size() != 16) {
		log.Printf("流文：%s 字形尺寸不符", path)
		return nil
	}
	for _, rows := range fd.Glyphs {
		if len(rows) != fd.size() {
			log.Printf("流文：%s 字形列數不符", path)
			return nil
		}
		for _, row := range rows {
			if uint32(row) >= uint32(1)<<fd.size() {
				log.Printf("流文：%s 字形位元超出尺寸", path)
				return nil
			}
		}
	}
	fd.byWord = map[string][]flowPos{}
	for pi, p := range fd.Paras {
		for a, w := range p.Words {
			fd.byWord[w] = append(fd.byWord[w], flowPos{int32(pi), int32(a)})
		}
	}
	log.Printf("流文：載入 %d 段落、%d 字形", len(fd.Paras), len(fd.Glyphs))
	return &fd
}

// stripCtl 直接刪除 <0x20 的字元（不補空白；bake_flow.py 的 strip_ctl 同義）。
func stripCtl(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x20 {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// flowResult 是一次比對的結果（R3）。
type flowResult struct {
	Zh     string
	Para   int   // 段落首個訊息 id
	A      int   // 頁在段落字詞序列的起點
	Lines  []int // 完成的行的訊息 id
	Cands  int
	Status string // 段落 id:a｜無譯文｜歧義｜缺字形
}

// match 依 R3 比對頁字詞序列 W。
func (fd *flowData) match(W []string) flowResult {
	res := flowResult{Status: "無譯文"}
	if fd == nil || len(W) == 0 {
		return res
	}
	n := len(W)
	type cand struct {
		zh    string
		para  int
		a     int
		lines []int
	}
	var cands []cand
	for _, pos := range fd.byWord[W[0]] {
		p := &fd.Paras[pos.P]
		a := int(pos.A)
		if a+n > len(p.Words) {
			continue
		}
		ok := true
		for k := 0; k < n; k++ {
			if p.Words[a+k] != W[k] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		c := cand{para: p.ID, a: a}
		empty := false
		for _, l := range p.Lines {
			if a <= l.End-1 && l.End-1 < a+n {
				c.lines = append(c.lines, l.ID)
				if l.Zh == "" {
					empty = true
				}
				c.zh += l.Zh
			}
		}
		if len(c.lines) == 0 || empty {
			c.zh = ""
		}
		cands = append(cands, c)
	}
	res.Cands = len(cands)
	if len(cands) == 0 {
		return res
	}
	same := true
	for _, c := range cands[1:] {
		if c.zh != cands[0].zh {
			same = false
		}
	}
	if !same {
		res.Status = "歧義"
		return res
	}
	c := cands[0]
	if c.zh == "" {
		return res
	}
	for _, r := range c.zh {
		if _, ok := fd.Glyphs[string(r)]; !ok {
			res.Status = "缺字形"
			return res
		}
	}
	res.Zh, res.Para, res.A, res.Lines = c.zh, c.para, c.a, c.lines
	res.Status = fmt.Sprintf("%d:%d", c.para, c.a)
	return res
}

// matchReads 先以當前頁時間限制來源情境，再以可見英文字詞比對（R3b）。
func (fd *flowData) matchReads(W []string, reads []flowRead, after, step uint64) flowResult {
	res := flowResult{Status: "無譯文"}
	if fd == nil || len(W) == 0 {
		return res
	}
	var ids []int
	for _, r := range reads {
		if after < r.Step && r.Step <= step {
			ids = append(ids, r.ID)
		}
	}
	type candidate struct {
		first int
		ids   []int
		zh    string
	}
	var cands []candidate
	seen := map[string]bool{}
	for s := range ids {
		var words []string
		var chain []int
		var completed []int
		zh := ""
		missing := false
		for _, id := range ids[s:] {
			m, ok := fd.Msgs[id]
			if !ok {
				break
			}
			chain = append(chain, id)
			words = append(words, strings.Fields(m.W)...)
			if len(words) <= len(W) {
				completed = append(completed, id)
				zh += m.Z
				missing = missing || m.Z == ""
			}
			if m.T == 1 {
				break
			}
		}
		if len(words) < len(W) {
			continue
		}
		ok := true
		for i, w := range W {
			if words[i] != w {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		key := fmt.Sprint(chain)
		if seen[key] {
			continue
		}
		seen[key] = true
		if missing || len(completed) == 0 {
			zh = ""
		}
		cands = append(cands, candidate{ids[s], completed, zh})
	}
	res.Cands = len(cands)
	if len(cands) == 0 {
		return res
	}
	for _, c := range cands[1:] {
		if c.zh != cands[0].zh {
			res.Status = "歧義"
			return res
		}
	}
	c := cands[0]
	if c.zh == "" {
		return res
	}
	for _, ch := range c.zh {
		if _, ok := fd.Glyphs[string(ch)]; !ok {
			res.Status = "缺字形"
			return res
		}
	}
	res.Zh, res.Para, res.Lines = c.zh, c.first, c.ids
	res.Status = fmt.Sprintf("取:%d", c.first)
	return res
}

func (ov *overlay) beginStory(ret oracle.Addr, sp uint16) {
	if ov.storyReturn != nil {
		ov.storyUnknown = true
		log.Printf("流文 邊界未知：前個故事返回未配對")
	}
	ov.storyReturn = &storyReturn{ret, sp}
}

func (ov *overlay) returnedStory(ret oracle.Addr, sp uint16, step uint64) {
	if r := ov.storyReturn; r != nil && r.ret == ret && sp == r.sp+2 {
		ov.storyLastReturn = step
		ov.storyReturn = nil
	}
}

// ── R4 排版 ───────────────────────────────────────────────────────────────

const (
	flowStartProhib = "，。、！？；：）」』》…—" // 行首禁用
	flowEndProhib   = "（「『《"          // 行尾禁用
)

type flowUnit struct {
	text string
	size int // 佔幾字
}

func flowUnits(zh string) []flowUnit {
	rs := []rune(zh)
	var out []flowUnit
	for i := 0; i < len(rs); i++ {
		if (rs[i] == '…' || rs[i] == '—') && i+1 < len(rs) && rs[i+1] == rs[i] {
			out = append(out, flowUnit{string(rs[i : i+2]), 2}) // ……、—— 不可拆
			i++
			continue
		}
		out = append(out, flowUnit{string(rs[i]), 1})
	}
	return out
}

// flowLayout 依 R4 把譯文排成列（每列至多 n 字）：串流貪婪填滿，避頭尾最多重複 3 次、列內至少留一個單位。
func flowLayout(zh string, n int) []string {
	if n < 1 {
		n = 1
	}
	u := flowUnits(zh)
	var rows []string
	i := 0
	for i < len(u) {
		j, w := i, 0
		for j < len(u) && w+u[j].size <= n {
			w += u[j].size
			j++
		}
		if j == i {
			j = i + 1
		}
		for t := 0; t < 3 && j-i > 1 && j < len(u); t++ {
			first := []rune(u[j].text)[0]
			last := []rune(u[j-1].text)[0]
			if !strings.ContainsRune(flowStartProhib, first) && !strings.ContainsRune(flowEndProhib, last) {
				break
			}
			j--
		}
		var sb strings.Builder
		for k := i; k < j; k++ {
			sb.WriteString(u[k].text)
		}
		rows = append(rows, sb.String())
		i = j
	}
	return rows
}

// flowGS 是旁白中文的字形放大（spec 006 R4，#41）：hdf≥2（Draw 倍率 4 以上）時等於 hdf；
// 否則倍率 3 以上放大 2 倍，倍率 2 維持 1 倍（列距只有 24 px，放不下 30 px 的字）。
func flowGS(k, hdf int) int {
	if hdf >= 2 {
		return hdf
	}
	if k >= 3 {
		return 2
	}
	return 1
}

// flowCols 回傳某個繪製端（倍率 k、字形放大 gs）每列字數上限 N。
func flowCols(k, gs int) int {
	n := flowWidth * k / (flowGlyphPitch * gs)
	if n < 1 {
		n = 1
	}
	return n
}

// flowRowPos 是某列在螢幕像素的位置（X0 置中於 160、Y0 依列）。
func flowRowPos(text string, r, k, gs int) (x0, y0 int) {
	n := len([]rune(text))
	pitch := flowGlyphPitch * gs
	return flowCenter*k - n*pitch/2, (flowY0 + flowRowStep*r) * k
}

// pageMatch 是頁事件的比對（spec 006 R5）：游標框之外，底圖像素（before＝after）與 after 完全相等，
// 且至少一半的筆畫像素與 before 不同（英文還在，色號不論：旁白筆畫會閃爍、淡入，F15）。
func (e *textEv) pageMatch(a []uint8, inCursor func(x, y int) bool) bool {
	stroke, present := 0, 0
	for y := 0; y < e.bh; y++ {
		for x := 0; x < e.bw; x++ {
			px, py := e.bx+x, e.by+y
			if inCursor(px, py) {
				continue
			}
			k := y*e.bw + x
			v := a[py*320+px]
			if e.before[k] == e.after[k] {
				if v != e.after[k] {
					return false
				}
				continue
			}
			stroke++
			if v != e.before[k] {
				present++
			}
		}
	}
	return stroke == 0 || present*2 >= stroke
}

// glyphOf 取字形（15 列，15 位元靠右）。
func (fd *flowData) glyphOf(r rune) []uint16 { return fd.Glyphs[string(r)] }

// ── 頁 ────────────────────────────────────────────────────────────────────

// flowPage 是引擎印出的一頁流動文字（R2）：至多 4 個片段。
type flowPage struct {
	site      oracle.Addr
	npc       bool // spec038：整個after精確發布，不採006部分墨點比對
	readAfter uint64
	words     []string
	frags     int       // 已進入的有效片段數
	lastY     int       // 上一片段的 y0
	evs       []*textEv // 還原事件（每個有效片段一個）
	lastStep  uint64    // 最後片段返回時的步數
	band      []uint8   // 最後片段返回時的文字帶離屏快照
	failed    bool
	res       flowResult
	rowsOK    bool // 以 flowK/flowHdf 排版後列數 ≤ 4
	counted   bool // 已計入 ov.shown
	bandLog   bool // 已 log 帶不符
}

var (
	flowK   = 3 // 排版用的倍率（compose 用來判斷列數；繪製端各自用自己的）
	flowHdf = 1
)

func (pg *flowPage) fail(why string, step uint64) {
	if !pg.failed {
		pg.failed = true
		log.Printf("流文 頁 %d 片段=%d 結果=無效(%s)", step, pg.frags, why)
	}
}

// retranslate 依目前的頁字詞序列重算譯文（每個片段返回後呼叫）。
func (pg *flowPage) retranslate(fd *flowData, reads []flowRead, step uint64) {
	if pg.site == storySite {
		pg.res = fd.matchReads(pg.words, reads, pg.readAfter, step)
	} else {
		pg.res = fd.match(pg.words)
	}
	pg.rowsOK = false
	if pg.res.Zh != "" {
		rows := flowLayout(pg.res.Zh, fd.cols(flowGeometry(pg.site), flowK, flowGS(flowK, flowHdf)))
		pg.rowsOK = len(rows) <= flowMaxRows
		if !pg.rowsOK {
			pg.res.Status = "超列"
		}
	}
	log.Printf("流文 頁 %d 片段=%d 字詞=%d 候選=%d 結果=%s", step, pg.frags, len(pg.words), pg.res.Cands, pg.res.Status)
}

// bandStatic 檢查 A000 的文字帶（游標框之外）與最後片段返回時的離屏帶快照，在頁所有事件的筆畫像素之外逐像素相等。
func (pg *flowPage) bandStatic(a []uint8, inCursor func(x, y int) bool) bool {
	site := flowGeometry(pg.site)
	if len(pg.band) != (site.x1-site.x0+1)*(flowBandY1-flowBandY0+1) {
		return false
	}
	w := site.x1 - site.x0 + 1
	skip := make([]bool, len(pg.band))
	for _, e := range pg.evs {
		for y := 0; y < e.bh; y++ {
			for x := 0; x < e.bw; x++ {
				k := y*e.bw + x
				bx, by := e.bx+x-site.x0, e.by+y-flowBandY0
				if e.before[k] != e.after[k] && bx >= 0 && bx < w && by >= 0 && by <= flowBandY1-flowBandY0 {
					skip[by*w+bx] = true
				}
			}
		}
	}
	k := 0
	for y := flowBandY0; y <= flowBandY1; y++ {
		for x := site.x0; x <= site.x1; x++ {
			if !skip[k] && !inCursor(x, y) && a[y*320+x] != pg.band[k] {
				return false
			}
			k++
		}
	}
	return true
}

// snapBand 取離屏的文字帶快照。
func snapBand(off []uint8, site flowSite) []uint8 {
	out := make([]uint8, 0, (site.x1-site.x0+1)*(flowBandY1-flowBandY0+1))
	for y := flowBandY0; y <= flowBandY1; y++ {
		out = append(out, off[y*320+site.x0:y*320+site.x1+1]...)
	}
	return out
}

// flowView 是交給繪製端的有效頁（快照節點產生，不可變）。
type flowView struct {
	Site oracle.Addr
	Zh   string
	Fg   [3]uint8
	Para int
}

func (pg *flowPage) view() flowView {
	fg := [3]uint8{}
	if len(pg.evs) > 0 {
		fg = pg.evs[0].fgRGB
	}
	return flowView{Site: pg.site, Zh: pg.res.Zh, Fg: fg, Para: pg.res.Para}
}

// ── 入口（R2）─────────────────────────────────────────────────────────────

// flowEntry 處理流動文字呼叫點的一次印字。回傳要記錄的待返回印字（nil＝不建事件）。呼叫端持有 evMu。
func (ov *overlay) flowEntry(o *oracle.Oracle, raw []byte, ret oracle.Addr, sp uint16, font uint16, x, y int) *pendingPrint {
	step := o.Steps()
	ov.fragSeq++
	seq := ov.fragSeq
	hasNUL := false
	n := len(raw)
	for i, b := range raw {
		if b == 0 {
			hasNUL, n = true, i
			break
		}
	}
	s := stripCtl(string(raw[:n]))
	log.Printf("流文 片段 %d ret=%04x:%04x font=%d y0=%d x0=%d %q", step, ret.Seg, ret.Off, font, y, x, s)

	// 內容防呆：64 bytes 內有 NUL、只含可印字元、至少一個字詞
	printable := true
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] > 0x7E {
			printable = false
		}
	}
	words := strings.Fields(s)
	if !hasNUL || !printable || len(words) == 0 {
		if ov.page != nil && !ov.page.failed {
			ov.page.fail("內容防呆", step)
		}
		return nil
	}

	// 頁：首片段 y0＝146；之後 y0＝上一片段＋12；至多 4 個片段
	pg := ov.page
	if pg == nil || pg.failed || pg.site != ret || y != pg.lastY+flowRowStep || pg.frags >= flowMaxFrags {
		pg = &flowPage{site: ret, readAfter: ov.storyLastReturn, lastY: y}
		ov.page = pg
		ov.pages = append(ov.pages, pg)
		if y != flowY0 {
			pg.fail("首列 y0 不是 146", step)
			return nil
		}
	}
	if ret == storySite && ov.storyUnknown {
		pg.fail("邊界未知", step)
		return nil
	}
	pg.words = append(pg.words, words...)
	pg.frags++
	pg.lastY = y
	return &pendingPrint{s: s, x0: x, y0: y, sp: sp, ret: ret, page: pg, font: font, seq: seq}
}

// ── 繪製端 ────────────────────────────────────────────────────────────────

// flowRowStrip 依資料原生尺寸把中文排成字條，字距2像素。
func (fd *flowData) flowRowStrip(text string, fg [3]uint8) *image.RGBA {
	rs := []rune(text)
	pitch := fd.size() + 2
	img := image.NewRGBA(image.Rect(0, 0, len(rs)*pitch, fd.size()))
	c := color.RGBA{fg[0], fg[1], fg[2], 255}
	for i, r := range rs {
		g := fd.glyphOf(r)
		for row, w := range g {
			for bit := 0; bit < fd.size(); bit++ {
				if (w>>(fd.size()-1-bit))&1 != 0 {
					img.Set(i*pitch+bit, row, c)
				}
			}
		}
	}
	return img
}

// drawFlow 在 Draw 疊上有效頁的中文列（字條放大 gs 倍，pitch＝17×gs）。
func (g *game) drawFlow(dst *ebiten.Image, views []flowView) {
	if g.flow == nil {
		return
	}
	hd := g.scale / 2
	if hd < 1 {
		hd = 1
	}
	if len(g.flowImgs) > 64 {
		g.flowImgs = map[string]*ebiten.Image{}
	}
	for _, v := range views {
		rows := flowLayout(v.Zh, g.flow.cols(flowGeometry(v.Site), g.scale, flowGS(g.scale, hd)))
		if len(rows) > flowMaxRows {
			continue
		}
		for r, text := range rows {
			key := fmt.Sprintf("%s|%v", text, v.Fg)
			img := g.flowImgs[key]
			if img == nil {
				img = ebiten.NewImageFromImage(g.flow.flowRowStrip(text, v.Fg))
				g.flowImgs[key] = img
			}
			gs := flowGS(g.scale, hd)
			x0, y0 := g.flow.rowPos(text, r, g.scale, gs)
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(gs), float64(gs))
			op.GeoM.Translate(float64(x0), float64(y0))
			op.Filter = ebiten.FilterNearest
			dst.DrawImage(img, &op)
		}
	}
}

// flowRowJSON 是 -hd-at 輸出的一列（R6）。
type flowRowJSON struct {
	X0, Y0 int
	Text   string
}

type flowJSON struct {
	Site      uint16
	Fragments []string
	Para, A   int
	N         int
	Lines     []int
	Zh        string
	Scale     int
	FgRGB     [3]uint8
	GS        int // 字形放大（R4）
	Rows      []flowRowJSON
}

// flowHD 在 hdFrame 的 img 上畫有效頁的中文列（字形像素放大 gs×gs、pitch 17×gs），並回傳 JSON 記錄。
func (fd *flowData) flowHD(img *image.RGBA, pages []*flowPage, k int) []flowJSON {
	var out []flowJSON
	gs := flowGS(k, 1)
	for _, pg := range pages {
		v := pg.view()
		rows := flowLayout(v.Zh, fd.cols(flowGeometry(v.Site), k, gs))
		if len(rows) > flowMaxRows {
			continue
		}
		j := flowJSON{Site: pg.site.Off, Para: pg.res.Para, A: pg.res.A, N: len(pg.words), Lines: pg.res.Lines, Zh: v.Zh, Scale: k, FgRGB: v.Fg, GS: gs}
		for _, e := range pg.evs {
			j.Fragments = append(j.Fragments, e.s)
		}
		c := color.RGBA{v.Fg[0], v.Fg[1], v.Fg[2], 255}
		for r, text := range rows {
			x0, y0 := fd.rowPos(text, r, k, gs)
			j.Rows = append(j.Rows, flowRowJSON{x0, y0, text})
			for i, ch := range []rune(text) {
				for row, w := range fd.glyphOf(ch) {
					for bit := 0; bit < fd.size(); bit++ {
						if (w>>(fd.size()-1-bit))&1 != 0 {
							for dy := 0; dy < gs; dy++ {
								for dx := 0; dx < gs; dx++ {
									img.Set(x0+(i*(fd.size()+2)+bit)*gs+dx, y0+row*gs+dy, c)
								}
							}
						}
					}
				}
			}
		}
		out = append(out, j)
	}
	return out
}
