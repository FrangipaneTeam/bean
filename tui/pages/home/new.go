package home

import (
	"fmt"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	ex "github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/internal/theme"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/dialogbox"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/FrangipaneTeam/bean/tui/pages/exlist"
	"github.com/FrangipaneTeam/bean/tui/pages/footer"
	"github.com/FrangipaneTeam/bean/tui/pages/header"
	"github.com/FrangipaneTeam/bean/tui/pages/k8s"
	"github.com/FrangipaneTeam/bean/tui/pages/md"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	keys *keymap.ListKeyMap

	// pages
	common     *common.Model
	header     *header.Model
	footer     *footer.Model
	errorPanel *errorpanel.Model
	markdown   *md.Model
	k8s        *k8s.Model

	config config.Provider

	k8sCurrentIDView string
	k8sProgressMsg   string

	dialogbox *dialogbox.Model

	pages     *exlist.Model
	pagesList map[common.PageID]*common.Page

	width        int
	height       int
	centerHeight int
}

type tickK8SGet time.Time

// New returns a new model of the examples page.
func New(e ex.LoadedExamples, width, height int, c config.Provider) model {
	h, v := theme.AppStyle.GetFrameSize()

	rootKeys := keymap.NewListKeyMap()
	dialogKeys := keymap.NewListKeyMap()
	version := lipgloss.NewStyle().Foreground(theme.NotificationColour).Render("v" + c.Version)

	if c.NewVersion != "" {
		newVersion := lipgloss.NewStyle().
			Foreground(theme.NotificationColour).
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

	pagesModel := exlist.New(rootKeys, e, width-h, height-v-headerHeight-footerHeight)
	errorPanel := errorpanel.New(width-h, height-v-headerHeight-footerHeight)

	footer.SetExamplesList(pagesModel)

	dialogbox := dialogbox.New(
		width-h,
		height-v-headerHeight-footerHeight,
		dialogKeys,
	)
	commonM := common.New(
		rootKeys,
		pagesModel,
	)
	footer.SetCommonModel(commonM)
	k8s := k8s.New(rootKeys, commonM, pagesModel)
	markdown := md.New(
		width-h,
		height-v-headerHeight-footerHeight,
		rootKeys,
		commonM,
		c,
	)

	return model{
		keys:       rootKeys,
		header:     header,
		footer:     footer,
		common:     commonM,
		errorPanel: errorPanel,

		markdown:  markdown,
		dialogbox: dialogbox,
		k8s:       k8s,
		config:    c,
		pages:     pagesModel,
		pagesList: common.BeanPages(),
	}
}
