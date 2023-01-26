package common

import (
	"context"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	PViewPort int = iota
	PRoot
	PRessources
	PPrintActions
	PK8S
	PDialogBox
)

type Model struct {
	viewName      int
	oldViewName   int
	keys          *tui.ListKeyMap
	contextToStop []context.CancelFunc
}

func New(keymap *tui.ListKeyMap) *Model {
	return &Model{
		viewName: PRoot,
		keys:     keymap,
	}
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			for _, cancel := range m.contextToStop {
				cancel()
			}
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

// SetView sets the view name.
func (m *Model) SetViewName(name int) {
	m.viewName = name
}

// GetView returns the view name.
func (m *Model) GetViewName() int {
	return m.viewName
}

// GetOldViewName returns the old view name.
func (m *Model) GetOldViewName() int {
	return m.oldViewName
}

// SetOldViewName sets the old view name.
func (m *Model) SetOldViewName(name int) {
	m.oldViewName = name
}

// AddContextToStop adds a context to stop.
func (m *Model) AddContextToStop(ctx context.CancelFunc) {
	m.contextToStop = append(m.contextToStop, ctx)
}

// ClearContextToStop clears the contexts to stop.
func (m *Model) ClearContextToStop() {
	m.contextToStop = nil
}
