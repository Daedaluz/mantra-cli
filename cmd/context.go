package cmd

import (
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage client contexts",
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
