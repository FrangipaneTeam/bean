// Package errorpanel provides a panel to display errors
package errorpanel

import (
	"strings"
	"time"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

// Model is the model of the error panel
type Model struct {
	tea.Model
	reason        string
	cause         string
	width, height int
	spinner       spinner.Model
}

// New returns a new model of the error panel
func New(w, h int) Model {
	s := spinner.NewModel()
	s.Spinner = spinner.Spinner{
		Frames: []string{
			"(●    ) W",
			"( ●   ) We",
			"(  ●  ) We'",
			"(   ● ) We'v",
			"(    ●) We've",
			"(    ●) We've ",
			"(   ● ) We've g",
			"(  ●  ) We've go",
			"( ●   ) We've got",
			"(●    ) We've got ",
			"( ●   ) We've got a",
			"(  ●  ) We've got a ",
			"(   ● ) We've got a p",
			"(    ●) We've got a pr",
			"(   ● ) We've got a pro",
			"(  ●  ) We've got a prob",
			"( ●   ) We've got a probl",
			"(●    ) We've got a proble",
			"( ●   ) We've got a problem",
			"(  ●  ) We've got a problem ",
			"(   ● ) We've got a problem !",
		},
		FPS: time.Second / 10,
	}
	s.Style = lipgloss.NewStyle().Foreground(tui.RedColour).Bold(true)

	return Model{
		spinner: s,
		width:   w,
		height:  h,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the model
func (m Model) View() string {
	var b strings.Builder

	reason := lipgloss.JoinVertical(
		lipgloss.Left,
		m.spinner.View(),
		tui.Reason.Margin(2, 0, 0, 0).Render(m.reason),
	)

	desc := tui.Cause.Render(wordwrap.String(m.cause, m.width))

	panel := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Render(wordwrap.String(reason, m.width)),
		lipgloss.NewStyle().Margin(0, 0, 2, 0).Render(desc),
	)

	b.WriteString(panel)
	return b.String()
}

// RaiseError raises an error
func (m Model) RaiseError(reason string, cause error) Model {
	m.reason = reason
	if cause != nil {
		m.cause = cause.Error()
	}

	return m
}

// Resize resizes the model
func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Width returns the width of the model
func (m Model) Width() int {
	return m.width
}

// Height returns the height of the model
func (m Model) Height() int {
	return lipgloss.Height(m.View())
}
