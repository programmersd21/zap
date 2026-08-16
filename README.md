# zap

fast file operations with real-time progress.

![demo](assets/demo.gif)

`zap` is a modern replacement for `cp`, `mv`, and `rm` with real-time transfer stats and progress display.

## features

- **fast** — 1 mb buffered streaming io for large file copies
- **real-time progress** — transfer speeds, file counts, and estimated time remaining
- **terminal aware** — dynamic width adjustments and automatic path truncation
- **safe defaults** — same-path identity protection, overwrite confirmations (`[Y/n]`), recursive delete checks, and root protection (`--no-preserve-root` to override)
- **script friendly** — automatically falls back to direct output when stdout is not a tty

## install

```bash
go install github.com/programmersd21/zap/cmd/zap@latest
```
```
yay -S zap-cli-bin
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

### trash (`-t`)
```bash
zap -t file.log
zap -t cache/
```

### shred (`-s`)
```bash
zap -s secret.db
```

## flags

| flag | description |
| :--- | :--- |
| `-c, --copy` | copy mode (default) |
| `-m, --move` | move mode (rename or copy+delete across devices) |
| `-d, --delete` | delete mode |
| `-t, --trash` | trash mode (system trash) |
| `-s, --shred` | shred mode (overwrite with random data, then delete) |
| `-f, --force` | overwrite without confirmation |
| `-r, --recursive` | required for deleting directories |
| `-v, --verbose` | verbosity: `-v` files, `-vv` operations, `-vvv` debug |
| `--no-preserve-root` | allow operating on the filesystem root |
| `--version` | show version |

## license

[mit](LICENSE)
