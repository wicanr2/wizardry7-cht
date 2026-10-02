package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type saveRejection struct{ Code string }

func (e *saveRejection) Error() string             { return e.Code }
func rejectSaveStorage(saveErrorCode string) error { return &saveRejection{saveErrorCode} }
func savePathAncestor(a, b string) bool {
	r, e := filepath.Rel(a, b)
	return e == nil && (r == "." || r != ".." && !strings.HasPrefix(r, ".."+string(os.PathSeparator)) && !filepath.IsAbs(r))
}

// Canonicalize missing leaves using the nearest existing savePathAncestor before mkdir.
func canonicalSavePath(p string) (string, error) {
	a, e := filepath.Abs(p)
	if e != nil {
		return "", e
	}
	a = filepath.Clean(a)
	tail := []string{}
	for {
		_, e = os.Lstat(a)
		if e == nil {
			break
		}
		if !os.IsNotExist(e) {
			return "", e
		}
		b := filepath.Dir(a)
		if b == a {
			return "", e
		}
		tail = append(tail, filepath.Base(a))
		a = b
	}
	a, e = filepath.EvalSymlinks(a)
	if e != nil {
		return "", e
	}
	for i := len(tail) - 1; i >= 0; i-- {
		a = filepath.Join(a, tail[i])
	}
	return a, nil
}
func probeSaveStorage(dir string) error {
	return probeSaveStorageCleanup(dir, os.Remove)
}
func probeSaveStorageCleanup(dir string, remove func(string) error) (result error) {
	f, e := os.CreateTemp(dir, ".r807-probeSaveStorage-")
	if e != nil {
		return e
	}
	name := f.Name()
	defer func() {
		if err := remove(name); result == nil && err != nil {
			result = err
		}
	}()
	if _, e = f.Write([]byte{1}); e != nil {
		f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	return f.Close()
}
func prepareSaveStorage(root, exe, save string, writeProbe func(string) error) (string, error) {
	if save == "" {
		return "", rejectSaveStorage("explicit_save_required")
	}
	r, e := canonicalSavePath(root)
	if e != nil {
		return "", rejectSaveStorage("root_path")
	}
	x, e := canonicalSavePath(exe)
	if e != nil {
		return "", rejectSaveStorage("exe_path")
	}
	s, e := canonicalSavePath(save)
	if e != nil {
		return "", rejectSaveStorage("save_path")
	}
	if savePathAncestor(s, r) || savePathAncestor(r, s) || savePathAncestor(s, x) || savePathAncestor(x, s) {
		return "", rejectSaveStorage("overlap")
	}
	// Physical savePathAncestor identity also covers case aliases on native filesystems.
	for _, pair := range [][2]string{{r, s}, {s, r}, {x, s}, {s, x}} {
		yes, err := savePhysicalAncestor(pair[0], pair[1])
		if err != nil {
			return "", rejectSaveStorage("save_path")
		}
		if yes {
			return "", rejectSaveStorage("overlap")
		}
	}
	// Existing alias components are canonicalized above; every save entry is Lstat checked.
	if e = os.MkdirAll(s, 0700); e != nil {
		return "", rejectSaveStorage("mkdir")
	}
	st, e := os.Lstat(s)
	if e != nil || !st.IsDir() {
		return "", rejectSaveStorage("not_directory")
	}
	dm, err := saveMetadataForPath(s, st)
	if err != nil {
		return "", rejectSaveStorage("platform_metadata_unknown")
	}
	if dm.Reparse {
		return "", rejectSaveStorage("not_directory")
	}
	if !dm.Owned {
		return "", rejectSaveStorage("directory_owner")
	}
	rootFiles := []os.FileInfo{}
	e = filepath.Walk(r, func(_ string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if i.Mode().IsRegular() {
			rootFiles = append(rootFiles, i)
		}
		return nil
	})
	if e != nil {
		return "", rejectSaveStorage("root_inventory")
	}
	xi, e := os.Stat(x)
	if e != nil {
		return "", rejectSaveStorage("exe_stat")
	}
	rootFiles = append(rootFiles, xi)
	entries, e := os.ReadDir(s)
	if e != nil {
		return "", rejectSaveStorage("entries")
	}
	for _, en := range entries {
		p := filepath.Join(s, en.Name())
		i, e := os.Lstat(p)
		if e != nil {
			return "", rejectSaveStorage("entry_stat")
		}
		if !i.Mode().IsRegular() {
			return "", rejectSaveStorage("nonregular_entry")
		}
		fm, err := saveMetadataForPath(p, i)
		if err != nil {
			return "", rejectSaveStorage("platform_metadata_unknown")
		}
		if fm.Reparse {
			return "", rejectSaveStorage("nonregular_entry")
		}
		if !fm.Owned {
			return "", rejectSaveStorage("entry_owner")
		}
		for _, j := range rootFiles {
			if os.SameFile(i, j) {
				return "", rejectSaveStorage("root_file_alias")
			}
		}
		if fm.Nlink != 1 {
			return "", rejectSaveStorage("multiple_hardlinks")
		}
		if !savePathWritable(p, i) {
			return "", rejectSaveStorage("entry_not_writable")
		}
	}
	if !savePathWritable(s, st) {
		return "", rejectSaveStorage("directory_not_writable")
	}
	if writeProbe == nil {
		writeProbe = probeSaveStorage
	}
	if e = writeProbe(s); e != nil {
		return "", rejectSaveStorage("initial_write_failed")
	}
	return s, nil
}
func saveErrorCode(e error) string {
	if e == nil {
		return "ok"
	}
	var r *saveRejection
	if errors.As(e, &r) {
		return r.Code
	}
	return "unclassified"
}

type saveMetadata struct {
	Owned   bool
	Nlink   uint64
	Reparse bool
}

func savePhysicalAncestor(a, b string) (bool, error) {
	ai, err := os.Stat(a)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	for {
		bi, e := os.Stat(b)
		if e == nil && os.SameFile(ai, bi) {
			return true, nil
		}
		if e != nil && !os.IsNotExist(e) {
			return false, e
		}
		parent := filepath.Dir(b)
		if parent == b {
			return false, nil
		}
		b = parent
	}
}
