package main

// R1003 V1–V5：只在既有owner快照抽五發布，Draw不讀Oracle。
import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/wicanr2/dosgolem/oracle"
)

var creationCaptureState struct {
	sync.Mutex
	stopped bool
	seen map[uint64]bool
}
var creationCaptureSteps = map[uint64]uint64{1519:332800000,1520:338800000,1585:363400000,1686:406800000,1712:417800000}

func creationCaptureFail(code string) {
	creationCaptureState.stopped = true
	fmt.Fprintln(os.Stderr, "R1003_CAPTURE_INVALID_"+code)
}
func creationCaptureJSON(path string, value any) bool {
	b, err := json.Marshal(value)
	if err != nil { return false }
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil { return false }
	b = append(b, '\n')
	n, err := f.Write(b)
	closeErr := f.Close()
	return err == nil && closeErr == nil && n == len(b)
}
func creationCapturePNG(path string, k int) bool {
	f, err := os.Open(path)
	if err != nil { return false }
	img, decodeErr := png.Decode(f)
	closeErr := f.Close()
	return decodeErr == nil && closeErr == nil && img.Bounds().Dx()==320*k && img.Bounds().Dy()==200*k
}
func creationCapture(o *oracle.Oracle, ov *overlay, am *automap, sh *shared, publication recordingPublication, k int) {
	root := os.Getenv("WIZ7_R1003_CAPTURE")
	if root == "" || publication.ID == 0 { return }
	pub := publication.ID
	expected, selected := creationCaptureSteps[pub]
	if !selected { return }
	creationCaptureState.Lock()
	defer creationCaptureState.Unlock()
	if creationCaptureState.stopped { return }
	if creationCaptureState.seen == nil { creationCaptureState.seen = map[uint64]bool{} }
	if creationCaptureState.seen[pub] { return }
	if len(creationCaptureState.seen)>=5 || o.Steps()!=expected || publication.Step!=expected || sh.recordingPublication.ID!=pub ||
		publication.Panel.HelpOpen || publication.Panel.MapOpen || sh.helpOpen!=publication.Panel.HelpOpen || sh.mapOpen!=publication.Panel.MapOpen || (k!=2 && k!=3) {
		creationCaptureFail("IDENTITY"); return
	}
	creationCaptureState.seen[pub] = true
	dir := filepath.Join(root,fmt.Sprintf("frame-%d",pub))
	if os.Mkdir(dir,0700)!=nil { creationCaptureFail("DIRECTORY"); return }
	hdAtDump(o,ov,am,dir,"hd")
	for name, size := range map[string]int{"hd.a000":64000,"hd.comp":64000,"hd.pal":768} {
		b, err := os.ReadFile(filepath.Join(dir,name))
		if err!=nil || len(b)!=size { creationCaptureFail("HD_FILE"); return }
	}
	b, err := os.ReadFile(filepath.Join(dir,"hd.json"))
	var hd struct{ Step uint64 }
	if err!=nil || json.Unmarshal(b,&hd)!=nil || hd.Step!=expected || !creationCapturePNG(filepath.Join(dir,"hd.png"),k) {
		creationCaptureFail("HD_FORMAT"); return
	}
	type sexEvent struct {
		Return oracle.Addr
		S string
		X,Y,BX,BY,BW,BH int
		Before,After []byte
		G [][]uint16
		Shown bool
	}
	var sex []sexEvent
	for _, e := range ov.evs {
		if !e.sex || e.dead { continue }
		shown := false
		for _, displayed := range sh.evs { if displayed.sex && displayed.originRet==e.originRet && displayed.s==e.s && displayed.x0==e.x0 && displayed.y0==e.y0 { shown=true;break } }
		sex=append(sex,sexEvent{e.originRet,e.s,e.x0,e.y0,e.bx,e.by,e.bw,e.bh,e.before,e.after,e.g,shown})
	}
	if !creationCaptureJSON(filepath.Join(dir,"capture.json"),map[string]any{"publication":pub,"step":expected,"scale":k,"sex":sex,"help_open":sh.helpOpen,"map_open":sh.mapOpen}) {
		creationCaptureFail("METADATA"); return
	}
	sh.drawReq=filepath.Join(dir,"draw.png")
}
func creationCaptureDraw(req string, publication recordingPublication, k int) {
	root:=os.Getenv("WIZ7_R1003_CAPTURE")
	if root=="" || req=="" { return }
	creationCaptureState.Lock()
	defer creationCaptureState.Unlock()
	if creationCaptureState.stopped { return }
	expected, selected:=creationCaptureSteps[publication.ID]
	dir:=filepath.Join(root,fmt.Sprintf("frame-%d",publication.ID))
	if !selected || !creationCaptureState.seen[publication.ID] || publication.Step!=expected || req!=filepath.Join(dir,"draw.png") || !creationCapturePNG(req,k) {
		creationCaptureFail("DRAW_PAIR"); return
	}
	if !creationCaptureJSON(filepath.Join(dir,"draw.json"),map[string]any{"publication":publication.ID,"step":publication.Step,"scale":k}) { creationCaptureFail("DRAW_METADATA") }
}
