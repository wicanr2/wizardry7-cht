// 013 READY：單檔 wiz7-input/1。完整驗證後才交出可重播文件。
package recording

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wicanr2/wizardry7-cht/initconfig"
	"hash"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

const MaxBytes = 64 << 20
const MaxLine = 1 << 20
const MaxEvents = 1000000
const MaxType = 1024
const MaxFiles = 4096

type File struct {
	Path    string `json:"path"`
	Size    uint64 `json:"size"`
	SHA     string `json:"sha256"`
	DOSDate uint16 `json:"dos_date"`
	DOSTime uint16 `json:"dos_time"`
	Attr    uint8  `json:"attributes"`
}

type Clock struct {
	CPUHz      uint64 `json:"cpu_hz"`
	CycleClock bool   `json:"cycle_clock"`
	IRQ0Base   uint64 `json:"irq0_base"`
	IRQ0Every  uint64 `json:"irq0_every"`
	Ticks      uint64 `json:"ticks"`
	ModelSHA   string `json:"model_sha256"`
}

type Header struct {
	Kind      string            `json:"kind"`
	Schema    string            `json:"schema"`
	Start     uint64            `json:"start_step"`
	ConfigSHA string            `json:"config_sha256"`
	Config    initconfig.Config `json:"config"`
}

type Panel struct {
	HelpOpen   bool   `json:"help_open"`
	HelpScreen string `json:"help_screen"`
	HelpPage   int64  `json:"help_page"`
	MapOpen    bool   `json:"map_open"`
}
type Params struct {
	Panel       *Panel  `json:"panel,omitempty"`
	Key         *uint8  `json:"key,omitempty"`
	Bytes       *[]byte `json:"bytes,omitempty"`
	X           *int    `json:"x,omitempty"`
	Y           *int    `json:"y,omitempty"`
	Button      *uint8  `json:"button,omitempty"`
	Steps       *uint64 `json:"steps,omitempty"`
	Publication *uint64 `json:"publication,omitempty"`
	FrontendSHA *string `json:"frontend_sha256,omitempty"`
}

type Event struct {
	Kind   string `json:"kind"`
	Seq    uint64 `json:"seq"`
	Before uint64 `json:"before_step"`
	After  uint64 `json:"after_step"`
	Error  string `json:"error_class"`
	Params Params `json:"params"`
}

type Footer struct {
	Kind     string `json:"kind"`
	Count    uint64 `json:"event_count"`
	End      uint64 `json:"end_step"`
	SHA      string `json:"content_sha256"`
	Accepted bool   `json:"accepted"`
}

type Document struct {
	Header Header
	Events []Event
	Footer Footer
}
type recordFile interface {
	Write([]byte) (int, error)
	Truncate(int64) error
	Stat() (os.FileInfo, error)
	Sync() error
	Close() error
}
type Writer struct {
	finishDeadline           time.Time
	f                        recordFile
	name, partial            string
	h                        hash.Hash
	bytes                    int64
	count, step, publication uint64
	frontend                 string
	failed, closed, hadError bool
	lastKind                 string
	identity                 os.FileInfo
	contentEnd               int64
	allowed                  func() bool
	link                     func(string, string) error
	remove                   func(string) error
}

func digest(s string) bool {
	b, e := hex.DecodeString(s)
	return e == nil && len(b) == 32 && strings.ToLower(s) == s
}

// 重複JSON欄位包括巢狀物件也拒絕；encoding/json的未知欄位拒絕另處理。
func uniqueJSON(raw []byte) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value func(int) error
	value = func(depth int) error {
		if depth > 16 {
			return errors.New("巢狀層數超限")
		}
		t, e := d.Token()
		if e != nil {
			return e
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return e
				}
				s, ok := k.(string)
				if !ok || seen[s] {
					return errors.New("重複或非法欄位")
				}
				seen[s] = true
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		case '[':
			for d.More() {
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		default:
			return errors.New("非法JSON邊界")
		}
		_, e = d.Token()
		return e
	}
	if e := value(0); e != nil {
		return e
	}
	if _, e := d.Token(); e != io.EOF {
		return errors.New("多餘JSON內容")
	}
	return nil
}

