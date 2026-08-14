package ui

import (
	"charm.land/lipgloss/v2"
)

type Theme struct {
	Name         string
	ProgressFrom string
	ProgressTo   string
	Unfilled     string
	Success      string
	Error        string
	Warning      string
	Primary      string
	Secondary    string
	Muted        string
	Accent       string
	CopyColor    string
	MoveColor    string
	DeleteColor  string
}

var ThemeMocha = Theme{
	Name:         "catppuccin-mocha",
	ProgressFrom: "#89b4fa",
	ProgressTo:   "#cba6f7",
	Unfilled:     "#313244",
	Success:      "#a6e3a1",
	Error:        "#f38ba8",
	Warning:      "#f9e2af",
	Primary:      "#cdd6f4",
	Secondary:    "#89dceb",
	Muted:        "#585b70",
	Accent:       "#f5c2e7",
	CopyColor:    "#89b4fa",
	MoveColor:    "#fab387",
	DeleteColor:  "#f38ba8",
}

type Styles struct {
	Primary      lipgloss.Style
	Secondary    lipgloss.Style
	Muted        lipgloss.Style
	Accent       lipgloss.Style
	Success      lipgloss.Style
	Error        lipgloss.Style
	Warning      lipgloss.Style
	Blue         lipgloss.Style
	Purple       lipgloss.Style
	Yellow       lipgloss.Style
	Pink         lipgloss.Style
	Green        lipgloss.Style
	Red          lipgloss.Style
	Peach        lipgloss.Style
	Faint        lipgloss.Style
	CopyStyle    lipgloss.Style
	MoveStyle    lipgloss.Style
	DeleteStyle  lipgloss.Style
	ProgressFrom string
	ProgressTo   string
	Unfilled     string
}

func NewStyles(theme Theme) Styles {
	return Styles{
		Primary:      lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Primary)),
		Secondary:    lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Secondary)),
		Muted:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)),
		Accent:       lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent)),
		Success:      lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Success)),
		Error:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Error)),
		Warning:      lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warning)),
		Blue:         lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ProgressFrom)),
		Purple:       lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ProgressTo)),
		Yellow:       lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warning)),
		Pink:         lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent)),
		Green:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Success)),
		Red:          lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Error)),
		Peach:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.MoveColor)),
		Faint:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)),
		CopyStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color(theme.CopyColor)).Bold(true),
		MoveStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color(theme.MoveColor)).Bold(true),
		DeleteStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color(theme.DeleteColor)).Bold(true),
		ProgressFrom: theme.ProgressFrom,
		ProgressTo:   theme.ProgressTo,
		Unfilled:     theme.Unfilled,
	}
}
