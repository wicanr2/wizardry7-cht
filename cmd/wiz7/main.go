// wiz7 是《巫術7》的可遊玩前端（#15）：原版 DS.EXE 跑在 dosgolem 上，
// 開視窗、接鍵盤與滑鼠。
//
// 架構：機器在背景 goroutine 用大塊 Run 推進（probe 同款快路徑——
// 60fps 小塊 Run 在 credits→主選單轉場會卡死，docs/re/009），
// 前端執行緒只抓畫面快照與轉發輸入（鎖保護）。
// 預設按原版虛擬時鐘同步；明示 -ips 可以加速，音樂維持音高（spec 009）。
//
//	wiz7 -root workplace/run [-scale 3] [-ips 24000000]
package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
	rc "github.com/wicanr2/wizardry7-cht/recording"
)

// scanCode 把前端會用到的鍵映射成 Set 1 掃描碼（make code）。
func scanCode(k ebiten.Key) (uint8, bool) {
	m := map[ebiten.Key]uint8{
		ebiten.KeyEscape:    0x01,
		ebiten.KeyEnter:     0x1C,
		ebiten.KeySpace:     0x39,
		ebiten.KeyUp:        0x48,
		ebiten.KeyLeft:      0x4B,
		ebiten.KeyRight:     0x4D,
		ebiten.KeyDown:      0x50,
		ebiten.KeyF3:        0x3D,
		ebiten.KeyF4:        0x3E,
		ebiten.KeyF5:        0x3F,
		ebiten.KeyA:         0x1E,
		ebiten.KeyB:         0x30,
		ebiten.KeyC:         0x2E,
		ebiten.KeyD:         0x20,
		ebiten.KeyE:         0x12,
		ebiten.KeyF:         0x21,
		ebiten.KeyG:         0x22,
		ebiten.KeyH:         0x23,
		ebiten.KeyI:         0x17,
		ebiten.KeyJ:         0x24,
		ebiten.KeyK:         0x25,
		ebiten.KeyL:         0x26,
		ebiten.KeyM:         0x32,
		ebiten.KeyN:         0x31,
		ebiten.KeyO:         0x18,
		ebiten.KeyP:         0x19,
		ebiten.KeyQ:         0x10,
		ebiten.KeyR:         0x13,
		ebiten.KeyS:         0x1F,
		ebiten.KeyT:         0x14,
		ebiten.KeyU:         0x16,
		ebiten.KeyV:         0x2F,
		ebiten.KeyW:         0x11,
		ebiten.KeyX:         0x2D,
		ebiten.KeyY:         0x15,
		ebiten.KeyZ:         0x2C,
		ebiten.KeyBackspace: 0x0E,
	}
	sc, ok := m[k]
	return sc, ok
}

// normalizeKey＝查表 key 正規化（與 bake_font.py 的 normalize 同一套）：
// 剝 <0x20 控制碼、連續數字→'$'、序數詞 FIRST..TENTH→'$'、'^'/'%' 佔位→'$'。
// token（數字/序數詞）另存，供 '$' 槽填回。
func normalizeKey(s string) (string, []string) {
	var b []byte
	var nums []string
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 {
			continue
		}
		if c >= '0' && c <= '9' {
			j := i
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			b = append(b, '$')
			nums = append(nums, s[i:j])
			i = j - 1
			continue
		}
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			j := i
			for j < len(s) && (s[j] >= 'A' && s[j] <= 'Z' || s[j] >= 'a' && s[j] <= 'z') {
				j++
			}
			word := strings.ToUpper(s[i:j])
			if _, ok := ordinalsPub[word]; ok {
				b = append(b, '$')
				nums = append(nums, word)
				i = j - 1
				continue
			}
		}
		if c == '^' || c == '%' {
			b = append(b, '$')
			i++
			continue
		}
		b = append(b, c)
	}
	return string(b), nums
}

// ordinalsPub＝序數詞→中文（與 xlate __ordinals__ 對應；查詞用）。
var ordinalsPub = map[string]bool{
	"FIRST": true, "SECOND": true, "THIRD": true, "FOURTH": true, "FIFTH": true,
	"SIXTH": true, "SEVENTH": true, "EIGHTH": true, "NINTH": true, "TENTH": true,
}

// leadingWordDollar＝^＝隊員名 pattern 的 variant key 產生器：
// runtime 字串開頭若是「全大寫詞」（隊員名，如 THESUS），換成 '$'。
// 混有大小寫的普通字（如 "Life"）不觸發。
func leadingWordDollar(s string) string {
	i := 0
	for i < len(s) && s[i] < 0x20 {
		i++
	}
	j := i
	for j < len(s) && s[j] >= 'A' && s[j] <= 'Z' {
		j++
	}
	if j-i < 2 || j >= len(s) && j == i {
		return ""
	}
	return "$" + s[j:]
}

// ── 不寫原版的印字覆繪（docs/spec/001，證據 docs/re/014 E1–E9）──────────
//
// 原版每一則字串照常畫出英文（R1）。轉譯層在印字入口與返回各取一張離屏
// 快照，差異矩形就是英文筆畫所在（E4）；合成時以 A000 判定該處是否仍是
// 這串英文（E7），是的話把筆畫還原成印字前的背景，再在放大畫布上畫中文。

// evStrings 取 shown 事件的英文原文（畫面辨識用，spec 003 R4）。
func evStrings(evs []*textEv) []string {
	out := make([]string, 0, len(evs))
	for _, e := range evs {
		out = append(out, e.s)
	}
	return out
}

// offscreen 回傳 VGA.DRV 離屏緩衝（blit 來源段＝740C:01E5 的值，E5）。
// 段值為 0 時回 nil（尚未初始化）。
func offscreen(o *oracle.Oracle) []uint8 {
	seg := o.Word(oracle.Addr{Seg: 0x740C, Off: 0x01E5})
	if seg == 0 {
		return nil
	}
	return o.Bytes(oracle.Addr{Seg: seg, Off: 0}, 320*200)
}

// textEv 是一則有譯文的印字事件（R2.4）。
type textEv struct {
	originRet      oracle.Addr
	originFont     uint16
	popupClip      image.Rectangle
	s              string
	x0, y0, dy     int
	bx, by, bw, bh int     // 英文筆畫框（原版座標）
	before, after  []uint8 // 框內離屏色號：印字前／後
	fg             uint8   // 筆畫色號（差異像素眾數）
	fgRGB          [3]uint8
	g              [][]uint16
	shown          bool               // 曾與 A000 相符（決定被刪除或等待，與頁有沒有譯文無關）
	match          bool               // 本次 compose 是否與 A000 相符（頁裁決用）
	dead           bool               // 已被刪除
	page           *flowPage          // 流動文字的還原事件所屬頁（spec 006）；nil＝一般事件
	img            *ebiten.Image      // Draw 時 lazily 預渲染的字形條
	compound       *compoundTrace     // spec012：完整兩行物品的原子事件與唯讀診斷
	roster         *rosterWitness     // spec022：碼帶事件的完整列身分
	sex            bool               // spec024：直接單字元性別事件
	profile        *profilePage       // spec024 R966：2倍人物整頁原子事件
	party          *partyClassWitness // spec036：側欄短碼及原姓名末行身份
}

// pendingPrint 是已進入、尚未返回的印字（R2.1）。
type pendingPrint struct {
	s      string
	x0, y0 int
	dy     int
	sp     uint16
	ret    oracle.Addr
	before []uint8
	g      [][]uint16
	page   *flowPage // 流動文字片段所屬頁（spec 006）；nil＝一般印字
	font   uint16
	seq    int // 流動文字片段序號
	roster *rosterRow
	party  *pendingPartyClass
}

// overlay 保存事件與統計；呼叫端持有 evMu。
type overlay struct {
	whoEntry          entry
	reads             []flowRead
	storyReturn       *storyReturn
	storyLastReturn   uint64
	storyUnknown      bool
	pend              *pendingPrint
	compound          *compoundPair                   // spec012：至多一件待配對物品
	flow              *flowData                       // flow.json（nil＝流動文字關閉）
	page              *flowPage                       // 目前的頁（流動文字呼叫點）
	pages             []*flowPage                     // 尚未被丟棄的頁
	fragSeq           int                             // 流動文字片段的序號（1 起，與 log「流文 片段」同序）
	onFrag            func(o *oracle.Oracle, seq int) // 測試用：片段返回後的回呼（-hd-after-frag）
	evs               []*textEv
	created           int
	shown             int
	dropped           int
	empty             int
	nested            int
	noseg             int
	roster            map[string]rosterEntry
	rosterEnabled     bool
	sexEntries        map[string]rosterEntry
	sexEnabled        bool
	sexPending        *pendingProfileSex
	profileEnabled    bool
	profile           *profileReflowState
	picture           *introPictureState
	brush             *brushMenuData
	clearance         *clearanceHint
	partyEnabled      bool
	partyScale        int
	partySize         int
	partyEntry        entry
	partyNames        [6]partyNameState
	partyPendingNames [6]*pendingPartyName
}

// cursorBox 是 A000 上游標可能佔的框（E8 量到 10×10，未量形狀先放寬到 16×16）。
const cursorBox = 16

