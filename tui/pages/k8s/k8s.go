package k8s

import (
	"fmt"
	"strings"

	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			m.SetTickRunning(false)

			// back to the list, cancel all k8sCmd
			for _, v := range m.CmdList {
				v.Cancel()
			}
			// running command °
			m.CmdList = make(map[string]*Cmd)
			common.RunningCommands = 0

		case key.Matches(msg, m.keys.Delete):
			cmd = func() tea.Msg {
				return Message{
					ShowDialogBox: true,
					PreviousPage:  m.common.GetViewName(),
				}
			}
			return m, cmd

		case key.Matches(msg, m.keys.Print):
			m.common.SetViewName(common.PPrintActions)
			m.keys.ShowDependanciesFiles.SetEnabled(true)
			return m, nil

		case key.Matches(msg, m.keys.ShowDependanciesFiles):
			m.SwitchDependenciesStatus()
			cmd = m.pages.CurrentList.NewStatusMessage(
				fmt.Sprintf("Show dependencies files →  %t", m.ShowDependenciesFiles),
			)
			common.ShowDependencies = m.ShowDependenciesFiles
			return m, cmd
		}
	}
	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	selected := m.pages.CurrentList.SelectedItem().(*exlist.Example)
	selectedFile := selected.FileWithPath()
	yamlFile := ""

	if selected.HaveSecretFile() {
		yamlFile = fmt.Sprintf("%s,%s.secret", selectedFile, selectedFile)
	} else {
		yamlFile = selectedFile
	}

	str := []string{
		"# Base file :",
		fmt.Sprintf("* kubectl apply -f %s", yamlFile),
		fmt.Sprintf("* kubectl delete -f %s", yamlFile),
		fmt.Sprintf("* kubectl get -f %s", yamlFile),
	}

	if selected.HaveExtraFile() {
		extraFile := selectedFile + ".extra"
		str = append(str,
			"# Extra file:",
			fmt.Sprintf("* kubectl apply -f %s", extraFile),
			fmt.Sprintf("* kubectl delete -f %s", extraFile),
			fmt.Sprintf("* kubectl get -f %s", extraFile),
		)
	}

	if m.ShowDependenciesFiles && selected.HaveDependenciesFiles() {
		files := strings.Join(selected.DependenciesFilesList(), ",")
		str = append(str,
			"# Dependencies :",
			fmt.Sprintf("* kubectl apply -f %s", files),
			fmt.Sprintf("* kubectl delete -f %s", files),
			fmt.Sprintf("* kubectl get -f %s", files),
		)
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.width-2),
		glamour.WithStylePath("dracula"),
	)
	ui, _ := renderer.Render(strings.Join(str, "\n"))
	return ui
}

func (k8sCmd *Cmd) JoinedFiles() string {
	return strings.Join(k8sCmd.Files, ",")
}

// IsTickRunning returns true if the tick is running.
func (m Model) IsTickRunning() bool {
	return m.tickRunning
}

// SetTickRunning sets the tick running state.
func (m *Model) SetTickRunning(state bool) {
	m.tickRunning = state
}

// GetRunningCmd returns the running command.
func (m Model) GetRunningCmd() int {
	return len(m.CmdList)
}

// GetDependenciesStatus returns the dependencies status.
func (m Model) GetDependenciesStatus() bool {
	return m.ShowDependenciesFiles
}

// SwitchDependenciesStatus switches the dependencies status.
func (m *Model) SwitchDependenciesStatus() {
	m.ShowDependenciesFiles = !m.ShowDependenciesFiles
}