func decode(raw []byte, out any) error {
	if len(raw) > MaxLine {
		return errors.New("單列超限")
	}
	if e := uniqueJSON(raw); e != nil {
		return e
	}
	if !validUnicodeJSON(raw) {
		return errors.New("非法Unicode編碼")
	}
	if e := exactFields(raw, reflect.TypeOf(out)); e != nil {
		return e
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if e := d.Decode(out); e != nil {
		return errors.New("未知欄位或非法型別")
	}
	switch out.(type) {
	case *Header:
		if e := requiredFields(raw, "kind", "schema", "start_step", "config_sha256", "config"); e != nil {
			return e
		}
	case *Event:
		if e := requiredFields(raw, "kind", "seq", "before_step", "after_step", "error_class", "params"); e != nil {
			return e
		}
	case *Footer:
		if e := requiredFields(raw, "kind", "event_count", "end_step", "content_sha256", "accepted"); e != nil {
			return e
		}
	}
	return nil
}

// 不允許JSON解碼器把非法UTF-8或孤立surrogate默默換成U+FFFD。
func validUnicodeJSON(raw []byte) bool {
	if !utf8.Valid(raw) {
		return false
	}
	hex4 := func(b []byte) (uint16, bool) {
		if len(b) < 4 {
			return 0, false
		}
		var n uint16
		for _, c := range b[:4] {
			n <<= 4
			switch {
			case c >= '0' && c <= '9':
				n += uint16(c - '0')
			case c >= 'a' && c <= 'f':
				n += uint16(c - 'a' + 10)
			case c >= 'A' && c <= 'F':
				n += uint16(c - 'A' + 10)
			default:
				return 0, false
			}
		}
		return n, true
	}
	inside := false
	for i := 0; i < len(raw); i++ {
		if raw[i] == '"' {
			inside = !inside
			continue
		}
		if !inside || raw[i] != '\\' {
			continue
		}
		i++
		if i >= len(raw) {
			return false
		}
		if raw[i] != 'u' {
			continue
		}
		n, ok := hex4(raw[i+1:])
		if !ok {
			return false
		}
		i += 4
		if n >= 0xdc00 && n <= 0xdfff {
			return false
		}
		if n >= 0xd800 && n <= 0xdbff {
			if i+6 >= len(raw) || raw[i+1] != '\\' || raw[i+2] != 'u' {
				return false
			}
			low, ok := hex4(raw[i+3:])
			if !ok || low < 0xdc00 || low > 0xdfff {
				return false
			}
			i += 6
		}
	}
	return !inside
}

// 欄位名稱依json tag逐字比對，不接受encoding/json的大小寫別名。
// omitempty只允許缺省，不允許null藉由nil指標隱藏額外引數。
func exactFields(raw []byte, t reflect.Type) error {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return errors.New("欄位不可為null")
	}
	switch t.Kind() {
	case reflect.Struct:
		var fields map[string]json.RawMessage
		if e := json.Unmarshal(raw, &fields); e != nil || fields == nil {
			return errors.New("非法物件")
		}
		known := make(map[string]reflect.StructField, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")
			if tag[0] == "" || tag[0] == "-" {
				continue
			}
			known[tag[0]] = f
			_, present := fields[tag[0]]
			optional := len(tag) == 2 && tag[1] == "omitempty"
			if !present && !optional {
				return errors.New("缺少必要欄位")
			}
		}
		for name, value := range fields {
			f, ok := known[name]
			if !ok {
				return errors.New("未知或大小寫不符欄位")
			}
			if e := exactFields(value, f.Type); e != nil {
				return e
			}
		}
	case reflect.Slice:
		// []byte允base64或uint8陣列；只有字串略過元素檢查。
		// JSON的null元素不能交標準解碼器默默轉成byte0。
		if t.Elem().Kind() == reflect.Uint8 && len(bytes.TrimSpace(raw)) > 0 && bytes.TrimSpace(raw)[0] == '"' {
			return nil
		}
		var values []json.RawMessage
		if e := json.Unmarshal(raw, &values); e != nil || values == nil {
			return errors.New("非法陣列")
		}
		for _, value := range values {
			if e := exactFields(value, t.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}

func requiredFields(raw []byte, names ...string) error {
	var fields map[string]json.RawMessage
	if e := json.Unmarshal(raw, &fields); e != nil {
		return errors.New("非法物件")
	}
	for _, name := range names {
		value, ok := fields[name]
		if !ok || bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return errors.New("缺少必要欄位")
		}
	}
	return nil
}

// Build驗證的是自述配置；不證明實際機器具有這個初態。
func NewHeader(c initconfig.Config) (Header, error) {
	validated, _, hash, err := initconfig.Validate(c)
	if err != nil {
		return Header{}, errors.New("invalid_config")
	}
	h := Header{"header", "wiz7-input/1", 0, hash, validated}
	if err = checkHeader(h); err != nil {
		return Header{}, err
	}
	return h, nil
}

func checkHeader(h Header) error {
	if h.Kind != "header" || h.Schema != "wiz7-input/1" || h.Start != 0 {
		return errors.New("header_schema")
	}
	_, _, hash, err := initconfig.Validate(h.Config)
	if err != nil {
		return errors.New("invalid_config")
	}
	if !digest(h.ConfigSHA) || h.ConfigSHA != hash {
		return errors.New("config_hash")
	}
	encoded, err := json.Marshal(h)
	if err != nil || len(encoded) > MaxLine {
		return errors.New("header_capacity")
	}
	var roundTrip Header
	if err = decode(encoded, &roundTrip); err != nil || !reflect.DeepEqual(h, roundTrip) {
		return errors.New("header_encoding")
	}
	return nil
}

// 僅在整檔驗完且與呼叫端的實際預期SHA相符時交出文件。
// 本套件沒有 LoadWith 或原版輸入投遞 API。
func LoadExpected(name, expected string) (Document, error) {
	if !digest(expected) {
		return Document{}, errors.New("expected_config_hash")
	}
	doc, err := Load(name)
	if err != nil {
		return Document{}, err
	}
	if doc.Header.ConfigSHA != expected {
		return Document{}, errors.New("config_mismatch")
	}
	return doc, nil
}

func checkEvent(e Event, seq, step, publication uint64, frontend string) (uint64, error) {
	if e.Seq != seq || e.Before != step || e.After < e.Before {
		return publication, errors.New("命令序號或步數不連續")
	}
	if e.Error != "none" && e.Error != "no_response" && e.Error != "execution" {
		return publication, errors.New("未知錯誤類別")
	}
	if e.Kind != "click" && e.Kind != "run" && (e.After != e.Before || e.Error != "none") {
		return publication, errors.New("非執行命令不可推進")
	}
	p := e.Params
	fields := []bool{p.Key != nil, p.Bytes != nil, p.X != nil, p.Y != nil, p.Button != nil, p.Steps != nil, p.Publication != nil, p.FrontendSHA != nil, p.Panel != nil}
	required := make([]bool, 9)
	switch e.Kind {
	case "press_key":
		required[0] = true
	case "type":
		required[1] = true
		if p.Bytes == nil || len(*p.Bytes) == 0 || len(*p.Bytes) > MaxType {
			return publication, errors.New("文字容量超限")
		}
	case "move_mouse", "click":
		required[2] = true
		required[3] = true
		if p.X == nil || p.Y == nil || *p.X < 0 || *p.X >= 320 || *p.Y < 0 || *p.Y >= 200 {
			return publication, errors.New("非法原版座標")
		}
	case "press_button", "release_button":
		required[4] = true
		if p.Button == nil || *p.Button != 1 {
			return publication, errors.New("非既有右鍵路徑")
		}
	case "run":
		required[5] = true
		if p.Steps == nil || *p.Steps == 0 || *p.Steps > 200000 || e.After-e.Before > *p.Steps || (e.Error == "none" && e.After-e.Before != *p.Steps) || e.Error == "no_response" {
			return publication, errors.New("非法Run水位")
		}
	case "snapshot":
		required[8] = true
		if p.Panel == nil || len(p.Panel.HelpScreen) > 64 || !utf8.ValidString(p.Panel.HelpScreen) || strings.ContainsAny(p.Panel.HelpScreen, "\x00\r\n") {
			return publication, errors.New("panel_state")
		}
		required[6] = true
		required[7] = true
		if p.Publication == nil || *p.Publication != publication+1 || p.FrontendSHA == nil || *p.FrontendSHA != frontend {
			return publication, errors.New("快照發布關聯不符")
		}
		publication = *p.Publication
	default:
		return publication, errors.New("未知命令")
	}
	for i := range fields {
		if fields[i] != required[i] {
			return publication, errors.New("命令引數不符")
		}
	}
	return publication, nil
}

func safeParent(name string) error {
	parent, e := filepath.Abs(filepath.Dir(name))
	if e != nil {
		return e
	}
	for {
		info, e := os.Lstat(parent)
		if e != nil {
			return e
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return errors.New("輸出父路徑不可含連結")
		}
		next := filepath.Dir(parent)
		if next == parent {
			return nil
		}
		parent = next
	}
}

func New(name string, h Header) (*Writer, error) {
	if e := checkHeader(h); e != nil {
		return nil, e
	}
	if e := safeParent(name); e != nil {
		return nil, e
	}
	if strings.HasSuffix(name, ".partial") {
		return nil, errors.New("輸出名稱不得為未完成檔")
	}
	if _, e := os.Lstat(name); !os.IsNotExist(e) {
		return nil, errors.New("錄製檔已存在或不可建立")
	}
	partial := name + ".partial"
	f, e := os.OpenFile(partial, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return nil, errors.New("錄製檔建立失敗")
	}
	w := &Writer{f: f, name: name, partial: partial, h: sha256.New(), step: h.Start, frontend: h.ConfigSHA}
	if e = w.line(h, true); e != nil {
		f.Close()
		return nil, e
	}
	return w, nil
}

func (w *Writer) line(v any, include bool) error {
	if w.closed || w.failed || !w.deadlineOK() {
		return errors.New("錄製器未完成")
	}
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	b = append(b, '\n')
	if len(b) > MaxLine+1 || w.bytes+int64(len(b)) > MaxBytes {
		w.failed = true
		return errors.New("錄製容量超限")
	}
	n, e := w.f.Write(b)
	if e != nil || n != len(b) || !w.deadlineOK() {
		w.failed = true
		return errors.New("錄製落盤失敗")
	}
	w.bytes += int64(n)
	if include {
		_, _ = w.h.Write(b)
	}
	return nil
}

func (w *Writer) Append(e Event) error {
	if w.count >= MaxEvents {
		w.failed = true
		return errors.New("錄製事件超限")
	}
	p, err := checkEvent(e, w.count+1, w.step, w.publication, w.frontend)
	if err != nil {
		w.failed = true
		return err
	}
	if err = w.line(e, true); err != nil {
		return err
	}
	w.lastKind = e.Kind
	w.count++
	w.step = e.After
	w.publication = p
	w.hadError = w.hadError || e.Error != "none"
	return nil
}

// FinishAt 使用同一絕對期限。Writer 僅供單一擁有者使用。
func (w *Writer) FinishAt(end uint64, deadline time.Time) error {
	if deadline.IsZero() {
		return errors.New("writer_deadline")
	}
	w.finishDeadline = deadline
	return w.Finish(end)
}

// FinishWhen 在每個磁碟操作返回後核對共同取消狀態。
// 回呼只讀生命週期狀態，不可呼叫 Writer。阻塞的系統呼叫仍由本擁有者收尾。
func (w *Writer) FinishWhen(end uint64, deadline time.Time, allowed func() bool) error {
	w.allowed = allowed
	return w.FinishAt(end, deadline)
}
func (w *Writer) deadlineOK() bool {
	return (w.finishDeadline.IsZero() || time.Now().Before(w.finishDeadline)) && (w.allowed == nil || w.allowed())
}
func (w *Writer) Finish(end uint64) error {
	if w.closed {
		return errors.New("錄製器已結束")
	}
	w.contentEnd = w.bytes
	identity, statErr := w.f.Stat()
	w.identity = identity
	bad := w.failed || w.hadError || w.count == 0 || w.lastKind != "snapshot" || end != w.step || statErr != nil || !w.deadlineOK()
	if !bad {
		bad = w.f.Sync() != nil || !w.deadlineOK()
	}
	if !bad {
		bad = w.line(Footer{"end", w.count, w.step, hex.EncodeToString(w.h.Sum(nil)), true}, false) != nil || !w.deadlineOK()
	}
	if !bad {
		bad = w.f.Sync() != nil || !w.deadlineOK()
	}
	if bad {
		_ = w.f.Truncate(w.contentEnd)
	}
	closeErr := w.f.Close()
	w.closed = true
	bad = bad || closeErr != nil || !w.deadlineOK()
	if !bad {
		bad = w.linkFile(w.partial, w.name) != nil || !w.deadlineOK()
	}
	if !bad {
		// 清理返回失敗也走同 owner 撤回，不能忽略 IO 故障而回報成功。
		bad = w.removeFile(w.partial) != nil || !w.deadlineOK()
	}
	if bad {
		if err := w.Revoke(); err != nil {
			return err
		}
		return errors.New("writer_finish")
	}
	return nil
}

func (w *Writer) linkFile(from, to string) error {
	if w.link != nil {
		return w.link(from, to)
	}
	return os.Link(from, to)
}
func (w *Writer) removeFile(name string) error {
	if w.remove != nil {
		return w.remove(name)
	}
	return os.Remove(name)
}

// Revoke 只由完成 Finish 的同一 worker 呼叫，絕不與阻塞 IO 並行。
// 同 inode 的成功名撤回前先保留 partial；替換檔、連結及其他 inode 一律不動。
func (w *Writer) Revoke() error {
	w.failed = true
	fail := func() error { return errors.New("cleanup_incomplete") }
	if !w.closed || w.identity == nil {
		return fail()
	}
	owned := func(info os.FileInfo) bool {
		return info != nil && info.Mode().IsRegular() && os.SameFile(w.identity, info)
	}
	nameInfo, nameErr := os.Lstat(w.name)
	partialInfo, partialErr := os.Lstat(w.partial)
	if owned(nameInfo) && os.IsNotExist(partialErr) {
		if w.linkFile(w.name, w.partial) != nil {
			return fail()
		}
		partialInfo, partialErr = os.Lstat(w.partial)
	}
	if partialErr != nil || !owned(partialInfo) {
		return fail()
	}
	f, err := os.OpenFile(w.partial, os.O_WRONLY, 0)
	if err != nil {
		return fail()
	}
	current, err := f.Stat()
	good := err == nil && owned(current)
	if good {
		good = f.Truncate(w.contentEnd) == nil
	}
	if good {
		good = f.Sync() == nil
	}
	if f.Close() != nil {
		good = false
	}
	if owned(nameInfo) {
		// 再核身份，避免移除他人換入的檔案。
		current, err := os.Lstat(w.name)
		if err != nil || !owned(current) || w.removeFile(w.name) != nil {
			good = false
		}
	} else if nameErr != nil && !os.IsNotExist(nameErr) {
		good = false
	}
	if !good {
		return fail()
	}
	return nil
}

func (w *Writer) Abort() error {
	if w.closed {
		return nil
	}
	w.failed = true
	w.closed = true
	return w.f.Close()
}

// 必須整檔驗完才回傳。
func Load(name string) (Document, error) {
	var out Document
	if strings.HasSuffix(name, ".partial") {
		return out, errors.New("未完成錄製")
	}
	if e := safeParent(name); e != nil {
		return out, errors.New("非法錄製父路徑")
	}
	st, e := os.Lstat(name)
	if e != nil || !st.Mode().IsRegular() || st.Size() > MaxBytes || st.Size() == 0 {
		return out, errors.New("非法錄製檔")
	}
	f, e := os.Open(name)
	if e != nil {
		return out, errors.New("錄製讀取失敗")
	}
	defer f.Close()
	r := bufio.NewReaderSize(io.LimitReader(f, MaxBytes+1), MaxLine+1)
	h := sha256.New()
	var step, publication uint64
	var total int64
	ended := false
	hadError := false
	line := 0
	for {
		raw, e := r.ReadSlice('\n')
		if e == io.EOF && len(raw) == 0 {
			break
		}
		if e != nil || len(raw) > MaxLine+1 {
			return Document{}, errors.New("不完整或超限列")
		}
		total += int64(len(raw))
		if total > MaxBytes {
			return Document{}, errors.New("錄製容量超限")
		}
		line++
		if ended {
			return Document{}, errors.New("結尾後多餘資料")
		}
		body := raw[:len(raw)-1]
		if line == 1 {
			if e = decode(body, &out.Header); e != nil {
				return Document{}, e
			}
			if e = checkHeader(out.Header); e != nil {
				return Document{}, e
			}
			step = out.Header.Start
			_, _ = h.Write(raw)
			continue
		}
		var tag struct {
			Kind string `json:"kind"`
		}
		if e = json.Unmarshal(body, &tag); e != nil {
			return Document{}, errors.New("非法JSON")
		}
		if tag.Kind == "end" {
			if e = decode(body, &out.Footer); e != nil {
				return Document{}, e
			}
			if out.Footer.Count != uint64(len(out.Events)) || out.Footer.End != step || out.Footer.SHA != hex.EncodeToString(h.Sum(nil)) || !out.Footer.Accepted || hadError || len(out.Events) == 0 || out.Events[len(out.Events)-1].Kind != "snapshot" {
				return Document{}, errors.New("結尾或驗收資格不符")
			}
			ended = true
			continue
		}
		if len(out.Events) >= MaxEvents {
			return Document{}, errors.New("事件數超限")
		}
		var event Event
		if e = decode(body, &event); e != nil {
			return Document{}, e
		}
		publication, e = checkEvent(event, uint64(len(out.Events))+1, step, publication, out.Header.ConfigSHA)
		if e != nil {
			return Document{}, e
		}
		// 結構可重播與成功驗收資格分開；實際API類別/Step仍須由重播者核對。
		hadError = hadError || event.Error != "none"
		step = event.After
		out.Events = append(out.Events, event)
		_, _ = h.Write(raw)
	}
	if !ended {
		return Document{}, fmt.Errorf("未完成錄製，列數%d", line)
	}
	return out, nil
}
