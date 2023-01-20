// Package md provides a markdown viewer
package md

import (
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the model of the markdown viewer
type Model struct {
	tea.Model
	Width, Height int
	Viewport      viewport.Model
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// New creates a new markdown viewer model
func New(w, h int) Model {
	vp := viewport.New(w, h)
	vp.Style = lipgloss.NewStyle().
		Width(w).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(tui.BorderColour)
	// PaddingRight(2)
	return Model{
		Width:    w,
		Height:   h,
		Viewport: vp,
	}
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	return m, nil
}

// View renders the model
func (m Model) View() string {
	return ""
}

// SetSize sets the size of the model
func (m *Model) SetSize(w, h int) {
	m.Width = w
	m.Height = h
	m.Viewport.Height = h
	m.Viewport.Width = w
	m.Viewport.Style = lipgloss.NewStyle().
		Width(w).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(tui.BorderColour)
	viewport.Sync(m.Viewport)
}
