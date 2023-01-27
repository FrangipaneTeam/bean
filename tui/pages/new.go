package pages

import (
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewName PageID
	listName string
	keys     *tui.ListKeyMap
	pages    map[PageID]*Page

	// list
	exampleList           map[string][]list.Item
	CurrentList           list.Model
	ShowDependenciesFiles bool
}

//go:generate stringer -type=PageID
const (
	PActual PageID = iota
	PViewPort
	PRoot
	PRessources
	PPrintActions
	PDialogBox
	PK8SGet
	PK8SGetFromRoot
	PError
)

type PageID int

type Page struct {
	Keys         *tui.ListKeyMap
	previousPage PageID
	oldIndex     int
}

func BeanPages() map[PageID]*Page {
	pages := make(map[PageID]*Page)
	rootKeys := tui.NewListKeyMap()
	kindKeys := tui.NewListKeyMap()
	viewportKeys := tui.NewListKeyMap()
	k8sGetRootKeys := tui.NewListKeyMap()
	k8sGetKeys := tui.NewListKeyMap()
	dialogBoxKeys := tui.NewListKeyMap()
	errorKeys := tui.NewListKeyMap()

	rootKeys.EnableRootKeys()
	kindKeys.EnableKindListKeys()
	viewportKeys.EnableViewPortKeys()
	k8sGetRootKeys.EnableGetRootKeys()
	k8sGetKeys.EnableGetKeys()
	dialogBoxKeys.EnableDialogBoxKeys()
	errorKeys.EnableErrorKeys()

	root := &Page{
		Keys:         rootKeys,
		previousPage: PRoot,
	}

	kind := &Page{
		Keys:         kindKeys,
		previousPage: PRoot,
	}

	viewport := &Page{
		Keys:         viewportKeys,
		previousPage: PRoot,
	}

	k8sPrint := &Page{
		Keys:         viewportKeys,
		previousPage: PRessources,
	}

	k8sGetRoot := &Page{
		Keys:         k8sGetRootKeys,
		previousPage: PRoot,
	}

	k8sGet := &Page{
		Keys:         k8sGetKeys,
		previousPage: PRessources,
	}

	dialogBox := &Page{
		Keys:         dialogBoxKeys,
		previousPage: PActual,
	}

	errorP := &Page{
		Keys:         errorKeys,
		previousPage: PRoot,
	}

	pages[PRoot] = root
	pages[PRessources] = kind
	pages[PViewPort] = viewport
	pages[PPrintActions] = k8sPrint
	pages[PK8SGetFromRoot] = k8sGetRoot
	pages[PK8SGet] = k8sGet
	pages[PDialogBox] = dialogBox
	pages[PError] = errorP

	return pages
}

func New(keymap *tui.ListKeyMap, exampleList tui.LoadedExamples, width int, height int) *Model {
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

	list := list.New(exampleList.Examples["-"],
		delegate,
		width,
		height,
	)
	list.Title = "Choose an example"
	list.DisableQuitKeybindings()
	list.SetShowHelp(false)
	list.SetStatusBarItemName("example", "examples")

	// list.SetSize()

	return &Model{
		viewName:              PRoot,
		keys:                  keymap,
		pages:                 BeanPages(),
		exampleList:           exampleList.Examples,
		CurrentList:           list,
		ShowDependenciesFiles: true,
	}
}

// RestorePreviousKeys restore the keys.
func (m *Model) RestorePreviousKeys() tea.Cmd {
	if _, ok := m.pages[m.viewName]; ok {
		previousPage := m.pages[m.viewName].previousPage
		*m.keys = *m.pages[previousPage].Keys
	} else {
		*m.keys = *m.pages[PRoot].Keys
	}
	return nil
}

// RestorePreviousViewName restore the view name.
func (m *Model) RestorePreviousView() tea.Cmd {
	var cmd tea.Cmd
	if newPage, ok := m.pages[m.viewName]; ok {
		switch m.viewName {
		case PRessources:
			m.viewName = newPage.previousPage
			m.CurrentList.Select(m.pages[m.viewName].oldIndex)
			cmd = m.CurrentList.SetItems(m.exampleList["-"])
			m.CurrentList.Title = "Choose an example"

		// case PDialogBox:
		// 	if m.keys.Print.Enabled() {
		// 		m.viewName = newPage.PreviousPage()
		// 	} else {
		// 		m.viewName = PK8SGet
		// 	}

		default:
			m.viewName = newPage.previousPage
		}
	} else {
		m.viewName = PRoot
	}

	return cmd
}

// SetView sets the view name.
func (m *Model) SetViewName(name PageID) {
	if m.viewName == PRoot {
		m.pages[m.viewName].oldIndex = m.CurrentList.Index()
	}
	m.viewName = name
	*m.keys = *m.pages[name].Keys

	// if m.CurrentList.FilterState() == list.FilterApplied {
	if m.viewName == PRessources {
		m.CurrentList.Select(0)
		m.CurrentList.ResetFilter()
		m.CurrentList.Title = "Choose a kind"
	}
	// }
}

// GetView returns the view name.
func (m Model) GetViewName() PageID {
	return m.viewName
}

// UpdateExamplesList updates the examples list.
func (m *Model) UpdateExamplesList(examples map[string][]list.Item) {
	m.exampleList = examples
}

func (m *Model) UpdateList(params ...string) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	title := m.listName
	if len(params) != 0 {
		title = params[0]
	}

	if _, ok := m.exampleList[title]; ok {
		i := m.exampleList[title]
		cmd = m.CurrentList.SetItems(i)
		m.listName = title
	}
	return m, cmd
}

// GetDependenciesStatus returns the dependencies status.
func (m Model) GetDependenciesStatus() bool {
	return m.ShowDependenciesFiles
}

// SwitchDependenciesStatus switches the dependencies status.
func (m *Model) SwitchDependenciesStatus() {
	m.ShowDependenciesFiles = !m.ShowDependenciesFiles
}

// SetPreviousViewName sets the previous view name.
func (m *Model) SetPreviousViewName(page, previous PageID) {
	m.pages[page].previousPage = previous
}
