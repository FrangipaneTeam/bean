// Package loading provides a loading page for the tui.
package loading

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/internal/examples"
	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/theme"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/home"
)

type errMsg error

type model struct {
	spinner       spinner.Model
	quitting      bool
	err           error
	width, height int
	errorPanel    *errorpanel.Model
	errorRaised   bool
	config        config.Provider
	theme         theme.Theme
}

// New returns a new model of the loading page.
// nolint: golint // model not used outside of this package
func New(c config.Provider) model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:    s,
		errorPanel: errorpanel.New(0, 0),
		config:     c,
		theme:      theme.Default(),
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.spinner.Tick, tea.EnterAltScreen, examples.LoadExamples(m.config))
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case exlist.LoadedExamples:
		e := home.New(msg, m.width, m.height, m.config)
		cmd = e.Init()
		return e, cmd

	case *errorpanel.ErrorMsg:
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		cmd = m.errorPanel.Init()
		return m, cmd

	case tea.WindowSizeMsg:
		top, right, bottom, left := m.theme.AppStyle.GetMargin()
		m.width, m.height = msg.Width-left-right, msg.Height-top-bottom

		m.errorPanel.SetSize(m.width, m.height)
		return m, nil
	}
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	if m.errorRaised {
		m.errorPanel, cmd = m.errorPanel.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var str string
	if m.errorRaised {
		str = lipgloss.JoinVertical(
			lipgloss.Top,
			m.errorPanel.View(),
			"Press q to quit !",
		)
	} else {
		str = fmt.Sprintf("%s Loading data... Press q to quit\n\n", m.spinner.View())
	}

	return m.theme.AppStyle.Render(str)
}
