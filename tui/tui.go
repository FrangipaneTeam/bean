package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/exp/maps"
)

// Example is a struct that holds the details of an example
type Example struct {
	FileName        string
	ExampleID       string
	Desc            string
	ExtraFileExist  bool
	ExtraFile       string
	SecretFileExist bool
	SecretFile      string
	Selectors       map[string]bool
	Refs            map[string]bool

	DependenciesFiles map[string]bool

	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Annotations struct {
			MetaUpboundIoExampleID string `yaml:"meta.upbound.io/example-id"`
		} `yaml:"annotations"`
		Labels struct {
			TestingUpboundIoExampleName string `yaml:"testing.upbound.io/example-name"`
		} `yaml:"labels"`
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		ForProvider map[interface{}]interface{} `yaml:"forProvider"`
	} `yaml:"spec"`
}

// LoadedExamples is a struct that holds the loaded examples
type LoadedExamples struct {
	Examples map[string][]list.Item
}

// ListTestedDone is a struct that holds the done message
type ListTestedDone struct{}

// ExamplesDetails is a map of Example
type ExamplesDetails map[string]*Example

// Title returns the title of the example
func (e Example) Title() string { return e.FileName }

// Description returns the description of the example
func (e Example) Description() string { return e.Desc }

// FilterValue returns the value to filter on
func (e Example) FilterValue() string { return e.FileName }

// HaveExtraFile returns true if the example has an extra file
func (e Example) HaveExtraFile() bool { return e.ExtraFileExist }

// HaveSecretFile returns true if the example has a secret file
func (e Example) HaveSecretFile() bool { return e.SecretFileExist }

// HaveDependenciesFiles returns true if the example has dependencies files
func (e Example) HaveDependenciesFiles() bool { return len(e.DependenciesFiles) > 0 }

// GetExampleID returns the example ID
func (e Example) GetExampleID() string { return e.ExampleID }

// DependenciesFilesList returns a list of dependencies files
func (e Example) DependenciesFilesList() []string {
	list := []string{}
	for k := range e.DependenciesFiles {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}

// FindForProvider returns a map of all the fields containing the pattern
func (e *Example) FindForProvider(pattern string) map[string]bool {
	// FindForProvider Recursively search all fields containing $parameter in the ForProvider.Spec

	maps := make(map[string]bool)

	for k, v := range e.Spec.ForProvider {
		// TODO use regex to find the pattern
		if strings.Contains(fmt.Sprintf("%v", k), pattern) {

			switch pattern {
			case "Selector":
				if v.(map[string]interface{})["matchLabels"].(map[string]interface{})["testing.upbound.io/example-name"] != nil {
					maps[v.(map[string]interface{})["matchLabels"].(map[string]interface{})["testing.upbound.io/example-name"].(string)] = true
				}
			case "Refs":
				for _, v2 := range v.([]interface{}) {
					m, ok := v2.(map[string]interface{})
					if !ok {
						continue
					}

					if name, ok := m["name"].(string); ok {
						maps[name] = true
					}
				}
			}

		}
	}
	return maps
}

// FindDependencies Find Selector and Ref in all files examples
func (e *ExamplesDetails) FindDependencies() {
	e.FindDependenciesSelector()
	e.FindDependenciesRefs()

	e.FindDependenciesSelectorFiles()
	e.FindDependenciesRefsFiles()
}

// FindDependenciesRefs Find Refs in all dependencies files examples
func (e *ExamplesDetails) FindDependenciesRefs() {
	for _, ex := range *e {
		for s := range ex.Refs {
			for _, ex2 := range *e {
				if s == ex2.Metadata.Name {
					maps.Copy(ex.Refs, ex2.Refs)
				}
			}
		}
	}
}

// FindDependenciesSelector Find Selectors in all dependencies files examples
func (e *ExamplesDetails) FindDependenciesSelector() {
	for _, ex := range *e {
		for s := range ex.Selectors {
			for _, ex2 := range *e {
				if s == ex2.Metadata.Labels.TestingUpboundIoExampleName {
					maps.Copy(ex.Selectors, ex2.Selectors)
				}
			}
		}
	}
}

// FindDependenciesSelectorFiles Find Selector in all files examples
func (e *ExamplesDetails) FindDependenciesSelectorFiles() {
	for _, ex := range *e {
		for s := range ex.Selectors {
			for _, ex2 := range *e {
				if ex2.Metadata.Labels.TestingUpboundIoExampleName == s {
					ex.DependenciesFiles[ex2.FileName] = true
				}
			}
		}
	}
}

// FindDependenciesRefsFiles Find Refs in all files examples
func (e *ExamplesDetails) FindDependenciesRefsFiles() {
	for _, ex := range *e {
		for s := range ex.Refs {
			for _, ex2 := range *e {
				if ex2.Metadata.Name == s {
					ex.DependenciesFiles[ex2.FileName] = true
				}
			}
		}
	}
}
