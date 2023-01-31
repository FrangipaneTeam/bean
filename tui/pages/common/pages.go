package common

import (
	"github.com/FrangipaneTeam/bean/internal/keymap"
	tea "github.com/charmbracelet/bubbletea"
)

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
	Keys         *keymap.ListKeyMap
	previousPage PageID
	oldIndex     int
}

func BeanPages() map[PageID]*Page {
	pages := make(map[PageID]*Page)
	rootKeys := keymap.NewListKeyMap()
	kindKeys := keymap.NewListKeyMap()
	viewportKeys := keymap.NewListKeyMap()
	k8sGetRootKeys := keymap.NewListKeyMap()
	k8sGetKeys := keymap.NewListKeyMap()
	dialogBoxKeys := keymap.NewListKeyMap()
	errorKeys := keymap.NewListKeyMap()

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
			m.ex.CurrentList.Select(m.pages[m.viewName].oldIndex)
			m.ex, cmd = m.ex.SetRootItems()
			m.ex.CurrentList.Title = "Choose an example"

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
		m.pages[m.viewName].oldIndex = m.ex.CurrentList.Index()
	}
	m.viewName = name
	*m.keys = *m.pages[name].Keys

	// if m.ex.CurrentList.FilterState() == list.FilterApplied {
	if m.viewName == PRessources {
		m.ex.CurrentList.Select(0)
		m.ex.CurrentList.ResetFilter()
		m.ex.CurrentList.Title = "Choose a kind"
	}
	// }
}

// GetView returns the view name.
func (m Model) GetViewName() PageID {
	return m.viewName
}

// SetPreviousViewName sets the previous view name.
func (m *Model) SetPreviousViewName(page, previous PageID) {
	m.pages[page].previousPage = previous
}
