package commands

import (
	"github.com/spf13/cobra"

	wsserver "github.com/gohive/demo-ws/cmd"
)

var wsServerCmd = &cobra.Command{
	Use:   "ws-server",
	Short: "Start the WebSocket server",
	Long:  `WebSocket server handles real-time bidirectional communication.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		wsserver.Run(configPath)
	},
}

func init() {
	wsServerCmd.Flags().StringP("config", "c", "./demo-ws/config.toml", "config file path")
	rootCmd.AddCommand(wsServerCmd)
}
