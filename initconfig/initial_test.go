// 自製設定與目錄，沒有載入或代製原版遊戲。
package initconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/wicanr2/dosgolem/oracle"
)

func writeFixture(t *testing.T, dir, name, body string) {
	t.Helper()
	p := filepath.Join(dir, name)
	if e := os.WriteFile(p, []byte(body), 0600); e != nil {
		t.Fatal(e)
	}
	m := time.Date(2026, 10, 3, 12, 34, 56, 0, time.Local)
	if e := os.Chtimes(p, m, m); e != nil {
		t.Fatal(e)
	}
}
func configuration(t *testing.T) (Config, string, string) {
	t.Helper()
	root, scratch := t.TempDir(), t.TempDir()
	writeFixture(t, root, "GAME.EXE", "self-authored executable fixture")
	writeFixture(t, root, "SAVE.DAT", "self-authored root progress")
	writeFixture(t, scratch, "SAVE.DAT", "self-authored scratch progress")
	r, e := CollectFiles(root)
	if e != nil {
		t.Fatal(e)
	}
	s, e := CollectFiles(scratch)
	if e != nil {
		t.Fatal(e)
	}
	a, e := AudioConfig(3004074.5544433594, false, false)
	if e != nil {
		t.Fatal(e)
	}
	c := Config{
		Schema:          "wiz7-startup-config/1",
		Program:         Program{Hash([]byte("self-authored program")), "go1.24.13", "linux", "amd64", "(devel)", Hash([]byte("self-authored build settings"))},
		Sources:         []Source{{"fixture.go", 17, Hash([]byte("self-authored source"))}},
		LoadWith:        LoadOptions{Adlib: true},
		Initial:         Scalars{"wiz7-initial-scalars/1", 0, 0, Clock{33000000, false, 165000, 165000}, oracle.DOSSettings{StoredMouseXScale: 0, FontFull: "END_S13.DAT", FontHalf: "END_S14.DAT", FontFullBytes: 30, FontHalfBytes: 15, Drive: 2, Dir: "RICH2"}, fixedPlatform()},
		SelectedProgram: SelectedProgram{"root", r[0]},
		Storage:         Storage{"scratch_then_root", "dos_archive_mtime_source_policy", r, s},
		Audio:           a,
		Frontend:        Frontend{3, true, true, true, []Resource{{Role: "xlate"}, {Role: "flow"}, {Role: "help"}, {Role: "roster"}, {Role: "profile-sex"}, {Role: "intro-picture"}}},
	}
	return c, root, scratch
}

func TestIndependentAudioAndDOSValues(t *testing.T) {
	c, _, _ := configuration(t)
	if c.Audio.H != (Rational{"12304689375", "4096"}) || c.Audio.K != (Rational{"1", "1"}) {
		t.Fatal("default clock")
	}
	a, e := AudioConfig(24000000, true, false)
	if e != nil {
		t.Fatal(e)
	}
	if a.K != (Rational{"52428800", "6562501"}) || a.IPSFloat64Bits != "4176e36000000000" {
		t.Fatal("accelerated clock")
	}
	if a.SampleRate != 44100 || a.DeviceChannels != 2 || a.SampleBits != 16 || a.RingFrames != 88200 || a.PrefillFrames != 4410 || a.MaximumPCMChunkFrames != 4410 || a.FadeFrames != 221 || a.DeviceBufferMilliseconds != 50 || a.ProgressQuantum != 200000 || a.ExitBudgetMilliseconds != 5000 {
		t.Fatal("source-bound audio constants")
	}
	for _, f := range append(c.Storage.Root, c.Storage.Scratch...) {
		if f.DOSDate != 23875 || f.DOSTime != 25692 || f.Attributes != 32 {
			t.Fatal("independent DOS date/time")
		}
	}
}

