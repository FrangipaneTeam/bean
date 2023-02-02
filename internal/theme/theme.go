package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

const (
	padding      = 1
	paddingRight = 1
)

type Theme struct {
	Colour             Colour         `yaml:"colour"`
	AppStyle           lipgloss.Style `yaml:"-"`
	CheckMark          string         `yaml:"-"`
	RunningMark        string         `yaml:"-"`
	ErrorMark          string         `yaml:"-"`
	ShortSeparator     lipgloss.Style `yaml:"-"`
	HelpTextStyle      lipgloss.Style `yaml:"-"`
	HelpFeintTextStyle lipgloss.Style `yaml:"-"`
	FeintTextStyle     lipgloss.Style `yaml:"-"`
	TextStyle          lipgloss.Style `yaml:"-"`
	Divider            string         `yaml:"-"`
	BorderBottom       lipgloss.Style `yaml:"-"`
	ErrorPanel         ErrorPanel     `yaml:"-"`
	List               List           `yaml:"list"`
	DialogBox          DialogBox      `yaml:"dialogBox"`
	Symbol             Symbol         `yaml:"symbol"`
}

type Symbol struct {
	CheckMark   string `yaml:"checkMark"`
	RunningMark string `yaml:"runningMark"`
	ErrorMark   string `yaml:"errorMark"`
	Divider     string `yaml:"divider"`
}

type Colour struct {
	Primary      lipgloss.Color `yaml:"primary"`
	Notification lipgloss.Color `yaml:"notification"`
	Spinner      lipgloss.Color `yaml:"spinner"`
	Border       lipgloss.Color `yaml:"border"`
	OK           lipgloss.Color `yaml:"ok"`
	Warning      lipgloss.Color `yaml:"warning"`
	Error        lipgloss.Color `yaml:"error"`
	Help         lipgloss.Color `yaml:"help"`
	Feint        lipgloss.Color `yaml:"feint"`
	Divider      lipgloss.Color `yaml:"divider"`
}

type ErrorPanel struct {
	Reason lipgloss.Style
	Cause  lipgloss.Style
}

type List struct {
	SelectedTitleColor       lipgloss.Color `yaml:"selectedTitle"`
	SelectedTitleBorderColor lipgloss.Color `yaml:"selectedTitleBorder"`
	SelectedDescColor        lipgloss.Color `yaml:"selectedDesc"`
	SelectedDescBorderColor  lipgloss.Color `yaml:"selectedDescBorder"`
}

type DialogBox struct {
	Style             lipgloss.Style `yaml:"-"`
	Border            lipgloss.Color `yaml:"border"`
	SubtleColour      lipgloss.Color `yaml:"subtle"`
	ButtonColour      ButtonColour   `yaml:"button"`
	ButtonStyle       lipgloss.Style `yaml:"-"`
	ActiveButtonStyle lipgloss.Style `yaml:"-"`
}

type ButtonColour struct {
	Normal lipgloss.Color `yaml:"normal"`
	Active lipgloss.Color `yaml:"active"`
	Text   lipgloss.Color `yaml:"text"`
}

