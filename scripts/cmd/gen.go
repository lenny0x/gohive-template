package main

import (
	"github.com/spf13/cobra"

	"github.com/gohive/scripts/gen"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate GORM query interfaces and model structs from the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		outPath, _ := cmd.Flags().GetString("out")
		modelPath, _ := cmd.Flags().GetString("model")
		return gen.Run(gen.Options{
			ConfigPath:   configPath,
			OutPath:      outPath,
			ModelPkgPath: modelPath,
		})
	},
}

func init() {
	genCmd.Flags().StringP("config", "c", defaultConfigPath, "config file path")
	genCmd.Flags().StringP("out", "o", "./models/query", "output directory for generated query code")
	genCmd.Flags().StringP("model", "m", "./models/query/model", "output directory for generated model structs")
	rootCmd.AddCommand(genCmd)
}
