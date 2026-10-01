package main

import (
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/wicanr2/dosgolem/opl2"
	"github.com/wicanr2/dosgolem/oracle"
)

type fakeAudioDevice struct{ plays, pauses, closes int }

func (d *fakeAudioDevice) Play()        { d.plays++ }
func (d *fakeAudioDevice) Pause()       { d.pauses++ }
func (d *fakeAudioDevice) Close() error { d.closes++; return nil }

func testPlayerAudio(t *testing.T, k *big.Rat, device audioDevice) *playerAudio {
	t.Helper()
	s, e := opl2.NewStream(opl2.StreamConfig{H: wizardryClock(), Tempo: k})
	if e != nil {
		t.Fatal(e)
	}
	stop := newStopRequest()
	return &playerAudio{stream: s, reader: newAudioReader(), device: device, stop: stop, clock: audioClock{epoch: time.Now().Add(-time.Hour), stop: stop}}
}

func TestAudioConfigurationAndFirstPlay(t *testing.T) {
	for _, bad := range []float64{0, -1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, e := audioTempo(bad, true); e == nil {
			t.Fatalf("接受無效ips%v", bad)
		}
	}
	k, e := audioTempo(0, false)
	if e != nil || k.Cmp(big.NewRat(1, 1)) != 0 {
		t.Fatal("預設不是精確1倍")
	}
	for _, speed := range []int64{1, 2, 4} {
		d := &fakeAudioDevice{}
		a := testPlayerAudio(t, big.NewRat(speed, 1), d)
		a.progress(oracle.CompletedOPL{Step: 200000, FirstSequence: 1})
		// 196875030000 / 65536 不是浮點預設值；獨立整數ceil核對。
		n := uint64(200000) * 44100 * 65536
		den := uint64(196875030000) * uint64(speed)
		want := (n + den - 1) / den
		if d.plays != 1 || a.reader.Stats().Queued != audioPrefillFrames+int(want) || a.stream.Stats().Frames != want {
			t.Fatalf("倍速%d device=%+v reader=%+v stream=%+v", speed, d, a.reader.Stats(), a.stream.Stats())
		}
		a.progress(oracle.CompletedOPL{Step: 400000, FirstSequence: 1, Final: true})
		if d.plays != 1 {
			t.Fatal("重複首次Play")
		}
	}
}

func TestAudioLowSpeedFirstChunkAndExplicitStop(t *testing.T) {
	d := &fakeAudioDevice{}
	a := testPlayerAudio(t, big.NewRat(1, 10), d)
	// 第一塊的絕對期限仍在等待，停止要求必須喚醒，不等整個低速批次。
	a.clock.epoch = time.Now()
	done := make(chan struct{})
	go func() { a.progress(oracle.CompletedOPL{Step: 200000, FirstSequence: 1}); close(done) }()
	readyDeadline := time.Now().Add(time.Second)
	for a.reader.Stats().Queued != 2*audioPrefillFrames {
		if time.Now().After(readyDeadline) {
			a.stop.Request()
			<-done
			t.Fatal("首塊未入隊")
		}
		time.Sleep(time.Millisecond)
	}
	a.stop.Request()
	select {
	case <-done:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("低速Pace沒有被停止喚醒")
	}
	if d.plays != 0 || !a.stream.Stats().Incomplete || a.reader.Stats().StopQueued != 2*audioPrefillFrames {
		t.Fatalf("device=%+v reader=%+v stream=%+v", d, a.reader.Stats(), a.stream.Stats())
	}
}

func TestAudioClockPauseAndFaultFallback(t *testing.T) {
	stop := newStopRequest()
	epoch := time.Now()
	c := audioClock{epoch: epoch, stop: stop}
	c.Pause(epoch.Add(time.Second))
	c.Resume(epoch.Add(4 * time.Second))
	if !c.epoch.Equal(epoch.Add(3 * time.Second)) {
		t.Fatal("暫停未平移epoch")
	}
	d := &fakeAudioDevice{}
	a := testPlayerAudio(t, big.NewRat(1, 1), d)
	a.fault = "device_init"
	a.clock.epoch = time.Now()
	before := time.Now()
	a.progress(oracle.CompletedOPL{Step: 200000, FirstSequence: 1})
	if time.Since(before) < 50*time.Millisecond || a.stream.Stats().Frames != 0 || d.plays != 0 {
		t.Fatal("故障沒有維持相同節流或仍播放")
	}
}

func TestAudioStopRequestSingleDeadline(t *testing.T) {
	s := newStopRequest()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); s.Request() }()
	}
	wg.Wait()
	deadline := s.Deadline()
	time.Sleep(time.Millisecond)
	s.Request()
	if !s.Requested() || s.Deadline() != deadline || deadline.Sub(s.at) != 5*time.Second {
		t.Fatal("退出期限被延長")
	}
}

func TestAudioHugePaceRemainsBoundedAndStopWakes(t *testing.T) {
	s := newStopRequest()
	c := audioClock{epoch: time.Now(), stop: s}
	phase := new(big.Rat).SetInt(new(big.Int).Lsh(big.NewInt(1), 200))
	done := make(chan error, 1)
	go func() { done <- c.Pace(phase) }()
	select {
	case <-done:
		t.Fatal("大目標沒有節流")
	case <-time.After(10 * time.Millisecond):
	}
	s.Request()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("大目標停止不能喚醒")
	}
}

func TestAudioVisualPublicationIndependent(t *testing.T) {
	original := &textEv{s: "A", x0: 3, g: [][]uint16{{0x1234, 0x5678}}}
	v := cloneVisualEvents([]*textEv{original})
	original.s = "B"
	original.g[0][0] = 0
	original.g = append(original.g, []uint16{1})
	if v[0].s != "A" || v[0].g[0][0] != 0x1234 || len(v[0].g) != 1 || v[0].page != nil || v[0].compound != nil {
		t.Fatal("顯示資料仍共享owner別名")
	}
}

func TestAudioStopPreservesSynthAndCountsLaterSource(t *testing.T) {
	d := &fakeAudioDevice{}
	a := testPlayerAudio(t, big.NewRat(1, 1), d)
	a.progress(oracle.CompletedOPL{Step: 200000, FirstSequence: 1,
		Writes: []oracle.OPLWrite{{Step: 100, Reg: 0x20, Val: 0}}})
	before := a.stream.Stats()
	a.stop.Request()
	// 來源仍須完成原操作；停止後不能把其事件送進已凍結的合成器。
	later := make([]oracle.OPLWrite, 16)
	for i := range later {
		later[i] = oracle.OPLWrite{Step: 200001 + uint64(i), Reg: 0x20}
	}
	a.progress(oracle.CompletedOPL{Step: 400000, FirstSequence: 2, Writes: later})
	a.progress(oracle.CompletedOPL{Step: 400000, FirstSequence: 18, Final: true})
	after := a.stream.Stats()
	if a.stoppedWrites != 16 || after.Received != before.Received || after.Applied != before.Applied ||
		after.Frames != before.Frames || after.EndStep != before.EndStep || d.plays != 1 || a.fault != "" {
		t.Fatalf("停止後合成或計數錯誤：before=%+v after=%+v skipped=%d device=%+v fault=%s", before, after, a.stoppedWrites, d, a.fault)
	}
}
