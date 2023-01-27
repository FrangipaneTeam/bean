// Package header provides a header for the tui
package header

import (
	"fmt"
	"strings"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	headerBlocks           = 2
	descriptionPaddingLeft = 2
)

// Model is the model of the header.
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
	width                  int
	Notification           string
	NotificationOK         string
	config                 config.Provider
	pages                  *pages.Model
	k8s                    *k8s.Model
}

// Init initializes the model.
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

// New creates a new header model.
func New(title string, desc string, w int, c config.Provider) Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(tui.SpinnerColour)

	return Model{
		Title:          title,
		Description:    desc,
		Notification:   "ready",
		NotificationOK: tui.RunningMark,
		spinner:        s,
		notifyCrds:     make(chan pages.NotifyActivity),
		notifyExamples: make(chan pages.NotifyActivity),
		width:          w,
		config:         c,
		pages:          &pages.Model{ShowDependenciesFiles: true},
		k8s:            &k8s.Model{CmdList: make(map[string]*k8s.Cmd)},
	}
}

// Update updates the model.
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
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case pages.ResponseCRDMsg:
		cmds = append(cmds, pages.WaitForCrdActivity(m.notifyCrds), tick)
		if !m.crdRecentActivity {
			m.hideNotify = 5
			m.crdRecentActivity = true
			m.activityFrom = msg
		}

	case pages.ResponseExamplesMsg:
		cmds = append(cmds, pages.WaitForExamplesActivity(m.notifyExamples), tick)
		if !m.examplesRecentActivity {
			m.hideNotify = 5
			m.examplesRecentActivity = true
			m.activityFrom = msg
		}
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	nameVersion := lipgloss.JoinHorizontal(
		lipgloss.Left,
		tui.TextStyle.Render(m.Title),
		tui.FeintTextStyle.Padding(0, 0, 0, descriptionPaddingLeft).Render(m.Description),
	)

	header := ""

	switch {
	case !m.crdRecentActivity && !m.examplesRecentActivity:
		header = fmt.Sprintf("%s Watch for new crd/examples files", m.spinner.View())
	case m.crdRecentActivity:
		c := lipgloss.NewStyle().Foreground(tui.SpinnerColour).Render("→")
		header = fmt.Sprintf(
			"%s New CRD files %s %s %ds ago",
			m.spinner.View(),
			c,
			"Updated",
			m.hideNotify,
		)
	case m.examplesRecentActivity:
		c := lipgloss.NewStyle().Foreground(tui.SpinnerColour).Render("→")
		header = fmt.Sprintf(
			"%s New examples %s %s %ds ago",
			m.spinner.View(),
			c,
			"Updated",
			m.hideNotify,
		)
	}

	notification := strings.Builder{}
	dependenciesStatus := strings.Builder{}

	t := strings.Trim(m.Notification, "\n")
	if m.k8s.GetRunningCmd() > 0 {
		fmt.Fprintf(
			&notification,
			"%s %s (%d r) %s",
			tui.Divider,
			t,
			m.k8s.GetRunningCmd(),
			m.NotificationOK,
		)
	} else {
		fmt.Fprintf(
			&notification,
			"%s %s %s",
			tui.Divider,
			t,
			m.NotificationOK,
		)
	}

	fmt.Fprintf(&dependenciesStatus, "")
	if m.pages.GetDependenciesStatus() {
		fmt.Fprintf(
			&dependenciesStatus,
			"%s dependencies %s",
			tui.Divider,
			lipgloss.NewStyle().Foreground(tui.NotificationColour).Render("⚠"),
		)
	}

	header = lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		lipgloss.NewStyle().Width(m.width/headerBlocks).Align(lipgloss.Left).Render(header),
		lipgloss.NewStyle().
			Width(m.width/headerBlocks).
			Align(lipgloss.Right).
			Render(dependenciesStatus.String()+notification.String()),
	)

	banner := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().
			MarginBottom(1).
			Width(m.width/headerBlocks).
			Align(lipgloss.Center).
			Render(nameVersion),
		header,
		tui.BorderBottom.Width(m.width).MarginBottom(1).String(),
		// tui.BorderBottom.Width(wP).String(),
	)
	// banner += "\n" + tui.BorderBottom.Width(m.width).String()

	return banner
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}

// Height return the height of the view.
func (m Model) Height() int {
	return lipgloss.Height(m.View())
}

// Width return the width of the view.
func (m Model) Width() int {
	return lipgloss.Width(m.View())
}

// SetWidth set the width of the view.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetPagesModel set the pages model.
func (m *Model) SetPagesModel(p *pages.Model) {
	m.pages = p
}

// SetK8sModel set the k8s model.
func (m *Model) SetK8SModel(k *k8s.Model) {
	m.k8s = k
}
