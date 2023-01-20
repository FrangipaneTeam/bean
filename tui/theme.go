package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColour        = lipgloss.Color("#3a1577")
	SecondaryColour      = lipgloss.Color("#5b1689")
	BorderColour         = lipgloss.Color("#807d8a")
	FeintColour          = lipgloss.Color("#807d8a")
	VeryFeintColour      = lipgloss.Color("#5e5e5e")
	TextColour           = lipgloss.Color("#f6f5fc")
	HighlightColour      = lipgloss.Color("#97DECE")
	HighlightFeintColour = lipgloss.Color("#CBEDD5")
	AmberColour          = lipgloss.Color("#e68a35")
	GreenColour          = lipgloss.Color("#26a621")
	RedColour            = lipgloss.Color("#FD8A8A")
	SpinnerColour        = lipgloss.Color("#FF5F87")
	DividerColour        = lipgloss.Color("#FF5F87")
	NotificationColour   = lipgloss.Color("#FF5F87")
	special              = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	AppStyle            = lipgloss.NewStyle().Padding(1)
	TextStyle           = lipgloss.NewStyle().Foreground(TextColour)
	FeintTextStyle      = lipgloss.NewStyle().Foreground(FeintColour)
	VeryFeintTextStyle  = lipgloss.NewStyle().Foreground(VeryFeintColour)
	HightlightTextStyle = lipgloss.NewStyle().Foreground(HighlightColour)

	Divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(DividerColour).
		String()

	Ellipsis      = NotificationStyle.Copy()
	HelpText      = TextStyle.Copy()
	HelpFeintText = FeintTextStyle.Copy()

	// Error

	CheckMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	RunningMark = lipgloss.NewStyle().SetString("?").
			Foreground(AmberColour).
			PaddingRight(1).
			String()

	ErrorMark = lipgloss.NewStyle().SetString("x").
			Foreground(AmberColour).
			PaddingRight(1).
			String()

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

	// StatusBarStyle = lipgloss.NewStyle().
	// Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
	// Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
	NotificationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(NotificationColour))
)
