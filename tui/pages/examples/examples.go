// Package examples provides the examples page.
package examples

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FrangipaneTeam/bean/tools"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	k8sDelete  = "delete"
	k8sApply   = "apply"
	k8sManaged = "managed"

	k8sProgressIncrement     = 0.1
	progressWidth        int = 10
)

func randSeq(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.header.Init(),
		m.footer.Init(),
		m.markdown.Init(),
	)
}

// Update updates the model.
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

		case key.Matches(msg, m.keys.Back):
			m.errorRaised = false
			m.tickRunning = false

			// back to the list, cancel all k8sCmd
			for _, v := range m.k8sCmdList {
				v.Cancel()
			}

			m.common.ClearContextToStop()
			m.header.RunningCommands = 0

			m.k8sCmdList = make(map[string]*k8s.Cmd)

			cmd = m.progressK8SGet.SetPercent(0)
			cmds = append(cmds, cmd)

			// error handling
			if m.errorRaised {
				m.errorRaised = false
				*m.keys = *m.oldKeys
				return m, nil
			}

			switch view := m.common.GetViewName(); view {
			case common.PRessources:
				m, cmd = m.rootView()
				return m, cmd

			case common.PViewPort:
				m.common.SetViewName(common.PRoot)
				m.keys.EnableRootKeys()
				return m, cmd

			case common.PPrintActions, common.PK8S, common.PDialogBox:
				if view == common.PK8S && !m.keys.Apply.Enabled() {
					m, cmd = m.rootView()
					return m, cmd
				}

				m.keys.EnableKindListKeys()
				m.common.SetViewName(common.PRessources)
				cmd = m.currentList.NewStatusMessage("back to " + m.listName)
				cmds = append(cmds, cmd)
				m, cmd = m.showYaml(m.listName)
			}

			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		case key.Matches(msg, m.keys.Select):

			switch view := m.common.GetViewName(); view {
			case common.PDialogBox:
				var newModel model
				var k8sCmd *k8s.Cmd

				if m.dialogbox.ActiveButton == dialogbox.GetCancelValue() {
					newModel = m
					newModel.header.Notification = "cancel delete"
					newModel.header.NotificationOK = tui.ErrorMark
					newModel.common.SetViewName(m.common.GetOldViewName())
					*newModel.keys = *m.oldKeys
				} else {
					m.header.NotificationOK = tui.RunningMark
					newModel, k8sCmd, cmd = m.generateK8SFiles()
					if cmd != nil {
						newModel.common.SetViewName(m.common.GetOldViewName())
						*newModel.keys = *m.oldKeys
						return newModel, cmd
					}

					newModel.k8sProgressMsg = "delete sent !"
					k8sCmd.Verb = k8sDelete

					newModel.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
					newModel.header.NotificationOK = tui.RunningMark

					ctx, cancel := context.WithCancel(context.Background())
					k8sCmd.Cancel = cancel
					m.common.AddContextToStop(cancel)
					m.k8sCmdList[k8sCmd.ID] = k8sCmd
					cmd = tools.Kubectl(ctx, k8sCmd)

					newModel.common.SetViewName(m.common.GetOldViewName())
					*newModel.keys = *m.oldKeys

					cmds = append(cmds, cmd)
				}

				if newModel.common.GetViewName() == common.PK8S {
					cmds = append(cmds, m.tickCmd())
				}

				return newModel, tea.Batch(cmds...)

			case common.PRoot:
				title := m.currentList.SelectedItem().(*tui.Example).Title()
				if m.currentList.FilterState() == list.FilterApplied {
					m.currentList.ResetFilter()
				}

				m.keys.EnableKindListKeys()

				m, cmd = m.showYaml(title)
				m.common.SetViewName(common.PRessources)

				m.currentList.Select(0)

				return m, cmd
			}

		case key.Matches(msg, m.keys.Help):
			if m.common.GetViewName() != common.PK8S {
				m.footer.Help.ShowAll = !m.footer.Help.ShowAll
				m.footer.Help.Width = m.width

				_, v := tui.AppStyle.GetFrameSize()
				listHeight := m.currentList.Height()
				helpHeight := lipgloss.Height(m.footer.Help.View(m.keys))

				if m.footer.Help.ShowAll {
					m.listOldHeight = listHeight
					m.centerHeight = listHeight - helpHeight + 1
				} else {
					m.centerHeight = m.listOldHeight
				}
				m.currentList.SetHeight(m.centerHeight)
				m.errorPanel.SetSize(m.width-v, m.centerHeight)
				m.markdown.SetSize(m.width, m.centerHeight)
			}
			return m, cmd

		case key.Matches(msg, m.keys.ShowRessources):
			m.keys.EnableViewPortKeys()

			if m.common.GetViewName() == common.PRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-resources.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.ShowTested):
			m.keys.EnableViewPortKeys()
			if m.common.GetViewName() == common.PRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-tested.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.GenerateListTested):
			cmd = pages.GenerateListTested(m.config)
			return m, cmd

		case key.Matches(msg, m.keys.ShowDependanciesFiles):
			m.showDependenciesFiles = !m.showDependenciesFiles
			m.header.DependenciesStatus = m.showDependenciesFiles
			cmd = m.currentList.NewStatusMessage(fmt.Sprintf("Show dependencies files →  %t", m.showDependenciesFiles))
			return m, cmd

		case key.Matches(msg, m.keys.Get), key.Matches(msg, m.keys.Apply):
			// TODO: in witch case m.currentList.SelectedItem() == nil ?
			var k8sCmd *k8s.Cmd
			m, k8sCmd, cmd = m.generateK8SFiles()
			if cmd != nil {
				return m, cmd
			}

			switch {
			case key.Matches(msg, m.keys.Get):
				m.common.SetViewName(common.PK8S)
				m.k8sCurrentIDView = k8sCmd.ID
				m.keys.EnableManagedKeys()
				k8sCmd.Verb = k8sManaged

			case key.Matches(msg, m.keys.Apply):
				m.k8sProgressMsg = "apply sent !"
				k8sCmd.Verb = k8sApply

			case key.Matches(msg, m.keys.Delete):
				m.k8sProgressMsg = "delete sent !"
				k8sCmd.Verb = k8sDelete
			}

			m.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
			m.header.NotificationOK = tui.RunningMark

			ctx, cancel := context.WithCancel(context.Background())
			k8sCmd.Cancel = cancel
			m.common.AddContextToStop(cancel)
			m.k8sCmdList[k8sCmd.ID] = k8sCmd
			m.header.RunningCommands++
			cmd = tools.Kubectl(ctx, k8sCmd)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		// m.header.Notification = "resizing"
		headerHeight := m.header.Height()
		footerHeight := m.footer.Height()

		h, v := tui.AppStyle.GetFrameSize()

		m.width, m.height = msg.Width-h, msg.Height-v
		centerH := m.height - headerHeight - footerHeight
		m.centerHeight = centerH

		m.header.SetWidth(m.width)
		m.footer.SetWidth(m.width)

		m.currentList.SetSize(m.width, centerH)
		m.dialogbox.SetSize(m.width, centerH)
		m.errorPanel.SetSize(m.width, centerH)
		m.markdown.SetSize(m.width, centerH)
		// m.header.NotificationOK = tui.CheckMark
		return m, nil

	case k8s.Message:
		question := "Delete all ressources ?"
		if m.showDependenciesFiles {
			question = "Delete all ressources WITH dependencies ?"
		}
		okValue := "No Fear !"
		cancelValue := "I'm scared !"
		m.dialogbox.SetDialogBox(question, okValue, cancelValue)
		m.common.SetOldViewName(m.common.GetViewName())
		m.oldKeys = &tui.ListKeyMap{}
		*m.oldKeys = *m.keys
		m.common.SetViewName(common.PDialogBox)
		m.keys.EnableDialogBoxKeys()
		return m, nil

	case tui.LoadedExamples:
		m.header.Notification = fmt.Sprintf("loaded new examples @ %s", time.Now().Format("15:04:05"))
		m.header.NotificationOK = tui.CheckMark
		m.exampleList = msg.Examples
		m.currentList.SetItems(msg.Examples[m.listName])

	case tools.Markdown:
		m.common.SetViewName(common.PViewPort)
		m.markdown.Viewport.SetContent(msg.Content)
		m.markdown.Viewport.GotoTop()
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		return m, cmd

	case tools.ErrorMsg:
		m.oldKeys = &tui.ListKeyMap{}
		*m.oldKeys = *m.keys
		m.keys.EnableErrorKeys()
		cmd = m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
		m.header.NotificationOK = tui.ErrorMark
		return m, cmd

	case tui.ListTestedDone:
		cmd = m.currentList.NewStatusMessage("List tested generated")
		return m, cmd

	case *k8s.Cmd:
		m.header.RunningCommands--
		m.k8sCmdList[msg.ID].Done = true
		m.header.NotificationOK = tui.CheckMark
		m.k8sProgressMsg = ""
		m.header.Notification = fmt.Sprintf("k %s @ %s", msg.Verb, time.Now().Format("15:04:05"))

		if msg.Verb == k8sManaged {
			m.common.SetViewName(common.PK8S)
			m.k8sCurrentIDView = msg.ID
			if !m.tickRunning {
				m.tickRunning = true
				cmd = m.tickCmd()
			}
		}
		return m, cmd

	case tickK8SGet:
		if m.common.GetViewName() == common.PDialogBox {
			return m, nil
		}
		if m.tickRunning {
			var kubectlCmd tea.Cmd
			if m.progressK8SGet.Percent() == 1.0 {
				m.progressK8SGet.SetPercent(0)

				ctx, cancel := context.WithCancel(context.Background())
				m.k8sCmdList[m.k8sCurrentIDView].Cancel = cancel
				m.common.AddContextToStop(cancel)
				m.header.RunningCommands++
				kubectlCmd = tools.Kubectl(ctx, m.k8sCmdList[m.k8sCurrentIDView])
			}

			// Note that you can also use progress.Model.SetPercent to set the
			// percentage value explicitly, too.
			cmd = m.progressK8SGet.IncrPercent(k8sProgressIncrement)
			return m, tea.Batch(m.tickCmd(), cmd, kubectlCmd)
		}

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var progressModel tea.Model
		progressModel, cmd = m.progressK8SGet.Update(msg)
		m.progressK8SGet = progressModel.(progress.Model)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	if m.common.GetViewName() == common.PRoot || m.common.GetViewName() == common.PRessources {
		newListModel, cmdList := m.currentList.Update(msg)
		m.currentList = newListModel
		cmds = append(cmds, cmdList)
	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	if m.errorRaised {
		m.errorPanel, cmd = m.errorPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.common.GetViewName() == common.PViewPort ||
		m.common.GetViewName() == common.PPrintActions {
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.k8s, cmd = m.k8s.Update(msg)
	cmds = append(cmds, cmd)

	if m.common.GetViewName() == common.PDialogBox {
		m.dialogbox, cmd = m.dialogbox.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.common, cmd = m.common.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns the string representation of the model.
func (m model) View() string {
	doc := strings.Builder{}
	header := strings.Builder{}
	footer := strings.Builder{}
	center := strings.Builder{}

	// header
	{
		header.WriteString(m.header.View())
	}

	// footer
	{
		footer.WriteString(m.footer.View())
	}

	if m.errorRaised {
		center.WriteString(m.errorPanel.View())
	} else {
		switch m.common.GetViewName() {
		case common.PDialogBox:
			center.WriteString(m.dialogbox.View())

		case common.PViewPort:
			center.WriteString(m.markdown.Viewport.View())

		case common.PK8S:
			cmd := m.k8sCmdList[m.k8sCurrentIDView]
			getOutput := "loading..."
			if cmd.Done {
				getOutput = cmd.Result
			}
			reloadOutput := ""

			h := "Using ressource : " + cmd.Kind
			if !m.keys.Apply.Enabled() {
				h = "kubectl get managed"
			}

			h = lipgloss.NewStyle().Background(tui.NotificationColour).Padding(0, 2, 0, 2).Margin(0, 0, 1, 0).Render(h)
			hHeight := lipgloss.Height(h)
			reloadHeight := lipgloss.Height(reloadOutput)

			boxHeight := m.centerHeight - hHeight - reloadHeight - 4
			reloadOutput = fmt.Sprintf("%s reloading... %s", m.progressK8SGet.View(), m.k8sProgressMsg)
			reloadOutput = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Margin(1, 0, 1, 0).Render(reloadOutput)
			cmdResult := lipgloss.NewStyle().Width(m.width - 2).MaxWidth(m.width).MaxHeight(boxHeight - 2).Padding(1).Render(getOutput)
			getOutput = lipgloss.NewStyle().Width(m.width - 2).Height(boxHeight).Border(lipgloss.RoundedBorder()).Render(cmdResult)

			ui := lipgloss.JoinVertical(lipgloss.Center, h, getOutput, reloadOutput)
			dialog := lipgloss.Place(m.width, m.centerHeight,
				lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Render(ui),
			)

			center.WriteString(dialog)

		case common.PRoot, common.PRessources:
			center.WriteString(lipgloss.NewStyle().Render(m.currentList.View()))

		case common.PPrintActions:
			selected := m.currentList.SelectedItem().(*tui.Example)
			selectedFile := selected.Title()
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

			if m.showDependenciesFiles && selected.HaveDependenciesFiles() {
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
			m.markdown.Viewport.SetContent(ui)

			center.WriteString(m.markdown.Viewport.View())
		}
	}

	// Render the document
	doc.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		header.String(),
		center.String(),
		footer.String(),
	))

	// if m.width > 0 {
	// 	// physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	// 	tui.AppStyle = tui.AppStyle.MaxWidth(m.width).MaxHeight(m.height)
	// }

	// Okay, let's print it
	return tui.AppStyle.Render(doc.String())
}

func (m model) showExamples() (model, tea.Cmd) {
	i := m.exampleList["-"]
	cmd := m.currentList.SetItems(i)
	m.currentList.Title = "Choose an example"
	m.listName = "-"

	m.currentList.Select(0)

	return m, cmd
}

func (m model) showYaml(title string) (model, tea.Cmd) {
	var cmd tea.Cmd
	if _, ok := m.exampleList[title]; ok {
		i := m.exampleList[title]
		cmd = m.currentList.SetItems(i)
		m.listName = title
	}
	return m, cmd
}

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		if !m.tickRunning {
			return nil
		}
		return tickK8SGet(t)
	})
}

