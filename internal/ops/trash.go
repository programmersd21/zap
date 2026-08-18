package ops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/ui"
)

type TrashOptions struct {
	Verbose    int
	Program    *tea.Program
	Errors     *errs.Collector
	CumulBytes *int64
	CumulFiles *int64
}

// Trash moves a file or directory to the system trash following the
// freedesktop.org trash specification on linux and ~/.Trash on macos.
func Trash(path string, opts TrashOptions) error {
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	filesDir, infoDir, err := trashDirs()
	if err != nil {
		return err
	}

	name, err := uniqueTrashName(filesDir, filepath.Base(path))
	if err != nil {
		return err
	}

	dst := filepath.Join(filesDir, name)
	count := countFiles(path)

	err = os.Rename(path, dst)
	if err == nil {
		opLog(opts.Verbose, "trash %s -> %s", path, dst)
		if writeErr := writeTrashInfo(infoDir, name, path); writeErr != nil {
			return writeErr
		}
		sendTrashProgress(opts, count, path)
		return nil
	}

	var linkErr *os.LinkError
	if !errors.As(err, &linkErr) {
		return fmt.Errorf("trash %s: %w", path, err)
	}

	// cross-device move: copy into the trash, then remove the source
	opLog(opts.Verbose, "trash %s: cross-device, copy+delete fallback", path)
	copyOpts := CopyOptions{
		Force:      true,
		Verbose:    opts.Verbose,
		Program:    opts.Program,
		Errors:     opts.Errors,
		CumulBytes: opts.CumulBytes,
		CumulFiles: opts.CumulFiles,
	}
	if err := Copy(path, dst, copyOpts); err != nil {
		return err
	}
	if writeErr := writeTrashInfo(infoDir, name, path); writeErr != nil {
		return writeErr
	}

	var deleted int64
	delOpts := DeleteOptions{
		Recursive:  true,
		Verbose:    opts.Verbose,
		Errors:     opts.Errors,
		CumulFiles: &deleted,
	}
	if err := Delete(path, delOpts); err != nil {
		return fmt.Errorf("cleanup source: %w", err)
	}
	if opts.Program != nil && opts.CumulFiles != nil {
		*opts.CumulFiles += countFiles(dst)
		opts.Program.Send(ui.ProgressMsg{
			FilesDone:   *opts.CumulFiles,
			CurrentFile: path,
		})
	}

	return nil
}

func trashDirs() (filesDir, infoDir string, err error) {
	var root string
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", err
		}
		root = filepath.Join(home, ".Trash")
	default:
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", "", err
			}
			dataHome = filepath.Join(home, ".local", "share")
		}
		root = filepath.Join(dataHome, "Trash")
	}

	filesDir = filepath.Join(root, "files")
	infoDir = filepath.Join(root, "info")
	if err := os.MkdirAll(filesDir, 0700); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(infoDir, 0700); err != nil {
		return "", "", err
	}
	return filesDir, infoDir, nil
}

func uniqueTrashName(dir, base string) (string, error) {
	if _, err := os.Lstat(filepath.Join(dir, base)); os.IsNotExist(err) {
		return base, nil
	}
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s.%d%s", stem, i, ext)
		if _, err := os.Lstat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate, nil
		}
	}
}

func writeTrashInfo(infoDir, name, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	info := fmt.Sprintf(
		"[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		abs, time.Now().Format("2006-01-02T15:04:05"),
	)
	if err := os.WriteFile(filepath.Join(infoDir, name+".trashinfo"), []byte(info), 0600); err != nil {
		return fmt.Errorf("trash info %s: %w", path, err)
	}
	return nil
}

func countFiles(path string) int64 {
	info, err := os.Lstat(path)
	if err != nil {
		return 1
	}
	if !info.IsDir() {
		return 1
	}

	var n int64
	_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			n++
		}
		return nil
	})
	return n
}

func sendTrashProgress(opts TrashOptions, n int64, path string) {
	if opts.Program == nil || opts.CumulFiles == nil {
		return
	}
	*opts.CumulFiles += n
	opts.Program.Send(ui.ProgressMsg{
		FilesDone:   *opts.CumulFiles,
		CurrentFile: path,
	})
}
