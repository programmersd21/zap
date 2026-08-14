package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func RenderBanner(theme Theme) string {
	logo := []string{
		` _____ _ _ __ `,
		`|_  / _` + "`" + ` | '_ \`,
		` /__\__,_| .__/`,
		`         |_|   `,
	}

	var sb strings.Builder
	sb.WriteString("\n")

	colors := []string{
		theme.ProgressFrom,
		"#9cbcfb",
		"#b3affc",
		theme.ProgressTo,
	}

	for i, line := range logo {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Bold(true)
		sb.WriteString(style.Render(line))
		sb.WriteString("\n")
	}

	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted)).
		Italic(true)
	sb.WriteString(tagStyle.Render("blazing fast file ops — gorgeous progress"))
	sb.WriteString("\n\n")

	return sb.String()
}

func RenderCompactBanner(theme Theme) string {
	bolt := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warning)).Bold(true).Render("⚡")
	name := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ProgressFrom)).Bold(true).Render("zap")
	return bolt + " " + name + "\n"
}
