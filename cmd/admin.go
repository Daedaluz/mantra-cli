package cmd

import (
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin api commands",
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.PersistentFlags().StringP("api-key", "t", envStr("API_KEY", "1234567890"), "api key for authentication [$API_KEY]")
}
