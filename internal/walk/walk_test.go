package walk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeStats(t *testing.T) {
	// create temporary directory structure
	tmpDir := t.TempDir()

	// create test files
	files := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "subdir", "file3.txt"),
	}

	// create files with known sizes
	for i, file := range files {
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		content := make([]byte, (i+1)*100) // 100, 200, 300 bytes
		if err := os.WriteFile(file, content, 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	tests := []struct {
		name          string
		paths         []string
		expectedFiles int64
		expectedDirs  int64
		minBytes      int64
	}{
		{
			name:          "single file",
			paths:         []string{files[0]},
			expectedFiles: 1,
			expectedDirs:  0,
			minBytes:      100,
		},
		{
			name:          "directory",
			paths:         []string{tmpDir},
			expectedFiles: 3,
			expectedDirs:  2,   // subdir + tmpDir itself
			minBytes:      600, // 100 + 200 + 300
		},
		{
			name:          "multiple files",
			paths:         []string{files[0], files[1]},
			expectedFiles: 2,
			expectedDirs:  0,
			minBytes:      300, // 100 + 200
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := ComputeStats(tt.paths)
			if err != nil {
				t.Fatalf("ComputeStats failed: %v", err)
			}

			if stats.TotalFiles != tt.expectedFiles {
				t.Errorf("expected %d files, got %d", tt.expectedFiles, stats.TotalFiles)
			}

			if stats.TotalDirs != tt.expectedDirs {
				t.Errorf("expected %d directories, got %d", tt.expectedDirs, stats.TotalDirs)
			}

			if stats.TotalBytes < tt.minBytes {
				t.Errorf("expected at least %d bytes, got %d", tt.minBytes, stats.TotalBytes)
			}
		})
	}
}

func TestComputeStatsNonExistent(t *testing.T) {
	_, err := ComputeStats([]string{"/nonexistent/path"})
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}

func TestComputeStatsSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// create a file
	targetFile := filepath.Join(tmpDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// create a symlink
	symlinkFile := filepath.Join(tmpDir, "link.txt")
	if err := os.Symlink(targetFile, symlinkFile); err != nil {
		t.Skipf("skipping symlink test: %v", err)
	}

	stats, err := ComputeStats([]string{symlinkFile})
	if err != nil {
		t.Fatalf("ComputeStats failed: %v", err)
	}

	if stats.TotalFiles != 1 {
		t.Errorf("expected 1 file (symlink), got %d", stats.TotalFiles)
	}
}
