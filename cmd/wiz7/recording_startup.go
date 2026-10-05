package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
)

// 建置器以 -X 注入同一 binary 的實際選定来源清單。普通 go build 可遊玩，
// 缺少這份可重現身份時拒絕錄放，不以執行期工作目錄猜測編譯來源。
var compiledRecordingSources, compiledRecordingSourceSHA string

type recordingStartup struct {
	root, program, scratch string
	scale                  int64
	ips                    float64
	explicitIPS, mute      bool
	resources              []ic.Resource
}
type recordingOwner struct {
	oracle   *oracle.Oracle
	config   []byte
	sha      string
	released bool
}
type recordingOriginal struct {
	Step, Ticks                      uint64
	RAM, Indexed, Palette, Registers string
}

func recordingOriginalHash(o *oracle.Oracle) recordingOriginal {
	palette := o.Palette()
	b := make([]byte, 0, len(palette)*3)
	for _, c := range palette {
		b = append(b, c[0], c[1], c[2])
	}
	registers, _ := json.Marshal(o.Regs())
	return recordingOriginal{o.Steps(), o.Ticks(), ic.Hash(o.Bytes(oracle.Far(0, 0), 1<<20)), ic.Hash(o.Indexed()), ic.Hash(b), ic.Hash(registers)}
}
func recordingSelectedFile(name, relative string) (ic.File, error) {
	items, err := ic.CollectSources(filepath.Dir(name), []string{filepath.Base(name)})
	if err != nil {
		return ic.File{}, err
	}
	info, err := os.Lstat(name)
	if err != nil {
		return ic.File{}, errors.New("selected_program")
	}
	m := info.ModTime()
	year := max(m.Year()-1980, 0)
	return ic.File{Path: relative, Size: items[0].Size, SHA256: items[0].SHA256, DOSDate: uint16(year)<<9 | uint16(m.Month())<<5 | uint16(m.Day()), DOSTime: uint16(m.Hour())<<11 | uint16(m.Minute())<<5 | uint16(m.Second()/2), Attributes: 32}, nil
}
func recordingSources() ([]ic.Source, error) {
	b, err := base64.StdEncoding.DecodeString(compiledRecordingSources)
	if err != nil || compiledRecordingSourceSHA == "" || ic.Hash(b) != compiledRecordingSourceSHA {
		return nil, errors.New("compiled_sources")
	}
	var sources []ic.Source
	if json.Unmarshal(b, &sources) != nil {
		return nil, errors.New("compiled_sources")
	}
	if len(sources) != len(recordingSourcePaths) {
		return nil, errors.New("compiled_sources")
	}
	for i, s := range sources {
		if s.Path != recordingSourcePaths[i] {
			return nil, errors.New("compiled_sources")
		}
	}
	return sources, nil
}

