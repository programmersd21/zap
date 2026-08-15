package ops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/programmersd21/zap/internal/errs"
)

func TestCopySamePathRejects(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "foo.txt")
	if err := os.WriteFile(file, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	ec := errs.NewCollector()
	err := Copy(file, file, CopyOptions{Force: true, Errors: ec})
	if err == nil {
		t.Fatal("expected error copying file to itself")
	}

	// verify file is untouched
	data, _ := os.ReadFile(file)
	if string(data) != "original" {
		t.Errorf("file was modified: got %q", data)
	}
}

func TestCopySamePathRelative(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "foo.txt")
	if err := os.WriteFile(file, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	// ./foo.txt vs foo.txt equivalent
	relSrc := file
	relDst := filepath.Join(filepath.Dir(file), ".", "foo.txt")

	ec := errs.NewCollector()
	err := Copy(relSrc, relDst, CopyOptions{Force: true, Errors: ec})
	if err == nil {
		t.Fatal("expected error copying file to itself via relative path")
	}

	data, _ := os.ReadFile(file)
	if string(data) != "original" {
		t.Errorf("file was modified: got %q", data)
	}
}

func TestMoveSamePathRejects(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "bar.txt")
	if err := os.WriteFile(file, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	ec := errs.NewCollector()
	err := Move(file, file, MoveOptions{Force: true, Errors: ec})
	if err == nil {
		t.Fatal("expected error moving file to itself")
	}

	data, _ := os.ReadFile(file)
	if string(data) != "original" {
		t.Errorf("file was destroyed: got %q", data)
	}
}
