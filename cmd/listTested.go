// Package cmd contains the commands to build the list of tested resources.
package cmd

import (
	"github.com/FrangipaneTeam/bean/tools"
	"github.com/spf13/cobra"
)

var (
	listTestedCmd = &cobra.Command{
		Use:   "listTested",
		Short: "Generate a list of tested resources",
		Long:  `Generate a list of tested resources and write it to a file in the current directory named "listTested.md"`,
		Run: func(cmd *cobra.Command, args []string) {
			tools.GenerateListTested(c)
		},
	}
)
