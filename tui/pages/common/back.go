package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Back() (*Model, []tea.Cmd) {
	var cmds []tea.Cmd

	if len(m.contextToStop) > 0 {
		m.ClearContextToStop()
	} else {
		m.RestorePreviousKeys()
		cmdView := m.RestorePreviousView()
		cmds = append(cmds, cmdView)
	}

	return m, cmds
}
