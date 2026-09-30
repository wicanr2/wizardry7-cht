package main

import (
	"encoding/json"
	"os"
	"testing"
)

// spec 003 R6 驗收 6：R4 辨識規則對實測 shown 事件的判定。資料在 workplace/（gitignore），缺檔就 skip。
func TestDetectScreenOnRecordedShown(t *testing.T) {
	cases := []struct {
		path, want string
	}{
		{"../../workplace/r663/zh3/hdat-272000000.json", "copyprot"},
		{"../../workplace/r663/zh3/hdat-286000000.json", "copyprot"},
		{"../../workplace/r665/maze/hdat-1470000000.json", "general"},
	}
	for _, c := range cases {
		b, err := os.ReadFile(c.path)
		if err != nil {
			t.Skipf("缺實測資料 %s", c.path)
		}
		var j struct{ Shown []struct{ S string } }
		if err := json.Unmarshal(b, &j); err != nil {
			t.Fatal(err)
		}
		var strs []string
		for _, e := range j.Shown {
			strs = append(strs, e.S)
		}
		if got := detectScreen(strs); got != c.want {
			t.Errorf("%s：判為 %s，應為 %s", c.path, got, c.want)
		}
	}
}

// 規則順序：主選單字串與防拷字串同時存在時必須判為防拷（docs/re/015 H6）。
func TestDetectScreenOrder(t *testing.T) {
	if got := detectScreen([]string{"CHARACTER MENU", "START NEW GAME", "SECURITY CLEARANCE"}); got != "copyprot" {
		t.Errorf("判為 %s", got)
	}
	if got := detectScreen([]string{"CHARACTER MENU", "ADD CHARACTER", "OK", "CANCEL"}); got != "addchar" {
		t.Errorf("判為 %s", got)
	}
	if got := detectScreen([]string{"CHARACTER MENU"}); got != "mainmenu" {
		t.Errorf("判為 %s", got)
	}
}
