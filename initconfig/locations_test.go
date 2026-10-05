package initconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocationUsesIdentityAndClosedOutput(t *testing.T) {
	base := t.TempDir()
	root, scratch, out := filepath.Join(base, "root"), filepath.Join(base, "scratch"), filepath.Join(base, "out")
	for _, dir := range []string{root, scratch, out} {
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, bad := range []string{filepath.Join(root, "input"), filepath.Join(scratch, "input")} {
		if CheckLocation(bad, root, scratch, nil, true) == nil {
			t.Fatal("DOS-visible output accepted")
		}
	}
	alias := filepath.Join(base, "alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Skip("symlink unavailable")
	}
	if CheckLocation(filepath.Join(alias, "input"), root, scratch, nil, true) == nil {
		t.Fatal("directory alias accepted")
	}
	name := filepath.Join(out, "input")
	if err := CheckLocation(name, root, scratch, nil, true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name+".partial", []byte("self-authored fragment"), 0600); err != nil {
		t.Fatal(err)
	}
	if CheckLocation(name, root, scratch, nil, true) == nil {
		t.Fatal("existing partial accepted")
	}
	protected := filepath.Join(root, "RESOURCE")
	if err := os.WriteFile(protected, []byte("self-authored resource"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(protected, name); err != nil {
		t.Skip("hardlink unavailable")
	}
	if CheckLocation(name, root, scratch, []string{protected}, false) == nil {
		t.Fatal("protected inode alias accepted")
	}
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte("self-authored wire"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := CheckLocation(name, root, scratch, []string{protected}, false); err != nil {
		t.Fatal(err)
	}
	if err := CheckLocation(filepath.Join(out, "absent"), root, scratch, nil, false); err == nil {
		t.Fatal("missing replay input accepted")
	}
}
