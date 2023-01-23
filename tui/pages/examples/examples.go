// Package examples provides the examples page.
package examples

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tools"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/footer"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/FrangipaneTeam/bean/tui/pages/md"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	pViewPort     = "viewport"
	pRoot         = "root"
	pRessources   = "ressources"
	pPrintActions = "printActions"
	pK8S          = "k8s"
	pDialogBox    = "dialogbox"
)

var (
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	k8sCmdList map[string]*k8s.Cmd
)

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type k8sCmd struct {
	ID        string
	Done      bool
	canceled  bool
	confirmed bool
	verb      string
	cmdOutput string
	Files     []string
	Kind      string
}

type model struct {
	exampleList map[string][]list.Item
	currentList list.Model

	listName      string
	width, height int
	keys          *tui.ListKeyMap
	oldKeys       *tui.ListKeyMap

	viewName    string
	oldViewName string

	showDependenciesFiles bool

	errorRaised bool

	header     header.Model
	footer     footer.Model
	errorPanel errorpanel.Model
	markdown   md.Model
	k8s        k8s.Model

	config config.Provider

	k8sCurrentIDView    string
	k8sProgressMsg      string
	k8sGetWithoutAction bool
	progressK8SGet      progress.Model
	tickRunning         bool

	previousItemPostion int

	listOldHeight int
	centerHeight  int

	dialogbox dialogbox.Model
}

type tickK8SGet time.Time

