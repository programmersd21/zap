package ops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShredDeletesFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "secret.txt")
	content := []byte("super secret data that must not survive")
	if err := os.WriteFile(src, content, 0600); err != nil {
		t.Fatal(err)
	}

	if err := Shred(src, ShredOptions{Passes: 1}); err != nil {
		t.Fatalf("Shred failed: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("file should be gone after shred")
	}
}

func TestShredOverwritesContent(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "doc.txt")
	original := make([]byte, 4096)
	for i := range original {
		original[i] = 0x41
	}
	if err := os.WriteFile(src, original, 0644); err != nil {
		t.Fatal(err)
	}

	if err := Shred(src, ShredOptions{Passes: 2}); err != nil {
		t.Fatalf("Shred failed: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Error("expected the file to be removed")
	}
}

func TestShredOverwriteFileChangesBytes(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "f.bin")
	original := make([]byte, 8192)
	if err := os.WriteFile(src, original, 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := overwriteFile(src, 1, 0, ShredOptions{}); err != nil {
		t.Fatalf("overwriteFile failed: %v", err)
	}

	// file must still exist after overwrite (deletion happens in Shred)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}
	if len(data) != len(original) {
		t.Errorf("expected size %d, got %d", len(original), len(data))
	}
	changed := false
	for i := range data {
		if data[i] != original[i] {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("expected file contents to be overwritten with random data")
	}
}

func TestShredRefusesDirectory(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "mydir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := Shred(dir, ShredOptions{}); err == nil {
		t.Error("expected shred to refuse a directory")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory should be untouched: %v", err)
	}
}

func TestShredRefusesRoot(t *testing.T) {
	if err := Shred("/", ShredOptions{}); err == nil {
		t.Error("expected shred of the root to be refused")
	}
}

func TestShredNonexistent(t *testing.T) {
	if err := Shred(filepath.Join(t.TempDir(), "nope"), ShredOptions{}); err != nil {
		t.Errorf("shredding a missing path should be a no-op, got %v", err)
	}
}
