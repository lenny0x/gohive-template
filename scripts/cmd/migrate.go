package main

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/migrate"
)

const defaultMigrationsDir = "./migrations"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		migrationsDir, _ := cmd.Flags().GetString("dir")
		return migrate.Up(configPath, migrationsDir)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		migrationsDir, _ := cmd.Flags().GetString("dir")
		return migrate.Down(configPath, migrationsDir)
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		migrationsDir, _ := cmd.Flags().GetString("dir")
		return migrate.Status(configPath, migrationsDir)
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		migrationsDir, _ := cmd.Flags().GetString("dir")
		migrationType, _ := cmd.Flags().GetString("type")
		return migrate.Create(migrationsDir, args[0], migrationType)
	},
}

func init() {
	// shared flags for up/down/status
	for _, cmd := range []*cobra.Command{migrateUpCmd, migrateDownCmd, migrateStatusCmd} {
		cmd.Flags().StringP("config", "c", defaultConfigPath, "config file path")
		cmd.Flags().StringP("dir", "d", defaultMigrationsDir, "migrations directory")
	}

	// create flags
	migrateCreateCmd.Flags().StringP("dir", "d", defaultMigrationsDir, "migrations directory")
	migrateCreateCmd.Flags().StringP("type", "t", "sql", "migration type: sql or go")

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateStatusCmd, migrateCreateCmd)
	rootCmd.AddCommand(migrateCmd)
}
