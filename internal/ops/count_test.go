package ops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/walk"
)

func TestCopyFileCount(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "dst")

	// create a tree resembling the reported structure:
	// src/file1.txt, src/sub/file2.txt, src/sub/deep/file3.txt
	paths := []string{
		filepath.Join(srcDir, "file1.txt"),
		filepath.Join(srcDir, "sub", "file2.txt"),
		filepath.Join(srcDir, "sub", "deep", "file3.txt"),
	}
	for _, p := range paths {
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// verify pre-scan count matches
	stats, err := walk.ComputeStats([]string{srcDir})
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalFiles != 3 {
		t.Errorf("ComputeStats: expected 3 files, got %d", stats.TotalFiles)
	}

	// verify copy operation count matches
	var cumulFiles int64
	var cumulBytes int64
	ec := errs.NewCollector()
	opts := CopyOptions{
		Force:      false,
		Errors:     ec,
		CumulBytes: &cumulBytes,
		CumulFiles: &cumulFiles,
	}
	if err := Copy(srcDir, dstDir, opts); err != nil {
		t.Fatal(err)
	}
	if cumulFiles != 3 {
		t.Errorf("CumulFiles after copy: expected 3, got %d", cumulFiles)
	}
	if cumulFiles != stats.TotalFiles {
		t.Errorf("count mismatch: pre-scan=%d, actual=%d", stats.TotalFiles, cumulFiles)
	}
}

func TestCopyFileCountLargerTree(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "dst")

	// create exactly 10 files in a nested structure
	files := []string{
		"build/please",
		"justfile",
		"src/cli.odin",
		"src/main.odin",
		"src/pam/context.odin",
		"src/pam/ffi.odin",
		"src/pam/pam.odin",
		"src/ticket.odin",
		"src/timestamp.odin",
		"src/utils.odin",
	}
	for _, f := range files {
		p := filepath.Join(srcDir, f)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	stats, err := walk.ComputeStats([]string{srcDir})
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalFiles != 10 {
		t.Errorf("ComputeStats: expected 10 files, got %d", stats.TotalFiles)
	}

	var cumulFiles int64
	var cumulBytes int64
	ec := errs.NewCollector()
	opts := CopyOptions{
		Force:      false,
		Errors:     ec,
		CumulBytes: &cumulBytes,
		CumulFiles: &cumulFiles,
	}
	if err := Copy(srcDir, dstDir, opts); err != nil {
		t.Fatal(err)
	}
	if cumulFiles != 10 {
		t.Errorf("CumulFiles: expected 10, got %d", cumulFiles)
	}
}

func TestDeleteFileCount(t *testing.T) {
	dir := t.TempDir()

	// create dir with 1 file (the exact reported bug scenario)
	subdir := filepath.Join(dir, "something")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "meow"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	stats, err := walk.ComputeStats([]string{subdir})
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalFiles != 1 {
		t.Errorf("ComputeStats: expected 1 file, got %d", stats.TotalFiles)
	}

	// verify delete counts the file
	var cumulFiles int64
	ec := errs.NewCollector()
	opts := DeleteOptions{
		Recursive:  true,
		Errors:     ec,
		CumulFiles: &cumulFiles,
	}
	if err := Delete(subdir, opts); err != nil {
		t.Fatal(err)
	}
	// should count the file + subdir + root dir = 3 entries
	// but the summary shows filesDone which is cumulFiles
	// sendDeleteProgress increments for files, dirs, and root
	if cumulFiles < 1 {
		t.Errorf("CumulFiles after delete: expected >= 1, got %d", cumulFiles)
	}
}
