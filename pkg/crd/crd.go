// Package crd provides a struct for parsing a CRD manifest.
package crd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CRD is a struct for parsing a CRD manifest.
type CRD struct {
	Spec struct {
		Group string `yaml:"group"`
		Names struct {
			Categories []string `yaml:"categories"`
			Kind       string   `yaml:"kind"`
			ListKind   string   `yaml:"listKind"`
			Plural     string   `yaml:"plural"`
			Singular   string   `yaml:"singular"`
		} `yaml:"names"`
		Scope string `yaml:"scope"`
	} `yaml:"spec"`
}

// GetCRDs returns a list of CRDs from the specified path.
func GetCRDs(path ...string) ([]CRD, error) {
	if len(path) == 0 {
		path = append(path, "package/crds/")
	}

	listFilesCrds := make([]string, 0)

	for _, p := range path {
		x, err := listFiles(p, ".yaml")
		if err != nil {
			return nil, err
		}

		listFilesCrds = append(listFilesCrds, x...)
	}

	crds := make([]CRD, 0)
	for _, file := range listFilesCrds {
		var w = &CRD{}

		y, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// Parse the file
		err = yaml.Unmarshal(y, w)
		if err != nil {
			return nil, err
		}

		crds = append(crds, *w)
	}

	return crds, nil
}

func listFiles(dir string, ext string) ([]string, error) {
	// Check if the directory exists
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("directory %s does not exist", dir)
	}

	// Get a list of files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Filter the list of files to only include files with the specified extension
	var filteredFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
			filteredFiles = append(filteredFiles, fmt.Sprintf("%s/%s", dir, file.Name()))
		}
	}

	return filteredFiles, nil
}
