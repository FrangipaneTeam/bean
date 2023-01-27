package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Model is the model of the error panel.
type Model struct {
	pages *pages.Model
	tea.Model
	keys  *tui.ListKeyMap
	width int
}

type Message struct {
	ShowDialogBox bool
	PreviousPage  pages.PageID
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
func New(keymap *tui.ListKeyMap, pages *pages.Model) Model {
	return Model{
		keys:  keymap,
		pages: pages,
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
		if m.pages.CurrentList.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.Delete):
			cmd = func() tea.Msg {
				return Message{
					ShowDialogBox: true,
					PreviousPage:  m.pages.GetViewName(),
				}
			}
			return m, cmd

		case key.Matches(msg, m.keys.Print):
			m.pages.SetViewName(pages.PPrintActions)
			m.keys.ShowDependanciesFiles.SetEnabled(true)
			return m, nil

		case key.Matches(msg, m.keys.ShowDependanciesFiles):
			m.pages.SwitchDependenciesStatus()
			cmd = m.pages.CurrentList.NewStatusMessage(
				fmt.Sprintf("Show dependencies files →  %t", m.pages.ShowDependenciesFiles),
			)
			return m, cmd
		}
	}
	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	selected := m.pages.CurrentList.SelectedItem().(*tui.Example)
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

	if m.pages.ShowDependenciesFiles && selected.HaveDependenciesFiles() {
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