// finish 在返回 hook 取 after 快照、算差異框並建立事件（R2.3–R2.4）。
func (ov *overlay) finish(o *oracle.Oracle) {
	p := ov.pend
	ov.pend = nil
	after := offscreen(o)
	if p.roster != nil {
		if e := finishRoster(p, after, o.Palette()); e != nil {
			ov.evs = append(ov.evs, e)
			ov.created++
		} else {
			log.Printf("職業碼缺譯：印字差異不符 y=%d", p.y0)
		}
		return
	}
	if p.party != nil {
		if e := finishPartyClass(p, after, o.Palette(), o.Steps()); e != nil {
			ov.evs = append(ov.evs, e)
			ov.created++
		} else {
			log.Printf("側欄職業缺譯：印字差異不符 y=%d", p.y0)
		}
		return
	}
	if after == nil || p.before == nil {
		ov.noseg++
		if p.page != nil {
			p.page.fail("無離屏", o.Steps())
		}
		return
	}
	x0, y0, x1, y1 := 320, 200, -1, -1
	for i := range after {
		if after[i] != p.before[i] {
			x, y := i%320, i/320
			x0, y0 = min(x0, x), min(y0, y)
			x1, y1 = max(x1, x), max(y1, y)
		}
	}
	if x1 < 0 {
		ov.empty++
		if p.page != nil {
			p.page.fail("無差異", o.Steps())
		}
		return
	}
	w, h := x1-x0+1, y1-y0+1
	e := &textEv{originRet: p.ret, originFont: p.font, s: p.s, x0: p.x0, y0: p.y0, dy: p.dy, bx: x0, by: y0, bw: w, bh: h, g: p.g, page: p.page}
	e.before = make([]uint8, w*h)
	e.after = make([]uint8, w*h)
	var cnt [256]int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y0+y)*320 + x0 + x
			k := y*w + x
			e.before[k], e.after[k] = p.before[i], after[i]
			if after[i] != p.before[i] {
				cnt[after[i]]++
			}
		}
	}
	for c := 1; c < 256; c++ {
		if cnt[c] > cnt[e.fg] {
			e.fg = uint8(c)
		}
	}
	e.fgRGB = o.Palette()[e.fg] // 返回 hook 當下的調色盤（R2.4）
	if pg := p.page; pg != nil {
		// spec 006 R2 片段閘門：字型 2、筆畫色 0x55、離屏差異框中心 x 在 159–161
		if pg.failed {
			return
		}
		site := flowSites[pg.site]
		if p.font != flowFont || e.fg != flowFg || 2*x0+w-1 < site.centerMin || 2*x0+w-1 > site.centerMax {
			pg.fail(fmt.Sprintf("閘門 font=%d fg=%02x 中心2x=%d", p.font, e.fg, 2*x0+w-1), o.Steps())
			return
		}
		pg.evs = append(pg.evs, e)
		pg.lastStep = o.Steps()
		pg.band = snapBand(after, site)
		pg.retranslate(ov.flow, ov.reads, o.Steps())
	}
	ov.evs = append(ov.evs, e)
	ov.created++
}

// compose 依 R3 判定事件狀態，回傳合成後的低解析色號、shown 事件與有效的流動文字頁（R4.1）。
// 兩趟（spec 006 R5）：第一趟判定每個事件的 match／pending／dropped 並處理一般事件；
// 第二趟依頁裁決，有效頁才還原筆畫、輸出譯文。可被 hdFrame 等在快照節點之外重複呼叫，
// 因此計數與 log 都用頁旗標，不靠呼叫次數。
func (ov *overlay) compose(o *oracle.Oracle) ([]uint8, []*textEv, []*flowPage) {
	a := o.Indexed()
	comp := append([]uint8(nil), a...)
	off := offscreen(o)
	mx, my := o.Mouse()
	step := o.Steps()
	inCursor := func(x, y int) bool {
		return x >= mx && x < mx+cursorBox && y >= my && y < my+cursorBox
	}
	popupActive := profilePopupActive(ov.evs, a, inCursor)
	restore := func(e *textEv) {
		if e.profile != nil {
			e.profile.restore(comp, inCursor)
			return
		}
		for y := 0; y < e.bh; y++ {
			for x := 0; x < e.bw; x++ {
				k := y*e.bw + x
				i := (e.by+y)*320 + e.bx + x
				if e.popupContains(e.bx+x, e.by+y) {
					continue
				}
				if e.before[k] == e.after[k] {
					continue
				}
				if e.page != nil || e.compound != nil {
					// 頁／完整物品：原子還原所有筆畫，跳過游標框。
					if !inCursor(e.bx+x, e.by+y) {
						comp[i] = e.before[k]
					}
				} else if a[i] == e.after[k] {
					comp[i] = e.before[k]
				}
			}
		}
	}
	// 第一趟
	kept := ov.evs[:0]
	var shown []*textEv
	for _, e := range ov.evs {
		if e.dead {
			ov.dropped++
			continue
		}
		e.popupClip = image.Rectangle{}
		if popupActive && profilePopupEquipment(e) {
			e.popupClip = profilePopupRect
		}
		match := e.rectMatches(a, func(x, y int) bool { return inCursor(x, y) || e.popupContains(x, y) })
		if e.profile != nil {
			match = e.profile.matches(a, inCursor)
			e.profile.view.pal = o.Palette()
		}
		if e.roster != nil {
			match = e.rosterMatches(a, inCursor)
		}
		if e.party != nil {
			match = match && e.partyNameMatches(a, inCursor)
		}
		if e.page != nil {
			match = e.pageMatch(a, inCursor)
		}
		e.match = match
		switch {
		case match:
			if !e.shown {
				e.shown = true
				if e.compound != nil {
					e.compound.State = "shown"
					log.Printf("複合物品 上屏 %s step=%d", e.s, step)
				}
				if e.page == nil {
					ov.shown++
				}
			}
			if e.page == nil {
				restore(e)
				shown = append(shown, e)
			}
			kept = append(kept, e)
		case e.shown:
			ov.dropped++
			e.dead = true
			if e.compound != nil {
				e.compound.State = "drop"
				log.Printf("複合物品 捨棄 %s 原因=上屏域不符 step=%d", e.s, step)
			}
		default:
			// 尚未上屏：離屏框仍是 after 才保留（R3 表第 3、4 列）
			still := e.rectMatches(off, nil)
			if e.profile != nil {
				still = e.profile.matches(off, nil)
			}
			if e.roster != nil {
				still = e.rosterMatches(off, nil)
			}
			if e.party != nil {
				still = still && e.partyNameMatches(off, nil)
			}
			if still {
				kept = append(kept, e)
			} else {
				ov.dropped++
				e.dead = true
				if e.compound != nil {
					e.compound.State = "drop"
					log.Printf("複合物品 捨棄 %s 原因=離屏域不符 step=%d", e.s, step)
				}
			}
		}
	}
	for i := len(kept); i < len(ov.evs); i++ {
		ov.evs[i] = nil
	}
	ov.evs = kept

	// 第二趟：頁裁決（R5）
	var active []*flowPage
	keptPages := ov.pages[:0]
	for _, pg := range ov.pages {
		dead := pg.failed
		for _, e := range pg.evs {
			if e.dead {
				dead = true
			}
		}
		if dead {
			continue // 丟棄頁；仍存活的事件是無效事件，隨 A000 改變被刪除
		}
		keptPages = append(keptPages, pg)
		if pg.res.Zh == "" || !pg.rowsOK || pg.frags == 0 || len(pg.evs) != pg.frags || step-pg.lastStep < flowQuiet {
			continue
		}
		all := true
		for _, e := range pg.evs {
			if !e.match {
				all = false
			}
		}
		if !all {
			continue
		}
		if !pg.bandStatic(a, inCursor) {
			if !pg.bandLog {
				pg.bandLog = true
				log.Printf("流文 帶不符 %d 段落=%d", step, pg.res.Para)
			}
			continue
		}
		for _, e := range pg.evs {
			restore(e)
			shown = append(shown, e)
		}
		if !pg.counted {
			pg.counted = true
			ov.shown++
		}
		active = append(active, pg)
	}
	for i := len(keptPages); i < len(ov.pages); i++ {
		ov.pages[i] = nil
	}
	ov.pages = keptPages
	return comp, shown, active
}

func (ov *overlay) stats() string {
	return fmt.Sprintf("事件 建立=%d 上屏=%d 刪除=%d 無差異=%d 巢狀捨棄=%d 無離屏=%d 存活=%d",
		ov.created, ov.shown, ov.dropped, ov.empty, ov.nested, ov.noseg, len(ov.evs))
}

func fnv(b []byte) uint32 {
	var h uint32 = 2166136261
	for _, v := range b {
		h = (h ^ uint32(v)) * 16777619
	}
	return h
}

// hdScale 是 hdFrame 的畫布倍率（與互動視窗 -scale 同義）。行距 6 的畫面
// （防拷三行）在 2× 只有 12 px，放不下 15 px 的中文字；3× 為 18 px。
var hdScale = 3

// hdFrame 在快照節點合成高解析幀：低解析合成（已還原英文筆畫）放大 hdScale 倍，
// 中文 15×15 原字形畫在 (x0×k, (y0+dy)×k)。呼叫端持有 evMu。
func hdFrame(o *oracle.Oracle, ov *overlay, path string) (comp []uint8, shown []*textEv, flows []flowJSON) {
	comp, shown, pages := ov.compose(o)
	pal := o.Palette()
	k := hdScale
	who := whoCursorFromSource(o, ov.whoEntry, k, o.Indexed())
	who.restore(comp)
	img := image.NewRGBA(image.Rect(0, 0, 320*k, 200*k))
	for py := 0; py < 200*k; py++ {
		for px := 0; px < 320*k; px++ {
			c := pal[comp[(py/k)*320+(px/k)]]
			img.Set(px, py, color.RGBA{c[0], c[1], c[2], 255})
		}
	}
	ov.picture.view(o.Indexed(), pal, o.VideoMode()).paintHD(img, k)
	ov.picture.paperView(o.Indexed(), pal, o.VideoMode(), k).paintHD(img, k)
	brush := ov.brush.view(shown, comp, pal)
	for _, e := range shown {
		if brush.suppress(e, k) {
			continue
		}
		if e.profile != nil {
			e.profile.view.paint(img)
			continue
		}
		fr := color.RGBA{e.fgRGB[0], e.fgRGB[1], e.fgRGB[2], 255}
		shift := worldSexShift(e, k)
		dx, dy := partyClassShift(e, k)
		cx := e.x0*k - shift + dx
		cy := buttonLabelY(e, shown, k) - shift + dy
		for _, glyph := range e.g {
			for row, w := range glyph {
				for bit := 0; bit < len(glyph); bit++ {
					if (w>>(len(glyph)-1-bit))&1 != 0 {
						if !e.popupContainsHD(cx+bit, cy+row, k) {
							img.Set(cx+bit, cy+row, fr)
						}
					}
				}
			}
			cx += len(glyph) + 2
		}
	}
	if ov.flow != nil {
		flows = ov.flow.flowHD(img, pages, k)
	}
	brush.paintHD(img, k)
	ov.clearance.view(shown, k).paint(img)
	who.paint(img, k)
	if f, err := os.Create(path); err == nil {
		png.Encode(f, img)
		f.Close()
	}
	return comp, shown, flows
}