func TestStorageClosedLayers(t *testing.T) {
	c, root, scratch := configuration(t)
	if len(c.Storage.Root) != 2 || len(c.Storage.Scratch) != 1 || c.Storage.Root[1].Path != c.Storage.Scratch[0].Path || c.Storage.Root[1].SHA256 == c.Storage.Scratch[0].SHA256 {
		t.Fatal("shadowed root lost")
	}
	if e := VerifyStorage(c.Storage, root, scratch); e != nil {
		t.Fatal(e)
	}
	writeFixture(t, scratch, "NEW.DAT", "new")
	if VerifyStorage(c.Storage, root, scratch) == nil {
		t.Fatal("new file accepted")
	}
	if e := os.Remove(filepath.Join(scratch, "NEW.DAT")); e != nil {
		t.Fatal(e)
	}
	writeFixture(t, root, "SAVE.DAT", "changed shadowed root")
	if VerifyStorage(c.Storage, root, scratch) == nil {
		t.Fatal("shadowed root mutation accepted")
	}
	writeFixture(t, root, "SAVE.DAT", "self-authored root progress")
	p := filepath.Join(scratch, "SAVE.DAT")
	m := time.Date(2026, 10, 4, 12, 34, 56, 0, time.Local)
	if e := os.Chtimes(p, m, m); e != nil {
		t.Fatal(e)
	}
	if VerifyStorage(c.Storage, root, scratch) == nil {
		t.Fatal("DOS date mutation accepted")
	}
	if e := os.Remove(p); e != nil {
		t.Fatal(e)
	}
	if VerifyStorage(c.Storage, root, scratch) == nil {
		t.Fatal("removed file accepted")
	}
	c.Storage.Scratch = []File{}
	if e := VerifyStorage(c.Storage, root, scratch); e != nil {
		t.Fatal(e)
	}
	if _, _, _, e := Validate(c); e != nil {
		t.Fatal("empty scratch", e)
	}
	if e := os.Symlink(filepath.Join(root, "GAME.EXE"), filepath.Join(scratch, "link")); e != nil {
		t.Fatal(e)
	}
	if _, e := CollectFiles(scratch); e == nil {
		t.Fatal("symlink accepted")
	}
}

func TestFingerprintSensitivityAndOwnership(t *testing.T) {
	c, _, _ := configuration(t)
	owned, b, h, e := Validate(c)
	if e != nil {
		t.Fatal(e)
	}
	if Hash(b) != h {
		t.Fatal("digest")
	}
	before := append([]byte{}, b...)
	c.Sources[0].Size++
	c.Storage.Root[0].Size++
	c.Storage.Scratch[0].Size++
	c.Frontend.Resources[0].Requested = true
	after, e := json.Marshal(owned)
	if e != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("caller slice alias")
	}
	cases := map[string]func(*Config){
		"binary":       func(c *Config) { c.Program.SHA256 = Hash([]byte("other binary")) },
		"source":       func(c *Config) { c.Sources[0].SHA256 = Hash([]byte("other source")) },
		"load_option":  func(c *Config) { c.LoadWith.DGROUP = 65537 },
		"DOS_setting":  func(c *Config) { c.Initial.DOS.Dir = "OTHER" },
		"resource":     func(c *Config) { c.Frontend.Resources[0] = Resource{"xlate", true, true, 3, Hash([]byte("abc"))} },
		"scale":        func(c *Config) { c.Frontend.Scale = 7 },
		"mute":         func(c *Config) { c.Audio, _ = AudioConfig(3004074.5544433594, false, true) },
		"explicit_one": func(c *Config) { c.Audio, _ = AudioConfig(3004074.5544433594, true, false) },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			x, _, _ := configuration(t)
			change(&x)
			_, _, newHash, e := Validate(x)
			if e != nil || newHash == h {
				t.Fatal("configuration difference lost", e)
			}
		})
	}
}

func TestRejectIncompleteBoundaries(t *testing.T) {
	cases := map[string]func(*Config){
		"step": func(c *Config) { c.Initial.Step = 1 }, "ticks": func(c *Config) { c.Initial.Ticks = 1 },
		"schema": func(c *Config) { c.Schema = "" }, "source_missing": func(c *Config) { c.Sources = nil },
		"scratch_null": func(c *Config) { c.Storage.Scratch = nil }, "root_null": func(c *Config) { c.Storage.Root = nil },
		"selected_other":       func(c *Config) { c.SelectedProgram.File.SHA256 = Hash([]byte("other")) },
		"audio_incoherent":     func(c *Config) { c.Audio.K.Numerator = "2" },
		"adlib":                func(c *Config) { c.LoadWith.Adlib = false },
		"extra_input_producer": func(c *Config) { c.Frontend.AutoplayDisabled = false },
		"resource_missing":     func(c *Config) { c.Frontend.Resources = nil },
		"loaded_no_sha":        func(c *Config) { c.Frontend.Resources[0].Requested = true; c.Frontend.Resources[0].Loaded = true },
		"duplicate_path":       func(c *Config) { c.Storage.Root = append(c.Storage.Root, c.Storage.Root[1]) },
		"escape":               func(c *Config) { c.Storage.Scratch[0].Path = "../other" },
		"DOS_date_invalid":     func(c *Config) { c.Storage.Scratch[0].DOSDate = 0 },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			c, _, _ := configuration(t)
			change(&c)
			if _, _, _, e := Validate(c); e == nil {
				t.Fatal("accepted")
			}
		})
	}
}
