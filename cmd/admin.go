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

	apiKeyDefault := "1234567890"
	if activeAPI != nil {
		apiKeyDefault = ctxStr(activeAPI.APIKey, apiKeyDefault)
	}

	adminCmd.PersistentFlags().StringP("api-key", "t", envStr("API_KEY", apiKeyDefault), "api key for authentication [$API_KEY]")
}
