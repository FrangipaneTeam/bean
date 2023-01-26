package k8s

import (
	"context"
	"strings"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the model of the error panel.
type Model struct {
	common *common.Model
	tea.Model
	keys *tui.ListKeyMap
}

type Message struct {
	ShowDialogBox bool
}

type Cmd struct {
	ID     string
	Done   bool
	Verb   string
	Files  []string
	Kind   string
	Result string
	Cancel context.CancelFunc
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// New returns a new model of the k8s page.
func New(keymap *tui.ListKeyMap, common *common.Model) Model {
	return Model{
		keys:   keymap,
		common: common,
	}
}

// Update updates the model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Delete):
			cmd = func() tea.Msg {
				return Message{ShowDialogBox: true}
			}
			return m, cmd

		case key.Matches(msg, m.keys.Print):
			m.keys.EnableViewPortKeys()
			m.keys.ShowDependanciesFiles.SetEnabled(true)
			m.common.SetViewName(common.PPrintActions)
			return m, nil
		}
	}
	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	return ""
}

func (k8sCmd *Cmd) JoinedFiles() string {
	return strings.Join(k8sCmd.Files, ",")
}
