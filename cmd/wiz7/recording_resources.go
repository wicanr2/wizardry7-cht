// 013 READY：同次讀取的顯示資源，同 bytes 驗證、解析及繪製。
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ic "github.com/wicanr2/wizardry7-cht/initconfig"
	"io"
	"log"
	"os"
	"path/filepath"
)

type recordingResources struct {
	metadata  []ic.Resource
	fault     bool
	afterRead func(string, string)
}
type resourceFile struct {
	*bytes.Reader
	info os.FileInfo
}

func (f *resourceFile) Stat() (os.FileInfo, error) { return f.info, nil }
func (f *resourceFile) Close() error               { return nil }

type recordingBundle struct {
	resources        *recordingResources
	trans            map[string]entry
	digits, ordinals map[string][]uint16
	flow             *flowData
	help             *helpData
	roster, sex      map[string]rosterEntry
	picture          *introPictureState
	brush            *brushMenuData
	paper            *introPaperData
}

func (r *recordingResources) index(role string) int {
	for i, v := range r.metadata {
		if v.Role == role {
			return i
		}
	}
	panic("resource_role")
}
func (r *recordingResources) read(role, path string) ([]byte, error) {
	f, e := r.openBytes(role, path, false)
	if e != nil {
		return nil, e
	}
	return io.ReadAll(f)
}
func (r *recordingResources) open(role, path string) (*resourceFile, error) {
	return r.openBytes(role, path, true)
}
func (r *recordingResources) openBytes(role, path string, sidecar bool) (*resourceFile, error) {
	m := &r.metadata[r.index(role)]
	m.Requested = path != ""
	if path == "" {
		return nil, os.ErrNotExist
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() {
		r.fault = true
		return nil, errors.New("resource_capture")
	}
	if sidecar && st.Size() > 1<<20 {
		r.fault = true
		return nil, errors.New("resource_capture")
	}
	b, e := io.ReadAll(io.LimitReader(f, st.Size()+1))
	if e != nil || int64(len(b)) != st.Size() {
		r.fault = true
		return nil, errors.New("resource_capture")
	}
	m.Size = uint64(len(b))
	m.SHA256 = ic.Hash(b)
	if r.afterRead != nil {
		r.afterRead(role, path)
	}
	return &resourceFile{bytes.NewReader(b), st}, nil
}
func (r *recordingResources) loaded(role string) { r.metadata[r.index(role)].Loaded = true }
func loadRecordingBundle(xlate, flow, help string) (recordingBundle, error) {
	r := &recordingResources{}
	for _, role := range []string{"xlate", "flow", "help", "roster", "profile-sex", "intro-picture"} {
		r.metadata = append(r.metadata, ic.Resource{Role: role})
	}
	b := recordingBundle{resources: r, trans: map[string]entry{}, digits: map[string][]uint16{}, ordinals: map[string][]uint16{}}
	raw, e := r.read("xlate", xlate)
	if e == nil {
		decoded, err := decodeXlate(raw)
		if err != nil {
			return b, errors.New("resource_parse")
		}
		b.trans, b.digits, b.ordinals = decoded.trans, decoded.digits, decoded.ordinals
		b.brush, b.paper = decoded.brush, decoded.paper
		r.loaded("xlate")
	}
	b.flow = r.loadFlow(flow)
	if b.flow != nil {
		r.loaded("flow")
	}
	b.help = r.loadHelp(help)
	if b.help != nil {
		r.loaded("help")
	}
	b.picture = r.loadIntroPicture(xlate, b.trans)
	b.picture.paper = b.paper
	if b.picture.data != nil {
		r.loaded("intro-picture")
	}
	b.roster = r.loadRoster(xlate, b.trans)
	if b.roster != nil {
		r.loaded("roster")
	}
	b.sex = r.loadProfileSex(xlate, b.trans)
	if b.sex != nil {
		r.loaded("profile-sex")
	}
	if r.fault {
		return b, errors.New("resource_capture")
	}
	return b, nil
}
func (r *recordingResources) loadFlow(path string) *flowData {
	if path == "" {
		log.Print("flow：資源檢查")
		return nil
	}
	b, err := r.read("flow", path)
	if err != nil {
		log.Print("flow：資源檢查")
		return nil
	}
	var fd flowData
	if err := json.Unmarshal(b, &fd); err != nil {
		log.Print("flow：資源檢查")
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return nil
	}
	_, explicitSize := fields["size"]
	if (explicitSize && fd.Size != 15 && fd.Size != 16) || (fd.size() != 15 && fd.size() != 16) {
		log.Print("flow：資源檢查")
		return nil
	}
	for _, rows := range fd.Glyphs {
		if len(rows) != fd.size() {
			log.Print("flow：資源檢查")
			return nil
		}
		for _, row := range rows {
			if uint32(row) >= uint32(1)<<fd.size() {
				log.Print("flow：資源檢查")
				return nil
			}
		}
	}
	fd.byWord = map[string][]flowPos{}
	for pi, p := range fd.Paras {
		for a, w := range p.Words {
			fd.byWord[w] = append(fd.byWord[w], flowPos{int32(pi), int32(a)})
		}
	}
	log.Print("flow：資源檢查")
	return &fd
}

func (r *recordingResources) loadHelp(path string) *helpData {
	b, err := r.read("help", path)
	if err != nil {
		log.Print("help：資源檢查")
		return nil
	}
	var h helpData
	if err := json.Unmarshal(b, &h); err != nil {
		log.Print("help：資源檢查")
		return nil
	}
	for _, g := range h.Glyphs {
		if (g.W != 8 && g.W != 15 && g.W != 16) || (len(g.Rows) != 15 && len(g.Rows) != 16) || (g.W != 8 && len(g.Rows) != g.W) {
			log.Print("help：資源檢查")
			return nil
		}
		for _, row := range g.Rows {
			if row < 0 || row >= 1<<g.W {
				log.Print("help：資源檢查")
				return nil
			}
		}
	}
	log.Print("help：資源檢查")
	return &h
}

func (r *recordingResources) loadRoster(xlate string, trans map[string]entry) map[string]rosterEntry {
	if xlate == "" || len(trans) == 0 {
		return nil
	}
	size := 0
	for _, e := range trans {
		if len(e.G) > 0 {
			size = len(e.G[0])
			break
		}
	}
	f, err := r.open("roster", filepath.Join(filepath.Dir(xlate), "roster.json"))
	if err == nil {
		defer f.Close()
		var st os.FileInfo
		st, err = f.Stat()
		if err == nil && st.Size() > 1<<20 {
			err = fmt.Errorf("職業碼資源過大")
		}
	}
	if err == nil {
		var entries map[string]rosterEntry
		entries, err = decodeRoster(f, size)
		if err == nil {
			log.Print("職業碼：載入六鍵")
			return entries
		}
	}
	log.Print("roster：資源檢查")
	return nil
}

func (r *recordingResources) loadProfileSex(xlate string, trans map[string]entry) map[string]rosterEntry {
	if xlate == "" || len(trans) == 0 {
		return nil
	}
	size := 0
	for _, e := range trans {
		for _, glyph := range e.G {
			h := len(glyph)
			if h != 15 && h != 16 || size != 0 && h != size {
				log.Print("性別資源停用：一般字模高度不一致")
				return nil
			}
			size = h
		}
	}
	f, err := r.open("profile-sex", filepath.Join(filepath.Dir(xlate), "profile-sex.json"))
	if err == nil {
		defer f.Close()
		var entries map[string]rosterEntry
		entries, err = decodeProfileSex(f, size)
		if err == nil {
			log.Print("性別資源：載入兩鍵")
			return entries
		}
	}
	log.Print("profile-sex：資源檢查")
	return nil
}

func (r *recordingResources) loadIntroPicture(path string, trans map[string]entry) *introPictureState {
	s := &introPictureState{}
	if path == "" || len(trans) == 0 {
		s.english = true
		return s
	}
	size := 0
	for _, e := range trans {
		for _, g := range e.G {
			h := len(g)
			if h != 15 && h != 16 || size != 0 && size != h {
				s.report("一般字模高度不一致")
				return s
			}
			size = h
		}
	}
	f, err := r.open("intro-picture", filepath.Join(filepath.Dir(path), "intro-picture.json"))
	if err == nil {
		defer f.Close()
		s.data, err = decodeIntroPicture(f, size)
	}
	if err != nil {
		s.report("資源不可用")
		return s
	}
	log.Print("片頭圖片資源：載入呈獻")
	return s
}
