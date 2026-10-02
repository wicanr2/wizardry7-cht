//go:build linux || darwin

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

type saveFixtureResult struct {
	Name       string `json:"name"`
	Code       string `json:"saveErrorCode"`
	ProbeCalls int    `json:"probe_calls"`
	Expected   string `json:"expected"`
	Pass       bool   `json:"pass"`
}

func mustSaveFixture(e error) {
	if e != nil {
		panic(e)
	}
}
func saveFixtureSHA(p string) string {
	b, e := os.ReadFile(p)
	mustSaveFixture(e)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func TestReviewedStorageGuard(t *testing.T) {
	base := filepath.Join(t.TempDir(), "fixtures")
	mustSaveFixture(os.MkdirAll(base, 0700))
	root := filepath.Join(base, "pristine")
	mustSaveFixture(os.Mkdir(root, 0700))
	exe := filepath.Join(root, "DS.EXE")
	mustSaveFixture(os.WriteFile(exe, []byte("self-authored pristine executable fixture"), 0600))
	before := saveFixtureSHA(exe)
	results := []saveFixtureResult{}
	test := func(name, save, want string, wp func(string) error) {
		calls := 0
		_, e := prepareSaveStorage(root, exe, save, func(s string) error {
			calls++
			if wp != nil {
				return wp(s)
			}
			return probeSaveStorage(s)
		})
		c := saveErrorCode(e)
		results = append(results, saveFixtureResult{name, c, calls, want, c == want})
	}
	test("empty", "", "explicit_save_required", nil)
	test("root", root, "overlap", nil)
	test("root_parent", base, "overlap", nil)
	test("root_child", filepath.Join(root, "new"), "overlap", nil)
	test("exe", exe, "overlap", nil)
	alias := filepath.Join(base, "alias")
	mustSaveFixture(os.Symlink(root, alias))
	test("symlink_root", alias, "overlap", nil)
	test("symlink_missing_child", filepath.Join(alias, "new"), "overlap", nil)
	plain := filepath.Join(base, "regular")
	mustSaveFixture(os.WriteFile(plain, []byte("fixture"), 0600))
	test("not_directory", plain, "mkdir", nil)
	mk := func(n string) string { p := filepath.Join(base, n); mustSaveFixture(os.Mkdir(p, 0700)); return p }
	sd := mk("entry_symlink")
	mustSaveFixture(os.Symlink(exe, filepath.Join(sd, "DS.EXE")))
	test("entry_symlink", sd, "nonregular_entry", nil)
	hd := mk("entry_hardlink")
	mustSaveFixture(os.Link(exe, filepath.Join(hd, "DS.EXE")))
	test("entry_hardlink", hd, "root_file_alias", nil)
	ext := filepath.Join(base, "other_source")
	mustSaveFixture(os.WriteFile(ext, []byte("self other source"), 0600))
	otherBefore := saveFixtureSHA(ext)
	hd2 := mk("other_hardlink")
	mustSaveFixture(os.Link(ext, filepath.Join(hd2, "OTHER.DAT")))
	test("other_hardlink", hd2, "multiple_hardlinks", nil)
	nd := mk("nested_directory")
	mustSaveFixture(os.Mkdir(filepath.Join(nd, "sub"), 0700))
	test("nested_directory", nd, "nonregular_entry", nil)
	fd := mk("fifo")
	mustSaveFixture(syscall.Mkfifo(filepath.Join(fd, "pipe"), 0600))
	test("fifo", fd, "nonregular_entry", nil)
	ro := mk("read_only_dir")
	mustSaveFixture(os.Chmod(ro, 0500))
	test("directory_unwritable", ro, "directory_not_writable", nil)
	rf := mk("read_only_file")
	mustSaveFixture(os.WriteFile(filepath.Join(rf, "SAVEGAME.DBS"), []byte("self save fixture"), 0400))
	test("file_unwritable", rf, "entry_not_writable", nil)
	fail := mk("write_fail")
	test("initial_write_failed", fail, "initial_write_failed", func(s string) error {
		f, e := os.Open(exe)
		if e != nil {
			return e
		}
		defer f.Close()
		_, e = f.Write([]byte("blocked"))
		return e
	})
	good := filepath.Join(base, "new_save")
	test("fresh", good, "ok", nil)
	mustSaveFixture(os.WriteFile(filepath.Join(good, "SAVEGAME.DBS"), []byte("self save fixture"), 0600))
	test("existing_regular", good, "ok", nil)
	exeAlias := filepath.Join(base, "exe_alias")
	mustSaveFixture(os.Link(exe, exeAlias))
	test("exe_path_hardlink", exeAlias, "overlap", nil)
	wo := mk("write_only_file")
	mustSaveFixture(os.WriteFile(filepath.Join(wo, "SAVEGAME.DBS"), []byte("self write-only fixture"), 0200))
	test("file_readwrite_denied", wo, "entry_not_writable", nil)
	// Unprotected OpenFile controls only touch separate self-authored replicas.
	controls := []map[string]any{}
	for _, kind := range []string{"symlink", "hardlink"} {
		d := mk("unsafe_" + kind)
		p := filepath.Join(d, "source")
		target := filepath.Join(d, "target")
		mustSaveFixture(os.WriteFile(p, []byte("self unsafe control"), 0600))
		a := saveFixtureSHA(p)
		if kind == "symlink" {
			mustSaveFixture(os.Symlink(p, target))
		} else {
			mustSaveFixture(os.Link(p, target))
		}
		f, e := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		mustSaveFixture(e)
		_, e = f.Write([]byte("changed"))
		mustSaveFixture(e)
		mustSaveFixture(f.Close())
		controls = append(controls, map[string]any{"kind": kind, "before": a, "after": saveFixtureSHA(p), "source_changed": a != saveFixtureSHA(p)})
	}
	pass := before == saveFixtureSHA(exe) && otherBefore == saveFixtureSHA(ext)
	for _, r := range results {
		pass = pass && r.Pass
	}
	summary := map[string]any{"pass": pass, "cases": results, "pristine_before": before, "pristine_after": saveFixtureSHA(exe), "other_source_before": otherBefore, "other_source_after": saveFixtureSHA(ext), "unsafe_OpenFile_controls": controls, "original_game_run": false, "default_save_location_selected": false}
	b, e := json.MarshalIndent(summary, "", "  ")
	mustSaveFixture(e)
	t.Log(string(b))
	if !pass {
		t.Fatal("guard fixture failed")
	}
}
