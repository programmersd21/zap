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
	Recursive  bool
	Program    *tea.Program
	Errors     *errs.Collector
	CumulFiles *int64
}

func Delete(path string, opts DeleteOptions) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		if !opts.Recursive {
			return fmt.Errorf("cannot delete directory %q: -r required", path)
		}
		return deleteDir(path, opts)
	}

	return deleteFile(path, opts)
}

func deleteFile(path string, opts DeleteOptions) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	sendDeleteProgress(opts, path)
	return nil
}

func deleteDir(path string, opts DeleteOptions) error {
	var files, dirs []string

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
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
			if opts.Errors != nil {
				opts.Errors.Add(f, err)
			}
		} else {
			sendDeleteProgress(opts, f)
		}
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		d := dirs[i]
		if err := os.Remove(d); err != nil {
			if opts.Errors != nil {
				opts.Errors.Add(d, err)
			}
		} else {
			sendDeleteProgress(opts, d)
		}
	}

	if err := os.Remove(path); err != nil {
		return err
	}
	sendDeleteProgress(opts, path)

	return nil
}

func sendDeleteProgress(opts DeleteOptions, path string) {
	if opts.Program == nil {
		return
	}
	var files int64
	if opts.CumulFiles != nil {
		*opts.CumulFiles++
		files = *opts.CumulFiles
	}
	opts.Program.Send(ui.ProgressMsg{
		FilesDone:   files,
		CurrentFile: path,
	})
}
