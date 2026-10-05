package recording

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/wicanr2/dosgolem/oracle"
	"github.com/wicanr2/wizardry7-cht/initconfig"
)

// 全部自製資料，不讀取、載入或代製任何原版程式。
func fixture(t *testing.T) Header {
	t.Helper()
	h := initconfig.Hash([]byte("self-authored fixture"))
	file := initconfig.File{Path: "GAME.EXE", Size: 17, SHA256: h, DOSDate: 23875, DOSTime: 25692, Attributes: 32}
	audio, err := initconfig.AudioConfig(24000000, true, false)
	if err != nil {
		t.Fatal(err)
	}
	c := initconfig.Config{
		Schema:   "wiz7-startup-config/1",
		Program:  initconfig.Program{SHA256: h, GoVersion: "go1.24.13", GOOS: "linux", GOARCH: "amd64", ModuleVersion: "(devel)", BuildSettingsSHA256: h},
		Sources:  []initconfig.Source{{Path: "fixture.go", Size: 17, SHA256: h}},
		LoadWith: initconfig.LoadOptions{Adlib: true},
		Initial: initconfig.Scalars{
			Schema:   "wiz7-initial-scalars/1",
			Clock:    initconfig.Clock{CPUHz: 33000000, CycleClock: false, IRQ0Base: 165000, IRQ0Every: 165000},
			DOS:      oracle.DOSSettings{FontFull: "END_S13.DAT", FontHalf: "END_S14.DAT", FontFullBytes: 30, FontHalfBytes: 15, Drive: 2, Dir: "RICH2"},
			Platform: initconfig.Platform{DOSDate: initconfig.Date{Year: 1993, Month: 1, Day: 1, Weekday: 5}, DOSTimeModel: "now_plus_ticks_5493_div_1000_mod_8640000", BDATicksObservation: "not_observed", EffectiveMouseScaleObservation: "not_observed"},
		},
		SelectedProgram: initconfig.SelectedProgram{SourceRole: "root", File: file},
		Storage:         initconfig.Storage{Resolution: "scratch_then_root", MetadataModel: "dos_archive_mtime_source_policy", Root: []initconfig.File{file}, Scratch: []initconfig.File{}},
		Audio:           audio,
		Frontend:        initconfig.Frontend{Scale: 3, HDComposeDisabled: true, AutoplayDisabled: true, TestInputsDisabled: true, Resources: []initconfig.Resource{{Role: "xlate"}, {Role: "flow"}, {Role: "help"}, {Role: "roster"}, {Role: "profile-sex"}, {Role: "intro-picture"}}},
	}
	header, err := NewHeader(c)
	if err != nil {
		t.Fatal(err)
	}
	return header
}

