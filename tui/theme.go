package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColour        = lipgloss.Color("#3a1577")
	SecondaryColour      = lipgloss.Color("#5b1689")
	BorderColour         = lipgloss.Color("#CBEDD5")
	FeintColour          = lipgloss.Color("#807d8a")
	VeryFeintColour      = lipgloss.Color("#5e5e5e")
	TextColour           = lipgloss.Color("#f6f5fc")
	HighlightColour      = lipgloss.Color("#97DECE")
	HighlightFeintColour = lipgloss.Color("#CBEDD5")
	AmberColour          = lipgloss.Color("#e68a35")
	GreenColour          = lipgloss.Color("#26a621")
	RedColour            = lipgloss.Color("#FD8A8A")
	SpinnerColour        = lipgloss.Color("#439A97")

	AppStyle            = lipgloss.NewStyle().Margin(2)
	TextStyle           = lipgloss.NewStyle().Foreground(TextColour)
	FeintTextStyle      = lipgloss.NewStyle().Foreground(FeintColour)
	VeryFeintTextStyle  = lipgloss.NewStyle().Foreground(VeryFeintColour)
	HightlightTextStyle = lipgloss.NewStyle().Foreground(HighlightColour)

	Ellipsis      = HightlightTextStyle.Copy()
	HelpText      = TextStyle.Copy()
	HelpFeintText = FeintTextStyle.Copy()

	// Error

	Label = lipgloss.NewStyle().
		Background(RedColour).
		Foreground(TextColour).
		Bold(true)

	Reason = lipgloss.NewStyle().
		Foreground(RedColour).
		Bold(true)

	Cause = lipgloss.NewStyle().
		Foreground(TextColour)

	BorderBottom = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(BorderColour)
	BorderTop = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(BorderColour)
)
