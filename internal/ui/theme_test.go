package ui

import (
	"testing"
)

func TestNewStyles(t *testing.T) {
	theme := ThemeMocha
	styles := NewStyles(theme)

	// verify styles are created with correct colors
	if styles.ProgressFrom != theme.ProgressFrom {
		t.Errorf("expected ProgressFrom %q, got %q", theme.ProgressFrom, styles.ProgressFrom)
	}

	if styles.ProgressTo != theme.ProgressTo {
		t.Errorf("expected ProgressTo %q, got %q", theme.ProgressTo, styles.ProgressTo)
	}

	if styles.Unfilled != theme.Unfilled {
		t.Errorf("expected Unfilled %q, got %q", theme.Unfilled, styles.Unfilled)
	}
}

func TestThemeColors(t *testing.T) {
	themes := []Theme{
		ThemeMocha,
	}

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			// verify all color fields are set
			if theme.ProgressFrom == "" {
				t.Error("ProgressFrom color is empty")
			}
			if theme.ProgressTo == "" {
				t.Error("ProgressTo color is empty")
			}
			if theme.Unfilled == "" {
				t.Error("Unfilled color is empty")
			}
			if theme.Success == "" {
				t.Error("Success color is empty")
			}
			if theme.Error == "" {
				t.Error("Error color is empty")
			}
			if theme.Warning == "" {
				t.Error("Warning color is empty")
			}
			if theme.Secondary == "" {
				t.Error("Secondary color is empty")
			}
			if theme.Primary == "" {
				t.Error("Primary color is empty")
			}
		})
	}
}