// 唯一建構入口。從 LoadWith 至凍結設定期間，Oracle 指標不會交出。
func newRecordingOwner(s recordingStartup, expected *ic.Config) (owner *recordingOwner, err error) {
	sources, err := recordingSources()
	if err != nil {
		return nil, err
	}
	binary, err := os.Executable()
	if err != nil {
		return nil, errors.New("binary_identity")
	}
	programIdentity, err := ic.CollectProgram(binary)
	if err != nil {
		return nil, err
	}
	rootFiles, err := ic.CollectFiles(s.root)
	if err != nil {
		return nil, err
	}
	scratchFiles, err := ic.CollectFiles(s.scratch)
	if err != nil {
		return nil, err
	}
	absRoot, err := filepath.Abs(s.root)
	if err != nil {
		return nil, errors.New("root_identity")
	}
	absProgram, err := filepath.Abs(s.program)
	if err != nil {
		return nil, errors.New("selected_program")
	}
	relative, err := filepath.Rel(absRoot, absProgram)
	if err != nil {
		return nil, errors.New("selected_program")
	}
	role, name := "external", filepath.Base(absProgram)
	if relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		role, name = "root", filepath.ToSlash(relative)
	}
	selected, err := recordingSelectedFile(absProgram, name)
	if err != nil {
		return nil, err
	}
	audio, err := ic.AudioConfig(s.ips, s.explicitIPS, s.mute)
	if err != nil {
		return nil, err
	}
	options := oracle.Options{Adlib: true}
	c := ic.Config{Schema: "wiz7-startup-config/1", Program: programIdentity, Sources: sources, LoadWith: ic.LoadOptions{Adlib: options.Adlib, DGROUP: options.DGROUP, MouseXScale: options.MouseXScale, FontFull: options.FontFull, FontHalf: options.FontHalf}, SelectedProgram: ic.SelectedProgram{SourceRole: role, File: selected}, Storage: ic.Storage{Resolution: "scratch_then_root", MetadataModel: "dos_archive_mtime_source_policy", Root: rootFiles, Scratch: scratchFiles}, Audio: audio, Frontend: ic.Frontend{Scale: s.scale, HDComposeDisabled: true, AutoplayDisabled: true, TestInputsDisabled: true, Resources: append([]ic.Resource{}, s.resources...)}}
	if expected != nil {
		if _, _, _, err := ic.Validate(*expected); err != nil {
			return nil, errors.New("expected_configuration")
		}
		static := *expected
		static.Initial = ic.Scalars{}
		if !reflect.DeepEqual(c, static) {
			return nil, errors.New("configuration_mismatch")
		}
	}
	o, err := oracle.LoadWith(absProgram, absRoot, options)
	if err != nil {
		return nil, errors.New("initialization_failed")
	}
	o.SetScratch(s.scratch)
	keep := false
	defer func() {
		if !keep {
			o.Close()
		}
	}()
	before := recordingOriginalHash(o)
	if o.Steps() != 0 || o.Ticks() != 0 {
		return nil, errors.New("initial_boundary")
	}
	clock := o.ClockSettings()
	c.Initial = ic.Scalars{Schema: "wiz7-initial-scalars/1", Step: o.Steps(), Ticks: o.Ticks(), Clock: ic.Clock{CPUHz: clock.CPUHz, CycleClock: clock.CycleClock, IRQ0Base: clock.IRQ0Base, IRQ0Every: clock.IRQ0Every}, DOS: o.DOSSettings(), Platform: ic.Platform{DOSDate: ic.Date{Year: 1993, Month: 1, Day: 1, Weekday: 5}, DOSTimeModel: "now_plus_ticks_5493_div_1000_mod_8640000", BDATicksObservation: "not_observed", EffectiveMouseScaleObservation: "not_observed"}}
	if err := ic.VerifyStorage(c.Storage, s.root, s.scratch); err != nil {
		return nil, err
	}
	selectedAfter, err := recordingSelectedFile(absProgram, name)
	if err != nil || selectedAfter != selected {
		return nil, errors.New("selected_program_changed")
	}
	if expected != nil && !reflect.DeepEqual(c.Initial, expected.Initial) {
		return nil, errors.New("configuration_mismatch")
	}
	_, encoded, hash, err := ic.Validate(c)
	if err != nil {
		return nil, err
	}
	if before != recordingOriginalHash(o) {
		return nil, errors.New("capture_changed_original")
	}
	owner = &recordingOwner{oracle: o, config: encoded, sha: hash}
	keep = true
	return owner, nil
}
func (o *recordingOwner) Config() ([]byte, string, error) {
	if o == nil || o.released || o.oracle == nil {
		return nil, "", errors.New("configuration_not_owned")
	}
	return append([]byte{}, o.config...), o.sha, nil
}
func (o *recordingOwner) Take() (*oracle.Oracle, error) {
	if o == nil || o.released || o.oracle == nil {
		return nil, errors.New("owner_already_released")
	}
	p := o.oracle
	o.oracle = nil
	o.released = true
	return p, nil
}
func (o *recordingOwner) Close() {
	if o != nil && o.oracle != nil {
		o.oracle.Close()
		o.oracle = nil
		o.released = true
	}
}
