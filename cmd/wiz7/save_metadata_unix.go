//go:build linux || darwin

package main

import (
	"errors"
	"os"
	"syscall"
)

func saveMetadataForPath(_ string, info os.FileInfo) (saveMetadata, error) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return saveMetadata{}, errors.New("saveMetadataForPath unavailable")
	}
	return saveMetadata{Owned: st.Uid == uint32(os.Getuid()), Nlink: uint64(st.Nlink)}, nil
}
func savePathWritable(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return info.Mode().Perm()&0300 == 0300
	}
	if info.Mode().Perm()&0200 == 0 {
		return false
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return false
	}
	return f.Close() == nil
}
