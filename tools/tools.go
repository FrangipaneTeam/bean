// Package tools provides a set of tools
package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

// KubectlResult is the result of a kubectl command.
type KubectlResult struct {
	Verb  string
	Out   string
	CmdID string
	Index int
	Item  *tui.Example
}

// Markdown is a struct that holds the content of a markdown file.
type Markdown struct {
	Content string
}

// ErrorMsg should be sent to notify a user of an unrecoverable error.
type ErrorMsg struct {
	Reason string
	Cause  error
	CmdID  string
	Index  int
	Item   *tui.Example
}

// RenderMarkdown renders a markdown file.
func RenderMarkdown(file string, wrap int) tea.Cmd {
	return func() tea.Msg {
		f, err := os.ReadFile(file)
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

func listDirExamples(path string) ([]os.DirEntry, error) {
	examplesList, err := os.ReadDir(fmt.Sprintf("%s/examples", path))
	return examplesList, err
}

func createExampleList(dir string) ([]*tui.Example, *ErrorMsg) {
	kindList, err := os.ReadDir(dir)
	if err != nil {
		return nil, &ErrorMsg{
			Reason: "could not read examples directory",
			Cause:  err,
		}
	}

	exampleList := make([]*tui.Example, 0)

	// for the sub list of examples dir
	for _, sf := range kindList {
		if !sf.Type().IsRegular() {
			continue
		}

		// get only yaml files
		if isYaml := IsYamlFile(sf.Name()); !isYaml {
			continue
		}

		// open and parse yaml file
		yfile, errReadFile := os.ReadFile(
			fmt.Sprintf("%s/%s", dir, sf.Name()),
		)
		if errReadFile != nil {
			return nil, &ErrorMsg{
				Reason: "could not read examples directory",
				Cause:  errReadFile,
			}
		}
		var k *tui.Example
		err = yaml.Unmarshal(yfile, &k)
		if err != nil {
			return nil, &ErrorMsg{
				Reason: "unmarshal error",
				Cause:  err,
			}
		}

		// continue if unmarshal empty yaml
		if k == nil {
			continue
		}

		k.FullPath = dir + "/" + sf.Name()
		k.FileName = sf.Name()
		k.Desc = k.Kind + " → " + k.APIVersion
		k.ExampleID = strings.ToLower(fmt.Sprintf("%s.%s", k.Kind, k.APIVersion))

		// check for selector
		k.Selectors, k.Refs = k.FindSelectorsAndRefs()
		k.DependenciesFiles = map[string]bool{}

		// check for extra files
		extraFileCount, errCheckExtra := checkForExtraFile(dir, sf.Name())
		if errCheckExtra != nil {
			return nil, errCheckExtra
		}
		if extraFileCount > 0 {
			k.ExtraFileExist = true
			k.Desc = fmt.Sprintf("%s + %d extra", k.Desc, extraFileCount)
		}

		// check for secret file
		extraSecretCount, errCheckSecret := checkForSecretFile(dir, sf.Name())
		if errCheckSecret != nil {
			return nil, errCheckSecret
		}
		if extraSecretCount > 0 {
			k.SecretFileExist = true
			k.Desc = fmt.Sprintf("%s + %d secret", k.Desc, extraSecretCount)
		}

		exampleList = append(exampleList, k)
	}
	return exampleList, nil
}

func checkForExtraFile(dir string, file string) (int, *ErrorMsg) {
	var (
		extraK8S  *tui.Example
		extraKind int
	)

	extraYFile, err := os.ReadFile(
		fmt.Sprintf("%s/%s.extra", dir, file),
	)
	if err == nil {
		extraY, errSplitYaml := splitYAML(extraYFile)
		if errSplitYaml != nil {
			return extraKind, &ErrorMsg{
				Reason: "split yaml error",
				Cause:  errSplitYaml,
			}
		}
		for _, f := range extraY {
			err = yaml.Unmarshal(f, &extraK8S)
			if err != nil {
				return extraKind, &ErrorMsg{
					Reason: "unmarshal error",
					Cause:  err,
				}
			}
			extraKind++
		}
	}
	return extraKind, nil
}

func checkForSecretFile(dir string, file string) (int, *ErrorMsg) {
	var (
		extraK8S  *tui.Example
		extraKind int
	)

	extraSFile, err := os.ReadFile(
		fmt.Sprintf("%s/%s.secret", dir, file),
	)
	if err == nil {
		extraY, errSplitYaml := splitYAML(extraSFile)
		if errSplitYaml != nil {
			return extraKind, &ErrorMsg{
				Reason: "split yaml error",
				Cause:  errSplitYaml,
			}
		}

		for _, f := range extraY {
			err = yaml.Unmarshal(f, &extraK8S)
			if err != nil {
				return extraKind, &ErrorMsg{
					Reason: "unmarshal error",
					Cause:  err,
				}
			}
			extraKind++
		}
	}
	return extraKind, nil
}

// GenerateExamplesList generates the list of examples.
func GenerateExamplesList(c config.Provider) tea.Msg {
	examplesList, err := listDirExamples(c.Path)
	if err != nil {
		return ErrorMsg{
			Reason: "could not read examples directory",
			Cause:  err,
		}
	}
	var s tui.LoadedExamples
	s.Examples = make(map[string][]list.Item)
	rootExamples := []list.Item{}

	for _, dir := range examplesList {
		if !dir.IsDir() {
			continue
		}
		// Generate file dependencies
		dirName := dir.Name()

		items, errMsg := createExampleList(c.Path + "/examples/" + dirName)
		if errMsg != nil {
			return errMsg
		}

		for _, item := range items {
			s.Examples[dirName] = append(s.Examples[dirName], item)
		}

		e := &tui.Example{
			FileName: dirName,
			FullPath: dirName,
			Desc:     fmt.Sprintf("%d examples", len(s.Examples[dirName])),
		}
		rootExamples = append(rootExamples, e)
	}

	s.Examples["-"] = rootExamples

	examplesWithDependencies := tui.ExamplesDetails{}
	for d, e := range s.Examples {
		if d == "-" {
			continue
		}
		for _, ex := range e {
			example, ok := ex.(*tui.Example)
			if ok {
				examplesWithDependencies[example.ExampleID] = example
			}
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

// processTemplate processes a template file and writes the output to a file.
func processTemplate(rawTemplate string, outputFile string, data interface{}) {
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
		//TODO: handle error in app
		log.Print(err)
	}
}

type testedStruct struct {
	crd.CRD
	Tested bool
}

type listTestedStruct map[string]map[string]testedStruct

// Init initializes the list.
func (l *listTestedStruct) Init() {
	*l = make(map[string]map[string]testedStruct)
}

// Add adds a CRD to the list.
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

// CheckIfTested checks if the resource is tested.
func (l *listTestedStruct) CheckIfTested(group, kind string) {
	if _, isMap := (*l)[group]; isMap {
		if entry, ok := (*l)[group][kind]; ok {
			entry.Tested = true
			(*l)[group][kind] = entry
		}
	}
}

// GenerateListTested generates the list of tested CRDs.
func GenerateListTested(c config.Provider) tea.Msg {
	crdS, err := crd.GetCRDs(c.Path + "/package/crds")
	if err != nil {
		return ErrorMsg{Reason: "can't get crds", Cause: err}
	}
	examples := GenerateExamplesList(c)
	switch ex := examples.(type) {
	case ErrorMsg:
		log.Print(ex.Cause.Error())
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
				e, ok := val.(*tui.Example)
				if ok {
					apiVersion := strings.Split(e.APIVersion, "/")
					data.CheckIfTested(apiVersion[0], e.Kind)
				}
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

		processTemplate(
			markdownTemplate,
			c.Path+"/list-tested.md",
			data,
		)
	}
	return tui.ListTestedDone{}
}

// IsYamlFile returns true if the file is a yaml file.
func IsYamlFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".yaml" || ext == ".yml"
}
