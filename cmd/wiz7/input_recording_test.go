package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
	rc "github.com/wicanr2/wizardry7-cht/recording"
)

// 只測自製錄放、Draw通知與生命週期，沒有載入／代製原版程式。
func fixtureInputRecording(t *testing.T) (*inputRecording, string) {
	t.Helper()
	hash := ic.Hash([]byte("self-authored configuration"))
	file := ic.File{Path: "FIXTURE", Size: 12, SHA256: hash, DOSDate: 23875, DOSTime: 25692, Attributes: 32}
	audio, err := ic.AudioConfig(3004074.5544433594, false, true)
	if err != nil {
		t.Fatal(err)
	}
	c := ic.Config{Schema: "wiz7-startup-config/1", Program: ic.Program{SHA256: hash, GoVersion: "go1.24.13", GOOS: "linux", GOARCH: "amd64", ModuleVersion: "(devel)", BuildSettingsSHA256: hash}, Sources: []ic.Source{{Path: "self.go", Size: 12, SHA256: hash}}, LoadWith: ic.LoadOptions{Adlib: true}, Initial: ic.Scalars{Schema: "wiz7-initial-scalars/1", Clock: ic.Clock{CPUHz: 33000000, IRQ0Base: 165000, IRQ0Every: 165000}, Platform: ic.Platform{DOSDate: ic.Date{Year: 1993, Month: 1, Day: 1, Weekday: 5}, DOSTimeModel: "now_plus_ticks_5493_div_1000_mod_8640000", BDATicksObservation: "not_observed", EffectiveMouseScaleObservation: "not_observed"}}, SelectedProgram: ic.SelectedProgram{SourceRole: "root", File: file}, Storage: ic.Storage{Resolution: "scratch_then_root", MetadataModel: "dos_archive_mtime_source_policy", Root: []ic.File{file}, Scratch: []ic.File{}}, Audio: audio, Frontend: ic.Frontend{Scale: 3, HDComposeDisabled: true, AutoplayDisabled: true, TestInputsDisabled: true, Resources: []ic.Resource{{Role: "xlate"}, {Role: "flow"}, {Role: "help"}, {Role: "roster"}, {Role: "profile-sex"}, {Role: "intro-picture"}}}}
	header, err := rc.NewHeader(c)
	if err != nil {
		t.Fatal(err)
	}
	name := filepath.Join(t.TempDir(), "input.jsonl")
	writer, err := rc.New(name, header)
	if err != nil {
		t.Fatal(err)
	}
	l := &inputRecording{writer: writer, configSHA: header.ConfigSHA, requested: make(chan struct{}), acks: make(chan recordingPublication, 1), stop: newStopRequest(), boundaryEligible: func() bool { return true }}
	return l, name
}
func assertRecordingFragment(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Fatal("unfinished recording published", err)
	}
	b, err := os.ReadFile(name + ".partial")
	if err != nil || bytes.Contains(b, []byte(`"kind":"end"`)) {
		t.Fatal("unfinished end", err)
	}
}
func TestFinalDrawEligibilityCheckedBeforeWriterHandoff(t *testing.T) {
	l, name := fixtureInputRecording(t)
	l.Request()
	l.boundaryEligible = func() bool { return false }
	l.finish(recordingPublication{ID: 1, Config: l.configSHA, Completed: time.Now()})
	if l.state.Load() != 3 || l.workerOwns || l.workerDone != nil {
		t.Fatal("ineligible final Draw entered writer worker")
	}
	assertRecordingFragment(t, name)
}
func TestRecordingFinalSnapshotAndCompletedCancellation(t *testing.T) {
	l, name := fixtureInputRecording(t)
	l.Request()
	l.finish(recordingPublication{ID: 1, Config: l.configSHA, Completed: time.Now()})
	if l.state.Load() != 2 {
		t.Fatal("successful finish rejected", l.failure.Load())
	}
	<-l.workerDone
	doc, err := rc.Load(name)
	if err != nil || len(doc.Events) != 1 || doc.Events[0].Kind != "snapshot" || !doc.Footer.Accepted {
		t.Fatal("completed file", err)
	}
	before, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	l.Cancel()
	l.Request()
	l.Poll()
	after, err := os.ReadFile(name)
	if err != nil || !bytes.Equal(before, after) || l.state.Load() != 2 {
		t.Fatal("completed recording cancelled", err)
	}
	l.stop.Request()
	if l.exitCode(l.stop.Deadline()) != 0 {
		t.Fatal("completed record exit failed")
	}
}
func TestCancelDuringDrawHasNoWriterRaceOrEnd(t *testing.T) {
	l, name := fixtureInputRecording(t)
	p := recordingPublication{ID: 1, Config: l.configSHA}
	done := make(chan struct{})
	go func() { l.afterPublish(p); close(done) }()
	l.Cancel()
	l.stop.Request()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Draw cancellation blocked owner")
	}
	if l.state.Load() != 3 || l.workerOwns {
		t.Fatal("cancelled Draw succeeded")
	}
	assertRecordingFragment(t, name)
}
func TestRecordingErrorClasses(t *testing.T) {
	if recordingError(nil) != "none" {
		t.Fatal("nil error class")
	}
	if recordingError(errors.New("boom")) != "execution" {
		t.Fatal("generic error class")
	}
	if recordingError(&oracle.NoResponseError{X: 1, Y: 2, Polls: 3}) != "no_response" {
		t.Fatal("no-response error class")
	}
}

