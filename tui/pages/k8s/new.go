package k8s

import (
	"context"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/elist"
	"github.com/charmbracelet/bubbles/progress"
)

const (
	progressWidth int = 10
)

// Model is the model of the error panel.
type Model struct {
	pages                 *elist.Model
	keys                  *tui.ListKeyMap
	width                 int
	tickRunning           bool
	CmdList               map[string]*Cmd
	GetProgress           progress.Model
	ShowDependenciesFiles bool
	common                *common.Model
}

type Message struct {
	ShowDialogBox bool
	PreviousPage  interface{}
}

type Cmd struct {
	ID       string
	Done     bool
	Verb     string
	Files    []string
	Kind     string
	Result   string
	Cancel   context.CancelFunc
	FromPage common.PageID
}

// New returns a new model of the k8s page.
func New(keymap *tui.ListKeyMap, common *common.Model, pages *elist.Model) *Model {
	cmdList := make(map[string]*Cmd)
	return &Model{
		keys:                  keymap,
		pages:                 pages,
		common:                common,
		CmdList:               cmdList,
		ShowDependenciesFiles: true,
		GetProgress: progress.New(
			progress.WithSolidFill("#CBEDD5"),
			progress.WithoutPercentage(),
			progress.WithWidth(progressWidth),
		),
	}
}
