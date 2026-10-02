package main

import (
	"errors"
	"os"
	"testing"
)

func TestFailedProbeCleanupRefusesStorage(t *testing.T) {
	dir := t.TempDir()
	sentinel := errors.New("self-authored cleanup fault")
	calls := 0
	err := probeSaveStorageCleanup(dir, func(name string) error {
		calls++
		if e := os.Remove(name); e != nil {
			t.Fatal(e)
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) || calls != 1 {
		t.Fatalf("error=%v calls=%d", err, calls)
	}
	entries, e := os.ReadDir(dir)
	if e != nil || len(entries) != 0 {
		t.Fatalf("residue=%d error=%v", len(entries), e)
	}
}
