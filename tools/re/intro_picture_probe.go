// 片頭唯讀觀察：只記開檔名稱、blit入口數值及原版畫格。
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/oracle"
)

type snapshot struct {
	Step                                       uint64
	Mem, A000, Pal                             uint32
	MemorySHA256, IndexedSHA256, PaletteSHA256 string
	Mouse                                      [2]int
	Mode                                       uint8
}
type openEvent struct {
	Step   uint64
	Name   string
	CS, IP uint16
}
type blitEvent struct {
	Step                                       uint64
	CS, IP, SP, AX, BX, CX, DX, SI, DI, DS, ES uint16
	StackWords                                 []uint16
}

func checked(e error) {
	if e != nil {
		panic(e)
	}
}
func hash(b []byte) uint32 { h := fnv.New32a(); h.Write(b); return h.Sum32() }
func sha(b []byte) string  { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func write(p string, b []byte) {
	f, e := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	checked(e)
	_, e = f.Write(b)
	checked(e)
	checked(f.Close())
}

type renderEvent struct {
	Step   uint64
	IP, SP uint16
	Kind   string
	Words  []uint16
}

type picOpen struct {
	Step           uint64
	Name           string
	SP, BP, SS     uint16
	Stack, BPWords []uint16
	FrameChain     [][]uint16
}

func main() {
	if len(os.Args) != 4 {
		panic("original-dir new-output-dir on|off")
	}
	root, out, mode := os.Args[1], os.Args[2], os.Args[3]
	if mode != "on" && mode != "off" {
		panic("bad observer mode")
	}
	checked(os.Mkdir(out, 0755))
	checked(os.Mkdir(filepath.Join(out, "saves"), 0755))
	o, err := oracle.LoadWith(filepath.Join(root, "DS.EXE"), root, oracle.Options{Adlib: true})
	checked(err)
	o.SetScratch(filepath.Join(out, "saves"))
	var renders []renderEvent
	var vinitActive = true
	var pendingSP uint16
	var pendingReturn uint16
	var pictureCaptures []snapshot
	capturePicture := func() {
		a := append([]byte(nil), o.Indexed()...)
		var pal []byte
		colors := o.Palette()
		for _, c := range colors {
			pal = append(pal, c[:]...)
		}
		mem := o.Bytes(oracle.Addr{}, 1<<20)
		x, y := o.Mouse()
		pictureCaptures = append(pictureCaptures, snapshot{o.Steps(), hash(mem), hash(a), hash(pal), sha(mem), sha(a), sha(pal), [2]int{x, y}, o.VideoMode()})
		base := filepath.Join(out, fmt.Sprintf("picture-%d", o.Steps()))
		write(base+".a000", a)
		write(base+".pal", pal)
		im := image.NewRGBA(image.Rect(0, 0, 320, 200))
		for i, v := range a {
			c := colors[v]
			im.SetRGBA(i%320, i/320, color.RGBA{c[0], c[1], c[2], 255})
		}
		f, e := os.OpenFile(base+".png", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		checked(e)
		checked(png.Encode(f, im))
		checked(f.Close())
	}
	var opens []openEvent
	var picOpens []picOpen
	printCalls := map[string]int{"strings": 0, "characters": 0}
	var blits []blitEvent
	if mode == "on" {
		for _, item := range []struct {
			ip    uint16
			kind  string
			words int
		}{{0x2663, "picture-load", 3}, {0x3ff1, "picture-load-service", 3}, {0x5a1a, "vinit-picture", 3}, {0x4018, "picture-render-service", 12}, {0x3921, "screen-blit", 5}} {
			item := item
			o.OnCall(oracle.Addr{Seg: 0x0110, Off: item.ip}, func(o *oracle.Oracle) {
				if !vinitActive {
					return
				}
				if len(renders) >= 512 {
					panic("picture event bound")
				}
				words := make([]uint16, item.words)
				for i := range words {
					words[i] = o.StackWord(i)
				}
				renders = append(renders, renderEvent{o.Steps(), item.ip, o.Regs().SP, item.kind, words})
				if item.ip == 0x3921 && (words[0] == 0x5b9b || words[0] == 0x5afd) {
					pendingSP = o.Regs().SP
					pendingReturn = words[0]
				}
			})
		}
		for _, ret := range []uint16{0x5b9b, 0x5afd} {
			ret := ret
			o.OnCall(oracle.Addr{Seg: 0x0110, Off: ret}, func(o *oracle.Oracle) {
				if vinitActive && pendingReturn == ret && o.Regs().SP == pendingSP+2 {
					capturePicture()
					pendingSP = 0
					pendingReturn = 0
				}
			})
		}
		o.OnFileOpen(func(name string) {
			if filepath.Ext(name) == ".OVR" && name != "VINIT.OVR" {
				vinitActive = false
			}
			ip := o.IP()
			opens = append(opens, openEvent{o.Steps(), name, ip.Seg, ip.Off})
			if name == "MON60.PIC" || name == "MON61.PIC" {
				r := o.Regs()
				stack, bp := make([]uint16, 16), make([]uint16, 10)
				for i := range stack {
					stack[i] = o.StackWord(i)
				}
				for i := range bp {
					bp[i] = o.Word(oracle.Addr{Seg: r.SS, Off: r.BP + uint16(i*2)})
				}
				var chain [][]uint16
				cursor := r.BP
				for frame := 0; frame < 8 && cursor >= r.SP && cursor <= 65520; frame++ {
					row := []uint16{cursor}
					for i := 0; i < 6; i++ {
						row = append(row, o.Word(oracle.Addr{Seg: r.SS, Off: cursor + uint16(i*2)}))
					}
					chain = append(chain, row)
					if row[1] <= cursor || row[1] > 65520 {
						break
					}
					cursor = row[1]
				}
				picOpens = append(picOpens, picOpen{o.Steps(), name, r.SP, r.BP, r.SS, stack, bp, chain})
			}
		})
		o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3895}, func(o *oracle.Oracle) { printCalls["strings"]++ })
		o.OnCall(oracle.Addr{Seg: 0x0110, Off: 0x3877}, func(o *oracle.Oracle) { printCalls["characters"]++ })
		o.OnCall(oracle.Addr{Seg: 0x740c, Off: 0x05b5}, func(o *oracle.Oracle) {
			if len(blits) >= 10000 {
				panic("blit event bound")
			}
			r := o.Regs()
			ip := o.IP()
			words := make([]uint16, 12)
			for i := range words {
				words[i] = o.StackWord(i)
			}
			blits = append(blits, blitEvent{o.Steps(), ip.Seg, ip.Off, r.SP, r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.DS, r.ES, words})
		})
	}
	var snapshots []snapshot
	for step := uint64(2000000); step <= 200000000; step += 2000000 {
		checked(o.Run(step - o.Steps()))
		if o.Steps() != step {
			panic("step overrun")
		}
		a := append([]byte(nil), o.Indexed()...)
		var pal []byte
		colors := o.Palette()
		for _, c := range colors {
			pal = append(pal, c[:]...)
		}
		mem := o.Bytes(oracle.Addr{}, 1<<20)
		x, y := o.Mouse()
		snapshots = append(snapshots, snapshot{step, hash(mem), hash(a), hash(pal), sha(mem), sha(a), sha(pal), [2]int{x, y}, o.VideoMode()})
		base := filepath.Join(out, fmt.Sprintf("original-%d", step))
		write(base+".a000", a)
		write(base+".pal", pal)
		im := image.NewRGBA(image.Rect(0, 0, 320, 200))
		for i, v := range a {
			c := colors[v]
			im.SetRGBA(i%320, i/320, color.RGBA{c[0], c[1], c[2], 255})
		}
		f, e := os.OpenFile(base+".png", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		checked(e)
		checked(png.Encode(f, im))
		checked(f.Close())
	}
	report := struct {
		Renders         []renderEvent
		PictureCaptures []snapshot
		Observer        string
		PrintCalls      map[string]int
		PicOpens        []picOpen
		Opened          []openEvent
		Blits           []blitEvent
		Snapshots       []snapshot
		Limits          string
	}{renders, pictureCaptures, mode, printCalls, picOpens, opens, blits, snapshots, "正常冷啟動無玩家輸入；只觀察檔名與繪圖數值；200M截止；blit stack未賦予語意；非正式中文GUI驗收。"}
	b, e := json.MarshalIndent(report, "", "  ")
	checked(e)
	write(filepath.Join(out, "observation.json"), append(b, '\n'))
	fmt.Println(mode, "complete", len(opens), "file opens", len(blits), "blit calls", len(snapshots), "snapshots")
}
