// Package pages initializes the examples list and watches for changes to the examples directory and the crd directory.
package examples

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/FrangipaneTeam/bean/config"
	yml "github.com/FrangipaneTeam/bean/pkg/yaml"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dietsche/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

// NotifyActivity is a struct that holds the name of the file that was changed.
type NotifyActivity struct {
	FileName string
}

// ResponseCRDMsg is a response to a crd activity.
type ResponseCRDMsg NotifyActivity

// ResponseExamplesMsg is a response to an examples activity.
type ResponseExamplesMsg NotifyActivity

// LoadExamples loads the examples from the examples directory.
func LoadExamples(c config.Provider) tea.Cmd {
	return func() tea.Msg {
		return GenerateExamplesList(c)
	}
}

// ListenForCRDActivity watches for crd file changes and sends a message to the channel.
func ListenForCRDActivity(ch chan NotifyActivity, c config.Provider) tea.Cmd {
	return func() tea.Msg {
		watcher, err := newWatcher()
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		defer watcher.Close()

		done := make(chan bool)
		go watchCRDFiles(watcher, done, ch)

		err = addCRDFolder(watcher, c.Path)
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		<-done
		return ch
	}
}

func newWatcher() (*fsnotify.Watcher, error) {
	return fsnotify.NewWatcher()
}

func watchCRDFiles(watcher *fsnotify.Watcher, done chan bool, ch chan NotifyActivity) tea.Msg {
	defer close(done)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errorpanel.ErrorMsg{
					Reason: "NewWatcher error",
					Cause:  errors.New("event not ok"),
				}
			}
			if yml.IsYamlFile(event.Name) {
				f := NotifyActivity{
					FileName: event.Name,
				}
				ch <- f
			}
		case errWatcher, ok := <-watcher.Errors:
			if !ok {
				return errorpanel.ErrorMsg{
					Reason: "NewWatcher error",
					Cause:  errWatcher,
				}
			}
		}
	}
}

func addCRDFolder(watcher *fsnotify.Watcher, path string) error {
	return watcher.Add(fmt.Sprintf("%s/package/crds", path))
}

// ListenForExamplesActivity watches for examples file changes and sends a message to the channel.
func ListenForExamplesActivity(ch chan NotifyActivity, c config.Provider) tea.Cmd {
	return func() tea.Msg {
		watcher, err := newRecursiveWatcher()
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		defer watcher.Close()

		done := make(chan bool)
		go watchExamplesFiles(watcher, done, ch)

		err = addExamplesFolder(watcher, c.Path)
		if err != nil {
			return errorpanel.ErrorMsg{
				Reason: "NewWatcher error",
				Cause:  err,
			}
		}
		<-done
		return ch
	}
}

func watchExamplesFiles(watcher *rfsnotify.RWatcher, done chan bool, ch chan NotifyActivity) {
	defer close(done)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if isExamplesFile(event.Name) {
				f := NotifyActivity{
					FileName: event.Name,
				}
				ch <- f
			}
		case errWatcher, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", errWatcher)
		}
	}
}

func isExamplesFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".yaml" || ext == ".yml" || ext == ".secret" || ext == ".extra"
}

func addExamplesFolder(watcher *rfsnotify.RWatcher, path string) error {
	return watcher.AddRecursive(fmt.Sprintf("%s/examples", path))
}

func newRecursiveWatcher() (*rfsnotify.RWatcher, error) {
	return rfsnotify.NewWatcher()
}

// WaitForCrdActivity waits for a message from the channel and returns a message.
func WaitForCrdActivity(sub <-chan NotifyActivity) tea.Cmd {
	return func() tea.Msg {
		return ResponseCRDMsg(<-sub)
	}
}

// WaitForExamplesActivity waits for a message from the channel and returns a message.
func WaitForExamplesActivity(sub <-chan NotifyActivity) tea.Cmd {
	return func() tea.Msg {
		return ResponseExamplesMsg(<-sub)
	}
}
