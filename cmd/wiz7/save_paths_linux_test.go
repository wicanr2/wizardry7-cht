package main

import (
	"errors"
	"testing"
)

func TestDefaultSaveLocation(t *testing.T) {
	cases := []struct {
		name, xdg, home, want string
		homeError, wantError  bool
		homeCalls             int
	}{
		{"absolute_xdg", "/data", "/unused", "/data/wizardry7-cht/save", false, false, 0},
		{"xdg_trailing_separator", "/data/", "/unused", "/data/wizardry7-cht/save", false, false, 0},
		{"xdg_root", "/", "/unused", "/wizardry7-cht/save", false, false, 0},
		{"empty_xdg", "", "/users/player", "/users/player/.local/share/wizardry7-cht/save", false, false, 1},
		{"relative_xdg", "relative", "/users/player", "/users/player/.local/share/wizardry7-cht/save", false, false, 1},
		{"relative_parent_xdg", "../data", "/users/player", "/users/player/.local/share/wizardry7-cht/save", false, false, 1},
		{"empty_home", "", "", "", false, true, 1},
		{"relative_home", "", "relative", "", false, true, 1},
		{"home_error", "", "", "", true, true, 1},
		{"unicode_home", "", "/玩家/資料", "/玩家/資料/.local/share/wizardry7-cht/save", false, false, 1},
		{"space_home", "", "/users/a b", "/users/a b/.local/share/wizardry7-cht/save", false, false, 1},
		{"xdg_bypasses_missing_home", "/data", "", "/data/wizardry7-cht/save", true, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			calls := 0
			p, e := linuxSaveDir(c.xdg, func() (string, error) {
				calls++
				if c.homeError {
					return "", errors.New("fixture")
				}
				return c.home, nil
			})
			if (e != nil) != c.wantError || p != c.want || calls != c.homeCalls {
				t.Fatalf("path=%q error=%v calls=%d", p, e, calls)
			}
		})
	}
}
func TestExplicitSaveOverridesDefault(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "relative")
	p, e := resolveSavePath("/explicit/玩家/save")
	if e != nil || p != "/explicit/玩家/save" {
		t.Fatalf("path=%q error=%v", p, e)
	}
}
