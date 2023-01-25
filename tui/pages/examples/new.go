package examples

import (
	"fmt"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/footer"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/FrangipaneTeam/bean/tui/pages/md"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	exampleList map[string][]list.Item
	currentList list.Model

	listName      string
	width, height int
	keys          *tui.ListKeyMap
	oldKeys       *tui.ListKeyMap

	showDependenciesFiles bool

	errorRaised bool

	// pages
	common     *common.Model
	header     header.Model
	footer     footer.Model
	errorPanel errorpanel.Model
	markdown   md.Model
	k8s        k8s.Model

	config config.Provider

	k8sCurrentIDView string
	k8sProgressMsg   string
	progressK8SGet   progress.Model
	tickRunning      bool

	listOldHeight int
	centerHeight  int

	dialogbox dialogbox.Model

	k8sCmdList map[string]*k8s.Cmd
}

type tickK8SGet time.Time

// New returns a new model of the examples page.
func New(e tui.LoadedExamples, width, height int, c config.Provider) model {
	h, v := tui.AppStyle.GetFrameSize()

	rootKeys := tui.NewListKeyMap()
	dialogKeys := tui.NewListKeyMap()
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

	// default activated keys
	rootKeys.EnableRootKeys()
	dialogKeys.EnableDialogBoxKeys()

	// common model
	common := common.New(rootKeys)

	return model{
		exampleList: e.Examples,
		currentList: list,
		listName:    "-",

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
		k8s:          k8s.New(rootKeys, common),
		common:       common,
		width:        width - h,
		height:       height - v,
		centerHeight: height - v - headerHeight - footerHeight,
		config:       c,

		showDependenciesFiles: true,
		progressK8SGet: progress.New(
			progress.WithSolidFill("#CBEDD5"),
			progress.WithoutPercentage(),
			progress.WithWidth(progressWidth),
		),

		k8sCmdList: make(map[string]*k8s.Cmd),
	}
}
