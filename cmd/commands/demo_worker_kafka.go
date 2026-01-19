package commands

import (
	"github.com/spf13/cobra"

	workerkafka "github.com/gohive/demo-worker-kafka/cmd"
)

var workerKafkaCmd = &cobra.Command{
	Use:   "worker-kafka",
	Short: "Start the Kafka consumer worker",
	Long:  `Kafka worker consumes messages from Kafka topics and processes them.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		workerkafka.Run(configPath)
	},
}

func init() {
	workerKafkaCmd.Flags().StringP("config", "c", "./demo-worker-kafka/config.toml", "config file path")
	rootCmd.AddCommand(workerKafkaCmd)
}