func (m model) rootView() (model, tea.Cmd) {
	var cmds []tea.Cmd
	m.common.SetViewName(common.PRoot)
	m.keys.EnableRootKeys()
	cmd := m.currentList.NewStatusMessage("back to home")
	cmds = append(cmds, cmd)

	if m.currentList.FilterState() == list.FilterApplied {
		m.currentList.ResetFilter()
	}

	m, cmd = m.showExamples()
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) generateK8SFiles() (model, *k8s.Cmd, tea.Cmd) {
	if m.currentList.SelectedItem() == nil {
		cmd := m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(
			"no item selected, empty list ?",
			errors.New("m.currentList.SelectedItem() == nil"),
		)
		m.errorRaised = true
		m.header.NotificationOK = tui.ErrorMark
		return m, nil, cmd
	}

	selectedItem := m.currentList.SelectedItem().(*tui.Example)

	file := selectedItem.Title()
	extra := selectedItem.HaveExtraFile()
	secret := selectedItem.HaveSecretFile()
	files := []string{file}

	if extra {
		f := fmt.Sprintf("%s.extra", file)
		files = append(files, f)
	}

	if secret {
		f := fmt.Sprintf("%s.secret", file)
		files = append(files, f)
	}

	if m.showDependenciesFiles && selectedItem.HaveDependenciesFiles() {
		d := selectedItem.DependenciesFilesList()
		files = append(files, d...)
	}

	cmd := &k8s.Cmd{
		ID:    randSeq(5),
		Done:  false,
		Files: files,
		Kind:  selectedItem.Description(),
	}

	return m, cmd, nil
}
