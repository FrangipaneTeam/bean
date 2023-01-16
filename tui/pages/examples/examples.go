// Package examples provides the examples page.
package examples

import (
	"fmt"
	"strings"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tools"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/footer"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/md"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	pViewPort     = "viewport"
	pRoot         = "root"
	pRessources   = "ressources"
	pPrintActions = "printActions"
	pK8S          = "k8s"
)

type model struct {
	exampleList map[string][]list.Item
	currentList list.Model

	listName      string
	width, height int
	keys          *tui.ListKeyMap

	viewName string

	showDependenciesFiles bool

	errorRaised bool

	header     header.Model
	footer     footer.Model
	errorPanel errorpanel.Model
	markdown   md.Model

	config config.Provider

	k8sOutput string
}

// New returns a new model of the examples page.
// nolint: golint // model not used outside of this package
func New(e tui.LoadedExamples, width, height int, c config.Provider) model {
	listKeys := tui.NewListKeyMap()
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(tui.HighlightColour).
		Foreground(tui.HighlightColour).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderForeground(tui.HighlightColour).
		Foreground(tui.HighlightFeintColour)

	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.
		Foreground(tui.FeintColour)

	delegate.Styles.FilterMatch = lipgloss.NewStyle().
		Underline(true).
		Bold(true)

	list := list.New(e.Examples["-"],
		delegate,
		width,
		int(float64(height)*0.6),
	)
	list.Title = "Choose an example"
	list.DisableQuitKeybindings()
	list.SetShowHelp(false)
	// list.Help = help.Model{}`

	header := header.New(
		"Bean",
		"A FrangipaneTeam bin",
		width,
		int(float64(height)*0.2),
		c,
	)
	footer := footer.New(width, int(float64(width)*0.2), listKeys)

	return model{
		exampleList: e.Examples,
		currentList: list,
		listName:    "-",
		viewName:    pRoot,

		keys:       listKeys,
		header:     header,
		footer:     footer,
		errorPanel: errorpanel.New(width, int(float64(height)*0.6)),
		markdown:   md.New(width, int(float64(height)*0.6)),
		width:      width,
		height:     height,
		config:     c,

		showDependenciesFiles: true,
	}
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	m.keys.RootHelp()
	return tea.Batch(
		tea.EnterAltScreen,
		m.header.Init(),
		m.footer.Init(),
		m.markdown.Init(),
	)
}

