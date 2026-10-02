package main

import (
	"fmt"

	"path/filepath"
)

func resolveSavePath(explicit string) (string, error) {
	if explicit != "" {
		return filepath.Abs(explicit)
	}
	return defaultSavePath()
}
func applicationSaveDir(base string) (string, error) {
	if base == "" || !filepath.IsAbs(base) {
		return "", fmt.Errorf("default_save_path")
	}
	return filepath.Join(base, "wizardry7-cht", "save"), nil
}
