package main

import (
	"encoding/binary"
	"io"
	"math/big"
	"testing"
)

func TestAudioIntegerFadeAllSamples(t *testing.T) {
	// 獨立有理oracle；窮舉int16與全部221個gain，包含舊浮點反例。
	for v := -32768; v <= 32767; v++ {
		for k := 0; k <= 220; k++ {
			n := int64(v) * int64(k)
			negative := n < 0
			if negative {
				n = -n
			}
			q := new(big.Int).Quo(big.NewInt(n*2+220), big.NewInt(440)).Int64()
			if negative {
				q = -q
			}
			if got := scaleSample(int16(v), k, 220); got != int16(q) {
				t.Fatalf("v=%d k=%d got=%d want=%d", v, k, got, q)
			}
		}
	}
}

func TestAudioStopPartialAndShort(t *testing.T) {
	for _, count := range []int{0, 1, 2, 17, 220, 221, 222} {
		for partial := 0; partial < 4; partial++ {
			r := newAudioReader()
			// 先消費人工預填，不把預填算進短尾段。
			_, _ = r.Read(make([]byte, audioPrefillFrames*4))
			if partial > 0 {
				_, _ = r.Read(make([]byte, partial))
			}
			values := make([]int16, count)
			for i := range values {
				values[i] = int16(23837 - i*71)
			}
			if err := r.Push(0, values); err != nil {
				t.Fatal(err)
			}
			r.Stop()
			var out []byte
			for i := 0; i < 2000; i++ {
				b := make([]byte, 1+i%7)
				n, e := r.Read(b)
				out = append(out, b[:n]...)
				if e == io.EOF {
					break
				}
				if e != nil {
					t.Fatal(e)
				}
			}
			prefix := 0
			if partial > 0 {
				prefix = 4 - partial
			}
			fade := min(count, 221)
			if len(out) != prefix+fade*4 {
				t.Fatalf("count=%d partial=%d bytes=%d", count, partial, len(out))
			}
			for i := 0; i < fade; i++ {
				v := values[i]
				want := 0
				if fade > 1 {
					n := int(v) * (fade - 1 - i)
					want = (n + (fade-1)/2) / (fade - 1)
				}
				left := int16(binary.LittleEndian.Uint16(out[prefix+i*4:]))
				right := int16(binary.LittleEndian.Uint16(out[prefix+i*4+2:]))
				if left != int16(want) || right != left {
					t.Fatalf("短尾段%d/%d=%d,%d want=%d", i, fade, left, right, want)
				}
			}
			s := r.Stats()
			if !s.ACK || !s.EOF || s.StopDiscard != uint64(count-fade) || s.FadeDeliveredBytes != uint64(4*fade) {
				t.Fatalf("%+v", s)
			}
			if err := r.Push(uint64(count), []int16{1}); err == nil {
				t.Fatal("停止後仍供應")
			}
		}
	}
}

func TestAudioBoundedOverflowAndRecovery(t *testing.T) {
	r := newAudioReader()
	values := make([]int16, 4410)
	for i := range values {
		values[i] = 23837
	}
	for i := 0; i < 23; i++ {
		if err := r.Push(uint64(i*len(values)), values); err != nil {
			t.Fatal(err)
		}
	}
	s := r.Stats()
	if s.Queued != audioRingFrames || s.MaximumQueued != audioRingFrames || s.Overflow != uint64(24*4410-audioRingFrames) {
		t.Fatalf("%+v", s)
	}
	_, _ = r.Read(make([]byte, audioRingFrames*4+4))
	if err := r.Push(23*4410, values[:221]); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 221*4)
	_, _ = r.Read(b)
	if got := int16(binary.LittleEndian.Uint16(b[30*4:])); got != 3251 {
		t.Fatalf("既有浮點反例=%d", got)
	}
	if r.Stats().Underflow != 1 || r.Stats().Recovery != 2 {
		t.Fatalf("%+v", r.Stats())
	}
}