// Update updates the model.
// nolint: gocyclo // TODO: refactor
// TODO show SetEnabled to enable/disable keys
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.currentList.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Back):
			m.errorRaised = false

			switch m.viewName {
			// case "root":
			// 	m.keys.RootHelp()
			// 	return m, nil
			case pRessources:
				m.viewName = pRoot
				m.keys.RootHelp()
				cmd = m.currentList.NewStatusMessage("back to home !")
				m.currentList.ResetSelected()
				cmds = append(cmds, cmd)
				m, cmd = m.showExamples()
			case pViewPort:
				m.viewName = pRoot
				m.keys.RootHelp()
			case pPrintActions, pK8S:
				m.viewName = pRessources
				m.keys.YamlHelp()
				cmd = m.currentList.NewStatusMessage("back to " + m.listName)
				m.currentList.ResetSelected()
				cmds = append(cmds, cmd)
				m, cmd = m.showYaml(m.listName)
			}

			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keys.Enter):
			if m.viewName != pRoot {
				return m, nil
			}
			title := m.currentList.SelectedItem().(*tui.Example).Title()
			m, cmd = m.showYaml(title)
			m.viewName = pRessources
			m.keys.YamlHelp()
			return m, cmd

		case key.Matches(msg, m.keys.Print):
			if m.viewName == pRessources {
				m.viewName = pK8S
				// m.k8sOutput = fmt.Sprintf("coucou %s", m.currentList.SelectedItem().(*tui.Example).Title())
				m.viewName = pPrintActions
				m.keys.YamlActionHelp()
			}
			return m, nil

		case key.Matches(msg, m.keys.Help):
			if m.viewName != pPrintActions && m.viewName != pK8S {
				m.footer.Help.ShowAll = !m.footer.Help.ShowAll
				m.footer.Help.Width = m.width
			}
			return m, nil

		case key.Matches(msg, m.keys.ShowRessources):
			if m.viewName == pRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-resources.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.ShowTested):
			if m.viewName == pRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-tested.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.GenerateListTested):
			cmd = pages.GenerateListTested(m.config)
			return m, cmd

		case key.Matches(msg, m.keys.ShowDependanciesFiles):
			m.showDependenciesFiles = !m.showDependenciesFiles
			cmd = m.currentList.NewStatusMessage(fmt.Sprintf("Show dependencies files →  %t", m.showDependenciesFiles))
			return m, cmd

		case key.Matches(msg, m.keys.Get), key.Matches(msg, m.keys.Apply), key.Matches(msg, m.keys.Delete):
			if m.viewName == pRessources {
				file := m.currentList.SelectedItem().(*tui.Example).Title()
				extra := m.currentList.SelectedItem().(*tui.Example).HaveExtraFile()
				secret := m.currentList.SelectedItem().(*tui.Example).HaveSecretFile()

				if extra {
					file += fmt.Sprintf(",%s.extra", file)
				}

				if secret {
					file += fmt.Sprintf(",%s.secret", file)
				}

				if m.showDependenciesFiles && m.currentList.SelectedItem().(*tui.Example).HaveDependenciesFiles() {
					d := strings.Join(m.currentList.SelectedItem().(*tui.Example).DependenciesFilesList(), ",")
					file += fmt.Sprintf(",%s", d)
				}

				verb := "unknown"
				switch {
				case key.Matches(msg, m.keys.Get):
					m.viewName = pK8S
					m.keys.OnlyBackQuit()
					verb = "get"

				case key.Matches(msg, m.keys.Apply):
					verb = "apply"
				case key.Matches(msg, m.keys.Delete):
					verb = "delete"
				}
				cmd = tools.Kubectl(verb, file)
				return m, cmd
			}

		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := tui.AppStyle.GetMargin()
		m.width, m.height = msg.Width-left-right, msg.Height-top-bottom

		m.header.Width = m.width
		m.header.Height = int(float64(m.height) * 0.2)
		m.footer.Width = m.width
		m.footer.Height = int(float64(m.height) * 0.2)

		m.markdown.Width = m.width
		m.markdown.Viewport.Width = m.width
		m.markdown.Viewport.Height = int(float64(m.height) * 0.6)
		// m.footer.Help.Height = int(float64(m.height) * 0.15)

		m.currentList.SetSize(m.width, int(float64(m.height)*0.6))
		m.errorPanel.Resize(m.width, int(float64(m.height)*0.6))
		// m.footer.Resize(m.width, int(float64(m.height)*0.15))
		// m.header.Resize(m.width, int(float64(m.height)*0.15))

	case tui.LoadedExamples:
		m.exampleList = msg.Examples
		cmd = m.currentList.SetItems(msg.Examples[m.listName])
		return m, cmd

	case tools.Markdown:
		m.viewName = pViewPort
		m.keys.ViewPortHelp()

		m.markdown.Viewport.SetContent(msg.Content)
		m.markdown.Viewport.GotoTop()
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		return m, cmd

	case tools.ErrorMsg:
		cmd = m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
		m.keys.ErrorHelp()
		return m, cmd

	case tui.ListTestedDone:
		cmd := m.currentList.NewStatusMessage("List tested generated")
		return m, cmd

	case tools.KubectlResult:
		switch msg.Verb {
		case "apply", "delete":
			cmd := m.currentList.NewStatusMessage(fmt.Sprintf("kubectl %s ok", msg.Verb))
			return m, cmd

		case "get":
			m.viewName = pK8S
			m.k8sOutput = msg.Out
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	if m.viewName == pRoot || m.viewName == pRessources {
		newListModel, cmd := m.currentList.Update(msg)
		m.currentList = newListModel
		cmds = append(cmds, cmd)
	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	if m.errorRaised {
		m.errorPanel, cmd = m.errorPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.viewName == pViewPort {
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns the string representation of the model.
func (m model) View() string {
	var view string
	if m.errorRaised {
		errorP := lipgloss.NewStyle().Height(m.currentList.Height()).Render(m.errorPanel.View())
		view = lipgloss.JoinVertical(
			lipgloss.Center,
			m.header.View(),
			errorP,
			m.footer.View(),
		)

		// return errorP
	} else {
		switch m.viewName {
		case pViewPort:
			view = lipgloss.JoinVertical(
				lipgloss.Center,
				m.header.View(),
				m.markdown.Viewport.View(),
				m.footer.View(),
			)

		case pK8S:
			get := lipgloss.NewStyle().Height(m.currentList.Height()).Render(m.k8sOutput)
			view = lipgloss.JoinVertical(
				lipgloss.Left,
				m.header.View(),
				get,
				m.footer.View(),
			)

		case pRoot, pRessources:
			view = lipgloss.JoinVertical(
				lipgloss.Left,
				m.header.View(),
				m.currentList.View(),
				m.footer.View(),
			)

		case pPrintActions:
			yamlFile := m.currentList.SelectedItem().(*tui.Example).Title()

			if m.currentList.SelectedItem().(*tui.Example).HaveSecretFile() {
				yamlFile = fmt.Sprintf("%s,%s.secret", yamlFile, yamlFile)
			}

			str := []string{
				lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl apply -f " + yamlFile),
				lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl delete -f " + yamlFile),
				lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl get -f " + yamlFile),
			}

			if m.currentList.SelectedItem().(*tui.Example).HaveExtraFile() {
				extraFile := yamlFile + ".extra"
				str = append(str,
					lipgloss.NewStyle().Padding(2, 0, 2, 0).Underline(true).Render("Extra file:"),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl apply -f "+extraFile),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl delete -f "+extraFile),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl get -f "+extraFile),
				)
			}

			if m.showDependenciesFiles && m.currentList.SelectedItem().(*tui.Example).HaveDependenciesFiles() {
				files := strings.Join(m.currentList.SelectedItem().(*tui.Example).DependenciesFilesList(), ",")
				str = append(str,
					lipgloss.NewStyle().Padding(2, 0, 2, 0).Underline(true).Render("Dependencies file:"),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl apply -f "+files),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl delete -f "+files),
					lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("kubectl get -f "+files),
				)
			}

			actions := lipgloss.JoinVertical(
				lipgloss.Center,
				str...,
			)

			actions = lipgloss.NewStyle().Copy().Align(lipgloss.Center, lipgloss.Center).Foreground(tui.HighlightColour).Height(m.currentList.Height()).Width(m.width).Render(actions)

			view = lipgloss.JoinVertical(
				lipgloss.Center,
				m.header.View(),
				actions,
				m.footer.View(),
			)
		}
	}
	return tui.AppStyle.Render(view)
}

func (m model) showExamples() (model, tea.Cmd) {
	i := m.exampleList["-"]
	cmd := m.currentList.SetItems(i)
	m.currentList.Title = "Choose an example"
	m.listName = "-"

	return m, cmd
}

func (m model) showYaml(title string) (model, tea.Cmd) {
	var cmd tea.Cmd
	if _, ok := m.exampleList[title]; ok {
		i := m.exampleList[title]
		m.currentList.ResetSelected()
		cmd = m.currentList.SetItems(i)
		m.currentList.Title = "Working on " + title
		m.listName = title
	}
	return m, cmd
}
