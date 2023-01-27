package k8s

import (
	"context"

	"github.com/FrangipaneTeam/bean/tui"
	"github.com/FrangipaneTeam/bean/tui/pages"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	progressWidth int = 10
)

// Model is the model of the error panel.
type Model struct {
	pages *pages.Model
	tea.Model
	keys        *tui.ListKeyMap
	width       int
	tickRunning bool
	CmdList     map[string]*Cmd
	GetProgress progress.Model
}

type Message struct {
	ShowDialogBox bool
	PreviousPage  pages.PageID
}

type Cmd struct {
	ID       string
	Done     bool
	Verb     string
	Files    []string
	Kind     string
	Result   string
	Cancel   context.CancelFunc
	FromPage pages.PageID
}

// New returns a new model of the k8s page.
func New(keymap *tui.ListKeyMap, pages *pages.Model) *Model {
	cmdList := make(map[string]*Cmd)
	return &Model{
		keys:    keymap,
		pages:   pages,
		CmdList: cmdList,
		GetProgress: progress.New(
			progress.WithSolidFill("#CBEDD5"),
			progress.WithoutPercentage(),
			progress.WithWidth(progressWidth),
		),
	}
}
