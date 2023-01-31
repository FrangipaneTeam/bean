package md

import (
	"os"

	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// markdown is a struct that holds the content of a markdown file.
type markdown struct {
	content string
}

// renderMarkdown renders a markdown file.
func renderMarkdown(file string, wrap int) tea.Cmd {
	return func() tea.Msg {
		f, err := os.ReadFile(file)
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "could not read markdown file",
				Cause:  err,
			}
		}

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(wrap),
			glamour.WithStylePath("dracula"),
		)
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "new render markdown failed",
				Cause:  err,
			}
		}
		str, err := renderer.Render(string(f))
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "render markdown failed",
				Cause:  err,
			}
		}
		return markdown{
			content: str,
		}
	}
}
