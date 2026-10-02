package main

import (
	"os"
	"path/filepath"
)

func linuxSaveDir(xdg string, home func() (string, error)) (string, error) {
	if xdg != "" && filepath.IsAbs(xdg) {
		return applicationSaveDir(xdg)
	}
	h, e := home()
	if e != nil {
		return "", e
	}
	if h == "" || !filepath.IsAbs(h) {
		return applicationSaveDir("")
	}
	return applicationSaveDir(filepath.Join(h, ".local", "share"))
}
func defaultSavePath() (string, error) {
	return linuxSaveDir(os.Getenv("XDG_DATA_HOME"), os.UserHomeDir)
}
