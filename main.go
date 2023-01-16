// Package main provides the entry point for the application
package main

import (
	_ "embed"
	"strings"

	"github.com/FrangipaneTeam/bean/cmd"
)

var (
	//go:embed version.txt
	version string
)

func main() {
	version = strings.TrimSuffix(version, "\n")
	cmd.Execute(version)
}
