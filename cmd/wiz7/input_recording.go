package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
	rc "github.com/wicanr2/wizardry7-cht/recording"
)

type recordingPublication struct {
	ID, Step  uint64
	Config    string
	Panel     rc.Panel
	Completed time.Time
}
type recordingFailure struct{ code string }
type inputRecording struct {
	replay           bool
	document         rc.Document
	writer           *rc.Writer
	configSHA        string
	font             *helpData
	state            atomic.Uint32 // 錄製中、收尾中、完成、失敗
	failure          atomic.Pointer[recordingFailure]
	replayDone       atomic.Bool
	requestOnce      sync.Once
	requested        chan struct{}
	deadline         time.Time // requested 關閉後不可變
	acks             chan recordingPublication
	workerDone       chan struct{}
	seq, publication uint64 // 機器 owner 專用
	cursor           int
	workerOwns       bool
	snapshot         func()
	boundaryEligible func() bool // 機器 owner 專用；worker 不碰原版或音訊。
	stop             *stopRequest
}

func newInputRecording(owner *recordingOwner, name string, document *rc.Document, font *helpData) (*inputRecording, error) {
	config, sha, err := owner.Config()
	if err != nil {
		return nil, err
	}
	l := &inputRecording{configSHA: sha, font: font, requested: make(chan struct{}), acks: make(chan recordingPublication, 1)}
	if document != nil {
		if document.Header.ConfigSHA != sha {
			return nil, errors.New("configuration_mismatch")
		}
		l.replay = true
		l.document = *document
		return l, nil
	}
	var c ic.Config
	if json.Unmarshal(config, &c) != nil {
		return nil, errors.New("configuration_mismatch")
	}
	header, err := rc.NewHeader(c)
	if err != nil {
		return nil, err
	}
	l.writer, err = rc.New(name, header)
	if err != nil {
		return nil, errors.New("recording_output_failed")
	}
	return l, nil
}
func recordingKey(k uint8) rc.Params    { return rc.Params{Key: &k} }
func recordingBytes(b []byte) rc.Params { v := append([]byte{}, b...); return rc.Params{Bytes: &v} }
func recordingXY(x, y int) rc.Params    { return rc.Params{X: &x, Y: &y} }
func recordingButton() rc.Params        { v := uint8(1); return rc.Params{Button: &v} }
func recordingSteps(n uint64) rc.Params { return rc.Params{Steps: &n} }
func recordingError(err error) string {
	if err == nil {
		return "none"
	}
	var response *oracle.NoResponseError
	if errors.As(err, &response) {
		return "no_response"
	}
	return "execution"
}
func (l *inputRecording) Request() {
	if l == nil || l.replay || l.state.Load() != 0 {
		return
	}
	l.requestOnce.Do(func() { l.deadline = time.Now().Add(5 * time.Second); close(l.requested) })
}
func (l *inputRecording) Requested() bool {
	select {
	case <-l.requested:
		return true
	default:
		return false
	}
}
func (l *inputRecording) fail(code string) {
	if code == "cleanup_incomplete" {
		l.failure.Store(&recordingFailure{code: code})
	} else {
		l.failure.CompareAndSwap(nil, &recordingFailure{code: code})
	}
	for {
		s := l.state.Load()
		if s == 2 || s == 3 {
			return
		}
		if l.state.CompareAndSwap(s, 3) {
			return
		}
	}
}

// Cancel 只改共同狀態，不碰阻塞中的 Writer。已完成錄製保持。
func (l *inputRecording) Cancel() {
	if l != nil && !l.replayDone.Load() {
		l.fail("recording_ineligible")
	}
}
func (l *inputRecording) Poll() {
	if l != nil && !l.replay && l.Requested() && !time.Now().Before(l.deadline) {
		l.fail("stop_timeout")
	}
}
func (l *inputRecording) abortOwned() {
	if !l.replay && !l.workerOwns && l.writer != nil {
		_ = l.writer.Abort()
	}
}
func (l *inputRecording) qualified(a *playerAudio, o *oracle.Oracle) bool {
	return a.fault == "" && o.OPLProgressFault() == "" && o.OPLSourceStats().FaultCode == "" && !a.stream.Stats().Incomplete
}

