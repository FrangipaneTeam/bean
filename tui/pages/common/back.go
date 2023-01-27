package common

import (
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Back() (*Model, []tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.k8s.SetTickRunning(false)

	// back to the list, cancel all k8sCmd
	for _, v := range m.k8s.CmdList {
		v.Cancel()
	}

	m.ClearContextToStop()
	// running command °
	m.k8s.CmdList = make(map[string]*k8s.Cmd)

	cmd = m.k8s.GetProgress.SetPercent(0)
	cmds = append(cmds, cmd)

	m.pages.RestorePreviousKeys()
	cmdView := m.pages.RestorePreviousView()

	// error handling
	if m.errorPanel.ErrorRaised() {
		m.errorPanel.Clear()
	}
	cmds = append(cmds, cmdView)

	return m, cmds
}