// hdAtDump 輸出 -hd-at 驗收資料（spec 001 R5 與驗收）：HD 合成 PNG、
// A000 原始色號、低解析合成色號、shown 事件與雜湊。呼叫端持有 evMu。
func hdAtDump(o *oracle.Oracle, ov *overlay, am *automap, dir string, name string) {
	base := fmt.Sprintf("%s/%s", dir, name)
	a := append([]uint8(nil), o.Indexed()...)
	comp, shown, flows := hdFrame(o, ov, base+".png")
	pal := o.Palette()
	palB := make([]byte, 0, 768)
	for _, c := range pal {
		palB = append(palB, c[0], c[1], c[2])
	}
	mem := fnv(o.Bytes(oracle.Addr{Seg: 0, Off: 0}, 0x100000))
	type evOut struct {
		S              string
		X0, Y0, Dy     int
		BX, BY, BW, BH int
		Before, After  []uint8
		Fg             uint8
		G              [][]uint16       // 字形（spec 002 R5 的序數檢查用）
		Compound       *compoundTrace   `json:",omitempty"`
		Profile        *profileMetadata `json:",omitempty"`
	}
	out := struct {
		Step           uint64
		Mem, A000, Pal uint32
		MemorySHA256   string
		Regs           any
		Stats          string
		Mouse          [2]int
		Shown          []evOut
		Map            mapJSON
		Flow           []flowJSON
	}{Step: o.Steps(), Mem: mem, A000: fnv(a), Pal: fnv(palB), Stats: ov.stats()}
	out.MemorySHA256 = fmt.Sprintf("%x", sha256.Sum256(o.Bytes(oracle.Addr{Seg: 0, Off: 0}, 0x100000)))
	out.Regs = o.Regs()
	out.Mouse[0], out.Mouse[1] = o.Mouse()
	out.Map = am.json()
	out.Flow = flows
	for _, e := range shown {
		out.Shown = append(out.Shown, evOut{e.s, e.x0, e.y0, e.dy, e.bx, e.by, e.bw, e.bh, e.before, e.after, e.fg, e.g, e.compound, e.profile.metadata()})
	}
	_ = os.WriteFile(base+".a000", a, 0o644)
	_ = os.WriteFile(base+".comp", comp, 0o644)
	_ = os.WriteFile(base+".pal", palB, 0o644) // spec007：獨立像素oracle的唯讀色盤收據
	if b, err := json.Marshal(out); err == nil {
		_ = os.WriteFile(base+".json", b, 0o644)
	}
	log.Printf("hdat %s steps=%d mem=%08x a000=%08x pal=%08x shown=%d %s",
		name, out.Step, out.Mem, out.A000, out.Pal, len(shown), out.Stats)
}

type entry struct {
	Dy int        `json:"dy"`
	G  [][]uint16 `json:"g"`
	D  []int      `json:"d,omitempty"` // '$' 數字槽索引（g 內位置，位圖全 0）
}

type shared struct {
	who                  *whoCursorView
	recordingPublication recordingPublication
	mu                   sync.Mutex
	rgba                 []byte  // 320×200×4 的畫面快照（Draw 用）
	keys                 []uint8 // 待注入的掃描碼（make＋break 成對排程）
	mouseX               int
	mouseY               int
	rbtn                 bool      // 右鍵狀態（左鍵走 clicks，spec 004）
	evs                  []*textEv // shown 事件（快照節點判定，Draw 只消費，R4.3）
	drawReq              string    // 非空＝下一次 Draw 完成後把視窗畫面存到這個路徑（-draw-at）
	helpOpen             bool      // F1 說明開啟中＝機器暫停、輸入不送原版（spec 003 R2）
	helpScreen           string    // 開啟當下辨識的畫面（R4）
	helpPage             int
	mapOpen              bool       // F2 自動地圖面板（spec 005 R3）
	mapView              mapView    // 快照節點產生的地圖狀態
	flow                 []flowView // 有效的流動文字頁（快照節點產生；Draw 消費，spec 006）
	picture              *introPictureView
	paper                *introPaperView
	brush                *brushMenuView
	clearance            *clearanceHintView
	stdin                []byte               // 待注入的 Stdin 字元（AH=3Fh 路徑）
	clicks               []struct{ x, y int } // 待執行的整包點擊（Click 路徑）
}

type game struct {
	whoImage          *ebiten.Image
	whoImageKey       string
	live              *inputRecording
	seenRecordingDraw uint64
	pictureGlyph      *image.RGBA
	pictureImage      *ebiten.Image
	paperView         *introPaperView
	paperImages       []*ebiten.Image
	brushView         *brushMenuView
	brushImages       [2][]*ebiten.Image
	clearanceView     *clearanceHintView
	clearanceImage    *ebiten.Image
	flow              *flowData                // flow.json（nil＝流動文字關閉）
	flowImgs          map[string]*ebiten.Image // 中文列字條快取（key＝列文字＋顏色）
	help              *helpData                // F1 說明資料（help.json；nil＝缺檔）
	helpImg           *ebiten.Image            // 目前說明頁的預渲染
	helpKey           string
	sh                *shared
	scale             int
	screen            *ebiten.Image
	lastMX            int
	lastMY            int
	lbtn              bool
	rbtn              bool
	stop              *stopRequest
	draws             atomic.Uint64
}

func (g *game) Update() error {
	if g.stop != nil && ebiten.IsWindowBeingClosed() {
		if g.live != nil {
			g.live.Cancel()
		}
		g.stop.Request()
	}
	if g.stop != nil && g.stop.Requested() {
		return ebiten.Termination
	}
	if g.live != nil {
		g.live.Poll()
		if g.live.replay {
			return nil
		}
	}
	footer := false
	if g.live != nil {
		footer = g.live.mouseInFooter(g.scale)
	}
	// ── F1 說明（spec 003 R1）：F1 永遠由前端處理；說明開啟時其他輸入一律不送原版 ──
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.sh.mu.Lock()
		if g.sh.helpOpen {
			g.sh.helpOpen = false
		} else {
			g.sh.helpOpen, g.sh.helpPage = true, 0
			g.sh.helpScreen = detectScreen(evStrings(g.sh.evs))
		}
		log.Printf("說明 %v 畫面=%s", g.sh.helpOpen, g.sh.helpScreen)
		g.sh.mu.Unlock()
	}
	g.sh.mu.Lock()
	open := g.sh.helpOpen
	if open {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
			g.sh.helpOpen = false
			log.Printf("說明 false（Esc）")
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), inpututil.IsKeyJustPressed(ebiten.KeyPageUp):
			g.sh.helpPage--
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight), inpututil.IsKeyJustPressed(ebiten.KeyPageDown):
			g.sh.helpPage++
		}
	}
	g.sh.mu.Unlock()
	if open {
		ebiten.AppendInputChars(nil) // 丟棄本幀字元，不排入 stdin
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) { // spec 005 R3：前端處理，不送原版
		g.sh.mu.Lock()
		g.sh.mapOpen = !g.sh.mapOpen
		log.Printf("地圖 面板=%v", g.sh.mapOpen)
		g.sh.mu.Unlock()
	}
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if sc, ok := scanCode(k); ok {
			g.sh.mu.Lock()
			g.sh.keys = append(g.sh.keys, sc, sc|0x80)
			g.sh.mu.Unlock()
		}
	}
	// 可列印字元 → Stdin（AH=3Fh 路徑：旁白推進/防拷/名字）
	for _, r := range ebiten.AppendInputChars(nil) {
		if r >= 32 && r < 127 {
			g.sh.mu.Lock()
			g.sh.stdin = append(g.sh.stdin, byte(r))
			g.sh.mu.Unlock()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.sh.mu.Lock()
		g.sh.stdin = append(g.sh.stdin, '\r')
		g.sh.mu.Unlock()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.sh.mu.Lock()
		g.sh.stdin = append(g.sh.stdin, '\b')
		g.sh.mu.Unlock()
	}
	if footer {
		return nil
	}
	mx, my := ebiten.CursorPosition()
	sx, sy := mx/g.scale, my/g.scale
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.lbtn = true
		// spec 004 R1：左鍵剛按下＝排一次整包點擊（遊戲座標，裁切到 320×200）
		cx, cy := min(max(sx, 0), 319), min(max(sy, 0), 199)
		g.sh.mu.Lock()
		g.sh.clicks = append(g.sh.clicks, struct{ x, y int }{cx, cy})
		g.sh.mu.Unlock()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.lbtn = false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.rbtn = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.rbtn = false
	}
	g.sh.mu.Lock()
	if sx != g.sh.mouseX || sy != g.sh.mouseY {
		g.sh.mouseX, g.sh.mouseY = sx, sy
	}
	g.sh.rbtn = g.rbtn
	g.sh.mu.Unlock()
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	g.draws.Add(1)
	g.sh.mu.Lock()
	rgba := append([]byte(nil), g.sh.rgba...)
	evs := append([]*textEv(nil), g.sh.evs...)
	flows := append([]flowView(nil), g.sh.flow...)
	picture := g.sh.picture
	paper, brush := g.sh.paper, g.sh.brush
	clearance := g.sh.clearance
	who := g.sh.who
	publication := g.sh.recordingPublication
	g.sh.mu.Unlock()
	if len(rgba) == 0 {
		return
	}
	gameArea := dst
	if g.live != nil && !g.live.replay {
		area := image.Rect(0, 0, 320*g.scale, 200*g.scale)
		if !area.In(dst.Bounds()) {
			return
		}
		gameArea = dst.SubImage(area).(*ebiten.Image)
	}
	if who.canDraw(gameArea.Bounds(), g.scale) {
		who.patchRGBA(rgba)
	}
	g.screen.WritePixels(rgba)
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(float64(g.scale), float64(g.scale))
	op.Filter = ebiten.FilterNearest
	dst.DrawImage(g.screen, &op)
	g.drawIntroPicture(gameArea, picture)
	g.drawIntroPaper(gameArea, paper)

	// ── 高解析中文合成 ──
	// 語意（AGENTS.md 高解析中文合成）：原版 320×200 終幀放大後，
	// 中文字以原字形尺寸畫在對應位置——字相對畫面變小、剛好放進按鈕。
	// hd＝max(1,scale/2)：2×與3×都保留後端15／16px原生字形。
	hd := g.scale / 2
	if hd < 1 {
		hd = 1
	}
	for _, e := range evs {
		if brush.suppress(e, g.scale) {
			continue
		}
		if e.profile != nil {
			if e.img == nil {
				img := image.NewRGBA(image.Rect(0, 0, 640, 400))
				e.profile.view.paint(img)
				e.img = ebiten.NewImageFromImage(img)
			}
			dst.DrawImage(e.img, nil)
			continue
		}
		if len(e.g) == 0 || len(e.g[0]) == 0 {
			continue // 空字形條（純控制碼/純 ASCII 譯文）跳過
		}
		if e.img == nil {
			// 字形條：整串排成一張（每字 15px＋2px 間距）
			n := 0
			for _, gl := range e.g {
				n += len(gl) + 2
			}
			rect := image.Rect(0, 0, n, len(e.g[0]))
			img := image.NewRGBA(rect)
			cx := 0
			for _, glyph := range e.g {
				for row, w := range glyph {
					for bit := 0; bit < len(glyph); bit++ {
						if (w>>(len(glyph)-1-bit))&1 != 0 {
							img.Set(cx+bit, row, color.RGBA{e.fgRGB[0], e.fgRGB[1], e.fgRGB[2], 255})
						}
					}
				}
				cx += len(glyph) + 2
			}
			e.img = ebiten.NewImageFromImage(img)
		}
		op2 := ebiten.DrawImageOptions{}
		op2.GeoM.Scale(float64(hd), float64(hd))
		// 平移在縮放之後，單位是視窗像素：原版座標 × 倍率
		shift := worldSexShift(e, g.scale)
		dx, dy := partyClassShift(e, g.scale)
		op2.GeoM.Translate(float64(e.x0*g.scale-shift+dx), float64(buttonLabelY(e, evs, g.scale)-shift+dy))
		op2.Filter = ebiten.FilterNearest
		e.drawPopupClipped(dst, e.img, &op2, g.scale)
	}
	g.drawFlow(dst, flows)
	g.drawBrushMenu(gameArea, brush)
	g.drawClearanceHint(gameArea, clearance)
	g.drawWhoCursor(gameArea, who)

	g.sh.mu.Lock()
	hOpen, hScreen, hPage := g.sh.helpOpen, g.sh.helpScreen, g.sh.helpPage
	mOpen, mView := g.sh.mapOpen, g.sh.mapView
	g.sh.mu.Unlock()
	if mOpen {
		g.drawMap(gameArea, mView)
	}
	if hOpen {
		g.drawHelp(gameArea, hScreen, hPage)
	}

	// -draw-at：把互動視窗實際畫面存檔（驗證 Draw 路徑，spec 001 R4.3）
	g.sh.mu.Lock()
	req := g.sh.drawReq
	g.sh.drawReq = ""
	g.sh.mu.Unlock()
	if req != "" {
		w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		dst.ReadPixels(img.Pix)
		if f, err := os.Create(req); err == nil {
			png.Encode(f, img)
			f.Close()
			log.Printf("drawat 輸出 %s（%dx%d）", req, w, h)
		}
	}
	if g.live != nil {
		g.live.drawControls(dst, g.scale)
		if publication.ID != 0 && publication.ID != g.seenRecordingDraw {
			g.seenRecordingDraw = publication.ID
			g.live.acknowledge(publication, rc.Panel{HelpOpen: hOpen, HelpScreen: hScreen, HelpPage: int64(hPage), MapOpen: mOpen})
		}
	}
}

