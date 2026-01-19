package commands

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/fixorder"
)

var fixOrderCmd = &cobra.Command{
	Use:   "fix-order",
	Short: "Fix expired pending orders",
	Long:  `Find and cancel orders that are pending but have expired.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		fixorder.Run(fixorder.Options{
			ConfigPath: configPath,
			DryRun:     dryRun,
		})
	},
}

func init() {
	fixOrderCmd.Flags().StringP("config", "c", "./config.yaml", "config file path")
	fixOrderCmd.Flags().Bool("dry-run", false, "show what would be changed without making changes")
	rootCmd.AddCommand(fixOrderCmd)
}
