// Package tools provides a set of tools
package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/pkg/crd"
	"github.com/FrangipaneTeam/bean/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"gopkg.in/yaml.v3"
)

// KubectlResult is the result of a kubectl command
type KubectlResult struct {
	Verb string
	Out  string
}

// Markdown is a struct that holds the content of a markdown file
type Markdown struct {
	Content string
}

// ErrorMsg should be sent to notify a user of an unrecoverable error
type ErrorMsg struct {
	Reason string
	Cause  error
}

// Kubectl runs a kubectl command
func Kubectl(verb string, file string) tea.Cmd {
	return func() tea.Msg {
		args := []string{verb, "-f", file}

		cmd := exec.Command("kubectl", args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return ErrorMsg{
				Reason: fmt.Sprintf("command kubectl %s failed", strings.Join(args, " ")),
				Cause:  errors.New(stderr.String()),
			}
		}
		return KubectlResult{
			Out:  stdout.String(),
			Verb: verb,
		}
	}
}

// RenderMarkdown renders a markdown file
func RenderMarkdown(file string, wrap int) tea.Cmd {
	return func() tea.Msg {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return ErrorMsg{
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
			return ErrorMsg{
				Reason: "new render markdown failed",
				Cause:  err,
			}
		}
		str, err := renderer.Render(string(f))
		if err != nil {
			return ErrorMsg{
				Reason: "render markdown failed",
				Cause:  err,
			}
		}
		return Markdown{
			Content: str,
		}

	}
}

// GenerateExamplesList generates the list of examples
// nolint:gocyclo // TODO need to refactor this function
func GenerateExamplesList(c config.Provider) tea.Msg {
	examplesList, err := ioutil.ReadDir(fmt.Sprintf("%s/examples", c.Path))
	if err != nil {
		return ErrorMsg{
			Reason: "could not read examples directory",
			Cause:  err,
		}
	}
	var s tui.LoadedExamples
	s.Examples = make(map[string][]list.Item)
	for _, exL := range examplesList {
		if exL.IsDir() {
			kindList, err := ioutil.ReadDir(fmt.Sprintf("%s/examples/%s", c.Path, exL.Name()))
			if err != nil {
				return ErrorMsg{
					Reason: "could not read examples directory",
					Cause:  err,
				}
			}
			// for the root list of examples dir
			e := &tui.Example{
				FileName: exL.Name(),
			}

			// for the sub list of examples dir
			for _, sf := range kindList {
				if !sf.Mode().IsRegular() {
					continue
				}

				// get only yaml files
				if sf.Name()[len(sf.Name())-4:] != "yaml" &&
					sf.Name()[len(sf.Name())-3:] != "yml" {
					continue
				}

				// open and parse yaml file
				yfile, err := ioutil.ReadFile(
					fmt.Sprintf("%s/examples/%s/%s", c.Path, exL.Name(), sf.Name()),
				)
				if err != nil {
					return ErrorMsg{
						Reason: "could not read examples directory",
						Cause:  err,
					}
				}
				var k *tui.Example
				err = yaml.Unmarshal(yfile, &k)
				if err != nil {
					return ErrorMsg{
						Reason: "unmarshal error",
						Cause:  err,
					}
				}

				// continue if unmarshal empty yaml
				if k == nil {
					continue
				}

				k.FileName = c.Path + "/examples/" + exL.Name() + "/" + sf.Name()
				k.Desc = k.Kind + " → " + k.APIVersion
				k.ExampleID = strings.ToLower(fmt.Sprintf("%s.%s", k.Kind, k.APIVersion))

				// check for selector
				k.Selectors = k.FindForProvider("Selector")
				k.Refs = k.FindForProvider("Refs")
				k.DependenciesFiles = map[string]bool{}

				// check for extra file examples
				var extraK8S *tui.Example
				extraYFile, err := ioutil.ReadFile(
					fmt.Sprintf("%s/examples/%s/%s.extra", c.Path, exL.Name(), sf.Name()),
				)
				if err == nil {
					extraY, err := splitYAML(extraYFile)
					if err != nil {
						return ErrorMsg{
							Reason: "split yaml error",
							Cause:  err,
						}
					}
					extraKind := 0
					for _, f := range extraY {
						err = yaml.Unmarshal(f, &extraK8S)
						if err != nil {
							return ErrorMsg{
								Reason: "unmarshal error",
								Cause:  err,
							}
						}
						extraKind++
					}
					k.ExtraFileExist = true
					k.Desc = fmt.Sprintf("%s + %d extra", k.Desc, extraKind)
				}

				// check for secret file
				extraSFile, err := ioutil.ReadFile(
					fmt.Sprintf("%s/examples/%s/%s.secret", c.Path, exL.Name(), sf.Name()),
				)
				if err == nil {
					extraY, err := splitYAML(extraSFile)
					if err != nil {
						return ErrorMsg{
							Reason: "split yaml error",
							Cause:  err,
						}
					}
					extraKind := 0
					for _, f := range extraY {
						err = yaml.Unmarshal(f, &extraK8S)
						if err != nil {
							return ErrorMsg{
								Reason: "unmarshal error",
								Cause:  err,
							}
						}
						extraKind++
					}

					k.SecretFileExist = true
					k.SecretFile = fmt.Sprintf(
						"%s/examples/%s/%s.secret",
						c.Path,
						exL.Name(),
						sf.Name(),
					)
					k.Desc = fmt.Sprintf("%s + %d secret", k.Desc, extraKind)
				}

				s.Examples[exL.Name()] = append(s.Examples[exL.Name()], k)
			}
			if _, ok := s.Examples[exL.Name()]; ok {
				s.Examples["-"] = append(s.Examples["-"], e)
				e.Desc = fmt.Sprintf("%d examples", len(s.Examples[exL.Name()]))
			}
		}
	}

	// Generate file dependencies
	examplesWithDependencies := tui.ExamplesDetails{}
	for d, e := range s.Examples {
		if d == "-" {
			continue
		}
		for _, ex := range e {
			example := ex.(*tui.Example)
			examplesWithDependencies[example.ExampleID] = example
		}

	}
	examplesWithDependencies.FindDependencies()
	return s
}

