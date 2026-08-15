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

	if !strings.Contains(string(output), "0.1.1") {
		t.Errorf("expected version 0.1.1, got: %s", output)
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
