package common

import (
	"context"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keys          *tui.ListKeyMap
	contextToStop []context.CancelFunc
	pages         *pages.Model
	header        *header.Model
	errorPanel    *errorpanel.Model
	k8s           *k8s.Model
}

func New(
	pages *pages.Model,
	header *header.Model,
	errorpanel *errorpanel.Model,
	k8s *k8s.Model,
) *Model {
	return &Model{
		keys:          tui.NewListKeyMap(),
		contextToStop: []context.CancelFunc{},
		pages:         pages,
		header:        header,
		errorPanel:    errorpanel,
		k8s:           k8s,
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