func fixtureReplayRecording(t *testing.T, events []rc.Event) *inputRecording {
	t.Helper()
	return &inputRecording{replay: true, document: rc.Document{Events: events}, configSHA: "test-config",
		requested: make(chan struct{}), acks: make(chan recordingPublication, 1), stop: newStopRequest()}
}

func snapshotEvent(step, publication uint64, panel rc.Panel) rc.Event {
	return rc.Event{Kind: "snapshot", Seq: 1, Before: step, After: step, Error: "none",
		Params: rc.Params{Publication: &publication, FrontendSHA: new(string), Panel: &panel}}
}

func TestReplayPublishRejectsExhaustedCursor(t *testing.T) {
	l := fixtureReplayRecording(t, nil)
	sh := &shared{}
	if p, ok := l.publish(100, sh); ok || p.ID != 1 {
		t.Fatal("exhausted cursor published")
	}
	if l.state.Load() != 3 || l.failure.Load().code != "input_step_mismatch" || !l.stop.Requested() {
		t.Fatal("exhausted cursor not rejected")
	}
}

func TestReplayPublishRejectsBeforeMismatch(t *testing.T) {
	l := fixtureReplayRecording(t, []rc.Event{snapshotEvent(100, 1, rc.Panel{})})
	if _, ok := l.publish(200, &shared{}); ok {
		t.Fatal("before mismatch published")
	}
	if l.state.Load() != 3 || l.failure.Load().code != "input_step_mismatch" {
		t.Fatal("before mismatch not rejected")
	}
}

func TestReplayPublishRejectsPublicationMismatch(t *testing.T) {
	l := fixtureReplayRecording(t, []rc.Event{snapshotEvent(100, 7, rc.Panel{})})
	if _, ok := l.publish(100, &shared{}); ok {
		t.Fatal("publication mismatch published")
	}
	if l.state.Load() != 3 || l.failure.Load().code != "input_step_mismatch" {
		t.Fatal("publication mismatch not rejected")
	}
}

func TestReplayPublishAcceptsAndAdvancesCursor(t *testing.T) {
	panel := rc.Panel{HelpOpen: true, HelpScreen: "general", HelpPage: 3, MapOpen: true}
	l := fixtureReplayRecording(t, []rc.Event{snapshotEvent(100, 1, panel)})
	sh := &shared{}
	p, ok := l.publish(100, sh)
	if !ok || p.ID != 1 {
		t.Fatal("matching snapshot rejected")
	}
	if !sh.helpOpen || sh.helpScreen != "general" || sh.helpPage != 3 || !sh.mapOpen {
		t.Fatal("panel not applied to shared")
	}
	l.acknowledge(p, panel)
	l.afterPublish(p)
	if l.cursor != 1 || !l.replayDone.Load() || l.state.Load() != 2 {
		t.Fatal("matching panel did not complete replay")
	}
}

func TestReplayPublishRejectsPanelMismatch(t *testing.T) {
	panel := rc.Panel{HelpOpen: true}
	l := fixtureReplayRecording(t, []rc.Event{snapshotEvent(100, 1, panel)})
	sh := &shared{}
	p, ok := l.publish(100, sh)
	if !ok {
		t.Fatal("matching snapshot rejected")
	}
	l.acknowledge(p, rc.Panel{HelpOpen: false})
	l.afterPublish(p)
	if l.state.Load() != 3 || l.failure.Load().code != "draw_mismatch" || !l.stop.Requested() {
		t.Fatal("panel mismatch not rejected")
	}
}

func TestDrawAckMismatchAndTimeout(t *testing.T) {
	l := fixtureReplayRecording(t, nil)
	p := recordingPublication{ID: 1, Step: 5, Config: "test-config"}
	l.acknowledge(p, rc.Panel{})
	if _, code := l.drawAck(recordingPublication{ID: 2, Step: 5, Config: "test-config"}, time.Now().Add(time.Second)); code != "draw_mismatch" {
		t.Fatal("wrong publication acknowledged", code)
	}
	l2 := fixtureReplayRecording(t, nil)
	if _, code := l2.drawAck(p, time.Now().Add(-time.Millisecond)); code != "draw_timeout" {
		t.Fatal("expired draw not timed out", code)
	}
	l3 := fixtureReplayRecording(t, nil)
	l3.stop.Request()
	if _, code := l3.drawAck(p, time.Now().Add(time.Second)); code != "recording_ineligible" {
		t.Fatal("stopped draw not ineligible", code)
	}
}

func TestFinishExpiredDeadlineKeepsFragment(t *testing.T) {
	l, name := fixtureInputRecording(t)
	l.Request()
	l.deadline = time.Now().Add(-time.Millisecond)
	l.finish(recordingPublication{ID: 1, Config: l.configSHA, Completed: time.Now()})
	if l.state.Load() != 3 || l.workerOwns {
		t.Fatal("expired finish entered writer worker")
	}
	if l.failure.Load().code != "draw_timeout" {
		t.Fatal("expired finish code", l.failure.Load().code)
	}
	assertRecordingFragment(t, name)
}

func TestRecordingRequestUsesOneAbsoluteDeadline(t *testing.T) {
	l, name := fixtureInputRecording(t)
	l.Request()
	deadline := l.deadline
	l.Request()
	if l.deadline != deadline {
		t.Fatal("repeated stop reset deadline")
	}
	l.deadline = time.Now().Add(-time.Millisecond)
	l.Poll()
	l.abortOwned()
	if l.state.Load() != 3 {
		t.Fatal("expired request remains active")
	}
	assertRecordingFragment(t, name)
}
