// Package pages initializes the examples list and watches for changes to the examples directory and the crd directory.
package pages

import (
	"fmt"
	"path/filepath"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tools"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dietsche/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

// Model is the model of the pages
type Model interface {
	tea.Model

	Resize(width, height int) Model
	Width() int
	Height() int
}

// NotifyActivity is a struct that holds the name of the file that was changed
type NotifyActivity struct {
	FileName string
}

// ResponseCRDMsg is a response to a crd activity
type ResponseCRDMsg NotifyActivity

// ResponseExamplesMsg is a response to an examples activity
type ResponseExamplesMsg NotifyActivity

// LoadExamples loads the examples from the examples directory
func LoadExamples(c config.Provider) tea.Cmd {
	return func() tea.Msg {
		return tools.GenerateExamplesList(c)
	}
}

// GenerateListTested generates the list of tested CRDs
func GenerateListTested(c config.Provider) tea.Cmd {
	return func() tea.Msg {
		return tools.GenerateListTested(c)
	}
}

// ListenForCRDActivity watches for crd file changes and sends a message to the channel
func ListenForCRDActivity(ch chan NotifyActivity, c config.Provider) tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return tools.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}

		defer watcher.Close()

		done := make(chan bool)

		// go routine to watch for crd file changes via fsnotify
		go func() tea.Msg {
			defer close(done)

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						if err != nil {
							return tools.ErrorMsg{
								Reason: "NewWatcher error",
								Cause:  err,
							}
						}
					}
					switch filepath.Ext(event.Name) {
					case ".yaml", ".yml":
						f := NotifyActivity{
							FileName: event.Name,
						}
						ch <- f
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return tools.ErrorMsg{
							Reason: "NewWatcher error",
							Cause:  err,
						}
					}
				}
			}
		}() // end of go routine

		err = watcher.Add(fmt.Sprintf("%s/package/crds", c.Path))
		if err != nil {
			return tools.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		<-done
		return ch
	}
}

// ListenForExamplesActivity watches for examples file changes and sends a message to the channel
// nolint: gocyclo // FIXME refactor
func ListenForExamplesActivity(ch chan NotifyActivity, c config.Provider) tea.Cmd {
	return func() tea.Msg {
		watcher, err := rfsnotify.NewWatcher()
		if err != nil {
			return tools.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		defer watcher.Close()

		done := make(chan bool)

		// go routine to watch for crd file changes via fsnotify
		go func() tea.Msg {
			defer close(done)

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return tools.ErrorMsg{
							Reason: "NewWatcher error",
							Cause:  err,
						}
					}
					// ignore changes if the file extension is not yaml, yml, secret or extra
					switch filepath.Ext(event.Name) {
					case ".yaml", ".yml", ".secret", ".extra":
						f := NotifyActivity{
							FileName: event.Name,
						}
						ch <- f
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return tools.ErrorMsg{
							Reason: "NewWatcher error",
							Cause:  err,
						}
					}
				}
			}
		}() // end of go routine

		err = watcher.AddRecursive(fmt.Sprintf("%s/examples", c.Path))
		if err != nil {
			return tools.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		<-done

		return ch
	}
}

// WaitForCrdActivity waits for a message from the channel and returns a message
func WaitForCrdActivity(sub <-chan NotifyActivity) tea.Cmd {
	return func() tea.Msg {
		return ResponseCRDMsg(<-sub)
	}
}

// WaitForExamplesActivity waits for a message from the channel and returns a message
func WaitForExamplesActivity(sub <-chan NotifyActivity) tea.Cmd {
	return func() tea.Msg {
		return ResponseExamplesMsg(<-sub)
	}
}