// New returns a new model of the examples page.
// nolint: golint // model not used outside of this package
func New(e tui.LoadedExamples, width, height int, c config.Provider) model {
	h, v := tui.AppStyle.GetFrameSize()
	// physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	// wP := physicalWidth - tui.AppStyle.GetHorizontalPadding()
	// hP := physicalHeight - tui.AppStyle.GetVerticalPadding()

	rootKeys := tui.NewListKeyMap()
	dialogKeys := tui.NewListKeyMap()
	// listKeys.EnableRootKeys()
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
		height,
	)
	list.Title = "Choose an example"
	list.DisableQuitKeybindings()
	list.SetShowHelp(false)
	list.SetStatusBarItemName("example", "examples")

	version := lipgloss.NewStyle().Foreground(tui.NotificationColour).Render("v" + c.Version)
	if c.NewVersion != "" {
		newVersion := lipgloss.NewStyle().
			Foreground(tui.NotificationColour).
			Render("v" + c.NewVersion)
		version = fmt.Sprintf("v%s (new version available: %s)", c.Version, newVersion)
	}

	header := header.New(
		"Bean "+version,
		"A FrangipaneTeam bin",
		width-h,
		c,
	)

	footer := footer.New(width-h, rootKeys)
	headerHeight := header.Height()
	footerHeight := footer.Height()

	list.SetSize(width-h, height-v-headerHeight-footerHeight)

	k8sCmdList = make(map[string]*k8s.Cmd)

	// default activated keys
	rootKeys.EnableRootKeys()
	dialogKeys.EnableDialogBoxKeys()

	return model{
		exampleList: e.Examples,
		currentList: list,
		listName:    "-",
		viewName:    pRoot,

		keys:       rootKeys,
		header:     header,
		footer:     footer,
		errorPanel: errorpanel.New(width-h, height-v-headerHeight-footerHeight),
		markdown:   md.New(width-h, height-v-headerHeight-footerHeight),
		dialogbox: dialogbox.New(
			width-h,
			height-v-headerHeight-footerHeight,
			dialogKeys,
		),
		k8s:          k8s.New(rootKeys),
		width:        width - h,
		height:       height - v,
		centerHeight: height - v - headerHeight - footerHeight,
		config:       c,

		showDependenciesFiles: true,
		progressK8SGet: progress.New(
			progress.WithSolidFill("#CBEDD5"),
			progress.WithoutPercentage(),
			progress.WithWidth(10),
		),
	}
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
			m.tickRunning = false

			// back to the list, cancel all k8sCmd
			for _, v := range k8sCmdList {
				v.Canceled = true
			}

			cmd = m.progressK8SGet.SetPercent(0)
			cmds = append(cmds, cmd)

			switch m.viewName {
			case pRessources:
				m, cmd = m.rootView()
				return m, cmd

			case pViewPort:
				m.viewName = pRoot
				m.keys.EnableRootKeys()
				return m, cmd

			case pPrintActions, pK8S, pDialogBox:
				if m.viewName == pK8S && !m.keys.Apply.Enabled() {
					m, cmd = m.rootView()
					return m, cmd
				}

				m.keys.EnableKindListKeys()
				m.viewName = pRessources
				cmd = m.currentList.NewStatusMessage("back to " + m.listName)
				// m.currentList.ResetSelected()
				cmds = append(cmds, cmd)
				m, cmd = m.showYaml(m.listName)
			}

			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		case key.Matches(msg, m.keys.Select):
			switch m.viewName {
			case pDialogBox:
				var newModel model
				var cmds []tea.Cmd
				var k8sCmd *k8s.Cmd

				if m.dialogbox.ActiveButton == 2 {
					newModel = m
					newModel.header.Notification = "cancel delete"
					newModel.header.NotificationOK = tui.ErrorMark
					newModel.viewName = m.oldViewName
					*newModel.keys = *m.oldKeys
				} else {
					m.header.NotificationOK = tui.RunningMark
					newModel, k8sCmd, cmd = m.generateK8SFiles()
					if cmd != nil {
						newModel.viewName = m.oldViewName
						*newModel.keys = *m.oldKeys
						return newModel, cmd
					}

					newModel.k8sProgressMsg = "delete sent !"
					k8sCmd.Verb = "delete"

					newModel.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
					newModel.header.NotificationOK = tui.RunningMark

					k8sCmdList[k8sCmd.ID] = k8sCmd
					cmd = tools.Kubectl(k8sCmd)

					newModel.viewName = m.oldViewName
					*newModel.keys = *m.oldKeys

					cmds = append(cmds, cmd)
				}

				if newModel.viewName == pK8S {
					cmds = append(cmds, m.tickCmd())
				}

				return newModel, tea.Batch(cmds...)

			case pRoot:
				title := m.currentList.SelectedItem().(*tui.Example).Title()
				if m.currentList.FilterState() == list.FilterApplied {
					m.currentList.ResetFilter()
				}

				m.keys.EnableKindListKeys()

				m, cmd = m.showYaml(title)
				m.viewName = pRessources

				m.currentList.Select(0)

				return m, cmd
			}

		case key.Matches(msg, m.keys.Print):
			m.keys.EnableViewPortKeys()
			m.keys.ShowDependanciesFiles.SetEnabled(true)
			m.viewName = pPrintActions
			return m, nil

		case key.Matches(msg, m.keys.Help):
			if m.viewName != pK8S {
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

			if m.viewName == pRoot {
				cmd = tools.RenderMarkdown(m.config.Path+"/list-resources.md", m.width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.ShowTested):
			m.keys.EnableViewPortKeys()
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

		case key.Matches(msg, m.keys.Get), key.Matches(msg, m.keys.Apply):
			// TODO: in witch case m.currentList.SelectedItem() == nil ?
			m, k8sCmd, cmd := m.generateK8SFiles()
			if cmd != nil {
				return m, cmd
			}

			switch {
			case key.Matches(msg, m.keys.Get):
				m.viewName = pK8S
				m.k8sCurrentIDView = k8sCmd.ID
				m.keys.Get.SetEnabled(false)
				m.keys.Print.SetEnabled(false)
				m.keys.ShowDependanciesFiles.SetEnabled(false)
				m.keys.Select.SetEnabled(false)
				m.keys.Back.SetEnabled(true)
				m.keys.Help.SetEnabled(false)
				k8sCmd.Verb = "managed"

			case key.Matches(msg, m.keys.Apply):
				m.k8sProgressMsg = "apply sent !"
				k8sCmd.Verb = "apply"

			case key.Matches(msg, m.keys.Delete):
				m.k8sProgressMsg = "delete sent !"
				k8sCmd.Verb = "delete"
			}

			m.header.Notification = fmt.Sprintf("k %s @ %s", k8sCmd.Verb, time.Now().Format("15:04:05"))
			m.header.NotificationOK = tui.RunningMark

			k8sCmdList[k8sCmd.ID] = k8sCmd
			cmd = tools.Kubectl(k8sCmd)
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
		okValue := "No Fear !"
		cancelValue := "I'm scared !"
		m.dialogbox.SetDialogBox(question, okValue, cancelValue)
		m.oldViewName = m.viewName
		m.oldKeys = &tui.ListKeyMap{}
		*m.oldKeys = *m.keys
		m.viewName = pDialogBox
		m.keys.EnableDialogBoxKeys()
		return m, nil

	case tui.LoadedExamples:
		m.header.Notification = fmt.Sprintf("loaded new examples @ %s", time.Now().Format("15:04:05"))
		m.header.NotificationOK = tui.CheckMark
		m.exampleList = msg.Examples
		m.currentList.SetItems(msg.Examples[m.listName])

	case tools.Markdown:
		m.viewName = pViewPort
		m.markdown.Viewport.SetContent(msg.Content)
		m.markdown.Viewport.GotoTop()
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		return m, cmd

	case tools.ErrorMsg:
		cmd = m.errorPanel.Init()
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
		m.header.NotificationOK = tui.ErrorMark
		return m, cmd

	case tui.ListTestedDone:
		cmd = m.currentList.NewStatusMessage("List tested generated")
		return m, cmd

	case *k8s.Cmd:
		if msg.Canceled {
			delete(k8sCmdList, msg.ID)
			return m, nil
		}

		switch msg.Verb {
		case "apply", "delete":
			m.k8sProgressMsg = ""
			m.header.Notification = fmt.Sprintf("k %s @ %s", msg.Verb, time.Now().Format("15:04:05"))
			m.header.NotificationOK = tui.CheckMark
			msg.Done = true
			return m, cmd

		case "managed":
			m.viewName = pK8S
			m.k8sCurrentIDView = msg.ID
			m.k8sProgressMsg = ""
			m.header.Notification = fmt.Sprintf("k %s @ %s", msg.Verb, time.Now().Format("15:04:05"))
			m.header.NotificationOK = tui.CheckMark
			msg.Done = true
			if !m.tickRunning {
				m.tickRunning = true
				cmd = m.tickCmd()
			}
			return m, cmd
		}

	case tickK8SGet:
		if m.viewName == pDialogBox {
			return m, nil
		}
		if m.tickRunning {
			var kCmd tea.Cmd
			if m.progressK8SGet.Percent() == 1.0 {
				m.progressK8SGet.SetPercent(0)

				kCmd = tools.Kubectl(k8sCmdList[m.k8sCurrentIDView])
			}

			// Note that you can also use progress.Model.SetPercent to set the
			// percentage value explicitly, too.
			cmd := m.progressK8SGet.IncrPercent(0.1)
			return m, tea.Batch(m.tickCmd(), cmd, kCmd)
		}

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progressK8SGet.Update(msg)
		m.progressK8SGet = progressModel.(progress.Model)
		return m, cmd
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

	if m.viewName == pViewPort || m.viewName == pPrintActions {
		m.markdown.Viewport, cmd = m.markdown.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.k8s, cmd = m.k8s.Update(msg)
	cmds = append(cmds, cmd)

	if m.viewName == pDialogBox {
		m.dialogbox, cmd = m.dialogbox.Update(msg)
		cmds = append(cmds, cmd)
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

	if m.errorRaised {
		center.WriteString(m.errorPanel.View())
	} else {
		switch m.viewName {
		case pDialogBox:
			center.WriteString(m.dialogbox.View())

		case pViewPort:
			center.WriteString(m.markdown.Viewport.View())

		case pK8S:
			cmd := k8sCmdList[m.k8sCurrentIDView]
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

		case pRoot, pRessources:
			center.WriteString(lipgloss.NewStyle().Render(m.currentList.View()))

		case pPrintActions:
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
	m.viewName = pRoot
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
