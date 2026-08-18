package ops

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/ui"
)

type DeleteOptions struct {
	Recursive      bool
	Verbose        int
	Program        *tea.Program
	Errors         *errs.Collector
	CumulFiles     *int64
	NoPreserveRoot bool
}

func Delete(path string, opts DeleteOptions) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !opts.NoPreserveRoot && isRoot(path) {
		opLog(opts.Verbose, "refuse %s: filesystem root", path)
		return fmt.Errorf("refusing to delete %q: this is the filesystem root (use --no-preserve-root to override)", path)
	}

	if info.IsDir() {
		if !opts.Recursive {
			opLog(opts.Verbose, "skip %s: directory requires -r", path)
			return fmt.Errorf("cannot delete directory %q: -r required", path)
		}
		return deleteDir(path, opts)
	}

	return deleteFile(path, opts)
}

func deleteFile(path string, opts DeleteOptions) error {
	if err := os.Remove(path); err != nil {
		opLog(opts.Verbose, "error deleting %s: %v", path, err)
		return err
	}
	opLog(opts.Verbose, "deleted %s", path)
	sendDeleteProgress(opts, path, true)
	return nil
}

func deleteDir(path string, opts DeleteOptions) error {
	var files, dirs []string

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			opLog(opts.Verbose, "skip %s: %v", p, err)
			if opts.Errors != nil {
				opts.Errors.Add(p, err)
			}
			return nil
		}
		if d.IsDir() {
			if p != path {
				dirs = append(dirs, p)
			}
		} else {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			opLog(opts.Verbose, "error deleting %s: %v", f, err)
			if opts.Errors != nil {
				opts.Errors.Add(f, err)
			}
		} else {
			opLog(opts.Verbose, "deleted %s", f)
			sendDeleteProgress(opts, f, true)
		}
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		d := dirs[i]
		if err := os.Remove(d); err != nil {
			opLog(opts.Verbose, "error deleting %s: %v", d, err)
			if opts.Errors != nil {
				opts.Errors.Add(d, err)
			}
		} else {
			opLog(opts.Verbose, "deleted %s", d)
			sendDeleteProgress(opts, d, false)
		}
	}

	if err := os.Remove(path); err != nil {
		return err
	}
	opLog(opts.Verbose, "deleted %s", path)
	sendDeleteProgress(opts, path, false)

	return nil
}

func sendDeleteProgress(opts DeleteOptions, path string, isFile bool) {
	var files int64
	if opts.CumulFiles != nil {
		if isFile {
			*opts.CumulFiles++
		}
		files = *opts.CumulFiles
	}
	if opts.Program == nil {
		return
	}
	opts.Program.Send(ui.ProgressMsg{
		FilesDone:   files,
		CurrentFile: path,
	})
}

func isRoot(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return filepath.Clean(abs) == string(filepath.Separator)
}
