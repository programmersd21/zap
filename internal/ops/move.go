package ops

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/programmersd21/zap/internal/errs"
)

type MoveOptions struct {
	Force      bool
	Verbose    int
	Program    *tea.Program
	Errors     *errs.Collector
	CumulBytes *int64
	CumulFiles *int64
}

func Move(src, dst string, opts MoveOptions) error {
	if samePath(src, dst) {
		opLog(opts.Verbose, "skip %s: same file as %s", src, dst)
		return fmt.Errorf("%s and %s are the same file", src, dst)
	}
	err := os.Rename(src, dst)
	if err == nil {
		opLog(opts.Verbose, "move %s -> %s", src, dst)
		return nil
	}

	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		if errno, ok := linkErr.Err.(syscall.Errno); ok && errno == syscall.EXDEV {
			opLog(opts.Verbose, "move %s -> %s: cross-device, copy+delete fallback", src, dst)
			return moveCrossDevice(src, dst, opts)
		}
	}
	return fmt.Errorf("move: %w", err)
}

func moveCrossDevice(src, dst string, opts MoveOptions) error {
	copyOpts := CopyOptions(opts)
	if err := Copy(src, dst, copyOpts); err != nil {
		return err
	}

	deleteOpts := DeleteOptions{
		Recursive: true,
		Program:   opts.Program,
		Errors:    opts.Errors,
	}
	if err := Delete(src, deleteOpts); err != nil && opts.Errors != nil {
		opts.Errors.Add(src, fmt.Errorf("cleanup source: %w", err))
	}

	return nil
}
