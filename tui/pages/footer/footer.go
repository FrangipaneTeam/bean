// Package footer provides a footer for the application
package footer

import (
	"fmt"
	"strings"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the model of the footer
type Model struct {
	tea.Model
	Message string
	width   int
	Help    help.Model
	Keymap  *tui.ListKeyMap
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// New creates a new footer model
func New(w int, km *tui.ListKeyMap) Model {
	help := help.New()
	help.Styles.ShortSeparator = tui.Ellipsis
	help.Styles.ShortKey = tui.HelpText
	help.Styles.ShortDesc = tui.HelpFeintText
	help.Styles.FullSeparator = tui.Ellipsis
	help.Styles.FullKey = tui.HelpText
	help.Styles.FullDesc = tui.HelpFeintText
	return Model{
		Message: "FrangipaneTeam",
		width:   w,
		Help:    help,
		Keymap:  km,
	}
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

// View renders the model
func (m Model) View() string {
	footer := strings.Builder{}
	message := fmt.Sprintf(
		"%s %s %s",
		tui.Divider, strings.Trim(m.Message, "\n"), tui.Divider,
	)

	f := lipgloss.NewStyle()

	ui := lipgloss.Place(
		m.width,
		0,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(tui.BorderColour).
			Foreground(tui.TextColour).
			Render(message),
	)

	banner := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().PaddingTop(1).Render(m.Help.View(m.Keymap)),
		ui,
	)
	footer.WriteString(f.Render(banner))
	return footer.String()
}

// Height return the height of the view
func (m Model) Height() int {
	return lipgloss.Height(m.View())
}

// Width return the width of the view
func (m Model) Width() int {
	return lipgloss.Width(m.View())
}

// SetWidth sets the width of the view
func (m *Model) SetWidth(w int) {
	m.width = w
}
