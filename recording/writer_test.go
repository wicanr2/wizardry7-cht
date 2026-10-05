package recording

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

type faultFile struct {
	*os.File
	syncCalls, failSync int
	failClose           bool
	shortWrite          bool
	entered, release    chan struct{}
}

func (f *faultFile) Write(b []byte) (int, error) {
	if f.shortWrite {
		f.shortWrite = false
		n, _ := f.File.Write(b[:len(b)/2])
		return n, errors.New("self-authored write fault")
	}
	return f.File.Write(b)
}
func (f *faultFile) Sync() error {
	f.syncCalls++
	if f.entered != nil && f.syncCalls == 1 {
		close(f.entered)
		<-f.release
	}
	if f.syncCalls == f.failSync {
		return errors.New("self-authored sync fault")
	}
	return f.File.Sync()
}
func (f *faultFile) Close() error {
	err := f.File.Close()
	if f.failClose {
		return errors.New("self-authored close fault")
	}
	return err
}

func readyWriter(t *testing.T) (*Writer, string) {
	t.Helper()
	h := fixture(t)
	name := filepath.Join(t.TempDir(), "input.jsonl")
	w, err := New(name, h)
	if err != nil {
		t.Fatal(err)
	}
	pub := uint64(1)
	if err := w.Append(Event{"snapshot", 1, 0, 0, "none", Params{Publication: &pub, FrontendSHA: &h.ConfigSHA, Panel: &Panel{}}}); err != nil {
		t.Fatal(err)
	}
	return w, name
}

func assertIncomplete(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Lstat(name); !os.IsNotExist(err) {
		t.Fatal("success name survived", err)
	}
	b, err := os.ReadFile(name + ".partial")
	if err != nil || bytes.Contains(b, []byte(`"kind":"end"`)) {
		t.Fatal("partial retained accepted end", err)
	}
	if _, err := Load(name + ".partial"); err == nil {
		t.Fatal("partial filename accepted")
	}
	copyName := filepath.Join(filepath.Dir(name), "renamed.jsonl")
	if err := os.WriteFile(copyName, b, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(copyName); err == nil {
		t.Fatal("renamed fragment accepted")
	}
}

func TestWriterFaultRollback(t *testing.T) {
	for _, point := range []string{"payload-sync", "end-sync", "close", "publish", "partial-cleanup", "post-publish-cancel", "post-cleanup-cancel", "deadline"} {
		t.Run(point, func(t *testing.T) {
			w, name := readyWriter(t)
			f := &faultFile{File: w.f.(*os.File)}
			w.f = f
			var cancelled atomic.Bool
			switch point {
			case "payload-sync":
				f.failSync = 1
			case "end-sync":
				f.failSync = 2
			case "close":
				f.failClose = true
			case "publish":
				w.link = func(string, string) error { return errors.New("self-authored link fault") }
			case "partial-cleanup":
				w.remove = func(p string) error {
					if p == w.partial {
						return errors.New("self-authored cleanup fault")
					}
					return os.Remove(p)
				}
			case "post-publish-cancel":
				w.link = func(a, b string) error { err := os.Link(a, b); cancelled.Store(true); return err }
			case "post-cleanup-cancel":
				w.remove = func(p string) error {
					err := os.Remove(p)
					if p == w.partial {
						cancelled.Store(true)
					}
					return err
				}
			case "deadline":
				w.link = func(a, b string) error {
					err := os.Link(a, b)
					w.finishDeadline = time.Now().Add(-time.Second)
					return err
				}
			}
			if err := w.FinishWhen(0, time.Now().Add(time.Second), func() bool { return !cancelled.Load() }); err == nil {
				t.Fatal("fault reported success")
			}
			assertIncomplete(t, name)
		})
	}
}

func TestBlockedWriterOwnsCleanup(t *testing.T) {
	w, name := readyWriter(t)
	entered, release := make(chan struct{}), make(chan struct{})
	w.f = &faultFile{File: w.f.(*os.File), entered: entered, release: release}
	var cancelled atomic.Bool
	done := make(chan error, 1)
	go func() { done <- w.FinishWhen(0, time.Now().Add(time.Second), func() bool { return !cancelled.Load() }) }()
	<-entered
	cancelled.Store(true)
	select {
	case <-done:
		t.Fatal("blocked syscall was falsely cancelled")
	default:
	}
	close(release)
	if err := <-done; err == nil {
		t.Fatal("cancelled writer succeeded")
	}
	assertIncomplete(t, name)
}

func TestCleanupFailureIsExplicitAndReplacementUntouched(t *testing.T) {
	w, name := readyWriter(t)
	var cancelled atomic.Bool
	w.link = func(a, b string) error { err := os.Link(a, b); cancelled.Store(true); return err }
	w.remove = func(string) error { return errors.New("self-authored cleanup fault") }
	if err := w.FinishWhen(0, time.Now().Add(time.Second), func() bool { return !cancelled.Load() }); err == nil || err.Error() != "cleanup_incomplete" {
		t.Fatal("cleanup failure hidden", err)
	}
	// 即使同 inode 的成功名稱未能移除，內容也已截去 end，改名不會被接受。
	if _, err := Load(name); err == nil {
		t.Fatal("unremoved invalid file accepted")
	}
	w.remove = nil
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
	replacement := []byte("self-authored unrelated file")
	if err := os.WriteFile(name, replacement, 0600); err != nil {
		t.Fatal(err)
	}
	if err := w.Revoke(); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(name)
	if err != nil || !bytes.Equal(actual, replacement) {
		t.Fatal("another inode was changed", err)
	}
}

func TestWriterRequiresSuccessfulFinalSnapshot(t *testing.T) {
	for _, scenario := range []string{"empty", "last-key", "api-error", "short-end-write"} {
		t.Run(scenario, func(t *testing.T) {
			h := fixture(t)
			name := filepath.Join(t.TempDir(), "input.jsonl")
			w, err := New(name, h)
			if err != nil {
				t.Fatal(err)
			}
			switch scenario {
			case "last-key":
				key := uint8(1)
				if err := w.Append(Event{"press_key", 1, 0, 0, "none", Params{Key: &key}}); err != nil {
					t.Fatal(err)
				}
			case "api-error":
				x, y, pub := 1, 1, uint64(1)
				if err := w.Append(Event{"click", 1, 0, 0, "no_response", Params{X: &x, Y: &y}}); err != nil {
					t.Fatal(err)
				}
				if err := w.Append(Event{"snapshot", 2, 0, 0, "none", Params{Publication: &pub, FrontendSHA: &h.ConfigSHA, Panel: &Panel{}}}); err != nil {
					t.Fatal(err)
				}
			case "short-end-write":
				pub := uint64(1)
				if err := w.Append(Event{"snapshot", 1, 0, 0, "none", Params{Publication: &pub, FrontendSHA: &h.ConfigSHA, Panel: &Panel{}}}); err != nil {
					t.Fatal(err)
				}
				w.f = &faultFile{File: w.f.(*os.File), shortWrite: true}
			}
			if err := w.FinishAt(0, time.Now().Add(time.Second)); err == nil {
				t.Fatal("invalid finish accepted")
			}
			assertIncomplete(t, name)
		})
	}
}
