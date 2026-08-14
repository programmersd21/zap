package ops

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/ui"
)

type CopyOptions struct {
	Force      bool
	Program    *tea.Program
	Errors     *errs.Collector
	CumulBytes *int64
	CumulFiles *int64
}

func Copy(src, dst string, opts CopyOptions) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("stat %s: %w", src, err)
	}

	if dstInfo, err := os.Lstat(dst); err == nil {
		if !opts.Force {
			return fmt.Errorf("destination %s exists", dst)
		}
		if dstInfo.IsDir() {
			if err := os.RemoveAll(dst); err != nil {
				return err
			}
		} else if err := os.Remove(dst); err != nil {
			return err
		}
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst, opts)
	}
	return copyFile(src, dst, srcInfo, opts)
}

func copyFile(src, dst string, srcInfo fs.FileInfo, opts CopyOptions) error {
	if srcInfo.Mode()&os.ModeSymlink != 0 {
		return copySymlink(src, dst, opts)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	buf := make([]byte, 1024*1024)
	for {
		n, err := srcFile.Read(buf)
		if n > 0 {
			if _, writeErr := dstFile.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			if opts.Program != nil && opts.CumulBytes != nil {
				*opts.CumulBytes += int64(n)
				var files int64
				if opts.CumulFiles != nil {
					files = *opts.CumulFiles
				}
				opts.Program.Send(ui.ProgressMsg{
					BytesDone:   *opts.CumulBytes,
					FilesDone:   files,
					CurrentFile: dst,
				})
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// increment file count and send final progress for this file
	if opts.CumulFiles != nil {
		*opts.CumulFiles++
		if opts.Program != nil && opts.CumulBytes != nil {
			opts.Program.Send(ui.ProgressMsg{
				BytesDone:   *opts.CumulBytes,
				FilesDone:   *opts.CumulFiles,
				CurrentFile: dst,
			})
		}
	}

	if err := dstFile.Chmod(srcInfo.Mode()); err != nil && opts.Errors != nil {
		opts.Errors.Add(dst, fmt.Errorf("chmod: %w", err))
	}
	if err := os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime()); err != nil && opts.Errors != nil {
		opts.Errors.Add(dst, fmt.Errorf("chtimes: %w", err))
	}

	return nil
}

func copySymlink(src, dst string, opts CopyOptions) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}
	if err := os.Symlink(target, dst); err != nil {
		return err
	}

	if opts.Program != nil && opts.CumulFiles != nil {
		*opts.CumulFiles++
		var bytes int64
		if opts.CumulBytes != nil {
			bytes = *opts.CumulBytes
		}
		opts.Program.Send(ui.ProgressMsg{
			BytesDone:   bytes,
			FilesDone:   *opts.CumulFiles,
			CurrentFile: dst,
		})
	}

	return nil
}

func copyDir(src, dst string, opts CopyOptions) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Errors != nil {
				opts.Errors.Add(path, err)
			}
			return nil
		}
		if path == src {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			if opts.Errors != nil {
				opts.Errors.Add(path, err)
			}
			return nil
		}
		dstPath := filepath.Join(dst, relPath)

		info, err := d.Info()
		if err != nil {
			if opts.Errors != nil {
				opts.Errors.Add(path, err)
			}
			return nil
		}

		if d.IsDir() {
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil && opts.Errors != nil {
				opts.Errors.Add(path, err)
			}
		} else if err := copyFile(path, dstPath, info, opts); err != nil && opts.Errors != nil {
			opts.Errors.Add(path, err)
		}

		return nil
	})
}
