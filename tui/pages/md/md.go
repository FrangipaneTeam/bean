// Package md provides a markdown viewer
package md

import (
	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/tui/pages/common"
	"github.com/FrangipaneTeam/bean/tui/pages/elist"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	paddingRight = 2
)

// Model is the model of the markdown viewer.
type Model struct {
	Viewport viewport.Model
	pages    *elist.Model
	keys     *keymap.ListKeyMap
	config   config.Provider
	common   *common.Model
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// New creates a new markdown viewer model.
func New(
	w, h int,
	keys *keymap.ListKeyMap,
	common *common.Model,
	pages *elist.Model,
	config config.Provider,
) *Model {
	vp := viewport.New(w, h)
	vp.Style = lipgloss.NewStyle().
		Width(w).
		PaddingRight(paddingRight)
	return &Model{
		Viewport: vp,
		pages:    pages,
		keys:     keys,
		config:   config,
		common:   common,
	}
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.pages.CurrentList.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.ShowRessources):
			m.keys.EnableViewPortKeys()

			if m.common.GetViewName() == common.PRoot {
				cmd = renderMarkdown(m.config.Path+"/list-resources.md", common.Width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.ShowTested):
			m.keys.EnableViewPortKeys()
			if m.common.GetViewName() == common.PRoot {
				cmd = renderMarkdown(m.config.Path+"/list-tested.md", common.Width)

				return m, cmd
			}

		case key.Matches(msg, m.keys.GenerateListTested):
			cmd = GenerateListTested(m.config)
			return m, cmd
		}
	case markdown:
		m.common.SetViewName(common.PViewPort)
		m.Viewport.SetContent(msg.content)
		m.Viewport.GotoTop()
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.SetSize(common.Width, common.CenterHeight)

	case common.ResizeMsg:
		m.SetSize(common.Width, common.CenterHeight)
	}

	if m.common.GetViewName() == common.PViewPort ||
		m.common.GetViewName() == common.PPrintActions {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	return ""
}

// SetSize sets the size of the model.
func (m *Model) SetSize(w, h int) {
	m.Viewport.Height = h
	m.Viewport.Width = w
	m.Viewport.Style = lipgloss.NewStyle().
		Width(w).
		PaddingRight(paddingRight)
	viewport.Sync(m.Viewport)
}
