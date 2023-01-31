package elist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// UpdateExamplesList updates the examples list.
func (m *Model) UpdateExamplesList(examples map[string][]list.Item) {
	m.exampleList = examples
}

func (m *Model) UpdateList(params ...string) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	title := m.listName
	if len(params) != 0 {
		title = params[0]
	}

	if _, ok := m.exampleList[title]; ok {
		i := m.exampleList[title]
		cmd = m.CurrentList.SetItems(i)
		m.listName = title
	}
	return m, cmd
}

func (m *Model) SetRootItems() (*Model, tea.Cmd) {
	cmd := m.CurrentList.SetItems(m.exampleList["-"])
	return m, cmd

}
