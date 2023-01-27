package common

import (
	"context"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keys          *tui.ListKeyMap
	contextToStop []context.CancelFunc
	pages         *pages.Model
}

func New(pages *pages.Model) *Model {
	return &Model{
		keys:          tui.NewListKeyMap(),
		contextToStop: []context.CancelFunc{},
		pages:         pages,
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
		if m.pages.CurrentList.FilterState() == list.Filtering {
			break
		}
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

// AddContextToStop adds a context to stop.
func (m *Model) AddContextToStop(ctx context.CancelFunc) {
	m.contextToStop = append(m.contextToStop, ctx)
}

// ClearContextToStop clears the contexts to stop.
func (m *Model) ClearContextToStop() {
	m.contextToStop = []context.CancelFunc{}
}