// Default is the default theme
func Default() (t Theme) {
	viper.SetDefault("theme.colour.primary", "#F6F5FC")
	viper.SetDefault("theme.colour.notification", "#FF5F87")
	viper.SetDefault("theme.colour.divider", "#FF5F87")
	viper.SetDefault("theme.colour.spinner", "#FF5F87")
	viper.SetDefault("theme.symbol.checkMark", "✓")
	viper.SetDefault("theme.symbol.runningMark", "⚙")
	viper.SetDefault("theme.symbol.errorMark", "✗")
	viper.SetDefault("theme.symbol.divider", "•")

	// viper.SetDefault("theme.colour.primary", "#F6F5FC")
	viper.SetDefault("theme.colour.warning", "#F0C674")
	viper.SetDefault("theme.colour.ok", "#B5BD68")
	viper.SetDefault("theme.colour.error", "#CC6666")
	viper.SetDefault("theme.colour.help", "#C5C8C6")
	viper.SetDefault("theme.colour.feint", "#969896")
	viper.SetDefault("theme.colour.borderBottom", "#969896")

	// list
	viper.SetDefault("theme.list.selectedTitle", "#FF5F87")
	viper.SetDefault("theme.list.selectedTitleBorder", "#FF5F87")
	viper.SetDefault("theme.list.selectedDesc", "#81A2BE")
	viper.SetDefault("theme.list.selectedDescBorder", "#FF5F87")

	// dialog box
	viper.SetDefault("theme.dialogBox.subtle", "#383838")
	viper.SetDefault("theme.dialogBox.border", "#874BFD")
	viper.SetDefault("theme.dialogBox.button.normal", "#888B7E")
	viper.SetDefault("theme.dialogBox.button.active", "#F25D94")
	viper.SetDefault("theme.dialogBox.button.text", "#FFF7DB")

	notificationColour := lipgloss.Color(viper.GetString("theme.colour.notification"))
	warningColour := lipgloss.Color(viper.GetString("theme.colour.warning"))
	okColour := lipgloss.Color(viper.GetString("theme.colour.ok"))
	errorColour := lipgloss.Color(viper.GetString("theme.colour.error"))
	helpColour := lipgloss.Color(viper.GetString("theme.colour.help"))
	feintColour := lipgloss.Color(viper.GetString("theme.colour.feint"))
	dividerColour := lipgloss.Color(viper.GetString("theme.colour.divider"))
	divider := viper.GetString("theme.symbol.divider")
	primaryColour := lipgloss.Color(viper.GetString("theme.colour.primary"))
	borderColour := lipgloss.Color(viper.GetString("theme.colour.borderBottom"))
	spinnerColour := lipgloss.Color(viper.GetString("theme.colour.spinner"))

	// list
	selectedTitleColor := lipgloss.Color(viper.GetString("theme.list.selectedTitle"))
	selectedTitleBorderColor := lipgloss.Color(viper.GetString("theme.list.selectedTitleBorder"))
	selectedDescColor := lipgloss.Color(viper.GetString("theme.list.selectedDesc"))
	selectedDescBorderColor := lipgloss.Color(viper.GetString("theme.list.selectedDescBorder"))

	// dialog box
	dialogBoxBorder := lipgloss.Color(viper.GetString("theme.dialogBox.border"))
	dialogBoxButtonText := lipgloss.Color(viper.GetString("theme.dialogBox.button.text"))
	dialogBoxButtonNormal := lipgloss.Color(viper.GetString("theme.dialogBox.button.normal"))
	dialogBoxButtonActive := lipgloss.Color(viper.GetString("theme.dialogBox.button.active"))

	// mark
	checkMark := viper.GetString("theme.symbol.checkMark")
	runningMark := viper.GetString("theme.symbol.runningMark")
	errorMark := viper.GetString("theme.symbol.errorMark")

	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(dialogBoxBorder).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	subtleColour := lipgloss.Color(viper.GetString("theme.dialogBox.subtle"))

	buttonStyle := lipgloss.NewStyle().
		Foreground(dialogBoxButtonText).
		Background(dialogBoxButtonNormal).
		Padding(0, 1).
		MarginTop(1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(dialogBoxButtonText).
		Background(dialogBoxButtonActive)

	t = Theme{
		AppStyle: lipgloss.NewStyle().Padding(padding),
		Colour: Colour{
			Primary:      primaryColour,
			Notification: notificationColour,
			Border:       borderColour,
			Spinner:      spinnerColour,
			OK:           okColour,
			Warning:      warningColour,
			Error:        errorColour,
			Help:         helpColour,
			Feint:        feintColour,
			Divider:      dividerColour,
		},
		TextStyle: lipgloss.NewStyle().Foreground(primaryColour),
		CheckMark: lipgloss.NewStyle().
			Foreground(okColour).
			PaddingRight(paddingRight).
			Render(checkMark),
		RunningMark: lipgloss.NewStyle().
			Foreground(warningColour).
			PaddingRight(paddingRight).
			Render(runningMark),
		ErrorMark: lipgloss.NewStyle().
			Foreground(errorColour).
			PaddingRight(paddingRight).
			Render(errorMark),
		ShortSeparator: lipgloss.NewStyle().
			Foreground(lipgloss.Color(notificationColour)),
		HelpTextStyle:      lipgloss.NewStyle().Foreground(helpColour),
		HelpFeintTextStyle: lipgloss.NewStyle().Foreground(feintColour),
		FeintTextStyle:     lipgloss.NewStyle().Foreground(feintColour),
		Divider: lipgloss.NewStyle().
			SetString(divider).
			Padding(0, 1).
			Foreground(dividerColour).
			String(),
		ErrorPanel: ErrorPanel{
			Reason: lipgloss.NewStyle().
				Foreground(errorColour).
				Bold(true),
			Cause: lipgloss.NewStyle().
				Foreground(primaryColour),
		},
		BorderBottom: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderColour),
		List: List{
			SelectedTitleColor:       selectedTitleColor,
			SelectedTitleBorderColor: selectedTitleBorderColor,
			SelectedDescColor:        selectedDescColor,
			SelectedDescBorderColor:  selectedDescBorderColor,
		},
		DialogBox: DialogBox{
			Style:             dialogBoxStyle,
			SubtleColour:      subtleColour,
			ButtonStyle:       buttonStyle,
			ActiveButtonStyle: activeButtonStyle,
			Border:            dialogBoxBorder,
			ButtonColour: ButtonColour{
				Normal: dialogBoxButtonText,
				Active: dialogBoxButtonActive,
				Text:   dialogBoxButtonText,
			},
		},
		Symbol: Symbol{
			CheckMark:   checkMark,
			RunningMark: runningMark,
			ErrorMark:   errorMark,
			Divider:     divider,
		},
	}

	return t
}
