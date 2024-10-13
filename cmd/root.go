package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "raftd",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(joinCmd)

	rootCmd.AddCommand(kvGetCmd)
	rootCmd.AddCommand(kvSetCmd)
	rootCmd.AddCommand(deleteCmd)
}
