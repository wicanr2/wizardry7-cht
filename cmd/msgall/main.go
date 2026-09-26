// msgall — MSG.DBS 全文抽取（#19 定案版）。
//
// 演算法＝遊戲 sub_10582/sub_13E02/sub_1045D/sub_104C6 的離線重做
// （docs/re/013 追記三十）：
//
//	HDR 項（6B）：[+0]w 群首訊息ID、[+2]w 塊內偏移、[+4]b 群內訊息數、[+5]b 塊號
//	記錄鏈：每筆 [L 壓縮長][D1 明文長][bitstream×(L-1)]，next = off + 1 + L
//	解壓：block[0]＝輸出長度、bitstream 自 [1] 起 MSB-first 連續、
//	      樹 4-byte 節點（左 word +0、右 word +2、內部 neg×4＝byte 偏移、葉輸出回根）
//
//	msgall -hdr <MSG.HDR> -dbs <MSG.DBS> -tree <tree.bin> -out <messages.json>
//
// 樹檔是 runtime dump（workplace/bake/trees/tree-73cb.bin，gitignore）；
// 本工具只產出原文表（原文欄位依使用者定案進版控）。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func word(t []byte, off int) uint16 {
	if off < 0 || off+1 >= len(t) {
		return 0x8000
	}
	return uint16(t[off]) | uint16(t[off+1])<<8
}

// decompress 解一筆記錄（rec[0]＝明文長度 D1，bitstream 自 rec[1] 起）。
func decompress(rec []byte, tree []byte) []byte {
	if len(rec) == 0 {
		return nil
	}
	want := int(rec[0])
	if want == 0 {
		return []byte{}
	}
	bx := 0
	al := 0
	cx := 8 // rec[1] 已載入 al，尚有 8 bits
	i := 1
	if len(rec) > 1 {
		al = int(rec[1])
		i = 2
	}
	out := []byte{}
	for len(out) < want {
		if i > len(rec)+1 {
			return nil // 位元流耗盡（不該發生）
		}
		if cx == 0 {
			cx = 8
			if i < len(rec) {
				al = int(rec[i])
			}
			i++
		}
		bit := (al >> 7) & 1
		al = (al << 1) & 0xFF
		cx--
		v := word(tree, bx+(bit*2))
		if v&0x8000 != 0 {
			bx = ((-int(v)) & 0xFFFF) * 4
		} else {
			out = append(out, byte(v))
			bx = 0
		}
	}
	return out
}

type entry struct {
	First int `json:"first"`
	Off   int `json:"off"`
	Span  int `json:"span"`
	Blk   int `json:"blk"`
}

type msg struct {
	ID  int    `json:"id"`
	Blk int    `json:"blk"`
	Off int    `json:"off"`
	En  string `json:"en"`
}

func main() {
	hdrPath := flag.String("hdr", "", "MSG.HDR")
	dbsPath := flag.String("dbs", "", "MSG.DBS")
	treePath := flag.String("tree", "", "tree.bin（runtime dump）")
	outPath := flag.String("out", "", "輸出 JSON")
	flag.Parse()
	if *hdrPath == "" || *dbsPath == "" || *treePath == "" || *outPath == "" {
		log.Fatal("用法: msgall -hdr MSG.HDR -dbs MSG.DBS -tree tree.bin -out messages.json")
	}

	hdr, err := os.ReadFile(*hdrPath)
	if err != nil {
		log.Fatal(err)
	}
	dbs, err := os.ReadFile(*dbsPath)
	if err != nil {
		log.Fatal(err)
	}
	tree, err := os.ReadFile(*treePath)
	if err != nil {
		log.Fatal(err)
	}
	if len(hdr) < 2 || (len(hdr)-2)%6 != 0 {
		log.Fatalf("MSG.HDR 大小異常: %d", len(hdr))
	}
	nEntries := (len(hdr) - 2) / 6
	count := int(uint16(hdr[0]) | uint16(hdr[1])<<8) // ds:100Ch＝遊戲 binary search 上界
	log.Printf("HDR: 檔大小 %d、項數 %d、遊戲計數 %d", len(hdr), nEntries, count)

	entries := make([]entry, 0, nEntries)
	for i := 0; i < nEntries; i++ {
		b := hdr[2+i*6 : 2+i*6+6]
		entries = append(entries, entry{
			First: int(uint16(b[0]) | uint16(b[1])<<8),
			Off:   int(uint16(b[2]) | uint16(b[3])<<8),
			Span:  int(b[4]),
			Blk:   int(b[5]),
		})
	}

	// 仿 sub_1045D/sub_10582：以訊息 ID binary search 找群；遊戲上界＝count（項 0..count-1）。
	// 覆蓋範圍＝[first, first+span]（含端點：found 判準 first+span >= msgID）。
	// 記錄鏈可跨塊：sub_10582 邊界處理＝載 blk+1 接著讀（loc_105D9 之後）。
	seen := map[int]bool{}
	msgs := []msg{}
	stat := struct{ groups, empty, hiBit, dup, cross int }{}
	maxID := 0
	for _, e := range entries[:count] {
		if e.Blk*1024+1024 > len(dbs) {
			continue
		}
		stat.groups++
		// 虛擬視窗＝本塊＋下一塊（記錄跨 0x400 邊界用）
		blkOff := e.Blk * 1024
		nxtOff := blkOff + 1024
		if nxtOff+1024 > len(dbs) {
			nxtOff = blkOff // 最後一塊沒有下一塊
		}
		window := dbs[blkOff : blkOff+1024]
		if nxtOff != blkOff {
			window = dbs[blkOff : nxtOff+1024]
		}
		o := e.Off
		for id := e.First; id <= e.First+e.Span; id++ {
			if o < 0 || o >= len(window) {
				break
			}
			L := int(window[o])
			if L == 0 || o+1+L > len(window) {
				break
			}
			if o+1+L > 1024 {
				stat.cross++
			}
			if seen[id] {
				stat.dup++
			} else {
				seen[id] = true
				txt := decompress(window[o+1:o+1+L], tree)
				if txt == nil {
					log.Fatalf("id %d: 位元流耗盡（演算法或樹不符）", id)
				}
				for _, c := range txt {
					if c >= 127 {
						stat.hiBit++
						break
					}
				}
				if len(txt) == 0 {
					stat.empty++
				}
				msgs = append(msgs, msg{ID: id, Blk: e.Blk, Off: o, En: string(txt)})
				if id > maxID {
					maxID = id
				}
			}
			o += 1 + L
		}
	}
	log.Printf("群 %d、訊息 %d（最大ID %d）、空 %d、含高位byte %d、重複 %d、跨塊 %d",
		stat.groups, len(msgs), maxID, stat.empty, stat.hiBit, stat.dup, stat.cross)

	w, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(map[string]any{
		"schema": "wiz7-messages/1",
		"note": "MSG.DBS 全量抽取（docs/re/013 追記三十）。id＝訊息ID（sub_10582 參數）、" +
			"blk＝MSG.DBS 塊號、off＝塊內記錄偏移、en＝原文（含控制碼，印字管線解讀）。",
		"tree":    *treePath,
		"count":   len(msgs),
		"maxID":   maxID,
		"records": msgs,
	}); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("done: %d 則 → %s\n", len(msgs), *outPath)
}
