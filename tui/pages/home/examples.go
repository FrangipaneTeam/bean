// Package home provides the home page.
package home

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/theme"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
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
		// tea.EnterAltScreen,
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
		case key.Matches(msg, m.keys.Select):
			switch view := m.common.GetViewName(); view {
			case common.PDialogBox:
				var k8sCmd *k8s.Cmd

				if m.dialogbox.ActiveButton == dialogbox.GetCancelValue() {
					m.header.Notification = "cancel delete"
					m.header.NotificationOK = theme.ErrorMark
					m.common.RestorePreviousKeys()
					m.common.RestorePreviousView()
				} else {
					m.header.NotificationOK = theme.RunningMark
					m, k8sCmd, cmd = m.generateK8SFiles()
					if cmd != nil {
						m.common.RestorePreviousKeys()
						m.common.RestorePreviousView()
						return m, cmd
					}

					m.k8sProgressMsg = "delete sent !"
					k8sCmd.Verb = k8sDelete

					m.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
					m.header.NotificationOK = theme.RunningMark

					ctx, cancel := context.WithCancel(context.Background())
					k8sCmd.Cancel = cancel
					m.common.AddContextToStop(cancel)
					m.k8s.CmdList[k8sCmd.ID] = k8sCmd
					cmd = k8s.Kubectl(ctx, k8sCmd)

					m.common.RestorePreviousKeys()
					m.common.RestorePreviousView()

					cmds = append(cmds, cmd)
				}

				if m.common.GetViewName() == common.PK8SGet || m.common.GetViewName() == common.PK8SGetFromRoot {
					cmds = append(cmds, m.tickCmd())
				}

				return m, tea.Batch(cmds...)

			case common.PRoot:
				title := m.pages.CurrentList.SelectedItem().(*exlist.Example).Title()

				m.pages, cmd = m.pages.UpdateList(title)
				m.common.SetViewName(common.PRessources)

				return m, cmd
			}

		case key.Matches(msg, m.keys.Get), key.Matches(msg, m.keys.Apply):
			var k8sCmd *k8s.Cmd
			m, k8sCmd, cmd = m.generateK8SFiles()
			if cmd != nil {
				return m, cmd
			}
			k8sCmd.FromPage = m.common.GetViewName()

			switch {
			case key.Matches(msg, m.keys.Get):
				switch m.common.GetViewName() {
				case common.PRoot:
					m.common.SetViewName(common.PK8SGetFromRoot)
				case common.PRessources:
					m.common.SetViewName(common.PK8SGet)
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
			m.header.NotificationOK = theme.RunningMark

			ctx, cancel := context.WithCancel(context.Background())
			k8sCmd.Cancel = cancel
			m.common.AddContextToStop(cancel)
			m.k8s.CmdList[k8sCmd.ID] = k8sCmd
			common.RunningCommands = len(m.k8s.CmdList)
			// m.header.RunningCommands++
			cmd = k8s.Kubectl(ctx, k8sCmd)
			return m, cmd
		}

	case k8s.Message:
		question := "Delete all ressources ?"
		if m.k8s.ShowDependenciesFiles {
			question = "Delete all ressources WITH dependencies ?"
		}
		okValue := "No Fear !"
		cancelValue := "I'm scared !"
		m.dialogbox.SetDialogBox(question, okValue, cancelValue)
		m.common.SetPreviousViewName(common.PDialogBox, msg.PreviousPage.(common.PageID))
		m.common.SetViewName(common.PDialogBox)
		return m, nil

	case exlist.LoadedExamples:
		m.header.Notification = fmt.Sprintf("loaded new examples @ %s", time.Now().Format("15:04:05"))
		m.header.NotificationOK = theme.CheckMark
		m.pages.UpdateExamplesList(msg.Examples)
		m.pages, cmd = m.pages.UpdateList()
		return m, cmd

	case errorpanel.ErrorMsg:
		cmd = m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.header.NotificationOK = theme.ErrorMark
		m.common.SetPreviousViewName(common.PError, msg.FromPage.(common.PageID))
		m.common.SetViewName(common.PError)
		if m.config.Debug {
			m.header.Notification = fmt.Sprintf("from %s", msg.FromPage.(common.PageID))
		}
		return m, cmd

	case exlist.ListTestedDone:
		cmd = m.pages.CurrentList.NewStatusMessage("List tested generated")
		return m, cmd

	case *k8s.Cmd:
		// m.header.RunningCommands--
		delete(m.k8s.CmdList, msg.ID)
		common.RunningCommands = len(m.k8s.CmdList)
		m.header.NotificationOK = theme.CheckMark
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
		if m.common.GetViewName() == common.PDialogBox {
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

	case tea.WindowSizeMsg:
		headerHeight := m.header.Height()
		footerHeight := m.footer.Height()

		h, v := theme.AppStyle.GetFrameSize()

		m.width, m.height = msg.Width-h, msg.Height-v
		centerH := m.height - headerHeight - footerHeight
		m.centerHeight = centerH

		m.header.SetWidth(m.width)
		m.footer.SetWidth(m.width)

		m.pages.CurrentList.SetSize(m.width, centerH)
		m.dialogbox.SetSize(m.width, centerH)

		common.Height = m.height
		common.Width = m.width
		common.CenterHeight = centerH

		cmd = func() tea.Msg {
			msg := common.ResizeMsg{}
			return msg
		}
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	if m.common.GetViewName() == common.PRoot || m.common.GetViewName() == common.PRessources {
		newListModel, cmdList := m.pages.CurrentList.Update(msg)
		m.pages.CurrentList = newListModel
		cmds = append(cmds, cmdList)
	}

	if m.pages.CurrentList.FilterState() != list.Filtering {
		m.footer, cmd = m.footer.Update(msg)
		cmds = append(cmds, cmd)

		m.errorPanel, cmd = m.errorPanel.Update(msg)
		cmds = append(cmds, cmd)

		m.markdown, cmd = m.markdown.Update(msg)
		cmds = append(cmds, cmd)

		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)

		m.k8s, cmd = m.k8s.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.common.GetViewName() == common.PDialogBox {
		m.dialogbox, cmd = m.dialogbox.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.common, cmd = m.common.Update(msg)
	cmds = append(cmds, cmd)

	// debug
	if m.config.Debug {
		m.header.Description = fmt.Sprintf("view: %s", m.common.GetViewName())
	}

	return m, tea.Batch(cmds...)
}

// View returns the string representation of the model.
func (m model) View() string {
	doc := strings.Builder{}
	header := strings.Builder{}
	footer := strings.Builder{}
	center := strings.Builder{}
	width := common.Width

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
		switch m.common.GetViewName() {
		case common.PDialogBox:
			center.WriteString(m.dialogbox.View())

		case common.PViewPort:
			center.WriteString(m.markdown.Viewport.View())

		case common.PK8SGet, common.PK8SGetFromRoot:
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

			h = lipgloss.NewStyle().Background(theme.NotificationColour).Padding(0, 2, 0, 2).Margin(0, 0, 1, 0).Render(h)
			hHeight := lipgloss.Height(h)
			reloadHeight := lipgloss.Height(reloadOutput)

			boxHeight := common.CenterHeight - hHeight - reloadHeight - 4
			reloadOutput = fmt.Sprintf("%s reloading... %s", m.k8s.GetProgress.View(), m.k8sProgressMsg)
			reloadOutput = lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Margin(1, 0, 1, 0).Render(reloadOutput)
			cmdResult := lipgloss.NewStyle().Width(width - 2).MaxWidth(width).MaxHeight(boxHeight - 2).Padding(1).Render(getOutput)
			getOutput = lipgloss.NewStyle().Width(width - 2).Height(boxHeight).Border(lipgloss.RoundedBorder()).Render(cmdResult)

			ui := lipgloss.JoinVertical(lipgloss.Center, h, getOutput, reloadOutput)
			dialog := lipgloss.Place(width, common.CenterHeight,
				lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Render(ui),
			)

			center.WriteString(dialog)

		case common.PRoot, common.PRessources:
			center.WriteString(lipgloss.NewStyle().Render(m.pages.CurrentList.View()))

		case common.PPrintActions:
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

	return theme.AppStyle.Render(doc.String())
}

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		if !m.k8s.IsTickRunning() {
			return nil
		}
		return tickK8SGet(t)
	})
}

func (m model) generateK8SFiles() (model, *k8s.Cmd, tea.Cmd) {
	if m.pages.CurrentList.SelectedItem() == nil {
		cmd := m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(
			"no item selected, empty list ?",
			errors.New("m.currentList.SelectedItem() == nil"),
		)
		m.header.NotificationOK = theme.ErrorMark
		return m, nil, cmd
	}

	selectedItem := m.pages.CurrentList.SelectedItem().(*exlist.Example)

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

	if m.k8s.ShowDependenciesFiles && selectedItem.HaveDependenciesFiles() {
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
