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
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	k8sDelete  = "delete"
	k8sApply   = "apply"
	k8sManaged = "managed"

	k8sProgressIncrement = 0.1
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
		if m.pages.CurrentList.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.Back):
			m.common, cmds = m.common.Back()
			return m, tea.Batch(cmds...)

		case key.Matches(msg, m.keys.Select):
			switch view := m.pages.GetViewName(); view {
			case pages.PDialogBox:
				var k8sCmd *k8s.Cmd

				if m.dialogbox.ActiveButton == dialogbox.GetCancelValue() {
					m.header.Notification = "cancel delete"
					m.header.NotificationOK = tui.ErrorMark
					m.pages.RestorePreviousKeys()
					m.pages.RestorePreviousView()
				} else {
					m.header.NotificationOK = tui.RunningMark
					m, k8sCmd, cmd = m.generateK8SFiles()
					if cmd != nil {
						m.pages.RestorePreviousKeys()
						m.pages.RestorePreviousView()
						return m, cmd
					}

					m.k8sProgressMsg = "delete sent !"
					k8sCmd.Verb = k8sDelete

					m.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
					m.header.NotificationOK = tui.RunningMark

					ctx, cancel := context.WithCancel(context.Background())
					k8sCmd.Cancel = cancel
					m.common.AddContextToStop(cancel)
					m.k8s.CmdList[k8sCmd.ID] = k8sCmd
					cmd = k8s.Kubectl(ctx, k8sCmd)

					m.pages.RestorePreviousKeys()
					m.pages.RestorePreviousView()

					cmds = append(cmds, cmd)
				}

				if m.pages.GetViewName() == pages.PK8SGet || m.pages.GetViewName() == pages.PK8SGetFromRoot {
					cmds = append(cmds, m.tickCmd())
				}

				return m, tea.Batch(cmds...)

			case pages.PRoot:
				title := m.pages.CurrentList.SelectedItem().(*tui.Example).Title()

				m.pages, cmd = m.pages.UpdateList(title)
				m.pages.SetViewName(pages.PRessources)

				return m, cmd
			}

		case key.Matches(msg, m.keys.Help):
			if m.pages.GetViewName() != pages.PK8SGet && m.pages.GetViewName() != pages.PK8SGetFromRoot {
				m.footer.Help.ShowAll = !m.footer.Help.ShowAll
				m.footer.Help.Width = m.width

				_, v := tui.AppStyle.GetFrameSize()
				listHeight := m.pages.CurrentList.Height()
				helpHeight := lipgloss.Height(m.footer.Help.View(m.keys))

				if m.footer.Help.ShowAll {
					m.listOldHeight = listHeight
					m.centerHeight = listHeight - helpHeight + 1
				} else {
					m.centerHeight = m.listOldHeight
				}
				m.pages.CurrentList.SetHeight(m.centerHeight)
				m.errorPanel.SetSize(m.width-v, m.centerHeight)
				m.markdown.SetSize(m.width, m.centerHeight)
			}
			return m, cmd

		case key.Matches(msg, m.keys.ShowRessources):
			m.keys.EnableViewPortKeys()

			if m.pages.GetViewName() == pages.PRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-resources.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.ShowTested):
			m.keys.EnableViewPortKeys()
			if m.pages.GetViewName() == pages.PRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-tested.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.GenerateListTested):
			cmd = pages.GenerateListTested(m.config)
			return m, cmd

		case key.Matches(msg, m.keys.Get), key.Matches(msg, m.keys.Apply):
			var k8sCmd *k8s.Cmd
			m, k8sCmd, cmd = m.generateK8SFiles()
			if cmd != nil {
				return m, cmd
			}
			k8sCmd.FromPage = m.pages.GetViewName()

			switch {
			case key.Matches(msg, m.keys.Get):
				switch m.pages.GetViewName() {
				case pages.PRoot:
					m.pages.SetViewName(pages.PK8SGetFromRoot)
				case pages.PRessources:
					m.pages.SetViewName(pages.PK8SGet)
				default:
					return m, nil
				}
				m.k8sCurrentIDView = k8sCmd.ID
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
			m.k8s.CmdList[k8sCmd.ID] = k8sCmd
			// m.header.RunningCommands++
			cmd = k8s.Kubectl(ctx, k8sCmd)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		headerHeight := m.header.Height()
		footerHeight := m.footer.Height()

		h, v := tui.AppStyle.GetFrameSize()

		m.width, m.height = msg.Width-h, msg.Height-v
		centerH := m.height - headerHeight - footerHeight
		m.centerHeight = centerH

		m.header.SetWidth(m.width)
		m.footer.SetWidth(m.width)

		m.pages.CurrentList.SetSize(m.width, centerH)
		m.dialogbox.SetSize(m.width, centerH)
		m.errorPanel.SetSize(m.width, centerH)
		m.markdown.SetSize(m.width, centerH)
		return m, nil

	case k8s.Message:
		question := "Delete all ressources ?"
		if m.pages.ShowDependenciesFiles {
			question = "Delete all ressources WITH dependencies ?"
		}
		okValue := "No Fear !"
		cancelValue := "I'm scared !"
		m.dialogbox.SetDialogBox(question, okValue, cancelValue)
		m.pages.SetPreviousViewName(pages.PDialogBox, msg.PreviousPage)
		m.pages.SetViewName(pages.PDialogBox)
		return m, nil

	case tui.LoadedExamples:
		m.header.Notification = fmt.Sprintf("loaded new examples @ %s", time.Now().Format("15:04:05"))
		m.header.NotificationOK = tui.CheckMark
		m.pages.UpdateExamplesList(msg.Examples)
		m.pages, cmd = m.pages.UpdateList()
		return m, cmd

	case tools.Markdown:
		m.pages.SetViewName(pages.PViewPort)
		m.markdown.Viewport.SetContent(msg.Content)
		m.markdown.Viewport.GotoTop()
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		return m, cmd

	case errorpanel.ErrorMsg:
		cmd = m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.header.NotificationOK = tui.ErrorMark
		m.pages.SetPreviousViewName(pages.PError, msg.FromPage.(pages.PageID))
		m.pages.SetViewName(pages.PError)
		if m.config.Debug {
			m.header.Notification = fmt.Sprintf("from %s", msg.FromPage.(pages.PageID))
		}
		return m, cmd

	case tui.ListTestedDone:
		cmd = m.pages.CurrentList.NewStatusMessage("List tested generated")
		return m, cmd

	case *k8s.Cmd:
		// m.header.RunningCommands--
		delete(m.k8s.CmdList, msg.ID)
		m.header.NotificationOK = tui.CheckMark
		m.k8sProgressMsg = ""
		m.header.Notification = fmt.Sprintf("k %s @ %s", msg.Verb, time.Now().Format("15:04:05"))

		if msg.Verb == k8sManaged {
			m.k8sCurrentIDView = msg.ID
			if !m.k8s.IsTickRunning() {
				m.k8s.SetTickRunning(true)
				cmd = m.tickCmd()
			}
		}
		return m, cmd

	case tickK8SGet:
		if m.pages.GetViewName() == pages.PDialogBox {
			return m, nil
		}
		if m.k8s.IsTickRunning() {
			var kubectlCmd tea.Cmd
			if m.k8s.GetProgress.Percent() == 1.0 {
				m.k8s.GetProgress.SetPercent(0)

				ctx, cancel := context.WithCancel(context.Background())
				m.k8s.CmdList[m.k8sCurrentIDView].Cancel = cancel
				m.common.AddContextToStop(cancel)
				// m.header.RunningCommands++
				kubectlCmd = k8s.Kubectl(ctx, m.k8s.CmdList[m.k8sCurrentIDView])
			}

			// Note that you can also use progress.Model.SetPercent to set the
			// percentage value explicitly, too.
			cmd = m.k8s.GetProgress.IncrPercent(k8sProgressIncrement)
			return m, tea.Batch(m.tickCmd(), cmd, kubectlCmd)
		}

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var progressModel tea.Model
		progressModel, cmd = m.k8s.GetProgress.Update(msg)
		m.k8s.GetProgress = progressModel.(progress.Model)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	if m.pages.GetViewName() == pages.PRoot || m.pages.GetViewName() == pages.PRessources {
		newListModel, cmdList := m.pages.CurrentList.Update(msg)
		m.pages.CurrentList = newListModel
		cmds = append(cmds, cmdList)
	}

	m.k8s, cmd = m.k8s.Update(msg)
	cmds = append(cmds, cmd)

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	if m.errorPanel.ErrorRaised() {
		m.errorPanel, cmd = m.errorPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.pages.GetViewName() == pages.PViewPort ||
		m.pages.GetViewName() == pages.PPrintActions {
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.pages.GetViewName() == pages.PDialogBox {
		m.dialogbox, cmd = m.dialogbox.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.common, cmd = m.common.Update(msg)
	cmds = append(cmds, cmd)

	// debug
	if m.config.Debug {
		m.header.Description = fmt.Sprintf("view: %s", m.pages.GetViewName())
	}

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

	if m.errorPanel.ErrorRaised() {
		center.WriteString(m.errorPanel.View())
	} else {
		switch m.pages.GetViewName() {
		case pages.PDialogBox:
			center.WriteString(m.dialogbox.View())

		case pages.PViewPort:
			center.WriteString(m.markdown.Viewport.View())

		case pages.PK8SGet, pages.PK8SGetFromRoot:
			cmd := m.k8s.CmdList[m.k8sCurrentIDView]
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
			reloadOutput = fmt.Sprintf("%s reloading... %s", m.k8s.GetProgress.View(), m.k8sProgressMsg)
			reloadOutput = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Margin(1, 0, 1, 0).Render(reloadOutput)
			cmdResult := lipgloss.NewStyle().Width(m.width - 2).MaxWidth(m.width).MaxHeight(boxHeight - 2).Padding(1).Render(getOutput)
			getOutput = lipgloss.NewStyle().Width(m.width - 2).Height(boxHeight).Border(lipgloss.RoundedBorder()).Render(cmdResult)

			ui := lipgloss.JoinVertical(lipgloss.Center, h, getOutput, reloadOutput)
			dialog := lipgloss.Place(m.width, m.centerHeight,
				lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Render(ui),
			)

			center.WriteString(dialog)

		case pages.PRoot, pages.PRessources:
			center.WriteString(lipgloss.NewStyle().Render(m.pages.CurrentList.View()))

		case pages.PPrintActions:
			ui := m.k8s.View()
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

// func (m model) showExamples() (model, tea.Cmd) {
// 	i := m.exampleList["-"]
// 	cmd := m.currentList.SetItems(i)
// 	m.currentList.Title = "Choose an example"
// 	m.listName = "-"

// 	m.currentList.Select(0)

// 	return m, cmd
// }

// func (m model) showYaml(title string) (model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	if _, ok := m.exampleList[title]; ok {
// 		i := m.exampleList[title]
// 		cmd = m.currentList.SetItems(i)
// 		m.listName = title
// 	}
// 	return m, cmd
// }

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		if !m.k8s.IsTickRunning() {
			return nil
		}
		return tickK8SGet(t)
	})
}

// func (m model) rootView() (model, tea.Cmd) {
// 	var cmds []tea.Cmd
// 	m.pages.SetViewName(pages.PRoot)
// 	m.keys.EnableRootKeys()
// 	cmd := m.currentList.NewStatusMessage("back to home")
// 	cmds = append(cmds, cmd)

// 	if m.currentList.FilterState() == list.FilterApplied {
// 		m.currentList.ResetFilter()
// 	}

// 	m, cmd = m.showExamples()
// 	cmds = append(cmds, cmd)

// 	return m, tea.Batch(cmds...)
// }

func (m model) generateK8SFiles() (model, *k8s.Cmd, tea.Cmd) {
	if m.pages.CurrentList.SelectedItem() == nil {
		cmd := m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(
			"no item selected, empty list ?",
			errors.New("m.currentList.SelectedItem() == nil"),
		)
		m.header.NotificationOK = tui.ErrorMark
		return m, nil, cmd
	}

	selectedItem := m.pages.CurrentList.SelectedItem().(*tui.Example)

	file := selectedItem.FileWithPath()
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

	if m.pages.ShowDependenciesFiles && selectedItem.HaveDependenciesFiles() {
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