// 原操作返回值保持，音訊資格在 AfterOperation 返回後檢查。
func (l *inputRecording) Operation(o *oracle.Oracle, a *playerAudio, kind string, p rc.Params, call func() error) error {
	before := o.Steps()
	err := call()
	after := o.Steps()
	if kind == "run" || kind == "click" {
		a.AfterOperation(o)
	}
	if l.replay {
		return err
	}
	if l.state.Load() == 0 {
		l.seq++
		if l.writer.Append(rc.Event{Kind: kind, Seq: l.seq, Before: before, After: after, Error: recordingError(err), Params: p}) != nil {
			l.fail("writer_append")
		}
		if err != nil || !l.qualified(a, o) {
			l.fail("recording_ineligible")
		}
	}
	if l.state.Load() == 3 {
		l.abortOwned()
	}
	if l.Requested() && l.state.Load() == 0 && l.snapshot != nil {
		l.snapshot()
	}
	return err
}

func (l *inputRecording) publish(step uint64, sh *shared) (recordingPublication, bool) {
	if l == nil || !l.replay && l.state.Load() != 0 {
		return recordingPublication{}, false
	}
	l.publication++
	p := recordingPublication{ID: l.publication, Step: step, Config: l.configSHA}
	if l.replay {
		if l.cursor >= len(l.document.Events) {
			l.fail("input_step_mismatch")
			l.stop.Request()
			return p, false
		}
		e := l.document.Events[l.cursor]
		if e.Kind != "snapshot" || e.Before != step || *e.Params.Publication != p.ID {
			l.fail("input_step_mismatch")
			l.stop.Request()
			return p, false
		}
		p.Panel = *e.Params.Panel
		sh.helpOpen, sh.helpScreen, sh.helpPage, sh.mapOpen = p.Panel.HelpOpen, p.Panel.HelpScreen, int(p.Panel.HelpPage), p.Panel.MapOpen
	}
	sh.recordingPublication = p
	return p, true
}
func (l *inputRecording) acknowledge(p recordingPublication, panel rc.Panel) {
	if p.ID == 0 {
		return
	}
	p.Panel = panel
	p.Completed = time.Now()
	select {
	case l.acks <- p:
	default:
	}
}
func (l *inputRecording) drawAck(p recordingPublication, deadline time.Time) (recordingPublication, string) {
	timer := time.NewTimer(max(time.Until(deadline), 0))
	defer timer.Stop()
	select {
	case got := <-l.acks:
		if got.ID != p.ID || got.Step != p.Step || got.Config != p.Config || got.Completed.IsZero() {
			return got, "draw_mismatch"
		}
		if !got.Completed.Before(deadline) || !time.Now().Before(deadline) {
			return got, "draw_timeout"
		}
		return got, ""
	case <-timer.C:
		return recordingPublication{}, "draw_timeout"
	case <-l.stop.done:
		return recordingPublication{}, "recording_ineligible"
	}
}
func (l *inputRecording) afterPublish(p recordingPublication) {
	deadline := time.Now().Add(5 * time.Second)
	if !l.replay && l.Requested() {
		deadline = l.deadline
	}
	ack, code := l.drawAck(p, deadline)
	if code != "" {
		l.fail(code)
		l.abortOwned()
		if l.replay {
			l.stop.Request()
		}
		return
	}
	if l.replay {
		if ack.Panel != *l.document.Events[l.cursor].Params.Panel {
			l.fail("draw_mismatch")
			l.stop.Request()
			return
		}
		l.cursor++
		if l.cursor == len(l.document.Events) {
			l.replayDone.Store(true)
			l.state.Store(2)
			log.Print("重播完成")
			l.stop.Request()
		}
		return
	}
	if l.state.Load() != 0 {
		l.abortOwned()
		return
	}
	if l.Requested() {
		l.finish(ack)
		return
	}
	l.seq++
	id, sha, panel := p.ID, p.Config, ack.Panel
	if l.writer.Append(rc.Event{Kind: "snapshot", Seq: l.seq, Before: p.Step, After: p.Step, Error: "none", Params: rc.Params{Publication: &id, FrontendSHA: &sha, Panel: &panel}}) != nil {
		l.fail("writer_append")
		l.abortOwned()
	}
}
func (l *inputRecording) finish(ack recordingPublication) {
	if l.boundaryEligible == nil || !l.boundaryEligible() {
		l.fail("recording_ineligible")
		l.abortOwned()
		return
	}
	if !l.state.CompareAndSwap(0, 1) {
		l.abortOwned()
		return
	}
	if !time.Now().Before(l.deadline) || !ack.Completed.Before(l.deadline) {
		l.fail("draw_timeout")
		l.abortOwned()
		return
	}
	count := l.seq
	l.workerOwns = true
	l.workerDone = make(chan struct{})
	done := make(chan string, 1)
	go func() {
		defer close(l.workerDone)
		allowed := func() bool { return l.state.Load() == 1 && time.Now().Before(l.deadline) }
		id, sha, panel := ack.ID, ack.Config, ack.Panel
		code := ""
		if !allowed() {
			code = "stop_timeout"
		} else if l.writer.Append(rc.Event{Kind: "snapshot", Seq: count + 1, Before: ack.Step, After: ack.Step, Error: "none", Params: rc.Params{Publication: &id, FrontendSHA: &sha, Panel: &panel}}) != nil {
			code = "writer_append"
		}
		if code == "" {
			if err := l.writer.FinishWhen(ack.Step, l.deadline, allowed); err != nil {
				code = "writer_finish"
				if err.Error() == "cleanup_incomplete" {
					code = "cleanup_incomplete"
				}
			} else if !allowed() || !l.state.CompareAndSwap(1, 2) {
				code = "stop_timeout"
				if l.writer.Revoke() != nil {
					code = "cleanup_incomplete"
				}
			}
		}
		if code != "" {
			l.fail(code)
			_ = l.writer.Abort()
		}
		done <- code
	}()
	timer := time.NewTimer(max(time.Until(l.deadline), 0))
	defer timer.Stop()
	select {
	case code := <-done:
		if code == "" {
			log.Print("錄製完成，繼續遊戲")
		} else {
			l.reportFailure()
		}
	case <-timer.C:
		l.fail("stop_timeout")
		l.reportFailure()
	case <-l.stop.done:
		l.fail("recording_ineligible")
	}
}
func (l *inputRecording) reportFailure() {
	code := "recording_ineligible"
	if f := l.failure.Load(); f != nil {
		code = f.code
	}
	if code == "cleanup_incomplete" {
		log.Print("錄製失敗，檔案收尾未完成：cleanup_incomplete")
	} else {
		log.Printf("錄製失敗，已保留未完成片段：%s", code)
	}
}
func (l *inputRecording) ReplayNext(o *oracle.Oracle, a *playerAudio, snapshot func()) {
	if l.cursor >= len(l.document.Events) {
		l.fail("input_step_mismatch")
		l.stop.Request()
		return
	}
	e := l.document.Events[l.cursor]
	if o.Steps() != e.Before {
		l.fail("input_step_mismatch")
		l.stop.Request()
		return
	}
	p := e.Params
	if e.Kind == "snapshot" {
		snapshot()
		return
	}
	var err error
	switch e.Kind {
	case "press_key":
		o.PressKey(oracle.Key(*p.Key))
	case "type":
		o.Type(string(*p.Bytes))
	case "move_mouse":
		o.MoveMouse(*p.X, *p.Y)
	case "click":
		err = o.Click(*p.X, *p.Y)
	case "press_button":
		o.PressButton(int(*p.Button))
	case "release_button":
		o.ReleaseButton(int(*p.Button))
	case "run":
		err = o.Run(*p.Steps)
	}
	if e.Kind == "run" || e.Kind == "click" {
		a.AfterOperation(o)
	}
	if o.Steps() != e.After || recordingError(err) != e.Error {
		l.fail("api_result_mismatch")
		l.stop.Request()
		return
	}
	if !l.qualified(a, o) {
		l.fail("recording_ineligible")
		l.stop.Request()
		return
	}
	l.cursor++
}
func (l *inputRecording) exitCode(deadline time.Time) int {
	if l == nil {
		return 0
	}
	l.Cancel()
	if l.workerOwns {
		timer := time.NewTimer(max(time.Until(deadline), 0))
		defer timer.Stop()
		select {
		case <-l.workerDone:
		case <-timer.C:
			return 2
		}
	} else {
		l.abortOwned()
	}
	if l.replay {
		if l.replayDone.Load() {
			return 0
		}
		log.Print("重播取消，尚未完成")
		return 2
	}
	if l.state.Load() == 2 {
		return 0
	}
	l.reportFailure()
	return 2
}
