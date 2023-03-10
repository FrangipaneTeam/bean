package exlist

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/internal/theme"
)

type Model struct {
	listName string
	keys     *keymap.ListKeyMap

	// list
	exampleList map[string][]list.Item
	CurrentList list.Model
}

func New(
	keymap *keymap.ListKeyMap,
	exampleList exlist.LoadedExamples,
	width int,
	height int,
) *Model {
	theme := theme.Default()
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(theme.List.SelectedTitleBorderColor).
		Foreground(theme.List.SelectedTitleColor).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderForeground(theme.List.SelectedDescBorderColor).
		Foreground(theme.List.SelectedDescColor)

	// delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.
	// 	Foreground(theme.List.SelectedDimDescColor)

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
