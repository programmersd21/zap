package errs

import (
	"errors"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestCollector(t *testing.T) {
	collector := NewCollector()

	if collector.HasErrors() {
		t.Error("new collector should not have errors")
	}

	if collector.Count() != 0 {
		t.Errorf("expected count 0, got %d", collector.Count())
	}
}

func TestCollectorAdd(t *testing.T) {
	collector := NewCollector()

	collector.Add("/path/to/file1.txt", errors.New("permission denied"))
	collector.Add("/path/to/file2.txt", errors.New("file not found"))

	if !collector.HasErrors() {
		t.Error("collector should have errors")
	}

	if collector.Count() != 2 {
		t.Errorf("expected count 2, got %d", collector.Count())
	}

	errs := collector.Errors()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}

	// verify error content
	if errs[0].Path != "/path/to/file1.txt" {
		t.Errorf("expected path %q, got %q", "/path/to/file1.txt", errs[0].Path)
	}

	if errs[0].Err.Error() != "permission denied" {
		t.Errorf("expected error %q, got %q", "permission denied", errs[0].Err.Error())
	}
}

func TestFileError(t *testing.T) {
	err := FileError{
		Path: "/path/to/file.txt",
		Err:  errors.New("test error"),
	}

	expected := "/path/to/file.txt: test error"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestFormatSummary(t *testing.T) {
	collector := NewCollector()
	style := lipgloss.NewStyle()

	// test with no errors
	summary := collector.FormatSummary(style)
	if summary != "" {
		t.Error("expected empty summary for no errors")
	}

	// test with one error
	collector.Add("/path/to/file1.txt", errors.New("error 1"))
	summary = collector.FormatSummary(style)

	if !strings.Contains(summary, "completed with 1 error") {
		t.Error("summary should contain singular 'error'")
	}

	if !strings.Contains(summary, "/path/to/file1.txt") {
		t.Error("summary should contain file path")
	}

	// test with multiple errors
	collector.Add("/path/to/file2.txt", errors.New("error 2"))
	summary = collector.FormatSummary(style)

	if !strings.Contains(summary, "completed with 2 errors") {
		t.Error("summary should contain plural 'errors'")
	}
}

func TestFormatSummaryTruncation(t *testing.T) {
	collector := NewCollector()
	style := lipgloss.NewStyle()

	// add more than 10 errors
	for i := 1; i <= 15; i++ {
		collector.Add("/path/to/file.txt", errors.New("error"))
	}

	summary := collector.FormatSummary(style)

	// should truncate at 10 and show "... and X more"
	if !strings.Contains(summary, "and 5 more error") {
		t.Error("summary should indicate truncation")
	}
}

func TestCollectorMultipleAdds(t *testing.T) {
	collector := NewCollector()

	for i := 0; i < 100; i++ {
		collector.Add("/path/to/file.txt", errors.New("error"))
	}

	if collector.Count() != 100 {
		t.Errorf("expected count 100, got %d", collector.Count())
	}

	if len(collector.Errors()) != 100 {
		t.Errorf("expected 100 errors, got %d", len(collector.Errors()))
	}
}
