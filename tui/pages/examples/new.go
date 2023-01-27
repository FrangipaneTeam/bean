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
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width, height int
	keys          *tui.ListKeyMap
	oldKeys       *tui.ListKeyMap

	// pages
	common     *common.Model
	header     header.Model
	footer     footer.Model
	errorPanel *errorpanel.Model
	markdown   md.Model
	k8s        *k8s.Model

	config config.Provider

	k8sCurrentIDView string
	k8sProgressMsg   string

	listOldHeight int
	centerHeight  int

	dialogbox dialogbox.Model

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
	pagesModel := pages.New(rootKeys, e, width-h, height-v-headerHeight-footerHeight)
	errorPanel := errorpanel.New(width-h, height-v-headerHeight-footerHeight)
	header.SetPagesModel(pagesModel)

	k8s := k8s.New(rootKeys, pagesModel)
	header.SetK8SModel(k8s)
	common := common.New(pagesModel, &header, errorPanel, k8s)

	return model{
		keys:       rootKeys,
		header:     header,
		footer:     footer,
		common:     common,
		errorPanel: errorPanel,

		markdown: md.New(width-h, height-v-headerHeight-footerHeight),
		dialogbox: dialogbox.New(
			width-h,
			height-v-headerHeight-footerHeight,
			dialogKeys,
		),
		k8s:          k8s,
		width:        width - h,
		height:       height - v,
		centerHeight: height - v - headerHeight - footerHeight,
		config:       c,
		pages:        pagesModel,
		pagesList:    pages.BeanPages(),
	}
}
