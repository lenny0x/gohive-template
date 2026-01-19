package commands

import (
	"github.com/spf13/cobra"

	demoapi "github.com/gohive/demo-api/cmd"
)

var apiAdminCmd = &cobra.Command{
	Use:   "api-admin",
	Short: "Start the Admin API service",
	Long:  `Admin API provides backend management endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		demoapi.Run(configPath)
	},
}

func init() {
	apiAdminCmd.Flags().StringP("config", "c", "./demo-api/config.toml", "config file path")
	rootCmd.AddCommand(apiAdminCmd)
}
