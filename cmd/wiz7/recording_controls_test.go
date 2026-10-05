package main

import "testing"

func TestRecordingEffectiveScaleClosesUncapturedSetting(t *testing.T) {
	for _, scale := range []int{2, 3} {
		for _, explicit := range []bool{false, true} {
			for _, hd := range []int{2, 3} {
				got, err := recordingEffectiveHDScale(scale, hd, explicit, true)
				if explicit && hd != scale {
					if err == nil {
						t.Fatal("uncaptured explicit scale accepted")
					}
				} else if err != nil || got != scale {
					t.Fatal("effective scale", got, err)
				}
				ordinary, err := recordingEffectiveHDScale(scale, hd, explicit, false)
				if err != nil || ordinary != hd {
					t.Fatal("ordinary/offline scale changed")
				}
			}
		}
	}
}
