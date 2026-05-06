package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// defaultConfigPath is the fallback config used by every subcommand when
// the user does not pass -c. Scripts run from the repo root and need a
// real, present config file with [database]/[log] sections.
const defaultConfigPath = "./demo-api/config.toml"

var rootCmd = &cobra.Command{
	Use:   "scripts",
	Short: "Gohive maintenance scripts",
	Long:  `A unified CLI to run database migrations, seeders, and one-off data fixes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
