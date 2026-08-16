package ops

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/ui"
)

const defaultShredPasses = 3

type ShredOptions struct {
	Passes         int
	Program        *tea.Program
	Errors         *errs.Collector
	CumulBytes     *int64
	CumulFiles     *int64
	NoPreserveRoot bool
}

// Shred overwrites a regular file with random data several times, then
// removes it. Directories, symlinks, and special files are refused.
func Shred(path string, opts ShredOptions) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("refusing to shred %q: not a regular file", path)
	}
	if !opts.NoPreserveRoot && isRoot(path) {
		return fmt.Errorf("refusing to shred %q: this is the filesystem root (use --no-preserve-root to override)", path)
	}

	passes := opts.Passes
	if passes <= 0 {
		passes = defaultShredPasses
	}

	baseline := int64(0)
	if opts.CumulBytes != nil {
		baseline = *opts.CumulBytes
	}

	size, err := overwriteFile(path, passes, baseline, opts)
	if err != nil {
		return err
	}

	if opts.Program != nil && opts.CumulBytes != nil {
		*opts.CumulBytes = baseline + size
		if opts.CumulFiles != nil {
			*opts.CumulFiles++
		}
		opts.Program.Send(ui.ProgressMsg{
			BytesDone:   *opts.CumulBytes,
			FilesDone:   *opts.CumulFiles,
			CurrentFile: path,
		})
	}

	delOpts := DeleteOptions{
		Errors:         opts.Errors,
		NoPreserveRoot: opts.NoPreserveRoot,
	}
	if err := Delete(path, delOpts); err != nil {
		return fmt.Errorf("shred %s: %w", path, err)
	}
	return nil
}

func overwriteFile(path string, passes int, baseline int64, opts ShredOptions) (int64, error) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	size := info.Size()
	if size == 0 {
		return 0, nil
	}

	buf := make([]byte, 1024*1024)
	for pass := 0; pass < passes; pass++ {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		var written int64
		for written < size {
			n := int64(len(buf))
			if remaining := size - written; remaining < n {
				n = remaining
			}
			if _, err := io.ReadFull(rand.Reader, buf[:n]); err != nil {
				return 0, err
			}
			if _, err := f.Write(buf[:n]); err != nil {
				return 0, err
			}
			written += n
		}
		if err := f.Sync(); err != nil {
			return 0, err
		}

		if opts.Program != nil && opts.CumulBytes != nil {
			*opts.CumulBytes = baseline + int64(math.Round(float64(size)*float64(pass+1)/float64(passes)))
			opts.Program.Send(ui.ProgressMsg{
				BytesDone:   *opts.CumulBytes,
				CurrentFile: path,
			})
		}
	}

	return size, nil
}
