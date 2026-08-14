package ui

import (
	"testing"
	"time"
)

func TestViewProgress(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 1000, 10)
	m.bytesDone = 500
	m.filesDone = 5
	m.currentFile = "/tmp/test.txt"
	m.bytesPerSec = 1024 * 1024 // 1 MB/s

	// Should not panic
	_ = m.View()
}

func TestViewCompletion(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 1000, 10)
	m.done = true
	m.startTime = time.Now().Add(-5 * time.Second)

	// Should not panic
	_ = m.View()
}

func TestViewVerboseMode(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, true, 1000, 10)
	m.completedFiles = []string{"/tmp/file1.txt", "/tmp/file2.txt"}

	// Should not panic
	_ = m.View()
}

func TestViewDeleteMode(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpDelete, false, 0, 10) // zero bytes = delete mode
	m.filesDone = 5

	// Should not panic
	_ = m.View()
}

func TestViewMoveMode(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpMove, false, 1000, 10)
	m.bytesDone = 500
	m.filesDone = 5

	// Should not panic
	_ = m.View()
}

func TestViewSpeedCalculation(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 10000000, 10)
	m.bytesDone = 5000000
	m.bytesPerSec = 500000

	// Should not panic
	_ = m.View()
}

func TestViewETACalculation(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 10000000, 10)
	m.bytesDone = 1000000
	m.totalBytes = 10000000
	m.bytesPerSec = 1000000 // 1 MB/s, should take 9 more seconds

	// Should not panic
	_ = m.View()
}

func TestViewWithCurrentFile(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 1000, 10)
	m.currentFile = "/very/long/path/to/some/file/that/should/be/truncated.txt"

	// Should not panic
	_ = m.View()
}

func TestViewCompletedFiles(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, true, 1000, 5)
	m.completedFiles = []string{
		"/tmp/file1.txt",
		"/tmp/file2.txt",
		"/tmp/file3.txt",
	}

	// Should not panic
	_ = m.View()
}

func TestViewWithZeroTotals(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 0, 0)

	// Should handle zero totals gracefully
	_ = m.View()
}

func TestViewWithLargeNumbers(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 1024*1024*1024*100, 10000) // 100 GB, 10k files
	m.bytesDone = 1024 * 1024 * 1024 * 50                          // 50 GB
	m.filesDone = 5000

	// Should handle large numbers
	_ = m.View()
}

func TestFormatDurationSubSecond(t *testing.T) {
	result := formatDuration(0)
	if result != "<1s" {
		t.Errorf("expected '<1s' for 0 seconds, got %s", result)
	}
}

func TestSpeedCalculation(t *testing.T) {
	theme := ThemeMocha
	m := NewModel(theme, OpCopy, false, 1000000, 10)
	m.lastUpdate = time.Now().Add(-1 * time.Second)

	// Simulate progress update after 1 second with 500KB transferred
	msg := ProgressMsg{
		BytesDone: 500000,
		FilesDone: 5,
	}

	updated, _ := m.Update(msg)
	m = updated.(Model)

	// Speed should be approximately 500KB/s
	if m.bytesPerSec < 400000 || m.bytesPerSec > 600000 {
		t.Errorf("expected speed around 500000, got %f", m.bytesPerSec)
	}
}
