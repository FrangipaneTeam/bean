package md

import (
	"log"
	"os"
	"strings"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/internal/crd"
	"github.com/FrangipaneTeam/bean/internal/examples"
	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/internal/template"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	tea "github.com/charmbracelet/bubbletea"
)

type testedStruct struct {
	crd.CRD
	Tested bool
}

type listTestedStruct map[string]map[string]testedStruct

// GenerateListTested generates the list of tested CRDs.
func GenerateListTested(c config.Provider) tea.Cmd {
	return func() tea.Msg {
		crdS, err := crd.GetCRDs(c.Path + "/package/crds")
		if err != nil {
			return errorpanel.ErrorMsg{Reason: "can't get crds", Cause: err}
		}
		examples := examples.GenerateExamplesList(c)
		switch ex := examples.(type) {
		case errorpanel.ErrorMsg:
			log.Print(ex.Cause.Error())
			os.Exit(1)

		case exlist.LoadedExamples:

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
					e, ok := val.(*exlist.Example)
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

			template.Process(
				markdownTemplate,
				c.Path+"/list-tested.md",
				data,
			)
		}
		return exlist.ListTestedDone{}
	}
}

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
