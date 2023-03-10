// Package errorpanel provides a panel to display errors
package errorpanel

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/internal/theme"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
)

const (
	spinnerFPS  = 10
	marginTop   = 2
	marginWidth = 2
)

// Model is the model of the error panel.
type Model struct {
	tea.Model
	reason        string
	cause         string
	width, height int
	spinner       spinner.Model
	errorRaised   bool
	keys          *keymap.ListKeyMap
	theme         theme.Theme
}

// ErrorMsg should be sent to notify a user of an unrecoverable error.
type ErrorMsg struct {
	Reason   string
	Cause    error
	CmdID    string
	Index    int
	Item     *exlist.Example
	FromPage interface{}
}

// New returns a new model of the error panel.
func New(w, h int) *Model {
	theme := theme.Default()
	s := spinner.New()
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
		FPS: time.Second / spinnerFPS,
	}
	s.Style = theme.ErrorPanel.Reason

	keys := keymap.NewListKeyMap()
	keys.EnableRootKeys()
	keys.Back.SetEnabled(true)

	return &Model{
		spinner: s,
		width:   w,
		height:  h,
		keys:    keys,
		theme:   theme,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			m.Clear()
			return m, nil
		}
	case common.ResizeMsg:
		m.SetSize(common.Width, common.CenterHeight)
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the model.
func (m Model) View() string {
	var b strings.Builder

	reason := lipgloss.JoinVertical(
		lipgloss.Left,
		m.spinner.View(),
		m.theme.ErrorPanel.Reason.Margin(marginTop, 0, 0, 0).Render(m.reason),
	)

	desc := m.theme.ErrorPanel.Cause.Render(wordwrap.String(m.cause, m.width))

	panel := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Render(wordwrap.String(reason, m.width)),
		lipgloss.NewStyle().Margin(0, 0, 0, 0).Render(desc),
	)

	panelWithBorder := lipgloss.NewStyle().
		Height(m.height - marginTop).
		Width(m.width - marginWidth).
		MaxHeight(m.height).
		Border(lipgloss.HiddenBorder()).
		Render(panel)

	b.WriteString(panelWithBorder)
	return b.String()
}

// RaiseError raises an error.
func (m *Model) RaiseError(reason string, cause error) *Model {
	m.errorRaised = true
	m.reason = reason
	if cause != nil {
		m.cause = cause.Error()
	}

	return m
}

// ErrorRaised returns true if an error has been raised.
func (m Model) ErrorRaised() bool {
	return m.errorRaised
}

// Clear clears the error.
func (m *Model) Clear() {
	m.errorRaised = false
}

// Resize resizes the model.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Width returns the width of the model.
func (m Model) Width() int {
	return m.width
}

// Height returns the height of the model.
func (m Model) Height() int {
	return lipgloss.Height(m.View())
}
