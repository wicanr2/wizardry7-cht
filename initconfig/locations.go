package initconfig

import (
	"errors"
	"os"
	"path/filepath"
)

// CheckLocation 先核實際父目錄身份，再允許建立或讀取錄放檔。
// 呼叫端須獨占不變動的起始輸入；此處不承諾抵抗外部同時換檔。
func CheckLocation(name, root, scratch string, protected []string, output bool) error {
	fail := func() error { return errors.New("recording_location") }
	name, err := filepath.Abs(name)
	if err != nil {
		return fail()
	}
	var roots []os.FileInfo
	for _, path := range []string{root, scratch} {
		info, err := os.Lstat(path)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fail()
		}
		roots = append(roots, info)
	}
	for current := filepath.Dir(name); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fail()
		}
		for _, rootInfo := range roots {
			if os.SameFile(info, rootInfo) {
				return fail()
			}
		}
		if filepath.Dir(current) == current {
			break
		}
	}
	info, err := os.Lstat(name)
	if output {
		if !os.IsNotExist(err) {
			return fail()
		}
		if _, err := os.Lstat(name + ".partial"); !os.IsNotExist(err) {
			return fail()
		}
		return nil
	}
	if err != nil || !info.Mode().IsRegular() {
		return fail()
	}
	for _, path := range protected {
		if path == "" {
			continue
		}
		candidate, err := os.Stat(path)
		if err == nil && os.SameFile(info, candidate) {
			return fail()
		}
	}
	return nil
}
