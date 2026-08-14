package walk

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Stats struct {
	TotalBytes int64
	TotalFiles int64
	TotalDirs  int64
}

func ComputeStats(paths []string) (Stats, error) {
	var stats Stats

	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			return stats, err
		}

		if info.IsDir() {
			err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				info, err := d.Info()
				if err != nil {
					return nil
				}

				if info.IsDir() {
					stats.TotalDirs++
				} else {
					stats.TotalFiles++
					if info.Mode().IsRegular() {
						stats.TotalBytes += info.Size()
					}
				}
				return nil
			})
			if err != nil {
				return stats, err
			}
		} else {
			stats.TotalFiles++
			if info.Mode().IsRegular() {
				stats.TotalBytes += info.Size()
			}
		}
	}

	return stats, nil
}
