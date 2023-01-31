package footer

import (
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) help() (*Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.common.GetViewName() != common.PK8SGet &&
		m.common.GetViewName() != common.PK8SGetFromRoot {
		m.Help.ShowAll = !m.Help.ShowAll
		m.Help.Width = m.width

		listHeight := m.ex.CurrentList.Height()
		helpHeight := lipgloss.Height(m.Help.View(m.Keymap))

		if m.Help.ShowAll {
			m.listOldHeight = listHeight
			m.centerHeight = listHeight - helpHeight + 1
		} else {
			m.centerHeight = m.listOldHeight
		}
		m.ex.CurrentList.SetHeight(m.centerHeight)

		common.Width = m.width
		common.CenterHeight = m.centerHeight

		cmd = func() tea.Msg {
			msg := common.ResizeMsg{}
			return msg
		}
	}
	return m, cmd
}
