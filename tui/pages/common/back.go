package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Back() (*Model, []tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	m.ClearContextToStop()

	m.RestorePreviousKeys()
	cmdView := m.RestorePreviousView()

	// error handling
	if m.errorPanel.ErrorRaised() {
		m.errorPanel.Clear()
	}
	cmds = append(cmds, cmdView)

	return m, cmds
}
