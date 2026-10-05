package recording

import "testing"

func TestWriterRejectsLossyUnicode(t *testing.T) {
	h := fixture(t)
	h.Config.Initial.DOS.FontFull = string([]byte{0xff})
	if _, err := NewHeader(h.Config); err == nil {
		t.Fatal("writer silently replaced invalid UTF-8")
	}
}
