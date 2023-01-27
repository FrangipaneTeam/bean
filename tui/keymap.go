// Package tui provides the keymap for the application
package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
)

// ListKeyMap is the keymap for the application.
type ListKeyMap struct {
	help.KeyMap
	ListKeyMap            list.KeyMap
	VpKM                  viewport.KeyMap
	Back                  key.Binding
	Select                key.Binding
	Quit                  key.Binding
	UpDown                key.Binding
	LeftRight             key.Binding
	Left                  key.Binding
	Right                 key.Binding
	Apply                 key.Binding
	Delete                key.Binding
	Print                 key.Binding
	Get                   key.Binding
	Help                  key.Binding
	ShowRessources        key.Binding
	ShowTested            key.Binding
	ShowDependanciesFiles key.Binding
	GenerateListTested    key.Binding
	ActiveShortHelp       []key.Binding
	ActiveFullHelp        [][]key.Binding
}

// NewListKeyMap creates a new keymap.
func NewListKeyMap() *ListKeyMap {
	return &ListKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		UpDown: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑↓", "up/down"),
		),
		LeftRight: key.NewBinding(
			key.WithKeys("left", "right"),
			key.WithHelp("← →", "prev/next page"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("⇐", "back"),
		),
		Apply: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "apply"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Print: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "print"),
		),
		Get: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "get"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↲", "select"),
		),
		ShowRessources: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "show ressources"),
		),
		ShowTested: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "show tested"),
		),
		ShowDependanciesFiles: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "dependencies"),
		),
		GenerateListTested: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "generate list tested"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		VpKM:       viewport.DefaultKeyMap(),
		ListKeyMap: list.DefaultKeyMap(),
	}
}

// ShortHelp returns the short help.
func (m ListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.VpKM.Up,
		m.VpKM.Down,
		m.VpKM.PageUp,
		m.VpKM.PageDown,
		m.Select,
		// m.ListKeyMap.Filter,
		m.Apply,
		m.Delete,
		m.Get,
		m.Print,
		// m.ShowDependanciesFiles,
		m.Help,
		m.Back,
		m.Quit,
	}
}

// FullHelp returns the full help.
func (m ListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.VpKM.Up, m.VpKM.Down, m.VpKM.PageUp, m.VpKM.PageDown},
		{m.UpDown, m.LeftRight, m.Back},
		{m.ListKeyMap.Filter, m.Select},
		{m.Help, m.Quit},
		{m.Apply, m.Delete, m.Print},
		{m.Get, m.ShowDependanciesFiles},
		{m.ShowRessources, m.ShowTested, m.GenerateListTested},
	}
}

// EnableViewPortKeys is the set of keys for the viewport.
func (m *ListKeyMap) EnableViewPortKeys() {
	m.disableK8SKeys()
	m.disableList()
	m.disableMD()
	m.Back.SetEnabled(true)
	m.Select.SetEnabled(false)

	m.VpKM.Up.SetEnabled(true)
	m.VpKM.Down.SetEnabled(true)
	m.VpKM.PageUp.SetEnabled(true)
	m.VpKM.PageDown.SetEnabled(true)
	m.VpKM.HalfPageUp.SetEnabled(true)
	m.VpKM.HalfPageDown.SetEnabled(true)
	m.Back.SetEnabled(true)
}

func (m *ListKeyMap) disableList() {
	m.UpDown.SetEnabled(false)
	m.LeftRight.SetEnabled(false)
	m.ListKeyMap.Filter.SetEnabled(false)
}

func (m *ListKeyMap) enableList() {
	m.UpDown.SetEnabled(true)
	m.LeftRight.SetEnabled(true)
	m.ListKeyMap.Filter.SetEnabled(true)
}

func (m *ListKeyMap) disableMD() {
	m.ShowRessources.SetEnabled(false)
	m.ShowTested.SetEnabled(false)
	m.GenerateListTested.SetEnabled(false)
}