func encode(t *testing.T, x any) []byte {
	t.Helper()
	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// 保持外層內容雜湊正確，負向案例才能驗證檔頭內部守門。
func document(t *testing.T, header []byte) string {
	t.Helper()
	content := append(append([]byte{}, header...), '\n')
	h := fixture(t)
	pub := uint64(1)
	event := Event{"snapshot", 1, 0, 0, "none", Params{Publication: &pub, FrontendSHA: &h.ConfigSHA, Panel: &Panel{}}}
	content = append(content, encode(t, event)...)
	content = append(content, '\n')
	end := Footer{"end", 1, 0, initconfig.Hash(content), true}
	all := append(content, encode(t, end)...)
	all = append(all, '\n')
	name := filepath.Join(t.TempDir(), "recording.jsonl")
	if err := os.WriteFile(name, all, 0600); err != nil {
		t.Fatal(err)
	}
	return name
}

func TestSingleFileRoundTrip(t *testing.T) {
	h := fixture(t)
	name := filepath.Join(t.TempDir(), "recording.jsonl")
	writer, err := New(name, h)
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Abort()
	steps, key, pub := uint64(200000), uint8(82), uint64(1)
	text := []byte("self-authored input")
	events := []Event{
		{"run", 1, 0, 200000, "none", Params{Steps: &steps}},
		{"press_key", 2, 200000, 200000, "none", Params{Key: &key}},
		{"type", 3, 200000, 200000, "none", Params{Bytes: &text}},
		{"snapshot", 4, 200000, 200000, "none", Params{Publication: &pub, FrontendSHA: &h.ConfigSHA, Panel: &Panel{}}},
	}
	for _, event := range events {
		if err := writer.Append(event); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Finish(200000); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadExpected(name, h.ConfigSHA)
	if err != nil || !reflect.DeepEqual(loaded.Header, h) || !reflect.DeepEqual(loaded.Events, events) || !loaded.Footer.Accepted {
		t.Fatal("roundtrip", err)
	}
	raw, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(bytes.Split(raw, []byte("\n"))[0], []byte(`"config":{"schema":`)) {
		t.Fatal("config not embedded")
	}
	if _, err := os.Stat(name + ".partial"); !os.IsNotExist(err) {
		t.Fatal("successful partial retained", err)
	}
	if out, err := LoadExpected(name, initconfig.Hash([]byte("other configuration"))); err == nil || !reflect.DeepEqual(out, Document{}) {
		t.Fatal("expected mismatch leaked document")
	}
	if out, err := LoadExpected(name, "bad"); err == nil || !reflect.DeepEqual(out, Document{}) {
		t.Fatal("invalid expectation accepted")
	}
}

func TestNestedHeaderAttacks(t *testing.T) {
	raw := encode(t, fixture(t))
	replace := func(from, to string) []byte {
		if bytes.Count(raw, []byte(from)) != 1 {
			t.Fatalf("ambiguous fixture edit %q", from)
		}
		return bytes.Replace(raw, []byte(from), []byte(to), 1)
	}
	attacks := map[string][]byte{
		"unknown-nested":      replace(`"font_full":"END_S13.DAT"`, `"font_full":"END_S13.DAT","unknown":0`),
		"duplicate-nested":    replace(`"font_full":"END_S13.DAT"`, `"font_full":"END_S13.DAT","font_full":"OTHER"`),
		"escaped-duplicate":   replace(`"font_full":"END_S13.DAT"`, `"font_full":"END_S13.DAT","font_\u0066ull":"OTHER"`),
		"case-alias":          replace(`"font_full":"END_S13.DAT"`, `"Font_Full":"END_S13.DAT"`),
		"missing-zero":        replace(`"stored_mouse_x_scale":0,`, ``),
		"null-zero":           replace(`"stored_mouse_x_scale":0`, `"stored_mouse_x_scale":null`),
		"null-object":         replace(`"load_with":{"adlib":true,"dgroup":0,"mouse_x_scale":0,"font_full":"","font_half":""}`, `"load_with":null`),
		"null-empty-scratch":  replace(`"scratch":[]`, `"scratch":null`),
		"null-array-entry":    replace(`"sources":[{`, `"sources":[null,{`),
		"missing-array-field": replace(`"sources":[{"path":"fixture.go","size":17,`, `"sources":[{"path":"fixture.go",`),
		"fractional-integer":  replace(`"stored_mouse_x_scale":0`, `"stored_mouse_x_scale":0.5`),
		"overflow-byte":       replace(`"drive":2`, `"drive":256`),
		"wrong-bool":          replace(`"adlib":true`, `"adlib":"true"`),
		"isolated-surrogate":  replace(`"dir":"RICH2"`, `"dir":"\ud800"`),
		"invalid-utf8":        replace(`"dir":"RICH2"`, string(append([]byte(`"dir":"`), append([]byte{0xff}, []byte(`"`)...)...))),
		"extra-value":         append(append([]byte{}, raw...), []byte(` {}`)...),
		"old-schema":          replace(`"schema":"wiz7-input/1"`, `"schema":"wiz7-input-prototype/1"`),
	}
	for name, header := range attacks {
		t.Run(name, func(t *testing.T) {
			out, err := Load(document(t, header))
			if err == nil || !reflect.DeepEqual(out, Document{}) {
				t.Fatal("malformed header accepted")
			}
		})
	}
}

func TestRequiredConfigFields(t *testing.T) {
	// 每一種物件及陣列元素的必要欄位都實際刪除／置null，包含零值欄位。
	raw := encode(t, fixture(t))
	var original map[string]any
	if err := json.Unmarshal(raw, &original); err != nil {
		t.Fatal(err)
	}
	type edit struct {
		path []any
		key  string
	}
	var edits []edit
	var visit func(any, []any)
	visit = func(v any, p []any) {
		switch x := v.(type) {
		case map[string]any:
			for key, value := range x {
				edits = append(edits, edit{append([]any{}, p...), key})
				visit(value, append(append([]any{}, p...), key))
			}
		case []any:
			for i, value := range x {
				visit(value, append(append([]any{}, p...), i))
			}
		}
	}
	visit(original, nil)
	for _, mutation := range edits {
		for _, mode := range []string{"missing", "null"} {
			var fresh map[string]any
			if err := json.Unmarshal(raw, &fresh); err != nil {
				t.Fatal(err)
			}
			var target any = fresh
			for _, part := range mutation.path {
				switch x := part.(type) {
				case string:
					target = target.(map[string]any)[x]
				case int:
					target = target.([]any)[x]
				}
			}
			m := target.(map[string]any)
			if mode == "missing" {
				delete(m, mutation.key)
			} else {
				m[mutation.key] = nil
			}
			if out, err := Load(document(t, encode(t, fresh))); err == nil || !reflect.DeepEqual(out, Document{}) {
				t.Fatalf("accepted %s at %v.%s", mode, mutation.path, mutation.key)
			}
		}
	}
	t.Logf("validated %d required paths with deletion and null", len(edits))
}

func TestConfigHashAndSemanticMismatch(t *testing.T) {
	h := fixture(t)
	h.Config.Frontend.Scale++ // 改成仍然有效的配置，保留舊SHA。
	if _, err := Load(document(t, encode(t, h))); err == nil {
		t.Fatal("modified config accepted with old SHA")
	}
	h = fixture(t)
	h.Config.Audio.K.Numerator = "2"
	h.ConfigSHA = initconfig.Hash(encode(t, h.Config))
	if _, err := Load(document(t, encode(t, h))); err == nil {
		t.Fatal("inconsistent tempo accepted with recomputed SHA")
	}
	h = fixture(t)
	h.Config.Storage.Root[0].Path = "../escape"
	h.ConfigSHA = initconfig.Hash(encode(t, h.Config))
	if _, err := Load(document(t, encode(t, h))); err == nil {
		t.Fatal("escape accepted with recomputed SHA")
	}
}

func TestCanonicalConfigHash(t *testing.T) {
	h := fixture(t)
	original := encode(t, h)
	var object any
	if err := json.Unmarshal(original, &object); err != nil {
		t.Fatal(err)
	}
	sorted, _ := json.MarshalIndent(object, "", " ")
	// JSON列不能含實際換行；使用空白及map排序改變欄序。
	sorted = bytes.ReplaceAll(sorted, []byte("\n"), []byte(" "))
	loaded, err := Load(document(t, sorted))
	if err != nil || !reflect.DeepEqual(loaded.Header, h) {
		t.Fatal("equivalent field order rejected", err)
	}
	if h.ConfigSHA == initconfig.Hash(sorted) {
		t.Fatal("test failed to distinguish header bytes from config hash")
	}
}
