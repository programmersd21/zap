package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d b", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"kb", "mb", "gb", "tb"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

func FormatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	str := fmt.Sprintf("%d", n)
	var result []rune
	for i, r := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, r)
	}
	return string(result)
}

func formatCount(current, total int64) string {
	return fmt.Sprintf("%s/%s", FormatNumber(current), FormatNumber(total))
}

func formatDuration(seconds int64) string {
	if seconds == 0 {
		return "<1s"
	}
	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	default:
		return fmt.Sprintf("%dh %dm", seconds/3600, (seconds%3600)/60)
	}
}

func shorten(path string, maxLen int) string {
	if maxLen <= 0 {
		maxLen = 60
	}
	if len(path) > maxLen {
		return "…" + path[len(path)-maxLen+1:]
	}
	return path
}

func (m Model) opBadge() string {
	switch m.op {
	case OpMove:
		return m.styles.MoveStyle.Render("move")
	case OpDelete:
		return m.styles.DeleteStyle.Render("delete")
	case OpTrash:
		return m.styles.TrashStyle.Render("trash")
	case OpShred:
		return m.styles.ShredStyle.Render("shred")
	default:
		return m.styles.CopyStyle.Render("copy")
	}
}

func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	if m.err != nil {
		icon := m.styles.Error.Bold(true).Render("✗")
		msg := m.styles.Error.Render(m.err.Error())
		return tea.NewView("\n" + icon + "  " + msg + "\n")
	}

	var sb strings.Builder
	sb.WriteString("\n")

	badge := m.opBadge()
	muted := m.styles.Muted
	accent := m.styles.Accent

	if !m.done {
		spinner := accent.Render(spinFrame[m.spinIdx])
		sb.WriteString("  ")
		sb.WriteString(badge)
		sb.WriteString("  ")
		sb.WriteString(spinner)
		sb.WriteString("  ")
		if m.currentFile != "" {
			pathWidth := m.termWidth - 20
			sb.WriteString(muted.Italic(true).Render(shorten(m.currentFile, pathWidth)))
		} else {
			sb.WriteString(muted.Render("scanning…"))
		}
		sb.WriteString("\n\n")
	}

	var percent float64
	switch {
	case m.totalBytes > 0:
		percent = float64(m.bytesDone) / float64(m.totalBytes)
	case m.totalFiles > 0:
		percent = float64(m.filesDone) / float64(m.totalFiles)
	}

	sb.WriteString("  ")
	sb.WriteString(m.progress.ViewAs(percent))
	sb.WriteString("\n")

	if !m.done {
		var parts []string
		if m.totalFiles > 0 {
			files := m.styles.Blue.Render(formatCount(m.filesDone, m.totalFiles))
			parts = append(parts, muted.Render("files")+lipgloss.NewStyle().Render(" ")+files)
		}
		if m.totalBytes > 0 {
			data := m.styles.Purple.Render(
				fmt.Sprintf("%s / %s", FormatBytes(m.bytesDone), FormatBytes(m.totalBytes)),
			)
			parts = append(parts, muted.Render("data")+lipgloss.NewStyle().Render(" ")+data)
		}
		if m.bytesPerSec > 0 {
			speed := m.styles.Yellow.Render(FormatBytes(int64(m.bytesPerSec)) + "/s")
			parts = append(parts, speed)

			if m.totalBytes > 0 {
				rem := float64(m.totalBytes-m.bytesDone) / m.bytesPerSec
				if rem > 0 && rem < 86400 {
					eta := muted.Render("eta ") + m.styles.Secondary.Render(formatDuration(int64(rem)))
					parts = append(parts, eta)
				}
			}
		}

		if len(parts) > 0 {
			dot := muted.Render(" · ")
			sb.WriteString("  ")
			sb.WriteString(strings.Join(parts, dot))
			sb.WriteString("\n")
		}
	}

	if m.verbose && len(m.completedFiles) > 0 {
		sb.WriteString("\n")
		for _, file := range m.completedFiles {
			tick := m.styles.Success.Render("✓")
			sb.WriteString("  ")
			sb.WriteString(tick)
			sb.WriteString("  ")
			sb.WriteString(muted.Render(shorten(file, m.termWidth-8)))
			sb.WriteString("\n")
		}
	}

	if m.prompt != nil {
		sb.WriteString("\n")
		warn := m.styles.Warning.Bold(true).Render("overwrite")
		path := m.styles.Primary.Render(shorten(m.prompt.path, m.termWidth-20))
		hint := muted.Render("[Y/n]")
		sb.WriteString("  ")
		sb.WriteString(warn)
		sb.WriteString("  ")
		sb.WriteString(path)
		sb.WriteString("  ")
		sb.WriteString(hint)
		sb.WriteString("\n")
	}

	if m.done {
		sb.WriteString("\n")
		elapsed := time.Since(m.startTime)
		icon := m.styles.Success.Bold(true).Render("✓")

		fileWord := "files"
		if m.filesDone == 1 {
			fileWord = "file"
		}

		var summary string
		switch {
		case m.totalBytes > 0 && m.totalFiles > 0:
			summary = fmt.Sprintf(
				"%s  %s %s · %s in %s",
				icon,
				FormatNumber(m.filesDone),
				fileWord,
				FormatBytes(m.bytesDone),
				formatDuration(int64(elapsed.Seconds())),
			)
		case m.totalFiles > 0:
			summary = fmt.Sprintf(
				"%s  %s %s in %s",
				icon,
				FormatNumber(m.filesDone),
				fileWord,
				formatDuration(int64(elapsed.Seconds())),
			)
		default:
			summary = fmt.Sprintf("%s  done in %s", icon, formatDuration(int64(elapsed.Seconds())))
		}

		sb.WriteString("  ")
		sb.WriteString(m.styles.Success.Render(summary))
		sb.WriteString("\n")
	}

	return tea.NewView(sb.String())
}
