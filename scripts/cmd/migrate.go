package main

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/migrate"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to create or update database tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		migrate.Run(configPath)
	},
}

func init() {
	migrateCmd.Flags().StringP("config", "c", "./config.yaml", "config file path")
	rootCmd.AddCommand(migrateCmd)
}
