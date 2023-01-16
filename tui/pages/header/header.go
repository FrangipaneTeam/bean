// Package header provides a header for the tui
package header

import (
	"fmt"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the model of the header
type Model struct {
	tea.Model
	Title                  string
	Description            string
	crdRecentActivity      bool
	examplesRecentActivity bool
	spinner                spinner.Model
	hideNotify             int
	activityFrom           interface{}
	notifyCrds             chan pages.NotifyActivity
	notifyExamples         chan pages.NotifyActivity
	Width, Height          int
	// errorPanel             errorpanel.Model
	// errorRaised            bool
	config config.Provider
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		pages.WaitForCrdActivity(m.notifyCrds),
		pages.WaitForExamplesActivity(m.notifyExamples),
		pages.ListenForCRDActivity(m.notifyCrds, m.config),
		pages.ListenForExamplesActivity(m.notifyExamples, m.config),
		m.spinner.Tick,
		tick,
	)
}

// New creates a new header model
func New(title string, desc string, w int, h int, c config.Provider) Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(tui.SpinnerColour)

	return Model{
		Title:          title,
		Description:    desc,
		spinner:        s,
		notifyCrds:     make(chan pages.NotifyActivity),
		notifyExamples: make(chan pages.NotifyActivity),
		Width:          w,
		Height:         h,
		// errorPanel:     errorpanel.New(w, h),
		config: c,
	}
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tickMsg:
		m.hideNotify--
		if m.hideNotify <= 0 {
			switch m.activityFrom.(type) {
			case pages.ResponseCRDMsg:
				m.crdRecentActivity = false
				// cmds = append(cmds, tools.Kubectl(m.config.Path+"/package/crds", "apply"))
			case pages.ResponseExamplesMsg:
				m.examplesRecentActivity = false
				cmds = append(cmds, pages.LoadExamples(m.config))
			}
		} else {
			return m, tick
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case pages.ResponseCRDMsg:
		if m.crdRecentActivity {
			return m, nil
		}
		m.hideNotify = 5
		m.crdRecentActivity = true
		m.activityFrom = msg
		cmds = append(cmds, pages.WaitForCrdActivity(m.notifyCrds), tick)

	case pages.ResponseExamplesMsg:
		if m.examplesRecentActivity {
			return m, nil
		}
		m.hideNotify = 5
		m.examplesRecentActivity = true
		m.activityFrom = msg
		cmds = append(cmds, pages.WaitForExamplesActivity(m.notifyExamples), tick)

		// case pages.ErrorMsg:
		// 	m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		// 	m.errorRaised = true

		// case tea.WindowSizeMsg:
		// 	m.Title = fmt.Sprintf("%s - %d", "coucou", msg.Width)
		// 	m.Width = msg.Width
		// 	return m, nil
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m Model) View() string {
	nameVersion := lipgloss.JoinHorizontal(
		lipgloss.Left,
		tui.TextStyle.Render(m.Title),
		tui.FeintTextStyle.Padding(0, 0, 0, 2).Render(m.Description),
	)

	var err string
	// if m.errorRaised {
	// 	err = "\n\n" + m.errorPanel.View()
	// }

	header := ""
	// nolint: gocritic // TODO: try to change the logic to avoid this
	if !m.crdRecentActivity && !m.examplesRecentActivity {
		header = fmt.Sprintf("%s Watching CRD files", m.spinner.View())
	} else if m.crdRecentActivity {
		c := lipgloss.NewStyle().Foreground(tui.SpinnerColour).Render("→")
		header = fmt.Sprintf("%s New CRD files %s %s %ds ago", m.spinner.View(), c, "Updated", m.hideNotify)
	} else if m.examplesRecentActivity {
		c := lipgloss.NewStyle().Foreground(tui.SpinnerColour).Render("→")
		header = fmt.Sprintf("%s New examples %s %s %ds ago", m.spinner.View(), c, "Updated", m.hideNotify)
	}
	banner := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().MarginBottom(1).Render(nameVersion),
		lipgloss.NewStyle().Render(header),
		err,
	)
	border := tui.BorderBottom.Width(m.Width).Render(banner)

	return border
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
