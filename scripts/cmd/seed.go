package main

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/seed"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Long:  `Seed the database with initial test data for development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		force, _ := cmd.Flags().GetBool("force")
		return seed.Run(seed.Options{
			ConfigPath: configPath,
			Force:      force,
		})
	},
}

func init() {
	seedCmd.Flags().StringP("config", "c", defaultConfigPath, "config file path")
	seedCmd.Flags().BoolP("force", "f", false, "force re-seed (clear existing data)")
	rootCmd.AddCommand(seedCmd)
}
