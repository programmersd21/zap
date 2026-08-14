package ui

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"zero bytes", 0, "0 b"},
		{"bytes", 500, "500 b"},
		{"kilobytes", 1024, "1.0 kb"},
		{"megabytes", 1024 * 1024, "1.0 mb"},
		{"gigabytes", 1024 * 1024 * 1024, "1.0 gb"},
		{"terabytes", 1024 * 1024 * 1024 * 1024, "1.0 tb"},
		{"large value", 1536 * 1024 * 1024, "1.5 gb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %q, expected %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		number   int64
		expected string
	}{
		{"zero", 0, "0"},
		{"small", 99, "99"},
		{"hundreds", 999, "999"},
		{"thousands", 1000, "1,000"},
		{"ten thousands", 10000, "10,000"},
		{"millions", 1000000, "1,000,000"},
		{"large", 1234567890, "1,234,567,890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumber(tt.number)
			if result != tt.expected {
				t.Errorf("FormatNumber(%d) = %q, expected %q", tt.number, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{"zero seconds", 0, "<1s"},
		{"seconds only", 45, "45s"},
		{"one minute", 60, "1m 0s"},
		{"minutes and seconds", 90, "1m 30s"},
		{"one hour", 3600, "1h 0m"},
		{"hours and minutes", 3661, "1h 1m"},
		{"complex", 7384, "2h 3m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.seconds)
			if result != tt.expected {
				t.Errorf("formatDuration(%d) = %q, expected %q", tt.seconds, result, tt.expected)
			}
		})
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		name     string
		current  int64
		total    int64
		expected string
	}{
		{"zero progress", 0, 100, "0/100"},
		{"partial progress", 50, 100, "50/100"},
		{"complete", 100, 100, "100/100"},
		{"large numbers", 1234, 5678, "1,234/5,678"},
		{"with thousands", 999, 9999, "999/9,999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCount(tt.current, tt.total)
			if result != tt.expected {
				t.Errorf("formatCount(%d, %d) = %q, expected %q", tt.current, tt.total, result, tt.expected)
			}
		})
	}
}
