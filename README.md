# zap

fast file operations with real-time progress.

![demo](assets/demo.gif)

`zap` is a modern replacement for `cp`, `mv`, and `rm` with real-time transfer stats and progress display.

## features

- **fast** — 1 mb buffered streaming io for large file copies
- **real-time progress** — transfer speeds, file counts, and estimated time remaining
- **terminal aware** — dynamic width adjustments and automatic path truncation
- **safe defaults** — same-path identity protection, overwrite confirmations (`[Y/n]`), and recursive delete checks
- **script friendly** — automatically falls back to direct output when stdout is not a tty

## install

```bash
go install github.com/programmersd21/zap/cmd/zap@latest
```

or build from source:

```bash
git clone https://github.com/programmersd21/zap.git
cd zap
make build
```

## usage

### copy (default)
```bash
zap file.txt backup/
zap photos/ backup/photos/
```

### move (`-m`)
```bash
zap -m old.txt new.txt
zap -m logs/ /mnt/backup/logs/
```

### delete (`-d`)
```bash
zap -d file.log
zap -d -r cache/
```

## flags

| flag | description |
| :--- | :--- |
| `-c, --copy` | copy mode (default) |
| `-m, --move` | move mode (rename or copy+delete across devices) |
| `-d, --delete` | delete mode |
| `-f, --force` | overwrite without confirmation |
| `-r, --recursive` | required for deleting directories |
| `-v, --verbose` | print each file as it completes |
| `--version` | show version |

## license

[mit](LICENSE)
