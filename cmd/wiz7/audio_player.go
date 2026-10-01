package main

import (
	"errors"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/wicanr2/dosgolem/opl2"
	"github.com/wicanr2/dosgolem/oracle"
)

// stopRequest 可由GUI提出；Oracle、Stream與播放器控制仍由owner操作。
type stopRequest struct {
	once sync.Once
	done chan struct{}
	at   time.Time
}

func newStopRequest() *stopRequest { return &stopRequest{done: make(chan struct{})} }
func (s *stopRequest) Request() {
	s.once.Do(func() { s.at = time.Now(); close(s.done) })
}
func (s *stopRequest) Requested() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}
func (s *stopRequest) Deadline() time.Time { return s.at.Add(5 * time.Second) }

func wizardryClock() *big.Rat { return big.NewRat(196875030000, 65536) }
func audioTempo(ips float64, explicit bool) (*big.Rat, error) {
	if !explicit {
		return big.NewRat(1, 1), nil
	}
	if math.IsNaN(ips) || math.IsInf(ips, 0) || ips <= 0 {
		return nil, errors.New("invalid_ips")
	}
	return new(big.Rat).Quo(new(big.Rat).SetFloat64(ips), wizardryClock()), nil
}

// 不格式化原版錯誤的Console；只保留有限分類供診斷與資格判定。
func frontendOperationCode(err error) string {
	if err == nil {
		return ""
	}
	var budget *oracle.BudgetError
	var exit *oracle.ExitError
	var response *oracle.NoResponseError
	switch {
	case errors.Is(err, oracle.ErrOPLProgressReentrant):
		return "reentrant"
	case errors.As(err, &budget):
		return "budget"
	case errors.As(err, &exit):
		return "exit"
	case errors.As(err, &response):
		return "no_response"
	default:
		return "execution"
	}
}

type audioDevice interface {
	Play()
	Pause()
	Close() error
}

// audioClock 不逐塊累積duration；每次以相同epoch計算絕對有理目標。
type audioClock struct {
	epoch   time.Time
	pauseAt time.Time
	stop    *stopRequest
}

func (c *audioClock) Pace(phase *big.Rat) error {
	if c.stop.Requested() {
		return nil
	}
	if phase == nil || phase.Sign() < 0 {
		return errors.New("pace_phase")
	}
	ns := new(big.Rat).Mul(phase, big.NewRat(int64(time.Second), audioRate))
	target := new(big.Int).Quo(ns.Num(), ns.Denom())
	// 極慢但有限正值速度的目標可能超過time.Duration；逐段等待同一
	// 大整數絕對目標，不限縮速度、不因溢位失去節流。
	for {
		left := new(big.Int).Sub(target, big.NewInt(int64(time.Since(c.epoch))))
		if left.Sign() <= 0 {
			return nil
		}
		delay := time.Second
		if left.IsInt64() {
			delay = min(delay, time.Duration(left.Int64()))
		}
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-c.stop.done:
			timer.Stop()
			return nil
		}
	}
}
func (c *audioClock) Pause(now time.Time) {
	if c.pauseAt.IsZero() {
		c.pauseAt = now
	}
}
func (c *audioClock) Resume(now time.Time) {
	if !c.pauseAt.IsZero() {
		c.epoch = c.epoch.Add(now.Sub(c.pauseAt))
		c.pauseAt = time.Time{}
	}
}

type playerAudio struct {
	stream                   *opl2.Stream
	reader                   *audioReader
	device                   audioDevice
	clock                    audioClock
	stop                     *stopRequest
	started, paused, stopped bool
	fault                    string
	stoppedWrites            uint64 // 停止要求後收到、未送入合成器的來源事件。
}

func newPlayerAudio(o *oracle.Oracle, k *big.Rat, mute bool, stop *stopRequest) (*playerAudio, error) {
	clock := o.ClockSettings()
	if clock.CycleClock || clock.IRQ0Base != 165000 {
		return nil, errors.New("clock_mismatch")
	}
	s, err := opl2.NewStream(opl2.StreamConfig{H: wizardryClock(), Tempo: k, InitialStep: o.Steps(), ExistingWrites: uint64(len(o.OPLWrites()))})
	if err != nil {
		return nil, err
	}
	a := &playerAudio{stream: s, stop: stop, clock: audioClock{stop: stop}, reader: newAudioReader()}
	if !mute {
		p, e := audio.NewContext(audioRate).NewPlayer(a.reader)
		if e != nil {
			a.fault = "device_init"
		} else {
			p.SetBufferSize(50 * time.Millisecond)
			a.device = p
		}
	}
	_, err = o.InstallOPLProgress(oracle.OPLProgressQuantum, a.progress, oracle.WithOPLStream())
	if err != nil {
		if a.device != nil {
			_ = a.device.Close()
		}
		return nil, err
	}
	// 玩家時間從資源／掛鉤準備完成後首次推進開始，裝置初始化不算遊戲時間。
	a.clock.epoch = time.Now()
	return a, nil
}

