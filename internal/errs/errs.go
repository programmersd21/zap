package errs

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type FileError struct {
	Path string
	Err  error
}

func (e FileError) Error() string {
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

type Collector struct {
	errors []FileError
}

func NewCollector() *Collector {
	return &Collector{errors: make([]FileError, 0)}
}

func (c *Collector) Add(path string, err error) {
	c.errors = append(c.errors, FileError{Path: path, Err: err})
}

func (c *Collector) HasErrors() bool {
	return len(c.errors) > 0
}

func (c *Collector) Count() int {
	return len(c.errors)
}

func (c *Collector) Errors() []FileError {
	return c.errors
}

func (c *Collector) FormatSummary(style lipgloss.Style) string {
	if len(c.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	noun := "error"
	if len(c.errors) != 1 {
		noun = "errors"
	}
	sb.WriteString(style.Bold(true).Render(fmt.Sprintf("✗ completed with %d %s", len(c.errors), noun)))
	sb.WriteString("\n\n")

	const maxDisplay = 10
	for i, e := range c.errors {
		if i >= maxDisplay {
			rem := len(c.errors) - maxDisplay
			remNoun := "error"
			if rem != 1 {
				remNoun = "errors"
			}
			sb.WriteString(style.Render(fmt.Sprintf("  … and %d more %s\n", rem, remNoun)))
			break
		}
		sb.WriteString(style.Render(fmt.Sprintf("  • %s: %v\n", e.Path, e.Err)))
	}

	return sb.String()
}
