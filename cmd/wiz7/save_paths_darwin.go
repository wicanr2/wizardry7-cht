package main

import (
	"os"
	"path/filepath"
)

func defaultSavePath() (string, error) {
	h, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	if h == "" || !filepath.IsAbs(h) {
		return applicationSaveDir("")
	}
	return applicationSaveDir(filepath.Join(h, "Library", "Application Support"))
}
