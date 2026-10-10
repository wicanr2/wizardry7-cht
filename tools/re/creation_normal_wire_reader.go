package main

// 只用正式Reader驗完整錄製，絕不輸出header／輸入內容或執行原版。
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	rc "github.com/wicanr2/wizardry7-cht/recording"
)

type safeWireSummary struct {
	Status                         string `json:"status"`
	Code                           string `json:"code,omitempty"`
	WireSHA256                     string `json:"wire_sha256"`
	ReaderAccepted                 bool   `json:"reader_accepted"`
	HeaderFooterSchemaHashVerified bool   `json:"header_footer_schema_hash_verified"`
	EventCount                     int    `json:"event_count"`
	PublicationCount               uint64 `json:"publication_count"`
	EndStep                        uint64 `json:"end_step"`
	LastSnapshotStep               uint64 `json:"last_snapshot_step"`
	RuntimeExecuted                bool   `json:"runtime_executed"`
}

func main() {
	wire := flag.String("wire", "", "完整錄製路徑")
	expected := flag.String("sha256", "", "固定wire SHA")
	out := flag.String("output", "", "新安全summary路徑")
	flag.Parse()
	result := safeWireSummary{Status: "FAIL", Code: "reader_rejected"}
	valid := func() bool {
		stat, err := os.Lstat(*wire)
		if err != nil || !stat.Mode().IsRegular() || stat.Size() <= 0 || stat.Size() > rc.MaxBytes {
			return false
		}
		f, err := os.Open(*wire)
		if err != nil {
			return false
		}
		h := sha256.New()
		_, readErr := io.Copy(h, io.LimitReader(f, rc.MaxBytes+1))
		closeErr := f.Close()
		if readErr != nil || closeErr != nil {
			return false
		}
		result.WireSHA256 = hex.EncodeToString(h.Sum(nil))
		if result.WireSHA256 != *expected {
			result.Code = "wire_hash_mismatch"
			return false
		}
		doc, err := rc.Load(*wire)
		if err != nil {
			return false
		}
		var publications uint64
		var lastStep uint64
		for _, event := range doc.Events {
			if event.Kind == "snapshot" {
				if event.Params.Publication == nil || *event.Params.Publication != publications+1 {
					return false
				}
				publications = *event.Params.Publication
				lastStep = event.After
			}
		}
		if uint64(len(doc.Events)) != doc.Footer.Count || lastStep != doc.Footer.End || !doc.Footer.Accepted {
			return false
		}
		result = safeWireSummary{Status: "PASS", WireSHA256: result.WireSHA256, ReaderAccepted: true,
			HeaderFooterSchemaHashVerified: true, EventCount: len(doc.Events), PublicationCount: publications,
			EndStep: doc.Footer.End, LastSnapshotStep: lastStep}
		return true
	}()
	f, err := os.OpenFile(*out, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "NORMAL_READER_OUTPUT_ERROR")
		os.Exit(1)
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		b = append(b, '\n')
		_, err = f.Write(b)
	}
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		fmt.Fprintln(os.Stderr, "NORMAL_READER_OUTPUT_ERROR")
		os.Exit(1)
	}
	fmt.Println(string(b))
	if !valid {
		os.Exit(1)
	}
}
