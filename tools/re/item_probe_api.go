// 私人驗收工具：沿現有前端派送，額外只讀原版雜湊。
package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/wicanr2/dosgolem/oracle"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
	rc "github.com/wicanr2/wizardry7-cht/recording"
)

var r982Trace *json.Encoder

func r982Control(o *oracle.Oracle, root, scratch string) *inputRecording {
	name := os.Getenv("WIZ7_R982_API_DOCUMENT")
	if name == "" {
		return nil
	}
	doc, err := rc.Load(name)
	if err != nil {
		panic("wire_rejected")
	}
	if err = ic.VerifyStorage(doc.Header.Config.Storage, root, scratch); err != nil {
		panic("storage_mismatch")
	}
	c := doc.Header.Config.Initial
	clk := o.ClockSettings()
	if o.Steps() != c.Step || o.Ticks() != c.Ticks || clk.CPUHz != c.Clock.CPUHz || clk.CycleClock != c.Clock.CycleClock || clk.IRQ0Base != c.Clock.IRQ0Base || clk.IRQ0Every != c.Clock.IRQ0Every || !reflect.DeepEqual(o.DOSSettings(), c.DOS) {
		panic("initial_scalars_mismatch")
	}
	output := os.Getenv("WIZ7_R982_TRACE")
	file, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	r982Trace = json.NewEncoder(file)
	return &inputRecording{replay: true, document: doc, configSHA: doc.Header.ConfigSHA, acks: make(chan recordingPublication, 1)}
}
func r982Observe(o *oracle.Oracle, l *inputRecording) {
	if r982Trace == nil || l == nil {
		return
	}
	v := recordingOriginalHash(o)
	row := struct {
		Publication uint64 `json:"publication"`
		Step        uint64 `json:"step"`
		Ticks       uint64 `json:"ticks"`
		RAM         string `json:"ram"`
		Indexed     string `json:"indexed"`
		Palette     string `json:"palette"`
		Registers   string `json:"registers"`
	}{l.publication + 1, v.Step, v.Ticks, v.RAM, v.Indexed, v.Palette, v.Registers}
	if err := r982Trace.Encode(row); err != nil {
		panic(err)
	}
}
