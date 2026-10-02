//go:build windows

package main

import (
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"runtime"
	"unsafe"
)

// These handles query existing objects. They never create, truncate, or write.
func openMetadata(path string, access uint32) (windows.Handle, error) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	return windows.CreateFile(p, access, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
}
func saveMetadataForPath(path string, _ os.FileInfo) (saveMetadata, error) {
	h, err := openMetadata(path, windows.FILE_READ_ATTRIBUTES|windows.READ_CONTROL)
	if err != nil {
		return saveMetadata{}, err
	}
	defer windows.CloseHandle(h)
	var fi windows.ByHandleFileInformation
	if err = windows.GetFileInformationByHandle(h, &fi); err != nil {
		return saveMetadata{}, err
	}
	if fi.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return saveMetadata{Reparse: true}, nil
	}
	sd, err := windows.GetSecurityInfo(h, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		return saveMetadata{}, err
	}
	owner, _, err := sd.Owner()
	if err != nil || owner == nil {
		return saveMetadata{}, errors.New("owner unavailable")
	}
	token := windows.GetCurrentProcessToken()
	user, err := token.GetTokenUser()
	if err != nil {
		return saveMetadata{}, err
	}
	defaultOwner, err := tokenOwner(token)
	if err != nil {
		return saveMetadata{}, err
	}
	// TokenOwner defines Windows' owner for objects newly created by this token.
	owned := owner.Equals(user.User.Sid) || owner.Equals(defaultOwner)
	return saveMetadata{Owned: owned, Nlink: uint64(fi.NumberOfLinks)}, nil
}
func tokenOwner(token windows.Token) (*windows.SID, error) {
	var n uint32
	err := windows.GetTokenInformation(token, windows.TokenOwner, nil, 0, &n)
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, errors.New("token owner size unavailable")
	}
	if n < uint32(unsafe.Sizeof(uintptr(0))) || n > 65536 {
		return nil, errors.New("token owner size invalid")
	}
	b := make([]byte, n)
	if err = windows.GetTokenInformation(token, windows.TokenOwner, &b[0], uint32(len(b)), &n); err != nil {
		return nil, err
	}
	// TOKEN_OWNER contains one SID pointer into the returned buffer.
	owner := *(**windows.SID)(unsafe.Pointer(&b[0]))
	if owner == nil {
		return nil, errors.New("token owner unavailable")
	}
	copied, err := owner.Copy()
	runtime.KeepAlive(b)
	return copied, err
}
func savePathWritable(path string, info os.FileInfo) (result bool) {
	// Directory writability is measured by the existing create/write/sync probeSaveStorage.
	if info.IsDir() {
		return true
	}
	h, err := openMetadata(path, windows.GENERIC_READ|windows.GENERIC_WRITE)
	if err != nil {
		return false
	}
	defer func() {
		if windows.CloseHandle(h) != nil {
			result = false
		}
	}()
	var fi windows.ByHandleFileInformation
	if windows.GetFileInformationByHandle(h, &fi) != nil {
		return false
	}
	return fi.FileAttributes&(windows.FILE_ATTRIBUTE_READONLY|windows.FILE_ATTRIBUTE_REPARSE_POINT) == 0
}
