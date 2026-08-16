package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestVersion tests --version flag
func TestVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// build the binary
	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	// test version
	cmd = exec.Command(binary, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run --version: %v", err)
	}

	if !strings.Contains(string(output), version) {
		t.Errorf("expected version %s, got: %s", version, output)
	}
}

// TestHelp tests --help flag
func TestHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	cmd = exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run --help: %v", err)
	}

	out := string(output)
	if !strings.Contains(out, "usage") {
		t.Error("help should contain usage")
	}
	if !strings.Contains(out, "examples") {
		t.Error("help should contain examples")
	}
}

// TestCopyIntegration tests basic copy operation
func TestCopyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// build binary
	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	// create test environment
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")

	// write source file
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	// run copy
	cmd = exec.Command(binary, "-f", srcFile, dstFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("copy failed: %v\noutput: %s", err, output)
	}

	// verify destination exists
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("expected 'test content', got: %s", content)
	}
}

// TestMoveIntegration tests move operation
func TestMoveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "move_src.txt")
	dstFile := filepath.Join(tmpDir, "move_dst.txt")

	if err := os.WriteFile(srcFile, []byte("move test"), 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	// run move
	cmd = exec.Command(binary, "-m", "-f", srcFile, dstFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("move failed: %v\noutput: %s", err, output)
	}

	// verify source is gone
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("source should not exist after move")
	}

	// verify destination exists
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	if string(content) != "move test" {
		t.Errorf("expected 'move test', got: %s", content)
	}
}

// TestDeleteIntegration tests delete operation
func TestDeleteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "delete_test.txt")

	if err := os.WriteFile(testFile, []byte("delete me"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// run delete
	cmd = exec.Command(binary, "-d", "-f", testFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("delete failed: %v\noutput: %s", err, output)
	}

	// verify file is gone
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("file should not exist after delete")
	}
}

// TestTrashIntegration tests trash operation
func TestTrashIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	tmpDir := t.TempDir()
	trashHome := filepath.Join(tmpDir, "trash-home")
	t.Setenv("XDG_DATA_HOME", trashHome)

	testFile := filepath.Join(tmpDir, "trash_test.txt")
	if err := os.WriteFile(testFile, []byte("trash me"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd = exec.Command(binary, "-t", testFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("trash failed: %v\noutput: %s", err, output)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("source should not exist after trash")
	}
	if _, err := os.Stat(filepath.Join(trashHome, "Trash", "files", "trash_test.txt")); err != nil {
		t.Errorf("file should be in the trash: %v", err)
	}
	if _, err := os.Stat(filepath.Join(trashHome, "Trash", "info", "trash_test.txt.trashinfo")); err != nil {
		t.Errorf("trashinfo should exist: %v", err)
	}
}

// TestDeletePreservesRoot tests root protection
func TestDeletePreservesRoot(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	cmd = exec.Command(binary, "-d", "-r", "/")
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Errorf("expected root delete to be refused, output: %s", output)
	}
}

// TestConflictingModes tests mode exclusivity
func TestConflictingModes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	cmd = exec.Command(binary, "-d", "-m", "/tmp/x", "/tmp/y")
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Errorf("expected conflicting modes to fail, output: %s", output)
	}
}

// TestCopyDirectory tests directory copy
func TestCopyDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := filepath.Join(t.TempDir(), "zap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "srcdir")
	dstDir := filepath.Join(tmpDir, "dstdir")

	// create source directory with files
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatalf("failed to create srcdir: %v", err)
	}

	for i := 1; i <= 3; i++ {
		file := filepath.Join(srcDir, "file"+string(rune('0'+i))+".txt")
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	// run copy
	cmd = exec.Command(binary, "-f", srcDir, dstDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("copy dir failed: %v\noutput: %s", err, output)
	}

	// verify destination directory exists
	if _, err := os.Stat(dstDir); err != nil {
		t.Fatalf("destination directory should exist: %v", err)
	}

	// verify files copied
	entries, err := os.ReadDir(dstDir)
	if err != nil {
		t.Fatalf("failed to read dstdir: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 files, got %d", len(entries))
	}
}
