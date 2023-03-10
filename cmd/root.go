// Package cmd contains the commands for the application.
package cmd

import (
	"log"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"

	"github.com/FrangipaneTeam/bean/config"
	"github.com/FrangipaneTeam/bean/tui/pages/loading"
)

var (
	c       = config.Provider{}
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "bean",
		Short: "A command-line tool to find username on websites",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().UTC().UnixNano())
			c.Viper = viper.GetViper()
			p := tea.NewProgram(loading.New(c))

			if _, err := p.Run(); err != nil {
				log.Printf("error running program: %v", err)
				os.Exit(1)
			}
		},
	}
)

// Execute executes the root command.
func Execute(version string) {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&c.Path, "path", "p", ".", "your provider path")
	rootCmd.PersistentFlags().BoolVarP(&c.Debug, "debug", "d", false, "debug mode")
	rootCmd.AddCommand(listTestedCmd)
	c.Version = version

	githubTag := &latest.GithubTag{
		Owner:      "FrangipaneTeam",
		Repository: "bean",
	}

	res, err := latest.Check(githubTag, c.Version)
	if err == nil {
		if res.Outdated {
			c.NewVersion = res.Current
		}
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bean")
	}

	// viper.AutomaticEnv()

	viper.ReadInConfig()
}
