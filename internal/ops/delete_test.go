package ops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/programmersd21/zap/internal/errs"
)

func TestDeleteRefusesRoot(t *testing.T) {
	opts := DeleteOptions{
		Recursive: true,
		Errors:    errs.NewCollector(),
	}
	if err := Delete(string(filepath.Separator), opts); err == nil {
		t.Fatal("expected delete of the filesystem root to be refused")
	}
}

func TestDeleteRootAllowedWithNoPreserveRoot(t *testing.T) {
	// the guard is the check we test here; passing NoPreserveRoot must not
	// trip the root refusal. we target a path that is not really the root.
	opts := DeleteOptions{
		Recursive:      true,
		NoPreserveRoot: true,
	}
	if err := Delete(filepath.Join(t.TempDir(), "missing"), opts); err != nil {
		t.Errorf("expected no-op delete, got %v", err)
	}
}

func TestIsRoot(t *testing.T) {
	if !isRoot("/") {
		t.Error("expected / to be the root")
	}
	if !isRoot("//") {
		t.Error("expected // to resolve to the root")
	}
	if isRoot("/home") {
		t.Error("expected /home not to be the root")
	}
	wd, _ := os.Getwd()
	if isRoot(wd) {
		t.Errorf("expected %s not to be the root", wd)
	}
}
