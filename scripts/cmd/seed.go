package main

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/seed"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Long:  `Seed the database with initial test data for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		force, _ := cmd.Flags().GetBool("force")

		seed.Run(seed.Options{
			ConfigPath: configPath,
			Force:      force,
		})
	},
}

func init() {
	seedCmd.Flags().StringP("config", "c", "./config.yaml", "config file path")
	seedCmd.Flags().BoolP("force", "f", false, "force re-seed (clear existing data)")
	rootCmd.AddCommand(seedCmd)
}
