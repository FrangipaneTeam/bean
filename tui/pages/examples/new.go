package examples

import (
	"fmt"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/footer"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/FrangipaneTeam/bean/tui/pages/md"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width, height int
	keys          *tui.ListKeyMap
	oldKeys       *tui.ListKeyMap

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

	pages     *pages.Model
	pagesList map[pages.PageID]*pages.Page
}

type tickK8SGet time.Time

// New returns a new model of the examples page.
func New(e tui.LoadedExamples, width, height int, c config.Provider) model {
	h, v := tui.AppStyle.GetFrameSize()

	rootKeys := tui.NewListKeyMap()
	dialogKeys := tui.NewListKeyMap()
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

	// default activated keys
	rootKeys.EnableRootKeys()
	dialogKeys.EnableDialogBoxKeys()

	// common model
	k8sCmdList := make(map[string]*k8s.Cmd)
	pagesModel := pages.New(rootKeys, e, width-h, height-v-headerHeight-footerHeight)
	header.SetPagesModel(pagesModel)

	common := common.New(pagesModel)

	return model{
		keys:   rootKeys,
		header: header,
		footer: footer,
		common: common,

		errorPanel: errorpanel.New(width-h, height-v-headerHeight-footerHeight),
		markdown:   md.New(width-h, height-v-headerHeight-footerHeight),
		dialogbox: dialogbox.New(
			width-h,
			height-v-headerHeight-footerHeight,
			dialogKeys,
		),
		k8s:          k8s.New(rootKeys, pagesModel),
		k8sCmdList:   k8sCmdList,
		width:        width - h,
		height:       height - v,
		centerHeight: height - v - headerHeight - footerHeight,
		config:       c,
		pages:        pagesModel,
		pagesList:    pages.BeanPages(),
		progressK8SGet: progress.New(
			progress.WithSolidFill("#CBEDD5"),
			progress.WithoutPercentage(),
			progress.WithWidth(progressWidth),
		),
	}
}
