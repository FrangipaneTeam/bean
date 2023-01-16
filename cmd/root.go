// Package cmd contains the commands for the application.
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui/pages/loading"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	c = config.Provider{}

	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "A command-line tool to find username on websites",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().UTC().UnixNano())

			p := tea.NewProgram(loading.New(c))

			if _, err := p.Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		},
	}
)

// Execute executes the root command.
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&c.Path, "path", "p", ".", "your provider path")
	rootCmd.AddCommand(listTestedCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
