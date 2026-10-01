package main

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

const audioRate = 44100
const audioRingFrames = 2 * audioRate
const audioPrefillFrames = audioRate / 10
const audioFadeFrames = 221

// scaleSample 是所有裝置淡化共用的整數 half-away 捨入；核心 PCM 不改動。
func scaleSample(v int16, numerator, denominator int) int16 {
	if denominator <= 0 {
		return 0
	}
	n := int64(v) * int64(numerator)
	negative := n < 0
	if negative {
		n = -n
	}
	q := n / int64(denominator)
	if (n%int64(denominator))*2 >= int64(denominator) {
		q++
	}
	if negative {
		q = -q
	}
	return int16(q)
}

type audioFrame struct {
	value  int16
	source uint64
	real   bool
}
type audioReaderStats struct {
	Queued, MaximumQueued                                       int
	Underflow, Overflow, Recovery, DeviceFrames, DeliveredBytes uint64
	SourceFrames, LastSource                                    uint64
	StopQueued, StopFade, StopDiscard                           uint64
	FadeDeliveredBytes                                          uint64
	Stopped, ACK, EOF                                           bool
}

// audioReader 只消費有界 PCM，不操作 Oracle，不寫檔、不保存逐幀審計。
// owner 推送／停止；裝置可依任意 byte 長度讀取。
type audioReader struct {
	mu                  sync.Mutex
	ring                [audioRingFrames]audioFrame
	head, count         int
	nextSource          uint64
	partial             [4]byte
	partialAt, partialN int
	partialFade         bool
	gap                 bool
	recovery            int
	stop                [audioFadeFrames]int16
	stopN, stopAt       int
	ack                 chan struct{}
	stats               audioReaderStats
}

func newAudioReader() *audioReader {
	r := &audioReader{count: audioPrefillFrames, gap: true, recovery: audioFadeFrames, ack: make(chan struct{})}
	r.stats.MaximumQueued = r.count
	return r
}

func (r *audioReader) Push(start uint64, samples []int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stats.Stopped {
		return errors.New("stopped")
	}
	if start != r.nextSource || uint64(len(samples)) > ^uint64(0)-start {
		return errors.New("pcm_sequence")
	}
	// 每次最多4410幀，由消費端逐塊節流，不可把整段低速PCM塞進ring。
	if len(samples) > audioPrefillFrames {
		return errors.New("pcm_capacity")
	}
	for i, v := range samples {
		if r.count == len(r.ring) {
			r.head = (r.head + 1) % len(r.ring)
			r.count--
			r.stats.Overflow++
		}
		at := (r.head + r.count) % len(r.ring)
		r.ring[at] = audioFrame{v, start + uint64(i), true}
		r.count++
	}
	r.nextSource += uint64(len(samples))
	r.stats.MaximumQueued = max(r.stats.MaximumQueued, r.count)
	return nil
}

func (r *audioReader) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stats.Stopped {
		return
	}
	r.stats.Stopped = true
	r.stopN = min(r.count, audioFadeFrames)
	r.stats.StopQueued = uint64(r.count)
	r.stats.StopFade = uint64(r.stopN)
	r.stats.StopDiscard = uint64(r.count - r.stopN)
	for i := 0; i < r.stopN; i++ {
		r.stop[i] = scaleSample(r.ring[(r.head+i)%len(r.ring)].value, r.stopN-1-i, r.stopN-1)
	}
	r.head, r.count = 0, 0
	r.checkACK()
}

func (r *audioReader) checkACK() {
	if r.stats.Stopped && r.stopAt == r.stopN && r.partialN == 0 && !r.stats.ACK {
		r.stats.ACK = true
		close(r.ack)
	}
}

func (r *audioReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for n < len(p) {
		if r.partialN == 0 {
			var v int16
			r.partialFade = false
			if r.stats.Stopped {
				if r.stopAt == r.stopN {
					break
				}
				v = r.stop[r.stopAt]
				r.stopAt++
				r.partialFade = true
			} else if r.count > 0 {
				f := r.ring[r.head]
				r.head = (r.head + 1) % len(r.ring)
				r.count--
				v = f.value
				if f.real {
					r.stats.SourceFrames++
					r.stats.LastSource = f.source
					if r.gap {
						r.gap = false
						r.recovery = 0
						r.stats.Recovery++
					}
					if r.recovery < audioFadeFrames {
						v = scaleSample(v, r.recovery, audioFadeFrames-1)
						r.recovery++
					}
				}
			} else {
				r.stats.Underflow++
				r.gap = true
			}
			binary.LittleEndian.PutUint16(r.partial[:2], uint16(v))
			copy(r.partial[2:], r.partial[:2])
			r.partialAt, r.partialN = 0, 4
			r.stats.DeviceFrames++
		}
		k := min(len(p)-n, r.partialN)
		copy(p[n:n+k], r.partial[r.partialAt:r.partialAt+k])
		r.partialAt += k
		r.partialN -= k
		n += k
		r.stats.DeliveredBytes += uint64(k)
		if r.partialFade {
			r.stats.FadeDeliveredBytes += uint64(k)
		}
	}
	r.checkACK()
	if n == 0 && r.stats.Stopped {
		r.stats.EOF = true
		return 0, io.EOF
	}
	return n, nil
}

func (r *audioReader) Stats() audioReaderStats {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.stats
	s.Queued = r.count
	return s
}
