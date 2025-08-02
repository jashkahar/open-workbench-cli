package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var templatesFS embed.FS

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(fs embed.FS) {
	templatesFS = fs
	rootCmd = &cobra.Command{
		Use:   "om",
		Short: "Open Workbench - A modern CLI for scaffolding web applications",
		Long: `Open Workbench (om) is a powerful command-line tool for scaffolding 
modern web applications with pre-configured templates and best practices.

The CLI supports multiple execution modes:
  - Interactive mode for guided project creation
  - Non-interactive CLI mode with command-line flags

Features:
  - Dynamic template system with conditional logic
  - Parameter validation and grouping
  - Post-scaffolding actions
  - Cross-platform support`,
	}

	// Add subcommands
	rootCmd.AddCommand(initCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.om.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
