package commands

import (
	"github.com/spf13/cobra"

	workerorder "github.com/gohive/demo-worker-order/cmd"
)

var workerOrderCmd = &cobra.Command{
	Use:   "worker-order",
	Short: "Start the order processing worker",
	Long:  `Order worker runs background tasks for order processing (similar to Celery).`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		workerorder.Run(configPath)
	},
}

func init() {
	workerOrderCmd.Flags().StringP("config", "c", "./demo-worker-order/config.toml", "config file path")
	rootCmd.AddCommand(workerOrderCmd)
}
