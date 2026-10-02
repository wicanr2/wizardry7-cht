package main

import "golang.org/x/sys/windows"

func defaultSavePath() (string, error) {
	p, e := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	if e != nil {
		return "", e
	}
	return applicationSaveDir(p)
}
