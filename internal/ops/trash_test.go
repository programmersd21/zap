package ops

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestTrashMovesToTrash(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("trash location differs on darwin")
	}

	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	src := filepath.Join(tmp, "src.txt")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Trash(src, TrashOptions{}); err != nil {
		t.Fatalf("Trash failed: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("source should be gone after trash")
	}

	trashFiles := filepath.Join(tmp, "Trash", "files")
	if _, err := os.Stat(filepath.Join(trashFiles, "src.txt")); err != nil {
		t.Errorf("file not found in trash files dir: %v", err)
	}

	infoPath := filepath.Join(tmp, "Trash", "info", "src.txt.trashinfo")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		t.Fatalf("trashinfo missing: %v", err)
	}
	if !strings.HasPrefix(string(data), "[Trash Info]\nPath=") {
		t.Errorf("unexpected trashinfo content: %q", data)
	}
	if !strings.Contains(string(data), "DeletionDate=") {
		t.Errorf("trashinfo missing DeletionDate: %q", data)
	}
}

func TestTrashNameCollision(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("trash location differs on darwin")
	}

	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	for _, name := range []string{"a.txt", "b.txt"} {
		src := filepath.Join(tmp, name)
		if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := Trash(filepath.Join(tmp, "a.txt"), TrashOptions{}); err != nil {
		t.Fatalf("first Trash failed: %v", err)
	}
	if err := Trash(filepath.Join(tmp, "b.txt"), TrashOptions{}); err != nil {
		t.Fatalf("second Trash failed: %v", err)
	}

	// both files must exist in trash with unique names
	trashFiles := filepath.Join(tmp, "Trash", "files")
	entries, err := os.ReadDir(trashFiles)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 trashed files, got %d", len(entries))
	}
}

func TestTrashDirectory(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("trash location differs on darwin")
	}

	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	dir := filepath.Join(tmp, "mydir")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "f.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Trash(dir, TrashOptions{}); err != nil {
		t.Fatalf("Trash failed: %v", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("source dir should be gone after trash")
	}
	if _, err := os.Stat(filepath.Join(tmp, "Trash", "files", "mydir", "sub", "f.txt")); err != nil {
		t.Errorf("directory tree not preserved in trash: %v", err)
	}
}

func TestTrashNonexistent(t *testing.T) {
	if err := Trash(filepath.Join(t.TempDir(), "nope"), TrashOptions{}); err != nil {
		t.Errorf("trashing a missing path should be a no-op, got %v", err)
	}
}

func TestTrashReportsErrors(t *testing.T) {
	// force an unusable trash root so the operation fails cleanly
	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "Trash")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_DATA_HOME", tmp)

	if err := Trash(blocker, TrashOptions{}); err == nil {
		t.Errorf("expected Trash to fail when the trash root is unusable")
	}
}
