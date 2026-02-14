package cmd

import (
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Manage API server configurations",
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