func (g *game) Layout(int, int) (int, int) {
	height := 200 * g.scale
	if g.live != nil && !g.live.replay {
		height += 40
	}
	return 320 * g.scale, height
}

func main() {
	root := flag.String("root", "workplace/run", "原版素材目錄（唯讀）")
	exe := flag.String("exe", "", "執行檔（預設 <root>/DS.EXE）")
	saveDir := flag.String("save-dir", "", "原版存檔目錄（預設系統使用者資料目錄）")
	scale := flag.Int("scale", 3, "畫面放大倍率")
	ips := flag.Float64("ips", 3004074.5544433594, "每秒指令數（未指定時採精確原版虛擬時鐘；指定正值可加速）")
	mute := flag.Bool("mute", false, "靜音（原版音訊事件與遊戲時間仍正常推進）")
	xlatePath := flag.String("xlate", "workplace/bake/xlate.json", "譯文表（空字串＝純英文模式）")
	hdShotF := flag.Bool("hdshot", false, "每 30 秒輸出 HD 合成幀（PNG）")
	autoplay := flag.Bool("autoplay", false, "內建自動遊玩時序（建隊→防拷→旁白→迷宮移動）")
	hdOnly := flag.Bool("hdonly", false, "頭less：不開視窗，只輸出 HD 幀（無 WM 環境驗證用）")
	hdAt := flag.String("hd-at", "", "步數逗號分隔：到該步數輸出 HD 合成幀與驗收資料 <hd-dir>/hdat-<步數>.*（spec 001 R5）")
	hdDir := flag.String("hd-dir", "workplace/bake", "-hd-at 的輸出目錄")
	flag.IntVar(&hdScale, "hd-scale", 3, "HD 幀（-hd-at、-hdshot、自動遊玩截圖）的畫布倍率（使用者定案 2026-09-30：預設 3）")
	helpPath := flag.String("help", "workplace/bake/help.json", "F1 說明資料（tools/re/bake_help.py 產生）")
	helpAt := flag.Uint64("help-at", 0, "測試用：到該步數開啟 F1 說明（spec 003 R6）")
	helpFor := flag.Int("help-for", 3000, "測試用：-help-at 開啟後幾毫秒（牆鐘）關閉")
	helpPageF := flag.Int("help-page", 0, "測試用：-help-at 開啟時停在第幾頁（0 起算）")
	clicksAt := flag.String("clicks-at", "", "測試用：step:x:y 逗號分隔，到該步數以 o.Click 點擊（遊戲座標；spec 005 R4）")
	mapAt := flag.Uint64("map-at", 0, "測試用：到該步數開啟自動地圖面板")
	flowPath := flag.String("flow", "workplace/bake/flow.json", "流動文字資料（tools/re/bake_flow.py 產生；空字串停用，spec 006）")
	hdAfterFrag := flag.String("hd-after-frag", "", "測試用：流動文字片段序號逗號分隔（1 起）；該片段返回的當下輸出 <hd-dir>/hdfrag-<序號>.*（spec 006 驗收 10）")
	autoplayStop := flag.Uint64("autoplay-stop", 0, "測試用：autoplay 只執行門檻 ≤ 此步數的動作（0＝全部），之後交給真實輸入")
	drawAt := flag.String("draw-at", "", "步數逗號分隔：到該步數後把互動視窗畫面存到 <hd-dir>/drawat-<步數>.png")
	recordPath := flag.String("record", "", "從啟動錄製輸入，按頁尾停止後繼續遊戲")
	replayPath := flag.String("replay", "", "重播相同版本、平台、資料及初始存檔的完整錄放檔")
	flag.Parse()
	recordingMode := *recordPath != "" || *replayPath != ""
	recordingFatal := func(code string) { log.Print("錄放啟動失敗：" + code); os.Exit(2) }
	if recordingMode {
		log.SetFlags(0)
		log.SetOutput(recordingPrivacyLog{})
		explicitHDScale := false
		flag.Visit(func(f *flag.Flag) {
			if f.Name == "hd-scale" {
				explicitHDScale = true
			}
		})
		var scaleErr error
		hdScale, scaleErr = recordingEffectiveHDScale(*scale, hdScale, explicitHDScale, true)
		if scaleErr != nil {
			recordingFatal("configuration_mismatch")
		}
		if *recordPath != "" && *replayPath != "" {
			recordingFatal("configuration_mismatch")
		}
		forbidden := map[string]bool{"autoplay": true, "autoplay-stop": true, "clicks-at": true, "help-at": true, "map-at": true, "hdonly": true, "hdshot": true, "hd-at": true, "hd-after-frag": true, "draw-at": true}
		flag.Visit(func(f *flag.Flag) {
			if forbidden[f.Name] {
				recordingFatal("configuration_mismatch")
			}
		})
	}
	explicitIPS := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "ips" {
			explicitIPS = true
		}
	})
	tempo, err := audioTempo(*ips, explicitIPS)
	if err != nil {
		if recordingMode {
			recordingFatal("configuration_mismatch")
		}
		log.Fatal("-ips 必須是有限正值")
	}

	program := *exe
	if program == "" {
		program = *root + "/DS.EXE"
	}
	resolvedSaveDir, err := resolveSavePath(*saveDir)
	if err != nil {
		if recordingMode {
			recordingFatal("initialization_failed")
		}
		log.Fatal("存檔設定錯誤：default_save_path")
	}
	preparedSaveDir, err := prepareSaveStorage(*root, program, resolvedSaveDir, nil)
	if err != nil {
		if recordingMode {
			recordingFatal("initialization_failed")
		}
		log.Fatalf("存檔設定錯誤：%s", saveErrorCode(err))
	}

	var live *inputRecording
	var bundle recordingBundle
	var hd *helpData
	var o *oracle.Oracle
	if recordingMode {
		binary, e := os.Executable()
		if e != nil {
			recordingFatal("initialization_failed")
		}
		protected := []string{program, binary, *xlatePath, *flowPath, *helpPath}
		if *xlatePath != "" {
			for _, name := range []string{"roster.json", "profile-sex.json", "intro-picture.json"} {
				protected = append(protected, filepath.Join(filepath.Dir(*xlatePath), name))
			}
		}
		location := *recordPath
		if location == "" {
			location = *replayPath
		}
		if ic.CheckLocation(location, *root, preparedSaveDir, protected, *recordPath != "") != nil {
			recordingFatal("configuration_mismatch")
		}
		var document *rc.Document
		if *replayPath != "" {
			decoded, e := rc.Load(*replayPath)
			if e != nil {
				recordingFatal("wire_rejected")
			}
			document = &decoded
		}
		bundle, err = loadRecordingBundle(*xlatePath, *flowPath, *helpPath)
		if err != nil {
			recordingFatal("configuration_mismatch")
		}
		var controlFont *helpData
		if *recordPath != "" {
			controlFont, err = recordingControlFont(bundle.flow, *scale)
			if err != nil {
				recordingFatal("configuration_mismatch")
			}
		}
		var expected *ic.Config
		if document != nil {
			expected = &document.Header.Config
		}
		owner, e := newRecordingOwner(recordingStartup{root: *root, program: program, scratch: preparedSaveDir, scale: int64(*scale), ips: *ips, explicitIPS: explicitIPS, mute: *mute, resources: bundle.resources.metadata}, expected)
		if e != nil {
			code := "configuration_mismatch"
			if e.Error() == "initialization_failed" {
				code = "initialization_failed"
			}
			recordingFatal(code)
		}
		live, e = newInputRecording(owner, *recordPath, document, controlFont)
		if e != nil {
			owner.Close()
			recordingFatal("recording_output_failed")
		}
		o, e = owner.Take()
		if e != nil {
			owner.Close()
			recordingFatal("initialization_failed")
		}
		hd = bundle.help
	} else {
		o, err = oracle.LoadWith(program, *root, oracle.Options{Adlib: true})
		if err != nil {
			log.Fatal(err)
		}
		o.SetScratch(preparedSaveDir)
	}
	log.Printf("原版存檔目錄：%s", preparedSaveDir)

	windowHeight := 200 * *scale
	if live != nil && !live.replay {
		windowHeight += 40
	}
	ebiten.SetWindowSize(320**scale, windowHeight)
	ebiten.SetWindowTitle("巫術7 繁體中文版 — dosgolem")
	ebiten.SetWindowClosingHandled(true)
	stop := newStopRequest()
	if live != nil {
		live.stop = stop
	}

	sh := &shared{rgba: make([]byte, 320*200*4)}
	// ── xlate 翻譯層 ──────────────────────────────────────────────
	// 譯文表（原文 → 每字 15×15 位圖，tools/re/bake_font.py 烘製）。
	trans := map[string]entry{}
	digits := map[string][]uint16{}
	ordinals := map[string][]uint16{}
	var brushData *brushMenuData
	var paperData *introPaperData
	if recordingMode {
		trans, digits, ordinals = bundle.trans, bundle.digits, bundle.ordinals
		brushData, paperData = bundle.brush, bundle.paper
	} else {
		if *xlatePath == "" {
			*xlatePath = "/nonexistent.json"
		}
		if b, err := os.ReadFile(*xlatePath); err == nil {
			decoded, err := decodeXlate(b)
			if err != nil {
				log.Fatal("xlate.json：", err)
			}
			trans, digits, ordinals = decoded.trans, decoded.digits, decoded.ordinals
			brushData, paperData = decoded.brush, decoded.paper
			log.Printf("xlate：載入 %d 串、數字 %d、序數 %d", len(trans), len(digits), len(ordinals))
		} else {
			log.Printf("xlate：無譯文表（%v）——純英文模式", xlatePath)
		}

	}
	var evMu sync.Mutex
	treeSeen := map[uint16]bool{}
	lastHD := time.Now()
	hdShot := *hdShotF
	_ = hdShot
	if !recordingMode {
		os.MkdirAll("workplace/bake/trees", 0o755)
	}
	// 印字串入口（0110:3895）：只讀，不寫原版（spec 001 R1、R2）。
	ov := &overlay{}
	ov.whoEntry = takeWhoCursorEntry(trans)
	ov.brush = brushData
	ov.partyEnabled, ov.partySize = partyClassConfig(trans, *scale, hdScale)
	ov.partyScale = *scale
	ov.partyEntry = takePartyClassEntry(trans)
	if recordingMode {
		ov.flow, ov.picture = bundle.flow, bundle.picture
	} else {
		ov.flow = loadFlow(*flowPath)
		ov.picture = loadIntroPicture(*xlatePath, trans)
	}
	ov.picture.paper = paperData
	if *scale != 2 && *scale != 3 {
		log.Print("片頭圖片缺譯：互動倍率未支援，保留原圖")
	}
	if hdScale != 2 && hdScale != 3 {
		log.Print("片頭圖片缺譯：HD倍率未支援，保留原圖")
	}
	if recordingMode {
		ov.roster = bundle.roster
	} else {
		ov.roster = loadRoster(*xlatePath, trans)
	}
	ov.rosterEnabled = rosterResourcesValid(ov.roster) && *scale == 3 && hdScale == 3
	if recordingMode {
		ov.sexEntries = bundle.sex
	} else {
		ov.sexEntries = loadProfileSex(*xlatePath, trans)
	}
	ov.sexEnabled = len(ov.sexEntries) == 2 && *scale == 3 && hdScale == 3
	ov.profileEnabled = len(trans) > 0 && *scale == 2 && hdScale == 2
	// 故事讀取端只讀記錄：id 限定來源情境，仍須比對英文（spec 006 R2/R3b）。
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x0582}, func(o *oracle.Oracle) {
		evMu.Lock()
		defer evMu.Unlock()
		ov.clearance.begin(o)
		if o.StackWord(0) != flowSites[storySite].reader {
			return
		}
		r := flowRead{o.Steps(), int(o.StackWord(1))}
		ov.reads = append(ov.reads, r)
		if len(ov.reads) > 32 {
			ov.reads = ov.reads[len(ov.reads)-32:]
		}
		log.Printf("流文 取 %d id=%d", r.Step, r.ID)
	})
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x0659}, func(o *oracle.Oracle) {
		evMu.Lock()
		defer evMu.Unlock()
		ov.clearance.finish(o)
	})
	fragDump := map[int]bool{}
	for _, f := range strings.Split(*hdAfterFrag, ",") {
		if f == "" {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			log.Fatalf("-hd-after-frag %q：%v", f, err)
		}
		fragDump[n] = true
	}
	// 排版倍率：互動視窗用 -scale（Draw 的字條放大 hd＝scale÷2），無頭用 -hd-scale（1:1）
	if *hdOnly {
		flowK, flowHdf = hdScale, 1
	} else {
		flowK, flowHdf = *scale, max(1, *scale/2)
	}
	am := newAutomap() // spec 005：自動地圖（只在快照節點更新）
	if len(fragDump) > 0 {
		ov.onFrag = func(o *oracle.Oracle, seq int) { // 呼叫端（返回 hook）持有 evMu
			if fragDump[seq] {
				hdAtDump(o, ov, am, *hdDir, fmt.Sprintf("hdfrag-%d", seq))
			}
		}
	}
	entry := oracle.Addr{Seg: 0x0110, Off: 0x3895}
	retHooked := map[uint32]bool{}
	// arm 在返回位址掛一次 hook：只認本次返回（IP＝返回位址且 SP＝進入 SP＋2，spec 001 R2.2）
	arm := func(ret oracle.Addr) {
		lin := ret.Linear()
		if retHooked[lin] {
			return
		}
		retHooked[lin] = true
		o.OnCall(ret, func(o *oracle.Oracle) {
			evMu.Lock()
			defer evMu.Unlock()
			ov.picture.returned(ret, o.Regs().SP, o.Indexed())
			ov.returnedStory(ret, o.Regs().SP, o.Steps())
			ov.returnPartyNames(ret, o.Regs().SP, o.Steps())
			ov.returnedProfileReflow(ret, o.Regs().SP,
				int(o.Word(oracle.Addr{Seg: 0x740c, Off: 0x0295})),
				int(o.Word(oracle.Addr{Seg: 0x740c, Off: 0x0297})), offscreen(o))
			if ov.compound != nil && ov.compound.phase != 1 {
				if ov.returnCompound(ret, o.Regs().SP, o.Steps(), offscreen(o), o.Palette()) {
					return
				}
			}
			if ov.pend == nil || ov.pend.ret != ret || o.Regs().SP != ov.pend.sp+2 {
				return
			}
			seq := ov.pend.seq
			ov.finish(o)
			if seq > 0 && ov.onFrag != nil {
				ov.onFrag(o, seq)
			}
		})
	}
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3921}, func(o *oracle.Oracle) {
		evMu.Lock()
		defer evMu.Unlock()
		r := o.Regs()
		ret := oracle.Addr{Seg: r.CS, Off: o.StackWord(0)}
		if ov.picture.enter(ret, r.SP, [4]uint16{o.StackWord(1), o.StackWord(2), o.StackWord(3), o.StackWord(4)}) {
			arm(ret)
		}
	})
	o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3877}, func(o *oracle.Oracle) {
		r := o.Regs()
		ret := oracle.Addr{Seg: r.CS, Off: o.StackWord(0)}
		if ret == worldSexReturn {
			if r.SP > 65530 {
				evMu.Lock()
				ov.clearWorldSex()
				evMu.Unlock()
				log.Print("世界性別缺譯：stack span越界")
				return
			}
			x, y := int(o.Word(oracle.Far(0x740c, 0x0295))), int(o.Word(oracle.Far(0x740c, 0x0297)))
			if x != 90 || y != 6 {
				return
			}
			evMu.Lock()
			defer evMu.Unlock()
			ov.enterProfileSex(ret, r.SP, o.StackWord(2), o.StackWord(1), x, y, offscreen(o))
			return
		}
		if ret != profileSexReturn {
			return
		}
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
		if x != 90 || y != 6 {
			return
		}
		evMu.Lock()
		defer evMu.Unlock()
		if ov.profileEnabled {
			ov.enterProfileReflowSex(ret, r.SP, o.StackWord(2), o.StackWord(1), x, y, offscreen(o))
		} else {
			ov.enterProfileSex(ret, r.SP, o.StackWord(2), o.StackWord(1), x, y, offscreen(o))
		}
	})
	o.OnCall(worldSexReturn, func(o *oracle.Oracle) {
		evMu.Lock()
		defer evMu.Unlock()
		if ov.sexPending != nil && ov.sexPending.ret == worldSexReturn {
			ov.returnProfileSex(o.Regs().SP, int(o.Word(oracle.Far(0x740c, 0x0295))), int(o.Word(oracle.Far(0x740c, 0x0297))), offscreen(o), o.Palette())
		}
	})
	o.OnCall(profileSexReturn, func(o *oracle.Oracle) {
		evMu.Lock()
		defer evMu.Unlock()
		if ov.profileEnabled {
			ov.returnedProfileReflow(profileSexReturn, o.Regs().SP,
				int(o.Word(oracle.Addr{Seg: 0x740c, Off: 0x0295})),
				int(o.Word(oracle.Addr{Seg: 0x740c, Off: 0x0297})), offscreen(o))
		}
		if ov.sexPending != nil && ov.sexPending.ret == profileSexReturn {
			ov.returnProfileSex(o.Regs().SP,
				int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295})),
				int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})), offscreen(o), o.Palette())
		}
	})
	o.OnCall(entry, func(o *oracle.Oracle) {
		evMu.Lock()
		ov.clearProfileSex("一般字串插入")
		evMu.Unlock()
		r0 := o.Regs()
		if r0.SP > 65530 {
			evMu.Lock()
			ov.retireRosterAt(int(o.Word(oracle.Far(0x740c, 0x0297))))
			evMu.Unlock()
			log.Print("職業碼缺譯：入口stack span")
			return
		}
		ret := oracle.Addr{Seg: r0.CS, Off: o.StackWord(0)} // near call（E2）
		off := o.StackWord(1)
		ds := r0.DS
		// spec012 C2：所有入口在flow/generic早退前先仲裁完整物品配對。
		// C1：配對專用讀取最多8byte，連NUL亦不跨16-bit offset邊界。
		partRaw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, min(8, 65536-int(off)))
		part := string(partRaw)
		if end := strings.IndexByte(part, 0); end >= 0 {
			part = part[:end]
		}
		in := compoundInput{S: part, DS: ds, Off: off, SP: r0.SP, Font: o.StackWord(2),
			Ret: ret, Step: o.Steps(), X: int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295})),
			Y: int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297})), raw: partRaw}
		evMu.Lock()
		handled := ov.enterCompound(in, func() []byte { return offscreen(o) }, trans, flowK)
		evMu.Unlock()
		if handled {
			arm(ret)
			return
		}
		evMu.Lock()
		profileHandled := false
		if ov.profileEnabled && ov.profile != nil {
			raw := partRaw
			if ret.Off != 0x032a {
				raw = o.Bytes(oracle.Addr{Seg: ds, Off: off}, min(64, 65536-int(off)))
			}
			profileHandled = ov.enterProfileReflowText(ret, r0.SP, o.StackWord(2), in.X, in.Y, raw, trans, func() []byte { return offscreen(o) })
		}
		evMu.Unlock()
		if profileHandled {
			arm(ret)
			return
		}
		if _, flow := flowSites[ret]; flow {
			raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64)
			if ret == storySite {
				arm(ret) // 返回通知先於內容防呆；沒有覆繪事件也要保留邊界。
			}
			// spec 006 R2：流動文字呼叫點不查表，整頁另行比對
			x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
			y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
			evMu.Lock()
			if ret == storySite {
				ov.beginStory(ret, r0.SP)
			}
			if ov.pend != nil {
				ov.nested++ // 防禦性：E2 實測未見巢狀
				if ov.pend.page != nil {
					ov.pend.page.fail("巢狀", o.Steps())
				}
			}
			p := ov.flowEntry(o, raw, ret, r0.SP, o.StackWord(2), x, y)
			if p != nil {
				if b := offscreen(o); b != nil {
					p.before = append([]byte(nil), b...)
				}
			}
			ov.pend = p
			evMu.Unlock()
			if p != nil {
				arm(ret)
			}
			return
		}
		if ov.rosterEnabled && ret == oracle.Far(0x0110, 0x0fcc) {
			evMu.Lock()
			ov.retireRosterAt(in.Y)
			evMu.Unlock()
			var evidence *rosterContextEvidence
			var raw, before []byte
			regs := rosterRegisters{r0.CS, r0.DS, r0.SS, r0.SP, r0.BP}
			if in.Font == 0 && in.X == 150 && rosterKnownY(in.Y) && off <= 65520 && rosterFrameBounds(regs) {
				raw = o.Bytes(oracle.Far(ds, off), 16)
				evidence, _ = readRosterInvocationContext(regs, ret, off, in.Font, in.X, in.Y, raw, o.Bytes)
			}
			if evidence != nil {
				if _, err := validateRosterContext(evidence, off, in.Font, in.X, in.Y); err == nil {
					before = offscreen(o)
				}
			}
			evMu.Lock()
			handled := ov.enterRoster(ret, r0.SP, in.Font, off, in.X, in.Y, evidence, raw, before)
			evMu.Unlock()
			if handled {
				if ov.pend != nil {
					arm(ret)
				}
				return
			}
		}
		raw := o.Bytes(oracle.Addr{Seg: ds, Off: off}, 64) // 一般路由
		evMu.Lock()
		partyName := ov.enterPartyName(ret, r0.SP, in.Font, off, in.X, in.Y, raw, o.Steps(), o.Bytes)
		evMu.Unlock()
		if partyName {
			arm(ret)
			return
		}
		s := string(raw)
		for i := 0; i < len(s); i++ {
			if s[i] == 0 {
				s = s[:i]
				break
			}
		}
		if s == "ALCHEMI" {
			evMu.Lock()
			var before []byte
			if ov.partyEnabled && ret == oracle.Far(0x0110, 0x1c71) && in.Font == 0 && r0.SP <= 65530 && off <= 65504 &&
				(in.X == 9 || in.X == 261) && (in.Y == 56 || in.Y == 101 || in.Y == 146) {
				before = offscreen(o)
			}
			handled := ov.enterPartyClass(ret, r0.SP, in.Font, off, in.X, in.Y, raw, o.Steps(), before, ov.partyEntry)
			pending := ov.pend != nil && ov.pend.party != nil
			evMu.Unlock()
			if handled {
				if pending {
					arm(ret)
				}
				return
			}
		}
		// 查表 key 正規化（與 bake_font.py 的 normalize 同一套）：
		// 剝 <0x20 控制碼、連續數字/序數詞 FIRST..TENTH → '$'（token 另存）。
		key, nums := normalizeKey(s)
		ent, ok := trans[key]
		if !ok {
			// ^＝隊員名 pattern（2040/2061/5010/1541 等）：runtime 開頭是
			// 隊員名（全大寫）——首詞換成 '$' 再查一次（表側 zh 不含名字）。
			if v := leadingWordDollar(s); v != "" {
				k2, n2 := normalizeKey(v)
				if e2, ok2 := trans[k2]; ok2 {
					ent, ok, nums = e2, true, n2
				}
			}
		}
		if !ok {
			log.Printf("無譯文 %d %q @(%d,%d) ret=%04x:%04x", o.Steps(), s,
				o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}), o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}), ret.Seg, ret.Off)
			return // 無譯文：放行原文
		}
		// spec015：所需槽位字模缺項時整則保留英文，不建立部分中文事件。
		ent, ok = expandSlots(ent, nums, digits, ordinals)
		if !ok {
			log.Printf("填槽缺字 %d @(%d,%d) ret=%04x:%04x", o.Steps(),
				o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}), o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}), ret.Seg, ret.Off)
			return
		}
		x := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0295}))
		y := int(o.Word(oracle.Addr{Seg: 0x740C, Off: 0x0297}))
		p := &pendingPrint{
			s: s, x0: x, y0: y, dy: ent.Dy, g: ent.G,
			sp:   r0.SP,
			ret:  ret,
			font: in.Font,
		}
		if b := offscreen(o); b != nil {
			p.before = append([]byte(nil), b...)
		}
		evMu.Lock()
		if ov.pend != nil {
			// 防禦性：E2 實測未見巢狀（R2.5）
			ov.nested++
			log.Printf("巢狀印字：捨棄未返回的 %q", ov.pend.s)
			if ov.pend.page != nil {
				ov.pend.page.fail("巢狀", o.Steps())
			}
		}
		ov.pend = p
		evMu.Unlock()
		arm(p.ret)
		// 動態樹捕捉：樹段變化時 dump（msgtree 模式）
		ts := o.Word(oracle.Addr{Seg: 0x0110, Off: 0x2D4E})
		if ts != 0 && !recordingMode {
			lastTree, have := treeSeen[ts]
			_ = lastTree
			if !have {
				treeSeen[ts] = true
				data := o.Bytes(oracle.Addr{Seg: ts, Off: 0}, 8192)
				fn := fmt.Sprintf("workplace/bake/trees/tree-%04x.bin", ts)
				if err := os.WriteFile(fn, data, 0o644); err == nil {
					log.Printf("樹段 %04x → %s", ts, fn)
				}
			}
		}
		log.Printf("翻譯 %d %q @(%d,%d)", o.Steps(), s, x, y)
	})
	if !recordingMode {
		hd = loadHelp(*helpPath)
	}
	ov.clearance = newClearanceHint(trans, hd)
	g := &game{sh: sh, scale: *scale, help: hd, flow: ov.flow, flowImgs: map[string]*ebiten.Image{}, stop: stop, live: live}
	if !*hdOnly {
		g.screen = ebiten.NewImage(320, 200)
	}

	// 自動遊玩步數門控表（oracle 層直接注入，繞過 Ebiten 輸入在 Xvfb 的失效）。
	// R10 改用 steps 觸發：wall-clock 在 R8（主選單 +45s）vs R9（約 +60s）間漂移，
	// steps 有 M1 決定性保證。錨點來自 R9 實測（docs/re/010 R9 追記）。
	var autos []struct {
		at uint64 // 機器步數門檻
		fn func(*oracle.Oracle)
	}
	if *autoplay {
		add := func(s uint64, fn func(*oracle.Oracle)) {
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{s, fn})
		}
		// R14：游標 park 觀測。每次動作後記 pos／sets／range。
		logMouse := func(o *oracle.Oracle, op string, x, y int, err error) {
			mx, my := o.Mouse()
			sets := o.MouseSets()
			last := "none"
			if len(sets) > 0 {
				l := sets[len(sets)-1]
				last = fmt.Sprintf("(%d,%d)@%d", l.X, l.Y, l.Step)
			}
			n0, n1, n2, n3 := o.MouseRange()
			log.Printf("AUTO滑鼠 %s (%d,%d) err=%s pos=(%d,%d) sets=%d last=%s range=[%d-%d,%d-%d]",
				op, x, y, frontendOperationCode(err), mx, my, len(sets), last, n0, n1, n2, n3)
		}
		cl := func(x, y int) func(*oracle.Oracle) {
			// R13：印出 Click error（含 Polls）——畫面不變時區分座標錯／沒準備好。
			// R14：附游標觀測（park 假說）。
			return func(o *oracle.Oracle) {
				err := o.Click(x, y)
				if err != nil {
					log.Printf("AUTO點擊 (%d,%d)：%s", x, y, frontendOperationCode(err))
				}
				logMouse(o, "click", x, y, err)
			}
		}
		domouse := func(op string, x, y int) func(*oracle.Oracle) {
			return func(o *oracle.Oracle) {
				var err error
				switch op {
				case "press":
					o.MoveMouse(x, y)
					o.PressButton(0)
				case "release":
					o.ReleaseButton(0)
				case "right":
					err = o.Click(x, y, oracle.Right())
				case "racepress":
					// R15：Move 後立即零 hover 按住 20M，賭第一次 park 前的 poll。
					o.MoveMouse(x, y)
					err = o.Press(oracle.Hold(20_000_000))
				case "parkpress":
					// R15 對照組：不 Move，原地按（讀到的必是 park 點，預期不動）。
					err = o.Press(oracle.Hold(20_000_000))
				case "log":
					// 純觀測，不動輸入。
				case "move":
					// R18：只搬移（按住另由 press/release 管），觀測 pos 釘住否。
					o.MoveMouse(x, y)
				default:
					err = o.Click(x, y)
				}
				logMouse(o, op, x, y, err)
			}
		}
		// START NEW GAME 路線（R9 拓樸修正：建隊對話框後直達七項選單，不繞 PCVW）：
		// 主選單 246M→建隊→THESUS 254M→七項選單 260M→START→巫6匯入詢問 CANCEL→
		// 防拷→旁白 40 次→迷宮八點。門檻＝錨點＋餘裕（約 4–5M steps）。
		add(250_000_000, cl(160, 62)) // 主選單「新增角色」直達建隊
		add(253_000_000, cl(170, 153))
		add(255_000_000, cl(90, 168))
		add(257_000_000, cl(90, 168))
		add(265_000_000, cl(160, 101))                                // 七項選單: START NEW GAME（y-13 熱點系統）
		add(276_000_000, cl(82, 169))                                 // 巫6 匯入詢問: CANCEL（語意待驗）
		add(288_000_000, func(o *oracle.Oracle) { o.Type("AAAA\r") }) // 防拷
		for i := 0; i < 133; i++ {                                    // 旁白推進（R11：引言延續到 1225M＋，每 8M 一次）
			add(295_000_000+uint64(i)*8_000_000, cl(25, 15))
		}
		// 迷宮精確箭頭（R20）：先點底文清場，再打實測幾何的箭頭。
		// driver 座標（x 已加倍；羅盤 game px 見 010 R19 追記）。
		type mvAct struct {
			s    uint64
			x, y int
			op   string
		}
		dirS := []mvAct{
			{1_430_000_000, 320, 170, "click"}, // 底文（清場，R19 confirmed 有效）
			{1_434_000_000, 320, 150, "click"}, // 底文上部（保險）
			{1_438_000_000, 310, 118, "click"}, // 北（上箭頭）
			{1_442_000_000, 230, 140, "click"}, // 西（左箭頭）
			{1_446_000_000, 390, 140, "click"}, // 東（右箭頭）
			{1_450_000_000, 310, 158, "click"}, // 南（下箭頭）
			{1_454_000_000, 310, 138, "click"}, // 中央
			{1_458_000_000, 320, 100, "click"}, // 中央視野
			{1_462_000_000, 310, 118, "click"}, // 北再一次
			{1_466_000_000, 230, 140, "click"}, // 西再一次
		}
		for _, d := range dirS {
			add(d.s, domouse(d.op, d.x, d.y))
		}
		// 每點後截圖（移動實證）
		for _, d := range dirS {
			shotS := d.s + 9_000_000
			add(shotS, func(o *oracle.Oracle) {
				evMu.Lock()
				hdFrame(o, ov, fmt.Sprintf("workplace/bake/mv-%d-%d-%d.png", d.s, d.x, d.y))
				evMu.Unlock()
			})
		}
	}
	if *autoplay && *autoplayStop > 0 {
		kept := autos[:0]
		for _, a := range autos {
			if a.at <= *autoplayStop {
				kept = append(kept, a)
			}
		}
		autos = kept
		log.Printf("autoplay 停在 %d 步之後（autoplay 剩 %d 個動作）", *autoplayStop, len(autos))
	}
	// -hd-at：依步數輸出 HD 合成幀（決定性，取代 wall-clock 的 -hdshot 做驗收）。
	if *hdAt != "" {
		for _, f := range strings.Split(*hdAt, ",") {
			st, err := strconv.ParseUint(strings.TrimSpace(f), 10, 64)
			if err != nil {
				log.Fatalf("-hd-at %q：%v", f, err)
			}
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{st, func(o *oracle.Oracle) {
				evMu.Lock()
				hdAtDump(o, ov, am, *hdDir, fmt.Sprintf("hdat-%d", st))
				evMu.Unlock()
			}})
		}
	}
	for _, f := range strings.Split(*clicksAt, ",") {
		var st uint64
		var x, y int
		if f == "" {
			continue
		}
		if _, err := fmt.Sscanf(f, "%d:%d:%d", &st, &x, &y); err != nil {
			log.Fatalf("-clicks-at %q：%v", f, err)
		}
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{st, func(o *oracle.Oracle) {
			err := o.Click(x, y)
			log.Printf("測試點擊 (%d,%d) steps=%d err=%s", x, y, o.Steps(), frontendOperationCode(err))
		}})
	}
	if *mapAt > 0 {
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{*mapAt, func(o *oracle.Oracle) {
			sh.mu.Lock()
			sh.mapOpen = true
			sh.mu.Unlock()
			log.Printf("地圖 面板=true（-map-at）")
		}})
	}
	if *helpAt > 0 {
		autos = append(autos, struct {
			at uint64
			fn func(*oracle.Oracle)
		}{*helpAt, func(o *oracle.Oracle) {
			sh.mu.Lock()
			sh.helpOpen, sh.helpPage = true, *helpPageF
			sh.helpScreen = detectScreen(evStrings(sh.evs))
			log.Printf("說明 true（-help-at）畫面=%s shown=%q", sh.helpScreen, evStrings(sh.evs))
			sh.mu.Unlock()
			time.AfterFunc(time.Duration(*helpFor)*time.Millisecond, func() {
				sh.mu.Lock()
				sh.helpOpen = false
				sh.mu.Unlock()
				log.Printf("說明 false（-help-for）")
			})
		}})
	}
	if *drawAt != "" {
		for _, f := range strings.Split(*drawAt, ",") {
			st, err := strconv.ParseUint(strings.TrimSpace(f), 10, 64)
			if err != nil {
				log.Fatalf("-draw-at %q：%v", f, err)
			}
			autos = append(autos, struct {
				at uint64
				fn func(*oracle.Oracle)
			}{st, func(o *oracle.Oracle) {
				sh.mu.Lock()
				sh.drawReq = fmt.Sprintf("%s/drawat-%d.png", *hdDir, st)
				sh.mu.Unlock()
			}})
		}
	}
	sort.SliceStable(autos, func(i, j int) bool { return autos[i].at < autos[j].at })
	autoIdx := 0

	// snapshot 是快照節點（spec 001 R4.3）：R3 判定一次，低解析合成（已還原英文筆畫）
	// 先在owner建立獨立資料，再短鎖發布；原版推進／合成／節流均不持sh.mu。
	// 每輪迴圈結尾與每個 auto 執行後都要呼叫：autoplay 的 o.Click 內部會跑數百萬步，
	// 多個 auto 可能在同一輪連鎖執行，只在輪末快照會讓 sh.evs 過期（R667 實測：
	// 240M→286M 之間零快照，-help-at 讀到 247M 的舊主選單）。
	snapshot := func() {
		am.update(o) // spec 005 R1／R2：只讀
		mapSnapshot := am.view()
		evMu.Lock()
		comp, shown, pages := ov.compose(o)
		picture := ov.picture.view(o.Indexed(), o.Palette(), o.VideoMode())
		paper := ov.picture.paperView(o.Indexed(), o.Palette(), o.VideoMode(), *scale)
		brush := ov.brush.view(shown, comp, o.Palette())
		clearance := ov.clearance.view(shown, *scale)
		who := whoCursorFromSource(o, ov.whoEntry, *scale, o.Indexed())
		evMu.Unlock()
		var views []flowView
		for _, pg := range pages {
			views = append(views, pg.view())
		}
		rgba := make([]byte, 320*200*4)
		pal := o.Palette()
		for i, v := range comp {
			c := pal[v]
			rgba[i*4], rgba[i*4+1], rgba[i*4+2], rgba[i*4+3] = c[0], c[1], c[2], 255
		}
		visual := cloneVisualEvents(shown)
		picture.patchRGBA(rgba, *scale)
		paper.patchRGBA(rgba, *scale)
		sh.mu.Lock()
		sh.mapView, sh.flow, sh.rgba, sh.evs = mapSnapshot, views, rgba, visual
		sh.picture = picture
		sh.paper, sh.brush = paper, brush
		sh.clearance = clearance
		sh.who = who
		var publication recordingPublication
		published := false
		if live != nil {
			publication, published = live.publish(o.Steps(), sh)
		}
		sh.mu.Unlock()
		if published {
			live.afterPublish(publication)
		}
	}

	audioOut, err := newPlayerAudio(o, tempo, *mute || *hdOnly, stop)
	if err != nil {
		if live != nil {
			live.Cancel()
			live.abortOwned()
			o.Close()
			recordingFatal("initialization_failed")
		}
		log.Fatal("音訊串流初始化失敗")
	}
	if live != nil {
		live.snapshot = snapshot
		live.boundaryEligible = func() bool { return live.qualified(audioOut, o) }
	}
	ownerDone := make(chan struct{})
	var sourceStats oracle.OPLSourceStats
	var finalStep uint64
	runStarted := time.Now()

	// 機器 goroutine：大塊 Run（probe 同款快路徑），節流＋吃前端輸入。
	// ⚠ MoveMouse 只在座標真的變了才送——遊戲把滑鼠活動當「有人在看」，
	// 每塊都送會重置故事頁的推進計時（docs/re/009 的教訓）。
	go func() {
		defer close(ownerDone)
		if live != nil {
			defer live.abortOwned()
		}
		defer stop.Request()
		defer func() { sourceStats = o.OPLSourceStats(); finalStep = o.Steps() }()
		chunk := uint64(200_000)
		lastR := false
		lastX, lastY := -1, -1
		paused := false
		for i := 0; ; i++ {
			if stop.Requested() {
				audioOut.stopAtBoundary()
				return
			}
			if live != nil {
				if live.replay {
					live.ReplayNext(o, audioOut, snapshot)
					continue
				}
				live.Poll()
				if live.state.Load() == 3 {
					live.abortOwned()
				}
				if live.Requested() && live.state.Load() == 0 {
					if !live.qualified(audioOut, o) {
						live.fail("recording_ineligible")
						live.abortOwned()
					} else {
						snapshot()
					}
				}
			}
			sh.mu.Lock()
			// spec 003 R2：說明開啟＝不推進、不投遞任何輸入（含游標位置）
			pauseChanged := sh.helpOpen != paused
			paused = sh.helpOpen
			if paused {
				sh.mu.Unlock()
				audioOut.Pause(true)
				if pauseChanged {
					log.Printf("說明 暫停=true steps=%d", o.Steps())
				}
				timer := time.NewTimer(16 * time.Millisecond)
				select {
				case <-timer.C:
				case <-stop.done:
				}
				timer.Stop()
				continue
			}
			keys := append([]uint8(nil), sh.keys...)
			sh.keys = nil
			stdin := append([]byte(nil), sh.stdin...)
			sh.stdin = nil
			bx, by := sh.mouseX, sh.mouseY
			r := sh.rbtn
			// 整包點擊（Click 路徑——分離的 Press/Release 在部分遊戲失效）
			var clicks []struct{ x, y int }
			for len(sh.clicks) > 0 {
				c := sh.clicks[0]
				sh.clicks = sh.clicks[1:]
				clicks = append(clicks, c)
			}
			sh.mu.Unlock()
			audioOut.Pause(false)
			if pauseChanged {
				log.Printf("說明 暫停=false steps=%d", o.Steps())
			}
			for _, sc := range keys {
				if live != nil {
					live.Operation(o, audioOut, "press_key", recordingKey(sc), func() error { o.PressKey(oracle.Key(sc)); return nil })
				} else {
					o.PressKey(oracle.Key(sc))
				}
				log.Printf("送鍵 %02x steps=%d", sc, o.Steps())
			}
			if len(stdin) > 0 {
				if live != nil {
					live.Operation(o, audioOut, "type", recordingBytes(stdin), func() error { o.Type(string(stdin)); return nil })
				} else {
					o.Type(string(stdin))
				}
			}
			// MoveMouse（游標位置）：只送真實變化，不重置原版故事計時。
			if bx != lastX || by != lastY {
				if live != nil {
					live.Operation(o, audioOut, "move_mouse", recordingXY(bx, by), func() error { o.MoveMouse(bx, by); return nil })
				} else {
					o.MoveMouse(bx, by)
				}
				lastX, lastY = bx, by
			}
			for _, c := range clicks {
				// spec 004 R2／R3：與 autoplay 同一 API；錯誤只記 log
				var err error
				if live != nil {
					err = live.Operation(o, audioOut, "click", recordingXY(c.x, c.y), func() error { return o.Click(c.x, c.y) })
				} else {
					err = o.Click(c.x, c.y)
					audioOut.AfterOperation(o)
				}
				log.Printf("左鍵 (%d,%d) steps=%d err=%s", c.x, c.y, o.Steps(), frontendOperationCode(err))
				snapshot() // spec 005 R1：每下點擊後快照，避免兩步併成一次
				if stop.Requested() {
					return
				}
			}
			// 右鍵維持 Press/Release 路徑
			if r != lastR {
				if r {
					if live != nil {
						live.Operation(o, audioOut, "press_button", recordingButton(), func() error { o.PressButton(1); return nil })
					} else {
						o.PressButton(1)
					}
				} else {
					if live != nil {
						live.Operation(o, audioOut, "release_button", recordingButton(), func() error { o.ReleaseButton(1); return nil })
					} else {
						o.ReleaseButton(1)
					}
				}
				lastR = r
				log.Printf("右鍵 按下=%v steps=%d", r, o.Steps())
			}
			var runErr error
			if live != nil {
				runErr = live.Operation(o, audioOut, "run", recordingSteps(chunk), func() error { return o.Run(chunk) })
			} else {
				runErr = o.Run(chunk)
				audioOut.AfterOperation(o)
			}
			if runErr != nil {
				log.Printf("機器停止於 steps=%d：%s", o.Steps(), frontendOperationCode(runErr))
				return
			}
			if stop.Requested() {
				return
			}
			// 自動遊玩：步數過門檻執行（R10 步數門控）
			for autoIdx < len(autos) {
				a := autos[autoIdx]
				if o.Steps() < a.at {
					break
				}
				fmt.Printf("AUTO %d 執行 steps=%d\n", a.at, o.Steps())
				a.fn(o)
				audioOut.AfterOperation(o)
				autoIdx++
				snapshot()
				if stop.Requested() {
					return
				}
			}
			// HD 幀輸出（每 30 秒）：A0000 2x 放大＋文字事件疊畫
			if hdShot && time.Since(lastHD) > 3*time.Second {
				lastHD = time.Now()
				evMu.Lock()
				hdFrame(o, ov, fmt.Sprintf("workplace/bake/hd-%d.png", time.Now().Unix()))
				evMu.Unlock()
			}
			if i%60 == 0 {
				nz := 0
				for _, v := range o.Indexed() {
					if v != 0 {
						nz++
					}
				}
				_, _, rgb := o.ScreenRGB()
				nz2 := 0
				for j := 0; j+2 < len(rgb); j += 3 {
					if rgb[j] != 0 || rgb[j+1] != 0 || rgb[j+2] != 0 {
						nz2++
					}
				}
				pal := o.Palette()
				log.Printf("chunk %d steps=%d mode=%d idxNZ=%d rgbNZ=%d pal[15]=%v", i, o.Steps(), o.VideoMode(), nz, nz2, pal[15])
			}
			snapshot()
		}
	}()

	if *hdOnly {
		// 頭less 驗證：不開視窗（無 WM 時 GLFW init 會 panic），
		// 只跑機器 goroutine＋hdFrame 輸出。
		select {}
	}
	gameErr := ebiten.RunGame(g)
	if live != nil {
		live.Cancel()
	}
	stop.Request()
	deadline := time.NewTimer(max(time.Until(stop.Deadline()), 0))
	defer deadline.Stop()
	select {
	case <-ownerDone:
		if !time.Now().Before(stop.Deadline()) {
			log.Print("退出失敗：owner_timeout")
			if live != nil {
				log.Print("錄製失敗，檔案收尾未完成：owner_timeout")
				os.Exit(2)
			}
			os.Exit(1)
		}
	case <-deadline.C:
		log.Print("退出失敗：owner_timeout")
		if live != nil {
			log.Print("錄製失敗，檔案收尾未完成：owner_timeout")
			os.Exit(2)
		}
		os.Exit(1)
	}
	closed := make(chan string, 1)
	go func() { closed <- audioOut.shutdownPlayer() }()
	select {
	case result := <-closed:
		if !time.Now().Before(stop.Deadline()) {
			log.Print("退出失敗：close_timeout")
			if live != nil {
				log.Print("錄製失敗，檔案收尾未完成：close_timeout")
				os.Exit(2)
			}
			os.Exit(1)
		}
		log.Printf("音訊結束 fault=%s close=%s stream=%+v reader=%+v 畫格率=%.2f", audioOut.fault, result,
			audioOut.stream.Stats(), audioOut.reader.Stats(), float64(g.draws.Load())/time.Since(runStarted).Seconds())
		log.Printf("原版完成Step=%d 來源=%+v；Received為合成器合法接受，送入合成器=Received+Rejected；停止後未送入=%d", finalStep, sourceStats, audioOut.stoppedWrites)
		if result != "" || audioOut.fault != "" {
			if live != nil {
				live.exitCode(stop.Deadline())
				o.Close()
				if live.state.Load() != 2 {
					os.Exit(2)
				}
			}
			os.Exit(1)
		}
	case <-deadline.C:
		log.Print("退出失敗：close_timeout")
		if live != nil {
			log.Print("錄製失敗，檔案收尾未完成：close_timeout")
			os.Exit(2)
		}
		os.Exit(1)
	}
	if live != nil {
		code := live.exitCode(stop.Deadline())
		o.Close()
		if code != 0 {
			os.Exit(code)
		}
	}
	if gameErr != nil {
		log.Printf("前端結束：%v", gameErr)
	}
}

// GUI只持有顯示資料；不共享owner仍可能修改的事件、頁面或字形slice。
func cloneVisualEvents(in []*textEv) []*textEv {
	out := make([]*textEv, 0, len(in))
	for _, e := range in {
		v := &textEv{originRet: e.originRet, sex: e.sex, popupClip: e.popupClip, s: e.s, x0: e.x0, y0: e.y0, dy: e.dy, fg: e.fg, fgRGB: e.fgRGB}
		if e.profile != nil {
			v.profile = e.profile.visualCopy()
		}
		if e.party != nil {
			party := *e.party
			v.party = &party
		}
		for _, glyph := range e.g {
			v.g = append(v.g, append([]uint16(nil), glyph...))
		}
		out = append(out, v)
	}
	return out
}
