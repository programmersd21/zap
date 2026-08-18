# changelog

all notable changes to this project will be documented in this file.

## [0.1.3] - 2026-08-16

### changes
- verbosity levels are now clearly separated: `-vv` reports operation details (sources, destinations, overwrites, skip reasons, permission and filesystem errors), `-vvv` adds internal debug output (resolved paths, tty decisions, computed stats, shred passes)
- operation diagnostics from the engine are logged to stderr with an `op:`/`debug:` prefix, so progress ui output on stdout stays clean

## [0.1.2] - 2026-08-16

### features
- trash mode (`-t`, `--trash`) moves files to the system trash — freedesktop spec on linux, `~/.Trash` on macos, with name collision handling and cross-device support
- shred mode (`-s`, `--shred`) overwrites regular files with random data before deletion, with live byte progress
- delete now refuses to operate on the filesystem root unless `--no-preserve-root` is given
- verbosity levels: `-v` files, `-vv` operation logging, `-vvv` debug output
- conflicting modes are rejected with a hint instead of behaving unexpectedly

## [0.1.1] - 2026-08-15

### fixes
- fixed delete operations reporting 0 files in completion summary
- fixed same-path operations (copy/move to self) potentially destroying files
- overwrite prompt now defaults to yes on enter (`[Y/n]`)
- improved error and cancellation handling in progress ui

## [0.1.0] - 2026-08-14

### features
- fast file operations with progress tracking for copy, move, and delete
- real-time transfer speed, eta, and file counts
- responsive layout with automatic path truncation and terminal width detection
- catppuccin mocha color palette with smooth animations
- safe overwrite prompts and directory delete validation
- non-tty fallback for script pipelines

[0.1.3]: https://github.com/programmersd21/zap/releases/tag/v0.1.3
[0.1.2]: https://github.com/programmersd21/zap/releases/tag/v0.1.2
[0.1.1]: https://github.com/programmersd21/zap/releases/tag/v0.1.1
[0.1.0]: https://github.com/programmersd21/zap/releases/tag/v0.1.0
