// Package tui provides the keymap for the application
package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

// ListKeyMap is the keymap for the application
type ListKeyMap struct {
	help.KeyMap
	VpKM                  viewport.KeyMap
	Back                  key.Binding
	Enter                 key.Binding
	Quit                  key.Binding
	UpDown                key.Binding
	LeftRight             key.Binding
	ForwardSlash          key.Binding
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

// NewListKeyMap creates a new keymap
func NewListKeyMap() *ListKeyMap {
	return &ListKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "quit"),
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
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↲", "select"),
		),
		ForwardSlash: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
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
		VpKM: viewport.DefaultKeyMap(),
	}
}

// ShortHelp returns the short help
func (m ListKeyMap) ShortHelp() []key.Binding {
	kb := make([]key.Binding, 0)
	kb = append(kb, m.ActiveShortHelp...)

	return kb
}

// FullHelp returns the full help
func (m ListKeyMap) FullHelp() [][]key.Binding {
	return m.ActiveFullHelp
}

// RootHelp sets the keymap for the root view
func (m *ListKeyMap) RootHelp() {
	kb := []key.Binding{
		m.Get,
		m.Help,
		m.Quit,
	}

	m.ActiveShortHelp = kb

	m.ActiveFullHelp = [][]key.Binding{
		{m.UpDown, m.LeftRight, m.Back},
		{m.ForwardSlash, m.Enter},
		{m.Help, m.Quit},
		{m.ShowRessources, m.ShowTested, m.GenerateListTested},
		{m.Get},
	}
}

// YamlHelp sets the keymap for the yaml view
func (m *ListKeyMap) YamlHelp() {
	kb := []key.Binding{
		m.Apply,
		m.Delete,
		m.Get,
		m.Print,
		m.ShowDependanciesFiles,
		m.Help,
		m.Quit,
	}

	m.ActiveShortHelp = kb

	m.ActiveFullHelp = [][]key.Binding{
		{m.UpDown, m.LeftRight, m.Back},
		{m.ForwardSlash, m.Enter},
		{m.Help, m.Quit},
		{m.Apply, m.Delete, m.Print},
		{m.Get, m.ShowDependanciesFiles},
	}
}

// GetHelp sets the keymap for the get view
func (m *ListKeyMap) GetHelp() {
	m.ActiveShortHelp = []key.Binding{
		m.Apply,
		m.Delete,
		m.Back,
		m.Quit,
	}

	m.ActiveFullHelp = [][]key.Binding{}
}

// YamlActionHelp sets the keymap for the yaml action view
func (m *ListKeyMap) YamlActionHelp() {

	m.ActiveShortHelp = []key.Binding{
		m.Back,
		m.Quit,
	}

	m.ActiveFullHelp = [][]key.Binding{}
}

// ViewPortHelp sets the keymap for the viewport view
func (m *ListKeyMap) ViewPortHelp() {
	m.ActiveShortHelp = []key.Binding{
		m.VpKM.Up,
		m.VpKM.Down,
		m.VpKM.PageUp,
		m.VpKM.PageDown,
		m.Back,
		m.Quit,
	}
}

// ErrorHelp sets the keymap for the error view
func (m *ListKeyMap) ErrorHelp() {
	m.ActiveShortHelp = []key.Binding{
		m.Back,
		m.Quit,
	}
}

// OnlyBackQuit sets the keymap for the simple view
func (m *ListKeyMap) OnlyBackQuit() {
	m.ActiveShortHelp = []key.Binding{
		m.Back,
		m.Quit,
	}

	m.ActiveFullHelp = [][]key.Binding{}
}
