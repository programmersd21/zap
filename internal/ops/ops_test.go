package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/programmersd21/zap/internal/errs"
)

func TestCopySingleFile(t *testing.T) {
	tmpDir := t.TempDir()

	// create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// copy file
	dstFile := filepath.Join(tmpDir, "dest.txt")
	opts := CopyOptions{
		Force:  false,
		Errors: errs.NewCollector(),
	}

	if err := Copy(srcFile, dstFile, opts); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// verify destination exists
	if _, err := os.Stat(dstFile); err != nil {
		t.Errorf("destination file not created: %v", err)
	}

	// verify content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("content mismatch: expected %q, got %q", content, dstContent)
	}
}

func TestCopyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// create source directory structure
	srcDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}

	files := map[string]string{
		filepath.Join(srcDir, "file1.txt"):           "content 1",
		filepath.Join(srcDir, "subdir", "file2.txt"): "content 2",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// copy directory
	dstDir := filepath.Join(tmpDir, "dest")
	opts := CopyOptions{
		Force:  false,
		Errors: errs.NewCollector(),
	}

	if err := Copy(srcDir, dstDir, opts); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// verify all files exist
	for srcPath, expectedContent := range files {
		relPath, _ := filepath.Rel(srcDir, srcPath)
		dstPath := filepath.Join(dstDir, relPath)

		content, err := os.ReadFile(dstPath)
		if err != nil {
			t.Errorf("failed to read %s: %v", dstPath, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("content mismatch for %s: expected %q, got %q", dstPath, expectedContent, content)
		}
	}
}

func TestMoveFile(t *testing.T) {
	tmpDir := t.TempDir()

	// create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// move file
	dstFile := filepath.Join(tmpDir, "dest.txt")
	opts := MoveOptions{
		Force:  false,
		Errors: errs.NewCollector(),
	}

	if err := Move(srcFile, dstFile, opts); err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	// verify source is gone
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("source file should not exist after move")
	}

	// verify destination exists and has correct content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("content mismatch: expected %q, got %q", content, dstContent)
	}
}

func TestDeleteFile(t *testing.T) {
	tmpDir := t.TempDir()

	// create file
	file := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// delete file
	opts := DeleteOptions{
		Recursive: false,
		Errors:    errs.NewCollector(),
	}

	if err := Delete(file, opts); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// verify file is gone
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Error("file should not exist after delete")
	}
}

func TestDeleteDirectoryRequiresRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// create directory
	dir := filepath.Join(tmpDir, "testdir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// try to delete without recursive flag
	opts := DeleteOptions{
		Recursive: false,
		Errors:    errs.NewCollector(),
	}

	err := Delete(dir, opts)
	if err == nil {
		t.Error("expected error when deleting directory without -r flag")
	}
}

func TestDeleteDirectoryRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// create directory structure
	dir := filepath.Join(tmpDir, "testdir")
	if err := os.MkdirAll(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// delete directory recursively
	opts := DeleteOptions{
		Recursive: true,
		Errors:    errs.NewCollector(),
	}

	if err := Delete(dir, opts); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// verify directory is gone
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("directory should not exist after delete")
	}
}

func TestCopySymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// create target file
	targetFile := filepath.Join(tmpDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	// create symlink
	srcLink := filepath.Join(tmpDir, "link.txt")
	if err := os.Symlink(targetFile, srcLink); err != nil {
		t.Skipf("skipping symlink test: %v", err)
	}

	// copy symlink
	dstLink := filepath.Join(tmpDir, "link_copy.txt")
	opts := CopyOptions{
		Force:  false,
		Errors: errs.NewCollector(),
	}

	if err := Copy(srcLink, dstLink, opts); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// verify destination is also a symlink
	linkInfo, err := os.Lstat(dstLink)
	if err != nil {
		t.Fatalf("failed to stat destination: %v", err)
	}

	if linkInfo.Mode()&os.ModeSymlink == 0 {
		t.Error("destination should be a symlink")
	}
}

func TestCopyPreservesPermissions(t *testing.T) {
	tmpDir := t.TempDir()

	// create source file with specific permissions
	srcFile := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("content"), 0755); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// copy file
	dstFile := filepath.Join(tmpDir, "dest.txt")
	opts := CopyOptions{
		Force:  false,
		Errors: errs.NewCollector(),
	}

	if err := Copy(srcFile, dstFile, opts); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// verify permissions are preserved
	srcInfo, _ := os.Stat(srcFile)
	dstInfo, _ := os.Stat(dstFile)

	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("permissions not preserved: expected %v, got %v", srcInfo.Mode(), dstInfo.Mode())
	}
}

// Benchmarks

func BenchmarkCopySmallFile(b *testing.B) {
	tmpDir := b.TempDir()
	src := filepath.Join(tmpDir, "small.dat")

	// 1 KB file
	data := make([]byte, 1024)
	if err := os.WriteFile(src, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := filepath.Join(tmpDir, fmt.Sprintf("copy%d.dat", i))
		opts := CopyOptions{Force: true}
		if err := Copy(src, dst, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCopyMediumFile(b *testing.B) {
	tmpDir := b.TempDir()
	src := filepath.Join(tmpDir, "medium.dat")

	// 1 MB file
	data := make([]byte, 1024*1024)
	if err := os.WriteFile(src, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := filepath.Join(tmpDir, fmt.Sprintf("copy%d.dat", i))
		opts := CopyOptions{Force: true}
		if err := Copy(src, dst, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCopyLargeFile(b *testing.B) {
	tmpDir := b.TempDir()
	src := filepath.Join(tmpDir, "large.dat")

	// 10 MB file
	data := make([]byte, 10*1024*1024)
	if err := os.WriteFile(src, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := filepath.Join(tmpDir, fmt.Sprintf("copy%d.dat", i))
		opts := CopyOptions{Force: true}
		if err := Copy(src, dst, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCopyDirectory(b *testing.B) {
	tmpDir := b.TempDir()
	src := filepath.Join(tmpDir, "srcdir")
	if err := os.Mkdir(src, 0755); err != nil {
		b.Fatal(err)
	}

	// create 10 files of 100 KB each
	for i := 0; i < 10; i++ {
		file := filepath.Join(src, fmt.Sprintf("file%d.dat", i))
		data := make([]byte, 100*1024)
		if err := os.WriteFile(file, data, 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := filepath.Join(tmpDir, fmt.Sprintf("dstdir%d", i))
		opts := CopyOptions{Force: true}
		if err := Copy(src, dst, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMove(b *testing.B) {
	tmpDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		src := filepath.Join(tmpDir, fmt.Sprintf("movesrc%d.dat", i))
		dst := filepath.Join(tmpDir, fmt.Sprintf("movedst%d.dat", i))

		data := make([]byte, 1024*1024) // 1 MB
		if err := os.WriteFile(src, data, 0644); err != nil {
			b.Fatal(err)
		}

		opts := MoveOptions{Force: true}
		if err := Move(src, dst, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	tmpDir := b.TempDir()

	// pre-create files
	files := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		file := filepath.Join(tmpDir, fmt.Sprintf("delete%d.dat", i))
		data := make([]byte, 1024) // 1 KB
		if err := os.WriteFile(file, data, 0644); err != nil {
			b.Fatal(err)
		}
		files[i] = file
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := DeleteOptions{}
		if err := Delete(files[i], opts); err != nil {
			b.Fatal(err)
		}
	}
}