func (m *ListKeyMap) enableMD() {
	m.ShowRessources.SetEnabled(true)
	m.ShowTested.SetEnabled(true)
	m.GenerateListTested.SetEnabled(true)
}

func (m *ListKeyMap) disableViewPortKeys() {
	m.VpKM.Up.SetEnabled(false)
	m.VpKM.Down.SetEnabled(false)
	m.VpKM.PageUp.SetEnabled(false)
	m.VpKM.PageDown.SetEnabled(false)
	m.VpKM.HalfPageUp.SetEnabled(false)
	m.VpKM.HalfPageDown.SetEnabled(false)
}

func (m *ListKeyMap) enableK8SKeys() {
	m.Apply.SetEnabled(true)
	m.Delete.SetEnabled(true)
	m.Print.SetEnabled(true)
	m.Get.SetEnabled(true)
	m.ShowDependanciesFiles.SetEnabled(true)
	m.disableMD()
	m.enableList()
}

func (m *ListKeyMap) disableK8SKeys() {
	m.Apply.SetEnabled(false)
	m.Delete.SetEnabled(false)
	m.Print.SetEnabled(false)
	m.Get.SetEnabled(false)
	m.ShowDependanciesFiles.SetEnabled(false)
}

// EnableRootKeys is the set of keys for the root.
func (m *ListKeyMap) EnableRootKeys() {
	m.disableK8SKeys()
	m.disableViewPortKeys()
	m.enableList()
	m.enableMD()
	m.Back.SetEnabled(false)
	m.Select.SetEnabled(true)
	m.Get.SetEnabled(true)
	m.Help.SetEnabled(true)
}

// EnableKindListKeys is the set of keys for the kind list.
func (m *ListKeyMap) EnableKindListKeys() {
	m.disableViewPortKeys()
	m.enableK8SKeys()
	m.Help.SetEnabled(true)
	m.ListKeyMap.Filter.SetEnabled(true)
	m.Back.SetEnabled(false)
	m.Select.SetEnabled(false)
	m.UpDown.SetEnabled(true)
	m.LeftRight.SetEnabled(true)
	m.Back.SetEnabled(true)
}

// EnablePrintK8SKeys is the set of keys for the k8s print view.
func (m *ListKeyMap) EnablePrintK8SKeys() {
	m.disableK8SKeys()
	m.disableViewPortKeys()
	m.disableList()
	m.ShowDependanciesFiles.SetEnabled(true)
	m.Back.SetEnabled(true)
	m.Select.SetEnabled(false)
}

// EnableDialogBoxKeys is the set of keys for the dialog box.
func (m *ListKeyMap) EnableDialogBoxKeys() {
	m.disableK8SKeys()
	m.disableViewPortKeys()
	m.Back.SetEnabled(false)
	m.ListKeyMap.Filter.SetEnabled(false)
	m.LeftRight.SetEnabled(true)
	m.Help.SetEnabled(false)
	m.Select.SetEnabled(true)
}

func (m *ListKeyMap) EnableGetRootKeys() {
	m.disableList()
	m.disableK8SKeys()
	m.disableViewPortKeys()
	// m.Get.SetEnabled(false)
	// m.Print.SetEnabled(false)
	// m.ShowDependanciesFiles.SetEnabled(false)
	m.Select.SetEnabled(false)
	// m.Back.SetEnabled(true)
	m.Help.SetEnabled(false)
}

func (m *ListKeyMap) EnableGetKeys() {
	m.enableList()
	m.disableK8SKeys()
	m.disableViewPortKeys()
	m.Apply.SetEnabled(true)
	m.Delete.SetEnabled(true)
	m.Select.SetEnabled(false)
	m.Help.SetEnabled(false)
}

func (m *ListKeyMap) EnableErrorKeys() {
	m.disableMD()
	m.disableList()
	m.disableK8SKeys()
	m.disableViewPortKeys()
	m.Back.SetEnabled(true)
	m.Select.SetEnabled(false)
	m.Help.SetEnabled(false)
}