func (a *playerAudio) fail(code string) {
	if a.fault == "" {
		a.fault = code
	}
}
func (a *playerAudio) stopAtBoundary() {
	if a.stopped {
		return
	}
	a.stopped = true
	a.stream.Stop()
	a.reader.Stop()
	// F1期間關窗也須消費已開始的frame與短淡出；未首次Play的裝置不新啟動。
	if a.started && a.paused && a.device != nil {
		a.device.Play()
		a.paused = false
	}
}
func (a *playerAudio) progress(batch oracle.CompletedOPL) {
	if a.stop.Requested() {
		a.stoppedWrites += uint64(len(batch.Writes))
		a.stopAtBoundary()
		return
	}
	if batch.ObserverFaultCode != "" {
		a.fail("observer_fault")
	}
	if a.fault == "" {
		err := a.stream.Accept(batch, func(chunk opl2.PCMChunk) error {
			if a.stop.Requested() {
				a.stopAtBoundary()
				return errors.New("stopped")
			}
			if a.device != nil {
				if err := a.reader.Push(chunk.Start, chunk.Samples); err != nil {
					return errors.New("reader_rejected")
				}
			}
			if err := a.clock.Pace(chunk.Phase); err != nil {
				return err
			}
			if a.stop.Requested() {
				a.stopAtBoundary()
				return errors.New("stopped")
			}
			if !a.started && batch.Step >= oracle.OPLProgressQuantum && a.device != nil {
				a.device.Play()
				a.started = true
			}
			return nil
		}, a.stop.Requested)
		if err != nil && !a.stop.Requested() {
			if e, ok := err.(*opl2.StreamError); ok {
				a.fail(e.Code)
			} else {
				a.fail("consumer_fault")
			}
		}
	}
	// 音訊故障不使原版變快。即使沒有接受PCM也使用同一精確相位／epoch。
	if a.fault != "" {
		a.paceAt(batch.Step)
	}
	if a.stop.Requested() {
		a.stopAtBoundary()
	}
}
func (a *playerAudio) paceAt(step uint64) {
	phase, err := a.stream.PhaseAt(step)
	if err != nil {
		a.fail("pace_phase")
		return
	}
	if err := a.clock.Pace(phase); err != nil {
		a.fail("pace_overflow")
	}
}

// AfterOperation 必須於每個原操作返回執行，觀察回呼被卸下也維持節流。
func (a *playerAudio) AfterOperation(o *oracle.Oracle) {
	if o.OPLProgressFault() != "" {
		a.fail("observer_fault")
	}
	if a.stop.Requested() {
		a.stopAtBoundary()
		return
	}
	a.paceAt(o.Steps())
}
func (a *playerAudio) Pause(paused bool) {
	if paused == a.paused {
		return
	}
	a.paused = paused
	if paused {
		a.clock.Pause(time.Now())
		if a.started && a.device != nil {
			a.device.Pause()
		}
	} else {
		a.clock.Resume(time.Now())
		if a.started && a.device != nil {
			a.device.Play()
		}
	}
}

// shutdownPlayer 在owner返回後執行；主程式仍以同一五秒期限監督此單次Close。
// Reader ACK只是byte交付，不能當作實體喇叭／ALSA關閉的ACK。
func (a *playerAudio) shutdownPlayer() string {
	a.stopAtBoundary()
	result := ""
	if a.started && a.device != nil {
		left := time.Until(a.stop.Deadline())
		if left <= 0 {
			return "timeout"
		}
		ackDeadline := time.Now().Add(min(left, 500*time.Millisecond))
		timer := time.NewTimer(max(time.Until(ackDeadline), 0))
		select {
		case <-a.reader.ack:
			if !time.Now().Before(ackDeadline) {
				result = "fade_timeout"
			}
		case <-timer.C:
			result = "fade_timeout"
		}
		timer.Stop()
	}
	if a.device != nil {
		if err := a.device.Close(); err != nil {
			result = "close_failure"
		}
	}
	return result
}
