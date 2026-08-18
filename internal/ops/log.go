package ops

import (
	"fmt"
	"os"
)

// opLog writes operation-level diagnostics (level >= 2, i.e. -vv) to stderr.
// These cover sources, destinations, permission issues, skip reasons, and
// filesystem errors — external causes a user can act on.
func opLog(verbose int, format string, args ...any) {
	if verbose >= 2 {
		fmt.Fprintf(os.Stderr, "op: "+format+"\n", args...)
	}
}

// debugLog writes internal diagnostics (level >= 3, i.e. -vvv) to stderr.
// These cover zap's own state and decisions for diagnosing internal bugs.
func debugLog(verbose int, format string, args ...any) {
	if verbose >= 3 {
		fmt.Fprintf(os.Stderr, "debug: "+format+"\n", args...)
	}
}