func splitYAML(resources []byte) ([][]byte, error) {
	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

// ProcessTemplate processes a template file and writes the output to a file
func ProcessTemplate(rawTemplate string, outputFile string, data interface{}) {
	var err error

	// Read the template file
	t := template.Must(template.New("listTested").Parse(rawTemplate)) // .Funcs(funcMap)

	// Parse the template
	t, err = t.Parse(rawTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// create a new file
	file, _ := os.Create(outputFile)
	defer file.Close()

	// Execute the template
	err = t.ExecuteTemplate(file, "listTested", data)
	if err != nil {
		log.Fatal(err)
	}

}

type testedStruct struct {
	crd.CRD
	Tested bool
}

type listTestedStruct map[string]map[string]testedStruct

// Init initializes the list
func (l *listTestedStruct) Init() {
	*l = make(map[string]map[string]testedStruct)
}

// Add adds a CRD to the list
func (l *listTestedStruct) Add(c crd.CRD) {
	// if Group or Kind is empty, skip
	if c.Spec.Group == "" || c.Spec.Names.Kind == "" {
		return
	}

	if _, ok := (*l)[c.Spec.Group]; !ok {
		(*l)[c.Spec.Group] = make(map[string]testedStruct)
	}

	(*l)[c.Spec.Group][c.Spec.Names.Kind] = testedStruct{
		CRD:    c,
		Tested: false,
	}

}

// CheckIfTested checks if the resource is tested
func (l *listTestedStruct) CheckIfTested(group, kind string) {
	if _, ok := (*l)[group]; ok {
		if entry, ok := (*l)[group][kind]; ok {
			entry.Tested = true
			(*l)[group][kind] = entry
		}
	}
}

// GenerateListTested generates the list of tested CRDs
func GenerateListTested(c config.Provider) tea.Msg {
	crdS, err := crd.GetCRDs(c.Path + "/package/crds")
	if err != nil {
		return ErrorMsg{Reason: "can't get crds", Cause: err}
	}
	examples := GenerateExamplesList(c)
	switch ex := examples.(type) {
	case ErrorMsg:
		fmt.Println(ex.Cause.Error())
		os.Exit(1)

	case tui.LoadedExamples:

		data := &listTestedStruct{}
		data.Init()

		for _, val := range crdS {
			data.Add(val)
		}

		for k, v := range ex.Examples {
			if k == "-" {
				continue
			}
			for _, val := range v {
				e := val.(*tui.Example)
				apiVersion := strings.Split(e.APIVersion, "/")
				data.CheckIfTested(apiVersion[0], e.Kind)
			}
		}

		var markdownTemplate = `
# List of tested resources

{{ range $group, $kinds := . }}
## {{ $group }}
| Kind | Tested |
| ---- | ------ |{{ range $kind, $crd := $kinds }}
| {{ $kind }} | {{ if $crd.Tested }} :white_check_mark: {{else}} :x: {{end}} |{{ end }}
{{ end }}
`

		ProcessTemplate(
			markdownTemplate,
			"list-tested.md",
			data,
		)

	}

	return tui.ListTestedDone{}
}
