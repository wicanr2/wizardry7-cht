// 013 READY：完整啟動設定的型別、驗證與檔案身份。沒有原版輸入入口。
package initconfig

import (
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/wicanr2/dosgolem/oracle"
)

type File struct {
	Path       string `json:"path"`
	Size       uint64 `json:"size"`
	SHA256     string `json:"sha256"`
	DOSDate    uint16 `json:"dos_date"`
	DOSTime    uint16 `json:"dos_time"`
	Attributes uint8  `json:"attributes"`
}
type Source struct {
	Path   string `json:"path"`
	Size   uint64 `json:"size"`
	SHA256 string `json:"sha256"`
}
type Program struct {
	SHA256              string `json:"sha256"`
	GoVersion           string `json:"go_version"`
	GOOS                string `json:"goos"`
	GOARCH              string `json:"goarch"`
	ModuleVersion       string `json:"module_version"`
	BuildSettingsSHA256 string `json:"build_settings_sha256"`
}
type LoadOptions struct {
	Adlib       bool   `json:"adlib"`
	DGROUP      uint32 `json:"dgroup"`
	MouseXScale uint16 `json:"mouse_x_scale"`
	FontFull    string `json:"font_full"`
	FontHalf    string `json:"font_half"`
}
type Clock struct {
	CPUHz      uint64 `json:"cpu_hz"`
	CycleClock bool   `json:"cycle_clock"`
	IRQ0Base   uint64 `json:"irq0_base"`
	IRQ0Every  uint64 `json:"irq0_every"`
}
type Date struct {
	Year    uint16 `json:"year"`
	Month   uint8  `json:"month"`
	Day     uint8  `json:"day"`
	Weekday uint8  `json:"weekday"`
}
type Platform struct {
	DOSDate                        Date   `json:"dos_date"`
	DOSTimeModel                   string `json:"dos_time_model"`
	RTCReadTime                    bool   `json:"bios_rtc_read_time_supported"`
	RTCReadDate                    bool   `json:"bios_rtc_read_date_supported"`
	BDATicksObservation            string `json:"bios_bda_ticks_observation"`
	EffectiveMouseScaleObservation string `json:"effective_mouse_x_scale_observation"`
}
type Scalars struct {
	Schema   string             `json:"schema"`
	Step     uint64             `json:"step"`
	Ticks    uint64             `json:"ticks"`
	Clock    Clock              `json:"clock"`
	DOS      oracle.DOSSettings `json:"dos"`
	Platform Platform           `json:"platform"`
}
type Rational struct {
	Numerator   string `json:"numerator"`
	Denominator string `json:"denominator"`
}
type Audio struct {
	H                        Rational `json:"h"`
	K                        Rational `json:"initial_k"`
	ExplicitIPS              bool     `json:"explicit_ips"`
	IPSFloat64Bits           string   `json:"ips_float64_bits"`
	Mute                     bool     `json:"mute"`
	DeviceRequested          bool     `json:"device_requested"`
	SampleRate               uint32   `json:"sample_rate"`
	DeviceChannels           uint8    `json:"device_channels"`
	SampleBits               uint8    `json:"sample_bits"`
	RingFrames               uint32   `json:"ring_frames"`
	PrefillFrames            uint32   `json:"prefill_frames"`
	MaximumPCMChunkFrames    uint32   `json:"maximum_pcm_chunk_frames"`
	FadeFrames               uint32   `json:"fade_frames"`
	DeviceBufferMilliseconds uint32   `json:"device_buffer_milliseconds"`
	ProgressQuantum          uint64   `json:"progress_quantum"`
	ExitBudgetMilliseconds   uint32   `json:"exit_budget_milliseconds"`
}
type Resource struct {
	Role      string `json:"role"`
	Requested bool   `json:"requested"`
	Loaded    bool   `json:"loaded"`
	Size      uint64 `json:"size"`
	SHA256    string `json:"sha256"`
}
type Frontend struct {
	Scale              int64      `json:"scale"`
	HDComposeDisabled  bool       `json:"hd_compose_disabled"`
	AutoplayDisabled   bool       `json:"autoplay_disabled"`
	TestInputsDisabled bool       `json:"test_inputs_disabled"`
	Resources          []Resource `json:"resources"`
}
type SelectedProgram struct {
	SourceRole string `json:"source_role"`
	File       File   `json:"file"`
}
type Storage struct {
	Resolution    string `json:"resolution"`
	MetadataModel string `json:"metadata_model"`
	Root          []File `json:"root"`
	Scratch       []File `json:"scratch"`
}
type Config struct {
	Schema          string          `json:"schema"`
	Program         Program         `json:"program"`
	Sources         []Source        `json:"sources"`
	LoadWith        LoadOptions     `json:"load_with"`
	Initial         Scalars         `json:"initial_scalars"`
	SelectedProgram SelectedProgram `json:"selected_program"`
	Storage         Storage         `json:"storage"`
	Audio           Audio           `json:"audio"`
	Frontend        Frontend        `json:"frontend"`
}

func Hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func hashOK(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
func relativeOK(s string) bool {
	if s == "" || len(s) > 260 || !utf8.ValidString(s) || path.IsAbs(s) || path.Clean(s) != s || s == "." || s == ".." || strings.HasPrefix(s, "../") || strings.ContainsAny(s, "\\:") {
		return false
	}
	for _, r := range s {
		if r < 0x20 {
			return false
		}
	}
	return true
}
func fixedPlatform() Platform {
	return Platform{Date{1993, 1, 1, 5}, "now_plus_ticks_5493_div_1000_mod_8640000", false, false, "not_observed", "not_observed"}
}

func rat(r *big.Rat) Rational { return Rational{r.Num().String(), r.Denom().String()} }
func AudioConfig(ips float64, explicit, mute bool) (Audio, error) {
	if math.IsNaN(ips) || math.IsInf(ips, 0) || ips <= 0 {
		return Audio{}, errors.New("invalid_ips")
	}
	h := big.NewRat(196875030000, 65536)
	k := big.NewRat(1, 1)
	if explicit {
		k.Quo(new(big.Rat).SetFloat64(ips), h)
	}
	return Audio{rat(h), rat(k), explicit, fmt.Sprintf("%016x", math.Float64bits(ips)), mute, !mute, 44100, 2, 16, 88200, 4410, 4410, 221, 50, oracle.OPLProgressQuantum, 5000}, nil
}
func CollectFiles(dir string) ([]File, error) {
	info, e := os.Lstat(dir)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("initial_directory")
	}
	files := make([]File, 0)
	e = filepath.WalkDir(dir, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return errors.New("initial_inventory")
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("initial_symlink")
		}
		if d.IsDir() {
			return nil
		}
		info, e := d.Info()
		if e != nil || !info.Mode().IsRegular() {
			return errors.New("initial_nonregular")
		}
		rel, e := filepath.Rel(dir, p)
		if e != nil || !relativeOK(filepath.ToSlash(rel)) {
			return errors.New("initial_path")
		}
		size, digest, e := fileHash(p)
		if e != nil || info.Size() < 0 || size != uint64(info.Size()) {
			return errors.New("initial_read")
		}
		m := info.ModTime()
		y := m.Year() - 1980
		if y < 0 {
			y = 0
		}
		files = append(files, File{filepath.ToSlash(rel), size, digest, uint16(y)<<9 | uint16(m.Month())<<5 | uint16(m.Day()), uint16(m.Hour())<<11 | uint16(m.Minute())<<5 | uint16(m.Second()/2), 0x20})
		if len(files) > 4096 {
			return errors.New("initial_capacity")
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}
func validFiles(files []File, empty bool) bool {
	if files == nil || len(files) > 4096 || !empty && len(files) == 0 {
		return false
	}
	prev := ""
	for _, f := range files {
		if !relativeOK(f.Path) || !hashOK(f.SHA256) || f.Path <= prev || f.Attributes != 0x20 || f.DOSDate&31 == 0 || (f.DOSDate>>5)&15 == 0 || (f.DOSDate>>5)&15 > 12 || f.DOSTime>>11 > 23 || (f.DOSTime>>5)&63 > 59 || f.DOSTime&31 > 29 {
			return false
		}
		prev = f.Path
	}
	return true
}
func Validate(c Config) (Config, []byte, string, error) {
	fail := func(code string) (Config, []byte, string, error) { return Config{}, nil, "", errors.New(code) }
	if c.Initial.Step != 0 || c.Initial.Ticks != 0 {
		return fail("not_pristine_initial_boundary")
	}
	if c.Schema != "wiz7-startup-config/1" || c.Initial.Schema != "wiz7-initial-scalars/1" {
		return fail("initial_schema")
	}
	if c.Initial.Clock.CycleClock || c.Initial.Clock.CPUHz == 0 || c.Initial.Clock.IRQ0Base != 165000 || c.Initial.Clock.IRQ0Every != 165000 || c.Initial.Platform != fixedPlatform() {
		return fail("initial_machine")
	}
	if !hashOK(c.Program.SHA256) || !hashOK(c.Program.BuildSettingsSHA256) || c.Program.GoVersion != "go1.24.13" || c.Program.GOOS == "" || c.Program.GOARCH == "" || c.Program.ModuleVersion == "" {
		return fail("initial_program")
	}
	if c.Sources == nil || len(c.Sources) == 0 || len(c.Sources) > 4096 {
		return fail("initial_sources")
	}
	previous := ""
	for _, s := range c.Sources {
		if !relativeOK(s.Path) || s.Path <= previous || !hashOK(s.SHA256) {
			return fail("initial_sources")
		}
		previous = s.Path
	}
	if !validFiles(c.Storage.Root, false) || !validFiles(c.Storage.Scratch, true) || c.Storage.Resolution != "scratch_then_root" || c.Storage.MetadataModel != "dos_archive_mtime_source_policy" {
		return fail("initial_storage")
	}
	if !validFiles([]File{c.SelectedProgram.File}, false) {
		return fail("initial_selected_program")
	}
	switch c.SelectedProgram.SourceRole {
	case "root":
		found := false
		for _, f := range c.Storage.Root {
			if f == c.SelectedProgram.File {
				found = true
			}
		}
		if !found {
			return fail("selected_program_not_in_root")
		}
	case "external":
		if path.Base(c.SelectedProgram.File.Path) != c.SelectedProgram.File.Path {
			return fail("initial_selected_program")
		}
	default:
		return fail("initial_selected_program")
	}
	bits, e := hex.DecodeString(c.Audio.IPSFloat64Bits)
	if e != nil || len(bits) != 8 || c.Audio.IPSFloat64Bits != strings.ToLower(c.Audio.IPSFloat64Bits) {
		return fail("initial_audio")
	}
	var word uint64
	for _, b := range bits {
		word = word<<8 | uint64(b)
	}
	expected, e := AudioConfig(math.Float64frombits(word), c.Audio.ExplicitIPS, c.Audio.Mute)
	if e != nil || expected != c.Audio || !c.LoadWith.Adlib {
		return fail("initial_audio")
	}
	// 正式flag.Int沒有2..6限制；只記錄原值，不由候選新增操作規則。
	if !c.Frontend.HDComposeDisabled || !c.Frontend.AutoplayDisabled || !c.Frontend.TestInputsDisabled || len(c.Frontend.Resources) != 6 {
		return fail("initial_frontend")
	}
	for i, role := range []string{"xlate", "flow", "help", "roster", "profile-sex", "intro-picture"} {
		r := c.Frontend.Resources[i]
		if r.Role != role || r.Loaded && !r.Requested || r.Loaded && !hashOK(r.SHA256) || r.SHA256 != "" && !hashOK(r.SHA256) || r.SHA256 == "" && r.Size != 0 || !r.Requested && (r.Loaded || r.Size != 0 || r.SHA256 != "") {
			return fail("initial_resource")
		}
	}
	b, e := json.Marshal(c)
	if e != nil {
		return fail("initial_marshal")
	}
	if len(b) > 1<<20 {
		return fail("initial_capacity")
	}
	// 複製所有slice，避免擷取後呼叫端改動而造成header與已算SHA分離。
	c.Sources = append([]Source{}, c.Sources...)
	c.Storage.Root = append([]File{}, c.Storage.Root...)
	c.Storage.Scratch = append([]File{}, c.Storage.Scratch...)
	c.Frontend.Resources = append([]Resource{}, c.Frontend.Resources...)
	return c, b, Hash(b), nil
}

// 下列掃描要求呼叫端獨占、不變動的起始目錄；並非抵抗並行換檔的安全邊界。
func fileHash(p string) (uint64, string, error) {
	info, e := os.Lstat(p)
	if e != nil || !info.Mode().IsRegular() {
		return 0, "", errors.New("initial_nonregular")
	}
	f, e := os.Open(p)
	if e != nil {
		return 0, "", errors.New("initial_read")
	}
	h := sha256.New()
	n, e := io.Copy(h, f)
	ce := f.Close()
	if e != nil || ce != nil || n != info.Size() {
		return 0, "", errors.New("initial_read")
	}
	return uint64(n), hex.EncodeToString(h.Sum(nil)), nil
}
func CollectSources(base string, paths []string) ([]Source, error) {
	result := make([]Source, 0, len(paths))
	for _, p := range paths {
		if !relativeOK(p) {
			return nil, errors.New("initial_sources")
		}
		current := base
		parts := strings.Split(p, "/")
		for _, part := range parts[:len(parts)-1] {
			current = filepath.Join(current, part)
			info, e := os.Lstat(current)
			if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil, errors.New("initial_sources")
			}
		}
		info, e := os.Lstat(base)
		if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, errors.New("initial_sources")
		}
		n, h, e := fileHash(filepath.Join(base, filepath.FromSlash(p)))
		if e != nil {
			return nil, e
		}
		result = append(result, Source{p, n, h})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	for i := 1; i < len(result); i++ {
		if result[i-1].Path == result[i].Path {
			return nil, errors.New("initial_sources")
		}
	}
	if len(result) == 0 || len(result) > 4096 {
		return nil, errors.New("initial_sources")
	}
	return result, nil
}
func CollectProgram(p string) (Program, error) {
	_, h, e := fileHash(p)
	if e != nil {
		return Program{}, e
	}
	bi, e := buildinfo.ReadFile(p)
	if e != nil {
		return Program{}, errors.New("initial_program_metadata")
	}
	type setting struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	settings := make([]setting, 0, len(bi.Settings))
	goos, goarch := "", ""
	for _, s := range bi.Settings {
		settings = append(settings, setting{s.Key, s.Value})
		if s.Key == "GOOS" {
			goos = s.Value
		}
		if s.Key == "GOARCH" {
			goarch = s.Value
		}
	}
	sort.Slice(settings, func(i, j int) bool { return settings[i].Key < settings[j].Key })
	b, e := json.Marshal(settings)
	if e != nil || goos == "" || goarch == "" || bi.Main.Version == "" {
		return Program{}, errors.New("initial_program_metadata")
	}
	return Program{h, bi.GoVersion, goos, goarch, bi.Main.Version, Hash(b)}, nil
}

// 重播啟動前，先核對兩層的完整集合；新增、刪除、內容或DOS中繼資料差異皆拒絕。
func VerifyStorage(expected Storage, root, scratch string) error {
	if expected.Resolution != "scratch_then_root" || expected.MetadataModel != "dos_archive_mtime_source_policy" || !validFiles(expected.Root, false) || !validFiles(expected.Scratch, true) {
		return errors.New("initial_storage")
	}
	r, e := CollectFiles(root)
	if e != nil {
		return e
	}
	s, e := CollectFiles(scratch)
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(r, expected.Root) || !reflect.DeepEqual(s, expected.Scratch) {
		return errors.New("initial_storage_mismatch")
	}
	return nil
}
