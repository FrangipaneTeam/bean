// Package footer provides a footer for the application
package footer

import (
	"fmt"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the model of the footer
type Model struct {
	tea.Model
	Message       string
	Width, Height int
	Help          help.Model
	Keymap        *tui.ListKeyMap
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// New creates a new footer model
func New(w, h int, km *tui.ListKeyMap) Model {
	help := help.New()
	help.Styles.ShortSeparator = tui.Ellipsis
	help.Styles.ShortKey = tui.HelpText
	help.Styles.ShortDesc = tui.HelpFeintText
	help.Styles.FullSeparator = tui.Ellipsis
	help.Styles.FullKey = tui.HelpText
	help.Styles.FullDesc = tui.HelpFeintText
	return Model{
		Message: "FrangipaneTeam",
		Width:   w,
		Height:  h,
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
	// panel := tui.BorderTop.Width(m.Width).Render(text)
	message := fmt.Sprintf("• %s •", m.Message)
	banner := lipgloss.JoinVertical(
		lipgloss.Center,
		m.Help.View(m.Keymap),
		tui.BorderTop.Width(m.Width).Render(""),
		tui.HightlightTextStyle.Render(message),
	)
	return banner
}

// Resize resizes the model
func (m Model) Resize(width, height int) Model {
	m.Width = width
	m.Height = height
	m.Help.Width = width

	return m
}
