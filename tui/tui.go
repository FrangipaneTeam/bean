package tui

import (
	"regexp"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/exp/maps"
)

var (
	reRefs, _     = regexp.Compile(`^\w+Refs$`)
	reSelector, _ = regexp.Compile(`^\w+Selector$`)
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
		ForProvider map[string]interface{} `yaml:"forProvider"`
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

// FindSelectorsAndRefs returns the selectors and refs of the example
func (e *Example) FindSelectorsAndRefs() (map[string]bool, map[string]bool) {
	mapsSelector := make(map[string]bool)
	mapsRefs := make(map[string]bool)

	for k, v := range e.Spec.ForProvider {
		if reSelector.MatchString(k) {
			if _, ok := v.(map[string]interface{}); ok {
				m := getSelector(v.(map[string]interface{}))
				maps.Copy(mapsSelector, m)
			}
		} else if reRefs.MatchString(k) {
			for _, v2 := range v.([]interface{}) {
				m, ok := v2.(map[string]interface{})
				if !ok {
					continue
				}

				if name, ok := m["name"].(string); ok {
					mapsRefs[name] = true
				}
			}
		} else {
			if _, ok := v.([]interface{}); ok {
				isArrayInterface(v.([]interface{}), mapsSelector, mapsRefs)
			}
		}
	}

	return mapsSelector, mapsRefs
}

func isArrayInterface(i []interface{}, mapsSelector map[string]bool, mapsRefs map[string]bool) {
	for _, a := range i {
		if _, ok := a.(map[string]interface{}); ok {
			for k, v := range a.(map[string]interface{}) {
				if reSelector.MatchString(k) {
					if _, ok := v.(map[string]interface{}); ok {
						m := getSelector(v.(map[string]interface{}))
						maps.Copy(mapsSelector, m)
					}
				} else if reRefs.MatchString(k) {
					for _, v2 := range v.([]interface{}) {
						m, ok := v2.(map[string]interface{})
						if !ok {
							continue
						}

						if name, ok := m["name"].(string); ok {
							mapsRefs[name] = true
						}
					}
				} else {
					if _, ok := v.([]interface{}); ok {
						isArrayInterface(v.([]interface{}), mapsSelector, mapsRefs)
					}
				}
			}
		}
	}
}

func getSelector(v map[string]interface{}) map[string]bool {
	m := make(map[string]bool)
	if _, ok := v["matchLabels"].(map[string]interface{}); !ok {
		return nil
	}
	matchLabels := v["matchLabels"].(map[string]interface{})
	exampleName := matchLabels["testing.upbound.io/example-name"]
	if exampleName != nil {
		if _, ok := exampleName.(string); ok {
			m[exampleName.(string)] = true
		} else {
			return nil
		}
	}
	return m
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
