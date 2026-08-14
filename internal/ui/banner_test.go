package ui

import (
	"strings"
	"testing"
)

func TestRenderBanner(t *testing.T) {
	banner := RenderBanner(ThemeMocha)

	if !strings.Contains(banner, "____") {
		t.Error("banner should contain the slanted figlet art")
	}

	if !strings.Contains(banner, "blazing fast") {
		t.Error("banner should contain the tagline")
	}
}

func TestRenderBannerGradient(t *testing.T) {
	banner := RenderBanner(ThemeMocha)

	// #89b4fa = rgb(137, 180, 250) — gradient start
	if !strings.Contains(banner, "38;2;137;180;250") {
		t.Error("banner should start the gradient at the theme ProgressFrom color")
	}

	// #cba6f7 = rgb(203, 166, 247) — gradient end
	if !strings.Contains(banner, "38;2;203;166;247") {
		t.Error("banner should end the gradient at the theme ProgressTo color")
	}
}

func TestRenderCompactBanner(t *testing.T) {
	banner := RenderCompactBanner(ThemeMocha)

	if !strings.Contains(banner, "zap") {
		t.Error("compact banner should contain 'zap'")
	}
}
