package yaml

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// SplitYAML splits a YAML file into multiple YAML files.
func SplitYAML(resources []byte) ([][]byte, error) {
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

// IsYamlFile returns true if the file is a yaml file regardless of the extension.
func IsYamlFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".yaml" || ext == ".yml"
}
