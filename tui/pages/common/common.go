package common

import (
	"context"

	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/elist"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keys          *keymap.ListKeyMap
	contextToStop []context.CancelFunc
	ex            *elist.Model
	errorPanel    *errorpanel.Model
	dialogbox     *dialogbox.Model
	viewName      PageID
	pages         map[PageID]*Page
}

type ResizeMsg struct{}

var (
	Width            int
	Height           int
	CenterHeight     int
	ShowDependencies bool = true
	RunningCommands  int
)

func New(
	keys *keymap.ListKeyMap,
	ex *elist.Model,
	errorpanel *errorpanel.Model,
	dialogbox *dialogbox.Model,
) *Model {
	return &Model{
		keys:          keys,
		contextToStop: []context.CancelFunc{},
		ex:            ex,
		errorPanel:    errorpanel,
		dialogbox:     dialogbox,
		viewName:      PRoot,
		pages:         BeanPages(),
	}
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.ex.CurrentList.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.Back):
			m, cmds = m.Back()
			return m, tea.Batch(cmds...)

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
