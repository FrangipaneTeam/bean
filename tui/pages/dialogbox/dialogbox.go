package dialogbox

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/FrangipaneTeam/bean/internal/keymap"
	"github.com/FrangipaneTeam/bean/internal/theme"
)

const (
	questionSize = 50
	dialogHeight = 9
	marginRight  = 2
)

const (
	okValue = 1 << iota
	cancelValue
)

// Model is the model of the error panel.
type Model struct {
	tea.Model
	keys         *keymap.ListKeyMap
	width        int
	height       int
	question     string
	okValue      string
	cancelValue  string
	ActiveButton int
	theme        theme.Theme
}

// New returns a new model of the error panel.
func New(w int, h int, keymap *keymap.ListKeyMap) *Model {
	return &Model{
		width:        w,
		height:       h,
		keys:         keymap,
		ActiveButton: cancelValue,
		theme:        theme.Default(),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	// cmd  tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.ActiveButton = okValue

		case key.Matches(msg, m.keys.Right):
			m.ActiveButton = cancelValue
		}
	}
	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
	var doc strings.Builder

	okButton := m.theme.DialogBox.ActiveButtonStyle.MarginRight(marginRight).Render(m.okValue)
	cancelButton := m.theme.DialogBox.ButtonStyle.Render(m.cancelValue)
	// okButton := buttonStyle.MarginRight(2).Render(m.okValue)
	// cancelButton := activeButtonStyle.Render(m.cancelValue)

	if m.ActiveButton == cancelValue {
		okButton = m.theme.DialogBox.ButtonStyle.MarginRight(marginRight).Render(m.okValue)
		cancelButton = m.theme.DialogBox.ActiveButtonStyle.Render(m.cancelValue)
	}

	question := lipgloss.NewStyle().
		Width(questionSize).
		Align(lipgloss.Center).
		Render(m.question)

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.theme.DialogBox.Style.Render(ui),
		lipgloss.WithWhitespaceChars("frangipane"),
		lipgloss.WithWhitespaceForeground(m.theme.DialogBox.SubtleColour),
	)

	// dialog = lipgloss.NewStyle().Width(m.width - 2).Render(dialog)

	doc.WriteString(dialog)

	return doc.String()
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetSize(w int, h int) {
	m.width = w
	m.height = h
}

func (m *Model) SetDialogBox(question string, okValue string, cancelValue string) {
	m.ActiveButton = 2
	m.question = question
	m.okValue = okValue
	m.cancelValue = cancelValue
}

func GetCancelValue() int {
	return cancelValue
}
