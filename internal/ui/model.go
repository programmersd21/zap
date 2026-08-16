package ui

import (
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Op int

const (
	OpCopy Op = iota
	OpMove
	OpDelete
	OpTrash
)

type ProgressMsg struct {
	BytesDone   int64
	FilesDone   int64
	CurrentFile string
}

type CompletedMsg struct{}

type ErrorMsg struct {
	Err error
}

type ConfirmMsg struct {
	Path   string
	Answer chan bool
}

type confirmState struct {
	path   string
	answer chan bool
}

type tickMsg time.Time
type doneQuitMsg struct{}

var spinFrame = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Model struct {
	op             Op
	styles         Styles
	verbose        bool
	totalBytes     int64
	totalFiles     int64
	bytesDone      int64
	filesDone      int64
	currentFile    string
	startTime      time.Time
	lastUpdate     time.Time
	bytesPerSec    float64
	progress       progress.Model
	completedFiles []string
	maxLogLines    int
	spinIdx        int
	done           bool
	err            error
	quitting       bool
	prompt         *confirmState
	termWidth      int
}

func NewModel(theme Theme, op Op, verbose bool, totalBytes, totalFiles int64) Model {
	styles := NewStyles(theme)
	prog := progress.New(
		progress.WithColors(lipgloss.Color(styles.ProgressFrom), lipgloss.Color(styles.ProgressTo)),
		progress.WithoutPercentage(),
		progress.WithWidth(44),
	)

	return Model{
		op:             op,
		styles:         styles,
		verbose:        verbose,
		totalBytes:     totalBytes,
		totalFiles:     totalFiles,
		startTime:      time.Now(),
		lastUpdate:     time.Now(),
		progress:       prog,
		completedFiles: make([]string, 0),
		maxLogLines:    8,
		termWidth:      80,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.prompt != nil {
			var ok, handled bool
			switch msg.String() {
			case "y", "Y", "enter":
				ok, handled = true, true
			case "n", "N", "esc", "ctrl+c":
				ok, handled = false, true
			}
			if handled {
				m.prompt.answer <- ok
				m.prompt = nil
				if msg.String() == "ctrl+c" {
					m.quitting = true
					return m, tea.Quit
				}
			}
			return m, nil
		}

		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	case ConfirmMsg:
		m.prompt = &confirmState{path: msg.Path, answer: msg.Answer}
		return m, nil

	case tickMsg:
		if !m.done {
			m.spinIdx = (m.spinIdx + 1) % len(spinFrame)
			return m, tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		barWidth := m.termWidth - 36
		if barWidth < 20 {
			barWidth = 20
		} else if barWidth > 60 {
			barWidth = 60
		}
		m.progress = progress.New(
			progress.WithColors(lipgloss.Color(m.styles.ProgressFrom), lipgloss.Color(m.styles.ProgressTo)),
			progress.WithoutPercentage(),
			progress.WithWidth(barWidth),
		)

	case ProgressMsg:
		if m.verbose && msg.CurrentFile != "" && msg.CurrentFile != m.currentFile && m.currentFile != "" {
			m.completedFiles = append(m.completedFiles, m.currentFile)
			if len(m.completedFiles) > m.maxLogLines {
				m.completedFiles = m.completedFiles[1:]
			}
		}

		now := time.Now()
		elapsed := now.Sub(m.lastUpdate).Seconds()
		if elapsed > 0 && m.totalBytes > 0 {
			instant := float64(msg.BytesDone-m.bytesDone) / elapsed
			if m.bytesPerSec == 0 {
				m.bytesPerSec = instant
			} else {
				m.bytesPerSec = 0.15*instant + 0.85*m.bytesPerSec
			}
		}
		m.lastUpdate = now

		m.bytesDone = msg.BytesDone
		m.filesDone = msg.FilesDone
		m.currentFile = msg.CurrentFile
		return m, nil

	case CompletedMsg:
		m.done = true
		return m, tea.Tick(350*time.Millisecond, func(t time.Time) tea.Msg {
			return doneQuitMsg{}
		})

	case doneQuitMsg:
		return m, tea.Quit

	case ErrorMsg:
		m.err = msg.Err
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}
