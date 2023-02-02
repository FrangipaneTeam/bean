//go:build ignore

// This program generates bean.yaml example config with default values. It can be invoked by running `go generate`.
package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/FrangipaneTeam/bean/internal/theme"
)

type config struct {
	Theme theme.Theme `yaml:"theme"`
}

func main() {
	theme := theme.Default()
	c := config{
		Theme: theme,
	}

	// yaml marshal the theme
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Printf("yaml marshal error : %v", err)
	}

	err = os.WriteFile("docs/bean.yaml", data, 0644)

	if err != nil {
		log.Printf("write file error : %v", err)
	}
}
