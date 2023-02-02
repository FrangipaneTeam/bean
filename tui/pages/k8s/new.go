package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/progress"

	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/exlist"
)

const (
	progressWidth int = 10
)

// Model is the model of the error panel.
type Model struct {
	pages                 *exlist.Model
	keys                  *keymap.ListKeyMap
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
	Debug    bool
}

// New returns a new model of the k8s page.
func New(keymap *keymap.ListKeyMap, common *common.Model, pages *exlist.Model) *Model {
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
