package elist

import (
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	listName string
	keys     *tui.ListKeyMap

	// list
	exampleList map[string][]list.Item
	CurrentList list.Model
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
		keys:        keymap,
		exampleList: exampleList.Examples,
		CurrentList: list,
	}
}
