package cmd

import (
	"fmt"

	"github.com/daedaluz/mantra-cli/internal/config"
	"github.com/spf13/cobra"
)

var apiAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add an API server configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		server, _ := cmd.Flags().GetString("server")
		apiKey, _ := cmd.Flags().GetString("api-key")
		plaintext, _ := cmd.Flags().GetBool("plaintext")
		skipVerify, _ := cmd.Flags().GetBool("skip-verify")

		cfg := config.Load()
		cfg.APIs[name] = &config.API{
			Server:     server,
			APIKey:     apiKey,
			Plaintext:  plaintext,
			SkipVerify: skipVerify,
		}
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		fmt.Printf("API %q added.\n", name)
	},
}

func init() {
	apiCmd.AddCommand(apiAddCmd)
	apiAddCmd.Flags().String("server", "mantra-api.inits.se:443", "Server address")
	apiAddCmd.Flags().String("api-key", "", "Optional API key for platform admin access")
	apiAddCmd.Flags().Bool("plaintext", false, "Use plaintext gRPC / h2c (no TLS)")
	apiAddCmd.Flags().Bool("skip-verify", false, "Skip TLS certificate verification")
}
