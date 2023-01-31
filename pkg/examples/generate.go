package examples

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/internal/exlist"
	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

// GenerateExamplesList generates the list of examples.
func GenerateExamplesList(c config.Provider) tea.Msg {
	examplesList, err := listDirExamples(c.Path)
	if err != nil {
		return errorpanel.ErrorMsg{
			Reason: "could not read examples directory",
			Cause:  err,
		}
	}
	var s exlist.LoadedExamples
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

		e := &exlist.Example{
			FileName: dirName,
			FullPath: dirName,
			Desc:     fmt.Sprintf("%d examples", len(s.Examples[dirName])),
		}
		rootExamples = append(rootExamples, e)
	}

	s.Examples["-"] = rootExamples

	examplesWithDependencies := exlist.ExamplesDetails{}
	for d, e := range s.Examples {
		if d == "-" {
			continue
		}
		for _, ex := range e {
			example, ok := ex.(*exlist.Example)
			if ok {
				examplesWithDependencies[example.ExampleID] = example
			}
		}
	}
	examplesWithDependencies.FindDependencies()
	return s
}

func listDirExamples(path string) ([]os.DirEntry, error) {
	examplesList, err := os.ReadDir(fmt.Sprintf("%s/examples", path))
	return examplesList, err
}

func createExampleList(dir string) ([]*exlist.Example, *errorpanel.ErrorMsg) {
	kindList, err := os.ReadDir(dir)
	if err != nil {
		return nil, &errorpanel.ErrorMsg{
			Reason: "could not read examples directory",
			Cause:  err,
		}
	}

	exampleList := make([]*exlist.Example, 0)

	// for the sub list of examples dir
	for _, sf := range kindList {
		if !sf.Type().IsRegular() {
			continue
		}

		// get only yaml files
		if isYaml := isYamlFile(sf.Name()); !isYaml {
			continue
		}

		// open and parse yaml file
		yfile, errReadFile := os.ReadFile(
			fmt.Sprintf("%s/%s", dir, sf.Name()),
		)
		if errReadFile != nil {
			return nil, &errorpanel.ErrorMsg{
				Reason: "could not read examples directory",
				Cause:  errReadFile,
			}
		}
		var k *exlist.Example
		err = yaml.Unmarshal(yfile, &k)
		if err != nil {
			return nil, &errorpanel.ErrorMsg{
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

// isYamlFile returns true if the file is a yaml file.
func isYamlFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".yaml" || ext == ".yml"
}

func checkForExtraFile(dir string, file string) (int, *errorpanel.ErrorMsg) {
	var (
		extraK8S  *exlist.Example
		extraKind int
	)

	extraYFile, err := os.ReadFile(
		fmt.Sprintf("%s/%s.extra", dir, file),
	)
	if err == nil {
		extraY, errSplitYaml := splitYAML(extraYFile)
		if errSplitYaml != nil {
			return extraKind, &errorpanel.ErrorMsg{
				Reason: "split yaml error",
				Cause:  errSplitYaml,
			}
		}
		for _, f := range extraY {
			err = yaml.Unmarshal(f, &extraK8S)
			if err != nil {
				return extraKind, &errorpanel.ErrorMsg{
					Reason: "unmarshal error",
					Cause:  err,
				}
			}
			extraKind++
		}
	}
	return extraKind, nil
}

func checkForSecretFile(dir string, file string) (int, *errorpanel.ErrorMsg) {
	var (
		extraK8S  *exlist.Example
		extraKind int
	)

	extraSFile, err := os.ReadFile(
		fmt.Sprintf("%s/%s.secret", dir, file),
	)
	if err == nil {
		extraY, errSplitYaml := splitYAML(extraSFile)
		if errSplitYaml != nil {
			return extraKind, &errorpanel.ErrorMsg{
				Reason: "split yaml error",
				Cause:  errSplitYaml,
			}
		}

		for _, f := range extraY {
			err = yaml.Unmarshal(f, &extraK8S)
			if err != nil {
				return extraKind, &errorpanel.ErrorMsg{
					Reason: "unmarshal error",
					Cause:  err,
				}
			}
			extraKind++
		}
	}
	return extraKind, nil
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
