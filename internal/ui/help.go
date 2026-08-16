package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

func RenderHelp(cmd *cobra.Command, theme Theme) string {
	var sb strings.Builder

	sb.WriteString(RenderBanner(theme))

	section := func(title, color string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).Render(title))
		sb.WriteString("\n")
	}

	value := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Render
	code := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Secondary)).Render
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)).Italic(true).Render
	flagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ProgressFrom))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Primary)).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted))

	section("usage", theme.ProgressFrom)
	sb.WriteString("  " + value("zap") + " " + descStyle.Render("[flags]") + " " + code("<src>...") + " " + code("<dst>") + "\n\n")

	section("modes", theme.ProgressTo)
	modes := []struct{ flag, name, desc string }{
		{"-c", "copy  ", "copy files or directories (default)"},
		{"-m", "move  ", "move — rename or copy+delete across devices"},
		{"-d", "delete", "delete files or directories"},
		{"-t", "trash ", "move to the system trash"},
	}
	for _, m := range modes {
		sb.WriteString("  " + flagStyle.Render(m.flag) + "  " + nameStyle.Render(m.name) + "  " + descStyle.Render(m.desc) + "\n")
	}
	sb.WriteString("\n")

	section("flags", theme.Warning)
	flags := []struct{ flag, desc string }{
		{"-f, --force       ", "skip confirmation prompts"},
		{"-r, --recursive   ", "required for deleting directories"},
		{"-v, --verbose     ", "repeat for more detail (-v, -vv, -vvv)"},
		{"    --no-preserve-root", "allow operating on the filesystem root"},
		{"    --version     ", "show version"},
	}
	for _, f := range flags {
		sb.WriteString("  " + flagStyle.Render(f.flag) + "  " + descStyle.Render(f.desc) + "\n")
	}
	sb.WriteString("\n")

	section("examples", theme.Success)
	examples := []struct{ cmd, note string }{
		{"zap file.txt backup/", "copy file to directory"},
		{"zap photos/ backup/", "copy directory recursively"},
		{"zap -m old.txt new.txt", "rename / move file"},
		{"zap -t photo.jpg", "move to trash"},
		{"zap -d -r logs/", "delete directory recursively"},
	}
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ProgressFrom))
	noteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)).Italic(true)
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)).Render("  ·  ")
	for _, e := range examples {
		sb.WriteString("  " + cmdStyle.Render(e.cmd) + sep + noteStyle.Render(e.note) + "\n")
	}
	sb.WriteString("\n")

	tipDot := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent)).Render("◆")
	tipText := hint("press ctrl+c at any time to interrupt safely")
	sb.WriteString("  " + tipDot + "  " + tipText + "\n\n")

	return sb.String()
}
